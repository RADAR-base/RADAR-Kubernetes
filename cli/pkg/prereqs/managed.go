package prereqs

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const pinnedHelmfileVersion = "0.169.1"
const pinnedHelmVersion = "3.16.3"

// ManagedBinDir returns the directory where radarctl stores its managed binaries.
func ManagedBinDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".radarctl", "bin")
}

// ManagedHelmfilePath returns the path to the radarctl-managed helmfile binary if it
// exists, or "" if it has not been downloaded yet.
func ManagedHelmfilePath() string {
	return managedToolPath("helmfile")
}

// ManagedHelmPath returns the path to the radarctl-managed helm binary if it
// exists, or "" if it has not been downloaded yet.
func ManagedHelmPath() string {
	return managedToolPath("helm")
}

func managedToolPath(name string) string {
	p := filepath.Join(ManagedBinDir(), name)
	if _, err := os.Stat(p); err != nil {
		return ""
	}
	return p
}

// ManagedBinEnv returns an ExtraEnv slice that prepends the managed bin dir to PATH
// if any managed binaries are present. Pass this to executor.ShellExecutor.ExtraEnv
// so that helmfile picks up the managed helm binary.
func ManagedBinEnv() []string {
	dir := ManagedBinDir()
	if _, err := os.Stat(dir); err != nil {
		return nil
	}
	return []string{"PATH=" + dir + ":" + os.Getenv("PATH")}
}

// DownloadHelmfile downloads the pinned helmfile v0.x binary from GitHub Releases into
// ~/.radarctl/bin/helmfile. It is called by AutoFix when the system helmfile is missing
// or at an incompatible version (>= v1).
func DownloadHelmfile() error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Map Go arch names to the names used in the release asset filenames.
	archMap := map[string]string{
		"amd64": "amd64",
		"arm64": "arm64",
		"386":   "386",
	}
	arch, ok := archMap[goarch]
	if !ok {
		return fmt.Errorf("unsupported architecture %s — download helmfile v%s manually from https://github.com/helmfile/helmfile/releases/tag/v%s", goarch, pinnedHelmfileVersion, pinnedHelmfileVersion)
	}

	osName := goos // darwin / linux match GitHub asset names directly
	if osName != "darwin" && osName != "linux" {
		return fmt.Errorf("unsupported OS %s — download helmfile v%s manually from https://github.com/helmfile/helmfile/releases/tag/v%s", osName, pinnedHelmfileVersion, pinnedHelmfileVersion)
	}

	assetName := fmt.Sprintf("helmfile_%s_%s_%s.tar.gz", pinnedHelmfileVersion, osName, arch)
	url := fmt.Sprintf("https://github.com/helmfile/helmfile/releases/download/v%s/%s", pinnedHelmfileVersion, assetName)

	fmt.Printf("Downloading helmfile v%s from GitHub...\n", pinnedHelmfileVersion)

	resp, err := http.Get(url) //nolint:gosec // URL is constructed from a hardcoded version string
	if err != nil {
		return fmt.Errorf("downloading helmfile: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading helmfile: HTTP %d from %s", resp.StatusCode, url)
	}

	binDir := ManagedBinDir()
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", binDir, err)
	}

	destPath := filepath.Join(binDir, "helmfile")
	if err := extractHelmfileBinary(resp.Body, destPath); err != nil {
		return fmt.Errorf("extracting helmfile binary: %w", err)
	}

	fmt.Printf("helmfile v%s installed to %s\n", pinnedHelmfileVersion, destPath)
	return nil
}

// DownloadHelm downloads the pinned helm v3.x binary from the official release URL
// into ~/.radarctl/bin/helm so it shadows a system-installed helm v4+.
func DownloadHelm() error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	archMap := map[string]string{"amd64": "amd64", "arm64": "arm64", "386": "386"}
	arch, ok := archMap[goarch]
	if !ok {
		return fmt.Errorf("unsupported architecture %s — download helm v%s manually", goarch, pinnedHelmVersion)
	}
	osName := goos
	if osName != "darwin" && osName != "linux" {
		return fmt.Errorf("unsupported OS %s — download helm v%s manually", osName, pinnedHelmVersion)
	}

	// Helm release asset: helm-v3.16.3-darwin-arm64.tar.gz
	assetName := fmt.Sprintf("helm-v%s-%s-%s.tar.gz", pinnedHelmVersion, osName, arch)
	url := fmt.Sprintf("https://get.helm.sh/%s", assetName)

	fmt.Printf("Downloading helm v%s from get.helm.sh...\n", pinnedHelmVersion)
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("downloading helm: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading helm: HTTP %d from %s", resp.StatusCode, url)
	}

	binDir := ManagedBinDir()
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", binDir, err)
	}
	destPath := filepath.Join(binDir, "helm")
	if err := extractNamedBinary(resp.Body, "helm", destPath); err != nil {
		return fmt.Errorf("extracting helm binary: %w", err)
	}
	fmt.Printf("helm v%s installed to %s\n", pinnedHelmVersion, destPath)
	return nil
}

// extractHelmfileBinary reads a .tar.gz stream, finds the "helmfile" entry, and writes
// it as an executable file to destPath.
func extractHelmfileBinary(r io.Reader, destPath string) error {
	return extractNamedBinary(r, "helmfile", destPath)
}

// extractNamedBinary reads a .tar.gz stream, finds an entry whose base name matches
// name, and writes it as an executable file to destPath.
func extractNamedBinary(r io.Reader, name, destPath string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) != name {
			continue
		}
		tmp, err := os.CreateTemp(filepath.Dir(destPath), "."+name+"-download-*")
		if err != nil {
			return err
		}
		if _, err := io.Copy(tmp, tr); err != nil { //nolint:gosec
			tmp.Close()
			os.Remove(tmp.Name())
			return err
		}
		tmp.Close()
		if err := os.Chmod(tmp.Name(), 0o755); err != nil {
			os.Remove(tmp.Name())
			return err
		}
		return os.Rename(tmp.Name(), destPath)
	}
	return fmt.Errorf("%s binary not found in archive", name)
}

package wizard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
)

// ensureLocalCluster makes sure a local single-node Kubernetes cluster is reachable
// for dev/demo profiles. It checks Docker first (required by k3d on macOS), then lists
// configured kubectl contexts filtered to local ones, and lets the user pick one. If
// none exist it offers to install k3d/k3s.
func ensureLocalCluster(a *Answers, repoRoot string) error {
	// Docker is required on macOS (k3d runs nodes as containers). On Linux, k3s is
	// self-contained, so we only require Docker there if the user already has k3d contexts.
	if runtime.GOOS == "darwin" {
		if err := ensureDocker(); err != nil {
			return err
		}
	}

	localCtxs := localKubeContexts()

	// On Linux, if no local contexts were found we'll install k3s (no Docker needed).
	// But if existing contexts imply k3d, Docker must be running there too.
	if runtime.GOOS == "linux" && hasK3dContext(localCtxs) {
		if err := ensureDocker(); err != nil {
			return err
		}
	}

	if len(localCtxs) > 0 {
		return selectLocalContext(a, localCtxs)
	}

	output.Warning("No local Kubernetes contexts found.")

	var confirm bool
	prompt := "Install a local single-node cluster now?"
	switch runtime.GOOS {
	case "darwin":
		prompt = "Install k3d (Docker-based K3s) locally and create a 'radar' cluster now?"
	case "linux":
		prompt = "Install k3s locally now? (requires sudo)"
	}
	if err := runForm(huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title(prompt).Value(&confirm),
	))); err != nil {
		return err
	}
	if !confirm {
		return fmt.Errorf("a local cluster is required for the %q profile", a.Profile)
	}

	switch runtime.GOOS {
	case "darwin":
		if err := installK3dMac(); err != nil {
			return err
		}
		if err := createK3dCluster("radar"); err != nil {
			return err
		}
		a.KubeContext = "k3d-radar"
	case "linux":
		if err := installK3sLinux(); err != nil {
			return err
		}
		a.KubeContext = "default"
	default:
		return fmt.Errorf("automatic local-cluster install is not supported on %s — install K3s/k3d manually and re-run", runtime.GOOS)
	}

	output.Success(fmt.Sprintf("Local cluster ready — kubectl context: %s", a.KubeContext))
	return nil
}

// selectLocalContext shows the list of local contexts and asks the user to confirm one.
func selectLocalContext(a *Answers, ctxs []string) error {
	current := currentKubeContext()

	// Pre-select: if a previously stored context is in the list keep it,
	// otherwise default to the active kubectl context.
	selected := a.KubeContext
	if selected == "" || !containsString(ctxs, selected) {
		if containsString(ctxs, current) {
			selected = current
		} else {
			selected = ctxs[0]
		}
	}

	opts := make([]huh.Option[string], len(ctxs))
	for i, c := range ctxs {
		label := c
		if c == current {
			label = c + "  (active)"
		}
		opts[i] = huh.NewOption(label, c)
	}

	if err := runForm(huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Select a local Kubernetes context").
			Description("Only contexts whose cluster server points to localhost or uses a known local\ndistribution (k3d, kind, minikube, docker-desktop, rancher-desktop) are shown.").
			Options(opts...).
			Value(&selected),
	))); err != nil {
		return err
	}

	a.KubeContext = selected
	output.Info(fmt.Sprintf("Using kubectl context: %s", selected))
	return ensureContextRunning(selected)
}

// ensureContextRunning checks that the cluster for the given kubectl context is
// reachable. If not, it tries to start it based on the distribution inferred from
// the context name. Returns an error only if the cluster can't be made ready.
func ensureContextRunning(ctx string) error {
	if isContextReachable(ctx) {
		return nil
	}

	output.Warning(fmt.Sprintf("Cluster for context %q is not responding.", ctx))

	// Infer distribution from context name to choose the right start command.
	ctxLower := strings.ToLower(ctx)
	switch {
	case strings.HasPrefix(ctxLower, "k3d-"):
		clusterName := ctx[4:] // strip "k3d-" prefix
		return startK3dCluster(clusterName)

	case strings.HasPrefix(ctxLower, "minikube") || ctxLower == "minikube":
		return startMinikube(ctx)

	case strings.HasPrefix(ctxLower, "kind-"):
		return fmt.Errorf(
			"kind cluster %q is not running — start it with: kind create cluster --name %s",
			ctx, ctx[5:])

	case ctxLower == "docker-desktop":
		return fmt.Errorf(
			"Docker Desktop Kubernetes is not running — enable it in Docker Desktop → Settings → Kubernetes → Enable Kubernetes")

	case ctxLower == "rancher-desktop":
		return fmt.Errorf(
			"Rancher Desktop Kubernetes is not running — open Rancher Desktop and wait for the cluster to start")

	default:
		return fmt.Errorf(
			"cluster for context %q is not reachable — start it manually and re-run", ctx)
	}
}

// isContextReachable returns true when `kubectl cluster-info` succeeds for the
// given context within a short timeout.
func isContextReachable(ctx string) bool {
	out, err := exec.Command(
		"kubectl", "cluster-info",
		"--context="+ctx,
		"--request-timeout=5s",
	).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), "running")
}

// startK3dCluster starts a stopped k3d cluster by name.
func startK3dCluster(name string) error {
	var confirm bool
	if err := runForm(huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("Start k3d cluster %q?", name)).
			Value(&confirm),
	))); err != nil {
		return err
	}
	if !confirm {
		return fmt.Errorf("cluster %q must be running to continue", name)
	}

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Starting k3d cluster %q...", name))
	if out, err := exec.Command("k3d", "cluster", "start", name).CombinedOutput(); err != nil {
		spinner.Fail("Failed to start cluster")
		return fmt.Errorf("k3d cluster start %s: %s", name, strings.TrimSpace(string(out)))
	}
	spinner.Success(fmt.Sprintf("k3d cluster %q started", name))

	// Wait up to 30 s for the API server to become ready.
	for range 15 {
		if isContextReachable("k3d-" + name) {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("cluster %q started but API server did not become ready in time", name)
}

// startMinikube starts or resumes a minikube cluster.
func startMinikube(ctx string) error {
	var confirm bool
	if err := runForm(huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Start minikube cluster?").
			Value(&confirm),
	))); err != nil {
		return err
	}
	if !confirm {
		return fmt.Errorf("minikube must be running to continue")
	}

	args := []string{"start"}
	// If the context name encodes a profile (minikube uses profile as context name)
	// and it's not the default "minikube", pass it explicitly.
	if ctx != "minikube" {
		args = append(args, "--profile", ctx)
	}

	spinner, _ := pterm.DefaultSpinner.Start("Starting minikube...")
	if out, err := exec.Command("minikube", args...).CombinedOutput(); err != nil {
		spinner.Fail("Failed to start minikube")
		return fmt.Errorf("minikube start: %s", strings.TrimSpace(string(out)))
	}
	spinner.Success("minikube started")
	return nil
}

// ensureDocker checks that Docker is installed and the daemon is reachable.
// If Docker is not installed it offers to install it (Homebrew cask on macOS,
// convenience script on Linux). If installed but not running it offers to start it.
func ensureDocker() error {
	installed := isDockerInstalled()
	if !installed {
		output.Warning("Docker is not installed. k3d requires Docker to run cluster nodes as containers.")
		var confirm bool
		switch runtime.GOOS {
		case "darwin":
			if err := runForm(huh.NewForm(huh.NewGroup(
				huh.NewConfirm().
					Title("Install Docker Desktop via Homebrew?").
					Value(&confirm),
			))); err != nil {
				return err
			}
			if !confirm {
				return fmt.Errorf("Docker is required for local clusters — install Docker Desktop and re-run")
			}
			if err := installDockerMac(); err != nil {
				return err
			}
		case "linux":
			if err := runForm(huh.NewForm(huh.NewGroup(
				huh.NewConfirm().
					Title("Install Docker Engine via the official convenience script? (requires sudo)").
					Value(&confirm),
			))); err != nil {
				return err
			}
			if !confirm {
				return fmt.Errorf("Docker is required for local clusters — install Docker Engine and re-run")
			}
			if err := installDockerLinux(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Docker is not installed — install Docker Desktop from https://www.docker.com and re-run")
		}
	}

	if isDockerRunning() {
		return nil
	}

	output.Warning("Docker is installed but the daemon is not running.")
	var startIt bool
	if err := runForm(huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Start Docker now?").
			Value(&startIt),
	))); err != nil {
		return err
	}
	if !startIt {
		return fmt.Errorf("Docker must be running for local clusters — start Docker and re-run")
	}
	return startDocker()
}

func isDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func isDockerRunning() bool {
	err := exec.Command("docker", "info").Run()
	return err == nil
}

func installDockerMac() error {
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew is required to install Docker Desktop — install it from https://brew.sh and re-run")
	}
	spinner, _ := pterm.DefaultSpinner.Start("Installing Docker Desktop via Homebrew (this may take a few minutes)...")
	if out, err := exec.Command("brew", "install", "--cask", "docker").CombinedOutput(); err != nil {
		spinner.Fail("Docker Desktop installation failed")
		return fmt.Errorf("brew install --cask docker: %s", strings.TrimSpace(string(out)))
	}
	spinner.Success("Docker Desktop installed")
	return startDocker()
}

func installDockerLinux() error {
	spinner, _ := pterm.DefaultSpinner.Start("Installing Docker Engine (requires sudo)...")
	cmd := exec.Command("sh", "-c", "curl -fsSL https://get.docker.com | sh")
	if out, err := cmd.CombinedOutput(); err != nil {
		spinner.Fail("Docker Engine installation failed")
		return fmt.Errorf("docker install script: %s", strings.TrimSpace(string(out)))
	}
	spinner.Success("Docker Engine installed")
	// Enable and start the service.
	_ = exec.Command("sudo", "systemctl", "enable", "--now", "docker").Run()
	return waitForDocker()
}

func startDocker() error {
	switch runtime.GOOS {
	case "darwin":
		spinner, _ := pterm.DefaultSpinner.Start("Starting Docker Desktop...")
		if err := exec.Command("open", "-a", "Docker").Run(); err != nil {
			spinner.Fail("Could not launch Docker Desktop")
			return fmt.Errorf("open Docker: %w", err)
		}
		if err := waitForDocker(); err != nil {
			spinner.Fail("Docker Desktop did not start in time")
			return err
		}
		spinner.Success("Docker Desktop is running")
	case "linux":
		spinner, _ := pterm.DefaultSpinner.Start("Starting Docker daemon (requires sudo)...")
		if err := exec.Command("sudo", "systemctl", "start", "docker").Run(); err != nil {
			spinner.Fail("Could not start Docker daemon")
			return fmt.Errorf("systemctl start docker: %w", err)
		}
		if err := waitForDocker(); err != nil {
			spinner.Fail("Docker daemon did not start in time")
			return err
		}
		spinner.Success("Docker daemon is running")
	default:
		return fmt.Errorf("please start Docker manually and re-run")
	}
	return nil
}

// waitForDocker polls `docker info` for up to 60 seconds.
func waitForDocker() error {
	for range 30 {
		if isDockerRunning() {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("Docker did not become ready within 60 seconds — start it manually and re-run")
}

func hasK3dContext(ctxs []string) bool {
	for _, c := range ctxs {
		if strings.HasPrefix(strings.ToLower(c), "k3d-") {
			return true
		}
	}
	return false
}

// localKubeContexts returns the names of kubectl contexts whose server URL resolves
// to a local address or whose name/type matches a well-known local distribution.
func localKubeContexts() []string {
	out, err := exec.Command("kubectl", "config", "get-contexts", "-o", "name").Output()
	if err != nil {
		return nil
	}
	all := strings.Fields(strings.TrimSpace(string(out)))

	var local []string
	for _, ctx := range all {
		if isLocalContext(ctx) {
			local = append(local, ctx)
		}
	}
	return local
}

// isLocalContext returns true when the context's cluster server points to localhost
// or the context name matches a known local distribution pattern.
func isLocalContext(ctx string) bool {
	// Well-known local distribution name prefixes / exact names.
	localPrefixes := []string{
		"k3d-", "kind-", "minikube", "docker-desktop", "rancher-desktop", "k3s",
	}
	ctxLower := strings.ToLower(ctx)
	for _, p := range localPrefixes {
		if ctxLower == p || strings.HasPrefix(ctxLower, p) {
			return true
		}
	}

	// Fall back to inspecting the server URL for this context.
	server := contextServerURL(ctx)
	return isLocalAddress(server)
}

// contextServerURL fetches the cluster.server field for the named context.
func contextServerURL(ctx string) string {
	out, err := exec.Command(
		"kubectl", "config", "view",
		"--context="+ctx,
		"--minify",
		"-o", "jsonpath={.clusters[0].cluster.server}",
	).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// isLocalAddress returns true for URLs whose host is a loopback address or hostname.
func isLocalAddress(server string) bool {
	// Strip scheme.
	host := server
	if i := strings.Index(host, "://"); i >= 0 {
		host = host[i+3:]
	}
	// Strip port.
	if i := strings.LastIndex(host, ":"); i >= 0 {
		host = host[:i]
	}
	// Strip brackets from IPv6.
	host = strings.Trim(host, "[]")

	switch {
	case host == "localhost",
		host == "::1",
		strings.HasPrefix(host, "127."),
		strings.HasSuffix(host, ".local"):
		return true
	}
	return false
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func currentKubeContext() string {
	out, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func installK3dMac() error {
	if _, err := exec.LookPath("k3d"); err == nil {
		return nil
	}
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew is required to install k3d on macOS — install it from https://brew.sh and re-run")
	}
	spinner, _ := pterm.DefaultSpinner.Start("Installing k3d via Homebrew...")
	cmd := exec.Command("brew", "install", "k3d")
	if out, err := cmd.CombinedOutput(); err != nil {
		spinner.Fail("Failed to install k3d")
		return fmt.Errorf("brew install k3d failed: %s", strings.TrimSpace(string(out)))
	}
	spinner.Success("k3d installed")
	return nil
}

func createK3dCluster(name string) error {
	// Idempotent: skip if a cluster with this name already exists.
	if out, _ := exec.Command("k3d", "cluster", "list", "-o", "json").Output(); strings.Contains(string(out), `"name":"`+name+`"`) {
		output.Info(fmt.Sprintf("k3d cluster %q already exists — reusing it", name))
		return nil
	}
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Creating k3d cluster %q (port 80/443 mapped to host)...", name))
	cmd := exec.Command("k3d", "cluster", "create", name,
		"--port", "80:80@loadbalancer",
		"--port", "443:443@loadbalancer",
		// Disable k3s's built-in traefik ingress controller so that
		// ingress-nginx can own ports 80/443 via the svclb LoadBalancer.
		"--k3s-arg", "--disable=traefik@server:0",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		spinner.Fail("Cluster creation failed")
		return fmt.Errorf("k3d cluster create: %s", strings.TrimSpace(string(out)))
	}
	spinner.Success(fmt.Sprintf("k3d cluster %q created", name))
	return nil
}

func installK3sLinux() error {
	if _, err := exec.LookPath("k3s"); err == nil {
		return nil
	}
	spinner, _ := pterm.DefaultSpinner.Start("Installing k3s (requires sudo)...")
	cmd := exec.Command("sh", "-c", "curl -sfL https://get.k3s.io | sh -")
	if out, err := cmd.CombinedOutput(); err != nil {
		spinner.Fail("Failed to install k3s")
		return fmt.Errorf("k3s install failed: %s", strings.TrimSpace(string(out)))
	}
	spinner.Success("k3s installed")
	return nil
}

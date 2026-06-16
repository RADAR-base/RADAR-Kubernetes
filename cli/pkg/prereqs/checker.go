package prereqs

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type CheckResult struct {
	Tool           string
	OK             bool
	Version        string
	Error          string
	InstallHint    string
	NeedsUpgrade   bool       // installed but below minVersion
	NeedsDowngrade bool       // installed but >= maxVersion
	FixCmds        [][]string // commands to run to install / upgrade
	FixFunc        func() error // called instead of FixCmds when non-nil
}

type toolDef struct {
	name         string
	args         []string
	parseVersion func(string) string
	installHint  string
	minVersion   string
	maxVersion   string     // exclusive upper bound; version must be < this; "" = no limit
	installCmds  [][]string
	upgradeCmds  [][]string
	fixFunc      func() error // overrides FixCmds for complex fixes (e.g. binary download)
}

var tools []toolDef

func init() {
	tools = commonTools
	if runtime.GOOS == "darwin" {
		tools = append(tools, macTools...)
	}
}

var commonTools = []toolDef{
	{
		name: "kubectl",
		args: []string{"version", "--client", "-o", "json"},
		parseVersion: func(out string) string {
			for _, line := range strings.Split(out, "\n") {
				if strings.Contains(line, "gitVersion") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						return strings.Trim(strings.TrimSpace(parts[1]), `",`)
					}
				}
			}
			return out
		},
		installHint: "https://kubernetes.io/docs/tasks/tools/",
		installCmds: [][]string{{"brew", "install", "kubectl"}},
		upgradeCmds: [][]string{{"brew", "upgrade", "kubectl"}},
	},
	{
		// Helmfile v0.169.1 uses `helm version --client --short` which was removed in
		// Helm v4. Pin to the last compatible v3 release.
		name: "helm",
		args: []string{"version", "--short"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "radarctl will download helm v3.16.3 automatically",
		minVersion:   "v3",
		maxVersion:   "v4.0.0",
		fixFunc:      DownloadHelm,
	},
	{
		// RADAR-Kubernetes environments.yaml uses Go template syntax that Helmfile v1
		// requires the .gotmpl extension for. Pin to the last compatible v0 release.
		name: "helmfile",
		args: []string{"--version"},
		parseVersion: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return out
		},
		installHint: "radarctl will download helmfile v0.169.1 automatically",
		minVersion:  "v0.169.1",
		maxVersion:  "v1.0.0",
		fixFunc:     DownloadHelmfile,
	},
	{
		name:         "helm diff",
		args:         []string{"diff", "version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "helm plugin install https://github.com/databus23/helm-diff",
		minVersion:   "v3.9.12",
		installCmds: [][]string{
			{"helm", "plugin", "remove", "diff"},
			{"helm", "plugin", "install", "--verify=false", "https://github.com/databus23/helm-diff"},
		},
		upgradeCmds: [][]string{{"helm", "plugin", "update", "diff"}},
	},
	{
		name: "yq",
		args: []string{"--version"},
		parseVersion: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return out
		},
		installHint: "https://github.com/mikefarah/yq#install",
		minVersion:  "v4.44.3",
		installCmds: [][]string{{"brew", "install", "yq"}},
		upgradeCmds: [][]string{{"brew", "upgrade", "yq"}},
	},
	{
		name: "java",
		args: []string{"-version"},
		parseVersion: func(out string) string {
			lines := strings.Split(out, "\n")
			if len(lines) > 0 {
				return strings.TrimSpace(lines[0])
			}
			return out
		},
		installHint: "https://adoptium.net/ or: brew install openjdk",
		installCmds: [][]string{{"brew", "install", "openjdk"}},
		upgradeCmds: [][]string{{"brew", "upgrade", "openjdk"}},
	},
	{
		name:         "openssl",
		args:         []string{"version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "brew install openssl  (macOS) or: apt install openssl",
		installCmds:  [][]string{{"brew", "install", "openssl"}},
		upgradeCmds:  [][]string{{"brew", "upgrade", "openssl"}},
	},
	{
		name:         "git",
		args:         []string{"--version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "https://git-scm.com/downloads",
		installCmds:  [][]string{{"brew", "install", "git"}},
		upgradeCmds:  [][]string{{"brew", "upgrade", "git"}},
	},
}

var macTools = []toolDef{
	{
		// bin/util.sh uses GNU sed syntax; macOS ships BSD sed which is incompatible.
		name:         "gsed",
		args:         []string{"--version"},
		parseVersion: func(out string) string { return strings.TrimSpace(strings.Split(out, "\n")[0]) },
		installHint:  "brew install gnu-sed",
		installCmds:  [][]string{{"brew", "install", "gnu-sed"}},
		upgradeCmds:  [][]string{{"brew", "upgrade", "gnu-sed"}},
	},
}

// semverCmp compares two version strings, returning -1, 0, or +1.
func semverCmp(a, b string) int {
	clean := func(s string) string {
		s = strings.TrimPrefix(strings.TrimSpace(s), "v")
		if i := strings.IndexAny(s, "+-"); i >= 0 {
			s = s[:i]
		}
		return s
	}
	aParts := strings.Split(clean(a), ".")
	bParts := strings.Split(clean(b), ".")
	for len(aParts) < len(bParts) {
		aParts = append(aParts, "0")
	}
	for len(bParts) < len(aParts) {
		bParts = append(bParts, "0")
	}
	for i := range aParts {
		av, _ := strconv.Atoi(aParts[i])
		bv, _ := strconv.Atoi(bParts[i])
		if av > bv {
			return 1
		}
		if av < bv {
			return -1
		}
	}
	return 0
}

func meetsMinVersion(version, min string) bool  { return semverCmp(version, min) >= 0 }
func belowMaxVersion(version, max string) bool   { return semverCmp(version, max) < 0 }

// Check runs all prerequisite checks. For helmfile it first checks the radarctl-managed
// binary (~/.radarctl/bin/helmfile) so a previously downloaded v0.x binary takes priority
// over any system-installed v1+ binary.
func Check(ex executor.Executor) []CheckResult {
	results := make([]CheckResult, 0, len(tools))
	for _, t := range tools {
		toolName := strings.Fields(t.name)[0]

		// Prefer managed binaries when present (helmfile and helm are pinned to v0.x/v3.x).
		bin := toolName
		switch toolName {
		case "helmfile":
			if p := ManagedHelmfilePath(); p != "" {
				bin = p
			}
		case "helm":
			if p := ManagedHelmPath(); p != "" {
				bin = p
			}
		}

		out, err := ex.Run(bin, t.args...)
		if err != nil {
			errMsg := "not found"
			for _, line := range strings.Split(err.Error(), "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "exit status") {
					errMsg = line
					break
				}
			}
			results = append(results, CheckResult{
				Tool:        t.name,
				OK:          false,
				Error:       errMsg,
				InstallHint: t.installHint,
				FixCmds:     t.installCmds,
				FixFunc:     t.fixFunc,
			})
			continue
		}

		version := t.parseVersion(out)
		r := CheckResult{Tool: t.name, OK: true, Version: version}

		switch {
		case t.maxVersion != "" && !belowMaxVersion(version, t.maxVersion):
			r.OK = false
			r.NeedsDowngrade = true
			r.Error = fmt.Sprintf("version %s is incompatible — need < %s (radarctl will install a compatible version to ~/.radarctl/bin/)", version, t.maxVersion)
			r.InstallHint = t.installHint
			r.FixFunc = t.fixFunc
		case t.minVersion != "" && !meetsMinVersion(version, t.minVersion):
			r.OK = false
			r.NeedsUpgrade = true
			r.Error = fmt.Sprintf("version %s found, need >= %s", version, t.minVersion)
			r.InstallHint = t.installHint
			r.FixCmds = t.upgradeCmds
			r.FixFunc = t.fixFunc
		}

		results = append(results, r)
	}
	return results
}

// AutoFix runs the fix for a failed check result. FixFunc takes priority over FixCmds.
func AutoFix(ex executor.Executor, result CheckResult) error {
	if result.FixFunc != nil {
		return result.FixFunc()
	}
	for i, cmd := range result.FixCmds {
		_, err := ex.Run(cmd[0], cmd[1:]...)
		if err != nil && i == len(result.FixCmds)-1 {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

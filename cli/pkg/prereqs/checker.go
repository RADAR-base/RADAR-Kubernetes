package prereqs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type CheckResult struct {
	Tool         string
	OK           bool
	Version      string
	Error        string
	InstallHint  string
	NeedsUpgrade bool       // true if installed but below minimum version
	FixCmds      [][]string // commands to run to install or upgrade
}

type toolDef struct {
	name         string
	args         []string
	parseVersion func(string) string
	installHint  string
	minVersion   string
	installCmds  [][]string // commands to install from scratch
	upgradeCmds  [][]string // commands to upgrade an existing install
}

var tools = []toolDef{
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
		name:         "helm",
		args:         []string{"version", "--short"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "https://helm.sh/docs/intro/install/",
		minVersion:   "v3",
		installCmds:  [][]string{{"brew", "install", "helm"}},
		upgradeCmds:  [][]string{{"brew", "upgrade", "helm"}},
	},
	{
		name: "helmfile",
		args: []string{"--version"},
		parseVersion: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return out
		},
		installHint: "https://github.com/helmfile/helmfile/releases",
		minVersion:  "v0.169.1",
		installCmds: [][]string{{"brew", "install", "helmfile"}},
		upgradeCmds: [][]string{{"brew", "upgrade", "helmfile"}},
	},
	{
		name:         "helm diff",
		args:         []string{"diff", "version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "helm plugin install https://github.com/databus23/helm-diff",
		minVersion:   "v3.9.12",
		// Remove first (best-effort) in case of broken install, then reinstall.
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

// meetsMinVersion returns true if version >= minVersion using simple semver comparison.
func meetsMinVersion(version, minVersion string) bool {
	stripV := func(s string) string {
		return strings.TrimPrefix(strings.TrimSpace(s), "v")
	}
	ver := strings.FieldsFunc(stripV(version), func(r rune) bool { return r == '+' || r == '-' })
	min := strings.FieldsFunc(stripV(minVersion), func(r rune) bool { return r == '+' || r == '-' })

	var verStr, minStr string
	if len(ver) > 0 {
		verStr = ver[0]
	}
	if len(min) > 0 {
		minStr = min[0]
	}

	vParts := strings.Split(verStr, ".")
	mParts := strings.Split(minStr, ".")
	for len(vParts) < len(mParts) {
		vParts = append(vParts, "0")
	}
	for len(mParts) < len(vParts) {
		mParts = append(mParts, "0")
	}

	for i := range vParts {
		v, _ := strconv.Atoi(vParts[i])
		m, _ := strconv.Atoi(mParts[i])
		if v > m {
			return true
		}
		if v < m {
			return false
		}
	}
	return true
}

// Check runs all prerequisite checks and returns one result per tool.
func Check(exec executor.Executor) []CheckResult {
	results := make([]CheckResult, 0, len(tools))
	for _, t := range tools {
		toolName := strings.Fields(t.name)[0]
		out, err := exec.Run(toolName, t.args...)
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
			})
			continue
		}
		version := t.parseVersion(out)
		r := CheckResult{
			Tool:    t.name,
			OK:      true,
			Version: version,
		}
		if t.minVersion != "" && !meetsMinVersion(version, t.minVersion) {
			r.OK = false
			r.NeedsUpgrade = true
			r.Error = fmt.Sprintf("version %s found, need >= %s", version, t.minVersion)
			r.InstallHint = t.installHint
			r.FixCmds = t.upgradeCmds
		}
		results = append(results, r)
	}
	return results
}

// AutoFix runs the fix commands for a failed check result.
// Intermediate commands are best-effort (errors ignored); only the last must succeed.
func AutoFix(exec executor.Executor, result CheckResult) error {
	for i, cmd := range result.FixCmds {
		_, err := exec.Run(cmd[0], cmd[1:]...)
		if err != nil && i == len(result.FixCmds)-1 {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

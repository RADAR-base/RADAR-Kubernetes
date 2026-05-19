package prereqs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type CheckResult struct {
	Tool        string
	OK          bool
	Version     string
	Error       string
	InstallHint string
}

type toolDef struct {
	name         string
	args         []string
	parseVersion func(string) string
	installHint  string
	minVersion   string // e.g. "v3.0.0" — empty means no minimum
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
	},
	{
		name:         "helm",
		args:         []string{"version", "--short"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "https://helm.sh/docs/intro/install/",
		minVersion:   "v3",
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
	},
	{
		name:         "helm diff",
		args:         []string{"diff", "version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "helm plugin install https://github.com/databus23/helm-diff",
		minVersion:   "v3.9.12",
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
	},
	{
		name:         "openssl",
		args:         []string{"version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "brew install openssl  (macOS) or: apt install openssl",
	},
	{
		name:         "git",
		args:         []string{"--version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "https://git-scm.com/downloads",
	},
}

// meetsMinVersion returns true if version >= minVersion using simple semver comparison.
// Both version and minVersion should be in the form "vX", "vX.Y", or "vX.Y.Z".
// The "v" prefix is stripped before comparison.
func meetsMinVersion(version, minVersion string) bool {
	stripV := func(s string) string {
		return strings.TrimPrefix(strings.TrimSpace(s), "v")
	}
	// Extract only the leading version token (e.g. "3.9.12" from "3.9.12+g...")
	ver := strings.FieldsFunc(stripV(version), func(r rune) bool {
		return r == '+' || r == '-'
	})
	min := strings.FieldsFunc(stripV(minVersion), func(r rune) bool {
		return r == '+' || r == '-'
	})

	var verStr, minStr string
	if len(ver) > 0 {
		verStr = ver[0]
	}
	if len(min) > 0 {
		minStr = min[0]
	}

	vParts := strings.Split(verStr, ".")
	mParts := strings.Split(minStr, ".")

	// Pad to same length
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
	return true // equal
}

// Check runs all prerequisite checks and returns one result per tool.
func Check(exec executor.Executor) []CheckResult {
	results := make([]CheckResult, 0, len(tools))
	for _, t := range tools {
		toolName := strings.Fields(t.name)[0]
		// For compound tool names like "helm diff", toolName is "helm" and
		// t.args already contains the subcommand (e.g. ["diff", "version"]).
		args := t.args
		out, err := exec.Run(toolName, args...)
		if err != nil {
			// Extract the most useful line from the error (skip generic "exit status N" lines)
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
			r.Error = fmt.Sprintf("version %s found, need >= %s", version, t.minVersion)
			r.InstallHint = t.installHint
		}
		results = append(results, r)
	}
	return results
}

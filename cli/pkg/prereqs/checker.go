package prereqs

import (
	"fmt"
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
	},
	{
		name: "helm diff",
		args: []string{"diff", "version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "helm plugin install https://github.com/databus23/helm-diff",
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

// Check runs all prerequisite checks and returns one result per tool.
func Check(exec executor.Executor) []CheckResult {
	results := make([]CheckResult, 0, len(tools))
	for _, t := range tools {
		toolName := strings.Fields(t.name)[0]
		args := t.args
		if strings.Contains(t.name, " ") {
			// e.g. "helm diff" → run as "helm" with ["diff", "version"]
			extraArgs := strings.Fields(t.name)[1:]
			args = append(extraArgs, t.args...)
		}
		out, err := exec.Run(toolName, args...)
		if err != nil {
			results = append(results, CheckResult{
				Tool:        t.name,
				OK:          false,
				Error:       fmt.Sprintf("not found or failed: %s", err.Error()),
				InstallHint: t.installHint,
			})
			continue
		}
		results = append(results, CheckResult{
			Tool:    t.name,
			OK:      true,
			Version: t.parseVersion(out),
		})
	}
	return results
}

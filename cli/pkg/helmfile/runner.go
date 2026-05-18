package helmfile

import (
	"fmt"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type Runner struct {
	exec        executor.Executor
	workDir     string
	environment string
}

func NewRunner(exec executor.Executor, workDir, environment string) *Runner {
	return &Runner{exec: exec, workDir: workDir, environment: environment}
}

func (r *Runner) baseArgs(selector string) []string {
	args := []string{"--environment", r.environment}
	if selector != "" {
		args = append(args, "--selector", selector)
	}
	return args
}

// Diff returns the helmfile diff output as a string.
func (r *Runner) Diff(selector string) (string, error) {
	args := append(r.baseArgs(selector), "diff")
	out, err := r.exec.Run("helmfile", args...)
	if err != nil {
		// helmfile diff exits non-zero when there are differences — treat as success
		return out, nil
	}
	return out, nil
}

// Sync runs helmfile sync.
func (r *Runner) Sync(selector string, atomic bool) error {
	args := r.baseArgs(selector)
	if atomic {
		args = append(args, "--atomic")
	}
	args = append(args, "sync")
	_, err := r.exec.Run("helmfile", args...)
	return err
}

// SyncWithCallback runs helmfile sync and calls onLine for each output line.
func (r *Runner) SyncWithCallback(selector string, atomic bool, onLine func(string)) error {
	args := r.baseArgs(selector)
	if atomic {
		args = append(args, "--atomic")
	}
	args = append(args, "sync")

	out, err := r.exec.Run("helmfile", args...)
	for _, line := range strings.Split(out, "\n") {
		if line != "" {
			onLine(line)
		}
	}
	return err
}

// ParseDiffSummary counts changed/new/removed releases from helmfile diff output.
func ParseDiffSummary(diff string) (updated, installed, removed int) {
	for _, line := range strings.Split(diff, "\n") {
		lower := strings.ToLower(line)
		switch {
		case strings.Contains(lower, "has changed"):
			updated++
		case strings.Contains(lower, "installing"):
			installed++
		case strings.Contains(lower, "deleting"):
			removed++
		}
	}
	return
}

// ParseReleaseFromLine extracts a release name from a helmfile log line.
func ParseReleaseFromLine(line string) string {
	lower := strings.ToLower(line)
	for _, prefix := range []string{"release=", "upgrading release ", "installing release "} {
		if idx := strings.Index(lower, prefix); idx >= 0 {
			rest := line[idx+len(prefix):]
			if end := strings.IndexAny(rest, ", \t\n"); end >= 0 {
				return strings.ToLower(rest[:end])
			}
			return strings.ToLower(strings.TrimSpace(rest))
		}
	}
	return ""
}

// Version returns the installed helmfile version string.
func (r *Runner) Version() (string, error) {
	out, err := r.exec.Run("helmfile", "--version")
	if err != nil {
		return "", fmt.Errorf("helmfile not found: %w", err)
	}
	return strings.TrimSpace(out), nil
}

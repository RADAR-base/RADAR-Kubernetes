package helmfile

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/prereqs"
)

type Runner struct {
	exec        executor.Executor
	workDir     string
	environment string
	helmfileBin string // absolute path to the helmfile binary to use
}

func NewRunner(exec executor.Executor, workDir, environment string) *Runner {
	bin := "helmfile"
	if p := prereqs.ManagedHelmfilePath(); p != "" {
		bin = p
	}
	return &Runner{exec: exec, workDir: workDir, environment: environment, helmfileBin: bin}
}

// helmfile returns the binary name / path to use for all helmfile invocations.
func (r *Runner) helmfile() string { return r.helmfileBin }

func (r *Runner) baseArgs(selector string) []string {
	args := []string{"--environment", r.environment}
	if selector != "" {
		args = append(args, "--selector", selector)
	}
	return args
}

// Diff returns the helmfile diff output as a string.
// helmfile diff exits with code 2 when there are differences — that is treated as success.
// Any other non-zero exit code is a real error and is returned.
func (r *Runner) Diff(selector string) (string, error) {
	args := append(r.baseArgs(selector), "diff")
	out, err := r.exec.Run(r.helmfile(), args...)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 2 {
			// exit 2 = has differences, not a real error
			return out, nil
		}
		return out, err
	}
	return out, nil
}

// Template runs helmfile template (render only, no apply).
func (r *Runner) Template(selector string) (string, error) {
	args := append(r.baseArgs(selector), "template")
	out, err := r.exec.Run(r.helmfile(), args...)
	return out, err
}

// Sync runs helmfile sync.
func (r *Runner) Sync(selector string, atomic bool) error {
	args := append(r.baseArgs(selector), "sync")
	if atomic {
		args = append(args, "--sync-args", "--atomic")
	}
	_, err := r.exec.Run(r.helmfile(), args...)
	return err
}

// SyncWithCallback runs helmfile sync and calls onLine for each output line as it arrives.
func (r *Runner) SyncWithCallback(selector string, atomic bool, onLine func(string)) error {
	args := append(r.baseArgs(selector), "sync")
	if atomic {
		args = append(args, "--sync-args", "--atomic")
	}
	return r.exec.RunStreaming(r.helmfile(), args, onLine)
}

// ParseDiffSummary counts releases with changes from helmfile diff output.
// helmfile v0.x diff format per resource: "<ns>, <release>, <Kind> has been added/changed/deleted:"
// We track which releases have at least one changed resource and whether every
// resource for that release is "added" (new install) or "deleted" (removal).
func ParseDiffSummary(diff string) (updated, installed, removed int) {
	type releaseState struct {
		hasAdded   bool
		hasDeleted bool
		hasChanged bool
	}
	states := map[string]*releaseState{}
	order := []string{}

	currentRelease := ""
	for _, line := range strings.Split(diff, "\n") {
		// "Comparing release=foo, chart=..."
		if idx := strings.Index(line, "Comparing release="); idx >= 0 {
			rest := line[idx+len("Comparing release="):]
			if end := strings.IndexAny(rest, ", \t"); end >= 0 {
				rest = rest[:end]
			}
			currentRelease = strings.TrimSpace(rest)
			if _, ok := states[currentRelease]; !ok {
				states[currentRelease] = &releaseState{}
				order = append(order, currentRelease)
			}
			continue
		}
		if currentRelease == "" {
			continue
		}
		lower := strings.ToLower(line)
		s := states[currentRelease]
		switch {
		case strings.Contains(lower, "has been added"):
			s.hasAdded = true
		case strings.Contains(lower, "has been changed"):
			s.hasChanged = true
		case strings.Contains(lower, "has been deleted"):
			s.hasDeleted = true
		}
	}

	for _, rel := range order {
		s := states[rel]
		if !s.hasAdded && !s.hasChanged && !s.hasDeleted {
			continue
		}
		switch {
		case s.hasAdded && !s.hasChanged && !s.hasDeleted:
			installed++
		case s.hasDeleted && !s.hasAdded && !s.hasChanged:
			removed++
		default:
			updated++
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

// Version returns the helmfile version string from whichever binary is in use.
func (r *Runner) Version() (string, error) {
	out, err := r.exec.Run(r.helmfile(), "--version")
	if err != nil {
		return "", fmt.Errorf("helmfile not found: %w", err)
	}
	return strings.TrimSpace(out), nil
}

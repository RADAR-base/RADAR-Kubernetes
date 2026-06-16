package executor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Executor runs shell commands and returns combined output.
type Executor interface {
	Run(name string, args ...string) (string, error)
	RunStreaming(name string, args []string, onLine func(string)) error
}

// ShellExecutor runs real system commands.
type ShellExecutor struct {
	WorkDir  string
	ExtraEnv []string // extra entries appended to os.Environ() for every command
}

func (e *ShellExecutor) env() []string {
	if len(e.ExtraEnv) == 0 {
		return nil // let exec.Command inherit the process environment
	}
	// Build env by merging os.Environ() with ExtraEnv. For keys that appear in
	// both, ExtraEnv wins — we replace rather than append so that e.g. a
	// PATH=~/.radarctl/bin:... entry in ExtraEnv actually takes precedence over
	// the inherited PATH (on macOS the first PATH entry wins when there are
	// duplicates, which would negate our managed-binary prefix otherwise).
	overrides := make(map[string]string, len(e.ExtraEnv))
	for _, kv := range e.ExtraEnv {
		if idx := strings.IndexByte(kv, '='); idx >= 0 {
			overrides[kv[:idx]] = kv
		}
	}
	base := os.Environ()
	result := make([]string, 0, len(base)+len(e.ExtraEnv))
	for _, kv := range base {
		key := kv
		if idx := strings.IndexByte(kv, '='); idx >= 0 {
			key = kv[:idx]
		}
		if replacement, ok := overrides[key]; ok {
			result = append(result, replacement)
			delete(overrides, key) // emit once
		} else {
			result = append(result, kv)
		}
	}
	// Append any ExtraEnv keys that weren't already in the base environment.
	for _, kv := range e.ExtraEnv {
		if idx := strings.IndexByte(kv, '='); idx >= 0 {
			if _, stillPending := overrides[kv[:idx]]; stillPending {
				result = append(result, kv)
			}
		}
	}
	return result
}

func (e *ShellExecutor) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if e.WorkDir != "" {
		cmd.Dir = e.WorkDir
	}
	cmd.Env = e.env()
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return strings.TrimSpace(out.String()), fmt.Errorf("%w\n%s", err, strings.TrimSpace(errOut.String()))
	}
	// Some tools (e.g. java -version) write to stderr even on success.
	if stdout := strings.TrimSpace(out.String()); stdout != "" {
		return stdout, nil
	}
	return strings.TrimSpace(errOut.String()), nil
}

func (e *ShellExecutor) RunStreaming(name string, args []string, onLine func(string)) error {
	cmd := exec.Command(name, args...)
	if e.WorkDir != "" {
		cmd.Dir = e.WorkDir
	}
	cmd.Env = e.env()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		onLine(scanner.Text())
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%w\n%s", err, strings.TrimSpace(errOut.String()))
	}
	return nil
}

// MockExecutor records calls and returns configured responses for testing.
type MockExecutor struct {
	Responses map[string]string // key: "name args...", value: output
	Errors    map[string]error
	Calls     []string
}

func (m *MockExecutor) Run(name string, args ...string) (string, error) {
	key := strings.Join(append([]string{name}, args...), " ")
	// Also try with the basename so tests keyed on "helmfile …" match when the runner
	// passes an absolute managed-binary path like "/home/user/.radarctl/bin/helmfile …".
	baseKey := strings.Join(append([]string{filepath.Base(name)}, args...), " ")
	m.Calls = append(m.Calls, baseKey)
	if err, ok := m.Errors[key]; ok {
		return "", err
	}
	if err, ok := m.Errors[baseKey]; ok {
		return "", err
	}
	if out, ok := m.Responses[key]; ok {
		return out, nil
	}
	if out, ok := m.Responses[baseKey]; ok {
		return out, nil
	}
	return "", nil
}

func (m *MockExecutor) RunStreaming(name string, args []string, onLine func(string)) error {
	key := strings.Join(append([]string{name}, args...), " ")
	baseKey := strings.Join(append([]string{filepath.Base(name)}, args...), " ")
	m.Calls = append(m.Calls, baseKey)
	if err, ok := m.Errors[key]; ok {
		return err
	}
	if err, ok := m.Errors[baseKey]; ok {
		return err
	}
	out := m.Responses[key]
	if out == "" {
		out = m.Responses[baseKey]
	}
	for _, line := range strings.Split(out, "\n") {
		if line != "" {
			onLine(line)
		}
	}
	return nil
}

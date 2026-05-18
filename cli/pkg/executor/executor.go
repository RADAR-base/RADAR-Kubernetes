package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor runs shell commands and returns combined output.
type Executor interface {
	Run(name string, args ...string) (string, error)
}

// ShellExecutor runs real system commands.
type ShellExecutor struct{}

func (e *ShellExecutor) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w\n%s", err, strings.TrimSpace(errOut.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

// MockExecutor records calls and returns configured responses for testing.
type MockExecutor struct {
	Responses map[string]string // key: "name args...", value: output
	Errors    map[string]error
	Calls     []string
}

func (m *MockExecutor) Run(name string, args ...string) (string, error) {
	key := strings.Join(append([]string{name}, args...), " ")
	m.Calls = append(m.Calls, key)
	if err, ok := m.Errors[key]; ok {
		return "", err
	}
	if out, ok := m.Responses[key]; ok {
		return out, nil
	}
	return "", nil
}

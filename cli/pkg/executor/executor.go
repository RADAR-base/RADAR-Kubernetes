package executor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor runs shell commands and returns combined output.
type Executor interface {
	Run(name string, args ...string) (string, error)
	RunStreaming(name string, args []string, onLine func(string)) error
}

// ShellExecutor runs real system commands.
type ShellExecutor struct {
	WorkDir string
}

func (e *ShellExecutor) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if e.WorkDir != "" {
		cmd.Dir = e.WorkDir
	}
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
	m.Calls = append(m.Calls, key)
	if err, ok := m.Errors[key]; ok {
		return "", err
	}
	if out, ok := m.Responses[key]; ok {
		return out, nil
	}
	return "", nil
}

func (m *MockExecutor) RunStreaming(name string, args []string, onLine func(string)) error {
	key := strings.Join(append([]string{name}, args...), " ")
	m.Calls = append(m.Calls, key)
	if err, ok := m.Errors[key]; ok {
		return err
	}
	if out, ok := m.Responses[key]; ok {
		for _, line := range strings.Split(out, "\n") {
			if line != "" {
				onLine(line)
			}
		}
	}
	return nil
}

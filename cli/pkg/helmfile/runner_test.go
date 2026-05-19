package helmfile_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/helmfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeExitError returns a real *exec.ExitError with the given exit code.
func makeExitError(code int) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("exit %d", code))
	return cmd.Run()
}

func TestDiff(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"helmfile --environment default diff": "Release=mongodb has changed\nRelease=kafka unchanged",
		},
	}
	r := helmfile.NewRunner(mock, "/repo", "default")
	out, err := r.Diff("")
	require.NoError(t, err)
	assert.Contains(t, out, "mongodb has changed")
}

func TestDiff_WithSelector(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"helmfile --environment default --selector name=mongodb diff": "Release=mongodb has changed",
		},
	}
	r := helmfile.NewRunner(mock, "/repo", "default")
	out, err := r.Diff("name=mongodb")
	require.NoError(t, err)
	assert.Contains(t, out, "mongodb")
}

func TestDiff_ExitCode2_TreatedAsSuccess(t *testing.T) {
	exitErr := makeExitError(2)
	mock := &executor.MockExecutor{
		Errors: map[string]error{
			"helmfile --environment default diff": exitErr,
		},
	}
	r := helmfile.NewRunner(mock, "/repo", "default")
	_, err := r.Diff("")
	require.NoError(t, err, "exit code 2 should be treated as success (has differences)")
}

func TestDiff_ExitCode1_ReturnsError(t *testing.T) {
	exitErr := makeExitError(1)
	mock := &executor.MockExecutor{
		Errors: map[string]error{
			"helmfile --environment default diff": exitErr,
		},
	}
	r := helmfile.NewRunner(mock, "/repo", "default")
	_, err := r.Diff("")
	require.Error(t, err, "exit code 1 should be returned as a real error")
}

func TestSyncArgs(t *testing.T) {
	mock := &executor.MockExecutor{}
	r := helmfile.NewRunner(mock, "/repo", "default")
	_ = r.Sync("", false)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, "helmfile --environment default sync", mock.Calls[0])
}

func TestSyncWithCallback_CallsOnLinePerLine(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"helmfile --environment default sync": "line one\nline two\nline three",
		},
	}
	r := helmfile.NewRunner(mock, "/repo", "default")
	var lines []string
	err := r.SyncWithCallback("", false, func(line string) {
		lines = append(lines, line)
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"line one", "line two", "line three"}, lines)
}

func TestParseDiffSummary(t *testing.T) {
	diff := "Release=mongodb has changed\ninstalling Release=redis\ndeleting Release=old"
	updated, installed, removed := helmfile.ParseDiffSummary(diff)
	assert.Equal(t, 1, updated)
	assert.Equal(t, 1, installed)
	assert.Equal(t, 1, removed)
}

func TestParseReleaseFromLine(t *testing.T) {
	assert.Equal(t, "mongodb", helmfile.ParseReleaseFromLine("Comparing release=mongodb, chart=radar/mongodb"))
	assert.Equal(t, "", helmfile.ParseReleaseFromLine("some other log line"))
}

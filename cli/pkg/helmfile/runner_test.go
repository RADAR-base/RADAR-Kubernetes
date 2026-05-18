package helmfile_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/helmfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestSyncArgs(t *testing.T) {
	mock := &executor.MockExecutor{}
	r := helmfile.NewRunner(mock, "/repo", "default")
	_ = r.Sync("", false)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, "helmfile --environment default sync", mock.Calls[0])
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

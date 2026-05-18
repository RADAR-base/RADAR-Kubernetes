package kubectl_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockPodJSON() string {
	return `{
  "items": [
    {
      "metadata": {"name": "mongodb-0", "namespace": "default"},
      "status": {
        "phase": "Running",
        "containerStatuses": [{"ready": true, "restartCount": 0}]
      }
    }
  ]
}`
}

func TestGetPods(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl get pods -n default -l app=mongodb -o json": mockPodJSON(),
		},
	}
	r := kubectl.NewRunner(mock, "")
	pods, err := r.GetPods("default", "app=mongodb")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	assert.Equal(t, "mongodb-0", pods[0].Name)
	assert.True(t, pods[0].Ready)
}

func TestGetLogs(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl logs mongodb-0 -n default --tail=20": "log line 1\nlog line 2",
		},
	}
	r := kubectl.NewRunner(mock, "")
	logs, err := r.GetLogs("mongodb-0", "default", 20)
	require.NoError(t, err)
	assert.Contains(t, logs, "log line 1")
}

package status_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runningPodJSON(name string) string {
	return `{"items":[{"metadata":{"name":"` + name + `","namespace":"default"},"status":{"phase":"Running","containerStatuses":[{"ready":true,"restartCount":0}]}}]}`
}

func failingPodJSON(name string) string {
	return `{"items":[{"metadata":{"name":"` + name + `","namespace":"default"},"status":{"phase":"Failed","containerStatuses":[{"ready":false,"restartCount":5}],"message":"OOMKilled"}}]}`
}

func TestCollect_HealthyRelease(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl get pods -n default -l app.kubernetes.io/name=mongodb -o json": runningPodJSON("mongodb-0"),
		},
	}
	kr := kubectl.NewRunner(mock, "")
	releases := []string{"mongodb"}
	report, err := status.Collect(kr, "default", releases)
	require.NoError(t, err)
	require.Len(t, report.Groups, 1)
	require.Len(t, report.Groups[0].Releases, 1)
	assert.Equal(t, status.Healthy, report.Groups[0].Releases[0].Health)
	assert.Equal(t, 1, report.Healthy)
}

func TestCollect_DegradedRelease(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl get pods -n default -l app.kubernetes.io/name=ksql-server -o json": failingPodJSON("ksql-0"),
			"kubectl logs ksql-0 -n default --tail=20":                                  "java.lang.OutOfMemoryError",
		},
	}
	kr := kubectl.NewRunner(mock, "")
	report, err := status.Collect(kr, "default", []string{"ksql-server"})
	require.NoError(t, err)
	assert.Equal(t, 1, report.Degraded)
	assert.Equal(t, status.Degraded, report.Groups[0].Releases[0].Health)
}

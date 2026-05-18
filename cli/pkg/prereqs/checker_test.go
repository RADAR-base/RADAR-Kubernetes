package prereqs_test

import (
	"fmt"
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/prereqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_AllPresent(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			`kubectl version --client -o json`: `{"clientVersion":{"gitVersion":"v1.30.2"}}`,
			"helm version --short":             "v3.15.1",
			"helmfile --version":               "helmfile version v0.169.1",
			"helm diff version":                "3.9.12",
			"yq --version":                     "yq (https://github.com/mikefarah/yq/) version v4.44.3",
			"java -version":                    `openjdk version "21.0.1"`,
			"openssl version":                  "OpenSSL 3.1.4 24 Oct 2023",
			"git --version":                    "git version 2.42.0",
		},
	}
	results := prereqs.Check(mock)
	for _, r := range results {
		assert.True(t, r.OK, "expected %s to pass, got error: %s", r.Tool, r.Error)
	}
}

func TestCheck_MissingTool(t *testing.T) {
	mock := &executor.MockExecutor{
		Errors: map[string]error{
			"kubectl version --client -o json": fmt.Errorf("executable file not found in $PATH"),
		},
		Responses: map[string]string{
			"helm version --short": "v3.15.1",
			"helmfile --version":   "helmfile version v0.169.1",
			"helm diff version":    "3.9.12",
			"yq --version":         "yq (https://github.com/mikefarah/yq/) version v4.44.3",
			"java -version":        `openjdk version "21.0.1"`,
			"openssl version":      "OpenSSL 3.1.4",
			"git --version":        "git version 2.42.0",
		},
	}
	results := prereqs.Check(mock)
	var kubectlResult *prereqs.CheckResult
	for i := range results {
		if results[i].Tool == "kubectl" {
			kubectlResult = &results[i]
			break
		}
	}
	require.NotNil(t, kubectlResult)
	assert.False(t, kubectlResult.OK)
	assert.NotEmpty(t, kubectlResult.Error)
}

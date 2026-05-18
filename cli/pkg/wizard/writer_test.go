package wizard_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/wizard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestWriteAnswers_ProductionYAML(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "etc"), 0755))

	answers := &wizard.Answers{
		ServerName:      "radar.example.com",
		MaintainerEmail: "ops@example.com",
		KubeContext:     "my-cluster",
		Profile:         "production",
		Features:        []string{},
		AppliedMods:     []string{},
		Secrets:         map[string]string{},
	}

	w := wizard.NewWriter(dir)
	require.NoError(t, w.WriteAnswers(answers))

	data, err := os.ReadFile(filepath.Join(dir, "etc", "production.yaml"))
	require.NoError(t, err)

	var out map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &out))
	assert.Equal(t, "radar.example.com", out["server_name"])
	assert.Equal(t, "ops@example.com", out["maintainer_email"])
}

func TestWriteAnswers_DevProfile_DisablesTLS(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "etc"), 0755))

	answers := &wizard.Answers{
		ServerName:      "localhost",
		MaintainerEmail: "dev@example.com",
		KubeContext:     "k3d-local",
		Profile:         "dev",
		AppliedMods:     []string{"mods/minimal.yaml", "mods/disable_tls.yaml"},
		Secrets:         map[string]string{},
	}

	w := wizard.NewWriter(dir)
	require.NoError(t, w.WriteAnswers(answers))

	data, err := os.ReadFile(filepath.Join(dir, "etc", "production.yaml"))
	require.NoError(t, err)
	var out map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &out))
	assert.Equal(t, false, out["enable_tls"])
}

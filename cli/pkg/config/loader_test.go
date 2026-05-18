package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	yaml := `
server_name: test.example.com
maintainer_email: ops@example.com
kubeContext: my-cluster
enable_tls: true
kafka_num_brokers: 3
`
	path := filepath.Join(dir, "production.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0644))

	cfg, err := config.LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, "test.example.com", cfg.ServerName)
	assert.Equal(t, "ops@example.com", cfg.MaintainerEmail)
	assert.Equal(t, "my-cluster", cfg.KubeContext)
	assert.True(t, cfg.EnableTLS)
	assert.Equal(t, 3, cfg.KafkaNumBrokers)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := config.LoadConfig("/nonexistent/path.yaml")
	assert.Error(t, err)
}

func TestWriteConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.yaml")

	cfg := &config.Config{
		ServerName:      "written.example.com",
		MaintainerEmail: "admin@example.com",
		EnableTLS:       true,
	}
	require.NoError(t, config.WriteConfig(cfg, path))

	loaded, err := config.LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, "written.example.com", loaded.ServerName)
}

func TestLoadSecrets(t *testing.T) {
	dir := t.TempDir()
	yaml := `
fitbit_client_id: my-id
fitbit_client_secret: my-secret
`
	path := filepath.Join(dir, "secrets.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0644))

	sec, err := config.LoadSecrets(path)
	require.NoError(t, err)
	assert.Equal(t, "my-id", sec.FitbitClientID)
	assert.Equal(t, "my-secret", sec.FitbitClientSecret)
}

package config_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFeature_Fitbit(t *testing.T) {
	cfg := &config.Config{}
	err := config.ApplyFeature(cfg, "fitbit")
	require.NoError(t, err)
	assert.True(t, cfg.EnableFitbit)
}

func TestApplyFeature_Unknown(t *testing.T) {
	cfg := &config.Config{}
	err := config.ApplyFeature(cfg, "unknown_feature")
	assert.Error(t, err)
}

func TestFeatureSecretPrompts_Fitbit(t *testing.T) {
	prompts := config.FeatureSecretPrompts("fitbit")
	require.Len(t, prompts, 2)
	assert.Equal(t, "fitbit_client_id", prompts[0].SecretKey)
	assert.Equal(t, "fitbit_client_secret", prompts[1].SecretKey)
	assert.True(t, prompts[1].Mask)
}

func TestApplyDeploymentProfile_Dev(t *testing.T) {
	cfg := &config.Config{}
	mods := config.ApplyDeploymentProfile(cfg, "dev")
	assert.False(t, cfg.EnableTLS)
	assert.True(t, cfg.DevDeployment)
	assert.Contains(t, mods, "mods/minimal.yaml")
	assert.Contains(t, mods, "mods/localdev.yaml")
	assert.Contains(t, mods, "mods/disable_tls.yaml")
	assert.Contains(t, mods, "mods/fast_deploy.yaml")
}

func TestApplyDeploymentProfile_Production(t *testing.T) {
	cfg := &config.Config{}
	mods := config.ApplyDeploymentProfile(cfg, "production")
	assert.True(t, cfg.EnableTLS)
	assert.False(t, cfg.DevDeployment)
	assert.Empty(t, mods)
}

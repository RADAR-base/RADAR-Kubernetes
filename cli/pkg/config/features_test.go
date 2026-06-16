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
	require.NotNil(t, cfg.RadarFitbitConnector)
	assert.True(t, cfg.RadarFitbitConnector.Install)
}

func TestApplyFeature_Oura(t *testing.T) {
	cfg := &config.Config{}
	require.NoError(t, config.ApplyFeature(cfg, "oura"))
	require.NotNil(t, cfg.RadarOuraConnector)
	assert.True(t, cfg.RadarOuraConnector.Install)
	assert.True(t, cfg.RadarRestSourcesAuthBackend.Install)
}

func TestApplyFeature_ARMT(t *testing.T) {
	cfg := &config.Config{}
	require.NoError(t, config.ApplyFeature(cfg, "armt"))
	require.NotNil(t, cfg.RadarAppserver)
	assert.True(t, cfg.RadarAppserver.Install)
}

func TestApplyFeature_RealtimeDashboards(t *testing.T) {
	cfg := &config.Config{}
	require.NoError(t, config.ApplyFeature(cfg, "realtime_dashboards"))
	require.NotNil(t, cfg.KsqlServer)
	require.NotNil(t, cfg.RadarGrafana)
	require.NotNil(t, cfg.RadarJdbcConnectorRealtimeDashboard)
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
	config.ApplyDeploymentProfile(cfg, "dev")
	assert.False(t, cfg.EnableTLS)
	assert.True(t, cfg.DevDeployment)
}

func TestApplyDeploymentProfile_Production(t *testing.T) {
	cfg := &config.Config{}
	config.ApplyDeploymentProfile(cfg, "production")
	assert.True(t, cfg.EnableTLS)
	assert.False(t, cfg.DevDeployment)
}

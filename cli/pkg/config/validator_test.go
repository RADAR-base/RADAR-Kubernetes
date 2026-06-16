package config_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/stretchr/testify/assert"
)

func validConfig() *config.Config {
	return &config.Config{
		ServerName:      "radar.example.com",
		MaintainerEmail: "ops@example.com",
		KubeContext:     "my-cluster",
		KafkaNumBrokers: 3,
	}
}

func validSecrets() *config.Secrets {
	return &config.Secrets{}
}

func TestValidate_ValidConfig(t *testing.T) {
	result := config.Validate(validConfig(), validSecrets())
	assert.Empty(t, result.Errors)
}

func TestValidate_DefaultServerName(t *testing.T) {
	cfg := validConfig()
	cfg.ServerName = "example.com"
	result := config.Validate(cfg, validSecrets())
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "server_name", result.Errors[0].Field)
}

func TestValidate_DefaultEmail(t *testing.T) {
	cfg := validConfig()
	cfg.MaintainerEmail = "MAINTAINER_EMAIL@example.com"
	result := config.Validate(cfg, validSecrets())
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "maintainer_email", result.Errors[0].Field)
}

func TestValidate_FitbitEnabledMissingSecret(t *testing.T) {
	cfg := validConfig()
	cfg.RadarFitbitConnector = &config.ChartToggle{Install: true}
	sec := validSecrets()
	sec.FitbitClientID = ""
	result := config.Validate(cfg, sec)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "fitbit_client_id", result.Errors[0].Field)
}

func TestValidate_PlaceholderSecretWarning(t *testing.T) {
	sec := validSecrets()
	sec.FitbitClientID = "change_me"
	result := config.Validate(validConfig(), sec)
	assert.Len(t, result.Warnings, 1)
	assert.Equal(t, "fitbit_client_id", result.Warnings[0].Field)
}

func TestValidate_ConfluentMissingBootstrap(t *testing.T) {
	cfg := validConfig()
	cfg.ConfluentCloud = config.ConfluentCloudConfig{Enabled: true}
	result := config.Validate(cfg, validSecrets())
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "confluent_cloud.bootstrapServerurl", result.Errors[0].Field)
}

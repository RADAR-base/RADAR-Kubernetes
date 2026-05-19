package config

import "fmt"

// SecretPrompt describes a secret value the wizard must collect from the user.
type SecretPrompt struct {
	SecretKey string
	Label     string
	Mask      bool
}

type feature struct {
	apply   func(*Config)
	secrets []SecretPrompt
}

var featureRegistry = map[string]feature{
	"fitbit": {
		apply: func(c *Config) { c.EnableFitbit = true },
		secrets: []SecretPrompt{
			{SecretKey: "fitbit_client_id", Label: "Fitbit OAuth Client ID", Mask: false},
			{SecretKey: "fitbit_client_secret", Label: "Fitbit OAuth Client Secret", Mask: true},
		},
	},
	"garmin": {
		apply: func(c *Config) { c.EnableGarmin = true },
		secrets: []SecretPrompt{
			{SecretKey: "garmin_consumer_key", Label: "Garmin Consumer Key", Mask: false},
			{SecretKey: "garmin_consumer_secret", Label: "Garmin Consumer Secret", Mask: true},
		},
	},
	"redcap": {
		apply: func(c *Config) { c.EnableRedcap = true },
		secrets: []SecretPrompt{
			{SecretKey: "redcap_token", Label: "REDCap API Token", Mask: true},
		},
	},
	"kratos": {
		apply:   func(c *Config) { c.EnableKratos = true; c.EnableHydra = true },
		secrets: nil,
	},
}

// ApplyFeature sets the config fields required by a named feature.
func ApplyFeature(cfg *Config, name string) error {
	f, ok := featureRegistry[name]
	if !ok {
		return fmt.Errorf("unknown feature: %s", name)
	}
	f.apply(cfg)
	return nil
}

// FeatureSecretPrompts returns the secret prompts needed for a feature.
func FeatureSecretPrompts(name string) []SecretPrompt {
	f, ok := featureRegistry[name]
	if !ok {
		return nil
	}
	return f.secrets
}

// ApplyDeploymentProfile sets config values for a named profile and returns the mods to apply.
func ApplyDeploymentProfile(cfg *Config, profile string) []string {
	switch profile {
	case "dev":
		cfg.EnableTLS = false
		cfg.DevDeployment = true
		cfg.KafkaNumBrokers = 1
		cfg.KafkaNumReplicas = 1
		return []string{
			"mods/minimal.yaml",
			"mods/localdev.yaml",
			"mods/disable_tls.yaml",
			"mods/fast_deploy.yaml",
		}
	case "staging":
		cfg.EnableTLS = true
		cfg.DevDeployment = false
		cfg.KafkaNumBrokers = 1
		return []string{"mods/minimal.yaml"}
	default: // production
		cfg.EnableTLS = true
		cfg.DevDeployment = false
		cfg.KafkaNumBrokers = 3
		return nil
	}
}

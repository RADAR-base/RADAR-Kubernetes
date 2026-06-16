package config

import "fmt"

// SecretPrompt describes a value the wizard must collect from the user. Goes to
// secrets.yaml (Mask=true) or production.yaml as appropriate.
type SecretPrompt struct {
	SecretKey string
	Label     string
	Mask      bool
}

type feature struct {
	apply   func(*Config)
	secrets []SecretPrompt
}

// featureRegistry maps wizard feature ids to chart-toggle mutations + secret prompts.
var featureRegistry = map[string]feature{
	"fitbit": {
		apply: func(c *Config) {
			c.RadarFitbitConnector = &ChartToggle{Install: true}
			c.RadarRestSourcesAuthBackend = &ChartToggle{Install: true}
			c.RadarRestSourcesAuthorizer = &ChartToggle{Install: true}
		},
		secrets: []SecretPrompt{
			{SecretKey: "fitbit_client_id", Label: "Fitbit OAuth Client ID", Mask: false},
			{SecretKey: "fitbit_client_secret", Label: "Fitbit OAuth Client Secret", Mask: true},
		},
	},
	"garmin": {
		apply: func(c *Config) { c.RadarGarminConnector = &ChartToggle{Install: true} },
		secrets: []SecretPrompt{
			{SecretKey: "garmin_consumer_key", Label: "Garmin Consumer Key", Mask: false},
			{SecretKey: "garmin_consumer_secret", Label: "Garmin Consumer Secret", Mask: true},
		},
	},
	"oura": {
		apply: func(c *Config) {
			c.RadarOuraConnector = &ChartToggle{Install: true}
			c.RadarRestSourcesAuthBackend = &ChartToggle{Install: true}
			c.RadarRestSourcesAuthorizer = &ChartToggle{Install: true}
		},
		secrets: []SecretPrompt{
			{SecretKey: "oura_api_client", Label: "Oura OAuth Client ID", Mask: false},
			{SecretKey: "oura_api_secret", Label: "Oura OAuth Client Secret", Mask: true},
		},
	},
	"armt": {
		// aRMT questionnaires/tasks + HealthKit integration: needs the appserver.
		// Firebase JSON is collected separately by the wizard (see armt_firebase_json prompt).
		apply: func(c *Config) {
			if c.RadarAppserver == nil {
				c.RadarAppserver = &Appserver{Install: true}
			} else {
				c.RadarAppserver.Install = true
			}
		},
		secrets: []SecretPrompt{
			{SecretKey: "armt_firebase_json_path", Label: "Path to your Firebase service-account JSON (firebase-adminsdk.json)", Mask: false},
		},
	},
	"realtime_dashboards": {
		apply: func(c *Config) {
			c.KsqlServer = &ChartToggle{Install: true}
			c.RadarGrafana = &ChartToggle{Install: true}
			c.RadarJdbcConnectorRealtimeDashboard = &ChartToggle{Install: true}
		},
	},
	"kratos": {
		apply: func(c *Config) {
			c.RadarKratos = &ChartToggle{Install: true}
			c.RadarHydra = &ChartToggle{Install: true}
		},
	},
	"upload": {
		// Upload portal needs no extra secrets — credentials come from the management portal.
		apply: func(c *Config) {},
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

// FeatureSecretPrompts returns the secret/value prompts needed for a feature.
func FeatureSecretPrompts(name string) []SecretPrompt {
	f, ok := featureRegistry[name]
	if !ok {
		return nil
	}
	return f.secrets
}

// ApplyDeploymentProfile sets config values for a named profile.
// environments.yaml.gotmpl reads these flags and conditionally includes the
// appropriate mods/ overlay files, so no explicit mod list is needed here.
func ApplyDeploymentProfile(cfg *Config, profile string) {
	switch profile {
	case "demo":
		cfg.EnableTLS = false
		cfg.EnableLogging = false
		cfg.DevDeployment = true
		cfg.KafkaNumBrokers = 1
		cfg.KafkaNumReplicas = 1
		cfg.AtomicInstall = false
	case "dev":
		cfg.EnableTLS = false
		cfg.DevDeployment = true
		cfg.KafkaNumBrokers = 1
		cfg.KafkaNumReplicas = 1
	case "staging":
		cfg.EnableTLS = true
		cfg.DevDeployment = false
		cfg.KafkaNumBrokers = 1
	default: // production
		cfg.EnableTLS = true
		cfg.DevDeployment = false
		cfg.KafkaNumBrokers = 3
	}
}

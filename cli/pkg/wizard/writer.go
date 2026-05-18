package wizard

import (
	"os"
	"path/filepath"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"gopkg.in/yaml.v3"
)

// Answers holds all values collected by the wizard.
type Answers struct {
	ServerName      string
	MaintainerEmail string
	KubeContext     string
	Profile         string // "production", "staging", "dev"
	UseConfluent    bool
	ConfluentURL    string
	ConfluentKey    string
	ConfluentSecret string
	Features        []string          // enabled feature names e.g. ["fitbit", "redcap"]
	UseExternalS3   bool
	S3Bucket        string
	S3Region        string
	AppliedMods     []string
	Secrets         map[string]string // secret key → value collected from prompts
}

// Writer persists wizard answers to config files.
type Writer struct {
	repoRoot string
}

func NewWriter(repoRoot string) *Writer {
	return &Writer{repoRoot: repoRoot}
}

// WriteAnswers writes production.yaml and secrets.yaml from the collected answers.
func (w *Writer) WriteAnswers(a *Answers) error {
	if err := os.MkdirAll(filepath.Join(w.repoRoot, "etc"), 0755); err != nil {
		return err
	}

	cfg := &config.Config{
		ServerName:      a.ServerName,
		MaintainerEmail: a.MaintainerEmail,
		KubeContext:     a.KubeContext,
		AtomicInstall:   true,
		BaseTimeout:     90,
	}

	config.ApplyDeploymentProfile(cfg, a.Profile)

	for _, f := range a.Features {
		_ = config.ApplyFeature(cfg, f)
	}

	if a.UseConfluent {
		cfg.ConfluentCloud = true
	}

	cfgPath := filepath.Join(w.repoRoot, "etc", "production.yaml")
	if err := config.WriteConfig(cfg, cfgPath); err != nil {
		return err
	}

	sec := &config.Secrets{}
	if a.UseConfluent {
		sec.ConfluentCloud.BootstrapServer = a.ConfluentURL
		sec.ConfluentCloud.APIKey = a.ConfluentKey
		sec.ConfluentCloud.APISecret = a.ConfluentSecret
	}
	if v, ok := a.Secrets["fitbit_client_id"]; ok {
		sec.FitbitClientID = v
	}
	if v, ok := a.Secrets["fitbit_client_secret"]; ok {
		sec.FitbitClientSecret = v
	}
	if v, ok := a.Secrets["garmin_consumer_key"]; ok {
		sec.GarminConsumerKey = v
	}
	if v, ok := a.Secrets["garmin_consumer_secret"]; ok {
		sec.GarminConsumerSecret = v
	}
	if v, ok := a.Secrets["redcap_token"]; ok {
		sec.RedcapToken = v
	}

	secPath := filepath.Join(w.repoRoot, "etc", "secrets.yaml")
	return config.WriteSecrets(sec, secPath)
}

// WriteEnvironmentsYAML writes environments.yaml with the selected mods.
func (w *Writer) WriteEnvironmentsYAML(mods []string) error {
	type envConfig struct {
		Environments map[string]struct {
			Values []string `yaml:"values"`
		} `yaml:"environments"`
	}

	values := []string{
		"../etc/base.yaml",
		"../etc/production.yaml",
		"../etc/secrets.yaml",
	}
	for _, mod := range mods {
		values = append(values, "../"+mod)
	}

	cfg := envConfig{
		Environments: map[string]struct {
			Values []string `yaml:"values"`
		}{
			"default": {Values: values},
		},
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(w.repoRoot, "environments.yaml"), data, 0644)
}

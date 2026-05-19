package wizard

import (
	"fmt"
	"path/filepath"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/charmbracelet/huh"
)

type Mode string

const (
	ModeWizard      Mode = "wizard"
	ModeInteractive Mode = "interactive"
	ModeExpert      Mode = "expert"
)

// Run executes the wizard in the given mode and returns the collected Answers.
func Run(mode Mode, repoRoot string) (*Answers, error) {
	defaults := loadExistingAnswers(repoRoot)
	switch mode {
	case ModeWizard:
		return runWizard(defaults)
	case ModeInteractive:
		return runInteractive(defaults)
	case ModeExpert:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown wizard mode: %s", mode)
	}
}

func runWizard(a *Answers) (*Answers, error) {
	if a == nil {
		a = &Answers{Secrets: make(map[string]string)}
	}
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}

	steps := []func(*Answers) error{
		collectBasics,
		collectProfile,
		func(a *Answers) error {
			if a.Profile == "demo" || a.Profile == "dev" {
				a.UseConfluent = false
				return nil
			}
			return collectKafka(a)
		},
		func(a *Answers) error {
			if a.UseConfluent {
				return collectConfluentCredentials(a)
			}
			return nil
		},
		collectFeatures,
		collectFeatureSecrets,
		collectStorage,
		func(a *Answers) error {
			// Demo profile disables monitoring automatically — skip the question.
			if a.Profile == "demo" {
				a.EnablePrometheus = false
				a.EnableGraylog = false
				return nil
			}
			return collectMonitoring(a)
		},
		collectConfirm,
	}

	for _, step := range steps {
		if err := step(a); err != nil {
			return nil, fmt.Errorf("wizard step failed: %w", err)
		}
	}

	// Kratos/Hydra is always enabled — not a user choice.
	a.EnableKratos = true
	a.Features = appendIfMissing(a.Features, "kratos")

	cfg := &config.Config{}
	a.AppliedMods = config.ApplyDeploymentProfile(cfg, a.Profile)

	return a, nil
}

// SelectMode asks the user to pick wizard/interactive/expert.
func SelectMode() (Mode, error) {
	var chosen string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How would you like to configure your deployment?").
				Options(
					huh.NewOption("Guided wizard (recommended for first-time setup)", string(ModeWizard)),
					huh.NewOption("Interactive (configure every option explicitly)", string(ModeInteractive)),
					huh.NewOption("Expert (validate existing config files only)", string(ModeExpert)),
				).
				Value(&chosen),
		),
	).Run()
	return Mode(chosen), err
}

// loadExistingAnswers attempts to load existing config and secrets to pre-fill the wizard.
func loadExistingAnswers(repoRoot string) *Answers {
	a := &Answers{Secrets: make(map[string]string)}

	cfg, err := config.LoadConfig(filepath.Join(repoRoot, "etc", "production.yaml"))
	if err == nil {
		a.ServerName = cfg.ServerName
		a.MaintainerEmail = cfg.MaintainerEmail
		a.KubeContext = cfg.KubeContext
		a.UseConfluent = cfg.ConfluentCloud
		a.Profile = profileFrom(cfg)
		a.UseExternalS3 = cfg.UseExternalS3
		a.S3Bucket = cfg.S3Bucket
		a.S3Region = cfg.S3Region
		a.EnablePrometheus = cfg.EnablePrometheus
		a.EnableGraylog = cfg.EnableGraylog
		a.EnableKratos = cfg.EnableKratos
	}

	sec, err := config.LoadSecrets(filepath.Join(repoRoot, "etc", "secrets.yaml"))
	if err == nil {
		if sec.ConfluentCloud.BootstrapServer != "" {
			a.ConfluentURL = sec.ConfluentCloud.BootstrapServer
			a.ConfluentKey = sec.ConfluentCloud.APIKey
			a.ConfluentSecret = sec.ConfluentCloud.APISecret
		}
		if sec.FitbitClientID != "" {
			a.Secrets["fitbit_client_id"] = sec.FitbitClientID
		}
		if sec.FitbitClientSecret != "" {
			a.Secrets["fitbit_client_secret"] = sec.FitbitClientSecret
		}
		if sec.GarminConsumerKey != "" {
			a.Secrets["garmin_consumer_key"] = sec.GarminConsumerKey
		}
		if sec.GarminConsumerSecret != "" {
			a.Secrets["garmin_consumer_secret"] = sec.GarminConsumerSecret
		}
		if sec.RedcapToken != "" {
			a.Secrets["redcap_token"] = sec.RedcapToken
		}
		a.S3AccessKey = sec.S3AccessKey
		a.S3SecretKey = sec.S3SecretKey
	}

	return a
}

func appendIfMissing(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}

func profileFrom(cfg *config.Config) string {
	if cfg.DevDeployment {
		return "dev"
	}
	if cfg.KafkaNumBrokers <= 1 {
		return "staging"
	}
	return "production"
}

package wizard

import (
	"fmt"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
)

func collectBasics(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server hostname (your domain name)").
				Placeholder("radar.example.com").
				Value(&a.ServerName).
				Validate(func(s string) error {
					if s == "" || s == "example.com" {
						return fmt.Errorf("enter your actual domain name")
					}
					return nil
				}),
			huh.NewInput().
				Title("Maintainer email (for TLS certificate notifications)").
				Placeholder("ops@example.com").
				Value(&a.MaintainerEmail).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("email is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Kubernetes context").
				Placeholder("default").
				Value(&a.KubeContext),
		),
	).Run()
}

func collectProfile(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Deployment profile").
				Options(
					huh.NewOption("Production — full stack, TLS enabled", "production"),
					huh.NewOption("Staging — minimal resources, TLS enabled", "staging"),
					huh.NewOption("Local dev — no TLS, single replicas, fast probes", "dev"),
				).
				Value(&a.Profile),
		),
	).Run()
}

func collectKafka(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Use Confluent Cloud instead of local Kafka?").
				Value(&a.UseConfluent),
		),
	).Run()
}

func collectConfluentCredentials(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Confluent Cloud bootstrap server URL").
				Value(&a.ConfluentURL),
			huh.NewInput().
				Title("Confluent Cloud API key").
				Value(&a.ConfluentKey),
			huh.NewInput().
				Title("Confluent Cloud API secret").
				Password(true).
				Value(&a.ConfluentSecret),
		),
	).Run()
}

func collectFeatures(a *Answers) error {
	featureOptions := []huh.Option[string]{
		huh.NewOption("Fitbit", "fitbit"),
		huh.NewOption("Garmin", "garmin"),
		huh.NewOption("REDCap", "redcap"),
		huh.NewOption("Upload portal", "upload"),
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Which data source integrations do you want to enable?").
				Options(featureOptions...).
				Value(&a.Features),
		),
	).Run()
}

func collectFeatureSecrets(a *Answers) error {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}
	for _, feature := range a.Features {
		prompts := config.FeatureSecretPrompts(feature)
		for _, p := range prompts {
			val := a.Secrets[p.SecretKey]
			input := huh.NewInput().
				Title(p.Label).
				Value(&val)
			if p.Mask {
				input = input.Password(true)
			}
			if err := huh.NewForm(huh.NewGroup(input)).Run(); err != nil {
				return err
			}
			a.Secrets[p.SecretKey] = val
		}
	}
	return nil
}

func collectAuth(a *Answers) error {
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Ory Hydra/Kratos (OAuth2/OIDC identity management)?").
				Value(&a.EnableKratos),
		),
	).Run(); err != nil {
		return err
	}
	if a.EnableKratos {
		a.Features = append(a.Features, "kratos")
	}
	return nil
}

func collectStorage(a *Answers) error {
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Use external S3 instead of local Minio?").
			Value(&a.UseExternalS3),
	)).Run(); err != nil {
		return err
	}
	if !a.UseExternalS3 {
		return nil
	}
	return huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("S3 bucket name").Value(&a.S3Bucket),
		huh.NewInput().Title("S3 region").Value(&a.S3Region),
		huh.NewInput().Title("S3 access key").Value(&a.S3AccessKey),
		huh.NewInput().Title("S3 secret key").Password(true).Value(&a.S3SecretKey),
	)).Run()
}

func collectMonitoring(a *Answers) error {
	return huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Enable Prometheus + Grafana monitoring?").
			Value(&a.EnablePrometheus),
		huh.NewConfirm().
			Title("Enable Graylog + Elasticsearch logging?").
			Value(&a.EnableGraylog),
	)).Run()
}

func kafkaSummary(a *Answers) string {
	if a.UseConfluent {
		return fmt.Sprintf("Confluent Cloud (%s)", a.ConfluentURL)
	}
	return "Local Kafka"
}

func featuresSummary(a *Answers) string {
	var enabled []string
	for _, f := range a.Features {
		if f != "kratos" {
			enabled = append(enabled, f)
		}
	}
	if len(enabled) == 0 {
		return "none"
	}
	return strings.Join(enabled, ", ")
}

func collectConfirm(a *Answers) error {
	summary := fmt.Sprintf(
		"Review your configuration:\n\n"+
			"  Server:        %s\n"+
			"  Email:         %s\n"+
			"  Profile:       %s\n"+
			"  Kafka:         %s\n"+
			"  Features:      %s\n"+
			"  Auth (Kratos): %v\n"+
			"  External S3:   %v\n"+
			"  Prometheus:    %v\n"+
			"  Graylog:       %v\n\n"+
			"Files to write: etc/production.yaml, etc/secrets.yaml, environments.yaml",
		a.ServerName, a.MaintainerEmail, a.Profile,
		kafkaSummary(a), featuresSummary(a), a.EnableKratos, a.UseExternalS3, a.EnablePrometheus, a.EnableGraylog,
	)

	pterm.Info.Println(summary)

	var proceed bool
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Proceed and write configuration files?").
			Value(&proceed),
	)).Run(); err != nil {
		return err
	}
	if !proceed {
		return fmt.Errorf("setup cancelled")
	}
	return nil
}

// runInteractive presents every config field individually with current values pre-filled.
// No branching logic — all options are shown explicitly.
func runInteractive(a *Answers) (*Answers, error) {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}

	// Basics
	if err := collectBasics(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Profile
	if err := collectProfile(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Kafka — always show Confluent fields
	if err := collectKafka(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}
	// Always show Confluent credentials (no branching)
	if err := collectConfluentCredentials(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Features
	if err := collectFeatures(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Always show all known feature secret fields
	if err := collectAllFeatureSecrets(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Auth
	if err := collectAuth(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Storage
	if err := collectStorage(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Monitoring
	if err := collectMonitoring(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	// Confirm
	if err := collectConfirm(a); err != nil {
		return nil, fmt.Errorf("wizard step failed: %w", err)
	}

	cfg := &config.Config{}
	a.AppliedMods = config.ApplyDeploymentProfile(cfg, a.Profile)

	return a, nil
}

// collectAllFeatureSecrets shows prompts for all known feature secrets regardless of selection.
func collectAllFeatureSecrets(a *Answers) error {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}
	allFeatures := []string{"fitbit", "garmin", "redcap"}
	for _, feature := range allFeatures {
		prompts := config.FeatureSecretPrompts(feature)
		for _, p := range prompts {
			val := a.Secrets[p.SecretKey]
			input := huh.NewInput().
				Title(p.Label + " (leave blank if not using)").
				Value(&val)
			if p.Mask {
				input = input.Password(true)
			}
			if err := huh.NewForm(huh.NewGroup(input)).Run(); err != nil {
				return err
			}
			if val != "" {
				a.Secrets[p.SecretKey] = val
			}
		}
	}
	return nil
}

package wizard

import (
	"fmt"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/charmbracelet/huh"
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
			val := ""
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
	var enable bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Ory Hydra/Kratos (OAuth2/OIDC identity management)?").
				Value(&enable),
		),
	).Run(); err != nil {
		return err
	}
	if enable {
		a.Features = append(a.Features, "kratos")
	}
	return nil
}

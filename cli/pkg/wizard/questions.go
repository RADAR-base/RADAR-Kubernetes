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
				Title("Server hostname").
				Description("The domain where your RADAR platform will be reachable, e.g. radar.myinstitution.org.\nUsed to configure ingress routing and TLS certificates.").
				Placeholder("radar.example.com").
				Value(&a.ServerName).
				Validate(func(s string) error {
					if err := validateHostname(s); err != nil {
						return err
					}
					if s == "example.com" || strings.HasSuffix(s, ".example.com") {
						return fmt.Errorf("enter your actual domain name, not the example placeholder")
					}
					return nil
				}),
			huh.NewInput().
				Title("Maintainer email").
				Description("Used for Let's Encrypt TLS certificate expiry notifications.\nMust be a real, monitored address.").
				Placeholder("ops@example.com").
				Value(&a.MaintainerEmail).
				Validate(validateEmail),
			huh.NewInput().
				Title("Kubernetes context").
				Description("The kubectl context to use for all deployments.\nRun 'kubectl config get-contexts' to list available contexts.\nLeave blank to use the currently active context.").
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
				Description("Controls resource sizing, replica counts, and TLS behaviour.\nProduction: full replicas, TLS on — use for live studies.\nStaging: reduced resources, TLS on — use for pre-production testing.\nLocal dev: no TLS, single replicas, fast probes — use for development.\nLocal demo: lightest possible — no TLS, no monitoring/logging, fastest startup for demonstrations.").
				Options(
					huh.NewOption("Production — full stack, TLS enabled", "production"),
					huh.NewOption("Staging — minimal resources, TLS enabled", "staging"),
					huh.NewOption("Local dev — no TLS, single replicas, fast probes", "dev"),
					huh.NewOption("Local demo — no TLS, no monitoring, fastest startup", "demo"),
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
				Description("Confluent Cloud is a fully managed Kafka service — no brokers to operate.\nRequires an active Confluent Cloud account with a cluster already provisioned.\nChoose 'No' to deploy Kafka inside your cluster (recommended for self-hosted setups).").
				Value(&a.UseConfluent),
		),
	).Run()
}

func collectConfluentCredentials(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Confluent Cloud bootstrap server URL").
				Description("Found under your Confluent Cloud cluster → Cluster Settings → Endpoints.\nFormat: pkc-xxxxx.us-east-1.aws.confluent.cloud:9092").
				Value(&a.ConfluentURL).
				Validate(validateConfluentURL),
			huh.NewInput().
				Title("Confluent Cloud API key").
				Description("An API key with produce/consume permissions on the cluster.\nCreate one under Confluent Cloud → API Keys (not your login credentials).").
				Value(&a.ConfluentKey),
			huh.NewInput().
				Title("Confluent Cloud API secret").
				Description("The secret corresponding to the API key above.\nOnly shown once when created — retrieve from your records if lost.").
				Password(true).
				Value(&a.ConfluentSecret),
		),
	).Run()
}

func collectFeatures(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Data source integrations").
				Description("Select the wearable/data integrations to enable.\nEach selection will prompt for the required API credentials on the next screen.\nYou can enable more integrations later by re-running radarctl init.").
				Options(
					huh.NewOption("Fitbit — wearable activity and health data via OAuth2", "fitbit"),
					huh.NewOption("Garmin — wearable activity data via consumer API", "garmin"),
					huh.NewOption("REDCap — survey and clinical data via REST API", "redcap"),
					huh.NewOption("Upload portal — manual file uploads (no extra credentials needed)", "upload"),
				).
				Value(&a.Features),
		),
	).Run()
}

func collectFeatureSecrets(a *Answers) error {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}

	descriptions := map[string]string{
		"fitbit_client_id":       "The OAuth2 Client ID from your Fitbit developer application.\nCreate one at dev.fitbit.com → Manage → Register an App.",
		"fitbit_client_secret":   "The OAuth2 Client Secret for your Fitbit application.\nFound alongside the Client ID in your Fitbit developer app settings.",
		"garmin_consumer_key":    "The Consumer Key from your Garmin Health API application.\nRequest access at developer.garmin.com → Health API.",
		"garmin_consumer_secret": "The Consumer Secret for your Garmin Health API application.\nFound alongside the Consumer Key in your Garmin developer settings.",
		"redcap_token":           "A REDCap API token with export rights on the target project.\nGenerate one in REDCap under My Profile → API → Request API Token.",
	}

	for _, feature := range a.Features {
		prompts := config.FeatureSecretPrompts(feature)
		for _, p := range prompts {
			val := a.Secrets[p.SecretKey]
			input := huh.NewInput().
				Title(p.Label).
				Value(&val)
			if desc, ok := descriptions[p.SecretKey]; ok {
				input = input.Description(desc)
			}
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
				Description("Deploys Ory Kratos (user identity) and Ory Hydra (OAuth2/OIDC provider).\nRequired for participant authentication in the RADAR mobile app and management portal.\nRecommended for all production deployments.").
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
			Title("Use external S3 instead of local MinIO?").
			Description("MinIO is an S3-compatible object store deployed inside your cluster — no external account needed.\nChoose external S3 if you want to store collected data in AWS S3 or another S3-compatible service (e.g. Wasabi, Backblaze B2).").
			Value(&a.UseExternalS3),
	)).Run(); err != nil {
		return err
	}
	if !a.UseExternalS3 {
		return nil
	}
	return huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("S3 bucket name").
			Description("The name of the pre-existing S3 bucket RADAR will write data to.\nThe bucket must already exist — radarctl does not create it.").
			Value(&a.S3Bucket).
			Validate(validateS3Bucket),
		huh.NewInput().
			Title("S3 region").
			Description("The AWS region where the bucket is hosted, e.g. us-east-1, eu-west-2.\nFor non-AWS providers, use the region identifier from their documentation.").
			Value(&a.S3Region).
			Validate(validateS3Region),
		huh.NewInput().
			Title("S3 access key ID").
			Description("AWS access key ID with s3:PutObject and s3:GetObject permissions on the bucket.\nCreate one under AWS IAM → Users → Security credentials.").
			Value(&a.S3AccessKey).
			Validate(validateS3Key),
		huh.NewInput().
			Title("S3 secret access key").
			Description("The secret access key corresponding to the access key ID above.").
			Password(true).
			Value(&a.S3SecretKey),
	)).Run()
}

func collectMonitoring(a *Answers) error {
	return huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Enable Prometheus + Grafana monitoring?").
			Description("Deploys kube-prometheus-stack: Prometheus for metrics collection and Grafana for dashboards.\nIncludes pre-built dashboards for Kafka, JVM, and Kubernetes cluster health.\nRecommended for production — adds ~2 CPU / 4 GB RAM to cluster requirements.").
			Value(&a.EnablePrometheus),
		huh.NewConfirm().
			Title("Enable Graylog + Elasticsearch log aggregation?").
			Description("Deploys Graylog with an Elasticsearch backend for centralised log search and alerting.\nUseful for debugging and audit trails across all services.\nNote: Elasticsearch is resource-heavy — adds ~4 CPU / 8 GB RAM to requirements.").
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

	steps := []func() error{
		func() error { return collectBasics(a) },
		func() error { return collectProfile(a) },
		func() error { return collectKafka(a) },
		func() error { return collectConfluentCredentials(a) },
		func() error { return collectFeatures(a) },
		func() error { return collectAllFeatureSecrets(a) },
		func() error { return collectAuth(a) },
		func() error { return collectStorage(a) },
		func() error { return collectMonitoring(a) },
		func() error { return collectConfirm(a) },
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return nil, fmt.Errorf("wizard step failed: %w", err)
		}
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

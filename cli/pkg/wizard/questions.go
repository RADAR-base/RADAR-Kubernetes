package wizard

import (
	"fmt"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
)

func collectBasics(a *Answers) error {
	return runForm(huh.NewForm(
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
	))
}

func collectProfile(a *Answers) error {
	return runForm(huh.NewForm(
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
	))
}

func collectKafka(a *Answers) error {
	return runForm(huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Use Confluent Cloud instead of local Kafka?").
				Description("Confluent Cloud is a fully managed Kafka service — no brokers to operate.\nRequires an active Confluent Cloud account with a cluster already provisioned.\nChoose 'No' to deploy Kafka inside your cluster (recommended for self-hosted setups).").
				Value(&a.UseConfluent),
		),
	))
}

func collectConfluentCredentials(a *Answers) error {
	return runForm(huh.NewForm(
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
	))
}

func collectFeatures(a *Answers) error {
	return runForm(huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Data sources and dashboards").
				Description("Select the integrations to enable.\nEach selection will prompt for any required API credentials on the next screen.\nYou can enable more integrations later by re-running radarctl init.").
				Options(
					huh.NewOption("Fitbit — wearable activity and health data via OAuth2", "fitbit"),
					huh.NewOption("Garmin — wearable activity data via consumer API", "garmin"),
					huh.NewOption("Oura — Oura ring data via OAuth2", "oura"),
					huh.NewOption("aRMT + HealthKit — questionnaires/tasks app + iOS health data (enables appserver)", "armt"),
					huh.NewOption("Realtime dashboards — kSQL + Grafana + JDBC connector for live analytics", "realtime_dashboards"),
					huh.NewOption("Upload portal — manual file uploads (no extra credentials needed)", "upload"),
				).
				Value(&a.Features),
		),
	))
}

func collectFeatureSecrets(a *Answers) error {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}

	descriptions := map[string]string{
		"fitbit_client_id":         "The OAuth2 Client ID from your Fitbit developer application.\nCreate one at dev.fitbit.com → Manage → Register an App.",
		"fitbit_client_secret":     "The OAuth2 Client Secret for your Fitbit application.\nFound alongside the Client ID in your Fitbit developer app settings.",
		"garmin_consumer_key":      "The Consumer Key from your Garmin Health API application.\nRequest access at developer.garmin.com → Health API.",
		"garmin_consumer_secret":   "The Consumer Secret for your Garmin Health API application.\nFound alongside the Consumer Key in your Garmin developer settings.",
		"oura_api_client":          "The OAuth2 Client ID from your Oura Cloud application.\nCreate one at cloud.ouraring.com → My Applications → Create New Application.",
		"oura_api_secret":          "The OAuth2 Client Secret for your Oura application.\nShown once when the application is created — retrieve from your records if lost.",
		"armt_firebase_json_path": "Absolute path to your Firebase service-account JSON file.\n\n" +
			"What this is: the private key that authorises radar-appserver to act on behalf of\nthe Firebase project backing your aRMT mobile app.\n\n" +
			"Why it's needed:\n" +
			"  • Push notifications — delivering scheduled questionnaires/tasks to participants.\n" +
			"  • Crashlytics — capturing aRMT app crashes for debugging user-facing issues.\n" +
			"  • Event tracking — measuring engagement (questionnaire opens, completion rates).\n\n" +
			"How to get it: Firebase Console → Project settings → Service accounts → Generate\nnew private key. Save the .json file locally and paste its absolute path here.\n" +
			"The file will be copied into etc/radar-appserver/firebase-adminsdk.json.",
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
			if err := runForm(huh.NewForm(huh.NewGroup(input))); err != nil {
				return err
			}
			a.Secrets[p.SecretKey] = val
		}
	}
	return nil
}

func collectAuth(a *Answers) error {
	if err := runForm(huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Ory Hydra/Kratos (OAuth2/OIDC identity management)?").
				Description("Deploys Ory Kratos (user identity) and Ory Hydra (OAuth2/OIDC provider).\nRequired for participant authentication in the RADAR mobile app and management portal.\nRecommended for all production deployments.").
				Value(&a.EnableKratos),
		),
	)); err != nil {
		return err
	}
	if a.EnableKratos {
		a.Features = append(a.Features, "kratos")
	}
	return nil
}

func collectStorage(a *Answers) error {
	if err := runForm(huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Use external S3 instead of local MinIO?").
			Description("MinIO is an S3-compatible object store deployed inside your cluster — no external account needed.\nChoose external S3 if you want to store collected data in AWS S3 or another S3-compatible service (e.g. Wasabi, Backblaze B2).").
			Value(&a.UseExternalS3),
	))); err != nil {
		return err
	}
	if !a.UseExternalS3 {
		return nil
	}
	return runForm(huh.NewForm(huh.NewGroup(
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
	)))
}

func collectMonitoring(a *Answers) error {
	return runForm(huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Enable Prometheus + Grafana monitoring?").
			Description("Deploys kube-prometheus-stack: Prometheus for metrics collection and Grafana for dashboards.\nIncludes pre-built dashboards for Kafka, JVM, and Kubernetes cluster health.\nRecommended for production — adds ~2 CPU / 4 GB RAM to cluster requirements.").
			Value(&a.EnablePrometheus),
		huh.NewConfirm().
			Title("Enable Graylog + Elasticsearch log aggregation?").
			Description("Deploys Graylog with an Elasticsearch backend for centralised log search and alerting.\nUseful for debugging and audit trails across all services.\nNote: Elasticsearch is resource-heavy — adds ~4 CPU / 8 GB RAM to requirements.").
			Value(&a.EnableGraylog),
	)))
}

func storageSummary(a *Answers) string {
	if a.UseExternalS3 {
		return fmt.Sprintf("external S3 (bucket: %s, region: %s)", a.S3Bucket, a.S3Region)
	}
	return "local MinIO"
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
	monitoringSummary := fmt.Sprintf("Prometheus: %v, Graylog: %v", a.EnablePrometheus, a.EnableGraylog)
	if a.Profile == "demo" {
		monitoringSummary = "disabled (demo profile)"
	}

	summary := fmt.Sprintf(
		"Review your configuration:\n\n"+
			"  Server:      %s\n"+
			"  Email:       %s\n"+
			"  Profile:     %s\n"+
			"  Kafka:       %s\n"+
			"  Features:    %s\n"+
			"  Auth:        Ory Kratos + Hydra (always enabled)\n"+
			"  Storage:     %s\n"+
			"  Monitoring:  %s\n\n"+
			"Files to write: etc/production.yaml, etc/secrets.yaml",
		a.ServerName, a.MaintainerEmail, a.Profile,
		kafkaSummary(a), featuresSummary(a), storageSummary(a), monitoringSummary,
	)

	pterm.Info.Println(summary)

	var proceed bool
	if err := runForm(huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Proceed and write configuration files?").
			Value(&proceed),
	))); err != nil {
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
		func() error {
			if a.Profile == "demo" || a.Profile == "dev" {
				a.UseConfluent = false
				return nil
			}
			return collectKafka(a)
		},
		func() error {
			if a.UseConfluent {
				return collectConfluentCredentials(a)
			}
			return nil
		},
		func() error { return collectFeatures(a) },
		func() error { return collectAllFeatureSecrets(a) },
		func() error {
			if a.Profile == "demo" || a.Profile == "dev" {
				a.UseExternalS3 = false
				return nil
			}
			return collectStorage(a)
		},
		func() error {
			if a.Profile == "demo" {
				a.EnablePrometheus = false
				a.EnableGraylog = false
				return nil
			}
			return collectMonitoring(a)
		},
		func() error { return collectConfirm(a) },
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return nil, fmt.Errorf("wizard step failed: %w", err)
		}
	}

	a.EnableKratos = true
	a.Features = appendIfMissing(a.Features, "kratos")

	return a, nil
}

// collectAllFeatureSecrets shows prompts for all known feature secrets regardless of selection.
func collectAllFeatureSecrets(a *Answers) error {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}
	allFeatures := []string{"fitbit", "garmin", "oura", "armt"}
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
			if err := runForm(huh.NewForm(huh.NewGroup(input))); err != nil {
				return err
			}
			if val != "" {
				a.Secrets[p.SecretKey] = val
			}
		}
	}
	return nil
}

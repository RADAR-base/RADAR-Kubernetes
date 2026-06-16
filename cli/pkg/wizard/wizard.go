package wizard

import (
	"errors"
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
		return runWizard(defaults, repoRoot)
	case ModeInteractive:
		return runInteractive(defaults)
	case ModeExpert:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown wizard mode: %s", mode)
	}
}

// isLocalProfile returns true for profiles that target a local single-node cluster.
func isLocalProfile(p string) bool {
	return p == "dev" || p == "demo"
}

// errCancelled is returned when the user picks "Cancel setup" from the inter-step
// navigation prompt. runInit converts this to a clean exit.
var errCancelled = errors.New("setup cancelled by user")

// step describes one wizard screen plus its visibility predicate. Each step is
// idempotent: re-running it with the same Answers should just re-prompt with the
// previous values pre-filled (huh does this automatically since fields bind to
// Answers pointers).
type step struct {
	name string
	show func(*Answers) bool                          // nil → always shown
	run  func(*Answers, string) error                 // repoRoot passed through for side effects
}

// runWizard walks the user through a sequence of screens, with a small navigation
// prompt between each that lets them go back one screen at a time.
func runWizard(a *Answers, repoRoot string) (*Answers, error) {
	if a == nil {
		a = &Answers{Secrets: make(map[string]string)}
	}
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}

	notLocal := func(a *Answers) bool { return !isLocalProfile(a.Profile) }
	isLocal := func(a *Answers) bool { return isLocalProfile(a.Profile) }
	useConfluent := func(a *Answers) bool { return a.UseConfluent }
	wantsMonitoring := func(a *Answers) bool { return a.Profile != "demo" && !isLocalProfile(a.Profile) }

	steps := []step{
		{name: "profile", run: func(a *Answers, _ string) error { return collectProfile(a) }},
		{name: "local-cluster", show: isLocal, run: func(a *Answers, r string) error {
			// Auto-fill local-profile defaults, then make sure a cluster exists.
			if a.ServerName == "" {
				a.ServerName = "localhost"
			}
			if a.MaintainerEmail == "" {
				a.MaintainerEmail = "dev@localhost"
			}
			a.UseConfluent = false
			a.UseExternalS3 = false
			if a.Profile == "demo" {
				a.EnablePrometheus = false
				a.EnableGraylog = false
			}
			return ensureLocalCluster(a, r)
		}},
		{name: "basics", show: notLocal, run: func(a *Answers, _ string) error { return collectBasics(a) }},
		{name: "kafka", show: notLocal, run: func(a *Answers, _ string) error { return collectKafka(a) }},
		{name: "confluent-credentials", show: useConfluent, run: func(a *Answers, _ string) error { return collectConfluentCredentials(a) }},
		{name: "features", run: func(a *Answers, _ string) error { return collectFeatures(a) }},
		{name: "feature-secrets", run: func(a *Answers, _ string) error { return collectFeatureSecrets(a) }},
		{name: "storage", show: notLocal, run: func(a *Answers, _ string) error { return collectStorage(a) }},
		{name: "monitoring", show: wantsMonitoring, run: func(a *Answers, _ string) error { return collectMonitoring(a) }},
		{name: "confirm", run: func(a *Answers, _ string) error { return collectConfirm(a) }},
	}

	i := 0
	for i < len(steps) {
		s := steps[i]
		if s.show != nil && !s.show(a) {
			i++
			continue
		}
		err := s.run(a, repoRoot)
		if errors.Is(err, huh.ErrUserAborted) {
			// Esc was pressed inside the form — interpret as "go back". From the
			// first visible step there is nowhere to go back to, so treat that
			// as a clean cancel.
			prev := previousVisibleStep(steps, a, i)
			if prev == i {
				return nil, errCancelled
			}
			i = prev
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("step %q failed: %w", s.name, err)
		}
		i++
	}

	// Kratos/Hydra is always enabled — not a user choice.
	a.EnableKratos = true
	a.Features = appendIfMissing(a.Features, "kratos")

	return a, nil
}

// previousVisibleStep returns the index of the most recent visible step before
// `from`. If we're already at the first visible step (or before it), we stay put.
func previousVisibleStep(steps []step, a *Answers, from int) int {
	for j := from - 1; j >= 0; j-- {
		if steps[j].show == nil || steps[j].show(a) {
			return j
		}
	}
	return from
}

// SelectMode asks the user to pick wizard/interactive/expert.
func SelectMode() (Mode, error) {
	var chosen string
	err := runForm(huh.NewForm(
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
	))
	if errors.Is(err, huh.ErrUserAborted) {
		return "", errCancelled
	}
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
		a.UseConfluent = cfg.ConfluentCloud.Enabled
		a.Profile = profileFrom(cfg)
		a.UseExternalS3 = cfg.UseExternalS3
		a.S3Bucket = cfg.S3Bucket
		a.S3Region = cfg.S3Region
		a.EnablePrometheus = cfg.KubePrometheusStack != nil && cfg.KubePrometheusStack.Install
		a.EnableGraylog = cfg.Graylog != nil && cfg.Graylog.Install
		a.EnableKratos = cfg.RadarKratos != nil && cfg.RadarKratos.Install
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
		if sec.OuraAPIClient != "" {
			a.Secrets["oura_api_client"] = sec.OuraAPIClient
		}
		if sec.OuraAPISecret != "" {
			a.Secrets["oura_api_secret"] = sec.OuraAPISecret
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

// IsCancelled reports whether err came from the user choosing "Cancel setup" in
// the navigation prompt (so cmd/init.go can exit cleanly without a stack trace).
func IsCancelled(err error) bool {
	return errors.Is(err, errCancelled)
}

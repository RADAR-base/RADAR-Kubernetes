package wizard

import (
	"fmt"

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
	switch mode {
	case ModeWizard:
		return runWizard()
	case ModeInteractive:
		return runWizard() // interactive uses same flow in v1
	case ModeExpert:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown wizard mode: %s", mode)
	}
}

func runWizard() (*Answers, error) {
	a := &Answers{Secrets: make(map[string]string)}

	steps := []func(*Answers) error{
		collectBasics,
		collectProfile,
		collectKafka,
		func(a *Answers) error {
			if a.UseConfluent {
				return collectConfluentCredentials(a)
			}
			return nil
		},
		collectFeatures,
		collectFeatureSecrets,
		collectAuth,
	}

	for _, step := range steps {
		if err := step(a); err != nil {
			return nil, fmt.Errorf("wizard step failed: %w", err)
		}
	}

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

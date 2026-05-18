package cmd

import (
	"fmt"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/prereqs"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/wizard"
	"github.com/spf13/cobra"
)

var skipPrereqs bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up a new RADAR-Kubernetes deployment interactively",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&skipPrereqs, "skip-prereqs", false, "Skip prerequisite checks (advanced users only)")
}

func runInit(_ *cobra.Command, _ []string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	if !skipPrereqs {
		output.Header("Checking prerequisites...")
		exec := &executor.ShellExecutor{}
		results := prereqs.Check(exec)
		failed := 0
		for _, r := range results {
			if r.OK {
				output.Success(fmt.Sprintf("%-12s %s", r.Tool, r.Version))
			} else {
				output.Error(fmt.Sprintf("%-12s not found — %s", r.Tool, r.InstallHint))
				failed++
			}
		}
		if failed > 0 {
			return fmt.Errorf("%d prerequisite(s) missing. Install them and re-run radarctl init", failed)
		}
	}

	mode, err := wizard.SelectMode()
	if err != nil {
		return err
	}

	if mode == wizard.ModeExpert {
		output.Info("Expert mode: running config validation only")
		return runValidate(nil, nil)
	}

	output.Header("Configuring RADAR-Kubernetes")
	answers, err := wizard.Run(mode, repoRoot)
	if err != nil {
		return fmt.Errorf("wizard failed: %w", err)
	}

	w := wizard.NewWriter(repoRoot)
	if err := w.WriteAnswers(answers); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	if err := w.WriteEnvironmentsYAML(answers.AppliedMods); err != nil {
		return fmt.Errorf("writing environments.yaml: %w", err)
	}

	output.Success("Configuration written to etc/production.yaml, etc/secrets.yaml, environments.yaml")
	output.Info("Run `radarctl deploy` to deploy the stack")
	return nil
}

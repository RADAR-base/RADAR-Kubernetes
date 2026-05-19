package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			return &ExitError{Code: 3, Message: fmt.Sprintf("%d prerequisite(s) missing. Install them and re-run radarctl init", failed)}
		}
	}

	// Check for saved wizard state and offer to resume.
	statePath := filepath.Join(repoRoot, ".radarctl-state.yaml")
	var answers *wizard.Answers
	if saved, err := wizard.LoadState(statePath); err == nil {
		output.Info("Previous wizard answers found. Resume from saved state? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "y" || line == "yes" {
			answers = saved
			output.Info("Resuming from saved state — skipping to confirmation step")
		} else {
			_ = os.Remove(statePath)
		}
	}

	if answers == nil {
		mode, err := wizard.SelectMode()
		if err != nil {
			return err
		}

		if mode == wizard.ModeExpert {
			output.Info("Expert mode: running config validation only")
			return validateConfig(repoRoot)
		}

		output.Header("Configuring RADAR-Kubernetes")
		answers, err = wizard.Run(mode, repoRoot)
		if err != nil {
			return fmt.Errorf("wizard failed: %w", err)
		}
	}

	// Save state after wizard completes successfully.
	if err := wizard.SaveState(answers, statePath); err != nil {
		output.Warning(fmt.Sprintf("Could not save wizard state: %s", err))
	}

	w := wizard.NewWriter(repoRoot)
	if err := w.WriteAnswers(answers); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	if err := w.WriteEnvironmentsYAML(answers.AppliedMods); err != nil {
		return fmt.Errorf("writing environments.yaml: %w", err)
	}

	output.Success("Configuration written to etc/production.yaml, etc/secrets.yaml, environments.yaml")

	// Run bin/keystore-init if it exists.
	keystoreInit := filepath.Join(repoRoot, "bin", "keystore-init")
	if _, err := os.Stat(keystoreInit); err == nil {
		output.Info("Running bin/keystore-init...")
		exec := &executor.ShellExecutor{WorkDir: repoRoot}
		if _, err := exec.Run(keystoreInit); err != nil {
			output.Warning(fmt.Sprintf("keystore-init failed: %s", err))
			output.Warning("You may need to run bin/keystore-init manually")
		} else {
			output.Success("Keystores initialized")
		}
	}

	// Remove state file after successful write.
	_ = os.Remove(statePath)

	output.Info("Run `radarctl deploy` to deploy the stack")
	return nil
}

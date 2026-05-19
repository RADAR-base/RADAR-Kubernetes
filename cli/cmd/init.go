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
	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
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

func printPrereqResults(results []prereqs.CheckResult) int {
	failed := 0
	for _, r := range results {
		if r.OK {
			output.Success(fmt.Sprintf("%-12s %s", r.Tool, r.Version))
		} else {
			msg := r.Error
			if msg == "" {
				msg = "not found"
			}
			output.Error(fmt.Sprintf("%-12s %s", r.Tool, msg))
			failed++
		}
	}
	return failed
}

func runPrereqsCheck(exec *executor.ShellExecutor) error {
	output.Header("Checking prerequisites...")
	results := prereqs.Check(exec)
	failed := printPrereqResults(results)
	if failed == 0 {
		return nil
	}

	// Count how many we can auto-fix.
	fixable := 0
	for _, r := range results {
		if !r.OK && len(r.FixCmds) > 0 {
			fixable++
		}
	}
	if fixable == 0 {
		return &ExitError{Code: 3, Message: fmt.Sprintf("%d prerequisite(s) missing", failed)}
	}

	var doInstall bool
	label := fmt.Sprintf("Install/upgrade %d missing prerequisite(s) now? (uses Homebrew)", fixable)
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title(label).Value(&doInstall),
	)).Run(); err != nil || !doInstall {
		return &ExitError{Code: 3, Message: fmt.Sprintf("%d prerequisite(s) missing", failed)}
	}

	for _, r := range results {
		if r.OK || len(r.FixCmds) == 0 {
			continue
		}
		action := "Installing"
		if r.NeedsUpgrade {
			action = "Upgrading"
		}
		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("%s %s...", action, r.Tool))
		if err := prereqs.AutoFix(exec, r); err != nil {
			spinner.Fail(fmt.Sprintf("%s %s failed: %s", action, r.Tool, err))
		} else {
			spinner.Success(fmt.Sprintf("%s %s complete", action, r.Tool))
		}
	}

	pterm.Println()
	output.Header("Re-checking prerequisites...")
	results = prereqs.Check(exec)
	failed = printPrereqResults(results)
	if failed > 0 {
		return &ExitError{Code: 3, Message: fmt.Sprintf("%d prerequisite(s) still missing after install", failed)}
	}
	return nil
}

func runInit(_ *cobra.Command, _ []string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	if !skipPrereqs {
		exec := &executor.ShellExecutor{}
		if err := runPrereqsCheck(exec); err != nil {
			return err
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

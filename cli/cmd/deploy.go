package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/helmfile"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	deployDiff     bool
	deployDryRun   bool
	deploySelector string
	deployYes      bool
	deployNoAtomic bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the RADAR stack (wraps helmfile sync)",
	RunE:  runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVar(&deployDiff, "diff", false, "Show what would change without applying")
	deployCmd.Flags().BoolVar(&deployDryRun, "dry-run", false, "Render templates only, do not apply")
	deployCmd.Flags().StringVar(&deploySelector, "selector", "", "Deploy only releases matching this selector (e.g. name=mongodb)")
	deployCmd.Flags().BoolVarP(&deployYes, "yes", "y", false, "Skip confirmation prompt")
	deployCmd.Flags().BoolVar(&deployNoAtomic, "no-atomic", false, "Disable automatic rollback on failure")
}

func runDeploy(_ *cobra.Command, _ []string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	exec := &executor.ShellExecutor{}
	hfRunner := helmfile.NewRunner(exec, repoRoot, "default")

	output.Info("Validating configuration...")
	cfg, err := config.LoadConfig(filepath.Join(repoRoot, "etc", "production.yaml"))
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	sec, err := config.LoadSecrets(filepath.Join(repoRoot, "etc", "secrets.yaml"))
	if err != nil {
		return fmt.Errorf("loading secrets: %w", err)
	}
	result := config.Validate(cfg, sec)
	for _, w := range result.Warnings {
		output.Warning(fmt.Sprintf("[%s] %s", w.Field, w.Message))
	}
	if !result.IsValid() {
		for _, e := range result.Errors {
			output.Error(fmt.Sprintf("[%s] %s", e.Field, e.Message))
		}
		os.Exit(2)
	}

	if deployDiff || deployDryRun {
		output.Info("Computing diff...")
		diff, err := hfRunner.Diff(deploySelector)
		if err != nil {
			return fmt.Errorf("helmfile diff: %w", err)
		}
		fmt.Println(diff)
		return nil
	}

	output.Info("Computing changes...")
	diff, _ := hfRunner.Diff(deploySelector)
	updated, installed, removed := helmfile.ParseDiffSummary(diff)
	pterm.Info.Printf("%d releases will be updated, %d installed, %d removed\n", updated, installed, removed)

	if !deployYes {
		pterm.Print("Proceed with deployment? [y/N] ")
		var input string
		fmt.Scanln(&input)
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			output.Info("Deployment cancelled")
			return nil
		}
	}

	output.Header("Deploying RADAR stack...")
	progress, _ := pterm.DefaultSpinner.Start("Starting deployment...")

	releaseStatuses := map[string]string{}

	err = hfRunner.SyncWithCallback(deploySelector, !deployNoAtomic, func(line string) {
		release := helmfile.ParseReleaseFromLine(line)
		if release != "" {
			releaseStatuses[release] = "syncing"
			progress.UpdateText(fmt.Sprintf("Syncing %s...", release))
		}
	})

	if err != nil {
		progress.Fail("Deployment failed")
		output.Error(err.Error())
		if outputFormat == string(output.JSON) {
			output.PrintJSON(map[string]any{"status": "failed", "error": err.Error()})
		}
		os.Exit(1)
	}

	progress.Success("Deployment complete")
	output.Success("All releases synced successfully")

	if outputFormat == string(output.JSON) {
		output.PrintJSON(map[string]any{"status": "healthy", "summary": map[string]int{
			"synced": len(releaseStatuses),
		}})
	}
	return nil
}

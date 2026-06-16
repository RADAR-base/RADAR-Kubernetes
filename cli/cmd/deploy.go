package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/helmfile"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/prereqs"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/status"
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

	// Guard: environments.yaml must exist (created by bin/init via radarctl init).
	if _, err := os.Stat(filepath.Join(repoRoot, "environments.yaml")); os.IsNotExist(err) {
		return &ExitError{Code: 2, Message: "environments.yaml not found — run `radarctl init` first to set up configuration files"}
	}

	exec := &executor.ShellExecutor{WorkDir: repoRoot, ExtraEnv: prereqs.ManagedBinEnv()}
	hfRunner := helmfile.NewRunner(exec, repoRoot, "default")

	// 1. Validate config
	if !isJSON() {
		output.Info("Validating configuration...")
	}
	cfg, err := config.LoadConfig(filepath.Join(repoRoot, "etc", "production.yaml"))
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	sec, err := config.LoadSecrets(filepath.Join(repoRoot, "etc", "secrets.yaml"))
	if err != nil {
		return fmt.Errorf("loading secrets: %w", err)
	}
	result := config.Validate(cfg, sec)
	if !isJSON() {
		for _, w := range result.Warnings {
			output.Warning(fmt.Sprintf("[%s] %s", w.Field, w.Message))
		}
	}
	if !result.IsValid() {
		if !isJSON() {
			for _, e := range result.Errors {
				output.Error(fmt.Sprintf("[%s] %s", e.Field, e.Message))
			}
		} else {
			output.PrintJSON(map[string]any{
				"status": "invalid",
				"errors": result.Errors,
			})
		}
		return &ExitError{Code: 2, Message: "configuration is invalid"}
	}

	// 2. Handle --diff
	if deployDiff {
		if !isJSON() {
			output.Info("Computing diff...")
		}
		diff, err := hfRunner.Diff(deploySelector)
		if err != nil {
			return fmt.Errorf("helmfile diff: %w", err)
		}
		fmt.Println(diff)
		return nil
	}

	// 2b. Handle --dry-run
	if deployDryRun {
		if !isJSON() {
			output.Info("Rendering templates...")
		}
		rendered, err := hfRunner.Template(deploySelector)
		if err != nil {
			return fmt.Errorf("helmfile template: %w", err)
		}
		fmt.Println(rendered)
		return nil
	}

	// 3. Show diff summary + confirmation (skip in JSON mode)
	if !isJSON() {
		output.Info("Computing changes...")
		diff, diffErr := hfRunner.Diff(deploySelector)
		if diffErr != nil {
			output.Warning(fmt.Sprintf("Could not compute diff: %s", diffErr))
		} else {
			updated, installed, removed := helmfile.ParseDiffSummary(diff)
			pterm.Info.Printf("%d releases will be updated, %d installed, %d removed\n", updated, installed, removed)
		}

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
	}

	// 4. Sync with spinner (only start spinner when not JSON mode)
	releaseStatuses := map[string]string{}

	var progress *pterm.SpinnerPrinter
	if !isJSON() {
		progress, _ = pterm.DefaultSpinner.Start("Starting deployment...")
	}

	err = hfRunner.SyncWithCallback(deploySelector, !deployNoAtomic, func(line string) {
		rel := helmfile.ParseReleaseFromLine(line)
		if rel != "" {
			releaseStatuses[rel] = "syncing"
			if progress != nil {
				progress.UpdateText(fmt.Sprintf("Syncing %s...", rel))
			}
		}
	})

	if err != nil {
		if progress != nil {
			progress.Fail("Deployment failed")
		}
		if !isJSON() {
			output.Error(err.Error())
		} else {
			output.PrintJSON(map[string]any{"status": "failed", "error": err.Error()})
		}
		return &ExitError{Code: 1, Message: fmt.Sprintf("deployment failed: %s", err)}
	}

	if progress != nil {
		progress.Success("Deployment complete")
	}

	// 5. Post-deploy health check
	kr := kubectl.NewRunner(exec, resolvedKubeContext(repoRoot))
	report, err := status.Collect(kr, "default", knownReleases)
	if err != nil {
		return &ExitError{Code: 4, Message: fmt.Sprintf("cluster unreachable: %s", err)}
	}

	// Flatten releases from all groups
	var releases []status.ReleaseStatus
	for _, g := range report.Groups {
		releases = append(releases, g.Releases...)
	}

	overallStatus := "healthy"
	if report.Degraded > 0 {
		overallStatus = "degraded"
	}

	if isJSON() {
		output.PrintJSON(map[string]any{
			"status":   overallStatus,
			"releases": releases,
			"summary": map[string]int{
				"healthy": report.Healthy,
				"failed":  report.Degraded,
				"pending": report.Warning,
			},
		})
	} else {
		output.Success("All releases synced successfully")
		pterm.Printf("✓ %d releases healthy\n", report.Healthy)
		for _, g := range report.Groups {
			for _, r := range g.Releases {
				if r.Health == status.Degraded {
					msg := r.Error
					if msg == "" {
						msg = r.Message
					}
					pterm.Printf("✗ %d degraded → %s (%s)\n", report.Degraded, r.Name, msg)
				}
			}
		}
	}

	if report.Degraded > 0 {
		return &ExitError{Code: 1, Message: "one or more releases degraded"}
	}
	return nil
}

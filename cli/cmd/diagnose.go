package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/status"
	"github.com/spf13/cobra"
)

type DiagnoseReport struct {
	Timestamp  string               `json:"timestamp"`
	Validation output.ValidationJSON `json:"validation"`
	Status     *status.Report       `json:"status"`
	Logs       map[string]string    `json:"logs"`
}

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Collect a full diagnostic snapshot for debugging or agentic loops",
	RunE:  runDiagnose,
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
}

func runDiagnose(_ *cobra.Command, _ []string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	report := DiagnoseReport{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Logs:      map[string]string{},
	}

	cfg, cfgErr := config.LoadConfig(filepath.Join(repoRoot, "etc", "production.yaml"))
	sec, secErr := config.LoadSecrets(filepath.Join(repoRoot, "etc", "secrets.yaml"))
	if cfgErr == nil && secErr == nil {
		result := config.Validate(cfg, sec)
		report.Validation.Valid = result.IsValid()
		for _, e := range result.Errors {
			report.Validation.Errors = append(report.Validation.Errors,
				fmt.Sprintf("%s: %s", e.Field, e.Message))
		}
		for _, w := range result.Warnings {
			report.Validation.Warnings = append(report.Validation.Warnings,
				fmt.Sprintf("%s: %s", w.Field, w.Message))
		}
	} else {
		report.Validation.Valid = false
		report.Validation.Errors = []string{"could not load config files"}
	}

	exec := &executor.ShellExecutor{}
	kr := kubectl.NewRunner(exec, kubeContext)
	statusReport, err := status.Collect(kr, "default", knownReleases)
	if err == nil {
		report.Status = statusReport
		for _, group := range statusReport.Groups {
			for _, r := range group.Releases {
				if r.Logs != "" {
					report.Logs[r.Name] = r.Logs
				}
			}
		}
	}

	output.PrintJSON(report)
	return nil
}

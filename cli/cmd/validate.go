package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate production.yaml and secrets.yaml before deploying",
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, _ []string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	cfgPath := filepath.Join(repoRoot, "etc", "production.yaml")
	secPath := filepath.Join(repoRoot, "etc", "secrets.yaml")

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("loading production.yaml: %w", err)
	}
	sec, err := config.LoadSecrets(secPath)
	if err != nil {
		return fmt.Errorf("loading secrets.yaml: %w", err)
	}

	result := config.Validate(cfg, sec)

	if outputFormat == string(output.JSON) {
		j := output.ValidationJSON{
			Valid: result.IsValid(),
		}
		for _, e := range result.Errors {
			j.Errors = append(j.Errors, fmt.Sprintf("%s: %s", e.Field, e.Message))
		}
		for _, w := range result.Warnings {
			j.Warnings = append(j.Warnings, fmt.Sprintf("%s: %s", w.Field, w.Message))
		}
		output.PrintJSON(j)
		if !result.IsValid() {
			os.Exit(2)
		}
		return nil
	}

	for _, e := range result.Errors {
		output.Error(fmt.Sprintf("[%s] %s", e.Field, e.Message))
	}
	for _, w := range result.Warnings {
		output.Warning(fmt.Sprintf("[%s] %s", w.Field, w.Message))
	}
	if result.IsValid() {
		output.Success("Configuration is valid")
		return nil
	}
	output.Error(fmt.Sprintf("%d error(s) found. Fix them before deploying.", len(result.Errors)))
	os.Exit(2)
	return nil
}

// findRepoRoot walks up from the current directory to find the RADAR-Kubernetes root
// (identified by the presence of helmfile.d/).
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "helmfile.d")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find RADAR-Kubernetes repo root (no helmfile.d/ found)")
		}
		dir = parent
	}
}

package cmd

import (
	"errors"
	"os"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/spf13/cobra"
)

var outputFormat string
var kubeContext string

var rootCmd = &cobra.Command{
	Use:   "radarctl",
	Short: "CLI for deploying and managing the RADAR-Kubernetes stack",
}

// ExitError signals radarctl should exit with a specific code.
type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string { return e.Message }

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "human", "Output format: human or json")
	rootCmd.PersistentFlags().StringVar(&kubeContext, "context", "", "Kubernetes context (overrides kubeContext in config)")
}

func isJSON() bool {
	return outputFormat == string(output.JSON)
}

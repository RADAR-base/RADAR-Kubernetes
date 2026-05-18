package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var outputFormat string
var kubeContext string

var rootCmd = &cobra.Command{
	Use:   "radarctl",
	Short: "CLI for deploying and managing the RADAR-Kubernetes stack",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "human", "Output format: human or json")
	rootCmd.PersistentFlags().StringVar(&kubeContext, "context", "", "Kubernetes context (overrides kubeContext in config)")
}

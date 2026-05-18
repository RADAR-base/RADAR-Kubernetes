package cmd

import (
	"fmt"
	"time"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/status"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	statusWatch     bool
	statusComponent string
	statusShowURLs  bool
)

var knownReleases = []string{
	"cert-manager", "kube-prometheus-stack", "nginx-ingress",
	"zookeeper", "kafka", "schema-registry", "ksql-server",
	"mongodb", "postgresql", "redis", "minio",
	"radar-appserver", "management-portal", "radar-fitbit-connector", "radar-s3-connector",
	"kratos", "hydra",
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show health status of all deployed RADAR releases",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&statusWatch, "watch", "w", false, "Refresh every 10 seconds")
	statusCmd.Flags().StringVar(&statusComponent, "component", "", "Filter to a specific component group (e.g. kafka)")
	statusCmd.Flags().BoolVar(&statusShowURLs, "show-urls", false, "Print ingress URLs for each service")
}

func runStatus(_ *cobra.Command, _ []string) error {
	exec := &executor.ShellExecutor{}
	kr := kubectl.NewRunner(exec, kubeContext)
	namespace := "default"

	for {
		report, err := status.Collect(kr, namespace, knownReleases)
		if err != nil {
			return err
		}

		if outputFormat == string(output.JSON) {
			output.PrintJSON(report)
		} else {
			printStatusReport(report)
		}

		if !statusWatch {
			break
		}
		time.Sleep(10 * time.Second)
		pterm.Println()
	}
	return nil
}

func printStatusReport(report *status.Report) {
	pterm.Printf("RADAR Stack Status\n")
	pterm.Printf("Last updated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	for _, group := range report.Groups {
		if statusComponent != "" && group.Name != statusComponent {
			continue
		}
		pterm.Bold.Println(group.Name)
		for _, r := range group.Releases {
			icon := iconFor(r.Health)
			podStr := fmt.Sprintf("%d/%d pods", r.ReadyPods, r.TotalPods)
			msg := r.Message
			if r.Error != "" {
				msg = r.Error
			}
			pterm.Printf("  %s %-30s %-10s %s", icon, r.Name, podStr, msg)
			if r.Logs != "" {
				logPreview := r.Logs
				if len(logPreview) > 80 {
					logPreview = logPreview[:80]
				}
				pterm.Printf("\n    └─ %s", logPreview)
			}
			pterm.Println()
		}
		pterm.Println()
	}

	pterm.Printf("Summary: %d healthy  %d degraded  %d warning\n",
		report.Healthy, report.Degraded, report.Warning)
}

func iconFor(h status.Health) string {
	switch h {
	case status.Healthy:
		return "✓"
	case status.Degraded:
		return "✗"
	case status.Warning:
		return "⚠"
	default:
		return "○"
	}
}

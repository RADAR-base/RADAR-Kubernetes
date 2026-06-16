package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/prereqs"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/wizard"
	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		if !r.OK && (len(r.FixCmds) > 0 || r.FixFunc != nil) {
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
		if r.OK || (len(r.FixCmds) == 0 && r.FixFunc == nil) {
			continue
		}
		action := "Installing"
		switch {
		case r.NeedsDowngrade:
			action = "Replacing"
		case r.NeedsUpgrade:
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

// runBinInit invokes the upstream bin/init script which creates the canonical
// scaffolding (environments.yaml, etc/production.yaml, etc/production.yaml.gotmpl,
// etc/secrets.yaml with random passwords, keystores). We connect stdio to the
// terminal so its output is visible, and pre-set DNAME so keystore-init doesn't
// prompt for certificate fields. bin/init is idempotent: it skips files that
// already exist, so re-running radarctl init won't clobber a configured setup.
func runBinInit(repoRoot, serverName string) error {
	scriptPath := filepath.Join(repoRoot, "bin", "init")
	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf("bin/init not found in %s — is this a RADAR-Kubernetes checkout?", repoRoot)
	}

	// bin/init refuses to proceed when etc/production.yaml exists but is older than
	// etc/base.yaml (prints a manual-merge warning and exits 1). Touching the file
	// updates its mtime so the check passes. bin/init is idempotent — it skips files
	// that already exist — so the content is never overwritten here.
	prodYAML := filepath.Join(repoRoot, "etc", "production.yaml")
	if _, err := os.Stat(prodYAML); err == nil {
		_ = exec.Command("touch", prodYAML).Run()
	}

	// If secrets.yaml exists but is missing the management_portal section it is
	// incomplete (written by a previous partial wizard run before bin/init could seed
	// it from base-secrets.yaml). Remove it so generate-secrets recreates it fully.
	secPath := filepath.Join(repoRoot, "etc", "secrets.yaml")
	if isSecretsIncomplete(secPath) {
		_ = os.Remove(secPath)
	}

	dname := "CN=" + serverName
	if serverName == "" {
		dname = "CN=radar-base"
	}
	env := append(os.Environ(), "DNAME="+dname)
	// bin/util.sh uses `sed -i` with GNU sed syntax, which fails on macOS BSD sed.
	// If gnu-sed is installed (brew install gnu-sed), shadow the system sed with it.
	if gnusedPath := gnuSedBinDir(); gnusedPath != "" {
		env = prependPATH(env, gnusedPath)
	}
	cmd := exec.Command(scriptPath)
	cmd.Dir = repoRoot
	cmd.Env = env
	// bin/generate-secrets (called by bin/init) prompts to reset existing secrets.
	// Auto-answer "n" so a re-run never rotates production passwords.
	cmd.Stdin = strings.NewReader("n\n")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runInit(_ *cobra.Command, _ []string) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	if !skipPrereqs {
		// Include managed bin dir in PATH so AutoFix commands (e.g. helm plugin install)
		// use the managed helm v3 rather than any system-installed helm v4.
		exec := &executor.ShellExecutor{ExtraEnv: prereqs.ManagedBinEnv()}
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
			if wizard.IsCancelled(err) {
				output.Info("Setup cancelled. Re-run radarctl init to resume.")
				return nil
			}
			return fmt.Errorf("wizard failed: %w", err)
		}
	}

	// Save state after wizard completes successfully.
	if err := wizard.SaveState(answers, statePath); err != nil {
		output.Warning(fmt.Sprintf("Could not save wizard state: %s", err))
	}

	// Run the upstream bin/init to set up the canonical scaffolding (env files,
	// production.yaml seeded from base.yaml, secrets.yaml with random passwords,
	// keystores). It's idempotent — files that already exist are left alone.
	output.Info("Running bin/init to scaffold config files...")
	if err := runBinInit(repoRoot, answers.ServerName); err != nil {
		return fmt.Errorf("bin/init failed: %w", err)
	}

	// Apply wizard answers as in-place overlays on top of what bin/init produced.
	w := wizard.NewWriter(repoRoot)
	if err := w.ApplyProductionOverlay(answers); err != nil {
		return fmt.Errorf("applying wizard answers to production.yaml: %w", err)
	}
	if err := w.MergeWizardSecrets(answers); err != nil {
		return fmt.Errorf("merging wizard secrets: %w", err)
	}
	output.Success("Applied wizard answers to etc/production.yaml and etc/secrets.yaml")

	// Remove state file after successful write.
	_ = os.Remove(statePath)

	output.Info("Run `radarctl deploy` to deploy the stack")
	return nil
}

// isSecretsIncomplete returns true when the file exists but is missing the
// management_portal key, which is seeded from base-secrets.yaml by generate-secrets.
// A missing key means the file was written by a partial wizard run and must be
// regenerated so helmfile templates can resolve management_portal.oauth_clients.*.
func isSecretsIncomplete(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false // doesn't exist — not our concern here
	}
	var root map[string]any
	if err := yaml.Unmarshal(data, &root); err != nil {
		return true // unparseable — regenerate to be safe
	}
	_, hasMP := root["management_portal"]
	return !hasMP
}

// gnuSedBinDir returns the directory that contains the GNU sed binary on macOS
// (installed via `brew install gnu-sed`), or "" if not found.
func gnuSedBinDir() string {
	// Homebrew on Apple Silicon installs to /opt/homebrew; Intel Macs use /usr/local.
	for _, prefix := range []string{"/opt/homebrew", "/usr/local"} {
		dir := prefix + "/opt/gnu-sed/libexec/gnubin"
		if info, err := os.Stat(dir + "/sed"); err == nil && !info.IsDir() {
			return dir
		}
	}
	return ""
}

// prependPATH returns a copy of env with dir inserted at the front of PATH.
func prependPATH(env []string, dir string) []string {
	result := make([]string, 0, len(env))
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			e = "PATH=" + dir + ":" + e[5:]
		}
		result = append(result, e)
	}
	return result
}

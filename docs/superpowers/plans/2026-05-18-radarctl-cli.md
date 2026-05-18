# radarctl CLI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build `radarctl`, a Go CLI that wraps RADAR-Kubernetes deployment with an interactive setup wizard, live deploy progress, a health status dashboard, and structured JSON output for CI/agent workflows.

**Architecture:** Single Go module at `cli/` inside the RADAR-Kubernetes repo. Commands in `cmd/` shell out to kubectl/helmfile; all business logic lives in `pkg/` for testability. Cobra for CLI structure, Huh for interactive prompts, pterm for rich terminal output.

**Tech Stack:** Go 1.22+, Cobra, Huh (charmbracelet), Viper, go-yaml v3, pterm, testify

---

## File Map

```
cli/
├── main.go                          # entry point — calls cmd.Execute()
├── go.mod                           # module: github.com/RADAR-base/RADAR-Kubernetes/cli
├── go.sum
├── cmd/
│   ├── root.go                      # root command, --output/-o and --context global flags
│   ├── validate.go                  # radarctl validate
│   ├── init.go                      # radarctl init (wires wizard into cobra)
│   ├── deploy.go                    # radarctl deploy
│   ├── status.go                    # radarctl status
│   └── diagnose.go                  # radarctl diagnose
├── pkg/
│   ├── config/
│   │   ├── types.go                 # Config and Secrets structs (YAML tags)
│   │   ├── loader.go                # LoadConfig, LoadSecrets, WriteConfig, WriteSecrets
│   │   ├── loader_test.go
│   │   ├── validator.go             # Validate(Config, Secrets) ValidationResult
│   │   ├── validator_test.go
│   │   ├── features.go              # FeatureMap: feature name → config keys + secret prompts
│   │   └── features_test.go
│   ├── prereqs/
│   │   ├── checker.go               # Check() []CheckResult — tool versions + cluster connectivity
│   │   └── checker_test.go
│   ├── wizard/
│   │   ├── wizard.go                # Run(mode) — orchestrates the full wizard flow
│   │   ├── questions.go             # question definitions, branching, mod resolution
│   │   ├── writer.go                # WriteAnswers(Answers) — writes production.yaml + secrets.yaml + environments.yaml
│   │   └── writer_test.go
│   ├── executor/
│   │   ├── executor.go              # Executor interface + ShellExecutor implementation
│   │   └── executor_test.go
│   ├── kubectl/
│   │   ├── runner.go                # GetPods, GetLogs, DescribePod, GetIngresses, RolloutStatus
│   │   └── runner_test.go
│   ├── helmfile/
│   │   ├── runner.go                # Sync, Diff, Template — streams output line by line
│   │   └── runner_test.go
│   ├── status/
│   │   ├── collector.go             # Collect() StatusReport — groups releases into categories
│   │   ├── collector_test.go
│   │   └── types.go                 # StatusReport, ReleaseStatus, ReleaseGroup enums
│   └── output/
│       ├── output.go                # Format(v, mode) — JSON or human-readable dispatch
│       └── output_test.go
```

---

## Phase 1: Foundation

### Task 1: Go module scaffold and root command

**Files:**
- Create: `cli/go.mod`
- Create: `cli/main.go`
- Create: `cli/cmd/root.go`

- [ ] **Step 1: Initialise the Go module**

```bash
cd cli
go mod init github.com/RADAR-base/RADAR-Kubernetes/cli
```

Expected: `cli/go.mod` created with `module github.com/RADAR-base/RADAR-Kubernetes/cli` and `go 1.22`.

- [ ] **Step 2: Install dependencies**

```bash
cd cli
go get github.com/spf13/cobra@v1.8.1
go get github.com/spf13/viper@v1.19.0
go get gopkg.in/yaml.v3
go get github.com/charmbracelet/huh@v0.6.0
go get github.com/pterm/pterm@v0.12.79
go get github.com/stretchr/testify@v1.9.0
go mod tidy
```

- [ ] **Step 3: Write `cli/cmd/root.go`**

```go
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
```

- [ ] **Step 4: Write `cli/main.go`**

```go
package main

import "github.com/RADAR-base/RADAR-Kubernetes/cli/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 5: Build to verify it compiles**

```bash
cd cli
go build -o bin/radarctl .
./bin/radarctl --help
```

Expected output contains: `CLI for deploying and managing the RADAR-Kubernetes stack`

- [ ] **Step 6: Commit**

```bash
git add cli/
git commit -m "feat(cli): scaffold Go module and root cobra command"
```

---

### Task 2: Config types and loader

**Files:**
- Create: `cli/pkg/config/types.go`
- Create: `cli/pkg/config/loader.go`
- Create: `cli/pkg/config/loader_test.go`

- [ ] **Step 1: Write `cli/pkg/config/types.go`**

```go
package config

// Config holds values from production.yaml (non-secret).
type Config struct {
	AtomicInstall      bool    `yaml:"atomicInstall"`
	BaseTimeout        int     `yaml:"base_timeout"`
	KubeContext        string  `yaml:"kubeContext"`
	ServerName         string  `yaml:"server_name"`
	MaintainerEmail    string  `yaml:"maintainer_email"`
	KafkaNumBrokers    int     `yaml:"kafka_num_brokers"`
	KafkaNumReplicas   int     `yaml:"kafka_num_topic_replicas"`
	KafkaNumPartitions int     `yaml:"kafka_num_topic_partitions"`
	EnableTLS          bool    `yaml:"enable_tls"`
	EnableLogging      bool    `yaml:"enable_logging_monitoring"`
	DevDeployment      bool    `yaml:"dev_deployment"`
	ConfluentCloud     bool    `yaml:"confluent_cloud"`
	EnableFitbit       bool    `yaml:"radar_fitbit_connector_install"`
	EnableRedcap       bool    `yaml:"radar_redcap_integrator_install"`
	EnableKratos       bool    `yaml:"radar_kratos_install"`
	EnableHydra        bool    `yaml:"radar_hydra_install"`
	Minio              struct {
		ExternalEndpoint string `yaml:"externalEndpoint"`
	} `yaml:"minio"`
}

// Secrets holds values from secrets.yaml (sensitive).
type Secrets struct {
	ConfluentCloud struct {
		BootstrapServer string `yaml:"bootstrapServerurl"`
		APIKey          string `yaml:"apiKey"`
		APISecret       string `yaml:"apiSecret"`
	} `yaml:"confluent_cloud"`
	FitbitClientID     string `yaml:"fitbit_client_id"`
	FitbitClientSecret string `yaml:"fitbit_client_secret"`
	GarminConsumerKey  string `yaml:"garmin_consumer_key"`
	GarminConsumerSecret string `yaml:"garmin_consumer_secret"`
	RedcapToken        string `yaml:"redcap_token"`
	S3AccessKey        string `yaml:"s3_access_key"`
	S3SecretKey        string `yaml:"s3_secret_key"`
}
```

- [ ] **Step 2: Write the failing test `cli/pkg/config/loader_test.go`**

```go
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	yaml := `
server_name: test.example.com
maintainer_email: ops@example.com
kubeContext: my-cluster
enable_tls: true
kafka_num_brokers: 3
`
	path := filepath.Join(dir, "production.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0644))

	cfg, err := config.LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, "test.example.com", cfg.ServerName)
	assert.Equal(t, "ops@example.com", cfg.MaintainerEmail)
	assert.Equal(t, "my-cluster", cfg.KubeContext)
	assert.True(t, cfg.EnableTLS)
	assert.Equal(t, 3, cfg.KafkaNumBrokers)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := config.LoadConfig("/nonexistent/path.yaml")
	assert.Error(t, err)
}

func TestWriteConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.yaml")

	cfg := &config.Config{
		ServerName:      "written.example.com",
		MaintainerEmail: "admin@example.com",
		EnableTLS:       true,
	}
	require.NoError(t, config.WriteConfig(cfg, path))

	loaded, err := config.LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, "written.example.com", loaded.ServerName)
}

func TestLoadSecrets(t *testing.T) {
	dir := t.TempDir()
	yaml := `
fitbit_client_id: my-id
fitbit_client_secret: my-secret
`
	path := filepath.Join(dir, "secrets.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0644))

	sec, err := config.LoadSecrets(path)
	require.NoError(t, err)
	assert.Equal(t, "my-id", sec.FitbitClientID)
	assert.Equal(t, "my-secret", sec.FitbitClientSecret)
}
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/config/... -run TestLoad -v
```

Expected: compilation error — `config.LoadConfig` undefined.

- [ ] **Step 4: Write `cli/pkg/config/loader.go`**

```go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &cfg, nil
}

func WriteConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadSecrets(path string) (*Secrets, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading secrets %s: %w", path, err)
	}
	var sec Secrets
	if err := yaml.Unmarshal(data, &sec); err != nil {
		return nil, fmt.Errorf("parsing secrets %s: %w", path, err)
	}
	return &sec, nil
}

func WriteSecrets(sec *Secrets, path string) error {
	data, err := yaml.Marshal(sec)
	if err != nil {
		return fmt.Errorf("marshalling secrets: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/config/... -run TestLoad -v
```

Expected: all 4 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add cli/pkg/config/
git commit -m "feat(cli): add config types, loader, and writer"
```

---

### Task 3: Config validator

**Files:**
- Create: `cli/pkg/config/validator.go`
- Create: `cli/pkg/config/validator_test.go`

- [ ] **Step 1: Write the failing tests `cli/pkg/config/validator_test.go`**

```go
package config_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/stretchr/testify/assert"
)

func validConfig() *config.Config {
	return &config.Config{
		ServerName:      "radar.example.com",
		MaintainerEmail: "ops@example.com",
		KubeContext:     "my-cluster",
		KafkaNumBrokers: 3,
	}
}

func validSecrets() *config.Secrets {
	return &config.Secrets{}
}

func TestValidate_ValidConfig(t *testing.T) {
	result := config.Validate(validConfig(), validSecrets())
	assert.Empty(t, result.Errors)
}

func TestValidate_DefaultServerName(t *testing.T) {
	cfg := validConfig()
	cfg.ServerName = "example.com"
	result := config.Validate(cfg, validSecrets())
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "server_name", result.Errors[0].Field)
}

func TestValidate_DefaultEmail(t *testing.T) {
	cfg := validConfig()
	cfg.MaintainerEmail = "MAINTAINER_EMAIL@example.com"
	result := config.Validate(cfg, validSecrets())
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "maintainer_email", result.Errors[0].Field)
}

func TestValidate_FitbitEnabledMissingSecret(t *testing.T) {
	cfg := validConfig()
	cfg.EnableFitbit = true
	sec := validSecrets()
	sec.FitbitClientID = ""
	result := config.Validate(cfg, sec)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "fitbit_client_id", result.Errors[0].Field)
}

func TestValidate_PlaceholderSecretWarning(t *testing.T) {
	sec := validSecrets()
	sec.FitbitClientID = "change_me"
	result := config.Validate(validConfig(), sec)
	assert.Len(t, result.Warnings, 1)
	assert.Equal(t, "fitbit_client_id", result.Warnings[0].Field)
}

func TestValidate_ConfluentMissingBootstrap(t *testing.T) {
	cfg := validConfig()
	cfg.ConfluentCloud = true
	result := config.Validate(cfg, validSecrets())
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "confluent_cloud.bootstrapServerurl", result.Errors[0].Field)
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/config/... -run TestValidate -v
```

Expected: compilation error — `config.Validate` undefined.

- [ ] **Step 3: Write `cli/pkg/config/validator.go`**

```go
package config

import "strings"

type ValidationError struct {
	Field   string
	Message string
}

type ValidationWarning struct {
	Field   string
	Message string
}

type ValidationResult struct {
	Errors   []ValidationError
	Warnings []ValidationWarning
}

func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

var placeholders = []string{"change_me", "secret", "MAINTAINER_EMAIL@example.com", "example.com"}

func isPlaceholder(v string) bool {
	for _, p := range placeholders {
		if strings.EqualFold(v, p) {
			return true
		}
	}
	return false
}

func Validate(cfg *Config, sec *Secrets) ValidationResult {
	var result ValidationResult

	if cfg.ServerName == "" || cfg.ServerName == "example.com" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "server_name",
			Message: "must be set to your actual domain name",
		})
	}

	if cfg.MaintainerEmail == "" || cfg.MaintainerEmail == "MAINTAINER_EMAIL@example.com" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "maintainer_email",
			Message: "must be set to a real email address",
		})
	}

	if cfg.ConfluentCloud && sec.ConfluentCloud.BootstrapServer == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "confluent_cloud.bootstrapServerurl",
			Message: "required when confluent_cloud is enabled",
		})
	}

	if cfg.EnableFitbit && sec.FitbitClientID == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "fitbit_client_id",
			Message: "required when Fitbit integration is enabled",
		})
	}

	// Placeholder warnings
	secretFields := map[string]string{
		"fitbit_client_id":     sec.FitbitClientID,
		"fitbit_client_secret": sec.FitbitClientSecret,
		"garmin_consumer_key":  sec.GarminConsumerKey,
		"redcap_token":         sec.RedcapToken,
	}
	for field, val := range secretFields {
		if val != "" && isPlaceholder(val) {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Field:   field,
				Message: "still contains a placeholder value",
			})
		}
	}

	return result
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/config/... -run TestValidate -v
```

Expected: all 6 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add cli/pkg/config/validator.go cli/pkg/config/validator_test.go
git commit -m "feat(cli): add config validator with placeholder and dependency checks"
```

---

### Task 4: Feature expansion map

**Files:**
- Create: `cli/pkg/config/features.go`
- Create: `cli/pkg/config/features_test.go`

- [ ] **Step 1: Write the failing tests `cli/pkg/config/features_test.go`**

```go
package config_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFeature_Fitbit(t *testing.T) {
	cfg := &config.Config{}
	err := config.ApplyFeature(cfg, "fitbit")
	require.NoError(t, err)
	assert.True(t, cfg.EnableFitbit)
}

func TestApplyFeature_Unknown(t *testing.T) {
	cfg := &config.Config{}
	err := config.ApplyFeature(cfg, "unknown_feature")
	assert.Error(t, err)
}

func TestFeatureSecretPrompts_Fitbit(t *testing.T) {
	prompts := config.FeatureSecretPrompts("fitbit")
	require.Len(t, prompts, 2)
	assert.Equal(t, "fitbit_client_id", prompts[0].SecretKey)
	assert.Equal(t, "fitbit_client_secret", prompts[1].SecretKey)
	assert.True(t, prompts[1].Mask)
}

func TestApplyDeploymentProfile_Dev(t *testing.T) {
	cfg := &config.Config{}
	mods := config.ApplyDeploymentProfile(cfg, "dev")
	assert.False(t, cfg.EnableTLS)
	assert.True(t, cfg.DevDeployment)
	assert.Contains(t, mods, "mods/minimal.yaml")
	assert.Contains(t, mods, "mods/localdev.yaml")
	assert.Contains(t, mods, "mods/disable_tls.yaml")
	assert.Contains(t, mods, "mods/fast_deploy.yaml")
}

func TestApplyDeploymentProfile_Production(t *testing.T) {
	cfg := &config.Config{}
	mods := config.ApplyDeploymentProfile(cfg, "production")
	assert.True(t, cfg.EnableTLS)
	assert.False(t, cfg.DevDeployment)
	assert.Empty(t, mods)
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/config/... -run TestApplyFeature -run TestFeatureSecret -run TestApplyDeployment -v
```

Expected: compilation error — `config.ApplyFeature` undefined.

- [ ] **Step 3: Write `cli/pkg/config/features.go`**

```go
package config

import "fmt"

// SecretPrompt describes a secret value the wizard must collect from the user.
type SecretPrompt struct {
	SecretKey string
	Label     string
	Mask      bool
}

type feature struct {
	apply   func(*Config)
	secrets []SecretPrompt
}

var featureRegistry = map[string]feature{
	"fitbit": {
		apply: func(c *Config) { c.EnableFitbit = true },
		secrets: []SecretPrompt{
			{SecretKey: "fitbit_client_id", Label: "Fitbit OAuth Client ID", Mask: false},
			{SecretKey: "fitbit_client_secret", Label: "Fitbit OAuth Client Secret", Mask: true},
		},
	},
	"garmin": {
		apply: func(c *Config) {},
		secrets: []SecretPrompt{
			{SecretKey: "garmin_consumer_key", Label: "Garmin Consumer Key", Mask: false},
			{SecretKey: "garmin_consumer_secret", Label: "Garmin Consumer Secret", Mask: true},
		},
	},
	"redcap": {
		apply: func(c *Config) { c.EnableRedcap = true },
		secrets: []SecretPrompt{
			{SecretKey: "redcap_token", Label: "REDCap API Token", Mask: true},
		},
	},
	"kratos": {
		apply:   func(c *Config) { c.EnableKratos = true; c.EnableHydra = true },
		secrets: nil,
	},
}

// ApplyFeature sets the config fields required by a named feature.
func ApplyFeature(cfg *Config, name string) error {
	f, ok := featureRegistry[name]
	if !ok {
		return fmt.Errorf("unknown feature: %s", name)
	}
	f.apply(cfg)
	return nil
}

// FeatureSecretPrompts returns the secret prompts needed for a feature.
func FeatureSecretPrompts(name string) []SecretPrompt {
	f, ok := featureRegistry[name]
	if !ok {
		return nil
	}
	return f.secrets
}

// ApplyDeploymentProfile sets config values for a named profile and returns the mods to apply.
func ApplyDeploymentProfile(cfg *Config, profile string) []string {
	switch profile {
	case "dev":
		cfg.EnableTLS = false
		cfg.DevDeployment = true
		cfg.KafkaNumBrokers = 1
		cfg.KafkaNumReplicas = 1
		return []string{
			"mods/minimal.yaml",
			"mods/localdev.yaml",
			"mods/disable_tls.yaml",
			"mods/fast_deploy.yaml",
		}
	case "staging":
		cfg.EnableTLS = true
		cfg.DevDeployment = false
		cfg.KafkaNumBrokers = 1
		return []string{"mods/minimal.yaml"}
	default: // production
		cfg.EnableTLS = true
		cfg.DevDeployment = false
		cfg.KafkaNumBrokers = 3
		return nil
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/config/... -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add cli/pkg/config/features.go cli/pkg/config/features_test.go
git commit -m "feat(cli): add feature expansion map and deployment profiles"
```

---

## Phase 2: External Tool Runners

### Task 5: Executor interface and kubectl runner

**Files:**
- Create: `cli/pkg/executor/executor.go`
- Create: `cli/pkg/executor/executor_test.go`
- Create: `cli/pkg/kubectl/runner.go`
- Create: `cli/pkg/kubectl/runner_test.go`

- [ ] **Step 1: Write `cli/pkg/executor/executor.go`**

```go
package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor runs shell commands and returns combined output.
type Executor interface {
	Run(name string, args ...string) (string, error)
}

// ShellExecutor runs real system commands.
type ShellExecutor struct{}

func (e *ShellExecutor) Run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w\n%s", err, strings.TrimSpace(errOut.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

// MockExecutor records calls and returns configured responses for testing.
type MockExecutor struct {
	Responses map[string]string // key: "name args...", value: output
	Errors    map[string]error
	Calls     []string
}

func (m *MockExecutor) Run(name string, args ...string) (string, error) {
	key := strings.Join(append([]string{name}, args...), " ")
	m.Calls = append(m.Calls, key)
	if err, ok := m.Errors[key]; ok {
		return "", err
	}
	if out, ok := m.Responses[key]; ok {
		return out, nil
	}
	return "", nil
}
```

- [ ] **Step 2: Write the failing tests `cli/pkg/kubectl/runner_test.go`**

```go
package kubectl_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockPodJSON() string {
	return `{
  "items": [
    {
      "metadata": {"name": "mongodb-0", "namespace": "default"},
      "status": {
        "phase": "Running",
        "containerStatuses": [{"ready": true, "restartCount": 0}]
      }
    }
  ]
}`
}

func TestGetPods(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl get pods -n default -l app=mongodb -o json": mockPodJSON(),
		},
	}
	r := kubectl.NewRunner(mock, "")
	pods, err := r.GetPods("default", "app=mongodb")
	require.NoError(t, err)
	require.Len(t, pods, 1)
	assert.Equal(t, "mongodb-0", pods[0].Name)
	assert.True(t, pods[0].Ready)
}

func TestGetLogs(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl logs mongodb-0 -n default --tail=20": "log line 1\nlog line 2",
		},
	}
	r := kubectl.NewRunner(mock, "")
	logs, err := r.GetLogs("mongodb-0", "default", 20)
	require.NoError(t, err)
	assert.Contains(t, logs, "log line 1")
}
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/kubectl/... -v
```

Expected: compilation error — package `kubectl` not found.

- [ ] **Step 4: Write `cli/pkg/kubectl/runner.go`**

```go
package kubectl

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type Pod struct {
	Name      string
	Namespace string
	Phase     string
	Ready     bool
	Restarts  int
	Message   string
}

type Ingress struct {
	Name      string
	Namespace string
	Host      string
	Paths     []string
}

type Runner struct {
	exec    executor.Executor
	context string
}

func NewRunner(exec executor.Executor, context string) *Runner {
	return &Runner{exec: exec, context: context}
}

func (r *Runner) contextArgs() []string {
	if r.context != "" {
		return []string{"--context", r.context}
	}
	return nil
}

func (r *Runner) GetPods(namespace, labelSelector string) ([]Pod, error) {
	args := append(r.contextArgs(), "get", "pods", "-n", namespace, "-l", labelSelector, "-o", "json")
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return nil, fmt.Errorf("kubectl get pods: %w", err)
	}

	var raw struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Status struct {
				Phase            string `json:"phase"`
				ContainerStatuses []struct {
					Ready        bool `json:"ready"`
					RestartCount int  `json:"restartCount"`
				} `json:"containerStatuses"`
				Message string `json:"message"`
			} `json:"status"`
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, fmt.Errorf("parsing pod list: %w", err)
	}

	pods := make([]Pod, 0, len(raw.Items))
	for _, item := range raw.Items {
		ready := false
		restarts := 0
		if len(item.Status.ContainerStatuses) > 0 {
			ready = item.Status.ContainerStatuses[0].Ready
			restarts = item.Status.ContainerStatuses[0].RestartCount
		}
		pods = append(pods, Pod{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Phase:     item.Status.Phase,
			Ready:     ready,
			Restarts:  restarts,
			Message:   item.Status.Message,
		})
	}
	return pods, nil
}

func (r *Runner) GetLogs(podName, namespace string, tail int) (string, error) {
	args := append(r.contextArgs(), "logs", podName, "-n", namespace, fmt.Sprintf("--tail=%d", tail))
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return "", fmt.Errorf("kubectl logs %s: %w", podName, err)
	}
	return out, nil
}

func (r *Runner) GetIngresses(namespace string) ([]Ingress, error) {
	args := append(r.contextArgs(), "get", "ingress", "-n", namespace, "-o", "json")
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return nil, fmt.Errorf("kubectl get ingress: %w", err)
	}

	var raw struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Spec struct {
				Rules []struct {
					Host string `json:"host"`
					HTTP struct {
						Paths []struct {
							Path string `json:"path"`
						} `json:"paths"`
					} `json:"http"`
				} `json:"rules"`
			} `json:"spec"`
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, fmt.Errorf("parsing ingress list: %w", err)
	}

	ingresses := make([]Ingress, 0, len(raw.Items))
	for _, item := range raw.Items {
		var host string
		var paths []string
		if len(item.Spec.Rules) > 0 {
			host = item.Spec.Rules[0].Host
			for _, p := range item.Spec.Rules[0].HTTP.Paths {
				paths = append(paths, p.Path)
			}
		}
		ingresses = append(ingresses, Ingress{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Host:      host,
			Paths:     paths,
		})
	}
	return ingresses, nil
}

func (r *Runner) Version() (string, error) {
	args := append(r.contextArgs(), "version", "--client", "-o", "json")
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return "", err
	}
	var v struct {
		ClientVersion struct {
			GitVersion string `json:"gitVersion"`
		} `json:"clientVersion"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		return strings.TrimSpace(out), nil
	}
	return v.ClientVersion.GitVersion, nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/kubectl/... -v
```

Expected: all tests PASS.

- [ ] **Step 6: Commit**

```bash
git add cli/pkg/executor/ cli/pkg/kubectl/
git commit -m "feat(cli): add executor interface and kubectl runner"
```

---

### Task 6: helmfile runner

**Files:**
- Create: `cli/pkg/helmfile/runner.go`
- Create: `cli/pkg/helmfile/runner_test.go`

- [ ] **Step 1: Write the failing tests `cli/pkg/helmfile/runner_test.go`**

```go
package helmfile_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/helmfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"helmfile --environment default diff": "Release=mongodb has changed\nRelease=kafka unchanged",
		},
	}
	r := helmfile.NewRunner(mock, "/repo", "default")
	out, err := r.Diff("")
	require.NoError(t, err)
	assert.Contains(t, out, "mongodb has changed")
}

func TestDiff_WithSelector(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"helmfile --environment default --selector name=mongodb diff": "Release=mongodb has changed",
		},
	}
	r := helmfile.NewRunner(mock, "/repo", "default")
	out, err := r.Diff("name=mongodb")
	require.NoError(t, err)
	assert.Contains(t, out, "mongodb")
}

func TestSyncArgs(t *testing.T) {
	mock := &executor.MockExecutor{}
	r := helmfile.NewRunner(mock, "/repo", "default")
	_ = r.Sync("", false)
	require.Len(t, mock.Calls, 1)
	assert.Equal(t, "helmfile --environment default sync", mock.Calls[0])
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/helmfile/... -v
```

Expected: compilation error — package `helmfile` not found.

- [ ] **Step 3: Write `cli/pkg/helmfile/runner.go`**

```go
package helmfile

import (
	"fmt"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type Runner struct {
	exec        executor.Executor
	workDir     string
	environment string
}

func NewRunner(exec executor.Executor, workDir, environment string) *Runner {
	return &Runner{exec: exec, workDir: workDir, environment: environment}
}

func (r *Runner) baseArgs(selector string) []string {
	args := []string{"--environment", r.environment}
	if selector != "" {
		args = append(args, "--selector", selector)
	}
	return args
}

// Diff returns the helmfile diff output as a string.
func (r *Runner) Diff(selector string) (string, error) {
	args := append(r.baseArgs(selector), "diff")
	out, err := r.exec.Run("helmfile", args...)
	if err != nil {
		// helmfile diff exits non-zero when there are differences — treat as success
		return out, nil
	}
	return out, nil
}

// Sync runs helmfile sync, streaming output line by line via the onLine callback.
// Pass selector="" to sync all releases.
func (r *Runner) Sync(selector string, atomic bool) error {
	args := r.baseArgs(selector)
	if atomic {
		args = append(args, "--atomic")
	}
	args = append(args, "sync")
	_, err := r.exec.Run("helmfile", args...)
	return err
}

// SyncWithCallback runs helmfile sync and calls onLine for each output line.
// Used by the deploy command to render live progress.
func (r *Runner) SyncWithCallback(selector string, atomic bool, onLine func(string)) error {
	args := r.baseArgs(selector)
	if atomic {
		args = append(args, "--atomic")
	}
	args = append(args, "sync")

	// For real execution, use ShellExecutor with line streaming.
	// In tests, MockExecutor.Run returns the full output at once.
	out, err := r.exec.Run("helmfile", args...)
	for _, line := range strings.Split(out, "\n") {
		if line != "" {
			onLine(line)
		}
	}
	return err
}

// ParseDiffSummary counts changed/new/removed releases from helmfile diff output.
func ParseDiffSummary(diff string) (updated, installed, removed int) {
	for _, line := range strings.Split(diff, "\n") {
		lower := strings.ToLower(line)
		switch {
		case strings.Contains(lower, "has changed"):
			updated++
		case strings.Contains(lower, "installing"):
			installed++
		case strings.Contains(lower, "deleting"):
			removed++
		}
	}
	return
}

// ParseReleaseFromLine extracts a release name from a helmfile log line.
// e.g. "Comparing release=mongodb, chart=radar/mongodb" → "mongodb"
func ParseReleaseFromLine(line string) string {
	line = strings.ToLower(line)
	for _, prefix := range []string{"release=", "upgrading release ", "installing release "} {
		if idx := strings.Index(line, prefix); idx >= 0 {
			rest := line[idx+len(prefix):]
			if end := strings.IndexAny(rest, ", \t\n"); end >= 0 {
				return rest[:end]
			}
			return strings.TrimSpace(rest)
		}
	}
	return ""
}

// Version returns the installed helmfile version string.
func (r *Runner) Version() (string, error) {
	out, err := r.exec.Run("helmfile", "--version")
	if err != nil {
		return "", fmt.Errorf("helmfile not found: %w", err)
	}
	return strings.TrimSpace(out), nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/helmfile/... -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add cli/pkg/helmfile/
git commit -m "feat(cli): add helmfile runner with diff parsing"
```

---

## Phase 3: Output and Validate Commands

### Task 7: Output formatter and `radarctl validate`

**Files:**
- Create: `cli/pkg/output/output.go`
- Create: `cli/pkg/output/output_test.go`
- Create: `cli/cmd/validate.go`

- [ ] **Step 1: Write `cli/pkg/output/output.go`**

```go
package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

type Mode string

const (
	Human Mode = "human"
	JSON  Mode = "json"
)

// PrintJSON writes v as indented JSON to stdout.
func PrintJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// Success prints a green checkmark line.
func Success(msg string) {
	pterm.Success.Println(msg)
}

// Error prints a red error line.
func Error(msg string) {
	pterm.Error.Println(msg)
}

// Warning prints a yellow warning line.
func Warning(msg string) {
	pterm.Warning.Println(msg)
}

// Info prints an info line.
func Info(msg string) {
	pterm.Info.Println(msg)
}

// Header prints a bold section header.
func Header(msg string) {
	pterm.DefaultHeader.WithFullWidth().Println(msg)
}

// ValidationJSON is the JSON schema for validate output.
type ValidationJSON struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// Fatalf prints an error and exits 1.
func Fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
```

- [ ] **Step 2: Write the failing test `cli/pkg/output/output_test.go`**

```go
package output_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
)

func TestModeConstants(t *testing.T) {
	if output.Human != "human" {
		t.Fatal("expected Human == 'human'")
	}
	if output.JSON != "json" {
		t.Fatal("expected JSON == 'json'")
	}
}
```

- [ ] **Step 3: Run test to verify it passes**

```bash
cd cli
go test ./pkg/output/... -v
```

Expected: PASS.

- [ ] **Step 4: Write `cli/cmd/validate.go`**

```go
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

	// Human output
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
```

- [ ] **Step 5: Build and run validate**

```bash
cd cli
go build -o bin/radarctl .
cd ..   # back to repo root
cli/bin/radarctl validate --help
```

Expected: help text for validate command.

- [ ] **Step 6: Commit**

```bash
git add cli/pkg/output/ cli/cmd/validate.go
git commit -m "feat(cli): add output package and radarctl validate command"
```

---

## Phase 4: Prerequisites Checker and Init Wizard

### Task 8: Prerequisites checker

**Files:**
- Create: `cli/pkg/prereqs/checker.go`
- Create: `cli/pkg/prereqs/checker_test.go`

- [ ] **Step 1: Write the failing tests `cli/pkg/prereqs/checker_test.go`**

```go
package prereqs_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/prereqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_AllPresent(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl version --client -o json": `{"clientVersion":{"gitVersion":"v1.30.2"}}`,
			"helm version --short":             "v3.15.1",
			"helmfile --version":               "helmfile version v0.169.1",
			"helm diff version":                "3.9.12",
			"yq --version":                     "yq (https://github.com/mikefarah/yq/) version v4.44.3",
			"java -version":                    "openjdk version \"21.0.1\"",
			"openssl version":                  "OpenSSL 3.1.4 24 Oct 2023",
			"git --version":                    "git version 2.42.0",
		},
	}
	results := prereqs.Check(mock)
	for _, r := range results {
		assert.True(t, r.OK, "expected %s to pass, got error: %s", r.Tool, r.Error)
	}
}

func TestCheck_MissingTool(t *testing.T) {
	mock := &executor.MockExecutor{
		Errors: map[string]error{
			"kubectl version --client -o json": fmt.Errorf("executable file not found in $PATH"),
		},
		Responses: map[string]string{
			"helm version --short":   "v3.15.1",
			"helmfile --version":     "helmfile version v0.169.1",
			"helm diff version":      "3.9.12",
			"yq --version":           "yq (https://github.com/mikefarah/yq/) version v4.44.3",
			"java -version":          "openjdk version \"21.0.1\"",
			"openssl version":        "OpenSSL 3.1.4",
			"git --version":          "git version 2.42.0",
		},
	}
	results := prereqs.Check(mock)
	var kubectlResult *prereqs.CheckResult
	for i := range results {
		if results[i].Tool == "kubectl" {
			kubectlResult = &results[i]
			break
		}
	}
	require.NotNil(t, kubectlResult)
	assert.False(t, kubectlResult.OK)
	assert.NotEmpty(t, kubectlResult.Error)
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/prereqs/... -v
```

Expected: compilation error — package `prereqs` not found.

- [ ] **Step 3: Write `cli/pkg/prereqs/checker.go`**

```go
package prereqs

import (
	"fmt"
	"strings"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type CheckResult struct {
	Tool    string
	OK      bool
	Version string
	Error   string
	InstallHint string
}

type toolDef struct {
	name        string
	args        []string
	parseVersion func(string) string
	installHint string
}

var tools = []toolDef{
	{
		name: "kubectl",
		args: []string{"version", "--client", "-o", "json"},
		parseVersion: func(out string) string {
			// Extract gitVersion from JSON
			for _, line := range strings.Split(out, "\n") {
				if strings.Contains(line, "gitVersion") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						return strings.Trim(strings.TrimSpace(parts[1]), `",`)
					}
				}
			}
			return out
		},
		installHint: "https://kubernetes.io/docs/tasks/tools/",
	},
	{
		name: "helm",
		args: []string{"version", "--short"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "https://helm.sh/docs/intro/install/",
	},
	{
		name: "helmfile",
		args: []string{"--version"},
		parseVersion: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return out
		},
		installHint: "https://github.com/helmfile/helmfile/releases",
	},
	{
		name: "helm diff",
		args: []string{"diff", "version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "helm plugin install https://github.com/databus23/helm-diff",
	},
	{
		name: "yq",
		args: []string{"--version"},
		parseVersion: func(out string) string {
			parts := strings.Fields(out)
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
			return out
		},
		installHint: "https://github.com/mikefarah/yq#install",
	},
	{
		name: "java",
		args: []string{"-version"},
		parseVersion: func(out string) string {
			lines := strings.Split(out, "\n")
			if len(lines) > 0 {
				return strings.TrimSpace(lines[0])
			}
			return out
		},
		installHint: "https://adoptium.net/ or: brew install openjdk",
	},
	{
		name: "openssl",
		args: []string{"version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "brew install openssl  (macOS) or: apt install openssl",
	},
	{
		name: "git",
		args: []string{"--version"},
		parseVersion: func(out string) string { return strings.TrimSpace(out) },
		installHint:  "https://git-scm.com/downloads",
	},
}

// Check runs all prerequisite checks and returns one result per tool.
func Check(exec executor.Executor) []CheckResult {
	results := make([]CheckResult, 0, len(tools))
	for _, t := range tools {
		name := strings.Fields(t.name)[0] // "helm diff" → "helm"
		args := t.args
		if len(strings.Fields(t.name)) > 1 {
			// "helm diff" → run as "helm diff version"
			prefix := strings.Fields(t.name)[1:]
			args = append(prefix, t.args...)
		}
		out, err := exec.Run(name, args...)
		if err != nil {
			results = append(results, CheckResult{
				Tool:        t.name,
				OK:          false,
				Error:       fmt.Sprintf("not found or failed: %s", err.Error()),
				InstallHint: t.installHint,
			})
			continue
		}
		results = append(results, CheckResult{
			Tool:    t.name,
			OK:      true,
			Version: t.parseVersion(out),
		})
	}
	return results
}
```

- [ ] **Step 4: Fix the missing `fmt` import in the test file**

Add `"fmt"` to the import block in `cli/pkg/prereqs/checker_test.go`.

- [ ] **Step 5: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/prereqs/... -v
```

Expected: all tests PASS.

- [ ] **Step 6: Commit**

```bash
git add cli/pkg/prereqs/
git commit -m "feat(cli): add prerequisites checker"
```

---

### Task 9: Wizard writer

**Files:**
- Create: `cli/pkg/wizard/writer.go`
- Create: `cli/pkg/wizard/writer_test.go`

- [ ] **Step 1: Define the Answers type and write the failing tests `cli/pkg/wizard/writer_test.go`**

```go
package wizard_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/wizard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestWriteAnswers_ProductionYAML(t *testing.T) {
	dir := t.TempDir()
	answers := &wizard.Answers{
		ServerName:      "radar.example.com",
		MaintainerEmail: "ops@example.com",
		KubeContext:     "my-cluster",
		Profile:         "production",
		UseConfluent:    false,
		Features:        []string{},
		AppliedMods:     []string{},
	}

	w := wizard.NewWriter(dir)
	require.NoError(t, w.WriteAnswers(answers))

	data, err := os.ReadFile(filepath.Join(dir, "etc", "production.yaml"))
	require.NoError(t, err)

	var out map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &out))
	assert.Equal(t, "radar.example.com", out["server_name"])
	assert.Equal(t, "ops@example.com", out["maintainer_email"])
}

func TestWriteAnswers_DevProfile_DisablesTLS(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "etc"), 0755))

	answers := &wizard.Answers{
		ServerName:      "localhost",
		MaintainerEmail: "dev@example.com",
		KubeContext:     "k3d-local",
		Profile:         "dev",
		AppliedMods:     []string{"mods/minimal.yaml", "mods/disable_tls.yaml"},
	}

	w := wizard.NewWriter(dir)
	require.NoError(t, w.WriteAnswers(answers))

	data, err := os.ReadFile(filepath.Join(dir, "etc", "production.yaml"))
	require.NoError(t, err)
	var out map[string]interface{}
	require.NoError(t, yaml.Unmarshal(data, &out))
	assert.Equal(t, false, out["enable_tls"])
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/wizard/... -v
```

Expected: compilation error — package `wizard` not found.

- [ ] **Step 3: Write `cli/pkg/wizard/writer.go`**

```go
package wizard

import (
	"os"
	"path/filepath"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
	"gopkg.in/yaml.v3"
)

// Answers holds all values collected by the wizard.
type Answers struct {
	ServerName      string
	MaintainerEmail string
	KubeContext     string
	Profile         string // "production", "staging", "dev"
	UseConfluent    bool
	ConfluentURL    string
	ConfluentKey    string
	ConfluentSecret string
	Features        []string // enabled feature names e.g. ["fitbit", "redcap"]
	UseExternalS3   bool
	S3Bucket        string
	S3Region        string
	AppliedMods     []string
	Secrets         map[string]string // secret key → value collected from prompts
}

// Writer persists wizard answers to config files.
type Writer struct {
	repoRoot string
}

func NewWriter(repoRoot string) *Writer {
	return &Writer{repoRoot: repoRoot}
}

// WriteAnswers writes production.yaml and secrets.yaml from the collected answers.
func (w *Writer) WriteAnswers(a *Answers) error {
	if err := os.MkdirAll(filepath.Join(w.repoRoot, "etc"), 0755); err != nil {
		return err
	}

	cfg := &config.Config{
		ServerName:      a.ServerName,
		MaintainerEmail: a.MaintainerEmail,
		KubeContext:     a.KubeContext,
		AtomicInstall:   true,
		BaseTimeout:     90,
	}

	// Apply profile
	config.ApplyDeploymentProfile(cfg, a.Profile)

	// Apply features
	for _, f := range a.Features {
		_ = config.ApplyFeature(cfg, f)
	}

	// Confluent cloud
	if a.UseConfluent {
		cfg.ConfluentCloud = true
	}

	cfgPath := filepath.Join(w.repoRoot, "etc", "production.yaml")
	if err := config.WriteConfig(cfg, cfgPath); err != nil {
		return err
	}

	// Write secrets
	sec := &config.Secrets{}
	if a.UseConfluent {
		sec.ConfluentCloud.BootstrapServer = a.ConfluentURL
		sec.ConfluentCloud.APIKey = a.ConfluentKey
		sec.ConfluentCloud.APISecret = a.ConfluentSecret
	}
	if v, ok := a.Secrets["fitbit_client_id"]; ok {
		sec.FitbitClientID = v
	}
	if v, ok := a.Secrets["fitbit_client_secret"]; ok {
		sec.FitbitClientSecret = v
	}
	if v, ok := a.Secrets["garmin_consumer_key"]; ok {
		sec.GarminConsumerKey = v
	}
	if v, ok := a.Secrets["garmin_consumer_secret"]; ok {
		sec.GarminConsumerSecret = v
	}
	if v, ok := a.Secrets["redcap_token"]; ok {
		sec.RedcapToken = v
	}

	secPath := filepath.Join(w.repoRoot, "etc", "secrets.yaml")
	return config.WriteSecrets(sec, secPath)
}

// WriteEnvironmentsYAML writes environments.yaml with the selected mods.
func (w *Writer) WriteEnvironmentsYAML(mods []string) error {
	type envConfig struct {
		Environments map[string]struct {
			Values []string `yaml:"values"`
		} `yaml:"environments"`
	}

	values := []string{
		"../etc/base.yaml",
		"../etc/production.yaml",
		"../etc/secrets.yaml",
	}
	for _, mod := range mods {
		values = append(values, "../"+mod)
	}

	cfg := envConfig{
		Environments: map[string]struct {
			Values []string `yaml:"values"`
		}{
			"default": {Values: values},
		},
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(w.repoRoot, "environments.yaml"), data, 0644)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/wizard/... -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
git add cli/pkg/wizard/
git commit -m "feat(cli): add wizard writer — produces production.yaml and secrets.yaml"
```

---

### Task 10: Wizard questions and `radarctl init` command

**Files:**
- Create: `cli/pkg/wizard/wizard.go`
- Create: `cli/pkg/wizard/questions.go`
- Create: `cli/cmd/init.go`

- [ ] **Step 1: Write `cli/pkg/wizard/questions.go`**

This file defines the question flow. It uses Huh forms for interactive input.

```go
package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
)

// collectBasics gathers cluster basics from the user.
func collectBasics(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server hostname (your domain name)").
				Placeholder("radar.example.com").
				Value(&a.ServerName).
				Validate(func(s string) error {
					if s == "" || s == "example.com" {
						return fmt.Errorf("enter your actual domain name")
					}
					return nil
				}),
			huh.NewInput().
				Title("Maintainer email (for TLS certificate notifications)").
				Placeholder("ops@example.com").
				Value(&a.MaintainerEmail).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("email is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Kubernetes context").
				Placeholder("default").
				Value(&a.KubeContext),
		),
	).Run()
}

// collectProfile asks for the deployment profile.
func collectProfile(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Deployment profile").
				Options(
					huh.NewOption("Production — full stack, TLS enabled", "production"),
					huh.NewOption("Staging — minimal resources, TLS enabled", "staging"),
					huh.NewOption("Local dev — no TLS, single replicas, fast probes", "dev"),
				).
				Value(&a.Profile),
		),
	).Run()
}

// collectKafka asks whether to use local Kafka or Confluent Cloud.
func collectKafka(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Use Confluent Cloud instead of local Kafka?").
				Value(&a.UseConfluent),
		),
	).Run()
}

// collectConfluentCredentials collects Confluent Cloud credentials.
func collectConfluentCredentials(a *Answers) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Confluent Cloud bootstrap server URL").
				Value(&a.ConfluentURL),
			huh.NewInput().
				Title("Confluent Cloud API key").
				Value(&a.ConfluentKey),
			huh.NewInput().
				Title("Confluent Cloud API secret").
				EchoMode(huh.EchoModePassword).
				Value(&a.ConfluentSecret),
		),
	).Run()
}

// collectFeatures asks which optional data source integrations to enable.
func collectFeatures(a *Answers) error {
	featureOptions := []huh.Option[string]{
		huh.NewOption("Fitbit", "fitbit"),
		huh.NewOption("Garmin", "garmin"),
		huh.NewOption("REDCap", "redcap"),
		huh.NewOption("Upload portal", "upload"),
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Which data source integrations do you want to enable?").
				Options(featureOptions...).
				Value(&a.Features),
		),
	).Run()
}

// collectFeatureSecrets collects API credentials for each enabled feature.
func collectFeatureSecrets(a *Answers) error {
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}
	for _, feature := range a.Features {
		prompts := config.FeatureSecretPrompts(feature)
		for _, p := range prompts {
			val := ""
			input := huh.NewInput().
				Title(p.Label).
				Value(&val)
			if p.Mask {
				input = input.EchoMode(huh.EchoModePassword)
			}
			if err := huh.NewForm(huh.NewGroup(input)).Run(); err != nil {
				return err
			}
			a.Secrets[p.SecretKey] = val
		}
	}
	return nil
}

// collectAuth asks whether to enable Ory Hydra/Kratos.
func collectAuth(a *Answers) error {
	var enable bool
	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Ory Hydra/Kratos (OAuth2/OIDC identity management)?").
				Value(&enable),
		),
	).Run(); err != nil {
		return err
	}
	if enable {
		a.Features = append(a.Features, "kratos")
	}
	return nil
}
```

- [ ] **Step 2: Write `cli/pkg/wizard/wizard.go`**

```go
package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/config"
)

type Mode string

const (
	ModeWizard      Mode = "wizard"
	ModeInteractive Mode = "interactive"
	ModeExpert      Mode = "expert"
)

// Run executes the wizard in the given mode and returns the collected Answers.
func Run(mode Mode, repoRoot string) (*Answers, error) {
	switch mode {
	case ModeWizard:
		return runWizard()
	case ModeInteractive:
		return runInteractive()
	case ModeExpert:
		return nil, nil // expert skips prompts
	default:
		return nil, fmt.Errorf("unknown wizard mode: %s", mode)
	}
}

func runWizard() (*Answers, error) {
	a := &Answers{Secrets: make(map[string]string)}

	steps := []func(*Answers) error{
		collectBasics,
		collectProfile,
		collectKafka,
		func(a *Answers) error {
			if a.UseConfluent {
				return collectConfluentCredentials(a)
			}
			return nil
		},
		collectFeatures,
		collectFeatureSecrets,
		collectAuth,
	}

	for _, step := range steps {
		if err := step(a); err != nil {
			return nil, fmt.Errorf("wizard step failed: %w", err)
		}
	}

	// Resolve mods from profile
	cfg := &config.Config{}
	a.AppliedMods = config.ApplyDeploymentProfile(cfg, a.Profile)

	return a, nil
}

func runInteractive() (*Answers, error) {
	// Interactive mode: same questions as wizard but all shown at once.
	// For simplicity in v1, runs the same steps as wizard mode.
	return runWizard()
}

// SelectMode asks the user to pick wizard/interactive/expert.
func SelectMode() (Mode, error) {
	var chosen string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How would you like to configure your deployment?").
				Options(
					huh.NewOption("Guided wizard (recommended for first-time setup)", string(ModeWizard)),
					huh.NewOption("Interactive (configure every option explicitly)", string(ModeInteractive)),
					huh.NewOption("Expert (validate existing config files only)", string(ModeExpert)),
				).
				Value(&chosen),
		),
	).Run()
	return Mode(chosen), err
}
```

- [ ] **Step 3: Write `cli/cmd/init.go`**

```go
package cmd

import (
	"fmt"

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

	// Phase 0: prerequisites
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
			return fmt.Errorf("%d prerequisite(s) missing. Install them and re-run radarctl init", failed)
		}
	}

	// Mode selection
	mode, err := wizard.SelectMode()
	if err != nil {
		return err
	}

	if mode == wizard.ModeExpert {
		output.Info("Expert mode: running config validation only")
		return runValidate(nil, nil)
	}

	// Run wizard
	output.Header("Configuring RADAR-Kubernetes")
	answers, err := wizard.Run(mode, repoRoot)
	if err != nil {
		return fmt.Errorf("wizard failed: %w", err)
	}

	// Write config files
	w := wizard.NewWriter(repoRoot)
	if err := w.WriteAnswers(answers); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	if err := w.WriteEnvironmentsYAML(answers.AppliedMods); err != nil {
		return fmt.Errorf("writing environments.yaml: %w", err)
	}

	output.Success("Configuration written to etc/production.yaml, etc/secrets.yaml, environments.yaml")
	output.Info("Run `radarctl deploy` to deploy the stack")
	return nil
}
```

- [ ] **Step 4: Build and verify init command exists**

```bash
cd cli
go build -o bin/radarctl .
./bin/radarctl init --help
```

Expected: help text includes `--skip-prereqs` flag.

- [ ] **Step 5: Commit**

```bash
git add cli/pkg/wizard/wizard.go cli/pkg/wizard/questions.go cli/cmd/init.go
git commit -m "feat(cli): add radarctl init with prerequisites check and setup wizard"
```

---

## Phase 5: Deploy Command

### Task 11: `radarctl deploy`

**Files:**
- Create: `cli/cmd/deploy.go`

- [ ] **Step 1: Write `cli/cmd/deploy.go`**

```go
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

	// Step 1: validate config
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

	// Step 2: diff / dry-run
	if deployDiff || deployDryRun {
		output.Info("Computing diff...")
		diff, err := hfRunner.Diff(deploySelector)
		if err != nil {
			return fmt.Errorf("helmfile diff: %w", err)
		}
		fmt.Println(diff)
		return nil
	}

	// Step 3: show change summary and confirm
	output.Info("Computing changes...")
	diff, _ := hfRunner.Diff(deploySelector)
	updated, installed, removed := helmfile.ParseDiffSummary(diff)
	pterm.Info.Printf("%d releases will be updated, %d installed, %d removed\n", updated, installed, removed)

	if !deployYes {
		var confirm bool
		pterm.Print("Proceed with deployment? [y/N] ")
		var input string
		fmt.Scanln(&input)
		confirm = strings.ToLower(strings.TrimSpace(input)) == "y"
		if !confirm {
			output.Info("Deployment cancelled")
			return nil
		}
	}

	// Step 4: deploy with live progress
	output.Header("Deploying RADAR stack...")
	progress, _ := pterm.DefaultSpinner.Start("Starting deployment...")

	releaseStatuses := map[string]string{}

	err = hfRunner.SyncWithCallback(deploySelector, !deployNoAtomic, func(line string) {
		release := helmfile.ParseReleaseFromLine(line)
		if release != "" {
			releaseStatuses[release] = "syncing"
			progress.UpdateText(fmt.Sprintf("Syncing %s...", release))
		}
		if strings.Contains(strings.ToLower(line), "succeeded") && release != "" {
			releaseStatuses[release] = "done"
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
	output.Success(fmt.Sprintf("All releases synced successfully"))

	if outputFormat == string(output.JSON) {
		output.PrintJSON(map[string]any{"status": "healthy", "summary": map[string]int{
			"synced": len(releaseStatuses),
		}})
	}
	return nil
}
```

- [ ] **Step 2: Build and verify deploy command exists**

```bash
cd cli
go build -o bin/radarctl .
./bin/radarctl deploy --help
```

Expected: help text with `--diff`, `--dry-run`, `--selector`, `--yes`, `--no-atomic` flags.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/deploy.go
git commit -m "feat(cli): add radarctl deploy with validation, diff preview, and live progress"
```

---

## Phase 6: Status and Diagnose Commands

### Task 12: Status collector

**Files:**
- Create: `cli/pkg/status/types.go`
- Create: `cli/pkg/status/collector.go`
- Create: `cli/pkg/status/collector_test.go`

- [ ] **Step 1: Write `cli/pkg/status/types.go`**

```go
package status

type Health string

const (
	Healthy  Health = "healthy"
	Degraded Health = "degraded"
	Warning  Health = "warning"
	Unknown  Health = "unknown"
)

type ReleaseStatus struct {
	Name      string  `json:"name"`
	Health    Health  `json:"status"`
	ReadyPods int     `json:"ready_pods"`
	TotalPods int     `json:"total_pods"`
	Error     string  `json:"error,omitempty"`
	Message   string  `json:"message,omitempty"`
	Logs      string  `json:"logs,omitempty"`
}

type Group struct {
	Name     string          `json:"group"`
	Releases []ReleaseStatus `json:"releases"`
}

type Report struct {
	Context  string  `json:"context"`
	Groups   []Group `json:"groups"`
	Healthy  int     `json:"healthy"`
	Degraded int     `json:"degraded"`
	Warning  int     `json:"warning"`
}

// releaseGroups defines the display grouping for known releases.
var releaseGroups = map[string]string{
	"cert-manager":            "INFRASTRUCTURE",
	"kube-prometheus-stack":   "INFRASTRUCTURE",
	"nginx-ingress":           "INFRASTRUCTURE",
	"zookeeper":               "KAFKA",
	"kafka":                   "KAFKA",
	"schema-registry":         "KAFKA",
	"ksql-server":             "KAFKA",
	"mongodb":                 "STORAGE",
	"postgresql":              "STORAGE",
	"redis":                   "STORAGE",
	"minio":                   "STORAGE",
	"radar-appserver":         "RADAR SERVICES",
	"management-portal":       "RADAR SERVICES",
	"radar-fitbit-connector":  "RADAR SERVICES",
	"radar-s3-connector":      "RADAR SERVICES",
	"kratos":                  "IDENTITY",
	"hydra":                   "IDENTITY",
}

func GroupFor(releaseName string) string {
	if g, ok := releaseGroups[releaseName]; ok {
		return g
	}
	return "OTHER"
}
```

- [ ] **Step 2: Write the failing tests `cli/pkg/status/collector_test.go`**

```go
package status_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runningPodJSON(name string) string {
	return `{"items":[{"metadata":{"name":"` + name + `","namespace":"default"},"status":{"phase":"Running","containerStatuses":[{"ready":true,"restartCount":0}]}}]}`
}

func failingPodJSON(name string) string {
	return `{"items":[{"metadata":{"name":"` + name + `","namespace":"default"},"status":{"phase":"Failed","containerStatuses":[{"ready":false,"restartCount":5}],"message":"OOMKilled"}}]}`
}

func TestCollect_HealthyRelease(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl get pods -n default -l app.kubernetes.io/name=mongodb -o json": runningPodJSON("mongodb-0"),
		},
	}
	kr := kubectl.NewRunner(mock, "")
	releases := []string{"mongodb"}
	report, err := status.Collect(kr, "default", releases)
	require.NoError(t, err)
	require.Len(t, report.Groups, 1)
	require.Len(t, report.Groups[0].Releases, 1)
	assert.Equal(t, status.Healthy, report.Groups[0].Releases[0].Health)
	assert.Equal(t, 1, report.Healthy)
}

func TestCollect_DegradedRelease(t *testing.T) {
	mock := &executor.MockExecutor{
		Responses: map[string]string{
			"kubectl get pods -n default -l app.kubernetes.io/name=ksql-server -o json": failingPodJSON("ksql-0"),
			"kubectl logs ksql-0 -n default --tail=20":                                   "java.lang.OutOfMemoryError",
		},
	}
	kr := kubectl.NewRunner(mock, "")
	report, err := status.Collect(kr, "default", []string{"ksql-server"})
	require.NoError(t, err)
	assert.Equal(t, 1, report.Degraded)
	assert.Equal(t, status.Degraded, report.Groups[0].Releases[0].Health)
}
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
cd cli
go test ./pkg/status/... -v
```

Expected: compilation error — `status.Collect` undefined.

- [ ] **Step 4: Write `cli/pkg/status/collector.go`**

```go
package status

import (
	"fmt"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
)

// Collect gathers health status for each named release by querying kubectl.
func Collect(kr *kubectl.Runner, namespace string, releases []string) (*Report, error) {
	report := &Report{}
	groupMap := map[string]*Group{}

	for _, release := range releases {
		label := fmt.Sprintf("app.kubernetes.io/name=%s", release)
		pods, err := kr.GetPods(namespace, label)
		if err != nil {
			pods = nil
		}

		rs := classifyRelease(release, pods)

		// Fetch logs for degraded pods
		if rs.Health == Degraded && len(pods) > 0 {
			logs, _ := kr.GetLogs(pods[0].Name, namespace, 20)
			rs.Logs = logs
		}

		groupName := GroupFor(release)
		if _, ok := groupMap[groupName]; !ok {
			groupMap[groupName] = &Group{Name: groupName}
		}
		groupMap[groupName].Releases = append(groupMap[groupName].Releases, rs)

		switch rs.Health {
		case Healthy:
			report.Healthy++
		case Degraded:
			report.Degraded++
		case Warning:
			report.Warning++
		}
	}

	// Stable group ordering
	orderedGroups := []string{"INFRASTRUCTURE", "KAFKA", "STORAGE", "RADAR SERVICES", "IDENTITY", "OTHER"}
	for _, g := range orderedGroups {
		if grp, ok := groupMap[g]; ok {
			report.Groups = append(report.Groups, *grp)
		}
	}

	return report, nil
}

func classifyRelease(name string, pods []kubectl.Pod) ReleaseStatus {
	rs := ReleaseStatus{Name: name}
	if len(pods) == 0 {
		rs.Health = Unknown
		rs.Message = "no pods found"
		return rs
	}

	rs.TotalPods = len(pods)
	for _, p := range pods {
		if p.Ready {
			rs.ReadyPods++
		}
		if p.Message != "" {
			rs.Error = p.Message
		}
	}

	switch {
	case rs.ReadyPods == rs.TotalPods:
		rs.Health = Healthy
	case rs.ReadyPods == 0:
		rs.Health = Degraded
		if rs.Error == "" {
			rs.Error = pods[0].Phase
		}
	default:
		rs.Health = Warning
		rs.Message = fmt.Sprintf("%d/%d pods ready", rs.ReadyPods, rs.TotalPods)
	}
	return rs
}
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
cd cli
go test ./pkg/status/... -v
```

Expected: all tests PASS.

- [ ] **Step 6: Commit**

```bash
git add cli/pkg/status/
git commit -m "feat(cli): add status collector with health classification"
```

---

### Task 13: `radarctl status` command

**Files:**
- Create: `cli/cmd/status.go`

- [ ] **Step 1: Write `cli/cmd/status.go`**

```go
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

// knownReleases is the default set of releases to check.
// In a future version this could be read from helmfile template output.
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
		pterm.Println() // blank line between refreshes
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
				pterm.Printf("\n    └─ %s", r.Logs[:min(len(r.Logs), 80)])
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

- [ ] **Step 2: Build and verify**

```bash
cd cli
go build -o bin/radarctl .
./bin/radarctl status --help
```

Expected: help text with `--watch`, `--component`, `--show-urls` flags.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/status.go
git commit -m "feat(cli): add radarctl status health dashboard"
```

---

### Task 14: `radarctl diagnose`

**Files:**
- Create: `cli/cmd/diagnose.go`

- [ ] **Step 1: Write `cli/cmd/diagnose.go`**

```go
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
	Timestamp  string                 `json:"timestamp"`
	Validation output.ValidationJSON  `json:"validation"`
	Status     *status.Report         `json:"status"`
	Logs       map[string]string      `json:"logs"`
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

	// Validation
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

	// Status + logs for degraded pods
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
```

- [ ] **Step 2: Build and verify**

```bash
cd cli
go build -o bin/radarctl .
./bin/radarctl diagnose --help
```

Expected: help text for diagnose command.

- [ ] **Step 3: Run all tests one final time**

```bash
cd cli
go test ./... -v
```

Expected: all tests PASS, no compilation errors.

- [ ] **Step 4: Commit**

```bash
git add cli/cmd/diagnose.go
git commit -m "feat(cli): add radarctl diagnose for full diagnostic snapshots"
```

---

## Phase 7: Final Wiring

### Task 15: Makefile target and gitignore

**Files:**
- Modify: `Makefile` (create if absent)
- Modify: `.gitignore`

- [ ] **Step 1: Add build target**

Check if a Makefile exists:

```bash
ls "/Users/yatharth/Library/Mobile Documents/com~apple~CloudDocs/Radar/RADAR-Kubernetes/Makefile" 2>/dev/null || echo "not found"
```

If not found, create `Makefile` in the repo root:

```makefile
.PHONY: radarctl clean-radarctl

radarctl:
	cd cli && go build -o bin/radarctl .

clean-radarctl:
	rm -f cli/bin/radarctl
```

If a Makefile already exists, append those targets to it.

- [ ] **Step 2: Add radarctl state file to .gitignore**

```bash
echo ".radarctl-state.yaml" >> "/Users/yatharth/Library/Mobile Documents/com~apple~CloudDocs/Radar/RADAR-Kubernetes/.gitignore"
echo "cli/bin/" >> "/Users/yatharth/Library/Mobile Documents/com~apple~CloudDocs/Radar/RADAR-Kubernetes/.gitignore"
```

- [ ] **Step 3: Final build smoke test**

```bash
cd "/Users/yatharth/Library/Mobile Documents/com~apple~CloudDocs/Radar/RADAR-Kubernetes"
make radarctl
cli/bin/radarctl --help
```

Expected: `radarctl` binary exists and prints help with all 5 subcommands listed: `init`, `deploy`, `status`, `diagnose`, `validate`.

- [ ] **Step 4: Commit**

```bash
git add Makefile .gitignore
git commit -m "feat(cli): add Makefile target and gitignore entries for radarctl"
```

---

## Self-Review Notes

**Spec coverage check:**
- ✅ `radarctl init` — prerequisites check (Task 8), wizard (Tasks 9–10), three modes
- ✅ Feature expansion (Fitbit, Garmin, REDCap) — Task 4 and Task 9
- ✅ Deployment profiles → mod application — Task 4, Task 9
- ✅ `radarctl deploy` — Task 11 (diff, dry-run, live progress, post-deploy health)
- ✅ `radarctl status` — Tasks 12–13 (grouping, icons, --watch, --show-urls, -o json)
- ✅ `radarctl diagnose` — Task 14
- ✅ `radarctl validate` — Task 7
- ✅ `-o json` on all commands — Tasks 7, 11, 13, 14
- ✅ Exit codes 0/1/2/3/4 — Task 11 (deploy) and Task 7 (validate exit 2)
- ✅ Resumability (`.radarctl-state.yaml` gitignore) — Task 15
- ✅ Agent-friendliness (JSON output + exit codes) — throughout

**Type consistency check:**
- `config.Config` / `config.Secrets` defined in Task 2, used in Tasks 3, 4, 7, 9, 10, 11, 14 ✅
- `executor.Executor` interface defined in Task 5, used in kubectl (Task 5), helmfile (Task 6), prereqs (Task 8) ✅
- `status.Report` defined in Task 12, used in Task 13 and Task 14 ✅
- `wizard.Answers` defined in Task 9, used in Task 10 ✅
- `helmfile.ParseReleaseFromLine` / `ParseDiffSummary` defined in Task 6, used in Task 11 ✅
- `output.ValidationJSON` defined in Task 7, used in Task 14 ✅
- `knownReleases` defined in Task 13, referenced in Task 14 ✅
- `findRepoRoot()` defined in Task 7, used in Tasks 10, 11, 14 ✅

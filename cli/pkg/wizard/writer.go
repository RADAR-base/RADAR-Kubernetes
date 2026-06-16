package wizard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Answers holds all values collected by the wizard.
type Answers struct {
	ServerName       string
	MaintainerEmail  string
	KubeContext      string
	Profile          string // "production", "staging", "dev", "demo"
	UseConfluent     bool
	ConfluentURL     string
	ConfluentKey     string
	ConfluentSecret  string
	Features         []string // enabled feature ids
	UseExternalS3    bool
	S3Bucket         string
	S3Region         string
	S3AccessKey      string
	S3SecretKey      string
	EnablePrometheus bool
	EnableGraylog    bool
	EnableKratos     bool
	Secrets          map[string]string
}

// Writer applies wizard answers as in-place overlays onto the files produced by
// bin/init (etc/production.yaml, etc/secrets.yaml, etc/production.yaml.gotmpl),
// preserving the upstream comments and structure. We use yq for production.yaml
// edits because yq is already a required prereq and preserves comments; secrets
// are merged via Go yaml.v3 because they're machine-only.
type Writer struct {
	repoRoot string
}

func NewWriter(repoRoot string) *Writer {
	return &Writer{repoRoot: repoRoot}
}

// ApplyProductionOverlay sets the wizard's chosen values inside etc/production.yaml,
// leaving all other keys (and the comments from base.yaml) intact. Requires that
// bin/init has already populated etc/production.yaml.
func (w *Writer) ApplyProductionOverlay(a *Answers) error {
	prodPath := filepath.Join(w.repoRoot, "etc", "production.yaml")
	if _, err := os.Stat(prodPath); err != nil {
		return fmt.Errorf("etc/production.yaml not found — run bin/init first: %w", err)
	}

	set := func(yqExpr string) error {
		out, err := exec.Command("yq", "-i", yqExpr, prodPath).CombinedOutput()
		if err != nil {
			return fmt.Errorf("yq %q failed: %s", yqExpr, string(out))
		}
		return nil
	}

	// Top-level deployment flags. Profile drives all four.
	flags := profileFlags(a.Profile)
	for expr, val := range flags {
		if err := set(fmt.Sprintf("%s = %s", expr, val)); err != nil {
			return err
		}
	}

	// User-supplied identifiers.
	if a.ServerName != "" {
		if err := set(fmt.Sprintf(".server_name = %q", a.ServerName)); err != nil {
			return err
		}
	}
	if a.MaintainerEmail != "" {
		if err := set(fmt.Sprintf(".maintainer_email = %q", a.MaintainerEmail)); err != nil {
			return err
		}
	}
	if a.KubeContext != "" {
		if err := set(fmt.Sprintf(".kubeContext = %q", a.KubeContext)); err != nil {
			return err
		}
	}

	if a.UseConfluent {
		if err := set(".confluent_cloud.enabled = true"); err != nil {
			return err
		}
	}

	if a.UseExternalS3 {
		if err := set(".external_s3 = true"); err != nil {
			return err
		}
		if err := set(fmt.Sprintf(".s3_bucket = %q", a.S3Bucket)); err != nil {
			return err
		}
		if err := set(fmt.Sprintf(".s3_region = %q", a.S3Region)); err != nil {
			return err
		}
	}

	// Monitoring master switch — drives mods/disable_monitoring_logging.yaml via environments.yaml.
	logMon := strconv.FormatBool(a.EnablePrometheus || a.EnableGraylog)
	if err := set(".enable_logging_monitoring = " + logMon); err != nil {
		return err
	}

	// Chart install toggles for enabled features.
	for _, expr := range featureChartExprs(a.Features, a.EnablePrometheus, a.EnableGraylog) {
		if err := set(expr); err != nil {
			return err
		}
	}

	// Firebase JSON for aRMT/appserver: copy the user-provided file into place
	// and uncomment the relevant block in production.yaml.gotmpl.
	if containsFeature(a.Features, "armt") {
		if jsonPath := a.Secrets["armt_firebase_json_path"]; jsonPath != "" {
			dst := filepath.Join(w.repoRoot, "etc", "radar-appserver", "firebase-adminsdk.json")
			if err := copyFile(jsonPath, dst); err != nil {
				return fmt.Errorf("copying firebase credentials: %w", err)
			}
			if err := w.uncommentAppserverFirebaseBlock(); err != nil {
				return fmt.Errorf("uncommenting appserver firebase block: %w", err)
			}
		}
	}

	return nil
}

// profileFlags returns the yq expression → value pairs that encode the deployment
// profile in production.yaml. Aligned with upstream flags (no minimal_install —
// that key isn't in base.yaml; dev_deployment cascades minimal+localdev+fast_deploy
// via environments.yaml.tmpl).
func profileFlags(profile string) map[string]string {
	switch profile {
	case "demo":
		return map[string]string{
			".enable_tls":               "false",
			".dev_deployment":           "true",
			".atomicInstall":            "false",
			".kafka_num_brokers":        "1",
			".kafka_num_topic_replicas": "1",
			".base_timeout":             "300", // k3d image pulls are slow on first run
		}
	case "dev":
		return map[string]string{
			".enable_tls":               "false",
			".dev_deployment":           "true",
			".kafka_num_brokers":        "1",
			".kafka_num_topic_replicas": "1",
			".base_timeout":             "300",
		}
	case "staging":
		return map[string]string{
			".enable_tls":        "true",
			".dev_deployment":    "false",
			".kafka_num_brokers": "1",
		}
	default: // production
		return map[string]string{
			".enable_tls":        "true",
			".dev_deployment":    "false",
			".kafka_num_brokers": "3",
		}
	}
}

// featureChartExprs returns the list of yq expressions that flip `_install: true`
// on each chart required by the user's selected features (plus monitoring/logging).
func featureChartExprs(features []string, prom, gray bool) []string {
	exprs := []string{}
	add := func(chart string) { exprs = append(exprs, "."+chart+"._install = true") }

	if prom {
		add("kube_prometheus_stack")
	}
	if gray {
		add("graylog")
		add("elasticsearch")
		add("mongodb")
		add("fluent_bit")
	}

	for _, f := range features {
		switch f {
		case "fitbit":
			add("radar_fitbit_connector")
			add("radar_rest_sources_auth_backend")
			add("radar_rest_sources_authorizer")
		case "garmin":
			add("radar_garmin_connector")
		case "oura":
			add("radar_oura_connector")
			add("radar_rest_sources_auth_backend")
			add("radar_rest_sources_authorizer")
		case "armt":
			add("radar_appserver")
		case "realtime_dashboards":
			add("ksql_server")
			add("radar_grafana")
			add("radar_jdbc_connector_realtime_dashboard")
		case "kratos":
			add("radar_kratos")
			add("radar_hydra")
		}
	}
	return exprs
}

// MergeWizardSecrets reads etc/secrets.yaml (seeded by bin/generate-secrets with
// strong random passwords) and overlays only the keys the wizard collected from
// the user. All other generated passwords are preserved as-is.
func (w *Writer) MergeWizardSecrets(a *Answers) error {
	secPath := filepath.Join(w.repoRoot, "etc", "secrets.yaml")

	root := make(map[string]any)
	if data, err := os.ReadFile(secPath); err == nil {
		if err := yaml.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("parsing existing %s: %w", secPath, err)
		}
	}

	setKey := func(path []string, value string) {
		if value == "" {
			return
		}
		cur := root
		for i, k := range path {
			if i == len(path)-1 {
				cur[k] = value
				return
			}
			next, ok := cur[k].(map[string]any)
			if !ok {
				next = make(map[string]any)
				cur[k] = next
			}
			cur = next
		}
	}

	if a.UseConfluent {
		setKey([]string{"confluent_cloud", "bootstrapServerurl"}, a.ConfluentURL)
		setKey([]string{"confluent_cloud", "apiKey"}, a.ConfluentKey)
		setKey([]string{"confluent_cloud", "apiSecret"}, a.ConfluentSecret)
	}
	setKey([]string{"fitbit_client_id"}, a.Secrets["fitbit_client_id"])
	setKey([]string{"fitbit_client_secret"}, a.Secrets["fitbit_client_secret"])
	setKey([]string{"garmin_consumer_key"}, a.Secrets["garmin_consumer_key"])
	setKey([]string{"garmin_consumer_secret"}, a.Secrets["garmin_consumer_secret"])
	setKey([]string{"oura_api_client"}, a.Secrets["oura_api_client"])
	setKey([]string{"oura_api_secret"}, a.Secrets["oura_api_secret"])
	if a.UseExternalS3 {
		setKey([]string{"s3_access_key"}, a.S3AccessKey)
		setKey([]string{"s3_secret_key"}, a.S3SecretKey)
	}

	out, err := yaml.Marshal(root)
	if err != nil {
		return fmt.Errorf("marshalling secrets: %w", err)
	}
	return os.WriteFile(secPath, out, 0600)
}

// uncommentAppserverFirebaseBlock flips the example block in
// etc/production.yaml.gotmpl from commented (lines wrapped in `{{/* ... */}}`)
// to active so radar_appserver loads firebase-adminsdk.json at install time.
func (w *Writer) uncommentAppserverFirebaseBlock() error {
	path := filepath.Join(w.repoRoot, "etc", "production.yaml.gotmpl")
	data, err := os.ReadFile(path)
	if err != nil {
		// File is optional in some flows; just warn via error return for caller to log.
		return fmt.Errorf("reading %s: %w", path, err)
	}
	// Idempotent: only edit if the block is still commented.
	// The upstream template wraps the block with `{{/*` and `*/}}` markers.
	out := uncommentFirebaseBlock(string(data))
	if out == string(data) {
		return nil
	}
	return os.WriteFile(path, []byte(out), 0644)
}

func containsFeature(features []string, name string) bool {
	for _, f := range features {
		if f == name {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading %s: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0600)
}

// SaveState serializes Answers to a YAML file for resumability.
func SaveState(a *Answers, path string) error {
	data, err := yaml.Marshal(a)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadState deserializes Answers from a YAML state file.
func LoadState(path string) (*Answers, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var a Answers
	if err := yaml.Unmarshal(data, &a); err != nil {
		return nil, err
	}
	if a.Secrets == nil {
		a.Secrets = make(map[string]string)
	}
	return &a, nil
}

package config

// ConfluentCloudConfig mirrors the confluent_cloud map used in base.yaml and helmfile
// templates (e.g. .Values.confluent_cloud.enabled). The cc sub-key holds credentials
// that are written by the wizard into secrets.yaml under the same yaml path.
type ConfluentCloudConfig struct {
	Enabled bool `yaml:"enabled"`
}

// ChartToggle is the standard `_install` nested key every helm chart in this repo
// reads to decide whether it should be installed. Used as a pointer in Config so
// `omitempty` skips charts the wizard didn't touch (chart's base.yaml default applies).
type ChartToggle struct {
	Install bool `yaml:"_install"`
}

// Appserver carries the chart toggle plus the Firebase credentials JSON string
// (consumed by radar_appserver.google_application_credentials).
type Appserver struct {
	Install                     bool   `yaml:"_install"`
	GoogleApplicationCredentials string `yaml:"google_application_credentials,omitempty"`
}

// Config holds values from production.yaml (non-secret).
type Config struct {
	AtomicInstall      bool   `yaml:"atomicInstall"`
	BaseTimeout        int    `yaml:"base_timeout"`
	KubeContext        string `yaml:"kubeContext"`
	ServerName         string `yaml:"server_name"`
	MaintainerEmail    string `yaml:"maintainer_email"`
	KafkaNumBrokers    int    `yaml:"kafka_num_brokers"`
	KafkaNumReplicas   int    `yaml:"kafka_num_topic_replicas"`
	KafkaNumPartitions int    `yaml:"kafka_num_topic_partitions,omitempty"`
	EnableTLS          bool   `yaml:"enable_tls"`
	EnableLogging      bool   `yaml:"enable_logging_monitoring"`
	DevDeployment      bool   `yaml:"dev_deployment"`
	ConfluentCloud     ConfluentCloudConfig `yaml:"confluent_cloud"`
	UseExternalS3      bool   `yaml:"external_s3"`
	S3Bucket           string `yaml:"s3_bucket,omitempty"`
	S3Region           string `yaml:"s3_region,omitempty"`

	// Chart install toggles — nested keys read by helm charts. nil => chart default
	// from etc/base.yaml applies (almost always _install: false).
	RadarFitbitConnector                *ChartToggle `yaml:"radar_fitbit_connector,omitempty"`
	RadarGarminConnector                *ChartToggle `yaml:"radar_garmin_connector,omitempty"`
	RadarOuraConnector                  *ChartToggle `yaml:"radar_oura_connector,omitempty"`
	RadarRestSourcesAuthBackend         *ChartToggle `yaml:"radar_rest_sources_auth_backend,omitempty"`
	RadarRestSourcesAuthorizer          *ChartToggle `yaml:"radar_rest_sources_authorizer,omitempty"`
	RadarAppserver                      *Appserver   `yaml:"radar_appserver,omitempty"`
	KsqlServer                          *ChartToggle `yaml:"ksql_server,omitempty"`
	RadarGrafana                        *ChartToggle `yaml:"radar_grafana,omitempty"`
	RadarJdbcConnectorRealtimeDashboard *ChartToggle `yaml:"radar_jdbc_connector_realtime_dashboard,omitempty"`
	RadarKratos                         *ChartToggle `yaml:"radar_kratos,omitempty"`
	RadarHydra                          *ChartToggle `yaml:"radar_hydra,omitempty"`
	KubePrometheusStack                 *ChartToggle `yaml:"kube_prometheus_stack,omitempty"`
	Graylog                             *ChartToggle `yaml:"graylog,omitempty"`

	Minio struct {
		ExternalEndpoint string `yaml:"externalEndpoint,omitempty"`
	} `yaml:"minio,omitempty"`
}

// Secrets holds values from secrets.yaml (sensitive).
type Secrets struct {
	ConfluentCloud struct {
		BootstrapServer string `yaml:"bootstrapServerurl,omitempty"`
		APIKey          string `yaml:"apiKey,omitempty"`
		APISecret       string `yaml:"apiSecret,omitempty"`
	} `yaml:"confluent_cloud,omitempty"`
	FitbitClientID       string `yaml:"fitbit_client_id,omitempty"`
	FitbitClientSecret   string `yaml:"fitbit_client_secret,omitempty"`
	GarminConsumerKey    string `yaml:"garmin_consumer_key,omitempty"`
	GarminConsumerSecret string `yaml:"garmin_consumer_secret,omitempty"`
	OuraAPIClient        string `yaml:"oura_api_client,omitempty"`
	OuraAPISecret        string `yaml:"oura_api_secret,omitempty"`
	S3AccessKey          string `yaml:"s3_access_key,omitempty"`
	S3SecretKey          string `yaml:"s3_secret_key,omitempty"`
}

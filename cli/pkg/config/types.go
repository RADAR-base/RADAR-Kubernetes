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
	UseExternalS3      bool    `yaml:"external_s3"`
	S3Bucket           string  `yaml:"s3_bucket"`
	S3Region           string  `yaml:"s3_region"`
	EnablePrometheus   bool    `yaml:"radar_prometheus_install"`
	EnableGraylog      bool    `yaml:"radar_graylog_install"`
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
	FitbitClientID       string `yaml:"fitbit_client_id"`
	FitbitClientSecret   string `yaml:"fitbit_client_secret"`
	GarminConsumerKey    string `yaml:"garmin_consumer_key"`
	GarminConsumerSecret string `yaml:"garmin_consumer_secret"`
	RedcapToken          string `yaml:"redcap_token"`
	S3AccessKey          string `yaml:"s3_access_key"`
	S3SecretKey          string `yaml:"s3_secret_key"`
}

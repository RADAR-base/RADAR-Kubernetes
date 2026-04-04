

# radar-iceberg-connector
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/radar-iceberg-connector)](https://artifacthub.io/packages/helm/radar-base/radar-iceberg-connector)

![Version: 0.7.3](https://img.shields.io/badge/Version-0.7.3-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 7.8.4](https://img.shields.io/badge/AppVersion-7.8.4-informational?style=flat-square)

A Helm chart for RADAR-base s3 connector. This connector uses Confluent s3 connector with a custom data transformers. These configurations enable a sink connector. See full list of properties here https://docs.confluent.io/kafka-connect-s3-sink/current/configuration_options.html#s3-configuration-options

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Pim van Nierop | <pim@thehyve.nl> | <https://www.thehyve.nl/experts/pim-van-nierop> |

## Source Code

* <https://github.com/RADAR-base/radar-helm-charts/tree/main/charts/radar-iceberg-connector>
* <https://github.com/RADAR-base/kafka-connect-transform-keyvalue>
* <https://docs.confluent.io/kafka-connect-s3-sink/current/configuration_options.html#s3-configuration-options>

## Prerequisites
* Kubernetes 1.28+
* Kubectl 1.28+
* Helm 3.1.0+

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://radar-base.github.io/radar-helm-charts | common | 2.x.x |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` | Number of radar-iceberg-connector replicas to deploy |
| image.registry | string | `"ghcr.io"` | Image registry |
| image.repository | string | `"radar-base/kafka-connect-transform-keyvalue/kafka-connect-transform-s3"` | Image repository |
| image.tag | string | `nil` | Image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.digest | string | `""` | Image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| image.pullSecrets | list | `[]` | Optionally specify an array of imagePullSecrets. Secrets must be manually created in the namespace. e.g: pullSecrets:   - myRegistryKeySecretName  |
| nameOverride | string | `""` | String to partially override radar-iceberg-connector.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override radar-iceberg-connector.fullname template with a string |
| podSecurityContext | object | `{}` | Configure radar-iceberg-connector pods' Security Context |
| kafkaHeapOpts | string | `"-Xms3g -Xmx4g"` | JVM Options |
| securityContext | object | `{}` | Configure radar-iceberg-connector containers' Security Context |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `8083` | radar-iceberg-connector port |
| resources.requests | object | `{"cpu":"100m","memory":"3Gi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| extraEnvVars | list | `[]` | Additional environment variables to pass to the connector. These can be used to pass supported kafka and connect specific [configs](https://docs.confluent.io/platform/current/installation/docker/config-reference.html#kconnect-long-configuration) |
| customLivenessProbe | object | `{}` | Custom livenessProbe that overrides the default one |
| livenessProbe.enabled | bool | `true` | Enable livenessProbe |
| livenessProbe.initialDelaySeconds | int | `5` | Initial delay seconds for livenessProbe |
| livenessProbe.periodSeconds | int | `30` | Period seconds for livenessProbe |
| livenessProbe.timeoutSeconds | int | `5` | Timeout seconds for livenessProbe |
| livenessProbe.successThreshold | int | `1` | Success threshold for livenessProbe |
| livenessProbe.failureThreshold | int | `3` | Failure threshold for livenessProbe |
| customReadinessProbe | object | `{}` | Custom readinessProbe that overrides the default one |
| readinessProbe.enabled | bool | `true` | Enable readinessProbe |
| readinessProbe.initialDelaySeconds | int | `5` | Initial delay seconds for readinessProbe |
| readinessProbe.periodSeconds | int | `30` | Period seconds for readinessProbe |
| readinessProbe.timeoutSeconds | int | `5` | Timeout seconds for readinessProbe |
| readinessProbe.successThreshold | int | `1` | Success threshold for readinessProbe |
| readinessProbe.failureThreshold | int | `3` | Failure threshold for readinessProbe |
| customStartupProbe | object | `{}` | Custom startupProbe that overrides the default one |
| startupProbe.enabled | bool | `true` | Enable startupProbe |
| startupProbe.initialDelaySeconds | int | `5` | Initial delay seconds for startupProbe |
| startupProbe.periodSeconds | int | `10` | Period seconds for startupProbe |
| startupProbe.timeoutSeconds | int | `10` | Timeout seconds for startupProbe |
| startupProbe.successThreshold | int | `1` | Success threshold for startupProbe |
| startupProbe.failureThreshold | int | `30` | Failure threshold for startupProbe |
| networkpolicy | object | check `values.yaml` | Network policy defines who can access this application and who this applications has access to |
| kafka_num_brokers | int | `3` | Number of deployed Kafka broker instances |
| kafka.url | string | `"SASL_PLAINTEXT://radar-kafka-kafka-bootstrap:9094"` | Kafka broker URLs |
| kafka.security | object | `{"env":[{"name":"CONNECT_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"},{"name":"CONNECT_SASL_MECHANISM","value":"SCRAM-SHA-512"},{"name":"CONNECT_PRODUCER_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"},{"name":"CONNECT_PRODUCER_SASL_MECHANISM","value":"SCRAM-SHA-512"},{"name":"CONNECT_CONSUMER_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"},{"name":"CONNECT_CONSUMER_SASL_MECHANISM","value":"SCRAM-SHA-512"}],"jaasSecret":{"key":"sasl.jaas.config","name":"shared-service-user"}}` | Security related env vars set in pods. Not applied when `cc.enabled` is true. |
| kafka.security.env | list | `[{"name":"CONNECT_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"},{"name":"CONNECT_SASL_MECHANISM","value":"SCRAM-SHA-512"},{"name":"CONNECT_PRODUCER_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"},{"name":"CONNECT_PRODUCER_SASL_MECHANISM","value":"SCRAM-SHA-512"},{"name":"CONNECT_CONSUMER_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"},{"name":"CONNECT_CONSUMER_SASL_MECHANISM","value":"SCRAM-SHA-512"}]` | Env vars set for authentication with Kafka brokers. Not applied when `cc.enabled` is true. |
| kafka.security.env[0] | object | `{"name":"CONNECT_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"}` | Protocol used to communicate with brokers. Valid values are: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL. |
| kafka.security.env[1] | object | `{"name":"CONNECT_SASL_MECHANISM","value":"SCRAM-SHA-512"}` | Mechanism used to authenticate with SASL. Valid values are: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512. |
| kafka.security.env[2] | object | `{"name":"CONNECT_PRODUCER_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"}` | Protocol used to communicate with brokers. Valid values are: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL. |
| kafka.security.env[3] | object | `{"name":"CONNECT_PRODUCER_SASL_MECHANISM","value":"SCRAM-SHA-512"}` | Mechanism used to authenticate with SASL. Valid values are: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512. |
| kafka.security.env[4] | object | `{"name":"CONNECT_CONSUMER_SECURITY_PROTOCOL","value":"SASL_PLAINTEXT"}` | Protocol used to communicate with brokers. Valid values are: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL. |
| kafka.security.env[5] | object | `{"name":"CONNECT_CONSUMER_SASL_MECHANISM","value":"SCRAM-SHA-512"}` | Mechanism used to authenticate with SASL. Valid values are: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512. |
| kafka.security.jaasSecret | object | `{"key":"sasl.jaas.config","name":"shared-service-user"}` | Secret for the Kafka SASL JAAS configuration. Not applied when `cc.enabled` is true. |
| schemaRegistry.url | string | `"http://radar-kafka-schema-registry:8081"` | Schema registry URL |
| catalogServer.url | string | `"http://catalog-server:9010"` | Catalog server URL |
| topics | string | `""` | List of topics to be consumed by the sink connector separated by comma. Topics defined in the catalog server will automatically be loaded if `initTopics.enabled` is true. |
| s3Endpoint | string | `"http://radar-seaweedfs-filer:8333/"` | Target S3 endpoint url |
| s3Tagging | bool | `false` | set to true, if S3 objects should be tagged with start and end offsets, as well as record count. |
| s3PartSize | int | `5242880` | The Part Size in S3 Multi-part Uploads. |
| s3Region | string | `nil` | The AWS region to be used the connector. Some compatibility layers require this. |
| flushSize | int | `10000` | Number of records written to store before invoking file commits. |
| rotateInterval | int | `900000` | The time interval in milliseconds to invoke file commits. |
| maxTasks | int | `4` | Number of tasks in the connector |
| bucketAccessKey | string | `"access_key"` | Access key of the target S3 bucket |
| bucketSecretKey | string | `"secret"` | Secret key of the target S3 bucket |
| bucketName | string | `"radar_intermediate_storage"` | Bucket name of the target S3 bucket |
| cc.enabled | bool | `false` | Set to true, if Confluent Cloud is used |
| cc.bootstrapServerurl | string | `""` | Confluent cloud based Kafka broker URL (if Confluent Cloud based Kafka cluster is used) |
| cc.schemaRegistryUrl | string | `""` | Confluent cloud based Schema registry URL (if Confluent Cloud based Schema registry is used) |
| cc.apiKey | string | `"ccApikey"` | API Key of the Confluent Cloud cluster |
| cc.apiSecret | string | `"ccApiSecret"` | API secret of the Confluent Cloud cluster |
| cc.schemaRegistryApiKey | string | `"srApiKey"` | API Key of the Confluent Cloud Schema registry |
| cc.schemaRegistryApiSecret | string | `"srApiSecret"` | API Key of the Confluent Cloud Schema registry |
| initTopics.enabled | bool | `true` | If true, fetch list of topics from catalog server |
| initTopics.image.repository | string | `"linuxserver/yq"` | Image repository to fetch topics with |
| initTopics.image.tag | string | `"3.2.2"` | Image tag to fetch topics with |
| initTopics.image.pullPolicy | string | `"IfNotPresent"` | Image pull policy to fetch topics with |

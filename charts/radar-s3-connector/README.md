

# radar-s3-connector

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for RADAR-base s3 connector. This connector uses Confluent s3 connector with a custom data transformers. These configurations enable a sink connector. See full list of properties here https://docs.confluent.io/kafka-connect-s3-sink/current/configuration_options.html#s3-configuration-options

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Source Code

* <https://github.com/RADAR-base/kafka-connect-transform-keyvalue>
* <https://docs.confluent.io/kafka-connect-s3-sink/current/configuration_options.html#s3-configuration-options>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` | Number of radar-s3-connector replicas to deploy |
| image.repository | string | `"radarbase/kafka-connect-transform-s3"` | radar-s3-connector image repository |
| image.tag | string | `"5.5.1"` | radar-s3-connector image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | radar-s3-connector image pull policy |
| nameOverride | string | `""` | String to partially override radar-s3-connector.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override radar-s3-connector.fullname template with a string |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `8083` | radar-s3-connector port |
| ingress.enabled | bool | `false` | Enable ingress controller resource |
| ingress.annotations | object | `{}` | Annotations to define default ingress class, certificate issuer |
| ingress.hosts | list | `[{"host":"chart-example.local","paths":[]}]` | Hosts to listen to incoming requests |
| ingress.tls | list | `[]` | Utilize TLS backend in ingress |
| resources.requests | object | `{"cpu":"100m","memory":"3Gi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| kafka.url | string | `"PLAINTEXT://cp-kafka-headless:9092"` | Kafka broker URLs |
| schemaRegistry.url | string | `"http://cp-schema-registry:8081"` |  |
| topics | string | check values.yaml | List of topics to be consumed by the sink connector separated by comma. |
| s3Endpoint | string | `"http://minio:9000/"` | Target S3 endpoint url |
| s3Tagging | bool | `false` | set to true, if S3 objects should be tagged with start and end offsets, as well as record count. |
| s3PartSize | int | `5242880` | The Part Size in S3 Multi-part Uploads. |
| s3Region | string | `nil` | The AWS region to be used the connector. |
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

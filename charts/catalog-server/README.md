

# catalog-server

![Version: 0.2.1](https://img.shields.io/badge/Version-0.2.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.7.1](https://img.shields.io/badge/AppVersion-0.7.1-informational?style=flat-square)

A Helm chart for RADAR-base catalogue server

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Source Code

* <https://github.com/RADAR-base/RADAR-Schemas/tree/master/java-sdk>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+
* PV provisioner support in the underlying infrastructure

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` | Number of catalog-server replicas to deploy |
| image.repository | string | `"radarbase/radar-schemas-tools"` | catalog-server image repository |
| image.tag | string | `"0.7.1"` | catalog-server image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | catalog-server image pull policy |
| nameOverride | string | `""` | String to partially override management-portal.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override management-portal.fullname template with a string |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `9010` | catalog-server port |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"256Mi"` |  |
| persistence.enabled | bool | `true` | Enable persistence using PVC |
| persistence.accessMode | string | `"ReadWriteOnce"` |  |
| persistence.size | string | `"5Mi"` | PVC Storage Request for catalog-server volume |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| kafka_num_brokers | int | `3` | number of Kafka brokers to look for |
| kafka | string | `"cp-kafka-headless:9092"` | URI of Kafka brokers |
| schema_registry | string | `"http://cp-schema-registry:8081"` | URL of the confluent schema registry |
| specificationsExclude | string | `nil` | List of paths of specifications relative to specifications folder, if any of the specifications should be excluded from automatically registering topics and schemas. |
| cc.enabled | bool | `false` | set to true if using Confluent Cloud for kafka cluster and schema registry |
| cc.bootstrapServerurl | string | `"confluent-url"` | URL of the bootstrap server of Confluent Cloud based kafka cluster |
| cc.apiKey | string | `"ccApikey"` | API key of the Confluent Cloud based kafka cluster |
| cc.apiSecret | string | `"ccApiSecret"` | API secret of the Confluent Cloud based kafka cluster |
| cc.schemaRegistryApiKey | string | `"srApiKey"` | API key of the Confluent Cloud based schema registry |
| cc.schemaRegistryApiSecret | string | `"srApiSecret"` | API secret of the Confluent Cloud based schema registry |

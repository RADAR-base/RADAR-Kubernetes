

# catalog-server

![Version: 0.2.1](https://img.shields.io/badge/Version-0.2.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.7.1](https://img.shields.io/badge/AppVersion-0.7.1-informational?style=flat-square)

A Helm chart for RADAR-base catalogue server

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

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
| replicaCount | int | `2` |  |
| image.repository | string | `"radarbase/radar-schemas-tools"` |  |
| image.tag | string | `"0.7.1"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| service.type | string | `"ClusterIP"` |  |
| service.port | int | `9010` |  |
| ingress.enabled | bool | `false` |  |
| ingress.annotations | object | `{}` |  |
| ingress.hosts[0].host | string | `"chart-example.local"` |  |
| ingress.hosts[0].paths | list | `[]` |  |
| ingress.tls | list | `[]` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"256Mi"` |  |
| persistence.enabled | bool | `true` |  |
| persistence.accessMode | string | `"ReadWriteOnce"` |  |
| persistence.size | string | `"5Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
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

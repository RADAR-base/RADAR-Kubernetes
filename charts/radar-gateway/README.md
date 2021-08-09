

# radar-gateway

![Version: 0.1.3](https://img.shields.io/badge/Version-0.1.3-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.6](https://img.shields.io/badge/AppVersion-0.5.6-informational?style=flat-square)

A Helm chart for RADAR-base gateway. For more details of the configurations, see https://github.com/RADAR-base/RADAR-Gateway/blob/master/gateway.yml.

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/RADAR-Gateway>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` |  |
| image.repository | string | `"radarbase/radar-gateway"` |  |
| image.tag | string | `"0.5.6"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| service.type | string | `"ClusterIP"` |  |
| service.port | int | `8080` |  |
| ingress.enabled | bool | `true` |  |
| ingress.annotations."kubernetes.io/ingress.class" | string | `"nginx"` |  |
| ingress.annotations."cert-manager.io/cluster-issuer" | string | `"letsencrypt-prod"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/rewrite-target" | string | `"/kafka/$1"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/server-snippet" | string | `"location ^~ /kafka/consumers {\n  deny all;\n}\nlocation ^~ /kafka/brokers {\n  deny all;\n}\nlocation ~* /kafka/topics/.+/partitions {\n  deny all;\n}\n"` |  |
| ingress.path | string | `"/kafka/?(.*)"` |  |
| ingress.hosts[0] | string | `"localhost"` |  |
| ingress.tls.secretName | string | `"radar-base-tls"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"128Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
| serviceMonitor.enabled | bool | `true` |  |
| managementportalHost | string | `"management-portal"` | Host name of the management portal application |
| schemaRegistry | string | `"http://cp-schema-registry:8081"` | Schema Registry URL |
| max_requests | int | `1000` | Not used. To be confirmed |
| bootstrapServers | string | `"cp-kafka-headless:9092"` | Kafka broker URLs |
| checkSourceId | bool | `true` | set to true, if sources in access token should be validated |
| adminProperties | object | `{}` | Additional Kafka Admin Client settings as key value pairs. Read from https://kafka.apache.org/documentation/#adminclientconfigs. |
| producerProperties | object | `{"compression.type":"lz4"}` | Kafka producer properties as key value pairs. Read from https://kafka.apache.org/documentation/#producerconfigs. |
| serializationProperties | object | `{}` | Additional Kafka serialization settings, used in KafkaAvroSerializer. Read from [io.confluent.kafka.serializers.AbstractKafkaSchemaSerDeConfig]. |
| cc.enabled | bool | `false` | set to true, if requests should be forwarded to Confluent Cloud based brokers. |
| cc.apiKey | string | `"ccApikey"` | Confluent Cloud cluster API key |
| cc.apiSecret | string | `"ccApiSecret"` | Confluent Cloud cluster API secret |
| cc.schemaRegistryApiKey | string | `"srApiKey"` | Confluent Cloud schema registry API key |
| cc.schemaRegistryApiSecret | string | `"srApiSecret"` | Confluent Cloud schema registry API secret |

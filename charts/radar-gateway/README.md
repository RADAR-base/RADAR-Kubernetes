

# radar-gateway

![Version: 0.1.3](https://img.shields.io/badge/Version-0.1.3-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.6](https://img.shields.io/badge/AppVersion-0.5.6-informational?style=flat-square)

A Helm chart for RADAR-base gateway.

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
| managementportalHost | string | `"management-portal"` |  |
| schemaRegistry | string | `"http://cp-schema-registry:8081"` |  |
| max_requests | int | `1000` |  |
| bootstrapServers | string | `"cp-kafka-headless:9092"` |  |
| checkSourceId | bool | `true` |  |
| adminProperties | object | `{}` |  |
| producerProperties."compression.type" | string | `"lz4"` |  |
| serializationProperties | object | `{}` |  |
| cc.enabled | bool | `false` |  |
| cc.apiKey | string | `"ccApikey"` |  |
| cc.apiSecret | string | `"ccApiSecret"` |  |
| cc.schemaRegistryApiKey | string | `"srApiKey"` |  |
| cc.schemaRegistryApiSecret | string | `"srApiSecret"` |  |

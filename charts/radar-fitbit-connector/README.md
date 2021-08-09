

# radar-fitbit-connector

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for RADAR-base fitbit connector

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/RADAR-REST-Connector>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+
* PV provisioner support in the underlying infrastructure

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` | Number of radar-fitbit-connector replicas to deploy |
| image.repository | string | `"radarbase/kafka-connect-rest-fitbit-source"` | radar-fitbit-connector image repository |
| image.tag | string | `"0.3.3"` | radar-fitbit-connector image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | radar-fitbit-connector image pull policy |
| nameOverride | string | `""` | String to partially override radar-fitbit-connector.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override radar-fitbit-connector.fullname template with a string |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `8083` | radar-fitbit-connector port |
| ingress.enabled | bool | `false` | Enable ingress controller resource |
| ingress.annotations | object | `{}` | Annotations to define default ingress class, certificate issuer |
| ingress.hosts | list | check values.yaml | Hosts to listen to incoming requests |
| ingress.tls | list | `[]` | TLS secrets for certificates |
| resources.requests | object | `{"cpu":"100m","memory":"1Gi"}` | CPU/Memory resource requests |
| persistence.enabled | bool | `true` | Enable persistence using PVC |
| persistence.accessMode | string | `"ReadWriteOnce"` | PVC Access Mode for radar-fitbit-connector volume |
| persistence.size | string | `"5Gi"` | PVC Storage Request for radar-fitbit-connector volume |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| zookeeper | string | `"cp-zookeeper-headless:2181"` | URI of Zookeeper instances of the cluster |
| kafka | string | `"PLAINTEXT://cp-kafka-headless:9092"` | URI of Kafka brokers of the cluster |
| kafka_num_brokers | string | `"3"` | Number of Kafka brokers. This is used to validate the cluster availability at connector init. |
| schema_registry | string | `"http://cp-schema-registry:8081"` | URL of the Kafka schema registry |
| radar_rest_sources_backend_url | string | `"http://radar-rest-sources-backend:8080/rest-sources/backend/"` | Base URL of the rest-sources-authorizer-backend service |
| connector_num_tasks | string | `"5"` | Number of connector tasks to be used in connector.properties |
| fitbit_api_client | string | `""` | Fitbit API client id. |
| fitbit_api_secret | string | `""` | Fitbit API client secret. |
| oauthClientId | string | `"radar_fitbit_connector"` | OAuth2 client id from Management Portal |
| oauthClientSecret | string | `"secret"` | OAuth2 client secret from Management Portal |
| managementportal_host | string | `"management-portal"` | Hostname of Management Portal. This will be used to create URLs to access Management Portal |
| includeIntradayData | bool | `true` | Set to true, if intraday access data should be collected by the connector. This will be set in connector.properties. |

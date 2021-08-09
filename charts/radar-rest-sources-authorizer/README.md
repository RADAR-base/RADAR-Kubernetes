

# radar-rest-sources-authorizer

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 3.2.0](https://img.shields.io/badge/AppVersion-3.2.0-informational?style=flat-square)

A Helm chart for the front-end application of RADAR-base Rest Sources Authorizer

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/RADAR-Rest-Source-Auth>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` | Number of radar-fitbit-connector replicas to deploy |
| image.repository | string | `"radarbase/radar-rest-source-authorizer"` | radar-rest-sources-authorizer image repository |
| image.tag | string | `"3.2.0"` | radar-rest-sources-authorizer image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | radar-rest-sources-authorizer image pull policy |
| nameOverride | string | `""` | String to partially override radar-rest-sources-authorizer.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override radar-rest-sources-authorizer.fullname template with a string |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `80` | radar-rest-sources-authorizer port |
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.annotations | object | check values.yaml | Annotations that define default ingress class, certificate issuer |
| ingress.path | string | `"/rest-sources/authorizer/?(.*)"` | Path within the url structure |
| ingress.hosts | list | `["localhost"]` | Hosts to accept requests from |
| ingress.tls.secretName | string | `"radar-base-tls"` | TLS Secret Name |
| resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| clientId | string | `"radar_rest_sources_authorizer"` | OAuth2 client id of the application registered in Management Portal. It is assumed that this is a public client with empty client secret. |
| serverName | string | `"localhost"` | Domain name of the server |

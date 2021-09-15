

# radar-upload-connect-frontend

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.9](https://img.shields.io/badge/AppVersion-0.5.9-informational?style=flat-square)

A Helm chart for RADAR-base upload connector frontend application that provides a UI for uploading files and sending them to the upload-backend.

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Source Code

* <https://github.com/RADAR-base/radar-upload-source-connector>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` | Number of radar-upload-connect-frontend replicas to deploy |
| image.repository | string | `"radarbase/radar-upload-connect-frontend"` | radar-upload-connect-frontend image repository |
| image.tag | string | `"0.5.9"` | radar-upload-connect-frontend image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | radar-upload-connect-frontend image pull policy |
| imagePullSecrets | list | `[]` | Docker registry secret names as an array |
| nameOverride | string | `""` | String to partially override radar-upload-connect-frontend.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override radar-upload-connect-frontend.fullname template with a string |
| podSecurityContext | object | `{}` | Configure radar-upload-connect-frontend pods' Security Context |
| securityContext | object | `{}` | Configure radar-upload-connect-frontend containers' Security Context |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `80` | radar-upload-connect-frontend port |
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.annotations | object | check values.yaml | Annotations that define default ingress class, certificate issuer |
| ingress.path | string | `"/upload/?(.*)"` | Path within the url structure |
| ingress.hosts | list | `["localhost"]` | Host to listen to requests to |
| ingress.tls.secretName | string | `"radar-base-tls"` | Name of the secret containing TLS certificates |
| resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| server_name | string | `"localhost"` | Server name or domain name |
| vue_app_client_id | string | `"radar_upload_frontend"` | OAuth2 client id of the upload connect frontend application |

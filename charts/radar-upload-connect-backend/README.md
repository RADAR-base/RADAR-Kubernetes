

# radar-upload-connect-backend

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.9](https://img.shields.io/badge/AppVersion-0.5.9-informational?style=flat-square)

A Helm chart for RADAR-base upload connector backend application.

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
| replicaCount | int | `2` | Number of radar-upload-connect-backend replicas to deploy |
| image.repository | string | `"radarbase/radar-upload-connect-backend"` | radar-upload-connect-backend image repository |
| image.tag | string | `"0.5.9"` | radar-upload-connect-backend image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | radar-upload-connect-backend image pull policy |
| nameOverride | string | `""` | String to partially override radar-upload-connect-backend.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override radar-upload-connect-backend.fullname template with a string |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `8085` | radar-upload-connect-backend port |
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.annotations | object | check values.yaml | Annotations that define default ingress class, certificate issuer and proxy settings |
| ingress.path | string | `"/upload/api/?(.*)"` | Path within the url structure |
| ingress.hosts | list | `["localhost"]` | Host to listen to requests to |
| ingress.tls.secretName | string | `"radar-base-tls"` | Name of the secret containing TLS certificates |
| resources.requests | object | `{"cpu":"100m","memory":"2Gi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| client_id | string | `"radar_upload_backend"` | OAuth2 client id of the upload connect backend application |
| client_secret | string | `"secret"` | OAuth2 client secret of the upload connect backend |
| postgres.host | string | `"radar-upload-postgresql-postgresql"` | Host name of the database to store uploaded data and metadata |
| postgres.user | string | `"postgres"` | Database username |
| postgres.password | string | `"password"` | Database password |
| managementportal_host | string | `"management-portal"` | Host name of the management portal application |
| serverName | string | `"localhost"` | Server name or domain name |

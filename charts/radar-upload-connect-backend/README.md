

# radar-upload-connect-backend

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.9](https://img.shields.io/badge/AppVersion-0.5.9-informational?style=flat-square)

A Helm chart for RADAR-base upload connector backend application.

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/radar-upload-source-connector>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+
* PV provisioner support in the underlying infrastructure

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` |  |
| image.repository | string | `"radarbase/radar-upload-connect-backend"` |  |
| image.tag | string | `"0.5.9"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| service.type | string | `"ClusterIP"` |  |
| service.port | int | `8085` |  |
| ingress.enabled | bool | `true` |  |
| ingress.annotations."kubernetes.io/ingress.class" | string | `"nginx"` |  |
| ingress.annotations."cert-manager.io/cluster-issuer" | string | `"letsencrypt-prod"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/rewrite-target" | string | `"/$1"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/proxy-body-size" | string | `"200m"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/proxy-request-buffering" | string | `"off"` |  |
| ingress.path | string | `"/upload/api/?(.*)"` |  |
| ingress.hosts[0] | string | `"localhost"` |  |
| ingress.tls.secretName | string | `"radar-base-tls"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"2Gi"` |  |
| persistence.enabled | bool | `true` |  |
| persistence.existingClaim | string | `"radar-output"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
| client_id | string | `"radar_upload_backend"` | OAuth2 client id of the upload connect backend application |
| client_secret | string | `"secret"` | OAuth2 client secret of the upload connect backend |
| postgres.host | string | `"radar-upload-postgresql-postgresql"` | Host name of the database to store uploaded data and metadata |
| postgres.user | string | `"postgres"` | Database username |
| postgres.password | string | `"password"` | Database password |
| managementportal_host | string | `"management-portal"` | Host name of the management portal application |
| serverName | string | `"localhost"` | Server name or domain name |


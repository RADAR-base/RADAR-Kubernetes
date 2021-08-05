
* PV provisioner support in the underlying infrastructure

# management-portal

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.6.4](https://img.shields.io/badge/AppVersion-0.6.4-informational?style=flat-square)

A Helm chart for RADAR-Base Management Portal

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/ManagementPortal>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` | Number of Management Portal replicas to deploy |
| image.repository | string | `"radarbase/management-portal"` | Management Portal image repository |
| image.tag | string | `"0.7.0"` | Management Portal image tag (immutable tags are recommended) |
| image.pullPolicy | string | `"IfNotPresent"` | Management Portal image pull policy |
| nameOverride | string | `""` | String to partially override management-portal.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override management-portal.fullname template with a string |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `8080` | Management Portal port |
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.annotations | object | check values.yaml | Annotations that define default ingress class, certificate issuer and rate limiter |
| ingress.path | string | `"/managementportal"` | Path within the url structure |
| ingress.hosts | list | `["localhost"]` | Hosts to accept requests from |
| ingress.tls.secretName | string | `"radar-base-tls"` | TLS Secret Name |
| resources.limits | object | `{"cpu":2,"memory":"1.7Gi"}` | CPU/Memory resource limits |
| resources.requests | object | `{"cpu":"100m","memory":"512Mi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| postgres.host | string | `"postgresql"` | host name of the postgres db |
| postgres.port | int | `5432` | post of the postgres db |
| postgres.database | string | `"managementportal"` | database name |
| postgres.connection_parameters | string | `""` | additional JDBC connection parameters e.g. sslmode=verify-full |
| postgres.user | string | `"postgres"` | postgres user |
| postgres.password | string | `"password"` | password of the postgres user |
| postgres.ssl.enabled | bool | `false` | set to true of the connecting to postgres using SSL |
| postgres.ssl.keystorepassword | string | `"keystorepassword"` |  |
| server_name | string | `"localhost"` | domain name of the server |
| catalogue_server | string | `"catalog-server"` | Hostname of the catalogue-server |
| managementportal.catalogue_server_enable_auto_import | string | `"false"` | set to true, if automatic source-type import from catalogue server should be enabled |
| managementportal.common_privacy_policy_url | string | `"http://info.thehyve.nl/radar-cns-privacy-policy"` | Override with a publicly resolvable url of the privacy-policy url for your set-up. This can be overridden on a project basis as well. |
| managementportal.oauth_checking_key_aliases_0 | string | `"radarbase-managementportal-ec"` | Keystore alias to sign JWT tokens from Management Portal |
| managementportal.oauth_checking_key_aliases_1 | string | `"selfsigned"` | Keystore alias to sign JWT tokens from Management Portal |
| managementportal.frontend_client_secret | string | `"xxx"` | OAuth Client secret of the Management Portal frontend application |
| managementportal.common_admin_password | string | `"xxx"` | Admin password of the default admin user created by the system |
| smtp.enabled | bool | `false` | set to true, if SMTP server should be enabled. Required to be true for production setup |
| smtp.host | string | `"smtp"` | Hostname of the SMTP server |
| smtp.port | int | `25` | Port of the SMTP server |
| smtp.username | string | `"username"` | Username of the SMTP server |
| smtp.password | string | `"secret"` | Password of the SMTP server |
| smtp.from | string | `"noreply@example.com"` | Email address which should be used to send activation emails |
| smtp.starttls | bool | `false` | set to true,if ttls should be enabled |
| smtp.auth | bool | `true` | set to true, if the account should be authenticated before sending emails |
| oauth_clients.pRMT | object | check values.yaml | Oauth Client configuration for pRMT |
| oauth_clients.aRMT | object | check values.yaml | Oauth Client configuration for aRMT |
| oauth_clients.THINC-IT | object | check values.yaml | Oauth Client configuration for THINC-IT |
| oauth_clients.radar_redcap_integrator | object | check values.yaml | Oauth Client configuration for REDCAP integrator |
| oauth_clients.radar_upload_backend | object | check values.yaml | Oauth Client configuration for Upload backend |
| oauth_clients.radar_upload_connect | object | check values.yaml | Oauth Client configuration for Upload connect |
| oauth_clients.radar_upload_frontend | object | check values.yaml | Oauth Client configuration for Upload frontend |
| oauth_clients.radar_rest_sources_auth_backend | object | check values.yaml | Oauth Client configuration for Rest Sources Backend |
| oauth_clients.radar_rest_sources_authorizer | object | check values.yaml | Oauth Client configuration for Rest Sources Authorizer |
| oauth_clients.radar_fitbit_connector | object | check values.yaml | Oauth Client configuration for Fitbit connector |
| oauth_clients.radar_appconfig | object | check values.yaml | Oauth Client configuration for Appconfig |
| oauth_clients.appconfig_frontend | object | check values.yaml | Oauth Client configuration for Appconfig Frontend |
| oauth_clients.grafana_dashboard | object | check values.yaml | Oauth Client configuration for Grafana Dashboard |

## OAuth Client Configuration
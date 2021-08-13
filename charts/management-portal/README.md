

# management-portal

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.6.4](https://img.shields.io/badge/AppVersion-0.6.4-informational?style=flat-square)

A Helm chart for RADAR-Base Management Portal to manage projects and participants throughout RADAR-base.

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

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
| imagePullSecrets | list | `[]` | Docker registry secret names as an array |
| nameOverride | string | `""` | String to partially override management-portal.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override management-portal.fullname template with a string |
| podSecurityContext | object | `{}` | Configure management-portal pods' Security Context |
| securityContext | object | `{}` | Configure management-portal containers' Security Context |
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
| managementportal.frontend_client_secret | string | `"xxx"` | OAuth2 Client secret of the Management Portal frontend application |
| managementportal.common_admin_password | string | `"xxx"` | Admin password of the default admin user created by the system |
| smtp.enabled | bool | `false` | set to true, if SMTP server should be enabled. Required to be true for production setup |
| smtp.host | string | `"smtp"` | Hostname of the SMTP server |
| smtp.port | int | `25` | Port of the SMTP server |
| smtp.username | string | `"username"` | Username of the SMTP server |
| smtp.password | string | `"secret"` | Password of the SMTP server |
| smtp.from | string | `"noreply@example.com"` | Email address which should be used to send activation emails |
| smtp.starttls | bool | `false` | set to true,if ttls should be enabled |
| smtp.auth | bool | `true` | set to true, if the account should be authenticated before sending emails |
| oauth_clients | object | check values.yaml | OAuth2 Client configuration |

## Requirements

Kubernetes: `<=1.17`

## OAuth Client Configuration
List of OAuth client configurations supported by RADAR-base. Each client should be enabled separately, if relevant and used in the installation.
Each client configuration has the following setup:
```
<client_id>: # client id
  enabled: false # set to true, if it should be enabled. Default is false.
  resource_ids: # list of resources that can be accessed by this client.
  client_secret: # Client secret. For public clients, the secret can be empty.
  scope: # List of permissions allowed for this client
  authorized_grant_types: # List of OAuth2 grant types that can be used by this client. e.g. authorization_code, refresh_token,client_credentials
  access_token_validity: # Lifetime of the access token for this client in seconds
  refresh_token_validity: # Lifetime of the refresh token for this client in seconds
  additional_information: # A JSON string containing additional meta-data of this client. e.g. {"dynamic_registration": true} should be set for clients which can automatically register a data source for a subject
  redirect_uri: # Redirect URL for clients which have authorization_code grant-type enabled.
  autoapprove: # List of permissions that can auto-approved when authorization-code flow succeeds.
```

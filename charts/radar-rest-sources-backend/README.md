

# radar-rest-sources-backend

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 3.2.0](https://img.shields.io/badge/AppVersion-3.2.0-informational?style=flat-square)

A Helm chart for the backend application of RADAR-base Rest Sources Authorizer

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
* PV provisioner support in the underlying infrastructure

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` |  |
| image.repository | string | `"radarbase/radar-rest-source-auth-backend"` |  |
| image.tag | string | `"3.2.0"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| service.type | string | `"ClusterIP"` |  |
| service.port | int | `8080` |  |
| ingress.enabled | bool | `true` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/affinity" | string | `"cookie"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/affinity-mode" | string | `"persistent"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/session-cookie-path" | string | `"/rest-sources/"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/session-cookie-samesite" | string | `"Strict"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/session-cookie-max-age" | string | `"900"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/session-cookie-expires" | string | `"900"` |  |
| ingress.annotations."kubernetes.io/ingress.class" | string | `"nginx"` |  |
| ingress.annotations."cert-manager.io/cluster-issuer" | string | `"letsencrypt-prod"` |  |
| ingress.path | string | `"/rest-sources/backend/?(.*)"` |  |
| ingress.hosts[0] | string | `"localhost"` |  |
| ingress.tls.secretName | string | `"radar-base-tls"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"400Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
| postgres.host | string | `"postgresql"` | host name of the postgres db |
| postgres.port | int | `5432` | post of the postgres db |
| postgres.database | string | `"restsourceauthorizer"` | database name |
| postgres.connection_parameters | string | `""` | additional JDBC connection parameters e.g. sslmode=verify-full |
| postgres.user | string | `"postgres"` | postgres user |
| postgres.password | string | `"password"` | password of the postgres user |
| postgres.ssl.enabled | bool | `false` | set to true of the connecting to postgres using SSL |
| postgres.ssl.keystorepassword | string | `"keystorepassword"` |  |
| managementportal_host | string | `"management-portal"` | hostname of the Management Portal |
| client_secret | string | `"secret"` | OAuth client secret of the radar-rest-sources-backend client from Management Portal |
| restSourceClients.fitbit.enable | bool | `false` | set to true, if Fitbit client should be used |
| restSourceClients.fitbit.sourceType | string | `"FitBit"` | Type of the data sources |
| restSourceClients.fitbit.authorizationEndpoint | string | `"https://www.fitbit.com/oauth2/authorize"` | Authorization endpoint for Fitbit authentication and authorization |
| restSourceClients.fitbit.tokenEndpoint | string | `"https://api.fitbit.com/oauth2/token"` | Token endpoint to request access-token from FitBit |
| restSourceClients.fitbit.clientId | string | `nil` | FitBit client id |
| restSourceClients.fitbit.clientSecret | string | `nil` | FitBit client secret |
| restSourceClients.fitbit.scope | string | `"activity heartrate sleep profile"` | List of scopes of the data that should be collected from Fitbit. For details, please refer to https://dev.fitbit.com/build/reference/web-api/developer-guide/application-design/#Scopes |


# management-portal

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.6.4](https://img.shields.io/badge/AppVersion-0.6.4-informational?style=flat-square)

A Helm chart for RADAR-Base Management Portal

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |

## Source Code

* <https://github.com/RADAR-base/ManagementPortal>

## Prerequisites
* Kubernetes 1.17+
* Helm 3.1.0
* PV provisioner support in the underlying infrastructure

## Requirements

Kubernetes: `<=1.17`

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` |  |
| image.repository | string | `"radarbase/management-portal"` | Management Portal image repository |
| image.tag | string | `"0.7.0"` | Management Portal image tag (immutable tags are recommended) |
| image.pullPolicy | string | `"IfNotPresent"` | Management Portal image pull policy |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| service.type | string | `"ClusterIP"` |  |
| service.port | int | `8080` |  |
| ingress.enabled | bool | `true` |  |
| ingress.annotations."kubernetes.io/ingress.class" | string | `"nginx"` |  |
| ingress.annotations."cert-manager.io/cluster-issuer" | string | `"letsencrypt-prod"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/server-snippet" | string | `"location /managementportal/oauth/ {\n  # Allow 20 fast-following requests, like when authorizing a user.\n  limit_req zone=login_limit burst=20;\n}\nlocation /managementportal/api/meta-token/ {\n  limit_req zone=login_limit;\n}\n"` |  |
| ingress.path | string | `"/managementportal"` |  |
| ingress.hosts[0] | string | `"localhost"` |  |
| ingress.tls.secretName | string | `"radar-base-tls"` |  |
| resources.limits.cpu | int | `2` |  |
| resources.limits.memory | string | `"1.7Gi"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"512Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
| postgres.host | string | `"postgresql"` |  |
| postgres.port | int | `5432` |  |
| postgres.database | string | `"managementportal"` |  |
| postgres.connection_parameters | string | `""` |  |
| postgres.user | string | `"postgres"` |  |
| postgres.password | string | `"password"` |  |
| postgres.ssl.enabled | bool | `false` |  |
| postgres.ssl.keystorepassword | string | `"keystorepassword"` |  |
| server_name | string | `"localhost"` |  |
| catalogue_server | string | `"catalog-server"` |  |
| managementportal.catalogue_server_enable_auto_import | string | `"false"` |  |
| managementportal.common_privacy_policy_url | string | `"http://info.thehyve.nl/radar-cns-privacy-policy"` |  |
| managementportal.oauth_checking_key_aliases_0 | string | `"radarbase-managementportal-ec"` |  |
| managementportal.oauth_checking_key_aliases_1 | string | `"selfsigned"` |  |
| managementportal.frontend_client_secret | string | `"xxx"` |  |
| managementportal.common_admin_password | string | `"xxx"` |  |
| oauth_clients.pRMT.enable | bool | `false` |  |
| oauth_clients.pRMT.resource_ids[0] | string | `"res_gateway"` |  |
| oauth_clients.pRMT.resource_ids[1] | string | `"res_ManagementPortal"` |  |
| oauth_clients.pRMT.resource_ids[2] | string | `"res_appconfig"` |  |
| oauth_clients.pRMT.client_secret | string | `""` |  |
| oauth_clients.pRMT.scope[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.pRMT.scope[1] | string | `"PROJECT.READ"` |  |
| oauth_clients.pRMT.scope[2] | string | `"ROLE.READ"` |  |
| oauth_clients.pRMT.scope[3] | string | `"SOURCE.READ"` |  |
| oauth_clients.pRMT.scope[4] | string | `"SOURCEDATA.READ"` |  |
| oauth_clients.pRMT.scope[5] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.pRMT.scope[6] | string | `"SUBJECT.READ"` |  |
| oauth_clients.pRMT.scope[7] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.pRMT.scope[8] | string | `"USER.READ"` |  |
| oauth_clients.pRMT.authorized_grant_types[0] | string | `"refresh_token"` |  |
| oauth_clients.pRMT.authorized_grant_types[1] | string | `"authorization_code"` |  |
| oauth_clients.pRMT.access_token_validity | int | `43200` |  |
| oauth_clients.pRMT.refresh_token_validity | int | `7948800` |  |
| oauth_clients.pRMT.additional_information | string | `"{\"dynamic_registration\": true}"` |  |
| oauth_clients.aRMT.enable | bool | `false` |  |
| oauth_clients.aRMT.resource_ids[0] | string | `"res_gateway"` |  |
| oauth_clients.aRMT.resource_ids[1] | string | `"res_ManagementPortal"` |  |
| oauth_clients.aRMT.resource_ids[2] | string | `"res_appconfig"` |  |
| oauth_clients.aRMT.client_secret | string | `""` |  |
| oauth_clients.aRMT.scope[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.aRMT.scope[1] | string | `"PROJECT.READ"` |  |
| oauth_clients.aRMT.scope[2] | string | `"ROLE.READ"` |  |
| oauth_clients.aRMT.scope[3] | string | `"SOURCE.READ"` |  |
| oauth_clients.aRMT.scope[4] | string | `"SOURCEDATA.READ"` |  |
| oauth_clients.aRMT.scope[5] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.aRMT.scope[6] | string | `"SUBJECT.READ"` |  |
| oauth_clients.aRMT.scope[7] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.aRMT.scope[8] | string | `"USER.READ"` |  |
| oauth_clients.aRMT.authorized_grant_types[0] | string | `"refresh_token"` |  |
| oauth_clients.aRMT.authorized_grant_types[1] | string | `"authorization_code"` |  |
| oauth_clients.aRMT.access_token_validity | int | `43200` |  |
| oauth_clients.aRMT.refresh_token_validity | int | `7948800` |  |
| oauth_clients.aRMT.additional_information | string | `"{\"dynamic_registration\": true}"` |  |
| oauth_clients.THINC-IT.enable | bool | `false` |  |
| oauth_clients.THINC-IT.resource_ids[0] | string | `"res_gateway"` |  |
| oauth_clients.THINC-IT.resource_ids[1] | string | `"res_ManagementPortal"` |  |
| oauth_clients.THINC-IT.resource_ids[2] | string | `"res_appconfig"` |  |
| oauth_clients.THINC-IT.client_secret | string | `""` |  |
| oauth_clients.THINC-IT.scope[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.THINC-IT.scope[1] | string | `"PROJECT.READ"` |  |
| oauth_clients.THINC-IT.scope[2] | string | `"ROLE.READ"` |  |
| oauth_clients.THINC-IT.scope[3] | string | `"SOURCE.READ"` |  |
| oauth_clients.THINC-IT.scope[4] | string | `"SOURCEDATA.READ"` |  |
| oauth_clients.THINC-IT.scope[5] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.THINC-IT.scope[6] | string | `"SUBJECT.READ"` |  |
| oauth_clients.THINC-IT.scope[7] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.THINC-IT.scope[8] | string | `"USER.READ"` |  |
| oauth_clients.THINC-IT.authorized_grant_types[0] | string | `"refresh_token"` |  |
| oauth_clients.THINC-IT.authorized_grant_types[1] | string | `"authorization_code"` |  |
| oauth_clients.THINC-IT.access_token_validity | int | `43200` |  |
| oauth_clients.THINC-IT.refresh_token_validity | int | `7948800` |  |
| oauth_clients.THINC-IT.additional_information | string | `"{\"dynamic_registration\": true}"` |  |
| oauth_clients.radar_redcap_integrator.enable | bool | `false` |  |
| oauth_clients.radar_redcap_integrator.resource_ids[0] | string | `"res_ManagementPortal"` |  |
| oauth_clients.radar_redcap_integrator.client_secret | string | `""` |  |
| oauth_clients.radar_redcap_integrator.scope[0] | string | `"PROJECT.READ"` |  |
| oauth_clients.radar_redcap_integrator.scope[1] | string | `"SUBJECT.CREATE"` |  |
| oauth_clients.radar_redcap_integrator.scope[2] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_redcap_integrator.scope[3] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.radar_redcap_integrator.authorized_grant_types[0] | string | `"client_credentials"` |  |
| oauth_clients.radar_redcap_integrator.access_token_validity | int | `900` |  |
| oauth_clients.radar_upload_backend.enable | bool | `false` |  |
| oauth_clients.radar_upload_backend.resource_ids[0] | string | `"res_ManagementPortal"` |  |
| oauth_clients.radar_upload_backend.client_secret | string | `""` |  |
| oauth_clients.radar_upload_backend.scope[0] | string | `"PROJECT.READ"` |  |
| oauth_clients.radar_upload_backend.scope[1] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_upload_backend.authorized_grant_types[0] | string | `"client_credentials"` |  |
| oauth_clients.radar_upload_backend.access_token_validity | int | `900` |  |
| oauth_clients.radar_upload_backend.additional_information | string | `"{\"dynamic_registration\": true}"` |  |
| oauth_clients.radar_upload_connect.enable | bool | `false` |  |
| oauth_clients.radar_upload_connect.resource_ids[0] | string | `"res_ManagementPortal"` |  |
| oauth_clients.radar_upload_connect.resource_ids[1] | string | `"res_upload"` |  |
| oauth_clients.radar_upload_connect.client_secret | string | `""` |  |
| oauth_clients.radar_upload_connect.scope[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.radar_upload_connect.scope[1] | string | `"PROJECT.READ"` |  |
| oauth_clients.radar_upload_connect.scope[2] | string | `"SOURCE.READ"` |  |
| oauth_clients.radar_upload_connect.scope[3] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.radar_upload_connect.scope[4] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_upload_connect.scope[5] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.radar_upload_connect.authorized_grant_types[0] | string | `"client_credentials"` |  |
| oauth_clients.radar_upload_connect.access_token_validity | int | `900` |  |
| oauth_clients.radar_upload_frontend.enable | bool | `false` |  |
| oauth_clients.radar_upload_frontend.resource_ids[0] | string | `"res_ManagementPortal"` |  |
| oauth_clients.radar_upload_frontend.resource_ids[1] | string | `"res_upload"` |  |
| oauth_clients.radar_upload_frontend.client_secret | string | `""` |  |
| oauth_clients.radar_upload_frontend.scope[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.radar_upload_frontend.scope[1] | string | `"PROJECT.READ"` |  |
| oauth_clients.radar_upload_frontend.scope[2] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.radar_upload_frontend.scope[3] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_upload_frontend.authorized_grant_types[0] | string | `"authorization_code"` |  |
| oauth_clients.radar_upload_frontend.access_token_validity | int | `900` |  |
| oauth_clients.radar_upload_frontend.redirect_uri[0] | string | `"/upload/login"` |  |
| oauth_clients.radar_upload_frontend.redirect_uri[1] | string | `"http://localhost:8080/upload/login"` |  |
| oauth_clients.radar_rest_sources_auth_backend.enable | bool | `false` |  |
| oauth_clients.radar_rest_sources_auth_backend.resource_ids[0] | string | `"res_ManagementPortal"` |  |
| oauth_clients.radar_rest_sources_auth_backend.resource_ids[1] | string | `"res_upload"` |  |
| oauth_clients.radar_rest_sources_auth_backend.client_secret | string | `""` |  |
| oauth_clients.radar_rest_sources_auth_backend.scope[0] | string | `"PROJECT.READ"` |  |
| oauth_clients.radar_rest_sources_auth_backend.scope[1] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_rest_sources_auth_backend.authorized_grant_types[0] | string | `"client_credentials"` |  |
| oauth_clients.radar_rest_sources_auth_backend.access_token_validity | int | `900` |  |
| oauth_clients.radar_rest_sources_authorizer.enable | bool | `false` |  |
| oauth_clients.radar_rest_sources_authorizer.resource_ids[0] | string | `"res_restAuthorizer"` |  |
| oauth_clients.radar_rest_sources_authorizer.client_secret | string | `""` |  |
| oauth_clients.radar_rest_sources_authorizer.scope[0] | string | `"PROJECT.READ"` |  |
| oauth_clients.radar_rest_sources_authorizer.scope[1] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.radar_rest_sources_authorizer.scope[2] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_rest_sources_authorizer.scope[3] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.radar_rest_sources_authorizer.authorized_grant_types[0] | string | `"authorization_code"` |  |
| oauth_clients.radar_rest_sources_authorizer.access_token_validity | int | `900` |  |
| oauth_clients.radar_rest_sources_authorizer.redirect_uri[0] | string | `"/rest-sources/authorizer/login"` |  |
| oauth_clients.radar_fitbit_connector.enable | bool | `false` |  |
| oauth_clients.radar_fitbit_connector.resource_ids[0] | string | `"res_restAuthorizer"` |  |
| oauth_clients.radar_fitbit_connector.client_secret | string | `""` |  |
| oauth_clients.radar_fitbit_connector.scope[0] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_fitbit_connector.scope[1] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.radar_fitbit_connector.authorized_grant_types[0] | string | `"client_credentials"` |  |
| oauth_clients.radar_fitbit_connector.access_token_validity | int | `900` |  |
| oauth_clients.radar_appconfig.enable | bool | `false` |  |
| oauth_clients.radar_appconfig.resource_ids[0] | string | `"res_ManagementPortal"` |  |
| oauth_clients.radar_appconfig.resource_ids[1] | string | `"res_appconfig"` |  |
| oauth_clients.radar_appconfig.client_secret | string | `""` |  |
| oauth_clients.radar_appconfig.scope[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.radar_appconfig.scope[1] | string | `"OAUTHCLIENTS.READ"` |  |
| oauth_clients.radar_appconfig.scope[2] | string | `"PROJECT.READ"` |  |
| oauth_clients.radar_appconfig.scope[3] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.radar_appconfig.scope[4] | string | `"SUBJECT.READ"` |  |
| oauth_clients.radar_appconfig.authorized_grant_types[0] | string | `"client_credentials"` |  |
| oauth_clients.radar_appconfig.access_token_validity | int | `900` |  |
| oauth_clients.appconfig_frontend.enable | bool | `false` |  |
| oauth_clients.appconfig_frontend.resource_ids[0] | string | `"res_appconfig"` |  |
| oauth_clients.appconfig_frontend.client_secret | string | `""` |  |
| oauth_clients.appconfig_frontend.scope[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.appconfig_frontend.scope[1] | string | `"OAUTHCLIENTS.READ"` |  |
| oauth_clients.appconfig_frontend.scope[2] | string | `"PROJECT.CREATE"` |  |
| oauth_clients.appconfig_frontend.scope[3] | string | `"PROJECT.READ"` |  |
| oauth_clients.appconfig_frontend.scope[4] | string | `"PROJECT.UPDATE"` |  |
| oauth_clients.appconfig_frontend.scope[5] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.appconfig_frontend.scope[6] | string | `"SUBJECT.READ"` |  |
| oauth_clients.appconfig_frontend.scope[7] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.appconfig_frontend.authorized_grant_types[0] | string | `"authorization_code"` |  |
| oauth_clients.appconfig_frontend.authorized_grant_types[1] | string | `"refresh_token"` |  |
| oauth_clients.appconfig_frontend.access_token_validity | int | `900` |  |
| oauth_clients.appconfig_frontend.refresh_token_validity | int | `78000` |  |
| oauth_clients.appconfig_frontend.redirect_uri[0] | string | `"/appconfig/login"` |  |
| oauth_clients.appconfig_frontend.autoapprove[0] | string | `"MEASUREMENT.CREATE"` |  |
| oauth_clients.appconfig_frontend.autoapprove[1] | string | `"OAUTHCLIENTS.READ"` |  |
| oauth_clients.appconfig_frontend.autoapprove[2] | string | `"PROJECT.CREATE"` |  |
| oauth_clients.appconfig_frontend.autoapprove[3] | string | `"PROJECT.READ"` |  |
| oauth_clients.appconfig_frontend.autoapprove[4] | string | `"PROJECT.UPDATE"` |  |
| oauth_clients.appconfig_frontend.autoapprove[5] | string | `"SOURCETYPE.READ"` |  |
| oauth_clients.appconfig_frontend.autoapprove[6] | string | `"SUBJECT.READ"` |  |
| oauth_clients.appconfig_frontend.autoapprove[7] | string | `"SUBJECT.UPDATE"` |  |
| oauth_clients.grafana_dashboard.enable | bool | `false` |  |
| oauth_clients.grafana_dashboard.resource_ids[0] | string | `"res_ManagementPortal"` |  |
| oauth_clients.grafana_dashboard.client_secret | string | `""` |  |
| oauth_clients.grafana_dashboard.scope[0] | string | `"USER.READ"` |  |
| oauth_clients.grafana_dashboard.authorized_grant_types[0] | string | `"authorization_code"` |  |
| oauth_clients.grafana_dashboard.authorized_grant_types[1] | string | `"refresh_token"` |  |
| oauth_clients.grafana_dashboard.access_token_validity | int | `900` |  |
| oauth_clients.grafana_dashboard.refresh_token_validity | int | `78000` |  |
| oauth_clients.grafana_dashboard.redirect_uri[0] | string | `"http://dashboard.localhost/login/generic_oauth"` |  |
| oauth_clients.grafana_dashboard.autoapprove[0] | string | `"USER.READ"` |  |
| smtp.enabled | bool | `false` |  |
| smtp.host | string | `"smtp"` |  |
| smtp.port | int | `25` |  |
| smtp.username | string | `"username"` |  |
| smtp.password | string | `"secret"` |  |
| smtp.from | string | `"noreply@example.com"` |  |
| smtp.starttls | bool | `false` |  |
| smtp.auth | bool | `true` |  |

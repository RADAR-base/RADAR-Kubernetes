
* PV provisioner support in the underlying infrastructure
# app-config

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![AppVersion: 0.3.2](https://img.shields.io/badge/AppVersion-0.3.2-informational?style=flat-square)

A Helm chart for Kubernetes

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` | Number of Appconfig replicas to deploy |
| image.repository | string | `"radarbase/radar-app-config"` | Appconfig image repository |
| image.tag | string | `"0.3.2"` | Appconfig image tag (immutable tags are recommended) |
| image.pullPolicy | string | `"IfNotPresent"` | Appconfig image pull policy |
| imagePullSecrets | list | `[]` | Docker registry secret names as an array |
| nameOverride | string | `""` | String to partially override management-portal.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override management-portal.fullname template with a string |
| namespace | string | `"default"` | Kubernetes namespace that Appconfig is going to be deployed on |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `nil` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| podSecurityContext | object | `{}` | Configure Appconfig pods' Security Context |
| securityContext | object | `{}` | Configure Appconfig containers' Security Context |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `8090` | Appconfig port |
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.annotations | object | check values.yaml | Annotations that define default ingress class, certificate issuer |
| ingress.path | string | `"/appconfig/api($|/)(.*)"` | Path within the url structure |
| ingress.hosts | list | `["localhost"]` | Hosts to accept requests from |
| ingress.tls.secretName | string | `"radar-base-tls"` | TLS Secret Name |
| resources.limits | object | `{"cpu":2}` | CPU/Memory resource limits |
| resources.requests | object | `{"cpu":"100m","memory":"768Mi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| javaOpts | string | `"-Xmx550m"` |  |
| clientId | string | `"radar_appconfig"` |  |
| clientSecret | string | `"secret"` |  |
| jdbc.driver | string | `"org.postgresql.Driver"` |  |
| jdbc.url | string | `"jdbc:postgresql://postgresql:5432/appconfig"` |  |
| jdbc.user | string | `"postgres"` |  |
| jdbc.password | string | `"password"` |  |
| jdbc.dialect | string | `"org.hibernate.dialect.PostgreSQLDialect"` |  |

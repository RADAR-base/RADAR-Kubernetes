# data-dashboard-backend

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

API for data in the data dashboard

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | <keyvan@thehyve.nl> | <https://www.thehyve.nl> |
| Joris Borgdorff | <joris@thehyve.nl> | <https://www.thehyve.nl/experts/joris-borgdorff> |
| Nivethika Mahasivam | <nivethika@thehyve.nl> | <https://www.thehyve.nl/experts/nivethika-mahasivam> |

## Source Code

* <https://github.com/IMI-H2O/outcomes-dashboard>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default                                            | Description |
|-----|------|----------------------------------------------------|-------------|
| replicaCount | int | `2`                                                | Number of replicas to deploy |
| image.repository | string | `"quay.io/imi-h2o/data-dashboard-backend"`         | docker image repository |
| image.pullPolicy | string | `"Always"`                                         | image pull policy |
| image.tag | string | `"dev"`                                            |  |
| imagePullSecrets | list | `[]`                                               | Docker registry secret names as an array |
| nameOverride | string | `""`                                               | String to partially override fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""`                                               | String to fully override fullname template with a string |
| podSecurityContext | object | `{}`                                               | Configure pod's Security Context |
| securityContext | object | `{}`                                               | Configure container's Security Context |
| service.type | string | `"ClusterIP"`                                      | Kubernetes Service type |
| service.port | int | `9000`                                             | data-dashboard-backend port |
| ingress.enabled | bool | `true`                                             | Enable ingress controller resource |
| ingress.className | string | `""`                                               | Ingress class name |
| ingress.annotations | object | check values.yaml                                  | Annotations that define default ingress class, certificate issuer |
| ingress.path | string | `"/api($                                           |/)(.*)"` | Path within the url structure |
| ingress.pathType | string | `"ImplementationSpecific"`                         |  |
| ingress.hosts | list | `["localhost"]`                                    | Hosts to accept requests from |
| ingress.tls.secretName | string | `"radar-base-tls"`                                 | TLS Secret Name |
| resources | object | `{}`                                               |  |
| autoscaling.enabled | bool | `false`                                            | Enable horizontal autoscaling |
| autoscaling.minReplicas | int | `1`                                                |  |
| autoscaling.maxReplicas | int | `100`                                              |  |
| autoscaling.targetCPUUtilizationPercentage | int | `80`                                               |  |
| nodeSelector | object | `{}`                                               | Node labels for pod assignment |
| tolerations | list | `[]`                                               | Toleration labels for pod assignment |
| affinity | object | `{}`                                               | Affinity labels for pod assignment |
| existingSecret | string | `""`                                               |  |
| javaOpts | string | `"-Xmx550m"`                                       | Standard JAVA_OPTS that should be passed to this service |
| managementPortal.url | string | `"http://management-portal:8080/managementportal"` | ManagementPortal URL |
| managementPortal.clientId | string | `"data_dashboard_backend"`                         | ManagementPortal OAuth 2.0 client ID, having grant type client_credentials |
| managementPortal.clientSecret | string | `"secret"`                                         | ManagementPortal OAuth 2.0 client secret |
| path | string | `"/api"`                                           |  |
| score | object | Diabetes score calculation table configurations    | Where to find lookup tables for score calculations, and how to calculate them. |
| jdbc.driver | string | `"org.postgresql.Driver"`                          | JDBC Driver to connect to the database. |
| jdbc.url | string | `"jdbc:postgresql://postgresql:5432/outcomes"`     | JDBC Connection url of the database. |
| jdbc.user | string | `"radarbase"`                                      | Username of the database |
| jdbc.password | string | `"password"`                                       | Password of the user |
| jdbc.dialect | string | `"org.hibernate.dialect.PostgreSQLDialect"`        | Hibernate dialect to use for JDBC Connection |

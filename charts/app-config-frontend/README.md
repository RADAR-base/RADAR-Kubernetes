

# app-config-frontend

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.3.2](https://img.shields.io/badge/AppVersion-0.3.2-informational?style=flat-square)

A Helm chart for the frontend application of RADAR-base application config (app-config).

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Source Code

* <https://github.com/RADAR-base/radar-app-config>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `2` | Number of Appconfig frontend replicas to deploy |
| image.repository | string | `"radarbase/radar-app-config-frontend"` | Appconfig frontend image repository |
| image.tag | string | `"0.3.2"` | Appconfig frontend image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | Appconfig frontend image pull policy |
| imagePullSecrets | list | `[]` | Docker registry secret names as an array |
| nameOverride | string | `""` | String to partially override app-config-frontend.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override app-config-frontend.fullname template with a string |
| podAnnotations | object | `{}` | Annotations for Appconfig frontend pods |
| podSecurityContext | object | `{}` | Configure Appconfig pods' Security Context |
| securityContext | object | `{}` | Configure Appconfig containers' Security Context |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `8080` | Appconfig frontend port |
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.annotations | object | check values.yaml | Annotations that define default ingress class, certificate issuer |
| ingress.hosts[0] | object | `{"paths":["/appconfig($|/)(.*)"]}` | Path within the url structure |
| ingress.tls | list | `[]` | Utilize TLS backend in ingress |
| resources.limits | object | `{"cpu":"200m","memory":"512Mi"}` | CPU/Memory resource limits |
| resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| authUrl | string | `"http://localhost/managementportal/oauth"` | Authorization URL of the IDP |
| authCallbackUrl | string | `"http://localhost/appconfig/login"` | Callback URL to where authorization-code should be returned |
| backendUrl | string | `"/appconfig/api"` | Base-URL of the App Config backend service |

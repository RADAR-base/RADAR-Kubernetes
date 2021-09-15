

# s3-proxy

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for S3 Proxy. It uses https://hub.docker.com/r/andrewgaul/s3proxy to proxy S3 API requests to any supported cloud provider. For more examples see Find some example configurations at https://github.com/gaul/s3proxy/wiki/Storage-backend-examples.

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Source Code

* <https://github.com/gaul/s3proxy>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` | Number of s3-proxy replicas to deploy |
| image.repository | string | `"andrewgaul/s3proxy"` | s3-proxy image repository |
| image.tag | string | `"travis-1430"` | s3-proxy image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy | string | `"IfNotPresent"` | s3-proxy image pull policy |
| imagePullSecrets | list | `[]` | Docker registry secret names as an array |
| nameOverride | string | `""` | String to partially override s3-proxy.fullname template with a string (will prepend the release name) |
| fullnameOverride | string | `""` | String to fully override s3-proxy.fullname template with a string |
| podSecurityContext | object | `{}` | Configure s3-proxy pods' Security Context |
| securityContext | object | `{}` | Configure s3-proxy containers' Security Context |
| service.type | string | `"ClusterIP"` | Kubernetes Service type |
| service.port | int | `80` | s3-proxy port |
| resources.requests | object | `{"cpu":"100m","memory":"128Mi"}` | CPU/Memory resource requests |
| nodeSelector | object | `{}` | Node labels for pod assignment |
| tolerations | list | `[]` | Toleration labels for pod assignment |
| affinity | object | `{}` | Affinity labels for pod assignment |
| s3.identity | string | `nil` | Credentials used to access this proxy |
| s3.credential | string | `""` | Credentials used to access this proxy |
| target | object | Check below | Where requests should be proxied to |
| target.provider | string | `nil` | Target provider |
| target.endpoint | string | `nil` | Target endpoint |
| target.identity | string | `nil` | Target identity |
| target.credential | string | `""` | Target credential |

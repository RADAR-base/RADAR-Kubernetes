

# s3-proxy

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for S3 Proxy. It uses https://hub.docker.com/r/andrewgaul/s3proxy to proxy S3 API requests to any supported cloud provider. For more examples see Find some example configurations at https://github.com/gaul/s3proxy/wiki/Storage-backend-examples.

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |
| Nivethika Mahasivam | nivethika@thehyve.nl |  |

## Source Code

* <https://github.com/gaul/s3proxy>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+
* PV provisioner support in the underlying infrastructure

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` |  |
| image.repository | string | `"andrewgaul/s3proxy"` |  |
| image.tag | string | `"travis-1430"` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| nameOverride | string | `""` |  |
| fullnameOverride | string | `""` |  |
| service.type | string | `"ClusterIP"` |  |
| service.port | int | `80` |  |
| ingress.enabled | bool | `false` |  |
| ingress.annotations | object | `{}` |  |
| ingress.tls | list | `[]` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"128Mi"` |  |
| nodeSelector | object | `{}` |  |
| tolerations | list | `[]` |  |
| affinity | object | `{}` |  |
| s3.identity | string | `nil` | Credentials used to access this proxy |
| s3.credential | string | `""` | Credentials used to access this proxy |
| target | object | `{"credential":"","endpoint":null,"identity":null,"provider":null}` | Where requests should be proxied to |
| target.provider | string | `nil` | Target provider |
| target.endpoint | string | `nil` | Target endpoint |
| target.identity | string | `nil` | Target identity |
| target.credential | string | `""` | Target credential |

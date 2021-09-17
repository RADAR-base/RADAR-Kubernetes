

# cc-schema-registry-proxy

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for Confluent Cloud schema registry proxy. This proxy service is used when RADAR-base platform is used with Confluent Cloud based schema registry. It forwards requests to schema registry with an additonal basic authentication header with Confluent Cloud schema registry credentials. This service will be enabled if `cc.enabled = true`.

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` | Number of cc-schema-registry-proxy replicas to deploy |
| service.type | string | `"ExternalName"` | Kubernetes Service type, |
| service.externalName | string | `"schema-registry-domain"` | Domain name used for pointing to actual schema registry instance |
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.path | string | `"/schema/?(.*)"` | Path within the url structure |
| ingress.hosts | list | `["localhost"]` | Hosts to listen requests from |
| ingress.tls.secretName | string | `"radar-base-tls"` | Name of the secret that contains TLS certificates |
| cc.schemaRegistryApiKey | string | `"srApiKey"` | Confluent cloud schema registry API key |
| cc.schemaRegistryApiSecret | string | `"srApiSecret"` | Confluent cloud schema registry API secret |



# cert-manager

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for cert-manager. This chart is an overly just to make sure `clusterissuer.yaml` is installed with the cluster and some default values. For more info refer to the cert-manager docs.

**Homepage:** <https://cert-manager.io>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl |  |
| Joris Borgdorff | joris@thehyve.nl |  |

## Source Code

* <https://github.com/jetstack/cert-manager>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Requirements

Kubernetes: `<=1.17`

| Repository | Name | Version |
|------------|------|---------|
| https://charts.jetstack.io | cert-manager | v1.1.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| server_name | string | `"localhost"` | Domain name of the server |
| maintainer_email | string | `"me@example.com"` | Email address of cluster maintainer |
| installCRDs | bool | `true` | Install CRDs that are needed by cert-manager |
| prometheus.enabled | bool | `true` | Enable Prometheus monitoring |
| prometheus.servicemonitor.enabled | bool | `true` | Enable Prometheus Operator ServiceMonitor monitoring |

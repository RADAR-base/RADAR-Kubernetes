

# minio

Minio object storage used for storing RADAR-Base intermediate and output data. Some values have changed to adapt needs of RADAR-Base for complete list of values visit the original chart.

**Homepage:** <https://min.io/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Source Code

* <https://github.com/minio/charts>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| ingress.enabled | bool | `true` | Enable ingress controller resource |
| ingress.hosts | list | `["s3.example.com"]` | Hosts to accept requests from |
| ingress.annotations | object | check values.yaml | Annotations that define default ingress class, certificate issuer |
| resources.requests | object | `{"cpu":"100m","memory":"400Mi"}` | CPU/Memory resource requests |
| replicas | int | `4` | Number of nodes, 4 is minimum for distributed mode |
| mode | string | `"distributed"` | MinIO server mode |
| persistence.size | string | `"20Gi"` | Size of persistent volume claim |
| accessKey | string | `"AKIAIOSFODNN7EXAMPLE"` | Default access key (5 to 20 characters) |
| secretKey | string | `"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"` | Default secret key (8 to 40 characters) |
| metrics.serviceMonitor | object | `{"enabled":true}` | Set this to true to create ServiceMonitor for Prometheus operator |
| buckets[0].name | string | `"radar-intermediate-storage"` | Bucket used for intermediate storage of RADAR-Base data |
| buckets[0].policy | string | `"none"` | Public access is disabled for this bucket |
| buckets[0].purge | bool | `false` | Don't purge the bucket if already exists |
| buckets[1].name | string | `"radar-output-storage"` | Bucket used for output storage of RADAR-Base data |
| buckets[1].policy | string | `"none"` | Public access is disabled for this bucket |
| buckets[1].purge | bool | `false` | Don't purge the bucket if already exists |

## Minio client
To create users and define access policies use [the guide on Minio website](https://docs.min.io/docs/minio-client-quickstart-guide.html) to install and configure Minio client.

## Access policy
Three JSON files in this directory `read-output.json`, `write-intermediate.json` and `write-output.json` are meant to allow read only access to buckets for users and provide write access to applications.

```
mc admin policy add <clustername> read-output read-output.json
mc admin policy add <clustername> write-intermediate write-intermediate.json
mc admin policy add <clustername> write-output write-output.json

# Assign a user to a policy
mc admin policy set <clustername> read-output user=example.user
```

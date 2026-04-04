

# radar-seaweedfs
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/radar-seaweedfs)](https://artifacthub.io/packages/helm/radar-base/radar-seaweedfs)

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 4.02](https://img.shields.io/badge/AppVersion-4.02-informational?style=flat-square)

Seaweedfs charts for Radar-base

**Homepage:** <https://radar-base.org>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Pim van Nierop | <pim@thehyve.nl> | <https://www.thehyve.nl/experts/pim-van-nierop> |

## Prerequisites
* Kubernetes 1.28+
* Kubectl 1.28+
* Helm 3.1.0+

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../../external/seaweedfs-operator | operator(seaweedfs-operator) | 0.1.9 |
| https://radar-base.github.io/radar-helm-charts | common | 2.x.x |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| serverName | string | `"localhost"` |  |
| operator.fullnameOverride | string | `"seaweedfs-operator"` |  |
| operator.webhook.enabled | bool | `false` |  |
| operator.enabled | bool | `true` |  |
| operator.serviceMonitor.enabled | bool | `true` |  |
| operator.grafanaDashboard.enabled | bool | `true` |  |
| image.registry | string | `"docker.io"` | Image registry |
| image.repository | string | `"chrislusf/seaweedfs"` | Image repository |
| image.tag | string | `nil` | Image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.digest | string | `""` | Image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| image.pullSecrets | list | `[]` | Optionally specify an array of imagePullSecrets. Secrets must be manually created in the namespace. e.g: pullSecrets:   - myRegistryKeySecretName  |
| master.replicas | int | `3` |  |
| master.volumeSizeLimitMB | string | `nil` | Max size of a single volume in MB. Default is ~30Gb. When full, new volumes will be created with this size. |
| master.resources | string | `nil` |  |
| volume.replicas | int | `2` |  |
| volume.storage.size | string | `"2Gi"` | Size of PVC for volume server. Note that PVC size must take into account replication factor (see 'replication'). e.g. 2Gi * 2 replicas in the server rack = 4Gi |
| volume.resources | string | `nil` |  |
| filer.replicas | int | `2` |  |
| filer.resources | string | `nil` |  |
| affinity | string | `nil` |  |
| replication | string | `"000"` | Replication factor for data volumes in format "XYZ" where X, Y and Z are the number of data replicates across data centers, racks and servers in a rack, respectively. X, Y and Z can be 0, 1 or 2 signifying the number of replicates at the corresponding level. For instance, "001" adds one replica on the servers in the same rack. Important: the RADAR-base volume servers are hard-coded to run in the same rack; the X and Y positions are not in effect, therefore. |
| createBuckets | list | `[]` | List of S3 buckets to create on the SeaweedFS filer. Buckets are created idempotently. If empty, no buckets will be created. |
| createBucketsWait | object | `{"bucketsApiTimeoutSeconds":120,"filerReadyTimeoutSeconds":120,"pollIntervalSeconds":2}` | Settings for the create-buckets Job waiting loops Configure timeouts (in seconds) for waiting on the filer and the buckets API. Defaults are 120s (2 minutes) each. Poll interval controls how often the checks run. |
| s3Config.identities.admin.actions[0] | string | `"Admin"` |  |
| s3Config.identities.admin.credentials.accessKey | string | `"secret"` |  |
| s3Config.identities.admin.credentials.secretKey | string | `"secret"` |  |
| s3Config.identities.user.actions[0] | string | `"Read"` |  |
| s3Config.identities.user.actions[1] | string | `"Write"` |  |
| s3Config.identities.user.actions[2] | string | `"List"` |  |
| s3Config.identities.user.credentials.accessKey | string | `"secret"` |  |
| s3Config.identities.user.credentials.secretKey | string | `"secret"` |  |
| networkpolicy | object | check `values.yaml` | Network policy defines who can access this application and who this applications has access to |

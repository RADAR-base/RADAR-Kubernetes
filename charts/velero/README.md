

# velero

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0](https://img.shields.io/badge/AppVersion-1.0-informational?style=flat-square)

A Helm chart for Velero, this chart is an overlay for Velero and adds some default values and a deployment to mirror the local object storage to a remote location.

**Homepage:** <https://velero.io>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Keyvan Hedayati | keyvan@thehyve.nl | https://www.thehyve.nl |
| Joris Borgdorff | joris@thehyve.nl | https://www.thehyve.nl/experts/joris-borgdorff |
| Nivethika Mahasivam | nivethika@thehyve.nl | https://www.thehyve.nl/experts/nivethika-mahasivam |

## Source Code

* <https://github.com/vmware-tanzu/helm-charts/tree/main/charts/velero>

## Prerequisites
* Kubernetes 1.17+
* Kubectl 1.17+
* Helm 3.1.0+
* S3-compatible object storage

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://vmware-tanzu.github.io/helm-charts | velero | 2.12.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| objectStorageBackupReplicaCount | int | `1` | Number of replicas for object storage backup pod, should be 1 |
| mc_image.repository | string | `"minio/mc"` | Object storage backup pod image repository |
| mc_image.tag | string | `"RELEASE.2020-09-03T00-08-28Z"` | Object storage backup pod image tag (immutable tags are recommended) |
| mc_image.pullPolicy | string | `"IfNotPresent"` | Object storage backup pod image pull policy |
| imagePullSecrets | list | `[]` | Docker registry secret names as an array |
| podSecurityContext | object | `{}` | Configure object storage backup pod pods' Security Context |
| securityContext | object | `{}` | Configure object storage backup pod containers' Security Context |
| local.address | string | `"minio.default:9000"` | Address of local object storage to backup data from |
| local.accessKey | string | `"accessKey"` | Access key of local object storage |
| local.secretKey | string | `"secretKey"` | Secret key of local object storage |
| local.intermediateBucketName | string | `"radar-intermediate-storage"` | Name of local intermediate data bucket |
| local.outputBucketName | string | `"radar-output-storage"` | Name of local output data bucket |
| backup.address | string | `"s3.example.com"` | Address of remote object storage to backup data to |
| backup.accessKey | string | `"accessKey"` | Access key of remote object storage |
| backup.secretKey | string | `"secretKey"` | Secret key of remote object storage |
| backup.intermediateBucketName | string | `"radar-intermediate-storage"` | Name of remote intermediate data bucket |
| backup.outputBucketName | string | `"radar-output-storage"` | Name of remote output data bucket |
| velero.initContainers | list | check values.yaml | Add plugins to enable using different storage systems, AWS plugin is needed to be able to push to S3-compatible object storages |
| velero.metrics.enabled | bool | `true` | Enable monitoring metrics to be collected |
| velero.metrics.serviceMonitor.enabled | bool | `true` | Enable prometheus-operator interface |
| velero.configuration.provider | string | `"aws"` | Cloud provider being used (e.g. aws, azure, gcp). |
| velero.configuration.backupStorageLocation | object | Check below | Parameters for the `default` BackupStorageLocation. See https://velero.io/docs/v1.0.0/api-types/backupstoragelocation/ |
| velero.configuration.backupStorageLocation.name | string | `"default"` | Cloud provider where backups should be stored. Usually should match `configuration.provider`. Required. |
| velero.configuration.backupStorageLocation.bucket | string | `"radar-base-backups"` | Bucket to store backups in. Required. |
| velero.configuration.backupStorageLocation.config | object | Check values.yaml | Additional provider-specific configuration. See link above for details of required/optional fields for your provider. |
| velero.credentials.secretContents.cloud | string | Check values.yaml | Check |
| velero.snapshotsEnabled | bool | `false` | Don't snapshot volumes where they're not supported |
| velero.deployRestic | bool | `true` | Deploy restic to backup Kubernetes volumes |
| velero.restic.podVolumePath | string | `"/var/lib/kubelet/pods"` | Path to find pod volumes |
| velero.restic.privileged | bool | `false` | Shouldn't need privilege to backup the volumes |
| velero.restic.priorityClassName | object | `{}` | Pod priority class name to use for the Restic daemonset. Optional. |
| velero.restic.resources | object | `{}` | Resource requests/limits to specify for the Restic daemonset deployment. Optional. |
| velero.restic.tolerations | list | `[]` | Tolerations to use for the Restic daemonset. Optional. |
| velero.restic.extraVolumes | list | `[]` | Extra volumes for the Restic daemonset. Optional. |
| velero.restic.extraVolumeMounts | list | `[]` | Extra volumeMounts for the Restic daemonset. Optional. |
| velero.restic.securityContext | object | `{}` | SecurityContext to use for the Velero deployment. Optional. Set fsGroup for `AWS IAM Roles for Service Accounts` see more informations at: https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html |
| velero.schedules.backup.schedule | string | `"0 3 * * *"` | Backup every day at 3:00 AM |
| velero.schedules.backup.template.ttl | string | `"240h"` | Keep backup for 10 days |
| velero.schedules.backup.template.includeClusterResources | bool | `true` | Backup cluster wide resources |
| velero.schedules.backup.template.snapshotVolumes | bool | `false` | Don't snapshot volumes where they're not supported |
| velero.schedules.backup.template.includedNamespaces | list | Check values.yaml | Namespaces to backup manifests and volumes from |

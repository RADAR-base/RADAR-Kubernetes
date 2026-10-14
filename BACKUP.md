# Backup
[Velero](https://velero.io/) have been integrated for backing up the Kubernetes cluster that is hosting the Radar-Base. The process backs up 3 resource types:
1. Kubernetes data and manifests: Resource definitions like `deployment`, `service` and other APIs.
2. Application data: The volumes that applications such a PostgreSQL and Redis use to store their data.
3. Radar-base output and intermediate data: The final outcome of the Radar-base stack after processing all of the incoming data.

The first two types are handled by Velero and the third type is handled by [Minio client CLI](https://docs.minio.io/docs/minio-client-complete-guide), `mc` command.


## Kubernetes data and manifests
In default configuration Velero will run at [3:00 AM everyday](https://github.com/RADAR-base/RADAR-Kubernetes/blob/dev/charts/velero/values.yaml#L94) and [backup all of the resource definitions in following namespaces](https://github.com/RADAR-base/RADAR-Kubernetes/blob/dev/charts/velero/values.yaml#L99):
 - cert-manager
 - default
 - graylog
 - kubernetes-dashboard
 - monitoring
 - velero

This process will also backup all of the [cluster wide definitions](https://github.com/RADAR-base/RADAR-Kubernetes/blob/dev/charts/velero/values.yaml#L97) such as `ClusterRole` and `ClusterRoleBinding`.

## Application data
In order to backup the data stored in persistent volumes Velero uses [Restic](https://restic.net/) to backup the volumes. In default configuration volumes for [following applications](https://github.com/RADAR-base/RADAR-Kubernetes/search?q=backup.velero.io%2Fbackup-volumes&type=) will be backed up:
- catalog-server
- postgresql
- radar-fitbit-connector
- radar-upload-postgresql
- redis
- smtp

In order to make sure that volume backups are atomic during the process `fsfreeze` will prevent applications to write any new data to file system.


## Radar-base output and intermediate data
Minio client CLI will run as a deployment on the cluster and [continuously streams](https://github.com/RADAR-base/RADAR-Kubernetes/blob/dev/charts/velero/templates/deployment.yaml#L42) the data from internal object storage to external backup location. At the moment this application [needs to be restarted](https://github.com/RADAR-base/RADAR-Kubernetes/issues/70) once in a while to keep working.

## Backup destination
Backup destination should be an Object storage like AWS S3, it should be configured separately for [Manifest and volume data backups](https://github.com/RADAR-base/RADAR-Kubernetes/blob/dev/charts/velero/values.yaml#L58) and [Radar-base output backup](https://github.com/RADAR-base/RADAR-Kubernetes/blob/dev/charts/velero/values.yaml#L20).

# Restore
In case the cluster goes down or there is a need to restore data from back up following steps should be followed:
- Remove the `schedule` section [from configuration](https://github.com/RADAR-base/RADAR-Kubernetes/blob/dev/charts/velero/values.yaml#L92) to disable new backup creation
- Install the Velero on new cluster with this command: `helmfile -f helmfile.d/99-velero.yaml sync`
- Follow the steps from 3 to 5 in [Disaster recovery page](https://velero.io/docs/v1.4/disaster-case/) of Velero documentation
- Use the following commands to restore Object storage


If everything works correctly *manifest* and *volume data* should be restored to the new cluster, if it fails restoration should be done manually.

Radar-base output and intermediate data needs to be manually restored.

## Manual restore

### Restore manifests
```
mc ls backup/radar-base-backups/backups/
```

```
mkdir backup
cd backup
mc cp backup/radar-base-backups/backups/velero-backup-XXXXXXXXXXXXXX/velero-backup-XXXXXXXXXXXXXX.tar.gz .
tar xzf velero-backup-XXXXXXXXXXXXXX.tar.gz
find . -name '*.json' -exec kubectl -f {} \;
```

### Restore volume data

## Output data restore
First install [Minio client CLI](https://docs.minio.io/docs/minio-client-complete-guide) accordingly to the linked guide and then add backup and production object storages with `mc config host add` command and then run these two commands to copy the data from backup server to new object storage:

```
mc mirror --overwrite --md5 backup/radar-intermediate-storage production/radar-intermediate-storage
mc mirror --overwrite --md5 backup/radar-output-storage production/radar-output-storage
````

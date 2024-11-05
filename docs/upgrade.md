## Upgrade instructions

Run the following instructions to upgrade an existing RADAR-Kubernetes cluster.

| :exclamation: Note                                                                                                                                                                                                                                |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Upgrading the major version of a PostgreSQL image is not supported. If necessary, we propose to use a `pg_dump` to dump the current data and a `pg_restore` to restore that data on a newer version. Please find instructions for this elsewhere. |

### Upgrade to RADAR-Kubernetes version 1.2

#### KSQL S erver and JDBC Connector service configuration
In `production.yaml` rename sections:

| Old Name                 | New Name                           |
|--------------------------|------------------------------------|
| radar_jdbc_connector     | radar_jdbc_connector_grafana       |
| radar_jdbc_connector_agg | radar_jdbc_connector_realtime_dashboard |

#### Data dashboard backend secret
You need to manually create this field in your secrets file: `management_portal.oauth_clients.radar_data_dashboard_backend.client_secret`.
If you're not using data dashboard backend you can put `null` in the value, otherwise make sure to generate a random secret and put it in that field.

#### MongoDB
MongoDB has been updated to a new version and it's not compatible with the current version that has been installed in the cluster. There are two pathes forward:
- Deleting the MongoDB and its volumes and then installing it and configuring Graylog again. This is the recommeneded approach since usually there is no important data is stored in MongoDB and the Graylog stack will be replaced in the next release.
- Upgrading MongoDB cluster. If you have configured Graylog significantly, then it might be better to upgrade the MongoDB instead of reinstalling it. Options are:
  - Manually exporting the databases and doing a clean reinstall as stated in the last step and the importing the data again.
  - Following official MongoDB upgrade instructions for version [5.0](https://www.mongodb.com/docs/manual/release-notes/5.0-upgrade-replica-set/), [6.0](https://www.mongodb.com/docs/manual/release-notes/6.0-upgrade-replica-set/) and [7.0](https://www.mongodb.com/docs/manual/release-notes/7.0-upgrade-replica-set/).

### Upgrade to RADAR-Kubernetes version 1.1.x
Before running the upgrade make sure to copy `environments.yaml.tmpl` to `environments.yaml` and if you've previously changed `environments.yaml` apply the changes again. This is necessary due to addition of `helmDefaults` and `repositories` configurations to this file.

### Upgrade to RADAR-Kubernetes version 1.0.0

Before running the upgrade, compare `etc/base.yaml` and `etc/base.yaml.gotmpl` with their `production.yaml` counterparts. Please ensure that all properties in `etc/base.yaml` are overridden in your `production.yaml` or that the `base.yaml` default value is fine, in which case no value needs to be provided in `production.yaml`.

To upgrade the initial services, run

```shell
kubectl delete -n monitoring deployments kube-prometheus-stack-kube-state-metrics
helm -n graylog uninstall mongodb
kubectl delete -n graylog pvc datadir-mongodb-0 datadir-mongodb-1
```

Note that this will remove your graylog settings but not your actual logs. This step is unfortunately needed to enable credentials on the Graylog database hosted by the mongodb chart. You will need to recreate the GELF TCP input source as during install.

Then run

```shell
helmfile -f helmfile.d/00-init.yaml apply --concurrency 1
helmfile -f helmfile.d/10-base.yaml --selector name=cert-manager-letsencrypt apply
```

To update the Kafka stack, run:

```shell
helmfile -f helmfile.d/10-base.yaml apply --concurrency 1
```

After this has succeeded, edit your `production.yaml` and change the `cp_kafka.customEnv.KAFKA_INTER_BROKER_PROTOCOL_VERSION` to the corresponding version documented in the [Confluent upgrade instructions](https://docs.confluent.io/platform/current/installation/upgrade.html) of your Kafka installation. Find the currently installed version of Kafka with `kubectl exec cp-kafka-0 -c cp-kafka-broker -- kafka-topics --version`.
When the `cp_kafka.customEnv.KAFKA_INTER_BROKER_PROTOCOL_VERSION` is updated, again run

```shell
helmfile -f helmfile.d/10-base.yaml apply
```

To upgrade to the latest PostgreSQL helm chart, in `production.yaml`, uncomment the line `postgresql.primary.persistence.existingClaim: "data-postgresql-postgresql-0"` to use the same data storage as previously. Then run

```shell
kubectl delete secrets postgresql
kubectl delete statefulsets postgresql-postgresql
helmfile -f helmfile.d/10-managementportal.yaml apply
```

If installed, `radar-appserver-postgresql`, uncomment the `production.yaml` line `radar_appserver_postgresql.primary.existingClaim: "data-radar-appserver-postgresql-postgresql-0"`. Then run

```shell
kubectl delete secrets radar-appserver-postgresql
kubectl delete statefulsets radar-appserver-postgresql-postgresql
helmfile -f helmfile.d/20-appserver.yaml apply
```

If installed, to upgrade `timescaledb`, uncomment the `production.yaml` line `timescaledb.primary.existingClaim: "data-timescaledb-postgresql-0"`. Then run

```shell
kubectl delete secrets timescaledb-postgresql
kubectl delete statefulsets timescaledb-postgresql
helmfile -f helmfile.d/20-grafana.yaml apply
```

If installed, to upgrade `radar-upload-postgresql`, uncomment the `production.yaml` line `radar_upload_postgresql.primary.existingClaim: "data-radar-upload-postgresql-postgresql-0"`. Then run

```shell
kubectl delete secrets radar-upload-postgresql
kubectl delete statefulsets radar-upload-postgresql-postgresql
helmfile -f helmfile.d/20-upload.yaml apply
```

If minio is installed, upgrade it with the following instructions:

```shell
# get minio PV and PVC
kubectl get pv | grep export-minio- | tr -s ' ' | cut -d ' ' -f 1,6 | tr '/' ' ' | cut -d ' ' -f 1,3 | tee minio-pv.list
# Uninstall the minio statefulset
helm uninstall minio
# Associate PV with the new PVC name
while read -r pv pvc
do
  # Don not delete PV
  kubectl patch pv $pv -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
  # Delete PVC
  kubectl delete pvc $pvc
  # Name of the new PVC
  newpvc=$(echo $pvc | sed 's/export-/data-/')
  # Associate PV with the new PVC name
  kubectl patch pv $pv -p '{"spec":{"claimRef":{"name": "'$newpvc'", "namespace": "default", "uid": null}}}'
  # Create new PVC
  cat <<EOF | sed "s/data-minio-i/$newpvc/" | kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  labels:
    app.kubernetes.io/instance: minio
    app.kubernetes.io/name: minio
  name: data-minio-i
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
EOF
done < minio-pv.list
# Do the new helm install.
helmfile -f helmfile.d/20-s3.yaml apply
```

Delete the redis stateful set (this will not delete the data on the volume)

```shell
kubectl delete statefulset redis-master
helmfile -f helmfile.d/20-s3.yaml sync --concurrency 1
```


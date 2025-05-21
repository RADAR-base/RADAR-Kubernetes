# Upgrade instructions

Run the following instructions to upgrade an existing RADAR-Kubernetes cluster.

## Upgrade to RADAR-Kubernetes version 1.3.0

This version introduces postgresql and TimescaleDB clusters managed by the CloudNativePG operator.

### Update `mods/migration/1.3.0.yaml` mods file

This file provides configuration for database migration. In the `cloundnative_postgresql` section, remove any database from the
`databases` list that is not needed. For instance:

```yaml
            ...
            databases:
              - managementportal
              # - appconfig
              # - appserver
              # - kratos
              - restsourceauthorizer
              # - uploadconnector
            source:
              ...
```

### Update `production.yaml` file

1. Add values for the number of Postgresql and TimescaleDB replicas. Change the number of replicas if desired:

```yaml
# Number of Postgres pods that will be installed
postgres_num_replicates: 2
# Number of TimescaleDB pods that will be installed
timescaledb_num_replicates: 2
```

2. If desired, add a section that changes the default storage size of the _PostgreSQL_ database:

```yaml
cloudnative_postgresql:
  cluster:
    cluster:
      storage:
        size: 10Gi
```

3. If desired, update the sections that change the default storage size of the respective _TimescaleDB_ databases:

```yaml
radar_jdbc_connector_grafana:
  radar-cloudnative-timescaledb:
    cluster:
      cluster:
        storage:
          # Change this to the desired size
          size: 50Gi
```

```yaml
radar_jdbc_connector_data_dashboard_backend:
  radar-cloudnative-timescaledb:
    cluster:
      cluster:
        storage:
          size: 508Gi
```

```yaml
radar_jdbc_connector_realtime_dashboard:
  radar-cloudnative-timescaledb:
    cluster:
      cluster:
        storage:
          size: 50Gi
```
 
4. Duplicate entry `grafana_metrics_db_user` ane rename to `grafana_metrics_endpoint_user` (keep the same value).

### Update `secrets.yaml` file

Duplicate secret `grafana_metrics_db_password` and rename to `grafana_metrics_endpoint_password` (keep the same value).

### Database migration

Important: before database migration the steps in the sections above must be completed.

The database migration process involves:

1. A manual import of existing _upload-connect-backend_ or _app-server_ databases into the CloudNativePG postgres cluster.
2. An automated import of the _management_portal_, _app-config_ and _rest-sources-authorizer_ databases into the CloudNativePG postgres cluster.
3. Post-migration cleanup.

#### 1. Manual import of _upload-connector_ and/or _app-server_ databases

Notes:
- Database passwords can be found in the `etc/secrets.yaml` file.
- Unless customized the username for all databases is `postgres`.

##### App-server database import

Perform when using the _app-server_ service. The username and password for the _app-server_ database are indicated with `<user>` and `<password>`, respectively. 
The username and password for the _management_portal_ database are indicated with `<mp-user>` and `<mp-password>`, respectively.

```shell
kubectl exec radar-appserver-postgresql-0 -- bash -c "PGPASSWORD=<password> pg_dump -U <user> appserver" > appserver.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -t -c 'CREATE DATABASE appserver'" 
cat appserver.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -d appserver" 
```

##### Upload-connector database import

Perform when using the _upload-connector_ service. The username and password for the _upload-connector_ database are indicated with `<user>` and `<password>`, respectively.
The username and password for the _management_portal_ database are indicated with `<mp-user>` and `<mp-password>`, respectively.

```shell
kubectl exec radar-upload-postgresql-0 -- bash -c "PGPASSWORD=<password> pg_dump -U <user> uploadconnector" > uploadconnector.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -t -c 'CREATE DATABASE uploadconnector'"
cat uploadconnector.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -d uploadconnector"
```

#### Update `environments.yaml` file

Add the _mods/migration/1.3.0.yaml_ file to the `values:` section, like so:

```yaml
environments:
  default:
    values:
      - ../etc/base.yaml
      - ../etc/production.yaml
      - ../etc/production.yaml.gotmpl
      - ../etc/secrets.yaml
      - ../mods/migration/1.3.0.yaml
```

#### 2. Automated import databases

Start the database migration of _management_portal_ and _TimescaleDB_ databases by using the auto-migration feature of
the CloudNativePG operator. Run the _helmfile_ command once with the `mods/migration/1.3.0.yaml` modification file.

```shell
helmfile sync 
```

### 3. Post migration cleanup

Perform these steps when the database migration is successful.

1. Remove any database passwords from the `secrets.yaml` file. An easy way to do this is to compare your `secrets.yaml`
   file to `base-secrets.yaml` file and remove any entries not present in `base-secrets.yaml`.

2. Turn of legacy database services. For this update the `production.yaml` file like so:

```yaml
postgresql:
  _install: false
...

radar_appserver_postgresql:
  _install: false
...

data_dashboard_timescaledb:
  _install: false
...

grafana_metrics_timescaledb:
  _install: false
...

realtime_dashboard_timescaledb:
  _install: false
...

radar_upload_postgresql:
  _install: false
...
```

3. Remove the _mods/migration/1.3.0.yaml_ file from the `environments.yaml` file.
4. Update the deployment:

```shell
helmfile sync
```

5. Remove any _pvc_ resource on the Kubernetes cluster associated with the old databases. Important: this step will permanently delete the data! Only perform this step when sure the migration completed successfully.
   The relevant _pvc_ names are:

- `data-postgresql-0`
- `data-data-dashboard-timescaledb-postgresql-0`
- `data-grafana-metrics-timescaledb-postgresql-0`
- `data-radar-appserver-postgresql-0`
- `data-radar-upload-postgresql-0`
- `data-realtime-dashboard-timescaledb-postgresql-0`

## Upgrade to RADAR-Kubernetes version 1.2.0

### Update `production.yaml` file

1. Remove any line beginning with `_chart_version:`.
2. Remove any line beginning with `imageTag:`.
3. Add email server config to `management_portal` and `radar_appserver` sections analogous to:

```yaml
    management_portal:
      smtp:
        enabled: true
        host: smtp
        port: 25
        from: noreply@example.com
        starttls: false
        auth: true
```

```yaml
    radar_appserver:
      smtp:
        enabled: true
        host: smtp
        port: 25
        from: noreply@example.com
        starttls: false
        auth: true
```

4. Update _TimescaleDB_ database configuration:

- Rename `timescaledb_username` to `grafana_metrics_db_username`
- Remove `grafana_metrics_username` and `timescaledb_db_name` variables.
- When using _realtime-dashboard_, add `realtime_dashboard_db_username` that points to the current value of
  `timescaledb_username`.

5. Refactor the _timescaledb_ database configuration:

Rename `timescaledb:` to `grafana_metrics_timescaledb:`.When using _data-dashboard-backend_ or
_realtime-dashboard-backend_, copy the `grafana_metrics_timescaledb:` entry to `realtime_dashboard_timescaledb:` and
`data_dashboard_timescaledb:` similar to:

```yaml
grafana_metrics_timescaledb:
  _install: true
  ...
  postgresql:
    replication:
      ...

data_dashboard_timescaledb:
  _install: true
  ...
  postgresql:
    replication:
      ...

realtime_dashboard_timescaledb:
  _install: true
  ...
  postgresql:
    replication:
      ...
```

IMPORTANT: For databases where data should persist after the update uncomment the respective `existingClaim` field. For
example:

```yaml
grafana_metrics_timescaledb:
  postgresql:
    primary:
      persistence:
        existingClaim: "data-timescaledb-postgresql-0"
```

6. Add `postgresql:` indent to _postgresql_ and _timescaledb_ related database configurations like so:

```yaml

postgresql:
  _install: true
  ...
  postgresql:
    replication:
      ...

grafana_metrics_timescaledb:
  _install: true
  ...
  postgresql:
    replication:
      ...

data_dashboard_timescaledb:
  _install: true
  ...
  postgresql:
    replication:
      ...

realtime_dashboard_timescaledb:
  _install: true
  ...
  postgresql:
    replication:
      ...

radar_appserver_postgresql:
  _install: true
  ...
  postgresql:
    replication:
    ...

radar_upload_postgresql:
  _install: true
  ...
  postgresql:
    replication:
    ...
```

7. Rename `radar_jdbc_connector:` to `radar_jdbc_connector_grafana:`.

8. Add:

```yaml
kratos:
  _install: false
  ...

kratos_ui:
  _install: false
  ...
```

### Update `secrets.yaml` file

1. Add the following new secrets to the `secrets.yaml` file to corresponding sections:

```yaml
management_portal:
  oauth_clients:
    radar_data_dashboard_backend:
      client_secret: <add your own random secret here>
```

```yaml
radar_appserver:
  smtp:
    username: <your smtp username>
    password: <your smtp password>
```

```yaml
data_dashboard_db_password: <same password as timescaledb_password>
realtime_dashboard_db_password: <same password as timescaledb_password>
```

2. Rename the `grafana_metrics_password` secret to `grafana_metrics_db_password` and `timescaledb_password` to
   `data_dashboard_db_password`.

3. Add `postgresql:` indent to `radar_appserver_postgresql:` like so:

```yaml
radar_appserver_postgresql:
  postgresql:
    global:
      postgresql:
        auth:
          postgresPassword: <password>
    auth:
```

### MongoDB

MongoDB has been updated to a new version and it's not compatible with the current version that has been installed in
the cluster. There are three paths forward:

- Keeping the current version. Add legacy chart version to your `production.yaml`:

```yaml
mongodb:
  _install: true
  # Keep the current version of MongoDB in order to be compatible with the stored data.
  _chart_version: 11.1.10
```

- Deleting the MongoDB and its volumes and then installing it and configuring Graylog again. This is the recommended
  approach since usually there is no important data is stored in MongoDB and the Graylog stack will be replaced in the
  next release.
- Upgrading MongoDB cluster. If you have configured Graylog significantly, then it might be better to upgrade the
  MongoDB instead of reinstalling it. Options are:
    - Manually exporting the databases and doing a clean reinstall as stated in the last step and the importing the data
      again.
    - Following official MongoDB upgrade instructions for
      version [5.0](https://www.mongodb.com/docs/manual/release-notes/5.0-upgrade-replica-set/), [6.0](https://www.mongodb.com/docs/manual/release-notes/6.0-upgrade-replica-set/)
      and [7.0](https://www.mongodb.com/docs/manual/release-notes/7.0-upgrade-replica-set/).

## Upgrade to RADAR-Kubernetes version 1.1.x

Before running the upgrade make sure to copy `environments.yaml.tmpl` to `environments.yaml` and if you've previously
changed `environments.yaml` apply the changes again. This is necessary due to addition of `helmDefaults` and
`repositories` configurations to this file.

## Upgrade to RADAR-Kubernetes version 1.0.0

Before running the upgrade, compare `etc/base.yaml` and `etc/base.yaml.gotmpl` with their `production.yaml`
counterparts. Please ensure that all properties in `etc/base.yaml` are overridden in your `production.yaml` or that the
`base.yaml` default value is fine, in which case no value needs to be provided in `production.yaml`.

To upgrade the initial services, run

```shell
kubectl delete -n monitoring deployments kube-prometheus-stack-kube-state-metrics
helm -n graylog uninstall mongodb
kubectl delete -n graylog pvc datadir-mongodb-0 datadir-mongodb-1
```

Note that this will remove your graylog settings but not your actual logs. This step is unfortunately needed to enable
credentials on the Graylog database hosted by the mongodb chart. You will need to recreate the GELF TCP input source as
during install.

Then run

```shell
helmfile -f helmfile.d/00-init.yaml apply --concurrency 1
helmfile -f helmfile.d/10-base.yaml --selector name=cert-manager-letsencrypt apply
```

To update the Kafka stack, run:

```shell
helmfile -f helmfile.d/10-base.yaml apply --concurrency 1
```

After this has succeeded, edit your `production.yaml` and change the
`cp_kafka.customEnv.KAFKA_INTER_BROKER_PROTOCOL_VERSION` to the corresponding version documented in
the [Confluent upgrade instructions](https://docs.confluent.io/platform/current/installation/upgrade.html) of your Kafka
installation. Find the currently installed version of Kafka with
`kubectl exec cp-kafka-0 -c cp-kafka-broker -- kafka-topics --version`.
When the `cp_kafka.customEnv.KAFKA_INTER_BROKER_PROTOCOL_VERSION` is updated, again run

```shell
helmfile -f helmfile.d/10-base.yaml apply
```

To upgrade to the latest PostgreSQL helm chart, in `production.yaml`, uncomment the line
`postgresql.primary.persistence.existingClaim: "data-postgresql-postgresql-0"` to use the same data storage as
previously. Then run

```shell
kubectl delete secrets postgresql
kubectl delete statefulsets postgresql-postgresql
helmfile -f helmfile.d/10-managementportal.yaml apply
```

If installed, `radar-appserver-postgresql`, uncomment the `production.yaml` line
`radar_appserver_postgresql.primary.existingClaim: "data-radar-appserver-postgresql-postgresql-0"`. Then run

```shell
kubectl delete secrets radar-appserver-postgresql
kubectl delete statefulsets radar-appserver-postgresql-postgresql
helmfile -f helmfile.d/20-appserver.yaml apply
```

If installed, to upgrade `TimescaleDB`, uncomment the `production.yaml` line
`timescaledb.primary.existingClaim: "data-timescaledb-postgresql-0"`. Then run

```shell
kubectl delete secrets timescaledb-postgresql
kubectl delete statefulsets timescaledb-postgresql
helmfile -f helmfile.d/20-grafana.yaml apply
```

If installed, to upgrade `radar-upload-postgresql`, uncomment the `production.yaml` line
`radar_upload_postgresql.primary.existingClaim: "data-radar-upload-postgresql-postgresql-0"`. Then run

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


# Upgrade instructions

Run the following instructions to upgrade an existing RADAR-Kubernetes cluster.

## Upgrade to RADAR-Kubernetes version 1.3.0

This version introduces postgresql and TimescaleDB clusters managed by the CloudNativePG operator.

### Update `production.yaml` file

1. Add values for the number of Postgresql and TimescaleDB replicas. Change the number of replicas if desired:

```yaml
# Number of Postgres pods that will be installed
postgres_num_replicates: 2
# Number of TimescaleDB pods that will be installed
timescaledb_num_replicates: 2
```

2. Add the CloudNativePG section. Note that the `_install` flag is set to `false`. Change the following:

- Change the disk size in `cloudnative.cluster.cluster.storage.size` if desired.
- Remove databases that are not needed in the `cloudnative.cluster.cluster.recovery.import.databases` section.

```yaml
cloudnative_postgresql:
  _install: false
  _extra_timeout: 0
  cluster:
    cluster:
      storage:
        # Change this to the desired size
        size: 10Gi
      cluster:
        mode: recovery
        recovery:
          method: import
          import:
            type: monolith
            # Update this list with the databases that need to be imported
            databases:
              - managementportal
              - appconfig
              - appserver
              - kratos
              - restsourceauthorizer
              - uploadconnector
            source:
              host: postgresql
              username: postgres
              database: postgres
              sslMode: prefer
              passwordSecret:
                create: false
                name: postgresql
                key: postgres-password
```

3. Update JDBC-connector sections to use the new CloudNativePG operator. Note that the `_install` flag is set to
   `false`. Change the disk size in `cloudnative.cluster.cluster.storage.size` if desired.

```yaml
radar_jdbc_connector_grafana:
  ...
  radar-cloudnative-timescaledb:
    enabled: false
    cluster:
      cluster:
        storage:
          # Change this to the desired size
          size: 8Gi
          mode: recovery
        import:
          databases:
            - grafana-metrics
          source:
            host: grafana-metrics-timescaledb-postgresql
            username: postgres
            database: postgres
            sslMode: prefer
            passwordSecret:
              create: false
              name: grafana-metrics-timescaledb-postgresql
              key: postgres-password
```

```yaml
radar_jdbc_connector_data_dashboard_backend:
  ...
  radar-cloudnative-timescaledb:
    enabled: false
    cluster:
      cluster:
        initdb:
          database: data-dashboard
          owner: data-dashboard
        storage:
          # Change this to the desired size
          size: 8Gi
        mode: recovery
        import:: recovery
        databases:
          - data-dashboard
        source:
          host: data-dashboard-timescaledb-postgresql
          username: postgres
          database: postgres
          sslMode: prefer
          passwordSecret:
            create: false
            name: data-dashboard-timescaledb-postgresql
            key: postgres-password
```

```yaml
radar_jdbc_connector_realtime_dashboard:
  ...
  radar-cloudnative-timescaledb:
    enabled: false
    cluster:
      cluster:
        initdb:
          database: realtime-dashboard
          owner: realtime-dashboard
        storage:
          # Change this to the desired size
          size: 8Gi
        mode: recovery
        import:
          databases:
            - realtime-dashboard
          source:
            host: realtime-dashboard-timescaledb-postgresql
            username: postgres
            database: postgres
            sslMode: prefer
            passwordSecret:
              create: false
              name: realtime-dashboard-timescaledb-postgresql
              key: postgres-password
```

4. Add a `grafana_username` entry. When you use a `grafana_metrics_db_username` username different from `postgres` duplicate (!)
   `grafana_metrics_db_username` to `grafana_metrics_endpoint_username` like so:

```yaml
grafana_username: grafana
grafana_metrics_db_username: my-username
grafana_metrics_endpoint_username: my-username
```

### Update `secrets.yaml` file

1. Duplicate (!) the `grafana_metrics_db_password` to a new `grafana_metrics_endpoint_password` entry similar to :

```yaml
grafana_metrics_db_password: my-password
grafana_metrics_endpoint_password: my-password
```

### Database migration

Important: before database migration the steps in the sections above must be completed.

The database migration process involves:

1. (optional) When using _upload-connect-backend_ , _kratos_ or _appserver_ services, perform a manual import of
   existing databases into the _management_portal_ postgres database.
2. Automated import of the _management_portal_ database into the new CloudNativePG postgres cluster(s).
3. Post-migration cleanup.

#### 1. Manual import of _upload-connect-backend_ , _kratos_ or _appserver_ databases

Note: the database passwords can be found in the `secrets.yaml` file.

##### App-server database import

Perform when using the _appserver_ service.

```shell
kubectl exec radar-appserver-postgresql-0 -- bash -c "PGPASSWORD=<appserver database password> pg_dump -U postgres appserver" > appserver.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -t -c 'CREATE DATABASE appserver'" 
cat appserver.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d appserver" 
````

##### Upload-connect-backend database import

Perform when using the _upload-connect-backend_ service.

```shell
kubectl exec radar-upload-postgresql-0 -- bash -c "PGPASSWORD=<upload-connect-backend database password> pg_dump -U postgres uploadconnector" > uploadconnector.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -t -c 'CREATE DATABASE uploadconnector'"
cat uploadconnector.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d uploadconnector"
```

##### Kratos database import

Perform when using the _kratos_ service.

```shell
kubectl exec radar-kratos-postgresql-0 -- bash -c "PGPASSWORD=<kratos database password> pg_dump -U postgres kratos" > kratos.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -t -c 'CREATE DATABASE kratos'"
cat kratos.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d kratos"
```

#### 2. Automated import databases

Start the database migration of _management_portal_ and _TimescalreDB_ databases by using the auto-migration featire of
the CloudNativePG operator. For this, enable the CloudNativePG databases by setting the `_install` flag to `true` in the
`production.yaml` file:

```yaml
cloudnative_postgresql:
  _install: true
  ...

radar_jdbc_connector_data_dashboard_backend:
  ...
  radar-cloudnative-timescaledb:
    enabled: true

radar_jdbc_connector_grafana:
  ...
  radar-cloudnative-timescaledb:
    enabled: true

radar_jdbc_connector_realtime_dashboard:
  ...
  radar-cloudnative-timescaledb:
    enabled: true
```

And run:

```shell
helmfile sync
```

### 3. Post migration cleanup

Perform these steps when the database migration is successful.

1. Remove any database passwords from the `secrets.yaml` file. An easy way to do this is to compare your `secrets.yaml`
   file to `base-secrets.yaml` file and remove any entries not present in `base-secrets.yaml`.

2. Remove the following keys in `production.yaml`:

- `cloudnative_postgresql.cluster.cluster.mode`
- `cloudnative_postgresql.cluster.cluster.recovery`
- `radar_jdbc_connector_grafana.radar-cloudnative-timescaledb.cluster.cluster.mode`
- `radar_jdbc_connector_grafana.radar-cloudnative-timescaledb.cluster.cluster.recovery`
- `radar_jdbc_connector_data_dashboard_backend.radar-cloudnative-timescaledb.cluster.cluster.mode`
- `radar_jdbc_connector_data_dashboard_backend.radar-cloudnative-timescaledb.cluster.cluster.recovery`
- `radar_jdbc_connector_realtime_dashboard.radar-cloudnative-timescaledb.cluster.cluster.mode`
- `radar_jdbc_connector_realtime_dashboard.radar-cloudnative-timescaledb.cluster.cluster.recovery`

3. Turn of legacy database services. For this update the `production.yaml` file like so:

```yaml
postgresql:
  _install: true
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
  _install: true
...
```

and run:

```shell
helmfile sync
```

5Remove any _pvc_ resource on the Kubernetes cluster associated with the old databases. The _pvc_ names for these are:

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
3. Add email server config to `management_portal` and `radar-appserver` sections analogous to:

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

5. For databases where data should persist after the update uncomment the respective `existingClaim` field. Example:

```yaml
realtime_dashboard_timescaledb:
  postgresql:
    primary:
      persistence:
        existingClaim: "data-timescaledb-postgresql-0"
```

### Update `secrets.yaml` file

1. Add the following new secrets to the `secrets.yaml` file to correctponding sections:

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

### MongoDB

MongoDB has been updated to a new version and it's not compatible with the current version that has been installed in
the cluster. There are two pathes forward:

- Deleting the MongoDB and its volumes and then installing it and configuring Graylog again. This is the recommeneded
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


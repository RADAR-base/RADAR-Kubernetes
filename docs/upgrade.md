# Upgrade instructions

Run the following instructions to upgrade an existing RADAR-Kubernetes cluster.

## Upgrade to RADAR-Kubernetes version 1.3.0

This version introduces postgresql and TimescaleDB clusters managed by the CloudNativePG operator.

### Update `production.yaml` file

1. Add values for the number of Postgresql and TimescaleDB replicas:

```yaml
# Number of Postgres pods that will be installed
postgres_num_replicates: 2
# Number of TimescaleDB pods that will be installed
timescaledb_num_replicates: 2
```

2. Add the following sections for CloudNativePG charts:

```yaml
cloudnativepg_operator:
  _chart_version: 0.23.2
  _extra_timeout: 0

cloudnative_postgresql:
  _install: true
  _chart_version: 0.1.0
  _extra_timeout: 0
  cluster:
    cluster:
      instances: 1
      storage:
        size: 10Gi
#  --  Uncomment this when migrating from the legacy postgresql databases
#  cluster:
#    mode: recovery
#    recovery:
#      method: import
#      import:
#        type: monolith
#        databases:
#          - managementportal
#          - appconfig
#          - kratos
#          - hydra
#          - restsourceauthorizer
#        source:
#          host: postgresql
#          username: postgres
#          database: postgres
#          sslMode: prefer
#          passwordSecret:
#            create: false
#            name: postgresql
#            key: postgres-password
```

3. Update JDBC-connector sections to use the new CloudNativePG operator:

```yaml
radar_jdbc_connector_grafana:
  ...
  radar-cloudnative-timescaledb:
    enabled: true
    cluster:
      cluster:
        initdb:
          database: grafana-metrics
          owner: grafana-metrics
        storage:
          size: 8Gi
#  --  Uncomment this when migrating from the legacy postgresql databases
#      mode: recovery
#      import:
#        databases:
#          - grafana-metrics
#        source:
#          host: grafana-metrics-timescaledb-postgresql
#          username: postgres
#          database: postgres
#          sslMode: prefer
#          passwordSecret:
#            create: false
#            name: grafana-metrics-timescaledb-postgresql
#            key: postgres-password
```

```yaml
radar_jdbc_connector_data_dashboard_backend:
  ...
  radar-cloudnative-timescaledb:
    enabled: true
    cluster:
      cluster:
        initdb:
          database: data-dashboard
          owner: data-dashboard
        storage:
          size: 8Gi
#  --  Uncomment this when migrating from the legacy postgresql databases
#      mode: recovery
#      import:
#        databases:
#          - data-dashboard
#        source:
#          host: data-dashboard-timescaledb-postgresql
#          username: postgres
#          database: postgres
#          sslMode: prefer
#          passwordSecret:
#            create: false
#            name: data-dashboard-timescaledb-postgresql
#            key: postgres-password
```

```yaml
radar_jdbc_connector_realtime_dashboard:
  ...
  radar-cloudnative-timescaledb:
    enabled: true
    cluster:
      cluster:
        initdb:
          database: realtime-dashboard
          owner: realtime-dashboard
        storage:
          size: 8Gi
#  --  Uncomment this when migrating from the legacy postgresql databases
#      mode: recovery
#      import:
#        databases:
#          - realtime-dashboard
#        source:
#          host: realtime-dashboard-timescaledb-postgresql
#          username: postgres
#          database: postgres
#          sslMode: prefer
#          passwordSecret:
#            create: false
#            name: realtime-dashboard-timescaledb-postgresql
#            key: postgres-password
```

4. Update grafana and grafana prometheus usernames. Add `grafana_username: grafana` and duplicate (!)
   `grafana_metrics_db_username` to
   `grafana_metrics_endpoint_username` like so:

```yaml
grafana_username: grafana
grafana_metrics_db_username: postgres
grafana_metrics_endpoint_username: postgres
```

### Database migration

The database migration process involves:

1. (optional) Manual import of existing _app-server_, _kratos_ and _upload-connect-backend_ databases into the
   _management_portal_
   postgres database. Only performed when _app_server_, _kratos_ or _upload-connect-backend_ databases are deployed.
2. Automated import of the _management_portal_ database into the new CloudNativePG postgres cluster.
3. (optional) Automated import of the timescaledb databases (connected to the JDBC-connector services) into the new
   CloudNativePG TimescaleDB cluster. Only performed when the _data-dashboard-backend_, _realtime_dashboard_ or
   _grafana-metrics_ databases are deployed.

### (optional) App-server database migration

```shell
kubectl exec radar-appserver-postgresql-0 -- bash -c "PGPASSWORD=<appserver database password> pg_dump -U postgres appserver" > appserver.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d appserver -t -c 'CREATE DATABASE appserver'" 
cat appserver.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d appserver" 
````

### (optional) Upload-connect-backend database migration

```shell
kubectl exec radar-upload-postgresql-0 -- bash -c "PGPASSWORD=<upload-connect-backend database password> pg_dump -U postgres upload_connect_backend" > uploadconnector.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d upload_connect_backend -t -c 'CREATE DATABASE uploadconnector'"
cat uploadconnector.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d uploadconnector"
```

### (optional) Kratos database migration

```shell
kubectl exec radar-kratos-postgresql-0 -- bash -c "PGPASSWORD=<kratos database password> pg_dump -U postgres kratos" > kratos.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d kratos -t -c 'CREATE DATABASE kratos'"
cat kratos.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<management-portal database password> psql -U postgres -d kratos"
```

Note: database passwords can be found in the `secrets.yaml` file.

2. Activate import of the _management-portal_ database in the `production.yaml` file. Make sure include all the
   databases that need to be imported. For example:

```yaml
cloudnative_postgresql:
  _install: true
  ...
  cluster:
    mode: recovery
    recovery:
      import:
        databases:
          - managementportal
          - restsourceauthorizer
          - appserver
          - uploadconnector
          - ...
```

Note: make sure to include all the databases that need to be imported, but to not include database nor present in the
_management_portal_ database.

3. Trigger the import of the _management_portal_ database enabling the _cloudnative-pg_ operator and
   _radar-cloudnative-postgresql_
   services in `production.yaml`:

```yaml
cloudnative_pg:
  _install: true
  ...
...

cloudnative_postgresql:
  _install: true
  ...
```

And run:

```shell
helmfile sync
```

4. (Optional) Activate the import of any deployed TimescaleDB database in the `production.yaml` file. This can be
   achieved by adding the import setting to the respective JDBC-connector definitions. For instance, to migrate the
   _data-dashboard-backend_ TimescaleDB, add the following to the `radar_jdbc_connector_data_dashboard_backend` section:

```yaml
radar_jdbc_connector_data_dashboard_backend:
  _install: true
  ...
  radar-cloudnative-timescaledb:
    cluster:
      mode: recovery
      import:
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

Equivalent changes should be made for the _realtime-dashboard_ and _grafana-metrics_ TimescaleDB databases. Make sure to
update the example above with the respective values for `datsbases:`, `host:` and `passwordSecret.name:` fields.

5. Redeploy the _radar-jdbc-connector_ services to trigger the import of the TimescaleDB databases:

```shell
helmfile sync
```

6. Remove any database passwords from the `secrets.yaml` file (new passwords will be generated by the operator, but
   nevertheless...). As reference, look at the changes of your `secrets.yaml` compared to the `base-secrets.yaml` file.

7. When migration was successful, stop the legacy database services and remove _pvc_ and _pv_ resources associated with
   the old databases.

### Update `environments.yaml` file

There has been a change in the way RADAR-base helm charts manage charts by exernal repositories. In file
`environments.yaml`
add the following in under the `repositories` section (`...` denotes repositories already present in the file):

```yaml
repositories:
    ...
      - name: cloudnative-pg
      url: https://cloudnative-pg.github.io/charts
      - name: grafana
        url: https://grafana.github.io/helm-charts
      - name: kratos
        url: https://k8s.ory.sh/helm/charts
      - name: hydra
        url: https://k8s.ory.sh/helm/charts
```

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


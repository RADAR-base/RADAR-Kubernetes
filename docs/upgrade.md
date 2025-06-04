# Upgrade instructions

<!-- TOC -->

* [Upgrade instructions](#upgrade-instructions)
    * [Upgrade to RADAR-Kubernetes version 1.3.0](#upgrade-to-radar-kubernetes-version-130)
        * [Update `mods/migration/1.3.0.yaml` mods file](#update-modsmigration130yaml-mods-file)
        * [Update `production.yaml` file](#update-productionyaml-file)
        * [Update `secrets.yaml` file](#update-secretsyaml-file)
        * [Database migration](#database-migration)
            * [1. Disable services that write to the databases](#1-disable-services-that-write-to-the-databases)
            * [2. Manual import of _upload-connector_ and/or
              _app-server_ databases](#2-manual-import-of-_upload-connector_-andor-_app-server_-databases)
                * [App-server database import](#app-server-database-import)
                * [Upload-connector database import](#upload-connector-database-import)
            * [3. Create users in the postgres database](#3-create-users-in-the-postgres-database)
            * [4. Automated import databases](#4-automated-import-databases)
                * [Update `environments.yaml` file](#update-environmentsyaml-file)
            * [5. Re-enable services that write to the databases](#5-re-enable-services-that-write-to-the-databases)
        * [Post-migration cleanup](#post-migration-cleanup)
    * [Upgrade to RADAR-Kubernetes version 1.2.0](#upgrade-to-radar-kubernetes-version-120)
        * [Update `production.yaml` file](#update-productionyaml-file-1)
        * [Update `secrets.yaml` file](#update-secretsyaml-file-1)
        * [MongoDB](#mongodb)
    * [Upgrade to RADAR-Kubernetes version 1.1.x](#upgrade-to-radar-kubernetes-version-11x)
    * [Upgrade to RADAR-Kubernetes version 1.0.0](#upgrade-to-radar-kubernetes-version-100)
    * [Supporting tasks](#supporting-tasks)
        * [Disable data ingestion](#disable-data-ingestion)
        * [Disable database changes](#disable-database-changes)

<!-- TOC -->

Run the following instructions to upgrade an existing RADAR-Kubernetes cluster.

## Upgrade to RADAR-Kubernetes version 1.3.0

This version introduces postgresql and Timescaledb clusters managed by the CloudNativePG operator.

### Update `mods/migration/1.3.0.yaml` mods file

This file provides configuration for database migration. In the `cloudnative_postgresql:` section, remove any database
from the `databases:` list that is not needed. For instance:

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

1. Remove any line beginning with `_chart_version:`.

2. Add values for the number of Postgresql and Timescaledb replicas. Change the number of replicas if desired:

```yaml
# Number of Postgres pods that will be installed
postgres_num_replicates: 2
```

3. If desired, add a section that changes the default storage size of the Postgresql cluster to be created:

```yaml
cloudnative_postgresql:
  cluster:
    cluster:
      storage:
        size: 10Gi
```

4. Set legacy versions for Timescaledb in jdbc-connector sections. If desired, change the storage size of the
   respective databases:

Note: major version upgrades performed by the CloudNativePG operator are currently under development. When v1.26 is
released, the Timescaledb databases can be upgraded to the a newer version.

```yaml
radar_jdbc_connector_grafana:
  timescaledb:
    cluster:
      # Do not remove: needed for legacy version. Can be removed after upgrade to new postgresql version (handled in future RADAR-base update).
      version:
        postgresql: "15"
        Timescaledb: "2.11"
      cluster:
        storage:
          # Change this to the desired size
          size: 50Gi
```

```yaml
radar_jdbc_connector_data_dashboard_backend:
  timescaledb:
    cluster:
      # Do not remove: needed for legacy version. Can be removed after upgrade to new postgresql version (handled in future RADAR-base update).
      version:
        postgresql: "15"
        Timescaledb: "2.11"
      cluster:
        storage:
          size: 508Gi
```

```yaml
radar_jdbc_connector_realtime_dashboard:
  timescaledb:
    cluster:
      # Do not remove: needed for legacy version. Can be removed after upgrade to new postgresql version (handled in future RADAR-base update).
      version:
        postgresql: "15"
        Timescaledb: "2.11"
      cluster:
        storage:
          size: 50Gi
```

5. Duplicate entry `grafana_metrics_db_user` ane rename to `grafana_metrics_endpoint_user` (keep the same value).

### Update `secrets.yaml` file

Duplicate secret `grafana_metrics_db_password` and rename to `grafana_metrics_endpoint_password` (keep the same value).

### Database migration

Important: before database migration the steps in the sections above must have been completed successfully.

Notes:

- Database passwords can be found in the `etc/secrets.yaml` file.
- Unless customized the username for all databases is `postgres`.

#### 1. Disable services that write to the databases

To prevent any changes to the databases during the migration, disable all services that write to the databases (
see [Disable database changes](#disable-data-ingestion)).

#### 2. Manual import of _upload-connector_ and/or _app-server_ databases

##### App-server database import

Perform when using the _app-server_ service. The username and password for the _app-server_ database are indicated with
`<user>` and `<password>`, respectively. The username and password for the _management_portal_ database are indicated
with `<mp-user>` and `<mp-password>`, respectively.

```shell
kubectl exec radar-appserver-postgresql-0 -- bash -c "PGPASSWORD=<password> pg_dump -U <user> appserver" > appserver.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -t -c 'CREATE DATABASE appserver'" 
cat appserver.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -d appserver" 
```

##### Upload-connector database import

Perform when using the _upload-connector_ service. The username and password for the _upload-connector_ database are
indicated with `<user>` and `<password>`, respectively.
The username and password for the _management_portal_ database are indicated with `<mp-user>` and `<mp-password>`,
respectively.

```shell
kubectl exec radar-upload-postgresql-0 -- bash -c "PGPASSWORD=<password> pg_dump -U <user> uploadconnector" > uploadconnector.sql
kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -t -c 'CREATE DATABASE uploadconnector'"
cat uploadconnector.sql | kubectl exec -i postgresql-0 -- bash -c "PGPASSWORD=<mp-password> psql -U <mp-user> -d uploadconnector"
```

#### 3. Create users in the postgres database

Log into the Postgresql database using the _psql_ utility. The username and password for the Postgresql database are
indicated with `<user>` and `<password>`, respectively. Note that the `--pset pager=off` option is used to disable
paging in the output to prevent problems with certain statements.

```shell
kubectl exec postgresql-0 -it -- sh -c 'PGPASSWORD=<password> psql -U <user> --pset pager=off'
```

Create database users and set ownership of the databases (remove any database that is not needed):

```sql
CREATE
USER managementportal;
CREATE
USER restsourceauthorizer;
CREATE
USER appconfig;
CREATE
USER kratos;
CREATE
USER appserver;
CREATE
USER uploadconnector;
ALTER
DATABASE managementportal OWNER to managementportal;
ALTER
DATABASE restsourceauthorizer OWNER to restsourceauthorizer;
ALTER
DATABASE appconfig OWNER to appconfig;
ALTER
DATABASE kratos OWNER to kratos;
ALTER
DATABASE appserver OWNER to appserver;
ALTER
DATABASE uploadconnector OWNER to uploadconnector;
```

Transfer ownership of all tables in respective databases to the new users. Ignore sections of the command for any
database that is not used.

```sql
\c
managementportal
CREATE
OR REPLACE FUNCTION exec(text) returns text language plpgsql volatile
  AS
$f$
BEGIN
EXECUTE $1;
RETURN $1;
END;
$f$;
SELECT exec ( 'ALTER TABLE ' || table_name || ' OWNER TO ' || table_catalog )
FROM information_schema.tables
WHERE table_schema = 'public';
SELECT exec ( 'ALTER SEQUENCE ' || sequence_name || ' OWNER TO ' || sequence_catalog )
FROM information_schema.sequences
WHERE sequence_schema = 'public';

\c
restsourceauthorizer
CREATE
OR REPLACE FUNCTION exec(text) returns void language plpgsql volatile
  AS
$f$
BEGIN
EXECUTE $1;
END;
$f$;
SELECT exec ( 'ALTER TABLE ' || table_name || ' OWNER TO ' || table_catalog )
FROM information_schema.tables
WHERE table_schema = 'public';
SELECT exec ( 'ALTER SEQUENCE ' || sequence_name || ' OWNER TO ' || sequence_catalog )
FROM information_schema.sequences
WHERE sequence_schema = 'public';

\c
appconfig
CREATE
OR REPLACE FUNCTION exec(text) returns text language plpgsql volatile
  AS
$f$
BEGIN
EXECUTE $1;
RETURN $1;
END;
$f$;
SELECT exec ( 'ALTER TABLE ' || table_name || ' OWNER TO ' || table_catalog )
FROM information_schema.tables
WHERE table_schema = 'public';
SELECT exec ( 'ALTER SEQUENCE ' || sequence_name || ' OWNER TO ' || sequence_catalog )
FROM information_schema.sequences
WHERE sequence_schema = 'public';

\c
kratos
CREATE
OR REPLACE FUNCTION exec(text) returns text language plpgsql volatile
  AS
$f$
BEGIN
EXECUTE $1;
RETURN $1;
END;
$f$;
SELECT exec ( 'ALTER TABLE ' || table_name || ' OWNER TO ' || table_catalog )
FROM information_schema.tables
WHERE table_schema = 'public';
SELECT exec ( 'ALTER SEQUENCE ' || sequence_name || ' OWNER TO ' || sequence_catalog )
FROM information_schema.sequences
WHERE sequence_schema = 'public';

\c
appserver
CREATE
OR REPLACE FUNCTION exec(text) returns text language plpgsql volatile
  AS
$f$
BEGIN
EXECUTE $1;
RETURN $1;
END;
$f$;
SELECT exec ( 'ALTER TABLE ' || table_name || ' OWNER TO ' || table_catalog )
FROM information_schema.tables
WHERE table_schema = 'public';
SELECT exec ( 'ALTER SEQUENCE ' || sequence_name || ' OWNER TO ' || sequence_catalog )
FROM information_schema.sequences
WHERE sequence_schema = 'public';

\c
uploadconnector
CREATE
OR REPLACE FUNCTION exec(text) returns text language plpgsql volatile
  AS
$f$
BEGIN
EXECUTE $1;
RETURN $1;
END ;
$f$;
SELECT exec ( 'ALTER TABLE ' || table_name || ' OWNER TO ' || table_catalog )
FROM information_schema.tables
WHERE table_schema = 'public';
SELECT exec ( 'ALTER SEQUENCE ' || sequence_name || ' OWNER TO ' || sequence_catalog )
FROM information_schema.sequences
WHERE sequence_schema = 'public';
```

#### 4. Automated import of the Postgresql database

##### Update `environments.yaml` file

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

Start the database migration of _management_portal_ and Timescaledb databases by using the auto-migration feature of
the CloudNativePG operator. Run the _helmfile_ command once with the `mods/migration/1.3.0.yaml` modification file.

```shell
helmfile sync 
```

Confirm that all database services initialize successfully.

#### 5. Manual import of Timescaledb databases

For the grafana and realtime-dashboard databases, the migration is performed manually due to current problems with the migration of Timescaledb hypertables.
The perform the following steps for any of these services when in use:

1. attach shell to grafana-timescaledb pod:

```shell
kubectl exec grafana-timescaledb-1 -it -- bash
```

and export/import the database into the new grafana-timescale cluster:

```
export source=postgres://<user>:<password>@grafana-metrics-timescaledb-postgresql/grafana-metrics
pg_dump -d "$source" -fc -f /controller/grafana.bak
psql -d grafana -t -c 'select timescaledb_pre_restore();'
pg_restore -d grafana -fc /controller/grafana.bak
psql -d grafana -t -c 'select timescaledb_post_restore();'
```

2. Attach shell to realtime-dashboard-timescaledb pod:

```shell
kubectl exec realtime-dashboard-timescaledb-1 -it -- bash
```

and export/import the database into the new realtime-dashboard timescaledb cluster:

```
export SOURCE=postgres://<user>:<password>@realtime-dashboard-timescaledb-postgresql/realtime-dashboard
pg_dump -d "$SOURCE" -Fc -f /controller/realtime-dashboard.bak
psql -d realtime-dashboard -t -c 'SELECT Timescaledb_pre_restore();'
pg_restore -d realtime-dashboard -Fc /controller/realtime-dashboard.bak
psql -d realtime-dashboard -t -c 'SELECT Timescaledb_post_restore();'
```

#### 5. Re-enable services that write to the databases

Re-enable all services that were disabled in the _Disable services that write to databases_ section above (see [Disable database changes](#disable-data-ingestion)).

### Post-migration cleanup

Perform these steps when the database migration is successful.

1. Remove any database passwords from the `secrets.yaml` file. An easy way to do this is to compare your `secrets.yaml`
   file to `base-secrets.yaml` file and remove any entries not present in `base-secrets.yaml`.

2. Turn off any deployed legacy database service. For this update the `production.yaml` file like so:

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

For esthetics, you can also remove any configuration passed under any of these sections. For instance, remove any of
indicated lines:

```yaml
grafana_metrics_timescaledb:
  _install: true
  _extra_timeout: 210
  replicaCount: 1
  postgresql: <-- remove this section
    replication:
      enable: false
      applicationName: radarGrafanaMetrics
    auth:
      database: grafana-metrics
    primary:
      persistence:
        size: 8Gi
```

3. Remove the _mods/migration/1.3.0.yaml_ file reference from the `environments.yaml` file.

4. Update the deployment:

```shell
helmfile sync
```

5. Remove any _pvc_ resource on the Kubernetes cluster associated with the old databases.

   Important: this step will permanently delete the data! Only perform this step when sure the migration completed
   successfully.

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

4. Update Timescaledb database configuration:

- Rename `timescaledb_username` to `grafana_metrics_db_username`
- Remove `grafana_metrics_username` and `timescaledb_db_name` variables.
- When using _realtime-dashboard_, add `realtime_dashboard_db_username` that points to the current value of
  `timescaledb_username`.

5. Refactor the Timescaledb database configuration:

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

6. Add `postgresql:` indent to Postgresql and Timescaledb related database configurations like so:

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
data_dashboard_db_password: <same password as Timescaledb_password>
realtime_dashboard_db_password: <same password as Timescaledb_password>
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
kubectl delete secrets Timescaledb-postgresql
kubectl delete statefulsets Timescaledb-postgresql
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

## Supporting tasks

### Disable data ingestion

To disable data ingestion, temporarily disable _gateway_, _rest-connector_ (Fitit, Garmin, Oura, ...) and
_upload-connector_ services in the `production.yaml` file:
This is accomplished by redeploying while setting the `replicaCount` to `0` for these services. For instance:

```yaml
radar_gateway:
  replicaCount: 0
...
radar_fitbit_connector:
  replicaCount: 0
...
radar_oura_connector:
  replicaCount: 0
...
radar_upload_source_connector:
  replicaCount: 0
...
```

Followed by:

```shell
helmfile sync
```

BEWARE: the list of services that is scaled down in this way may vary depending on the services that are used in your
RADAR-Kubernetes deployment.The example above is not exhaustive and will not be updated with future versions of
RADAR-Kubernetes.

In order to re-enable data ingestion, set the `replicaCount` back to the original value and redeploy.

### Disable database changes

In prevent databases from receiving any updates (e.g., when performing an _off-line_ database update), we should disable
all services that perform write operations.
This is accomplished by redeploying while setting the `replicaCount` to `0` for these services. At the moment of this writing this would include:

```yaml
management_portal:
  replicaCount: 0
...
radar_appserver:
  replicaCount: 0
...
radar_appconfig:
  replicaCount: 0
...
radar_rest_sources_backend:
  replicaCount: 0
...
radar_upload_connect_backend:
  replicaCount: 0
...
kratos:
  replicaCount: 0
...
hydra:
  replicaCount: 0
...
radar_jdbc_connector_grafana:
  replicaCount: 0
...
radar_jdbc_connector_data_dashboard_backend:
  replicaCount: 0
...
radar_jdbc_connector_realtime_dashboard:
  replicaCount: 0
...
```

followed by:

```shell
helmfile sync
```

BEWARE: the list of services that is scaled down in this way may vary depending on the services that are used in your
RADAR-Kubernetes deployment.The example above is not exhaustive and will not be updated with future versions of
RADAR-Kubernetes.

In order to re-enable services, set the `replicaCount` back to the original value and redeploy.

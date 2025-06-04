# Configuration for external PostgreSQL databases

<!-- TOC -->
* [Configuration for external PostgreSQL databases](#configuration-for-external-postgresql-databases)
  * [Postgresql database configuration](#postgresql-database-configuration)
    * [Kratos and Hydra configuration](#kratos-and-hydra-configuration)
  * [TimescaleDB configuration](#timescaledb-configuration)
<!-- TOC -->

It is possible to use _Postgresql and TimescaleDB databases that are not provided by the default RADAR-base Kubernetes
deployment.
We refer to these as _external_ databases. This document provides some guidance on how to configure external databases.
In short, this is done by overriding the default values in the RADAR-base helm charts with the values of the external
databases via helmfile.

This document assumes sufficient knowledge of helm and helmfile to be able to apply the changes, as well as sufficient
expertise with Kubernetes, for instance to create custom password secrets.

## Postgresql database configuration

Add the following configuration in your `production.yaml` file (for more information on configuration of the :

```yaml
management_portal:
  postgres:
    # -- host of the postgres db
    host: external-postgres-host
    # -- post of the postgres db
    port: 5432
    # -- database name
    database: external-management-portal-db
    # -- user to connect to the database
    user: external-db-user

radar_appserver:
  postgres:
    host: external-postgres-host
    port: 5432
    database: external-radar-appserver-db
    user: external-db-user

radar_appconfig:
  jdbc:
    # -- URI of the postgres db (e.g. jdbc:postgresql://external-postgres-host:5432/external-radar-appconfig)
    url: external-postgres-uri
    user: external-db-user

radar_kratos:
  postgres:
    host: external-postgres-host
    port: 5432
    database: external-radar-kratos-db
    user:  external-db-user

radar_hydra:
  postgres:
    host: external-postgres-host
    port: 5432
    database: external-radar-hydra-db
    user:  external-db-user

radar_rest_sources_backend:
  postgres:
    host: external-postgres-host
    port: 5432
    database: external-radar-rest-sources-backend-db
    user: external-db-user

radar_upload_connect_backend:
  postgres:
    host: external-postgres-host
    port: 5432
    database: external-radar-upload-connect-backend-db
    user: external-db-user
```

Add the password to the `secrets.yaml` file:

```yaml
management_portal:
  postgres:
    # -- password to connect to the database
    password: external-db-password
radar_appserver:
  postgres:
    password: external-db-password
radar_appconfig:
  jdbc:
    password: external-db-password
radar_rest_sources_backend:
  postgres:
    password: external-db-password
radar_upload_connect_backend:
  postgres:
    password: external-db-password
```

### Kratos and Hydra configuration

For Ory Kratos and Hydra, Kubernetes Secret resources should be created according the following specifications:

```yaml
apiVersion: v1
data:
  uri: postgresql://external-db-user:external-db-password@external-postgres-host:external-db-port/kratos-db-name
kind: Secret
metadata:
  annotations:
    meta.helm.sh/release-name: radar-cloudnative-postgresql
    meta.helm.sh/release-namespace: default
  labels:
    app.kubernetes.io/managed-by: Helm
    cnpg.io/reload: "true"
  name: radar-cloudnative-postgresql-kratos
  namespace: default
type: kubernetes.io/basic-auth

---

apiVersion: v1
data:
  uri: postgresql://external-db-user:external-db-password@external-postgres-host:external-db-port/hydra-db-name
kind: Secret
metadata:
  annotations:
    meta.helm.sh/release-name: radar-cloudnative-postgresql
    meta.helm.sh/release-namespace: default
  labels:
    app.kubernetes.io/managed-by: Helm
    cnpg.io/reload: "true"
  name: radar-cloudnative-postgresql-hydra
  namespace: default
type: kubernetes.io/basic-auth
```

## TimescaleDB configuration

The configuration for TimescaleDB is similar to that of Postgresql, but it is applied under the _jdbc-connector_ sections:

```yaml
radar_jdbc_connector_grafana:
  jdbc:
    # -- URI of the TimescaleDB (e.g. jdbc:postgresql://external-postgres-host:5432/external-timescaledb)
    url: external-timescaledb-uri
    user: external-db-user

radar_jdbc_connector_data_dashboard_backend:
  jdbc:
    url: external-timescaledb-uri
    user: external-db-user

radar_jdbc_connector_realtime_dashboard:
  jdbc:
    url: external-timescaledb-uri
    user: external-db-user
```

Add the password to the `secrets.yaml` file:

```yaml
radar_jdbc_connector_grafana:
  jdbc:
    password: external-db-password

radar_jdbc_connector_data_dashboard_backend:
  jdbc:
    password: external-db-password

radar_jdbc_connector_realtime_dashboard:
  jdbc:
    password: external-db-password
```


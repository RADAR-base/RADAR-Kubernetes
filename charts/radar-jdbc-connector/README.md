# radar-jdbc-connector

![Version: 0.1.1](https://img.shields.io/badge/Version-0.1.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.2.1](https://img.shields.io/badge/AppVersion-1.2.1-informational?style=flat-square)

A Helm chart for RADAR-base JDBC Kafka connector. This is a fork of the Kafka JDBC connector which allows data from topics to be imported into JDBC databases (including TimescaleDB databases which is used in the dashboard pipeline). For more details see "https://github.com/RADAR-base/RADAR-JDBC-Connector".

**Homepage:** <https://radar-base.org>

## Maintainers

| Name                | Email                | Url |
| ------------------- | -------------------- | --- |
| Keyvan Hedayati     | keyvan@thehyve.nl    |     |
| Joris Borgdorff     | joris@thehyve.nl     |     |
| Nivethika Mahasivam | nivethika@thehyve.nl |     |

## Source Code

- <https://github.com/RADAR-base/RADAR-JDBC-Connector>

## Prerequisites

- Kubernetes 1.17+
- Kubectl 1.17+
- Helm 3.1.0+

## Values

| Key                       | Type   | Default                                                                                                                                                             | Description                                                                                                              |
| ------------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| replicaCount              | int    | `1`                                                                                                                                                                 | Number of JDBC connector replicas to deploy                                                                              |
| image.repository          | string | `"radarbase/radar-jdbc-connector"`                                                                                                                                  | JDBC connector image repository                                                                                          |
| image.tag                 | string | `"latest"`                                                                                                                                                          | JDBC connector image tag (immutable tags are recommended) Overrides the image tag whose default is the chart appVersion. |
| image.pullPolicy          | string | `"IfNotPresent"`                                                                                                                                                    | JDBC connector image pull policy                                                                                         |
| imagePullSecrets          | list   | `[]`                                                                                                                                                                | Docker registry secret names as an array                                                                                 |
| nameOverride              | string | `""`                                                                                                                                                                | String to partially override radar-jdbc-connector.fullname template with a string (will prepend the release name)        |
| fullnameOverride          | string | `""`                                                                                                                                                                | String to fully override radar-jdbc-connector.fullname template with a string                                            |
| resources.limits          | object | `{"cpu":"1000m"}`                                                                                                                                                   | CPU/Memory resource limits                                                                                               |
| resources.requests        | object | `{"cpu":"100m","memory":"128Mi"}`                                                                                                                                   | CPU/Memory resource requests                                                                                             |
| nodeSelector              | object | `{}`                                                                                                                                                                | Node labels for pod assignment                                                                                           |
| tolerations               | list   | `[]`                                                                                                                                                                | Toleration labels for pod assignment                                                                                     |
| affinity                  | object | `{}`                                                                                                                                                                | Affinity labels for pod assignment                                                                                       |
| zookeeper                 | string | `"cp-zookeeper-headless:2181"`                                                                                                                                      | URI of Zookeeper instances of the cluster                                                                                |
| kafka                     | string | `"PLAINTEXT://cp-kafka-headless:9092"`                                                                                                                              | URI of Kafka brokers of the cluster                                                                                      |
| kafka_num_brokers         | string | `"3"`                                                                                                                                                               | Number of Kafka brokers. This is used to validate the cluster availability at connector init.                            |
| schema_registry           | string | `"http://cp-schema-registry:8081"`                                                                                                                                  | URL of the Kafka schema registry                                                                                         |
| topics                    | string | `"android_phone_relative_location, android_phone_battery_level, connect_upload_altoida_summary, connect_fitbit_intraday_heart_rate, connect_fitbit_intraday_steps"` | Comma-separated list of topics the connector will read from and ingest into the database                                 |
| timescaledb.host          | string | `"timescaledb-postgresql-headless"`                                                                                                                                 | Host of the TimescaleDB database                                                                                         |
| timescaledb.username      | string | `"grafana"`                                                                                                                                                         | TimescaleDB database username                                                                                            |
| timescaledb.password      | string | `"password"`                                                                                                                                                        | TimescaleDB database password                                                                                            |
| timescaledb.database_name | string | `"grafana-metrics"`                                                                                                                                                 | TimescaleDB database name                                                                                                |

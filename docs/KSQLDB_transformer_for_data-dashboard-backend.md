# kafka-data transformer (KSQLDB) for data-dashboard-backend service

Reference: https://docs.ksqldb.io/

The data-dashboard-backend service uses data from Kafka topics to the _observation_ table in the RADAR-base Data
Dashboard backend service database. The data in the Kafka topics is transformed by the KSQLDB Kafka data transformer to
be imported into the _observation_ table.

The KSQLDB Kafka data transformer is able to register Consumer/Producers to Kafka that transform data in a topic and
publish the results to another topic.

The provided KSQLDB _questionnaire_response_observations.sql_ and _questionnaire_app_events_observation.sql_ SQL files
transform, respectively, the _questionnaire_response_ and _questionnaire_app_event_ topics and publish the data to the
_ksql_observations_ topic. The _ksql_observations_ topic is consumed by the Kafka-JDBC-connector used for the by the
RADAR-base Data Dashboard backend service (see: [20-data-dashboard.yaml](../../helmfile.d/20-dashboard.yaml)).

When transformation of other topics is required, new SQL files can be added to this directory. These new files should be
referenced in the _kafka_data_transformer -> ksql -> queries_ section of the `etc/base.yaml.gotmpl` file. New KSQLDB SQL
files should transform towards the following format of the _ksql_observations_ topic:

```
    TOPIC KEY:
      PROJECT: the project identifier
      SOURCE: the source identifier
      SUBJECT: the subject/study participant identifier
    TOPIC VALUE:
      TOPIC: the topic identifier
      CATEGORY: the category identifier (optional)
      VARIABLE: the variable identifier
      DATE: the date of the observation
      END_DATE: the end date of the observation (optional)
      TYPE: the type of the observation (STRING, STRING_JSON, INTEGER, DOUBLE)
      VALUE_TEXTUAL: the textual value of the observation (optional, must be set when VALUE_NUMERIC is NULL)
      VALUE_NUMERIC: the numeric value of the observation (optional, must be set when VALUE_TEXTUAL is NULL)
```

New messages are added to the _ksql_observations_ topic by inserting into the _observations_ stream (
see [_base_observations_stream.sql](_base_observations_stream.sql)):

```
INSERT INTO observations
SELECT
...
PARTITION BY q.projectId, q.userId, q.sourceId
EMIT CHANGES;
```
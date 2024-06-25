# Kafka-data transformer (ksql-server) for data-dashboard-backend service

Reference: https://docs.ksql-server.io/

The data-dashboard-backend service uses data derived from Kafka topics imported into the _observation_ table
data-dashboard-backend service database. The data in the Kafka topics is transformed by the ksql-server data transformer
to be imported into the _observation_ table.

The ksql-server data transformer is able to register Consumer/Producers to Kafka that transform data in a topic and
publish the results to another topic.

The provided ksql-server _questionnaire_response_observations.sql_ and _questionnaire_app_events_observation.sql_ SQL files
transform, respectively, the _questionnaire_response_ and _questionnaire_app_event_ topics and publish the data to the
_ksql_observations_ topic. The _ksql_observations_ topic is consumed by the radar-jdbc-connector service deployed for the
data-dashboard-backend service (see: [20-data-dashboard.yaml](../helmfile.d/20-dashboard.yaml)).

When transformation of other topics is required, new SQL files can be added to this directory. These new files should be
referenced in the _cp-ksql-server_ -> ksql -> queries_ section of the `etc/base.yaml.gotmpl` file. New ksql-server SQL
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
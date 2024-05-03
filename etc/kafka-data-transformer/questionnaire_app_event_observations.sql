-- This is the ksqlDB script that reads the questionnaire_response and questionnaire_app_event topics and writes
-- the observations to a new topic named 'ksql_observations_<topic>'. More topics transformations can be added.
-- Every topic must be transformed to the common format:
-- KEY:
--   PROJECT: the project identifier
--   SOURCE: the source identifier
--   SUBJECT: the subject/study participant identifier
-- VALUE:
--   TOPIC: the topic identifier
--   CATEGORY: the category identifier (optional)
--   VARIABLE: the variable identifier
--   DATE: the date of the observation
--   END_DATE: the end date of the observation (optional)
--   TYPE: the type of the observation (STRING, STRING_JSON, INTEGER, DOUBLE)
--   VALUE_TEXTUAL: the textual value of the observation (optional, must be set when VALUE_NUMERIC is NULL)
--   VALUE_NUMERIC: the numeric value of the observation (optional, must be set when VALUE_TEXTUAL is NULL)

SET 'auto.offset.reset' = 'earliest';

-- * -- * -- topic: QUESTIONNAIRE_APP_EVENT -- * -- * --

CREATE STREAM questionnaire_app_event (
    projectId VARCHAR KEY,    -- 'KEY' means that this field is part of the kafka message key
    userId VARCHAR KEY,
    sourceId VARCHAR KEY,
    questionnaireName VARCHAR,
    eventType VARCHAR,
    time DOUBLE,
    metadata MAP<VARCHAR, VARCHAR>
) WITH (
    kafka_topic = 'questionnaire_app_event',
    partitions = 3,
    format = 'avro'
);

CREATE STREAM questionnaire_app_event_observations
WITH (
    kafka_topic = 'ksql_observations_questionnaire_app_event',
    partitions = 3,
    format = 'avro'
)
AS SELECT
    FROM_UNIXTIME(CAST(q.time * 1000 AS BIGINT)) as DATE,
    q.projectId AS PROJECT,
    q.userId AS SUBJECT,
    q.sourceId AS SOURCE,
    'questionnaire_app_event' as `TOPIC`,
    CAST(NULL as VARCHAR) as CATEGORY,
    CAST(NULL as TIMESTAMP) as END_DATE,
    q.questionnaireName as VARIABLE,
    'STRING_JSON' as TYPE,
    CAST(NULL as DOUBLE) as VALUE_NUMERIC,
    TO_JSON_STRING(q.metadata) as VALUE_TEXTUAL
FROM questionnaire_app_event q
PARTITION BY q.projectId, q.userId, q.sourceId -- this sets the fields in the kafka message key
EMIT CHANGES;

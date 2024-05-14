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
    kafka_topic = 'ksql_observations',
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

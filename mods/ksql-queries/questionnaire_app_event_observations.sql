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
    partitions = 1,
    format = 'avro'
);

INSERT INTO observations
WITH (QUERY_ID='questionnaire_app_event_observations')
SELECT
    q.projectId AS PROJECT,
    q.userId AS SUBJECT,
    q.sourceId AS SOURCE,
    'questionnaire_app_event' as TOPIC_NAME,
    q.questionnaireName as CATEGORY,
    q.eventType as VARIABLE,
    FROM_UNIXTIME(CAST(q.time * 1000 AS BIGINT)) as OBSERVATION_TIME,
    CAST(NULL as TIMESTAMP) as OBSERVATION_TIME_END,
    'STRING_JSON' as TYPE,
    CAST(NULL as DOUBLE) as VALUE_NUMERIC,
    TO_JSON_STRING(q.metadata) as VALUE_TEXTUAL
FROM questionnaire_app_event q
PARTITION BY q.projectId, q.userId, q.sourceId -- this sets the fields in the kafka message key
EMIT CHANGES;

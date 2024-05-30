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

INSERT INTO observations
SELECT
    q.projectId AS PROJECT,
    q.userId AS SUBJECT,
    q.sourceId AS SOURCE,
    'questionnaire_app_event' as `TOPIC`,
    CAST(NULL as VARCHAR) as CATEGORY,
    q.questionnaireName as VARIABLE,
    FROM_UNIXTIME(CAST(q.time * 1000 AS BIGINT)) as DATE,
    CAST(NULL as TIMESTAMP) as END_DATE,
    'STRING_JSON' as TYPE,
    CAST(NULL as DOUBLE) as VALUE_NUMERIC,
    TO_JSON_STRING(q.metadata) as VALUE_TEXTUAL
FROM questionnaire_app_event q
PARTITION BY q.projectId, q.userId, q.sourceId -- this sets the fields in the kafka message key
EMIT CHANGES;

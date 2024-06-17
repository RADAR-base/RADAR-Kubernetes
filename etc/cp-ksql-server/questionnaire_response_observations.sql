CREATE STREAM questionnaire_response (
    projectId VARCHAR KEY, -- 'KEY' means that this field is part of the kafka message key
    userId VARCHAR KEY,
    sourceId VARCHAR KEY,
    time DOUBLE,
    timeCompleted DOUBLE,
    timeNotification DOUBLE,
    name VARCHAR,
    version VARCHAR,
    answers ARRAY<STRUCT<questionId VARCHAR, value STRUCT<int INT, string VARCHAR, double DOUBLE>, startTime DOUBLE, endTime DOUBLE>>
) WITH (
    kafka_topic = 'questionnaire_response',
    partitions = 3,
    format = 'avro'
);

CREATE STREAM questionnaire_response_exploded
AS SELECT
    EXPLODE(TRANSFORM(q.answers, a => a->questionId)) as VARIABLE,
    FROM_UNIXTIME(CAST(q.time * 1000 AS BIGINT)) as OBSERVATION_TIME,
    q.projectId,
    q.userId,
    q.sourceId,
    'questionnaire_response' as TOPIC_NAME,
    q.name as CATEGORY,
    CAST(NULL as TIMESTAMP) as OBSERVATION_TIME_END,
    -- WARNING!!! The cast from VARCHAR (string) to DOUBLE will throw an JAVA exception if the string is not a number.
    -- This does not mean that the message will be lost. The value will be present in the VALUE_TEXTUAL_OPTIONAL field.
    EXPLODE(TRANSFORM(q.answers, a => COALESCE(a->value->double, CAST(a->value->int as DOUBLE), CAST(a->value->string as DOUBLE)))) as VALUE_NUMERIC,
    EXPLODE(TRANSFORM(q.answers, a => CASE
        WHEN a->value->int IS NOT NULL THEN 'INTEGER'
        WHEN a->value->double IS NOT NULL THEN 'DOUBLE'
        ELSE NULL
    END)) as TYPE,
    -- Note: When cast to double works for the string value, the VALUE_TEXTUAL_OPTIONAL will also be set.
    EXPLODE(TRANSFORM(q.answers, a => a->value->string)) as VALUE_TEXTUAL_OPTIONAL
FROM questionnaire_response q
EMIT CHANGES;

INSERT INTO observations
WITH (QUERY_ID='questionnaire_response_observations')
SELECT
   q.projectId as PROJECT,
   q.sourceId as SOURCE,
   q.userId as SUBJECT,
   TOPIC_NAME, CATEGORY, VARIABLE, OBSERVATION_TIME, OBSERVATION_TIME_END,
   CASE
       WHEN TYPE IS NULL AND VALUE_NUMERIC IS NOT NULL THEN 'DOUBLE' -- must have been derived from a string cast
       WHEN TYPE IS NULL AND VALUE_NUMERIC IS NULL THEN 'STRING'
       ELSE TYPE                                                     -- keep the original type when TYPE is not NULL
   END as TYPE,
   VALUE_NUMERIC,
    CASE
        WHEN VALUE_NUMERIC IS NOT NULL THEN NULL                     -- When cast to double has worked for the string value, set VALUE_TEXTUAL to NULL.
        ELSE VALUE_TEXTUAL_OPTIONAL
    END as VALUE_TEXTUAL
FROM questionnaire_response_exploded q
PARTITION BY q.projectId, q.userId, q.sourceId -- this sets the fields in the kafka message key
EMIT CHANGES;

-- TODO: exploding the 'select:' questions is not yet fully designed.
-- I keep the code here for future reference.
-- Multi-select questionnaire questions are stored as a single 'value' string with the
-- names of the selected options separated by comma's. Multiselect questions are prefixed
-- by 'select:' in the questionId.
-- When 'questionId' is like 'select:%' create a new stream with the select options.
-- The options in the value field split commas and added as separate VARIABLE records.
-- The VALUE_NUMERIC is set to 1 and VALUE_TEXTUAL is set to NULL.
-- INSERT INTO observations
-- SELECT
--     EXPLODE(SPLIT(VALUE_TEXTUAL, ',')) as VARIABLE,
--     PROJECT, SOURCE, SUBJECT, TOPIC_NAME, CATEGORY, OBSERVATION_TIME, OBSERVATION_TIME_END,
--     'INTEGER' as TYPE,
--     CAST(1 as DOUBLE) VALUE_NUMERIC,
--     CAST(NULL as VARCHAR) as VALUE_TEXTUAL
-- FROM questionnaire_response_observations
-- WHERE
--  VARIABLE IS NOT NULL
--  AND VARIABLE LIKE 'select:%'
--  AND VALUE_TEXTUAL IS NOT NULL
--  AND VALUE_TEXTUAL != ''
-- PARTITION BY SUBJECT, PROJECT, SOURCE
-- EMIT CHANGES;

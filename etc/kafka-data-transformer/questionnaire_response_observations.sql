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

-- * -- * -- topic: QUESTIONNAIRE_RESPONSE -- * -- * --

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
    FROM_UNIXTIME(CAST(q.time * 1000 AS BIGINT)) as DATE,
    q.projectId,
    q.userId,
    q.sourceId,
    'questionnaire_response' as `TOPIC`,
    q.name as CATEGORY,
    CAST(NULL as TIMESTAMP) as END_DATE,
    CASE
        WHEN a->value->int IS NOT NULL THEN 'INTEGER'
        WHEN a->value->double IS NOT NULL THEN 'DOUBLE'
        ELSE NULL
    END as TYPE,
    -- WARNING!!! The cast from VARCHAR (string) to DOUBLE will throw an JAVA exception if the string is not a number.
    -- This does not mean that the message will be lost. The value will be present in the VALUE_TEXTUAL_OPTIONAL field.
    EXPLODE(TRANSFORM(q.answers, a => COALESCE(a->value->double, CAST(a->value->int as DOUBLE), CAST(a->value->string as DOUBLE)))) as VALUE_NUMERIC,
    -- Note: When cast to double works for the string value, the VALUE_TEXTUAL_OPTIONAL will also be set.
    EXPLODE(TRANSFORM(q.answers, a => a->value->string)) as VALUE_TEXTUAL_OPTIONAL
FROM questionnaire_response q
EMIT CHANGES;

CREATE STREAM questionnaire_response_observations
WITH (
    kafka_topic = 'ksql_observations',
    partitions = 3,
    format = 'avro'
)
AS SELECT
   q.projectId as PROJECT,
   q.sourceId as SOURCE,
   q.userId as SUBJECT,
   `TOPIC`, CATEGORY, VARIABLE, DATE, END_DATE,
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
-- CREATE STREAM questionnaire_response_observations_select
-- WITH (
--     kafka_topic = 'ksql_observations',
--     partitions = 3,
--     format = 'avro'
-- )
-- AS SELECT
--     EXPLODE(SPLIT(VALUE_TEXTUAL, ',')) as VARIABLE,
--     PROJECT, SOURCE, SUBJECT, `TOPIC`, CATEGORY, DATE, END_DATE,
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

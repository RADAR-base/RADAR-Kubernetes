SET 'auto.offset.reset' = 'earliest';

-- Define stream from questionnaire_response topic
-- This way it can be read by ksqlDB
CREATE STREAM questionnaire_response (
    projectId VARCHAR KEY,
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

-- Convert all questionnaires answers to single observation records
CREATE STREAM outcomes_observations_unindexed
WITH (
    kafka_topic = 'ksql_outcomes_observations_unindexed',
    partitions = 3,
    format = 'avro'
)
AS SELECT
    EXPLODE(TRANSFORM(q.answers, a => a->questionId)) as `variable`,
    FROM_UNIXTIME(CAST(q.time * 1000 AS BIGINT)) as `date`,
    q.userId as `subject_id`,
    CAST(NULL as TIMESTAMP) as `end_date`,
    -- if present, take the int response, otherwise try to convert if the string response to a numeric value.
    EXPLODE(TRANSFORM(q.answers, a => COALESCE(a->value->int, CAST(a->value->string as INT)))) as `value_numeric`,
    EXPLODE(TRANSFORM(q.answers, a => a->value->string)) as `value_textual`
FROM questionnaire_response q
PARTITION BY q.userId
EMIT CHANGES;

-- Read select statements.
CREATE STREAM outcomes_observations_unindexed_select
WITH (
    kafka_topic = 'ksql_outcomes_observations_unindexed',
    partitions = 3,
    format = 'avro'
)
AS SELECT
    EXPLODE(SPLIT(`value_textual`, ',')) as `variable`,
    `subject_id`,
    `date`,
    `end_date`,
    1 as `value_numeric`,
    CAST(NULL as VARCHAR) as `value_textual`
FROM outcomes_observations_unindexed
WHERE `variable` LIKE 'select:%'
  AND `value_textual` IS NOT NULL
  AND `value_textual` != ''
EMIT CHANGES;

-- Read the correct variable ID from the database and join it with the question ID.
-- Non-existing variables will not be joined so questionIds that are not used in H2O
-- will not be converted.
CREATE STREAM outcomes_observations
WITH (
    kafka_topic = 'ksql_outcomes_observations',
    partitions = 3,
    format = 'avro'
)
AS SELECT
    o.`variable`,
    v.id as `variable_id`,
    o.`subject_id`,
    o.`date`,
    o.`end_date`,
    o.`value_numeric`,
    o.`value_textual`
FROM outcomes_observations_unindexed o
-- this will force o.variable to be the record key
JOIN outcomes_variable v ON o.`variable` = v.name
EMIT CHANGES;

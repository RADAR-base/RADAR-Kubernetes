SET 'auto.offset.reset' = 'earliest';

-- Register the 'ksql_observations' topic (is created when not exists).
CREATE STREAM observations (
    PROJECT VARCHAR KEY,    -- 'KEY' means that this field is part of the kafka message key
    SUBJECT VARCHAR KEY,
    SOURCE VARCHAR KEY,
    TOPIC_NAME VARCHAR,
    CATEGORY VARCHAR,
    VARIABLE VARCHAR,
    OBSERVATION_TIME TIMESTAMP,
    OBSERVATION_TIME_END TIMESTAMP,
    TYPE VARCHAR,
    VALUE_NUMERIC DOUBLE,
    VALUE_TEXTUAL VARCHAR
) WITH (
    kafka_topic = 'ksql_observations',
    partitions = 3,
    format = 'avro'
);

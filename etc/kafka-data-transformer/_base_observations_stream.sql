SET 'auto.offset.reset' = 'earliest';

-- Register the 'observations' stream (is created when not exists).
CREATE STREAM observations (
    PROJECT VARCHAR KEY,    -- 'KEY' means that this field is part of the kafka message key
    SUBJECT VARCHAR KEY,
    SOURCE VARCHAR KEY,
    `TOPIC` VARCHAR,
    CATEGORY VARCHAR,
    VARIABLE VARCHAR,
    DATE TIMESTAMP,
    END_DATE TIMESTAMP,
    TYPE VARCHAR,
    VALUE_NUMERIC DOUBLE,
    VALUE_TEXTUAL VARCHAR
) WITH (
    kafka_topic = 'ksql_observations',
    partitions = 3,
    format = 'avro'
);

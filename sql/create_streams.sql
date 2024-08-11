DROP STREAM IF EXISTS brawl_stream;
DROP STREAM IF EXISTS not_on_list_stream;
DROP STREAM IF EXISTS accident_stream;
DROP STREAM IF EXISTS dirty_table_stream;
DROP STREAM IF EXISTS broken_items_stream;
DROP STREAM IF EXISTS bad_food_stream;
DROP STREAM IF EXISTS music_stream;
DROP STREAM IF EXISTS feeling_ill_stream;
DROP STREAM IF EXISTS bride_stream;
DROP STREAM IF EXISTS groom_stream;

CREATE STREAM brawl_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='brawl',
  VALUE_FORMAT='json'
);

CREATE STREAM not_on_list_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='not_on_list',
  VALUE_FORMAT='json'
);

CREATE STREAM accident_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='accident',
  VALUE_FORMAT='json'
);

CREATE STREAM dirty_table_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='dirty_table',
  VALUE_FORMAT='json'
);

CREATE STREAM broken_items_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='broken_items',
  VALUE_FORMAT='json'
);

CREATE STREAM bad_food_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='bad_food',
  VALUE_FORMAT='json'
);

CREATE STREAM music_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='music',
  VALUE_FORMAT='json'
);

CREATE STREAM feeling_ill_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='feeling_ill',
  VALUE_FORMAT='json'
);

CREATE STREAM bride_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='bride',
  VALUE_FORMAT='json'
);

CREATE STREAM groom_stream (
  id STRING,
  event_id   STRING,
  event_type STRING,
  event_time STRING,
  priority STRING,
  description STRING,
  status STRING
) WITH (
  KAFKA_TOPIC='groom',
  VALUE_FORMAT='json'
);

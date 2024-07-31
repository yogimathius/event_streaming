CREATE STREAM brawl_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='brawl',
  VALUE_FORMAT='json'
);

CREATE STREAM not_on_list_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='not_on_list',
  VALUE_FORMAT='json'
);

CREATE STREAM accident_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='accident',
  VALUE_FORMAT='json'
);

CREATE STREAM dirty_table_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='dirty_table',
  VALUE_FORMAT='json'
);

CREATE STREAM broken_items_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='broken_items',
  VALUE_FORMAT='json'
);

CREATE STREAM bad_food_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='bad_food',
  VALUE_FORMAT='json'
);

CREATE STREAM music_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='music',
  VALUE_FORMAT='json'
);

CREATE STREAM feeling_ill_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='feeling_ill',
  VALUE_FORMAT='json'
);

CREATE STREAM bride_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='bride',
  VALUE_FORMAT='json'
);

CREATE STREAM groom_stream (
  id STRING,
  event_time TIMESTAMP
) WITH (
  KAFKA_TOPIC='groom',
  VALUE_FORMAT='json'
);

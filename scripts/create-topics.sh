#!/bin/bash

# Define your topics here
KAFKA_TOPICS=(
  "brawl"
  "not_on_list"
  "accident"
  "dirty_table"
  "broken_items"
  "bad_food"
  "music"
  "feeling_ill"
  "bride"
  "groom"
)

# Loop over each topic and create it
for TOPIC in "${KAFKA_TOPICS[@]}"; do
  /usr/bin/kafka-topics --create --if-not-exists --topic "$TOPIC" --replication-factor 1 --partitions 3 --bootstrap-server kafka:29092
done

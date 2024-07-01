#!/bin/bash

# Define your topics here
KAFKA_TOPICS=(
  "topic1"
  "topic2"
  "topic3"
)

# Loop over each topic and create it
for TOPIC in "${KAFKA_TOPICS[@]}"; do
  /opt/bitnami/kafka/bin/kafka-topics.sh --create --if-not-exists --topic "$TOPIC" --replication-factor 1 --partitions 3 --bootstrap-server kafka:9092
done

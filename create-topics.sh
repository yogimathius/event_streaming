#!/bin/bash

# Define your topics here
KAFKA_TOPICS=(
  "brawl"
  "not-on-list"
  "person-fell"
  "injured-kid"
  "dirty-table"
  "broken-glass-clean-up"
  "bad-food"
  "music-too-loud"
  "music-too-low"
  "feeling-ill-catering"
  "missing-rings"
  "missing-bride"
  "missing-groom"
  "broken-glass-waiters"
  "person-fell-waiters"
  "injured-kid-waiters"
  "feeling-ill-waiters"
)

# Loop over each topic and create it
for TOPIC in "${KAFKA_TOPICS[@]}"; do
  /opt/bitnami/kafka/bin/kafka-topics.sh --create --if-not-exists --topic "$TOPIC" --replication-factor 1 --partitions 3 --bootstrap-server kafka:9092
done

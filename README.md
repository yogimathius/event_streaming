# Event Streaming (Wedding Simulation)

Distributed event-streaming simulation using Go, Kafka, Rust workers, and a Python event generator.

## Purpose
- Model a high-volume event-processing pipeline with multiple worker teams and priority levels.
- Validate topic-based routing and asynchronous processing behavior in a realistic simulation scenario.

## Current Implementation
- Go producer/API (`go-server/`) accepts events and publishes to Kafka.
- Rust consumer workers (`rust-consumer/`) process messages by topic.
- Python simulator (`simulator/simulator.py`) generates randomized or dataset-driven events.
- Infrastructure orchestration via Docker Compose (`docker-compose-services.yml`, `docker-compose.yml`).

## Interfaces
- HTTP endpoint exposed by Go service:
- `POST /message` (default local target in simulator: `http://localhost:8080/message`)
- Kafka topics created by `scripts/create-topics.sh`:
- `event_created`
- `brawl`
- `not_on_list`
- `accident`
- `dirty_table`
- `broken_items`
- `bad_food`
- `music`
- `feeling_ill`
- `bride`
- `groom`

## Running Locally
1. Create topics script executable:
- `chmod +x scripts/create-topics.sh`
2. Start Kafka/Zookeeper stack:
- `docker compose -f docker-compose-services.yml up -d`
3. Start application services:
- `docker compose up --build`
4. Run simulator:
- `python3 simulator/simulator.py`
or with dataset:
- `python3 simulator/simulator.py assets/dataset_1.json`

## Current Status
- Functional simulation pipeline with producer, broker, and consumers is present.
- This repository is a scenario-driven event-streaming demo, not a general-purpose production platform.

## Next Steps
- Add automated integration checks for end-to-end message flow.
- Add runbooks for failure modes (broker restart, consumer lag, malformed payloads).
- Remove/ignore local env files and document required variables in `.env.example`.

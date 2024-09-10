# Decision Log

## Overview

This document records the significant decisions made during the development of the project. Each decision includes the context, options considered, the final decision, and any relevant notes.

---

## Decision 1: Prototyped Kafka Consumer/Producer Logic in Python

**Date**: 2024-03-15

### Context

- Initial prototyping of the consumer/producer logic was done using Python to quickly validate the concept.

### Options Considered

1. **Prototype in Python**: Use Python for quick prototyping.
2. **Prototype in Another Language**: Use another language for prototyping.

### Decision

- Prototype the consumer/producer logic in Python.

### Consequences

- Quick validation of the concept with minimal development time.

### Notes

- Python was chosen for its simplicity and rapid development capabilities.

---

## Decision 2: Prototyped Consumer/Producer Logic in Go

**Date**: 2024-04-01

### Context

- After validating the concept in Python, the consumer/producer logic was prototyped in Go as an effort to learn Golang.

### Options Considered

1. **Prototype in Go**: Use Go as an effort to learn Golang.
2. **Continue with Python**: Continue using Python.

### Decision

- Prototype the consumer/producer logic in Go.

### Consequences

- Improved performance and concurrency handling. Proficiency in Golang increased.

### Notes

- Go was chosen for its performance benefits and efficient concurrency model, as well as to learn a new language.

---

## Decision 3: Prototyped Go Server with Gin Framework

**Date**: 2024-04-15

### Context

- A Go server was prototyped using the Gin framework to handle incoming requests.

### Options Considered

1. **Use Gin Framework**: Use the Gin framework for the Go server.
2. **Use Another Framework**: Use another framework for the Go server.

### Decision

- Prototype the Go server with the Gin framework.

### Consequences

- Efficient handling of HTTP requests with minimal overhead.

### Notes

- Gin was chosen for its simplicity and performance.

---

## Decision 4: Added Producer to Go Server

**Date**: 2024-04-30

### Context

- A Kafka producer was added to the Go server to handle message production.

### Options Considered

1. **Add Kafka Producer**: Add a Kafka producer to the Go server.

### Decision

- Add a Kafka producer to the Go server.

### Consequences

- Improved integration and handling of message production.

### Notes

- This change helps in emitting messages to Kafka from the server.

---

## Decision 5: Prototyped Database Structure with Postgres

**Date**: 2024-05-10

### Context

- The database structure was prototyped using Postgres to store and manage data efficiently.

### Options Considered

1. **Use Postgres**: Use Postgres for the database.
2. **Use Another Database**: Use another database for the project.

### Decision

- Prototype the database structure with Postgres.

### Consequences

- Reliable and efficient data storage and management.

### Notes

- Postgres was chosen for its robustness and feature set, as well as familiarity with the technology.

---

## Decision 6: Added Debug Logging to Go Consumer

**Date**: 2024-05-20

### Context

- Debug logging was added to the Go consumer to facilitate troubleshooting and monitoring.

### Options Considered

1. **Add Debug Logging**: Add debug logging to the Go consumer.
2. **No Debug Logging**: Continue without debug logging.

### Decision

- Add debug logging to the Go consumer.

### Consequences

- Improved ability to troubleshoot and monitor the consumer.

### Notes

- This change helps in identifying and resolving issues more efficiently.

---

## Decision 7: Added Debug Intermediary Status Updates via Kafka in Go Consumer/Producer

**Date**: 2024-06-01

### Context

- Intermediary status updates were added via Kafka to track the progress of messages through the system.

### Options Considered

1. **Add Status Updates**: Add intermediary status updates via Kafka.
2. **No Status Updates**: Continue without status updates.

### Decision

- Add intermediary status updates via Kafka.

### Consequences

- Improved visibility into the message processing pipeline.

### Notes

- This change helps in monitoring the progress of messages and identifying bottlenecks.

---

## Decision 8: Implemented Python Simulator Using Mock Data Sets

**Date**: 2024-06-15

### Context

- A Python simulator was implemented using mock data sets to test the system under various scenarios.

### Options Considered

1. **Implement Simulator**: Implement a simulator using mock data sets.
2. **No Simulator**: Continue without a simulator.

### Decision

- Implement a Python simulator using mock data sets.

### Consequences

- Improved ability to test and validate the system under different conditions.

### Notes

- This change helps in ensuring the system's robustness and reliability.

---

## Decision 9: Set Up ksqlDB Container to Consume Streams

**Date**: 2024-07-01

### Context

- A ksqlDB container was set up to consume and process streams for real-time analytics.

### Options Considered

1. **Set Up ksqlDB**: Use ksqlDB for stream processing.
2. **Use Another Solution**: Use another solution for stream processing.

### Decision

- Set up a ksqlDB container to consume streams.

### Consequences

- Improved ability to perform real-time analytics on streaming data.

### Notes

- ksqlDB was chosen for its powerful stream processing capabilities.

---

## Decision 10: Prototyped Rust Consumer/Worker Logic with mspc Channel Crate

**Date**: 2024-07-15

### Context

- The consumer/worker logic was prototyped in Rust using the mspc channel crate for better performance and concurrency.

### Options Considered

1. **Prototype in Rust**: Use Rust for better performance and concurrency.
2. **Continue with Go**: Continue using Go.

### Decision

- Prototype the consumer/worker logic in Rust with the mspc channel crate.

### Consequences

- Improved performance and concurrency handling.

### Notes

- Rust was chosen for its performance benefits and efficient concurrency model.

---

## Decision 11: Prototyped Go Workers Using Async Go Routines

**Date**: 2024-08-01

### Context

- Go workers were prototyped using async go routines to handle tasks concurrently.

### Options Considered

1. **Use Async Go Routines**: Use async go routines for concurrency.
2. **Use Another Concurrency Model**: Use another concurrency model.

### Decision

- Prototype Go workers using async go routines.

### Consequences

- Improved concurrency handling and performance.

### Notes

- Go routines were chosen for their simplicity and efficiency.

---

## Decision 12: Completed Go Implementation of Producer/Consumer Logic

**Date**: 2024-08-15

### Context

- The implementation of the producer/consumer logic in Go was completed to provide a robust and efficient solution.

### Options Considered

1. **Complete Go Implementation**: Complete the implementation in Go.
2. **Switch to Another Language**: Switch to another language.

### Decision

- Complete the Go implementation of the producer/consumer logic.

### Consequences

- Provided a robust and efficient solution for message processing.

### Notes

- Go was chosen for its performance benefits and efficient concurrency model.

---

## Decision 13: Implemented Rust Threaded Workers on One Channel

**Date**: 2024-08-25

### Context

- Rust threaded workers were implemented on one channel to handle tasks concurrently.

### Options Considered

1. **Use One Channel**: Implement threaded workers on one channel.
2. **Use Multiple Channels**: Implement threaded workers on multiple channels.

### Decision

- Implement Rust threaded workers on one channel.

### Consequences

- Improved concurrency handling and performance.

### Notes

- This change helps in better managing and processing tasks.

---

## Decision 14: Implemented Rust Threaded Workers on Multiple Channels (High, Medium, Low Priority Queueing)

**Date**: 2024-09-01

### Context

- Rust threaded workers were implemented on multiple channels to handle tasks with different priorities.

### Options Considered

1. **Use One Channel**: Continue using one channel.
2. **Use Multiple Channels**: Implement threaded workers on multiple channels.

### Decision

- Implement Rust threaded workers on multiple channels (high, medium, low priority queueing).

### Consequences

- Improved handling of tasks with different priorities.

### Notes

- This change helps in better managing and processing tasks based on priority.

---

## Decision 15: Made Rust Consumer Scalable to Handle Dynamic Topics/Teams/Routines Through Env Vars

**Date**: 2024-09-05

### Context

- The Rust consumer was made scalable to handle dynamic topics, teams, and routines through environment variables.

### Options Considered

1. **Hardcode Configurations**: Continue hardcoding configurations.
2. **Use Env Vars**: Use environment variables for dynamic configurations.

### Decision

- Make the Rust consumer scalable to handle dynamic topics, teams, and routines through environment variables.

### Consequences

- Improved scalability and flexibility in managing configurations.

### Notes

- This change helps in better managing and scaling the system.

---

## Decision 16: Migrated All Consumer/Worker Logic to Rust and Removed Go Consumer/Worker Logic

**Date**: 2024-09-07

### Context

- All consumer/worker logic was migrated to Rust, and the Go consumer/worker logic was removed to standardize on Rust.

### Options Considered

1. **Keep Mixed Logic**: Continue using a mix of Go and Rust.
2. **Standardize on Rust**: Migrate all logic to Rust.

### Decision

- Migrate all consumer/worker logic to Rust and remove Go consumer/worker logic.

### Consequences

- Simplified the technology stack and improved maintainability.

### Notes

- This change helps in reducing the complexity of the system.

---

## Decision 17: Removed ksqlDB from Main, TBD When to Bring Back

**Date**: 2024-09-08

### Context

- ksqlDB was removed from the main setup to simplify the architecture, with a decision to bring it back later if needed.

### Options Considered

1. **Keep ksqlDB**: Continue using ksqlDB.
2. **Remove ksqlDB**: Remove ksqlDB and simplify the architecture.

### Decision

- Remove ksqlDB from the main setup.

### Consequences

- Simplified the architecture and reduced complexity.

### Notes

- ksqlDB may be reintroduced in the future based on requirements.

---

## Decision 18: Configured Priority Queues to be Staggered Dynamically Based on Routine Type and Priority

**Date**: 2024-09-09

### Context

- Priority queues were configured to be staggered dynamically based on routine type and priority to improve task handling.

### Options Considered

1. **Static Configuration**: Use a static configuration for priority queues.
2. **Dynamic Configuration**: Configure priority queues dynamically based on routine type and priority.

### Decision

- Configure priority queues to be staggered dynamically based on routine type and priority.

### Consequences

- Improved handling of tasks with different priorities and routine types.

### Notes

- This change helps in better managing and processing tasks dynamically.

## Decision 19: Refactor Python Simulator to Execute POST Requests According to Provided JSON Datasets

**Date**: 2024-09-10

### Context

- The Python simulator was refactored to execute POST requests according to provided JSON datasets, with set times between 00:00 and 06:00 in MM:SS format.

### Options Considered

1. **Keep Existing Simulator**: Continue using the existing simulator.
2. **Refactor Simulator**: Refactor the simulator to execute POST requests according to provided JSON datasets.

### Decision

- Refactor the Python simulator to execute POST requests according to provided JSON datasets.

### Consequences

- Improved ability to simulate real-world scenarios and test the system under different conditions.
- Python simulator consumes more resources with the real-time simulator implementation.

### Notes

- This change helps in better testing and validating the system with realistic data and timing.

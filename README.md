### Wedding Event Simulation Project

This project simulates a wedding event where various teams handle different types of events. The simulation aims to process these events efficiently within specified time frames to maintain guest satisfaction.

#### Tech Stack

- **Python Event Simulator**: Makes POST requests to the Go server.
- **Go Gin Server**: Receives requests and produces Kafka topics.
- **Kafka**: Manages event messaging between the producer (API) and consumers (teams).
- **Rust Consumer/Worker**: Handles topics and inserts results into the Postgres database.
- **Postgres DB**: Stores event statuses.
- **Zookeeper**: Coordinates and manages Kafka brokers.
- **Kafka Manager**: (Optional) Provides a user interface to manage Kafka.

#### Getting Started

Follow these steps to set up and run the project:

1. **Make the script executable:**

   ```bash
   chmod +x scripts/create-topics.sh
   ```

2. **Start Kafka and Zookeeper services:**

   ```bash
   docker compose -f docker-compose-services.yml up -d
   ```

3. **Build and start the application containers:**

   ```bash
   docker compose up --build
   ```

4. **Run python simulator**:

   To run the simulator with a randomized set of events, execute the following command:

   ```bash
   python3 simluator/simulator.py
   ```

   To run the simulator with Gabriel's datasets, execute the following command:

   ```bash
   python3 simluator/simulator.py assets/dataset_1.json # or 2,3,4,5
   ```

#### Project Structure

- **Kafka:** Manages event messaging between the producer (API) and consumers (teams).
- **Zookeeper:** Coordinates and manages Kafka brokers.
- **Kafka Manager:** (Optional) Provides a user interface to manage Kafka.
- **Go Server (API):** Receives POST requests and produces events to Kafka topics.
- **Scalable Rust Consumers:** Multiple consumers, one for each team, consuming specific topics from Kafka.

#### Topics and Teams

The Kafka topics are created based on different teams handling various events:

- **Security:** `brawl`, `not-on-list`, `person-fell`, `injured-kid`
- **Clean-up:** `dirty-table`, `broken-glass`
- **Catering:** `bad-food`, `music-too-loud`, `music-too-low`, `feeling-ill`
- **Officiant:** `missing-rings`, `missing-bride`, `missing-groom`
- **Waiters:** `broken-glass`, `feeling-ill`

#### Simulation

The simulation runs for 6 minutes, during which events are generated and processed. Each event has a type and priority, which determines the processing time frame:

- **High Priority:** 5 seconds
- **Medium Priority:** 10 seconds
- **Low Priority:** 15 seconds

Teams handle events based on their specified routines:

- **Standard:** 20 seconds working, 5 seconds idle
- **Intermittent:** 5 seconds working, 5 seconds idle
- **Concentrated:** 60 seconds working, 60 seconds idle

#### Troubleshooting

- Kafka tends to to boot up successfully on the first go, so you may need to restart the services.
- If the Go server fails to start, check the logs for any errors.
- If the Rust consumers fail to start, ensure the Kafka topics are created and the server is running.

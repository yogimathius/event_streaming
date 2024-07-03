### Wedding Event Simulation Project

This project simulates a wedding event where various teams handle different types of events. The simulation aims to process these events efficiently within specified time frames to maintain guest satisfaction.

#### Getting Started

Follow these steps to set up and run the project:

1. **Make the script executable:**

   ```bash
   chmod +x create-topics.sh
   ```

2. **Start Kafka and Zookeeper services:**

   ```bash
   docker compose -f docker-compose-kafka.yml up -d
   ```

3. **Build and start the application containers:**
   ```bash
   docker compose up --build
   ```

#### Project Structure

- **Kafka:** Manages event messaging between the producer (API) and consumers (teams).
- **Zookeeper:** Coordinates and manages Kafka brokers.
- **Kafka Manager:** (Optional) Provides a user interface to manage Kafka.
- **Go Server (API):** Receives POST requests and produces events to Kafka topics.
- **Scalable Consumers:** Multiple consumers, one for each team, consuming specific topics from Kafka.

#### Topics and Teams

The Kafka topics are created based on different teams handling various events:

- **Security:** `brawl`, `not-on-list`, `person-fell`, `injured-kid`
- **Clean-up:** `dirty-table`, `broken-glass-clean-up`
- **Catering:** `bad-food`, `music-too-loud`, `music-too-low`, `feeling-ill-catering`
- **Officiant:** `missing-rings`, `missing-bride`, `missing-groom`
- **Waiters:** `broken-glass-waiters`, `feeling-ill-waiters`

#### Simulation

The simulation runs for 6 minutes, during which events are generated and processed. Each event has a type and priority, which determines the processing time frame:

- **High Priority:** 5 seconds
- **Medium Priority:** 10 seconds
- **Low Priority:** 15 seconds

Teams handle events based on their specified routines:

- **Standard:** 20 seconds working, 5 seconds idle
- **Intermittent:** 5 seconds working, 5 seconds idle
- **Concentrated:** 60 seconds working, 60 seconds idle

#### Additional Setup

Ensure you have the necessary dependencies installed and your environment is properly configured before running the project. If using a Python script for simulation, generate the `requirements.txt` file with:

```bash
pip freeze > requirements.txt
```

To install dependencies:

```bash
pip install -r requirements.txt
```

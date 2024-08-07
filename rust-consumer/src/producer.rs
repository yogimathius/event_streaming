use crate::worker::Event;
use rdkafka::config::ClientConfig;
use rdkafka::producer::{BaseProducer, BaseRecord};
use std::time::{SystemTime, UNIX_EPOCH};

pub struct KafkaProducer {
    producer: BaseProducer,
}

impl KafkaProducer {
    pub fn new() -> Self {
        let producer: BaseProducer = ClientConfig::new()
            .set("bootstrap.servers", "kafka:29092")
            .create()
            .expect("Failed to create Kafka producer");

        KafkaProducer { producer }
    }

    pub fn send(&mut self, payload: Event) {
        let topic = &payload.event_type;
        let payload = serde_json::to_string(&payload).unwrap();

        // Get the current timestamp
        let current_timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("Time went backwards")
            .as_millis() as i64;

        println!("Sending payload: {:?}", payload);

        let record = BaseRecord::<(), String>::to(topic)
            .payload(&payload)
            .timestamp(current_timestamp);

        if let Err(e) = self.producer.send(record) {
            println!("Failed to send message to Kafka: {:?}", e);
        }

        // Poll to handle delivery reports
        self.producer.poll(std::time::Duration::from_millis(100));
    }
}

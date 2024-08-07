use crate::worker::Event;
use kafka::producer::{Producer, Record, RequiredAcks};
use std::time::Duration;

pub struct KafkaProducer {
    producer: Producer,
}

impl KafkaProducer {
    pub fn new() -> Self {
        let producer = Producer::from_hosts(vec!["kafka:29092".to_owned()])
            .with_ack_timeout(Duration::from_secs(1))
            .with_required_acks(RequiredAcks::One)
            .create()
            .expect("Failed to create Kafka producer");

        KafkaProducer { producer }
    }

    pub fn send(&mut self, payload: Event) {
        let topic = payload.clone().event_type;
        let payload = serde_json::to_string(&payload).unwrap();
        let record = Record::from_value(&topic, payload.as_bytes()).with_partition(1);
        println!("Sending payload: {:?}", payload);

        if let Err(e) = self.producer.send(&record) {
            println!("Failed to send message to Kafka: {:?}", e);
        }
    }
}

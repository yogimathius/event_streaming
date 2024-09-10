use crate::{channel::transmitter::Transmitter, models::event::Event, models::priority::Priority};
use dotenv::dotenv;
use kafka::consumer::{Consumer, FetchOffset, GroupOffsetStorage};
use std::env;

use super::producer::KafkaProducer;

pub struct KafkaConsumer {
    consumer: Consumer,
    producer: KafkaProducer,
}

impl KafkaConsumer {
    pub fn new() -> Self {
        dotenv().ok();
        let kafka_broker = env::var("KAFKA_BROKER").expect("KAFKA_HOSTS must be set");
        let kafka_group = env::var("KAFKA_GROUP").expect("KAFKA_GROUP must be set");
        let kafka_topics: Vec<String> = env::var("KAFKA_TOPICS")
            .expect("KAFKA_TOPICS must be set")
            .split(',')
            .map(|s| s.to_owned())
            .collect();

        let consumer = {
            let mut consumer_builder = Consumer::from_hosts(vec![kafka_broker])
                .with_group(kafka_group)
                .with_fallback_offset(FetchOffset::Earliest)
                .with_offset_storage(Some(GroupOffsetStorage::Kafka));
            for topic in kafka_topics {
                consumer_builder = consumer_builder.with_topic(topic);
            }
            consumer_builder.create().unwrap()
        };

        let producer = KafkaProducer::new();
        KafkaConsumer { consumer, producer }
    }

    pub fn poll(&mut self, high_tx: &Transmitter, medium_tx: &Transmitter, low_tx: &Transmitter) {
        loop {
            for ms in self.consumer.poll().unwrap().iter() {
                for m in ms.messages() {
                    let payload = std::str::from_utf8(m.value).unwrap();
                    match serde_json::from_str::<Event>(payload) {
                        Ok(mut event) => {
                            if event.status == "message produced" {
                                event.status = "message consumed".to_owned();
                                event.event_time = chrono::Utc::now().to_rfc3339();
                                self.producer.send(event.clone());
                                self.process_event(event, high_tx, medium_tx, low_tx);
                            }
                        }
                        Err(e) => {
                            println!("Failed to deserialize message: {:?}", e);
                        }
                    }
                }
                let result = self.consumer.consume_messageset(ms);

                if let Err(e) = result {
                    println!("Error consuming messageset: {:?}", e);
                }

                self.consumer.commit_consumed().unwrap();
            }
            self.consumer.commit_consumed().unwrap();
        }
    }

    fn process_event(
        &mut self,
        event: Event,
        high_tx: &Transmitter,
        medium_tx: &Transmitter,
        low_tx: &Transmitter,
    ) {
        match event.priority {
            Priority::High => self.delegate_to_high_worker(event.clone(), high_tx),
            Priority::Medium => self.delegate_to_med_worker(event.clone(), medium_tx),
            Priority::Low => self.delegate_to_low_worker(event.clone(), low_tx),
        }
    }

    fn delegate_to_high_worker(&mut self, event: Event, tx: &Transmitter) {
        tx.send(event.clone());
        self.update_event_status(event, "delegated high priority");
    }

    fn delegate_to_med_worker(&mut self, event: Event, tx: &Transmitter) {
        tx.send(event.clone());
        self.update_event_status(event, "delegated medium priority");
    }

    fn delegate_to_low_worker(&mut self, event: Event, tx: &Transmitter) {
        tx.send(event.clone());
        self.update_event_status(event, "delegated low priority");
    }

    fn update_event_status(&mut self, mut event: Event, status: &str) {
        event.status = status.to_owned();
        event.event_time = chrono::Utc::now().to_rfc3339();
        self.producer.send(event);
    }
}

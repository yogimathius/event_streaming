use crate::{event::Event, priority::Priority, producer::KafkaProducer, transmitter::Transmitter};
use kafka::consumer::{Consumer, FetchOffset, GroupOffsetStorage};
// use uuid::Uuid;

pub struct KafkaConsumer {
    consumer: Consumer,
    producer: KafkaProducer,
}

impl KafkaConsumer {
    pub fn new() -> Self {
        let consumer = Consumer::from_hosts(vec!["kafka:29092".to_owned()])
            .with_group("security".to_owned())
            .with_topic("brawl".to_owned())
            .with_topic("not_on_list".to_owned())
            .with_topic("accident".to_owned())
            .with_fallback_offset(FetchOffset::Earliest)
            .with_offset_storage(Some(GroupOffsetStorage::Kafka))
            .create()
            .expect("Failed to create Kafka consumer");
        let producer = KafkaProducer::new();
        KafkaConsumer { consumer, producer }
    }

    pub fn poll(&mut self, tx: &Transmitter) {
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
                                self.process_event(event, tx);
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

    fn process_event(&mut self, event: Event, tx: &Transmitter) {
        match event.priority {
            Priority::High => self.delegate_to_high_worker(event.clone(), tx),
            Priority::Medium => self.delegate_to_med_worker(event.clone(), tx),
            Priority::Low => self.delegate_to_low_worker(event.clone(), tx),
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

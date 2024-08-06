use crate::{
    producer::KafkaProducer,
    worker::{Event, Priority, Transmitter},
};
use kafka::consumer::{Consumer, FetchOffset, GroupOffsetStorage};

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
            .unwrap();
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
                                continue;
                            }
                            event.status = "message consumed".to_owned();
                            event.event_time = chrono::Utc::now().to_rfc3339();
                            self.producer.send(event.clone());
                            self.process_event(event, tx);
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

    fn delegate_to_high_worker(&mut self, mut event: Event, tx: &Transmitter) {
        println!("Delegated to high priority worker: {:?}", event);
        tx.send(event.clone());
        event.status = "delegated high priority".to_owned();
        event.event_time = chrono::Utc::now().to_rfc3339();
        self.producer.send(event);
    }

    fn delegate_to_med_worker(&mut self, mut event: Event, tx: &Transmitter) {
        println!("Delegated to medium priority worker: {:?}", event);
        tx.send(event.clone());
        event.status = "delegated medium priority".to_owned();
        event.event_time = chrono::Utc::now().to_rfc3339();
        self.producer.send(event);
    }

    fn delegate_to_low_worker(&mut self, mut event: Event, tx: &Transmitter) {
        println!("Delegated to low priority worker: {:?}", event);
        tx.send(event.clone());
        event.status = "delegated low priority".to_owned();
        event.event_time = chrono::Utc::now().to_rfc3339();
        self.producer.send(event);
    }
}

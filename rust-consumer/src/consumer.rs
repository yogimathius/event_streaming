use crate::worker::{Event, Priority, Transmitter};
use kafka::consumer::{Consumer, FetchOffset, GroupOffsetStorage};

pub struct KafkaConsumer {
    consumer: Consumer,
}

impl KafkaConsumer {
    pub fn new() -> Self {
        let consumer = Consumer::from_hosts(vec!["kafka:9092".to_owned()])
            .with_group("security".to_owned())
            .with_topic("brawl".to_owned())
            .with_topic("not_on_list".to_owned())
            .with_topic("accident".to_owned())
            .with_fallback_offset(FetchOffset::Earliest)
            .with_offset_storage(Some(GroupOffsetStorage::Kafka))
            .create()
            .unwrap();

        KafkaConsumer { consumer }
    }

    pub fn poll(&mut self, tx: &Transmitter) {
        loop {
            for ms in self.consumer.poll().unwrap().iter() {
                for m in ms.messages() {
                    let payload = std::str::from_utf8(m.value).unwrap();
                    match serde_json::from_str::<Event>(payload) {
                        Ok(event) => {
                            tx.send(event);
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

    fn process_event(&self, event: Event) {
        match event.priority {
            Priority::High => self.delegate_to_high_worker(event),
            Priority::Medium => self.delegate_to_med_worker(event),
            Priority::Low => self.delegate_to_low_worker(event),
        }
    }

    fn delegate_to_high_worker(&self, event: Event) {
        println!("Delegated high priority event: {:?}", event);
    }

    fn delegate_to_med_worker(&self, event: Event) {
        println!("Delegated medium priority event: {:?}", event);
    }

    fn delegate_to_low_worker(&self, event: Event) {
        println!("Delegated low priority event: {:?}", event);
    }
}

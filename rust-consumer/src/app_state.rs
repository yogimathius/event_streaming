use crate::{
    consumer::KafkaConsumer,
    worker::{Channel, RoutineType},
};
use std::{sync::Arc, thread};

pub struct AppState {
    channel: Arc<Channel>,
}

impl AppState {
    pub fn new() -> Self {
        let channel = Arc::new(Channel::new());

        AppState { channel }
    }

    pub fn start_workers(&self) -> Vec<thread::JoinHandle<()>> {
        (0..3)
            .map(|worker_id| {
                let channel = Arc::clone(&self.channel);
                thread::spawn(move || {
                    println!("Worker {} started", worker_id);
                    channel.start_worker(worker_id, RoutineType::Standard);
                })
            })
            .collect()
    }

    pub fn start_kafka_consumer(&self) -> thread::JoinHandle<()> {
        let mut kafka_consumer = KafkaConsumer::new();
        let tx_clone = self.channel.tx.clone();
        thread::spawn(move || {
            kafka_consumer.poll(&tx_clone);
        })
    }
}

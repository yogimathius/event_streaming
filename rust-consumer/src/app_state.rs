use crate::{
    channel::channel::Channel, kafka::consumer::KafkaConsumer, models::routine_type::RoutineType,
};
use std::{env, sync::Arc, thread};

pub struct AppState {
    high_priority_channel: Arc<Channel>,
    medium_priority_channel: Arc<Channel>,
    low_priority_channel: Arc<Channel>,
}

impl AppState {
    pub fn new() -> Self {
        let high_priority_channel = Arc::new(Channel::new());
        let medium_priority_channel = Arc::new(Channel::new());
        let low_priority_channel = Arc::new(Channel::new());

        AppState {
            high_priority_channel,
            medium_priority_channel,
            low_priority_channel,
        }
    }

    pub fn start_workers(&self) -> Vec<thread::JoinHandle<()>> {
        let routine_env = env::var("ROUTINE").expect("ROUTINE must be set");
        let routine_type = RoutineType::new(&routine_env);

        let mut handles = Vec::new();

        let high_priority_channel = Arc::clone(&self.high_priority_channel);
        let medium_priority_channel = Arc::clone(&self.medium_priority_channel);
        let low_priority_channel = Arc::clone(&self.low_priority_channel);

        let high_handle = thread::spawn(move || {
            high_priority_channel.start_workers("High", routine_type.clone(), 6)
        });
        handles.push(high_handle);

        let medium_handle = thread::spawn(move || {
            medium_priority_channel.start_workers("Medium", routine_type.clone(), 4)
        });
        handles.push(medium_handle);

        let low_handle = thread::spawn(move || {
            low_priority_channel.start_workers("Low", routine_type.clone(), 4)
        });
        handles.push(low_handle);

        let mut all_handles = Vec::new();
        for handle in handles {
            let worker_handles = handle.join().expect("Failed to join worker thread");
            all_handles.extend(worker_handles);
        }

        all_handles
    }

    pub fn start_kafka_consumer(&self) -> thread::JoinHandle<()> {
        let mut kafka_consumer = KafkaConsumer::new();
        let high_tx = self.high_priority_channel.tx.clone();
        let medium_tx = self.medium_priority_channel.tx.clone();
        let low_tx = self.low_priority_channel.tx.clone();
        thread::spawn(move || {
            kafka_consumer.poll(&high_tx, &medium_tx, &low_tx);
        })
    }
}

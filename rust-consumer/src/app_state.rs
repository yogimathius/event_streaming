use crate::{consumer::KafkaConsumer, models::routine_type::RoutineType, worker::Channel};
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
        let mut handles = Vec::new();
        let routine_env = env::var("ROUTINE").expect("ROUTINE must be set");

        let routine_type = RoutineType::new(&routine_env);

        for worker_id in 0..3 {
            let channel = Arc::clone(&self.high_priority_channel);
            handles.push(thread::spawn(move || {
                println!("High priority worker {} started", worker_id);
                channel.start_worker(worker_id, routine_type);
            }));
        }

        for worker_id in 0..3 {
            let channel = Arc::clone(&self.medium_priority_channel);
            handles.push(thread::spawn(move || {
                println!("Medium priority worker {} started", worker_id);
                channel.start_worker(worker_id, routine_type);
            }));
        }

        for worker_id in 0..3 {
            let channel = Arc::clone(&self.low_priority_channel);
            handles.push(thread::spawn(move || {
                println!("Low priority worker {} started", worker_id);
                channel.start_worker(worker_id, routine_type);
            }));
        }

        handles
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

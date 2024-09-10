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

    fn start_priority_workers(
        &self,
        channel: Arc<Channel>,
        priority: String,
    ) -> Vec<thread::JoinHandle<()>> {
        let routine_env = env::var("ROUTINE").expect("ROUTINE must be set");

        let routine_type = RoutineType::new(&routine_env);
        (0..3)
            .map(|worker_id| {
                let channel = Arc::clone(&channel);
                let priority = priority.clone();
                thread::spawn(move || {
                    println!("{} priority worker {} started", priority, worker_id);
                    channel.start_worker(worker_id, routine_type);
                })
            })
            .collect()
    }

    pub fn start_workers(&self) -> Vec<thread::JoinHandle<()>> {
        let mut handles = Vec::new();

        handles.extend(
            self.start_priority_workers(
                Arc::clone(&self.high_priority_channel),
                "High".to_string(),
            ),
        );
        handles.extend(self.start_priority_workers(
            Arc::clone(&self.medium_priority_channel),
            "Medium".to_string(),
        ));
        handles.extend(
            self.start_priority_workers(Arc::clone(&self.low_priority_channel), "Low".to_string()),
        );

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

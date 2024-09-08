use std::{sync::Arc, thread};

use workers::{
    consumer::KafkaConsumer,
    worker::{Channel, RoutineType},
};

fn main() {
    let channel = Arc::new(Channel::new());
    let tx = channel.tx.clone();

    let worker_handles: Vec<_> = (0..8)
        .map(|worker_id| {
            let channel = Arc::clone(&channel);
            thread::spawn(move || {
                println!("Worker {} started", worker_id);
                channel.start_worker(worker_id, RoutineType::Standard);
            })
        })
        .collect();

    let mut kafka_consumer = KafkaConsumer::new();
    let tx_clone = tx.clone();
    let consumer_handle = thread::spawn(move || {
        kafka_consumer.poll(&tx_clone);
    });
    for handle in worker_handles {
        handle.join().unwrap();
    }
    consumer_handle.join().unwrap();
}

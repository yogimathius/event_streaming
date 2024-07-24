use std::thread;

use workers::{
    consumer::KafkaConsumer,
    worker::{Channel, RoutineType},
};

fn main() {
    let channel = Channel::new();
    let tx = channel.tx.clone();

    let worker = channel.start_worker(RoutineType::Standard);

    let mut kafka_consumer = KafkaConsumer::new();
    let tx_clone = tx.clone();
    let consumer_handle = thread::spawn(move || {
        kafka_consumer.poll(&tx_clone);
    });
    // Ensure the worker thread completes
    worker.join().unwrap();

    // Ensure the consumer thread completes
    consumer_handle.join().unwrap();
}

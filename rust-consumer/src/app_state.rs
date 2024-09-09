use crate::{
    consumer::KafkaConsumer,
    server::run_server,
    worker::{Channel, RoutineType},
};
use rdkafka::producer::FutureProducer;
use rdkafka::ClientConfig;
use std::{sync::Arc, thread};

pub struct AppState {
    channel: Arc<Channel>,
    // producer: Arc<FutureProducer>,
}

impl AppState {
    pub fn new() -> Self {
        let channel = Arc::new(Channel::new());
        // let producer: FutureProducer = ClientConfig::new()
        //     .set("bootstrap.servers", "your_kafka_broker")
        //     .create()
        //     .expect("Producer creation error");

        AppState {
            channel,
            // producer: Arc::new(producer),
        }
    }

    pub fn start_workers(&self) -> Vec<thread::JoinHandle<()>> {
        (0..45)
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

    // pub fn start_http_server(&self) -> thread::JoinHandle<()> {
    //     let producer = Arc::clone(&self.producer);
    //     thread::spawn(move || {
    //         let rt = tokio::runtime::Builder::new_current_thread()
    //             .enable_all()
    //             .build()
    //             .unwrap();
    //         rt.block_on(run_server(producer));
    //     })
    // }
}

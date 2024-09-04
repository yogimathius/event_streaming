use chrono::Utc;
use crossbeam_channel::unbounded;
use std::{
    sync::{Arc, Mutex},
    thread,
};

use crate::{producer::KafkaProducer, receiver::Receiver, transmitter::Transmitter};

#[derive(Debug, Copy, Clone)]
pub enum RoutineType {
    Standard,
    Intermittent,
    Concentrated,
}

pub struct Channel {
    pub tx: Transmitter,
    pub rx: Receiver,
    pub producer: Arc<Mutex<KafkaProducer>>,
}

impl Channel {
    pub fn new() -> Self {
        let (sender, receiver) = unbounded();
        let producer = KafkaProducer::new();

        Channel {
            tx: Transmitter { sender },
            rx: Receiver { receiver },
            producer: Arc::new(Mutex::new(producer)),
        }
    }

    pub fn start_worker(
        &self,
        worker_id: usize,
        _routine_type: RoutineType,
    ) -> thread::JoinHandle<()> {
        let rx = self.rx.clone();
        let builder = thread::Builder::new();
        let producer: Arc<Mutex<KafkaProducer>> = Arc::clone(&self.producer);

        let handle = builder
            .spawn(move || {
                loop {
                    let job = rx.recv(); // we could use try_recv too

                    match job {
                        Some(mut job) => {
                            job.status = "picked up by worker".to_owned();
                            job.event_time = Utc::now().to_rfc3339();
                            let mut producer = producer.lock().expect("Failed to lock producer");

                            producer.send(job.clone());
                            println!("Worker {} picked up job: {:?}", worker_id, job);
                        }
                        None => break,
                    }
                }
            })
            .expect("Failed to spawn thread");

        handle
    }
}

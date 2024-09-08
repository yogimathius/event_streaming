use chrono::{DateTime, Utc};
use crossbeam_channel::unbounded;
use r2d2::{Pool, PooledConnection};
use r2d2_postgres::{postgres::NoTls, PostgresConnectionManager};
use rand::Rng;
use std::{
    sync::{Arc, Mutex},
    thread,
    time::Duration,
};

use crate::{
    event::Event, priority::Priority, producer::KafkaProducer, receiver::Receiver,
    transmitter::Transmitter,
};

#[derive(Debug, Copy, Clone)]
pub enum RoutineType {
    Standard,
    Intermittent,
    Concentrated,
}

impl RoutineType {
    pub fn working_duration(&self) -> Duration {
        match self {
            RoutineType::Standard => Duration::from_secs(10),
            RoutineType::Intermittent => Duration::from_secs(5),
            RoutineType::Concentrated => Duration::from_secs(60),
        }
    }

    pub fn idle_duration(&self) -> Duration {
        match self {
            RoutineType::Standard => Duration::from_secs(5),
            RoutineType::Intermittent => Duration::from_secs(10),
            RoutineType::Concentrated => Duration::from_secs(60),
        }
    }

    pub fn time_to_complete(priority: Priority) -> u64 {
        match priority {
            Priority::High => 5,
            Priority::Medium => 10,
            Priority::Low => 15,
        }
    }
}

pub struct Channel {
    pub tx: Transmitter,
    pub rx: Receiver,
    pub producer: Arc<Mutex<KafkaProducer>>,
    pub pool: Pool<PostgresConnectionManager<NoTls>>,
}

impl Channel {
    pub fn new() -> Self {
        let (sender, receiver) = unbounded();
        let producer = KafkaProducer::new();
        let manager = PostgresConnectionManager::new(
            "postgres://postgres:pass123@postgres:5432/event_streaming?sslmode=disable"
                .parse()
                .unwrap(),
            NoTls,
        );
        let pool = Pool::new(manager).unwrap();

        Channel {
            tx: Transmitter { sender },
            rx: Receiver { receiver },
            producer: Arc::new(Mutex::new(producer)),
            pool,
        }
    }

    pub fn start_worker(
        &self,
        worker_id: usize,
        routine_type: RoutineType,
    ) -> thread::JoinHandle<()> {
        let rx = self.rx.clone();
        let builder = thread::Builder::new();
        let producer: Arc<Mutex<KafkaProducer>> = Arc::clone(&self.producer);
        let pool = self.pool.clone();

        let handle = builder
            .spawn(move || {
                // Introduce a random initial delay to stagger the schedules
                let initial_delay = rand::thread_rng().gen_range(0..10);
                println!(
                    "Worker {} initial delay: {} seconds",
                    worker_id, initial_delay
                );
                thread::sleep(Duration::from_secs(initial_delay));

                loop {
                    println!("Worker {} is working", worker_id);
                    thread::sleep(routine_type.working_duration());

                    println!("Worker {} is idle", worker_id);
                    let idle_duration = routine_type.idle_duration();
                    let start_idle = std::time::Instant::now();

                    while start_idle.elapsed() < idle_duration {
                        if let Some(mut job) = rx.try_recv() {
                            process_job(&pool, worker_id, &mut job, producer.clone());
                        } else {
                            // Sleep for a short duration to avoid busy-waiting
                            thread::sleep(Duration::from_millis(100));
                        }
                    }
                }
            })
            .expect("Failed to spawn thread");

        handle
    }
}

fn process_job(
    pool: &Pool<PostgresConnectionManager<NoTls>>,
    worker_id: usize,
    job: &mut Event,
    producer: Arc<Mutex<KafkaProducer>>,
) {
    let mut client = pool
        .get()
        .expect("Failed to get a connection from the pool");
    let event_time = DateTime::parse_from_rfc3339(&job.event_time)
        .expect("Failed to parse event_time")
        .with_timezone(&Utc);

    let now = Utc::now();
    let time_to_complete = RoutineType::time_to_complete(job.priority.clone());
    if now.signed_duration_since(event_time)
        > chrono::Duration::from_std(Duration::from_secs(time_to_complete)).unwrap()
    {
        println!(
            "❌❌❌❌❌ Worker {} dropped job: {:?} ❌❌❌❌❌",
            worker_id, job
        );
        job.status = format!("message unresolved");
        add_event_message(&mut client, 1, job).expect("Failed to add event message");
        return;
    }

    job.status = format!("picked up by worker {}", worker_id);
    job.event_time = Utc::now().to_rfc3339();
    let mut producer = producer.lock().expect("Failed to lock producer");

    producer.send(job.clone());
    println!(
        "✅✅✅✅✅ Worker {} completing job: {:?} ✅✅✅✅✅",
        worker_id, job
    );

    add_event_message(&mut client, 1, job).expect("Failed to add event message");
    thread::sleep(Duration::from_secs(3));
}

fn add_event_message(
    client: &mut PooledConnection<PostgresConnectionManager<NoTls>>,
    event_id: i32,
    message: &Event,
) -> Result<(), postgres::Error> {
    let query = "INSERT INTO event_messages (event_id, event_type, priority, description, status) VALUES ($1, $2, $3, $4, $5)";
    client.execute(
        query,
        &[
            &event_id,
            &message.event_type,
            &message.priority.as_str(),
            &message.description,
            &message.status,
        ],
    )?;
    Ok(())
}

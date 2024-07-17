use serde::{Deserialize, Serialize};
use std::{
    collections::VecDeque,
    sync::{Arc, Condvar, Mutex},
    thread,
};
#[derive(Debug, Deserialize, Serialize)]
pub struct Event {
    event_type: String,
    priority: String,
    description: String,
}

pub struct Transmitter {
    store: Arc<Mutex<VecDeque<Event>>>,
    emitter: Arc<Condvar>,
}

pub struct Receiver {
    store: Arc<Mutex<VecDeque<Event>>>,
    emitter: Arc<Condvar>,
}

pub struct Channel {
    pub tx: Transmitter,
    pub rx: Receiver,
}

impl Event {
    pub fn new(event_type: &str, priority: &str, description: &str) -> Self {
        Event {
            event_type: event_type.to_string(),
            priority: priority.to_string(),
            description: description.to_string(),
        }
    }
}

impl Transmitter {
    pub fn send(&self, data: Event) {
        self.store.lock().unwrap().push_back(data);
        self.emitter.notify_one();
    }
}

impl Receiver {
    pub fn recv(&self) -> Option<Event> {
        let mut store = self.store.lock().unwrap();

        while store.is_empty() {
            store = self.emitter.wait(store).unwrap();
        }

        store.pop_front()
    }

    // fn try_recv(&self) -> Option<Event> {
    //     self.store.lock().unwrap().pop_front()
    // }
}

impl Channel {
    pub fn new() -> Self {
        let store = Arc::new(Mutex::new(VecDeque::new()));
        let emitter = Arc::new(Condvar::new());

        Channel {
            tx: Transmitter {
                store: Arc::clone(&store),
                emitter: Arc::clone(&emitter),
            },
            rx: Receiver {
                store: Arc::clone(&store),
                emitter: Arc::clone(&emitter),
            },
        }
    }

    pub fn start_worker(&self) -> thread::JoinHandle<()> {
        let rx = self.rx.clone();

        thread::spawn(move || {
            loop {
                let job = rx.recv(); // we could use try_recv too

                match job {
                    Some(job) => println!("Job: {:?}", job),
                    None => break,
                }
            }
        })
    }
}

impl Clone for Receiver {
    fn clone(&self) -> Self {
        Receiver {
            store: Arc::clone(&self.store),
            emitter: Arc::clone(&self.emitter),
        }
    }
}

impl Clone for Transmitter {
    fn clone(&self) -> Self {
        Transmitter {
            store: Arc::clone(&self.store),
            emitter: Arc::clone(&self.emitter),
        }
    }
}

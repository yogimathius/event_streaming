use crossbeam_channel::Sender as CrossbeamSender;

use crate::event::Event;

pub struct Transmitter {
    pub sender: CrossbeamSender<Event>,
}

impl Transmitter {
    pub fn send(&self, data: Event) {
        self.sender.send(data).unwrap();
    }
}

impl Clone for Transmitter {
    fn clone(&self) -> Self {
        Transmitter {
            sender: self.sender.clone(),
        }
    }
}

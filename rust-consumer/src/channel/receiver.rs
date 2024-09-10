use crossbeam_channel::Receiver as CrossbeamReceiver;

use crate::event::Event;

pub struct Receiver {
    pub receiver: CrossbeamReceiver<Event>,
}

impl Receiver {
    pub fn recv(&self) -> Option<Event> {
        self.receiver.recv().ok()
    }

    pub fn try_recv(&self) -> Option<Event> {
        self.receiver.try_recv().ok()
    }
}

impl Clone for Receiver {
    fn clone(&self) -> Self {
        Receiver {
            receiver: self.receiver.clone(),
        }
    }
}

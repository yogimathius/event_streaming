use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::priority::Priority;

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct Event {
    pub event_type: String,
    pub priority: Priority,
    description: String,
    pub status: String,
    pub event_time: String,
    pub id: Uuid,
}

impl Event {
    pub fn new(event_type: &str, priority: &str, description: &str, status: &str) -> Self {
        Event {
            event_type: event_type.to_string(),
            priority: match priority {
                "High" => Priority::High,
                "Medium" => Priority::Medium,
                "Low" => Priority::Low,
                _ => Priority::Low,
            },
            description: description.to_string(),
            status: status.to_string(),
            event_time: chrono::Utc::now().to_rfc3339(),
            id: Uuid::new_v4(),
        }
    }
}

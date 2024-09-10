use std::time::Duration;

use super::priority::Priority;

#[derive(Debug, Copy, Clone)]
pub enum RoutineType {
    Standard,
    Intermittent,
    Concentrated,
}

impl RoutineType {
    pub fn new(routine_type: &str) -> Self {
        match routine_type {
            "Standard" => RoutineType::Standard,
            "Intermittent" => RoutineType::Intermittent,
            "Concentrated" => RoutineType::Concentrated,
            _ => panic!("Invalid routine type"),
        }
    }

    pub fn staggered_duration(&self) -> Duration {
        match self {
            RoutineType::Standard => Duration::from_secs(5),
            RoutineType::Intermittent => Duration::from_secs(1),
            RoutineType::Concentrated => Duration::from_secs(10),
        }
    }

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

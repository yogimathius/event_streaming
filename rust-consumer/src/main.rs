use workers::app_state::AppState;

fn main() {
    let app_state = AppState::new();

    let worker_handles = app_state.start_workers();
    let consumer_handle = app_state.start_kafka_consumer();

    for handle in worker_handles {
        handle.join().unwrap();
    }
    consumer_handle.join().unwrap();
}

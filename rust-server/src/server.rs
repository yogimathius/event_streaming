use crate::event::Event;
use rdkafka::producer::{FutureProducer, FutureRecord};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use warp::Filter;

#[derive(Deserialize, Serialize)]
struct Message {
    id: String,
    event_type: String,
    priority: String,
    description: String,
}

pub async fn run_server(producer: Arc<FutureProducer>) {
    let producer_filter = warp::any().map(move || Arc::clone(&producer));

    let message_route = warp::path("message")
        .and(warp::post())
        .and(warp::body::json())
        .and(producer_filter)
        .and_then(handle_message);

    warp::serve(message_route).run(([127, 0, 0, 1], 3030)).await;
}

async fn handle_message(
    msg: Message,
    producer: Arc<FutureProducer>,
) -> Result<impl warp::Reply, warp::Rejection> {
    let event = Event::new(&msg.event_type, &msg.priority, &msg.description, "Received");

    let payload = serde_json::to_string(&event).expect("Failed to serialize event");

    let record = FutureRecord::to("your_topic")
        .payload(&payload)
        .key("some_key");

    producer
        .send(record, rdkafka::util::Timeout::Never)
        .await
        .expect("Failed to send event");

    Ok(warp::reply::json(&"Message received"))
}

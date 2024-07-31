from confluent_kafka import Consumer, KafkaError

conf = {
    'bootstrap.servers': 'localhost:29092',
    'group.id': 'my_consumer_group',
    'auto.offset.reset': 'earliest'
}

consumer = Consumer(conf)

# Subscribe to a topic
consumer.subscribe(['hello-world-topic-one'])

try:
    while True:
        msg = consumer.poll(timeout=1.0)  # Wait for messages
        if msg is None:
            continue
        if msg.error():
            if msg.error().code() == KafkaError._PARTITION_EOF:
                continue
            else:
                print(msg.error())
                break
        # Process the message
        print(f"Received message: {msg.value().decode('utf-8')}")

except KeyboardInterrupt:
    pass

finally:
    consumer.close()

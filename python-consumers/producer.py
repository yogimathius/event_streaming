import six
import sys
if sys.version_info >= (3, 12, 0):
    sys.modules['kafka.vendor.six.moves'] = six.moves
from kafka import KafkaProducer
import json

producer = KafkaProducer(
    bootstrap_servers='localhost:9092',
    value_serializer=lambda v: json.dumps(v).encode('utf-8')
)

# Simple "Hello World" message for real-time events
real_time_message = {
    'event_type': 'HelloWorld',
    'timestamp': '2024-06-13T12:00:00Z',
    'priority': 'Low',
    'source': 'Guest123',
    'location': 'Main Hall'
}

# Send message to hello-world-topic
producer.send('hello-world-topic', real_time_message)
print("Real-time message sent to hello-world-topic")

producer.flush()

analytics_message = {
    'event_type': 'HelloWorld',
    'timestamp': '2024-06-13T12:00:00Z',
    'source': 'Guest123',
    'location': 'Main Hall',
    'resolved_time': '2024-06-13T12:01:00Z',
    'handler': 'Staff456',
    'resolution': 'Acknowledged',
    'guest_feedback': 'Positive',
    'related_events': ['Event456', 'Event789'],
    'metadata': {
        'device': 'iPhone',
        'network': 'WiFi',
        'app_version': '1.2.3'
    }
}

# Send message to hello-world-analytics
producer.send('hello-world-analytics', analytics_message)
print("Analytics message sent to hello-world-analytics")

producer.flush()

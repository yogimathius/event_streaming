import requests
import random
import json
import time
import datetime
import pytz

# Define the possible events
events = [
    {"event_type": "brawl", "priority": "High", "description": "Fight broke out"},
    {"event_type": "not_on_list", "priority": "Medium", "description": "Person not on list"},
    {"event_type": "accident", "priority": "Low", "description": "Person fell"},
    {"event_type": "dirty_table", "priority": "Low", "description": "Table is dirty"},
    {"event_type": "broken_items", "priority": "Medium", "description": "Glass broken"},
    {"event_type": "bad_food", "priority": "High", "description": "Bad food served"},
    {"event_type": "music", "priority": "Low", "description": "Music is too loud"},
    {"event_type": "music", "priority": "Low", "description": "Music is too low"},
    {"event_type": "feeling_ill", "priority": "Medium", "description": "Guest feeling ill"},
    {"event_type": "bride", "priority": "High", "description": "Bride is missing"},
    {"event_type": "groom", "priority": "High", "description": "Groom is missing"},
    {"event_type": "broken_items", "priority": "Medium", "description": "Glass broken"},
]

# Simulation parameters
simulation_duration = 6 * 60  # 6 minutes in seconds
api_url = "http://localhost:8080/message"  # Update with your actual API endpoint

# Function to generate and send events
def send_event():
    event = random.choice(events)
    current_time_utc = datetime.datetime.now(pytz.utc)

    # Format the time to include the timezone information with a colon
    event["timestamp"] = current_time_utc.strftime("%Y-%m-%dT%H:%M:%S%z")

    # Insert a colon into the timezone part of the timestamp to match the expected format
    event["timestamp"] = event["timestamp"][:-2] + ":" + event["timestamp"][-2:]
    headers = {'Content-Type': 'application/json'}
    
    response = requests.post(api_url, headers=headers, data=json.dumps(event))
    
    if response.status_code == 200:
      print(f"Event sent successfully: {event}")
    else:
      print(f"Failed to send event: {event}")
      print(f"Error: {response.text}")

# Run the simulation
start_time = time.time()

while time.time() - start_time < simulation_duration:
    send_event()
    # Wait for a random duration between events (1 to 10 seconds)
    time.sleep(random.uniform(0.001, 3.0))

print("Simulation completed.")

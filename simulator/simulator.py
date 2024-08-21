import requests
import random
import json
import time
import datetime
import pytz
import sys
import uuid

print ('argument list', sys.argv)

filePath = sys.argv[1] if len(sys.argv) > 1 else ""
f = open(filePath, "r") if filePath else None
 
# Define the possible events
mock_data = [
    {"event_type": "brawl", "priority": "High", "description": "Fight broke out"},
    {"event_type": "brawl", "priority": "Medium", "description": "Minor altercation"},
    {"event_type": "brawl", "priority": "Low", "description": "Verbal argument"},

    {"event_type": "not_on_list", "priority": "High", "description": "VIP not on list"},
    {"event_type": "not_on_list", "priority": "Medium", "description": "Guest not on list"},
    {"event_type": "not_on_list", "priority": "Low", "description": "Staff not on list"},

    {"event_type": "accident", "priority": "High", "description": "Severe injury"},
    {"event_type": "accident", "priority": "Medium", "description": "Moderate injury"},
    {"event_type": "accident", "priority": "Low", "description": "Minor injury"},

    {"event_type": "dirty_table", "priority": "High", "description": "Table extremely dirty"},
    {"event_type": "dirty_table", "priority": "Medium", "description": "Table moderately dirty"},
    {"event_type": "dirty_table", "priority": "Low", "description": "Table slightly dirty"},

    {"event_type": "broken_items", "priority": "High", "description": "Multiple items broken"},
    {"event_type": "broken_items", "priority": "Medium", "description": "Several items broken"},
    {"event_type": "broken_items", "priority": "Low", "description": "Few items broken"},

    {"event_type": "bad_food", "priority": "High", "description": "Food poisoning"},
    {"event_type": "bad_food", "priority": "Medium", "description": "Food undercooked"},
    {"event_type": "bad_food", "priority": "Low", "description": "Food not tasty"},

    {"event_type": "music", "priority": "High", "description": "Music stopped"},
    {"event_type": "music", "priority": "Medium", "description": "Music too loud"},
    {"event_type": "music", "priority": "Low", "description": "Music too low"},

    {"event_type": "feeling_ill", "priority": "High", "description": "Guest fainted"},
    {"event_type": "feeling_ill", "priority": "Medium", "description": "Guest nauseous"},
    {"event_type": "feeling_ill", "priority": "Low", "description": "Guest feeling dizzy"},

    {"event_type": "bride", "priority": "High", "description": "Bride is missing"},
    {"event_type": "bride", "priority": "Medium", "description": "Bride is late"},
    {"event_type": "bride", "priority": "Low", "description": "Bride is upset"},

    {"event_type": "groom", "priority": "High", "description": "Groom is missing"},
    {"event_type": "groom", "priority": "Medium", "description": "Groom is late"},
    {"event_type": "groom", "priority": "Low", "description": "Groom is upset"},
]

events = json.loads(f.read()) if f else mock_data

# Simulation parameters
simulation_duration = 180 # 6 minutes in seconds
api_url = "http://localhost:8080/message"  # Update with your actual API endpoint

# Function to generate and send events
def send_event(event):
    current_time_utc = datetime.datetime.now(pytz.utc)

    # Format the time to include the timezone information with a colon
    event["event_time"] = current_time_utc.strftime("%Y-%m-%dT%H:%M:%S%z")
    id = uuid.uuid1()
    event["id"] = str(id)
    # Insert a colon into the timezone part of the event_time to match the expected format
    event["event_time"] = event["event_time"][:-2] + ":" + event["event_time"][-2:]
    headers = {'Content-Type': 'application/json'}
    
    response = requests.post(api_url, headers=headers, data=json.dumps(event))
    
    if response.status_code == 200:
      print(f"Event sent successfully: {event}")
    else:
      print(f"Failed to send event: {event}")
      print(f"Error: {response.text}")

# Run the simulation
start_time = time.time()

if filePath:
   for event in events:
      print(f"Event loaded: {event}")
      send_event(event)
else:
  while time.time() - start_time < simulation_duration:
      event = random.choice(events)
      send_event(event)
      # Wait for a random duration between events (1 to 10 seconds)
      time.sleep(random.uniform(0.2, 0.5))

print("Simulation completed.")

import requests
import json
import time

url = "http://localhost:5000/ai"
headers = {"Content-Type": "application/json"}
data = {
    "outlet_id": "4",
    "outlet_location": "Nice",
    "central_hub_location": "Paris",
    "date": "2024-01-07T00:00:00Z",
    "event": "No event",
    "event_description": "No event",
    "client_preferences": "Strong demand for Olive Oil and Baguette, moderate interest in Black Tea, minimal preference for Manchego Cheese.",
    "outlet_inventory": {
        "baguette": {
            "current_storage_amount": 300,
            "daily_replenishment_without_envent_from_central_hub": 50,
            "max_warehouse_capacity": 300
        },
        "black_tea": {
            "current_storage_amount": 120,
            "daily_replenishment_without_envent_from_central_hub": 20,
            "max_warehouse_capacity": 250
        },
        "manchego_cheese": {
            "current_storage_amount": 230,
            "daily_replenishment_without_envent_from_central_hub": 40,
            "max_warehouse_capacity": 400
        },
        "olive_oil": {
            "current_storage_amount": 60,
            "daily_replenishment_without_envent_from_central_hub": 30,
            "max_warehouse_capacity": 500
        }
    }
}

print("Sleeping for 2 seconds...")
time.sleep(2)
print(f"Sending request to {url}...")
try:
    response = requests.post(url, headers=headers, json=data)
    print(f"Status Code: {response.status_code}")
    print("Response JSON:")
    print(json.dumps(response.json(), indent=4))
except Exception as e:
    print(f"Error: {e}")

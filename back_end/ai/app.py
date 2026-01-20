from flask import Flask, request, jsonify
import asyncio
import functools
import json
import os
import threading
import time
import websockets

from camel.types import ModelType
from multi_agent_communication_supply_chain import role_playing, messages_queue


global central_hub_json

central_hub_json = {
    "central_hub_inventory": {
        "baguette": {
            "current_storage_amount": 1000,
        },
        "black_tea": {
            "current_storage_amount": 1000,
        },
        "manchego_cheese": {
            "current_storage_amount": 1000,
        },
        "olive_oil": {
            "current_storage_amount": 1000,
        },
    }
}

# Example of request JSON
# response_json = {
#     "outlet_inventory": {
#         "olive_oil": {
#             "changed_replenishment_amount_from_central_hub": 1000,
#         },
#         "baguette": {
#             "changed_replenishment_amount_from_central_hub": 1000,
#         },
#         "manchego_cheese": {
#             "changed_replenishment_amount_from_central_hub": 1000,
#         },
#         "black_tea": {
#             "changed_replenishment_amount_from_central_hub": 1000,
#         }
#     },
#     "central_hub_inventory": {
#         "olive_oil": {
#             "current_storage_amount": 1000,
#         },
#         "baguette": {
#             "current_storage_amount": 1000,
#         },
#         "manchego_cheese": {
#             "current_storage_amount": 1000,
#         },
#         "black_tea": {
#             "current_storage_amount": 1000,
#         },
#     },
#     "transportation_duration": 1
# }


app = Flask(__name__)

# Store websockets with their paths: {websocket: path}
connected_websockets = {}
messages_queue = messages_queue  # From import

async def broadcast_messages():
    print("Broadcast task started, waiting for messages...")
    while True:
        try:
            message = await asyncio.to_thread(messages_queue.get, timeout=1)
        except:
            # Timeout, just continue checking
            await asyncio.sleep(0.1)
            continue
            
        sender_id = message["sender_id"]
        print(f"Broadcasting message for outlet {sender_id} to {len(connected_websockets)} clients")
        if message is None:
            break
            
        user_message = message["user_message"] + "\n\n"
        assistant_message = message["assistant_message"] + "\n\n"

        # NOTE: Current logic in multi_agent_communication_supply_chain.py:
        # user_response (Hub) -> "user_message"
        # assistant_response (Outlet) -> "assistant_message"
        
        # Determine which WebSocket paths should receive this message
        # Messages for outlet N should go to /messageN
        target_path = f"/message{sender_id}"
        
        chunk_delay = 0.02  # Slow down to 20ms per character for readability

        # Broadcast Hub message (Speaker 0) to relevant clients
        for char in user_message:
            msg_to_send = {
                "SpeakerID": "0",  # Hub
                "ReceiverID": sender_id, # Target Outlet
                "text": char,
            }
            json_msg = json.dumps(msg_to_send)
            for ws, ws_path in list(connected_websockets.items()):
                # Only send to clients connected to the target outlet's message path
                if ws_path == target_path:
                    try:
                        await ws.send(json_msg)
                    except Exception as e:
                        print(f"Error sending to websocket {ws_path}: {e}")
            await asyncio.sleep(chunk_delay)

        # Broadcast Outlet message (Speaker sender_id) to relevant clients
        for char in assistant_message:
            msg_to_send = {
                "SpeakerID": sender_id, # Outlet
                "ReceiverID": "0", # Target Hub
                "text": char,
            }
            json_msg = json.dumps(msg_to_send)
            for ws, ws_path in list(connected_websockets.items()):
                # Only send to clients connected to the target outlet's message path
                if ws_path == target_path:
                    try:
                        await ws.send(json_msg)
                    except Exception as e:
                        print(f"Error sending to websocket {ws_path}: {e}")
            await asyncio.sleep(chunk_delay)

        messages_queue.task_done()
        print(f"Done broadcasting message for outlet {sender_id} to path {target_path}")

async def handle_websocket_connection(websocket):
    path = websocket.request.path if hasattr(websocket, 'request') else websocket.path
    print(f"New WebSocket connection from {path}")
    connected_websockets[websocket] = path
    print(f"Total connected clients: {len(connected_websockets)}")
    try:
        await websocket.wait_closed()
    finally:
        if websocket in connected_websockets:
            del connected_websockets[websocket]
        print(f"Connection closed from {path}, remaining: {len(connected_websockets)}")

def run_websocket_server():
    async def main_server():
        print("Starting WebSocket server on port 8000...")
        async with websockets.serve(handle_websocket_connection, 'localhost', 8000):
            print("WebSocket server started!")
            # Start the broadcast task in background
            broadcast_task = asyncio.create_task(broadcast_messages())
            # Keep server running forever
            await asyncio.Future()
    asyncio.run(main_server())

# Clenup the chat record, path 'back_end/ai/chat_record'
def cleanup_chat_record():
    directory_path = os.path.join(os.path.dirname(__file__), "chat_record")
    for file_name in os.listdir(directory_path):
        file_path = os.path.join(directory_path, file_name)
        os.remove(file_path)

websocket_server_thread = None
current_messages_queue = None

# Define a route for the AI request
@app.route('/ai', methods=['POST'])
def handle_ai_request():
    # Get JSON data from the request
    request_data = request.get_json()

    def format_product_names(json_data):
        if "outlet_inventory" in json_data:
            formatted_inventory = {}
            for product_name, details in json_data["outlet_inventory"].items():
                # Convert the product name to lowercase and replace spaces with underscores
                formatted_name = product_name.lower().replace(" ", "_")
                formatted_inventory[formatted_name] = details

            json_data["outlet_inventory"] = formatted_inventory
        return json_data

    request_data = format_product_names(request_data)

    # Perform some AI-related processing with role_playing
    global central_hub_json
    try:
        cleanup_chat_record()  # Cleanup the chat record
        response_json, updated_central_hub_json = role_playing(request_json=request_data, central_hub_json=central_hub_json, model_type=ModelType.DEEPSEEK_R1, messages_queue=messages_queue)
    except Exception as e:
        print(f"Error in role_playing: {e}")
        import traceback
        traceback.print_exc()
        # If the role_playing function fails, return a default response
        response_json = {
            "outlet_inventory": {
                "baguette": {
                    "future_storage_amount": 50,
                    "specific_reason_of_replenishment": "to meet the moderate demand as per the client\"s preferences"
                },
                "black_tea": {
                    "future_storage_amount": 20,
                    "specific_reason_of_replenishment": "to maintain a minimal stock level due to the client\"s minimal interest"
                },
                "manchego_cheese": {
                    "future_storage_amount": 40,
                    "specific_reason_of_replenishment": "to meet the strong demand as per the client\"s preferences"
                },
                "olive_oil": {
                    "future_storage_amount": 30,
                    "specific_reason_of_replenishment": "to meet the strong demand as per the client\"s preferences"
                }
            },
            "central_hub_inventory": {
                "baguette": {
                    "current_storage_amount": 530
                },
                "black_tea": {
                    "current_storage_amount": 364
                },
                "manchego_cheese": {
                    "current_storage_amount": 530
                },
                "olive_oil": {
                    "current_storage_amount": 180
                }
            },
            "transportation_duration": 1
        }
        updated_central_hub_json = {
            "central_hub_inventory": {
                "baguette": {
                    "current_storage_amount": 530
                },
                "black_tea": {
                    "current_storage_amount": 364
                },
                "manchego_cheese": {
                    "current_storage_amount": 530
                },
                "olive_oil": {
                    "current_storage_amount": 180
                }
            }
        }

    for product in updated_central_hub_json["central_hub_inventory"]:
        product_account = int(updated_central_hub_json["central_hub_inventory"][product]["current_storage_amount"])
        if product_account <= 0:
            updated_central_hub_json["central_hub_inventory"][product]["current_storage_amount"] = 1000
            print(f"Product {product} is out of stock, replenished to 1000.")
    central_hub_json = updated_central_hub_json

    # Return the response from role_playing
    print(f"Response JSON:\n{response_json}")
    return jsonify(response_json)

def run_flask_app():
    # Running on http://0.0.0.0:5000/ without threading even in debug mode
    app.run(debug=False, host='0.0.0.0', port=5000)

def main():
    flask_thread = threading.Thread(target=run_flask_app)
    flask_thread.start()

    run_websocket_server()

if __name__ == "__main__":
    main()

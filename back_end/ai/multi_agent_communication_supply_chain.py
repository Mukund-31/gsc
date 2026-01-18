import json
import os
import queue
import re

from colorama import Fore

from camel.configs import ChatGPTConfig, FunctionCallingConfig
from camel.societies import RolePlaying
from camel.functions import MATH_FUNCS
from camel.types import ModelType, TaskType

from format_agent import FormatAgent

messages_queue = queue.Queue()


def role_playing(model_type=ModelType.GPT_3_5_TURBO, chat_turn_limit=5, request_json=None, central_hub_json=None, messages_queue=None) -> None:
    if request_json is None:
        # Default request json
        request_json = {
            "outlet_id": "1",
            "outlet_location": "Delhi",
            "central_hub_location": "Delhi",
            "date": "2023-12-07T15:04:05Z",
            "event": "Deepavali",
            "event_description": "The upcoming festival is expected to significantly increase the demand for festive foods. We anticipate a higher demand for Ghee and Naan. Preparing additional stock of these items is advised.",
            "client_preferences": "Customers in Delhi show a strong preference for traditional products. Ghee and Naan are particularly popular.",
            "weather": "clear",
            "outlet_inventory": {
                "ghee": {
                    "current_storage_amount": 100,
                    "daily_replenishment_without_envent_from_central_hub": 30,
                    "max_warehouse_capacity": 500
                },
                "naan": {
                    "current_storage_amount": 200,
                    "daily_replenishment_without_envent_from_central_hub": 50,
                    "max_warehouse_capacity": 300
                },
                "paneer": {
                    "current_storage_amount": 150,
                    "daily_replenishment_without_envent_from_central_hub": 40,
                    "max_warehouse_capacity": 400
                },
                "masala_chai": {
                    "current_storage_amount": 150,
                    "daily_replenishment_without_envent_from_central_hub": 40,
                    "max_warehouse_capacity": 500
                }
            }
        }
    # Copy the central hub json
    _central_hub_json = central_hub_json.copy()
    user_id = request_json["outlet_id"]
    response_json = {
        "outlet_inventory": {
            "naan": {
                "future_storage_amount": "<NUM>",
                "specific_reason_of_replenishment": "<STRING>",
            },
            "masala_chai": {
                "future_storage_amount": "<NUM>",
                "specific_reason_of_replenishment": "<STRING>",
            },
            "paneer": {
                "future_storage_amount": "<NUM>",
                "specific_reason_of_replenishment": "<STRING>",
            },
            "ghee": {
                "future_storage_amount": "<NUM>",
                "specific_reason_of_replenishment": "<STRING>",
            },
        },
        "transportation_duration": "<NUM> day"
    }
    # Add the central hub inventory to the request
    input_json = request_json | central_hub_json

    context_text = "===== CONTEXT =====\n" + json.dumps(input_json, indent=4) + """
The \"historical_daily_replenishment_amount_from_central_hub\" means the average daily replenishment amount from the central hub to the outlet in the past. So it could be used as a reference for the replenishment amount in the future.
The \"max_warehouse_capacity\" means the maximum capacity of the warehouse of the outlet.
The \"specific_reason_of_replenishment\" means the specific reason of replenishment for the outlet (the decisions made by the central hub) at present.
The current storage amount of the outlet should be less than the maximum capacity of the warehouse of the outlet.
While making decisions, the central hub should first consider the neccessary information in the context, and then predict what is the unknown demand of outlet in the event.
"""
    task_prompt = "Fill the JSON template with ACTUAL NUMBERS. For each product:\n" + \
        "1. If event implies HIGH demand: future_storage_amount = current + (2-3x daily_replenishment)\n" + \
        "2. If event implies NORMAL demand: future_storage_amount = current + daily_replenishment\n" + \
        "3. Never exceed max_warehouse_capacity\n" + \
        "4. Provide a brief reason (1 sentence)\n\n" + \
        "Example: If current=100, daily=30, high demand event → future=160 to 190\n" + \
        "DO NOT explain steps. ONLY output the filled JSON template."

    answer_template = "===== JSON TEMPLATE =====\n"
    answer_template += json.dumps(response_json, indent=4)
    assistant_answer_template = answer_template

    chat_record = context_text

    ai_user_role = "Inventory Management Specialist of Central Hub"
    ai_user_description = "This expert has strong organizational skills, attention to detail, and a deep understanding of supply chain and inventory management systems, who should be proficient in inventory tracking has the ability to analyze stock levels to ensure availability for events. Their duties would include categorizing goods, forecasting demand based on the event description, and configuring the inventory system to reflect accurate information for event-specific requirements."
    ai_assistant_role = "Event Logistics Coordinator of Outlet"
    ai_assistant_description = "This expert has experience in event planning and logistics, with a knack for coordinating with multiple stakeholders, who should have competencies in project management, mathematical calculation (Need to determine the final value of the equations) and problem-solving. Their duties involve understanding the event description to determine the necessary goods, liaising with the Inventory Management Specialist to ensure proper stock levels, and overseeing the setup to meet the event's needs."

    # You can use the following code to play the role-playing game
    function_list = [*MATH_FUNCS]
    
    if model_type.is_open_source:
        from camel.configs import OpenSourceConfig
        if model_type == ModelType.DEEPSEEK_R1:
            model_path = "deepseek-ai/DeepSeek-R1-Distill-Qwen-1.5B"  # For tokenizer
        elif model_type == ModelType.QWEN:
            model_path = "Qwen/Qwen2.5-0.5B"
        else:
            model_path = "vicuna-7b-v1.5"
        server_url = "http://localhost:11434/v1"
        assistant_model_config = OpenSourceConfig(
            model_path=model_path,
            server_url=server_url,
            api_params=ChatGPTConfig(temperature=0.7),
        )
        user_model_config = OpenSourceConfig(
            model_path=model_path,
            server_url=server_url,
            api_params=ChatGPTConfig(temperature=0.7),
        )
    else:
        assistant_model_config = FunctionCallingConfig.from_openai_function_list(
            function_list=function_list,
            kwargs=dict(temperature=0.7),
        )
        assistant_model_config = ChatGPTConfig(temperature=0.7)
        user_model_config = ChatGPTConfig(temperature=0.7)

    sys_msg_meta_dicts = [
        dict(
            assistant_role=ai_assistant_role, user_role=ai_user_role,
            assistant_description=ai_assistant_description + "\n" +
            context_text, user_description=ai_user_description)
        for _ in range(2)
    ]
    role_play_session = RolePlaying(
        assistant_role_name=ai_assistant_role,
        user_role_name=ai_user_role,
        assistant_agent_kwargs=dict(
            model_type=model_type,
            model_config=assistant_model_config,
            function_list=function_list,
        ),
        user_agent_kwargs=dict(
            model_type=model_type,
            model_config=user_model_config,
        ),
        model_type=model_type,
        task_type=TaskType.ROLE_DESCRIPTION,
        task_prompt=task_prompt + "\n" + assistant_answer_template,
        with_task_specify=False,
        extend_sys_msg_meta_dicts=sys_msg_meta_dicts,
    )
    print(Fore.YELLOW + f"Original task prompt:\n{task_prompt}\n")
    print(Fore.CYAN + f"Assistant prompt:\n{role_play_session.assistant_sys_msg.content}\n")
    print(Fore.MAGENTA + f"User prompt:\n{role_play_session.user_sys_msg.content}\n")
 
    n = 0
    input_assistant_msg, _ = role_play_session.init_chat()
    while n < chat_turn_limit:
        n += 1
        assistant_response, user_response = role_play_session.step(
            input_assistant_msg)

        if assistant_response.terminated:
            print(Fore.GREEN +
                  ("AI Assistant terminated. Reason: "
                   f"{assistant_response.info['termination_reasons']}."))
            break
        if user_response.terminated:
            print(Fore.GREEN +
                  ("AI User terminated. "
                   f"Reason: {user_response.info['termination_reasons']}."))
            break

        print(Fore.BLUE + f"{ai_user_role}:\n\n{user_response.msg.content}\n")
        print(Fore.GREEN + f"{ai_assistant_role}:\n\n{assistant_response.msg.content}\n")

        # Output the msg to the markdown file: chat_record.md, and if the file does not exist, create it
        event_name = request_json['event'].replace(" ", "_")  # Replace spaces with underscores for filename
        # 'back_end/ai/chat_record/chat_record_<event_name>.md'
        file_path = os.path.join(os.path.dirname(__file__), "chat_record", f"chat_record_{event_name}.md")
        with open(file_path, "a") as f:
            user_msg_md = user_response.msg.content.replace('\n', '\n\n')
            assistant_msg_md = assistant_response.msg.content.replace('\n', '\n\n')
            f.write(f"[{ai_user_role}]:\n\n{user_msg_md}\n\n\n")
            f.write(f"[{ai_assistant_role}]:\n\n{assistant_msg_md}\n\n\n")

        if messages_queue:
             # Filter out technical content - only send clean conversation
             def clean_message(msg):
                 import re
                 
                 # Remove <think>...</think> blocks but keep what comes after
                 msg = re.sub(r'<think>.*?</think>', '', msg, flags=re.DOTALL)
                 
                 # Remove JSON code blocks
                 msg = re.sub(r'```json.*?```', '', msg, flags=re.DOTALL)
                 
                 # Extract key decision phrases
                 decision_keywords = ['order', 'replenish', 'stock', 'maintain', 'monitor', 'strategy', 'plan']
                 
                 lines = msg.split('\n')
                 decision_lines = []
                 
                 for line in lines:
                     stripped = line.strip()
                     
                     # Skip empty, JSON, or technical lines
                     if not stripped or stripped.startswith(('{', '}', '"', '#', '**', '===', 'Instruction', 'Input', 'Output')):
                         continue
                     
                     # Skip lines that are mostly JSON-like
                     if ':' in stripped and '"' in stripped:
                         continue
                     
                     # Keep lines with decision keywords or substantial content
                     if any(keyword in stripped.lower() for keyword in decision_keywords) or len(stripped) > 30:
                         # Clean up markdown formatting
                         cleaned = re.sub(r'\*\*', '', stripped)
                         decision_lines.append(cleaned)
                 
                 # If we have decision lines, return them
                 if decision_lines:
                     result = '\n'.join(decision_lines[:8])  # Max 8 lines
                     return result[:600]  # Max 600 chars
                 
                 # Fallback: return first substantial paragraph
                 paragraphs = [p.strip() for p in msg.split('\n\n') if len(p.strip()) > 50]
                 if paragraphs:
                     return paragraphs[0][:400]
                 
                 return msg.strip()[:300] if msg.strip() else "[Processing...]"
             
             user_clean = clean_message(user_response.msg.content)
             assistant_clean = clean_message(assistant_response.msg.content)
             
             # Only send if there's actual content after filtering
             if user_clean or assistant_clean:
                 messages_queue.put({
                     "sender_id": user_id, 
                     "user_message": user_clean if user_clean else "[Processing...]", 
                     "assistant_message": assistant_clean if assistant_clean else "[Analyzing...]"
                 })
                 print(Fore.WHITE + f"The length of the messages_queue is {messages_queue.qsize()}\n")

        chat_record += (f"[{ai_user_role}]:{user_response.msg.content}\n\n" + \
            f"[{ai_assistant_role}]:{assistant_response.msg.content}\n\n")

        if "CAMEL_TASK_DONE" in user_response.msg.content or \
            "CAMEL_TASK_DONE" in assistant_response.msg.content:

            if model_type.is_open_source:
                 format_agent = FormatAgent(model_type=model_type, model_config=assistant_model_config)
            else:
                 format_agent = FormatAgent(model_type=ModelType.GPT_4_TURBO, model_config=ChatGPTConfig(temperature=0.0))  # To make the output more readable, we use GPT-4 only
            output_text = format_agent.run(
                user_role_name=ai_user_role,
                assistant_role_name=ai_assistant_role,
                chat_record=chat_record,
                answer_template=response_json,
            ).replace("\'", "\"")
            output_text = re.sub(r'(\w)"(\w)', r'\1\"\2', output_text)
            print(Fore.BLUE + f"output_text:\n{output_text}\n")

            # Extract the json format in the output_text
            # Robust extraction logic to find the JSON block containing "outlet_inventory"
            try:
                # Find the position of "outlet_inventory"
                key_pos = output_text.find('"outlet_inventory"')
                if key_pos == -1:
                    raise ValueError("outlet_inventory not found")
                
                # Find the nearest opening brace before the key
                start_pos = output_text.rfind('{', 0, key_pos)
                if start_pos == -1:
                    raise ValueError("Opening brace not found")
                
                # Try to parse from start_pos to the end, handling potential trailing text
                # We'll try to find the matching closing brace by incrementally shortening from the end
                candidate = output_text[start_pos:]
                found_json = None
                
                # Clean up markdown code blocks if present
                candidate = candidate.strip().replace('```json', '').replace('```', '')
                
                # First try direct parse
                try:
                    found_json = json.loads(candidate)
                except json.JSONDecodeError:
                    # If failed, try to find the last closing brace and slice
                    last_brace = candidate.rfind('}')
                    while last_brace != -1:
                        try:
                            sub_candidate = candidate[:last_brace+1]
                            found_json = json.loads(sub_candidate)
                            break
                        except json.JSONDecodeError:
                            last_brace = candidate.rfind('}', 0, last_brace)
                
                if found_json:
                    role_playing_output_json = found_json
                else:
                    # Fallback to regex
                    role_playing_output_json = json.loads(re.search(r'\{.*"outlet_inventory".*\}', output_text, re.DOTALL).group())

            except Exception as e:
                print(Fore.RED + f"JSON Parsing Error: {e}")
                # Last resort fallback
                role_playing_output_json = json.loads(re.search(r'{.*}', output_text, re.DOTALL).group())

            print(Fore.BLUE + f"role_playing_output_json:\n{json.dumps(role_playing_output_json, indent=4)}\n")
            try:
                # Handle numeric extraction if needed
                if isinstance(role_playing_output_json.get("transportation_duration"), str):
                     nums = [int(s) for s in role_playing_output_json["transportation_duration"].split() if s.isdigit()]
                     role_playing_output_json["transportation_duration"] = nums[0] if nums else 1
                
                if role_playing_output_json.get("transportation_duration", 1) >= 7:
                    role_playing_output_json["transportation_duration"] = 3
            except:
                role_playing_output_json["transportation_duration"] = 1
            # Example of role_playing_output_json
            # {
            #     "outlet_inventory": {
            #         "olive_oil": {
            #             "future_storage_amount": "150",
            #             "specific_reason_of_replenishment": "Expected high demand for Bastille Day event"
            #         },
            #         "baguette": {
            #             "future_storage_amount": "250",
            #             "specific_reason_of_replenishment": "Moderate demand expected for Bastille Day event"
            #         },
            #         "manchego_cheese": {
            #             "future_storage_amount": "400",
            #             "specific_reason_of_replenishment": "Strong preference and high demand expected for Bastille Day event"
            #         },
            #         "black_tea": {
            #             "future_storage_amount": "50",
            #             "specific_reason_of_replenishment": "Minimal interest expected for Bastille Day event"
            #         }
            #     },
            #     "transportation_duration": 2
            # }
            break

        input_assistant_msg = assistant_response.msg

    # Convert string into JSON
    outlet_inventory_json = role_playing_output_json["outlet_inventory"]
    trasportation_duration_json = role_playing_output_json["transportation_duration"]

    # Calculate the changed replenishment amount from central hub
    for product in outlet_inventory_json:
        current_storage_amount = request_json["outlet_inventory"][product]["current_storage_amount"]
        future_storage_amount = outlet_inventory_json[product]["future_storage_amount"]
        changed_replenishment_amount_from_central_hub = int(future_storage_amount) - int(current_storage_amount)
        if changed_replenishment_amount_from_central_hub < 0:
            changed_replenishment_amount_from_central_hub = 0
        central_hub_json["central_hub_inventory"][product]["current_storage_amount"] -= changed_replenishment_amount_from_central_hub

        outlet_inventory_json[product]["future_storage_amount"] = int(future_storage_amount)

    # Format the final answer json
    final_answer_json = {
        "outlet_inventory": outlet_inventory_json,
        "central_hub_inventory": central_hub_json["central_hub_inventory"],
        "transportation_duration": trasportation_duration_json
    }

    print(Fore.RED + f"original_outlet_json:\n{json.dumps(request_json, indent=4)}\n")
    print(Fore.RED + f"original_central_hub_json:\n{json.dumps(_central_hub_json, indent=4)}\n")
    print(Fore.RED + f"final_answer_json:\n{json.dumps(final_answer_json, indent=4)}\n")
    print(Fore.RESET)
    return final_answer_json, central_hub_json

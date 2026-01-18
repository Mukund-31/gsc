# Responsive AI Clusters in Supply Chain (Indian Context)

A sophisticated multi-agent system simulating an autonomous supply chain across major Indian cities. The system uses AI agents powered by **DeepSeek-R1** (via Ollama) to make real-time inventory management decisions based on simulated events.

## 🚀 Project Overview

This project simulates a supply chain network with a Central Hub and multiple Outlets in **Mumbai, Delhi, Bangalore, and Kolkata**. When simulated events occur (e.g., "Diwali" or "New Year's"), the system triggers a multi-agent AI conversation to determine optimal inventory replenishment strategies for items like **Ghee, Naan, Paneer, and Masala Chai**.

### Example Scenario
-   **Event**: "Diwali Festival" in Mumbai.
-   **Impact**: High demand for **Ghee** (for sweets) and **Paneer**.
-   **AI Action**:
    1.  **Hub Manager** (Agent 1) informs the Mumbai outlet about the upcoming festival and suggests stocking up on Ghee.
    2.  **Outlet Coordinator** (Agent 2) checks current stock (e.g., 10 units) and calculates the need (e.g., +100 units).
    3.  **Visual Result**: A delivery truck (animated box) carrying 100 units of Ghee moves from the Hub to Mumbai.
    4.  **Inventory Update**: Upon arrival, the Mumbai inventory table updates exclusively.

### Key Components

1.  **Frontend (Vue.js + D3.js)**: 
    -   Visualizes the map of India (abstracted), inventory tables, and transportation animations.
    -   Displays real-time AI agent chat showing their reasoning process.
    
2.  **Go Backend (Simulation Engine)**:
    -   Manages the physical world state (inventory levels, events, time).
    -   Handles logistics, delivery scheduling (`days_left` countdown), and randomized initial stock.
    -   Communicates with the AI backend.

3.  **Python AI Backend (Intelligence Core)**:
    -   Powered by **DeepSeek-R1:1.5b** (Reasoning Model).
    -   Uses the **CAMEL** framework to orchestrate autonomous dialogue.

## 🔄 System Flowchart

```mermaid
graph TD
    A[Go Simulation Engine] -->|1. Event Detected (e.g., Diwali)| B(Central Hub)
    B -->|2. Send State (Inventory + Event)| C{Python AI Backend}
    C -->|3. Trigger Agents| D[DeepSeek-R1 Model]
    
    subgraph "AI Agent Conversation"
    D --> E[Hub Manager Agent]
    E -->|Suggestion| F[Outlet Coordinator Agent]
    F -->|Reasoning & Decision| E
    end
    
    F -->|4. Final Decision (JSON)| B
    B -->|5. Dispatch Goods| G[Frontend Visualization]
    G -->|6. Show Chat & Moving Boxes| H(User)
    G -->|7. Update Inventory Table| H
```

## 🤖 AI Agents & Communication Protocol

The core intelligence relies on a Role-Playing mechanism where two AI agents collaborate.

### The Protocol (CAMEL Framework)
-   **Initiator**: The Go backend detects an event.
-   **Agent 1 (User)**: `Inventory Management Specialist of Central Hub`
    -   *Role*: Strategic Advisor. "We have plenty of Ghee in the hub, and Diwali is coming. I suggest increasing stock by 20%."
-   **Agent 2 (Assistant)**: `Event Logistics Coordinator of Outlet`
    -   *Role*: Decision Maker. "Agreed. Current stock is 50. I will order 80 more units of Ghee and 40 units of Paneer."
-   **Handshake**: The agents iterate (up to 5 turns) to refine the plan before returning the final JSON.

## 🛠️ Prerequisites

-   **Go** (1.19+)
-   **Python** (3.9+)
-   **Node.js** & **npm**
-   **Ollama** (for running the local LLM)

## 📦 Installation & Setup

### 1. Install & Prepare Ollama
Install Ollama from [ollama.com](https://ollama.com). Then pull the reasoning model:

```bash
ollama pull deepseek-r1:1.5b
```

### 2. Run the System (4 Terminals)

**Terminal 1: AI Model Server**
```bash
ollama serve
```

**Terminal 2: Python AI Backend**
```bash
cd back_end/ai
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

**Terminal 3: Go Simulation Backend**
```bash
cd back_end/go_routine
go run main.go
```

**Terminal 4: Frontend UI**
```bash
cd front_end
npm install
npm run serve
```

## 🖥️ Usage

1.  Open `http://localhost:8080`.
2.  Click **Start**.
3.  Observe the specific needs of **Mumbai, Delhi, Bangalore, and Kolkata** being met by AI-driven decisions.

## 📂 Project Structure

```
├── back_end
│   ├── ai               # Python: DeepSeek-R1 Agents
│   │   └── chat_record  # Logs of agent conversations
│   └── go_routine       # Go: Simulation (Events, Products: Naan, Ghee...)
├── front_end            # Vue.js: Visualization of Indian supply chain
└── setup_env.sh         # Environment setup script
```

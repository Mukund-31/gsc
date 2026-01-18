# Responsive AI Clusters in Supply Chain

A sophisticated multi-agent system simulating an autonomous supply chain. The system uses AI agents powered by **DeepSeek-R1** (via Ollama) to make real-time inventory management decisions based on simulated events.

## 🚀 Project Overview

This project simulates a supply chain network with a Central Hub and multiple Outlets (Paris, Lyon, Marseille, Lille). When simulate events occur (e.g., "New Year's Day"), the system triggers a multi-agent AI conversation to determine optimal inventory replenishment strategies.

### Key Components

1.  **Frontend (Vue.js + D3.js)**: 
    -   Visualizes the supply chain map, inventory levels, and transportation of goods.
    -   Displays real-time AI agent conversations and reasoning.
    
2.  **Go Backend (Simulation Engine)**:
    -   Manages the physical world state (inventory, events, time).
    -   Handles logistics and delivery scheduling.
    -   Communicates with the AI backend to request decisions.

3.  **Python AI Backend (Intelligence Core)**:
    -   Powered by **DeepSeek-R1:1.5b** (Reasoning Model).
    -   Uses the **CAMEL** framework (Communicative Agents for "Mind" Exploration of Large Language Model Society).
    -   Orchestrates autonomous dialogue between two agents:
        1.  **Hub Manager**: Suggests stock adjustments based on global storage.
        2.  **Outlet Coordinator**: Determines specific order quantities based on local events.

## 🤖 AI Agents & Communication Protocol

The core intelligence relies on a Role-Playing mechanism where two AI agents collaborate to solve the inventory problem.

### The Protocol (CAMEL Framework)
The system implements a **Role-Playing** session where agents exchange messages to fulfill a specific task.

-   **Initiator**: The Go backend detects an event and sends a JSON payload to the Python backend.
-   **Agent 1 (User Role)**: `Inventory Management Specialist of Central Hub`
    -   *Responsibility*: Analyzes the event description and Central Hub's available stock. Suggests a replenishment strategy.
-   **Agent 2 (Assistant Role)**: `Event Logistics Coordinator of Outlet`
    -   *Responsibility*: Calculates precise replenishment numbers (e.g., "Order 50 Baguettes"). Fills the structured JSON response.
-   **Model**: `DeepSeek-R1:1.5b` running locally via Ollama. It uses "Chain of Thought" reasoning (visible in logs) to derive decisions.
-   **Handshake**: The agents iterate (up to 5 turns) to refine the plan before returning the final JSON to the Go simulation engine.

### Data Flow
1.  **Event Trigger**: Go Simulation -> Python AI (`POST /ai`)
2.  **Reasoning**: Python AI agents converse (DeepSeek-R1)
3.  **Decision**: Python AI -> Go Simulation (JSON Response)
4.  **Visualization**:
    -   Python AI -> Frontend (WebSocket `:8000`) -> Displays Agent Chat
    -   Go Simulation -> Frontend (WebSocket `:8001`) -> Displays Delivery Animation

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

You need to run 4 separate terminal processes to start the full system.

#### Terminal 1: AI Model Server
Start the Ollama server (keep this running).
```bash
ollama serve
```

#### Terminal 2: Python AI Backend
Handles the intelligent conversation and inventory decisions.
```bash
cd back_end/ai
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```
*Note: Ensure `setup_env.sh` is sourced or paths are correct if you have custom setups.*

#### Terminal 3: Go Simulation Backend
Runs the physical simulation.
```bash
cd back_end/go_routine
go run main.go
```

#### Terminal 4: Frontend UI
Runs the visualization dashboard.
```bash
cd front_end
npm install
npm run serve
```

## 🖥️ Usage

1.  Open your browser to `http://localhost:8080`.
2.  Click the **Start** button in the top right.
3.  Watch as:
    -   Time progresses in the simulation.
    -   Events (like "New Year's Day") trigger agent conversations in the message boxes.
    -   AI agents decide on inventory numbers.
    -   Goods are physically transported (animated boxes) from the Hub to Outlets.
    -   Inventory tables update automatically upon delivery arrival.

## 📂 Project Structure

```
├── back_end
│   ├── ai               # Python: AI Agents, Flask App, WebSockets
│   │   ├── camel        # CAMEL Multi-Agent Framework implementation
│   │   └── app.py       # Main entry point for AI service
│   └── go_routine       # Go: Simulation Engine
│       ├── central_hub  # Hub logic
│       ├── outlet       # Outlet logic
│       └── main.go      # Entry point for Simulation
├── front_end            # Vue.js Application
│   ├── src
│   │   └── App.vue      # Main UI and Visualization logic
│   └── public
└── setup_env.sh         # Environment setup script
```

## 🤝 Contributing
1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

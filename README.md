# Pharma AI Supply Chain (Indian Cities)

A sophisticated multi-agent system simulating an autonomous pharmaceutical supply chain across major Indian cities. The system uses AI agents powered by **DeepSeek-R1** (via Ollama) to make real-time inventory management decisions based on simulated medical events.

## 🚀 Project Overview

This project simulates a pharmaceutical distribution network with a Central Hub and multiple Pharmacies/Hospitals in **Mumbai, Delhi, Bangalore, and Kolkata**. When simulated events occur (e.g., "Flu Season", "Dengue Outbreak"), the system triggers a multi-agent AI conversation to determine optimal inventory replenishment strategies for critical medicines like **Paracetamol, Antibiotics, Insulin, and Vitamins**.

### Example Scenario
-   **Event**: "Flu Season Peak" in Mumbai (Jan 01).
-   **Impact**: Surge in demand for **Paracetamol** (for fever) and **Antibiotics** (for infections).
-   **AI Action**:
    1.  **Hub Manager** (Agent 1) analyzes the event and central stock. "Mumbai is facing a flu surge. We must prioritize Paracetamol."
    2.  **Outlet Coordinator** (Agent 2) calculates precise needs. "Agreed. Current stock is low. Requesting 300 units of Paracetamol."
    3.  **Visual Result**: A delivery truck (animated box) carrying meds moves from the Central Hub to the Mumbai pharmacy.
    4.  **Inventory Update**: The Mumbai outlet's inventory table updates immediately upon delivery.

### Key Components

1.  **Frontend (Vue.js + D3.js)**: 
    -   Visualizes the supply chain network.
    -   Displays inventory tables for medicines.
    -   Shows real-time AI agent chat and reasoning.
    
2.  **Go Backend (Simulation Engine)**:
    -   Manages world state (Inventory: Insulin, Vitamins...).
    -   Triggers medical events (Dengue, Flu, Diabetes Camps).
    -   Handles logistics and delivery scheduling.

3.  **Python AI Backend (Intelligence Core)**:
    -   Powered by **DeepSeek-R1:1.5b**.
    -   Uses **CAMEL** framework for agent role-playing.
    -   **Agents**:
        -   `Role: Central Hub Inventory Specialist`
        -   `Role: Pharmacy Outlet Coordinator`

## 🔄 System Flowchart

```mermaid
graph TD
    A[Go Simulation Engine] -->|1. Medical Event (e.g., Dengue Alert)| B(Central Hub)
    B -->|2. Send State (Medicine Stock)| C{Python AI Backend}
    C -->|3. Trigger Agents| D[DeepSeek-R1 Model]
    
    subgraph "AI Agent Conversation"
    D --> E[Hub Manager Agent]
    E -->|Suggestion| F[Pharmacy Coordinator Agent]
    F -->|Reasoning & Decision| E
    end
    
    F -->|4. Final Order (JSON)| B
    B -->|5. Dispatch Medicines| G[Frontend Visualization]
    G -->|6. Show Chat & Delivery Animation| H(User)
    G -->|7. Update Stock Table| H
```

## 🛠️ Prerequisites

-   **Go** (1.19+)
-   **Python** (3.9+)
-   **Node.js** & **npm**
-   **Ollama** (pull `deepseek-r1:1.5b`)

## 📦 Installation & Setup

### 1. Install & Prepare Ollama
```bash
ollama serve
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
3.  Watch as AI agents respond to **Flu Seasons, Dengue Outbreaks, and Medical Camps** by dispatching items like **Insulin** and **Paracetamol** to **Mumbai, Delhi, Bangalore, and Kolkata**.

## 📂 Project Structure

```
├── back_end
│   ├── ai               # Python: Pharma AI Agents
│   │   └── chat_record  # Logs of logic
│   └── go_routine       # Go: Simulation (Events: Flu, Dengue...)
├── front_end            # Vue.js: Pharma Supply Chain Visualization
└── setup_env.sh         # Environment setup script
```

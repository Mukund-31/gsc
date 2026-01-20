# Responsive AI Clusters in Supply Chain

A real-time supply chain simulation featuring **Google's Agent-to-Agent (A2A) Protocol** with multi-turn AI conversations between warehouse and outlet agents.

## 🌟 Features

- **Google A2A Protocol Implementation** - Full multi-turn agent conversations
- **Bilateral AI Reasoning** - Both Warehouse and Outlets use Ollama AI (deepseek-r1:1.5b)
- **Real-time Negotiation** - Agents can accept, reject, or counter-propose
- **Agent Cards** - JSON capability discovery for all agents
- **Separate Chat Boxes** - Individual conversation panels for each outlet
- **Animated Transfers** - Visual package movement from warehouse to outlets
- **Year-round Events** - 27 events throughout 2024 (Jan-Dec)

## 🏗️ Architecture

```
Backend (Go)          Frontend (React + Vite)
├── Agent Cards       ├── Dashboard UI
├── Multi-turn A2A    ├── WebSocket Client
├── Ollama Client     ├── Framer Motion Animations
└── Event System      └── Tailwind CSS
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.25.6+**
- **Node.js 18+**
- **Ollama** with `deepseek-r1:1.5b` model

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Mukund-31/gsc.git
   cd gsc
   ```

2. **Install Ollama and pull the model:**
   ```bash
   # Install Ollama from https://ollama.ai
   ollama pull deepseek-r1:1.5b
   ```

3. **Build the frontend:**
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

4. **Run the backend:**
   ```bash
   cd backend
   export PATH=$PWD/../tools/go/bin:$PATH
   go run main.go
   ```

5. **Open your browser:**
   ```
   http://localhost:8080
   ```

## 📋 How It Works

### Multi-Turn A2A Conversations

Each event triggers a 3-turn negotiation:

1. **🏪 Turn 1 - Outlet Analysis**
   - Outlet agent analyzes the event
   - Formulates inventory request
   - Sends to Warehouse

2. **🏭 Turn 2 - Warehouse Proposal**
   - Warehouse evaluates request
   - Considers available stock
   - Proposes solution

3. **🏪 Turn 3 - Negotiation**
   - Outlet accepts or negotiates
   - Can counter-propose
   - Final agreement reached

4. **✅ Transfer Complete**
   - Inventory updated
   - Animated package delivery

### Example Conversation

```
🏪 Outlet North
  Event: New Year's Day - High demand for party supplies
  Analysis: Expected 50% increase in demand
  Request: Electronics: 10, Groceries: 20

🏭 Warehouse → Outlet North
  Reasoning: Sufficient stock available, fair distribution
  Proposal: Electronics: 10, Groceries: 20

🏪 Outlet North: accept - Proposal accepted

✅ Outlet North: Agreement reached - Transfer complete
```

## 🎯 Events Calendar

- **January**: New Year's Day, First Weekend, Monday Rush, Mid-month Sale
- **February**: Month Start, Valentine's Day
- **March**: Spring Season, Mid-month Promotion
- **April-December**: Seasonal events including Independence Day, Back to School, Halloween, Black Friday, Christmas

## 🛠️ Technology Stack

### Backend
- **Go 1.25.6** - High-performance backend
- **Ollama** - Local AI model inference
- **WebSockets** - Real-time communication
- **deepseek-r1:1.5b** - Lightweight reasoning model

### Frontend
- **React 18** - UI framework
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Framer Motion** - Animations
- **WebSocket API** - Real-time updates

## 📁 Project Structure

```
gsc/
├── backend/
│   ├── main.go              # Main server & A2A logic
│   ├── agent/
│   │   ├── agent.go         # Agent structures
│   │   ├── agent_cards.go   # A2A Agent Cards
│   │   └── event.go         # Event definitions
│   └── server/
│       └── server.go        # WebSocket server
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   └── Dashboard.jsx # Main UI
│   │   ├── App.jsx
│   │   └── index.css
│   └── dist/                # Built files
└── tools/
    ├── go/                  # Go 1.25.6
    └── node/                # Node.js
```

## 🎨 UI Features

- **Color-coded Messages**
  - 🏪 Outlet messages (Blue)
  - 🏭 Warehouse messages (Orange)
  - ✅ Success messages (Green)

- **Separate Chat Boxes** - Each outlet has its own conversation panel
- **Animated Packages** - 5-second smooth animations showing inventory transfers
- **Real-time Updates** - Live inventory levels and conversation history

## 🔧 Configuration

### Simulation Speed
- **1 Day = 10 seconds** (configurable in `main.go`)

### Animation Duration
- **Package animation: 5 seconds** (configurable in `Dashboard.jsx`)

### AI Model
- **Model**: deepseek-r1:1.5b
- **Endpoint**: http://localhost:11434
- **Format**: JSON responses

## 📊 Performance

- **Concurrent AI calls**: Up to 6 per event (3 turns × 2 agents)
- **WebSocket latency**: <50ms
- **Animation FPS**: 60fps (Framer Motion)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License

## 🙏 Acknowledgments

- **Google A2A Protocol** - Agent-to-Agent communication standard
- **Ollama** - Local AI model inference
- **deepseek-r1** - Efficient reasoning model

## 📞 Contact

For questions or support, please open an issue on GitHub.

---

**Built with ❤️ using Google's A2A Protocol**

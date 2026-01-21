# A2A Supply Chain Demo Script

Use this script to demonstrate the capabilities of the Autonomous Agent-to-Agent (A2A) protocol in the Pharmaceutical Supply Chain simulation.

## Key Concepts to Highlight
1.  **Decentralized Intelligence**: Unlike traditional systems where a central server just pushes updates, here **Every Outlet is an AI** that analyzes its own local situation.
2.  **Autonomous Negotiation**: The Warehouse and Outlets **negotiate** using JSON-RPC. They don't just update a database; they exchange reasoned proposals.
3.  **Real-time Adaptation**: The system handles random events (Flu Season, Recalls) dynamically.

## Demo Walkthrough Steps

### 1. The "Pulse" of the System (Autonomous Events)
*   **Show**: The Dashboard. Point out the **Header** badge: "PROTOCOL: A2A (JSON-RPC) | AUTONOMOUS".
*   **Explain**: "The system runs on a 45-second cycle per day. Every facility is constantly monitoring for events."
*   **Wait**: Watch for a **Purple** `Event:` log to appear in one of the chat boxes.
*   **Point Out**: "See! Outlet-X just detected 'Flu Season'. This wasn't hardcoded; the AI Agent read the event and triggered a response."

### 2. The Analysis Phase (Cyan Logs)
*   **Show**: The **Cyan** text in the chat box starting with "Analysis:".
*   **Read**: "Look at the reasoning. The agent says, 'Current stock is low... expecting spike'. This is **Local Intelligence**. The central warehouse doesn't know this context; the edge agent does."

### 3. The Negotiation Handshake (Orange -> Yellow)
*   **Show**: The **Orange** "Request:" followed by the Warehouse's **Yellow** "Proposal:".
*   **Explain**: "This is the A2A Handshake. The Outlet sends a specialized JSON request. The Warehouse validates it against its own stock (checking for shortages) and returns a specific Proposal. It refused to send 'Antisotes' because they don't exist, and it capped the quantity because it's managing global stock."

### 4. The Agreement & Action (Green + Animation)
*   **Show**: The **Green** "Agreement reached" and the **Package Animation** moving across the screen.
*   **Explain**: "Once both agents agree on the JSON contract, the transaction executes physically (animation) and the ledger updates instantly."

### 5. Advanced Feature: Reverse Logistics (Red Logs)
*   **Scenario**: Wait for a "Recall" event or explain the capability.
*   **Explain**: "If a 'Recall' event occurs, the Outlet Agent creates a **Negative Request**. The Warehouse Agent understands this as a return flow and accepts the stock back. This bidirectional flexibility is unique to agent-based semantic protocols."

## Why is this better?
*   **Resilience**: If the central server goes down, edge agents can still analyze local needs (even if they can't ship immediately).
*   **Precision**: Allocations are based on *reasoning* (e.g., "Urgent care priority"), not just "First Come First Served".
*   **Scalability**: New agents (e.g., a "Vaccine Center") can be added without rewriting the central code—they just need to speak the A2A protocol.

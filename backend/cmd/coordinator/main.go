package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Appointat/Responsive-AI-Clusters-in-Supply-Chain/agent"
	"github.com/a2aproject/a2a-go/a2a"
	"github.com/a2aproject/a2a-go/a2aclient"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Coordinator struct {
	clients         map[string]*a2aclient.Client
	warehouseClient *a2aclient.Client
	wsClients       map[*websocket.Conn]bool
	broadcast       chan []byte
	mutex           sync.Mutex
	
	// Track inventory state for frontend
	warehouseInv map[string]int
	outletInvs   map[string]map[string]int
	
	// Track active conversations to prevent overlap
	activeOutlets map[string]bool
	
	CurrentDate string
}

var randomScenarios = []string{
	"Routine Daily Restock - Maintain safety stock levels",
	"Emergency Walk-ins - Higher than expected patient intake",
	"Antibiotics Usage Spike - Urgent care requirements",
	"General Ward Refill - Standard consumption",
	"Surgical Unit Request - Pre-op preparations",
	"Vaccination Clinic - Scheduled community batch",
	"Depot Shortage - Urgent Recall of Antibiotics",
}

var validItems = map[string]bool{
	"Antibiotics": true,
	"Painkillers": true,
	"Vaccines":    true,
	"Bandages":    true,
}

func (c *Coordinator) generateRandomEvent(date time.Time) agent.Event {
	outlets := []string{"Outlet-1", "Outlet-2", "Outlet-3", "Outlet-4"}
	outletID := outlets[rand.Intn(len(outlets))] // Pick random outlet
	desc := randomScenarios[rand.Intn(len(randomScenarios))] // Pick random scenario
	
	return agent.Event{
		Date:        date,
		Description: desc,
		OutletID:    outletID,
	}
}

func main() {
	coord := &Coordinator{
		clients:       make(map[string]*a2aclient.Client),
		wsClients:     make(map[*websocket.Conn]bool),
		broadcast:     make(chan []byte),
		activeOutlets: make(map[string]bool),
		
		warehouseInv: map[string]int{
			"Antibiotics": 5000,
			"Painkillers": 10000,
			"Vaccines":    2000,
			"Bandages":    5000,
		},
		outletInvs: map[string]map[string]int{
			"Outlet-1": {"Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200},
			"Outlet-2": {"Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200},
			"Outlet-3": {"Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200},
			"Outlet-4": {"Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200},
		},
	}
	
	// Create custom HTTP client with no timeout for A2A calls
	noTimeoutClient := &http.Client{
		Timeout: 0, // No timeout - let context handle it
	}

	// Start WebSocket server
	http.HandleFunc("/ws", coord.handleWS)
	go coord.handleMessages()

	go func() {
		log.Println("Coordinator UI Server starting on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("Server error:", err)
		}
	}()

	// Wait for agent servers to start
	time.Sleep(5 * time.Second)

	// Connect to warehouse
	ctx := context.Background()
	whCard, err := fetchAgentCard(ctx, "http://localhost:8081/.well-known/agent-card.json")
	if err != nil {
		log.Printf("Failed to fetch warehouse card: %v", err)
	} else {
		coord.warehouseClient, err = a2aclient.NewFromCard(ctx, whCard, a2aclient.WithJSONRPCTransport(noTimeoutClient))
		if err != nil {
			log.Printf("Failed to create warehouse client: %v", err)
		} else {
			log.Println("Connected to Central Medical Depot")
		}
	}

	// Connect to Outlets
	outlets := []struct{ ID, Port string }{
		{"Outlet-1", "8082"},
		{"Outlet-2", "8083"},
		{"Outlet-3", "8084"},
		{"Outlet-4", "8085"},
	}

	for _, o := range outlets {
		cardURL := fmt.Sprintf("http://localhost:%s/.well-known/agent-card.json", o.Port)
		card, err := fetchAgentCard(ctx, cardURL)
		if err == nil {
			client, err := a2aclient.NewFromCard(ctx, card, a2aclient.WithJSONRPCTransport(noTimeoutClient))
			if err == nil {
				coord.clients[o.ID] = client
				log.Printf("Connected to %s", o.ID)
			} else {
				log.Printf("Failed to create client for %s: %v", o.ID, err)
			}
		} else {
			log.Printf("Failed to fetch card for %s: %v", o.ID, err)
		}
	}

	// Start simulation loop
	log.Println("Starting A2A Simulation Loop...")
	// Start simulation timer
	ticker := time.NewTicker(45 * time.Second) // 1 day every 45 seconds
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for range ticker.C {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format("2006-01-02")

		coord.mutex.Lock()
		coord.CurrentDate = dateStr
		coord.mutex.Unlock()

		coord.broadcastState()


		// Check for events
		events := getEventsForDate(date)
		
		// If no events scheduled, generate a random one to keep things lively
		if len(events) == 0 {
			evt := coord.generateRandomEvent(date)
			events = append(events, evt)
		}

		for _, evt := range events {
			go coord.runConversation(evt)
		}


	}
}

func (c *Coordinator) handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer ws.Close()

	c.mutex.Lock()
	c.wsClients[ws] = true
	c.mutex.Unlock()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			c.mutex.Lock()
			delete(c.wsClients, ws)
			c.mutex.Unlock()
			break
		}
	}
}

func (c *Coordinator) handleMessages() {
	for {
		msg := <-c.broadcast
		c.mutex.Lock()
		for client := range c.wsClients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Write error: %v", err)
				client.Close()
				delete(c.wsClients, client)
			}
		}
		c.mutex.Unlock()
	}
}

func (c *Coordinator) sendUpdate(data interface{}) {
	jsonMsg, _ := json.Marshal(data)
	c.broadcast <- jsonMsg
}

func (c *Coordinator) broadcastState() {
	c.mutex.Lock()
	dateStr := c.CurrentDate
	
	whInvCopy := make(map[string]int)
	for k, v := range c.warehouseInv {
		whInvCopy[k] = v
	}

	var outletList []map[string]interface{}
	for id, inv := range c.outletInvs {
		name := "Unknown Facility"
		if id == "Outlet-1" { name = "City General Hospital" }
		if id == "Outlet-2" { name = "Community Clinic South" }
		if id == "Outlet-3" { name = "University Medical Center" }
		if id == "Outlet-4" { name = "Metro Pharmacy" }
		
		invCopy := make(map[string]int)
		for k, v := range inv {
			invCopy[k] = v
		}

		outletList = append(outletList, map[string]interface{}{
			"id":        id,
			"Name":      name,
			"Inventory": invCopy,
		})
	}
	
	sort.Slice(outletList, func(i, j int) bool {
		return outletList[i]["id"].(string) < outletList[j]["id"].(string)
	})
	
	c.mutex.Unlock()

	c.sendUpdate(map[string]interface{}{
		"type": "state",
		"date": dateStr,
		"warehouse": map[string]interface{}{
			"Name":      "Central Medical Depot",
			"Inventory": whInvCopy,
		},
		"outlets": outletList,
	})
}

func (c *Coordinator) runConversation(evt agent.Event) {
	outletID := evt.OutletID

	// Check if outlet is already busy
	c.mutex.Lock()
	if c.activeOutlets[outletID] {
		c.mutex.Unlock()
		log.Printf("Skipping event for %s - busy with previous conversation", outletID)
		return
	}
	c.activeOutlets[outletID] = true
	c.mutex.Unlock()

	// Ensure we clear the busy flag when done
	defer func() {
		c.mutex.Lock()
		delete(c.activeOutlets, outletID)
		c.mutex.Unlock()
	}()

	// Use explicit 10 minute timeout for AI operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	c.sendUpdate(map[string]interface{}{
		"type":    "log",
		"message": fmt.Sprintf("🏪 [%s] %s\n  Event: %s", outletID, evt.Date.Format("2006-01-02"), evt.Description),
	})

	client, ok := c.clients[outletID]
	if !ok {
		log.Printf("Client not found for %s", outletID)
		return
	}

	// TURN 1: Analyze Event
	input1 := map[string]interface{}{
		"type":              "analyze_event",
		"event_description": evt.Description,
	}
	inputBytes1, _ := json.Marshal(input1)

	msg1 := a2a.NewMessage(a2a.MessageRoleUser, &a2a.TextPart{Text: string(inputBytes1)})
	result1, err := client.SendMessage(ctx, &a2a.MessageSendParams{Message: msg1})
	if err != nil {
		log.Printf("Turn 1 failed: %v", err)
		return
	}

	task1, _ := result1.(*a2a.Task)
	var output1 map[string]interface{}
	if task1 != nil && len(task1.Artifacts) > 0 {
		log.Printf("[Coordinator] Turn 1: Got %d artifacts", len(task1.Artifacts))
		for _, p := range task1.Artifacts[0].Parts {
			// Try pointer type
			if txt, ok := p.(*a2a.TextPart); ok {
				json.Unmarshal([]byte(txt.Text), &output1)
				break
			}
			// Try value type
			if txt, ok := p.(a2a.TextPart); ok {
				json.Unmarshal([]byte(txt.Text), &output1)
				break
			}
		}
	} else {
		log.Printf("[Coordinator] Turn 1: No artifacts in result")
	}

	analysis, _ := output1["analysis"].(string)
	request, _ := output1["request"].(map[string]interface{})

	log.Printf("[Coordinator] Turn 1 complete: analysis=%s, request=%v", analysis, request)

	c.sendUpdate(map[string]interface{}{
		"type":    "log",
		"message": fmt.Sprintf("🏪 %s Analysis: %s\nRequest: %v", outletID, analysis, request),
	})

	// TURN 2: Warehouse Proposal
	if c.warehouseClient == nil {
		log.Println("Warehouse not connected")
		return
	}

	input2 := map[string]interface{}{
		"type":        "proposal_request",
		"outlet_name": outletID,
		"request":     request,
		"analysis":    analysis,
	}
	inputBytes2, _ := json.Marshal(input2)

	msg2 := a2a.NewMessage(a2a.MessageRoleUser, &a2a.TextPart{Text: string(inputBytes2)})
	result2, err := c.warehouseClient.SendMessage(ctx, &a2a.MessageSendParams{Message: msg2})
	if err != nil {
		log.Printf("Turn 2 failed: %v", err)
		return
	}

	task2, _ := result2.(*a2a.Task)
	var output2 map[string]interface{}
	if task2 != nil && len(task2.Artifacts) > 0 {
		for _, p := range task2.Artifacts[0].Parts {
			if txt, ok := p.(*a2a.TextPart); ok {
				json.Unmarshal([]byte(txt.Text), &output2)
				break
			}
			if txt, ok := p.(a2a.TextPart); ok {
				json.Unmarshal([]byte(txt.Text), &output2)
				break
			}
		}
	}

	proposal, _ := output2["proposal"].(map[string]interface{})
	reasoning, _ := output2["reasoning"].(string)

	c.sendUpdate(map[string]interface{}{
		"type":    "log",
		"message": fmt.Sprintf("🏭 Warehouse → %s\n  Reasoning: %s\n  Proposal: %v", outletID, reasoning, proposal),
	})

	// TURN 3: Negotiation
	input3 := map[string]interface{}{
		"type":             "negotiate_proposal",
		"proposal":         proposal,
		"reasoning":        reasoning,
		"original_request": request,
	}
	inputBytes3, _ := json.Marshal(input3)

	msg3 := a2a.NewMessage(a2a.MessageRoleUser, &a2a.TextPart{Text: string(inputBytes3)})
	result3, err := client.SendMessage(ctx, &a2a.MessageSendParams{Message: msg3})
	if err != nil {
		log.Printf("Turn 3 failed: %v", err)
		return
	}

	task3, _ := result3.(*a2a.Task)
	var output3 map[string]interface{}
	if task3 != nil && len(task3.Artifacts) > 0 {
		for _, p := range task3.Artifacts[0].Parts {
			if txt, ok := p.(*a2a.TextPart); ok {
				json.Unmarshal([]byte(txt.Text), &output3)
				break
			}
			if txt, ok := p.(a2a.TextPart); ok {
				json.Unmarshal([]byte(txt.Text), &output3)
				break
			}
		}
	}

	decision, _ := output3["decision"].(string)
	
	log.Printf("[Coordinator] Turn 3: decision=%s, full output=%v", decision, output3)

	c.sendUpdate(map[string]interface{}{
		"type":    "log",
		"message": fmt.Sprintf("🏪 %s: %s", outletID, decision),
	})

	// Check for accept (case insensitive) or assume accept if proposal exists
	isAccepted := strings.EqualFold(decision, "accept") || (decision == "" && proposal != nil)
	
	if isAccepted {
		log.Printf("[Coordinator] Triggering shipment for %s with proposal: %v", outletID, proposal)
		
		// Send shipments directly from proposal
		for item, val := range proposal {
			// Skip invalid/hallucinated items
			if !validItems[item] {
				continue
			}

			qtyFloat, ok := val.(float64)
			if !ok {
				continue
			}
			qty := int(qtyFloat)
			
			// Update Coordinator's local inventory state
			c.mutex.Lock()
			
			// Validate quantities to prevent negative stock (handling both inflow and outflow)
			whVal := c.warehouseInv[item]
			var outVal int
			if c.outletInvs[outletID] != nil {
				outVal = c.outletInvs[outletID][item]
			}
			
			if qty > 0 {
				// Outflow: Warehouse -> Outlet
				if whVal < qty {
					qty = whVal
				}
			} else {
				// Inflow: Outlet -> Warehouse (Return)
				// qty is negative (e.g. -50). outVal must be >= 50
				if outVal < -qty {
					qty = -outVal
				}
			}

			// Apply updates
			c.warehouseInv[item] -= qty
			if c.outletInvs[outletID] == nil {
				c.outletInvs[outletID] = make(map[string]int)
			}
			c.outletInvs[outletID][item] += qty
			
			c.mutex.Unlock()

			c.sendUpdate(map[string]interface{}{
				"type": "shipment",
				"from": "Warehouse",
				"to":   outletID,
				"item": item,
				"qty":  qty,
			})
		}

		c.sendUpdate(map[string]interface{}{
			"type":    "log",
			"message": fmt.Sprintf("✅ %s: Agreement reached - Transfer complete\n---", outletID),
		})
		
		// Immediately broadcast new inventory state
		c.broadcastState()
	}
}

func fetchAgentCard(ctx context.Context, url string) (*a2a.AgentCard, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var card a2a.AgentCard
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, err
	}
	return &card, nil
}

func getEventsForDate(date time.Time) []agent.Event {
	// Use existing event system
	return agent.GetEventsForDate(date)
}

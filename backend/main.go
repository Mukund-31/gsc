package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Appointat/Responsive-AI-Clusters-in-Supply-Chain/agent"
	"github.com/Appointat/Responsive-AI-Clusters-in-Supply-Chain/server"
)

func main() {
	// Initialize Agents
	ollamaClient := agent.NewOllamaClient("http://localhost:11434", "deepseek-r1:1.5b")
	
	warehouse := &agent.Warehouse{
		Agent: agent.Agent{ID: "Warehouse-Central", Name: "Central Warehouse", Card: agent.GetWarehouseCard()},
		Inventory: map[string]int{
			"Electronics": 1000,
			"Groceries":   5000,
			"Clothing":    2000,
		},
		Client: ollamaClient,
	}

	outlets := []*agent.Outlet{
		{Agent: agent.Agent{ID: "Outlet-1", Name: "Outlet North", Card: agent.GetOutletCard("Outlet-1", "Outlet North", "North District")}, Inventory: map[string]int{"Electronics": 50, "Groceries": 100}, Events: agent.GlobalEvents, Client: ollamaClient},
		{Agent: agent.Agent{ID: "Outlet-2", Name: "Outlet South", Card: agent.GetOutletCard("Outlet-2", "Outlet South", "South District")}, Inventory: map[string]int{"Electronics": 60, "Groceries": 200}, Events: agent.GlobalEvents, Client: ollamaClient},
		{Agent: agent.Agent{ID: "Outlet-3", Name: "Outlet East", Card: agent.GetOutletCard("Outlet-3", "Outlet East", "East District")},  Inventory: map[string]int{"Electronics": 40, "Groceries": 150}, Events: agent.GlobalEvents, Client: ollamaClient},
		{Agent: agent.Agent{ID: "Outlet-4", Name: "Outlet West", Card: agent.GetOutletCard("Outlet-4", "Outlet West", "West District")},  Inventory: map[string]int{"Electronics": 70, "Groceries": 120}, Events: agent.GlobalEvents, Client: ollamaClient},
	}

	srv := server.NewServer(warehouse, outlets)
	go srv.HandleMessages()

	// Simulation Loop
	go func() {
		// Start Date slightly before the first event so user catches it
		simDate := time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC)
		ticker := time.NewTicker(10 * time.Second) // Slower: 10 seconds per day for better observation

		for range ticker.C {
			simDate = simDate.AddDate(0, 0, 1)
			log.Printf("Sim Date: %s", simDate.Format("2006-01-02"))

			// Check for Events
			for _, outlet := range outlets {
				for _, event := range agent.GlobalEvents {
					// Check if event matches both Date and OutletID
					if event.Date.Equal(simDate) && event.OutletID == outlet.ID {
						log.Printf("Event triggered for %s: %s", outlet.Name, event.Description)
						
						// === GOOGLE A2A PROTOCOL: MULTI-TURN CONVERSATION ===
						
						// TURN 1: Outlet Agent analyzes event and formulates request
						outletPrompt := fmt.Sprintf(`You are the Event Logistics Coordinator for %s.
Date: %s | Event: %s
Your Current Inventory: Electronics=%d, Groceries=%d

Analyze this event and formulate a request to the Central Warehouse. Consider:
1. Expected demand increase from this event
2. Your current stock levels
3. What items you'll need

Return ONLY valid JSON in this exact format (no extra text):
{"analysis": "your brief analysis here", "request": {"Electronics": 10, "Groceries": 20}}`,
							outlet.Name, simDate.Format("2006-01-02"), event.Description,
							outlet.Inventory["Electronics"], outlet.Inventory["Groceries"])
						
						outletResponse, err := outlet.Client.Generate(outletPrompt)
						if err != nil {
							log.Printf("Outlet AI Error: %v", err)
							continue
						}
						
						// Parse outlet request
						outletJSON := extractJSON(outletResponse)
						var outletReq struct {
							Analysis string           `json:"analysis"`
							Request  map[string]int `json:"request"`
						}
						if err := json.Unmarshal([]byte(outletJSON), &outletReq); err != nil {
							log.Printf("[ERROR] Failed to parse outlet JSON: %v, JSON: %s", err, outletJSON)
							continue
						}
						
						// Debug logging
						log.Printf("[DEBUG] Outlet %s - Analysis: %s, Request: %v", outlet.Name, outletReq.Analysis, outletReq.Request)
						
						// Log Turn 1
						srv.SendUpdate(map[string]interface{}{
							"type": "log",
							"message": fmt.Sprintf("🏪 [%s] %s\n  Event: %s\n  Analysis: %s\n  Request: %v", 
								simDate.Format("2006-01-02"), outlet.Name, event.Description, outletReq.Analysis, outletReq.Request),
						})
						
						// TURN 2: Warehouse Agent evaluates and proposes solution
						warehousePrompt := fmt.Sprintf(`You are the Central Warehouse Manager.
Outlet %s requests: %v
Their analysis: %s
Warehouse Stock: Electronics=%d, Groceries=%d

Evaluate this request and propose a solution. Consider:
1. Available warehouse inventory
2. Fairness to other outlets
3. Optimal replenishment amount

Return ONLY valid JSON in this exact format (no extra text):
{"reasoning": "your brief reasoning here", "proposal": {"Electronics": 10, "Groceries": 20}}`,
							outlet.Name, outletReq.Request, outletReq.Analysis,
							warehouse.Inventory["Electronics"], warehouse.Inventory["Groceries"])
						
						warehouseResponse, err := warehouse.Client.Generate(warehousePrompt)
						if err != nil {
							log.Printf("Warehouse AI Error: %v", err)
							continue
						}
						
						// Parse warehouse proposal
						warehouseJSON := extractJSON(warehouseResponse)
						var warehouseProposal struct {
							Reasoning string           `json:"reasoning"`
							Proposal  map[string]int `json:"proposal"`
						}
						if err := json.Unmarshal([]byte(warehouseJSON), &warehouseProposal); err != nil {
							log.Printf("[ERROR] Failed to parse warehouse JSON: %v, JSON: %s", err, warehouseJSON)
							continue
						}
						
						// Debug logging
						log.Printf("[DEBUG] Warehouse - Reasoning: %s, Proposal: %v", warehouseProposal.Reasoning, warehouseProposal.Proposal)
						
						// Log Turn 2
						srv.SendUpdate(map[string]interface{}{
							"type": "log",
							"message": fmt.Sprintf("🏭 Warehouse → %s\n  Reasoning: %s\n  Proposal: %v", 
								outlet.Name, warehouseProposal.Reasoning, warehouseProposal.Proposal),
						})
						
						// TURN 3: Outlet Agent accepts or negotiates
						negotiationPrompt := fmt.Sprintf(`You requested: %v
Warehouse proposed: %v
Their reasoning: %s

Do you accept this proposal? If not, negotiate. Return ONLY valid JSON (no extra text):
{"decision": "accept", "message": "brief message", "counter_offer": {"Electronics": 0, "Groceries": 0}}`,
							outletReq.Request, warehouseProposal.Proposal, warehouseProposal.Reasoning)
						
						negotiationResponse, err := outlet.Client.Generate(negotiationPrompt)
						if err != nil {
							log.Printf("Negotiation AI Error: %v", err)
							continue
						}
						
						negotiationJSON := extractJSON(negotiationResponse)
						var negotiation struct {
							Decision     string           `json:"decision"`
							Message      string           `json:"message"`
							CounterOffer map[string]int `json:"counter_offer"`
						}
						if err := json.Unmarshal([]byte(negotiationJSON), &negotiation); err != nil {
							log.Printf("[ERROR] Failed to parse negotiation JSON: %v, JSON: %s", err, negotiationJSON)
							continue
						}
						
						// Debug logging
						log.Printf("[DEBUG] Negotiation - Decision: %s, Message: %s", negotiation.Decision, negotiation.Message)
						
						// Log Turn 3
						srv.SendUpdate(map[string]interface{}{
							"type": "log",
							"message": fmt.Sprintf("🏪 %s: %s - %s", outlet.Name, negotiation.Decision, negotiation.Message),
						})
						
						// Final decision
						finalDecision := warehouseProposal.Proposal
						if negotiation.Decision == "negotiate" && len(negotiation.CounterOffer) > 0 {
							// Simple: warehouse accepts counter-offer if reasonable
							finalDecision = negotiation.CounterOffer
							srv.SendUpdate(map[string]interface{}{
								"type": "log",
								"message": fmt.Sprintf("🏭 Warehouse: Accepting counter-offer %v", finalDecision),
							})
						}
						
						// Execute inventory transfer
						for item, qty := range finalDecision {
							if qty > 0 && warehouse.Inventory[item] >= qty {
								warehouse.Inventory[item] -= qty
								outlet.Inventory[item] += qty
								
								// Send Animation Event
								srv.SendUpdate(map[string]interface{}{
									"type": "shipment",
									"from": "Warehouse",
									"to": outlet.ID,
									"item": item,
									"qty": qty,
								})
							}
						}
						
						srv.SendUpdate(map[string]interface{}{
							"type": "log",
							"message": fmt.Sprintf("✅ %s: Agreement reached - Transfer complete\n---", outlet.Name),
						})
					}
				}
			}

			// Periodic updates
			srv.SendUpdate(map[string]interface{}{
				"type": "state",
				"warehouse": warehouse,
				"outlets": outlets,
				"date": simDate.Format("2006-01-02"),
			})
		}
	}()

	// Serve Frontend
	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", srv.HandleConnections)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// extractJSON extracts JSON from AI response (handles markdown, think tags, etc.)
func extractJSON(response string) string {
	// Log raw response for debugging
	log.Printf("[DEBUG] Raw AI response (first 200 chars): %s", truncate(response, 200))
	
	// Remove <think> tags if present (deepseek-r1 specific)
	response = strings.ReplaceAll(response, "<think>", "")
	response = strings.ReplaceAll(response, "</think>", "")
	
	// Remove markdown code blocks
	response = strings.ReplaceAll(response, "```json", "")
	response = strings.ReplaceAll(response, "```", "")
	
	// Trim whitespace
	response = strings.TrimSpace(response)
	
	// Find ALL JSON objects and merge them
	// deepseek-r1 sometimes returns multiple JSON objects
	var allJSON []string
	for {
		startIdx := strings.Index(response, "{")
		if startIdx == -1 {
			break
		}
		
		// Find matching closing brace
		braceCount := 0
		endIdx := -1
		for i := startIdx; i < len(response); i++ {
			if response[i] == '{' {
				braceCount++
			} else if response[i] == '}' {
				braceCount--
				if braceCount == 0 {
					endIdx = i
					break
				}
			}
		}
		
		if endIdx != -1 {
			allJSON = append(allJSON, response[startIdx:endIdx+1])
			response = response[endIdx+1:]
		} else {
			break
		}
	}
	
	// If we found multiple JSON objects, try to merge them
	if len(allJSON) > 1 {
		merged := mergeJSON(allJSON)
		log.Printf("[DEBUG] Merged %d JSON objects into: %s", len(allJSON), merged)
		return merged
	} else if len(allJSON) == 1 {
		return allJSON[0]
	}
	
	log.Printf("[WARN] No JSON found in response, returning empty object")
	return "{}"
}

// mergeJSON merges multiple JSON objects into one
func mergeJSON(jsonObjects []string) string {
	merged := make(map[string]interface{})
	
	for _, jsonStr := range jsonObjects {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &obj); err == nil {
			for k, v := range obj {
				merged[k] = v
			}
		}
	}
	
	result, _ := json.Marshal(merged)
	return string(result)
}

// truncate helper function
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

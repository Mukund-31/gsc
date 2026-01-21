package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Appointat/Responsive-AI-Clusters-in-Supply-Chain/agent"
	"github.com/a2aproject/a2a-go/a2a"
	"github.com/a2aproject/a2a-go/a2asrv"
	"github.com/a2aproject/a2a-go/a2asrv/eventqueue"
)

type WarehouseExecutor struct {
	inventory map[string]int
	client    *agent.OllamaClient
	card      *a2a.AgentCard
}

func NewWarehouseExecutor() *WarehouseExecutor {
	card := &a2a.AgentCard{
		URL:  "http://localhost:8081/invoke",
		Name: "Central Medical Depot",
		Description: "Central hub for pharmaceutical and medical supply distribution.",
		PreferredTransport: a2a.TransportProtocolJSONRPC,
		ProtocolVersion: "0.3.4",
		Skills: []a2a.AgentSkill{
			{Name: "inventory_management", Description: "Monitor medical stockpile levels"},
			{Name: "replenishment", Description: "Dispatch emergency and routine medical supplies"},
			{Name: "negotiation", Description: "Prioritize allocation based on medical urgency"},
		},
	}
	return &WarehouseExecutor{
		inventory: map[string]int{
			"Antibiotics": 5000,
			"Painkillers": 10000,
			"Vaccines":    2000,
			"Bandages":    5000,
		},
		client: agent.NewOllamaClient("http://localhost:11434", "deepseek-r1:1.5b"),
		card:   card,
	}
}

func (w *WarehouseExecutor) Execute(ctx context.Context, reqCtx *a2asrv.RequestContext, queue eventqueue.Queue) error {
	// 1. Get Input
	var inputJSON string
	if reqCtx.Message != nil {
		log.Printf("[Warehouse] Message has %d parts", len(reqCtx.Message.Parts))
		for i, part := range reqCtx.Message.Parts {
			log.Printf("[Warehouse] Part %d type: %T", i, part)
			// Try pointer type
			if txtPart, ok := part.(*a2a.TextPart); ok {
				inputJSON = txtPart.Text
				break
			}
			// Try value type
			if txtPart, ok := part.(a2a.TextPart); ok {
				inputJSON = txtPart.Text
				break
			}
		}
	} else {
		log.Printf("[Warehouse] Message is nil")
	}

	if inputJSON == "" {
		return fmt.Errorf("no input text found")
	}

	// 2. Parse Logic
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}

	taskType, _ := input["type"].(string) 
	log.Printf("[Warehouse] Executing task type: %s", taskType)

	var result map[string]interface{}
	var err error

	switch taskType {
	case "proposal_request":
		result, err = w.handleProposal(ctx, input)
	case "transfer_request":
		result, err = w.handleTransfer(ctx, input)
	default:
		log.Printf("Unknown task type: %s", taskType)
		result = map[string]interface{}{"status": "unknown"}
	}

	if err != nil {
		return err
	}

	// 3. Send Response
	resultBytes, _ := json.Marshal(result)
	
	respPart := &a2a.TextPart{
		Text: string(resultBytes),
	}

	evt := a2a.NewArtifactEvent(reqCtx, respPart)

	if err := queue.Write(ctx, evt); err != nil {
		log.Printf("[Warehouse] Failed to write artifact: %v", err)
		return err
	}
	
	log.Printf("[Warehouse] Artifact sent, signaling completion")
	
	// Signal task completion with Final=true to stop event processing
	doneEvt := a2a.NewStatusUpdateEvent(reqCtx, a2a.TaskStateCompleted, nil)
	doneEvt.Final = true
	if err := queue.Write(ctx, doneEvt); err != nil {
		log.Printf("[Warehouse] Failed to write status: %v", err)
	}
	
	return nil
}

// Cancel implements AgentExecutor interface
func (w *WarehouseExecutor) Cancel(ctx context.Context, reqCtx *a2asrv.RequestContext, queue eventqueue.Queue) error {
	log.Printf("[Warehouse] Task cancelled: %s", reqCtx.TaskID)
	// Implement any cleanup or cancellation logic here
	return nil
}

func (w *WarehouseExecutor) handleProposal(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	outletName, _ := input["outlet_name"].(string)
	requestData, _ := input["request"].(map[string]interface{})

	if len(requestData) == 0 {
		log.Printf("[Warehouse] Empty request from %s, skipping AI", outletName)
		return map[string]interface{}{
			"reasoning": "No items requested.",
			"proposal":  map[string]int{},
		}, nil
	}


	log.Printf("[Warehouse] Proposal request from %s: %v", outletName, requestData)

	prompt := fmt.Sprintf(`You are the Inventory Fulfillment Manager.
Goal: Satisfy the Facility Request using Depot Stock.

Directives:
1. If Request > 0, propose EXACTLY the requested amount. Do NOT send extra unless Critical.
2. If Request > Depot Stock, send Depot Stock (Partial fulfillment).
3. Do NOT send items that were not requested (Request count 0 = Send 0).
4. If Request < 0 (Return), accept it.

Examples:
- Request: {"A": 50}, Stock: {"A": 500} -> Proposal: {"A": 50}
- Request: {"A": 500}, Stock: {"A": 100} -> Proposal: {"A": 100}

Current Context:
Facility: %s
Request: %v
Depot Stock: %v

Return ONLY valid JSON:
{"reasoning": "Fulfilling exact request...", "proposal": {"Antibiotics": 50}}`, 
		outletName, requestData, w.inventory)

	resp, err := w.client.Generate(prompt)
	if err != nil {
		log.Printf("AI error: %v", err)
		return map[string]interface{}{
			"reasoning": "AI unavailable",
			"proposal": requestData,
		}, nil
	}

	jsonStr := extractJSON(resp)
	var aiResult map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &aiResult)

	return map[string]interface{}{
		"type": "proposal",
		"reasoning": aiResult["reasoning"],
		"proposal": aiResult["proposal"],
	}, nil
}

func (w *WarehouseExecutor) handleTransfer(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	final, _ := input["final_decision"].(map[string]interface{})
	transferred := make(map[string]int)
	
	for k, v := range final {
		qty := int(v.(float64))
		if qty > 0 && w.inventory[k] >= qty {
			w.inventory[k] -= qty
			transferred[k] = qty
		}
	}
	log.Printf("[Warehouse] Transferred: %v", transferred)
	return map[string]interface{}{
		"type": "transfer_complete",
		"transferred": transferred,
		"inventory": w.inventory,
	}, nil
}

func extractJSON(s string) string {
	// Remove thinking tags
	s = removeThinkingTags(s)
	
	start, end := -1, -1
	brace := 0
	for i, c := range s {
		if c == '{' {
			if start == -1 { start = i }
			brace++
		}
		if c == '}' {
			brace--
			if brace == 0 && start != -1 {
				end = i + 1
				break
			}
		}
	}
	if start != -1 && end != -1 { 
		result := s[start:end]
		log.Printf("[Warehouse] Extracted JSON: %s", result[:min(100, len(result))])
		return result
	}
	return "{}"
}

func removeThinkingTags(s string) string {
	for {
		start := strings.Index(s, "<think>")
		if start == -1 {
			break
		}
		end := strings.Index(s, "</think>")
		if end == -1 {
			break
		}
		s = s[:start] + s[end+8:]
	}
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	exec := NewWarehouseExecutor()
	handler := a2asrv.NewHandler(exec)
	// NewJSONRPCHandler returns http.Handler
	jsonRPC := a2asrv.NewJSONRPCHandler(handler)

	mux := http.NewServeMux()
	mux.Handle("/.well-known/agent-card.json", a2asrv.NewStaticAgentCardHandler(exec.card))
	mux.Handle("/invoke", jsonRPC)

	port := "8081"
	log.Printf("Warehouse running on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Appointat/Responsive-AI-Clusters-in-Supply-Chain/agent"
	"github.com/a2aproject/a2a-go/a2a"
	"github.com/a2aproject/a2a-go/a2asrv"
	"github.com/a2aproject/a2a-go/a2asrv/eventqueue"
)

type OutletExecutor struct {
	ID        string
	Name      string
	Inventory map[string]int
	client    *agent.OllamaClient
	card      *a2a.AgentCard
}

func NewOutletExecutor(id, name, port string) *OutletExecutor {
	card := &a2a.AgentCard{
		URL:  fmt.Sprintf("http://localhost:%s/invoke", port),
		Name: name,
		Description: "Medical facility managing local pharmaceutical inventory.",
		PreferredTransport: a2a.TransportProtocolJSONRPC,
		ProtocolVersion: "0.3.4",
		Skills: []a2a.AgentSkill{
			{Name: "inventory_analysis", Description: "Track critical medicine usage and predict shortages"},
			{Name: "negotiation", Description: "Request urgent medical supplies from Depot"},
		},
	}
	
	inv := map[string]int{
		"Antibiotics": 50,
		"Painkillers": 100,
		"Vaccines":    20,
		"Bandages":    200,
	}

	return &OutletExecutor{
		ID:        id,
		Name:      name,
		Inventory: inv,
		client:    agent.NewOllamaClient("http://localhost:11434", "deepseek-r1:1.5b"),
		card:      card,
	}
}

func (o *OutletExecutor) Execute(ctx context.Context, reqCtx *a2asrv.RequestContext, queue eventqueue.Queue) error {
	// 1. Get Input
	var inputJSON string
	if reqCtx.Message != nil {
		log.Printf("[%s] Message has %d parts", o.Name, len(reqCtx.Message.Parts))
		for i, part := range reqCtx.Message.Parts {
			log.Printf("[%s] Part %d type: %T", o.Name, i, part)
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
		log.Printf("[%s] Message is nil", o.Name)
	}
	
	if inputJSON == "" {
		return fmt.Errorf("no input text found")
	}

	var input map[string]interface{}
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return fmt.Errorf("failed to input: %w", err)
	}

	taskType, _ := input["type"].(string)
	log.Printf("[%s] Executing task: %s", o.Name, taskType)

	var result map[string]interface{}
	var err error

	switch taskType {
	case "analyze_event":
		result, err = o.handleAnalysis(ctx, input)
	case "negotiate_proposal":
		result, err = o.handleNegotiation(ctx, input)
	default:
		log.Printf("Unknown task: %s", taskType)
		result = map[string]interface{}{"status": "unknown"}
	}

	if err != nil {
		log.Printf("[%s] Task error: %v", o.Name, err)
		return err
	}

	resultBytes, _ := json.Marshal(result)
	log.Printf("[%s] Sending response: %s", o.Name, string(resultBytes)[:min(100, len(resultBytes))])
	
	respPart := &a2a.TextPart{Text: string(resultBytes)}
	evt := a2a.NewArtifactEvent(reqCtx, respPart)
	
	if err := queue.Write(ctx, evt); err != nil {
		log.Printf("[%s] Failed to write artifact: %v", o.Name, err)
		return err
	}
	
	log.Printf("[%s] Artifact sent, closing queue", o.Name)
	
	// Signal task completion with Final=true to stop event processing
	doneEvt := a2a.NewStatusUpdateEvent(reqCtx, a2a.TaskStateCompleted, nil)
	doneEvt.Final = true
	if err := queue.Write(ctx, doneEvt); err != nil {
		log.Printf("[%s] Failed to write status: %v", o.Name, err)
	}
	
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Cancel implements AgentExecutor interface
func (o *OutletExecutor) Cancel(ctx context.Context, reqCtx *a2asrv.RequestContext, queue eventqueue.Queue) error {
	log.Printf("[%s] Task cancelled: %s", o.Name, reqCtx.TaskID)
	return nil
}

func (o *OutletExecutor) handleAnalysis(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	evt, _ := input["event_description"].(string)
	
	prompt := fmt.Sprintf(`You are Medical Facility Manager: %s.
Event: %s
Current Stock: %v

Analyze specific needs (Antibiotics, Painkillers, Vaccines, Bandages).
If Event mentions 'Recall' or 'Shortage', quantify items to RETURN to Depot as NEGATIVE numbers (e.g., "Antibiotics": -50).
Return ONLY valid JSON:
{"analysis": "...", "request": {"Antibiotics": 100, "Bandages": 50}}`,
		o.Name, evt, o.Inventory)

	resp, err := o.client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	jsonStr := extractJSON(resp)
	var aiResult map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &aiResult)

	return map[string]interface{}{
		"type": "analysis_result",
		"analysis": aiResult["analysis"],
		"request": aiResult["request"],
	}, nil
}

func (o *OutletExecutor) handleNegotiation(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	proposal, _ := input["proposal"].(map[string]interface{})
	reasoning, _ := input["reasoning"].(string)
	myReq, _ := input["original_request"].(map[string]interface{})

	prompt := fmt.Sprintf(`You are Medical Facility Manager: %s.
You requested: %v
Depot proposed: %v
Reasoning: %s

Accept (or counter if shortage remains). Return ONLY valid JSON:
{"decision": "accept", "message": "...", "counter_offer": {"Antibiotics": 0, "Bandages": 0}}`,
		o.Name, myReq, proposal, reasoning)

	resp, err := o.client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	jsonStr := extractJSON(resp)
	var aiResult map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &aiResult)

	return map[string]interface{}{
		"type": "negotiation_result",
		"decision": aiResult["decision"],
		"message": aiResult["message"],
		"counter_offer": aiResult["counter_offer"],
	}, nil
}

func extractJSON(s string) string {
	// Remove thinking tags first
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
		log.Printf("[extractJSON] Extracted: %s", result[:min(100, len(result))])
		return result
	}
	log.Printf("[extractJSON] Failed to find JSON in: %s", s[:min(200, len(s))])
	return "{}"
}

func removeThinkingTags(s string) string {
	// Remove <think>...</think> tags that deepseek-r1 produces
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

func main() {
	id := flag.String("id", "Outlet-1", "Outlet ID")
	name := flag.String("name", "Outlet North", "Outlet Name")
	port := flag.String("port", "8082", "Port to listen on")
	flag.Parse()

	exec := NewOutletExecutor(*id, *name, *port)
	handler := a2asrv.NewHandler(exec)
	jsonRPC := a2asrv.NewJSONRPCHandler(handler)

	mux := http.NewServeMux()
	mux.Handle("/.well-known/agent-card.json", a2asrv.NewStaticAgentCardHandler(exec.card))
	mux.Handle("/invoke", jsonRPC)
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("🏪 Outlet %s running on :%s", *name, *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatal(err)
	}
}

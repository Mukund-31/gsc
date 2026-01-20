package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// InventoryItem represents an item in the inventory
type InventoryItem struct {
	Name     string
	Quantity int
}

// ConversationTurn represents one turn in an A2A conversation
type ConversationTurn struct {
	Role    string `json:"role"` // "outlet" or "warehouse"
	Message string `json:"message"`
}

// Agent represents a base agent in the supply chain
type Agent struct {
	ID   string
	Name string
	Card *AgentCard
}

// Warehouse represents the central warehouse agent
type Warehouse struct {
	Agent
	Inventory           map[string]int
	Client              *OllamaClient
	ConversationHistory []ConversationTurn
}

// Outlet represents a retail outlet agent
type Outlet struct {
	Agent
	Inventory           map[string]int
	Events              []Event
	Client              *OllamaClient // Now outlets also have AI!
	ConversationHistory []ConversationTurn
}

// OllamaClient handles communication with Ollama API
type OllamaClient struct {
	BaseURL string
	Model   string
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		BaseURL: baseURL,
		Model:   model,
	}
}

// Generate sends a prompt to Ollama and returns the response
func (c *OllamaClient) Generate(prompt string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", c.BaseURL)
	
	payload := map[string]interface{}{
		"model":  c.Model,
		"prompt": prompt,
		"stream": false,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	
	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}
	
	return response, nil
}

package agent

import (
	"encoding/json"
	"time"
)

// AgentCard represents an agent's capabilities and identity (Google A2A Protocol)
type AgentCard struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"` // "warehouse" or "outlet"
	Capabilities []string               `json:"capabilities"`
	Location     string                 `json:"location,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
}

// GetWarehouseCard returns the agent card for the central warehouse
func GetWarehouseCard() *AgentCard {
	return &AgentCard{
		ID:   "Warehouse-Central",
		Name: "Central Warehouse",
		Type: "warehouse",
		Capabilities: []string{
			"inventory_management",
			"supply_coordination",
			"demand_forecasting",
			"multi_outlet_optimization",
		},
		Location: "Central Hub",
		Metadata: map[string]interface{}{
			"max_capacity":     10000,
			"ai_model":         "deepseek-r1:1.5b",
			"response_time_ms": 5000,
		},
		CreatedAt: time.Now(),
	}
}

// GetOutletCard returns the agent card for an outlet
func GetOutletCard(id, name, location string) *AgentCard {
	return &AgentCard{
		ID:   id,
		Name: name,
		Type: "outlet",
		Capabilities: []string{
			"event_analysis",
			"demand_prediction",
			"inventory_monitoring",
			"negotiation",
		},
		Location: location,
		Metadata: map[string]interface{}{
			"ai_model":         "deepseek-r1:1.5b",
			"response_time_ms": 5000,
		},
		CreatedAt: time.Now(),
	}
}

// ToJSON serializes the agent card to JSON
func (ac *AgentCard) ToJSON() (string, error) {
	data, err := json.MarshalIndent(ac, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

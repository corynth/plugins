package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Metadata represents plugin metadata
type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}

// InputSpec represents input parameter specification
type InputSpec struct {
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// OutputSpec represents output parameter specification
type OutputSpec struct {
	Type string `json:"type"`
}

// ActionSpec represents an action specification
type ActionSpec struct {
	Description string                `json:"description"`
	Inputs      map[string]InputSpec  `json:"inputs"`
	Outputs     map[string]OutputSpec `json:"outputs"`
}

// SlackPlugin represents the Slack plugin
type SlackPlugin struct {
	metadata   Metadata
	token      string
	webhookURL string
	client     *http.Client
}

// NewSlackPlugin creates a new Slack plugin instance
func NewSlackPlugin() *SlackPlugin {
	return &SlackPlugin{
		metadata: Metadata{
			Name:        "slack",
			Version:     "1.0.0",
			Description: "Slack workspace messaging and notifications",
			Author:      "Corynth Team",
			Tags:        []string{"slack", "messaging", "notifications"},
		},
		token:      os.Getenv("SLACK_BOT_TOKEN"),
		webhookURL: os.Getenv("SLACK_WEBHOOK_URL"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetMetadata returns plugin metadata
func (s *SlackPlugin) GetMetadata() Metadata {
	return s.metadata
}

// GetActions returns available actions
func (s *SlackPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"message": {
			Description: "Send message to channel",
			Inputs: map[string]InputSpec{
				"channel": {
					Type:        "string",
					Required:    true,
					Description: "Channel name or ID",
				},
				"text": {
					Type:        "string",
					Required:    true,
					Description: "Message text",
				},
				"username": {
					Type:        "string",
					Required:    false,
					Description: "Bot username",
				},
				"icon_emoji": {
					Type:        "string",
					Required:    false,
					Description: "Bot emoji icon",
				},
			},
			Outputs: map[string]OutputSpec{
				"success": {Type: "boolean"},
				"timestamp": {Type: "string"},
			},
		},
		"webhook": {
			Description: "Send webhook message",
			Inputs: map[string]InputSpec{
				"text": {
					Type:        "string",
					Required:    true,
					Description: "Message text",
				},
				"username": {
					Type:        "string",
					Required:    false,
					Description: "Bot username",
				},
				"channel": {
					Type:        "string",
					Required:    false,
					Description: "Override channel",
				},
			},
			Outputs: map[string]OutputSpec{
				"success": {Type: "boolean"},
			},
		},
	}
}

// Execute executes the specified action
func (s *SlackPlugin) Execute(action string, params map[string]interface{}) map[string]interface{} {
	switch action {
	case "message":
		return s.sendMessage(params)
	case "webhook":
		return s.sendWebhook(params)
	default:
		return map[string]interface{}{
			"error": fmt.Sprintf("Unknown action: %s", action),
		}
	}
}

// sendMessage sends a message using Slack Bot API
func (s *SlackPlugin) sendMessage(params map[string]interface{}) map[string]interface{} {
	if s.token == "" {
		return map[string]interface{}{
			"error": "SLACK_BOT_TOKEN not configured",
		}
	}

	// Extract parameters with defaults
	channel, _ := params["channel"].(string)
	text, _ := params["text"].(string)
	username, ok := params["username"].(string)
	if !ok {
		username = "Corynth Bot"
	}
	iconEmoji, ok := params["icon_emoji"].(string)
	if !ok {
		iconEmoji = ":robot_face:"
	}

	// Prepare request data
	data := map[string]interface{}{
		"channel":    channel,
		"text":       text,
		"username":   username,
		"icon_emoji": iconEmoji,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to marshal request data: %v", err),
		}
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to send request: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to read response: %v", err),
		}
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to parse response: %v", err),
		}
	}

	// Extract success and timestamp
	success, _ := result["ok"].(bool)
	timestamp, _ := result["ts"].(string)

	return map[string]interface{}{
		"success":   success,
		"timestamp": timestamp,
	}
}

// sendWebhook sends a message using Slack webhook
func (s *SlackPlugin) sendWebhook(params map[string]interface{}) map[string]interface{} {
	if s.webhookURL == "" {
		return map[string]interface{}{
			"error": "SLACK_WEBHOOK_URL not configured",
		}
	}

	// Extract parameters with defaults
	text, _ := params["text"].(string)
	username, ok := params["username"].(string)
	if !ok {
		username = "Corynth Bot"
	}

	// Prepare request data
	data := map[string]interface{}{
		"text":     text,
		"username": username,
	}

	// Add channel if specified
	if channel, ok := params["channel"].(string); ok && channel != "" {
		data["channel"] = channel
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to marshal request data: %v", err),
		}
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to create request: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to send request: %v", err),
		}
	}
	defer resp.Body.Close()

	return map[string]interface{}{
		"success": resp.StatusCode == 200,
	}
}

func main() {
	if len(os.Args) < 2 {
		result := map[string]interface{}{
			"error": "action required",
		}
		json.NewEncoder(os.Stdout).Encode(result)
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewSlackPlugin()

	var result interface{}

	switch action {
	case "metadata":
		result = plugin.GetMetadata()
	case "actions":
		result = plugin.GetActions()
	default:
		// Read parameters from stdin
		var params map[string]interface{}
		decoder := json.NewDecoder(os.Stdin)
		if err := decoder.Decode(&params); err != nil {
			// If no input or invalid JSON, use empty params
			params = make(map[string]interface{})
		}
		result = plugin.Execute(action, params)
	}

	json.NewEncoder(os.Stdout).Encode(result)
}
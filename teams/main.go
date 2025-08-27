package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/corynth/corynth/pkg/plugin"
	pluginv2 "github.com/corynth/corynth/pkg/plugin/v2"
)

type TeamsPlugin struct{}

func (p *TeamsPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "teams",
		Version:     "1.0.0",
		Description: "Microsoft Teams integration for sending messages, cards, and notifications via webhooks",
		Author:      "Corynth Team",
		Tags:        []string{"communication", "notifications", "microsoft", "collaboration"},
		License:     "MIT",
	}
}

func (p *TeamsPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "send_message",
			Description: "Send a simple text message to Teams channel",
			Inputs: map[string]plugin.InputSpec{
				"webhook_url": {
					Type:        "string",
					Description: "Teams channel webhook URL",
					Required:    true,
				},
				"text": {
					Type:        "string",
					Description: "Message text",
					Required:    true,
				},
				"title": {
					Type:        "string",
					Description: "Message title",
					Required:    false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"status": {
					Type:        "string",
					Description: "Response status",
				},
			},
		},
		{
			Name:        "send_alert",
			Description: "Send a formatted alert notification",
			Inputs: map[string]plugin.InputSpec{
				"webhook_url": {
					Type:        "string",
					Description: "Teams channel webhook URL",
					Required:    true,
				},
				"alert_type": {
					Type:        "string",
					Description: "Alert type: info, warning, error, success",
					Required:    false,
					Default:     "info",
				},
				"title": {
					Type:        "string",
					Description: "Alert title",
					Required:    true,
				},
				"message": {
					Type:        "string",
					Description: "Alert message",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"status": {
					Type:        "string",
					Description: "Response status",
				},
			},
		},
	}
}

func (p *TeamsPlugin) Validate(params map[string]interface{}) error {
	webhookURL, hasURL := params["webhook_url"].(string)
	if !hasURL || webhookURL == "" {
		return fmt.Errorf("webhook_url is required")
	}
	
	if !strings.Contains(webhookURL, "outlook.office.com") && !strings.Contains(webhookURL, "outlook.office365.com") {
		return fmt.Errorf("invalid Teams webhook URL format")
	}
	
	return nil
}

func (p *TeamsPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "send_message":
		return p.sendMessage(ctx, params)
	case "send_alert":
		return p.sendAlert(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *TeamsPlugin) sendMessage(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	webhookURL := params["webhook_url"].(string)
	text := params["text"].(string)
	
	title := ""
	if t, ok := params["title"].(string); ok {
		title = t
	}

	message := map[string]interface{}{
		"@type":    "MessageCard",
		"@context": "http://schema.org/extensions",
		"text":     text,
	}
	
	if title != "" {
		message["title"] = title
	}

	return p.sendToTeams(webhookURL, message)
}

func (p *TeamsPlugin) sendAlert(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	webhookURL := params["webhook_url"].(string)
	title := params["title"].(string)
	message := params["message"].(string)
	
	alertType := "info"
	if at, ok := params["alert_type"].(string); ok {
		alertType = at
	}

	var themeColor string
	var titlePrefix string
	
	switch alertType {
	case "error":
		themeColor = "FF0000"
		titlePrefix = "ERROR: "
	case "warning":
		themeColor = "FFA500"
		titlePrefix = "WARNING: "
	case "success":
		themeColor = "00FF00"
		titlePrefix = "SUCCESS: "
	default:
		themeColor = "0078D4"
		titlePrefix = "INFO: "
	}

	card := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "http://schema.org/extensions",
		"title":      titlePrefix + title,
		"text":       message,
		"themeColor": themeColor,
	}

	return p.sendToTeams(webhookURL, card)
}

func (p *TeamsPlugin) sendToTeams(webhookURL string, payload map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Teams API error: %s", resp.Status)
	}

	return map[string]interface{}{
		"status":      "success",
		"status_code": resp.StatusCode,
	}, nil
}

var ExportedPlugin plugin.Plugin = &TeamsPlugin{}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		if err := pluginv2.ServePlugin(&TeamsPlugin{}); err != nil {
			fmt.Printf("Error serving plugin: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Corynth Teams Plugin v1.0.0\n")
		fmt.Printf("Usage: %s serve\n", os.Args[0])
	}
}
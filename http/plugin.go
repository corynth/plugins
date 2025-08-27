package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	
	"github.com/corynth/corynth-dist/pkg/plugin"
)

type HttpPlugin struct{}

func (p *HttpPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "http",
		Version:     "1.0.0",
		Description: "HTTP client for REST API calls and web requests",
		Author:      "Corynth Team",
		Tags:        []string{"http", "api", "rest", "web", "client"},
		License:     "Apache-2.0",
	}
}

func (p *HttpPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "get",
			Description: "Make an HTTP GET request",
			Inputs: map[string]plugin.InputSpec{
				"url": {
					Type:        "string",
					Description: "URL to request",
					Required:    true,
				},
				"headers": {
					Type:        "object",
					Description: "HTTP headers",
					Required:    false,
				},
				"timeout": {
					Type:        "number",
					Description: "Request timeout in seconds",
					Required:    false,
					Default:     30,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"status_code": {
					Type:        "number",
					Description: "HTTP status code",
				},
				"body": {
					Type:        "string",
					Description: "Response body",
				},
				"headers": {
					Type:        "object",
					Description: "Response headers",
				},
			},
		},
		{
			Name:        "post",
			Description: "Make an HTTP POST request",
			Inputs: map[string]plugin.InputSpec{
				"url": {
					Type:        "string",
					Description: "URL to post to",
					Required:    true,
				},
				"body": {
					Type:        "string",
					Description: "Request body",
					Required:    false,
				},
				"headers": {
					Type:        "object",
					Description: "HTTP headers",
					Required:    false,
				},
				"timeout": {
					Type:        "number",
					Description: "Request timeout in seconds",
					Required:    false,
					Default:     30,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"status_code": {
					Type:        "number",
					Description: "HTTP status code",
				},
				"body": {
					Type:        "string",
					Description: "Response body",
				},
			},
		},
	}
}

func (p *HttpPlugin) Validate(params map[string]interface{}) error {
	return nil
}

func (p *HttpPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "get":
		return p.executeGet(ctx, params)
	case "post":
		return p.executePost(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *HttpPlugin) executeGet(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url, ok := params["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("url parameter is required")
	}

	// Set timeout
	timeout := 30
	if t, ok := params["timeout"].(float64); ok {
		timeout = int(t)
	}

	// Create client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers if provided
	if headers, ok := params["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			req.Header.Set(key, fmt.Sprintf("%v", value))
		}
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert response headers to map
	responseHeaders := make(map[string]interface{})
	for key, values := range resp.Header {
		if len(values) == 1 {
			responseHeaders[key] = values[0]
		} else {
			responseHeaders[key] = values
		}
	}

	return map[string]interface{}{
		"status_code": resp.StatusCode,
		"body":        string(body),
		"headers":     responseHeaders,
		"url":         url,
		"method":      "GET",
	}, nil
}

func (p *HttpPlugin) executePost(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	url, ok := params["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("url parameter is required")
	}

	// Get request body
	requestBody := ""
	if b, ok := params["body"].(string); ok {
		requestBody = b
	}

	// Set timeout
	timeout := 30
	if t, ok := params["timeout"].(float64); ok {
		timeout = int(t)
	}

	// Create client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Create request with body
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default Content-Type if not provided
	req.Header.Set("Content-Type", "application/json")

	// Set headers if provided
	if headers, ok := params["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			req.Header.Set(key, fmt.Sprintf("%v", value))
		}
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert response headers to map
	responseHeaders := make(map[string]interface{})
	for key, values := range resp.Header {
		if len(values) == 1 {
			responseHeaders[key] = values[0]
		} else {
			responseHeaders[key] = values
		}
	}

	return map[string]interface{}{
		"status_code": resp.StatusCode,
		"body":        string(responseBody),
		"headers":     responseHeaders,
		"url":         url,
		"method":      "POST",
	}, nil
}

var ExportedPlugin plugin.Plugin = &HttpPlugin{}
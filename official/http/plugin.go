package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}

type IOSpec struct {
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
}

type ActionSpec struct {
	Description string            `json:"description"`
	Inputs      map[string]IOSpec `json:"inputs"`
	Outputs     map[string]IOSpec `json:"outputs"`
}

type HTTPPlugin struct {
	client *http.Client
}

func NewHTTPPlugin() *HTTPPlugin {
	return &HTTPPlugin{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *HTTPPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "http",
		Version:     "1.0.0",
		Description: "HTTP client for REST API calls and web requests",
		Author:      "Corynth Team",
		Tags:        []string{"http", "web", "api", "rest"},
	}
}

func (p *HTTPPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"get": {
			Description: "Make HTTP GET requests with headers",
			Inputs: map[string]IOSpec{
				"url":     {Type: "string", Required: true, Description: "Request URL"},
				"headers": {Type: "object", Required: false, Description: "HTTP headers"},
				"timeout": {Type: "number", Required: false, Default: 30, Description: "Request timeout in seconds"},
				"auth":    {Type: "object", Required: false, Description: "Basic auth with username/password"},
			},
			Outputs: map[string]IOSpec{
				"status_code": {Type: "number", Description: "HTTP status code"},
				"headers":     {Type: "object", Description: "Response headers"},
				"content":     {Type: "string", Description: "Response body"},
				"json":        {Type: "object", Description: "Parsed JSON response (if applicable)"},
			},
		},
		"post": {
			Description: "Make HTTP POST requests with JSON data",
			Inputs: map[string]IOSpec{
				"url":          {Type: "string", Required: true, Description: "Request URL"},
				"headers":      {Type: "object", Required: false, Description: "HTTP headers"},
				"body":         {Type: "string", Required: false, Description: "Request body as string"},
				"json":         {Type: "object", Required: false, Description: "Request body as JSON"},
				"timeout":      {Type: "number", Required: false, Default: 30, Description: "Request timeout in seconds"},
				"auth":         {Type: "object", Required: false, Description: "Basic auth with username/password"},
				"content_type": {Type: "string", Required: false, Default: "application/json", Description: "Content-Type header"},
			},
			Outputs: map[string]IOSpec{
				"status_code": {Type: "number", Description: "HTTP status code"},
				"headers":     {Type: "object", Description: "Response headers"},
				"content":     {Type: "string", Description: "Response body"},
				"json":        {Type: "object", Description: "Parsed JSON response (if applicable)"},
			},
		},
	}
}

func (p *HTTPPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "get":
		return p.makeGetRequest(params)
	case "post":
		return p.makePostRequest(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *HTTPPlugin) makeGetRequest(params map[string]interface{}) (map[string]interface{}, error) {
	url, ok := params["url"].(string)
	if !ok || url == "" {
		return map[string]interface{}{"error": "url is required"}, nil
	}

	// Set timeout
	if timeout, ok := params["timeout"].(float64); ok {
		p.client.Timeout = time.Duration(timeout) * time.Second
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to create request: %v", err)}, nil
	}

	// Set headers
	if headers, ok := params["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	// Set authentication
	if auth, ok := params["auth"].(map[string]interface{}); ok {
		if username, hasUser := auth["username"].(string); hasUser {
			if password, hasPass := auth["password"].(string); hasPass {
				req.SetBasicAuth(username, password)
			}
		}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("request failed: %v", err)}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to read response: %v", err)}, nil
	}

	result := map[string]interface{}{
		"status_code": resp.StatusCode,
		"content":     string(body),
		"headers":     convertHeaders(resp.Header),
	}

	// Try to parse JSON
	if len(body) > 0 && strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var jsonData interface{}
		if json.Unmarshal(body, &jsonData) == nil {
			result["json"] = jsonData
		}
	}

	return result, nil
}

func (p *HTTPPlugin) makePostRequest(params map[string]interface{}) (map[string]interface{}, error) {
	url, ok := params["url"].(string)
	if !ok || url == "" {
		return map[string]interface{}{"error": "url is required"}, nil
	}

	// Set timeout
	if timeout, ok := params["timeout"].(float64); ok {
		p.client.Timeout = time.Duration(timeout) * time.Second
	}

	// Prepare request body
	var body io.Reader
	contentType := getStringParam(params, "content_type", "application/json")

	if jsonData, hasJSON := params["json"]; hasJSON {
		// JSON body
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to marshal JSON: %v", err)}, nil
		}
		body = bytes.NewReader(jsonBytes)
	} else if bodyStr, hasBody := params["body"].(string); hasBody {
		// String body
		body = strings.NewReader(bodyStr)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to create request: %v", err)}, nil
	}

	// Set Content-Type
	req.Header.Set("Content-Type", contentType)

	// Set headers
	if headers, ok := params["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	// Set authentication
	if auth, ok := params["auth"].(map[string]interface{}); ok {
		if username, hasUser := auth["username"].(string); hasUser {
			if password, hasPass := auth["password"].(string); hasPass {
				req.SetBasicAuth(username, password)
			}
		}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("request failed: %v", err)}, nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to read response: %v", err)}, nil
	}

	result := map[string]interface{}{
		"status_code": resp.StatusCode,
		"content":     string(respBody),
		"headers":     convertHeaders(resp.Header),
	}

	// Try to parse JSON response
	if len(respBody) > 0 && strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var jsonData interface{}
		if json.Unmarshal(respBody, &jsonData) == nil {
			result["json"] = jsonData
		}
	}

	return result, nil
}

// Helper functions
func getStringParam(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultValue
}

func convertHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewHTTPPlugin()

	var result interface{}

	switch action {
	case "metadata":
		result = plugin.GetMetadata()
	case "actions":
		result = plugin.GetActions()
	default:
		var params map[string]interface{}
		inputData, err := io.ReadAll(os.Stdin)
		if err != nil {
			result = map[string]interface{}{"error": fmt.Sprintf("failed to read input: %v", err)}
		} else if len(inputData) > 0 {
			if err := json.Unmarshal(inputData, &params); err != nil {
				result = map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}
			} else {
				result, err = plugin.Execute(action, params)
				if err != nil {
					result = map[string]interface{}{"error": err.Error()}
				}
			}
		} else {
			result, err = plugin.Execute(action, map[string]interface{}{})
			if err != nil {
				result = map[string]interface{}{"error": err.Error()}
			}
		}
	}

	json.NewEncoder(os.Stdout).Encode(result)
}
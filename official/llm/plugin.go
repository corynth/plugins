package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// LLMPlugin represents the LLM plugin
type LLMPlugin struct{}

// Metadata represents plugin metadata
type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}

// ActionInput represents an input parameter
type ActionInput struct {
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
}

// ActionOutput represents an output parameter
type ActionOutput struct {
	Type string `json:"type"`
}

// Action represents a plugin action
type Action struct {
	Description string                  `json:"description"`
	Inputs      map[string]ActionInput  `json:"inputs"`
	Outputs     map[string]ActionOutput `json:"outputs"`
}

// OpenAIRequest represents an OpenAI API request
type OpenAIRequest struct {
	Model       string              `json:"model"`
	Messages    []map[string]string `json:"messages"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
	Temperature float64             `json:"temperature,omitempty"`
}

// OpenAIResponse represents an OpenAI API response
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage map[string]interface{} `json:"usage,omitempty"`
}

// OllamaRequest represents an Ollama API request
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents an Ollama API response
type OllamaResponse struct {
	Response string `json:"response"`
}

// GetMetadata returns plugin metadata
func (p *LLMPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "llm",
		Version:     "1.0.0",
		Description: "Large Language Model integration (OpenAI, Ollama)",
		Author:      "Corynth Team",
		Tags:        []string{"llm", "ai", "gpt", "openai", "ollama"},
	}
}

// GetActions returns available actions
func (p *LLMPlugin) GetActions() map[string]Action {
	return map[string]Action{
		"generate": {
			Description: "Generate text using LLM",
			Inputs: map[string]ActionInput{
				"prompt": {
					Type:        "string",
					Required:    true,
					Description: "Input prompt",
				},
				"model": {
					Type:        "string",
					Required:    false,
					Default:     "gpt-3.5-turbo",
					Description: "Model name",
				},
				"max_tokens": {
					Type:        "number",
					Required:    false,
					Default:     150,
					Description: "Max tokens",
				},
				"temperature": {
					Type:        "number",
					Required:    false,
					Default:     0.7,
					Description: "Temperature",
				},
			},
			Outputs: map[string]ActionOutput{
				"text":  {Type: "string"},
				"usage": {Type: "object"},
			},
		},
		"chat": {
			Description: "Chat conversation",
			Inputs: map[string]ActionInput{
				"messages": {
					Type:        "array",
					Required:    true,
					Description: "Message history",
				},
				"model": {
					Type:        "string",
					Required:    false,
					Default:     "gpt-3.5-turbo",
					Description: "Model name",
				},
			},
			Outputs: map[string]ActionOutput{
				"response": {Type: "string"},
				"usage":    {Type: "object"},
			},
		},
		"ollama": {
			Description: "Use local Ollama model",
			Inputs: map[string]ActionInput{
				"prompt": {
					Type:        "string",
					Required:    true,
					Description: "Input prompt",
				},
				"model": {
					Type:        "string",
					Required:    false,
					Default:     "llama2",
					Description: "Ollama model name",
				},
			},
			Outputs: map[string]ActionOutput{
				"response": {Type: "string"},
			},
		},
	}
}

// Execute performs the specified action
func (p *LLMPlugin) Execute(action string, params map[string]interface{}) map[string]interface{} {
	switch action {
	case "generate":
		return p.openaiGenerate(params)
	case "chat":
		return p.openaiChat(params)
	case "ollama":
		return p.ollamaGenerate(params)
	default:
		return map[string]interface{}{"error": fmt.Sprintf("Unknown action: %s", action)}
	}
}

// openaiGenerate generates text using OpenAI API
func (p *LLMPlugin) openaiGenerate(params map[string]interface{}) map[string]interface{} {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return map[string]interface{}{"error": "OPENAI_API_KEY not configured"}
	}

	prompt, ok := params["prompt"].(string)
	if !ok {
		return map[string]interface{}{"error": "prompt is required"}
	}

	model := "gpt-3.5-turbo"
	if m, ok := params["model"].(string); ok {
		model = m
	}

	maxTokens := 150
	if mt, ok := params["max_tokens"]; ok {
		switch v := mt.(type) {
		case float64:
			maxTokens = int(v)
		case int:
			maxTokens = v
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				maxTokens = parsed
			}
		}
	}

	temperature := 0.7
	if temp, ok := params["temperature"]; ok {
		switch v := temp.(type) {
		case float64:
			temperature = v
		case int:
			temperature = float64(v)
		case string:
			if parsed, err := strconv.ParseFloat(v, 64); err == nil {
				temperature = parsed
			}
		}
	}

	request := OpenAIRequest{
		Model: model,
		Messages: []map[string]string{
			{"role": "user", "content": prompt},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to marshal request: %v", err)}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to create request: %v", err)}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("HTTP request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return map[string]interface{}{"error": fmt.Sprintf("API error (%d): %s", resp.StatusCode, string(body))}
	}

	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to decode response: %v", err)}
	}

	if len(openaiResp.Choices) == 0 {
		return map[string]interface{}{"error": "No response choices returned"}
	}

	return map[string]interface{}{
		"text":  openaiResp.Choices[0].Message.Content,
		"usage": openaiResp.Usage,
	}
}

// openaiChat handles chat conversations using OpenAI API
func (p *LLMPlugin) openaiChat(params map[string]interface{}) map[string]interface{} {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return map[string]interface{}{"error": "OPENAI_API_KEY not configured"}
	}

	messagesParam, ok := params["messages"]
	if !ok {
		return map[string]interface{}{"error": "messages are required"}
	}

	// Convert messages to the correct format
	var messages []map[string]string
	if msgSlice, ok := messagesParam.([]interface{}); ok {
		for _, msg := range msgSlice {
			if msgMap, ok := msg.(map[string]interface{}); ok {
				convertedMsg := make(map[string]string)
				for k, v := range msgMap {
					if str, ok := v.(string); ok {
						convertedMsg[k] = str
					}
				}
				messages = append(messages, convertedMsg)
			}
		}
	} else {
		return map[string]interface{}{"error": "messages must be an array"}
	}

	model := "gpt-3.5-turbo"
	if m, ok := params["model"].(string); ok {
		model = m
	}

	request := OpenAIRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to marshal request: %v", err)}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to create request: %v", err)}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("HTTP request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return map[string]interface{}{"error": fmt.Sprintf("API error (%d): %s", resp.StatusCode, string(body))}
	}

	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to decode response: %v", err)}
	}

	if len(openaiResp.Choices) == 0 {
		return map[string]interface{}{"error": "No response choices returned"}
	}

	return map[string]interface{}{
		"response": openaiResp.Choices[0].Message.Content,
		"usage":    openaiResp.Usage,
	}
}

// ollamaGenerate generates text using Ollama API
func (p *LLMPlugin) ollamaGenerate(params map[string]interface{}) map[string]interface{} {
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	prompt, ok := params["prompt"].(string)
	if !ok {
		return map[string]interface{}{"error": "prompt is required"}
	}

	model := "llama2"
	if m, ok := params["model"].(string); ok {
		model = m
	}

	request := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to marshal request: %v", err)}
	}

	client := &http.Client{Timeout: 120 * time.Second} // Longer timeout for local models
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/generate", ollamaURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to create request: %v", err)}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("HTTP request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return map[string]interface{}{"error": fmt.Sprintf("Ollama API error (%d): %s", resp.StatusCode, string(body))}
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("Failed to decode response: %v", err)}
	}

	return map[string]interface{}{
		"response": ollamaResp.Response,
	}
}

func main() {
	if len(os.Args) < 2 {
		result := map[string]interface{}{"error": "action required"}
		json.NewEncoder(os.Stdout).Encode(result)
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := &LLMPlugin{}

	var params map[string]interface{}
	
	// Always try to read from stdin
	input, err := io.ReadAll(os.Stdin)
	if err == nil && len(input) > 0 {
		// Trim whitespace and parse JSON
		trimmed := bytes.TrimSpace(input)
		if len(trimmed) > 0 {
			if err := json.Unmarshal(trimmed, &params); err != nil {
				// If JSON parsing fails, create error response
				result := map[string]interface{}{"error": fmt.Sprintf("Invalid JSON input: %v", err)}
				json.NewEncoder(os.Stdout).Encode(result)
				return
			}
		}
	}

	if params == nil {
		params = make(map[string]interface{})
	}

	var result interface{}
	switch action {
	case "metadata":
		result = plugin.GetMetadata()
	case "actions":
		result = plugin.GetActions()
	default:
		result = plugin.Execute(action, params)
	}

	json.NewEncoder(os.Stdout).Encode(result)
}
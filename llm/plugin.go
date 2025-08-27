package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os/exec"
    "strings"
    "time"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type LLMPlugin struct{}

func (p *LLMPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "llm",
        Version:     "1.0.0",
        Description: "Large Language Model integration with commercial APIs and self-hosted Ollama",
        Author:      "Corynth Team",
        Tags:        []string{"llm", "ai", "gpt", "claude", "ollama", "openai", "anthropic", "generation"},
        License:     "Apache-2.0",
    }
}

func (p *LLMPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "generate",
            Description: "Generate text using LLM",
            Inputs: map[string]plugin.InputSpec{
                "prompt": {
                    Type:        "string",
                    Description: "Input prompt for text generation",
                    Required:    true,
                },
                "provider": {
                    Type:        "string",
                    Description: "LLM provider (openai, anthropic, ollama)",
                    Required:    false,
                    Default:     "ollama",
                },
                "model": {
                    Type:        "string",
                    Description: "Model name (gpt-4, claude-3, llama2, etc.)",
                    Required:    false,
                    Default:     "llama2",
                },
                "temperature": {
                    Type:        "number",
                    Description: "Sampling temperature (0.0-1.0)",
                    Required:    false,
                    Default:     0.7,
                },
                "max_tokens": {
                    Type:        "number",
                    Description: "Maximum tokens to generate",
                    Required:    false,
                    Default:     1000,
                },
                "system_prompt": {
                    Type:        "string",
                    Description: "System prompt/instructions",
                    Required:    false,
                },
                "context": {
                    Type:        "array",
                    Description: "Previous conversation context",
                    Required:    false,
                },
                "api_key": {
                    Type:        "string",
                    Description: "API key for commercial providers",
                    Required:    false,
                },
                "base_url": {
                    Type:        "string",
                    Description: "Custom API base URL",
                    Required:    false,
                },
                "options": {
                    Type:        "object",
                    Description: "Additional provider-specific options",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "response": {
                    Type:        "string",
                    Description: "Generated text response",
                },
                "tokens_used": {
                    Type:        "number",
                    Description: "Number of tokens used",
                },
                "model": {
                    Type:        "string",
                    Description: "Model used for generation",
                },
                "provider": {
                    Type:        "string",
                    Description: "Provider used",
                },
                "context": {
                    Type:        "array",
                    Description: "Updated conversation context",
                },
            },
        },
        {
            Name:        "chat",
            Description: "Interactive chat with LLM",
            Inputs: map[string]plugin.InputSpec{
                "messages": {
                    Type:        "array",
                    Description: "Array of chat messages",
                    Required:    true,
                },
                "provider": {
                    Type:        "string",
                    Description: "LLM provider (openai, anthropic, ollama)",
                    Required:    false,
                    Default:     "ollama",
                },
                "model": {
                    Type:        "string",
                    Description: "Model name",
                    Required:    false,
                    Default:     "llama2",
                },
                "temperature": {
                    Type:        "number",
                    Description: "Sampling temperature (0.0-1.0)",
                    Required:    false,
                    Default:     0.7,
                },
                "max_tokens": {
                    Type:        "number",
                    Description: "Maximum tokens to generate",
                    Required:    false,
                    Default:     1000,
                },
                "api_key": {
                    Type:        "string",
                    Description: "API key for commercial providers",
                    Required:    false,
                },
                "base_url": {
                    Type:        "string",
                    Description: "Custom API base URL",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "response": {
                    Type:        "string",
                    Description: "Chat response",
                },
                "messages": {
                    Type:        "array",
                    Description: "Updated chat history",
                },
                "tokens_used": {
                    Type:        "number",
                    Description: "Number of tokens used",
                },
            },
        },
        {
            Name:        "embed",
            Description: "Generate embeddings for text",
            Inputs: map[string]plugin.InputSpec{
                "text": {
                    Type:        "string",
                    Description: "Text to embed",
                    Required:    true,
                },
                "provider": {
                    Type:        "string",
                    Description: "Embedding provider (openai, ollama)",
                    Required:    false,
                    Default:     "ollama",
                },
                "model": {
                    Type:        "string",
                    Description: "Embedding model name",
                    Required:    false,
                    Default:     "nomic-embed-text",
                },
                "api_key": {
                    Type:        "string",
                    Description: "API key for commercial providers",
                    Required:    false,
                },
                "base_url": {
                    Type:        "string",
                    Description: "Custom API base URL",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "embedding": {
                    Type:        "array",
                    Description: "Embedding vector",
                },
                "dimensions": {
                    Type:        "number",
                    Description: "Embedding dimensions",
                },
                "model": {
                    Type:        "string",
                    Description: "Model used for embedding",
                },
            },
        },
        {
            Name:        "list_models",
            Description: "List available models",
            Inputs: map[string]plugin.InputSpec{
                "provider": {
                    Type:        "string",
                    Description: "LLM provider (ollama, openai, anthropic)",
                    Required:    false,
                    Default:     "ollama",
                },
                "api_key": {
                    Type:        "string",
                    Description: "API key for commercial providers",
                    Required:    false,
                },
                "base_url": {
                    Type:        "string",
                    Description: "Custom API base URL",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "models": {
                    Type:        "array",
                    Description: "List of available models",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of models available",
                },
            },
        },
        {
            Name:        "pull_model",
            Description: "Pull/download model (Ollama only)",
            Inputs: map[string]plugin.InputSpec{
                "model": {
                    Type:        "string",
                    Description: "Model name to pull",
                    Required:    true,
                },
                "base_url": {
                    Type:        "string",
                    Description: "Ollama base URL",
                    Required:    false,
                    Default:     "http://localhost:11434",
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "status": {
                    Type:        "string",
                    Description: "Pull operation status",
                },
                "model": {
                    Type:        "string",
                    Description: "Pulled model name",
                },
            },
        },
    }
}

func (p *LLMPlugin) Validate(params map[string]interface{}) error {
    return nil
}

func (p *LLMPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "generate":
        return p.executeGenerate(ctx, params)
    case "chat":
        return p.executeChat(ctx, params)
    case "embed":
        return p.executeEmbed(ctx, params)
    case "list_models":
        return p.executeListModels(ctx, params)
    case "pull_model":
        return p.executePullModel(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *LLMPlugin) executeGenerate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    prompt, _ := params["prompt"].(string)
    provider, _ := params["provider"].(string)
    if provider == "" {
        provider = "ollama"
    }
    
    model, _ := params["model"].(string)
    if model == "" {
        if provider == "openai" {
            model = "gpt-3.5-turbo"
        } else if provider == "anthropic" {
            model = "claude-3-haiku-20240307"
        } else {
            model = "llama2"
        }
    }
    
    temperature, _ := params["temperature"].(float64)
    if temperature == 0 {
        temperature = 0.7
    }
    
    maxTokens, _ := params["max_tokens"].(float64)
    if maxTokens == 0 {
        maxTokens = 1000
    }
    
    systemPrompt, _ := params["system_prompt"].(string)
    apiKey, _ := params["api_key"].(string)
    baseURL, _ := params["base_url"].(string)
    
    switch provider {
    case "openai":
        return p.callOpenAI(ctx, prompt, model, temperature, int(maxTokens), systemPrompt, apiKey, baseURL)
    case "anthropic":
        return p.callAnthropic(ctx, prompt, model, temperature, int(maxTokens), systemPrompt, apiKey, baseURL)
    case "ollama":
        return p.callOllama(ctx, prompt, model, temperature, int(maxTokens), systemPrompt, baseURL)
    default:
        return nil, fmt.Errorf("unsupported provider: %s", provider)
    }
}

func (p *LLMPlugin) executeChat(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    messages, _ := params["messages"].([]interface{})
    provider, _ := params["provider"].(string)
    if provider == "" {
        provider = "ollama"
    }
    
    model, _ := params["model"].(string)
    if model == "" {
        model = "llama2"
    }
    
    temperature, _ := params["temperature"].(float64)
    if temperature == 0 {
        temperature = 0.7
    }
    
    maxTokens, _ := params["max_tokens"].(float64)
    if maxTokens == 0 {
        maxTokens = 1000
    }
    
    apiKey, _ := params["api_key"].(string)
    baseURL, _ := params["base_url"].(string)
    
    // Convert messages to chat format
    var prompt string
    if len(messages) > 0 {
        var chatBuilder strings.Builder
        for _, msg := range messages {
            if m, ok := msg.(map[string]interface{}); ok {
                role, _ := m["role"].(string)
                content, _ := m["content"].(string)
                chatBuilder.WriteString(fmt.Sprintf("%s: %s\n", role, content))
            }
        }
        prompt = chatBuilder.String()
    }
    
    result, err := p.executeGenerate(ctx, map[string]interface{}{
        "prompt":      prompt,
        "provider":    provider,
        "model":       model,
        "temperature": temperature,
        "max_tokens":  maxTokens,
        "api_key":     apiKey,
        "base_url":    baseURL,
    })
    
    if err != nil {
        return nil, err
    }
    
    // Add assistant response to messages
    response, _ := result["response"].(string)
    updatedMessages := append(messages, map[string]interface{}{
        "role":    "assistant",
        "content": response,
    })
    
    result["messages"] = updatedMessages
    return result, nil
}

func (p *LLMPlugin) executeEmbed(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    text, _ := params["text"].(string)
    provider, _ := params["provider"].(string)
    if provider == "" {
        provider = "ollama"
    }
    
    model, _ := params["model"].(string)
    if model == "" {
        model = "nomic-embed-text"
    }
    
    baseURL, _ := params["base_url"].(string)
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    
    switch provider {
    case "ollama":
        return p.callOllamaEmbed(ctx, text, model, baseURL)
    case "openai":
        apiKey, _ := params["api_key"].(string)
        return p.callOpenAIEmbed(ctx, text, model, apiKey, baseURL)
    default:
        return nil, fmt.Errorf("unsupported embedding provider: %s", provider)
    }
}

func (p *LLMPlugin) executeListModels(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    provider, _ := params["provider"].(string)
    if provider == "" {
        provider = "ollama"
    }
    
    baseURL, _ := params["base_url"].(string)
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    
    switch provider {
    case "ollama":
        return p.listOllamaModels(ctx, baseURL)
    default:
        // Static list for commercial providers
        var models []string
        if provider == "openai" {
            models = []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"}
        } else if provider == "anthropic" {
            models = []string{"claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"}
        }
        
        return map[string]interface{}{
            "models": models,
            "count":  len(models),
        }, nil
    }
}

func (p *LLMPlugin) executePullModel(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    model, _ := params["model"].(string)
    baseURL, _ := params["base_url"].(string)
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    
    return p.pullOllamaModel(ctx, model, baseURL)
}

func (p *LLMPlugin) callOpenAI(ctx context.Context, prompt, model string, temperature float64, maxTokens int, systemPrompt, apiKey, baseURL string) (map[string]interface{}, error) {
    if apiKey == "" {
        return nil, fmt.Errorf("OpenAI API key required")
    }
    
    if baseURL == "" {
        baseURL = "https://api.openai.com/v1"
    }
    
    messages := []map[string]interface{}{
        {"role": "user", "content": prompt},
    }
    
    if systemPrompt != "" {
        messages = append([]map[string]interface{}{
            {"role": "system", "content": systemPrompt},
        }, messages...)
    }
    
    payload := map[string]interface{}{
        "model":       model,
        "messages":    messages,
        "temperature": temperature,
        "max_tokens":  maxTokens,
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)
    
    client := &http.Client{Timeout: 60 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("OpenAI API error: %s", string(body))
    }
    
    choices, _ := result["choices"].([]interface{})
    if len(choices) == 0 {
        return nil, fmt.Errorf("no response from OpenAI")
    }
    
    choice, _ := choices[0].(map[string]interface{})
    message, _ := choice["message"].(map[string]interface{})
    content, _ := message["content"].(string)
    
    usage, _ := result["usage"].(map[string]interface{})
    totalTokens, _ := usage["total_tokens"].(float64)
    
    return map[string]interface{}{
        "response":    content,
        "tokens_used": int(totalTokens),
        "model":       model,
        "provider":    "openai",
    }, nil
}

func (p *LLMPlugin) callAnthropic(ctx context.Context, prompt, model string, temperature float64, maxTokens int, systemPrompt, apiKey, baseURL string) (map[string]interface{}, error) {
    if apiKey == "" {
        return nil, fmt.Errorf("Anthropic API key required")
    }
    
    if baseURL == "" {
        baseURL = "https://api.anthropic.com/v1"
    }
    
    payload := map[string]interface{}{
        "model":      model,
        "max_tokens": maxTokens,
        "messages": []map[string]interface{}{
            {"role": "user", "content": prompt},
        },
        "temperature": temperature,
    }
    
    if systemPrompt != "" {
        payload["system"] = systemPrompt
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/messages", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")
    
    client := &http.Client{Timeout: 60 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Anthropic API error: %s", string(body))
    }
    
    content, _ := result["content"].([]interface{})
    if len(content) == 0 {
        return nil, fmt.Errorf("no response from Anthropic")
    }
    
    contentBlock, _ := content[0].(map[string]interface{})
    text, _ := contentBlock["text"].(string)
    
    usage, _ := result["usage"].(map[string]interface{})
    inputTokens, _ := usage["input_tokens"].(float64)
    outputTokens, _ := usage["output_tokens"].(float64)
    
    return map[string]interface{}{
        "response":    text,
        "tokens_used": int(inputTokens + outputTokens),
        "model":       model,
        "provider":    "anthropic",
    }, nil
}

func (p *LLMPlugin) callOllama(ctx context.Context, prompt, model string, temperature float64, maxTokens int, systemPrompt, baseURL string) (map[string]interface{}, error) {
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    
    payload := map[string]interface{}{
        "model":  model,
        "prompt": prompt,
        "stream": false,
        "options": map[string]interface{}{
            "temperature":   temperature,
            "num_predict":   maxTokens,
        },
    }
    
    if systemPrompt != "" {
        payload["system"] = systemPrompt
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/generate", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{Timeout: 120 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Ollama API error: %s", string(body))
    }
    
    response, _ := result["response"].(string)
    
    return map[string]interface{}{
        "response":    response,
        "tokens_used": 0, // Ollama doesn't return token count in simple API
        "model":       model,
        "provider":    "ollama",
    }, nil
}

func (p *LLMPlugin) callOllamaEmbed(ctx context.Context, text, model, baseURL string) (map[string]interface{}, error) {
    payload := map[string]interface{}{
        "model": model,
        "prompt": text,
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/embeddings", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{Timeout: 60 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Ollama embedding error: %s", string(body))
    }
    
    embedding, _ := result["embedding"].([]interface{})
    
    return map[string]interface{}{
        "embedding":  embedding,
        "dimensions": len(embedding),
        "model":      model,
    }, nil
}

func (p *LLMPlugin) callOpenAIEmbed(ctx context.Context, text, model, apiKey, baseURL string) (map[string]interface{}, error) {
    if baseURL == "" {
        baseURL = "https://api.openai.com/v1"
    }
    
    payload := map[string]interface{}{
        "model": model,
        "input": text,
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/embeddings", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)
    
    client := &http.Client{Timeout: 60 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("OpenAI embedding error: %s", string(body))
    }
    
    data, _ := result["data"].([]interface{})
    if len(data) == 0 {
        return nil, fmt.Errorf("no embedding data returned")
    }
    
    embeddingData, _ := data[0].(map[string]interface{})
    embedding, _ := embeddingData["embedding"].([]interface{})
    
    return map[string]interface{}{
        "embedding":  embedding,
        "dimensions": len(embedding),
        "model":      model,
    }, nil
}

func (p *LLMPlugin) listOllamaModels(ctx context.Context, baseURL string) (map[string]interface{}, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/tags", nil)
    if err != nil {
        return nil, err
    }
    
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Ollama API error: %s", string(body))
    }
    
    models, _ := result["models"].([]interface{})
    modelNames := []string{}
    
    for _, model := range models {
        if m, ok := model.(map[string]interface{}); ok {
            if name, ok := m["name"].(string); ok {
                modelNames = append(modelNames, name)
            }
        }
    }
    
    return map[string]interface{}{
        "models": modelNames,
        "count":  len(modelNames),
    }, nil
}

func (p *LLMPlugin) pullOllamaModel(ctx context.Context, model, baseURL string) (map[string]interface{}, error) {
    payload := map[string]interface{}{
        "name":   model,
        "stream": false,
    }
    
    jsonPayload, _ := json.Marshal(payload)
    
    req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/pull", bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{Timeout: 300 * time.Second} // Longer timeout for model pulling
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    if resp.StatusCode == 200 {
        return map[string]interface{}{
            "status": "success",
            "model":  model,
        }, nil
    }
    
    return map[string]interface{}{
        "status": "failed",
        "model":  model,
        "error":  string(body),
    }, nil
}

var ExportedPlugin plugin.Plugin = &LLMPlugin{}
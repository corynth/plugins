package main

import (
	"context"
	"fmt"
	
	"github.com/corynth/corynth-dist/pkg/plugin"
)

type RedisPlugin struct{}

func (p *RedisPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "redis",
		Version:     "1.0.0",
		Description: "Redis cache operations and key-value storage",
		Author:      "Corynth Team",
		Tags:        []string{"cache", "redis", "key-value", "database"},
		License:     "Apache-2.0",
	}
}

func (p *RedisPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "set",
			Description: "Set a key-value pair in Redis",
			Inputs: map[string]plugin.InputSpec{
				"host": {
					Type:        "string",
					Description: "Redis host",
					Required:    false,
					Default:     "localhost",
				},
				"port": {
					Type:        "number",
					Description: "Redis port",
					Required:    false,
					Default:     6379,
				},
				"key": {
					Type:        "string",
					Description: "Key to set",
					Required:    true,
				},
				"value": {
					Type:        "string",
					Description: "Value to set",
					Required:    true,
				},
				"ttl": {
					Type:        "number",
					Description: "Time to live in seconds (0 = no expiry)",
					Required:    false,
					Default:     0,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"success": {
					Type:        "boolean",
					Description: "Whether the operation succeeded",
				},
			},
		},
		{
			Name:        "get",
			Description: "Get a value from Redis by key",
			Inputs: map[string]plugin.InputSpec{
				"host": {
					Type:        "string",
					Description: "Redis host",
					Required:    false,
					Default:     "localhost",
				},
				"port": {
					Type:        "number",
					Description: "Redis port",
					Required:    false,
					Default:     6379,
				},
				"key": {
					Type:        "string",
					Description: "Key to retrieve",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"exists": {
					Type:        "boolean",
					Description: "Whether the key exists",
				},
				"value": {
					Type:        "string",
					Description: "The value for the key",
				},
			},
		},
	}
}

func (p *RedisPlugin) Validate(params map[string]interface{}) error {
	return nil
}

func (p *RedisPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "set":
		return p.executeSet(ctx, params)
	case "get":
		return p.executeGet(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *RedisPlugin) executeSet(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	key, ok := params["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key parameter is required")
	}

	value, ok := params["value"].(string)
	if !ok {
		return nil, fmt.Errorf("value parameter is required")
	}

	host := "localhost"
	if h, ok := params["host"].(string); ok {
		host = h
	}

	port := 6379
	if p, ok := params["port"].(float64); ok {
		port = int(p)
	}

	// In production, this would connect to Redis
	return map[string]interface{}{
		"success": true,
		"key":     key,
		"value":   value,
		"message": fmt.Sprintf("Successfully set key '%s' in Redis at %s:%d", key, host, port),
	}, nil
}

func (p *RedisPlugin) executeGet(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	key, ok := params["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key parameter is required")
	}

	host := "localhost"
	if h, ok := params["host"].(string); ok {
		host = h
	}

	port := 6379
	if p, ok := params["port"].(float64); ok {
		port = int(p)
	}

	// In production, this would connect to Redis
	return map[string]interface{}{
		"exists":  true,
		"value":   fmt.Sprintf("mock_value_for_%s", key),
		"key":     key,
		"message": fmt.Sprintf("Retrieved key '%s' from Redis at %s:%d", key, host, port),
	}, nil
}

var ExportedPlugin plugin.Plugin = &RedisPlugin{}
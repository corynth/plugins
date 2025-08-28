// Corynth Go Plugin Template
// This is the official template for creating Corynth JSON protocol plugins in Go.
// Copy this file and modify it to create your own plugin.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Plugin metadata structure
type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}

// Action input/output specification
type IOSpec struct {
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
}

// Action specification
type ActionSpec struct {
	Description string            `json:"description"`
	Inputs      map[string]IOSpec `json:"inputs"`
	Outputs     map[string]IOSpec `json:"outputs"`
}

// YourPlugin - Replace with your plugin name
type YourPlugin struct {
	// Add any plugin state/configuration here
}

// NewYourPlugin creates a new plugin instance
func NewYourPlugin() *YourPlugin {
	return &YourPlugin{
		// Initialize any state here
	}
}

// GetMetadata returns plugin metadata
func (p *YourPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "yourplugin", // Change this
		Version:     "1.0.0",
		Description: "Your plugin description here",
		Author:      "Your Name",
		Tags:        []string{"tag1", "tag2"}, // Update tags
	}
}

// GetActions returns available actions
func (p *YourPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"example_action": {
			Description: "Example action that does something",
			Inputs: map[string]IOSpec{
				"input_param": {
					Type:        "string",
					Required:    true,
					Description: "An example input parameter",
				},
				"optional_param": {
					Type:        "number",
					Required:    false,
					Default:     10,
					Description: "An optional parameter with default",
				},
			},
			Outputs: map[string]IOSpec{
				"result": {
					Type:        "string",
					Description: "The action result",
				},
				"success": {
					Type:        "boolean",
					Description: "Whether the action succeeded",
				},
			},
		},
		// Add more actions here
	}
}

// Execute runs the specified action with given parameters
func (p *YourPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "example_action":
		return p.exampleAction(params)
	// Add more action cases here
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// exampleAction implements the example action
func (p *YourPlugin) exampleAction(params map[string]interface{}) (map[string]interface{}, error) {
	// Extract parameters with type checking
	inputParam, ok := params["input_param"].(string)
	if !ok {
		return map[string]interface{}{
			"error": "input_param must be a string",
		}, nil
	}

	// Get optional parameter with default
	optionalParam := 10.0
	if val, exists := params["optional_param"]; exists {
		if floatVal, ok := val.(float64); ok {
			optionalParam = floatVal
		}
	}

	// Implement your action logic here
	result := fmt.Sprintf("Processed %s with value %.2f", inputParam, optionalParam)

	// Return results
	return map[string]interface{}{
		"result":  result,
		"success": true,
	}, nil
}

// Helper function to get string parameter
func getStringParam(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultValue
}

// Helper function to get float parameter
func getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key].(float64); ok {
		return val
	}
	return defaultValue
}

// Helper function to get bool parameter
func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultValue
}

// Helper function to get array parameter
func getArrayParam(params map[string]interface{}, key string) []interface{} {
	if val, ok := params[key].([]interface{}); ok {
		return val
	}
	return []interface{}{}
}

// Helper function to return error response
func errorResponse(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

func main() {
	// Check for required action argument
	if len(os.Args) < 2 {
		result := map[string]interface{}{"error": "action required"}
		json.NewEncoder(os.Stdout).Encode(result)
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewYourPlugin()

	var result interface{}

	switch action {
	case "metadata":
		result = plugin.GetMetadata()
	case "actions":
		result = plugin.GetActions()
	default:
		// Read parameters from stdin
		var params map[string]interface{}
		inputData, err := io.ReadAll(os.Stdin)
		if err != nil {
			result = errorResponse(fmt.Errorf("failed to read input: %v", err))
		} else if len(inputData) > 0 {
			if err := json.Unmarshal(inputData, &params); err != nil {
				result = errorResponse(fmt.Errorf("failed to parse JSON input: %v", err))
			} else {
				result, err = plugin.Execute(action, params)
				if err != nil {
					result = errorResponse(err)
				}
			}
		} else {
			// No input data, execute with empty params
			result, err = plugin.Execute(action, map[string]interface{}{})
			if err != nil {
				result = errorResponse(err)
			}
		}
	}

	// Output result as JSON
	json.NewEncoder(os.Stdout).Encode(result)
}
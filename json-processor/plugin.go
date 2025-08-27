package main

import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type JsonProcessorPlugin struct{}

func (p *JsonProcessorPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "json-processor",
        Version:     "1.0.0",
        Description: "JSON parsing, manipulation, and validation",
        Author:      "Corynth Team",
        Tags:        []string{"json", "data", "parsing", "validation", "transform"},
        License:     "Apache-2.0",
    }
}

func (p *JsonProcessorPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "parse",
            Description: "Parse JSON string into structured data",
            Inputs: map[string]plugin.InputSpec{
                "json": {
                    Type:        "string",
                    Description: "JSON string to parse",
                    Required:    true,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "data": {
                    Type:        "object",
                    Description: "Parsed JSON data",
                },
                "valid": {
                    Type:        "boolean",
                    Description: "Whether the JSON was valid",
                },
            },
        },
        {
            Name:        "query",
            Description: "Query JSON data using simple path syntax",
            Inputs: map[string]plugin.InputSpec{
                "data": {
                    Type:        "object",
                    Description: "JSON data to query (or JSON string)",
                    Required:    true,
                },
                "path": {
                    Type:        "string",
                    Description: "JSON path (e.g., 'user.name', 'items[0].title')",
                    Required:    true,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "result": {
                    Type:        "object",
                    Description: "Query result",
                },
                "found": {
                    Type:        "boolean",
                    Description: "Whether the path was found",
                },
            },
        },
        {
            Name:        "validate",
            Description: "Validate JSON structure",
            Inputs: map[string]plugin.InputSpec{
                "json": {
                    Type:        "string",
                    Description: "JSON string to validate",
                    Required:    true,
                },
                "required_fields": {
                    Type:        "array",
                    Description: "List of required field paths",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "valid": {
                    Type:        "boolean",
                    Description: "Whether the JSON is valid",
                },
                "errors": {
                    Type:        "array", 
                    Description: "Validation errors",
                },
            },
        },
        {
            Name:        "transform",
            Description: "Transform JSON data",
            Inputs: map[string]plugin.InputSpec{
                "data": {
                    Type:        "object",
                    Description: "JSON data to transform",
                    Required:    true,
                },
                "mappings": {
                    Type:        "object",
                    Description: "Field mappings (old_path: new_path)",
                    Required:    true,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "result": {
                    Type:        "object",
                    Description: "Transformed data",
                },
            },
        },
    }
}

func (p *JsonProcessorPlugin) Validate(params map[string]interface{}) error {
    return nil
}

func (p *JsonProcessorPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "parse":
        return p.executeParse(ctx, params)
    case "query":
        return p.executeQuery(ctx, params)
    case "validate":
        return p.executeValidate(ctx, params)
    case "transform":
        return p.executeTransform(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *JsonProcessorPlugin) executeParse(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    jsonStr, ok := params["json"].(string)
    if !ok {
        return nil, fmt.Errorf("json parameter is required")
    }
    
    var data interface{}
    err := json.Unmarshal([]byte(jsonStr), &data)
    if err != nil {
        return map[string]interface{}{
            "data":  nil,
            "valid": false,
            "error": err.Error(),
        }, nil
    }
    
    return map[string]interface{}{
        "data":  data,
        "valid": true,
    }, nil
}

func (p *JsonProcessorPlugin) executeQuery(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    data := params["data"]
    path, ok := params["path"].(string)
    if !ok {
        return nil, fmt.Errorf("path parameter is required")
    }
    
    // If data is a string, try to parse it as JSON first
    if jsonStr, isStr := data.(string); isStr {
        var parsed interface{}
        if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
            return nil, fmt.Errorf("failed to parse JSON data: %w", err)
        }
        data = parsed
    }
    
    result, found := p.queryPath(data, path)
    
    return map[string]interface{}{
        "result": result,
        "found":  found,
        "path":   path,
    }, nil
}

func (p *JsonProcessorPlugin) executeValidate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    jsonStr, ok := params["json"].(string)
    if !ok {
        return nil, fmt.Errorf("json parameter is required")
    }
    
    var data interface{}
    err := json.Unmarshal([]byte(jsonStr), &data)
    
    errors := []string{}
    valid := err == nil
    
    if !valid {
        errors = append(errors, fmt.Sprintf("Invalid JSON: %s", err.Error()))
    }
    
    // Check required fields if specified
    if requiredFields, ok := params["required_fields"].([]interface{}); ok && valid {
        for _, field := range requiredFields {
            if fieldStr, ok := field.(string); ok {
                if _, found := p.queryPath(data, fieldStr); !found {
                    errors = append(errors, fmt.Sprintf("Required field missing: %s", fieldStr))
                    valid = false
                }
            }
        }
    }
    
    return map[string]interface{}{
        "valid":  valid,
        "errors": errors,
    }, nil
}

func (p *JsonProcessorPlugin) executeTransform(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    data := params["data"]
    mappings, ok := params["mappings"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("mappings parameter is required")
    }
    
    result := make(map[string]interface{})
    
    // Apply each mapping
    for oldPath, newPathInterface := range mappings {
        newPath, ok := newPathInterface.(string)
        if !ok {
            continue
        }
        
        value, found := p.queryPath(data, oldPath)
        if found {
            p.setPath(result, newPath, value)
        }
    }
    
    return map[string]interface{}{
        "result": result,
    }, nil
}

func (p *JsonProcessorPlugin) queryPath(data interface{}, path string) (interface{}, bool) {
    if path == "" {
        return data, true
    }
    
    parts := strings.Split(path, ".")
    current := data
    
    for _, part := range parts {
        // Handle array indexing like "items[0]"
        if strings.Contains(part, "[") && strings.Contains(part, "]") {
            arrayPart := part[:strings.Index(part, "[")]
            indexPart := part[strings.Index(part, "[")+1 : strings.Index(part, "]")]
            
            // Get the array
            if obj, ok := current.(map[string]interface{}); ok {
                if arr, exists := obj[arrayPart]; exists {
                    current = arr
                } else {
                    return nil, false
                }
            } else {
                return nil, false
            }
            
            // Get the indexed element
            if arr, ok := current.([]interface{}); ok {
                if index, err := strconv.Atoi(indexPart); err == nil && index >= 0 && index < len(arr) {
                    current = arr[index]
                } else {
                    return nil, false
                }
            } else {
                return nil, false
            }
        } else {
            // Handle regular object field access
            if obj, ok := current.(map[string]interface{}); ok {
                if value, exists := obj[part]; exists {
                    current = value
                } else {
                    return nil, false
                }
            } else {
                return nil, false
            }
        }
    }
    
    return current, true
}

func (p *JsonProcessorPlugin) setPath(data map[string]interface{}, path string, value interface{}) {
    if path == "" {
        return
    }
    
    parts := strings.Split(path, ".")
    current := data
    
    for i, part := range parts {
        if i == len(parts)-1 {
            // Set the final value
            current[part] = value
        } else {
            // Create intermediate objects if they don't exist
            if _, exists := current[part]; !exists {
                current[part] = make(map[string]interface{})
            }
            if next, ok := current[part].(map[string]interface{}); ok {
                current = next
            } else {
                // Can't traverse further
                return
            }
        }
    }
}

var ExportedPlugin plugin.Plugin = &JsonProcessorPlugin{}
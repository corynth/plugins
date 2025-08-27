package main

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type CalculatorPlugin struct{}

func (p *CalculatorPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "calculator",
        Version:     "1.0.0",
        Description: "Mathematical calculations and unit conversions",
        Author:      "Corynth Team",
        Tags:        []string{"math", "calculation", "utility", "converter"},
        License:     "Apache-2.0",
    }
}

func (p *CalculatorPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "calculate",
            Description: "Perform mathematical calculations",
            Inputs: map[string]plugin.InputSpec{
                "expression": {
                    Type:        "string",
                    Description: "Mathematical expression to evaluate",
                    Required:    true,
                },
                "precision": {
                    Type:        "number",
                    Description: "Number of decimal places for result",
                    Required:    false,
                    Default:     2,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "result": {
                    Type:        "number",
                    Description: "Calculation result",
                },
                "expression": {
                    Type:        "string", 
                    Description: "Original expression",
                },
            },
        },
        {
            Name:        "convert",
            Description: "Convert between units",
            Inputs: map[string]plugin.InputSpec{
                "value": {
                    Type:        "number",
                    Description: "Value to convert",
                    Required:    true,
                },
                "from": {
                    Type:        "string",
                    Description: "Source unit (celsius, fahrenheit, meters, feet)",
                    Required:    true,
                },
                "to": {
                    Type:        "string",
                    Description: "Target unit (celsius, fahrenheit, meters, feet)",
                    Required:    true,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "result": {
                    Type:        "number",
                    Description: "Converted value",
                },
                "unit": {
                    Type:        "string",
                    Description: "Target unit",
                },
            },
        },
    }
}

func (p *CalculatorPlugin) Validate(params map[string]interface{}) error {
    return nil
}

func (p *CalculatorPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "calculate":
        return p.executeCalculate(ctx, params)
    case "convert":
        return p.executeConvert(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *CalculatorPlugin) executeCalculate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    expression, ok := params["expression"].(string)
    if !ok {
        return nil, fmt.Errorf("expression parameter is required")
    }
    
    // Simple calculation - in a real implementation you'd use a proper math parser
    result, err := p.evaluateSimpleExpression(expression)
    if err != nil {
        return nil, fmt.Errorf("failed to evaluate expression: %w", err)
    }
    
    precision := 2.0
    if p, ok := params["precision"].(float64); ok {
        precision = p
    }
    
    // Round to specified precision
    multiplier := 1.0
    for i := 0; i < int(precision); i++ {
        multiplier *= 10
    }
    result = float64(int(result*multiplier+0.5)) / multiplier
    
    return map[string]interface{}{
        "result":     result,
        "expression": expression,
        "precision":  int(precision),
    }, nil
}

func (p *CalculatorPlugin) executeConvert(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    value, ok := params["value"].(float64)
    if !ok {
        return nil, fmt.Errorf("value parameter is required")
    }
    
    from, ok := params["from"].(string)
    if !ok {
        return nil, fmt.Errorf("from parameter is required")
    }
    
    to, ok := params["to"].(string)
    if !ok {
        return nil, fmt.Errorf("to parameter is required")
    }
    
    result, err := p.convertUnits(value, from, to)
    if err != nil {
        return nil, err
    }
    
    return map[string]interface{}{
        "result": result,
        "unit":   to,
        "from":   from,
        "value":  value,
    }, nil
}

func (p *CalculatorPlugin) evaluateSimpleExpression(expr string) (float64, error) {
    // Simple expression evaluation - handles basic operations
    expr = strings.ReplaceAll(expr, " ", "")
    
    // Handle simple cases first
    if f, err := strconv.ParseFloat(expr, 64); err == nil {
        return f, nil
    }
    
    // Handle simple addition
    if strings.Contains(expr, "+") {
        parts := strings.Split(expr, "+")
        if len(parts) == 2 {
            a, err1 := strconv.ParseFloat(parts[0], 64)
            b, err2 := strconv.ParseFloat(parts[1], 64)
            if err1 == nil && err2 == nil {
                return a + b, nil
            }
        }
    }
    
    // Handle simple subtraction
    if strings.Contains(expr, "-") && !strings.HasPrefix(expr, "-") {
        parts := strings.Split(expr, "-")
        if len(parts) == 2 {
            a, err1 := strconv.ParseFloat(parts[0], 64)
            b, err2 := strconv.ParseFloat(parts[1], 64)
            if err1 == nil && err2 == nil {
                return a - b, nil
            }
        }
    }
    
    // Handle simple multiplication
    if strings.Contains(expr, "*") {
        parts := strings.Split(expr, "*")
        if len(parts) == 2 {
            a, err1 := strconv.ParseFloat(parts[0], 64)
            b, err2 := strconv.ParseFloat(parts[1], 64)
            if err1 == nil && err2 == nil {
                return a * b, nil
            }
        }
    }
    
    // Handle simple division
    if strings.Contains(expr, "/") {
        parts := strings.Split(expr, "/")
        if len(parts) == 2 {
            a, err1 := strconv.ParseFloat(parts[0], 64)
            b, err2 := strconv.ParseFloat(parts[1], 64)
            if err1 == nil && err2 == nil {
                if b == 0 {
                    return 0, fmt.Errorf("division by zero")
                }
                return a / b, nil
            }
        }
    }
    
    return 0, fmt.Errorf("unsupported expression: %s", expr)
}

func (p *CalculatorPlugin) convertUnits(value float64, from, to string) (float64, error) {
    // Temperature conversions
    if from == "celsius" && to == "fahrenheit" {
        return (value * 9.0 / 5.0) + 32, nil
    }
    if from == "fahrenheit" && to == "celsius" {
        return (value - 32) * 5.0 / 9.0, nil
    }
    
    // Length conversions
    if from == "meters" && to == "feet" {
        return value * 3.28084, nil
    }
    if from == "feet" && to == "meters" {
        return value / 3.28084, nil
    }
    
    // Same unit
    if from == to {
        return value, nil
    }
    
    return 0, fmt.Errorf("unsupported conversion from %s to %s", from, to)
}

var ExportedPlugin plugin.Plugin = &CalculatorPlugin{}
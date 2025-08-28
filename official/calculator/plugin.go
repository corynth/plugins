package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"math"
	"os"
	"strconv"
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

type CalculatorPlugin struct{}

func NewCalculatorPlugin() *CalculatorPlugin {
	return &CalculatorPlugin{}
}

func (p *CalculatorPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "calculator",
		Version:     "1.0.0",
		Description: "Mathematical calculations with safe expression evaluation",
		Author:      "Corynth Team",
		Tags:        []string{"math", "calculation", "utility", "ast"},
	}
}

func (p *CalculatorPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"calculate": {
			Description: "Perform safe mathematical calculations using AST parsing",
			Inputs: map[string]IOSpec{
				"expression": {
					Type:        "string",
					Required:    true,
					Description: "Mathematical expression to evaluate (supports +, -, *, /, %, parentheses)",
				},
				"precision": {
					Type:        "number",
					Required:    false,
					Default:     2,
					Description: "Decimal precision for results",
				},
			},
			Outputs: map[string]IOSpec{
				"result":     {Type: "number", Description: "Calculation result"},
				"expression": {Type: "string", Description: "Original expression"},
			},
		},
	}
}

func (p *CalculatorPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "calculate":
		return p.calculate(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *CalculatorPlugin) calculate(params map[string]interface{}) (map[string]interface{}, error) {
	expression, ok := params["expression"].(string)
	if !ok || expression == "" {
		return map[string]interface{}{"error": "expression parameter is required"}, nil
	}

	precision := 2
	if prec, ok := params["precision"].(float64); ok {
		precision = int(prec)
		if precision < 0 {
			precision = 0
		}
	}

	// Parse and evaluate the expression using AST
	result, err := p.evaluateExpression(expression)
	if err != nil {
		return map[string]interface{}{
			"error":      fmt.Sprintf("Invalid expression: %v", err),
			"expression": expression,
		}, nil
	}

	// Apply precision
	if precision > 0 {
		multiplier := math.Pow(10, float64(precision))
		result = math.Round(result*multiplier) / multiplier
	} else {
		result = math.Round(result)
	}

	return map[string]interface{}{
		"result":     result,
		"expression": expression,
	}, nil
}

func (p *CalculatorPlugin) evaluateExpression(expr string) (float64, error) {
	// Parse the expression into an AST
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return 0, fmt.Errorf("syntax error: %v", err)
	}

	// Evaluate the AST
	return p.evalNode(node)
}

func (p *CalculatorPlugin) evalNode(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		return p.evalBinaryExpr(n)
	case *ast.UnaryExpr:
		return p.evalUnaryExpr(n)
	case *ast.ParenExpr:
		return p.evalNode(n.X)
	case *ast.BasicLit:
		return p.evalBasicLit(n)
	case *ast.Ident:
		// Only allow math constants
		switch n.Name {
		case "pi":
			return math.Pi, nil
		case "e":
			return math.E, nil
		default:
			return 0, fmt.Errorf("undefined identifier: %s", n.Name)
		}
	default:
		return 0, fmt.Errorf("unsupported expression type: %T", node)
	}
}

func (p *CalculatorPlugin) evalBinaryExpr(expr *ast.BinaryExpr) (float64, error) {
	left, err := p.evalNode(expr.X)
	if err != nil {
		return 0, err
	}

	right, err := p.evalNode(expr.Y)
	if err != nil {
		return 0, err
	}

	switch expr.Op {
	case token.ADD:
		return left + right, nil
	case token.SUB:
		return left - right, nil
	case token.MUL:
		return left * right, nil
	case token.QUO:
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	case token.REM:
		if right == 0 {
			return 0, fmt.Errorf("modulo by zero")
		}
		return math.Mod(left, right), nil
	default:
		return 0, fmt.Errorf("unsupported binary operator: %v", expr.Op)
	}
}

func (p *CalculatorPlugin) evalUnaryExpr(expr *ast.UnaryExpr) (float64, error) {
	operand, err := p.evalNode(expr.X)
	if err != nil {
		return 0, err
	}

	switch expr.Op {
	case token.ADD:
		return +operand, nil
	case token.SUB:
		return -operand, nil
	default:
		return 0, fmt.Errorf("unsupported unary operator: %v", expr.Op)
	}
}

func (p *CalculatorPlugin) evalBasicLit(lit *ast.BasicLit) (float64, error) {
	switch lit.Kind {
	case token.INT:
		val, err := strconv.ParseInt(lit.Value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid integer: %v", err)
		}
		return float64(val), nil
	case token.FLOAT:
		val, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid float: %v", err)
		}
		return val, nil
	default:
		return 0, fmt.Errorf("unsupported literal type: %v", lit.Kind)
	}
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewCalculatorPlugin()

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
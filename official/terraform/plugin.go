package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type TerraformPlugin struct {
	WorkingDir string
}

type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}

type ActionSpec struct {
	Description string                 `json:"description"`
	Inputs      map[string]interface{} `json:"inputs"`
	Outputs     map[string]interface{} `json:"outputs"`
}

func (p *TerraformPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "terraform",
		Version:     "1.0.0",
		Description: "Terraform Infrastructure as Code operations",
		Author:      "Corynth Team",
		Tags:        []string{"terraform", "iac", "infrastructure", "cloud", "provisioning"},
	}
}

func (p *TerraformPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"init": {
			Description: "Initialize Terraform working directory",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
				"upgrade": map[string]interface{}{
					"type":        "boolean",
					"required":    false,
					"default":     false,
					"description": "Upgrade modules and plugins",
				},
			},
			Outputs: map[string]interface{}{
				"success": map[string]interface{}{"type": "boolean"},
				"output":  map[string]interface{}{"type": "string"},
			},
		},
		"plan": {
			Description: "Create Terraform execution plan with change analysis",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
				"var_file": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Variables file path",
				},
				"vars": map[string]interface{}{
					"type":        "object",
					"required":    false,
					"description": "Variable key-value pairs",
				},
				"out": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Plan output file",
				},
				"destroy": map[string]interface{}{
					"type":        "boolean",
					"required":    false,
					"default":     false,
					"description": "Create destroy plan",
				},
			},
			Outputs: map[string]interface{}{
				"success":    map[string]interface{}{"type": "boolean"},
				"output":     map[string]interface{}{"type": "string"},
				"plan_file":  map[string]interface{}{"type": "string"},
				"changes":    map[string]interface{}{"type": "number"},
				"adds":       map[string]interface{}{"type": "number"},
				"changes_op": map[string]interface{}{"type": "number"},
				"destroys":   map[string]interface{}{"type": "number"},
			},
		},
		"apply": {
			Description: "Apply infrastructure changes",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
				"plan_file": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Plan file to apply",
				},
				"var_file": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Variables file path",
				},
				"vars": map[string]interface{}{
					"type":        "object",
					"required":    false,
					"description": "Variable key-value pairs",
				},
				"auto_approve": map[string]interface{}{
					"type":        "boolean",
					"required":    false,
					"default":     false,
					"description": "Skip interactive approval",
				},
			},
			Outputs: map[string]interface{}{
				"success": map[string]interface{}{"type": "boolean"},
				"output":  map[string]interface{}{"type": "string"},
				"outputs": map[string]interface{}{"type": "object"},
			},
		},
		"destroy": {
			Description: "Destroy managed infrastructure",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
				"var_file": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Variables file path",
				},
				"vars": map[string]interface{}{
					"type":        "object",
					"required":    false,
					"description": "Variable key-value pairs",
				},
				"auto_approve": map[string]interface{}{
					"type":        "boolean",
					"required":    false,
					"default":     false,
					"description": "Skip interactive approval",
				},
			},
			Outputs: map[string]interface{}{
				"success": map[string]interface{}{"type": "boolean"},
				"output":  map[string]interface{}{"type": "string"},
			},
		},
		"validate": {
			Description: "Validate configuration syntax",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
			},
			Outputs: map[string]interface{}{
				"success": map[string]interface{}{"type": "boolean"},
				"output":  map[string]interface{}{"type": "string"},
				"valid":   map[string]interface{}{"type": "boolean"},
				"errors":  map[string]interface{}{"type": "array"},
			},
		},
		"output": {
			Description: "Extract output values",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Specific output name",
				},
			},
			Outputs: map[string]interface{}{
				"success": map[string]interface{}{"type": "boolean"},
				"outputs": map[string]interface{}{"type": "object"},
			},
		},
		"workspace": {
			Description: "Manage workspaces (list, new, select, delete)",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
				"operation": map[string]interface{}{
					"type":        "string",
					"required":    true,
					"description": "Operation: list, new, select, delete",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Workspace name",
				},
			},
			Outputs: map[string]interface{}{
				"success":    map[string]interface{}{"type": "boolean"},
				"workspaces": map[string]interface{}{"type": "array"},
				"current":    map[string]interface{}{"type": "string"},
			},
		},
		"import": {
			Description: "Import existing resources",
			Inputs: map[string]interface{}{
				"working_dir": map[string]interface{}{
					"type":        "string",
					"required":    false,
					"description": "Working directory path",
				},
				"address": map[string]interface{}{
					"type":        "string",
					"required":    true,
					"description": "Resource address in Terraform",
				},
				"id": map[string]interface{}{
					"type":        "string",
					"required":    true,
					"description": "Resource ID in provider",
				},
			},
			Outputs: map[string]interface{}{
				"success": map[string]interface{}{"type": "boolean"},
				"output":  map[string]interface{}{"type": "string"},
			},
		},
	}
}

func (p *TerraformPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	// Set working directory
	if wd, ok := params["working_dir"].(string); ok && wd != "" {
		p.WorkingDir = wd
	} else {
		p.WorkingDir, _ = os.Getwd()
	}

	switch action {
	case "init":
		return p.terraformInit(params)
	case "plan":
		return p.terraformPlan(params)
	case "apply":
		return p.terraformApply(params)
	case "destroy":
		return p.terraformDestroy(params)
	case "validate":
		return p.terraformValidate(params)
	case "output":
		return p.terraformOutput(params)
	case "workspace":
		return p.terraformWorkspace(params)
	case "import":
		return p.terraformImport(params)
	default:
		return map[string]interface{}{
			"error": fmt.Sprintf("Unknown action: %s", action),
		}, nil
	}
}

func (p *TerraformPlugin) runTerraformCommand(args []string, input string) (string, int, error) {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = p.WorkingDir

	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	output, err := cmd.CombinedOutput()
	exitCode := 0

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		} else {
			return "", -1, err
		}
	}

	return string(output), exitCode, nil
}

func (p *TerraformPlugin) terraformInit(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"init", "-no-color"}

	if upgrade, ok := params["upgrade"].(bool); ok && upgrade {
		args = append(args, "-upgrade")
	}

	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	return map[string]interface{}{
		"success": exitCode == 0,
		"output":  output,
	}, nil
}

func (p *TerraformPlugin) terraformPlan(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"plan", "-no-color", "-detailed-exitcode"}

	if varFile, ok := params["var_file"].(string); ok && varFile != "" {
		args = append(args, "-var-file", varFile)
	}

	if vars, ok := params["vars"].(map[string]interface{}); ok {
		for key, value := range vars {
			args = append(args, "-var", fmt.Sprintf("%s=%v", key, value))
		}
	}

	if outFile, ok := params["out"].(string); ok && outFile != "" {
		args = append(args, "-out", outFile)
	}

	if destroy, ok := params["destroy"].(bool); ok && destroy {
		args = append(args, "-destroy")
	}

	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	result := map[string]interface{}{
		"success": exitCode == 0 || exitCode == 2, // 2 means changes present
		"output":  output,
	}

	if outFile, ok := params["out"].(string); ok && outFile != "" {
		result["plan_file"] = filepath.Join(p.WorkingDir, outFile)
	}

	// Parse plan output for changes count
	changes, adds, changesOp, destroys := p.parsePlanOutput(output)
	result["changes"] = changes
	result["adds"] = adds
	result["changes_op"] = changesOp
	result["destroys"] = destroys

	return result, nil
}

func (p *TerraformPlugin) terraformApply(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"apply", "-no-color"}

	if autoApprove, ok := params["auto_approve"].(bool); ok && autoApprove {
		args = append(args, "-auto-approve")
	}

	if planFile, ok := params["plan_file"].(string); ok && planFile != "" {
		args = append(args, planFile)
	} else {
		if varFile, ok := params["var_file"].(string); ok && varFile != "" {
			args = append(args, "-var-file", varFile)
		}

		if vars, ok := params["vars"].(map[string]interface{}); ok {
			for key, value := range vars {
				args = append(args, "-var", fmt.Sprintf("%s=%v", key, value))
			}
		}
	}

	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	result := map[string]interface{}{
		"success": exitCode == 0,
		"output":  output,
	}

	// Get outputs after successful apply
	if exitCode == 0 {
		if outputs, err := p.getTerraformOutputs(); err == nil {
			result["outputs"] = outputs
		}
	}

	return result, nil
}

func (p *TerraformPlugin) terraformDestroy(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"destroy", "-no-color"}

	if autoApprove, ok := params["auto_approve"].(bool); ok && autoApprove {
		args = append(args, "-auto-approve")
	}

	if varFile, ok := params["var_file"].(string); ok && varFile != "" {
		args = append(args, "-var-file", varFile)
	}

	if vars, ok := params["vars"].(map[string]interface{}); ok {
		for key, value := range vars {
			args = append(args, "-var", fmt.Sprintf("%s=%v", key, value))
		}
	}

	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	return map[string]interface{}{
		"success": exitCode == 0,
		"output":  output,
	}, nil
}

func (p *TerraformPlugin) terraformValidate(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"validate", "-json"}

	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	result := map[string]interface{}{
		"success": exitCode == 0,
		"output":  output,
	}

	// Parse JSON validation output
	var validation map[string]interface{}
	if err := json.Unmarshal([]byte(output), &validation); err == nil {
		if valid, ok := validation["valid"].(bool); ok {
			result["valid"] = valid
		}
		if errorCount, ok := validation["error_count"].(float64); ok && errorCount > 0 {
			if diagnostics, ok := validation["diagnostics"].([]interface{}); ok {
				result["errors"] = diagnostics
			}
		}
	}

	return result, nil
}

func (p *TerraformPlugin) terraformOutput(params map[string]interface{}) (map[string]interface{}, error) {
	outputs, err := p.getTerraformOutputs()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	result := map[string]interface{}{
		"success": true,
		"outputs": outputs,
	}

	// If specific output name requested, return just that value
	if name, ok := params["name"].(string); ok && name != "" {
		if value, exists := outputs[name]; exists {
			result["outputs"] = map[string]interface{}{name: value}
		}
	}

	return result, nil
}

func (p *TerraformPlugin) terraformWorkspace(params map[string]interface{}) (map[string]interface{}, error) {
	operation, ok := params["operation"].(string)
	if !ok {
		return map[string]interface{}{"error": "operation parameter is required"}, nil
	}

	var args []string
	switch operation {
	case "list":
		args = []string{"workspace", "list"}
	case "new":
		name, ok := params["name"].(string)
		if !ok {
			return map[string]interface{}{"error": "name parameter required for new workspace"}, nil
		}
		args = []string{"workspace", "new", name}
	case "select":
		name, ok := params["name"].(string)
		if !ok {
			return map[string]interface{}{"error": "name parameter required for select workspace"}, nil
		}
		args = []string{"workspace", "select", name}
	case "delete":
		name, ok := params["name"].(string)
		if !ok {
			return map[string]interface{}{"error": "name parameter required for delete workspace"}, nil
		}
		args = []string{"workspace", "delete", name}
	default:
		return map[string]interface{}{"error": "invalid operation: " + operation}, nil
	}

	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	result := map[string]interface{}{
		"success": exitCode == 0,
	}

	// Parse workspace list output
	if operation == "list" && exitCode == 0 {
		workspaces := []string{}
		current := ""
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				if strings.HasPrefix(line, "* ") {
					current = strings.TrimPrefix(line, "* ")
					workspaces = append(workspaces, current)
				} else {
					workspaces = append(workspaces, line)
				}
			}
		}
		result["workspaces"] = workspaces
		result["current"] = current
	}

	return result, nil
}

func (p *TerraformPlugin) terraformImport(params map[string]interface{}) (map[string]interface{}, error) {
	address, ok := params["address"].(string)
	if !ok {
		return map[string]interface{}{"error": "address parameter is required"}, nil
	}

	id, ok := params["id"].(string)
	if !ok {
		return map[string]interface{}{"error": "id parameter is required"}, nil
	}

	args := []string{"import", "-no-color", address, id}

	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	return map[string]interface{}{
		"success": exitCode == 0,
		"output":  output,
	}, nil
}

func (p *TerraformPlugin) getTerraformOutputs() (map[string]interface{}, error) {
	args := []string{"output", "-json"}
	output, exitCode, err := p.runTerraformCommand(args, "")
	if err != nil || exitCode != 0 {
		return make(map[string]interface{}), nil // Return empty map if no outputs
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal([]byte(output), &outputs); err != nil {
		return make(map[string]interface{}), nil
	}

	// Extract just the values from Terraform's output format
	result := make(map[string]interface{})
	for key, value := range outputs {
		if valueMap, ok := value.(map[string]interface{}); ok {
			if val, exists := valueMap["value"]; exists {
				result[key] = val
			}
		}
	}

	return result, nil
}

func (p *TerraformPlugin) parsePlanOutput(output string) (int, int, int, int) {
	changes, adds, changesOp, destroys := 0, 0, 0, 0

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Plan:") {
			// Parse plan summary line
			if strings.Contains(line, "to add") {
				fmt.Sscanf(line, "Plan: %d to add, %d to change, %d to destroy.", &adds, &changesOp, &destroys)
				changes = adds + changesOp + destroys
			}
			break
		}
	}

	return changes, adds, changesOp, destroys
}

func main() {
	if len(os.Args) < 2 {
		result := map[string]interface{}{"error": "action required"}
		json.NewEncoder(os.Stdout).Encode(result)
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := &TerraformPlugin{}

	var params map[string]interface{}
	if input, err := io.ReadAll(os.Stdin); err == nil && len(input) > 0 {
		json.Unmarshal(input, &params)
	}
	if params == nil {
		params = make(map[string]interface{})
	}

	var result map[string]interface{}
	switch action {
	case "metadata":
		result = map[string]interface{}{
			"name":        plugin.GetMetadata().Name,
			"version":     plugin.GetMetadata().Version,
			"description": plugin.GetMetadata().Description,
			"author":      plugin.GetMetadata().Author,
			"tags":        plugin.GetMetadata().Tags,
		}
	case "actions":
		result = make(map[string]interface{})
		for name, spec := range plugin.GetActions() {
			result[name] = spec
		}
	default:
		var err error
		result, err = plugin.Execute(action, params)
		if err != nil {
			result = map[string]interface{}{"error": err.Error()}
		}
	}

	json.NewEncoder(os.Stdout).Encode(result)
}
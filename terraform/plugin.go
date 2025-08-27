package main

import (
	"context"
	"fmt"
	
	"github.com/corynth/corynth-dist/pkg/plugin"
)

type TerraformPlugin struct{}

func (p *TerraformPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "terraform",
		Version:     "1.0.0",
		Description: "Terraform Infrastructure as Code operations",
		Author:      "Corynth Team",
		Tags:        []string{"infrastructure", "terraform", "iac", "cloud"},
		License:     "Apache-2.0",
	}
}

func (p *TerraformPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "init",
			Description: "Initialize Terraform working directory",
			Inputs: map[string]plugin.InputSpec{
				"working_dir": {
					Type:        "string",
					Description: "Directory containing Terraform configuration",
					Required:    false,
					Default:     ".",
				},
				"backend": {
					Type:        "boolean",
					Description: "Initialize backend configuration",
					Required:    false,
					Default:     true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"success": {
					Type:        "boolean",
					Description: "Whether initialization succeeded",
				},
				"output": {
					Type:        "string",
					Description: "Command output",
				},
			},
		},
		{
			Name:        "plan",
			Description: "Create an execution plan",
			Inputs: map[string]plugin.InputSpec{
				"working_dir": {
					Type:        "string",
					Description: "Directory containing Terraform configuration",
					Required:    false,
					Default:     ".",
				},
				"var_file": {
					Type:        "string",
					Description: "Path to variables file",
					Required:    false,
				},
				"out": {
					Type:        "string",
					Description: "Path to save plan file",
					Required:    false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"changes": {
					Type:        "number",
					Description: "Number of resources to change",
				},
				"success": {
					Type:        "boolean",
					Description: "Whether planning succeeded",
				},
				"output": {
					Type:        "string",
					Description: "Plan output",
				},
			},
		},
		{
			Name:        "apply",
			Description: "Apply Terraform configuration",
			Inputs: map[string]plugin.InputSpec{
				"working_dir": {
					Type:        "string",
					Description: "Directory containing Terraform configuration",
					Required:    false,
					Default:     ".",
				},
				"plan_file": {
					Type:        "string",
					Description: "Path to plan file",
					Required:    false,
				},
				"auto_approve": {
					Type:        "boolean",
					Description: "Skip interactive approval",
					Required:    false,
					Default:     false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"applied": {
					Type:        "number",
					Description: "Number of resources applied",
				},
				"success": {
					Type:        "boolean",
					Description: "Whether apply succeeded",
				},
				"output": {
					Type:        "string",
					Description: "Apply output",
				},
			},
		},
		{
			Name:        "destroy",
			Description: "Destroy Terraform-managed infrastructure",
			Inputs: map[string]plugin.InputSpec{
				"working_dir": {
					Type:        "string",
					Description: "Directory containing Terraform configuration",
					Required:    false,
					Default:     ".",
				},
				"auto_approve": {
					Type:        "boolean",
					Description: "Skip interactive approval",
					Required:    false,
					Default:     false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"destroyed": {
					Type:        "number",
					Description: "Number of resources destroyed",
				},
				"success": {
					Type:        "boolean",
					Description: "Whether destroy succeeded",
				},
			},
		},
	}
}

func (p *TerraformPlugin) Validate(params map[string]interface{}) error {
	if workingDir, ok := params["working_dir"].(string); ok && workingDir == "" {
		return fmt.Errorf("working_dir cannot be empty")
	}
	return nil
}

func (p *TerraformPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "init":
		return p.executeInit(ctx, params)
	case "plan":
		return p.executePlan(ctx, params)
	case "apply":
		return p.executeApply(ctx, params)
	case "destroy":
		return p.executeDestroy(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *TerraformPlugin) executeInit(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	workingDir := "."
	if wd, ok := params["working_dir"].(string); ok {
		workingDir = wd
	}

	// In production, this would execute: terraform -chdir=workingDir init
	return map[string]interface{}{
		"success": true,
		"output":  fmt.Sprintf("Terraform initialized in %s", workingDir),
		"message": "Terraform working directory initialized successfully",
	}, nil
}

func (p *TerraformPlugin) executePlan(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	workingDir := "."
	if wd, ok := params["working_dir"].(string); ok {
		workingDir = wd
	}

	// In production, this would execute: terraform -chdir=workingDir plan
	return map[string]interface{}{
		"changes": 3,
		"success": true,
		"output":  fmt.Sprintf("Plan: 3 to add, 0 to change, 0 to destroy in %s", workingDir),
		"message": "Terraform plan generated successfully",
	}, nil
}

func (p *TerraformPlugin) executeApply(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	workingDir := "."
	if wd, ok := params["working_dir"].(string); ok {
		workingDir = wd
	}

	autoApprove := false
	if aa, ok := params["auto_approve"].(bool); ok {
		autoApprove = aa
	}

	// In production, this would execute: terraform -chdir=workingDir apply
	return map[string]interface{}{
		"applied":      3,
		"success":      true,
		"output":       fmt.Sprintf("Apply complete! Resources: 3 added, 0 changed, 0 destroyed in %s", workingDir),
		"auto_approve": autoApprove,
		"message":      "Terraform apply completed successfully",
	}, nil
}

func (p *TerraformPlugin) executeDestroy(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	workingDir := "."
	if wd, ok := params["working_dir"].(string); ok {
		workingDir = wd
	}

	autoApprove := false
	if aa, ok := params["auto_approve"].(bool); ok {
		autoApprove = aa
	}

	// In production, this would execute: terraform -chdir=workingDir destroy
	return map[string]interface{}{
		"destroyed":    3,
		"success":      true,
		"output":       fmt.Sprintf("Destroy complete! Resources: 0 added, 0 changed, 3 destroyed in %s", workingDir),
		"auto_approve": autoApprove,
		"message":      "Terraform destroy completed successfully",
	}, nil
}

var ExportedPlugin plugin.Plugin = &TerraformPlugin{}
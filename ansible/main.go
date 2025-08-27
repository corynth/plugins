package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corynth/corynth/pkg/plugin"
	pluginv2 "github.com/corynth/corynth/pkg/plugin/v2"
)

type AnsiblePlugin struct{}

func (p *AnsiblePlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "ansible",
		Version:     "1.0.0",
		Description: "Ansible automation plugin for configuration management, deployments, and infrastructure orchestration",
		Author:      "Corynth Team",
		Tags:        []string{"automation", "configuration", "deployment", "infrastructure", "orchestration"},
		License:     "MIT",
	}
}

func (p *AnsiblePlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "playbook",
			Description: "Execute an Ansible playbook",
			Inputs: map[string]plugin.InputSpec{
				"playbook_path": {
					Type:        "string",
					Description: "Path to playbook YAML file",
					Required:    true,
				},
				"inventory": {
					Type:        "string",
					Description: "Inventory file path or host list",
					Required:    false,
					Default:     "localhost,",
				},
				"extra_vars": {
					Type:        "object",
					Description: "Extra variables to pass to playbook",
					Required:    false,
				},
				"tags": {
					Type:        "string",
					Description: "Comma-separated list of tags to run",
					Required:    false,
				},
				"limit": {
					Type:        "string",
					Description: "Limit execution to specific hosts",
					Required:    false,
				},
				"check_mode": {
					Type:        "boolean",
					Description: "Run in check mode (dry-run)",
					Required:    false,
					Default:     false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"status": {
					Type:        "string",
					Description: "Execution status",
				},
				"output": {
					Type:        "string",
					Description: "Ansible output",
				},
				"changed": {
					Type:        "boolean",
					Description: "Whether any changes were made",
				},
				"failed": {
					Type:        "boolean",
					Description: "Whether execution failed",
				},
			},
		},
		{
			Name:        "adhoc",
			Description: "Run an ad-hoc Ansible command",
			Inputs: map[string]plugin.InputSpec{
				"hosts": {
					Type:        "string",
					Description: "Target hosts pattern",
					Required:    true,
				},
				"module": {
					Type:        "string",
					Description: "Ansible module to execute",
					Required:    true,
				},
				"args": {
					Type:        "string",
					Description: "Module arguments",
					Required:    false,
				},
				"inventory": {
					Type:        "string",
					Description: "Inventory file or host list",
					Required:    false,
					Default:     "localhost,",
				},
				"become": {
					Type:        "boolean",
					Description: "Use privilege escalation",
					Required:    false,
					Default:     false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"status": {
					Type:        "string",
					Description: "Command execution status",
				},
				"output": {
					Type:        "string",
					Description: "Command output",
				},
			},
		},
		{
			Name:        "inventory_list",
			Description: "List hosts from inventory",
			Inputs: map[string]plugin.InputSpec{
				"inventory": {
					Type:        "string",
					Description: "Inventory file path",
					Required:    true,
				},
				"group": {
					Type:        "string",
					Description: "Specific group to list",
					Required:    false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"hosts": {
					Type:        "array",
					Description: "List of hosts",
				},
			},
		},
		{
			Name:        "vault_decrypt",
			Description: "Decrypt Ansible vault file",
			Inputs: map[string]plugin.InputSpec{
				"vault_file": {
					Type:        "string",
					Description: "Path to encrypted vault file",
					Required:    true,
				},
				"vault_password": {
					Type:        "string",
					Description: "Vault password",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"content": {
					Type:        "string",
					Description: "Decrypted content",
				},
			},
		},
	}
}

func (p *AnsiblePlugin) Validate(params map[string]interface{}) error {
	// Check if ansible is installed
	if _, err := exec.LookPath("ansible"); err != nil {
		return fmt.Errorf("ansible not found in PATH. Please install Ansible")
	}
	
	if _, err := exec.LookPath("ansible-playbook"); err != nil {
		return fmt.Errorf("ansible-playbook not found in PATH. Please install Ansible")
	}
	
	return nil
}

func (p *AnsiblePlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "playbook":
		return p.runPlaybook(ctx, params)
	case "adhoc":
		return p.runAdHoc(ctx, params)
	case "inventory_list":
		return p.listInventory(ctx, params)
	case "vault_decrypt":
		return p.decryptVault(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *AnsiblePlugin) runPlaybook(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	playbookPath := params["playbook_path"].(string)
	
	// Check if playbook exists
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("playbook not found: %s", playbookPath)
	}

	// Build ansible-playbook command
	args := []string{"ansible-playbook"}
	
	// Add inventory
	inventory := "localhost,"
	if inv, ok := params["inventory"].(string); ok && inv != "" {
		inventory = inv
	}
	args = append(args, "-i", inventory)
	
	// Add extra vars
	if extraVars, ok := params["extra_vars"].(map[string]interface{}); ok {
		varsJSON, _ := json.Marshal(extraVars)
		args = append(args, "--extra-vars", string(varsJSON))
	}
	
	// Add tags
	if tags, ok := params["tags"].(string); ok && tags != "" {
		args = append(args, "--tags", tags)
	}
	
	// Add limit
	if limit, ok := params["limit"].(string); ok && limit != "" {
		args = append(args, "--limit", limit)
	}
	
	// Add check mode
	if checkMode, ok := params["check_mode"].(bool); ok && checkMode {
		args = append(args, "--check")
	}
	
	// Add playbook path
	args = append(args, playbookPath)
	
	// Execute command
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}
	
	// Parse output for status
	outputStr := string(output)
	changed := strings.Contains(outputStr, "changed=") && !strings.Contains(outputStr, "changed=0")
	failed := exitCode != 0 || strings.Contains(outputStr, "failed=") && !strings.Contains(outputStr, "failed=0")
	
	status := "success"
	if failed {
		status = "failed"
	} else if changed {
		status = "changed"
	}

	return map[string]interface{}{
		"status":  status,
		"output":  outputStr,
		"changed": changed,
		"failed":  failed,
		"exit_code": exitCode,
	}, nil
}

func (p *AnsiblePlugin) runAdHoc(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	hosts := params["hosts"].(string)
	module := params["module"].(string)
	
	// Build ansible command
	args := []string{"ansible", hosts}
	
	// Add inventory
	inventory := "localhost,"
	if inv, ok := params["inventory"].(string); ok && inv != "" {
		inventory = inv
	}
	args = append(args, "-i", inventory)
	
	// Add module
	args = append(args, "-m", module)
	
	// Add module arguments
	if moduleArgs, ok := params["args"].(string); ok && moduleArgs != "" {
		args = append(args, "-a", moduleArgs)
	}
	
	// Add become
	if become, ok := params["become"].(bool); ok && become {
		args = append(args, "--become")
	}
	
	// Execute command
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	
	status := "success"
	if err != nil {
		status = "failed"
	}

	return map[string]interface{}{
		"status": status,
		"output": string(output),
	}, nil
}

func (p *AnsiblePlugin) listInventory(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	inventory := params["inventory"].(string)
	
	args := []string{"ansible-inventory", "-i", inventory, "--list"}
	
	if group, ok := params["group"].(string); ok && group != "" {
		args = append(args, "--host", group)
	}
	
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	
	var inventoryData map[string]interface{}
	if err := json.Unmarshal(output, &inventoryData); err != nil {
		return nil, fmt.Errorf("failed to parse inventory JSON: %w", err)
	}

	// Extract hosts list
	var hosts []string
	if all, ok := inventoryData["_meta"].(map[string]interface{}); ok {
		if hostvars, ok := all["hostvars"].(map[string]interface{}); ok {
			for host := range hostvars {
				hosts = append(hosts, host)
			}
		}
	}

	return map[string]interface{}{
		"hosts":     hosts,
		"inventory": inventoryData,
	}, nil
}

func (p *AnsiblePlugin) decryptVault(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	vaultFile := params["vault_file"].(string)
	vaultPassword := params["vault_password"].(string)
	
	// Check if vault file exists
	if _, err := os.Stat(vaultFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault file not found: %s", vaultFile)
	}
	
	// Create temporary password file
	tempDir := os.TempDir()
	passwordFile := filepath.Join(tempDir, "ansible_vault_pass")
	defer os.Remove(passwordFile)
	
	if err := os.WriteFile(passwordFile, []byte(vaultPassword), 0600); err != nil {
		return nil, fmt.Errorf("failed to create password file: %w", err)
	}
	
	// Run ansible-vault view
	cmd := exec.CommandContext(ctx, "ansible-vault", "view", "--vault-password-file", passwordFile, vaultFile)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault: %w", err)
	}

	return map[string]interface{}{
		"status":  "success",
		"content": string(output),
	}, nil
}

var ExportedPlugin plugin.Plugin = &AnsiblePlugin{}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		if err := pluginv2.ServePlugin(&AnsiblePlugin{}); err != nil {
			fmt.Printf("Error serving plugin: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Corynth Ansible Plugin v1.0.0\n")
		fmt.Printf("Usage: %s serve\n", os.Args[0])
	}
}
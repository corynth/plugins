package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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

type AnsiblePlugin struct{}

func NewAnsiblePlugin() *AnsiblePlugin {
	return &AnsiblePlugin{}
}

func (p *AnsiblePlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "ansible",
		Version:     "1.0.0",
		Description: "Ansible configuration management and automation",
		Author:      "Corynth Team",
		Tags:        []string{"ansible", "configuration", "automation", "playbook"},
	}
}

func (p *AnsiblePlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"playbook": {
			Description: "Run Ansible playbook",
			Inputs: map[string]IOSpec{
				"playbook":  {Type: "string", Required: true, Description: "Playbook YAML content or file path"},
				"inventory": {Type: "string", Required: false, Description: "Inventory content or file path"},
				"vars":      {Type: "object", Required: false, Description: "Extra variables"},
				"limit":     {Type: "string", Required: false, Description: "Limit to specific hosts"},
				"tags":      {Type: "string", Required: false, Description: "Run specific tags"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Operation success"},
				"output":  {Type: "string", Description: "Command output"},
				"stats":   {Type: "object", Description: "Ansible execution statistics"},
			},
		},
		"ad_hoc": {
			Description: "Run ad-hoc command",
			Inputs: map[string]IOSpec{
				"hosts":     {Type: "string", Required: true, Description: "Target hosts"},
				"module":    {Type: "string", Required: true, Description: "Ansible module"},
				"args":      {Type: "string", Required: false, Description: "Module arguments"},
				"inventory": {Type: "string", Required: false, Description: "Inventory file"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Operation success"},
				"output":  {Type: "string", Description: "Command output"},
			},
		},
	}
}

func (p *AnsiblePlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "playbook":
		return p.runPlaybook(params)
	case "ad_hoc":
		return p.runAdHoc(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *AnsiblePlugin) runPlaybook(params map[string]interface{}) (map[string]interface{}, error) {
	playbook, ok := params["playbook"].(string)
	if !ok || playbook == "" {
		return map[string]interface{}{"error": "playbook is required"}, nil
	}

	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "ansible-")
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to create temp dir: %v", err)}, nil
	}
	defer os.RemoveAll(tmpDir)

	// Handle playbook file
	var playbookFile string
	if _, err := os.Stat(playbook); err == nil {
		// It's an existing file path
		playbookFile = playbook
	} else {
		// It's YAML content, write to temp file
		playbookFile = filepath.Join(tmpDir, "playbook.yml")
		if err := ioutil.WriteFile(playbookFile, []byte(playbook), 0644); err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to write playbook: %v", err)}, nil
		}
	}

	// Handle inventory
	inventoryFile := "localhost,"
	if inventory, ok := params["inventory"].(string); ok && inventory != "" {
		if _, err := os.Stat(inventory); err == nil {
			// It's an existing file path
			inventoryFile = inventory
		} else {
			// It's inventory content, write to temp file
			invFile := filepath.Join(tmpDir, "inventory")
			if err := ioutil.WriteFile(invFile, []byte(inventory), 0644); err != nil {
				return map[string]interface{}{"error": fmt.Sprintf("failed to write inventory: %v", err)}, nil
			}
			inventoryFile = invFile
		}
	}

	// Build ansible-playbook command
	args := []string{"ansible-playbook", playbookFile, "-i", inventoryFile}

	// Add extra vars
	if vars, ok := params["vars"].(map[string]interface{}); ok && len(vars) > 0 {
		varsJSON, err := json.Marshal(vars)
		if err == nil {
			args = append(args, "--extra-vars", string(varsJSON))
		}
	}

	// Add limit
	if limit, ok := params["limit"].(string); ok && limit != "" {
		args = append(args, "--limit", limit)
	}

	// Add tags
	if tags, ok := params["tags"].(string); ok && tags != "" {
		args = append(args, "--tags", tags)
	}

	// Execute command
	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	output, err := cmd.CombinedOutput()
	
	success := err == nil
	outputStr := string(output)
	stats := p.parseAnsibleStats(outputStr)

	return map[string]interface{}{
		"success": success,
		"output":  outputStr,
		"stats":   stats,
	}, nil
}

func (p *AnsiblePlugin) runAdHoc(params map[string]interface{}) (map[string]interface{}, error) {
	hosts, ok := params["hosts"].(string)
	if !ok || hosts == "" {
		return map[string]interface{}{"error": "hosts is required"}, nil
	}

	module, ok := params["module"].(string)
	if !ok || module == "" {
		return map[string]interface{}{"error": "module is required"}, nil
	}

	// Build ansible command
	args := []string{"ansible", hosts}

	// Add inventory
	inventory := "localhost,"
	if inv, ok := params["inventory"].(string); ok && inv != "" {
		inventory = inv
	}
	args = append(args, "-i", inventory, "-m", module)

	// Add module arguments
	if moduleArgs, ok := params["args"].(string); ok && moduleArgs != "" {
		args = append(args, "-a", moduleArgs)
	}

	// Execute command
	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	output, err := cmd.CombinedOutput()
	
	success := err == nil
	outputStr := string(output)

	return map[string]interface{}{
		"success": success,
		"output":  outputStr,
	}, nil
}

func (p *AnsiblePlugin) parseAnsibleStats(output string) map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Look for PLAY RECAP section
	lines := strings.Split(output, "\n")
	inRecap := false
	
	for _, line := range lines {
		if strings.Contains(line, "PLAY RECAP") {
			inRecap = true
			continue
		}
		
		if inRecap && strings.TrimSpace(line) != "" {
			// Parse stats lines like: "localhost : ok=2 changed=0 unreachable=0 failed=0"
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					host := strings.TrimSpace(parts[0])
					statsStr := strings.TrimSpace(parts[1])
					
					hostStats := make(map[string]interface{})
					
					// Parse individual stats using regex
					re := regexp.MustCompile(`(\w+)=(\d+)`)
					matches := re.FindAllStringSubmatch(statsStr, -1)
					
					for _, match := range matches {
						if len(match) == 3 {
							key := match[1]
							value := match[2]
							hostStats[key] = value
						}
					}
					
					if len(hostStats) > 0 {
						stats[host] = hostStats
					}
				}
			}
		}
	}
	
	return stats
}

// Helper function to get string parameter
func getStringParam(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key].(string); ok {
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

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewAnsiblePlugin()

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
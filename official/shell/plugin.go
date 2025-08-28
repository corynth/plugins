package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
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

type ShellPlugin struct{}

func NewShellPlugin() *ShellPlugin {
	return &ShellPlugin{}
}

func (p *ShellPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "shell",
		Version:     "1.0.0",
		Description: "Execute shell commands and scripts with support for various interpreters",
		Author:      "Corynth Team",
		Tags:        []string{"shell", "command", "script", "execution", "bash", "python"},
	}
}

func (p *ShellPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"exec": {
			Description: "Execute a shell command",
			Inputs: map[string]IOSpec{
				"command": {
					Type:        "string",
					Required:    true,
					Description: "Shell command to execute",
				},
				"working_dir": {
					Type:        "string",
					Required:    false,
					Description: "Working directory for command execution",
				},
				"timeout": {
					Type:        "number",
					Required:    false,
					Default:     300,
					Description: "Timeout in seconds",
				},
				"shell": {
					Type:        "boolean",
					Required:    false,
					Default:     true,
					Description: "Use shell for execution",
				},
				"env": {
					Type:        "object",
					Required:    false,
					Description: "Environment variables as key-value pairs",
				},
			},
			Outputs: map[string]IOSpec{
				"output":    {Type: "string", Description: "Combined stdout and stderr output"},
				"stdout":    {Type: "string", Description: "Standard output"},
				"stderr":    {Type: "string", Description: "Standard error"},
				"exit_code": {Type: "number", Description: "Process exit code"},
				"success":   {Type: "boolean", Description: "Whether command succeeded (exit code 0)"},
			},
		},
		"script": {
			Description: "Execute a script with specified interpreter",
			Inputs: map[string]IOSpec{
				"script": {
					Type:        "string",
					Required:    true,
					Description: "Script content to execute",
				},
				"working_dir": {
					Type:        "string",
					Required:    false,
					Description: "Working directory for script execution",
				},
				"timeout": {
					Type:        "number",
					Required:    false,
					Default:     300,
					Description: "Timeout in seconds",
				},
				"shell_type": {
					Type:        "string",
					Required:    false,
					Default:     "bash",
					Description: "Shell/interpreter type (bash, sh, python, python3, node, etc.)",
				},
				"env": {
					Type:        "object",
					Required:    false,
					Description: "Environment variables as key-value pairs",
				},
			},
			Outputs: map[string]IOSpec{
				"output":    {Type: "string", Description: "Combined stdout and stderr output"},
				"stdout":    {Type: "string", Description: "Standard output"},
				"stderr":    {Type: "string", Description: "Standard error"},
				"exit_code": {Type: "number", Description: "Process exit code"},
				"success":   {Type: "boolean", Description: "Whether script succeeded (exit code 0)"},
			},
		},
	}
}

func (p *ShellPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "exec":
		return p.executeCommand(params)
	case "script":
		return p.executeScript(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *ShellPlugin) executeCommand(params map[string]interface{}) (map[string]interface{}, error) {
	command, ok := params["command"].(string)
	if !ok || command == "" {
		return map[string]interface{}{"error": "command parameter is required"}, nil
	}

	// Extract parameters with defaults
	workingDir := p.getStringParam(params, "working_dir", "")
	timeout := p.getFloatParam(params, "timeout", 300)
	useShell := p.getBoolParam(params, "shell", true)
	envVars := p.getMapParam(params, "env")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if useShell {
		// Use shell to execute command
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", command)
	} else {
		// Split command into parts for direct execution
		parts := strings.Fields(command)
		if len(parts) == 0 {
			return map[string]interface{}{"error": "empty command"}, nil
		}
		cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
	}

	// Set working directory if specified
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	// Set environment variables
	if len(envVars) > 0 {
		env := os.Environ()
		for key, value := range envVars {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	// Execute command and capture output
	stdout, stderr, exitCode := p.runCommand(cmd)

	return map[string]interface{}{
		"output":    stdout + stderr,
		"stdout":    stdout,
		"stderr":    stderr,
		"exit_code": exitCode,
		"success":   exitCode == 0,
	}, nil
}

func (p *ShellPlugin) executeScript(params map[string]interface{}) (map[string]interface{}, error) {
	script, ok := params["script"].(string)
	if !ok || script == "" {
		return map[string]interface{}{"error": "script parameter is required"}, nil
	}

	// Extract parameters with defaults
	workingDir := p.getStringParam(params, "working_dir", "")
	timeout := p.getFloatParam(params, "timeout", 300)
	shellType := p.getStringParam(params, "shell_type", "bash")
	envVars := p.getMapParam(params, "env")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Determine the command based on shell type
	var cmd *exec.Cmd
	switch shellType {
	case "bash":
		cmd = exec.CommandContext(ctx, "bash", "-c", script)
	case "sh":
		cmd = exec.CommandContext(ctx, "sh", "-c", script)
	case "python", "python3":
		// Create temporary file for Python script
		tmpFile, err := ioutil.TempFile("", "corynth_script_*.py")
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to create temp file: %v", err)}, nil
		}
		defer os.Remove(tmpFile.Name())
		
		if _, err := tmpFile.WriteString(script); err != nil {
			tmpFile.Close()
			return map[string]interface{}{"error": fmt.Sprintf("failed to write script: %v", err)}, nil
		}
		tmpFile.Close()
		
		cmd = exec.CommandContext(ctx, shellType, tmpFile.Name())
	case "node", "nodejs":
		// Create temporary file for Node.js script
		tmpFile, err := ioutil.TempFile("", "corynth_script_*.js")
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to create temp file: %v", err)}, nil
		}
		defer os.Remove(tmpFile.Name())
		
		if _, err := tmpFile.WriteString(script); err != nil {
			tmpFile.Close()
			return map[string]interface{}{"error": fmt.Sprintf("failed to write script: %v", err)}, nil
		}
		tmpFile.Close()
		
		cmd = exec.CommandContext(ctx, "node", tmpFile.Name())
	default:
		// For other interpreters, try to execute directly with -c flag
		cmd = exec.CommandContext(ctx, shellType, "-c", script)
	}

	// Set working directory if specified
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	// Set environment variables
	if len(envVars) > 0 {
		env := os.Environ()
		for key, value := range envVars {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	// Execute script and capture output
	stdout, stderr, exitCode := p.runCommand(cmd)

	return map[string]interface{}{
		"output":    stdout + stderr,
		"stdout":    stdout,
		"stderr":    stderr,
		"exit_code": exitCode,
		"success":   exitCode == 0,
	}, nil
}

func (p *ShellPlugin) runCommand(cmd *exec.Cmd) (stdout, stderr string, exitCode int) {
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// Command failed to start or other error
			exitCode = -1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		exitCode = 0
	}

	return stdout, stderr, exitCode
}

// Helper functions to extract parameters with type safety
func (p *ShellPlugin) getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultValue
}

func (p *ShellPlugin) getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key].(float64); ok {
		return val
	}
	return defaultValue
}

func (p *ShellPlugin) getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultValue
}

func (p *ShellPlugin) getMapParam(params map[string]interface{}, key string) map[string]string {
	result := make(map[string]string)
	if val, ok := params[key].(map[string]interface{}); ok {
		for k, v := range val {
			if str, ok := v.(string); ok {
				result[k] = str
			}
		}
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewShellPlugin()

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
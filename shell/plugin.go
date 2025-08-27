package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	
	"github.com/corynth/corynth-dist/pkg/plugin"
)

type ShellPlugin struct{}

func (p *ShellPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "shell",
		Version:     "1.0.0",
		Description: "Execute shell commands and scripts",
		Author:      "Corynth Team",
		Tags:        []string{"shell", "bash", "command", "script"},
		License:     "Apache-2.0",
	}
}

func (p *ShellPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "exec",
			Description: "Execute a shell command",
			Inputs: map[string]plugin.InputSpec{
				"command": {
					Type:        "string",
					Description: "Command to execute",
					Required:    true,
				},
				"args": {
					Type:        "array",
					Description: "Command arguments",
					Required:    false,
				},
				"env": {
					Type:        "object",
					Description: "Environment variables",
					Required:    false,
				},
				"working_dir": {
					Type:        "string",
					Description: "Working directory",
					Required:    false,
				},
				"timeout": {
					Type:        "number",
					Description: "Command timeout in seconds",
					Required:    false,
					Default:     300,
				},
				"shell": {
					Type:        "string",
					Description: "Shell to use (bash, sh, zsh)",
					Required:    false,
					Default:     "bash",
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"stdout": {
					Type:        "string",
					Description: "Standard output",
				},
				"stderr": {
					Type:        "string",
					Description: "Standard error",
				},
				"exit_code": {
					Type:        "number",
					Description: "Exit code",
				},
			},
		},
		{
			Name:        "script",
			Description: "Execute a shell script",
			Inputs: map[string]plugin.InputSpec{
				"script": {
					Type:        "string",
					Description: "Script content to execute",
					Required:    true,
				},
				"env": {
					Type:        "object",
					Description: "Environment variables",
					Required:    false,
				},
				"working_dir": {
					Type:        "string",
					Description: "Working directory",
					Required:    false,
				},
				"timeout": {
					Type:        "number",
					Description: "Script timeout in seconds",
					Required:    false,
					Default:     300,
				},
				"shell": {
					Type:        "string",
					Description: "Shell to use (bash, sh, zsh)",
					Required:    false,
					Default:     "bash",
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"stdout": {
					Type:        "string",
					Description: "Standard output",
				},
				"stderr": {
					Type:        "string",
					Description: "Standard error",
				},
				"exit_code": {
					Type:        "number",
					Description: "Exit code",
				},
			},
		},
	}
}

func (p *ShellPlugin) Validate(params map[string]interface{}) error {
	if _, ok := params["command"]; !ok {
		if _, ok := params["script"]; !ok {
			return fmt.Errorf("either 'command' or 'script' is required")
		}
	}
	return nil
}

func (p *ShellPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "exec":
		return p.executeCommand(ctx, params)
	case "script":
		return p.executeScript(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *ShellPlugin) executeCommand(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	command, ok := params["command"].(string)
	if !ok || command == "" {
		return nil, fmt.Errorf("command parameter is required")
	}

	shell := "bash"
	if s, ok := params["shell"].(string); ok {
		shell = s
	}

	// Set timeout
	timeout := 300
	if t, ok := params["timeout"].(float64); ok {
		timeout = int(t)
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Build command with arguments if provided
	var cmd *exec.Cmd
	if args, ok := params["args"].([]interface{}); ok {
		argStrings := make([]string, len(args))
		for i, arg := range args {
			argStrings[i] = fmt.Sprintf("%v", arg)
		}
		fullCommand := command + " " + strings.Join(argStrings, " ")
		cmd = exec.CommandContext(timeoutCtx, shell, "-c", fullCommand)
	} else {
		cmd = exec.CommandContext(timeoutCtx, shell, "-c", command)
	}

	// Set environment variables
	cmd.Env = os.Environ()
	if env, ok := params["env"].(map[string]interface{}); ok {
		for key, value := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%v", key, value))
		}
	}

	// Set working directory
	if workingDir, ok := params["working_dir"].(string); ok {
		cmd.Dir = workingDir
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute command
	err := cmd.Run()
	
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return nil, fmt.Errorf("command execution failed: %w", err)
		}
	}

	return map[string]interface{}{
		"stdout":    stdout.String(),
		"stderr":    stderr.String(),
		"exit_code": exitCode,
	}, nil
}

func (p *ShellPlugin) executeScript(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	script, ok := params["script"].(string)
	if !ok || script == "" {
		return nil, fmt.Errorf("script parameter is required")
	}

	shell := "bash"
	if s, ok := params["shell"].(string); ok {
		shell = s
	}

	// Set timeout
	timeout := 300
	if t, ok := params["timeout"].(float64); ok {
		timeout = int(t)
	}

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Create temporary script file
	tmpFile, err := os.CreateTemp("", "corynth-script-*.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary script file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write script content
	if _, err := tmpFile.WriteString(script); err != nil {
		return nil, fmt.Errorf("failed to write script content: %w", err)
	}
	tmpFile.Close()

	// Make script executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return nil, fmt.Errorf("failed to make script executable: %w", err)
	}

	// Build command
	cmd := exec.CommandContext(timeoutCtx, shell, tmpFile.Name())

	// Set environment variables
	cmd.Env = os.Environ()
	if env, ok := params["env"].(map[string]interface{}); ok {
		for key, value := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%v", key, value))
		}
	}

	// Set working directory
	if workingDir, ok := params["working_dir"].(string); ok {
		cmd.Dir = workingDir
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute script
	err = cmd.Run()
	
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return nil, fmt.Errorf("script execution failed: %w", err)
		}
	}

	return map[string]interface{}{
		"stdout":    stdout.String(),
		"stderr":    stderr.String(),
		"exit_code": exitCode,
	}, nil
}

var ExportedPlugin plugin.Plugin = &ShellPlugin{}
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
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

type DockerPlugin struct{}

func NewDockerPlugin() *DockerPlugin {
	return &DockerPlugin{}
}

func (p *DockerPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "docker",
		Version:     "1.0.0",
		Description: "Docker container operations and image management",
		Author:      "Corynth Team",
		Tags:        []string{"docker", "containers", "containerization", "devops"},
	}
}

func (p *DockerPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"run": {
			Description: "Run Docker container with ports, volumes, env vars",
			Inputs: map[string]IOSpec{
				"image":    {Type: "string", Required: true, Description: "Container image name"},
				"name":     {Type: "string", Required: false, Description: "Container name"},
				"detach":   {Type: "boolean", Required: false, Default: true, Description: "Run in detached mode"},
				"ports":    {Type: "array", Required: false, Description: "Port mappings (e.g., ['8080:80'])"},
				"volumes":  {Type: "array", Required: false, Description: "Volume mappings (e.g., ['/host:/container'])"},
				"env":      {Type: "object", Required: false, Description: "Environment variables"},
				"command":  {Type: "string", Required: false, Description: "Command to run"},
				"network":  {Type: "string", Required: false, Description: "Network to connect to"},
				"remove":   {Type: "boolean", Required: false, Default: false, Description: "Remove container when it exits"},
			},
			Outputs: map[string]IOSpec{
				"container_id": {Type: "string", Description: "Container ID"},
				"name":         {Type: "string", Description: "Container name"},
				"success":      {Type: "boolean", Description: "Operation success"},
			},
		},
		"ps": {
			Description: "List containers with filtering",
			Inputs: map[string]IOSpec{
				"all":    {Type: "boolean", Required: false, Default: false, Description: "Show all containers"},
				"filter": {Type: "string", Required: false, Description: "Filter containers (e.g., 'status=running')"},
			},
			Outputs: map[string]IOSpec{
				"containers": {Type: "array", Description: "List of container information"},
			},
		},
		"stop": {
			Description: "Stop running containers",
			Inputs: map[string]IOSpec{
				"container": {Type: "string", Required: true, Description: "Container ID or name"},
				"timeout":   {Type: "number", Required: false, Default: 10, Description: "Seconds to wait before killing"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Operation success"},
			},
		},
		"start": {
			Description: "Start stopped containers",
			Inputs: map[string]IOSpec{
				"container": {Type: "string", Required: true, Description: "Container ID or name"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Operation success"},
			},
		},
		"logs": {
			Description: "Get container logs with tail/follow",
			Inputs: map[string]IOSpec{
				"container": {Type: "string", Required: true, Description: "Container ID or name"},
				"tail":      {Type: "number", Required: false, Description: "Number of lines to show from end"},
				"follow":    {Type: "boolean", Required: false, Default: false, Description: "Follow log output"},
			},
			Outputs: map[string]IOSpec{
				"logs": {Type: "string", Description: "Container logs"},
			},
		},
		"exec": {
			Description: "Execute commands in containers",
			Inputs: map[string]IOSpec{
				"container":   {Type: "string", Required: true, Description: "Container ID or name"},
				"command":     {Type: "string", Required: true, Description: "Command to execute"},
				"interactive": {Type: "boolean", Required: false, Default: false, Description: "Interactive mode"},
			},
			Outputs: map[string]IOSpec{
				"output":    {Type: "string", Description: "Command output"},
				"exit_code": {Type: "number", Description: "Exit code"},
			},
		},
		"build": {
			Description: "Build Docker images from context",
			Inputs: map[string]IOSpec{
				"path":       {Type: "string", Required: true, Description: "Build context path"},
				"tag":        {Type: "string", Required: false, Description: "Image tag"},
				"dockerfile": {Type: "string", Required: false, Description: "Dockerfile path"},
				"args":       {Type: "object", Required: false, Description: "Build arguments"},
			},
			Outputs: map[string]IOSpec{
				"image_id": {Type: "string", Description: "Built image ID"},
				"success":  {Type: "boolean", Description: "Build success"},
			},
		},
		"images": {
			Description: "List Docker images",
			Inputs: map[string]IOSpec{
				"all": {Type: "boolean", Required: false, Default: false, Description: "Show all images"},
			},
			Outputs: map[string]IOSpec{
				"images": {Type: "array", Description: "List of image information"},
			},
		},
	}
}

func (p *DockerPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "run":
		return p.runContainer(params)
	case "ps":
		return p.listContainers(params)
	case "stop":
		return p.stopContainer(params)
	case "start":
		return p.startContainer(params)
	case "logs":
		return p.getContainerLogs(params)
	case "exec":
		return p.execCommand(params)
	case "build":
		return p.buildImage(params)
	case "images":
		return p.listImages(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *DockerPlugin) runContainer(params map[string]interface{}) (map[string]interface{}, error) {
	image, ok := params["image"].(string)
	if !ok || image == "" {
		return map[string]interface{}{"error": "image is required"}, nil
	}

	args := []string{"run"}
	
	if getBoolParam(params, "detach", true) {
		args = append(args, "-d")
	}
	
	if name, ok := params["name"].(string); ok && name != "" {
		args = append(args, "--name", name)
	}
	
	if getBoolParam(params, "remove", false) {
		args = append(args, "--rm")
	}
	
	if network, ok := params["network"].(string); ok && network != "" {
		args = append(args, "--network", network)
	}
	
	// Add port mappings
	if ports, ok := params["ports"].([]interface{}); ok {
		for _, port := range ports {
			if p, ok := port.(string); ok {
				args = append(args, "-p", p)
			}
		}
	}
	
	// Add volume mappings
	if volumes, ok := params["volumes"].([]interface{}); ok {
		for _, volume := range volumes {
			if v, ok := volume.(string); ok {
				args = append(args, "-v", v)
			}
		}
	}
	
	// Add environment variables
	if envVars, ok := params["env"].(map[string]interface{}); ok {
		for key, value := range envVars {
			args = append(args, "-e", fmt.Sprintf("%s=%v", key, value))
		}
	}
	
	args = append(args, image)
	
	// Add command if provided
	if command, ok := params["command"].(string); ok && command != "" {
		args = append(args, "sh", "-c", command)
	}
	
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return map[string]interface{}{
			"error":   err.Error(),
			"output":  string(output),
			"success": false,
		}, nil
	}
	
	containerID := strings.TrimSpace(string(output))
	
	// Get container name if not provided
	containerName := ""
	if name, ok := params["name"].(string); ok {
		containerName = name
	} else if containerID != "" {
		// Get name from docker inspect
		inspectCmd := exec.Command("docker", "inspect", "--format={{.Name}}", containerID)
		if nameOutput, err := inspectCmd.Output(); err == nil {
			containerName = strings.TrimPrefix(strings.TrimSpace(string(nameOutput)), "/")
		}
	}
	
	return map[string]interface{}{
		"container_id": containerID,
		"name":         containerName,
		"success":      true,
	}, nil
}

func (p *DockerPlugin) listContainers(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"ps", "--format", "json"}
	
	if getBoolParam(params, "all", false) {
		args = append(args, "-a")
	}
	
	if filter, ok := params["filter"].(string); ok && filter != "" {
		args = append(args, "--filter", filter)
	}
	
	cmd := exec.Command("docker", args...)
	output, err := cmd.Output()
	
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}, nil
	}
	
	containers := []map[string]interface{}{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	
	for scanner.Scan() {
		var container map[string]interface{}
		if err := json.Unmarshal([]byte(scanner.Text()), &container); err == nil {
			containers = append(containers, container)
		}
	}
	
	return map[string]interface{}{
		"containers": containers,
	}, nil
}

func (p *DockerPlugin) stopContainer(params map[string]interface{}) (map[string]interface{}, error) {
	container, ok := params["container"].(string)
	if !ok || container == "" {
		return map[string]interface{}{"error": "container is required"}, nil
	}
	
	args := []string{"stop"}
	
	if timeout, ok := params["timeout"].(float64); ok {
		args = append(args, "-t", fmt.Sprintf("%.0f", timeout))
	}
	
	args = append(args, container)
	
	cmd := exec.Command("docker", args...)
	err := cmd.Run()
	
	return map[string]interface{}{
		"success": err == nil,
	}, nil
}

func (p *DockerPlugin) startContainer(params map[string]interface{}) (map[string]interface{}, error) {
	container, ok := params["container"].(string)
	if !ok || container == "" {
		return map[string]interface{}{"error": "container is required"}, nil
	}
	
	cmd := exec.Command("docker", "start", container)
	err := cmd.Run()
	
	return map[string]interface{}{
		"success": err == nil,
	}, nil
}

func (p *DockerPlugin) getContainerLogs(params map[string]interface{}) (map[string]interface{}, error) {
	container, ok := params["container"].(string)
	if !ok || container == "" {
		return map[string]interface{}{"error": "container is required"}, nil
	}
	
	args := []string{"logs"}
	
	if tail, ok := params["tail"].(float64); ok {
		args = append(args, "--tail", fmt.Sprintf("%.0f", tail))
	}
	
	if getBoolParam(params, "follow", false) {
		args = append(args, "-f")
	}
	
	args = append(args, container)
	
	cmd := exec.Command("docker", args...)
	output, err := cmd.Output()
	
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}, nil
	}
	
	return map[string]interface{}{
		"logs": string(output),
	}, nil
}

func (p *DockerPlugin) execCommand(params map[string]interface{}) (map[string]interface{}, error) {
	container, ok := params["container"].(string)
	if !ok || container == "" {
		return map[string]interface{}{"error": "container is required"}, nil
	}
	
	command, ok := params["command"].(string)
	if !ok || command == "" {
		return map[string]interface{}{"error": "command is required"}, nil
	}
	
	args := []string{"exec"}
	
	if getBoolParam(params, "interactive", false) {
		args = append(args, "-it")
	}
	
	args = append(args, container, "sh", "-c", command)
	
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	
	return map[string]interface{}{
		"output":    string(output),
		"exit_code": exitCode,
	}, nil
}

func (p *DockerPlugin) buildImage(params map[string]interface{}) (map[string]interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return map[string]interface{}{"error": "path is required"}, nil
	}
	
	args := []string{"build"}
	
	if tag, ok := params["tag"].(string); ok && tag != "" {
		args = append(args, "-t", tag)
	}
	
	if dockerfile, ok := params["dockerfile"].(string); ok && dockerfile != "" {
		args = append(args, "-f", dockerfile)
	}
	
	// Add build arguments
	if buildArgs, ok := params["args"].(map[string]interface{}); ok {
		for key, value := range buildArgs {
			args = append(args, "--build-arg", fmt.Sprintf("%s=%v", key, value))
		}
	}
	
	args = append(args, path)
	
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return map[string]interface{}{
			"error":   err.Error(),
			"output":  string(output),
			"success": false,
		}, nil
	}
	
	// Extract image ID from output
	imageID := ""
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Successfully built") {
			parts := strings.Fields(line)
			if len(parts) > 2 {
				imageID = parts[len(parts)-1]
			}
		}
	}
	
	return map[string]interface{}{
		"image_id": imageID,
		"success":  true,
	}, nil
}

func (p *DockerPlugin) listImages(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"images", "--format", "json"}
	
	if getBoolParam(params, "all", false) {
		args = append(args, "-a")
	}
	
	cmd := exec.Command("docker", args...)
	output, err := cmd.Output()
	
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}, nil
	}
	
	images := []map[string]interface{}{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	
	for scanner.Scan() {
		var image map[string]interface{}
		if err := json.Unmarshal([]byte(scanner.Text()), &image); err == nil {
			images = append(images, image)
		}
	}
	
	return map[string]interface{}{
		"images": images,
	}, nil
}

// Helper functions
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
	plugin := NewDockerPlugin()

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
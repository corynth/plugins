package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
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

type KubernetesPlugin struct{}

func NewKubernetesPlugin() *KubernetesPlugin {
	return &KubernetesPlugin{}
}

func (p *KubernetesPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "kubernetes",
		Version:     "1.0.0",
		Description: "Kubernetes cluster management and resource operations",
		Author:      "Corynth Team",
		Tags:        []string{"kubernetes", "k8s", "container", "orchestration", "cloud-native"},
	}
}

func (p *KubernetesPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"apply": {
			Description: "Apply Kubernetes manifests",
			Inputs: map[string]IOSpec{
				"manifest":  {Type: "string", Required: false, Description: "YAML manifest content"},
				"file":      {Type: "string", Required: false, Description: "Path to manifest file"},
				"namespace": {Type: "string", Required: false, Description: "Target namespace"},
				"dry_run":   {Type: "boolean", Required: false, Default: false, Description: "Dry run mode"},
			},
			Outputs: map[string]IOSpec{
				"success":   {Type: "boolean", Description: "Operation success"},
				"resources": {Type: "array", Description: "Applied resources"},
			},
		},
		"get": {
			Description: "Get Kubernetes resources",
			Inputs: map[string]IOSpec{
				"resource":       {Type: "string", Required: true, Description: "Resource type (pods, services, etc.)"},
				"name":           {Type: "string", Required: false, Description: "Resource name"},
				"namespace":      {Type: "string", Required: false, Description: "Target namespace"},
				"all_namespaces": {Type: "boolean", Required: false, Default: false, Description: "All namespaces"},
				"selector":       {Type: "string", Required: false, Description: "Label selector"},
				"output":         {Type: "string", Required: false, Default: "json", Description: "Output format"},
			},
			Outputs: map[string]IOSpec{
				"resources": {Type: "array", Description: "Resource information"},
			},
		},
		"describe": {
			Description: "Describe Kubernetes resources",
			Inputs: map[string]IOSpec{
				"resource":  {Type: "string", Required: true, Description: "Resource type"},
				"name":      {Type: "string", Required: true, Description: "Resource name"},
				"namespace": {Type: "string", Required: false, Description: "Target namespace"},
			},
			Outputs: map[string]IOSpec{
				"description": {Type: "string", Description: "Resource description"},
			},
		},
		"scale": {
			Description: "Scale deployments or replica sets",
			Inputs: map[string]IOSpec{
				"resource":  {Type: "string", Required: true, Description: "Resource type (deployment, replicaset)"},
				"name":      {Type: "string", Required: true, Description: "Resource name"},
				"replicas":  {Type: "number", Required: true, Description: "Number of replicas"},
				"namespace": {Type: "string", Required: false, Description: "Target namespace"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Scaling success"},
			},
		},
		"logs": {
			Description: "Get pod logs",
			Inputs: map[string]IOSpec{
				"pod":       {Type: "string", Required: true, Description: "Pod name"},
				"container": {Type: "string", Required: false, Description: "Container name"},
				"namespace": {Type: "string", Required: false, Description: "Target namespace"},
				"tail":      {Type: "number", Required: false, Description: "Number of lines"},
				"follow":    {Type: "boolean", Required: false, Default: false, Description: "Follow logs"},
				"previous":  {Type: "boolean", Required: false, Default: false, Description: "Previous container logs"},
			},
			Outputs: map[string]IOSpec{
				"logs": {Type: "string", Description: "Pod logs"},
			},
		},
		"exec": {
			Description: "Execute commands in pods",
			Inputs: map[string]IOSpec{
				"pod":       {Type: "string", Required: true, Description: "Pod name"},
				"container": {Type: "string", Required: false, Description: "Container name"},
				"command":   {Type: "string", Required: true, Description: "Command to execute"},
				"namespace": {Type: "string", Required: false, Description: "Target namespace"},
			},
			Outputs: map[string]IOSpec{
				"output":    {Type: "string", Description: "Command output"},
				"exit_code": {Type: "number", Description: "Exit code"},
			},
		},
		"port_forward": {
			Description: "Forward local ports to pod",
			Inputs: map[string]IOSpec{
				"pod":          {Type: "string", Required: true, Description: "Pod name"},
				"port_mapping": {Type: "string", Required: true, Description: "Port mapping (e.g., '8080:80')"},
				"namespace":    {Type: "string", Required: false, Description: "Target namespace"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Port forward success"},
			},
		},
		"delete": {
			Description: "Delete Kubernetes resources",
			Inputs: map[string]IOSpec{
				"resource":  {Type: "string", Required: true, Description: "Resource type"},
				"name":      {Type: "string", Required: false, Description: "Resource name"},
				"file":      {Type: "string", Required: false, Description: "Manifest file to delete"},
				"namespace": {Type: "string", Required: false, Description: "Target namespace"},
				"selector":  {Type: "string", Required: false, Description: "Label selector"},
				"force":     {Type: "boolean", Required: false, Default: false, Description: "Force deletion"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Deletion success"},
			},
		},
	}
}

func (p *KubernetesPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "apply":
		return p.applyManifest(params)
	case "get":
		return p.getResources(params)
	case "describe":
		return p.describeResource(params)
	case "scale":
		return p.scaleResource(params)
	case "logs":
		return p.getLogs(params)
	case "exec":
		return p.execCommand(params)
	case "port_forward":
		return p.portForward(params)
	case "delete":
		return p.deleteResources(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

// runKubectlCommand runs kubectl command with proper error handling
func (p *KubernetesPlugin) runKubectlCommand(args []string, inputData string) (string, string, error) {
	cmd := exec.Command("kubectl", args...)

	if inputData != "" {
		cmd.Stdin = strings.NewReader(inputData)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("failed to start kubectl: %v", err)
	}

	stdoutBytes, _ := io.ReadAll(stdout)
	stderrBytes, _ := io.ReadAll(stderr)

	err = cmd.Wait()

	return string(stdoutBytes), string(stderrBytes), err
}

func (p *KubernetesPlugin) applyManifest(params map[string]interface{}) (map[string]interface{}, error) {
	manifest, _ := params["manifest"].(string)
	filePath, _ := params["file"].(string)
	namespace, _ := params["namespace"].(string)
	dryRun := getBoolParam(params, "dry_run", false)

	args := []string{"apply"}

	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if dryRun {
		args = append(args, "--dry-run=client")
	}

	var inputData string
	if manifest != "" {
		args = append(args, "-f", "-")
		inputData = manifest
	} else if filePath != "" {
		args = append(args, "-f", filePath)
	} else {
		return map[string]interface{}{"error": "Either manifest or file parameter is required"}, nil
	}

	stdout, stderr, err := p.runKubectlCommand(args, inputData)

	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   stderr,
		}, nil
	}

	resources := []string{}
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && (strings.Contains(line, "configured") || strings.Contains(line, "created") || strings.Contains(line, "unchanged")) {
			resources = append(resources, line)
		}
	}

	return map[string]interface{}{
		"success":   true,
		"resources": resources,
	}, nil
}

func (p *KubernetesPlugin) getResources(params map[string]interface{}) (map[string]interface{}, error) {
	resource, ok := params["resource"].(string)
	if !ok || resource == "" {
		return map[string]interface{}{"error": "resource is required"}, nil
	}

	name, _ := params["name"].(string)
	namespace, _ := params["namespace"].(string)
	allNamespaces := getBoolParam(params, "all_namespaces", false)
	selector, _ := params["selector"].(string)
	outputFormat := getStringParam(params, "output", "json")

	args := []string{"get", resource}

	if name != "" {
		args = append(args, name)
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else if allNamespaces {
		args = append(args, "--all-namespaces")
	}
	if selector != "" {
		args = append(args, "-l", selector)
	}

	args = append(args, "-o", outputFormat)

	stdout, stderr, err := p.runKubectlCommand(args, "")

	if err != nil {
		return map[string]interface{}{"error": stderr}, nil
	}

	if outputFormat == "json" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(stdout), &data); err == nil {
			if items, ok := data["items"].([]interface{}); ok {
				return map[string]interface{}{"resources": items}, nil
			} else {
				return map[string]interface{}{"resources": []interface{}{data}}, nil
			}
		} else {
			return map[string]interface{}{"resources": stdout}, nil
		}
	} else {
		return map[string]interface{}{"resources": stdout}, nil
	}
}

func (p *KubernetesPlugin) describeResource(params map[string]interface{}) (map[string]interface{}, error) {
	resource, ok := params["resource"].(string)
	if !ok || resource == "" {
		return map[string]interface{}{"error": "resource is required"}, nil
	}

	name, ok := params["name"].(string)
	if !ok || name == "" {
		return map[string]interface{}{"error": "name is required"}, nil
	}

	namespace, _ := params["namespace"].(string)

	args := []string{"describe", resource, name}

	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	stdout, stderr, err := p.runKubectlCommand(args, "")

	if err != nil {
		return map[string]interface{}{"error": stderr}, nil
	}

	return map[string]interface{}{"description": stdout}, nil
}

func (p *KubernetesPlugin) scaleResource(params map[string]interface{}) (map[string]interface{}, error) {
	resource, ok := params["resource"].(string)
	if !ok || resource == "" {
		return map[string]interface{}{"error": "resource is required"}, nil
	}

	name, ok := params["name"].(string)
	if !ok || name == "" {
		return map[string]interface{}{"error": "name is required"}, nil
	}

	replicas, ok := params["replicas"].(float64)
	if !ok {
		return map[string]interface{}{"error": "replicas is required"}, nil
	}

	namespace, _ := params["namespace"].(string)

	args := []string{"scale", resource, name, fmt.Sprintf("--replicas=%v", int(replicas))}

	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	_, _, err := p.runKubectlCommand(args, "")

	return map[string]interface{}{
		"success": err == nil,
	}, nil
}

func (p *KubernetesPlugin) getLogs(params map[string]interface{}) (map[string]interface{}, error) {
	pod, ok := params["pod"].(string)
	if !ok || pod == "" {
		return map[string]interface{}{"error": "pod is required"}, nil
	}

	container, _ := params["container"].(string)
	namespace, _ := params["namespace"].(string)
	tail, _ := params["tail"].(float64)
	follow := getBoolParam(params, "follow", false)
	previous := getBoolParam(params, "previous", false)

	args := []string{"logs", pod}

	if container != "" {
		args = append(args, "-c", container)
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if tail > 0 {
		args = append(args, "--tail", strconv.Itoa(int(tail)))
	}
	if follow {
		args = append(args, "-f")
	}
	if previous {
		args = append(args, "-p")
	}

	stdout, stderr, err := p.runKubectlCommand(args, "")

	if err != nil {
		return map[string]interface{}{"error": stderr}, nil
	}

	return map[string]interface{}{"logs": stdout}, nil
}

func (p *KubernetesPlugin) execCommand(params map[string]interface{}) (map[string]interface{}, error) {
	pod, ok := params["pod"].(string)
	if !ok || pod == "" {
		return map[string]interface{}{"error": "pod is required"}, nil
	}

	command, ok := params["command"].(string)
	if !ok || command == "" {
		return map[string]interface{}{"error": "command is required"}, nil
	}

	container, _ := params["container"].(string)
	namespace, _ := params["namespace"].(string)

	args := []string{"exec", pod}

	if container != "" {
		args = append(args, "-c", container)
	}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	args = append(args, "--")
	args = append(args, strings.Fields(command)...)

	stdout, _, err := p.runKubectlCommand(args, "")

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	return map[string]interface{}{
		"output":    stdout,
		"exit_code": exitCode,
	}, nil
}

func (p *KubernetesPlugin) portForward(params map[string]interface{}) (map[string]interface{}, error) {
	pod, ok := params["pod"].(string)
	if !ok || pod == "" {
		return map[string]interface{}{"error": "pod is required"}, nil
	}

	portMapping, ok := params["port_mapping"].(string)
	if !ok || portMapping == "" {
		return map[string]interface{}{"error": "port_mapping is required"}, nil
	}

	namespace, _ := params["namespace"].(string)

	args := []string{"port-forward", pod, portMapping}

	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	// Note: This is a basic implementation. In practice, port-forward runs continuously
	// For workflow use, you might want to run this in background or with timeout
	_, _, err := p.runKubectlCommand(args, "")

	return map[string]interface{}{
		"success": err == nil,
	}, nil
}

func (p *KubernetesPlugin) deleteResources(params map[string]interface{}) (map[string]interface{}, error) {
	resource, ok := params["resource"].(string)
	if !ok || resource == "" {
		return map[string]interface{}{"error": "resource is required"}, nil
	}

	name, _ := params["name"].(string)
	filePath, _ := params["file"].(string)
	namespace, _ := params["namespace"].(string)
	selector, _ := params["selector"].(string)
	force := getBoolParam(params, "force", false)

	args := []string{"delete"}

	if filePath != "" {
		args = append(args, "-f", filePath)
	} else if name != "" {
		args = append(args, resource, name)
	} else if selector != "" {
		args = append(args, resource, "-l", selector)
	} else {
		return map[string]interface{}{"error": "name, file, or selector parameter is required"}, nil
	}

	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	if force {
		args = append(args, "--force")
	}

	_, _, err := p.runKubectlCommand(args, "")

	return map[string]interface{}{
		"success": err == nil,
	}, nil
}

// Helper functions
func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultValue
}

func getStringParam(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key].(string); ok {
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
	plugin := NewKubernetesPlugin()

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

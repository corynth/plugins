package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type KubernetesPlugin struct{}

func (p *KubernetesPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "kubernetes",
        Version:     "1.0.0",
        Description: "Kubernetes cluster management and resource operations",
        Author:      "Corynth Team",
        Tags:        []string{"kubernetes", "k8s", "container", "orchestration", "cloud-native"},
        License:     "Apache-2.0",
    }
}

func (p *KubernetesPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "apply",
            Description: "Apply Kubernetes manifests to cluster",
            Inputs: map[string]plugin.InputSpec{
                "manifest": {
                    Type:        "string",
                    Description: "YAML manifest content or file path",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace (default: default)",
                    Required:    false,
                    Default:     "default",
                },
                "kubeconfig": {
                    Type:        "string",
                    Description: "Path to kubeconfig file",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "status": {
                    Type:        "string",
                    Description: "Apply operation status",
                },
                "resources": {
                    Type:        "array",
                    Description: "Applied resources",
                },
            },
        },
        {
            Name:        "get",
            Description: "Get Kubernetes resources",
            Inputs: map[string]plugin.InputSpec{
                "resource": {
                    Type:        "string",
                    Description: "Resource type (pods, services, deployments, etc.)",
                    Required:    true,
                },
                "name": {
                    Type:        "string",
                    Description: "Resource name (optional for listing all)",
                    Required:    false,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace (default: default)",
                    Required:    false,
                    Default:     "default",
                },
                "kubeconfig": {
                    Type:        "string",
                    Description: "Path to kubeconfig file",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "resources": {
                    Type:        "array",
                    Description: "Retrieved resources",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of resources found",
                },
            },
        },
        {
            Name:        "delete",
            Description: "Delete Kubernetes resources",
            Inputs: map[string]plugin.InputSpec{
                "resource": {
                    Type:        "string",
                    Description: "Resource type (pods, services, deployments, etc.)",
                    Required:    true,
                },
                "name": {
                    Type:        "string",
                    Description: "Resource name",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace (default: default)",
                    Required:    false,
                    Default:     "default",
                },
                "kubeconfig": {
                    Type:        "string",
                    Description: "Path to kubeconfig file",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "status": {
                    Type:        "string",
                    Description: "Delete operation status",
                },
                "deleted": {
                    Type:        "boolean",
                    Description: "Whether resource was deleted",
                },
            },
        },
        {
            Name:        "scale",
            Description: "Scale deployments or replica sets",
            Inputs: map[string]plugin.InputSpec{
                "resource": {
                    Type:        "string",
                    Description: "Resource type (deployment, replicaset)",
                    Required:    true,
                },
                "name": {
                    Type:        "string",
                    Description: "Resource name",
                    Required:    true,
                },
                "replicas": {
                    Type:        "number",
                    Description: "Target replica count",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace (default: default)",
                    Required:    false,
                    Default:     "default",
                },
                "kubeconfig": {
                    Type:        "string",
                    Description: "Path to kubeconfig file",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "status": {
                    Type:        "string",
                    Description: "Scale operation status",
                },
                "replicas": {
                    Type:        "number",
                    Description: "Current replica count",
                },
            },
        },
        {
            Name:        "logs",
            Description: "Get pod logs",
            Inputs: map[string]plugin.InputSpec{
                "pod": {
                    Type:        "string",
                    Description: "Pod name",
                    Required:    true,
                },
                "container": {
                    Type:        "string",
                    Description: "Container name (optional)",
                    Required:    false,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace (default: default)",
                    Required:    false,
                    Default:     "default",
                },
                "tail": {
                    Type:        "number",
                    Description: "Number of recent lines to show",
                    Required:    false,
                    Default:     100,
                },
                "kubeconfig": {
                    Type:        "string",
                    Description: "Path to kubeconfig file",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "logs": {
                    Type:        "string",
                    Description: "Pod logs content",
                },
                "lines": {
                    Type:        "number",
                    Description: "Number of log lines retrieved",
                },
            },
        },
        {
            Name:        "wait",
            Description: "Wait for resource condition",
            Inputs: map[string]plugin.InputSpec{
                "resource": {
                    Type:        "string",
                    Description: "Resource type",
                    Required:    true,
                },
                "name": {
                    Type:        "string",
                    Description: "Resource name",
                    Required:    true,
                },
                "condition": {
                    Type:        "string",
                    Description: "Condition to wait for (available, ready, etc.)",
                    Required:    true,
                },
                "timeout": {
                    Type:        "number",
                    Description: "Timeout in seconds (default: 300)",
                    Required:    false,
                    Default:     300,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace (default: default)",
                    Required:    false,
                    Default:     "default",
                },
                "kubeconfig": {
                    Type:        "string",
                    Description: "Path to kubeconfig file",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "ready": {
                    Type:        "boolean",
                    Description: "Whether condition was met",
                },
                "status": {
                    Type:        "string",
                    Description: "Wait operation result",
                },
            },
        },
    }
}

func (p *KubernetesPlugin) Validate(params map[string]interface{}) error {
    // Check if kubectl is available
    if _, err := exec.LookPath("kubectl"); err != nil {
        return fmt.Errorf("kubectl is not installed or not in PATH")
    }
    return nil
}

func (p *KubernetesPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "apply":
        return p.executeApply(ctx, params)
    case "get":
        return p.executeGet(ctx, params)
    case "delete":
        return p.executeDelete(ctx, params)
    case "scale":
        return p.executeScale(ctx, params)
    case "logs":
        return p.executeLogs(ctx, params)
    case "wait":
        return p.executeWait(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *KubernetesPlugin) executeApply(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    manifest, _ := params["manifest"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := []string{"apply", "-n", namespace}
    
    if kubeconfig, ok := params["kubeconfig"].(string); ok && kubeconfig != "" {
        args = append([]string{"--kubeconfig", kubeconfig}, args...)
    }
    
    // Check if manifest is a file path or content
    var manifestFile string
    if strings.Contains(manifest, "\n") || strings.HasPrefix(manifest, "apiVersion:") {
        // Content - write to temp file
        tmpFile, err := os.CreateTemp("", "k8s-manifest-*.yaml")
        if err != nil {
            return nil, fmt.Errorf("failed to create temp file: %w", err)
        }
        defer os.Remove(tmpFile.Name())
        
        if _, err := tmpFile.WriteString(manifest); err != nil {
            return nil, fmt.Errorf("failed to write manifest: %w", err)
        }
        tmpFile.Close()
        manifestFile = tmpFile.Name()
    } else {
        // File path
        if _, err := os.Stat(manifest); err != nil {
            return nil, fmt.Errorf("manifest file not found: %s", manifest)
        }
        manifestFile = manifest
    }
    
    args = append(args, "-f", manifestFile)
    
    cmd := exec.CommandContext(ctx, "kubectl", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "status": "failed",
            "error":  string(output),
        }, nil
    }
    
    // Parse applied resources from output
    resources := []string{}
    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" && (strings.Contains(line, "created") || strings.Contains(line, "configured") || strings.Contains(line, "unchanged")) {
            resources = append(resources, line)
        }
    }
    
    return map[string]interface{}{
        "status":    "success",
        "output":    string(output),
        "resources": resources,
    }, nil
}

func (p *KubernetesPlugin) executeGet(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    resource, _ := params["resource"].(string)
    name, _ := params["name"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := []string{"get", resource, "-n", namespace, "-o", "json"}
    
    if kubeconfig, ok := params["kubeconfig"].(string); ok && kubeconfig != "" {
        args = append([]string{"--kubeconfig", kubeconfig}, args...)
    }
    
    if name != "" {
        args = append(args, name)
    }
    
    cmd := exec.CommandContext(ctx, "kubectl", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "resources": []interface{}{},
            "count":     0,
            "error":     err.Error(),
        }, nil
    }
    
    // Parse JSON output
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse kubectl output: %w", err)
    }
    
    var resources []interface{}
    var count int
    
    if items, ok := result["items"].([]interface{}); ok {
        // List response
        resources = items
        count = len(items)
    } else {
        // Single resource response
        resources = []interface{}{result}
        count = 1
    }
    
    return map[string]interface{}{
        "resources": resources,
        "count":     count,
    }, nil
}

func (p *KubernetesPlugin) executeDelete(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    resource, _ := params["resource"].(string)
    name, _ := params["name"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := []string{"delete", resource, name, "-n", namespace}
    
    if kubeconfig, ok := params["kubeconfig"].(string); ok && kubeconfig != "" {
        args = append([]string{"--kubeconfig", kubeconfig}, args...)
    }
    
    cmd := exec.CommandContext(ctx, "kubectl", args...)
    output, err := cmd.CombinedOutput()
    
    deleted := err == nil
    status := "success"
    if err != nil {
        status = "failed"
    }
    
    return map[string]interface{}{
        "status":  status,
        "deleted": deleted,
        "output":  string(output),
    }, nil
}

func (p *KubernetesPlugin) executeScale(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    resource, _ := params["resource"].(string)
    name, _ := params["name"].(string)
    replicas, _ := params["replicas"].(float64)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := []string{"scale", resource, name, fmt.Sprintf("--replicas=%d", int(replicas)), "-n", namespace}
    
    if kubeconfig, ok := params["kubeconfig"].(string); ok && kubeconfig != "" {
        args = append([]string{"--kubeconfig", kubeconfig}, args...)
    }
    
    cmd := exec.CommandContext(ctx, "kubectl", args...)
    output, err := cmd.CombinedOutput()
    
    status := "success"
    if err != nil {
        status = "failed"
    }
    
    return map[string]interface{}{
        "status":   status,
        "replicas": int(replicas),
        "output":   string(output),
    }, nil
}

func (p *KubernetesPlugin) executeLogs(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    pod, _ := params["pod"].(string)
    container, _ := params["container"].(string)
    namespace, _ := params["namespace"].(string)
    tail, _ := params["tail"].(float64)
    if namespace == "" {
        namespace = "default"
    }
    if tail == 0 {
        tail = 100
    }
    
    args := []string{"logs", pod, "-n", namespace, fmt.Sprintf("--tail=%d", int(tail))}
    
    if kubeconfig, ok := params["kubeconfig"].(string); ok && kubeconfig != "" {
        args = append([]string{"--kubeconfig", kubeconfig}, args...)
    }
    
    if container != "" {
        args = append(args, "-c", container)
    }
    
    cmd := exec.CommandContext(ctx, "kubectl", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "logs":  "",
            "lines": 0,
            "error": err.Error(),
        }, nil
    }
    
    logs := string(output)
    lines := len(strings.Split(strings.TrimSpace(logs), "\n"))
    
    return map[string]interface{}{
        "logs":  logs,
        "lines": lines,
    }, nil
}

func (p *KubernetesPlugin) executeWait(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    resource, _ := params["resource"].(string)
    name, _ := params["name"].(string)
    condition, _ := params["condition"].(string)
    timeout, _ := params["timeout"].(float64)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    if timeout == 0 {
        timeout = 300
    }
    
    args := []string{"wait", fmt.Sprintf("%s/%s", resource, name), 
                    fmt.Sprintf("--for=condition=%s", condition), 
                    "-n", namespace, 
                    fmt.Sprintf("--timeout=%ds", int(timeout))}
    
    if kubeconfig, ok := params["kubeconfig"].(string); ok && kubeconfig != "" {
        args = append([]string{"--kubeconfig", kubeconfig}, args...)
    }
    
    cmd := exec.CommandContext(ctx, "kubectl", args...)
    output, err := cmd.CombinedOutput()
    
    ready := err == nil
    status := "ready"
    if err != nil {
        status = "timeout"
    }
    
    return map[string]interface{}{
        "ready":  ready,
        "status": status,
        "output": string(output),
    }, nil
}

var ExportedPlugin plugin.Plugin = &KubernetesPlugin{}
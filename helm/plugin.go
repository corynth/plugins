package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type HelmPlugin struct{}

func (p *HelmPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "helm",
        Version:     "1.0.0",
        Description: "Helm package manager for Kubernetes applications",
        Author:      "Corynth Team",
        Tags:        []string{"helm", "kubernetes", "package-manager", "charts", "cloud-native"},
        License:     "Apache-2.0",
    }
}

func (p *HelmPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "install",
            Description: "Install a Helm chart",
            Inputs: map[string]plugin.InputSpec{
                "chart": {
                    Type:        "string",
                    Description: "Chart name or path (local/remote)",
                    Required:    true,
                },
                "name": {
                    Type:        "string",
                    Description: "Release name",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace",
                    Required:    false,
                    Default:     "default",
                },
                "values": {
                    Type:        "object",
                    Description: "Values to override in chart",
                    Required:    false,
                },
                "values_file": {
                    Type:        "string",
                    Description: "Path to values YAML file",
                    Required:    false,
                },
                "version": {
                    Type:        "string",
                    Description: "Chart version to install",
                    Required:    false,
                },
                "repository": {
                    Type:        "string",
                    Description: "Chart repository URL",
                    Required:    false,
                },
                "create_namespace": {
                    Type:        "boolean",
                    Description: "Create namespace if it doesn't exist",
                    Required:    false,
                    Default:     false,
                },
                "wait": {
                    Type:        "boolean",
                    Description: "Wait for resources to be ready",
                    Required:    false,
                    Default:     true,
                },
                "timeout": {
                    Type:        "string",
                    Description: "Timeout duration (e.g., '5m')",
                    Required:    false,
                    Default:     "5m",
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
                    Description: "Installation status",
                },
                "revision": {
                    Type:        "number",
                    Description: "Release revision number",
                },
                "notes": {
                    Type:        "string",
                    Description: "Chart installation notes",
                },
            },
        },
        {
            Name:        "upgrade",
            Description: "Upgrade a Helm release",
            Inputs: map[string]plugin.InputSpec{
                "name": {
                    Type:        "string",
                    Description: "Release name",
                    Required:    true,
                },
                "chart": {
                    Type:        "string",
                    Description: "Chart name or path",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace",
                    Required:    false,
                    Default:     "default",
                },
                "values": {
                    Type:        "object",
                    Description: "Values to override in chart",
                    Required:    false,
                },
                "values_file": {
                    Type:        "string",
                    Description: "Path to values YAML file",
                    Required:    false,
                },
                "version": {
                    Type:        "string",
                    Description: "Chart version to upgrade to",
                    Required:    false,
                },
                "install": {
                    Type:        "boolean",
                    Description: "Install if release doesn't exist",
                    Required:    false,
                    Default:     false,
                },
                "wait": {
                    Type:        "boolean",
                    Description: "Wait for resources to be ready",
                    Required:    false,
                    Default:     true,
                },
                "timeout": {
                    Type:        "string",
                    Description: "Timeout duration (e.g., '5m')",
                    Required:    false,
                    Default:     "5m",
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
                    Description: "Upgrade status",
                },
                "revision": {
                    Type:        "number",
                    Description: "New release revision number",
                },
            },
        },
        {
            Name:        "uninstall",
            Description: "Uninstall a Helm release",
            Inputs: map[string]plugin.InputSpec{
                "name": {
                    Type:        "string",
                    Description: "Release name",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace",
                    Required:    false,
                    Default:     "default",
                },
                "keep_history": {
                    Type:        "boolean",
                    Description: "Keep release history",
                    Required:    false,
                    Default:     false,
                },
                "timeout": {
                    Type:        "string",
                    Description: "Timeout duration (e.g., '5m')",
                    Required:    false,
                    Default:     "5m",
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
                    Description: "Uninstall status",
                },
                "removed": {
                    Type:        "boolean",
                    Description: "Whether release was removed",
                },
            },
        },
        {
            Name:        "list",
            Description: "List Helm releases",
            Inputs: map[string]plugin.InputSpec{
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace (all namespaces if empty)",
                    Required:    false,
                },
                "all_namespaces": {
                    Type:        "boolean",
                    Description: "List releases across all namespaces",
                    Required:    false,
                    Default:     false,
                },
                "status": {
                    Type:        "string",
                    Description: "Filter by status (deployed, failed, pending)",
                    Required:    false,
                },
                "kubeconfig": {
                    Type:        "string",
                    Description: "Path to kubeconfig file",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "releases": {
                    Type:        "array",
                    Description: "List of releases",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of releases found",
                },
            },
        },
        {
            Name:        "status",
            Description: "Get status of a Helm release",
            Inputs: map[string]plugin.InputSpec{
                "name": {
                    Type:        "string",
                    Description: "Release name",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace",
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
                    Description: "Release status",
                },
                "revision": {
                    Type:        "number",
                    Description: "Current revision number",
                },
                "chart": {
                    Type:        "string",
                    Description: "Chart name and version",
                },
                "namespace": {
                    Type:        "string",
                    Description: "Release namespace",
                },
            },
        },
        {
            Name:        "template",
            Description: "Render chart templates locally",
            Inputs: map[string]plugin.InputSpec{
                "chart": {
                    Type:        "string",
                    Description: "Chart name or path",
                    Required:    true,
                },
                "name": {
                    Type:        "string",
                    Description: "Release name for templating",
                    Required:    true,
                },
                "namespace": {
                    Type:        "string",
                    Description: "Target namespace",
                    Required:    false,
                    Default:     "default",
                },
                "values": {
                    Type:        "object",
                    Description: "Values to use for templating",
                    Required:    false,
                },
                "values_file": {
                    Type:        "string",
                    Description: "Path to values YAML file",
                    Required:    false,
                },
                "version": {
                    Type:        "string",
                    Description: "Chart version to template",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "manifests": {
                    Type:        "string",
                    Description: "Rendered Kubernetes manifests",
                },
                "resources": {
                    Type:        "array",
                    Description: "List of resource types in manifests",
                },
            },
        },
        {
            Name:        "repo_add",
            Description: "Add a Helm repository",
            Inputs: map[string]plugin.InputSpec{
                "name": {
                    Type:        "string",
                    Description: "Repository name",
                    Required:    true,
                },
                "url": {
                    Type:        "string",
                    Description: "Repository URL",
                    Required:    true,
                },
                "username": {
                    Type:        "string",
                    Description: "Repository username",
                    Required:    false,
                },
                "password": {
                    Type:        "string",
                    Description: "Repository password",
                    Required:    false,
                },
                "force_update": {
                    Type:        "boolean",
                    Description: "Replace existing repository",
                    Required:    false,
                    Default:     false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "status": {
                    Type:        "string",
                    Description: "Repository add status",
                },
                "added": {
                    Type:        "boolean",
                    Description: "Whether repository was added",
                },
            },
        },
        {
            Name:        "repo_update",
            Description: "Update Helm repository index",
            Inputs: map[string]plugin.InputSpec{},
            Outputs: map[string]plugin.OutputSpec{
                "status": {
                    Type:        "string",
                    Description: "Update status",
                },
                "updated": {
                    Type:        "boolean",
                    Description: "Whether repositories were updated",
                },
            },
        },
    }
}

func (p *HelmPlugin) Validate(params map[string]interface{}) error {
    // Check if helm is available
    if _, err := exec.LookPath("helm"); err != nil {
        return fmt.Errorf("helm is not installed or not in PATH")
    }
    return nil
}

func (p *HelmPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "install":
        return p.executeInstall(ctx, params)
    case "upgrade":
        return p.executeUpgrade(ctx, params)
    case "uninstall":
        return p.executeUninstall(ctx, params)
    case "list":
        return p.executeList(ctx, params)
    case "status":
        return p.executeStatus(ctx, params)
    case "template":
        return p.executeTemplate(ctx, params)
    case "repo_add":
        return p.executeRepoAdd(ctx, params)
    case "repo_update":
        return p.executeRepoUpdate(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *HelmPlugin) buildBaseArgs(params map[string]interface{}) []string {
    var args []string
    
    if kubeconfig, ok := params["kubeconfig"].(string); ok && kubeconfig != "" {
        args = append(args, "--kubeconfig", kubeconfig)
    }
    
    return args
}

func (p *HelmPlugin) executeInstall(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    chart, _ := params["chart"].(string)
    name, _ := params["name"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := p.buildBaseArgs(params)
    args = append(args, "install", name, chart, "--namespace", namespace)
    
    // Handle optional parameters
    if version, ok := params["version"].(string); ok && version != "" {
        args = append(args, "--version", version)
    }
    
    if repo, ok := params["repository"].(string); ok && repo != "" {
        args = append(args, "--repo", repo)
    }
    
    if createNs, ok := params["create_namespace"].(bool); ok && createNs {
        args = append(args, "--create-namespace")
    }
    
    if wait, ok := params["wait"].(bool); ok && wait {
        args = append(args, "--wait")
    }
    
    if timeout, ok := params["timeout"].(string); ok && timeout != "" {
        args = append(args, "--timeout", timeout)
    }
    
    // Handle values
    if values, ok := params["values"].(map[string]interface{}); ok {
        for key, value := range values {
            args = append(args, "--set", fmt.Sprintf("%s=%v", key, value))
        }
    }
    
    if valuesFile, ok := params["values_file"].(string); ok && valuesFile != "" {
        args = append(args, "--values", valuesFile)
    }
    
    args = append(args, "--output", "json")
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "status": "failed",
            "error":  string(output),
        }, nil
    }
    
    // Parse JSON output
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        // If JSON parsing fails, return raw output
        return map[string]interface{}{
            "status": "success",
            "output": string(output),
        }, nil
    }
    
    revision := 1
    if info, ok := result["info"].(map[string]interface{}); ok {
        if rev, ok := info["revision"].(float64); ok {
            revision = int(rev)
        }
    }
    
    notes := ""
    if info, ok := result["info"].(map[string]interface{}); ok {
        if n, ok := info["notes"].(string); ok {
            notes = n
        }
    }
    
    return map[string]interface{}{
        "status":   "success",
        "revision": revision,
        "notes":    notes,
        "output":   string(output),
    }, nil
}

func (p *HelmPlugin) executeUpgrade(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    name, _ := params["name"].(string)
    chart, _ := params["chart"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := p.buildBaseArgs(params)
    args = append(args, "upgrade", name, chart, "--namespace", namespace)
    
    if install, ok := params["install"].(bool); ok && install {
        args = append(args, "--install")
    }
    
    if version, ok := params["version"].(string); ok && version != "" {
        args = append(args, "--version", version)
    }
    
    if wait, ok := params["wait"].(bool); ok && wait {
        args = append(args, "--wait")
    }
    
    if timeout, ok := params["timeout"].(string); ok && timeout != "" {
        args = append(args, "--timeout", timeout)
    }
    
    // Handle values
    if values, ok := params["values"].(map[string]interface{}); ok {
        for key, value := range values {
            args = append(args, "--set", fmt.Sprintf("%s=%v", key, value))
        }
    }
    
    if valuesFile, ok := params["values_file"].(string); ok && valuesFile != "" {
        args = append(args, "--values", valuesFile)
    }
    
    args = append(args, "--output", "json")
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "status": "failed",
            "error":  string(output),
        }, nil
    }
    
    // Parse JSON output for revision
    var result map[string]interface{}
    revision := 1
    if err := json.Unmarshal(output, &result); err == nil {
        if info, ok := result["info"].(map[string]interface{}); ok {
            if rev, ok := info["revision"].(float64); ok {
                revision = int(rev)
            }
        }
    }
    
    return map[string]interface{}{
        "status":   "success",
        "revision": revision,
        "output":   string(output),
    }, nil
}

func (p *HelmPlugin) executeUninstall(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    name, _ := params["name"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := p.buildBaseArgs(params)
    args = append(args, "uninstall", name, "--namespace", namespace)
    
    if keepHistory, ok := params["keep_history"].(bool); ok && keepHistory {
        args = append(args, "--keep-history")
    }
    
    if timeout, ok := params["timeout"].(string); ok && timeout != "" {
        args = append(args, "--timeout", timeout)
    }
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.CombinedOutput()
    
    removed := err == nil
    status := "success"
    if err != nil {
        status = "failed"
    }
    
    return map[string]interface{}{
        "status":  status,
        "removed": removed,
        "output":  string(output),
    }, nil
}

func (p *HelmPlugin) executeList(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := p.buildBaseArgs(params)
    args = append(args, "list", "--output", "json")
    
    if namespace, ok := params["namespace"].(string); ok && namespace != "" {
        args = append(args, "--namespace", namespace)
    }
    
    if allNs, ok := params["all_namespaces"].(bool); ok && allNs {
        args = append(args, "--all-namespaces")
    }
    
    if status, ok := params["status"].(string); ok && status != "" {
        args = append(args, "--filter", status)
    }
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "releases": []interface{}{},
            "count":    0,
            "error":    err.Error(),
        }, nil
    }
    
    var releases []interface{}
    if err := json.Unmarshal(output, &releases); err != nil {
        return nil, fmt.Errorf("failed to parse helm list output: %w", err)
    }
    
    return map[string]interface{}{
        "releases": releases,
        "count":    len(releases),
    }, nil
}

func (p *HelmPlugin) executeStatus(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    name, _ := params["name"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := p.buildBaseArgs(params)
    args = append(args, "status", name, "--namespace", namespace, "--output", "json")
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "status":    "not found",
            "error":     err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse helm status output: %w", err)
    }
    
    status := "unknown"
    revision := 0
    chart := ""
    ns := namespace
    
    if info, ok := result["info"].(map[string]interface{}); ok {
        if s, ok := info["status"].(string); ok {
            status = s
        }
        if r, ok := info["revision"].(float64); ok {
            revision = int(r)
        }
    }
    
    if chartInfo, ok := result["chart"].(map[string]interface{}); ok {
        if metadata, ok := chartInfo["metadata"].(map[string]interface{}); ok {
            if name, ok := metadata["name"].(string); ok {
                if version, ok := metadata["version"].(string); ok {
                    chart = fmt.Sprintf("%s-%s", name, version)
                } else {
                    chart = name
                }
            }
        }
    }
    
    if resultNs, ok := result["namespace"].(string); ok {
        ns = resultNs
    }
    
    return map[string]interface{}{
        "status":    status,
        "revision":  revision,
        "chart":     chart,
        "namespace": ns,
    }, nil
}

func (p *HelmPlugin) executeTemplate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    chart, _ := params["chart"].(string)
    name, _ := params["name"].(string)
    namespace, _ := params["namespace"].(string)
    if namespace == "" {
        namespace = "default"
    }
    
    args := p.buildBaseArgs(params)
    args = append(args, "template", name, chart, "--namespace", namespace)
    
    if version, ok := params["version"].(string); ok && version != "" {
        args = append(args, "--version", version)
    }
    
    // Handle values
    if values, ok := params["values"].(map[string]interface{}); ok {
        for key, value := range values {
            args = append(args, "--set", fmt.Sprintf("%s=%v", key, value))
        }
    }
    
    if valuesFile, ok := params["values_file"].(string); ok && valuesFile != "" {
        args = append(args, "--values", valuesFile)
    }
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "manifests": "",
            "resources": []string{},
            "error":     err.Error(),
        }, nil
    }
    
    manifests := string(output)
    
    // Extract resource types from manifests
    resources := []string{}
    lines := strings.Split(manifests, "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "kind:") {
            kind := strings.TrimSpace(strings.TrimPrefix(line, "kind:"))
            if kind != "" {
                found := false
                for _, existing := range resources {
                    if existing == kind {
                        found = true
                        break
                    }
                }
                if !found {
                    resources = append(resources, kind)
                }
            }
        }
    }
    
    return map[string]interface{}{
        "manifests": manifests,
        "resources": resources,
    }, nil
}

func (p *HelmPlugin) executeRepoAdd(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    name, _ := params["name"].(string)
    url, _ := params["url"].(string)
    
    args := p.buildBaseArgs(params)
    args = append(args, "repo", "add", name, url)
    
    if username, ok := params["username"].(string); ok && username != "" {
        args = append(args, "--username", username)
    }
    
    if password, ok := params["password"].(string); ok && password != "" {
        args = append(args, "--password", password)
    }
    
    if forceUpdate, ok := params["force_update"].(bool); ok && forceUpdate {
        args = append(args, "--force-update")
    }
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.CombinedOutput()
    
    added := err == nil
    status := "success"
    if err != nil {
        status = "failed"
    }
    
    return map[string]interface{}{
        "status": status,
        "added":  added,
        "output": string(output),
    }, nil
}

func (p *HelmPlugin) executeRepoUpdate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := p.buildBaseArgs(params)
    args = append(args, "repo", "update")
    
    cmd := exec.CommandContext(ctx, "helm", args...)
    output, err := cmd.CombinedOutput()
    
    updated := err == nil
    status := "success"
    if err != nil {
        status = "failed"
    }
    
    return map[string]interface{}{
        "status":  status,
        "updated": updated,
        "output":  string(output),
    }, nil
}

var ExportedPlugin plugin.Plugin = &HelmPlugin{}
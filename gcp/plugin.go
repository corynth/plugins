package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
    "strings"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type GCPPlugin struct{}

func (p *GCPPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "gcp",
        Version:     "1.0.0",
        Description: "Google Cloud Platform operations and resource management",
        Author:      "Corynth Team",
        Tags:        []string{"gcp", "google-cloud", "gce", "gke", "gcs", "cloud-functions", "cloud-native"},
        License:     "Apache-2.0",
    }
}

func (p *GCPPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "compute_list",
            Description: "List Compute Engine instances",
            Inputs: map[string]plugin.InputSpec{
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
                "zone": {
                    Type:        "string",
                    Description: "GCP zone (e.g., us-central1-a)",
                    Required:    false,
                },
                "filter": {
                    Type:        "string",
                    Description: "Filter expression for instances",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "instances": {
                    Type:        "array",
                    Description: "List of compute instances",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of instances found",
                },
            },
        },
        {
            Name:        "compute_create",
            Description: "Create Compute Engine instance",
            Inputs: map[string]plugin.InputSpec{
                "name": {
                    Type:        "string",
                    Description: "Instance name",
                    Required:    true,
                },
                "machine_type": {
                    Type:        "string",
                    Description: "Machine type (e.g., e2-medium, n1-standard-1)",
                    Required:    true,
                },
                "image": {
                    Type:        "string",
                    Description: "Boot disk image",
                    Required:    true,
                },
                "zone": {
                    Type:        "string",
                    Description: "GCP zone",
                    Required:    true,
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
                "network": {
                    Type:        "string",
                    Description: "VPC network name",
                    Required:    false,
                    Default:     "default",
                },
                "subnet": {
                    Type:        "string",
                    Description: "Subnet name",
                    Required:    false,
                },
                "labels": {
                    Type:        "object",
                    Description: "Instance labels",
                    Required:    false,
                },
                "metadata": {
                    Type:        "object",
                    Description: "Instance metadata",
                    Required:    false,
                },
                "startup_script": {
                    Type:        "string",
                    Description: "Startup script content",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "instance_name": {
                    Type:        "string",
                    Description: "Created instance name",
                },
                "status": {
                    Type:        "string",
                    Description: "Instance status",
                },
                "internal_ip": {
                    Type:        "string",
                    Description: "Internal IP address",
                },
                "external_ip": {
                    Type:        "string",
                    Description: "External IP address",
                },
            },
        },
        {
            Name:        "compute_delete",
            Description: "Delete Compute Engine instances",
            Inputs: map[string]plugin.InputSpec{
                "names": {
                    Type:        "array",
                    Description: "List of instance names to delete",
                    Required:    true,
                },
                "zone": {
                    Type:        "string",
                    Description: "GCP zone",
                    Required:    true,
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "deleted": {
                    Type:        "array",
                    Description: "List of deleted instance names",
                },
                "status": {
                    Type:        "string",
                    Description: "Deletion status",
                },
            },
        },
        {
            Name:        "gke_list",
            Description: "List GKE clusters",
            Inputs: map[string]plugin.InputSpec{
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
                "location": {
                    Type:        "string",
                    Description: "GCP location (region or zone)",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "clusters": {
                    Type:        "array",
                    Description: "List of GKE clusters",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of clusters found",
                },
            },
        },
        {
            Name:        "gke_get_credentials",
            Description: "Get GKE cluster credentials",
            Inputs: map[string]plugin.InputSpec{
                "cluster": {
                    Type:        "string",
                    Description: "Cluster name",
                    Required:    true,
                },
                "location": {
                    Type:        "string",
                    Description: "Cluster location (region or zone)",
                    Required:    true,
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "configured": {
                    Type:        "boolean",
                    Description: "Whether credentials were configured",
                },
                "context": {
                    Type:        "string",
                    Description: "Kubernetes context name",
                },
            },
        },
        {
            Name:        "storage_list",
            Description: "List Cloud Storage buckets or objects",
            Inputs: map[string]plugin.InputSpec{
                "bucket": {
                    Type:        "string",
                    Description: "Bucket name (if listing objects)",
                    Required:    false,
                },
                "prefix": {
                    Type:        "string",
                    Description: "Object prefix filter",
                    Required:    false,
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "buckets": {
                    Type:        "array",
                    Description: "List of buckets (if no bucket specified)",
                },
                "objects": {
                    Type:        "array",
                    Description: "List of objects (if bucket specified)",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of items found",
                },
            },
        },
        {
            Name:        "storage_upload",
            Description: "Upload file to Cloud Storage",
            Inputs: map[string]plugin.InputSpec{
                "file_path": {
                    Type:        "string",
                    Description: "Local file path",
                    Required:    true,
                },
                "bucket": {
                    Type:        "string",
                    Description: "GCS bucket name",
                    Required:    true,
                },
                "object_name": {
                    Type:        "string",
                    Description: "Object name in bucket",
                    Required:    true,
                },
                "content_type": {
                    Type:        "string",
                    Description: "Content type",
                    Required:    false,
                },
                "metadata": {
                    Type:        "object",
                    Description: "Object metadata",
                    Required:    false,
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "url": {
                    Type:        "string",
                    Description: "Object GCS URL",
                },
                "size": {
                    Type:        "number",
                    Description: "Object size in bytes",
                },
                "md5": {
                    Type:        "string",
                    Description: "Object MD5 hash",
                },
            },
        },
        {
            Name:        "storage_download",
            Description: "Download file from Cloud Storage",
            Inputs: map[string]plugin.InputSpec{
                "bucket": {
                    Type:        "string",
                    Description: "GCS bucket name",
                    Required:    true,
                },
                "object_name": {
                    Type:        "string",
                    Description: "Object name in bucket",
                    Required:    true,
                },
                "file_path": {
                    Type:        "string",
                    Description: "Local file path to save",
                    Required:    true,
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "size": {
                    Type:        "number",
                    Description: "Downloaded file size",
                },
                "md5": {
                    Type:        "string",
                    Description: "File MD5 hash",
                },
            },
        },
        {
            Name:        "functions_deploy",
            Description: "Deploy Cloud Function",
            Inputs: map[string]plugin.InputSpec{
                "name": {
                    Type:        "string",
                    Description: "Function name",
                    Required:    true,
                },
                "source": {
                    Type:        "string",
                    Description: "Source code directory",
                    Required:    true,
                },
                "entry_point": {
                    Type:        "string",
                    Description: "Function entry point",
                    Required:    true,
                },
                "runtime": {
                    Type:        "string",
                    Description: "Runtime (e.g., nodejs18, python39, go119)",
                    Required:    true,
                },
                "trigger": {
                    Type:        "string",
                    Description: "Trigger type (http, pubsub, storage)",
                    Required:    true,
                },
                "region": {
                    Type:        "string",
                    Description: "Deployment region",
                    Required:    false,
                    Default:     "us-central1",
                },
                "memory": {
                    Type:        "string",
                    Description: "Memory allocation (e.g., 256MB, 1GB)",
                    Required:    false,
                    Default:     "256MB",
                },
                "timeout": {
                    Type:        "string",
                    Description: "Timeout duration (e.g., 60s)",
                    Required:    false,
                    Default:     "60s",
                },
                "env_vars": {
                    Type:        "object",
                    Description: "Environment variables",
                    Required:    false,
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "url": {
                    Type:        "string",
                    Description: "Function trigger URL",
                },
                "status": {
                    Type:        "string",
                    Description: "Deployment status",
                },
                "version": {
                    Type:        "string",
                    Description: "Function version",
                },
            },
        },
        {
            Name:        "functions_invoke",
            Description: "Invoke Cloud Function",
            Inputs: map[string]plugin.InputSpec{
                "name": {
                    Type:        "string",
                    Description: "Function name",
                    Required:    true,
                },
                "data": {
                    Type:        "object",
                    Description: "Function input data",
                    Required:    false,
                },
                "region": {
                    Type:        "string",
                    Description: "Function region",
                    Required:    false,
                    Default:     "us-central1",
                },
                "project": {
                    Type:        "string",
                    Description: "GCP project ID",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "response": {
                    Type:        "object",
                    Description: "Function response",
                },
                "execution_id": {
                    Type:        "string",
                    Description: "Execution ID",
                },
                "duration": {
                    Type:        "string",
                    Description: "Execution duration",
                },
            },
        },
    }
}

func (p *GCPPlugin) Validate(params map[string]interface{}) error {
    // Check if gcloud is available
    if _, err := exec.LookPath("gcloud"); err != nil {
        return fmt.Errorf("gcloud CLI is not installed or not in PATH")
    }
    return nil
}

func (p *GCPPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "compute_list":
        return p.executeComputeList(ctx, params)
    case "compute_create":
        return p.executeComputeCreate(ctx, params)
    case "compute_delete":
        return p.executeComputeDelete(ctx, params)
    case "gke_list":
        return p.executeGKEList(ctx, params)
    case "gke_get_credentials":
        return p.executeGKEGetCredentials(ctx, params)
    case "storage_list":
        return p.executeStorageList(ctx, params)
    case "storage_upload":
        return p.executeStorageUpload(ctx, params)
    case "storage_download":
        return p.executeStorageDownload(ctx, params)
    case "functions_deploy":
        return p.executeFunctionsDeploy(ctx, params)
    case "functions_invoke":
        return p.executeFunctionsInvoke(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *GCPPlugin) buildBaseArgs(params map[string]interface{}) []string {
    var args []string
    
    if project, ok := params["project"].(string); ok && project != "" {
        args = append(args, "--project", project)
    }
    
    args = append(args, "--format", "json")
    return args
}

func (p *GCPPlugin) executeComputeList(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := []string{"compute", "instances", "list"}
    args = append(args, p.buildBaseArgs(params)...)
    
    if zone, ok := params["zone"].(string); ok && zone != "" {
        args = append(args, "--zones", zone)
    }
    
    if filter, ok := params["filter"].(string); ok && filter != "" {
        args = append(args, "--filter", filter)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "instances": []interface{}{},
            "count":     0,
            "error":     err.Error(),
        }, nil
    }
    
    var instances []interface{}
    if err := json.Unmarshal(output, &instances); err != nil {
        return nil, fmt.Errorf("failed to parse gcloud output: %w", err)
    }
    
    return map[string]interface{}{
        "instances": instances,
        "count":     len(instances),
    }, nil
}

func (p *GCPPlugin) executeComputeCreate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    name, _ := params["name"].(string)
    machineType, _ := params["machine_type"].(string)
    image, _ := params["image"].(string)
    zone, _ := params["zone"].(string)
    
    args := []string{"compute", "instances", "create", name,
                     "--machine-type", machineType,
                     "--image", image,
                     "--zone", zone}
    
    if network, ok := params["network"].(string); ok && network != "" {
        args = append(args, "--network", network)
    }
    
    if subnet, ok := params["subnet"].(string); ok && subnet != "" {
        args = append(args, "--subnet", subnet)
    }
    
    if labels, ok := params["labels"].(map[string]interface{}); ok {
        labelArgs := []string{}
        for k, v := range labels {
            labelArgs = append(labelArgs, fmt.Sprintf("%s=%v", k, v))
        }
        if len(labelArgs) > 0 {
            args = append(args, "--labels", strings.Join(labelArgs, ","))
        }
    }
    
    if metadata, ok := params["metadata"].(map[string]interface{}); ok {
        for k, v := range metadata {
            args = append(args, "--metadata", fmt.Sprintf("%s=%v", k, v))
        }
    }
    
    if startupScript, ok := params["startup_script"].(string); ok && startupScript != "" {
        args = append(args, "--metadata", fmt.Sprintf("startup-script=%s", startupScript))
    }
    
    args = append(args, p.buildBaseArgs(params)...)
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "instance_name": "",
            "status":        "failed",
            "error":         err.Error(),
        }, nil
    }
    
    var result []map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse gcloud output: %w", err)
    }
    
    instanceName := name
    status := "RUNNING"
    internalIP := ""
    externalIP := ""
    
    if len(result) > 0 {
        inst := result[0]
        if s, ok := inst["status"].(string); ok {
            status = s
        }
        if networkInterfaces, ok := inst["networkInterfaces"].([]interface{}); ok && len(networkInterfaces) > 0 {
            if ni, ok := networkInterfaces[0].(map[string]interface{}); ok {
                if ip, ok := ni["networkIP"].(string); ok {
                    internalIP = ip
                }
                if accessConfigs, ok := ni["accessConfigs"].([]interface{}); ok && len(accessConfigs) > 0 {
                    if ac, ok := accessConfigs[0].(map[string]interface{}); ok {
                        if ip, ok := ac["natIP"].(string); ok {
                            externalIP = ip
                        }
                    }
                }
            }
        }
    }
    
    return map[string]interface{}{
        "instance_name": instanceName,
        "status":        status,
        "internal_ip":   internalIP,
        "external_ip":   externalIP,
    }, nil
}

func (p *GCPPlugin) executeComputeDelete(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    names, ok := params["names"].([]interface{})
    if !ok || len(names) == 0 {
        return nil, fmt.Errorf("names parameter is required")
    }
    
    zone, _ := params["zone"].(string)
    
    nameList := []string{}
    for _, n := range names {
        if nameStr, ok := n.(string); ok {
            nameList = append(nameList, nameStr)
        }
    }
    
    args := []string{"compute", "instances", "delete"}
    args = append(args, nameList...)
    args = append(args, "--zone", zone, "--quiet")
    args = append(args, p.buildBaseArgs(params)...)
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "deleted": []string{},
            "status":  "failed",
            "error":   string(output),
        }, nil
    }
    
    return map[string]interface{}{
        "deleted": nameList,
        "status":  "success",
        "output":  string(output),
    }, nil
}

func (p *GCPPlugin) executeGKEList(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := []string{"container", "clusters", "list"}
    args = append(args, p.buildBaseArgs(params)...)
    
    if location, ok := params["location"].(string); ok && location != "" {
        args = append(args, "--location", location)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "clusters": []interface{}{},
            "count":    0,
            "error":    err.Error(),
        }, nil
    }
    
    var clusters []interface{}
    if err := json.Unmarshal(output, &clusters); err != nil {
        return nil, fmt.Errorf("failed to parse gcloud output: %w", err)
    }
    
    return map[string]interface{}{
        "clusters": clusters,
        "count":    len(clusters),
    }, nil
}

func (p *GCPPlugin) executeGKEGetCredentials(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    cluster, _ := params["cluster"].(string)
    location, _ := params["location"].(string)
    
    args := []string{"container", "clusters", "get-credentials", cluster,
                     "--location", location}
    
    if project, ok := params["project"].(string); ok && project != "" {
        args = append(args, "--project", project)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.CombinedOutput()
    
    configured := err == nil
    context := fmt.Sprintf("gke_%s_%s_%s", params["project"], location, cluster)
    
    return map[string]interface{}{
        "configured": configured,
        "context":    context,
        "output":     string(output),
    }, nil
}

func (p *GCPPlugin) executeStorageList(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    bucket, hasBucket := params["bucket"].(string)
    
    var args []string
    if hasBucket && bucket != "" {
        // List objects in bucket
        args = []string{"storage", "ls", fmt.Sprintf("gs://%s", bucket)}
        
        if prefix, ok := params["prefix"].(string); ok && prefix != "" {
            args[2] = fmt.Sprintf("gs://%s/%s", bucket, prefix)
        }
    } else {
        // List buckets
        args = []string{"storage", "buckets", "list"}
        args = append(args, p.buildBaseArgs(params)...)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "buckets": []interface{}{},
            "objects": []interface{}{},
            "count":   0,
            "error":   err.Error(),
        }, nil
    }
    
    if hasBucket && bucket != "" {
        // Parse object list (plain text output)
        lines := strings.Split(string(output), "\n")
        objects := []interface{}{}
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line != "" {
                objects = append(objects, map[string]interface{}{
                    "name": strings.TrimPrefix(line, fmt.Sprintf("gs://%s/", bucket)),
                    "url":  line,
                })
            }
        }
        
        return map[string]interface{}{
            "objects": objects,
            "count":   len(objects),
        }, nil
    } else {
        // Parse bucket list (JSON output)
        var buckets []interface{}
        if err := json.Unmarshal(output, &buckets); err != nil {
            return nil, fmt.Errorf("failed to parse gcloud output: %w", err)
        }
        
        return map[string]interface{}{
            "buckets": buckets,
            "count":   len(buckets),
        }, nil
    }
}

func (p *GCPPlugin) executeStorageUpload(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    filePath, _ := params["file_path"].(string)
    bucket, _ := params["bucket"].(string)
    objectName, _ := params["object_name"].(string)
    
    gcsPath := fmt.Sprintf("gs://%s/%s", bucket, objectName)
    args := []string{"storage", "cp", filePath, gcsPath}
    
    if contentType, ok := params["content_type"].(string); ok && contentType != "" {
        args = append(args, "--content-type", contentType)
    }
    
    if metadata, ok := params["metadata"].(map[string]interface{}); ok {
        for k, v := range metadata {
            args = append(args, "--metadata", fmt.Sprintf("%s=%v", k, v))
        }
    }
    
    if project, ok := params["project"].(string); ok && project != "" {
        args = append(args, "--project", project)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "url":   "",
            "size":  0,
            "md5":   "",
            "error": string(output),
        }, nil
    }
    
    return map[string]interface{}{
        "url":    gcsPath,
        "size":   0, // Would need additional call to get size
        "md5":    "", // Would need additional call to get hash
        "output": string(output),
    }, nil
}

func (p *GCPPlugin) executeStorageDownload(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    bucket, _ := params["bucket"].(string)
    objectName, _ := params["object_name"].(string)
    filePath, _ := params["file_path"].(string)
    
    gcsPath := fmt.Sprintf("gs://%s/%s", bucket, objectName)
    args := []string{"storage", "cp", gcsPath, filePath}
    
    if project, ok := params["project"].(string); ok && project != "" {
        args = append(args, "--project", project)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "size":  0,
            "md5":   "",
            "error": string(output),
        }, nil
    }
    
    return map[string]interface{}{
        "size":   0, // Would need additional call to get size
        "md5":    "", // Would need additional call to get hash
        "output": string(output),
    }, nil
}

func (p *GCPPlugin) executeFunctionsDeploy(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    name, _ := params["name"].(string)
    source, _ := params["source"].(string)
    entryPoint, _ := params["entry_point"].(string)
    runtime, _ := params["runtime"].(string)
    trigger, _ := params["trigger"].(string)
    region, _ := params["region"].(string)
    if region == "" {
        region = "us-central1"
    }
    
    args := []string{"functions", "deploy", name,
                     "--source", source,
                     "--entry-point", entryPoint,
                     "--runtime", runtime,
                     "--region", region}
    
    // Add trigger configuration
    switch trigger {
    case "http":
        args = append(args, "--trigger-http", "--allow-unauthenticated")
    case "pubsub":
        if topic, ok := params["topic"].(string); ok {
            args = append(args, "--trigger-topic", topic)
        }
    case "storage":
        if bucket, ok := params["bucket"].(string); ok {
            args = append(args, "--trigger-bucket", bucket)
        }
    }
    
    if memory, ok := params["memory"].(string); ok && memory != "" {
        args = append(args, "--memory", memory)
    }
    
    if timeout, ok := params["timeout"].(string); ok && timeout != "" {
        args = append(args, "--timeout", timeout)
    }
    
    if envVars, ok := params["env_vars"].(map[string]interface{}); ok {
        envArgs := []string{}
        for k, v := range envVars {
            envArgs = append(envArgs, fmt.Sprintf("%s=%v", k, v))
        }
        if len(envArgs) > 0 {
            args = append(args, "--set-env-vars", strings.Join(envArgs, ","))
        }
    }
    
    if project, ok := params["project"].(string); ok && project != "" {
        args = append(args, "--project", project)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "url":     "",
            "status":  "failed",
            "version": "",
            "error":   string(output),
        }, nil
    }
    
    url := ""
    if trigger == "http" {
        url = fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", region, params["project"], name)
    }
    
    return map[string]interface{}{
        "url":     url,
        "status":  "deployed",
        "version": "latest",
        "output":  string(output),
    }, nil
}

func (p *GCPPlugin) executeFunctionsInvoke(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    name, _ := params["name"].(string)
    region, _ := params["region"].(string)
    if region == "" {
        region = "us-central1"
    }
    
    args := []string{"functions", "call", name, "--region", region}
    
    if data, ok := params["data"].(map[string]interface{}); ok {
        dataJSON, _ := json.Marshal(data)
        args = append(args, "--data", string(dataJSON))
    }
    
    if project, ok := params["project"].(string); ok && project != "" {
        args = append(args, "--project", project)
    }
    
    cmd := exec.CommandContext(ctx, "gcloud", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "response":     nil,
            "execution_id": "",
            "duration":     "",
            "error":        err.Error(),
        }, nil
    }
    
    // Parse response (usually includes executionId and result)
    var response map[string]interface{}
    if err := json.Unmarshal(output, &response); err != nil {
        response = map[string]interface{}{
            "raw_output": string(output),
        }
    }
    
    executionID := ""
    if id, ok := response["executionId"].(string); ok {
        executionID = id
    }
    
    return map[string]interface{}{
        "response":     response,
        "execution_id": executionID,
        "duration":     "N/A",
    }, nil
}

var ExportedPlugin plugin.Plugin = &GCPPlugin{}
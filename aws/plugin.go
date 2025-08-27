package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
    "strings"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type AWSPlugin struct{}

func (p *AWSPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "aws",
        Version:     "1.0.0",
        Description: "Amazon Web Services cloud operations and resource management",
        Author:      "Corynth Team",
        Tags:        []string{"aws", "cloud", "ec2", "s3", "lambda", "vpc", "iam", "cloud-native"},
        License:     "Apache-2.0",
    }
}

func (p *AWSPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "ec2_list",
            Description: "List EC2 instances",
            Inputs: map[string]plugin.InputSpec{
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                    Default:     "us-east-1",
                },
                "state": {
                    Type:        "string", 
                    Description: "Instance state filter (running, stopped, etc.)",
                    Required:    false,
                },
                "tags": {
                    Type:        "object",
                    Description: "Tag filters (key-value pairs)",
                    Required:    false,
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "instances": {
                    Type:        "array",
                    Description: "List of EC2 instances",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of instances found",
                },
            },
        },
        {
            Name:        "ec2_create",
            Description: "Create EC2 instance",
            Inputs: map[string]plugin.InputSpec{
                "image_id": {
                    Type:        "string",
                    Description: "AMI ID",
                    Required:    true,
                },
                "instance_type": {
                    Type:        "string",
                    Description: "Instance type (t2.micro, m5.large, etc.)",
                    Required:    true,
                },
                "key_name": {
                    Type:        "string",
                    Description: "Key pair name",
                    Required:    false,
                },
                "security_groups": {
                    Type:        "array",
                    Description: "Security group IDs",
                    Required:    false,
                },
                "subnet_id": {
                    Type:        "string",
                    Description: "Subnet ID",
                    Required:    false,
                },
                "user_data": {
                    Type:        "string",
                    Description: "User data script",
                    Required:    false,
                },
                "tags": {
                    Type:        "object",
                    Description: "Instance tags",
                    Required:    false,
                },
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                    Default:     "us-east-1",
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "instance_id": {
                    Type:        "string",
                    Description: "Created instance ID",
                },
                "state": {
                    Type:        "string",
                    Description: "Instance state",
                },
                "public_ip": {
                    Type:        "string",
                    Description: "Public IP address",
                },
                "private_ip": {
                    Type:        "string",
                    Description: "Private IP address",
                },
            },
        },
        {
            Name:        "ec2_terminate",
            Description: "Terminate EC2 instances",
            Inputs: map[string]plugin.InputSpec{
                "instance_ids": {
                    Type:        "array",
                    Description: "List of instance IDs to terminate",
                    Required:    true,
                },
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                    Default:     "us-east-1",
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "terminated": {
                    Type:        "array",
                    Description: "List of terminated instance IDs",
                },
                "status": {
                    Type:        "string",
                    Description: "Termination status",
                },
            },
        },
        {
            Name:        "s3_list",
            Description: "List S3 buckets or objects",
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
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
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
            Name:        "s3_upload",
            Description: "Upload file to S3",
            Inputs: map[string]plugin.InputSpec{
                "file_path": {
                    Type:        "string",
                    Description: "Local file path",
                    Required:    true,
                },
                "bucket": {
                    Type:        "string",
                    Description: "S3 bucket name",
                    Required:    true,
                },
                "key": {
                    Type:        "string",
                    Description: "S3 object key",
                    Required:    true,
                },
                "content_type": {
                    Type:        "string",
                    Description: "Content type",
                    Required:    false,
                },
                "acl": {
                    Type:        "string",
                    Description: "Access control list (private, public-read, etc.)",
                    Required:    false,
                    Default:     "private",
                },
                "metadata": {
                    Type:        "object",
                    Description: "Object metadata",
                    Required:    false,
                },
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "etag": {
                    Type:        "string",
                    Description: "Object ETag",
                },
                "url": {
                    Type:        "string",
                    Description: "Object URL",
                },
                "size": {
                    Type:        "number",
                    Description: "Object size in bytes",
                },
            },
        },
        {
            Name:        "s3_download",
            Description: "Download file from S3",
            Inputs: map[string]plugin.InputSpec{
                "bucket": {
                    Type:        "string",
                    Description: "S3 bucket name",
                    Required:    true,
                },
                "key": {
                    Type:        "string",
                    Description: "S3 object key",
                    Required:    true,
                },
                "file_path": {
                    Type:        "string",
                    Description: "Local file path to save",
                    Required:    true,
                },
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "size": {
                    Type:        "number",
                    Description: "Downloaded file size",
                },
                "etag": {
                    Type:        "string",
                    Description: "Object ETag",
                },
                "last_modified": {
                    Type:        "string",
                    Description: "Last modified timestamp",
                },
            },
        },
        {
            Name:        "lambda_invoke",
            Description: "Invoke Lambda function",
            Inputs: map[string]plugin.InputSpec{
                "function_name": {
                    Type:        "string",
                    Description: "Lambda function name or ARN",
                    Required:    true,
                },
                "payload": {
                    Type:        "object",
                    Description: "Function payload",
                    Required:    false,
                },
                "invocation_type": {
                    Type:        "string",
                    Description: "Invocation type (RequestResponse, Event, DryRun)",
                    Required:    false,
                    Default:     "RequestResponse",
                },
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                    Default:     "us-east-1",
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "response": {
                    Type:        "object",
                    Description: "Function response",
                },
                "status_code": {
                    Type:        "number",
                    Description: "HTTP status code",
                },
                "log_result": {
                    Type:        "string",
                    Description: "Function logs",
                },
            },
        },
        {
            Name:        "lambda_list",
            Description: "List Lambda functions",
            Inputs: map[string]plugin.InputSpec{
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                    Default:     "us-east-1",
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "functions": {
                    Type:        "array",
                    Description: "List of Lambda functions",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of functions found",
                },
            },
        },
        {
            Name:        "iam_list_users",
            Description: "List IAM users",
            Inputs: map[string]plugin.InputSpec{
                "path_prefix": {
                    Type:        "string",
                    Description: "Path prefix filter",
                    Required:    false,
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "users": {
                    Type:        "array",
                    Description: "List of IAM users",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of users found",
                },
            },
        },
        {
            Name:        "vpc_list",
            Description: "List VPCs",
            Inputs: map[string]plugin.InputSpec{
                "region": {
                    Type:        "string",
                    Description: "AWS region",
                    Required:    false,
                    Default:     "us-east-1",
                },
                "filters": {
                    Type:        "object",
                    Description: "VPC filters",
                    Required:    false,
                },
                "profile": {
                    Type:        "string",
                    Description: "AWS profile name",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "vpcs": {
                    Type:        "array",
                    Description: "List of VPCs",
                },
                "count": {
                    Type:        "number",
                    Description: "Number of VPCs found",
                },
            },
        },
    }
}

func (p *AWSPlugin) Validate(params map[string]interface{}) error {
    // Check if AWS CLI is available
    if _, err := exec.LookPath("aws"); err != nil {
        return fmt.Errorf("aws CLI is not installed or not in PATH")
    }
    return nil
}

func (p *AWSPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "ec2_list":
        return p.executeEC2List(ctx, params)
    case "ec2_create":
        return p.executeEC2Create(ctx, params)
    case "ec2_terminate":
        return p.executeEC2Terminate(ctx, params)
    case "s3_list":
        return p.executeS3List(ctx, params)
    case "s3_upload":
        return p.executeS3Upload(ctx, params)
    case "s3_download":
        return p.executeS3Download(ctx, params)
    case "lambda_invoke":
        return p.executeLambdaInvoke(ctx, params)
    case "lambda_list":
        return p.executeLambdaList(ctx, params)
    case "iam_list_users":
        return p.executeIAMListUsers(ctx, params)
    case "vpc_list":
        return p.executeVPCList(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *AWSPlugin) buildBaseArgs(params map[string]interface{}) []string {
    var args []string
    
    if region, ok := params["region"].(string); ok && region != "" {
        args = append(args, "--region", region)
    }
    
    if profile, ok := params["profile"].(string); ok && profile != "" {
        args = append(args, "--profile", profile)
    }
    
    return args
}

func (p *AWSPlugin) executeEC2List(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := p.buildBaseArgs(params)
    args = append(args, "ec2", "describe-instances", "--output", "json")
    
    // Add state filter if provided
    if state, ok := params["state"].(string); ok && state != "" {
        args = append(args, "--filters", fmt.Sprintf("Name=instance-state-name,Values=%s", state))
    }
    
    // Add tag filters if provided
    if tags, ok := params["tags"].(map[string]interface{}); ok {
        for key, value := range tags {
            args = append(args, "--filters", fmt.Sprintf("Name=tag:%s,Values=%v", key, value))
        }
    }
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "instances": []interface{}{},
            "count":     0,
            "error":     err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    instances := []interface{}{}
    if reservations, ok := result["Reservations"].([]interface{}); ok {
        for _, reservation := range reservations {
            if res, ok := reservation.(map[string]interface{}); ok {
                if insts, ok := res["Instances"].([]interface{}); ok {
                    instances = append(instances, insts...)
                }
            }
        }
    }
    
    return map[string]interface{}{
        "instances": instances,
        "count":     len(instances),
    }, nil
}

func (p *AWSPlugin) executeEC2Create(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    imageID, _ := params["image_id"].(string)
    instanceType, _ := params["instance_type"].(string)
    
    args := p.buildBaseArgs(params)
    args = append(args, "ec2", "run-instances", 
                  "--image-id", imageID,
                  "--instance-type", instanceType,
                  "--count", "1",
                  "--output", "json")
    
    if keyName, ok := params["key_name"].(string); ok && keyName != "" {
        args = append(args, "--key-name", keyName)
    }
    
    if secGroups, ok := params["security_groups"].([]interface{}); ok && len(secGroups) > 0 {
        sgArgs := []string{}
        for _, sg := range secGroups {
            if sgStr, ok := sg.(string); ok {
                sgArgs = append(sgArgs, sgStr)
            }
        }
        if len(sgArgs) > 0 {
            args = append(args, "--security-group-ids", strings.Join(sgArgs, " "))
        }
    }
    
    if subnetID, ok := params["subnet_id"].(string); ok && subnetID != "" {
        args = append(args, "--subnet-id", subnetID)
    }
    
    if userData, ok := params["user_data"].(string); ok && userData != "" {
        args = append(args, "--user-data", userData)
    }
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "instance_id": "",
            "state":       "failed",
            "error":       err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    instanceID := ""
    state := ""
    publicIP := ""
    privateIP := ""
    
    if instances, ok := result["Instances"].([]interface{}); ok && len(instances) > 0 {
        if inst, ok := instances[0].(map[string]interface{}); ok {
            if id, ok := inst["InstanceId"].(string); ok {
                instanceID = id
            }
            if stateInfo, ok := inst["State"].(map[string]interface{}); ok {
                if s, ok := stateInfo["Name"].(string); ok {
                    state = s
                }
            }
            if pubIP, ok := inst["PublicIpAddress"].(string); ok {
                publicIP = pubIP
            }
            if privIP, ok := inst["PrivateIpAddress"].(string); ok {
                privateIP = privIP
            }
        }
    }
    
    return map[string]interface{}{
        "instance_id": instanceID,
        "state":       state,
        "public_ip":   publicIP,
        "private_ip":  privateIP,
    }, nil
}

func (p *AWSPlugin) executeEC2Terminate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    instanceIDs, ok := params["instance_ids"].([]interface{})
    if !ok || len(instanceIDs) == 0 {
        return nil, fmt.Errorf("instance_ids parameter is required")
    }
    
    ids := []string{}
    for _, id := range instanceIDs {
        if idStr, ok := id.(string); ok {
            ids = append(ids, idStr)
        }
    }
    
    args := p.buildBaseArgs(params)
    args = append(args, "ec2", "terminate-instances",
                  "--instance-ids", strings.Join(ids, " "),
                  "--output", "json")
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.CombinedOutput()
    
    if err != nil {
        return map[string]interface{}{
            "terminated": []string{},
            "status":     "failed",
            "error":      string(output),
        }, nil
    }
    
    return map[string]interface{}{
        "terminated": ids,
        "status":     "success",
        "output":     string(output),
    }, nil
}

func (p *AWSPlugin) executeS3List(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := p.buildBaseArgs(params)
    
    bucket, hasBucket := params["bucket"].(string)
    
    if hasBucket && bucket != "" {
        // List objects in bucket
        args = append(args, "s3api", "list-objects-v2", "--bucket", bucket, "--output", "json")
        
        if prefix, ok := params["prefix"].(string); ok && prefix != "" {
            args = append(args, "--prefix", prefix)
        }
    } else {
        // List buckets
        args = append(args, "s3api", "list-buckets", "--output", "json")
    }
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "buckets": []interface{}{},
            "objects": []interface{}{},
            "count":   0,
            "error":   err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    if hasBucket && bucket != "" {
        objects := []interface{}{}
        if contents, ok := result["Contents"].([]interface{}); ok {
            objects = contents
        }
        
        return map[string]interface{}{
            "objects": objects,
            "count":   len(objects),
        }, nil
    } else {
        buckets := []interface{}{}
        if bucketList, ok := result["Buckets"].([]interface{}); ok {
            buckets = bucketList
        }
        
        return map[string]interface{}{
            "buckets": buckets,
            "count":   len(buckets),
        }, nil
    }
}

func (p *AWSPlugin) executeS3Upload(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    filePath, _ := params["file_path"].(string)
    bucket, _ := params["bucket"].(string)
    key, _ := params["key"].(string)
    
    args := p.buildBaseArgs(params)
    args = append(args, "s3api", "put-object",
                  "--bucket", bucket,
                  "--key", key,
                  "--body", filePath,
                  "--output", "json")
    
    if contentType, ok := params["content_type"].(string); ok && contentType != "" {
        args = append(args, "--content-type", contentType)
    }
    
    if acl, ok := params["acl"].(string); ok && acl != "" {
        args = append(args, "--acl", acl)
    }
    
    if metadata, ok := params["metadata"].(map[string]interface{}); ok {
        metaArgs := []string{}
        for k, v := range metadata {
            metaArgs = append(metaArgs, fmt.Sprintf("%s=%v", k, v))
        }
        if len(metaArgs) > 0 {
            args = append(args, "--metadata", strings.Join(metaArgs, ","))
        }
    }
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "etag": "",
            "url":  "",
            "size": 0,
            "error": err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    etag := ""
    if e, ok := result["ETag"].(string); ok {
        etag = strings.Trim(e, "\"")
    }
    
    url := fmt.Sprintf("s3://%s/%s", bucket, key)
    
    return map[string]interface{}{
        "etag": etag,
        "url":  url,
        "size": 0, // Would need additional call to get size
    }, nil
}

func (p *AWSPlugin) executeS3Download(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    bucket, _ := params["bucket"].(string)
    key, _ := params["key"].(string)
    filePath, _ := params["file_path"].(string)
    
    args := p.buildBaseArgs(params)
    args = append(args, "s3api", "get-object",
                  "--bucket", bucket,
                  "--key", key,
                  filePath,
                  "--output", "json")
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "size":          0,
            "etag":          "",
            "last_modified": "",
            "error":         err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    size := 0
    if s, ok := result["ContentLength"].(float64); ok {
        size = int(s)
    }
    
    etag := ""
    if e, ok := result["ETag"].(string); ok {
        etag = strings.Trim(e, "\"")
    }
    
    lastModified := ""
    if lm, ok := result["LastModified"].(string); ok {
        lastModified = lm
    }
    
    return map[string]interface{}{
        "size":          size,
        "etag":          etag,
        "last_modified": lastModified,
    }, nil
}

func (p *AWSPlugin) executeLambdaInvoke(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    functionName, _ := params["function_name"].(string)
    invocationType, _ := params["invocation_type"].(string)
    if invocationType == "" {
        invocationType = "RequestResponse"
    }
    
    args := p.buildBaseArgs(params)
    args = append(args, "lambda", "invoke",
                  "--function-name", functionName,
                  "--invocation-type", invocationType,
                  "--output", "json")
    
    if payload, ok := params["payload"].(map[string]interface{}); ok {
        payloadJSON, _ := json.Marshal(payload)
        args = append(args, "--payload", string(payloadJSON))
    }
    
    // Need temp file for response
    args = append(args, "/tmp/lambda-response.json")
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "response":    nil,
            "status_code": 0,
            "log_result":  "",
            "error":       err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    statusCode := 0
    if sc, ok := result["StatusCode"].(float64); ok {
        statusCode = int(sc)
    }
    
    logResult := ""
    if lr, ok := result["LogResult"].(string); ok {
        logResult = lr
    }
    
    return map[string]interface{}{
        "response":    result,
        "status_code": statusCode,
        "log_result":  logResult,
    }, nil
}

func (p *AWSPlugin) executeLambdaList(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := p.buildBaseArgs(params)
    args = append(args, "lambda", "list-functions", "--output", "json")
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "functions": []interface{}{},
            "count":     0,
            "error":     err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    functions := []interface{}{}
    if funcs, ok := result["Functions"].([]interface{}); ok {
        functions = funcs
    }
    
    return map[string]interface{}{
        "functions": functions,
        "count":     len(functions),
    }, nil
}

func (p *AWSPlugin) executeIAMListUsers(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := p.buildBaseArgs(params)
    args = append(args, "iam", "list-users", "--output", "json")
    
    if pathPrefix, ok := params["path_prefix"].(string); ok && pathPrefix != "" {
        args = append(args, "--path-prefix", pathPrefix)
    }
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "users": []interface{}{},
            "count": 0,
            "error": err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    users := []interface{}{}
    if userList, ok := result["Users"].([]interface{}); ok {
        users = userList
    }
    
    return map[string]interface{}{
        "users": users,
        "count": len(users),
    }, nil
}

func (p *AWSPlugin) executeVPCList(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    args := p.buildBaseArgs(params)
    args = append(args, "ec2", "describe-vpcs", "--output", "json")
    
    if filters, ok := params["filters"].(map[string]interface{}); ok {
        for key, value := range filters {
            args = append(args, "--filters", fmt.Sprintf("Name=%s,Values=%v", key, value))
        }
    }
    
    cmd := exec.CommandContext(ctx, "aws", args...)
    output, err := cmd.Output()
    
    if err != nil {
        return map[string]interface{}{
            "vpcs":  []interface{}{},
            "count": 0,
            "error": err.Error(),
        }, nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse AWS output: %w", err)
    }
    
    vpcs := []interface{}{}
    if vpcList, ok := result["Vpcs"].([]interface{}); ok {
        vpcs = vpcList
    }
    
    return map[string]interface{}{
        "vpcs":  vpcs,
        "count": len(vpcs),
    }, nil
}

var ExportedPlugin plugin.Plugin = &AWSPlugin{}
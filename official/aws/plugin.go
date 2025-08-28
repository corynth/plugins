package main

import (
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

type AWSPlugin struct{}

func NewAWSPlugin() *AWSPlugin {
	return &AWSPlugin{}
}

func (p *AWSPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "aws",
		Version:     "1.0.0",
		Description: "Amazon Web Services cloud operations and resource management",
		Author:      "Corynth Team",
		Tags:        []string{"aws", "cloud", "ec2", "s3", "lambda", "vpc", "iam", "cloud-native"},
	}
}

func (p *AWSPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"ec2_list": {
			Description: "List EC2 instances with filters",
			Inputs: map[string]IOSpec{
				"region":  {Type: "string", Required: false, Description: "AWS region"},
				"filters": {Type: "object", Required: false, Description: "Instance filters"},
				"state":   {Type: "string", Required: false, Description: "Instance state filter"},
			},
			Outputs: map[string]IOSpec{
				"instances": {Type: "array", Description: "EC2 instances"},
			},
		},
		"ec2_launch": {
			Description: "Launch EC2 instance with full configuration",
			Inputs: map[string]IOSpec{
				"image_id":        {Type: "string", Required: true, Description: "AMI ID"},
				"instance_type":   {Type: "string", Required: false, Default: "t2.micro", Description: "Instance type"},
				"key_name":        {Type: "string", Required: false, Description: "Key pair name"},
				"security_groups": {Type: "array", Required: false, Description: "Security group IDs"},
				"subnet_id":       {Type: "string", Required: false, Description: "Subnet ID"},
				"user_data":       {Type: "string", Required: false, Description: "User data script"},
				"count":           {Type: "number", Required: false, Default: 1, Description: "Number of instances"},
				"region":          {Type: "string", Required: false, Description: "AWS region"},
			},
			Outputs: map[string]IOSpec{
				"instances": {Type: "array", Description: "Launched instances"},
			},
		},
		"ec2_terminate": {
			Description: "Terminate EC2 instances",
			Inputs: map[string]IOSpec{
				"instance_ids": {Type: "array", Required: true, Description: "Instance IDs to terminate"},
				"region":       {Type: "string", Required: false, Description: "AWS region"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Termination success"},
			},
		},
		"s3_list": {
			Description: "List S3 buckets and objects",
			Inputs: map[string]IOSpec{
				"bucket": {Type: "string", Required: false, Description: "Bucket name (list objects) or empty (list buckets)"},
				"prefix": {Type: "string", Required: false, Description: "Object prefix filter"},
			},
			Outputs: map[string]IOSpec{
				"items": {Type: "array", Description: "Buckets or objects"},
			},
		},
		"s3_upload": {
			Description: "Upload files to S3 buckets",
			Inputs: map[string]IOSpec{
				"bucket":    {Type: "string", Required: true, Description: "S3 bucket name"},
				"key":       {Type: "string", Required: true, Description: "S3 object key"},
				"file_path": {Type: "string", Required: true, Description: "Local file path to upload"},
				"metadata":  {Type: "object", Required: false, Description: "Object metadata"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Upload success"},
				"url":     {Type: "string", Description: "S3 object URL"},
			},
		},
		"s3_download": {
			Description: "Download files from S3 buckets",
			Inputs: map[string]IOSpec{
				"bucket":    {Type: "string", Required: true, Description: "S3 bucket name"},
				"key":       {Type: "string", Required: true, Description: "S3 object key"},
				"file_path": {Type: "string", Required: true, Description: "Local file path to save"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Download success"},
			},
		},
		"lambda_invoke": {
			Description: "Invoke Lambda functions with payload",
			Inputs: map[string]IOSpec{
				"function_name":     {Type: "string", Required: true, Description: "Lambda function name"},
				"payload":           {Type: "object", Required: false, Description: "Function payload"},
				"invocation_type":   {Type: "string", Required: false, Default: "RequestResponse", Description: "Synchronous or Event"},
				"region":            {Type: "string", Required: false, Description: "AWS region"},
			},
			Outputs: map[string]IOSpec{
				"response":    {Type: "object", Description: "Function response"},
				"status_code": {Type: "number", Description: "HTTP status code"},
			},
		},
		"lambda_list": {
			Description: "List Lambda functions",
			Inputs: map[string]IOSpec{
				"prefix": {Type: "string", Required: false, Description: "Function name prefix"},
				"region": {Type: "string", Required: false, Description: "AWS region"},
			},
			Outputs: map[string]IOSpec{
				"functions": {Type: "array", Description: "Lambda functions"},
			},
		},
	}
}

func (p *AWSPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "ec2_list":
		return p.ec2List(params)
	case "ec2_launch":
		return p.ec2Launch(params)
	case "ec2_terminate":
		return p.ec2Terminate(params)
	case "s3_list":
		return p.s3List(params)
	case "s3_upload":
		return p.s3Upload(params)
	case "s3_download":
		return p.s3Download(params)
	case "lambda_invoke":
		return p.lambdaInvoke(params)
	case "lambda_list":
		return p.lambdaList(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *AWSPlugin) ec2List(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"ec2", "describe-instances", "--output", "json"}
	
	if region, ok := params["region"].(string); ok && region != "" {
		args = append(args, "--region", region)
	}
	
	if state, ok := params["state"].(string); ok && state != "" {
		args = append(args, "--filters", fmt.Sprintf("Name=instance-state-name,Values=%s", state))
	}
	
	output, err := exec.Command("aws", args...).Output()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err)}, nil
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}, nil
	}
	
	instances := []map[string]interface{}{}
	if reservations, ok := result["Reservations"].([]interface{}); ok {
		for _, reservation := range reservations {
			if reservationMap, ok := reservation.(map[string]interface{}); ok {
				if instancesList, ok := reservationMap["Instances"].([]interface{}); ok {
					for _, instance := range instancesList {
						if instanceMap, ok := instance.(map[string]interface{}); ok {
							instances = append(instances, instanceMap)
						}
					}
				}
			}
		}
	}
	
	return map[string]interface{}{"instances": instances}, nil
}

func (p *AWSPlugin) ec2Launch(params map[string]interface{}) (map[string]interface{}, error) {
	imageId, ok := params["image_id"].(string)
	if !ok || imageId == "" {
		return map[string]interface{}{"error": "image_id is required"}, nil
	}
	
	args := []string{"ec2", "run-instances", "--image-id", imageId, "--output", "json"}
	
	if instanceType, ok := params["instance_type"].(string); ok && instanceType != "" {
		args = append(args, "--instance-type", instanceType)
	} else {
		args = append(args, "--instance-type", "t2.micro")
	}
	
	if count, ok := params["count"].(float64); ok {
		countStr := fmt.Sprintf("%.0f", count)
		args = append(args, "--count", countStr)
	} else {
		args = append(args, "--count", "1")
	}
	
	if region, ok := params["region"].(string); ok && region != "" {
		args = append(args, "--region", region)
	}
	
	if keyName, ok := params["key_name"].(string); ok && keyName != "" {
		args = append(args, "--key-name", keyName)
	}
	
	if userData, ok := params["user_data"].(string); ok && userData != "" {
		args = append(args, "--user-data", userData)
	}
	
	output, err := exec.Command("aws", args...).Output()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err)}, nil
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}, nil
	}
	
	return result, nil
}

func (p *AWSPlugin) ec2Terminate(params map[string]interface{}) (map[string]interface{}, error) {
	instanceIds, ok := params["instance_ids"].([]interface{})
	if !ok || len(instanceIds) == 0 {
		return map[string]interface{}{"error": "instance_ids is required"}, nil
	}
	
	ids := make([]string, len(instanceIds))
	for i, id := range instanceIds {
		if idStr, ok := id.(string); ok {
			ids[i] = idStr
		} else {
			return map[string]interface{}{"error": "invalid instance ID format"}, nil
		}
	}
	
	args := []string{"ec2", "terminate-instances", "--instance-ids"}
	args = append(args, ids...)
	args = append(args, "--output", "json")
	
	if region, ok := params["region"].(string); ok && region != "" {
		args = append(args, "--region", region)
	}
	
	_, err := exec.Command("aws", args...).Output()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err), "success": false}, nil
	}
	
	return map[string]interface{}{"success": true}, nil
}

func (p *AWSPlugin) s3List(params map[string]interface{}) (map[string]interface{}, error) {
	bucket, hasBucket := params["bucket"].(string)
	
	if !hasBucket || bucket == "" {
		// List buckets
		output, err := exec.Command("aws", "s3api", "list-buckets", "--output", "json").Output()
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err)}, nil
		}
		
		var result map[string]interface{}
		if err := json.Unmarshal(output, &result); err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}, nil
		}
		
		return map[string]interface{}{"items": result["Buckets"]}, nil
	} else {
		// List objects in bucket
		args := []string{"s3api", "list-objects-v2", "--bucket", bucket, "--output", "json"}
		
		if prefix, ok := params["prefix"].(string); ok && prefix != "" {
			args = append(args, "--prefix", prefix)
		}
		
		output, err := exec.Command("aws", args...).Output()
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err)}, nil
		}
		
		var result map[string]interface{}
		if err := json.Unmarshal(output, &result); err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}, nil
		}
		
		contents := result["Contents"]
		if contents == nil {
			contents = []interface{}{}
		}
		
		return map[string]interface{}{"items": contents}, nil
	}
}

func (p *AWSPlugin) s3Upload(params map[string]interface{}) (map[string]interface{}, error) {
	bucket, ok := params["bucket"].(string)
	if !ok || bucket == "" {
		return map[string]interface{}{"error": "bucket is required"}, nil
	}
	
	key, ok := params["key"].(string)
	if !ok || key == "" {
		return map[string]interface{}{"error": "key is required"}, nil
	}
	
	filePath, ok := params["file_path"].(string)
	if !ok || filePath == "" {
		return map[string]interface{}{"error": "file_path is required"}, nil
	}
	
	args := []string{"s3", "cp", filePath, fmt.Sprintf("s3://%s/%s", bucket, key)}
	
	err := exec.Command("aws", args...).Run()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err), "success": false}, nil
	}
	
	return map[string]interface{}{
		"success": true,
		"url":     fmt.Sprintf("s3://%s/%s", bucket, key),
	}, nil
}

func (p *AWSPlugin) s3Download(params map[string]interface{}) (map[string]interface{}, error) {
	bucket, ok := params["bucket"].(string)
	if !ok || bucket == "" {
		return map[string]interface{}{"error": "bucket is required"}, nil
	}
	
	key, ok := params["key"].(string)
	if !ok || key == "" {
		return map[string]interface{}{"error": "key is required"}, nil
	}
	
	filePath, ok := params["file_path"].(string)
	if !ok || filePath == "" {
		return map[string]interface{}{"error": "file_path is required"}, nil
	}
	
	args := []string{"s3", "cp", fmt.Sprintf("s3://%s/%s", bucket, key), filePath}
	
	err := exec.Command("aws", args...).Run()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err), "success": false}, nil
	}
	
	return map[string]interface{}{"success": true}, nil
}

func (p *AWSPlugin) lambdaInvoke(params map[string]interface{}) (map[string]interface{}, error) {
	functionName, ok := params["function_name"].(string)
	if !ok || functionName == "" {
		return map[string]interface{}{"error": "function_name is required"}, nil
	}
	
	args := []string{"lambda", "invoke", "--function-name", functionName, "--output", "json"}
	
	if region, ok := params["region"].(string); ok && region != "" {
		args = append(args, "--region", region)
	}
	
	if invocationType, ok := params["invocation_type"].(string); ok && invocationType != "" {
		args = append(args, "--invocation-type", invocationType)
	}
	
	args = append(args, "/tmp/lambda-response.json")
	
	if payload, ok := params["payload"]; ok {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to marshal payload: %v", err)}, nil
		}
		args = append(args, "--payload", string(payloadBytes))
	}
	
	output, err := exec.Command("aws", args...).Output()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err)}, nil
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}, nil
	}
	
	// Read response payload
	responseData, err := os.ReadFile("/tmp/lambda-response.json")
	if err == nil {
		var responsePayload interface{}
		if json.Unmarshal(responseData, &responsePayload) == nil {
			result["response"] = responsePayload
		}
		os.Remove("/tmp/lambda-response.json")
	}
	
	return result, nil
}

func (p *AWSPlugin) lambdaList(params map[string]interface{}) (map[string]interface{}, error) {
	args := []string{"lambda", "list-functions", "--output", "json"}
	
	if region, ok := params["region"].(string); ok && region != "" {
		args = append(args, "--region", region)
	}
	
	output, err := exec.Command("aws", args...).Output()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("aws command failed: %v", err)}, nil
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}, nil
	}
	
	functions := result["Functions"]
	if functions == nil {
		functions = []interface{}{}
	}
	
	// Filter by prefix if provided
	if prefix, ok := params["prefix"].(string); ok && prefix != "" {
		if functionsList, ok := functions.([]interface{}); ok {
			filtered := []interface{}{}
			for _, fn := range functionsList {
				if fnMap, ok := fn.(map[string]interface{}); ok {
					if name, ok := fnMap["FunctionName"].(string); ok && strings.HasPrefix(name, prefix) {
						filtered = append(filtered, fn)
					}
				}
			}
			functions = filtered
		}
	}
	
	return map[string]interface{}{"functions": functions}, nil
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}
	
	action := os.Args[1]
	plugin := NewAWSPlugin()
	
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
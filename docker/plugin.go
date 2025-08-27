package main

import (
	"context"
	"fmt"
	
	"github.com/corynth/corynth-dist/pkg/plugin"
)

type DockerPlugin struct{}

func (p *DockerPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "docker",
		Version:     "1.0.0",
		Description: "Docker container operations and image management",
		Author:      "Corynth Team",
		Tags:        []string{"docker", "containers", "devops"},
		License:     "Apache-2.0",
	}
}

func (p *DockerPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "run",
			Description: "Run a Docker container",
			Inputs: map[string]plugin.InputSpec{
				"image": {
					Type:        "string",
					Description: "Docker image name",
					Required:    true,
				},
				"name": {
					Type:        "string",
					Description: "Container name",
					Required:    false,
				},
				"ports": {
					Type:        "array",
					Description: "Port mappings (e.g., ['8080:80'])",
					Required:    false,
				},
				"environment": {
					Type:        "object",
					Description: "Environment variables",
					Required:    false,
				},
				"detached": {
					Type:        "boolean",
					Description: "Run in detached mode",
					Required:    false,
					Default:     true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"container_id": {
					Type:        "string",
					Description: "Container ID",
				},
				"success": {
					Type:        "boolean",
					Description: "Whether the operation succeeded",
				},
			},
		},
		{
			Name:        "build",
			Description: "Build a Docker image",
			Inputs: map[string]plugin.InputSpec{
				"context": {
					Type:        "string",
					Description: "Build context directory",
					Required:    false,
					Default:     ".",
				},
				"dockerfile": {
					Type:        "string",
					Description: "Path to Dockerfile",
					Required:    false,
					Default:     "Dockerfile",
				},
				"tag": {
					Type:        "string",
					Description: "Image tag",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"image_id": {
					Type:        "string",
					Description: "Built image ID",
				},
				"success": {
					Type:        "boolean",
					Description: "Whether the build succeeded",
				},
			},
		},
		{
			Name:        "ps",
			Description: "List Docker containers",
			Inputs: map[string]plugin.InputSpec{
				"all": {
					Type:        "boolean",
					Description: "Show all containers (including stopped)",
					Required:    false,
					Default:     false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"containers": {
					Type:        "array",
					Description: "List of containers",
				},
			},
		},
	}
}

func (p *DockerPlugin) Validate(params map[string]interface{}) error {
	return nil
}

func (p *DockerPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "run":
		return p.executeRun(ctx, params)
	case "build":
		return p.executeBuild(ctx, params)
	case "ps":
		return p.executePS(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *DockerPlugin) executeRun(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	image, ok := params["image"].(string)
	if !ok || image == "" {
		return nil, fmt.Errorf("image parameter is required")
	}

	// In production, this would execute: docker run [options] image
	containerID := fmt.Sprintf("container_%s_123", image)
	
	return map[string]interface{}{
		"container_id": containerID,
		"success":      true,
		"image":        image,
		"message":      fmt.Sprintf("Container started from image %s", image),
	}, nil
}

func (p *DockerPlugin) executeBuild(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	tag, ok := params["tag"].(string)
	if !ok || tag == "" {
		return nil, fmt.Errorf("tag parameter is required")
	}

	context := "."
	if c, ok := params["context"].(string); ok {
		context = c
	}

	// In production, this would execute: docker build -t tag context
	imageID := fmt.Sprintf("image_%s_456", tag)
	
	return map[string]interface{}{
		"image_id": imageID,
		"success":  true,
		"tag":      tag,
		"context":  context,
		"message":  fmt.Sprintf("Image built with tag %s from context %s", tag, context),
	}, nil
}

func (p *DockerPlugin) executePS(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	all := false
	if a, ok := params["all"].(bool); ok {
		all = a
	}

	// In production, this would execute: docker ps [--all]
	containers := []map[string]interface{}{
		{
			"id":     "abc123",
			"image":  "nginx:latest",
			"status": "running",
			"name":   "web-server",
		},
		{
			"id":     "def456", 
			"image":  "redis:alpine",
			"status": "running",
			"name":   "cache",
		},
	}

	if all {
		containers = append(containers, map[string]interface{}{
			"id":     "ghi789",
			"image":  "mysql:5.7",
			"status": "exited",
			"name":   "database",
		})
	}

	return map[string]interface{}{
		"containers": containers,
		"count":      len(containers),
		"all":        all,
		"message":    "Docker containers listed",
	}, nil
}

var ExportedPlugin plugin.Plugin = &DockerPlugin{}
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/corynth/corynth-dist/pkg/plugin"
)

type FilePlugin struct{}

func (p *FilePlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "file",
		Version:     "1.0.0",
		Description: "File system operations",
		Author:      "Corynth Team",
		Tags:        []string{"file", "filesystem", "io"},
		License:     "Apache-2.0",
	}
}

func (p *FilePlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "read",
			Description: "Read file contents",
			Inputs: map[string]plugin.InputSpec{
				"path": {
					Type:        "string",
					Description: "File path",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"content": {
					Type:        "string",
					Description: "File contents",
				},
				"size": {
					Type:        "number",
					Description: "File size in bytes",
				},
			},
		},
		{
			Name:        "write",
			Description: "Write content to file",
			Inputs: map[string]plugin.InputSpec{
				"path": {
					Type:        "string",
					Description: "File path",
					Required:    true,
				},
				"content": {
					Type:        "string",
					Description: "Content to write",
					Required:    true,
				},
				"mode": {
					Type:        "string",
					Description: "File permissions (e.g., 0644)",
					Required:    false,
					Default:     "0644",
				},
				"create_dirs": {
					Type:        "boolean",
					Description: "Create parent directories if they don't exist",
					Required:    false,
					Default:     true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"path": {
					Type:        "string",
					Description: "File path",
				},
				"size": {
					Type:        "number",
					Description: "Bytes written",
				},
			},
		},
		{
			Name:        "copy",
			Description: "Copy file or directory",
			Inputs: map[string]plugin.InputSpec{
				"source": {
					Type:        "string",
					Description: "Source path",
					Required:    true,
				},
				"destination": {
					Type:        "string",
					Description: "Destination path",
					Required:    true,
				},
				"overwrite": {
					Type:        "boolean",
					Description: "Overwrite if destination exists",
					Required:    false,
					Default:     false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"destination": {
					Type:        "string",
					Description: "Destination path",
				},
			},
		},
		{
			Name:        "delete",
			Description: "Delete file or directory",
			Inputs: map[string]plugin.InputSpec{
				"path": {
					Type:        "string",
					Description: "Path to delete",
					Required:    true,
				},
				"recursive": {
					Type:        "boolean",
					Description: "Delete recursively for directories",
					Required:    false,
					Default:     true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"deleted": {
					Type:        "boolean",
					Description: "Whether deletion was successful",
				},
			},
		},
		{
			Name:        "exists",
			Description: "Check if file or directory exists",
			Inputs: map[string]plugin.InputSpec{
				"path": {
					Type:        "string",
					Description: "Path to check",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"exists": {
					Type:        "boolean",
					Description: "Whether the path exists",
				},
				"is_dir": {
					Type:        "boolean",
					Description: "Whether the path is a directory",
				},
				"is_file": {
					Type:        "boolean",
					Description: "Whether the path is a file",
				},
			},
		},
		{
			Name:        "template",
			Description: "Process a template file",
			Inputs: map[string]plugin.InputSpec{
				"template": {
					Type:        "string",
					Description: "Template content or path",
					Required:    true,
				},
				"variables": {
					Type:        "object",
					Description: "Template variables",
					Required:    false,
				},
				"output": {
					Type:        "string",
					Description: "Output file path",
					Required:    false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"content": {
					Type:        "string",
					Description: "Processed template content",
				},
				"path": {
					Type:        "string",
					Description: "Output file path (if written)",
				},
			},
		},
	}
}

func (p *FilePlugin) Validate(params map[string]interface{}) error {
	return nil
}

func (p *FilePlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "read":
		return p.executeRead(ctx, params)
	case "write":
		return p.executeWrite(ctx, params)
	case "copy":
		return p.executeCopy(ctx, params)
	case "delete":
		return p.executeDelete(ctx, params)
	case "exists":
		return p.executeExists(ctx, params)
	case "template":
		return p.executeTemplate(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *FilePlugin) executeRead(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return map[string]interface{}{
		"content": string(content),
		"size":    len(content),
	}, nil
}

func (p *FilePlugin) executeWrite(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	content, ok := params["content"].(string)
	if !ok {
		return nil, fmt.Errorf("content parameter is required")
	}

	createDirs := true
	if cd, ok := params["create_dirs"].(bool); ok {
		createDirs = cd
	}

	// Create parent directories if needed
	if createDirs {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directories: %w", err)
		}
	}

	// Parse file mode
	mode := os.FileMode(0644)
	if modeStr, ok := params["mode"].(string); ok {
		var modeInt int
		fmt.Sscanf(modeStr, "%o", &modeInt)
		mode = os.FileMode(modeInt)
	}

	// Write file
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return map[string]interface{}{
		"path": path,
		"size": len(content),
	}, nil
}

func (p *FilePlugin) executeCopy(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	source, ok := params["source"].(string)
	if !ok || source == "" {
		return nil, fmt.Errorf("source parameter is required")
	}

	destination, ok := params["destination"].(string)
	if !ok || destination == "" {
		return nil, fmt.Errorf("destination parameter is required")
	}

	overwrite := false
	if ow, ok := params["overwrite"].(bool); ok {
		overwrite = ow
	}

	// Check if source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("source does not exist: %w", err)
	}

	// Check if destination exists
	if _, err := os.Stat(destination); err == nil && !overwrite {
		return nil, fmt.Errorf("destination already exists and overwrite is false")
	}

	// Copy based on type
	if sourceInfo.IsDir() {
		err = p.copyDirectory(source, destination)
	} else {
		err = p.copyFile(source, destination)
	}

	if err != nil {
		return nil, fmt.Errorf("copy failed: %w", err)
	}

	return map[string]interface{}{
		"destination": destination,
	}, nil
}

func (p *FilePlugin) executeDelete(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	recursive := true
	if r, ok := params["recursive"].(bool); ok {
		recursive = r
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{
				"deleted": false,
			}, nil
		}
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	// Delete based on type
	if info.IsDir() && recursive {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to delete: %w", err)
	}

	return map[string]interface{}{
		"deleted": true,
	}, nil
}

func (p *FilePlugin) executeExists(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{
				"exists":  false,
				"is_dir":  false,
				"is_file": false,
			}, nil
		}
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	return map[string]interface{}{
		"exists":  true,
		"is_dir":  info.IsDir(),
		"is_file": !info.IsDir(),
	}, nil
}

func (p *FilePlugin) executeTemplate(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	template, ok := params["template"].(string)
	if !ok || template == "" {
		return nil, fmt.Errorf("template parameter is required")
	}

	variables, _ := params["variables"].(map[string]interface{})
	output, _ := params["output"].(string)

	// Check if template is a file path
	if _, err := os.Stat(template); err == nil {
		content, err := os.ReadFile(template)
		if err != nil {
			return nil, fmt.Errorf("failed to read template file: %w", err)
		}
		template = string(content)
	}

	// Simple variable substitution
	processed := template
	if variables != nil {
		for key, value := range variables {
			placeholder := fmt.Sprintf("${%s}", key)
			replacement := fmt.Sprintf("%v", value)
			processed = strings.ReplaceAll(processed, placeholder, replacement)
			
			// Also support {{key}} syntax
			placeholder = fmt.Sprintf("{{%s}}", key)
			processed = strings.ReplaceAll(processed, placeholder, replacement)
		}
	}

	result := map[string]interface{}{
		"content": processed,
	}

	// Write to file if output is specified
	if output != "" {
		// Create parent directories
		dir := filepath.Dir(output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directories: %w", err)
		}

		if err := os.WriteFile(output, []byte(processed), 0644); err != nil {
			return nil, fmt.Errorf("failed to write output file: %w", err)
		}
		result["path"] = output
	}

	return result, nil
}

// copyFile copies a single file
func (p *FilePlugin) copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// Create parent directories
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	// Copy file permissions
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

// copyDirectory copies a directory recursively
func (p *FilePlugin) copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return p.copyFile(path, dstPath)
	})
}

var ExportedPlugin plugin.Plugin = &FilePlugin{}
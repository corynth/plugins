package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

type FilePlugin struct{}

func NewFilePlugin() *FilePlugin {
	return &FilePlugin{}
}

func (p *FilePlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "file",
		Version:     "1.0.0",
		Description: "File system operations (read, write, copy, move)",
		Author:      "Corynth Team",
		Tags:        []string{"file", "filesystem", "io"},
	}
}

func (p *FilePlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"read": {
			Description: "Read file contents",
			Inputs: map[string]IOSpec{
				"path": {Type: "string", Required: true, Description: "File path to read"},
			},
			Outputs: map[string]IOSpec{
				"content": {Type: "string", Description: "File content"},
				"size":    {Type: "number", Description: "File size in bytes"},
			},
		},
		"write": {
			Description: "Write content to files with directory creation",
			Inputs: map[string]IOSpec{
				"path":        {Type: "string", Required: true, Description: "File path to write"},
				"content":     {Type: "string", Required: true, Description: "Content to write"},
				"create_dirs": {Type: "boolean", Required: false, Default: false, Description: "Create directories if they don't exist"},
				"append":      {Type: "boolean", Required: false, Default: false, Description: "Append to file instead of overwriting"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Write success"},
				"size":    {Type: "number", Description: "Bytes written"},
			},
		},
		"copy": {
			Description: "Copy files and directories",
			Inputs: map[string]IOSpec{
				"source":      {Type: "string", Required: true, Description: "Source path"},
				"destination": {Type: "string", Required: true, Description: "Destination path"},
				"create_dirs": {Type: "boolean", Required: false, Default: false, Description: "Create destination directories"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Copy success"},
			},
		},
		"move": {
			Description: "Move or rename files",
			Inputs: map[string]IOSpec{
				"source":      {Type: "string", Required: true, Description: "Source path"},
				"destination": {Type: "string", Required: true, Description: "Destination path"},
				"create_dirs": {Type: "boolean", Required: false, Default: false, Description: "Create destination directories"},
			},
			Outputs: map[string]IOSpec{
				"success": {Type: "boolean", Description: "Move success"},
			},
		},
	}
}

func (p *FilePlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "read":
		return p.readFile(params)
	case "write":
		return p.writeFile(params)
	case "copy":
		return p.copyFile(params)
	case "move":
		return p.moveFile(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *FilePlugin) readFile(params map[string]interface{}) (map[string]interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return map[string]interface{}{"error": "path is required"}, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to read file: %v", err)}, nil
	}

	return map[string]interface{}{
		"content": string(content),
		"size":    len(content),
	}, nil
}

func (p *FilePlugin) writeFile(params map[string]interface{}) (map[string]interface{}, error) {
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return map[string]interface{}{"error": "path is required"}, nil
	}

	content, ok := params["content"].(string)
	if !ok {
		return map[string]interface{}{"error": "content is required"}, nil
	}

	createDirs := getBoolParam(params, "create_dirs", false)
	appendMode := getBoolParam(params, "append", false)

	// Create directories if requested
	if createDirs {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return map[string]interface{}{
				"error":   fmt.Sprintf("failed to create directories: %v", err),
				"success": false,
			}, nil
		}
	}

	// Determine file mode
	flag := os.O_CREATE | os.O_WRONLY
	if appendMode {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return map[string]interface{}{
			"error":   fmt.Sprintf("failed to open file: %v", err),
			"success": false,
		}, nil
	}
	defer file.Close()

	bytesWritten, err := file.WriteString(content)
	if err != nil {
		return map[string]interface{}{
			"error":   fmt.Sprintf("failed to write file: %v", err),
			"success": false,
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"size":    bytesWritten,
	}, nil
}

func (p *FilePlugin) copyFile(params map[string]interface{}) (map[string]interface{}, error) {
	source, ok := params["source"].(string)
	if !ok || source == "" {
		return map[string]interface{}{"error": "source is required"}, nil
	}

	destination, ok := params["destination"].(string)
	if !ok || destination == "" {
		return map[string]interface{}{"error": "destination is required"}, nil
	}

	createDirs := getBoolParam(params, "create_dirs", false)

	// Create destination directories if requested
	if createDirs {
		dir := filepath.Dir(destination)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return map[string]interface{}{
				"error":   fmt.Sprintf("failed to create directories: %v", err),
				"success": false,
			}, nil
		}
	}

	// Check if source is a file or directory
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return map[string]interface{}{
			"error":   fmt.Sprintf("failed to stat source: %v", err),
			"success": false,
		}, nil
	}

	if sourceInfo.IsDir() {
		// Copy directory recursively
		err = copyDir(source, destination)
	} else {
		// Copy single file
		err = copyFile(source, destination)
	}

	if err != nil {
		return map[string]interface{}{
			"error":   fmt.Sprintf("failed to copy: %v", err),
			"success": false,
		}, nil
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (p *FilePlugin) moveFile(params map[string]interface{}) (map[string]interface{}, error) {
	source, ok := params["source"].(string)
	if !ok || source == "" {
		return map[string]interface{}{"error": "source is required"}, nil
	}

	destination, ok := params["destination"].(string)
	if !ok || destination == "" {
		return map[string]interface{}{"error": "destination is required"}, nil
	}

	createDirs := getBoolParam(params, "create_dirs", false)

	// Create destination directories if requested
	if createDirs {
		dir := filepath.Dir(destination)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return map[string]interface{}{
				"error":   fmt.Sprintf("failed to create directories: %v", err),
				"success": false,
			}, nil
		}
	}

	err := os.Rename(source, destination)
	if err != nil {
		return map[string]interface{}{
			"error":   fmt.Sprintf("failed to move file: %v", err),
			"success": false,
		}, nil
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

// Helper functions
func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultValue
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func copyDir(src, dst string) error {
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

		return copyFile(path, dstPath)
	})
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewFilePlugin()

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
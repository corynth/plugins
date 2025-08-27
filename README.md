# Corynth Plugin Sources

This repository contains source code for Corynth plugins that are compiled on-demand when installed.

## Plugin Structure

Each plugin should be in its own directory with the following structure:

```
plugin-name/
├── plugin.go         # Main plugin implementation
├── go.mod           # Go module file (optional - will be generated if missing)
├── README.md        # Plugin documentation
└── examples/        # Example usage (optional)
```

## Available Plugins

### Calculator
Mathematical calculations and unit conversions
- **Directory**: `calculator/`
- **Actions**: `calculate`, `convert`

### Weather
Weather information and forecasts
- **Directory**: `weather/`  
- **Actions**: `current`, `forecast`

### JSON Processor
JSON parsing, manipulation, and validation
- **Directory**: `json-processor/`
- **Actions**: `parse`, `query`, `validate`, `transform`

## Installation

Plugins are automatically compiled and installed from source when you use:

```bash
corynth plugin install <plugin-name>
```

The system will:
1. Clone this repository
2. Find the plugin source directory
3. Compile the plugin with the correct interface
4. Install the compiled plugin locally

## Creating New Plugins

1. Create a new directory with your plugin name
2. Implement the Corynth plugin interface in `plugin.go`
3. Test your plugin locally
4. Submit a pull request

### Plugin Template

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/corynth/corynth-dist/src/pkg/plugin"
)

type YourPlugin struct{}

func (p *YourPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "your-plugin",
        Version:     "1.0.0", 
        Description: "Description of your plugin",
        Author:      "Your Name",
        Tags:        []string{"tag1", "tag2"},
        License:     "Apache-2.0",
    }
}

func (p *YourPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "action-name",
            Description: "Description of the action",
            Inputs: map[string]plugin.InputSpec{
                "param1": {
                    Type:        "string",
                    Description: "Parameter description", 
                    Required:    true,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "result": {
                    Type:        "string",
                    Description: "Result description",
                },
            },
        },
    }
}

func (p *YourPlugin) Validate(params map[string]interface{}) error {
    return nil // Add validation logic
}

func (p *YourPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "action-name":
        return p.executeAction(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *YourPlugin) executeAction(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    // Implementation here
    return map[string]interface{}{
        "result": "success",
    }, nil
}

var ExportedPlugin plugin.Plugin = &YourPlugin{}
```

## Requirements

- Go 1.21+
- Must implement the Corynth plugin interface
- Must export plugin as `ExportedPlugin`
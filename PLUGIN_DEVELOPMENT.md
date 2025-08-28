# Corynth Plugin Development Guide

## Overview

Corynth uses **Go-first JSON protocol plugins** for maximum performance, security, and cross-platform compatibility. All plugins communicate via JSON stdin/stdout and are implemented in Go, providing fast execution and consistent behavior across all environments.

## Architecture

### JSON Protocol Communication
- **Input**: JSON via stdin for parameters
- **Output**: JSON via stdout for results  
- **Metadata**: Special commands `metadata` and `actions`
- **Execution**: Plugin receives action name as first argument

### Go-First Approach
- **Primary Language**: Go for all official plugins
- **Performance**: Native compilation, fast startup
- **Security**: Type safety, memory safety
- **Portability**: Single binary deployment

## Quick Start with Bootstrap

Use our automated bootstrap script to create a new plugin:

```bash
# Generate complete plugin structure
./official/bootstrap-plugin.sh my-new-plugin

# This creates:
# my-new-plugin/
# ├── plugin.go      # Full implementation with helper functions
# ├── plugin         # Executable wrapper script  
# ├── go.mod         # Go module configuration
# └── samples/       # Example HCL workflows
```

The bootstrap generates a complete, working plugin that you can immediately test and customize.

## Manual Plugin Development

### 1. Directory Structure

```
your-plugin/
├── plugin.go         # Main Go implementation
├── plugin           # Executable bash wrapper  
├── go.mod           # Go module file
├── README.md        # Documentation
└── samples/         # Example workflows
    └── basic.hcl    # Sample HCL usage
```

### 2. Plugin Interface

Every plugin must implement these core functions:

```go
package main

import (
    "encoding/json"
    "fmt" 
    "io"
    "os"
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

type YourPlugin struct{}

func NewYourPlugin() *YourPlugin {
    return &YourPlugin{}
}

func (p *YourPlugin) GetMetadata() Metadata {
    return Metadata{
        Name:        "your-plugin",
        Version:     "1.0.0", 
        Description: "Your plugin description",
        Author:      "Your Name",
        Tags:        []string{"category", "functionality"},
    }
}

func (p *YourPlugin) GetActions() map[string]ActionSpec {
    return map[string]ActionSpec{
        "action_name": {
            Description: "Action description",
            Inputs: map[string]IOSpec{
                "param1": {Type: "string", Required: true, Description: "Parameter description"},
                "param2": {Type: "number", Required: false, Default: 10, Description: "Optional parameter"},
            },
            Outputs: map[string]IOSpec{
                "result": {Type: "string", Description: "Result description"},
                "success": {Type: "boolean", Description: "Operation success"},
            },
        },
    }
}

func (p *YourPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "action_name":
        return p.executeAction(params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *YourPlugin) executeAction(params map[string]interface{}) (map[string]interface{}, error) {
    // Your implementation here
    param1, ok := params["param1"].(string)
    if !ok {
        return map[string]interface{}{"error": "param1 is required"}, nil
    }
    
    return map[string]interface{}{
        "result": fmt.Sprintf("Processed: %s", param1),
        "success": true,
    }, nil
}
```

### 3. Main Function and JSON Protocol

```go
func main() {
    if len(os.Args) < 2 {
        json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
        os.Exit(1)
    }
    
    action := os.Args[1]
    plugin := NewYourPlugin()
    
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
```

### 4. Bash Wrapper Script

Create an executable `plugin` file:

```bash
#!/usr/bin/env bash
# Corynth Your-Plugin
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec go run "$DIR/plugin.go" "$@"
```

### 5. Go Module Setup

Create `go.mod`:

```go
module github.com/corynth/plugins/official/your-plugin

go 1.21
```

## Best Practices

### Security
- **Input validation**: Always validate and sanitize inputs
- **Error handling**: Return errors in JSON format, never panic
- **Resource limits**: Implement timeouts for long operations
- **No arbitrary execution**: Avoid `eval` or dynamic code execution

### Performance  
- **Fast startup**: Minimize initialization time
- **Memory efficiency**: Clean up resources properly
- **Concurrent safety**: Handle concurrent execution if needed
- **Binary size**: Use only necessary dependencies

### Error Handling
```go
// Always return errors in JSON format
if err != nil {
    return map[string]interface{}{
        "error": fmt.Sprintf("operation failed: %v", err),
        "success": false,
    }, nil  // Note: nil error, error is in response
}
```

### Helper Functions
```go
// String parameter with default
func getStringParam(params map[string]interface{}, key string, defaultValue string) string {
    if val, ok := params[key].(string); ok {
        return val
    }
    return defaultValue
}

// Boolean parameter with default  
func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
    if val, ok := params[key].(bool); ok {
        return val
    }
    return defaultValue
}

// Number parameter with default
func getNumberParam(params map[string]interface{}, key string, defaultValue float64) float64 {
    if val, ok := params[key].(float64); ok {
        return val
    }
    return defaultValue
}
```

## Testing Your Plugin

### 1. Test Metadata
```bash
cd your-plugin
go run plugin.go metadata
```

### 2. Test Actions
```bash  
go run plugin.go actions
```

### 3. Test Execution
```bash
echo '{"param1": "test"}' | go run plugin.go action_name
```

### 4. Integration Testing
Create a sample HCL workflow:

```hcl
# samples/basic.hcl
step "test" {
  plugin = "your-plugin"
  action = "action_name" 
  params = {
    param1 = "hello world"
    param2 = 42
  }
}
```

## Official Plugin Examples

Study these production plugins for reference:

- **http**: REST API calls with authentication
- **file**: File system operations with safety checks  
- **docker**: Container management with subprocess calls
- **kubernetes**: kubectl integration with JSON parsing
- **sql**: Database operations with multiple drivers
- **shell**: Command execution with security controls

## Plugin Categories

### Infrastructure
- **docker**: Container operations
- **kubernetes**: Cluster management  
- **terraform**: Infrastructure as Code
- **ansible**: Configuration management

### Cloud Providers
- **aws**: Amazon Web Services
- **gcp**: Google Cloud Platform (planned)
- **azure**: Microsoft Azure (planned)

### Communication
- **slack**: Team messaging
- **email**: SMTP notifications
- **teams**: Microsoft Teams (planned)

### Data & Storage  
- **sql**: Database operations
- **file**: File system management
- **redis**: Caching (planned)

### AI & Analytics
- **llm**: Large Language Models
- **reporting**: Report generation

### Utilities
- **shell**: Command execution
- **http**: Web requests
- **calculator**: Mathematical operations

## Submission Guidelines

### 1. Code Quality
- Follow Go conventions and best practices
- Include comprehensive error handling
- Add input validation for all parameters
- Write clear, documented code

### 2. Testing
- Test all actions with various inputs
- Verify error conditions return proper JSON
- Test with real-world scenarios  
- Include sample workflows

### 3. Documentation
- Complete README with usage examples
- Document all actions and parameters
- Include troubleshooting section
- Provide sample HCL workflows

### 4. Pull Request
- Use descriptive commit messages
- Include tests and documentation
- Follow the established directory structure
- Test with the bootstrap framework

## Migration from Legacy Plugins

If migrating from Python or other implementations:

1. **Use Bootstrap**: Generate Go skeleton with `./bootstrap-plugin.sh`
2. **Port Logic**: Convert core functionality to Go
3. **Test Thoroughly**: Verify all actions work identically  
4. **Update Registry**: Ensure registry.json reflects new implementation
5. **Document Changes**: Note any behavioral differences

## Advanced Topics

### External Dependencies
```go
// go.mod with external packages
module github.com/corynth/plugins/official/your-plugin

go 1.21

require (
    github.com/gorilla/websocket v1.5.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### Configuration Files
```go
// Load config from environment or files
func (p *YourPlugin) loadConfig() (*Config, error) {
    config := &Config{
        APIKey: os.Getenv("PLUGIN_API_KEY"),
        Timeout: 30,
    }
    
    if config.APIKey == "" {
        return nil, fmt.Errorf("PLUGIN_API_KEY environment variable required")
    }
    
    return config, nil
}
```

### Long-Running Operations
```go
func (p *YourPlugin) longOperation(params map[string]interface{}) (map[string]interface{}, error) {
    timeout := getNumberParam(params, "timeout", 300)
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
    defer cancel()
    
    // Use context for cancellation
    select {
    case result := <-doWork(ctx):
        return result, nil
    case <-ctx.Done():
        return map[string]interface{}{
            "error": "operation timed out",
            "success": false,
        }, nil
    }
}
```

## Plugin Installation Process

### How Corynth Installs Remote Plugins

When you run `corynth plugin install <plugin-name>`, here's exactly what happens:

#### 1. Repository Configuration
Corynth automatically configures the official plugin repository:
```go
// Default repository added if none configured
{
    Name:     "official",
    URL:      "https://github.com/corynth/plugins", 
    Branch:   "main",
    Priority: 1,
}
```

#### 2. Plugin Discovery
Corynth searches for plugins in this order:
1. **Registry Lookup**: Fetches `registry.json` from GitHub
   - `https://raw.githubusercontent.com/corynth/plugins/main/registry.json`
   - `https://raw.githubusercontent.com/corynth/plugins/main/plugins/registry.json`

2. **Repository Clone**: Clones/updates the plugin repository to local cache
   - Location: `.corynth/cache/repos/official/`
   - Uses Git with automatic branch tracking

#### 3. Plugin Resolution Priority
Corynth looks for plugins in this specific order:

1. **JSON Protocol Plugins** (Current approach):
   ```bash
   # Plugin located at: official/<plugin-name>/plugin.go
   # Bash wrapper at: official/<plugin-name>/plugin
   ```

2. **Compiled .so Files** (Legacy):
   ```bash
   # Pre-compiled: corynth-plugin-<name>.so
   ```

3. **Source Go Files**:
   ```bash
   # Root level: <plugin-name>.go
   # Directory: <plugin-name>/plugin.go
   ```

#### 4. Installation Process

For our Go JSON protocol plugins:
```bash
# 1. Clone repository
git clone https://github.com/corynth/plugins.git ~/.corynth/cache/repos/official/

# 2. Copy plugin directory
cp -r official/<plugin-name>/ ~/.corynth/plugins/<plugin-name>/

# 3. Plugin is ready for execution via JSON protocol
# No compilation needed - uses `go run` via wrapper script
```

#### 5. Plugin Execution Model

Our plugins use the **JSON stdin/stdout protocol**:
```bash
# Metadata command
./official/http/plugin metadata

# Action execution  
echo '{"url": "https://api.github.com/user"}' | ./official/http/plugin get
```

### Installation Commands

#### Install Plugin
```bash
# Install from official repository
corynth plugin install http

# Install with auto-install enabled (default)  
# Plugin will be automatically installed when first used
```

#### Discovery and Information
```bash
# List all available plugins
corynth plugin discover

# Search plugins by name/tags
corynth plugin search docker
corynth plugin search --tags cloud

# Get detailed plugin info
corynth plugin info kubernetes

# List installed plugins
corynth plugin list
```

#### Plugin Management
```bash
# Update plugin to latest version
corynth plugin update terraform

# Remove installed plugin  
corynth plugin remove ansible

# Show plugin categories
corynth plugin categories
```

### Configuration Options

#### Custom Repository Configuration
Add to `corynth.hcl`:
```hcl
plugins {
  auto_install = true
  
  repository "official" {
    url      = "https://github.com/corynth/plugins"
    branch   = "main" 
    priority = 1
  }
  
  repository "custom" {
    url      = "https://github.com/yourorg/corynth-plugins"
    branch   = "main"
    token    = "your-github-token"  # For private repos
    priority = 2
  }
}
```

#### Registry Format
Our registry.json follows this structure:
```json
{
  "version": "2.0.0",
  "updated": "2025-01-01T00:00:00Z", 
  "plugins": [
    {
      "name": "http",
      "version": "1.0.0",
      "description": "HTTP client for REST API calls",
      "format": "json-protocol",
      "language": "go",
      "actions": ["get", "post"],
      "requirements": {
        "corynth": ">=1.2.0",
        "runtime": "go"
      }
    }
  ]
}
```

### Troubleshooting Installation

#### Common Issues

1. **Plugin Not Found**:
   ```bash
   # Check available plugins
   corynth plugin discover
   
   # Verify repository access
   curl -s https://raw.githubusercontent.com/corynth/plugins/main/registry.json
   ```

2. **Installation Failed**:
   ```bash
   # Clear plugin cache
   rm -rf ~/.corynth/cache/
   
   # Reinstall with verbose output
   corynth plugin install <plugin-name> --verbose
   ```

3. **Permission Issues**:
   ```bash
   # Ensure proper permissions
   chmod +x ~/.corynth/plugins/<plugin-name>/plugin
   ```

### Plugin Directory Structure After Installation
```
~/.corynth/
├── cache/
│   └── repos/
│       └── official/           # Cloned repository
│           ├── registry.json
│           └── official/
│               ├── http/
│               ├── docker/
│               └── ...
└── plugins/                    # Installed plugins (if needed)
    └── <plugin-name>/
```

### Development Testing

Test your plugin installation process:
```bash
# 1. Test metadata access
curl -s https://raw.githubusercontent.com/corynth/plugins/main/registry.json | jq '.plugins[] | select(.name=="your-plugin")'

# 2. Test local installation
cd ~/.corynth/cache/repos/official/official/your-plugin
./plugin metadata
./plugin actions

# 3. Test via Corynth
corynth plugin install your-plugin
corynth plugin info your-plugin
```

## Support and Contribution

- **Issues**: Report bugs via GitHub issues
- **Discussions**: Use GitHub discussions for questions  
- **Contributing**: Follow the contribution guidelines
- **Bootstrap**: Use the automated bootstrap for consistency
- **Installation**: Reference this guide for plugin deployment

The Go-first approach ensures that Corynth plugins are fast, secure, and maintainable while providing a consistent development experience across all platforms. The JSON protocol enables language-agnostic execution while maintaining the performance benefits of Go implementation.
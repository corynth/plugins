# Corynth Plugins - Go-First JSON Protocol

This is the official plugin repository for Corynth workflow orchestration. All plugins use Go implementations with JSON stdin/stdout communication for maximum performance, security, and cross-platform compatibility.

## 🚀 Quick Start

### Install Plugins
```bash
# Install a plugin from the official repository
corynth plugin install http

# Auto-install during workflow execution (default)
corynth apply workflow.hcl  # Plugins installed automatically when needed
```

### Use in Workflows
```hcl
step "api_call" {
  plugin = "http"
  action = "get" 
  params = {
    url = "https://api.github.com/user"
    headers = {
      "Authorization" = "Bearer ${var.github_token}"
    }
  }
}
```

## 📦 Available Plugins (14 Total)

### 🌐 Infrastructure & Cloud
- **docker** - Container operations and image management  
- **kubernetes** - Cluster management and resource operations
- **terraform** - Infrastructure as Code operations
- **ansible** - Configuration management and automation
- **aws** - Amazon Web Services operations (EC2, S3, Lambda)

### 💬 Communication & Notifications  
- **slack** - Workspace messaging and notifications
- **email** - SMTP email notifications
- **http** - REST API calls and web requests

### 💾 Data & Storage
- **sql** - Database operations (SQLite, PostgreSQL, MySQL)  
- **file** - File system operations (read, write, copy, move)

### 🤖 AI & Analytics
- **llm** - Large Language Model integration (OpenAI, Ollama)
- **reporting** - Report generation (Markdown, HTML, text)

### 🛠️ System & Utilities
- **shell** - Command execution with various interpreters  
- **calculator** - Mathematical calculations with AST parsing

## 🔧 Plugin Management

### Discovery Commands
```bash
# List all available plugins 
corynth plugin discover

# Search plugins by name or tags
corynth plugin search docker
corynth plugin search --tags cloud

# Get detailed plugin information
corynth plugin info kubernetes

# Show plugin categories
corynth plugin categories
```

### Installation & Updates
```bash
# Install specific plugin
corynth plugin install terraform

# Update plugin to latest version  
corynth plugin update ansible

# Remove installed plugin
corynth plugin remove slack

# List installed plugins
corynth plugin list
```

## 🏗️ Development

### Bootstrap New Plugin
```bash
# Generate complete plugin structure
./official/bootstrap-plugin.sh my-new-plugin

# Creates:
# my-new-plugin/
# ├── plugin.go      # Full Go implementation  
# ├── plugin         # Executable wrapper
# ├── go.mod         # Module configuration
# └── samples/       # Example workflows
```

### Plugin Architecture

All plugins follow the **Go-first JSON protocol**:

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
)

type YourPlugin struct{}

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
            },
            Outputs: map[string]IOSpec{
                "result": {Type: "string", Description: "Result description"},
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
```

### JSON Protocol Communication

Plugins communicate via JSON stdin/stdout:

```bash
# Get plugin metadata
./plugin metadata

# Get available actions  
./plugin actions

# Execute action with parameters
echo '{"param1": "value"}' | ./plugin action_name
```

## 📁 Repository Structure

```
corynth/plugins/
├── official/              # Official plugins directory
│   ├── http/             # HTTP client plugin
│   │   ├── plugin.go     # Go implementation  
│   │   └── plugin        # Bash wrapper script
│   ├── docker/           # Docker operations
│   ├── kubernetes/       # K8s management
│   └── ...              # All 14 plugins
├── registry.json         # Plugin registry metadata
├── PLUGIN_DEVELOPMENT.md # Complete development guide
└── README.md            # This file
```

## 🔄 Installation Process

When you run `corynth plugin install <plugin-name>`:

1. **Registry Lookup**: Fetches plugin metadata from `registry.json`
2. **Repository Clone**: Downloads plugin source to local cache  
3. **Plugin Resolution**: Locates Go source and wrapper scripts
4. **Ready to Execute**: Plugin available via JSON protocol

No compilation needed - plugins use `go run` for maximum flexibility.

## 🚀 Plugin Features

### ✅ Production Ready
- **No Placeholders**: All plugins fully implemented and tested
- **Error Handling**: Comprehensive error responses in JSON format  
- **Input Validation**: Type checking and parameter validation
- **Security**: Safe execution with input sanitization

### ⚡ High Performance
- **Go Implementation**: Native performance and fast startup
- **JSON Protocol**: Efficient stdin/stdout communication
- **No Dependencies**: Self-contained with minimal external requirements
- **Concurrent Safe**: Support for parallel execution

### 🌍 Cross-Platform
- **Universal**: Works on Linux, macOS, and Windows
- **Standard Interface**: Consistent JSON protocol across all plugins
- **Portable**: Single repository deployment model

## 🎯 Example Usage

### HTTP API Integration
```hcl
step "fetch_data" {
  plugin = "http"
  action = "get"
  params = {
    url = "https://api.service.com/data"
    headers = {
      "Authorization" = "Bearer ${var.api_token}"
      "Content-Type" = "application/json"
    }
    timeout = 30
  }
}
```

### Docker Container Management  
```hcl
step "deploy_app" {
  plugin = "docker" 
  action = "run"
  params = {
    image = "nginx:latest"
    ports = ["80:8080"]
    env = {
      "ENVIRONMENT" = "production"
    }
    volumes = ["/data:/usr/share/nginx/html"]
  }
}
```

### Kubernetes Deployment
```hcl
step "deploy_k8s" {
  plugin = "kubernetes"
  action = "apply" 
  params = {
    manifest = file("deployment.yaml")
    namespace = "production"
  }
}
```

## 📖 Documentation

- **[Development Guide](PLUGIN_DEVELOPMENT.md)**: Complete plugin development documentation
- **[Plugin Examples](official/)**: Browse all available plugins and their implementations
- **[Bootstrap Tool](official/bootstrap-plugin.sh)**: Automated plugin creation utility

## 🤝 Contributing

1. **Fork** this repository
2. **Create** a new plugin using the bootstrap tool  
3. **Implement** using the Go-first JSON protocol
4. **Test** thoroughly with real-world scenarios
5. **Submit** a pull request with comprehensive documentation

### Development Standards
- ✅ Go implementation with JSON protocol
- ✅ Comprehensive error handling  
- ✅ Input validation and sanitization
- ✅ Sample workflows included
- ✅ Production-ready (no placeholders)

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/corynth/plugins/issues)
- **Discussions**: [GitHub Discussions](https://github.com/corynth/plugins/discussions)  
- **Documentation**: [Plugin Development Guide](PLUGIN_DEVELOPMENT.md)

---

**🚀 Ready to build with Corynth?** Start with `corynth plugin install http` and explore the ecosystem!
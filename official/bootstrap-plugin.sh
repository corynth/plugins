#!/bin/bash

# Corynth Plugin Bootstrap Script
# Creates a new Go plugin from the official template

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }
print_info() { echo -e "${YELLOW}→${NC} $1"; }

# Check if plugin name is provided
if [ -z "$1" ]; then
    echo "Usage: ./bootstrap-plugin.sh <plugin-name> [author]"
    echo ""
    echo "Example: ./bootstrap-plugin.sh myawesome \"John Doe\""
    exit 1
fi

PLUGIN_NAME="$1"
AUTHOR="${2:-Corynth Team}"
PLUGIN_DIR="./$PLUGIN_NAME"

# Validate plugin name (alphanumeric and hyphen only)
if [[ ! "$PLUGIN_NAME" =~ ^[a-z0-9-]+$ ]]; then
    print_error "Plugin name must be lowercase alphanumeric with hyphens only"
    exit 1
fi

# Check if directory already exists
if [ -d "$PLUGIN_DIR" ]; then
    print_error "Directory $PLUGIN_DIR already exists"
    exit 1
fi

print_info "Creating new Corynth plugin: $PLUGIN_NAME"

# Create plugin directory
mkdir -p "$PLUGIN_DIR"
print_success "Created directory: $PLUGIN_DIR"

# Copy and customize template
cat > "$PLUGIN_DIR/plugin.go" << 'EOF'
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

type PLUGIN_STRUCT struct {
	// Add state here
}

func NewPLUGIN_STRUCT() *PLUGIN_STRUCT {
	return &PLUGIN_STRUCT{}
}

func (p *PLUGIN_STRUCT) GetMetadata() Metadata {
	return Metadata{
		Name:        "PLUGIN_NAME",
		Version:     "1.0.0",
		Description: "PLUGIN_NAME plugin for Corynth",
		Author:      "AUTHOR_NAME",
		Tags:        []string{"PLUGIN_NAME"},
	}
}

func (p *PLUGIN_STRUCT) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"hello": {
			Description: "Example hello action",
			Inputs: map[string]IOSpec{
				"name": {
					Type:        "string",
					Required:    false,
					Default:     "World",
					Description: "Name to greet",
				},
			},
			Outputs: map[string]IOSpec{
				"message": {
					Type:        "string",
					Description: "Greeting message",
				},
			},
		},
	}
}

func (p *PLUGIN_STRUCT) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "hello":
		return p.hello(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *PLUGIN_STRUCT) hello(params map[string]interface{}) (map[string]interface{}, error) {
	name := "World"
	if val, ok := params["name"].(string); ok && val != "" {
		name = val
	}
	
	return map[string]interface{}{
		"message": fmt.Sprintf("Hello, %s! From PLUGIN_NAME plugin.", name),
	}, nil
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewPLUGIN_STRUCT()
	
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
EOF

# Convert plugin name to proper case for struct name
STRUCT_NAME=$(echo "$PLUGIN_NAME" | sed 's/-/_/g' | sed 's/\b\(.\)/\u\1/g')Plugin

# Replace placeholders
sed -i '' "s/PLUGIN_NAME/$PLUGIN_NAME/g" "$PLUGIN_DIR/plugin.go"
sed -i '' "s/PLUGIN_STRUCT/${STRUCT_NAME}/g" "$PLUGIN_DIR/plugin.go"
sed -i '' "s/AUTHOR_NAME/$AUTHOR/g" "$PLUGIN_DIR/plugin.go"

print_success "Created plugin.go from template"

# Create the compiled plugin script
cat > "$PLUGIN_DIR/plugin" << EOF
#!/usr/bin/env bash
# Auto-generated Corynth plugin runner
DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
exec go run "\$DIR/plugin.go" "\$@"
EOF

chmod +x "$PLUGIN_DIR/plugin"
print_success "Created executable plugin wrapper"

# Create go.mod for the plugin
cat > "$PLUGIN_DIR/go.mod" << EOF
module github.com/corynth/plugins/$PLUGIN_NAME

go 1.21

// Add your dependencies here
EOF

print_success "Created go.mod file"

# Create README
cat > "$PLUGIN_DIR/README.md" << EOF
# $PLUGIN_NAME Plugin

A Corynth workflow orchestration plugin.

## Description

$PLUGIN_NAME plugin for Corynth, created by $AUTHOR.

## Actions

### hello
Example action that returns a greeting message.

**Inputs:**
- \`name\` (string, optional): Name to greet (default: "World")

**Outputs:**
- \`message\` (string): Greeting message

## Usage

\`\`\`hcl
step "greet" {
  plugin = "$PLUGIN_NAME"
  action = "hello"
  params = {
    name = "Corynth"
  }
}
\`\`\`

## Development

1. Modify \`plugin.go\` to add your actions
2. Test locally:
   \`\`\`bash
   ./plugin metadata
   ./plugin actions
   echo '{"name":"Test"}' | ./plugin hello
   \`\`\`

3. Build for distribution:
   \`\`\`bash
   go build -o plugin plugin.go
   \`\`\`

## Requirements

- Go 1.21+
- Any additional dependencies in go.mod

## License

Apache 2.0
EOF

print_success "Created README.md"

# Create .gitignore
cat > "$PLUGIN_DIR/.gitignore" << EOF
# Compiled binary
/plugin
!plugin.go

# Go build artifacts
*.exe
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Dependency directories
vendor/

# Go workspace
go.work
EOF

print_success "Created .gitignore"

# Test the plugin
print_info "Testing plugin..."
cd "$PLUGIN_DIR"

if go run plugin.go metadata > /dev/null 2>&1; then
    print_success "Plugin metadata test passed"
else
    print_error "Plugin metadata test failed"
fi

if go run plugin.go actions > /dev/null 2>&1; then
    print_success "Plugin actions test passed"
else
    print_error "Plugin actions test failed"
fi

if echo '{"name":"Bootstrap"}' | go run plugin.go hello > /dev/null 2>&1; then
    print_success "Plugin execution test passed"
else
    print_error "Plugin execution test failed"
fi

echo ""
print_success "Plugin '$PLUGIN_NAME' created successfully!"
echo ""
echo "Next steps:"
echo "  1. cd $PLUGIN_DIR"
echo "  2. Edit plugin.go to implement your actions"
echo "  3. Test: ./plugin metadata"
echo "  4. Test: echo '{\"param\":\"value\"}' | ./plugin your-action"
echo "  5. Build: go build -o plugin plugin.go"
echo ""
echo "Happy plugin development! 🚀"
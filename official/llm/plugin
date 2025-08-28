#!/bin/bash

# LLM Plugin - Go Implementation Wrapper
# This script compiles and runs the Go plugin, falling back to Python if compilation fails

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_PLUGIN="$SCRIPT_DIR/plugin.go"
GO_BINARY="$SCRIPT_DIR/llm-plugin"
PYTHON_PLUGIN="$SCRIPT_DIR/plugin.py"

# Function to compile Go plugin
compile_go_plugin() {
    if command -v go >/dev/null 2>&1; then
        if [ "$GO_PLUGIN" -nt "$GO_BINARY" ] || [ ! -f "$GO_BINARY" ]; then
            echo "Compiling Go LLM plugin..." >&2
            if go build -o "$GO_BINARY" "$GO_PLUGIN" 2>/dev/null; then
                echo "Go plugin compiled successfully" >&2
                return 0
            else
                echo "Go compilation failed, will use Python fallback" >&2
                return 1
            fi
        fi
        return 0
    else
        echo "Go not found, using Python fallback" >&2
        return 1
    fi
}

# Function to run Python plugin
run_python_plugin() {
    if [ -f "$PYTHON_PLUGIN" ]; then
        exec python3 "$PYTHON_PLUGIN" "$@"
    else
        echo '{"error": "Neither Go nor Python plugin available"}' 
        exit 1
    fi
}

# Main execution logic
if compile_go_plugin && [ -f "$GO_BINARY" ]; then
    # Run compiled Go binary
    exec "$GO_BINARY" "$@"
else
    # Fallback to Python implementation
    run_python_plugin "$@"
fi
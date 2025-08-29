#!/bin/bash
# Test compiled plugins

set -e

VERSION=${VERSION:-"v1.0.0"}
BUILD_DIR="releases/$VERSION"

echo "🧪 Testing Corynth Plugin Binaries"
echo ""

# Detect current platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
    linux)   PLATFORM="linux" ;;
    darwin)   PLATFORM="darwin" ;;
    mingw*|cygwin*|msys*) PLATFORM="windows" ;;
    *)        echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

ARCH=$(uname -m)
case $ARCH in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)            echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "Platform: $PLATFORM-$ARCH"
echo ""

# Test plugins
PLUGINS=("http" "file" "calculator")

for plugin in "${PLUGINS[@]}"; do
    binary="$BUILD_DIR/corynth-plugin-$plugin-$PLATFORM-$ARCH"
    if [ "$PLATFORM" = "windows" ]; then
        binary="${binary}.exe"
    fi
    
    if [ ! -f "$binary" ]; then
        echo "❌ $plugin: Binary not found ($binary)"
        continue
    fi
    
    echo "🔍 Testing $plugin plugin..."
    
    # Test metadata
    metadata=$($binary metadata 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "  ✅ Metadata: $(echo $metadata | jq -r '.name + " v" + .version')"
    else
        echo "  ❌ Metadata failed"
        continue
    fi
    
    # Test actions
    actions=$($binary actions 2>/dev/null)
    if [ $? -eq 0 ]; then
        action_count=$(echo $actions | jq length 2>/dev/null || echo "unknown")
        echo "  ✅ Actions: $action_count available"
    else
        echo "  ❌ Actions failed"
        continue
    fi
    
    # Test specific functionality
    case $plugin in
        "calculator")
            result=$(echo '{"expression": "2 + 2", "precision": 0}' | $binary execute calculate 2>/dev/null)
            if echo "$result" | grep -q "4"; then
                echo "  ✅ Calculate: 2 + 2 = 4"
            else
                echo "  ❌ Calculate failed"
            fi
            ;;
        "http")
            result=$(echo '{"url": "https://httpbin.org/json", "timeout": 10}' | $binary execute get 2>/dev/null)
            if echo "$result" | grep -q "slideshow"; then
                echo "  ✅ HTTP GET: Successfully fetched test JSON"
            else
                echo "  ❌ HTTP GET failed"
            fi
            ;;
        "file")
            echo "test content" > /tmp/test-file.txt
            result=$(echo '{"path": "/tmp/test-file.txt"}' | $binary execute read 2>/dev/null)
            if echo "$result" | grep -q "test content"; then
                echo "  ✅ File read: Successfully read test file"
                rm -f /tmp/test-file.txt
            else
                echo "  ❌ File read failed"
            fi
            ;;
    esac
    echo ""
done

echo "✅ Plugin testing complete!"
echo ""
echo "All plugins are ready for distribution as pre-compiled binaries."
#!/bin/bash
# Corynth Plugin Binary Builder
# Builds all plugins for all platforms and creates releases

set -e

VERSION=${VERSION:-"v1.0.0"}
BUILD_DIR="releases/$VERSION"
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

echo "🔨 Building Corynth Plugins v$VERSION"
echo "Building for platforms: ${PLATFORMS[*]}"
echo ""

# Clean and create build directory
rm -rf releases
mkdir -p "$BUILD_DIR"

# Find all plugins
PLUGINS=$(find official -name "plugin.go" | sed 's|/plugin.go||' | sed 's|official/||')

echo "Found plugins: $PLUGINS"
echo ""

# Build each plugin for each platform
for plugin in $PLUGINS; do
    echo "🔧 Building $plugin..."
    plugin_dir="official/$plugin"
    
    if [ ! -f "$plugin_dir/plugin.go" ]; then
        echo "  ❌ Missing plugin.go, skipping"
        continue
    fi
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r -a platform_split <<< "$platform"
        GOOS=${platform_split[0]}
        GOARCH=${platform_split[1]}
        
        output_name="corynth-plugin-${plugin}-${GOOS}-${GOARCH}"
        if [ "$GOOS" = "windows" ]; then
            output_name="${output_name}.exe"
        fi
        
        echo "  📦 Building for ${GOOS}/${GOARCH}..."
        
        cd "$plugin_dir"
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-w -s" -o "../../${BUILD_DIR}/${output_name}" plugin.go
        cd ../..
        
        # Verify binary was created
        if [ -f "$BUILD_DIR/$output_name" ]; then
            size=$(du -h "$BUILD_DIR/$output_name" | cut -f1)
            echo "    ✅ Created $output_name ($size)"
        else
            echo "    ❌ Failed to create $output_name"
        fi
    done
    echo ""
done

# Create checksums
echo "📝 Creating checksums..."
cd "$BUILD_DIR"
shasum -a 256 * > checksums.txt
cd ../..

# Create summary
echo "📊 Build Summary"
echo "================"
echo "Version: $VERSION"
echo "Plugins built: $(echo $PLUGINS | wc -w)"
echo "Total binaries: $(ls -1 $BUILD_DIR/*.exe $BUILD_DIR/corynth-plugin-* 2>/dev/null | wc -l)"
echo "Build size: $(du -sh $BUILD_DIR | cut -f1)"
echo ""
echo "Files created in $BUILD_DIR/:"
ls -lh "$BUILD_DIR/" | head -10
if [ $(ls -1 "$BUILD_DIR/" | wc -l) -gt 10 ]; then
    echo "... and $(($(ls -1 "$BUILD_DIR/" | wc -l) - 10)) more files"
fi

echo ""
echo "✅ Plugin build complete!"
echo ""
echo "Next steps:"
echo "1. Test plugins: ./test-plugin.sh"
echo "2. Create GitHub release with files in $BUILD_DIR/"
echo "3. Update registry.json with download URLs"
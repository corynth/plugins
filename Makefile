# Corynth Plugins Production Makefile
# Builds, tests, and releases production-ready plugin binaries

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build configuration
LDFLAGS := -ldflags "-w -s -X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.Commit=$(COMMIT)"
BUILD_DIR := releases/$(VERSION)
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Plugin discovery
PLUGINS := $(shell find official -name "plugin.go" | sed 's|official/||' | sed 's|/plugin.go||')

.PHONY: all build test clean release install help security-scan

## Build all plugins for all platforms
all: clean build test

## Build plugin binaries for all platforms
build:
	@echo "🔨 Building Corynth Plugins $(VERSION)"
	@echo "Plugins: $(PLUGINS)"
	@echo "Platforms: $(PLATFORMS)"
	@mkdir -p $(BUILD_DIR)
	
	@for plugin in $(PLUGINS); do \
		echo "Building $$plugin..."; \
		if [ ! -f "official/$$plugin/plugin.go" ]; then \
			echo "  ❌ Missing plugin.go for $$plugin"; \
			continue; \
		fi; \
		\
		for platform in $(PLATFORMS); do \
			os=$$(echo $$platform | cut -d'/' -f1); \
			arch=$$(echo $$platform | cut -d'/' -f2); \
			output="corynth-plugin-$$plugin-$$os-$$arch"; \
			if [ "$$os" = "windows" ]; then \
				output="$$output.exe"; \
			fi; \
			\
			echo "  📦 Building for $$os/$$arch..."; \
			cd official/$$plugin && \
			GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o "../../$(BUILD_DIR)/$$output" plugin.go && \
			cd ../..; \
			\
			if [ -f "$(BUILD_DIR)/$$output" ]; then \
				size=$$(du -h "$(BUILD_DIR)/$$output" | cut -f1); \
				echo "    ✅ $$output ($$size)"; \
			else \
				echo "    ❌ Failed to build $$output"; \
				exit 1; \
			fi; \
		done; \
		echo ""; \
	done
	
	@echo "📝 Creating checksums and manifest..."
	@cd $(BUILD_DIR) && \
		echo "$(VERSION)" > version.txt && \
		echo "$(BUILD_DATE)" > build-date.txt && \
		shasum -a 256 corynth-plugin-* version.txt build-date.txt > checksums.txt
	
	@echo "✅ Build complete: $(BUILD_DIR)"

## Test all built plugins
test: build
	@echo "🧪 Testing plugin binaries..."
	@./test-plugins-production.sh $(VERSION)

## Security scan all binaries  
security-scan: build
	@echo "🔒 Running security scans..."
	@for binary in $(BUILD_DIR)/corynth-plugin-*; do \
		if [ -x "$$binary" ]; then \
			echo "Scanning $$(basename $$binary)..."; \
			file "$$binary" | grep -q "executable" || echo "  ❌ Not executable"; \
			strings "$$binary" | grep -q "corynth" && echo "  ✅ Contains expected strings"; \
		fi; \
	done

## Create GitHub release
release: build test security-scan
	@echo "📦 Creating production release..."
	@if [ ! -f "$(BUILD_DIR)/checksums.txt" ]; then \
		echo "❌ Checksums missing"; \
		exit 1; \
	fi
	
	@echo "Release artifacts:"
	@ls -lah $(BUILD_DIR)/ | head -10
	@echo ""
	@echo "🎯 Ready for deployment!"
	@echo "📋 Next steps:"
	@echo "  1. Create GitHub release: gh release create $(VERSION) $(BUILD_DIR)/*"
	@echo "  2. Update registry.json with new version"
	@echo "  3. Test end-to-end plugin installation"

## Install plugins locally for testing
install: build
	@echo "📥 Installing plugins locally..."
	@platform=$$(uname -s | tr '[:upper:]' '[:lower:]')
	@arch=$$(uname -m)
	@case $$arch in \
		x86_64) arch="amd64" ;; \
		aarch64|arm64) arch="arm64" ;; \
		*) echo "❌ Unsupported architecture: $$arch"; exit 1 ;; \
	esac
	
	@mkdir -p ~/.corynth/plugins
	@for plugin in $(PLUGINS); do \
		binary="$(BUILD_DIR)/corynth-plugin-$$plugin-$$platform-$$arch"; \
		if [ -f "$$binary" ]; then \
			cp "$$binary" "~/.corynth/plugins/corynth-plugin-$$plugin"; \
			chmod +x "~/.corynth/plugins/corynth-plugin-$$plugin"; \
			echo "  ✅ Installed $$plugin"; \
		fi; \
	done

## Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf releases/
	@echo "✅ Clean complete"

## Show build statistics
stats: build
	@echo "📊 Build Statistics"
	@echo "=================="
	@echo "Version: $(VERSION)"
	@echo "Plugins: $$(echo '$(PLUGINS)' | wc -w)"
	@echo "Platforms: $$(echo '$(PLATFORMS)' | wc -w)" 
	@echo "Total binaries: $$(ls -1 $(BUILD_DIR)/corynth-plugin-* | wc -l)"
	@echo "Build size: $$(du -sh $(BUILD_DIR) | cut -f1)"
	@echo ""
	@echo "Plugin sizes:"
	@for plugin in $(PLUGINS); do \
		echo -n "  $$plugin: "; \
		ls -lah $(BUILD_DIR)/corynth-plugin-$$plugin-* | head -1 | awk '{print $$5}'; \
	done

## Development mode - build and test single plugin
dev:
	@if [ -z "$(PLUGIN)" ]; then \
		echo "❌ Usage: make dev PLUGIN=http"; \
		echo "Available plugins: $(PLUGINS)"; \
		exit 1; \
	fi
	
	@echo "🚀 Development build for $(PLUGIN)..."
	@cd official/$(PLUGIN) && go build -o plugin plugin.go
	@echo "✅ Built: official/$(PLUGIN)/plugin"
	@echo "🧪 Testing..."
	@cd official/$(PLUGIN) && ./plugin metadata && ./plugin actions
	@echo "✅ $(PLUGIN) plugin working"

## Show help
help:
	@echo "Corynth Plugins Production Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make              # Build all plugins"
	@echo "  make test         # Build and test"
	@echo "  make release      # Full production release"
	@echo "  make dev PLUGIN=http  # Quick dev build"
	@echo "  make install      # Install locally"

# Default target
.DEFAULT_GOAL := help
# LLM Plugin - Go Implementation

A high-performance Go implementation of the LLM plugin for Large Language Model integration, supporting both OpenAI API and local Ollama models.

## Features

- **Native Go HTTP Client**: Uses Go's standard `net/http` package for optimal performance
- **OpenAI Integration**: Full support for OpenAI GPT models via REST API
- **Local Ollama Support**: Connect to local Ollama instances for privacy-focused LLM usage
- **Automatic Compilation**: Wrapper script automatically compiles Go code and falls back to Python if needed
- **Type Safety**: Strong typing for better reliability and performance
- **Error Handling**: Comprehensive error handling with detailed error messages

## Architecture

### Files Structure
```
llm/
├── plugin.go          # Go implementation (main plugin code)
├── plugin             # Bash wrapper script (auto-compiles and runs)
├── plugin.py          # Python fallback (backup of original)
├── llm-plugin         # Compiled Go binary (auto-generated)
└── samples/
    └── test-generate.sh # Test script demonstrating all actions
```

### Plugin Interface
The Go plugin implements the same interface as the Python version:
- `metadata` - Returns plugin information
- `actions` - Lists available actions with parameters
- `generate` - Text generation using OpenAI models
- `chat` - Conversational chat with message history
- `ollama` - Local Ollama model integration

## Usage

### Environment Variables
- `OPENAI_API_KEY` - Required for OpenAI actions (generate, chat)
- `OLLAMA_URL` - Optional, defaults to `http://localhost:11434`

### Basic Commands

#### Get Plugin Information
```bash
./plugin metadata
./plugin actions
```

#### Generate Text (OpenAI)
```bash
echo '{"prompt": "Write a haiku about programming", "max_tokens": 100, "temperature": 0.7}' | ./plugin generate
```

#### Chat Conversation (OpenAI)
```bash
echo '{"messages": [{"role": "user", "content": "What is Go?"}], "model": "gpt-3.5-turbo"}' | ./plugin chat
```

#### Local Ollama Model
```bash
echo '{"prompt": "Explain microservices", "model": "llama2"}' | ./plugin ollama
```

### Supported Parameters

#### Generate Action
- `prompt` (string, required) - Input prompt
- `model` (string, optional) - Model name (default: "gpt-3.5-turbo")
- `max_tokens` (number, optional) - Maximum tokens (default: 150)
- `temperature` (number, optional) - Creativity level (default: 0.7)

#### Chat Action
- `messages` (array, required) - Message history with role/content structure
- `model` (string, optional) - Model name (default: "gpt-3.5-turbo")

#### Ollama Action
- `prompt` (string, required) - Input prompt
- `model` (string, optional) - Ollama model name (default: "llama2")

## Performance Benefits

### Go vs Python
- **Faster HTTP Requests**: Native HTTP client with connection pooling
- **Lower Memory Usage**: Compiled binary with optimized memory management
- **Better Error Handling**: Type-safe error handling with detailed messages
- **Concurrent Safety**: Built-in goroutine safety for concurrent requests
- **No Runtime Dependencies**: Self-contained binary, no Python interpreter needed

### Automatic Optimization
- **Auto-compilation**: Wrapper script compiles Go source only when changed
- **Fallback System**: Gracefully falls back to Python if Go compilation fails
- **Connection Pooling**: HTTP client reuses connections for better performance

## Error Handling

The Go implementation provides detailed error messages for common issues:
- Missing API keys
- Network connectivity problems
- Invalid JSON input
- API rate limiting
- Model availability issues

## Testing

Run the comprehensive test script:
```bash
cd samples
./test-generate.sh
```

This tests all plugin actions and demonstrates proper usage patterns.

## Development

### Building Manually
```bash
go build -o llm-plugin plugin.go
```

### Wrapper Script Behavior
1. Checks if Go is available
2. Compiles Go source if newer than binary
3. Runs compiled Go binary
4. Falls back to Python on compilation failure

## Compatibility

- **Go Version**: Requires Go 1.16+ (uses embed and other modern features)
- **API Compatibility**: 100% compatible with Python plugin interface
- **JSON Schema**: Identical input/output format
- **Environment Variables**: Same configuration as Python version
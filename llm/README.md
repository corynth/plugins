# LLM Plugin

Production-ready Large Language Model plugin supporting commercial APIs (OpenAI, Anthropic) and self-hosted Ollama with full temperature control and parameter customization.

## Features

- **Multi-Provider Support** - OpenAI, Anthropic Claude, and Ollama
- **Temperature Control** - Precise control over response randomness (0.0-1.0)
- **Token Management** - Configurable token limits and usage tracking
- **Chat History** - Maintain conversation context across interactions
- **Embeddings** - Generate text embeddings for similarity and search
- **Model Management** - List and download models (Ollama)
- **Custom Parameters** - Provider-specific options and configurations

## Prerequisites

- **For Ollama**: Ollama installed and running (`ollama serve`)
- **For OpenAI**: Valid OpenAI API key
- **For Anthropic**: Valid Anthropic API key

## Actions

### generate
Generate text using an LLM with custom parameters.

```hcl
step "generate_docs" {
  plugin = "llm"
  action = "generate"
  
  params = {
    prompt = "Write a technical overview of Kubernetes networking"
    provider = "openai"
    model = "gpt-4"
    temperature = 0.3
    max_tokens = 1500
    system_prompt = "You are a technical documentation expert. Write clear, accurate, and comprehensive explanations."
    api_key = "${vault_secret}"
  }
}
```

### chat
Interactive chat with conversation history.

```hcl
step "technical_discussion" {
  plugin = "llm"
  action = "chat"
  
  params = {
    messages = [
      {
        role = "user"
        content = "How does Kubernetes service discovery work?"
      },
      {
        role = "assistant" 
        content = "Kubernetes service discovery works through..."
      },
      {
        role = "user"
        content = "Can you explain DNS resolution in more detail?"
      }
    ]
    provider = "anthropic"
    model = "claude-3-sonnet-20240229"
    temperature = 0.4
    max_tokens = 2000
    api_key = "${anthropic_key}"
  }
}
```

### embed
Generate embeddings for text similarity and search.

```hcl
step "create_embeddings" {
  plugin = "llm"
  action = "embed"
  
  params = {
    text = "Kubernetes is a container orchestration platform"
    provider = "ollama"
    model = "nomic-embed-text"
    base_url = "http://localhost:11434"
  }
}
```

### list_models
List available models for a provider.

```hcl
step "check_models" {
  plugin = "llm"
  action = "list_models"
  
  params = {
    provider = "ollama"
    base_url = "http://localhost:11434"
  }
}
```

### pull_model
Download and install models (Ollama only).

```hcl
step "install_model" {
  plugin = "llm"
  action = "pull_model"
  
  params = {
    model = "llama2:13b"
    base_url = "http://localhost:11434"
  }
}
```

## Parameters

### Common Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `provider` | string | No | "ollama" | LLM provider (openai, anthropic, ollama) |
| `model` | string | No | provider-specific | Model name |
| `temperature` | number | No | 0.7 | Sampling temperature (0.0-1.0) |
| `max_tokens` | number | No | 1000 | Maximum tokens to generate |
| `api_key` | string | No | - | API key for commercial providers |
| `base_url` | string | No | provider-specific | Custom API base URL |

### Generation Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `prompt` | string | Yes | - | Input prompt for generation |
| `system_prompt` | string | No | - | System instructions |
| `context` | array | No | - | Previous conversation context |
| `options` | object | No | - | Provider-specific options |

### Chat Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `messages` | array | Yes | - | Chat message history |

### Embedding Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `text` | string | Yes | - | Text to embed |

### Model Management Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `model` | string | Yes | - | Model name to pull (pull_model) |

## Outputs

### generate/chat
- `response` - Generated text response
- `tokens_used` - Number of tokens consumed
- `model` - Model used for generation
- `provider` - Provider used
- `context` - Updated conversation context (generate only)
- `messages` - Updated chat history (chat only)

### embed
- `embedding` - Embedding vector array
- `dimensions` - Number of embedding dimensions
- `model` - Model used for embedding

### list_models
- `models` - Array of available model names
- `count` - Number of models available

### pull_model
- `status` - Pull operation status (success/failed)
- `model` - Model name that was pulled

## Provider Configuration

### OpenAI
```hcl
params = {
  provider = "openai"
  model = "gpt-4"  # or gpt-3.5-turbo, gpt-4-turbo
  api_key = "${openai_api_key}"
  base_url = "https://api.openai.com/v1"  # optional
  temperature = 0.7
  max_tokens = 1500
}
```

### Anthropic Claude
```hcl
params = {
  provider = "anthropic"
  model = "claude-3-sonnet-20240229"  # or claude-3-opus, claude-3-haiku
  api_key = "${anthropic_api_key}"
  base_url = "https://api.anthropic.com/v1"  # optional
  temperature = 0.5
  max_tokens = 2000
}
```

### Ollama (Self-hosted)
```hcl
params = {
  provider = "ollama"
  model = "llama2"  # or codellama, mistral, etc.
  base_url = "http://localhost:11434"  # default
  temperature = 0.8
  max_tokens = 1000
}
```

## Examples

### Code Generation Workflow
```hcl
workflow "ai-code-generation" {
  description = "Generate and review code using AI"
  
  step "generate_function" {
    plugin = "llm"
    action = "generate"
    
    params = {
      prompt = "Write a Python function that validates email addresses using regex"
      provider = "openai"
      model = "gpt-4"
      temperature = 0.2
      max_tokens = 500
      system_prompt = "You are an expert Python developer. Write clean, well-documented code with proper error handling."
      api_key = "${openai_key}"
    }
  }
  
  step = "review_code" {
    plugin = "llm"
    action = "generate"
    depends_on = ["generate_function"]
    
    params = {
      prompt = "Review this Python code for security and best practices:\n\n${generate_function.response}"
      provider = "anthropic"
      model = "claude-3-sonnet-20240229"
      temperature = 0.1
      max_tokens = 800
      system_prompt = "You are a senior code reviewer. Focus on security, performance, and maintainability."
      api_key = "${anthropic_key}"
    }
  }
  
  step = "save_results" {
    plugin = "file"
    action = "write"
    depends_on = ["review_code"]
    
    params = {
      path = "./generated/email_validator.py"
      content = "${generate_function.response}"
    }
  }
}
```

### Documentation Generation
```hcl
workflow "generate-documentation" {
  description = "Auto-generate technical documentation"
  
  step "analyze_codebase" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "find ./src -name '*.go' -exec head -20 {} \\;"
    }
  }
  
  step "generate_overview" {
    plugin = "llm"
    action = "generate"
    depends_on = ["analyze_codebase"]
    
    params = {
      prompt = "Based on this Go codebase structure, write a technical overview:\n\n${analyze_codebase.output}"
      provider = "ollama"
      model = "codellama:13b"
      temperature = 0.3
      max_tokens = 2000
      system_prompt = "You are a technical writer. Create clear, structured documentation."
      base_url = "http://localhost:11434"
    }
  }
  
  step "create_api_docs" {
    plugin = "llm"
    action = "chat"
    depends_on = ["generate_overview"]
    
    params = {
      messages = [
        {
          role = "user"
          content = "Now create API documentation for the main functions you identified"
        }
      ]
      provider = "ollama"
      model = "codellama:13b"
      temperature = 0.2
      max_tokens = 1500
    }
  }
  
  step "save_documentation" {
    plugin = "reporting"
    action = "create_report"
    depends_on = ["create_api_docs"]
    
    params = {
      title = "API Documentation"
      format = "markdown"
      output_file = "./docs/api-reference.md"
      sections = [
        {
          heading = "Overview"
          content = "${generate_overview.response}"
        },
        {
          heading = "API Reference"
          content = "${create_api_docs.response}"
        }
      ]
    }
  }
}
```

### Multi-Model Comparison
```hcl
workflow "model-comparison" {
  description = "Compare responses from multiple LLM providers"
  
  step "gpt_response" {
    plugin = "llm"
    action = "generate"
    
    params = {
      prompt = "Explain the benefits of microservices architecture"
      provider = "openai"
      model = "gpt-4"
      temperature = 0.5
      max_tokens = 800
      api_key = "${openai_key}"
    }
  }
  
  step "claude_response" {
    plugin = "llm"
    action = "generate"
    
    params = {
      prompt = "Explain the benefits of microservices architecture"
      provider = "anthropic"
      model = "claude-3-sonnet-20240229"
      temperature = 0.5
      max_tokens = 800
      api_key = "${anthropic_key}"
    }
  }
  
  step "ollama_response" {
    plugin = "llm"
    action = "generate"
    
    params = {
      prompt = "Explain the benefits of microservices architecture"
      provider = "ollama"
      model = "llama2:13b"
      temperature = 0.5
      max_tokens = 800
    }
  }
  
  step "comparison_report" {
    plugin = "reporting"
    action = "create_report"
    depends_on = ["gpt_response", "claude_response", "ollama_response"]
    
    params = {
      title = "LLM Model Comparison: Microservices"
      format = "markdown"
      output_file = "./analysis/model-comparison.md"
      sections = [
        {
          heading = "GPT-4 Response"
          content = "${gpt_response.response}"
        },
        {
          heading = "Claude-3 Sonnet Response"
          content = "${claude_response.response}"
        },
        {
          heading = "Llama2-13B Response"
          content = "${ollama_response.response}"
        }
      ]
    }
  }
}
```

### RAG (Retrieval-Augmented Generation)
```hcl
workflow "rag-search" {
  description = "Search documentation using embeddings and generate contextual responses"
  
  step "create_query_embedding" {
    plugin = "llm"
    action = "embed"
    
    params = {
      text = "How do I configure Kubernetes ingress controllers?"
      provider = "ollama"
      model = "nomic-embed-text"
    }
  }
  
  step "search_documentation" {
    plugin = "shell"
    action = "exec"
    depends_on = ["create_query_embedding"]
    
    params = {
      command = "grep -r 'ingress' ./docs/ | head -10"
    }
  }
  
  step "generate_contextual_answer" {
    plugin = "llm"
    action = "generate"
    depends_on = ["search_documentation"]
    
    params = {
      prompt = "Using this documentation context, answer the question about Kubernetes ingress controllers:\n\nContext:\n${search_documentation.output}\n\nQuestion: How do I configure Kubernetes ingress controllers?"
      provider = "anthropic"
      model = "claude-3-haiku-20240307"
      temperature = 0.3
      max_tokens = 1000
      system_prompt = "Answer based on the provided context. If the context doesn't contain enough information, say so."
      api_key = "${anthropic_key}"
    }
  }
}
```

## Temperature Guidelines

- **0.0-0.2**: Deterministic, factual responses (code generation, technical docs)
- **0.3-0.5**: Balanced creativity and consistency (explanations, tutorials)
- **0.6-0.8**: Creative but coherent (brainstorming, creative writing)
- **0.9-1.0**: Highly creative, diverse responses (creative tasks)

## Security Best Practices

- Store API keys in environment variables or secure vaults
- Use minimal token limits to control costs
- Validate and sanitize all inputs
- Monitor usage and implement rate limiting
- Use system prompts to constrain behavior

## Integration

Works seamlessly with other Corynth plugins:
- **reporting plugin** - Generate formatted AI responses
- **vault plugin** - Secure API key management
- **file plugin** - Save generated content
- **http plugin** - Webhook integrations with AI responses
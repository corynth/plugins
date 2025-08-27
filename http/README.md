# HTTP Plugin

## Overview

The HTTP Plugin provides comprehensive HTTP client functionality for REST API interactions, webhooks, and web service integration within Corynth workflows.

## Features

- Full HTTP method support (GET, POST, PUT, DELETE, PATCH)
- Custom headers and authentication
- JSON and form data handling
- Request/response timeout management
- Retry logic with configurable delays
- SSL/TLS configuration options

## Actions

### get
Performs HTTP GET request

**Parameters:**
- `url` (string, required): Target URL
- `headers` (map, optional): HTTP headers
- `timeout` (number, optional): Request timeout in seconds (default: 30)
- `retry` (number, optional): Number of retry attempts (default: 0)
- `retry_delay` (number, optional): Delay between retries in seconds (default: 1)

**Returns:**
- `status_code`: HTTP response status code
- `headers`: Response headers
- `body`: Response body
- `response_time`: Request duration in milliseconds

### post
Performs HTTP POST request

**Parameters:**
- `url` (string, required): Target URL
- `headers` (map, optional): HTTP headers
- `body` (map/string, optional): Request body (JSON object or string)
- `timeout` (number, optional): Request timeout in seconds
- `retry` (number, optional): Number of retry attempts
- `retry_delay` (number, optional): Delay between retries

### put
Performs HTTP PUT request (same parameters as POST)

### delete
Performs HTTP DELETE request

**Parameters:**
- `url` (string, required): Target URL
- `headers` (map, optional): HTTP headers
- `timeout` (number, optional): Request timeout in seconds

### patch
Performs HTTP PATCH request (same parameters as POST)

## Usage Examples

### Basic GET Request
```hcl
step "api_health_check" {
  plugin = "http"
  action = "get"
  params = {
    url = "https://api.example.com/health"
    headers = {
      "User-Agent" = "Corynth-Workflow/1.0"
    }
    timeout = 30
  }
}
```

### POST Request with JSON Body
```hcl
step "create_user" {
  plugin = "http"
  action = "post"
  params = {
    url = "https://api.example.com/users"
    headers = {
      "Content-Type"  = "application/json"
      "Authorization" = "Bearer ${var.api_token}"
    }
    body = {
      name  = "John Doe"
      email = "john@example.com"
      role  = "user"
    }
    timeout = 60
    retry   = 3
  }
}
```

### Webhook Call with Authentication
```hcl
step "send_webhook" {
  plugin = "http"
  action = "post"
  params = {
    url = var.webhook_url
    headers = {
      "Content-Type"     = "application/json"
      "X-Webhook-Secret" = var.webhook_secret
    }
    body = {
      event      = "deployment.completed"
      environment = var.environment
      version     = var.app_version
      timestamp   = timestamp()
    }
  }
}
```

### API Call with Retry Logic
```hcl
step "reliable_api_call" {
  plugin = "http"
  action = "get"
  params = {
    url         = "https://api.unreliable-service.com/data"
    timeout     = 120
    retry       = 5
    retry_delay = 10
    headers = {
      "Accept" = "application/json"
    }
  }
}
```

## Authentication Examples

### Bearer Token Authentication
```hcl
step "authenticated_request" {
  plugin = "http"
  action = "get"
  params = {
    url = "https://api.example.com/protected"
    headers = {
      "Authorization" = "Bearer ${var.access_token}"
    }
  }
}
```

### Basic Authentication
```hcl
step "basic_auth_request" {
  plugin = "http"
  action = "get"
  params = {
    url = "https://api.example.com/data"
    headers = {
      "Authorization" = "Basic ${base64encode('${var.username}:${var.password}')}"
    }
  }
}
```

### API Key Authentication
```hcl
step "api_key_request" {
  plugin = "http"
  action = "get"
  params = {
    url = "https://api.example.com/data?api_key=${var.api_key}"
    headers = {
      "X-API-Key" = var.api_key
    }
  }
}
```

## Error Handling

The HTTP plugin automatically handles common HTTP errors and provides detailed error information:

```hcl
step "handle_api_errors" {
  plugin = "http"
  action = "get"
  params = {
    url = "https://api.example.com/data"
  }
}

step "check_response" {
  plugin = "shell"
  action = "exec"
  depends_on = ["handle_api_errors"]
  condition = "${handle_api_errors.status_code >= 400}"
  params = {
    command = "echo 'API request failed with status: ${handle_api_errors.status_code}'"
  }
}
```

## Response Processing

Access response data in subsequent steps:

```hcl
step "get_user_data" {
  plugin = "http"
  action = "get"
  params = {
    url = "https://api.example.com/users/123"
  }
}

step "process_user_data" {
  plugin = "shell"
  action = "exec"
  depends_on = ["get_user_data"]
  params = {
    command = "echo 'User name: ${get_user_data.body.name}'"
  }
}
```

## Configuration

### Environment Variables
- `HTTP_TIMEOUT`: Default timeout for all requests (seconds)
- `HTTP_RETRY_ATTEMPTS`: Default number of retry attempts
- `HTTP_USER_AGENT`: Default User-Agent header

### Plugin Configuration
```hcl
# Configure default timeouts globally
variable "http_timeout" {
  type    = number
  default = 60
}

step "api_call" {
  plugin = "http"
  action = "get"
  params = {
    url     = "https://api.example.com/data"
    timeout = var.http_timeout
  }
}
```

## Best Practices

1. **Always set timeouts** to prevent hanging workflows
2. **Use retry logic** for unreliable services
3. **Handle authentication securely** using variables, not hardcoded values
4. **Validate responses** before using data in subsequent steps
5. **Use appropriate HTTP methods** (GET for retrieval, POST for creation, etc.)
6. **Set User-Agent headers** to identify your application
7. **Implement error handling** for non-200 status codes

## Troubleshooting

### Common Issues

**Connection Timeout**
```
Error: request failed: context deadline exceeded
```
- Increase timeout value
- Check network connectivity
- Verify URL accessibility

**Authentication Failure**
```
Error: HTTP 401 Unauthorized
```
- Verify authentication credentials
- Check token expiration
- Ensure correct authentication method

**Rate Limiting**
```
Error: HTTP 429 Too Many Requests
```
- Add retry logic with delays
- Implement exponential backoff
- Check API rate limits

### Debug Mode

Enable detailed logging:
```bash
export CORYNTH_DEBUG=true
corynth run workflow.hcl
```

## Sample Workflows

See the `/samples` directory for complete workflow examples:
- `api-health-check.hcl`: API monitoring and health checks
- `webhook-integration.hcl`: Webhook processing and notifications
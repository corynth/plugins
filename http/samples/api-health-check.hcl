workflow "api-health-check" {
  description = "Monitor API endpoints and validate responses"
  version     = "1.0.0"

  variable "api_base_url" {
    type        = string
    default     = "https://jsonplaceholder.typicode.com"
    description = "Base URL for API endpoints"
  }

  step "check_api_status" {
    plugin = "http"
    action = "get"
    
    params = {
      url = "${var.api_base_url}/posts/1"
      timeout = 10
      headers = {
        "Accept" = "application/json"
        "User-Agent" = "Corynth-HealthCheck/1.0"
      }
    }
  }

  step "validate_response" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["check_api_status"]
    
    params = {
      command = "echo 'API Status: ${check_api_status.status_code}' && echo 'Response received successfully'"
    }
  }

  step "test_api_post" {
    plugin = "http"
    action = "post"
    
    depends_on = ["validate_response"]
    
    params = {
      url = "${var.api_base_url}/posts"
      body = "{\"title\": \"Test Post\", \"body\": \"Created by Corynth HTTP Plugin\", \"userId\": 1}"
      headers = {
        "Content-Type" = "application/json"
        "Accept" = "application/json"
      }
      timeout = 15
    }
  }
}
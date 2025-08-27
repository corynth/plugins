workflow "webhook-integration" {
  description = "Send webhooks and process responses"
  version     = "1.0.0"

  variable "webhook_url" {
    type        = string
    description = "Webhook endpoint URL"
  }

  variable "payload_data" {
    type        = string
    default     = "{\"event\": \"deployment\", \"status\": \"success\", \"app\": \"my-app\"}"
    description = "JSON payload to send"
  }

  step "send_webhook" {
    plugin = "http"
    action = "post"
    
    params = {
      url = var.webhook_url
      body = var.payload_data
      headers = {
        "Content-Type" = "application/json"
        "X-Event-Type" = "deployment"
        "X-Source" = "corynth"
      }
      timeout = 30
    }
  }

  step "verify_delivery" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["send_webhook"]
    
    params = {
      command = "echo 'Webhook delivered with status: ${send_webhook.status_code}' && echo 'Response: ${send_webhook.body}'"
    }
  }

  step "log_webhook_result" {
    plugin = "file"
    action = "write"
    
    depends_on = ["verify_delivery"]
    
    params = {
      path = "/tmp/webhook-log.txt"
      content = "Webhook sent at $(date)\nURL: ${var.webhook_url}\nStatus: ${send_webhook.status_code}\nResponse: ${send_webhook.body}"
    }
  }
}
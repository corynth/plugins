workflow "deployment-notification" {
  description = "Send deployment notifications to Slack"
  version     = "1.0.0"

  variable "slack_webhook" {
    type        = string
    description = "Slack webhook URL for notifications"
  }

  variable "app_name" {
    type        = string
    default     = "my-application"
    description = "Application name being deployed"
  }

  variable "environment" {
    type        = string
    default     = "production"
    description = "Deployment environment"
  }

  step "notify_deployment_start" {
    plugin = "slack"
    action = "workflow_notification"
    
    params = {
      webhook_url = var.slack_webhook
      workflow_name = "${var.app_name} Deployment"
      status = "started"
      details = "Starting deployment to ${var.environment} environment"
      channel = "#deployments"
    }
  }

  step "simulate_deployment" {
    plugin = "shell"
    action = "script"
    
    depends_on = ["notify_deployment_start"]
    
    params = {
      script = <<-EOF
        #!/bin/bash
        echo "Deploying ${var.app_name} to ${var.environment}..."
        sleep 5
        echo "Deployment process completed"
      EOF
      timeout = 30
    }
  }

  step "notify_deployment_success" {
    plugin = "slack"
    action = "send_rich_message"
    
    depends_on = ["simulate_deployment"]
    
    params = {
      webhook_url = var.slack_webhook
      text = "${var.app_name} deployed successfully to ${var.environment}!"
      color = "good"
      fields = {
        "Environment" = var.environment
        "Application" = var.app_name
        "Status" = "Success"
        "Duration" = "5 seconds"
      }
      channel = "#deployments"
      username = "Corynth Deploy Bot"
    }
  }

  step "send_team_notification" {
    plugin = "slack"
    action = "send_message"
    
    depends_on = ["notify_deployment_success"]
    
    params = {
      webhook_url = var.slack_webhook
      text = ":rocket: ${var.app_name} is now live in ${var.environment}! :tada:"
      channel = "#general"
      username = "Corynth"
      icon_emoji = ":rocket:"
    }
  }
}
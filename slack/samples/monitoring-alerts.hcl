workflow "monitoring-alerts" {
  description = "Send monitoring alerts and status updates to Slack"
  version     = "1.0.0"

  variable "slack_webhook" {
    type        = string
    description = "Slack webhook URL for alerts"
  }

  variable "service_name" {
    type        = string
    default     = "web-service"
    description = "Service being monitored"
  }

  variable "alert_type" {
    type        = string
    default     = "warning"
    description = "Alert severity: warning, failed, success"
  }

  step "check_service_health" {
    plugin = "http"
    action = "get"
    
    params = {
      url = "http://localhost:8080/health"
      timeout = 5
    }
  }

  step "send_status_alert" {
    plugin = "slack"
    action = "workflow_notification"
    
    depends_on = ["check_service_health"]
    
    params = {
      webhook_url = var.slack_webhook
      workflow_name = "${var.service_name} Health Check"
      status = var.alert_type
      duration = "5s"
      details = "Service responded with status ${check_service_health.status_code}"
      channel = "#alerts"
    }
  }

  step "detailed_system_check" {
    plugin = "shell"
    action = "script"
    
    depends_on = ["send_status_alert"]
    
    params = {
      script = <<-EOF
        #!/bin/bash
        echo "=== System Health Check ==="
        echo "CPU Usage: $(top -l 1 | grep "CPU usage" || echo "N/A")"
        echo "Memory: $(free -h 2>/dev/null || vm_stat | head -5)"
        echo "Disk: $(df -h / | tail -1)"
        echo "Network: $(netstat -an | grep LISTEN | wc -l) listening ports"
      EOF
    }
  }

  step "send_detailed_report" {
    plugin = "slack"
    action = "send_rich_message"
    
    depends_on = ["detailed_system_check"]
    
    params = {
      webhook_url = var.slack_webhook
      text = "System Health Report for ${var.service_name}"
      color = "${var.alert_type == 'failed' ? 'danger' : (var.alert_type == 'warning' ? 'warning' : 'good')}"
      fields = {
        "Service" = var.service_name
        "Health Status" = "${check_service_health.status_code}"
        "Alert Level" = var.alert_type
        "System Info" = "${detailed_system_check.stdout}"
      }
      channel = "#monitoring"
      username = "Corynth Monitor"
    }
  }
}
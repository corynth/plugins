# Microsoft Teams Plugin

Microsoft Teams integration plugin for Corynth that enables sending messages, alerts, and notifications via Teams webhooks.

## Features

- Send text messages to Teams channels
- Send formatted alert notifications with color coding
- Support for webhook-based integration
- Professional formatting and styling
- Real Teams API integration

## Prerequisites

- Microsoft Teams workspace
- Incoming webhook connector configured in target channel
- Internet connection for API calls

## Installation

```bash
cd teams
go build -o corynth-plugin-teams main.go
```

## Actions

### send_message

Send a simple text message to Teams channel.

**Parameters:**
- `webhook_url` (required): Teams channel webhook URL
- `text` (required): Message text
- `title` (optional): Message title

**Example:**
```hcl
step "notify_team" {
  plugin = "teams"
  action = "send_message"
  params = {
    webhook_url = "https://outlook.office.com/webhook/..."
    title = "Deployment Update"
    text = "Application v1.2.3 deployed successfully"
  }
}
```

### send_alert

Send formatted alert with automatic color coding.

**Parameters:**
- `webhook_url` (required): Teams channel webhook URL
- `alert_type` (optional): info, warning, error, success (default: "info")
- `title` (required): Alert title
- `message` (required): Alert message

**Example:**
```hcl
step "error_alert" {
  plugin = "teams"
  action = "send_alert"
  params = {
    webhook_url = "https://outlook.office.com/webhook/..."
    alert_type = "error"
    title = "Database Connection Failed"
    message = "Unable to connect to primary database server"
  }
}
```

## Sample Workflows

### Deployment Notifications
```hcl
workflow "deployment_notifications" {
  description = "CI/CD pipeline with Teams notifications"
  
  step "start_notification" {
    plugin = "teams"
    action = "send_message"
    params = {
      webhook_url = "${teams_webhook}"
      title = "Pipeline Started"
      text = "CI/CD pipeline initiated for ${repository}"
    }
  }
  
  step "success_notification" {
    plugin = "teams"
    action = "send_alert"
    params = {
      webhook_url = "${teams_webhook}"
      alert_type = "success"
      title = "Deployment Successful"
      message = "Application deployed to ${environment}"
    }
  }
}
```

## Alert Types

- **info**: Blue theme for general information
- **warning**: Orange theme for warnings
- **error**: Red theme for critical issues  
- **success**: Green theme for successful operations

## Dependencies

- Standard Go HTTP client
- JSON encoding/decoding

## License

MIT License
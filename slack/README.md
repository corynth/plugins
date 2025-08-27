# Slack Plugin

A comprehensive Slack integration plugin for Corynth that enables sending messages, uploading files, creating channels, and managing workspace interactions.

## Features

- Send messages to channels and direct messages
- Upload files with custom metadata
- Create public and private channels
- Set user status messages
- Full Slack API integration
- Error handling and validation

## Prerequisites

- Slack workspace with bot permissions
- Slack bot token (starts with `xoxb-`)
- Internet connection for API calls

## Installation

Build and install the plugin:

```bash
cd slack
go build -o corynth-plugin-slack main.go
```

## Actions

### send_message

Send a message to any Slack channel or user.

**Parameters:**
- `token` (required): Slack bot token (xoxb-...)
- `channel` (required): Channel ID or name (#general, @username, C1234567890)
- `text` (required): Message text (supports Slack markdown)
- `username` (optional): Bot username override (default: "Corynth Bot")
- `icon_emoji` (optional): Bot icon emoji (default: ":robot_face:")

**Returns:**
- `ts`: Message timestamp ID
- `channel`: Channel ID where message was sent

**Example:**
```hcl
step "notify_deployment" {
  plugin = "slack"
  action = "send_message"
  params = {
    token = "xoxb-your-bot-token"
    channel = "#deployments"
    text = "Deployment completed successfully\nVersion: v1.2.3\nDuration: 45 seconds"
    username = "Deploy Bot"
    icon_emoji = ":rocket:"
  }
}
```

### send_file

Upload a file to Slack with optional metadata.

**Parameters:**
- `token` (required): Slack bot token
- `channels` (required): Comma-separated channel IDs
- `file_path` (required): Local file path to upload
- `filename` (optional): Custom filename for upload
- `title` (optional): File title
- `initial_comment` (optional): Comment to post with file

**Returns:**
- `file_id`: Uploaded file ID

**Example:**
```hcl
step "upload_report" {
  plugin = "slack"
  action = "send_file"
  params = {
    token = "xoxb-your-bot-token"
    channels = "C1234567890"
    file_path = "/tmp/deployment-report.pdf"
    title = "Weekly Deployment Report"
    initial_comment = "Here's this week's deployment summary"
  }
}
```

### create_channel

Create a new Slack channel.

**Parameters:**
- `token` (required): Slack admin token
- `name` (required): Channel name (lowercase, no spaces, no #)
- `is_private` (optional): Create as private channel (default: false)

**Returns:**
- `channel_id`: Created channel ID
- `channel_name`: Channel name

**Example:**
```hcl
step "create_project_channel" {
  plugin = "slack"
  action = "create_channel"
  params = {
    token = "xoxp-admin-token"
    name = "project-alpha"
    is_private = false
  }
}
```

### set_status

Set user status message and emoji.

**Parameters:**
- `token` (required): User token with users.profile:write scope
- `status_text` (required): Status message text
- `status_emoji` (optional): Status emoji (default: ":speech_balloon:")
- `status_expiration` (optional): Expiration timestamp (0 = no expiration)

**Example:**
```hcl
step "set_deployment_status" {
  plugin = "slack"
  action = "set_status"
  params = {
    token = "xoxp-user-token"
    status_text = "Deploying v1.2.3"
    status_emoji = ":construction:"
    status_expiration = 3600
  }
}
```

## Sample Workflows

### Deployment Notification Pipeline
```hcl
workflow "deployment_notifications" {
  description = "Comprehensive deployment notification system"
  
  step "notify_start" {
    plugin = "slack"
    action = "send_message"
    params = {
      token = "xoxb-your-bot-token"
      channel = "#deployments"
      text = "Starting deployment of application v${version}"
      icon_emoji = ":construction:"
    }
  }
  
  step "set_deployer_status" {
    plugin = "slack"
    action = "set_status"
    params = {
      token = "xoxp-user-token"
      status_text = "Deploying v${version}"
      status_emoji = ":rocket:"
    }
  }
  
  step "upload_logs" {
    plugin = "slack"
    action = "send_file"
    depends_on = ["notify_start"]
    params = {
      token = "xoxb-your-bot-token"
      channels = "#deployments"
      file_path = "/tmp/deployment.log"
      title = "Deployment Logs v${version}"
      initial_comment = "Detailed deployment logs attached"
    }
  }
  
  step "notify_completion" {
    plugin = "slack"
    action = "send_message"
    depends_on = ["upload_logs"]
    params = {
      token = "xoxb-your-bot-token"
      channel = "#deployments"
      text = "Deployment completed successfully\nVersion: v${version}\nDuration: ${duration}\nStatus: Healthy"
      icon_emoji = ":white_check_mark:"
    }
  }
}
```

### Incident Response Workflow
```hcl
workflow "incident_response" {
  description = "Automated incident response with Slack integration"
  
  step "create_incident_channel" {
    plugin = "slack"
    action = "create_channel"
    params = {
      token = "xoxp-admin-token"
      name = "incident-${incident_id}"
      is_private = false
    }
  }
  
  step "alert_oncall" {
    plugin = "slack"
    action = "send_message"
    depends_on = ["create_incident_channel"]
    params = {
      token = "xoxb-your-bot-token"
      channel = "@oncall-engineer"
      text = "INCIDENT ALERT\nSeverity: ${severity}\nDescription: ${description}\nChannel: #incident-${incident_id}"
      icon_emoji = ":rotating_light:"
    }
  }
  
  step "update_team_channel" {
    plugin = "slack"
    action = "send_message"
    depends_on = ["create_incident_channel"]
    params = {
      token = "xoxb-your-bot-token"  
      channel = "#ops-team"
      text = "Incident ${incident_id} reported\nSeverity: ${severity}\nResponse channel: <#${create_incident_channel.channel_id}>"
      icon_emoji = ":warning:"
    }
  }
}
```

## Token Configuration

### Bot Token (xoxb-)
Required for:
- Sending messages
- Uploading files
- Basic channel operations

**Scopes needed:**
- `chat:write`
- `files:write`
- `channels:read`

### User Token (xoxp-)
Required for:
- Creating channels (admin permissions)
- Setting user status
- Advanced workspace management

**Scopes needed:**
- `channels:manage` 
- `users.profile:write`

## Error Handling

The plugin provides detailed error messages for:
- Invalid token formats
- Missing required parameters
- Slack API errors (rate limits, permissions, etc.)
- Network connectivity issues
- File system errors

## Security

- Tokens are never logged or stored
- All API calls use HTTPS
- Supports workspace-level token restrictions
- Validates token format before API calls

## Dependencies

- Standard Go HTTP client
- JSON encoding/decoding
- No external dependencies required

## License

MIT License - See LICENSE file for details.
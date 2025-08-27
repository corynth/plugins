# Shell Plugin

## Overview

The Shell Plugin enables command execution and script running within Corynth workflows, providing powerful system integration capabilities for DevOps automation.

## Features

- Command execution with real-time output capture
- Script execution with multi-line support
- Environment variable management
- Working directory control
- Process timeout and cancellation
- Exit code handling and error management

## Actions

### exec
Executes a single command

**Parameters:**
- `command` (string, required): Command to execute
- `env` (map, optional): Environment variables
- `working_dir` (string, optional): Working directory for command execution
- `timeout` (number, optional): Command timeout in seconds (default: 300)
- `ignore_error` (bool, optional): Continue on non-zero exit codes (default: false)

**Returns:**
- `stdout`: Command standard output
- `stderr`: Command standard error
- `exit_code`: Process exit code
- `execution_time`: Command duration in milliseconds

### script
Executes a multi-line script

**Parameters:**
- `script` (string, required): Script content (supports heredoc syntax)
- `shell` (string, optional): Shell interpreter (default: "/bin/bash")
- `env` (map, optional): Environment variables
- `working_dir` (string, optional): Working directory
- `timeout` (number, optional): Script timeout in seconds

**Returns:**
- `stdout`: Script standard output
- `stderr`: Script standard error
- `exit_code`: Process exit code
- `execution_time`: Script duration in milliseconds

## Usage Examples

### Basic Command Execution
```hcl
step "list_files" {
  plugin = "shell"
  action = "exec"
  params = {
    command = "ls -la /tmp"
  }
}
```

### Command with Environment Variables
```hcl
step "deploy_with_config" {
  plugin = "shell"
  action = "exec"
  params = {
    command = "kubectl apply -f deployment.yaml"
    env = {
      KUBECONFIG  = "/path/to/kubeconfig"
      ENVIRONMENT = var.environment
      IMAGE_TAG   = var.app_version
    }
    working_dir = "/deployments"
    timeout     = 300
  }
}
```

### Multi-line Script Execution
```hcl
step "complex_deployment" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      set -e
      
      echo "Starting deployment process..."
      
      # Build application
      npm install
      npm run build
      
      # Run tests
      npm test
      
      # Deploy to production
      if [ "$ENVIRONMENT" = "production" ]; then
        echo "Deploying to production..."
        kubectl apply -f k8s/production/
      else
        echo "Deploying to staging..."
        kubectl apply -f k8s/staging/
      fi
      
      echo "Deployment completed successfully!"
    EOF
    env = {
      ENVIRONMENT = var.environment
      NODE_ENV    = "production"
    }
    working_dir = "/app"
    timeout     = 600
  }
}
```

### Conditional Error Handling
```hcl
step "optional_operation" {
  plugin = "shell"
  action = "exec"
  params = {
    command      = "risky-command --might-fail"
    ignore_error = true
  }
}

step "handle_failure" {
  plugin = "shell"
  action = "exec"
  depends_on = ["optional_operation"]
  condition  = "${optional_operation.exit_code != 0}"
  params = {
    command = "echo 'Operation failed, running cleanup...'"
  }
}
```

### System Information Gathering
```hcl
step "gather_system_info" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      echo "=== System Information ==="
      echo "Hostname: $(hostname)"
      echo "OS: $(uname -a)"
      echo "CPU: $(nproc) cores"
      echo "Memory: $(free -h | grep Mem | awk '{print $2}')"
      echo "Disk: $(df -h / | tail -1 | awk '{print $4}')"
      echo "Uptime: $(uptime)"
      echo "Docker: $(docker --version 2>/dev/null || echo 'Not installed')"
      echo "Kubernetes: $(kubectl version --client --short 2>/dev/null || echo 'Not installed')"
    EOF
  }
}
```

## Advanced Usage

### Database Operations
```hcl
step "database_backup" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      BACKUP_FILE="/backups/db_backup_$(date +%Y%m%d_%H%M%S).sql"
      
      echo "Creating database backup: $BACKUP_FILE"
      mysqldump -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME > $BACKUP_FILE
      
      if [ $? -eq 0 ]; then
        echo "Backup completed successfully: $BACKUP_FILE"
        # Compress backup
        gzip $BACKUP_FILE
        echo "Backup compressed: $BACKUP_FILE.gz"
      else
        echo "Backup failed!"
        exit 1
      fi
    EOF
    env = {
      DB_HOST     = var.database_host
      DB_USER     = var.database_user
      DB_PASSWORD = var.database_password
      DB_NAME     = var.database_name
    }
    timeout = 1800  # 30 minutes
  }
}
```

### Docker Operations
```hcl
step "docker_build_and_push" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      set -e
      
      IMAGE_NAME="${APP_NAME}:${BUILD_TAG}"
      
      echo "Building Docker image: $IMAGE_NAME"
      docker build -t $IMAGE_NAME .
      
      echo "Tagging image..."
      docker tag $IMAGE_NAME $REGISTRY_URL/$IMAGE_NAME
      
      echo "Pushing to registry..."
      docker push $REGISTRY_URL/$IMAGE_NAME
      
      echo "Cleaning up local image..."
      docker rmi $IMAGE_NAME
      
      echo "Build and push completed: $REGISTRY_URL/$IMAGE_NAME"
    EOF
    env = {
      APP_NAME     = var.application_name
      BUILD_TAG    = var.build_tag
      REGISTRY_URL = var.docker_registry
    }
    working_dir = "/src"
  }
}
```

### Infrastructure Validation
```hcl
step "validate_infrastructure" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      ERRORS=0
      
      echo "=== Infrastructure Validation ==="
      
      # Check database connectivity
      echo "Checking database connectivity..."
      if ! mysqladmin ping -h $DB_HOST -u $DB_USER -p$DB_PASSWORD --silent; then
        echo "❌ Database connection failed"
        ERRORS=$((ERRORS + 1))
      else
        echo "✅ Database connection successful"
      fi
      
      # Check Redis connectivity
      echo "Checking Redis connectivity..."
      if ! redis-cli -h $REDIS_HOST ping | grep -q PONG; then
        echo "❌ Redis connection failed"
        ERRORS=$((ERRORS + 1))
      else
        echo "✅ Redis connection successful"
      fi
      
      # Check Kubernetes cluster
      echo "Checking Kubernetes cluster..."
      if ! kubectl cluster-info > /dev/null 2>&1; then
        echo "❌ Kubernetes cluster not accessible"
        ERRORS=$((ERRORS + 1))
      else
        echo "✅ Kubernetes cluster accessible"
      fi
      
      echo "=== Validation Summary ==="
      if [ $ERRORS -eq 0 ]; then
        echo "✅ All infrastructure checks passed"
        exit 0
      else
        echo "❌ $ERRORS infrastructure checks failed"
        exit 1
      fi
    EOF
    env = {
      DB_HOST    = var.database_host
      DB_USER    = var.database_user
      DB_PASSWORD = var.database_password
      REDIS_HOST = var.redis_host
    }
  }
}
```

## File Operations

### Log Processing
```hcl
step "process_logs" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      LOG_DIR="/var/log/application"
      REPORT_FILE="/tmp/log_analysis_$(date +%Y%m%d).txt"
      
      echo "=== Log Analysis Report ===" > $REPORT_FILE
      echo "Generated: $(date)" >> $REPORT_FILE
      echo "" >> $REPORT_FILE
      
      # Count error messages
      ERROR_COUNT=$(grep -r "ERROR" $LOG_DIR | wc -l)
      echo "Total errors: $ERROR_COUNT" >> $REPORT_FILE
      
      # Count warning messages
      WARN_COUNT=$(grep -r "WARN" $LOG_DIR | wc -l)
      echo "Total warnings: $WARN_COUNT" >> $REPORT_FILE
      
      # Top error messages
      echo "" >> $REPORT_FILE
      echo "Top 10 error messages:" >> $REPORT_FILE
      grep -r "ERROR" $LOG_DIR | cut -d':' -f3- | sort | uniq -c | sort -nr | head -10 >> $REPORT_FILE
      
      echo "Log analysis completed: $REPORT_FILE"
    EOF
  }
}
```

### Backup and Archive
```hcl
step "create_backup_archive" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
      BACKUP_DIR="/backups/$BACKUP_DATE"
      ARCHIVE_NAME="backup_$BACKUP_DATE.tar.gz"
      
      echo "Creating backup directory: $BACKUP_DIR"
      mkdir -p $BACKUP_DIR
      
      # Copy application files
      echo "Backing up application files..."
      cp -r /var/www/html $BACKUP_DIR/
      
      # Copy configuration files
      echo "Backing up configuration files..."
      cp -r /etc/nginx $BACKUP_DIR/
      cp -r /etc/ssl $BACKUP_DIR/
      
      # Create archive
      echo "Creating archive: $ARCHIVE_NAME"
      cd /backups
      tar -czf $ARCHIVE_NAME $BACKUP_DATE/
      
      # Clean up temporary directory
      rm -rf $BACKUP_DIR
      
      # Upload to S3 (if configured)
      if [ -n "$AWS_S3_BUCKET" ]; then
        echo "Uploading to S3..."
        aws s3 cp $ARCHIVE_NAME s3://$AWS_S3_BUCKET/backups/
      fi
      
      echo "Backup completed: $ARCHIVE_NAME"
    EOF
    env = {
      AWS_S3_BUCKET = var.backup_s3_bucket
    }
  }
}
```

## Error Handling and Debugging

### Comprehensive Error Handling
```hcl
step "robust_operation" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      set -e  # Exit on error
      set -u  # Exit on undefined variable
      set -o pipefail  # Exit on pipe failure
      
      # Function for error handling
      error_handler() {
        echo "Error occurred in script at line $1"
        echo "Exit code: $2"
        # Cleanup operations
        cleanup_resources
        exit $2
      }
      
      # Set trap for error handling
      trap 'error_handler $LINENO $?' ERR
      
      cleanup_resources() {
        echo "Cleaning up resources..."
        # Add cleanup logic here
      }
      
      # Main script logic
      echo "Starting operation..."
      
      # Your operations here
      perform_operation
      
      echo "Operation completed successfully"
    EOF
  }
}
```

### Debug Mode Execution
```hcl
step "debug_script" {
  plugin = "shell"
  action = "script"
  params = {
    script = <<-EOF
      #!/bin/bash
      # Enable debug mode if DEBUG is set
      if [ "$DEBUG" = "true" ]; then
        set -x  # Print commands as they execute
      fi
      
      echo "Debug mode: $DEBUG"
      echo "Environment: $ENVIRONMENT"
      
      # Your script logic here
    EOF
    env = {
      DEBUG       = var.debug_mode
      ENVIRONMENT = var.environment
    }
  }
}
```

## Best Practices

1. **Always set timeouts** for long-running operations
2. **Use error handling** with `set -e` and trap functions
3. **Validate inputs** before executing commands
4. **Use environment variables** instead of hardcoded values
5. **Implement cleanup logic** for temporary resources
6. **Log operations** for debugging and auditing
7. **Use working directories** to isolate operations
8. **Quote variables** to handle spaces and special characters

## Security Considerations

1. **Never log sensitive information** like passwords or API keys
2. **Use secure environment variable handling**
3. **Validate and sanitize inputs** to prevent injection attacks
4. **Use least privilege** when running commands
5. **Avoid running as root** when possible

## Troubleshooting

### Common Issues

**Command Not Found**
```
Error: command not found: kubectl
```
- Verify command is installed and in PATH
- Use full path to executable
- Install required dependencies

**Permission Denied**
```
Error: permission denied
```
- Check file permissions
- Ensure user has necessary privileges
- Use sudo if required (with caution)

**Timeout Errors**
```
Error: context deadline exceeded
```
- Increase timeout value
- Optimize command performance
- Break down into smaller operations

### Debug Mode

Enable detailed logging:
```bash
export CORYNTH_DEBUG=true
corynth run workflow.hcl
```

## Sample Workflows

See the `/samples` directory for complete workflow examples:
- `system-info.hcl`: System information gathering
- `build-deploy.hcl`: Application build and deployment automation
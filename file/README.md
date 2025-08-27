# File Plugin

## Overview

The File Plugin provides comprehensive file system operations for Corynth workflows, enabling file manipulation, template processing, and directory management for automation tasks.

## Features

- File read/write operations with encoding support
- Directory operations (create, copy, move, delete)
- File existence and metadata checking
- Template processing with variable substitution
- Recursive file operations
- Permission and ownership management
- Path validation and security controls

## Actions

### read
Reads file content

**Parameters:**
- `path` (string, required): File path to read
- `encoding` (string, optional): File encoding (default: "utf-8")

**Returns:**
- `content`: File content as string
- `size`: File size in bytes
- `modified_time`: Last modification timestamp

### write
Writes content to file

**Parameters:**
- `path` (string, required): File path to write
- `content` (string, required): Content to write
- `encoding` (string, optional): File encoding (default: "utf-8")
- `mode` (string, optional): File permissions (default: "0644")
- `create_dirs` (bool, optional): Create parent directories (default: true)

**Returns:**
- `path`: Written file path
- `size`: Written file size in bytes

### copy
Copies files or directories

**Parameters:**
- `source` (string, required): Source path
- `destination` (string, required): Destination path
- `recursive` (bool, optional): Copy directories recursively (default: false)
- `overwrite` (bool, optional): Overwrite existing files (default: false)

**Returns:**
- `source`: Source path
- `destination`: Destination path
- `files_copied`: Number of files copied

### move
Moves/renames files or directories

**Parameters:**
- `source` (string, required): Source path
- `destination` (string, required): Destination path
- `overwrite` (bool, optional): Overwrite existing files (default: false)

**Returns:**
- `source`: Original source path
- `destination`: New destination path

### delete
Deletes files or directories

**Parameters:**
- `path` (string, required): Path to delete
- `recursive` (bool, optional): Delete directories recursively (default: false)
- `force` (bool, optional): Force deletion of read-only files (default: false)

**Returns:**
- `path`: Deleted path
- `files_deleted`: Number of files deleted

### exists
Checks if file or directory exists

**Parameters:**
- `path` (string, required): Path to check

**Returns:**
- `exists`: Boolean indicating existence
- `is_file`: Boolean indicating if path is a file
- `is_directory`: Boolean indicating if path is a directory
- `size`: File size (if file exists)

### template
Processes template files with variable substitution

**Parameters:**
- `template` (string, required): Template file path or inline template
- `output` (string, required): Output file path
- `variables` (map, required): Variables for substitution
- `template_engine` (string, optional): Template engine ("go", "mustache") (default: "go")

**Returns:**
- `template`: Template path used
- `output`: Generated output file path
- `variables_used`: Number of variables substituted

## Usage Examples

### Basic File Operations
```hcl
step "read_config" {
  plugin = "file"
  action = "read"
  params = {
    path = "/etc/app/config.yaml"
  }
}

step "write_log" {
  plugin = "file"
  action = "write"
  params = {
    path    = "/var/log/deployment.log"
    content = "Deployment started at ${timestamp()}\n"
    mode    = "0644"
  }
}
```

### File Backup and Archival
```hcl
step "backup_config" {
  plugin = "file"
  action = "copy"
  params = {
    source      = "/etc/nginx/nginx.conf"
    destination = "/backups/nginx.conf.backup.${formatdate('YYYY-MM-DD', timestamp())}"
  }
}

step "archive_logs" {
  plugin = "file"
  action = "copy"
  params = {
    source      = "/var/log/application/"
    destination = "/archives/logs-${formatdate('YYYY-MM-DD', timestamp())}/"
    recursive   = true
  }
}
```

### Template Processing
```hcl
step "generate_nginx_config" {
  plugin = "file"
  action = "template"
  params = {
    template = "templates/nginx.conf.tpl"
    output   = "/etc/nginx/sites-available/${var.site_name}"
    variables = {
      server_name = var.domain_name
      port        = var.port
      ssl_cert    = var.ssl_certificate_path
      ssl_key     = var.ssl_private_key_path
      root_dir    = var.web_root
    }
  }
}
```

### Configuration Management
```hcl
step "update_application_config" {
  plugin = "file"
  action = "template"
  params = {
    template = <<-EOF
      database:
        host: {{.db_host}}
        port: {{.db_port}}
        name: {{.db_name}}
        user: {{.db_user}}
        password: {{.db_password}}
      
      redis:
        host: {{.redis_host}}
        port: {{.redis_port}}
      
      app:
        environment: {{.environment}}
        debug: {{.debug_mode}}
        log_level: {{.log_level}}
    EOF
    output = "/app/config/production.yaml"
    variables = {
      db_host     = var.database_host
      db_port     = var.database_port
      db_name     = var.database_name
      db_user     = var.database_user
      db_password = var.database_password
      redis_host  = var.redis_host
      redis_port  = var.redis_port
      environment = var.environment
      debug_mode  = var.debug_enabled
      log_level   = var.log_level
    }
  }
}
```

### Conditional File Operations
```hcl
step "check_config_exists" {
  plugin = "file"
  action = "exists"
  params = {
    path = "/etc/app/config.yaml"
  }
}

step "create_default_config" {
  plugin = "file"
  action = "write"
  depends_on = ["check_config_exists"]
  condition  = "!${check_config_exists.exists}"
  params = {
    path = "/etc/app/config.yaml"
    content = <<-EOF
      # Default configuration
      app:
        name: "default-app"
        port: 8080
        debug: false
    EOF
  }
}
```

## Advanced Usage

### Directory Structure Creation
```hcl
step "create_project_structure" {
  plugin = "file"
  action = "write"
  params = {
    path    = "/projects/${var.project_name}/README.md"
    content = "# ${var.project_name}\n\nProject created on ${timestamp()}"
    create_dirs = true
  }
}

step "create_src_directory" {
  plugin = "file"
  action = "write"
  depends_on = ["create_project_structure"]
  params = {
    path    = "/projects/${var.project_name}/src/.gitkeep"
    content = ""
  }
}

step "create_docs_directory" {
  plugin = "file"
  action = "write"
  depends_on = ["create_project_structure"]
  params = {
    path    = "/projects/${var.project_name}/docs/.gitkeep"
    content = ""
  }
}
```

### Log File Processing
```hcl
step "read_application_logs" {
  plugin = "file"
  action = "read"
  params = {
    path = "/var/log/application/app.log"
  }
}

step "extract_errors" {
  plugin = "shell"
  action = "exec"
  depends_on = ["read_application_logs"]
  params = {
    command = "echo '${read_application_logs.content}' | grep ERROR > /tmp/errors.log"
  }
}

step "create_error_report" {
  plugin = "file"
  action = "template"
  depends_on = ["extract_errors"]
  params = {
    template = <<-EOF
      # Error Report
      Generated: {{.timestamp}}
      Log File: {{.log_file}}
      
      ## Error Summary
      {{.error_content}}
    EOF
    output = "/reports/error-report-${formatdate('YYYY-MM-DD', timestamp())}.md"
    variables = {
      timestamp     = timestamp()
      log_file      = "/var/log/application/app.log"
      error_content = file("/tmp/errors.log")
    }
  }
}
```

### Deployment Artifacts Management
```hcl
step "prepare_deployment_directory" {
  plugin = "file"
  action = "delete"
  params = {
    path      = "/tmp/deployment"
    recursive = true
    force     = true
  }
}

step "create_deployment_structure" {
  plugin = "file"
  action = "write"
  depends_on = ["prepare_deployment_directory"]
  params = {
    path        = "/tmp/deployment/config/app.yaml"
    content     = ""
    create_dirs = true
  }
}

step "copy_application_files" {
  plugin = "file"
  action = "copy"
  depends_on = ["create_deployment_structure"]
  params = {
    source      = "/src/build/"
    destination = "/tmp/deployment/app/"
    recursive   = true
  }
}

step "generate_deployment_manifest" {
  plugin = "file"
  action = "template"
  depends_on = ["copy_application_files"]
  params = {
    template = "k8s/deployment.yaml.tpl"
    output   = "/tmp/deployment/k8s/deployment.yaml"
    variables = {
      app_name      = var.application_name
      image_tag     = var.build_tag
      replicas      = var.replica_count
      environment   = var.environment
      resource_cpu  = var.cpu_limit
      resource_mem  = var.memory_limit
    }
  }
}
```

### Backup Rotation
```hcl
step "create_timestamped_backup" {
  plugin = "file"
  action = "copy"
  params = {
    source      = "/data/important/"
    destination = "/backups/data-${formatdate('YYYY-MM-DD-HH-mm', timestamp())}/"
    recursive   = true
  }
}

step "cleanup_old_backups" {
  plugin = "shell"
  action = "exec"
  depends_on = ["create_timestamped_backup"]
  params = {
    command = "find /backups -type d -name 'data-*' -mtime +7 -exec rm -rf {} +"
  }
}
```

## Security and Permissions

### Secure File Operations
```hcl
step "create_secure_config" {
  plugin = "file"
  action = "write"
  params = {
    path    = "/etc/app/secrets.conf"
    content = "api_key=${var.api_key}\ndb_password=${var.db_password}"
    mode    = "0600"  # Read/write for owner only
  }
}
```

### Validation Before Operations
```hcl
step "validate_source_exists" {
  plugin = "file"
  action = "exists"
  params = {
    path = var.source_file
  }
}

step "safe_file_copy" {
  plugin = "file"
  action = "copy"
  depends_on = ["validate_source_exists"]
  condition  = "${validate_source_exists.exists}"
  params = {
    source      = var.source_file
    destination = var.destination_file
    overwrite   = false
  }
}
```

## Error Handling

### File Operation Error Handling
```hcl
step "safe_file_operation" {
  plugin = "file"
  action = "read"
  params = {
    path = "/path/that/might/not/exist.txt"
  }
}

step "handle_missing_file" {
  plugin = "file"
  action = "write"
  depends_on = ["safe_file_operation"]
  condition  = "${safe_file_operation.error != ''}"
  params = {
    path    = "/path/that/might/not/exist.txt"
    content = "Default content created by workflow"
  }
}
```

## Performance Considerations

1. **Use appropriate operations**: Don't read large files into memory unnecessarily
2. **Batch operations**: Combine multiple file operations when possible
3. **Set proper permissions**: Use restrictive permissions for security
4. **Clean up temporary files**: Remove temporary files after use
5. **Use recursive operations carefully**: Large directory trees can be slow

## Best Practices

1. **Always validate paths** before operations
2. **Use relative paths** when possible for portability
3. **Set appropriate file permissions** for security
4. **Create parent directories** when needed
5. **Handle errors gracefully** with conditional logic
6. **Use templates** for configuration management
7. **Implement backup strategies** before destructive operations
8. **Quote file paths** that might contain spaces

## Troubleshooting

### Common Issues

**Permission Denied**
```
Error: permission denied: /etc/sensitive/config
```
- Check file permissions and ownership
- Ensure workflow runs with appropriate privileges
- Use sudo if necessary (with caution)

**File Not Found**
```
Error: file not found: /path/to/file
```
- Verify file path is correct
- Check if file exists using exists action
- Create file if needed

**Template Parsing Error**
```
Error: template parsing failed
```
- Verify template syntax
- Check variable names match template placeholders
- Ensure all required variables are provided

### Debug Mode

Enable detailed logging:
```bash
export CORYNTH_DEBUG=true
corynth run workflow.hcl
```

## Sample Workflows

See the `/samples` directory for complete workflow examples:
- `config-management.hcl`: Configuration file generation and management
- `log-processing.hcl`: Log file analysis and reporting
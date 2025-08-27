workflow "log-processing" {
  description = "Process and analyze application logs"
  version     = "1.0.0"

  variable "log_path" {
    type        = string
    default     = "/var/log/app.log"
    description = "Path to application log file"
  }

  step "check_log_exists" {
    plugin = "file"
    action = "exists"
    
    params = {
      path = var.log_path
    }
  }

  step "read_log_file" {
    plugin = "file"
    action = "read"
    
    depends_on = ["check_log_exists"]
    
    params = {
      path = var.log_path
    }
  }

  step "create_error_log" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["read_log_file"]
    
    params = {
      command = "echo '${read_log_file.content}' | grep -i error > /tmp/errors.log || echo 'No errors found'"
    }
  }

  step "generate_log_summary" {
    plugin = "file"
    action = "template"
    
    depends_on = ["create_error_log"]
    
    params = {
      template = <<-EOF
        Log Analysis Report
        ===================
        Generated: $(date)
        Source: ${log_path}
        Size: ${log_size} bytes
        
        Summary:
        - Total lines processed
        - Error count: (see /tmp/errors.log)
        - Analysis completed successfully
        
        Next Steps:
        1. Review error log at /tmp/errors.log
        2. Monitor for recurring issues
        3. Update log rotation if needed
      EOF
      variables = {
        "log_path" = var.log_path
        "log_size" = "${read_log_file.size}"
      }
      output = "/tmp/log-analysis-report.txt"
    }
  }

  step "cleanup_temp_files" {
    plugin = "file"
    action = "delete"
    
    depends_on = ["generate_log_summary"]
    
    params = {
      path = "/tmp/errors.log"
      recursive = false
    }
  }
}
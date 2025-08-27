workflow "system-info" {
  description = "Gather comprehensive system information"
  version     = "1.0.0"

  step "get_hostname" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "hostname"
      shell = "bash"
    }
  }

  step "get_os_info" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["get_hostname"]
    
    params = {
      command = "uname -a"
      shell = "bash"
    }
  }

  step "get_disk_usage" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["get_os_info"]
    
    params = {
      command = "df -h /"
      shell = "bash"
    }
  }

  step "get_memory_info" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["get_disk_usage"]
    
    params = {
      command = "free -h || vm_stat"
      shell = "bash"
    }
  }

  step "generate_report" {
    plugin = "shell"
    action = "script"
    
    depends_on = ["get_memory_info"]
    
    params = {
      script = <<-EOF
        #!/bin/bash
        echo "=== System Information Report ==="
        echo "Generated: $(date)"
        echo "Hostname: ${get_hostname.stdout}"
        echo "OS: ${get_os_info.stdout}"
        echo ""
        echo "=== Disk Usage ==="
        echo "${get_disk_usage.stdout}"
        echo ""
        echo "=== Memory Usage ==="
        echo "${get_memory_info.stdout}"
      EOF
      shell = "bash"
    }
  }
}
workflow "config-management" {
  description = "Manage application configuration files"
  version     = "1.0.0"

  variable "app_name" {
    type        = string
    default     = "webapp"
    description = "Application name"
  }

  variable "environment" {
    type        = string
    default     = "staging"
    description = "Target environment"
  }

  step "create_config_dir" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "mkdir -p /tmp/configs/${var.app_name}"
    }
  }

  step "generate_app_config" {
    plugin = "file"
    action = "template"
    
    depends_on = ["create_config_dir"]
    
    params = {
      template = <<-EOF
        {
          "app_name": "${app_name}",
          "environment": "${environment}",
          "database": {
            "host": "${app_name}-db.${environment}.local",
            "port": 5432,
            "name": "${app_name}_${environment}"
          },
          "redis": {
            "host": "${app_name}-cache.${environment}.local",
            "port": 6379
          },
          "logging": {
            "level": "${environment == 'production' ? 'warn' : 'debug'}",
            "output": "/var/log/${app_name}/${environment}.log"
          }
        }
      EOF
      variables = {
        "app_name" = var.app_name
        "environment" = var.environment
      }
      output = "/tmp/configs/${var.app_name}/config.json"
    }
  }

  step "create_nginx_config" {
    plugin = "file"
    action = "write"
    
    depends_on = ["generate_app_config"]
    
    params = {
      path = "/tmp/configs/${var.app_name}/nginx.conf"
      content = <<-EOF
        server {
            listen 80;
            server_name ${var.app_name}.${var.environment}.local;
            
            location / {
                proxy_pass http://localhost:3000;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
            }
            
            location /health {
                access_log off;
                return 200 "healthy\n";
            }
        }
      EOF
    }
  }

  step "verify_configs" {
    plugin = "file"
    action = "exists"
    
    depends_on = ["create_nginx_config"]
    
    params = {
      path = "/tmp/configs/${var.app_name}"
    }
  }

  step "backup_configs" {
    plugin = "file"
    action = "copy"
    
    depends_on = ["verify_configs"]
    
    params = {
      source = "/tmp/configs/${var.app_name}"
      destination = "/tmp/configs-backup/${var.app_name}-$(date +%Y%m%d)"
      overwrite = true
    }
  }
}
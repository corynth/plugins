workflow "multi-environment" {
  description = "Deploy to multiple environments with Terraform"
  version     = "1.0.0"

  variable "environments" {
    type        = list(string)
    default     = ["dev", "staging", "prod"]
    description = "List of environments to deploy"
  }

  variable "tf_workspace" {
    type        = string
    default     = "/tmp/terraform-multi-env"
    description = "Terraform working directory"
  }

  step "prepare_terraform" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "mkdir -p ${var.tf_workspace} && cd ${var.tf_workspace}"
    }
  }

  step "create_main_config" {
    plugin = "file"
    action = "template"
    
    depends_on = ["prepare_terraform"]
    
    params = {
      template = <<-EOF
        terraform {
          required_version = ">= 1.0"
        }

        variable "environment" {
          description = "Environment name"
          type        = string
        }

        resource "local_file" "env_config" {
          filename = "/tmp/config-${environment}.json"
          content = jsonencode({
            environment = var.environment
            deployed_at = timestamp()
            managed_by  = "corynth-terraform"
          })
        }

        output "config_file" {
          value = local_file.env_config.filename
        }
      EOF
      output = "${var.tf_workspace}/main.tf"
    }
  }

  step "init_terraform" {
    plugin = "terraform"
    action = "init"
    
    depends_on = ["create_main_config"]
    
    params = {
      working_dir = var.tf_workspace
    }
  }

  # Deploy to dev environment
  step "plan_dev" {
    plugin = "terraform"
    action = "plan"
    
    depends_on = ["init_terraform"]
    
    params = {
      working_dir = var.tf_workspace
      variables = {
        "environment" = "dev"
      }
    }
  }

  step "apply_dev" {
    plugin = "terraform"
    action = "apply"
    
    depends_on = ["plan_dev"]
    
    params = {
      working_dir = var.tf_workspace
      variables = {
        "environment" = "dev"
      }
      auto_approve = true
    }
  }

  # Deploy to staging environment
  step "plan_staging" {
    plugin = "terraform"
    action = "plan"
    
    depends_on = ["apply_dev"]
    
    params = {
      working_dir = var.tf_workspace
      variables = {
        "environment" = "staging"
      }
    }
  }

  step "apply_staging" {
    plugin = "terraform"
    action = "apply"
    
    depends_on = ["plan_staging"]
    
    params = {
      working_dir = var.tf_workspace
      variables = {
        "environment" = "staging"
      }
      auto_approve = true
    }
  }

  step "verify_deployments" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["apply_staging"]
    
    params = {
      command = "ls -la /tmp/config-*.json && echo 'Multi-environment deployment completed'"
    }
  }
}
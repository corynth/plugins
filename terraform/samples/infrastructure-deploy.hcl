workflow "infrastructure-deploy" {
  description = "Deploy infrastructure using Terraform"
  version     = "1.0.0"

  variable "tf_workspace" {
    type        = string
    default     = "/tmp/terraform-workspace"
    description = "Terraform working directory"
  }

  variable "environment" {
    type        = string
    default     = "staging"
    description = "Target environment"
  }

  step "prepare_workspace" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "mkdir -p ${var.tf_workspace} && cd ${var.tf_workspace}"
    }
  }

  step "create_terraform_config" {
    plugin = "file"
    action = "write"
    
    depends_on = ["prepare_workspace"]
    
    params = {
      path = "${var.tf_workspace}/main.tf"
      content = <<-EOF
        terraform {
          required_version = ">= 1.0"
          required_providers {
            local = {
              source  = "hashicorp/local"
              version = "~> 2.0"
            }
          }
        }

        variable "environment" {
          description = "Environment name"
          type        = string
          default     = "${var.environment}"
        }

        resource "local_file" "deployment_marker" {
          filename = "/tmp/deployment-${var.environment}.txt"
          content  = "Deployed at: $(date)\nEnvironment: ${var.environment}\nManaged by: Corynth Terraform Plugin"
        }

        output "deployment_file" {
          value = local_file.deployment_marker.filename
        }
      EOF
    }
  }

  step "terraform_init" {
    plugin = "terraform"
    action = "init"
    
    depends_on = ["create_terraform_config"]
    
    params = {
      working_dir = var.tf_workspace
    }
  }

  step "terraform_plan" {
    plugin = "terraform"
    action = "plan"
    
    depends_on = ["terraform_init"]
    
    params = {
      working_dir = var.tf_workspace
    }
  }

  step "terraform_apply" {
    plugin = "terraform"
    action = "apply"
    
    depends_on = ["terraform_plan"]
    
    params = {
      working_dir = var.tf_workspace
      auto_approve = true
    }
  }

  step "verify_deployment" {
    plugin = "file"
    action = "exists"
    
    depends_on = ["terraform_apply"]
    
    params = {
      path = "/tmp/deployment-${var.environment}.txt"
    }
  }
}
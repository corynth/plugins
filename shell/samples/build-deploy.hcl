workflow "build-deploy" {
  description = "Build and deploy application with shell commands"
  version     = "1.0.0"

  variable "app_name" {
    type        = string
    default     = "my-app"
    description = "Application name"
  }

  variable "build_env" {
    type        = string
    default     = "production"
    description = "Build environment"
  }

  step "clean_build" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "rm -rf dist/ build/ node_modules/.cache"
      working_dir = "/tmp/build-workspace"
    }
  }

  step "install_dependencies" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["clean_build"]
    
    params = {
      command = "npm install --production"
      working_dir = "/tmp/build-workspace"
      env = {
        "NODE_ENV" = var.build_env
        "APP_NAME" = var.app_name
      }
      timeout = 300
    }
  }

  step "run_build" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["install_dependencies"]
    
    params = {
      command = "npm run build"
      working_dir = "/tmp/build-workspace"
      env = {
        "NODE_ENV" = var.build_env
        "APP_NAME" = var.app_name
      }
      timeout = 600
    }
  }

  step "run_tests" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["run_build"]
    
    params = {
      command = "npm test -- --coverage"
      working_dir = "/tmp/build-workspace"
      env = {
        "CI" = "true"
        "NODE_ENV" = "test"
      }
      timeout = 300
    }
  }

  step "deploy_application" {
    plugin = "shell"
    action = "script"
    
    depends_on = ["run_tests"]
    
    params = {
      script = <<-EOF
        #!/bin/bash
        set -e
        
        echo "Deploying ${var.app_name} to ${var.build_env}"
        
        # Create deployment package
        tar -czf ${var.app_name}-$(date +%Y%m%d-%H%M%S).tar.gz dist/
        
        # Upload to deployment server (mock)
        echo "Uploading deployment package..."
        echo "Deployment completed successfully!"
        
        # Health check
        echo "Running post-deployment health check..."
        sleep 2
        echo "Health check passed!"
      EOF
      working_dir = "/tmp/build-workspace"
      env = {
        "DEPLOY_ENV" = var.build_env
        "APP_NAME" = var.app_name
      }
      timeout = 180
    }
  }
}
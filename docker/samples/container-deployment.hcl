workflow "container-deployment" {
  description = "Build and deploy Docker containers"
  version     = "1.0.0"

  variable "app_name" {
    type        = string
    default     = "webapp"
    description = "Application name"
  }

  variable "version" {
    type        = string
    default     = "latest"
    description = "Application version tag"
  }

  variable "port_mapping" {
    type        = string
    default     = "8080:80"
    description = "Port mapping for container"
  }

  step "prepare_build_context" {
    plugin = "shell"
    action = "exec"
    
    params = {
      command = "mkdir -p /tmp/docker-build && cd /tmp/docker-build"
    }
  }

  step "create_dockerfile" {
    plugin = "file"
    action = "write"
    
    depends_on = ["prepare_build_context"]
    
    params = {
      path = "/tmp/docker-build/Dockerfile"
      content = <<-EOF
        FROM nginx:alpine
        
        # Copy application files
        COPY . /usr/share/nginx/html/
        
        # Create simple index page
        RUN echo '<h1>${var.app_name} v${var.version}</h1><p>Deployed via Corynth Docker Plugin</p>' > /usr/share/nginx/html/index.html
        
        EXPOSE 80
        
        CMD ["nginx", "-g", "daemon off;"]
      EOF
    }
  }

  step "build_image" {
    plugin = "docker"
    action = "build"
    
    depends_on = ["create_dockerfile"]
    
    params = {
      context = "/tmp/docker-build"
      dockerfile = "Dockerfile"
      tag = "${var.app_name}:${var.version}"
    }
  }

  step "run_container" {
    plugin = "docker"
    action = "run"
    
    depends_on = ["build_image"]
    
    params = {
      image = "${var.app_name}:${var.version}"
      name = "${var.app_name}-container"
      ports = [var.port_mapping]
      detached = true
      environment = {
        "APP_NAME" = var.app_name
        "APP_VERSION" = var.version
        "DEPLOYED_BY" = "corynth"
      }
    }
  }

  step "verify_container" {
    plugin = "docker"
    action = "ps"
    
    depends_on = ["run_container"]
    
    params = {
      all = false
    }
  }

  step "test_application" {
    plugin = "http"
    action = "get"
    
    depends_on = ["verify_container"]
    
    params = {
      url = "http://localhost:8080"
      timeout = 10
    }
  }
}
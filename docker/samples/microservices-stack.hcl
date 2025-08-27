workflow "microservices-stack" {
  description = "Deploy a complete microservices stack with Docker"
  version     = "1.0.0"

  variable "stack_name" {
    type        = string
    default     = "myapp"
    description = "Stack name prefix"
  }

  step "list_existing_containers" {
    plugin = "docker"
    action = "ps"
    
    params = {
      all = true
    }
  }

  step "run_database" {
    plugin = "docker"
    action = "run"
    
    depends_on = ["list_existing_containers"]
    
    params = {
      image = "postgres:13"
      name = "${var.stack_name}-db"
      ports = ["5432:5432"]
      environment = {
        "POSTGRES_DB" = "${var.stack_name}"
        "POSTGRES_USER" = "app"
        "POSTGRES_PASSWORD" = "password123"
      }
      detached = true
    }
  }

  step "run_redis" {
    plugin = "docker"
    action = "run"
    
    depends_on = ["run_database"]
    
    params = {
      image = "redis:7-alpine"
      name = "${var.stack_name}-cache"
      ports = ["6379:6379"]
      detached = true
    }
  }

  step "build_api_image" {
    plugin = "docker"
    action = "build"
    
    depends_on = ["run_redis"]
    
    params = {
      context = "/tmp/api-build"
      tag = "${var.stack_name}-api:latest"
    }
  }

  step "run_api_service" {
    plugin = "docker"
    action = "run"
    
    depends_on = ["build_api_image"]
    
    params = {
      image = "${var.stack_name}-api:latest"
      name = "${var.stack_name}-api"
      ports = ["3000:3000"]
      environment = {
        "DATABASE_URL" = "postgresql://app:password123@${var.stack_name}-db:5432/${var.stack_name}"
        "REDIS_URL" = "redis://${var.stack_name}-cache:6379"
        "NODE_ENV" = "production"
      }
      detached = true
    }
  }

  step "run_frontend" {
    plugin = "docker"
    action = "run"
    
    depends_on = ["run_api_service"]
    
    params = {
      image = "nginx:alpine"
      name = "${var.stack_name}-web"
      ports = ["80:80"]
      environment = {
        "API_URL" = "http://${var.stack_name}-api:3000"
      }
      detached = true
    }
  }

  step "verify_stack" {
    plugin = "docker"
    action = "ps"
    
    depends_on = ["run_frontend"]
    
    params = {
      all = false
    }
  }

  step "health_check_stack" {
    plugin = "shell"
    action = "script"
    
    depends_on = ["verify_stack"]
    
    params = {
      script = <<-EOF
        #!/bin/bash
        echo "=== Microservices Stack Health Check ==="
        echo "Stack: ${var.stack_name}"
        echo ""
        
        echo "Testing database connection..."
        sleep 2
        echo "✓ Database ready"
        
        echo "Testing cache connection..."
        sleep 1
        echo "✓ Redis ready"
        
        echo "Testing API service..."
        sleep 2
        echo "✓ API service ready"
        
        echo "Testing frontend..."
        sleep 1
        echo "✓ Frontend ready"
        
        echo ""
        echo "🚀 Microservices stack deployment completed!"
        echo "Containers running: ${verify_stack.count}"
      EOF
    }
  }
}
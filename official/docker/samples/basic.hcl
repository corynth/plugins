workflow "docker-basic-operations" {
  description = "Basic Docker container operations"
  version     = "1.0.0"

  step "run_nginx" {
    plugin = "docker"
    action = "run"
    
    params = {
      image   = "nginx:alpine"
      name    = "test-nginx"
      detach  = true
      ports   = ["8080:80"]
      remove  = false
    }
  }

  step "check_containers" {
    plugin = "docker"
    action = "ps"
    
    depends_on = ["run_nginx"]
    
    params = {
      all = false
    }
  }

  step "get_logs" {
    plugin = "docker"
    action = "logs"
    
    depends_on = ["run_nginx"]
    
    params = {
      container = "test-nginx"
      tail      = 50
    }
  }

  step "exec_command" {
    plugin = "docker"
    action = "exec"
    
    depends_on = ["run_nginx"]
    
    params = {
      container = "test-nginx"
      command   = "nginx -v"
    }
  }

  step "stop_container" {
    plugin = "docker"
    action = "stop"
    
    depends_on = ["get_logs", "exec_command"]
    
    params = {
      container = "test-nginx"
      timeout   = 10
    }
  }
}
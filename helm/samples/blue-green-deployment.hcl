workflow "helm-blue-green-deployment" {
  description = "Blue-green deployment pattern using Helm"
  version = "1.0.0"

  step "deploy_green_version" {
    plugin = "helm"
    action = "install"
    
    params = {
      name = "app-green"
      chart = "stable/webapp"
      namespace = "production"
      create_namespace = true
      values = {
        image = {
          tag = "v2.0.0"
        }
        service = {
          name = "app-green"
          selector = {
            version = "green"
          }
        }
        labels = {
          version = "green"
        }
        replicaCount = 3
      }
      wait = true
      timeout = "15m"
    }
  }

  step "test_green_deployment" {
    plugin = "helm"
    action = "status"
    depends_on = ["deploy_green_version"]
    
    params = {
      name = "app-green"
      namespace = "production"
    }
  }

  step "switch_service_to_green" {
    plugin = "helm"
    action = "upgrade"
    depends_on = ["test_green_deployment"]
    
    params = {
      name = "app-service"
      chart = "stable/service"
      namespace = "production"
      install = true
      values = {
        selector = {
          version = "green"
        }
        port = 80
        targetPort = 8080
      }
      wait = true
    }
  }

  step "verify_traffic_switch" {
    plugin = "helm"
    action = "status"
    depends_on = ["switch_service_to_green"]
    
    params = {
      name = "app-service"
      namespace = "production"
    }
  }

  step "cleanup_blue_version" {
    plugin = "helm"
    action = "uninstall"
    depends_on = ["verify_traffic_switch"]
    
    params = {
      name = "app-blue"
      namespace = "production"
      keep_history = false
    }
  }

  step "final_status_check" {
    plugin = "helm"
    action = "list"
    depends_on = ["cleanup_blue_version"]
    
    params = {
      namespace = "production"
      status = "deployed"
    }
  }
}
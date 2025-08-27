workflow "helm-application-deployment" {
  description = "Deploy application using Helm with full lifecycle management"
  version = "1.0.0"

  step "add_repository" {
    plugin = "helm"
    action = "repo_add"
    
    params = {
      name = "bitnami"
      url = "https://charts.bitnami.com/bitnami"
    }
  }

  step "update_repositories" {
    plugin = "helm"
    action = "repo_update"
    depends_on = ["add_repository"]
  }

  step "install_nginx" {
    plugin = "helm"
    action = "install"
    depends_on = ["update_repositories"]
    
    params = {
      name = "web-server"
      chart = "bitnami/nginx"
      namespace = "production"
      create_namespace = true
      version = "15.4.4"
      values = {
        replicaCount = 3
        service = {
          type = "ClusterIP"
          port = 80
        }
        resources = {
          requests = {
            cpu = "100m"
            memory = "128Mi"
          }
          limits = {
            cpu = "250m"
            memory = "256Mi"
          }
        }
      }
      wait = true
      timeout = "10m"
    }
  }

  step "verify_installation" {
    plugin = "helm"
    action = "status"
    depends_on = ["install_nginx"]
    
    params = {
      name = "web-server"
      namespace = "production"
    }
  }

  step "list_releases" {
    plugin = "helm"
    action = "list"
    depends_on = ["verify_installation"]
    
    params = {
      namespace = "production"
    }
  }
}
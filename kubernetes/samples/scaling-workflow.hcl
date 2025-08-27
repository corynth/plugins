workflow "kubernetes-scaling" {
  description = "Scale Kubernetes deployment and monitor results"
  version = "1.0.0"

  step "initial_scale" {
    plugin = "kubernetes"
    action = "scale"
    
    params = {
      resource = "deployment"
      name = "nginx-deployment"
      replicas = 5
      namespace = "default"
    }
  }

  step "wait_scale_complete" {
    plugin = "kubernetes"
    action = "wait"
    depends_on = ["initial_scale"]
    
    params = {
      resource = "deployment"
      name = "nginx-deployment"
      condition = "available"
      timeout = 180
      namespace = "default"
    }
  }

  step "verify_replica_count" {
    plugin = "kubernetes"
    action = "get"
    depends_on = ["wait_scale_complete"]
    
    params = {
      resource = "deployment"
      name = "nginx-deployment"
      namespace = "default"
    }
  }

  step "scale_down" {
    plugin = "kubernetes"
    action = "scale"
    depends_on = ["verify_replica_count"]
    
    params = {
      resource = "deployment"
      name = "nginx-deployment"
      replicas = 2
      namespace = "default"
    }
  }

  step "final_verification" {
    plugin = "kubernetes"
    action = "get"
    depends_on = ["scale_down"]
    
    params = {
      resource = "pods"
      namespace = "default"
    }
  }
}
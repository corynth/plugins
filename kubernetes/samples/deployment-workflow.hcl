workflow "kubernetes-deployment" {
  description = "Deploy application to Kubernetes with health checks"
  version = "1.0.0"

  step "deploy_application" {
    plugin = "kubernetes"
    action = "apply"
    
    params = {
      manifest = <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
  type: ClusterIP
EOF
      namespace = "default"
    }
  }

  step "wait_deployment_ready" {
    plugin = "kubernetes"
    action = "wait"
    depends_on = ["deploy_application"]
    
    params = {
      resource = "deployment"
      name = "nginx-deployment"
      condition = "available"
      timeout = 300
      namespace = "default"
    }
  }

  step "verify_pods" {
    plugin = "kubernetes"
    action = "get"
    depends_on = ["wait_deployment_ready"]
    
    params = {
      resource = "pods"
      namespace = "default"
    }
  }

  step "check_logs" {
    plugin = "kubernetes"
    action = "logs"
    depends_on = ["verify_pods"]
    
    params = {
      pod = "nginx-deployment"  # Will be replaced with actual pod name
      namespace = "default"
      tail = 50
    }
  }
}
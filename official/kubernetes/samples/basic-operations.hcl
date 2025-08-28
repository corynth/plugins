# Kubernetes Plugin Basic Operations Example
# This example demonstrates common Kubernetes operations using the plugin

step "list_pods" {
  plugin = "kubernetes"
  action = "get"
  params = {
    resource = "pods"
    all_namespaces = true
  }
}

step "list_services" {
  plugin = "kubernetes"
  action = "get"
  params = {
    resource = "services"
    namespace = "default"
  }
}

step "describe_nodes" {
  plugin = "kubernetes"
  action = "get"
  params = {
    resource = "nodes"
    output = "wide"
  }
}

# Example of applying a manifest
step "apply_manifest" {
  plugin = "kubernetes"
  action = "apply"
  params = {
    manifest = <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
spec:
  containers:
  - name: busybox
    image: busybox:latest
    command: ["sleep", "3600"]
EOF
    dry_run = true  # Use dry run for safety
  }
}

# Example of scaling a deployment
step "scale_deployment" {
  plugin = "kubernetes"
  action = "scale"
  params = {
    resource = "deployment"
    name = "my-app"
    replicas = 3
    namespace = "default"
  }
}

# Example of getting pod logs
step "get_pod_logs" {
  plugin = "kubernetes"
  action = "logs"
  params = {
    pod = "test-pod"
    namespace = "default"
    tail = 50
  }
}
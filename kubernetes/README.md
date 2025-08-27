# Kubernetes Plugin

Production-ready Kubernetes cluster management and resource operations plugin for Corynth workflows.

## Features

- **Apply Manifests** - Deploy YAML manifests or direct content
- **Resource Management** - Get, delete, scale Kubernetes resources
- **Pod Operations** - Retrieve logs, wait for conditions
- **Multi-cluster** - Support for multiple kubeconfig files
- **Namespace Support** - Operations across different namespaces

## Prerequisites

- `kubectl` installed and configured
- Valid kubeconfig file (default: `~/.kube/config`)
- Kubernetes cluster access

## Actions

### apply
Deploy Kubernetes manifests to cluster.

```hcl
step "deploy_app" {
  plugin = "kubernetes"
  action = "apply"
  
  params = {
    manifest = file("deployment.yaml")
    namespace = "production"
  }
}
```

### get
Retrieve Kubernetes resources.

```hcl
step "list_pods" {
  plugin = "kubernetes"
  action = "get"
  
  params = {
    resource = "pods"
    namespace = "default"
  }
}
```

### delete
Delete Kubernetes resources.

```hcl
step "cleanup" {
  plugin = "kubernetes"
  action = "delete"
  
  params = {
    resource = "deployment"
    name = "my-app"
    namespace = "staging"
  }
}
```

### scale
Scale deployments or replica sets.

```hcl
step "scale_up" {
  plugin = "kubernetes"
  action = "scale"
  
  params = {
    resource = "deployment"
    name = "web-server"
    replicas = 5
    namespace = "production"
  }
}
```

### logs
Get pod logs.

```hcl
step "check_logs" {
  plugin = "kubernetes"
  action = "logs"
  
  params = {
    pod = "web-server-abc123"
    namespace = "production"
    tail = 100
  }
}
```

### wait
Wait for resource condition.

```hcl
step "wait_ready" {
  plugin = "kubernetes"
  action = "wait"
  
  params = {
    resource = "deployment"
    name = "my-app"
    condition = "available"
    timeout = 600
  }
}
```

## Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `manifest` | string | Yes | - | YAML content or file path (apply) |
| `resource` | string | Yes | - | Resource type (pods, services, deployments) |
| `name` | string | Varies | - | Resource name |
| `namespace` | string | No | "default" | Target namespace |
| `kubeconfig` | string | No | - | Path to kubeconfig file |
| `replicas` | number | Yes | - | Target replica count (scale) |
| `condition` | string | Yes | - | Wait condition (available, ready) |
| `timeout` | number | No | 300 | Timeout in seconds |
| `tail` | number | No | 100 | Number of log lines |
| `container` | string | No | - | Container name for logs |

## Outputs

### apply
- `status` - Operation result (success/failed)
- `resources` - List of applied resources
- `output` - kubectl command output

### get
- `resources` - Retrieved resource objects
- `count` - Number of resources found

### delete
- `status` - Operation result (success/failed)
- `deleted` - Whether resource was deleted
- `output` - kubectl command output

### scale
- `status` - Operation result (success/failed)
- `replicas` - Target replica count
- `output` - kubectl command output

### logs
- `logs` - Pod logs content
- `lines` - Number of log lines retrieved

### wait
- `ready` - Whether condition was met
- `status` - Wait result (ready/timeout)
- `output` - kubectl command output

## Examples

### Blue-Green Deployment
```hcl
workflow "blue-green-deployment" {
  description = "Blue-green deployment with Kubernetes"
  
  step "deploy_green" {
    plugin = "kubernetes"
    action = "apply"
    
    params = {
      manifest = file("green-deployment.yaml")
      namespace = "production"
    }
  }
  
  step "wait_green_ready" {
    plugin = "kubernetes"
    action = "wait"
    depends_on = ["deploy_green"]
    
    params = {
      resource = "deployment"
      name = "app-green"
      condition = "available"
      namespace = "production"
      timeout = 600
    }
  }
  
  step "switch_traffic" {
    plugin = "kubernetes"
    action = "apply"
    depends_on = ["wait_green_ready"]
    
    params = {
      manifest = file("service-green.yaml")
      namespace = "production"
    }
  }
  
  step "cleanup_blue" {
    plugin = "kubernetes"
    action = "delete"
    depends_on = ["switch_traffic"]
    
    params = {
      resource = "deployment"
      name = "app-blue"
      namespace = "production"
    }
  }
}
```

### Auto-scaling Workflow
```hcl
workflow "auto-scale" {
  description = "Auto-scale based on conditions"
  
  step "check_load" {
    plugin = "kubernetes"
    action = "get"
    
    params = {
      resource = "pods"
      namespace = "production"
    }
  }
  
  step "scale_up" {
    plugin = "kubernetes"  
    action = "scale"
    depends_on = ["check_load"]
    
    params = {
      resource = "deployment"
      name = "web-server"
      replicas = 10
      namespace = "production"
    }
  }
  
  step "verify_scaling" {
    plugin = "kubernetes"
    action = "wait"
    depends_on = ["scale_up"]
    
    params = {
      resource = "deployment"
      name = "web-server"
      condition = "available"
      namespace = "production"
      timeout = 300
    }
  }
}
```

## Error Handling

The plugin handles common Kubernetes errors gracefully:

- **Connection errors** - Reports cluster connectivity issues
- **Resource not found** - Returns appropriate status for missing resources
- **Permission errors** - Reports RBAC or authentication failures
- **Timeout errors** - Handles long-running operations with configurable timeouts

## Security

- Supports multiple kubeconfig files for different clusters
- Respects Kubernetes RBAC permissions
- No credential storage in workflow files
- Secure handling of manifest content

## Performance

- Efficient JSON parsing of kubectl output
- Minimal memory footprint for large resource lists
- Configurable timeouts for long-running operations
- Temporary file cleanup for manifest operations
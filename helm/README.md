# Helm Plugin

Production-ready Helm package manager plugin for Kubernetes application deployment and management.

## Features

- **Chart Management** - Install, upgrade, uninstall Helm charts
- **Repository Operations** - Add and update Helm repositories  
- **Templating** - Render charts locally for validation
- **Release Management** - List and check status of releases
- **Values Override** - Support for values files and inline values
- **Multi-cluster** - Support for multiple kubeconfig files

## Prerequisites

- `helm` v3.0+ installed and configured
- `kubectl` configured for target cluster
- Valid kubeconfig file (default: `~/.kube/config`)

## Actions

### install
Install a Helm chart as a new release.

```hcl
step "deploy_nginx" {
  plugin = "helm"
  action = "install"
  
  params = {
    name = "my-nginx"
    chart = "nginx"
    namespace = "production"
    repository = "https://charts.bitnami.com/bitnami"
    values = {
      replicaCount = 3
      service = {
        type = "LoadBalancer"
      }
    }
    wait = true
    timeout = "10m"
  }
}
```

### upgrade
Upgrade an existing Helm release.

```hcl
step "upgrade_app" {
  plugin = "helm"
  action = "upgrade"
  
  params = {
    name = "my-app"
    chart = "my-chart"
    namespace = "production"
    version = "2.1.0"
    install = true  # Install if doesn't exist
    values_file = "production-values.yaml"
  }
}
```

### uninstall
Remove a Helm release.

```hcl
step "cleanup" {
  plugin = "helm"
  action = "uninstall"
  
  params = {
    name = "old-app"
    namespace = "staging"
    keep_history = false
  }
}
```

### list
List all Helm releases.

```hcl
step "list_releases" {
  plugin = "helm"
  action = "list"
  
  params = {
    all_namespaces = true
    status = "deployed"
  }
}
```

### status
Get detailed status of a release.

```hcl
step "check_status" {
  plugin = "helm"
  action = "status"
  
  params = {
    name = "my-app"
    namespace = "production"
  }
}
```

### template
Render chart templates locally.

```hcl
step "validate_manifests" {
  plugin = "helm"
  action = "template"
  
  params = {
    name = "test-release"
    chart = "./my-chart"
    namespace = "staging"
    values = {
      environment = "test"
      debug = true
    }
  }
}
```

### repo_add
Add a Helm repository.

```hcl
step "add_bitnami_repo" {
  plugin = "helm"
  action = "repo_add"
  
  params = {
    name = "bitnami"
    url = "https://charts.bitnami.com/bitnami"
  }
}
```

### repo_update
Update repository indexes.

```hcl
step "update_repos" {
  plugin = "helm"
  action = "repo_update"
}
```

## Parameters

### install/upgrade Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `name` | string | Yes | - | Release name |
| `chart` | string | Yes | - | Chart name or path |
| `namespace` | string | No | "default" | Target namespace |
| `values` | object | No | - | Values to override |
| `values_file` | string | No | - | Path to values YAML file |
| `version` | string | No | - | Chart version |
| `repository` | string | No | - | Chart repository URL |
| `create_namespace` | boolean | No | false | Create namespace if missing |
| `wait` | boolean | No | true | Wait for resources to be ready |
| `timeout` | string | No | "5m" | Operation timeout |
| `install` | boolean | No | false | Install if release doesn't exist (upgrade only) |
| `kubeconfig` | string | No | - | Path to kubeconfig file |

### Other Action Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `all_namespaces` | boolean | No | false | List across all namespaces |
| `status` | string | No | - | Filter releases by status |
| `keep_history` | boolean | No | false | Keep release history on uninstall |
| `url` | string | Yes | - | Repository URL (repo_add) |
| `username` | string | No | - | Repository username |
| `password` | string | No | - | Repository password |
| `force_update` | boolean | No | false | Replace existing repo |

## Outputs

### install/upgrade
- `status` - Operation result (success/failed)
- `revision` - Release revision number
- `notes` - Chart installation notes (install only)

### uninstall
- `status` - Operation result (success/failed)
- `removed` - Whether release was removed

### list
- `releases` - Array of release objects
- `count` - Number of releases found

### status
- `status` - Release status (deployed, failed, etc.)
- `revision` - Current revision number
- `chart` - Chart name and version
- `namespace` - Release namespace

### template
- `manifests` - Rendered Kubernetes YAML
- `resources` - List of resource types

### repo_add/repo_update
- `status` - Operation result (success/failed)
- `added`/`updated` - Whether operation succeeded

## Examples

### Complete Application Deployment
```hcl
workflow "deploy-app-with-helm" {
  description = "Deploy application using Helm with repository setup"
  
  step "add_repo" {
    plugin = "helm"
    action = "repo_add"
    
    params = {
      name = "my-charts"
      url = "https://charts.example.com"
    }
  }
  
  step "update_repos" {
    plugin = "helm"
    action = "repo_update"
    depends_on = ["add_repo"]
  }
  
  step "deploy_app" {
    plugin = "helm"
    action = "install"
    depends_on = ["update_repos"]
    
    params = {
      name = "web-app"
      chart = "my-charts/webapp"
      namespace = "production"
      create_namespace = true
      version = "1.5.2"
      values = {
        image = {
          tag = "v1.5.2"
        }
        ingress = {
          enabled = true
          host = "app.example.com"
        }
        resources = {
          requests = {
            cpu = "100m"
            memory = "128Mi"
          }
          limits = {
            cpu = "500m"
            memory = "512Mi"
          }
        }
      }
      wait = true
      timeout = "15m"
    }
  }
  
  step "verify_deployment" {
    plugin = "helm"
    action = "status"
    depends_on = ["deploy_app"]
    
    params = {
      name = "web-app"
      namespace = "production"
    }
  }
}
```

### Blue-Green Deployment
```hcl
workflow "blue-green-helm-deployment" {
  description = "Blue-green deployment using Helm"
  
  step "deploy_green" {
    plugin = "helm"
    action = "install"
    
    params = {
      name = "app-green"
      chart = "./charts/myapp"
      namespace = "production"
      values = {
        version = "green"
        image = {
          tag = "v2.0.0"
        }
        service = {
          selector = {
            version = "green"
          }
        }
      }
      wait = true
      timeout = "10m"
    }
  }
  
  step "test_green" {
    plugin = "helm"
    action = "status"
    depends_on = ["deploy_green"]
    
    params = {
      name = "app-green"
      namespace = "production"
    }
  }
  
  step "switch_traffic" {
    plugin = "helm"
    action = "upgrade"
    depends_on = ["test_green"]
    
    params = {
      name = "app-service"
      chart = "./charts/service"
      namespace = "production"
      values = {
        selector = {
          version = "green"
        }
      }
    }
  }
  
  step "cleanup_blue" {
    plugin = "helm"
    action = "uninstall"
    depends_on = ["switch_traffic"]
    
    params = {
      name = "app-blue"
      namespace = "production"
    }
  }
}
```

### Multi-Environment Deployment
```hcl
workflow "multi-env-deployment" {
  description = "Deploy to multiple environments with environment-specific values"
  
  step "deploy_staging" {
    plugin = "helm"
    action = "install"
    
    params = {
      name = "myapp"
      chart = "stable/myapp"
      namespace = "staging"
      create_namespace = true
      values_file = "values-staging.yaml"
      wait = true
    }
  }
  
  step "test_staging" {
    plugin = "helm"
    action = "status"
    depends_on = ["deploy_staging"]
    
    params = {
      name = "myapp"
      namespace = "staging"
    }
  }
  
  step "deploy_production" {
    plugin = "helm"
    action = "install"
    depends_on = ["test_staging"]
    
    params = {
      name = "myapp"
      chart = "stable/myapp"
      namespace = "production"
      create_namespace = true
      values_file = "values-production.yaml"
      wait = true
      timeout = "20m"
    }
  }
}
```

### Chart Development Workflow
```hcl
workflow "chart-development" {
  description = "Validate and test Helm chart development"
  
  step "template_chart" {
    plugin = "helm"
    action = "template"
    
    params = {
      name = "test-release"
      chart = "./charts/myapp"
      namespace = "development"
      values = {
        environment = "test"
        debug = true
        replicaCount = 1
      }
    }
  }
  
  step "install_dev" {
    plugin = "helm"
    action = "install"
    depends_on = ["template_chart"]
    
    params = {
      name = "myapp-dev"
      chart = "./charts/myapp"
      namespace = "development"
      create_namespace = true
      values = {
        environment = "development"
        debug = true
        replicaCount = 1
      }
      wait = true
    }
  }
  
  step "test_installation" {
    plugin = "helm"
    action = "status"
    depends_on = ["install_dev"]
    
    params = {
      name = "myapp-dev"
      namespace = "development"
    }
  }
  
  step "cleanup_dev" {
    plugin = "helm"
    action = "uninstall"
    depends_on = ["test_installation"]
    
    params = {
      name = "myapp-dev"
      namespace = "development"
    }
  }
}
```

## Error Handling

The plugin handles common Helm errors gracefully:

- **Chart not found** - Reports missing charts with suggestions
- **Release conflicts** - Handles existing release name conflicts  
- **Values errors** - Validates values format and required fields
- **Timeout errors** - Configurable timeouts for long-running operations
- **Permission errors** - Reports RBAC or authentication failures

## Security

- Supports authentication via kubeconfig files
- No credential storage in workflow definitions  
- Secure handling of repository credentials
- Values can be passed via secure files or environment variables
- Respects Kubernetes RBAC permissions

## Performance

- Efficient JSON parsing of Helm output
- Minimal memory footprint for large charts
- Parallel repository operations
- Configurable timeouts for operations
- Optimized template rendering

## Integration

Works seamlessly with other Corynth plugins:
- **kubernetes plugin** - For additional resource management
- **docker plugin** - For image building and registry operations
- **vault plugin** - For secure credential management
- **terraform plugin** - For infrastructure provisioning
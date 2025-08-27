# GCP Plugin

Production-ready Google Cloud Platform plugin for comprehensive cloud resource management and operations.

## Features

- **Compute Engine** - Create, list, delete VM instances with full configuration
- **GKE Management** - List clusters and manage credentials
- **Cloud Storage** - Upload, download, list buckets and objects
- **Cloud Functions** - Deploy and invoke serverless functions
- **Multi-region** - Support for all GCP regions and zones
- **Project Support** - Multiple GCP project management

## Prerequisites

- `gcloud` CLI installed and configured
- Valid GCP credentials (via gcloud auth, service account, or ADC)
- Appropriate IAM permissions for target services

## Actions

### compute_list
List Compute Engine instances with optional filtering.

```hcl
step "list_instances" {
  plugin = "gcp"
  action = "compute_list"
  
  params = {
    project = "my-project"
    zone = "us-central1-a"
    filter = "status=RUNNING"
  }
}
```

### compute_create
Create new Compute Engine instances.

```hcl
step "create_vm" {
  plugin = "gcp"
  action = "compute_create"
  
  params = {
    name = "web-server-01"
    machine_type = "e2-medium"
    image = "debian-cloud/debian-11"
    zone = "us-central1-a"
    project = "my-project"
    network = "default"
    labels = {
      environment = "production"
      application = "web"
    }
    startup_script = <<EOF
#!/bin/bash
apt-get update
apt-get install -y nginx
systemctl start nginx
EOF
  }
}
```

### compute_delete
Delete Compute Engine instances.

```hcl
step "cleanup_instances" {
  plugin = "gcp"
  action = "compute_delete"
  
  params = {
    names = ["instance-1", "instance-2"]
    zone = "us-central1-a"
    project = "my-project"
  }
}
```

### gke_list
List GKE clusters.

```hcl
step "list_clusters" {
  plugin = "gcp"
  action = "gke_list"
  
  params = {
    project = "my-project"
    location = "us-central1"
  }
}
```

### gke_get_credentials
Get GKE cluster credentials for kubectl.

```hcl
step "get_cluster_creds" {
  plugin = "gcp"
  action = "gke_get_credentials"
  
  params = {
    cluster = "production-cluster"
    location = "us-central1-a"
    project = "my-project"
  }
}
```

### storage_list
List Cloud Storage buckets or objects.

```hcl
step "list_buckets" {
  plugin = "gcp"
  action = "storage_list"
  
  params = {
    project = "my-project"
  }
}

step "list_objects" {
  plugin = "gcp"
  action = "storage_list"
  
  params = {
    bucket = "my-data-bucket"
    prefix = "logs/2024/"
  }
}
```

### storage_upload
Upload files to Cloud Storage.

```hcl
step "upload_artifact" {
  plugin = "gcp"
  action = "storage_upload"
  
  params = {
    file_path = "./build/app.jar"
    bucket = "deployment-artifacts"
    object_name = "releases/v1.2.3/app.jar"
    content_type = "application/java-archive"
    metadata = {
      version = "1.2.3"
      build_date = "2024-01-15"
    }
  }
}
```

### storage_download
Download files from Cloud Storage.

```hcl
step "download_config" {
  plugin = "gcp"
  action = "storage_download"
  
  params = {
    bucket = "config-bucket"
    object_name = "production/app-config.json"
    file_path = "./config/app-config.json"
    project = "my-project"
  }
}
```

### functions_deploy
Deploy Cloud Functions.

```hcl
step "deploy_function" {
  plugin = "gcp"
  action = "functions_deploy"
  
  params = {
    name = "data-processor"
    source = "./functions/processor"
    entry_point = "processData"
    runtime = "nodejs18"
    trigger = "http"
    region = "us-central1"
    memory = "512MB"
    timeout = "120s"
    env_vars = {
      ENVIRONMENT = "production"
      API_KEY = "${vault_secret}"
    }
    project = "my-project"
  }
}
```

### functions_invoke
Invoke Cloud Functions.

```hcl
step "process_data" {
  plugin = "gcp"
  action = "functions_invoke"
  
  params = {
    name = "data-processor"
    data = {
      input_bucket = "raw-data"
      output_bucket = "processed-data"
      batch_id = "batch-001"
    }
    region = "us-central1"
    project = "my-project"
  }
}
```

## Parameters

### Common Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `project` | string | No | - | GCP project ID |

### Compute Engine Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `name` | string | Yes | - | Instance name |
| `machine_type` | string | Yes | - | Machine type (e.g., e2-medium) |
| `image` | string | Yes | - | Boot disk image |
| `zone` | string | Yes | - | GCP zone |
| `network` | string | No | "default" | VPC network name |
| `subnet` | string | No | - | Subnet name |
| `labels` | object | No | - | Instance labels |
| `metadata` | object | No | - | Instance metadata |
| `startup_script` | string | No | - | Startup script content |
| `names` | array | Yes | - | Instance names (delete only) |
| `filter` | string | No | - | Filter expression |

### GKE Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `cluster` | string | Yes | - | Cluster name |
| `location` | string | Yes | - | GCP location (region or zone) |

### Cloud Storage Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `bucket` | string | Yes | - | GCS bucket name |
| `object_name` | string | Yes | - | Object name in bucket |
| `file_path` | string | Yes | - | Local file path |
| `prefix` | string | No | - | Object prefix filter |
| `content_type` | string | No | - | MIME content type |
| `metadata` | object | No | - | Object metadata |

### Cloud Functions Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `name` | string | Yes | - | Function name |
| `source` | string | Yes | - | Source code directory |
| `entry_point` | string | Yes | - | Function entry point |
| `runtime` | string | Yes | - | Runtime (nodejs18, python39, go119) |
| `trigger` | string | Yes | - | Trigger type (http, pubsub, storage) |
| `region` | string | No | "us-central1" | Deployment region |
| `memory` | string | No | "256MB" | Memory allocation |
| `timeout` | string | No | "60s" | Timeout duration |
| `env_vars` | object | No | - | Environment variables |
| `data` | object | No | - | Function input data (invoke) |

## Outputs

### Compute Engine Outputs
- `instances` - List of compute instance objects
- `count` - Number of instances found
- `instance_name` - Created instance name
- `status` - Instance status
- `internal_ip` - Internal IP address
- `external_ip` - External IP address
- `deleted` - List of deleted instance names

### GKE Outputs
- `clusters` - List of GKE clusters
- `count` - Number of clusters found
- `configured` - Whether credentials were configured
- `context` - Kubernetes context name

### Cloud Storage Outputs
- `buckets` - List of GCS buckets
- `objects` - List of GCS objects
- `count` - Number of items found
- `url` - Object GCS URL
- `size` - Object size in bytes
- `md5` - Object MD5 hash

### Cloud Functions Outputs
- `url` - Function trigger URL (HTTP functions)
- `status` - Deployment status
- `version` - Function version
- `response` - Function response payload
- `execution_id` - Execution ID
- `duration` - Execution duration

## Examples

### Complete Application Deployment
```hcl
workflow "gcp-application-deployment" {
  description = "Deploy application infrastructure on GCP"
  
  step "create_instance" {
    plugin = "gcp"
    action = "compute_create"
    
    params = {
      name = "app-server-01"
      machine_type = "n2-standard-2"
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
      zone = "us-central1-a"
      project = "production-project"
      labels = {
        environment = "production"
        application = "webapp"
      }
      startup_script = <<EOF
#!/bin/bash
apt-get update
apt-get install -y docker.io
docker run -d -p 80:8080 myapp:latest
EOF
    }
  }
  
  step "upload_application" {
    plugin = "gcp"
    action = "storage_upload"
    depends_on = ["create_instance"]
    
    params = {
      file_path = "./dist/app.tar.gz"
      bucket = "app-deployments"
      object_name = "releases/v1.0.0/app.tar.gz"
      project = "production-project"
    }
  }
  
  step "deploy_function" {
    plugin = "gcp"
    action = "functions_deploy"
    depends_on = ["upload_application"]
    
    params = {
      name = "app-api"
      source = "./functions/api"
      entry_point = "handleRequest"
      runtime = "nodejs18"
      trigger = "http"
      region = "us-central1"
      memory = "1GB"
      timeout = "300s"
      project = "production-project"
    }
  }
  
  step "verify_deployment" {
    plugin = "gcp"
    action = "compute_list"
    depends_on = ["deploy_function"]
    
    params = {
      project = "production-project"
      zone = "us-central1-a"
      filter = "labels.application=webapp"
    }
  }
}
```

### GKE Cluster Operations
```hcl
workflow "gke-operations" {
  description = "Manage GKE cluster and deploy applications"
  
  step "list_clusters" {
    plugin = "gcp"
    action = "gke_list"
    
    params = {
      project = "k8s-project"
      location = "us-central1"
    }
  }
  
  step "get_credentials" {
    plugin = "gcp"
    action = "gke_get_credentials"
    depends_on = ["list_clusters"]
    
    params = {
      cluster = "production-cluster"
      location = "us-central1-a"
      project = "k8s-project"
    }
  }
  
  step "deploy_to_gke" {
    plugin = "kubernetes"
    action = "apply"
    depends_on = ["get_credentials"]
    
    params = {
      manifest = file("k8s-deployment.yaml")
      namespace = "production"
    }
  }
}
```

### Serverless Data Pipeline
```hcl
workflow "gcp-serverless-pipeline" {
  description = "Serverless data processing pipeline"
  
  step "list_input_files" {
    plugin = "gcp"
    action = "storage_list"
    
    params = {
      bucket = "raw-data"
      prefix = "incoming/"
      project = "data-project"
    }
  }
  
  step "deploy_processor" {
    plugin = "gcp"
    action = "functions_deploy"
    depends_on = ["list_input_files"]
    
    params = {
      name = "batch-processor"
      source = "./functions/processor"
      entry_point = "processBatch"
      runtime = "python39"
      trigger = "storage"
      bucket = "raw-data"
      memory = "2GB"
      timeout = "540s"
      env_vars = {
        OUTPUT_BUCKET = "processed-data"
        LOG_LEVEL = "INFO"
      }
      project = "data-project"
    }
  }
  
  step "trigger_processing" {
    plugin = "gcp"
    action = "functions_invoke"
    depends_on = ["deploy_processor"]
    
    params = {
      name = "batch-processor"
      data = {
        batch_size = 1000
        parallel = true
      }
      region = "us-central1"
      project = "data-project"
    }
  }
  
  step "verify_output" {
    plugin = "gcp"
    action = "storage_list"
    depends_on = ["trigger_processing"]
    
    params = {
      bucket = "processed-data"
      prefix = "output/"
      project = "data-project"
    }
  }
}
```

## Error Handling

The plugin handles common GCP errors gracefully:

- **Authentication errors** - Reports credential or permission issues
- **Resource not found** - Returns appropriate status for missing resources
- **Quota errors** - Reports quota and limit violations
- **Region/zone errors** - Validates location availability
- **API errors** - Handles GCP API errors with detailed messages

## Security

- Supports all GCP authentication methods (gcloud auth, service accounts, ADC)
- No credential storage in workflow definitions
- Respects GCP IAM permissions and policies
- Secure handling of sensitive parameters
- Audit trail through Cloud Audit Logs

## Performance

- Efficient JSON parsing of gcloud output
- Minimal memory footprint for large resource lists
- Parallel operations across regions
- Configurable timeouts for long-running operations
- Optimized for batch operations

## Integration

Works seamlessly with other Corynth plugins:
- **kubernetes plugin** - For GKE cluster management
- **terraform plugin** - For infrastructure provisioning
- **vault plugin** - For secure credential management
- **docker plugin** - For container operations with GCR
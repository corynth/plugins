# AWS Plugin

Production-ready Amazon Web Services plugin for comprehensive cloud resource management and operations.

## Features

- **EC2 Management** - Create, list, terminate instances with full configuration
- **S3 Operations** - Upload, download, list buckets and objects  
- **Lambda Functions** - Invoke functions and list deployments
- **IAM Management** - List users and manage access
- **VPC Operations** - List and manage Virtual Private Clouds
- **Multi-region** - Support for all AWS regions
- **Profile Support** - Multiple AWS credential profiles

## Prerequisites

- `aws` CLI v2.0+ installed and configured
- Valid AWS credentials (via AWS CLI, environment variables, or IAM roles)
- Appropriate IAM permissions for target services

## Actions

### ec2_list
List EC2 instances with optional filtering.

```hcl
step "list_instances" {
  plugin = "aws"
  action = "ec2_list"
  
  params = {
    region = "us-west-2"
    state = "running"
    tags = {
      Environment = "production"
      Application = "web-server"
    }
  }
}
```

### ec2_create
Create new EC2 instances.

```hcl
step "create_server" {
  plugin = "aws"
  action = "ec2_create"
  
  params = {
    image_id = "ami-0abcdef1234567890"
    instance_type = "t3.medium"
    key_name = "my-keypair"
    security_groups = ["sg-12345678"]
    subnet_id = "subnet-12345678"
    region = "us-east-1"
    tags = {
      Name = "web-server-01"
      Environment = "production"
    }
    user_data = <<EOF
#!/bin/bash
yum update -y
yum install -y httpd
systemctl start httpd
systemctl enable httpd
EOF
  }
}
```

### ec2_terminate
Terminate EC2 instances.

```hcl
step "cleanup_instances" {
  plugin = "aws"
  action = "ec2_terminate"
  
  params = {
    instance_ids = ["i-1234567890abcdef0", "i-0987654321fedcba0"]
    region = "us-east-1"
  }
}
```

### s3_list
List S3 buckets or objects.

```hcl
step "list_buckets" {
  plugin = "aws"
  action = "s3_list"
}

step "list_objects" {
  plugin = "aws"
  action = "s3_list"
  
  params = {
    bucket = "my-data-bucket"
    prefix = "logs/2024/"
  }
}
```

### s3_upload
Upload files to S3.

```hcl
step "upload_artifact" {
  plugin = "aws"
  action = "s3_upload"
  
  params = {
    file_path = "./build/app.zip"
    bucket = "deployment-artifacts"
    key = "releases/v1.2.3/app.zip"
    content_type = "application/zip"
    acl = "private"
    metadata = {
      version = "1.2.3"
      build_date = "2024-01-15"
    }
  }
}
```

### s3_download
Download files from S3.

```hcl
step "download_config" {
  plugin = "aws"
  action = "s3_download"
  
  params = {
    bucket = "config-bucket"
    key = "production/app-config.json"
    file_path = "./config/app-config.json"
    region = "us-west-1"
  }
}
```

### lambda_invoke
Invoke Lambda functions.

```hcl
step "process_data" {
  plugin = "aws"
  action = "lambda_invoke"
  
  params = {
    function_name = "data-processor"
    payload = {
      input_bucket = "raw-data"
      output_bucket = "processed-data"
      batch_id = "batch-001"
    }
    invocation_type = "RequestResponse"
    region = "us-east-1"
  }
}
```

### lambda_list
List Lambda functions.

```hcl
step "list_functions" {
  plugin = "aws"
  action = "lambda_list"
  
  params = {
    region = "us-east-1"
  }
}
```

### iam_list_users
List IAM users.

```hcl
step "audit_users" {
  plugin = "aws"
  action = "iam_list_users"
  
  params = {
    path_prefix = "/developers/"
  }
}
```

### vpc_list
List VPCs with optional filtering.

```hcl
step "list_vpcs" {
  plugin = "aws"
  action = "vpc_list"
  
  params = {
    region = "us-west-2"
    filters = {
      "state" = "available"
      "tag:Environment" = "production"
    }
  }
}
```

## Parameters

### Common Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `region` | string | No | "us-east-1" | AWS region |
| `profile` | string | No | - | AWS profile name |

### EC2 Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `image_id` | string | Yes | - | AMI ID for instance creation |
| `instance_type` | string | Yes | - | EC2 instance type |
| `key_name` | string | No | - | SSH key pair name |
| `security_groups` | array | No | - | Security group IDs |
| `subnet_id` | string | No | - | VPC subnet ID |
| `user_data` | string | No | - | Instance initialization script |
| `tags` | object | No | - | Resource tags |
| `state` | string | No | - | Instance state filter |
| `instance_ids` | array | Yes | - | Instance IDs (terminate only) |

### S3 Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `bucket` | string | Yes | - | S3 bucket name |
| `key` | string | Yes | - | S3 object key |
| `file_path` | string | Yes | - | Local file path |
| `prefix` | string | No | - | Object prefix filter |
| `content_type` | string | No | - | MIME content type |
| `acl` | string | No | "private" | Access control list |
| `metadata` | object | No | - | Object metadata |

### Lambda Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `function_name` | string | Yes | - | Function name or ARN |
| `payload` | object | No | - | Function input payload |
| `invocation_type` | string | No | "RequestResponse" | Invocation type |

### IAM Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `path_prefix` | string | No | - | User path prefix filter |

### VPC Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `filters` | object | No | - | VPC filters (key-value pairs) |

## Outputs

### EC2 Outputs
- `instances` - List of EC2 instance objects
- `count` - Number of instances found
- `instance_id` - Created instance ID (create only)
- `state` - Instance state
- `public_ip` - Public IP address
- `private_ip` - Private IP address
- `terminated` - List of terminated instance IDs

### S3 Outputs
- `buckets` - List of S3 buckets
- `objects` - List of S3 objects
- `count` - Number of items found
- `etag` - Object ETag
- `url` - Object S3 URL
- `size` - Object size in bytes
- `last_modified` - Last modified timestamp

### Lambda Outputs
- `functions` - List of Lambda functions
- `count` - Number of functions found
- `response` - Function response payload
- `status_code` - HTTP status code
- `log_result` - Function execution logs

### IAM Outputs
- `users` - List of IAM users
- `count` - Number of users found

### VPC Outputs
- `vpcs` - List of VPC objects
- `count` - Number of VPCs found

## Examples

### Complete Web Application Deployment
```hcl
workflow "aws-web-app-deployment" {
  description = "Deploy web application infrastructure on AWS"
  
  step "create_web_server" {
    plugin = "aws"
    action = "ec2_create"
    
    params = {
      image_id = "ami-0abcdef1234567890"
      instance_type = "t3.medium"
      key_name = "production-key"
      security_groups = ["sg-web-servers"]
      subnet_id = "subnet-public-1a"
      region = "us-east-1"
      tags = {
        Name = "web-server-01"
        Environment = "production"
        Application = "webapp"
      }
      user_data = <<EOF
#!/bin/bash
yum update -y
yum install -y httpd
systemctl start httpd
systemctl enable httpd
echo "<h1>Web Server Ready</h1>" > /var/www/html/index.html
EOF
    }
  }
  
  step "upload_application" {
    plugin = "aws"
    action = "s3_upload"
    depends_on = ["create_web_server"]
    
    params = {
      file_path = "./dist/webapp.zip"
      bucket = "webapp-deployments"
      key = "releases/v1.0.0/webapp.zip"
      content_type = "application/zip"
      region = "us-east-1"
      metadata = {
        version = "1.0.0"
        deployed_by = "corynth"
      }
    }
  }
  
  step "process_deployment" {
    plugin = "aws"
    action = "lambda_invoke"
    depends_on = ["upload_application"]
    
    params = {
      function_name = "deployment-processor"
      region = "us-east-1"
      payload = {
        instance_id = "${create_web_server.instance_id}"
        artifact_url = "${upload_application.url}"
        environment = "production"
      }
    }
  }
  
  step "verify_deployment" {
    plugin = "aws"
    action = "ec2_list"
    depends_on = ["process_deployment"]
    
    params = {
      region = "us-east-1"
      tags = {
        Name = "web-server-01"
      }
    }
  }
}
```

### Multi-Region Data Backup
```hcl
workflow "aws-multi-region-backup" {
  description = "Backup data across multiple AWS regions"
  
  step "list_source_objects" {
    plugin = "aws"
    action = "s3_list"
    
    params = {
      bucket = "production-data"
      prefix = "daily-backups/"
      region = "us-east-1"
    }
  }
  
  step "create_backup_instance_west" {
    plugin = "aws"
    action = "ec2_create"
    depends_on = ["list_source_objects"]
    
    params = {
      image_id = "ami-backup-tools"
      instance_type = "t3.large"
      key_name = "backup-key"
      region = "us-west-2"
      tags = {
        Name = "backup-processor-west"
        Purpose = "data-backup"
      }
    }
  }
  
  step "process_backup" {
    plugin = "aws"
    action = "lambda_invoke"
    depends_on = ["create_backup_instance_west"]
    
    params = {
      function_name = "cross-region-backup"
      region = "us-west-2"
      payload = {
        source_bucket = "production-data"
        source_region = "us-east-1"
        target_bucket = "backup-data-west"
        target_region = "us-west-2"
        object_count = "${list_source_objects.count}"
      }
    }
  }
  
  step "cleanup_backup_instance" {
    plugin = "aws"
    action = "ec2_terminate"
    depends_on = ["process_backup"]
    
    params = {
      instance_ids = ["${create_backup_instance_west.instance_id}"]
      region = "us-west-2"
    }
  }
}
```

### Infrastructure Audit Workflow
```hcl
workflow "aws-infrastructure-audit" {
  description = "Audit AWS infrastructure across services"
  
  step "audit_ec2_instances" {
    plugin = "aws"
    action = "ec2_list"
    
    params = {
      region = "us-east-1"
      state = "running"
    }
  }
  
  step "audit_s3_buckets" {
    plugin = "aws"
    action = "s3_list"
    depends_on = ["audit_ec2_instances"]
  }
  
  step "audit_lambda_functions" {
    plugin = "aws"
    action = "lambda_list"
    depends_on = ["audit_s3_buckets"]
    
    params = {
      region = "us-east-1"
    }
  }
  
  step "audit_iam_users" {
    plugin = "aws"
    action = "iam_list_users"
    depends_on = ["audit_lambda_functions"]
  }
  
  step "audit_vpcs" {
    plugin = "aws"
    action = "vpc_list"
    depends_on = ["audit_iam_users"]
    
    params = {
      region = "us-east-1"
    }
  }
  
  step "upload_audit_report" {
    plugin = "aws"
    action = "s3_upload"
    depends_on = ["audit_vpcs"]
    
    params = {
      file_path = "./reports/infrastructure-audit.json"
      bucket = "compliance-reports"
      key = "audits/infrastructure/$(date +%Y-%m-%d).json"
      content_type = "application/json"
      acl = "private"
      metadata = {
        audit_date = "$(date -Iseconds)"
        instances_count = "${audit_ec2_instances.count}"
        buckets_count = "${audit_s3_buckets.count}"
        functions_count = "${audit_lambda_functions.count}"
        users_count = "${audit_iam_users.count}"
        vpcs_count = "${audit_vpcs.count}"
      }
    }
  }
}
```

### Auto-Scaling Application Deployment
```hcl
workflow "aws-auto-scaling-deployment" {
  description = "Deploy application with auto-scaling infrastructure"
  
  step "create_load_balancer_instance" {
    plugin = "aws"
    action = "ec2_create"
    
    params = {
      image_id = "ami-nginx-lb"
      instance_type = "t3.medium"
      key_name = "production-key"
      security_groups = ["sg-load-balancers"]
      subnet_id = "subnet-public-1a"
      region = "us-east-1"
      tags = {
        Name = "load-balancer-01"
        Role = "load-balancer"
        Environment = "production"
      }
    }
  }
  
  step "create_app_servers" {
    plugin = "aws"
    action = "ec2_create"
    depends_on = ["create_load_balancer_instance"]
    
    params = {
      image_id = "ami-app-server"
      instance_type = "t3.large"
      key_name = "production-key"
      security_groups = ["sg-app-servers"]
      subnet_id = "subnet-private-1a"
      region = "us-east-1"
      tags = {
        Name = "app-server-01"
        Role = "application"
        Environment = "production"
        LoadBalancer = "${create_load_balancer_instance.instance_id}"
      }
    }
  }
  
  step "deploy_application_code" {
    plugin = "aws"
    action = "lambda_invoke"
    depends_on = ["create_app_servers"]
    
    params = {
      function_name = "application-deployer"
      region = "us-east-1"
      payload = {
        load_balancer_id = "${create_load_balancer_instance.instance_id}"
        app_server_id = "${create_app_servers.instance_id}"
        deployment_bucket = "app-deployments"
        config_key = "production/app-config.json"
      }
    }
  }
  
  step "health_check" {
    plugin = "aws"
    action = "ec2_list"
    depends_on = ["deploy_application_code"]
    
    params = {
      region = "us-east-1"
      tags = {
        Environment = "production"
        Role = "application"
      }
      state = "running"
    }
  }
}
```

## Error Handling

The plugin handles common AWS errors gracefully:

- **Authentication errors** - Reports credential or permission issues
- **Resource not found** - Returns appropriate status for missing resources  
- **Rate limiting** - Handles AWS API throttling
- **Region errors** - Validates region availability
- **Service limits** - Reports quota and limit violations

## Security

- Supports all AWS authentication methods (CLI profiles, environment variables, IAM roles)
- No credential storage in workflow definitions
- Respects AWS IAM permissions and policies
- Secure handling of sensitive parameters
- Audit trail through AWS CloudTrail integration

## Performance

- Efficient JSON parsing of AWS CLI output
- Minimal memory footprint for large resource lists
- Parallel operations across regions
- Configurable timeouts for long-running operations
- Optimized for batch operations

## Integration

Works seamlessly with other Corynth plugins:
- **kubernetes plugin** - For EKS cluster management
- **terraform plugin** - For infrastructure provisioning
- **vault plugin** - For secure credential management  
- **docker plugin** - For container operations with ECR
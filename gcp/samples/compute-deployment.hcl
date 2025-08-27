workflow "gcp-compute-deployment" {
  description = "Deploy and manage Compute Engine instances"
  version = "1.0.0"

  step "list_existing_vms" {
    plugin = "gcp"
    action = "compute_list"
    
    params = {
      project = "my-project"
      zone = "us-central1-a"
      filter = "status=RUNNING"
    }
  }

  step "create_web_server" {
    plugin = "gcp"
    action = "compute_create"
    depends_on = ["list_existing_vms"]
    
    params = {
      name = "web-server-prod"
      machine_type = "e2-standard-2"
      image = "debian-cloud/debian-11"
      zone = "us-central1-a"
      project = "my-project"
      network = "default"
      labels = {
        environment = "production"
        application = "web"
        team = "platform"
      }
      metadata = {
        enable-oslogin = "TRUE"
      }
      startup_script = <<EOF
#!/bin/bash
apt-get update
apt-get install -y nginx
cat > /var/www/html/index.html <<HTML
<!DOCTYPE html>
<html>
<head><title>GCP Web Server</title></head>
<body>
  <h1>Production Web Server on GCP</h1>
  <p>Instance: $(hostname)</p>
  <p>Zone: us-central1-a</p>
</body>
</html>
HTML
systemctl restart nginx
EOF
    }
  }

  step "upload_logs" {
    plugin = "gcp"
    action = "storage_upload"
    depends_on = ["create_web_server"]
    
    params = {
      file_path = "./logs/deployment.log"
      bucket = "deployment-logs"
      object_name = "compute/$(date +%Y-%m-%d)/web-server-prod.log"
      content_type = "text/plain"
      project = "my-project"
      metadata = {
        instance_name = "${create_web_server.instance_name}"
        deployment_date = "$(date -Iseconds)"
      }
    }
  }

  step "verify_deployment" {
    plugin = "gcp"
    action = "compute_list"
    depends_on = ["upload_logs"]
    
    params = {
      project = "my-project"
      zone = "us-central1-a"
      filter = "name=web-server-prod AND status=RUNNING"
    }
  }
}
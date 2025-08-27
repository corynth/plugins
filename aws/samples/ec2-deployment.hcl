workflow "aws-ec2-deployment" {
  description = "Deploy and manage EC2 instances with full lifecycle"
  version = "1.0.0"

  step "list_existing_instances" {
    plugin = "aws"
    action = "ec2_list"
    
    params = {
      region = "us-east-1"
      state = "running"
      tags = {
        Environment = "production"
        Application = "web-server"
      }
    }
  }

  step "create_web_server" {
    plugin = "aws"
    action = "ec2_create"
    depends_on = ["list_existing_instances"]
    
    params = {
      image_id = "ami-0abcdef1234567890"
      instance_type = "t3.medium"
      key_name = "production-keypair"
      security_groups = ["sg-12345678"]
      subnet_id = "subnet-12345678"
      region = "us-east-1"
      tags = {
        Name = "web-server-01"
        Environment = "production"
        Application = "web-server"
        CreatedBy = "corynth-workflow"
      }
      user_data = <<EOF
#!/bin/bash
yum update -y
yum install -y httpd
systemctl start httpd
systemctl enable httpd
echo "<h1>Production Web Server</h1>" > /var/www/html/index.html
echo "<p>Instance ID: $(curl -s http://169.254.169.254/latest/meta-data/instance-id)</p>" >> /var/www/html/index.html
EOF
    }
  }

  step "verify_instance_creation" {
    plugin = "aws"
    action = "ec2_list"
    depends_on = ["create_web_server"]
    
    params = {
      region = "us-east-1"
      tags = {
        Name = "web-server-01"
      }
      state = "running"
    }
  }

  step "upload_application_logs" {
    plugin = "aws"
    action = "s3_upload"
    depends_on = ["verify_instance_creation"]
    
    params = {
      file_path = "./logs/deployment.log"
      bucket = "application-logs"
      key = "deployments/web-server-01/$(date +%Y-%m-%d)/deployment.log"
      content_type = "text/plain"
      acl = "private"
      region = "us-east-1"
      metadata = {
        instance_id = "${create_web_server.instance_id}"
        deployment_date = "$(date -Iseconds)"
        application = "web-server"
      }
    }
  }
}
workflow "infrastructure-status-report" {
  description = "Generate comprehensive infrastructure status report"
  version = "1.0.0"

  step "collect_aws_instances" {
    plugin = "aws"
    action = "ec2_list"
    
    params = {
      region = "us-east-1"
      state = "running"
    }
  }

  step "collect_gcp_instances" {
    plugin = "gcp"
    action = "compute_list"
    depends_on = ["collect_aws_instances"]
    
    params = {
      project = "my-project"
      zone = "us-central1-a"
    }
  }

  step "create_instance_table" {
    plugin = "reporting"
    action = "create_table"
    depends_on = ["collect_gcp_instances"]
    
    params = {
      title = "Cloud Instance Inventory"
      headers = ["Provider", "Instance ID", "Type", "State", "Region"]
      rows = [
        ["AWS", "i-1234567890", "t3.medium", "running", "us-east-1"],
        ["AWS", "i-0987654321", "t3.large", "running", "us-east-1"],
        ["GCP", "web-server-01", "e2-medium", "RUNNING", "us-central1-a"]
      ]
      format = "markdown"
      alignment = ["left", "center", "center", "center", "center"]
    }
  }

  step "create_usage_chart" {
    plugin = "reporting"
    action = "create_chart"
    depends_on = ["create_instance_table"]
    
    params = {
      type = "bar"
      title = "Resource Usage by Service"
      data = {
        labels = ["Web Servers", "Databases", "Cache", "Load Balancers"]
        values = [65.3, 78.9, 42.1, 23.7]
      }
      output_file = "./charts/resource-usage.txt"
    }
  }

  step "generate_full_report" {
    plugin = "reporting"
    action = "create_report"
    depends_on = ["create_usage_chart"]
    
    params = {
      title = "Weekly Infrastructure Status Report"
      format = "html"
      output_file = "./reports/infrastructure-status.html"
      template = "technical"
      metadata = {
        author = "DevOps Team"
        date = "$(date +%Y-%m-%d)"
        classification = "Internal"
        version = "1.0"
      }
      sections = [
        {
          heading = "Executive Summary"
          content = "Infrastructure overview for week ending $(date +%Y-%m-%d). All critical services operational with ${collect_aws_instances.count} AWS instances and ${collect_gcp_instances.count} GCP instances running."
        },
        {
          heading = "Instance Inventory"
          content = "${create_instance_table.table}"
        },
        {
          heading = "Resource Utilization"
          content = "Current resource usage across all services:\n\n$(cat ./charts/resource-usage.txt)"
        },
        {
          heading = "Key Metrics"
          list = [
            "Total Cloud Instances: ${collect_aws_instances.count + collect_gcp_instances.count}",
            "Average CPU Utilization: 62%", 
            "Average Memory Usage: 74%",
            "Network Throughput: 2.3 Gbps",
            "Storage Usage: 1.2 TB"
          ]
        },
        {
          heading = "Action Items"
          content = "## Immediate Actions\n- Scale web servers during peak hours\n- Review cache configuration\n- Update security groups\n\n## Planned Activities\n- Migrate legacy instances\n- Implement auto-scaling policies\n- Upgrade monitoring system"
        },
        {
          heading = "Compliance Status"
          table = {
            headers = ["Check", "Status", "Last Audit", "Next Review"]
            rows = [
              ["Security Groups", "✅ Pass", "2024-01-10", "2024-02-10"],
              ["Backup Validation", "✅ Pass", "2024-01-12", "2024-02-12"],
              ["SSL Certificates", "⚠️ Expiring", "2024-01-15", "2024-01-30"],
              ["Access Reviews", "✅ Pass", "2024-01-08", "2024-02-08"]
            ]
          }
        }
      ]
    }
  }

  step "convert_to_pdf" {
    plugin = "reporting"
    action = "convert_format"
    depends_on = ["generate_full_report"]
    
    params = {
      input_file = "./reports/infrastructure-status.html"
      output_format = "pdf"
      output_file = "./reports/infrastructure-status.pdf"
    }
  }

  step "display_summary" {
    plugin = "reporting"
    action = "display"
    depends_on = ["convert_to_pdf"]
    
    params = {
      content = "📊 Infrastructure Report Generated\n\n✅ HTML Report: ${generate_full_report.file_path}\n✅ PDF Report: ${convert_to_pdf.output_file}\n\n📈 Summary:\n- Total Instances: ${collect_aws_instances.count + collect_gcp_instances.count}\n- Report Format: ${generate_full_report.format}\n- Charts Created: 1"
      format = "text"
    }
  }
}
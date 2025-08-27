# Reporting Plugin

Production-ready reporting plugin for generating formatted reports, tables, and charts with multiple output formats.

## Features

- **Multi-format Reports** - Generate reports in Markdown, HTML, PDF, and text formats
- **Table Generation** - Create formatted tables with alignment and styling options
- **Chart Visualization** - Generate ASCII charts for terminal display
- **Format Conversion** - Convert between different document formats
- **Template Support** - Multiple report templates (default, technical, executive)
- **Screen Display** - Display formatted content directly in terminal

## Prerequisites

- Optional: `pandoc` for advanced format conversion
- Optional: `wkhtmltopdf` for PDF generation

## Actions

### create_report
Generate a comprehensive formatted report.

```hcl
step "generate_report" {
  plugin = "reporting"
  action = "create_report"
  
  params = {
    title = "Infrastructure Audit Report"
    format = "markdown"
    output_file = "./reports/audit-report.md"
    template = "technical"
    metadata = {
      author = "DevOps Team"
      date = "2024-01-15"
      version = "1.0"
    }
    sections = [
      {
        heading = "Executive Summary"
        content = "This report provides an overview of our infrastructure status."
      },
      {
        heading = "Server Inventory"
        table = {
          headers = ["Server", "Status", "CPU", "Memory"]
          rows = [
            ["web-01", "Running", "2.3 GHz", "8GB"],
            ["db-01", "Running", "3.1 GHz", "16GB"],
            ["cache-01", "Stopped", "2.1 GHz", "4GB"]
          ]
        }
      },
      {
        heading = "Recommendations"
        list = [
          "Restart cache-01 server",
          "Upgrade memory on web servers",
          "Implement monitoring alerts"
        ]
      }
    ]
  }
}
```

### create_table
Generate formatted tables.

```hcl
step "create_status_table" {
  plugin = "reporting"
  action = "create_table"
  
  params = {
    title = "Service Status"
    headers = ["Service", "Status", "Uptime", "Response Time"]
    rows = [
      ["API Gateway", "Healthy", "99.9%", "45ms"],
      ["Database", "Healthy", "99.7%", "12ms"],
      ["Cache", "Warning", "98.2%", "8ms"]
    ]
    format = "markdown"
    alignment = ["left", "center", "right", "right"]
  }
}
```

### create_chart
Generate chart visualizations.

```hcl
step "create_usage_chart" {
  plugin = "reporting"
  action = "create_chart"
  
  params = {
    type = "bar"
    title = "CPU Usage by Server"
    data = {
      labels = ["web-01", "web-02", "db-01", "cache-01"]
      values = [45.2, 38.7, 72.1, 23.5]
    }
    width = 800
    height = 400
    output_file = "./charts/cpu-usage.txt"
  }
}
```

### convert_format
Convert reports between formats.

```hcl
step "convert_to_pdf" {
  plugin = "reporting"
  action = "convert_format"
  
  params = {
    input_file = "./reports/audit-report.md"
    output_format = "pdf"
    output_file = "./reports/audit-report.pdf"
    options = {
      page_size = "A4"
      margin = "1in"
    }
  }
}
```

### display
Display content on screen.

```hcl
step "show_summary" {
  plugin = "reporting"
  action = "display"
  
  params = {
    content = "${generate_report.content}"
    format = "markdown"
    style = {
      highlight = true
      colors = true
    }
  }
}
```

## Parameters

### create_report Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `title` | string | Yes | - | Report title |
| `sections` | array | Yes | - | Report sections with content |
| `format` | string | No | "markdown" | Output format (markdown, html, pdf, text) |
| `output_file` | string | No | - | Output file path |
| `template` | string | No | "default" | Report template |
| `metadata` | object | No | - | Report metadata (author, date, version) |

### create_table Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `headers` | array | Yes | - | Table column headers |
| `rows` | array | Yes | - | Table data rows |
| `format` | string | No | "markdown" | Table format (markdown, ascii, html, csv) |
| `alignment` | array | No | - | Column alignments (left, center, right) |
| `title` | string | No | - | Table title |

### create_chart Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | string | Yes | - | Chart type (bar, line, pie, scatter) |
| `data` | object | Yes | - | Chart data (labels and values) |
| `title` | string | No | - | Chart title |
| `width` | number | No | 800 | Chart width in pixels |
| `height` | number | No | 600 | Chart height in pixels |
| `output_file` | string | No | - | Output image file path |

## Outputs

### create_report
- `content` - Generated report content
- `file_path` - Path to saved report file
- `format` - Report format used

### create_table
- `table` - Formatted table output
- `row_count` - Number of data rows
- `column_count` - Number of columns

### create_chart
- `chart_url` - Chart image URL or path
- `format` - Chart output format

### convert_format
- `output_file` - Path to converted file
- `status` - Conversion status

### display
- `displayed` - Whether content was displayed

## Examples

### Complete Infrastructure Report
```hcl
workflow "infrastructure-report" {
  description = "Generate comprehensive infrastructure status report"
  
  step "collect_server_data" {
    plugin = "aws"
    action = "ec2_list"
    
    params = {
      region = "us-east-1"
      state = "running"
    }
  }
  
  step "create_server_table" {
    plugin = "reporting"
    action = "create_table"
    depends_on = ["collect_server_data"]
    
    params = {
      title = "Server Inventory"
      headers = ["Instance ID", "Type", "State", "Launch Time"]
      rows = "${collect_server_data.instances}"
      format = "markdown"
    }
  }
  
  step "generate_full_report" {
    plugin = "reporting"
    action = "create_report"
    depends_on = ["create_server_table"]
    
    params = {
      title = "Weekly Infrastructure Report"
      format = "html"
      output_file = "./reports/weekly-report.html"
      template = "technical"
      metadata = {
        author = "Infrastructure Team"
        date = "$(date +%Y-%m-%d)"
        classification = "Internal"
      }
      sections = [
        {
          heading = "Executive Summary"
          content = "Infrastructure status for the week ending $(date +%Y-%m-%d)"
        },
        {
          heading = "Server Status"
          content = "${create_server_table.table}"
        },
        {
          heading = "Metrics Overview"
          list = [
            "Total Servers: ${collect_server_data.count}",
            "Running Instances: ${collect_server_data.count}",
            "Average CPU: 45%",
            "Average Memory: 62%"
          ]
        },
        {
          heading = "Action Items"
          content = "- Review scaling policies\n- Update security groups\n- Plan maintenance windows"
        }
      ]
    }
  }
  
  step "convert_to_pdf" {
    plugin = "reporting"
    action = "convert_format"
    depends_on = ["generate_full_report"]
    
    params = {
      input_file = "./reports/weekly-report.html"
      output_format = "pdf"
      output_file = "./reports/weekly-report.pdf"
    }
  }
}
```

### Performance Dashboard
```hcl
workflow "performance-dashboard" {
  description = "Generate performance metrics dashboard"
  
  step "create_cpu_chart" {
    plugin = "reporting"
    action = "create_chart"
    
    params = {
      type = "bar"
      title = "CPU Usage by Service"
      data = {
        labels = ["Web", "API", "Database", "Cache"]
        values = [65.3, 42.1, 78.9, 23.7]
      }
      output_file = "./charts/cpu-usage.txt"
    }
  }
  
  step "create_memory_chart" {
    plugin = "reporting"
    action = "create_chart"
    depends_on = ["create_cpu_chart"]
    
    params = {
      type = "pie"
      title = "Memory Distribution"
      data = {
        labels = ["Application", "Cache", "Buffer", "Free"]
        values = [45.2, 23.1, 8.7, 23.0]
      }
      output_file = "./charts/memory-dist.txt"
    }
  }
  
  step "generate_dashboard" {
    plugin = "reporting"
    action = "create_report"
    depends_on = ["create_memory_chart"]
    
    params = {
      title = "Performance Dashboard"
      format = "markdown"
      output_file = "./dashboards/performance.md"
      sections = [
        {
          heading = "System Overview"
          content = "Real-time performance metrics as of $(date)"
        },
        {
          heading = "CPU Usage"
          content = "$(cat ./charts/cpu-usage.txt)"
        },
        {
          heading = "Memory Usage"
          content = "$(cat ./charts/memory-dist.txt)"
        }
      ]
    }
  }
  
  step "display_dashboard" {
    plugin = "reporting"
    action = "display"
    depends_on = ["generate_dashboard"]
    
    params = {
      content = "${generate_dashboard.content}"
      format = "markdown"
    }
  }
}
```

## Report Templates

### Default Template
- Clean, simple formatting
- Standard section headings
- Basic styling

### Technical Template
- Detailed table of contents
- Code block formatting
- Technical metadata
- Appendix support

### Executive Template
- Executive summary focus
- High-level metrics
- Minimal technical details
- Action items emphasis

## Format Support

### Markdown
- GitHub-flavored markdown
- Table support with alignment
- Code syntax highlighting
- Metadata frontmatter

### HTML
- Responsive design
- CSS styling included
- Table formatting
- Print-friendly layout

### PDF
- Professional formatting
- Page breaks
- Headers and footers
- Table of contents

### CSV
- Data export format
- Excel compatibility
- Clean table structure

## Integration

Works seamlessly with other Corynth plugins:
- **aws/gcp plugins** - For infrastructure data collection
- **kubernetes plugin** - For cluster status reporting
- **docker plugin** - For container metrics
- **vault plugin** - For security compliance reports
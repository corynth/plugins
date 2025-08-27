# Reporting Plugin

A comprehensive reporting plugin for Corynth that generates reports in multiple formats including PDF, HTML, Markdown, JSON, and YAML.

## Features

- Generate reports in 5 different formats
- Support for metadata inclusion
- Format conversion capabilities
- Clean, professional output styling
- Cross-platform compatibility

## Prerequisites

- Go 1.21 or later
- Internet connection for PDF generation dependencies

## Installation

The plugin is automatically built and installed when you run:

```bash
make build-plugins
```

Or build manually:

```bash
cd plugins-src/reporting
go build -o ../../bin/plugins/corynth-plugin-reporting main.go
```

## Actions

### generate

Generate a report in the specified format.

**Parameters:**
- `title` (required): Report title
- `content` (required): Report content (markdown supported)  
- `format` (optional): Output format - pdf, html, markdown, json, yaml (default: markdown)
- `output_path` (optional): Output file path (default: ./report)
- `metadata` (optional): Additional metadata object to include

**Returns:**
- `file_path`: Path to generated report file
- `size`: File size in bytes
- `format`: Generated format

**Example:**
```hcl
step "create_report" {
  plugin = "reporting"
  action = "generate"
  params = {
    title = "System Status Report"
    content = "## Overview\nAll systems operational.\n\n## Details\n- CPU: 15%\n- Memory: 45%\n- Disk: 67%"
    format = "pdf"
    output_path = "/tmp/status-report.pdf"
    metadata = {
      author = "Corynth Automation"
      department = "DevOps"
      classification = "Internal"
    }
  }
}
```

### convert

Convert an existing report from one format to another.

**Parameters:**
- `input_file` (required): Input file path
- `output_format` (required): Target format - pdf, html, markdown, json, yaml
- `output_path` (optional): Output file path (defaults to input path with new extension)

**Returns:**
- `file_path`: Path to converted file

**Example:**
```hcl
step "convert_to_pdf" {
  plugin = "reporting"
  action = "convert"
  params = {
    input_file = "/tmp/report.md"
    output_format = "pdf"
    output_path = "/tmp/report.pdf"
  }
}
```

## Sample Workflows

### Multi-Format Report Generation
```hcl
workflow "comprehensive_reporting" {
  description = "Generate reports in multiple formats"
  
  step "generate_markdown" {
    plugin = "reporting"
    action = "generate"
    params = {
      title = "Infrastructure Health Report"
      content = "## Summary\nInfrastructure is healthy and operational.\n\n## Metrics\n- Uptime: 99.9%\n- Response Time: 250ms\n- Error Rate: 0.01%"
      format = "markdown"
      output_path = "/tmp/health-report.md"
    }
  }
  
  step "convert_to_pdf" {
    plugin = "reporting"
    action = "convert"
    depends_on = ["generate_markdown"]
    params = {
      input_file = "/tmp/health-report.md"
      output_format = "pdf"
      output_path = "/tmp/health-report.pdf"
    }
  }
  
  step "convert_to_html" {
    plugin = "reporting"
    action = "convert"
    depends_on = ["generate_markdown"]
    params = {
      input_file = "/tmp/health-report.md"
      output_format = "html"
      output_path = "/tmp/health-report.html"
    }
  }
}
```

### Data Processing Report
```hcl
workflow "data_analysis_report" {
  description = "Generate data analysis report with structured output"
  
  step "analysis_report" {
    plugin = "reporting"
    action = "generate"
    params = {
      title = "Daily Analytics Report"
      content = "## Key Metrics\n- Users: 1,250\n- Sessions: 3,400\n- Revenue: $12,500\n\n## Trends\n- User growth: +5.2%\n- Conversion rate: 3.1%"
      format = "json"
      output_path = "/tmp/analytics.json"
      metadata = {
        date = "2024-01-15"
        source = "Analytics System"
        period = "24h"
      }
    }
  }
}
```

## Output Formats

### PDF
Professional PDF reports with:
- Title headers
- Metadata sections
- Formatted content
- Generated timestamps

### HTML
Web-ready HTML with:
- Responsive design
- Clean typography
- Metadata display
- Professional styling

### Markdown
Clean markdown format with:
- Structured headers
- Metadata lists
- Proper formatting
- Timestamp footer

### JSON
Structured JSON output with:
- Title and content fields
- Metadata object
- Generation timestamps
- Format identification

### YAML
YAML format with:
- Hierarchical structure
- Readable formatting
- Complete metadata
- Standard compliance

## Error Handling

The plugin provides detailed error messages for:
- Invalid format specifications
- File system permission issues
- Content parsing errors
- Dependency availability

## Dependencies

- `github.com/go-pdf/fpdf` - PDF generation
- `gopkg.in/yaml.v3` - YAML processing
- Standard Go libraries for HTML/JSON/Markdown

## License

MIT License - See LICENSE file for details.
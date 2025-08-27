package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"gopkg.in/yaml.v3"

	"github.com/corynth/corynth/pkg/plugin"
	pluginv2 "github.com/corynth/corynth/pkg/plugin/v2"
)

type ReportingPlugin struct{}

func (p *ReportingPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "reporting",
		Version:     "1.0.0",
		Description: "Generate reports in multiple formats (PDF, HTML, Markdown, JSON, YAML)",
		Author:      "Corynth Team",
		Tags:        []string{"reporting", "documentation", "export", "formats"},
		License:     "MIT",
	}
}

func (p *ReportingPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "generate",
			Description: "Generate a report in specified format",
			Inputs: map[string]plugin.InputSpec{
				"title": {
					Type:        "string",
					Description: "Report title",
					Required:    true,
				},
				"content": {
					Type:        "string",
					Description: "Report content (markdown supported)",
					Required:    true,
				},
				"format": {
					Type:        "string",
					Description: "Output format (pdf, html, markdown, json, yaml)",
					Required:    false,
					Default:     "markdown",
				},
				"output_path": {
					Type:        "string",
					Description: "Output file path",
					Required:    false,
					Default:     "./report",
				},
				"metadata": {
					Type:        "object",
					Description: "Additional metadata to include",
					Required:    false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"file_path": {
					Type:        "string",
					Description: "Path to generated report file",
				},
				"size": {
					Type:        "number",
					Description: "File size in bytes",
				},
				"format": {
					Type:        "string",
					Description: "Generated format",
				},
			},
		},
		{
			Name:        "convert",
			Description: "Convert report from one format to another",
			Inputs: map[string]plugin.InputSpec{
				"input_file": {
					Type:        "string",
					Description: "Input file path",
					Required:    true,
				},
				"output_format": {
					Type:        "string",
					Description: "Target format (pdf, html, markdown, json, yaml)",
					Required:    true,
				},
				"output_path": {
					Type:        "string",
					Description: "Output file path",
					Required:    false,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"file_path": {
					Type:        "string",
					Description: "Path to converted file",
				},
			},
		},
	}
}

func (p *ReportingPlugin) Validate(params map[string]interface{}) error {
	format, ok := params["format"].(string)
	if ok {
		validFormats := []string{"pdf", "html", "markdown", "json", "yaml"}
		valid := false
		for _, vf := range validFormats {
			if format == vf {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid format '%s', must be one of: %s", format, strings.Join(validFormats, ", "))
		}
	}
	return nil
}

func (p *ReportingPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "generate":
		return p.generateReport(params)
	case "convert":
		return p.convertReport(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *ReportingPlugin) generateReport(params map[string]interface{}) (map[string]interface{}, error) {
	title := params["title"].(string)
	content := params["content"].(string)
	format := "markdown"
	if f, ok := params["format"].(string); ok {
		format = f
	}
	outputPath := "./report"
	if op, ok := params["output_path"].(string); ok {
		outputPath = op
	}

	var metadata map[string]interface{}
	if m, ok := params["metadata"].(map[string]interface{}); ok {
		metadata = m
	}

	// Generate report based on format
	var filePath string
	var err error

	switch format {
	case "pdf":
		filePath, err = p.generatePDF(title, content, outputPath, metadata)
	case "html":
		filePath, err = p.generateHTML(title, content, outputPath, metadata)
	case "markdown":
		filePath, err = p.generateMarkdown(title, content, outputPath, metadata)
	case "json":
		filePath, err = p.generateJSON(title, content, outputPath, metadata)
	case "yaml":
		filePath, err = p.generateYAML(title, content, outputPath, metadata)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate %s report: %w", format, err)
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return map[string]interface{}{
		"file_path": filePath,
		"size":      fileInfo.Size(),
		"format":    format,
		"status":    "success",
	}, nil
}

func (p *ReportingPlugin) convertReport(params map[string]interface{}) (map[string]interface{}, error) {
	inputFile := params["input_file"].(string)
	outputFormat := params["output_format"].(string)
	
	outputPath := inputFile
	if op, ok := params["output_path"].(string); ok {
		outputPath = op
	}

	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	// Convert to target format
	filePath, err := p.convertToFormat(string(content), outputFormat, outputPath)
	if err != nil {
		return nil, fmt.Errorf("conversion failed: %w", err)
	}

	return map[string]interface{}{
		"file_path": filePath,
		"format":    outputFormat,
		"status":    "success",
	}, nil
}

func (p *ReportingPlugin) generatePDF(title, content, outputPath string, metadata map[string]interface{}) (string, error) {
	if !strings.HasSuffix(outputPath, ".pdf") {
		outputPath += ".pdf"
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	
	// Add title
	pdf.Cell(190, 10, title)
	pdf.Ln(15)
	
	// Add metadata if provided
	if metadata != nil {
		pdf.SetFont("Arial", "I", 10)
		for key, value := range metadata {
			pdf.Cell(190, 5, fmt.Sprintf("%s: %v", key, value))
			pdf.Ln(5)
		}
		pdf.Ln(5)
	}
	
	// Add content
	pdf.SetFont("Arial", "", 12)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		pdf.Cell(190, 5, line)
		pdf.Ln(5)
	}
	
	return outputPath, pdf.OutputFileAndClose(outputPath)
}

func (p *ReportingPlugin) generateHTML(title, content, outputPath string, metadata map[string]interface{}) (string, error) {
	if !strings.HasSuffix(outputPath, ".html") {
		outputPath += ".html"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #333; border-bottom: 2px solid #007acc; }
        .metadata { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 4px; }
        .content { line-height: 1.6; }
        pre { background: #f8f8f8; padding: 10px; border-radius: 4px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>%s</h1>`, title, title)

	if metadata != nil {
		html += `<div class="metadata"><h3>Report Metadata</h3>`
		for key, value := range metadata {
			html += fmt.Sprintf("<p><strong>%s:</strong> %v</p>", key, value)
		}
		html += `</div>`
	}

	html += fmt.Sprintf(`<div class="content">
        <pre>%s</pre>
    </div>
    <footer>
        <p><em>Generated on %s</em></p>
    </footer>
</body>
</html>`, content, time.Now().Format("2006-01-02 15:04:05"))

	return outputPath, os.WriteFile(outputPath, []byte(html), 0644)
}

func (p *ReportingPlugin) generateMarkdown(title, content, outputPath string, metadata map[string]interface{}) (string, error) {
	if !strings.HasSuffix(outputPath, ".md") {
		outputPath += ".md"
	}

	md := fmt.Sprintf("# %s\n\n", title)
	
	if metadata != nil {
		md += "## Metadata\n\n"
		for key, value := range metadata {
			md += fmt.Sprintf("- **%s**: %v\n", key, value)
		}
		md += "\n"
	}
	
	md += "## Content\n\n"
	md += content
	md += fmt.Sprintf("\n\n---\n*Generated on %s*\n", time.Now().Format("2006-01-02 15:04:05"))

	return outputPath, os.WriteFile(outputPath, []byte(md), 0644)
}

func (p *ReportingPlugin) generateJSON(title, content, outputPath string, metadata map[string]interface{}) (string, error) {
	if !strings.HasSuffix(outputPath, ".json") {
		outputPath += ".json"
	}

	report := map[string]interface{}{
		"title":       title,
		"content":     content,
		"metadata":    metadata,
		"generated_at": time.Now().Format(time.RFC3339),
		"format":      "json",
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return outputPath, os.WriteFile(outputPath, data, 0644)
}

func (p *ReportingPlugin) generateYAML(title, content, outputPath string, metadata map[string]interface{}) (string, error) {
	if !strings.HasSuffix(outputPath, ".yaml") && !strings.HasSuffix(outputPath, ".yml") {
		outputPath += ".yaml"
	}

	report := map[string]interface{}{
		"title":       title,
		"content":     content,
		"metadata":    metadata,
		"generated_at": time.Now().Format(time.RFC3339),
		"format":      "yaml",
	}

	data, err := yaml.Marshal(report)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return outputPath, os.WriteFile(outputPath, data, 0644)
}

func (p *ReportingPlugin) convertToFormat(content, format, outputPath string) (string, error) {
	// Basic conversion logic
	switch format {
	case "json":
		return p.generateJSON("Converted Report", content, outputPath, nil)
	case "yaml":
		return p.generateYAML("Converted Report", content, outputPath, nil)
	case "html":
		return p.generateHTML("Converted Report", content, outputPath, nil)
	case "markdown":
		return p.generateMarkdown("Converted Report", content, outputPath, nil)
	case "pdf":
		return p.generatePDF("Converted Report", content, outputPath, nil)
	default:
		return "", fmt.Errorf("unsupported target format: %s", format)
	}
}

var ExportedPlugin plugin.Plugin = &ReportingPlugin{}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		if err := pluginv2.ServePlugin(&ReportingPlugin{}); err != nil {
			fmt.Printf("Error serving plugin: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Corynth Reporting Plugin v1.0.0\n")
		fmt.Printf("Usage: %s serve\n", os.Args[0])
	}
}
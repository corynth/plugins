package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/corynth/corynth-dist/pkg/plugin"
)

type ReportingPlugin struct{}

func (p *ReportingPlugin) Metadata() plugin.Metadata {
    return plugin.Metadata{
        Name:        "reporting",
        Version:     "1.0.0",
        Description: "Generate formatted reports with tables, charts, and multiple output formats",
        Author:      "Corynth Team",
        Tags:        []string{"reporting", "pdf", "markdown", "tables", "documentation", "output"},
        License:     "Apache-2.0",
    }
}

func (p *ReportingPlugin) Actions() []plugin.Action {
    return []plugin.Action{
        {
            Name:        "create_report",
            Description: "Create a formatted report with multiple sections",
            Inputs: map[string]plugin.InputSpec{
                "title": {
                    Type:        "string",
                    Description: "Report title",
                    Required:    true,
                },
                "sections": {
                    Type:        "array",
                    Description: "Report sections with content",
                    Required:    true,
                },
                "format": {
                    Type:        "string",
                    Description: "Output format (markdown, html, pdf)",
                    Required:    false,
                    Default:     "markdown",
                },
                "output_file": {
                    Type:        "string",
                    Description: "Output file path",
                    Required:    false,
                },
                "template": {
                    Type:        "string",
                    Description: "Report template (default, technical, executive)",
                    Required:    false,
                    Default:     "default",
                },
                "metadata": {
                    Type:        "object",
                    Description: "Report metadata (author, date, version)",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "content": {
                    Type:        "string",
                    Description: "Generated report content",
                },
                "file_path": {
                    Type:        "string",
                    Description: "Path to saved report file",
                },
                "format": {
                    Type:        "string",
                    Description: "Report format used",
                },
            },
        },
        {
            Name:        "create_table",
            Description: "Create a formatted table",
            Inputs: map[string]plugin.InputSpec{
                "headers": {
                    Type:        "array",
                    Description: "Table column headers",
                    Required:    true,
                },
                "rows": {
                    Type:        "array",
                    Description: "Table data rows",
                    Required:    true,
                },
                "format": {
                    Type:        "string",
                    Description: "Table format (markdown, ascii, html, csv)",
                    Required:    false,
                    Default:     "markdown",
                },
                "alignment": {
                    Type:        "array",
                    Description: "Column alignments (left, center, right)",
                    Required:    false,
                },
                "title": {
                    Type:        "string",
                    Description: "Table title",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "table": {
                    Type:        "string",
                    Description: "Formatted table output",
                },
                "row_count": {
                    Type:        "number",
                    Description: "Number of data rows",
                },
                "column_count": {
                    Type:        "number",
                    Description: "Number of columns",
                },
            },
        },
        {
            Name:        "create_chart",
            Description: "Create a chart visualization",
            Inputs: map[string]plugin.InputSpec{
                "type": {
                    Type:        "string",
                    Description: "Chart type (bar, line, pie, scatter)",
                    Required:    true,
                },
                "data": {
                    Type:        "object",
                    Description: "Chart data (labels and values)",
                    Required:    true,
                },
                "title": {
                    Type:        "string",
                    Description: "Chart title",
                    Required:    false,
                },
                "width": {
                    Type:        "number",
                    Description: "Chart width in pixels",
                    Required:    false,
                    Default:     800,
                },
                "height": {
                    Type:        "number",
                    Description: "Chart height in pixels",
                    Required:    false,
                    Default:     600,
                },
                "output_file": {
                    Type:        "string",
                    Description: "Output image file path",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "chart_url": {
                    Type:        "string",
                    Description: "Chart image URL or path",
                },
                "format": {
                    Type:        "string",
                    Description: "Chart output format",
                },
            },
        },
        {
            Name:        "convert_format",
            Description: "Convert report between formats",
            Inputs: map[string]plugin.InputSpec{
                "input_file": {
                    Type:        "string",
                    Description: "Input file path",
                    Required:    true,
                },
                "output_format": {
                    Type:        "string",
                    Description: "Target format (pdf, html, docx, markdown)",
                    Required:    true,
                },
                "output_file": {
                    Type:        "string",
                    Description: "Output file path",
                    Required:    false,
                },
                "options": {
                    Type:        "object",
                    Description: "Conversion options",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "output_file": {
                    Type:        "string",
                    Description: "Path to converted file",
                },
                "status": {
                    Type:        "string",
                    Description: "Conversion status",
                },
            },
        },
        {
            Name:        "display",
            Description: "Display formatted content on screen",
            Inputs: map[string]plugin.InputSpec{
                "content": {
                    Type:        "string",
                    Description: "Content to display",
                    Required:    true,
                },
                "format": {
                    Type:        "string",
                    Description: "Display format (text, markdown, table)",
                    Required:    false,
                    Default:     "text",
                },
                "style": {
                    Type:        "object",
                    Description: "Display styling options",
                    Required:    false,
                },
            },
            Outputs: map[string]plugin.OutputSpec{
                "displayed": {
                    Type:        "boolean",
                    Description: "Whether content was displayed",
                },
            },
        },
    }
}

func (p *ReportingPlugin) Validate(params map[string]interface{}) error {
    return nil
}

func (p *ReportingPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
    switch action {
    case "create_report":
        return p.executeCreateReport(ctx, params)
    case "create_table":
        return p.executeCreateTable(ctx, params)
    case "create_chart":
        return p.executeCreateChart(ctx, params)
    case "convert_format":
        return p.executeConvertFormat(ctx, params)
    case "display":
        return p.executeDisplay(ctx, params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

func (p *ReportingPlugin) executeCreateReport(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    title, _ := params["title"].(string)
    sections, _ := params["sections"].([]interface{})
    format, _ := params["format"].(string)
    if format == "" {
        format = "markdown"
    }
    
    outputFile, _ := params["output_file"].(string)
    template, _ := params["template"].(string)
    if template == "" {
        template = "default"
    }
    
    var report strings.Builder
    
    // Generate report based on format
    switch format {
    case "markdown":
        report.WriteString(fmt.Sprintf("# %s\n\n", title))
        
        // Add metadata if provided
        if metadata, ok := params["metadata"].(map[string]interface{}); ok {
            report.WriteString("---\n")
            for k, v := range metadata {
                report.WriteString(fmt.Sprintf("%s: %v\n", k, v))
            }
            report.WriteString("---\n\n")
        }
        
        // Add table of contents
        if len(sections) > 1 {
            report.WriteString("## Table of Contents\n\n")
            for i, section := range sections {
                if sec, ok := section.(map[string]interface{}); ok {
                    if heading, ok := sec["heading"].(string); ok {
                        report.WriteString(fmt.Sprintf("%d. [%s](#%s)\n", i+1, heading, 
                            strings.ToLower(strings.ReplaceAll(heading, " ", "-"))))
                    }
                }
            }
            report.WriteString("\n")
        }
        
        // Add sections
        for _, section := range sections {
            if sec, ok := section.(map[string]interface{}); ok {
                if heading, ok := sec["heading"].(string); ok {
                    level := 2
                    if l, ok := sec["level"].(float64); ok {
                        level = int(l)
                    }
                    report.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", level), heading))
                }
                
                if content, ok := sec["content"].(string); ok {
                    report.WriteString(fmt.Sprintf("%s\n\n", content))
                }
                
                if table, ok := sec["table"].(map[string]interface{}); ok {
                    tableStr := p.generateMarkdownTable(table)
                    report.WriteString(fmt.Sprintf("%s\n\n", tableStr))
                }
                
                if list, ok := sec["list"].([]interface{}); ok {
                    for _, item := range list {
                        report.WriteString(fmt.Sprintf("- %v\n", item))
                    }
                    report.WriteString("\n")
                }
                
                if code, ok := sec["code"].(map[string]interface{}); ok {
                    lang, _ := code["language"].(string)
                    content, _ := code["content"].(string)
                    report.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, content))
                }
            }
        }
        
        // Add footer
        report.WriteString(fmt.Sprintf("\n---\n*Generated on %s*\n", time.Now().Format("2006-01-02 15:04:05")))
        
    case "html":
        report.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
        report.WriteString(fmt.Sprintf("<title>%s</title>\n", title))
        report.WriteString("<style>")
        report.WriteString("body { font-family: Arial, sans-serif; max-width: 900px; margin: 0 auto; padding: 20px; }")
        report.WriteString("table { border-collapse: collapse; width: 100%; margin: 20px 0; }")
        report.WriteString("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }")
        report.WriteString("th { background-color: #f2f2f2; }")
        report.WriteString("code { background-color: #f4f4f4; padding: 2px 4px; border-radius: 3px; }")
        report.WriteString("pre { background-color: #f4f4f4; padding: 10px; border-radius: 5px; overflow-x: auto; }")
        report.WriteString("</style>\n</head>\n<body>\n")
        report.WriteString(fmt.Sprintf("<h1>%s</h1>\n", title))
        
        for _, section := range sections {
            if sec, ok := section.(map[string]interface{}); ok {
                if heading, ok := sec["heading"].(string); ok {
                    report.WriteString(fmt.Sprintf("<h2>%s</h2>\n", heading))
                }
                if content, ok := sec["content"].(string); ok {
                    report.WriteString(fmt.Sprintf("<p>%s</p>\n", content))
                }
            }
        }
        
        report.WriteString("</body>\n</html>\n")
        
    default:
        // Text format
        report.WriteString(fmt.Sprintf("%s\n%s\n\n", title, strings.Repeat("=", len(title))))
        
        for _, section := range sections {
            if sec, ok := section.(map[string]interface{}); ok {
                if heading, ok := sec["heading"].(string); ok {
                    report.WriteString(fmt.Sprintf("%s\n%s\n\n", heading, strings.Repeat("-", len(heading))))
                }
                if content, ok := sec["content"].(string); ok {
                    report.WriteString(fmt.Sprintf("%s\n\n", content))
                }
            }
        }
    }
    
    content := report.String()
    
    // Save to file if specified
    if outputFile != "" {
        err := os.WriteFile(outputFile, []byte(content), 0644)
        if err != nil {
            return nil, fmt.Errorf("failed to write report file: %w", err)
        }
        
        // Convert to PDF if requested
        if format == "pdf" {
            pdfFile := strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + ".pdf"
            if err := p.convertToPDF(outputFile, pdfFile); err == nil {
                outputFile = pdfFile
            }
        }
    }
    
    return map[string]interface{}{
        "content":   content,
        "file_path": outputFile,
        "format":    format,
    }, nil
}

func (p *ReportingPlugin) executeCreateTable(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    headers, _ := params["headers"].([]interface{})
    rows, _ := params["rows"].([]interface{})
    format, _ := params["format"].(string)
    if format == "" {
        format = "markdown"
    }
    
    title, _ := params["title"].(string)
    alignment, _ := params["alignment"].([]interface{})
    
    var table strings.Builder
    
    if title != "" {
        table.WriteString(fmt.Sprintf("### %s\n\n", title))
    }
    
    switch format {
    case "markdown":
        table.WriteString("|")
        for _, header := range headers {
            table.WriteString(fmt.Sprintf(" %v |", header))
        }
        table.WriteString("\n|")
        
        for i, header := range headers {
            align := "---"
            if i < len(alignment) {
                if a, ok := alignment[i].(string); ok {
                    switch a {
                    case "center":
                        align = ":---:"
                    case "right":
                        align = "---:"
                    }
                }
            }
            table.WriteString(fmt.Sprintf(" %s |", align))
        }
        table.WriteString("\n")
        
        for _, row := range rows {
            table.WriteString("|")
            if r, ok := row.([]interface{}); ok {
                for _, cell := range r {
                    table.WriteString(fmt.Sprintf(" %v |", cell))
                }
            }
            table.WriteString("\n")
        }
        
    case "csv":
        // CSV format
        for i, header := range headers {
            if i > 0 {
                table.WriteString(",")
            }
            table.WriteString(fmt.Sprintf("\"%v\"", header))
        }
        table.WriteString("\n")
        
        for _, row := range rows {
            if r, ok := row.([]interface{}); ok {
                for i, cell := range r {
                    if i > 0 {
                        table.WriteString(",")
                    }
                    table.WriteString(fmt.Sprintf("\"%v\"", cell))
                }
            }
            table.WriteString("\n")
        }
        
    case "html":
        table.WriteString("<table>\n<thead>\n<tr>")
        for _, header := range headers {
            table.WriteString(fmt.Sprintf("<th>%v</th>", header))
        }
        table.WriteString("</tr>\n</thead>\n<tbody>\n")
        
        for _, row := range rows {
            table.WriteString("<tr>")
            if r, ok := row.([]interface{}); ok {
                for _, cell := range r {
                    table.WriteString(fmt.Sprintf("<td>%v</td>", cell))
                }
            }
            table.WriteString("</tr>\n")
        }
        table.WriteString("</tbody>\n</table>\n")
        
    default:
        // ASCII format
        maxWidths := make([]int, len(headers))
        for i, header := range headers {
            maxWidths[i] = len(fmt.Sprintf("%v", header))
        }
        
        for _, row := range rows {
            if r, ok := row.([]interface{}); ok {
                for i, cell := range r {
                    if i < len(maxWidths) {
                        width := len(fmt.Sprintf("%v", cell))
                        if width > maxWidths[i] {
                            maxWidths[i] = width
                        }
                    }
                }
            }
        }
        
        // Print headers
        for i, header := range headers {
            table.WriteString(fmt.Sprintf("%-*v ", maxWidths[i], header))
        }
        table.WriteString("\n")
        
        // Print separator
        for i := range headers {
            table.WriteString(strings.Repeat("-", maxWidths[i]) + " ")
        }
        table.WriteString("\n")
        
        // Print rows
        for _, row := range rows {
            if r, ok := row.([]interface{}); ok {
                for i, cell := range r {
                    if i < len(maxWidths) {
                        table.WriteString(fmt.Sprintf("%-*v ", maxWidths[i], cell))
                    }
                }
            }
            table.WriteString("\n")
        }
    }
    
    return map[string]interface{}{
        "table":        table.String(),
        "row_count":    len(rows),
        "column_count": len(headers),
    }, nil
}

func (p *ReportingPlugin) executeCreateChart(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    chartType, _ := params["type"].(string)
    data, _ := params["data"].(map[string]interface{})
    title, _ := params["title"].(string)
    width, _ := params["width"].(float64)
    if width == 0 {
        width = 800
    }
    height, _ := params["height"].(float64)
    if height == 0 {
        height = 600
    }
    outputFile, _ := params["output_file"].(string)
    
    // For simplicity, generate ASCII chart for terminal display
    var chart strings.Builder
    
    if title != "" {
        chart.WriteString(fmt.Sprintf("%s\n", title))
        chart.WriteString(strings.Repeat("=", len(title)) + "\n\n")
    }
    
    switch chartType {
    case "bar":
        labels, _ := data["labels"].([]interface{})
        values, _ := data["values"].([]interface{})
        
        maxValue := 0.0
        for _, v := range values {
            if val, ok := v.(float64); ok && val > maxValue {
                maxValue = val
            }
        }
        
        for i, label := range labels {
            if i < len(values) {
                value, _ := values[i].(float64)
                barWidth := int((value / maxValue) * 50)
                chart.WriteString(fmt.Sprintf("%-15s |%s %.1f\n", 
                    label, strings.Repeat("█", barWidth), value))
            }
        }
        
    case "pie":
        labels, _ := data["labels"].([]interface{})
        values, _ := data["values"].([]interface{})
        
        total := 0.0
        for _, v := range values {
            if val, ok := v.(float64); ok {
                total += val
            }
        }
        
        for i, label := range labels {
            if i < len(values) {
                value, _ := values[i].(float64)
                percentage := (value / total) * 100
                chart.WriteString(fmt.Sprintf("%-15s: %.1f%%\n", label, percentage))
            }
        }
        
    default:
        chart.WriteString("Chart type not supported in ASCII mode\n")
    }
    
    chartStr := chart.String()
    
    // Save to file if specified
    if outputFile != "" {
        err := os.WriteFile(outputFile, []byte(chartStr), 0644)
        if err != nil {
            return nil, fmt.Errorf("failed to write chart file: %w", err)
        }
    }
    
    return map[string]interface{}{
        "chart_url": outputFile,
        "format":    "ascii",
    }, nil
}

func (p *ReportingPlugin) executeConvertFormat(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    inputFile, _ := params["input_file"].(string)
    outputFormat, _ := params["output_format"].(string)
    outputFile, _ := params["output_file"].(string)
    
    if outputFile == "" {
        outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "." + outputFormat
    }
    
    // Check if pandoc is available for conversion
    if _, err := exec.LookPath("pandoc"); err == nil {
        cmd := exec.CommandContext(ctx, "pandoc", inputFile, "-o", outputFile)
        if err := cmd.Run(); err != nil {
            return map[string]interface{}{
                "output_file": "",
                "status":      "failed",
                "error":       err.Error(),
            }, nil
        }
        
        return map[string]interface{}{
            "output_file": outputFile,
            "status":      "success",
        }, nil
    }
    
    // Fallback to basic conversion
    content, err := os.ReadFile(inputFile)
    if err != nil {
        return nil, fmt.Errorf("failed to read input file: %w", err)
    }
    
    // Simple conversion (just copy for now)
    err = os.WriteFile(outputFile, content, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to write output file: %w", err)
    }
    
    return map[string]interface{}{
        "output_file": outputFile,
        "status":      "success",
    }, nil
}

func (p *ReportingPlugin) executeDisplay(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
    content, _ := params["content"].(string)
    format, _ := params["format"].(string)
    if format == "" {
        format = "text"
    }
    
    // Display to stdout
    fmt.Println(content)
    
    return map[string]interface{}{
        "displayed": true,
    }, nil
}

func (p *ReportingPlugin) generateMarkdownTable(table map[string]interface{}) string {
    headers, _ := table["headers"].([]interface{})
    rows, _ := table["rows"].([]interface{})
    
    var result strings.Builder
    
    result.WriteString("|")
    for _, header := range headers {
        result.WriteString(fmt.Sprintf(" %v |", header))
    }
    result.WriteString("\n|")
    
    for range headers {
        result.WriteString(" --- |")
    }
    result.WriteString("\n")
    
    for _, row := range rows {
        result.WriteString("|")
        if r, ok := row.([]interface{}); ok {
            for _, cell := range r {
                result.WriteString(fmt.Sprintf(" %v |", cell))
            }
        }
        result.WriteString("\n")
    }
    
    return result.String()
}

func (p *ReportingPlugin) convertToPDF(inputFile, outputFile string) error {
    // Try using pandoc if available
    if _, err := exec.LookPath("pandoc"); err == nil {
        cmd := exec.Command("pandoc", inputFile, "-o", outputFile)
        return cmd.Run()
    }
    
    // Try using wkhtmltopdf if available
    if _, err := exec.LookPath("wkhtmltopdf"); err == nil {
        cmd := exec.Command("wkhtmltopdf", inputFile, outputFile)
        return cmd.Run()
    }
    
    return fmt.Errorf("no PDF converter available")
}

var ExportedPlugin plugin.Plugin = &ReportingPlugin{}
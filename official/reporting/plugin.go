package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	htmltemplate "html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	texttemplate "text/template"
	"time"
)

type Metadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
}

type IOSpec struct {
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description"`
}

type ActionSpec struct {
	Description string            `json:"description"`
	Inputs      map[string]IOSpec `json:"inputs"`
	Outputs     map[string]IOSpec `json:"outputs"`
}

type ReportingPlugin struct{}

func NewReportingPlugin() *ReportingPlugin {
	return &ReportingPlugin{}
}

func (p *ReportingPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "reporting",
		Version:     "1.0.0",
		Description: "Generate formatted reports with tables, charts, and multiple output formats",
		Author:      "Corynth Team",
		Tags:        []string{"reporting", "pdf", "markdown", "tables", "documentation", "output"},
	}
}

func (p *ReportingPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"create_report": {
			Description: "Create formatted report",
			Inputs: map[string]IOSpec{
				"title":       {Type: "string", Required: true, Description: "Report title"},
				"content":     {Type: "string", Required: true, Description: "Report content"},
				"format":      {Type: "string", Required: false, Default: "markdown", Description: "Output format: markdown, html, text"},
				"output_path": {Type: "string", Required: false, Description: "Output file path"},
				"metadata":    {Type: "object", Required: false, Description: "Report metadata"},
			},
			Outputs: map[string]IOSpec{
				"report":    {Type: "string", Description: "Generated report"},
				"file_path": {Type: "string", Description: "Output file path"},
			},
		},
		"create_table": {
			Description: "Create formatted table",
			Inputs: map[string]IOSpec{
				"data":    {Type: "array", Required: true, Description: "Table data"},
				"headers": {Type: "array", Required: false, Description: "Column headers"},
				"format":  {Type: "string", Required: false, Default: "markdown", Description: "Table format"},
				"title":   {Type: "string", Required: false, Description: "Table title"},
			},
			Outputs: map[string]IOSpec{
				"table": {Type: "string", Description: "Formatted table"},
			},
		},
		"create_chart": {
			Description: "Create ASCII chart",
			Inputs: map[string]IOSpec{
				"data":  {Type: "object", Required: true, Description: "Chart data"},
				"type":  {Type: "string", Required: false, Default: "bar", Description: "Chart type: bar, line"},
				"title": {Type: "string", Required: false, Description: "Chart title"},
				"width": {Type: "number", Required: false, Default: 60, Description: "Chart width"},
			},
			Outputs: map[string]IOSpec{
				"chart": {Type: "string", Description: "ASCII chart"},
			},
		},
	}
}

func (p *ReportingPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "create_report":
		return p.createReport(params)
	case "create_table":
		return p.createTable(params)
	case "create_chart":
		return p.createChart(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *ReportingPlugin) createReport(params map[string]interface{}) (map[string]interface{}, error) {
	title, ok := params["title"].(string)
	if !ok || title == "" {
		return map[string]interface{}{"error": "title is required"}, nil
	}

	content, ok := params["content"].(string)
	if !ok || content == "" {
		return map[string]interface{}{"error": "content is required"}, nil
	}

	format := getStringParam(params, "format", "markdown")
	outputPath := getStringParam(params, "output_path", "")
	metadata := getMapParam(params, "metadata", make(map[string]interface{}))

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var report string
	var err error

	switch format {
	case "markdown":
		report, err = p.generateMarkdownReport(title, content, metadata, timestamp)
	case "html":
		report, err = p.generateHTMLReport(title, content, metadata, timestamp)
	default: // text
		report, err = p.generateTextReport(title, content, metadata, timestamp)
	}

	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	// Write to file if path specified
	if outputPath != "" {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to create directory: %v", err)}, nil
		}

		if err := os.WriteFile(outputPath, []byte(report), 0644); err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to write file: %v", err)}, nil
		}
	}

	return map[string]interface{}{
		"report":    report,
		"file_path": outputPath,
	}, nil
}

func (p *ReportingPlugin) createTable(params map[string]interface{}) (map[string]interface{}, error) {
	dataRaw, ok := params["data"]
	if !ok {
		return map[string]interface{}{"error": "data is required"}, nil
	}

	data, ok := dataRaw.([]interface{})
	if !ok {
		return map[string]interface{}{"error": "data must be an array"}, nil
	}

	if len(data) == 0 {
		return map[string]interface{}{"table": "No data provided"}, nil
	}

	var headers []string
	if headersRaw, ok := params["headers"].([]interface{}); ok {
		headers = make([]string, len(headersRaw))
		for i, h := range headersRaw {
			if s, ok := h.(string); ok {
				headers[i] = s
			} else {
				headers[i] = fmt.Sprintf("%v", h)
			}
		}
	} else {
		// Auto-detect headers
		if len(data) > 0 {
			if rowMap, ok := data[0].(map[string]interface{}); ok {
				headers = make([]string, 0, len(rowMap))
				for key := range rowMap {
					headers = append(headers, key)
				}
				sort.Strings(headers)
			} else if rowSlice, ok := data[0].([]interface{}); ok {
				headers = make([]string, len(rowSlice))
				for i := range headers {
					headers[i] = fmt.Sprintf("Column %d", i+1)
				}
			} else {
				headers = []string{"Value"}
			}
		}
	}

	format := getStringParam(params, "format", "markdown")
	title := getStringParam(params, "title", "")

	var table string
	if format == "markdown" {
		table = p.generateMarkdownTable(data, headers, title)
	} else {
		table = p.generateTextTable(data, headers, title)
	}

	return map[string]interface{}{
		"table": table,
	}, nil
}

func (p *ReportingPlugin) createChart(params map[string]interface{}) (map[string]interface{}, error) {
	dataRaw, ok := params["data"]
	if !ok {
		return map[string]interface{}{"error": "data is required"}, nil
	}

	data, ok := dataRaw.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"error": "data must be an object"}, nil
	}

	if len(data) == 0 {
		return map[string]interface{}{"chart": "No data provided"}, nil
	}

	chartType := getStringParam(params, "type", "bar")
	title := getStringParam(params, "title", "")
	width := int(getFloatParam(params, "width", 60))

	var chart string
	switch chartType {
	case "bar", "line": // Both use bar chart for simplicity
		chart = p.generateBarChart(data, title, width)
	default:
		chart = p.generateBarChart(data, title, width)
	}

	return map[string]interface{}{
		"chart": chart,
	}, nil
}

func (p *ReportingPlugin) generateMarkdownReport(title, content string, metadata map[string]interface{}, timestamp string) (string, error) {
	tmplStr := `# {{.Title}}

**Generated:** {{.Timestamp}}

{{if .Metadata}}## Metadata
{{range $key, $value := .Metadata}}- **{{$key}}:** {{$value}}
{{end}}
{{end}}## Report Content
{{.Content}}`

	tmpl, err := texttemplate.New("markdown").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"Title":     title,
		"Content":   content,
		"Metadata":  metadata,
		"Timestamp": timestamp,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *ReportingPlugin) generateHTMLReport(title, content string, metadata map[string]interface{}, timestamp string) (string, error) {
	tmplStr := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        .metadata { background: #f5f5f5; padding: 15px; margin: 20px 0; }
        .timestamp { color: #666; font-style: italic; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="timestamp">Generated: {{.Timestamp}}</div>
    {{if .Metadata}}<div class="metadata">
        <h3>Metadata</h3>
        <ul>
        {{range $key, $value := .Metadata}}<li><strong>{{$key}}:</strong> {{$value}}</li>{{end}}
        </ul>
    </div>{{end}}
    <div class="content">{{.Content}}</div>
</body>
</html>`

	tmpl, err := htmltemplate.New("html").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"Title":     title,
		"Content":   content,
		"Metadata":  metadata,
		"Timestamp": timestamp,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *ReportingPlugin) generateTextReport(title, content string, metadata map[string]interface{}, timestamp string) (string, error) {
	var lines []string

	// Title with underline
	lines = append(lines, strings.Repeat("=", len(title)))
	lines = append(lines, title)
	lines = append(lines, strings.Repeat("=", len(title)))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Generated: %s", timestamp))
	lines = append(lines, "")

	// Metadata
	if len(metadata) > 0 {
		lines = append(lines, "METADATA:")
		lines = append(lines, strings.Repeat("-", 20))
		for key, value := range metadata {
			lines = append(lines, fmt.Sprintf("%s: %v", key, value))
		}
		lines = append(lines, "")
	}

	// Content
	lines = append(lines, "CONTENT:")
	lines = append(lines, strings.Repeat("-", 20))
	lines = append(lines, content)

	return strings.Join(lines, "\n"), nil
}

func (p *ReportingPlugin) generateMarkdownTable(data []interface{}, headers []string, title string) string {
	var lines []string

	if title != "" {
		lines = append(lines, fmt.Sprintf("### %s", title))
		lines = append(lines, "")
	}

	// Headers
	lines = append(lines, "| "+strings.Join(headers, " | ")+" |")
	lines = append(lines, "| "+strings.Join(func() []string {
		result := make([]string, len(headers))
		for i, h := range headers {
			result[i] = strings.Repeat("-", len(h))
		}
		return result
	}(), " | ")+" |")

	// Data rows
	for _, rowRaw := range data {
		var rowData []string
		if rowMap, ok := rowRaw.(map[string]interface{}); ok {
			rowData = make([]string, len(headers))
			for i, header := range headers {
				if val, exists := rowMap[header]; exists {
					rowData[i] = fmt.Sprintf("%v", val)
				} else {
					rowData[i] = ""
				}
			}
		} else if rowSlice, ok := rowRaw.([]interface{}); ok {
			rowData = make([]string, len(headers))
			for i, val := range rowSlice {
				if i < len(rowData) {
					rowData[i] = fmt.Sprintf("%v", val)
				}
			}
		} else {
			rowData = []string{fmt.Sprintf("%v", rowRaw)}
		}

		lines = append(lines, "| "+strings.Join(rowData, " | ")+" |")
	}

	return strings.Join(lines, "\n")
}

func (p *ReportingPlugin) generateTextTable(data []interface{}, headers []string, title string) string {
	var lines []string

	if title != "" {
		lines = append(lines, title)
		lines = append(lines, strings.Repeat("=", len(title)))
		lines = append(lines, "")
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}

	// Convert data to string matrix and update widths
	var rows [][]string
	for _, rowRaw := range data {
		var rowData []string
		if rowMap, ok := rowRaw.(map[string]interface{}); ok {
			rowData = make([]string, len(headers))
			for i, header := range headers {
				if val, exists := rowMap[header]; exists {
					rowData[i] = fmt.Sprintf("%v", val)
				} else {
					rowData[i] = ""
				}
			}
		} else if rowSlice, ok := rowRaw.([]interface{}); ok {
			rowData = make([]string, len(headers))
			for i, val := range rowSlice {
				if i < len(rowData) {
					rowData[i] = fmt.Sprintf("%v", val)
				}
			}
		} else {
			rowData = []string{fmt.Sprintf("%v", rowRaw)}
		}

		// Update column widths
		for i, cell := range rowData {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
		rows = append(rows, rowData)
	}

	// Header
	headerParts := make([]string, len(headers))
	for i, h := range headers {
		headerParts[i] = fmt.Sprintf("%-*s", colWidths[i], h)
	}
	headerLine := strings.Join(headerParts, " | ")
	lines = append(lines, headerLine)
	lines = append(lines, strings.Repeat("-", len(headerLine)))

	// Data
	for _, row := range rows {
		rowParts := make([]string, len(headers))
		for i, cell := range row {
			if i < len(rowParts) {
				rowParts[i] = fmt.Sprintf("%-*s", colWidths[i], cell)
			}
		}
		lines = append(lines, strings.Join(rowParts, " | "))
	}

	return strings.Join(lines, "\n")
}

func (p *ReportingPlugin) generateBarChart(data map[string]interface{}, title string, width int) string {
	var lines []string

	if title != "" {
		lines = append(lines, title)
		lines = append(lines, strings.Repeat("=", len(title)))
		lines = append(lines, "")
	}

	if len(data) == 0 {
		return "No data provided"
	}

	// Find max value for scaling
	maxVal := 0.0
	for _, v := range data {
		if val, err := convertToFloat(v); err == nil && val > maxVal {
			maxVal = val
		}
	}

	if maxVal == 0 {
		maxVal = 1
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, label := range keys {
		value := data[label]
		val, err := convertToFloat(value)
		if err != nil {
			continue
		}

		barLength := int((val / maxVal) * float64(width))
		bar := strings.Repeat("█", barLength)
		lines = append(lines, fmt.Sprintf("%15s | %s %.2f", label, bar, val))
	}

	return strings.Join(lines, "\n")
}

// Helper functions
func getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key].(string); ok {
		return val
	}
	return defaultValue
}

func getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key].(float64); ok {
		return val
	}
	return defaultValue
}

func getMapParam(params map[string]interface{}, key string, defaultValue map[string]interface{}) map[string]interface{} {
	if val, ok := params[key].(map[string]interface{}); ok {
		return val
	}
	return defaultValue
}

func convertToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewReportingPlugin()

	var result interface{}

	switch action {
	case "metadata":
		result = plugin.GetMetadata()
	case "actions":
		result = plugin.GetActions()
	default:
		var params map[string]interface{}
		inputData, err := io.ReadAll(os.Stdin)
		if err != nil {
			result = map[string]interface{}{"error": fmt.Sprintf("failed to read input: %v", err)}
		} else if len(inputData) > 0 {
			if err := json.Unmarshal(inputData, &params); err != nil {
				result = map[string]interface{}{"error": fmt.Sprintf("failed to parse JSON: %v", err)}
			} else {
				result, err = plugin.Execute(action, params)
				if err != nil {
					result = map[string]interface{}{"error": err.Error()}
				}
			}
		} else {
			result, err = plugin.Execute(action, map[string]interface{}{})
			if err != nil {
				result = map[string]interface{}{"error": err.Error()}
			}
		}
	}

	json.NewEncoder(os.Stdout).Encode(result)
}
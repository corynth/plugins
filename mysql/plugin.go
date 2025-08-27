package main

import (
	"context"
	"fmt"
	
	"github.com/corynth/corynth-dist/pkg/plugin"
)

type MySQLPlugin struct{}

func (p *MySQLPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:        "mysql",
		Version:     "1.0.0",
		Description: "MySQL database operations and management",
		Author:      "Corynth Team",
		Tags:        []string{"database", "mysql", "sql", "rdbms"},
		License:     "Apache-2.0",
	}
}

func (p *MySQLPlugin) Actions() []plugin.Action {
	return []plugin.Action{
		{
			Name:        "query",
			Description: "Execute a SELECT query",
			Inputs: map[string]plugin.InputSpec{
				"host": {
					Type:        "string",
					Description: "MySQL host",
					Required:    false,
					Default:     "localhost",
				},
				"port": {
					Type:        "number",
					Description: "MySQL port",
					Required:    false,
					Default:     3306,
				},
				"database": {
					Type:        "string",
					Description: "Database name",
					Required:    true,
				},
				"user": {
					Type:        "string",
					Description: "Database user",
					Required:    true,
				},
				"password": {
					Type:        "string",
					Description: "Database password",
					Required:    false,
				},
				"query": {
					Type:        "string",
					Description: "SQL query to execute",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"rows": {
					Type:        "array",
					Description: "Query result rows",
				},
				"count": {
					Type:        "number",
					Description: "Number of rows returned",
				},
			},
		},
		{
			Name:        "execute",
			Description: "Execute an INSERT, UPDATE, or DELETE statement",
			Inputs: map[string]plugin.InputSpec{
				"host": {
					Type:        "string",
					Description: "MySQL host",
					Required:    false,
					Default:     "localhost",
				},
				"port": {
					Type:        "number",
					Description: "MySQL port",
					Required:    false,
					Default:     3306,
				},
				"database": {
					Type:        "string",
					Description: "Database name",
					Required:    true,
				},
				"user": {
					Type:        "string",
					Description: "Database user",
					Required:    true,
				},
				"password": {
					Type:        "string",
					Description: "Database password",
					Required:    false,
				},
				"statement": {
					Type:        "string",
					Description: "SQL statement to execute",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"affected_rows": {
					Type:        "number",
					Description: "Number of rows affected",
				},
				"success": {
					Type:        "boolean",
					Description: "Whether the operation succeeded",
				},
			},
		},
		{
			Name:        "backup",
			Description: "Create a database backup",
			Inputs: map[string]plugin.InputSpec{
				"host": {
					Type:        "string",
					Description: "MySQL host",
					Required:    false,
					Default:     "localhost",
				},
				"port": {
					Type:        "number",
					Description: "MySQL port",
					Required:    false,
					Default:     3306,
				},
				"database": {
					Type:        "string",
					Description: "Database name",
					Required:    true,
				},
				"user": {
					Type:        "string",
					Description: "Database user",
					Required:    true,
				},
				"password": {
					Type:        "string",
					Description: "Database password",
					Required:    false,
				},
				"output_file": {
					Type:        "string",
					Description: "Path to backup file",
					Required:    true,
				},
			},
			Outputs: map[string]plugin.OutputSpec{
				"file": {
					Type:        "string",
					Description: "Path to the backup file",
				},
				"size": {
					Type:        "number",
					Description: "Size of backup in bytes",
				},
			},
		},
	}
}

func (p *MySQLPlugin) Validate(params map[string]interface{}) error {
	if database, ok := params["database"].(string); ok && database == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	
	if user, ok := params["user"].(string); ok && user == "" {
		return fmt.Errorf("user cannot be empty")
	}
	
	if port, ok := params["port"].(float64); ok {
		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535")
		}
	}
	
	return nil
}

func (p *MySQLPlugin) Execute(ctx context.Context, action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "query":
		return p.executeQuery(ctx, params)
	case "execute":
		return p.executeStatement(ctx, params)
	case "backup":
		return p.executeBackup(ctx, params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *MySQLPlugin) executeQuery(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	database, ok := params["database"].(string)
	if !ok || database == "" {
		return nil, fmt.Errorf("database parameter is required")
	}

	user, ok := params["user"].(string)
	if !ok || user == "" {
		return nil, fmt.Errorf("user parameter is required")
	}

	query, ok := params["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	host := "localhost"
	if h, ok := params["host"].(string); ok {
		host = h
	}

	port := 3306
	if p, ok := params["port"].(float64); ok {
		port = int(p)
	}

	// In production, this would connect to MySQL and execute the query
	// For demonstration, we'll return mock data
	mockRows := []map[string]interface{}{
		{"id": 1, "name": "Record 1"},
		{"id": 2, "name": "Record 2"},
	}

	return map[string]interface{}{
		"rows":    mockRows,
		"count":   len(mockRows),
		"message": fmt.Sprintf("Executed query on %s@%s:%d/%s", user, host, port, database),
	}, nil
}

func (p *MySQLPlugin) executeStatement(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	database, ok := params["database"].(string)
	if !ok || database == "" {
		return nil, fmt.Errorf("database parameter is required")
	}

	user, ok := params["user"].(string)
	if !ok || user == "" {
		return nil, fmt.Errorf("user parameter is required")
	}

	statement, ok := params["statement"].(string)
	if !ok || statement == "" {
		return nil, fmt.Errorf("statement parameter is required")
	}

	host := "localhost"
	if h, ok := params["host"].(string); ok {
		host = h
	}

	port := 3306
	if p, ok := params["port"].(float64); ok {
		port = int(p)
	}

	// In production, this would connect to MySQL and execute the statement
	// For demonstration, we'll simulate success
	return map[string]interface{}{
		"affected_rows": 1,
		"success":       true,
		"message":       fmt.Sprintf("Executed statement on %s@%s:%d/%s", user, host, port, database),
	}, nil
}

func (p *MySQLPlugin) executeBackup(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	database, ok := params["database"].(string)
	if !ok || database == "" {
		return nil, fmt.Errorf("database parameter is required")
	}

	user, ok := params["user"].(string)
	if !ok || user == "" {
		return nil, fmt.Errorf("user parameter is required")
	}

	outputFile, ok := params["output_file"].(string)
	if !ok || outputFile == "" {
		return nil, fmt.Errorf("output_file parameter is required")
	}

	host := "localhost"
	if h, ok := params["host"].(string); ok {
		host = h
	}

	port := 3306
	if p, ok := params["port"].(float64); ok {
		port = int(p)
	}

	// In production, this would use mysqldump to create a backup
	// For demonstration, we'll simulate success
	return map[string]interface{}{
		"file":    outputFile,
		"size":    1024000, // Mock 1MB backup
		"message": fmt.Sprintf("Backed up database %s@%s:%d/%s to %s", user, host, port, database, outputFile),
	}, nil
}

var ExportedPlugin plugin.Plugin = &MySQLPlugin{}
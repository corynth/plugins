package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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

type SQLPlugin struct{}

func NewSQLPlugin() *SQLPlugin {
	return &SQLPlugin{}
}

func (p *SQLPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "sql",
		Version:     "1.0.0",
		Description: "SQL database operations for SQLite, PostgreSQL, and MySQL",
		Author:      "Corynth Team",
		Tags:        []string{"sql", "database", "query", "sqlite", "postgresql", "mysql"},
	}
}

func (p *SQLPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"query": {
			Description: "Execute SELECT query and return results",
			Inputs: map[string]IOSpec{
				"connection_string": {
					Type:        "string",
					Required:    true,
					Description: "Database connection string (sqlite://path, postgres://user:pass@host/db, mysql://user:pass@host/db)",
				},
				"query": {
					Type:        "string",
					Required:    true,
					Description: "SQL SELECT query to execute",
				},
				"params": {
					Type:        "array",
					Required:    false,
					Description: "Query parameters for prepared statements",
				},
			},
			Outputs: map[string]IOSpec{
				"rows":      {Type: "array", Description: "Query result rows as array of objects"},
				"columns":   {Type: "array", Description: "Column names"},
				"row_count": {Type: "number", Description: "Number of rows returned"},
			},
		},
		"execute": {
			Description: "Execute INSERT/UPDATE/DELETE statement",
			Inputs: map[string]IOSpec{
				"connection_string": {
					Type:        "string",
					Required:    true,
					Description: "Database connection string",
				},
				"statement": {
					Type:        "string",
					Required:    true,
					Description: "SQL statement to execute (INSERT/UPDATE/DELETE)",
				},
				"params": {
					Type:        "array",
					Required:    false,
					Description: "Statement parameters for prepared statements",
				},
			},
			Outputs: map[string]IOSpec{
				"affected_rows": {Type: "number", Description: "Number of rows affected"},
				"last_insert_id": {Type: "number", Description: "Last inserted ID (if applicable)"},
				"success":       {Type: "boolean", Description: "Operation success status"},
			},
		},
		"schema": {
			Description: "Get database schema information",
			Inputs: map[string]IOSpec{
				"connection_string": {
					Type:        "string",
					Required:    true,
					Description: "Database connection string",
				},
				"table_name": {
					Type:        "string",
					Required:    false,
					Description: "Specific table name to get schema for",
				},
			},
			Outputs: map[string]IOSpec{
				"tables":  {Type: "array", Description: "List of table names"},
				"columns": {Type: "object", Description: "Column information by table name"},
			},
		},
	}
}

func (p *SQLPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "query":
		return p.executeQuery(params)
	case "execute":
		return p.executeStatement(params)
	case "schema":
		return p.getSchema(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *SQLPlugin) parseConnectionString(connStr string) (string, string, error) {
	u, err := url.Parse(connStr)
	if err != nil {
		return "", "", fmt.Errorf("invalid connection string: %v", err)
	}

	switch u.Scheme {
	case "sqlite":
		// sqlite://path/to/db.sqlite or sqlite:///absolute/path/to/db.sqlite
		path := u.Path
		if u.Host != "" {
			path = u.Host + path
		}
		return "sqlite3", path, nil

	case "postgres", "postgresql":
		// postgres://user:password@host:port/dbname?sslmode=disable
		return "postgres", connStr, nil

	case "mysql":
		// mysql://user:password@host:port/dbname
		// Convert to MySQL DSN format: user:password@tcp(host:port)/dbname
		userInfo := u.User
		if userInfo == nil {
			return "", "", fmt.Errorf("mysql connection requires user credentials")
		}
		
		username := userInfo.Username()
		password, _ := userInfo.Password()
		host := u.Host
		if host == "" {
			host = "localhost:3306"
		}
		dbname := strings.TrimPrefix(u.Path, "/")
		
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbname)
		
		// Add query parameters
		if u.RawQuery != "" {
			dsn += "?" + u.RawQuery
		}
		
		return "mysql", dsn, nil

	default:
		return "", "", fmt.Errorf("unsupported database type: %s", u.Scheme)
	}
}

func (p *SQLPlugin) executeQuery(params map[string]interface{}) (map[string]interface{}, error) {
	connStr, ok := params["connection_string"].(string)
	if !ok || connStr == "" {
		return map[string]interface{}{"error": "connection_string is required"}, nil
	}

	query, ok := params["query"].(string)
	if !ok || query == "" {
		return map[string]interface{}{"error": "query is required"}, nil
	}

	driverName, dataSource, err := p.parseConnectionString(connStr)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to connect: %v", err)}, nil
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to ping database: %v", err)}, nil
	}

	// Get parameters
	var queryParams []interface{}
	if paramsVal, ok := params["params"]; ok {
		if paramsList, ok := paramsVal.([]interface{}); ok {
			queryParams = paramsList
		}
	}

	rows, err := db.Query(query, queryParams...)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("query failed: %v", err)}, nil
	}
	defer rows.Close()

	// Get column information
	columns, err := rows.Columns()
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to get columns: %v", err)}, nil
	}

	// Prepare result storage
	var result []map[string]interface{}
	columnCount := len(columns)
	
	for rows.Next() {
		// Create a slice of interface{} to hold the column values
		values := make([]interface{}, columnCount)
		valuePtrs := make([]interface{}, columnCount)
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the result into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("scan failed: %v", err)}, nil
		}

		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			
			// Convert []byte to string for better JSON serialization
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			
			row[col] = val
		}
		
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("rows error: %v", err)}, nil
	}

	return map[string]interface{}{
		"rows":      result,
		"columns":   columns,
		"row_count": len(result),
	}, nil
}

func (p *SQLPlugin) executeStatement(params map[string]interface{}) (map[string]interface{}, error) {
	connStr, ok := params["connection_string"].(string)
	if !ok || connStr == "" {
		return map[string]interface{}{"error": "connection_string is required"}, nil
	}

	statement, ok := params["statement"].(string)
	if !ok || statement == "" {
		return map[string]interface{}{"error": "statement is required"}, nil
	}

	driverName, dataSource, err := p.parseConnectionString(connStr)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to connect: %v", err)}, nil
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to ping database: %v", err)}, nil
	}

	// Get parameters
	var stmtParams []interface{}
	if paramsVal, ok := params["params"]; ok {
		if paramsList, ok := paramsVal.([]interface{}); ok {
			stmtParams = paramsList
		}
	}

	result, err := db.Exec(statement, stmtParams...)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("execution failed: %v", err)}, nil
	}

	affectedRows, _ := result.RowsAffected()
	lastInsertID, _ := result.LastInsertId()

	return map[string]interface{}{
		"affected_rows":   affectedRows,
		"last_insert_id":  lastInsertID,
		"success":         true,
	}, nil
}

func (p *SQLPlugin) getSchema(params map[string]interface{}) (map[string]interface{}, error) {
	connStr, ok := params["connection_string"].(string)
	if !ok || connStr == "" {
		return map[string]interface{}{"error": "connection_string is required"}, nil
	}

	driverName, dataSource, err := p.parseConnectionString(connStr)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to connect: %v", err)}, nil
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to ping database: %v", err)}, nil
	}

	tableName, _ := params["table_name"].(string)

	switch driverName {
	case "sqlite3":
		return p.getSQLiteSchema(db, tableName)
	case "postgres":
		return p.getPostgreSQLSchema(db, tableName)
	case "mysql":
		return p.getMySQLSchema(db, tableName)
	default:
		return map[string]interface{}{"error": "unsupported database type for schema"}, nil
	}
}

func (p *SQLPlugin) getSQLiteSchema(db *sql.DB, tableName string) (map[string]interface{}, error) {
	if tableName != "" {
		// Get specific table schema
		query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
		rows, err := db.Query(query)
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to get table info: %v", err)}, nil
		}
		defer rows.Close()

		var columns []map[string]interface{}
		for rows.Next() {
			var cid int
			var name, dataType string
			var notNull, pk int
			var defaultValue sql.NullString

			if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
				return map[string]interface{}{"error": fmt.Sprintf("scan failed: %v", err)}, nil
			}

			column := map[string]interface{}{
				"name":         name,
				"type":         dataType,
				"not_null":     notNull == 1,
				"primary_key":  pk == 1,
				"default":      nil,
			}

			if defaultValue.Valid {
				column["default"] = defaultValue.String
			}

			columns = append(columns, column)
		}

		return map[string]interface{}{
			"tables":  []string{tableName},
			"columns": map[string]interface{}{tableName: columns},
		}, nil
	} else {
		// Get all tables
		rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to get tables: %v", err)}, nil
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return map[string]interface{}{"error": fmt.Sprintf("scan failed: %v", err)}, nil
			}
			tables = append(tables, name)
		}

		return map[string]interface{}{
			"tables":  tables,
			"columns": map[string]interface{}{},
		}, nil
	}
}

func (p *SQLPlugin) getPostgreSQLSchema(db *sql.DB, tableName string) (map[string]interface{}, error) {
	if tableName != "" {
		// Get specific table schema
		query := `
			SELECT column_name, data_type, is_nullable, column_default
			FROM information_schema.columns 
			WHERE table_name = $1
			ORDER BY ordinal_position`

		rows, err := db.Query(query, tableName)
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to get table info: %v", err)}, nil
		}
		defer rows.Close()

		var columns []map[string]interface{}
		for rows.Next() {
			var columnName, dataType, isNullable string
			var columnDefault sql.NullString

			if err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault); err != nil {
				return map[string]interface{}{"error": fmt.Sprintf("scan failed: %v", err)}, nil
			}

			column := map[string]interface{}{
				"name":     columnName,
				"type":     dataType,
				"not_null": isNullable == "NO",
				"default":  nil,
			}

			if columnDefault.Valid {
				column["default"] = columnDefault.String
			}

			columns = append(columns, column)
		}

		return map[string]interface{}{
			"tables":  []string{tableName},
			"columns": map[string]interface{}{tableName: columns},
		}, nil
	} else {
		// Get all tables
		query := `
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_type = 'BASE TABLE'`

		rows, err := db.Query(query)
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to get tables: %v", err)}, nil
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return map[string]interface{}{"error": fmt.Sprintf("scan failed: %v", err)}, nil
			}
			tables = append(tables, name)
		}

		return map[string]interface{}{
			"tables":  tables,
			"columns": map[string]interface{}{},
		}, nil
	}
}

func (p *SQLPlugin) getMySQLSchema(db *sql.DB, tableName string) (map[string]interface{}, error) {
	if tableName != "" {
		// Get specific table schema
		query := `
			SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_KEY
			FROM INFORMATION_SCHEMA.COLUMNS 
			WHERE TABLE_NAME = ?
			ORDER BY ORDINAL_POSITION`

		rows, err := db.Query(query, tableName)
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to get table info: %v", err)}, nil
		}
		defer rows.Close()

		var columns []map[string]interface{}
		for rows.Next() {
			var columnName, dataType, isNullable, columnKey string
			var columnDefault sql.NullString

			if err := rows.Scan(&columnName, &dataType, &isNullable, &columnDefault, &columnKey); err != nil {
				return map[string]interface{}{"error": fmt.Sprintf("scan failed: %v", err)}, nil
			}

			column := map[string]interface{}{
				"name":        columnName,
				"type":        dataType,
				"not_null":    isNullable == "NO",
				"primary_key": columnKey == "PRI",
				"default":     nil,
			}

			if columnDefault.Valid {
				column["default"] = columnDefault.String
			}

			columns = append(columns, column)
		}

		return map[string]interface{}{
			"tables":  []string{tableName},
			"columns": map[string]interface{}{tableName: columns},
		}, nil
	} else {
		// Get all tables
		rows, err := db.Query("SHOW TABLES")
		if err != nil {
			return map[string]interface{}{"error": fmt.Sprintf("failed to get tables: %v", err)}, nil
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				return map[string]interface{}{"error": fmt.Sprintf("scan failed: %v", err)}, nil
			}
			tables = append(tables, name)
		}

		return map[string]interface{}{
			"tables":  tables,
			"columns": map[string]interface{}{},
		}, nil
	}
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewSQLPlugin()

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
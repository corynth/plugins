# SQL Plugin

A comprehensive SQL database plugin for Corynth that supports SQLite, PostgreSQL, and MySQL databases.

## Features

- **Multi-database support**: SQLite, PostgreSQL, and MySQL
- **Full SQL operations**: SELECT queries, INSERT/UPDATE/DELETE statements
- **Schema introspection**: Get table and column information
- **Prepared statements**: Support for parameterized queries
- **Connection string parsing**: Flexible database connection formats

## Actions

### `query`
Execute SELECT queries and return structured results.

**Inputs:**
- `connection_string` (string, required): Database connection string
- `query` (string, required): SQL SELECT query to execute
- `params` (array, optional): Query parameters for prepared statements

**Outputs:**
- `rows` (array): Query result rows as array of objects
- `columns` (array): Column names
- `row_count` (number): Number of rows returned

### `execute`
Execute INSERT/UPDATE/DELETE statements.

**Inputs:**
- `connection_string` (string, required): Database connection string
- `statement` (string, required): SQL statement to execute
- `params` (array, optional): Statement parameters for prepared statements

**Outputs:**
- `affected_rows` (number): Number of rows affected
- `last_insert_id` (number): Last inserted ID (if applicable)
- `success` (boolean): Operation success status

### `schema`
Get database schema information.

**Inputs:**
- `connection_string` (string, required): Database connection string
- `table_name` (string, optional): Specific table name to get schema for

**Outputs:**
- `tables` (array): List of table names
- `columns` (object): Column information by table name

## Connection Strings

### SQLite
```
sqlite:///absolute/path/to/database.db
sqlite://relative/path/to/database.db
```

### PostgreSQL
```
postgres://username:password@hostname:port/database?sslmode=disable
postgresql://username:password@hostname:port/database?sslmode=require
```

### MySQL
```
mysql://username:password@hostname:port/database
mysql://username:password@hostname:port/database?charset=utf8mb4
```

## Examples

See the `samples/` directory for complete workflow examples:

- `sqlite-example.hcl`: SQLite operations with employee database
- `postgres-example.hcl`: PostgreSQL operations with products database
- `mysql-example.hcl`: MySQL operations with orders database

## Dependencies

The plugin automatically manages its Go dependencies:
- `github.com/mattn/go-sqlite3` for SQLite support
- `github.com/lib/pq` for PostgreSQL support  
- `github.com/go-sql-driver/mysql` for MySQL support

## Usage in Workflows

```hcl
step "create_users_table" {
  plugin = "sql"
  action = "execute"
  params = {
    connection_string = "sqlite:///tmp/app.db"
    statement = <<EOF
      CREATE TABLE users (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL
      )
    EOF
  }
}

step "insert_user" {
  plugin = "sql"
  action = "execute"
  depends_on = ["create_users_table"]
  params = {
    connection_string = "sqlite:///tmp/app.db"
    statement = "INSERT INTO users (name, email) VALUES (?, ?)"
    params = ["John Doe", "john@example.com"]
  }
}

step "query_users" {
  plugin = "sql"
  action = "query"
  depends_on = ["insert_user"]
  params = {
    connection_string = "sqlite:///tmp/app.db"
    query = "SELECT * FROM users WHERE email LIKE ?"
    params = ["%@example.com"]
  }
}
```

## Security Considerations

- Always use parameterized queries to prevent SQL injection
- Use appropriate connection string options (SSL, timeouts, etc.)
- Ensure database credentials are properly secured
- Validate user inputs before passing to database queries
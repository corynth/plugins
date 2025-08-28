// SQLite Database Operations Example
// This workflow demonstrates SQLite database operations using the SQL plugin

step "create_table" {
  plugin = "sql"
  action = "execute"
  params = {
    connection_string = "sqlite:///tmp/example.db"
    statement = <<EOF
      CREATE TABLE IF NOT EXISTS employees (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        department TEXT NOT NULL,
        salary REAL NOT NULL,
        hire_date TEXT NOT NULL
      )
    EOF
  }
}

step "insert_employees" {
  plugin = "sql"
  action = "execute"
  depends_on = ["create_table"]
  params = {
    connection_string = "sqlite:///tmp/example.db"
    statement = "INSERT INTO employees (name, department, salary, hire_date) VALUES (?, ?, ?, ?)"
    params = ["Alice Johnson", "Engineering", 95000.0, "2023-01-15"]
  }
}

step "insert_more_employees" {
  plugin = "sql"
  action = "execute"
  depends_on = ["insert_employees"]
  params = {
    connection_string = "sqlite:///tmp/example.db"
    statement = <<EOF
      INSERT INTO employees (name, department, salary, hire_date) VALUES 
      ('Bob Smith', 'Marketing', 75000.0, '2023-02-10'),
      ('Carol Davis', 'Engineering', 105000.0, '2022-11-20'),
      ('David Wilson', 'Sales', 65000.0, '2023-03-05')
    EOF
  }
}

step "query_all_employees" {
  plugin = "sql"
  action = "query"
  depends_on = ["insert_more_employees"]
  params = {
    connection_string = "sqlite:///tmp/example.db"
    query = "SELECT * FROM employees ORDER BY salary DESC"
  }
}

step "query_engineering" {
  plugin = "sql"
  action = "query"
  depends_on = ["insert_more_employees"]
  params = {
    connection_string = "sqlite:///tmp/example.db"
    query = "SELECT name, salary FROM employees WHERE department = ? ORDER BY salary DESC"
    params = ["Engineering"]
  }
}

step "get_schema" {
  plugin = "sql"
  action = "schema"
  depends_on = ["create_table"]
  params = {
    connection_string = "sqlite:///tmp/example.db"
    table_name = "employees"
  }
}

step "update_salary" {
  plugin = "sql"
  action = "execute"
  depends_on = ["query_all_employees"]
  params = {
    connection_string = "sqlite:///tmp/example.db"
    statement = "UPDATE employees SET salary = ? WHERE name = ?"
    params = [110000.0, "Alice Johnson"]
  }
}
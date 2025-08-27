workflow "database-operations" {
  description = "MySQL database operations and schema management"
  version     = "1.0.0"

  variable "mysql_host" {
    type        = string
    default     = "localhost"
    description = "MySQL server host"
  }

  variable "mysql_port" {
    type        = number
    default     = 3306
    description = "MySQL server port"
  }

  variable "database_name" {
    type        = string
    default     = "testdb"
    description = "Database name to manage"
  }

  variable "mysql_user" {
    type        = string
    default     = "admin"
    description = "MySQL username"
  }

  variable "mysql_password" {
    type        = string
    default     = "password123"
    description = "MySQL password"
  }

  step "create_database" {
    plugin = "mysql"
    action = "exec"
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      query    = "CREATE DATABASE IF NOT EXISTS ${var.database_name};"
    }
  }

  step "create_users_table" {
    plugin = "mysql"
    action = "exec"
    
    depends_on = ["create_database"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        CREATE TABLE IF NOT EXISTS users (
          id INT AUTO_INCREMENT PRIMARY KEY,
          username VARCHAR(50) NOT NULL UNIQUE,
          email VARCHAR(100) NOT NULL,
          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        );
      EOF
    }
  }

  step "insert_sample_users" {
    plugin = "mysql"
    action = "exec"
    
    depends_on = ["create_users_table"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        INSERT INTO users (username, email) VALUES
        ('john_doe', 'john@example.com'),
        ('jane_smith', 'jane@example.com'),
        ('admin_user', 'admin@example.com')
        ON DUPLICATE KEY UPDATE email = VALUES(email);
      EOF
    }
  }

  step "query_users" {
    plugin = "mysql"
    action = "query"
    
    depends_on = ["insert_sample_users"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = "SELECT id, username, email, created_at FROM users ORDER BY created_at DESC;"
    }
  }

  step "create_backup" {
    plugin = "mysql"
    action = "dump"
    
    depends_on = ["query_users"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      output   = "/tmp/${var.database_name}_backup_$(date +%Y%m%d_%H%M%S).sql"
    }
  }

  step "verify_backup" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["create_backup"]
    
    params = {
      command = "ls -la /tmp/${var.database_name}_backup_*.sql && echo 'Database operations completed successfully. Users found: ${query_users.row_count}'"
    }
  }
}
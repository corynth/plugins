workflow "schema-migration" {
  description = "MySQL schema migration and data transformation"
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
    default     = "migrationdb"
    description = "Database name for migration"
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

  step "create_migration_database" {
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

  step "create_migrations_table" {
    plugin = "mysql"
    action = "exec"
    
    depends_on = ["create_migration_database"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        CREATE TABLE IF NOT EXISTS schema_migrations (
          id INT AUTO_INCREMENT PRIMARY KEY,
          version VARCHAR(20) NOT NULL UNIQUE,
          description VARCHAR(255),
          applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
      EOF
    }
  }

  step "migration_001_products_table" {
    plugin = "mysql"
    action = "exec"
    
    depends_on = ["create_migrations_table"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        CREATE TABLE IF NOT EXISTS products (
          id INT AUTO_INCREMENT PRIMARY KEY,
          name VARCHAR(100) NOT NULL,
          price DECIMAL(10,2) NOT NULL,
          category_id INT,
          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        INSERT INTO schema_migrations (version, description) 
        VALUES ('001', 'Create products table')
        ON DUPLICATE KEY UPDATE description = VALUES(description);
      EOF
    }
  }

  step "migration_002_categories_table" {
    plugin = "mysql"
    action = "exec"
    
    depends_on = ["migration_001_products_table"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        CREATE TABLE IF NOT EXISTS categories (
          id INT AUTO_INCREMENT PRIMARY KEY,
          name VARCHAR(50) NOT NULL UNIQUE,
          description TEXT,
          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        INSERT INTO categories (name, description) VALUES
        ('Electronics', 'Electronic devices and accessories'),
        ('Books', 'Physical and digital books'),
        ('Clothing', 'Apparel and fashion items')
        ON DUPLICATE KEY UPDATE description = VALUES(description);
        
        INSERT INTO schema_migrations (version, description) 
        VALUES ('002', 'Create categories table and seed data')
        ON DUPLICATE KEY UPDATE description = VALUES(description);
      EOF
    }
  }

  step "migration_003_add_foreign_key" {
    plugin = "mysql"
    action = "exec"
    
    depends_on = ["migration_002_categories_table"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        ALTER TABLE products 
        ADD CONSTRAINT fk_products_category 
        FOREIGN KEY (category_id) REFERENCES categories(id)
        ON DELETE SET NULL;
        
        INSERT INTO schema_migrations (version, description) 
        VALUES ('003', 'Add foreign key constraint to products table')
        ON DUPLICATE KEY UPDATE description = VALUES(description);
      EOF
    }
  }

  step "seed_sample_data" {
    plugin = "mysql"
    action = "exec"
    
    depends_on = ["migration_003_add_foreign_key"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        INSERT INTO products (name, price, category_id) VALUES
        ('Laptop Pro', 1299.99, 1),
        ('Programming Guide', 49.99, 2),
        ('Casual T-Shirt', 24.99, 3),
        ('Wireless Headphones', 199.99, 1),
        ('Mystery Novel', 14.99, 2)
        ON DUPLICATE KEY UPDATE price = VALUES(price);
      EOF
    }
  }

  step "verify_migration_status" {
    plugin = "mysql"
    action = "query"
    
    depends_on = ["seed_sample_data"]
    
    params = {
      host     = var.mysql_host
      port     = var.mysql_port
      username = var.mysql_user
      password = var.mysql_password
      database = var.database_name
      query    = <<-EOF
        SELECT 
          m.version,
          m.description,
          m.applied_at,
          CASE 
            WHEN m.version = '001' THEN (SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '${var.database_name}' AND table_name = 'products')
            WHEN m.version = '002' THEN (SELECT COUNT(*) FROM categories)
            WHEN m.version = '003' THEN (SELECT COUNT(*) FROM information_schema.key_column_usage WHERE table_schema = '${var.database_name}' AND constraint_name = 'fk_products_category')
            ELSE 0
          END as verification_count
        FROM schema_migrations m
        ORDER BY m.version;
      EOF
    }
  }

  step "migration_summary" {
    plugin = "shell"
    action = "exec"
    
    depends_on = ["verify_migration_status"]
    
    params = {
      command = "echo '=== Schema Migration Summary ===' && echo 'Database: ${var.database_name}' && echo 'Migrations applied: ${verify_migration_status.row_count}' && echo 'Migration completed successfully'"
    }
  }
}
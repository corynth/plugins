// PostgreSQL Database Operations Example
// This workflow demonstrates PostgreSQL database operations using the SQL plugin
// Note: Requires a running PostgreSQL instance

step "create_table" {
  plugin = "sql"
  action = "execute"
  params = {
    connection_string = "postgres://username:password@localhost:5432/testdb?sslmode=disable"
    statement = <<EOF
      CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        category VARCHAR(100) NOT NULL,
        price DECIMAL(10,2) NOT NULL,
        stock_quantity INTEGER NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      )
    EOF
  }
}

step "insert_products" {
  plugin = "sql"
  action = "execute"
  depends_on = ["create_table"]
  params = {
    connection_string = "postgres://username:password@localhost:5432/testdb?sslmode=disable"
    statement = <<EOF
      INSERT INTO products (name, category, price, stock_quantity) VALUES 
      ($1, $2, $3, $4),
      ($5, $6, $7, $8),
      ($9, $10, $11, $12)
    EOF
    params = [
      "Laptop Pro", "Electronics", 1299.99, 25,
      "Wireless Mouse", "Electronics", 49.99, 100,
      "Office Chair", "Furniture", 299.99, 15
    ]
  }
}

step "query_products" {
  plugin = "sql"
  action = "query"
  depends_on = ["insert_products"]
  params = {
    connection_string = "postgres://username:password@localhost:5432/testdb?sslmode=disable"
    query = "SELECT id, name, category, price, stock_quantity FROM products WHERE price > $1 ORDER BY price DESC"
    params = [100.0]
  }
}

step "get_schema" {
  plugin = "sql"
  action = "schema"
  depends_on = ["create_table"]
  params = {
    connection_string = "postgres://username:password@localhost:5432/testdb?sslmode=disable"
    table_name = "products"
  }
}

step "update_stock" {
  plugin = "sql"
  action = "execute"
  depends_on = ["insert_products"]
  params = {
    connection_string = "postgres://username:password@localhost:5432/testdb?sslmode=disable"
    statement = "UPDATE products SET stock_quantity = stock_quantity - $1 WHERE name = $2"
    params = [5, "Laptop Pro"]
  }
}

step "aggregate_query" {
  plugin = "sql"
  action = "query"
  depends_on = ["insert_products"]
  params = {
    connection_string = "postgres://username:password@localhost:5432/testdb?sslmode=disable"
    query = <<EOF
      SELECT 
        category,
        COUNT(*) as product_count,
        AVG(price) as avg_price,
        SUM(stock_quantity) as total_stock
      FROM products 
      GROUP BY category 
      ORDER BY avg_price DESC
    EOF
  }
}
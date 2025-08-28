// MySQL Database Operations Example
// This workflow demonstrates MySQL database operations using the SQL plugin
// Note: Requires a running MySQL instance

step "create_table" {
  plugin = "sql"
  action = "execute"
  params = {
    connection_string = "mysql://username:password@localhost:3306/testdb"
    statement = <<EOF
      CREATE TABLE IF NOT EXISTS orders (
        id INT AUTO_INCREMENT PRIMARY KEY,
        customer_name VARCHAR(255) NOT NULL,
        product_name VARCHAR(255) NOT NULL,
        quantity INT NOT NULL DEFAULT 1,
        unit_price DECIMAL(10,2) NOT NULL,
        total_amount DECIMAL(10,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
        order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        status ENUM('pending', 'processing', 'shipped', 'delivered') DEFAULT 'pending'
      )
    EOF
  }
}

step "insert_orders" {
  plugin = "sql"
  action = "execute"
  depends_on = ["create_table"]
  params = {
    connection_string = "mysql://username:password@localhost:3306/testdb"
    statement = <<EOF
      INSERT INTO orders (customer_name, product_name, quantity, unit_price, status) VALUES 
      (?, ?, ?, ?, ?),
      (?, ?, ?, ?, ?),
      (?, ?, ?, ?, ?)
    EOF
    params = [
      "John Smith", "Laptop Pro", 1, 1299.99, "processing",
      "Alice Johnson", "Wireless Mouse", 2, 49.99, "shipped",
      "Bob Wilson", "Office Chair", 1, 299.99, "pending"
    ]
  }
}

step "query_orders" {
  plugin = "sql"
  action = "query"
  depends_on = ["insert_orders"]
  params = {
    connection_string = "mysql://username:password@localhost:3306/testdb"
    query = "SELECT id, customer_name, product_name, quantity, unit_price, total_amount, status FROM orders WHERE status = ? ORDER BY total_amount DESC"
    params = ["processing"]
  }
}

step "get_all_tables" {
  plugin = "sql"
  action = "schema"
  depends_on = ["create_table"]
  params = {
    connection_string = "mysql://username:password@localhost:3306/testdb"
  }
}

step "get_orders_schema" {
  plugin = "sql"
  action = "schema"
  depends_on = ["create_table"]
  params = {
    connection_string = "mysql://username:password@localhost:3306/testdb"
    table_name = "orders"
  }
}

step "update_order_status" {
  plugin = "sql"
  action = "execute"
  depends_on = ["insert_orders"]
  params = {
    connection_string = "mysql://username:password@localhost:3306/testdb"
    statement = "UPDATE orders SET status = ? WHERE customer_name = ? AND status = ?"
    params = ["delivered", "Alice Johnson", "shipped"]
  }
}

step "sales_report" {
  plugin = "sql"
  action = "query"
  depends_on = ["insert_orders"]
  params = {
    connection_string = "mysql://username:password@localhost:3306/testdb"
    query = <<EOF
      SELECT 
        status,
        COUNT(*) as order_count,
        SUM(total_amount) as total_revenue,
        AVG(total_amount) as avg_order_value
      FROM orders 
      GROUP BY status 
      ORDER BY total_revenue DESC
    EOF
  }
}
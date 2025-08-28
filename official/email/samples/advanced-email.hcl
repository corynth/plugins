// Advanced email features example
// Shows HTML email, attachments, and multiple recipients

// Environment variables needed:
// SMTP_SERVER, SMTP_PORT, SMTP_USER, SMTP_PASSWORD, SMTP_FROM_EMAIL
// Optional: SMTP_TLS (default: true)

step "create_report_file" {
  plugin = "file"
  action = "write"
  params = {
    path = "/tmp/monthly_report.txt"
    content = "Monthly Report\n=============\n\nSales: $10,000\nExpenses: $3,000\nProfit: $7,000"
  }
}

step "send_monthly_report" {
  plugin = "email"
  action = "send"
  depends_on = ["create_report_file"]
  params = {
    to = [
      "manager@company.com",
      "accounting@company.com",
      "ceo@company.com"
    ]
    subject = "Monthly Financial Report - ${formatdate("YYYY-MM", timestamp())}"
    body = <<-EOT
      <html>
        <body style="font-family: Arial, sans-serif;">
          <h2>Monthly Financial Report</h2>
          <p>Dear Team,</p>
          <p>Please find attached the monthly financial report for your review.</p>
          
          <div style="background-color: #f0f0f0; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h3 style="color: #2c3e50;">Key Highlights:</h3>
            <ul>
              <li><strong>Revenue:</strong> Exceeded target by 15%</li>
              <li><strong>Expenses:</strong> Reduced by 8% from last month</li>
              <li><strong>New Customers:</strong> 45 new acquisitions</li>
            </ul>
          </div>
          
          <p>If you have any questions, please don't hesitate to reach out.</p>
          <p>Best regards,<br>Finance Team</p>
        </body>
      </html>
    EOT
    html = true
    attachments = ["/tmp/monthly_report.txt"]
  }
}

step "send_plain_text_followup" {
  plugin = "email"
  action = "send"
  depends_on = ["send_monthly_report"]
  params = {
    to = ["manager@company.com"]
    subject = "Follow-up: Monthly Report"
    body = "Hi Manager,\n\nJust wanted to follow up on the monthly report I sent earlier. Please let me know if you need any additional information or clarification.\n\nThanks!"
    html = false
  }
}
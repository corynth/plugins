// Basic email sending example
// Set these environment variables before running:
// export SMTP_SERVER="smtp.gmail.com"
// export SMTP_PORT="587"
// export SMTP_USER="your-email@gmail.com"
// export SMTP_PASSWORD="your-app-password"
// export SMTP_FROM_EMAIL="your-email@gmail.com"

step "send_notification" {
  plugin = "email"
  action = "send"
  params = {
    to = ["recipient@example.com"]
    subject = "Hello from Corynth"
    body = "This is a test email sent from the Corynth email plugin!"
    html = false
  }
}

step "send_html_email" {
  plugin = "email"
  action = "send"
  params = {
    to = ["recipient@example.com", "another@example.com"]
    subject = "HTML Email Test"
    body = "<h1>Hello World!</h1><p>This is an <strong>HTML</strong> email.</p>"
    html = true
    attachments = ["/path/to/document.pdf", "/path/to/image.png"]
  }
}
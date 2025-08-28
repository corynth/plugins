package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

type EmailPlugin struct{}

func NewEmailPlugin() *EmailPlugin {
	return &EmailPlugin{}
}

func (p *EmailPlugin) GetMetadata() Metadata {
	return Metadata{
		Name:        "email",
		Version:     "1.0.0",
		Description: "Email notifications and communication",
		Author:      "Corynth Team",
		Tags:        []string{"email", "smtp", "notifications"},
	}
}

func (p *EmailPlugin) GetActions() map[string]ActionSpec {
	return map[string]ActionSpec{
		"send": {
			Description: "Send email",
			Inputs: map[string]IOSpec{
				"to": {
					Type:        "array",
					Required:    true,
					Description: "Recipient emails",
				},
				"subject": {
					Type:        "string",
					Required:    true,
					Description: "Email subject",
				},
				"body": {
					Type:        "string",
					Required:    true,
					Description: "Email body",
				},
				"from_email": {
					Type:        "string",
					Required:    false,
					Description: "Sender email",
				},
				"attachments": {
					Type:        "array",
					Required:    false,
					Description: "File paths to attach",
				},
				"html": {
					Type:        "boolean",
					Required:    false,
					Default:     false,
					Description: "HTML email",
				},
			},
			Outputs: map[string]IOSpec{
				"success":    {Type: "boolean", Description: "Email sent successfully"},
				"message_id": {Type: "string", Description: "Message ID"},
			},
		},
	}
}

func (p *EmailPlugin) Execute(action string, params map[string]interface{}) (map[string]interface{}, error) {
	switch action {
	case "send":
		return p.sendEmail(params)
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func (p *EmailPlugin) sendEmail(params map[string]interface{}) (map[string]interface{}, error) {
	// Parse and validate parameters
	toEmails, err := p.parseToEmails(params["to"])
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}

	subject, ok := params["subject"].(string)
	if !ok || subject == "" {
		return map[string]interface{}{"error": "subject is required"}, nil
	}

	body, ok := params["body"].(string)
	if !ok || body == "" {
		return map[string]interface{}{"error": "body is required"}, nil
	}

	// Get from email - use parameter or environment variable
	fromEmail := ""
	if fe, ok := params["from_email"].(string); ok && fe != "" {
		fromEmail = fe
	} else {
		fromEmail = os.Getenv("SMTP_FROM_EMAIL")
	}
	if fromEmail == "" {
		return map[string]interface{}{"error": "from_email is required (parameter or SMTP_FROM_EMAIL env var)"}, nil
	}

	// Parse attachments
	attachments := p.parseAttachments(params["attachments"])

	// Check if HTML email
	isHTML := getBoolParam(params, "html", false)

	// Get SMTP configuration from environment
	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		smtpServer = "localhost"
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		smtpPortStr = "587"
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("invalid SMTP_PORT: %v", err)}, nil
	}

	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpTLS := getBoolFromEnv("SMTP_TLS", true)

	// Build email message
	message, messageID, err := p.buildMessage(fromEmail, toEmails, subject, body, isHTML, attachments)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to build message: %v", err)}, nil
	}

	// Send email
	err = p.sendSMTP(smtpServer, smtpPort, smtpUser, smtpPass, smtpTLS, fromEmail, toEmails, message)
	if err != nil {
		return map[string]interface{}{"error": fmt.Sprintf("failed to send email: %v", err)}, nil
	}

	return map[string]interface{}{
		"success":    true,
		"message_id": messageID,
	}, nil
}

func (p *EmailPlugin) parseToEmails(to interface{}) ([]string, error) {
	if to == nil {
		return nil, fmt.Errorf("to is required")
	}

	switch v := to.(type) {
	case []interface{}:
		emails := make([]string, 0, len(v))
		for i, email := range v {
			if emailStr, ok := email.(string); ok && emailStr != "" {
				emails = append(emails, emailStr)
			} else {
				return nil, fmt.Errorf("to[%d] must be a non-empty string", i)
			}
		}
		if len(emails) == 0 {
			return nil, fmt.Errorf("at least one recipient email is required")
		}
		return emails, nil
	case string:
		if v == "" {
			return nil, fmt.Errorf("to email cannot be empty")
		}
		return []string{v}, nil
	default:
		return nil, fmt.Errorf("to must be a string or array of strings")
	}
}

func (p *EmailPlugin) parseAttachments(attachments interface{}) []string {
	if attachments == nil {
		return nil
	}

	if arr, ok := attachments.([]interface{}); ok {
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
		return result
	}
	return nil
}

func (p *EmailPlugin) buildMessage(fromEmail string, toEmails []string, subject, body string, isHTML bool, attachments []string) ([]byte, string, error) {
	// Generate a message ID
	messageID := fmt.Sprintf("<%d@corynth-email-plugin>", generateTimestamp())

	var message strings.Builder

	// Headers
	message.WriteString(fmt.Sprintf("From: %s\r\n", fromEmail))
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(toEmails, ", ")))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString(fmt.Sprintf("Message-ID: %s\r\n", messageID))
	message.WriteString("MIME-Version: 1.0\r\n")

	boundary := fmt.Sprintf("boundary_%d", generateTimestamp())

	if len(attachments) > 0 {
		// Multipart message with attachments
		message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		message.WriteString("\r\n")

		// Body part
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		if isHTML {
			message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		} else {
			message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		}
		message.WriteString("Content-Transfer-Encoding: 8bit\r\n")
		message.WriteString("\r\n")
		message.WriteString(body)
		message.WriteString("\r\n")

		// Attachment parts
		for _, filePath := range attachments {
			if err := p.addAttachment(&message, boundary, filePath); err != nil {
				return nil, "", fmt.Errorf("failed to add attachment %s: %v", filePath, err)
			}
		}

		message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// Simple message without attachments
		if isHTML {
			message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		} else {
			message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		}
		message.WriteString("Content-Transfer-Encoding: 8bit\r\n")
		message.WriteString("\r\n")
		message.WriteString(body)
		message.WriteString("\r\n")
	}

	return []byte(message.String()), messageID, nil
}

func (p *EmailPlugin) addAttachment(message *strings.Builder, boundary, filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Read file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Get filename and detect content type
	filename := filepath.Base(filePath)
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Write attachment headers
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filename))
	message.WriteString("\r\n")

	// Encode file content in base64
	encoded := base64.StdEncoding.EncodeToString(fileContent)
	
	// Write base64 content with line breaks every 76 characters (RFC 2045)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		message.WriteString(encoded[i:end])
		message.WriteString("\r\n")
	}

	return nil
}

func (p *EmailPlugin) sendSMTP(server string, port int, username, password string, useTLS bool, from string, to []string, message []byte) error {
	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", server, port)
	
	var client *smtp.Client
	var err error

	if useTLS && port == 465 {
		// SMTP over SSL (SMTPS)
		tlsConfig := &tls.Config{
			ServerName: server,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %v", err)
		}
		client, err = smtp.NewClient(conn, server)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %v", err)
		}
	} else {
		// Regular SMTP connection
		client, err = smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %v", err)
		}
	}
	defer client.Close()

	// Start TLS if needed and not already using SSL
	if useTLS && port != 465 {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName: server,
			}
			if err = client.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("failed to start TLS: %v", err)
			}
		}
	}

	// Authenticate if credentials provided
	if username != "" && password != "" {
		auth := smtp.PlainAuth("", username, password, server)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %v", err)
		}
	}

	// Send email
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient %s: %v", addr, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	return nil
}

// Helper functions
func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key].(bool); ok {
		return val
	}
	return defaultValue
}

func getBoolFromEnv(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultValue
}

func generateTimestamp() int64 {
	return time.Now().UnixNano()
}

func main() {
	if len(os.Args) < 2 {
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{"error": "action required"})
		os.Exit(1)
	}

	action := os.Args[1]
	plugin := NewEmailPlugin()

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
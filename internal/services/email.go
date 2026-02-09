package services

import (
	"fmt"
	"net/smtp"
	"strconv"
)

// EmailService handles sending emails.
type EmailService struct {
	host     string
	port     int
	user     string
	pass     string
	from     string
}

// NewEmailService creates a new email service.
func NewEmailService(host string, port int, user, pass, from string) *EmailService {
	return &EmailService{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
	}
}

// SendOrderConfirmation sends an order confirmation email to the customer.
func (e *EmailService) SendOrderConfirmation(to, customerName string, orderID int64, total int64) error {
	if e.user == "" || e.pass == "" {
		// Skip sending if credentials not configured
		return nil
	}

	subject := "Order Confirmation - Bhomanshah"
	body := fmt.Sprintf(`Dear %s,

Thank you for your order!

Order ID: %d
Total: PKR %d

Your order has been received and is being processed. You will receive updates on your order status.

Best regards,
Bhomanshah Team
`, customerName, orderID, total)

	return e.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP.
func (e *EmailService) sendEmail(to, subject, body string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", e.user, e.pass, e.host)

	// Construct the message
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	// Send the email
	addr := e.host + ":" + strconv.Itoa(e.port)
	err := smtp.SendMail(addr, auth, e.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
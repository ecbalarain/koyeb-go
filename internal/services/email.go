package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// EmailService handles sending emails via Brevo API.
type EmailService struct {
	apiKey string
	from   string
}

// NewEmailService creates a new email service.
func NewEmailService(apiKey, from string) *EmailService {
	return &EmailService{
		apiKey: apiKey,
		from:   from,
	}
}

// SendOrderConfirmation sends an order confirmation email to the customer.
func (e *EmailService) SendOrderConfirmation(to, customerName string, orderID int64, total int64) error {
	if e.apiKey == "" {
		// Skip sending if API key not configured
		return nil
	}

	subject := "Order Confirmation - Bhomanshah"
	htmlContent := fmt.Sprintf(`<html><head></head><body>
<p>Dear %s,</p>
<p>Thank you for your order!</p>
<p>Order ID: %d<br>
Total: PKR %d</p>
<p>Your order has been received and is being processed. You will receive updates on your order status.</p>
<p>Best regards,<br>
Bhomanshah Team</p>
</body></html>`, customerName, orderID, total)

	return e.sendEmail(to, subject, htmlContent)
}

// sendEmail sends an email using Brevo API.
func (e *EmailService) sendEmail(to, subject, htmlContent string) error {
	url := "https://api.brevo.com/v3/smtp/email"

	payload := map[string]interface{}{
		"sender": map[string]string{
			"email": e.from,
			"name":  "Bhomanshah",
		},
		"to": []map[string]string{
			{
				"email": to,
			},
		},
		"subject":     subject,
		"htmlContent": htmlContent,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("api-key", e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return nil
}
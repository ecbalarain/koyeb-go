package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/koyeb/example-golang/internal/models"
)

// EmailService handles sending emails via Brevo API.
type EmailService struct {
	apiKey string
	from   string
	orderStatusBase string
}

// NewEmailService creates a new email service.
func NewEmailService(apiKey, from, orderStatusBase string) *EmailService {
	return &EmailService{
		apiKey: apiKey,
		from:   from,
		orderStatusBase: orderStatusBase,
	}
}

// SendOrderConfirmation sends an order confirmation email to the customer.
func (e *EmailService) SendOrderConfirmation(to, customerName string, orderID int64, total int64, items []models.OrderItem) error {
	if e.apiKey == "" {
		// Skip sending if API key not configured
		return nil
	}

	subject := "Order Confirmation - Bhomanshah"
	statusURL := e.orderStatusURL(orderID, to)
	itemsHTML := e.renderOrderItems(items)
	statusLinkHTML := ""
	if statusURL != "" {
		statusLinkHTML = fmt.Sprintf(`<p><a href="%s" style="color:#18181b;text-decoration:underline;">View order status</a></p>`, statusURL)
	}

	htmlContent := fmt.Sprintf(`<html><head></head><body style="font-family:Arial,Helvetica,sans-serif;color:#111827;">
<p>Dear %s,</p>
<p>Thank you for your order!</p>
<p><strong>Order ID:</strong> %d<br>
<strong>Total:</strong> PKR %d</p>
%s
%s
<p>Your order has been received and is being processed. You will receive updates on your order status.</p>
<p>Best regards,<br>
Bhomanshah Team</p>
</body></html>`, html.EscapeString(customerName), orderID, total, itemsHTML, statusLinkHTML)

	return e.sendEmail(to, subject, htmlContent)
}

func (e *EmailService) orderStatusURL(orderID int64, email string) string {
	base := strings.TrimSpace(e.orderStatusBase)
	if base == "" {
		return ""
	}
	base = strings.TrimRight(base, "/")
	return fmt.Sprintf("%s/order-status.html?order_id=%d&email=%s", base, orderID, email)
}

func (e *EmailService) renderOrderItems(items []models.OrderItem) string {
	if len(items) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<table style="width:100%;border-collapse:collapse;margin:16px 0;">`)
	sb.WriteString(`<thead><tr>`)
	sb.WriteString(`<th align="left" style="border-bottom:1px solid #e5e7eb;padding:6px 0;">Item</th>`)
	sb.WriteString(`<th align="left" style="border-bottom:1px solid #e5e7eb;padding:6px 0;">Variant</th>`)
	sb.WriteString(`<th align="right" style="border-bottom:1px solid #e5e7eb;padding:6px 0;">Qty</th>`)
	sb.WriteString(`<th align="right" style="border-bottom:1px solid #e5e7eb;padding:6px 0;">Line total</th>`)
	sb.WriteString(`</tr></thead><tbody>`)

	for _, item := range items {
		lineTotal := item.PriceAtPurchase * int64(item.Qty)
		name := html.EscapeString(item.ProductName)
		variant := html.EscapeString(item.VariantLabel)
		if variant == "" {
			variant = "-"
		}
		sb.WriteString(`<tr>`) 
		sb.WriteString(fmt.Sprintf(`<td style="padding:6px 0;">%s</td>`, name))
		sb.WriteString(fmt.Sprintf(`<td style="padding:6px 0;">%s</td>`, variant))
		sb.WriteString(fmt.Sprintf(`<td align="right" style="padding:6px 0;">%d</td>`, item.Qty))
		sb.WriteString(fmt.Sprintf(`<td align="right" style="padding:6px 0;">PKR %d</td>`, lineTotal))
		sb.WriteString(`</tr>`)
	}

	sb.WriteString(`</tbody></table>`)
	return sb.String()
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
package models

import "time"

// Order represents a customer order (COD).
type Order struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customer_name"`
	CustomerPhone   string    `json:"customer_phone"`
	CustomerAddress string    `json:"customer_address"`
	CustomerCity    string    `json:"customer_city"`
	Total           int64     `json:"total"` // total in smallest currency unit
	Status          string    `json:"status"`  // new, confirmed, shipped, canceled
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// OrderItem represents a single line item in an order.
type OrderItem struct {
	ID              int64  `json:"id"`
	OrderID         int64  `json:"order_id"`
	VariantID       int64  `json:"variant_id"`
	ProductName     string `json:"product_name"`
	VariantLabel    string `json:"variant_label"` // e.g. "Red / M"
	PriceAtPurchase int64  `json:"price_at_purchase"`
	Qty             int    `json:"qty"`
}

// CreateOrderRequest is the request body for creating an order.
type CreateOrderRequest struct {
	CustomerName    string             `json:"customer_name" validate:"required,min=2,max=100"`
	CustomerPhone   string             `json:"customer_phone" validate:"required,min=10,max=15"`
	CustomerAddress string             `json:"customer_address" validate:"required,min=5,max=500"`
	CustomerCity    string             `json:"customer_city" validate:"required,min=2,max=100"`
	Notes           string             `json:"notes" validate:"max=500"`
	Items           []CreateOrderItem  `json:"items" validate:"required,min=1,dive"`
}

// CreateOrderItem represents a single item in the order creation request.
type CreateOrderItem struct {
	VariantID int64 `json:"variant_id" validate:"required,gt=0"`
	Qty       int   `json:"qty" validate:"required,gt=0,lte=10"`
}

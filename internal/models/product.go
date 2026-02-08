package models

import "time"

// Product represents a product in the store.
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Images      string    `json:"images"` // JSON array stored as text
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Variant represents a size/color variant of a product.
type Variant struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	Price     int64     `json:"price"`  // price in smallest currency unit (paisa/cents)
	Stock     int       `json:"stock"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

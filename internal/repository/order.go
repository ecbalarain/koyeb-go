package repository

import (
	"database/sql"
	"fmt"

	"github.com/koyeb/example-golang/internal/models"
)

// OrderRepository handles database operations for orders.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order with its items in a transaction.
func (r *OrderRepository) Create(order *models.Order, items []models.OrderItem) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	query := `
		INSERT INTO orders (customer_name, customer_phone, customer_address, customer_city, customer_email, total, status, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query,
		order.CustomerName, order.CustomerPhone, order.CustomerAddress,
		order.CustomerCity, order.CustomerEmail, order.Total, order.Status, order.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	order.ID = orderID

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (order_id, variant_id, product_name, variant_label, price_at_purchase, qty)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	for _, item := range items {
		_, err := tx.Exec(itemQuery,
			orderID, item.VariantID, item.ProductName, item.VariantLabel,
			item.PriceAtPurchase, item.Qty,
		)
		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}

		// Decrement stock atomically
		stockQuery := `
			UPDATE variants
			SET stock = stock - ?
			WHERE id = ? AND stock >= ?
		`
		result, err := tx.Exec(stockQuery, item.Qty, item.VariantID, item.Qty)
		if err != nil {
			return fmt.Errorf("failed to decrement stock for variant %d: %w", item.VariantID, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to check stock decrement result: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("insufficient stock for variant %d", item.VariantID)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAll returns all orders, optionally filtered by status.
func (r *OrderRepository) GetAll(status string) ([]models.Order, error) {
	query := `
		SELECT id, customer_name, customer_phone, customer_address, customer_city, customer_email,
		       total, status, notes, created_at, updated_at
		FROM orders
	`

	var args []interface{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID, &o.CustomerName, &o.CustomerPhone, &o.CustomerAddress, &o.CustomerCity, &o.CustomerEmail,
			&o.Total, &o.Status, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return orders, nil
}

// GetByID returns an order by its ID.
func (r *OrderRepository) GetByID(id int64) (*models.Order, error) {
	query := `
		SELECT id, customer_name, customer_phone, customer_address, customer_city, customer_email,
		       total, status, notes, created_at, updated_at
		FROM orders
		WHERE id = ?
	`

	var o models.Order
	err := r.db.QueryRow(query, id).Scan(
		&o.ID, &o.CustomerName, &o.CustomerPhone, &o.CustomerAddress, &o.CustomerCity, &o.CustomerEmail,
		&o.Total, &o.Status, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query order by ID: %w", err)
	}

	return &o, nil
}

// GetItemsByOrderID returns all items for an order.
func (r *OrderRepository) GetItemsByOrderID(orderID int64) ([]models.OrderItem, error) {
	query := `
		SELECT id, order_id, variant_id, product_name, variant_label, price_at_purchase, qty
		FROM order_items
		WHERE order_id = ?
		ORDER BY id
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.VariantID, &item.ProductName,
			&item.VariantLabel, &item.PriceAtPurchase, &item.Qty,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return items, nil
}

// UpdateStatus updates the status of an order.
func (r *OrderRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

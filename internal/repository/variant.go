package repository

import (
	"database/sql"
	"fmt"

	"github.com/koyeb/example-golang/internal/models"
)

// VariantRepository handles database operations for variants.
type VariantRepository struct {
	db *sql.DB
}

// NewVariantRepository creates a new variant repository.
func NewVariantRepository(db *sql.DB) *VariantRepository {
	return &VariantRepository{db: db}
}

// GetByProductID returns all variants for a product.
func (r *VariantRepository) GetByProductID(productID int64, activeOnly bool) ([]models.Variant, error) {
	query := `
		SELECT id, product_id, size, color, price, stock, active, created_at, updated_at
		FROM variants
		WHERE product_id = ?
	`
	if activeOnly {
		query += " AND active = TRUE"
	}
	query += " ORDER BY size, color"

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	var variants []models.Variant
	for rows.Next() {
		var v models.Variant
		err := rows.Scan(
			&v.ID, &v.ProductID, &v.Size, &v.Color, &v.Price,
			&v.Stock, &v.Active, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}
		variants = append(variants, v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return variants, nil
}

// GetByID returns a variant by its ID.
func (r *VariantRepository) GetByID(id int64) (*models.Variant, error) {
	query := `
		SELECT id, product_id, size, color, price, stock, active, created_at, updated_at
		FROM variants
		WHERE id = ?
	`

	var v models.Variant
	err := r.db.QueryRow(query, id).Scan(
		&v.ID, &v.ProductID, &v.Size, &v.Color, &v.Price,
		&v.Stock, &v.Active, &v.CreatedAt, &v.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query variant by ID: %w", err)
	}

	return &v, nil
}

// Create creates a new variant.
func (r *VariantRepository) Create(v *models.Variant) error {
	query := `
		INSERT INTO variants (product_id, size, color, price, stock, active)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, v.ProductID, v.Size, v.Color, v.Price, v.Stock, v.Active)
	if err != nil {
		return fmt.Errorf("failed to create variant: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	v.ID = id
	return nil
}

// Update updates an existing variant.
func (r *VariantRepository) Update(v *models.Variant) error {
	query := `
		UPDATE variants
		SET size = ?, color = ?, price = ?, stock = ?, active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(query, v.Size, v.Color, v.Price, v.Stock, v.Active, v.ID)
	if err != nil {
		return fmt.Errorf("failed to update variant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}

// Delete deletes a variant.
func (r *VariantRepository) Delete(id int64) error {
	query := `DELETE FROM variants WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}

// DecrementStock atomically decrements the stock of a variant.
// Returns an error if stock is insufficient.
func (r *VariantRepository) DecrementStock(id int64, qty int) error {
	query := `
		UPDATE variants
		SET stock = stock - ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND stock >= ? AND active = TRUE
	`

	result, err := r.db.Exec(query, qty, id, qty)
	if err != nil {
		return fmt.Errorf("failed to decrement stock: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("insufficient stock or variant not active")
	}

	return nil
}

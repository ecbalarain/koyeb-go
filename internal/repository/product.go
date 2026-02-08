package repository

import (
	"database/sql"
	"fmt"

	"github.com/koyeb/example-golang/internal/models"
)

// ProductRepository handles database operations for products.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll returns all products (optionally filtered by active status).
func (r *ProductRepository) GetAll(activeOnly bool) ([]models.Product, error) {
	query := `
		SELECT id, name, slug, description, category, images, active, created_at, updated_at
		FROM products
	`
	if activeOnly {
		query += " WHERE active = TRUE"
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.Category,
			&p.Images, &p.Active, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, nil
}

// GetByID returns a product by its ID.
func (r *ProductRepository) GetByID(id int64) (*models.Product, error) {
	query := `
		SELECT id, name, slug, description, category, images, active, created_at, updated_at
		FROM products
		WHERE id = ?
	`

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.Category,
		&p.Images, &p.Active, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product by ID: %w", err)
	}

	return &p, nil
}

// GetBySlug returns a product by its slug.
func (r *ProductRepository) GetBySlug(slug string) (*models.Product, error) {
	query := `
		SELECT id, name, slug, description, category, images, active, created_at, updated_at
		FROM products
		WHERE slug = ?
	`

	var p models.Product
	err := r.db.QueryRow(query, slug).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.Category,
		&p.Images, &p.Active, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product by slug: %w", err)
	}

	return &p, nil
}

// SetActive toggles the active status of a product.
func (r *ProductRepository) SetActive(id int64, active bool) error {
	query := `UPDATE products SET active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`

	result, err := r.db.Exec(query, active, id)
	if err != nil {
		return fmt.Errorf("failed to update product active status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Create creates a new product.
func (r *ProductRepository) Create(name, slug, description, category, images string) (int64, error) {
	query := `
		INSERT INTO products (name, slug, description, category, images, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := r.db.Exec(query, name, slug, description, category, images)
	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// Update updates an existing product.
func (r *ProductRepository) Update(id int64, name, slug, description, category, images string) error {
	query := `
		UPDATE products 
		SET name = ?, slug = ?, description = ?, category = ?, images = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(query, name, slug, description, category, images, id)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

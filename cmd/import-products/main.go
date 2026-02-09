package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/koyeb/example-golang/internal/config"
	"github.com/koyeb/example-golang/internal/database"
	_ "github.com/go-sql-driver/mysql"
)

type productInput struct {
	ID               int64             `json:"id"`
	Name             string            `json:"name"`
	Slug             string            `json:"slug"`
	Description      string            `json:"description"`
	ShortDescription string            `json:"shortDescription"`
	Category         string            `json:"category"`
	Categories       []string          `json:"categories"`
	Tags             []string          `json:"tags"`
	Images           []string          `json:"images"`
	Sizes            []string          `json:"sizes"`
	Attributes       []attributeInput  `json:"attributes"`
	Prices           productPrices     `json:"prices"`
	SourceURL        string            `json:"sourceUrl"`
	InStock          bool              `json:"inStock"`
}

type attributeInput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type productPrices struct {
	Currency     string `json:"currency"`
	Price        *int64 `json:"price"`
	RegularPrice *int64 `json:"regularPrice"`
	SalePrice    *int64 `json:"salePrice"`
}

const (
	defaultStock = 50
	defaultColor = "Default"
)

func main() {
	productsPath := flag.String("products", "cloudflare-pages-frontend/products.json", "Path to products.json")
	deactivateMissing := flag.Bool("deactivate-missing", false, "Deactivate DB products not present in products.json")
	flag.Parse()

	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	inputs, err := readProducts(*productsPath)
	if err != nil {
		log.Fatalf("Failed to read products: %v", err)
	}

	if err := importProducts(db, inputs, *deactivateMissing); err != nil {
		log.Fatalf("Import failed: %v", err)
	}

	log.Printf("Imported %d products", len(inputs))
}

func readProducts(path string) ([]productInput, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var inputs []productInput
	if err := json.Unmarshal(data, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func importProducts(db *sql.DB, inputs []productInput, deactivateMissing bool) error {
	seen := make(map[string]struct{})

	for _, p := range inputs {
		slug := p.Slug
		if slug == "" {
			return fmt.Errorf("product %d missing slug", p.ID)
		}
		seen[slug] = struct{}{}

		productID, err := upsertProduct(db, p)
		if err != nil {
			return fmt.Errorf("upsert product %s: %w", slug, err)
		}

		if err := replaceVariants(db, productID, p); err != nil {
			return fmt.Errorf("replace variants for %s: %w", slug, err)
		}
	}

	if deactivateMissing {
		if err := deactivateMissingProducts(db, seen); err != nil {
			return fmt.Errorf("deactivate missing products: %w", err)
		}
	}

	return nil
}

func upsertProduct(db *sql.DB, p productInput) (int64, error) {
	imagesJSON, err := json.Marshal(p.Images)
	if err != nil {
		return 0, err
	}

	description := p.Description
	if description == "" {
		description = p.ShortDescription
	}

	var productID int64
	row := db.QueryRow("SELECT id FROM products WHERE slug = ?", p.Slug)
	switch err := row.Scan(&productID); err {
	case nil:
		_, err := db.Exec(
			`UPDATE products SET name = ?, description = ?, category = ?, images = ?, active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
			p.Name, description, p.Category, string(imagesJSON), true, productID,
		)
		if err != nil {
			return 0, err
		}
		return productID, nil
	case sql.ErrNoRows:
		result, err := db.Exec(
			`INSERT INTO products (name, slug, description, category, images, active, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
			p.Name, p.Slug, description, p.Category, string(imagesJSON),
		)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	default:
		return 0, err
	}
}

func replaceVariants(db *sql.DB, productID int64, p productInput) error {
	sizes := p.Sizes
	if len(sizes) == 0 {
		sizes = []string{"One Size"}
	}

	price := resolvePrice(p.Prices)
	stock := 0
	if p.InStock {
		stock = defaultStock
	}

	if _, err := db.Exec("DELETE FROM variants WHERE product_id = ?", productID); err != nil {
		return err
	}

	for _, size := range sizes {
		_, err := db.Exec(
			`INSERT INTO variants (product_id, size, color, price, stock, active)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			productID, size, defaultColor, price, stock, p.InStock,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func resolvePrice(prices productPrices) int64 {
	if prices.Price != nil {
		return *prices.Price
	}
	if prices.SalePrice != nil {
		return *prices.SalePrice
	}
	if prices.RegularPrice != nil {
		return *prices.RegularPrice
	}
	return 0
}

func deactivateMissingProducts(db *sql.DB, seen map[string]struct{}) error {
	rows, err := db.Query("SELECT id, slug FROM products")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var slug string
		if err := rows.Scan(&id, &slug); err != nil {
			return err
		}
		if _, ok := seen[slug]; ok {
			continue
		}
		if _, err := db.Exec("UPDATE products SET active = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id = ?", id); err != nil {
			return err
		}
	}

	return rows.Err()
}

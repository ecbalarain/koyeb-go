-- Migration: Add customer_email to orders table
ALTER TABLE orders ADD COLUMN customer_email VARCHAR(255) NOT NULL AFTER customer_city;
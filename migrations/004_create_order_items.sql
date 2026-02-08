-- Migration: Create order_items table
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    variant_id BIGINT NOT NULL,
    product_name VARCHAR(255) NOT NULL, -- Snapshot at time of order
    variant_label VARCHAR(100) NOT NULL, -- e.g. "Red / M"
    price_at_purchase BIGINT NOT NULL, -- Price at time of purchase
    qty INT NOT NULL,

    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (variant_id) REFERENCES variants(id),
    INDEX idx_order_id (order_id)
);
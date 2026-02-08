-- Migration: Create variants table
CREATE TABLE IF NOT EXISTS variants (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    size VARCHAR(10) NOT NULL, -- S, M, L
    color VARCHAR(50) NOT NULL,
    price BIGINT NOT NULL, -- Price in smallest currency unit (paisa/cents)
    stock INT NOT NULL DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE KEY unique_variant (product_id, size, color),
    INDEX idx_product_id (product_id),
    INDEX idx_active (active)
);
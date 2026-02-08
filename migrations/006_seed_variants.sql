-- Seed: Demo variants (only if not exists)
-- Canvas Tote variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'S', 'Black', 2800, 15, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'S' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'S', 'Navy', 2800, 12, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'S' AND color = 'Navy');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'S', 'Gray', 2800, 8, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'S' AND color = 'Gray');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'M', 'Black', 2800, 20, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'M' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'M', 'Navy', 2800, 18, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'M' AND color = 'Navy');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'M', 'Gray', 2800, 10, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'M' AND color = 'Gray');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'L', 'Black', 2800, 25, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'L' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'L', 'Navy', 2800, 22, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'L' AND color = 'Navy');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 1, 'L', 'Gray', 2800, 15, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 1 AND size = 'L' AND color = 'Gray');

-- Stone Mug variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 2, 'S', 'White', 1400, 30, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 2 AND size = 'S' AND color = 'White');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 2, 'S', 'Gray', 1400, 25, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 2 AND size = 'S' AND color = 'Gray');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 2, 'S', 'Black', 1400, 20, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 2 AND size = 'S' AND color = 'Black');

-- Minimal Tee variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'S', 'White', 2200, 12, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'S' AND color = 'White');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'S', 'Black', 2200, 15, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'S' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'S', 'Gray', 2200, 10, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'S' AND color = 'Gray');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'M', 'White', 2200, 18, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'M' AND color = 'White');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'M', 'Black', 2200, 20, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'M' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'M', 'Gray', 2200, 14, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'M' AND color = 'Gray');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'L', 'White', 2200, 16, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'L' AND color = 'White');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'L', 'Black', 2200, 22, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'L' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 3, 'L', 'Gray', 2200, 12, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 3 AND size = 'L' AND color = 'Gray');

-- Desk Lamp variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 4, 'S', 'Black', 3600, 8, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 4 AND size = 'S' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 4, 'S', 'White', 3600, 10, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 4 AND size = 'S' AND color = 'White');

-- Wireless Charger variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 5, 'S', 'Black', 2900, 15, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 5 AND size = 'S' AND color = 'Black');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 5, 'S', 'White', 2900, 12, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 5 AND size = 'S' AND color = 'White');

-- Key Organizer variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 6, 'S', 'Brown', 1800, 20, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 6 AND size = 'S' AND color = 'Brown');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 6, 'S', 'Black', 1800, 18, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 6 AND size = 'S' AND color = 'Black');

-- Scented Candle variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 7, 'S', 'Lavender', 1900, 25, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 7 AND size = 'S' AND color = 'Lavender');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 7, 'S', 'Vanilla', 1900, 22, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 7 AND size = 'S' AND color = 'Vanilla');

-- Notebook Set variants
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 8, 'S', 'Blue', 1600, 30, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 8 AND size = 'S' AND color = 'Blue');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 8, 'S', 'Green', 1600, 28, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 8 AND size = 'S' AND color = 'Green');
INSERT INTO variants (product_id, size, color, price, stock, active)
SELECT 8, 'S', 'Red', 1600, 25, TRUE WHERE NOT EXISTS (SELECT 1 FROM variants WHERE product_id = 8 AND size = 'S' AND color = 'Red');
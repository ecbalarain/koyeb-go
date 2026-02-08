-- Seed: Demo products (only if not exists)
INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Canvas Tote', 'canvas-tote', 'A durable canvas tote bag perfect for everyday use.', 'Basics', '["/images/canvas-tote.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'canvas-tote');

INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Stone Mug', 'stone-mug', 'Ceramic mug with a stone-like finish, holds 12oz.', 'Home', '["/images/stone-mug.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'stone-mug');

INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Minimal Tee', 'minimal-tee', 'Soft cotton t-shirt with a clean, minimal design.', 'Basics', '["/images/minimal-tee.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'minimal-tee');

INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Desk Lamp', 'desk-lamp', 'LED desk lamp with adjustable brightness and USB charging.', 'Home', '["/images/desk-lamp.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'desk-lamp');

INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Wireless Charger', 'wireless-charger', 'Fast wireless charging pad compatible with all Qi devices.', 'Tech', '["/images/wireless-charger.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'wireless-charger');

INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Key Organizer', 'key-organizer', 'Leather key organizer with multiple compartments.', 'Tech', '["/images/key-organizer.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'key-organizer');

INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Scented Candle', 'scented-candle', 'Soy wax candle with lavender scent, burns for 40 hours.', 'Home', '["/images/scented-candle.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'scented-candle');

INSERT INTO products (name, slug, description, category, images, active)
SELECT 'Notebook Set', 'notebook-set', 'Set of 3 lined notebooks with recycled paper covers.', 'Basics', '["/images/notebook-set.jpg"]', TRUE
WHERE NOT EXISTS (SELECT 1 FROM products WHERE slug = 'notebook-set');
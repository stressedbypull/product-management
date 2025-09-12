-- Insert 3 categories 
INSERT INTO categories (code, name) VALUES
('CAT001', 'CLOTHING'),
('CAT002', 'SHOES'),
('CAT003', 'ACCESSORIES');

-- link categories to products

-- Link products to "Clothing"
UPDATE products SET category_id = (SELECT id FROM categories WHERE name = 'CLOTHING')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

-- Link products to "Shoes"
UPDATE products SET category_id = (SELECT id FROM categories WHERE name = 'SHOES')
WHERE code IN ('PROD002', 'PROD006');

-- Link products to "Accessories"
UPDATE products SET category_id = (SELECT id FROM categories WHERE name = 'ACCESSORIES')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');

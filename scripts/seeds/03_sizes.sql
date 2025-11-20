-- =============================================
-- SEED: Tallas
-- =============================================
-- Crear tallas para camisetas, pantalones y zapatos

-- Tallas de Camisetas (SHIRT)
INSERT INTO sizes (type, value, "order", is_active, created_at, updated_at)
VALUES 
    ('SHIRT', 'XS', 1, true, NOW(), NOW()),
    ('SHIRT', 'S', 2, true, NOW(), NOW()),
    ('SHIRT', 'M', 3, true, NOW(), NOW()),
    ('SHIRT', 'L', 4, true, NOW(), NOW()),
    ('SHIRT', 'XL', 5, true, NOW(), NOW()),
    ('SHIRT', 'XXL', 6, true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Tallas de Pantalones (PANTS) - Cintura en pulgadas
INSERT INTO sizes (type, value, "order", is_active, created_at, updated_at)
VALUES 
    ('PANTS', '24', 1, true, NOW(), NOW()),
    ('PANTS', '26', 2, true, NOW(), NOW()),
    ('PANTS', '28', 3, true, NOW(), NOW()),
    ('PANTS', '30', 4, true, NOW(), NOW()),
    ('PANTS', '32', 5, true, NOW(), NOW()),
    ('PANTS', '34', 6, true, NOW(), NOW()),
    ('PANTS', '36', 7, true, NOW(), NOW()),
    ('PANTS', '38', 8, true, NOW(), NOW()),
    ('PANTS', '40', 9, true, NOW(), NOW()),
    ('PANTS', '42', 10, true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Tallas de Zapatos (SHOES) - Sistema US
INSERT INTO sizes (type, value, "order", is_active, created_at, updated_at)
VALUES 
    ('SHOES', '5', 1, true, NOW(), NOW()),
    ('SHOES', '5.5', 2, true, NOW(), NOW()),
    ('SHOES', '6', 3, true, NOW(), NOW()),
    ('SHOES', '6.5', 4, true, NOW(), NOW()),
    ('SHOES', '7', 5, true, NOW(), NOW()),
    ('SHOES', '7.5', 6, true, NOW(), NOW()),
    ('SHOES', '8', 7, true, NOW(), NOW()),
    ('SHOES', '8.5', 8, true, NOW(), NOW()),
    ('SHOES', '9', 9, true, NOW(), NOW()),
    ('SHOES', '9.5', 10, true, NOW(), NOW()),
    ('SHOES', '10', 11, true, NOW(), NOW()),
    ('SHOES', '10.5', 12, true, NOW(), NOW()),
    ('SHOES', '11', 13, true, NOW(), NOW()),
    ('SHOES', '11.5', 14, true, NOW(), NOW()),
    ('SHOES', '12', 15, true, NOW(), NOW()),
    ('SHOES', '13', 16, true, NOW(), NOW()),
    ('SHOES', '14', 17, true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Resumen:
-- Total de tallas: 33
--   - Camisetas: 6
--   - Pantalones: 10
--   - Zapatos: 17

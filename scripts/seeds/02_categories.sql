-- =============================================
-- SEED: Categorías
-- =============================================
-- Crear categorías iniciales para productos de cuero

INSERT INTO categories (name, description, is_active, created_at, updated_at)
VALUES 
    ('Chaquetas', 'Chaquetas y abrigos de cuero', true, NOW(), NOW()),
    ('Pantalones', 'Pantalones y jeans de cuero', true, NOW(), NOW()),
    ('Camisas', 'Camisas y blusas', true, NOW(), NOW()),
    ('Accesorios', 'Cinturones, carteras y más', true, NOW(), NOW()),
    ('Calzado', 'Zapatos y botas de cuero', true, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

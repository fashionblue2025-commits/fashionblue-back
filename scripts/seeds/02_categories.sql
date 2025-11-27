-- =============================================
-- SEED: Categorías
-- =============================================
-- Crear categorías iniciales para productos de cuero

INSERT INTO categories (name, description, is_active, created_at, updated_at)
VALUES 
    ('Chaquetas', 'Chaquetas y abrigos de cuero', true, NOW(), NOW()),
ON CONFLICT (name) DO NOTHING;

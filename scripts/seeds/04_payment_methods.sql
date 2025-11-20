-- =============================================
-- SEED: Métodos de Pago
-- =============================================
-- Crear métodos de pago disponibles

INSERT INTO payment_methods (name, is_active, created_at, updated_at)
VALUES 
    ('NEQUI Sonia', true, NOW(), NOW()),
    ('NEQUI Jhon', true, NOW(), NOW()),
    ('Daviplata', true, NOW(), NOW()),
    ('Efectivo', true, NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Resumen:
-- Total de métodos de pago: 4

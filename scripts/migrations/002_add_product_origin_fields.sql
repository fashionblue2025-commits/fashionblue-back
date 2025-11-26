-- Migración: Agregar campos de origen a la tabla products
-- Fecha: 2024-11-20
-- Descripción: Rastrear si un producto fue creado por cotización o para inventario

-- Agregar columna origin_type
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS origin_type VARCHAR(20) DEFAULT 'INVENTORY';

-- Agregar columna origin_order_id
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS origin_order_id INTEGER;

-- Agregar columna is_custom
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS is_custom BOOLEAN DEFAULT FALSE;

-- Agregar foreign key constraint
ALTER TABLE products 
ADD CONSTRAINT IF NOT EXISTS fk_products_origin_order 
FOREIGN KEY (origin_order_id) REFERENCES orders(id) ON DELETE SET NULL;

-- Agregar índices
CREATE INDEX IF NOT EXISTS idx_products_origin_type ON products(origin_type);
CREATE INDEX IF NOT EXISTS idx_products_origin_order_id ON products(origin_order_id);
CREATE INDEX IF NOT EXISTS idx_products_is_custom ON products(is_custom);

-- Comentarios
COMMENT ON COLUMN products.origin_type IS 'Origen del producto: CUSTOM (por cotización), INVENTORY (para stock)';
COMMENT ON COLUMN products.origin_order_id IS 'ID de la orden que originó este producto (si es CUSTOM)';
COMMENT ON COLUMN products.is_custom IS 'true: producto único/personalizado, false: producto estándar';

-- ============================================================================
-- Migración 005: Refactorizar Productos a Producto Base + Variantes
-- Fecha: 2024-11-21
-- Descripción: 
--   - Crea tabla product_variants para manejar stock por color/talla
--   - Migra datos de products a product_variants
--   - Elimina columnas de variante de products (color, size_id, stock, reserved_stock)
--   - Actualiza order_items para referenciar product_variant_id
-- ============================================================================

BEGIN;

-- ============================================================================
-- PASO 1: Crear tabla product_variants
-- ============================================================================

CREATE TABLE IF NOT EXISTS product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    color VARCHAR(50) NOT NULL,
    size_id INTEGER,
    stock INTEGER NOT NULL DEFAULT 0,
    reserved_stock INTEGER NOT NULL DEFAULT 0,
    unit_price DECIMAL(10,2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_product_variants_product 
        FOREIGN KEY (product_id) 
        REFERENCES products(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_product_variants_size 
        FOREIGN KEY (size_id) 
        REFERENCES sizes(id) 
        ON DELETE SET NULL
);

-- Índices para product_variants
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_color ON product_variants(color);
CREATE INDEX idx_product_variants_size_id ON product_variants(size_id);
CREATE INDEX idx_product_variants_is_active ON product_variants(is_active);

-- Índice compuesto para búsqueda de variantes únicas
CREATE UNIQUE INDEX idx_product_variants_unique 
    ON product_variants(product_id, color, COALESCE(size_id, 0));

COMMENT ON TABLE product_variants IS 'Variantes de productos por color y talla con stock independiente';
COMMENT ON COLUMN product_variants.product_id IS 'Referencia al producto base';
COMMENT ON COLUMN product_variants.color IS 'Color de esta variante';
COMMENT ON COLUMN product_variants.size_id IS 'Talla de esta variante (puede ser NULL)';
COMMENT ON COLUMN product_variants.stock IS 'Stock disponible de esta variante';
COMMENT ON COLUMN product_variants.reserved_stock IS 'Stock reservado por órdenes aprobadas';
COMMENT ON COLUMN product_variants.unit_price IS 'Precio de esta variante (puede diferir del producto base)';

-- ============================================================================
-- PASO 2: Migrar datos existentes de products a product_variants
-- ============================================================================

-- Solo migrar productos que tengan color o size_id definidos
-- (productos que ya son "variantes" en el sistema actual)
INSERT INTO product_variants (
    product_id,
    color,
    size_id,
    stock,
    reserved_stock,
    unit_price,
    is_active,
    created_at,
    updated_at
)
SELECT 
    id as product_id,
    COALESCE(color, 'Sin Color') as color,  -- Si no tiene color, usar 'Sin Color'
    size_id,
    COALESCE(stock, 0) as stock,
    COALESCE(reserved_stock, 0) as reserved_stock,
    unit_price,
    is_active,
    created_at,
    updated_at
FROM products
WHERE color IS NOT NULL OR size_id IS NOT NULL OR stock > 0 OR reserved_stock > 0;

-- Log de migración
DO $$
DECLARE
    migrated_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO migrated_count FROM product_variants;
    RAISE NOTICE 'Migrated % product variants from products table', migrated_count;
END $$;

-- ============================================================================
-- PASO 3: Agregar columna product_variant_id a order_items
-- ============================================================================

-- Agregar nueva columna
ALTER TABLE order_items 
ADD COLUMN IF NOT EXISTS product_variant_id INTEGER;

-- Crear índice
CREATE INDEX IF NOT EXISTS idx_order_items_product_variant_id 
    ON order_items(product_variant_id);

-- Agregar foreign key
ALTER TABLE order_items
ADD CONSTRAINT fk_order_items_product_variant
    FOREIGN KEY (product_variant_id)
    REFERENCES product_variants(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;

COMMENT ON COLUMN order_items.product_variant_id IS 'Referencia a la variante específica (color + talla)';

-- ============================================================================
-- PASO 4: Migrar product_id a product_variant_id en order_items
-- ============================================================================

-- Intentar mapear order_items existentes a product_variants
-- Esto es un "best effort" - algunos items pueden no tener variante exacta
UPDATE order_items oi
SET product_variant_id = pv.id
FROM product_variants pv
WHERE oi.product_id = pv.product_id
  AND COALESCE(oi.color, 'Sin Color') = pv.color
  AND (
    (oi.size_id IS NULL AND pv.size_id IS NULL) OR
    (oi.size_id = pv.size_id)
  )
  AND oi.product_variant_id IS NULL;

-- Log de migración
DO $$
DECLARE
    mapped_count INTEGER;
    unmapped_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO mapped_count 
    FROM order_items 
    WHERE product_variant_id IS NOT NULL;
    
    SELECT COUNT(*) INTO unmapped_count 
    FROM order_items 
    WHERE product_variant_id IS NULL AND product_id IS NOT NULL;
    
    RAISE NOTICE 'Mapped % order items to product variants', mapped_count;
    RAISE NOTICE '% order items could not be mapped (will be created as new variants)', unmapped_count;
END $$;

-- ============================================================================
-- PASO 5: Hacer color NOT NULL en order_items
-- ============================================================================

-- Actualizar valores NULL a 'Sin Color'
UPDATE order_items 
SET color = 'Sin Color' 
WHERE color IS NULL OR color = '';

-- Hacer la columna NOT NULL
ALTER TABLE order_items 
ALTER COLUMN color SET NOT NULL;

-- ============================================================================
-- PASO 6: Eliminar columnas de variante de products
-- ============================================================================

-- Eliminar columnas que ahora están en product_variants
ALTER TABLE products DROP COLUMN IF EXISTS color;
ALTER TABLE products DROP COLUMN IF EXISTS size_id;
ALTER TABLE products DROP COLUMN IF EXISTS stock;
ALTER TABLE products DROP COLUMN IF EXISTS reserved_stock;

-- Actualizar valores por defecto para productos base
ALTER TABLE products 
ALTER COLUMN material_cost SET DEFAULT 0,
ALTER COLUMN labor_cost SET DEFAULT 0,
ALTER COLUMN production_cost SET DEFAULT 0;

COMMENT ON TABLE products IS 'Productos base (maestros) sin variantes de color/talla';

-- ============================================================================
-- PASO 7: Crear vista para compatibilidad (opcional)
-- ============================================================================

-- Vista que combina productos con sus variantes para consultas legacy
CREATE OR REPLACE VIEW products_with_variants AS
SELECT 
    p.id as product_id,
    p.name as product_name,
    p.description,
    p.category_id,
    p.material_cost,
    p.labor_cost,
    p.production_cost,
    p.unit_price as base_unit_price,
    p.wholesale_price,
    p.min_wholesale_qty,
    p.min_stock,
    p.is_active as product_is_active,
    pv.id as variant_id,
    pv.color,
    pv.size_id,
    pv.stock,
    pv.reserved_stock,
    pv.stock - pv.reserved_stock as available_stock,
    pv.unit_price as variant_unit_price,
    pv.is_active as variant_is_active,
    p.created_at,
    p.updated_at
FROM products p
LEFT JOIN product_variants pv ON p.id = pv.product_id;

COMMENT ON VIEW products_with_variants IS 'Vista de compatibilidad que combina productos base con sus variantes';

-- ============================================================================
-- PASO 8: Actualizar triggers de updated_at
-- ============================================================================

-- Trigger para product_variants
CREATE OR REPLACE FUNCTION update_product_variants_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_product_variants_updated_at
    BEFORE UPDATE ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_product_variants_updated_at();

-- ============================================================================
-- VERIFICACIÓN FINAL
-- ============================================================================

DO $$
DECLARE
    product_count INTEGER;
    variant_count INTEGER;
    order_item_count INTEGER;
    mapped_items INTEGER;
BEGIN
    SELECT COUNT(*) INTO product_count FROM products;
    SELECT COUNT(*) INTO variant_count FROM product_variants;
    SELECT COUNT(*) INTO order_item_count FROM order_items;
    SELECT COUNT(*) INTO mapped_items FROM order_items WHERE product_variant_id IS NOT NULL;
    
    RAISE NOTICE '========================================';
    RAISE NOTICE 'MIGRATION 005 COMPLETED SUCCESSFULLY';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Products (base): %', product_count;
    RAISE NOTICE 'Product Variants: %', variant_count;
    RAISE NOTICE 'Order Items (total): %', order_item_count;
    RAISE NOTICE 'Order Items (mapped to variants): %', mapped_items;
    RAISE NOTICE '========================================';
END $$;

COMMIT;

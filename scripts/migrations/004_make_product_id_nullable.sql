-- Migración: Hacer product_id nullable y agregar quantity_reserved en order_items
-- Fecha: 2024-11-21
-- Descripción: 
--   1. Permite crear órdenes CUSTOM/INVENTORY donde el producto se crea después de aprobar la cotización
--   2. Agrega campo quantity_reserved para trackear stock reservado vs stock a fabricar

-- 1. Modificar la columna product_id para que sea nullable
ALTER TABLE order_items 
ALTER COLUMN product_id DROP NOT NULL;

COMMENT ON COLUMN order_items.product_id IS 'ID del producto. Puede ser NULL para órdenes CUSTOM/INVENTORY donde el producto se crea después de aprobar la cotización.';

-- 2. Verificar que product_name sea NOT NULL (ya debería serlo)
ALTER TABLE order_items 
ALTER COLUMN product_name SET NOT NULL;

COMMENT ON COLUMN order_items.product_name IS 'Nombre del producto. Requerido para poder crear el producto cuando se apruebe la orden.';

-- 3. Agregar columna quantity_reserved
ALTER TABLE order_items 
ADD COLUMN IF NOT EXISTS quantity_reserved INTEGER DEFAULT 0 NOT NULL;

COMMENT ON COLUMN order_items.quantity_reserved IS 'Cantidad reservada del stock existente. La diferencia (quantity - quantity_reserved) es lo que se debe fabricar.';

-- 4. Actualizar registros existentes (por defecto no tienen reservas)
UPDATE order_items 
SET quantity_reserved = 0 
WHERE quantity_reserved IS NULL;

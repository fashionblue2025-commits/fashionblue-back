-- Agregar columna reserved_quantity a order_items
-- Esta columna trackea cuántas unidades del item fueron reservadas del stock existente

ALTER TABLE order_items 
ADD COLUMN IF NOT EXISTS reserved_quantity INTEGER NOT NULL DEFAULT 0;

-- Comentario de la columna
COMMENT ON COLUMN order_items.reserved_quantity IS 'Cantidad de unidades reservadas del stock existente (el resto debe fabricarse)';

-- Nota: Para items existentes, reserved_quantity será 0
-- Esto significa que se asumirá que necesitan fabricación completa
-- Lo cual es seguro para órdenes antiguas

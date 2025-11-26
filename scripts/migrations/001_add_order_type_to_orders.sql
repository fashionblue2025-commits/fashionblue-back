-- Migración: Agregar order_type a la tabla orders
-- Fecha: 2024-11-20
-- Descripción: Agregar campo para diferenciar entre CUSTOM, INVENTORY y SALE

-- Agregar columna order_type
ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS order_type VARCHAR(20) DEFAULT 'CUSTOM';

-- Actualizar órdenes existentes según su tipo anterior
-- INTERNAL y EXTERNAL se convierten en CUSTOM (producción por demanda)
UPDATE orders 
SET order_type = 'CUSTOM' 
WHERE type IN ('INTERNAL', 'EXTERNAL') OR type IS NULL;

-- Agregar índice para mejorar consultas
CREATE INDEX IF NOT EXISTS idx_orders_order_type ON orders(order_type);

-- Comentario en la columna
COMMENT ON COLUMN orders.order_type IS 'Tipo de orden: CUSTOM (cotización), INVENTORY (producción para stock), SALE (venta de existente)';

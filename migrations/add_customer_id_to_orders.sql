-- Agregar customer_id a orders para clientes internos
-- Cuando una orden se completa para un cliente interno, se crea una transacción contable

ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_id INTEGER DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);

-- Comentario de la columna
COMMENT ON COLUMN orders.customer_id IS 'ID del cliente interno (opcional). Si está presente, se creará una transacción contable al completar la venta';

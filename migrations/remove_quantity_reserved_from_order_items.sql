-- Eliminar columna quantity_reserved de order_items
-- La cantidad reservada ahora se obtiene directamente de product_variants.reserved_stock

ALTER TABLE order_items DROP COLUMN IF EXISTS quantity_reserved;

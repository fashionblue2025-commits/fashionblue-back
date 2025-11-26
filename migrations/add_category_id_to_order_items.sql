-- Agregar columna category_id a order_items
-- Esta columna almacena un snapshot de la categoría del producto
-- para poder crear productos con la categoría correcta desde órdenes de inventario

ALTER TABLE order_items 
ADD COLUMN category_id INTEGER NOT NULL DEFAULT 1;

-- Crear índice para mejorar performance en búsquedas por categoría
CREATE INDEX idx_order_items_category_id ON order_items(category_id);

-- Agregar comentario a la columna
COMMENT ON COLUMN order_items.category_id IS 'Snapshot de la categoría del producto al momento de crear el item';

-- Seed: Product Photos
-- Fecha: 2024-11-21
-- Descripción: Datos de ejemplo para fotos de productos
-- Nota: Este seed asume que ya existen productos en la base de datos

-- Limpiar datos existentes (opcional - comentar si no quieres eliminar datos)
-- TRUNCATE product_photos CASCADE;

-- Insertar fotos de ejemplo para productos existentes
-- Nota: Ajusta los product_id según los productos que tengas en tu base de datos

-- Ejemplo: Fotos para producto ID 1 (si existe)
INSERT INTO product_photos (product_id, photo_url, description, is_primary, display_order, uploaded_at)
SELECT 1, 'https://via.placeholder.com/800x800/8B4513/FFFFFF?text=Chaqueta+Cuero+1', 'Vista frontal de la chaqueta', true, 1, NOW()
WHERE EXISTS (SELECT 1 FROM products WHERE id = 1)
ON CONFLICT DO NOTHING;

INSERT INTO product_photos (product_id, photo_url, description, is_primary, display_order, uploaded_at)
SELECT 1, 'https://via.placeholder.com/800x800/8B4513/FFFFFF?text=Chaqueta+Cuero+2', 'Vista lateral de la chaqueta', false, 2, NOW()
WHERE EXISTS (SELECT 1 FROM products WHERE id = 1)
ON CONFLICT DO NOTHING;

INSERT INTO product_photos (product_id, photo_url, description, is_primary, display_order, uploaded_at)
SELECT 1, 'https://via.placeholder.com/800x800/8B4513/FFFFFF?text=Chaqueta+Cuero+3', 'Detalle de costuras', false, 3, NOW()
WHERE EXISTS (SELECT 1 FROM products WHERE id = 1)
ON CONFLICT DO NOTHING;

-- Ejemplo: Fotos para producto ID 2 (si existe)
INSERT INTO product_photos (product_id, photo_url, description, is_primary, display_order, uploaded_at)
SELECT 2, 'https://via.placeholder.com/800x800/654321/FFFFFF?text=Pantalon+1', 'Vista frontal del pantalón', true, 1, NOW()
WHERE EXISTS (SELECT 1 FROM products WHERE id = 2)
ON CONFLICT DO NOTHING;

INSERT INTO product_photos (product_id, photo_url, description, is_primary, display_order, uploaded_at)
SELECT 2, 'https://via.placeholder.com/800x800/654321/FFFFFF?text=Pantalon+2', 'Vista trasera del pantalón', false, 2, NOW()
WHERE EXISTS (SELECT 1 FROM products WHERE id = 2)
ON CONFLICT DO NOTHING;

-- Ejemplo: Fotos para producto ID 3 (si existe)
INSERT INTO product_photos (product_id, photo_url, description, is_primary, display_order, uploaded_at)
SELECT 3, 'https://via.placeholder.com/800x800/2F4F4F/FFFFFF?text=Camisa+1', 'Vista frontal de la camisa', true, 1, NOW()
WHERE EXISTS (SELECT 1 FROM products WHERE id = 3)
ON CONFLICT DO NOTHING;

INSERT INTO product_photos (product_id, photo_url, description, is_primary, display_order, uploaded_at)
SELECT 3, 'https://via.placeholder.com/800x800/2F4F4F/FFFFFF?text=Camisa+2', 'Detalle del cuello', false, 2, NOW()
WHERE EXISTS (SELECT 1 FROM products WHERE id = 3)
ON CONFLICT DO NOTHING;

-- Verificar cuántas fotos se insertaron
DO $$
DECLARE
    photo_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO photo_count FROM product_photos;
    RAISE NOTICE 'Total de fotos insertadas: %', photo_count;
END $$;

-- Mostrar resumen de fotos por producto
SELECT 
    p.id as product_id,
    p.name as product_name,
    COUNT(pp.id) as total_photos,
    SUM(CASE WHEN pp.is_primary THEN 1 ELSE 0 END) as primary_photos
FROM products p
LEFT JOIN product_photos pp ON p.id = pp.product_id
GROUP BY p.id, p.name
ORDER BY p.id;

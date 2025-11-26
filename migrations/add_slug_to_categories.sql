-- Agregar columna slug a categories
-- El slug identifica el tipo de tallas que usa cada categoría

ALTER TABLE categories 
ADD COLUMN slug VARCHAR(50) NOT NULL DEFAULT 'clothing';

-- Crear índice único para el slug
CREATE UNIQUE INDEX idx_categories_slug ON categories(slug);

-- Actualizar slugs según las categorías existentes
-- Ajusta estos valores según tus categorías reales
UPDATE categories SET slug = 'clothing' WHERE name ILIKE '%chaqueta%' OR name ILIKE '%camisa%' OR name ILIKE '%pantalon%' OR name ILIKE '%vestido%';
UPDATE categories SET slug = 'shoes' WHERE name ILIKE '%zapato%' OR name ILIKE '%calzado%';
UPDATE categories SET slug = 'accessories' WHERE name ILIKE '%accesorio%' OR name ILIKE '%bolso%' OR name ILIKE '%cartera%';

-- Agregar comentario a la columna
COMMENT ON COLUMN categories.slug IS 'Identificador único para el tipo de tallas (clothing, shoes, accessories, etc)';

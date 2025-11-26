-- Migración: Crear tabla product_photos
-- Fecha: 2024-11-21
-- Descripción: Tabla para almacenar fotos de productos

-- Crear tabla product_photos
CREATE TABLE IF NOT EXISTS product_photos (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    photo_url VARCHAR(500) NOT NULL,
    description TEXT,
    is_primary BOOLEAN DEFAULT FALSE,
    display_order INTEGER DEFAULT 0,
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_product_photos_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);

-- Crear índices para mejorar el rendimiento
CREATE INDEX IF NOT EXISTS idx_product_photos_product_id ON product_photos(product_id);
CREATE INDEX IF NOT EXISTS idx_product_photos_is_primary ON product_photos(is_primary);

-- Comentarios de la tabla
COMMENT ON TABLE product_photos IS 'Almacena las fotos de los productos';
COMMENT ON COLUMN product_photos.id IS 'ID único de la foto';
COMMENT ON COLUMN product_photos.product_id IS 'ID del producto al que pertenece la foto';
COMMENT ON COLUMN product_photos.photo_url IS 'URL de la foto (local o Cloudinary)';
COMMENT ON COLUMN product_photos.description IS 'Descripción opcional de la foto';
COMMENT ON COLUMN product_photos.is_primary IS 'Indica si es la foto principal del producto';
COMMENT ON COLUMN product_photos.display_order IS 'Orden de visualización de la foto';
COMMENT ON COLUMN product_photos.uploaded_at IS 'Fecha y hora de carga de la foto';
COMMENT ON COLUMN product_photos.created_at IS 'Fecha de creación del registro';
COMMENT ON COLUMN product_photos.updated_at IS 'Fecha de última actualización';

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_product_photos_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_product_photos_updated_at
    BEFORE UPDATE ON product_photos
    FOR EACH ROW
    EXECUTE FUNCTION update_product_photos_updated_at();

-- Mensaje de confirmación
DO $$
BEGIN
    RAISE NOTICE 'Tabla product_photos creada exitosamente';
END $$;

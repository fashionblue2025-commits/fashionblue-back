-- Crear tabla de permisos de usuario por categoría
CREATE TABLE IF NOT EXISTS user_category_permissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    can_view BOOLEAN DEFAULT true,
    can_create BOOLEAN DEFAULT false,
    can_edit BOOLEAN DEFAULT false,
    can_delete BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Un usuario no puede tener múltiples permisos para la misma categoría
    UNIQUE(user_id, category_id)
);

-- Índices para mejorar el rendimiento
CREATE INDEX idx_user_category_permissions_user_id ON user_category_permissions(user_id);
CREATE INDEX idx_user_category_permissions_category_id ON user_category_permissions(category_id);

-- Trigger para actualizar updated_at
CREATE OR REPLACE FUNCTION update_user_category_permissions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_category_permissions_updated_at
    BEFORE UPDATE ON user_category_permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_user_category_permissions_updated_at();

-- Comentarios
COMMENT ON TABLE user_category_permissions IS 'Permisos granulares de usuarios sobre categorías específicas';
COMMENT ON COLUMN user_category_permissions.can_view IS 'Usuario puede ver productos de esta categoría';
COMMENT ON COLUMN user_category_permissions.can_create IS 'Usuario puede crear productos en esta categoría';
COMMENT ON COLUMN user_category_permissions.can_edit IS 'Usuario puede editar productos de esta categoría';
COMMENT ON COLUMN user_category_permissions.can_delete IS 'Usuario puede eliminar productos de esta categoría';

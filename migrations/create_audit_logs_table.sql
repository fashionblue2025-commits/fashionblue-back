-- Crear tabla de logs de auditoría
-- Esta tabla almacena todos los eventos del sistema para trazabilidad y compliance

CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    order_id INTEGER NOT NULL,
    order_number VARCHAR(50),
    user_id INTEGER,
    user_name VARCHAR(100),
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    description TEXT,
    metadata JSONB,
    ip_address VARCHAR(45),
    user_agent VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Índices para mejorar rendimiento de consultas
CREATE INDEX IF NOT EXISTS idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_order_id ON audit_logs(order_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_order_created ON audit_logs(order_id, created_at DESC);

-- Comentarios de la tabla
COMMENT ON TABLE audit_logs IS 'Registro completo de auditoría de eventos del sistema';
COMMENT ON COLUMN audit_logs.event_type IS 'Tipo de evento (ej: order.approved, order.cancelled)';
COMMENT ON COLUMN audit_logs.order_id IS 'ID de la orden relacionada';
COMMENT ON COLUMN audit_logs.order_number IS 'Número de orden para referencia rápida';
COMMENT ON COLUMN audit_logs.user_id IS 'ID del usuario que generó el evento (si aplica)';
COMMENT ON COLUMN audit_logs.old_status IS 'Estado anterior de la orden';
COMMENT ON COLUMN audit_logs.new_status IS 'Nuevo estado de la orden';
COMMENT ON COLUMN audit_logs.description IS 'Descripción legible del evento';
COMMENT ON COLUMN audit_logs.metadata IS 'Datos adicionales del evento en formato JSON';
COMMENT ON COLUMN audit_logs.ip_address IS 'Dirección IP del usuario (si aplica)';
COMMENT ON COLUMN audit_logs.user_agent IS 'User agent del navegador (si aplica)';
COMMENT ON COLUMN audit_logs.created_at IS 'Fecha y hora del evento';

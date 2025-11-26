-- Crear tabla de transacciones financieras (unifica inyecciones de capital y gastos)
CREATE TABLE IF NOT EXISTS financial_transactions (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    category VARCHAR(50) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    description TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para mejorar el rendimiento
CREATE INDEX IF NOT EXISTS idx_financial_transactions_type ON financial_transactions(type);
CREATE INDEX IF NOT EXISTS idx_financial_transactions_category ON financial_transactions(category);
CREATE INDEX IF NOT EXISTS idx_financial_transactions_date ON financial_transactions(date);
CREATE INDEX IF NOT EXISTS idx_financial_transactions_type_date ON financial_transactions(type, date);

-- Comentarios
COMMENT ON TABLE financial_transactions IS 'Tabla unificada para transacciones financieras (ingresos y gastos)';
COMMENT ON COLUMN financial_transactions.type IS 'Tipo de transacción: INCOME (ingreso) o EXPENSE (gasto)';
COMMENT ON COLUMN financial_transactions.category IS 'Categoría de la transacción (ej: INVESTMENT, OPERATIONAL, etc)';
COMMENT ON COLUMN financial_transactions.amount IS 'Monto de la transacción (siempre positivo)';
COMMENT ON COLUMN financial_transactions.description IS 'Descripción detallada de la transacción';
COMMENT ON COLUMN financial_transactions.date IS 'Fecha de la transacción';

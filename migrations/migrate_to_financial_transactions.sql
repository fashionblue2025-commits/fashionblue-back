-- Migración para consolidar capital_injections y expenses en financial_transactions

-- 1. Crear la nueva tabla unificada
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

-- 2. Migrar datos de capital_injections (si existe)
INSERT INTO financial_transactions (type, category, amount, description, date, created_at, updated_at)
SELECT 
    'INCOME' as type,
    COALESCE(type, 'OTHER') as category,  -- Usar el campo type como category
    amount,
    description,
    date,
    created_at,
    updated_at
FROM capital_injections
WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'capital_injections')
ON CONFLICT DO NOTHING;

-- 3. Migrar datos de expenses (si existe)
INSERT INTO financial_transactions (type, category, amount, description, date, created_at, updated_at)
SELECT 
    'EXPENSE' as type,
    category,
    amount,
    description,
    date,
    created_at,
    updated_at
FROM expenses
WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'expenses')
ON CONFLICT DO NOTHING;

-- 4. Crear índices para mejorar el rendimiento
CREATE INDEX IF NOT EXISTS idx_financial_transactions_type ON financial_transactions(type);
CREATE INDEX IF NOT EXISTS idx_financial_transactions_category ON financial_transactions(category);
CREATE INDEX IF NOT EXISTS idx_financial_transactions_date ON financial_transactions(date);
CREATE INDEX IF NOT EXISTS idx_financial_transactions_type_date ON financial_transactions(type, date);

-- 5. Comentarios
COMMENT ON TABLE financial_transactions IS 'Tabla unificada para transacciones financieras (ingresos y gastos)';
COMMENT ON COLUMN financial_transactions.type IS 'Tipo de transacción: INCOME (ingreso) o EXPENSE (gasto)';
COMMENT ON COLUMN financial_transactions.category IS 'Categoría de la transacción (ej: INVESTMENT, OPERATIONAL, etc)';
COMMENT ON COLUMN financial_transactions.amount IS 'Monto de la transacción (siempre positivo)';
COMMENT ON COLUMN financial_transactions.description IS 'Descripción detallada de la transacción';
COMMENT ON COLUMN financial_transactions.date IS 'Fecha de la transacción';

-- 6. OPCIONAL: Eliminar tablas antiguas (comentado por seguridad)
-- DROP TABLE IF EXISTS capital_injections;
-- DROP TABLE IF EXISTS expenses;

-- Nota: Descomentar las líneas DROP solo después de verificar que la migración fue exitosa
-- y que todos los datos se migraron correctamente a financial_transactions

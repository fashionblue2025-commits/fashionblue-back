-- =============================================
-- SEED: Usuarios
-- =============================================
-- Crear usuario Super Admin
-- Email: admin@fashionblue.com
-- Password: admin123

INSERT INTO users (email, password, first_name, last_name, role, is_active, created_at, updated_at)
VALUES (
    'admin@fashionblue.com',
    '$2a$10$YourHashedPasswordHere', -- Se actualizará con el hash real
    'Super',
    'Admin',
    'SUPER_ADMIN',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Nota: El hash de bcrypt para "admin123" es:
-- $2a$10$90e/EkkrasxGkCqaKj8vPe5cPgXm5IF7zELZWSxStEBFG/IxPd2VG
-- Actualizar con el hash correcto:
UPDATE users 
SET password = '$2a$10$90e/EkkrasxGkCqaKj8vPe5cPgXm5IF7zELZWSxStEBFG/IxPd2VG'
WHERE email = 'admin@fashionblue.com';

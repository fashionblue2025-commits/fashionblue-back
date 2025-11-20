-- =============================================
-- SEED COMPLETO - Fashion Blue
-- =============================================
-- Este archivo ejecuta todos los seeds en orden
-- Ejecutar: psql -U fashionblue -d fashionblue_db -f scripts/seeds/seed_all.sql

\echo '🌱 Starting database seeding...'
\echo '================================'

-- 1. Usuarios
\echo ''
\echo '👤 Creating Super Admin...'
\i scripts/seeds/01_users.sql

-- 2. Categorías
\echo ''
\echo '📁 Creating Categories...'
\i scripts/seeds/02_categories.sql

-- 3. Tallas
\echo ''
\echo '📏 Creating Sizes...'
\i scripts/seeds/03_sizes.sql

-- 4. Métodos de Pago
\echo ''
\echo '💳 Creating Payment Methods...'
\i scripts/seeds/04_payment_methods.sql

\echo ''
\echo '================================'
\echo '✅ DATABASE SEEDED SUCCESSFULLY!'
\echo '================================'
\echo ''
\echo '📊 Summary:'
\echo '   👤 Users: 1 (Super Admin)'
\echo '   📁 Categories: 5'
\echo '   📏 Sizes: 33 total'
\echo '      - Shirts: 6'
\echo '      - Pants: 10'
\echo '      - Shoes: 17'
\echo '   💳 Payment Methods: 4'
\echo ''
\echo '🔐 Login credentials:'
\echo '   Email: admin@fashionblue.com'
\echo '   Password: admin123'
\echo ''

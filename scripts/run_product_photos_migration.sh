#!/bin/bash

# Script para ejecutar la migración de product_photos
# Autor: Fashion Blue Team
# Fecha: 2024-11-21

set -e  # Salir si hay algún error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Migración: Product Photos Table${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Cargar variables de entorno desde .env
if [ -f .env ]; then
    echo -e "${GREEN}✓${NC} Cargando variables de entorno desde .env"
    export $(cat .env | grep -v '^#' | xargs)
else
    echo -e "${RED}✗${NC} Archivo .env no encontrado"
    exit 1
fi

# Configurar variables de conexión
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-fashionblue}"
DB_NAME="${DB_NAME:-fashionblue_db}"

echo -e "${BLUE}Configuración:${NC}"
echo -e "  Host: ${DB_HOST}"
echo -e "  Port: ${DB_PORT}"
echo -e "  Database: ${DB_NAME}"
echo -e "  User: ${DB_USER}"
echo ""

# Verificar conexión a PostgreSQL
echo -e "${YELLOW}⏳${NC} Verificando conexión a PostgreSQL..."
if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c '\q' 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Conexión exitosa a PostgreSQL"
else
    echo -e "${RED}✗${NC} No se pudo conectar a PostgreSQL"
    echo -e "${YELLOW}Verifica que PostgreSQL esté corriendo y las credenciales sean correctas${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}⏳${NC} Ejecutando migración..."

# Ejecutar la migración
MIGRATION_FILE="scripts/migrations/003_create_product_photos_table.sql"

if [ ! -f "$MIGRATION_FILE" ]; then
    echo -e "${RED}✗${NC} Archivo de migración no encontrado: $MIGRATION_FILE"
    exit 1
fi

if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f "$MIGRATION_FILE"; then
    echo -e "${GREEN}✓${NC} Migración ejecutada exitosamente"
else
    echo -e "${RED}✗${NC} Error al ejecutar la migración"
    exit 1
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Migración completada${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Verificar que la tabla se creó correctamente
echo -e "${YELLOW}⏳${NC} Verificando tabla product_photos..."
PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c "
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'product_photos'
ORDER BY ordinal_position;
"

echo ""
echo -e "${GREEN}✓ Tabla product_photos creada y verificada${NC}"
echo ""
echo -e "${BLUE}Siguiente paso:${NC}"
echo -e "  Puedes ejecutar el seed de fotos de productos con:"
echo -e "  ${YELLOW}./scripts/run_product_photos_seed.sh${NC}"
echo ""

#!/bin/bash

# Script para ejecutar el seed de product_photos
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
echo -e "${BLUE}  Seed: Product Photos${NC}"
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

# Verificar que la tabla product_photos existe
echo ""
echo -e "${YELLOW}⏳${NC} Verificando que la tabla product_photos existe..."
if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c "SELECT 1 FROM product_photos LIMIT 1;" > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Tabla product_photos encontrada"
else
    echo -e "${RED}✗${NC} La tabla product_photos no existe"
    echo -e "${YELLOW}Ejecuta primero: ./scripts/run_product_photos_migration.sh${NC}"
    exit 1
fi

# Verificar que existen productos
echo ""
echo -e "${YELLOW}⏳${NC} Verificando que existen productos..."
PRODUCT_COUNT=$(PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -t -c "SELECT COUNT(*) FROM products;")
PRODUCT_COUNT=$(echo $PRODUCT_COUNT | xargs)  # Trim whitespace

if [ "$PRODUCT_COUNT" -eq "0" ]; then
    echo -e "${YELLOW}⚠${NC} No hay productos en la base de datos"
    echo -e "${YELLOW}El seed se ejecutará pero no insertará fotos${NC}"
    echo -e "${YELLOW}Crea productos primero o ejecuta: ./scripts/run-seeds.sh${NC}"
    echo ""
else
    echo -e "${GREEN}✓${NC} Encontrados $PRODUCT_COUNT productos"
fi

echo ""
echo -e "${YELLOW}⏳${NC} Ejecutando seed de product_photos..."

# Ejecutar el seed
SEED_FILE="scripts/seeds/05_product_photos.sql"

if [ ! -f "$SEED_FILE" ]; then
    echo -e "${RED}✗${NC} Archivo de seed no encontrado: $SEED_FILE"
    exit 1
fi

if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f "$SEED_FILE"; then
    echo -e "${GREEN}✓${NC} Seed ejecutado exitosamente"
else
    echo -e "${RED}✗${NC} Error al ejecutar el seed"
    exit 1
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Seed completado${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Mostrar estadísticas
echo -e "${BLUE}Estadísticas:${NC}"
PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c "
SELECT 
    COUNT(*) as total_photos,
    COUNT(DISTINCT product_id) as products_with_photos,
    SUM(CASE WHEN is_primary THEN 1 ELSE 0 END) as primary_photos
FROM product_photos;
"

echo ""
echo -e "${GREEN}✓ Proceso completado${NC}"
echo ""

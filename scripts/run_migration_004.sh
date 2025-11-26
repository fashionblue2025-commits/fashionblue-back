#!/bin/bash

# Script para ejecutar la migración 004: Make product_id nullable
# Fecha: 2024-11-21

set -e  # Detener en caso de error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Migración 004: Make product_id nullable${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Cargar variables de entorno desde .env si existe
if [ -f .env ]; then
    echo -e "${YELLOW}📄 Cargando variables de entorno desde .env${NC}"
    export $(cat .env | grep -v '^#' | xargs)
else
    echo -e "${YELLOW}⚠️  Archivo .env no encontrado, usando valores por defecto${NC}"
fi

# Variables de conexión a la base de datos
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-fashionblue}
DB_PASSWORD=${DB_PASSWORD:-fashionblue123}
DB_NAME=${DB_NAME:-fashionblue_db}

echo -e "${YELLOW}🔧 Configuración de base de datos:${NC}"
echo "   Host: $DB_HOST"
echo "   Port: $DB_PORT"
echo "   Database: $DB_NAME"
echo "   User: $DB_USER"
echo ""

# Verificar si PostgreSQL está disponible
echo -e "${YELLOW}🔍 Verificando conexión a PostgreSQL...${NC}"
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    echo -e "${RED}❌ Error: No se puede conectar a PostgreSQL${NC}"
    echo -e "${YELLOW}   Asegúrate de que PostgreSQL esté corriendo y las credenciales sean correctas${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Conexión exitosa${NC}"
echo ""

# Ejecutar la migración
echo -e "${YELLOW}🚀 Ejecutando migración...${NC}"
MIGRATION_FILE="scripts/migrations/004_make_product_id_nullable.sql"

if [ ! -f "$MIGRATION_FILE" ]; then
    echo -e "${RED}❌ Error: Archivo de migración no encontrado: $MIGRATION_FILE${NC}"
    exit 1
fi

echo -e "${YELLOW}   Archivo: $MIGRATION_FILE${NC}"
echo ""

# Ejecutar la migración
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $MIGRATION_FILE; then
    echo ""
    echo -e "${GREEN}✅ Migración ejecutada exitosamente${NC}"
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Cambios aplicados:${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo -e "  • product_id en order_items ahora es NULLABLE"
    echo -e "  • Permite crear órdenes sin producto existente"
    echo -e "  • El producto se crea cuando se aprueba la cotización"
    echo ""
else
    echo ""
    echo -e "${RED}❌ Error al ejecutar la migración${NC}"
    exit 1
fi

# Verificar el cambio
echo -e "${YELLOW}🔍 Verificando cambios en la tabla order_items...${NC}"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'order_items' 
  AND column_name IN ('product_id', 'product_name')
ORDER BY ordinal_position;
"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ Migración completada${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}📝 Próximos pasos:${NC}"
echo -e "  1. Reinicia la aplicación para que GORM reconozca los cambios"
echo -e "  2. Ahora puedes crear órdenes CUSTOM/INVENTORY sin ProductID"
echo -e "  3. El producto se creará automáticamente en el estado PLANNED/APPROVED"
echo ""

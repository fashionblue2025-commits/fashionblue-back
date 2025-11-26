#!/bin/bash

# Script para ejecutar migraciones de base de datos
# Uso: ./scripts/run_migrations.sh

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "Ejecutando Migraciones de Base de Datos"
echo "========================================="
echo ""

# Cargar variables de entorno
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Verificar que las variables estén configuradas
if [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ] || [ -z "$DB_USER" ]; then
    echo -e "${RED}Error: Variables de entorno no configuradas${NC}"
    echo "Asegúrate de tener un archivo .env con:"
    echo "  DB_HOST=localhost"
    echo "  DB_PORT=5432"
    echo "  DB_NAME=fashionblue_db"
    echo "  DB_USER=postgres"
    echo "  DB_PASSWORD=tu_password"
    exit 1
fi

# Construir string de conexión
DB_URL="postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

echo -e "${YELLOW}Conectando a:${NC} $DB_HOST:$DB_PORT/$DB_NAME"
echo ""

# Ejecutar cada migración
MIGRATIONS_DIR="./scripts/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}Error: Directorio de migraciones no encontrado${NC}"
    exit 1
fi

# Ordenar archivos numéricamente
for migration in $(ls $MIGRATIONS_DIR/*.sql | sort -V); do
    filename=$(basename "$migration")
    echo -e "${YELLOW}Ejecutando:${NC} $filename"
    
    # Ejecutar migración
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $filename completada"
    else
        echo -e "${RED}✗${NC} Error en $filename"
        exit 1
    fi
    echo ""
done

echo "========================================="
echo -e "${GREEN}Todas las migraciones completadas${NC}"
echo "========================================="

#!/bin/bash

###############################################################################
# Script para ejecutar migraciones en Producción (VPS)
# Ejecuta las migraciones directamente en el contenedor de PostgreSQL
###############################################################################

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║         Fashion Blue - Migraciones en Producción              ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo -e "${NC}\n"

# Cargar variables de entorno
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo -e "${RED}❌ Archivo .env no encontrado${NC}"
    exit 1
fi

# Verificar que estamos en el directorio correcto
if [ ! -d "migrations" ]; then
    echo -e "${RED}❌ Directorio 'migrations' no encontrado${NC}"
    echo -e "${YELLOW}Ejecuta este script desde la raíz del proyecto${NC}"
    exit 1
fi

CONTAINER_NAME="fashionblue-postgres-prod"
DB_USER=${DB_USER:-fashionblue_user}
DB_NAME=${DB_NAME:-fashionblue_prod}

echo -e "${YELLOW}📋 Configuración:${NC}"
echo -e "   Base de datos: ${DB_NAME}"
echo -e "   Usuario: ${DB_USER}"
echo -e "   Contenedor: ${CONTAINER_NAME}"
echo ""

# Verificar que el contenedor existe y está corriendo
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${RED}❌ Contenedor ${CONTAINER_NAME} no está corriendo${NC}"
    echo -e "${YELLOW}Ejecuta: docker compose -f docker-compose.prod.yml up -d${NC}"
    exit 1
fi

# Función para ejecutar SQL
execute_sql() {
    local sql_file=$1
    local filename=$(basename "$sql_file")
    
    echo -e "${YELLOW}📝 Ejecutando: ${filename}${NC}"
    
    if docker exec -i "${CONTAINER_NAME}" psql -U "${DB_USER}" -d "${DB_NAME}" < "$sql_file"; then
        echo -e "${GREEN}✅ ${filename} ejecutado correctamente${NC}"
    else
        echo -e "${RED}❌ Error ejecutando ${filename}${NC}"
        return 1
    fi
}

# Crear tabla de control de migraciones si no existe
echo -e "\n${YELLOW}🔧 Verificando tabla de control de migraciones...${NC}"
docker exec -i "${CONTAINER_NAME}" psql -U "${DB_USER}" -d "${DB_NAME}" <<-EOSQL
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
EOSQL

echo -e "${GREEN}✅ Tabla de control verificada${NC}"

# Listar migraciones
echo -e "\n${YELLOW}📂 Migraciones disponibles:${NC}"
migrations=($(ls -1 migrations/*.sql | sort))

if [ ${#migrations[@]} -eq 0 ]; then
    echo -e "${YELLOW}⚠️  No se encontraron archivos de migración${NC}"
    exit 0
fi

for migration in "${migrations[@]}"; do
    echo -e "   • $(basename "$migration")"
done

echo -e "\n${YELLOW}¿Deseas ejecutar todas las migraciones? (yes/no):${NC}"
read -r confirm

if [ "$confirm" != "yes" ]; then
    echo -e "${YELLOW}❌ Operación cancelada${NC}"
    exit 0
fi

# Ejecutar migraciones
echo -e "\n${BLUE}🚀 Iniciando ejecución de migraciones...${NC}\n"

success_count=0
error_count=0

for migration in "${migrations[@]}"; do
    filename=$(basename "$migration" .sql)
    
    # Verificar si ya fue aplicada
    already_applied=$(docker exec "${CONTAINER_NAME}" psql -U "${DB_USER}" -d "${DB_NAME}" -t -c \
        "SELECT COUNT(*) FROM schema_migrations WHERE version = '${filename}';" | tr -d ' ')
    
    if [ "$already_applied" -eq "1" ]; then
        echo -e "${BLUE}⏭️  ${filename} ya fue aplicada (omitiendo)${NC}"
        continue
    fi
    
    # Ejecutar migración
    if execute_sql "$migration"; then
        # Registrar migración como aplicada
        docker exec "${CONTAINER_NAME}" psql -U "${DB_USER}" -d "${DB_NAME}" -c \
            "INSERT INTO schema_migrations (version) VALUES ('${filename}');" > /dev/null
        ((success_count++))
    else
        ((error_count++))
        echo -e "${RED}❌ Fallo al ejecutar migración, deteniéndose${NC}"
        break
    fi
    
    echo ""
done

# Resumen
echo -e "\n${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    Resumen de Migraciones                      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}\n"

echo -e "${GREEN}✅ Ejecutadas exitosamente: ${success_count}${NC}"
if [ $error_count -gt 0 ]; then
    echo -e "${RED}❌ Errores: ${error_count}${NC}"
fi

# Mostrar migraciones aplicadas
echo -e "\n${YELLOW}📊 Migraciones en la base de datos:${NC}"
docker exec "${CONTAINER_NAME}" psql -U "${DB_USER}" -d "${DB_NAME}" -c \
    "SELECT version, applied_at FROM schema_migrations ORDER BY applied_at;"

echo -e "\n${GREEN}✨ Proceso completado${NC}\n"

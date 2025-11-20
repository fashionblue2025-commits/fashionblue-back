#!/bin/bash

# Script para desarrollo local - Fashion Blue
# Este script levanta solo la base de datos en Docker y ejecuta la app localmente

set -e

echo "🚀 Fashion Blue - Desarrollo Local"
echo "=================================="

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Función para limpiar al salir
cleanup() {
    echo -e "\n${YELLOW}🛑 Deteniendo servicios...${NC}"
    docker-compose down
    exit 0
}

trap cleanup SIGINT SIGTERM

# 1. Levantar solo PostgreSQL y pgAdmin
echo -e "${BLUE}📦 Levantando PostgreSQL y pgAdmin...${NC}"
docker-compose up -d postgres pgadmin

# 2. Esperar a que PostgreSQL esté listo
echo -e "${BLUE}⏳ Esperando a que PostgreSQL esté listo...${NC}"
until docker exec fashionblue-postgres pg_isready -U fashionblue > /dev/null 2>&1; do
    echo -e "${YELLOW}   Esperando PostgreSQL...${NC}"
    sleep 2
done
echo -e "${GREEN}✅ PostgreSQL está listo!${NC}"

# 3. Mostrar información
echo ""
echo -e "${GREEN}✅ Servicios levantados:${NC}"
echo -e "   📊 PostgreSQL: ${BLUE}localhost:5432${NC}"
echo -e "   🔧 pgAdmin:    ${BLUE}http://localhost:5050${NC}"
echo -e "      Email:     ${YELLOW}admin@fashionblue.com${NC}"
echo -e "      Password:  ${YELLOW}admin123${NC}"
echo ""
echo -e "${GREEN}🔧 Para conectar desde la app:${NC}"
echo -e "   Host:     ${BLUE}localhost${NC}"
echo -e "   Port:     ${BLUE}5432${NC}"
echo -e "   User:     ${BLUE}fashionblue${NC}"
echo -e "   Password: ${BLUE}fashionblue123${NC}"
echo -e "   Database: ${BLUE}fashionblue_db${NC}"
echo ""
echo -e "${GREEN}🚀 Ahora puedes ejecutar la aplicación:${NC}"
echo -e "   ${BLUE}go run cmd/api/main.go${NC}"
echo -e "   o usar el debugger de VS Code/GoLand"
echo ""
echo -e "${YELLOW}💡 Presiona Ctrl+C para detener los servicios${NC}"
echo ""

# Mantener el script corriendo
tail -f /dev/null

#!/bin/bash

# Script para ejecutar seeds SQL en la base de datos
# Uso: ./scripts/run-seeds.sh

set -e

# Colores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🌱 Fashion Blue - Database Seeding${NC}"
echo "=================================="
echo ""

# Cargar variables de entorno
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
    echo -e "${GREEN}✅ Environment variables loaded${NC}"
else
    echo -e "${YELLOW}⚠️  .env file not found, using defaults${NC}"
    DB_HOST=${DB_HOST:-localhost}
    DB_PORT=${DB_PORT:-5432}
    DB_USER=${DB_USER:-fashionblue}
    DB_PASSWORD=${DB_PASSWORD:-fashionblue123}
    DB_NAME=${DB_NAME:-fashionblue_db}
fi

echo ""
echo -e "${BLUE}📊 Database connection:${NC}"
echo "   Host: $DB_HOST"
echo "   Port: $DB_PORT"
echo "   User: $DB_USER"
echo "   Database: $DB_NAME"
echo ""

# Verificar si PostgreSQL está disponible
echo -e "${BLUE}🔍 Checking PostgreSQL connection...${NC}"
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    echo -e "${RED}❌ Error: Cannot connect to PostgreSQL${NC}"
    echo "   Make sure PostgreSQL is running and credentials are correct"
    exit 1
fi
echo -e "${GREEN}✅ PostgreSQL connection successful${NC}"
echo ""

# Ejecutar seeds
echo -e "${BLUE}🌱 Running seeds...${NC}"
echo "=================================="

# 1. Users
echo ""
echo -e "${YELLOW}👤 Creating Super Admin...${NC}"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/seeds/01_users.sql

# 2. Categories
echo ""
echo -e "${YELLOW}📁 Creating Categories...${NC}"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/seeds/02_categories.sql

# 3. Sizes
echo ""
echo -e "${YELLOW}📏 Creating Sizes...${NC}"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/seeds/03_sizes.sql

echo ""
echo "=================================="
echo -e "${GREEN}✅ DATABASE SEEDED SUCCESSFULLY!${NC}"
echo "=================================="
echo ""
echo -e "${BLUE}📊 Summary:${NC}"
echo "   👤 Users: 1 (Super Admin)"
echo "   📁 Categories: 5"
echo "   📏 Sizes: 33 total"
echo "      - Shirts: 6"
echo "      - Pants: 10"
echo "      - Shoes: 17"
echo ""
echo -e "${GREEN}🔐 Login credentials:${NC}"
echo "   Email: ${YELLOW}admin@fashionblue.com${NC}"
echo "   Password: ${YELLOW}admin123${NC}"
echo ""

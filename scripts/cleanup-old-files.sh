#!/bin/bash

# Script para limpiar archivos antiguos después de la refactorización

set -e

echo "🧹 Limpieza de Archivos Antiguos - Fashion Blue"
echo "==============================================="
echo ""
echo "⚠️  Este script eliminará archivos que ya no se usan después de la refactorización"
echo ""

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para confirmar
confirm() {
    read -p "$(echo -e ${YELLOW}$1${NC}) (escribe 'SI' para continuar): " response
    if [ "$response" != "SI" ]; then
        echo -e "${RED}❌ Operación cancelada${NC}"
        exit 0
    fi
}

echo -e "${BLUE}📋 Archivos que se eliminarán:${NC}"
echo ""

# 1. Handlers antiguos (en la raíz de handlers/)
echo -e "${YELLOW}Handlers antiguos:${NC}"
OLD_HANDLERS=(
    "internal/adapters/http/handlers/auth_handler.go"
    "internal/adapters/http/handlers/user_handler.go"
    "internal/adapters/http/handlers/capital_injection_handler.go"
    "internal/adapters/http/handlers/category_handler.go"
    "internal/adapters/http/handlers/product_handler.go"
    "internal/adapters/http/handlers/customer_handler.go"
    "internal/adapters/http/handlers/sale_handler.go"
    "internal/adapters/http/handlers/supplier_handler.go"
    "internal/adapters/http/handlers/purchase_handler.go"
)

for file in "${OLD_HANDLERS[@]}"; do
    if [ -f "$file" ]; then
        echo "  - $file"
    fi
done

# 2. Carpeta de servicios antiguos
echo ""
echo -e "${YELLOW}Servicios antiguos:${NC}"
if [ -d "internal/application/services" ]; then
    echo "  - internal/application/services/ (carpeta completa)"
fi

# 3. Repositorios antiguos (postgres/)
echo ""
echo -e "${YELLOW}Repositorios antiguos:${NC}"
if [ -d "internal/adapters/postgres" ]; then
    echo "  - internal/adapters/postgres/ (carpeta completa)"
fi

# 4. Ports antiguos
echo ""
echo -e "${YELLOW}Interfaces antiguas:${NC}"
OLD_PORTS=(
    "internal/ports/repositories.go"
    "internal/ports/services.go"
)

for file in "${OLD_PORTS[@]}"; do
    if [ -f "$file" ]; then
        echo "  - $file"
    fi
done

# 5. Domain antiguos (con GORM)
echo ""
echo -e "${YELLOW}Entidades de dominio antiguas (con GORM):${NC}"
OLD_DOMAIN=(
    "internal/domain/user.go"
    "internal/domain/capital_injection.go"
    "internal/domain/category.go"
    "internal/domain/product.go"
    "internal/domain/customer.go"
    "internal/domain/sale.go"
    "internal/domain/supplier.go"
    "internal/domain/purchase.go"
)

for file in "${OLD_DOMAIN[@]}"; do
    if [ -f "$file" ]; then
        echo "  - $file"
    fi
done

echo ""
echo -e "${BLUE}📁 Archivos que se mantienen (nuevos):${NC}"
echo "  ✅ internal/domain/entities/*.go"
echo "  ✅ internal/domain/ports/*.go"
echo "  ✅ internal/adapters/persistence/models/*.go"
echo "  ✅ internal/adapters/persistence/repositories/*/*.go"
echo "  ✅ internal/application/usecases/*/*.go"
echo "  ✅ internal/adapters/http/handlers/*/*.go"
echo ""

# Confirmar
confirm "¿Deseas continuar con la eliminación?"

echo ""
echo -e "${BLUE}🗑️  Eliminando archivos...${NC}"

# Crear backup antes de eliminar
BACKUP_DIR="backup_old_files_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo -e "${YELLOW}💾 Creando backup en: $BACKUP_DIR${NC}"

# Función para mover a backup y eliminar
backup_and_remove() {
    local file=$1
    if [ -f "$file" ]; then
        local backup_path="$BACKUP_DIR/$file"
        mkdir -p "$(dirname "$backup_path")"
        cp "$file" "$backup_path"
        rm "$file"
        echo -e "  ${GREEN}✓${NC} $file"
    fi
}

# Función para mover carpeta a backup y eliminar
backup_and_remove_dir() {
    local dir=$1
    if [ -d "$dir" ]; then
        local backup_path="$BACKUP_DIR/$dir"
        mkdir -p "$(dirname "$backup_path")"
        cp -r "$dir" "$backup_path"
        rm -rf "$dir"
        echo -e "  ${GREEN}✓${NC} $dir/"
    fi
}

# Eliminar handlers antiguos
echo ""
echo -e "${YELLOW}Eliminando handlers antiguos...${NC}"
for file in "${OLD_HANDLERS[@]}"; do
    backup_and_remove "$file"
done

# Eliminar servicios antiguos
echo ""
echo -e "${YELLOW}Eliminando servicios antiguos...${NC}"
backup_and_remove_dir "internal/application/services"

# Eliminar repositorios antiguos
echo ""
echo -e "${YELLOW}Eliminando repositorios antiguos...${NC}"
backup_and_remove_dir "internal/adapters/postgres"

# Eliminar ports antiguos
echo ""
echo -e "${YELLOW}Eliminando interfaces antiguas...${NC}"
for file in "${OLD_PORTS[@]}"; do
    backup_and_remove "$file"
done

# Eliminar carpeta ports si está vacía
if [ -d "internal/ports" ]; then
    if [ -z "$(ls -A internal/ports)" ]; then
        rmdir "internal/ports"
        echo -e "  ${GREEN}✓${NC} internal/ports/ (carpeta vacía)"
    fi
fi

# Eliminar domain antiguos
echo ""
echo -e "${YELLOW}Eliminando entidades de dominio antiguas...${NC}"
for file in "${OLD_DOMAIN[@]}"; do
    backup_and_remove "$file"
done

echo ""
echo -e "${GREEN}✅ Limpieza completada!${NC}"
echo ""
echo -e "${BLUE}📊 Resumen:${NC}"
echo "  - Backup creado en: ${YELLOW}$BACKUP_DIR${NC}"
echo "  - Archivos eliminados: ${GREEN}$(find "$BACKUP_DIR" -type f | wc -l)${NC}"
echo ""
echo -e "${YELLOW}💡 Si algo sale mal, puedes restaurar desde el backup:${NC}"
echo "  cp -r $BACKUP_DIR/internal/* internal/"
echo ""
echo -e "${GREEN}🎉 Tu proyecto ahora solo tiene la nueva arquitectura limpia!${NC}"

#!/bin/bash

# Script de despliegue automatizado para Fashion Blue API
# Uso: ./scripts/deploy.sh

set -e  # Salir si hay error

echo "🚀 Iniciando despliegue de Fashion Blue API..."

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verificar que estamos en el directorio correcto
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: Debes ejecutar este script desde la raíz del proyecto${NC}"
    exit 1
fi

# 1. Compilar la aplicación
echo -e "${YELLOW}📦 Compilando aplicación...${NC}"
go build -o fashion-blue-api ./cmd/api
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Compilación exitosa${NC}"
else
    echo -e "${RED}❌ Error en compilación${NC}"
    exit 1
fi

# 2. Ejecutar tests (opcional)
echo -e "${YELLOW}🧪 Ejecutando tests...${NC}"
go test ./... -v
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Tests pasados${NC}"
else
    echo -e "${RED}⚠️  Algunos tests fallaron - ¿Continuar? (y/n)${NC}"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 3. Verificar variables de entorno
echo -e "${YELLOW}🔍 Verificando configuración...${NC}"
if [ ! -f ".env" ]; then
    echo -e "${RED}❌ Error: Archivo .env no encontrado${NC}"
    echo "Crea un archivo .env basado en .env.example"
    exit 1
fi

# Verificar variables críticas
required_vars=("DB_HOST" "DB_USER" "DB_PASSWORD" "DB_NAME" "JWT_SECRET")
for var in "${required_vars[@]}"; do
    if ! grep -q "^${var}=" .env; then
        echo -e "${RED}❌ Variable requerida ${var} no encontrada en .env${NC}"
        exit 1
    fi
done
echo -e "${GREEN}✅ Configuración verificada${NC}"

# 4. Crear backup del binario actual (si existe)
if [ -f "fashion-blue-api.backup" ]; then
    echo -e "${YELLOW}📦 Creando backup del binario anterior...${NC}"
    mv fashion-blue-api.backup fashion-blue-api.backup.old
    mv fashion-blue-api fashion-blue-api.backup
    echo -e "${GREEN}✅ Backup creado${NC}"
fi

# 5. Verificar conexión a base de datos
echo -e "${YELLOW}🔌 Verificando conexión a base de datos...${NC}"
source .env
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Conexión a BD exitosa${NC}"
else
    echo -e "${RED}❌ Error: No se pudo conectar a la base de datos${NC}"
    exit 1
fi

# 6. Ejecutar migraciones pendientes
echo -e "${YELLOW}📊 ¿Ejecutar migraciones? (y/n)${NC}"
read -r run_migrations
if [[ "$run_migrations" =~ ^[Yy]$ ]]; then
    echo "Ejecutando migraciones..."
    for migration in migrations/*.sql; do
        echo "  - Aplicando: $(basename $migration)"
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$migration" 2>&1 | grep -v "already exists" | grep -v "skipping"
    done
    echo -e "${GREEN}✅ Migraciones completadas${NC}"
fi

# 7. Reiniciar servicio (si está en producción)
echo -e "${YELLOW}🔄 ¿Reiniciar servicio systemd? (y/n)${NC}"
read -r restart_service
if [[ "$restart_service" =~ ^[Yy]$ ]]; then
    echo "Reiniciando fashion-blue service..."
    sudo systemctl restart fashion-blue
    sleep 3
    sudo systemctl status fashion-blue --no-pager
    echo -e "${GREEN}✅ Servicio reiniciado${NC}"
fi

# 8. Health check
echo -e "${YELLOW}🏥 Ejecutando health check...${NC}"
sleep 2
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ "$response" = "200" ]; then
    echo -e "${GREEN}✅ API funcionando correctamente${NC}"
else
    echo -e "${RED}❌ Error: API no responde correctamente (HTTP $response)${NC}"
    echo "Ver logs con: sudo journalctl -u fashion-blue -n 50"
    exit 1
fi

# 9. Resumen
echo ""
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ Despliegue completado exitosamente!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo ""
echo "📊 Información del despliegue:"
echo "  - Binario: fashion-blue-api"
echo "  - Backup: fashion-blue-api.backup"
echo "  - Health: http://localhost:8080/health"
echo ""
echo "📝 Comandos útiles:"
echo "  - Ver logs: sudo journalctl -u fashion-blue -f"
echo "  - Estado: sudo systemctl status fashion-blue"
echo "  - Reiniciar: sudo systemctl restart fashion-blue"
echo ""
echo -e "${GREEN}🎉 ¡Todo listo!${NC}"

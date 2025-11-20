#!/bin/bash

# Script para desplegar en ambiente de producción

set -e

echo "🚀 Desplegando Fashion Blue - PRODUCCIÓN"
echo "========================================="
echo ""
echo "⚠️  ADVERTENCIA: Estás a punto de desplegar en PRODUCCIÓN"
echo ""
read -p "¿Estás seguro? (escribe 'SI' para continuar): " confirm

if [ "$confirm" != "SI" ]; then
    echo "❌ Despliegue cancelado"
    exit 0
fi

# Verificar que existe el archivo .env.production
if [ ! -f .env.production ]; then
    echo "❌ Error: .env.production no existe"
    echo "📝 Copia .env.production.example a .env.production y configura los valores"
    exit 1
fi

# Verificar que las credenciales no sean las de ejemplo
if grep -q "CHANGE_ME" .env.production; then
    echo "❌ Error: .env.production contiene valores de ejemplo (CHANGE_ME)"
    echo "📝 Por favor configura todas las credenciales antes de desplegar"
    exit 1
fi

# Cargar variables de entorno
export $(cat .env.production | grep -v '^#' | xargs)

echo "✅ Variables de entorno cargadas desde .env.production"

# Backup de la base de datos (si existe)
if docker ps | grep -q fashionblue-postgres; then
    echo "💾 Creando backup de la base de datos..."
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    docker exec fashionblue-postgres pg_dump -U ${DB_USER} ${DB_NAME} > "backups/${BACKUP_FILE}"
    echo "✅ Backup creado: backups/${BACKUP_FILE}"
fi

# Levantar servicios
echo "🐳 Desplegando servicios de producción..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Esperar a que los servicios estén listos
echo "⏳ Esperando a que los servicios estén listos..."
sleep 10

# Verificar health check
echo "🔍 Verificando health check..."
if curl -f http://localhost:${APP_PORT:-8080}/health > /dev/null 2>&1; then
    echo "✅ API está respondiendo correctamente"
else
    echo "❌ Error: API no está respondiendo"
    echo "📝 Ver logs: docker-compose logs api"
    exit 1
fi

echo ""
echo "✅ Despliegue completado exitosamente!"
echo ""
echo "📊 Servicios disponibles:"
echo "   API:       http://localhost:${APP_PORT:-8080}"
echo "   Health:    http://localhost:${APP_PORT:-8080}/health"
echo ""
echo "📝 Ver logs: docker-compose logs -f api"
echo "🛑 Detener:  docker-compose -f docker-compose.yml -f docker-compose.prod.yml down"
echo ""
echo "⚠️  Recuerda:"
echo "   - Monitorear los logs regularmente"
echo "   - Configurar backups automáticos"
echo "   - Revisar métricas de rendimiento"

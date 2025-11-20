#!/bin/bash

# Script para desplegar en ambiente de desarrollo

set -e

echo "🚀 Desplegando Fashion Blue - DESARROLLO"
echo "========================================"

# Verificar que existe el archivo .env.development
if [ ! -f .env.development ]; then
    echo "❌ Error: .env.development no existe"
    echo "📝 Copia .env.development.example a .env.development y configura los valores"
    exit 1
fi

# Cargar variables de entorno
export $(cat .env.development | grep -v '^#' | xargs)

echo "✅ Variables de entorno cargadas desde .env.development"

# Levantar servicios
echo "🐳 Levantando servicios de desarrollo..."
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

echo ""
echo "✅ Servicios levantados exitosamente!"
echo ""
echo "📊 Servicios disponibles:"
echo "   API:       http://localhost:${APP_PORT:-8080}"
echo "   Health:    http://localhost:${APP_PORT:-8080}/health"
echo "   PostgreSQL: localhost:${DB_PORT:-5432}"
echo "   pgAdmin:   http://localhost:${PGADMIN_PORT:-5050}"
echo ""
echo "📝 Ver logs: docker-compose logs -f"
echo "🛑 Detener:  docker-compose down"

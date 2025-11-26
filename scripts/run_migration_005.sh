#!/bin/bash

# ============================================================================
# Script para ejecutar la migración 005: Refactorizar Productos a Variantes
# Uso: ./scripts/run_migration_005.sh
# ============================================================================

set -e

echo "=============================================="
echo "  Migración 005: Product Base + Variants"
echo "=============================================="
echo ""

# Cargar variables de entorno
if [ -f .env ]; then
    echo "📄 Cargando variables de entorno desde .env"
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  Archivo .env no encontrado, usando variables de entorno del sistema"
fi

# Configuración de la base de datos
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-fashionblue_db}
DB_USER=${DB_USER:-fashionblue}
DB_PASSWORD=${DB_PASSWORD}

echo "🔧 Configuración de base de datos:"
echo "   Host: $DB_HOST"
echo "   Port: $DB_PORT"
echo "   Database: $DB_NAME"
echo "   User: $DB_USER"
echo ""

# Verificar que PostgreSQL esté disponible
echo "🔍 Verificando conexión a PostgreSQL..."
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    echo "❌ Error: No se puede conectar a PostgreSQL"
    echo "   Asegúrate de que PostgreSQL esté corriendo y las credenciales sean correctas"
    exit 1
fi

echo "✅ Conexión exitosa a PostgreSQL"
echo ""

# Advertencia
echo "⚠️  ADVERTENCIA: Esta migración realizará cambios importantes:"
echo "   1. Creará tabla product_variants"
echo "   2. Migrará datos de products a product_variants"
echo "   3. Eliminará columnas de products (color, size_id, stock, reserved_stock)"
echo "   4. Actualizará order_items para usar product_variant_id"
echo ""
echo "   Se recomienda hacer un backup de la base de datos antes de continuar."
echo ""

# Confirmar
read -p "¿Deseas continuar? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Migración cancelada"
    exit 1
fi

echo ""
echo "🚀 Ejecutando migración 005..."
echo ""

# Ejecutar la migración
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/migrations/005_refactor_product_variants.sql

if [ $? -eq 0 ]; then
    echo ""
    echo "=============================================="
    echo "✅ Migración 005 completada exitosamente"
    echo "=============================================="
    echo ""
    echo "📊 Cambios aplicados:"
    echo "   ✓ Tabla 'product_variants' creada"
    echo "   ✓ Datos migrados de 'products' a 'product_variants'"
    echo "   ✓ Columna 'product_variant_id' agregada a 'order_items'"
    echo "   ✓ Columnas eliminadas de 'products': color, size_id, stock, reserved_stock"
    echo "   ✓ Vista 'products_with_variants' creada"
    echo ""
    echo "🎯 Próximos pasos:"
    echo "   1. Reiniciar la aplicación"
    echo "   2. Verificar que las órdenes existentes funcionen correctamente"
    echo "   3. Probar creación de nuevas órdenes con variantes"
    echo ""
else
    echo ""
    echo "❌ Error al ejecutar la migración"
    echo "   Revisa los logs arriba para más detalles"
    exit 1
fi

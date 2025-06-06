#!/bin/bash
# migrate-prod.sh - Ejecutar migraciones en producción

set -e

echo "🗄️  Database migrations for production..."

# Verificar que estamos conectados a Clever Cloud
if ! clever status > /dev/null 2>&1; then
    echo "❌ Error: Not connected to Clever Cloud app"
    echo "Run 'clever link' to connect to your app"
    exit 1
fi

# Verificar que las variables de DB existen
echo "🔍 Checking database configuration..."
if ! clever env | grep -q "MYSQL_ADDON_HOST"; then
    echo "❌ Error: MySQL addon not found or not linked"
    echo "Run: clever addon link gocms-mysql"
    exit 1
fi

# Usar el Makefile que ya tienes configurado
echo "📊 Current migration status:"
make migrate-status-prod

echo ""
read -p "Apply pending migrations? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "⬆️  Applying migrations..."
    make migrate-prod
    
    echo ""
    echo "📊 Final migration status:"
    make migrate-status-prod
    
    echo ""
    echo "✅ Migrations completed successfully!"
else
    echo "ℹ️  No migrations applied."
fi
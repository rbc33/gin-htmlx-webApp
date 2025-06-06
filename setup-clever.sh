#!/bin/bash
# clever-startup.sh - Setup inicial de la aplicación en Clever Cloud

set -e

echo "🔧 Setting up gocms in Clever Cloud..."

# 1. Crear aplicación
echo "📱 Creating application..."
clever create --type go gocms-app

# 2. Crear y enlazar base de datos MySQL
echo "🗄️  Creating MySQL database..."
clever addon create mysql-addon --name gocms-mysql --plan dev
clever addon link gocms-mysql

# 3. Configurar variables de entorno
echo "⚙️  Setting environment variables..."
clever env set ENVIRONMENT production
clever env set MEDIA_DIR media

# 4. Generar CSS y hacer primer deploy
echo "🎨 Generating CSS..."
make css

echo "🚀 Initial deployment..."
clever deploy

echo ""
echo "✅ Initial setup completed!"
echo "🌍 App URL: $(clever domain)"
echo ""
echo "Next steps:"
echo "  1. Run: ./migrate-prod.sh (to apply database migrations)"
echo "  2. Check: clever logs (to see application logs)"
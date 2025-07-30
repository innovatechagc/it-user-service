#!/bin/bash

# Script para ejecutar el servicio localmente
echo "🚀 Iniciando it-user-service en modo local..."

# Cargar variables de entorno (filtrando comentarios)
set -a
source .env
set +a

# Verificar que las variables estén cargadas
echo "📊 Configuración:"
echo "  - DB_HOST: $DB_HOST"
echo "  - DB_PORT: $DB_PORT"
echo "  - DB_NAME: $DB_NAME"
echo "  - PORT: $PORT"
echo "  - ENVIRONMENT: $ENVIRONMENT"
echo ""

# Crear directorio bin si no existe
mkdir -p bin

# Compilar y ejecutar
echo "🔨 Compilando aplicación..."
go build -o bin/it-user-service ./cmd

if [ $? -eq 0 ]; then
    echo "✅ Compilación exitosa"
    echo "🌐 Iniciando servidor en http://localhost:$PORT"
    echo "📋 Health check: http://localhost:$PORT/api/v1/health"
    echo ""
    ./bin/it-user-service
else
    echo "❌ Error en la compilación"
    exit 1
fi
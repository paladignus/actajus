#!/bin/bash

# Script para subir os serviços de observabilidade e aplicação

echo "🚀 Iniciando os serviços de observabilidade e aplicação..."

# Verifica se o Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker não está rodando. Por favor inicie o Docker e tente novamente."
    exit 1
fi

# Verifica se o docker-compose está disponível
if ! command -v docker-compose &> /dev/null; then
    echo "⚠️  docker-compose não encontrado. Tentando usar docker compose (v2)..."
    if ! command -v docker &> /dev/null || ! docker compose version > /dev/null 2>&1; then
        echo "❌ Nenhum comando docker-compose encontrado."
        exit 1
    fi
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

echo "✅ Usando: $COMPOSE_CMD"

# Subir os serviços em modo detached
echo "🐳 Subindo todos os serviços..."
$COMPOSE_CMD up -d

echo ""
echo "📊 Serviços disponíveis:"
echo "   - Aplicação:      http://localhost:8080"
echo "   - Jaeger UI:      http://localhost:16686"
echo "   - Prometheus:     http://localhost:9090"
echo "   - Grafana:        http://localhost:3000 (admin/admin123)"
echo "   - Loki:           http://localhost:3100"
echo "   - PostgreSQL:     localhost:5433 (interno: 5432)"
echo "   - pgAdmin4:       http://localhost:5050 (admin@actajus.com.br/admin123)"
echo "   - NATS:           localhost:4224 (interno: 4222)"
echo ""
echo "📈 Para ver os logs: $COMPOSE_CMD logs -f"
echo "📈 Para parar os serviços: $COMPOSE_CMD down"
echo ""
echo "✅ Todos os serviços foram iniciados com sucesso!"
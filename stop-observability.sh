#!/bin/bash

# Script para parar os serviços de observabilidade

echo "🛑 Parando os serviços de observabilidade..."

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

# Parar e remover os serviços
$COMPOSE_CMD down

echo "✅ Todos os serviços foram parados e removidos!"
echo "💡 Para remover volumes e dados persistentes, use: $COMPOSE_CMD down -v"
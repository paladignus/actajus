#!/bin/bash

# Script para validar os arquivos YAML do Kubernetes

echo "🔍 Validando arquivos YAML do Kubernetes..."

# Encontrar todos os arquivos YAML
find k8s -name "*.yaml" -type f | while read file; do
    echo "Validando: $file"
    if python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null; then
        echo "  ✅ Válido"
    else
        echo "  ❌ Inválido"
        echo "  Conteúdo do erro:"
        python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>&1 | head -5
        echo ""
    fi
done

echo ""
echo "✅ Validação concluída!"
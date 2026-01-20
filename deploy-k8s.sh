#!/bin/bash

# Script para deploy do ActaJus no Kubernetes
# Este script aplica todos os recursos necessários para rodar a aplicação com observabilidade completa

set -e  # Sai se algum comando falhar

echo "🚀 Iniciando deploy do ActaJus no Kubernetes..."

# Verifica se o kubectl está disponível
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl não encontrado. Por favor instale o kubectl e configure o acesso ao cluster."
    exit 1
fi

# Verifica se estamos conectados a um cluster
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Não foi possível conectar ao cluster Kubernetes."
    exit 1
fi

echo "✅ kubectl está configurado e conectado ao cluster"

# Aplica os namespaces primeiro
echo "📁 Aplicando namespaces..."
kubectl apply -f k8s/00-namespaces.yaml

echo "⏰ Aguardando criação dos namespaces..."
sleep 5

# Aplica os componentes de observabilidade
echo "🔍 Aplicando componentes de observabilidade..."

# Aplica o Jaeger
kubectl apply -f k8s/observability/01-namespace.yaml
kubectl apply -f k8s/observability/02-jaeger.yaml

# Aguarda o Jaeger estar pronto
echo "⏰ Aguardando Jaeger estar pronto..."
kubectl wait --for=condition=ready pod -l app=jaeger -n observability --timeout=120s

# Aplica os componentes de monitoramento
echo "📊 Aplicando componentes de monitoramento..."

# Aplica Prometheus, Loki e Grafana
kubectl apply -f k8s/monitoring/01-namespace.yaml
kubectl apply -f k8s/monitoring/02-prometheus.yaml
kubectl apply -f k8s/monitoring/03-loki.yaml
kubectl apply -f k8s/monitoring/04-grafana.yaml

# Aguarda os componentes de monitoramento estarem prontos
echo "⏰ Aguardando Prometheus, Loki e Grafana estarem prontos..."
kubectl wait --for=condition=ready pod -l app=prometheus -n monitoring --timeout=120s
kubectl wait --for=condition=ready pod -l app=loki -n monitoring --timeout=120s
kubectl wait --for=condition=ready pod -l app=grafana -n monitoring --timeout=120s

# Aplica os componentes da aplicação
echo "📦 Aplicando componentes da aplicação..."

# Cria o ConfigMap com os scripts de inicialização do PostgreSQL
if [ -f "createdb.sql" ]; then
    kubectl create configmap postgres-initdb \
        --from-file=1-createdb.sql=createdb.sql \
        --namespace=actajus \
        --dry-run=client -o yaml | kubectl apply -f -
else
    echo "⚠️  Arquivo createdb.sql não encontrado. Criando ConfigMap vazio..."
    kubectl create configmap postgres-initdb \
        --from-literal=1-createdb.sql="# Script de criação do banco de dados" \
        --namespace=actajus \
        --dry-run=client -o yaml | kubectl apply -f -
fi

if [ -f "inserts.sql" ]; then
    kubectl create configmap postgres-inserts \
        --from-file=2-inserts.sql=inserts.sql \
        --namespace=actajus \
        --dry-run=client -o yaml | kubectl apply -f -
else
    echo "⚠️  Arquivo inserts.sql não encontrado. Criando ConfigMap vazio..."
    kubectl create configmap postgres-inserts \
        --from-literal=2-inserts.sql="# Scripts de inserção de dados" \
        --namespace=actajus \
        --dry-run=client -o yaml | kubectl apply -f -
fi

# Ajusta o deployment do PostgreSQL para usar os ConfigMaps corretos
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-initdb
  namespace: actajus
data:
  1-createdb.sql: |
$(sed 's/^/    /' createdb.sql 2>/dev/null || echo "    -- Script de criação do banco de dados")
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-inserts
  namespace: actajus
data:
  2-inserts.sql: |
$(sed 's/^/    /' inserts.sql 2>/dev/null || echo "    -- Scripts de inserção de dados")
EOF

# Aplica os componentes da aplicação
kubectl apply -f k8s/app/01-namespace.yaml
kubectl apply -f k8s/app/02-postgres.yaml
kubectl apply -f k8s/app/03-nats.yaml
kubectl apply -f k8s/app/04-actajus-app.yaml

# Aguarda os componentes da aplicação estarem prontos
echo "⏰ Aguardando PostgreSQL estar pronto..."
kubectl wait --for=condition=ready pod -l app=postgres -n actajus --timeout=180s

echo "⏰ Aguardando NATS estar pronto..."
kubectl wait --for=condition=ready pod -l app=nats -n actajus --timeout=120s

echo "⏰ Aguardando ActaJus App estar pronto..."
kubectl wait --for=condition=ready pod -l app=actajus-app -n actajus --timeout=180s

# Exibe informações sobre os serviços
echo ""
echo "✅ Deploy concluído com sucesso!"
echo ""
echo "📍 Serviços disponíveis:"
echo "   - ActaJus API:           $(kubectl get svc actajus-app -n actajus -o jsonpath='{.status.loadBalancer.ingress[0].ip}:{.spec.ports[0].port}' 2>/dev/null || echo 'service IP:8080')"
echo "   - Jaeger UI:             $(kubectl get svc jaeger -n observability -o jsonpath='{.spec.clusterIP}:{.spec.ports[0].port}')"
echo "   - Prometheus:            $(kubectl get svc prometheus -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].ip}:{.spec.ports[0].port}' 2>/dev/null || echo 'service IP:9090')"
echo "   - Grafana:               $(kubectl get svc grafana -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].ip}:{.spec.ports[0].port}' 2>/dev/null || echo 'service IP:3000')"
echo "   - Loki:                  $(kubectl get svc loki -n monitoring -o jsonpath='{.spec.clusterIP}:{.spec.ports[0].port}')"
echo "   - PostgreSQL:            $(kubectl get svc postgres -n actajus -o jsonpath='{.spec.clusterIP}:{.spec.ports[0].port}')"
echo "   - NATS:                  $(kubectl get svc nats -n actajus -o jsonpath='{.spec.clusterIP}:{.spec.ports[0].port}')"
echo ""
echo "📊 Para verificar os logs da aplicação:"
echo "   kubectl logs -f deployment/actajus-app -n actajus"
echo ""
echo "🔍 Para ver todos os recursos criados:"
echo "   kubectl get all -n actajus"
echo "   kubectl get all -n observability"
echo "   kubectl get all -n monitoring"
echo ""

echo "💡 Para remover todos os recursos:"
echo "   kubectl delete -f k8s/00-namespaces.yaml"
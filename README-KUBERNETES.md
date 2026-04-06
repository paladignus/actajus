# ActaJus - Deploy no Kubernetes

Este diretório contém todos os arquivos necessários para fazer o deploy da aplicação ActaJus no Kubernetes com sistema de observabilidade completo.

## 🏗️ Arquitetura no Kubernetes

A aplicação é dividida em 3 namespaces:

- `actajus`: Componentes da aplicação principal (API, PostgreSQL, NATS)
- `observability`: Jaeger para tracing distribuído
- `monitoring`: Prometheus, Loki e Grafana para monitoramento

## 📁 Estrutura de diretórios

```
k8s/
├── 00-namespaces.yaml              # Criação dos namespaces
├── app/
│   ├── 01-namespace.yaml          # Namespace actajus
│   ├── 02-postgres.yaml           # PostgreSQL com PVC
│   ├── 03-nats.yaml               # NATS messaging
│   └── 04-actajus-app.yaml        # Deployment e Service do app
├── observability/
│   ├── 01-namespace.yaml          # Namespace observability
│   └── 02-jaeger.yaml             # Jaeger all-in-one
└── monitoring/
    ├── 01-namespace.yaml          # Namespace monitoring
    ├── 02-prometheus.yaml         # Prometheus com ServiceMonitor
    ├── 03-loki.yaml               # Loki para logs
    └── 04-grafana.yaml            # Grafana com datasources provisionados
├── app/                          # Componentes da aplicação (incluindo pgAdmin4)
    ├── 01-namespace.yaml
    ├── 02-postgres.yaml
    ├── 03-nats.yaml
    ├── 04-actajus-app.yaml
    ├── 05-pgadmin.yaml
    └── 06-pgadmin-config.yaml
```

## 🚀 Como fazer o deploy

### 1. Pré-requisitos

- Kubernetes cluster (Minikube, Kind, EKS, GKE, AKS, etc.)
- `kubectl` configurado para acessar o cluster
- Docker (para build da imagem da aplicação)

### 2. Build da imagem da aplicação

Primeiro, você precisa construir e enviar a imagem da aplicação para um registry:

```bash
# Build da imagem
docker build -t actajus:latest .

# (Opcional) Se estiver usando Minikube
eval $(minikube docker-env)
docker build -t actajus:latest .

# Ou envie para um registry público/privado
docker tag actajus:latest seu-registry/actajus:latest
docker push seu-registry/actajus:latest
```

### 3. Atualizar a imagem no deployment

Edite o arquivo `k8s/app/04-actajus-app.yaml` e atualize a imagem:

```yaml
containers:
- name: app
  image: seu-registry/actajus:latest  # Atualize esta linha
```

### 4. Executar o deploy

```bash
./deploy-k8s.sh
```

### 5. Acessar os serviços

Após o deploy, você pode acessar os serviços:

| Serviço | Comando para obter URL | Padrão |
|---------|------------------------|--------|
| ActaJus API | `kubectl get svc actajus-app -n actajus` | LoadBalancer |
| Jaeger UI | `kubectl get svc jaeger -n observability` | ClusterIP:16686 |
| Prometheus | `kubectl get svc prometheus -n monitoring` | LoadBalancer |
| Grafana | `kubectl get svc grafana -n monitoring` | LoadBalancer |
| Loki | `kubectl get svc loki -n monitoring` | ClusterIP:3100 |
| pgAdmin4 | `kubectl get svc pgadmin -n actajus` | LoadBalancer |

## 🔧 Configurações avançadas

### Configuração do PostgreSQL

O PostgreSQL é configurado com:

- Volume persistente de 1Gi
- Scripts de inicialização via ConfigMaps
- Configuração de segredos para credenciais

### Configuração da aplicação ActaJus

A aplicação está configurada com:

- 2 réplicas para alta disponibilidade
- Configuração de tracing para Jaeger
- Configuração de métricas para Prometheus
- Health checks
- Limites e requisições de recursos
- Configuração para coleta automática de métricas pelo Prometheus

### Observabilidade

#### Jaeger

- Service discovery automático no namespace `actajus`
- Configurado para receber traces via agente
- UI acessível via LoadBalancer ou port-forward

#### Prometheus

- Service discovery automático para pods com anotações Prometheus
- Coleta de métricas HTTP da aplicação
- Configuração para métricas Go e HTTP

#### Loki e Grafana

- Loki configurado para armazenamento local
- Grafana com datasources provisionados (Prometheus e Loki)
- Credenciais padrão: `admin / admin123`

## 🔍 Troubleshooting

### Verificar status dos pods

```bash
kubectl get pods -n actajus
kubectl get pods -n observability  
kubectl get pods -n monitoring
```

### Verificar logs

```bash
# Logs da aplicação
kubectl logs -f deployment/actajus-app -n actajus

# Logs do Prometheus
kubectl logs -f deployment/prometheus -n monitoring

# Logs do Jaeger
kubectl logs -f deployment/jaeger -n observability
```

### Acessar serviços localmente

```bash
# Port forward para Grafana
kubectl port-forward svc/grafana -n monitoring 3000:3000

# Port forward para Jaeger
kubectl port-forward svc/jaeger -n observability 16686:16686

# Port forward para ActaJus API
kubectl port-forward svc/actajus-app -n actajus 8080:8080
```

### Atualizar a aplicação

```bash
# Após fazer alterações e build da nova imagem
kubectl set image deployment/actajus-app app=seu-novo-registry/actajus:nova-versao -n actajus
```

## 📊 Boas práticas implementadas

- **Namespaces**: Separação lógica dos componentes
- **Service Accounts**: Para melhor gerenciamento de permissões
- **Persistent Volumes**: Para dados persistentes
- **Health Checks**: Liveness e readiness probes
- **Resource Limits**: Para melhor gerenciamento de recursos
- **Secrets**: Para armazenamento seguro de credenciais
- **ConfigMaps**: Para configurações externalizadas
- **Service Discovery**: Automático para Prometheus e tracing

## 🚢 Deploy em produção

Para usar em ambiente de produção, considere:

1. **Segurança**:
   - Usar certificados TLS
   - Configurar Network Policies
   - Usar RBAC apropriado
   - Secrets externos (HashiCorp Vault, AWS Secrets Manager, etc.)

2. **Observabilidade**:
   - Armazenamento permanente para Loki e Prometheus
   - Configuração de alertas
   - Retenção de dados adequada

3. **Escalabilidade**:
   - Horizontal Pod Autoscaler (HPA)
   - Configurar PDBs (Pod Disruption Budgets)
   - Clustering para PostgreSQL (usando operators como Zalando Spilo ou Crunchy Data)

4. **CI/CD**:
   - Integração com ferramentas como ArgoCD, Flux ou Jenkins X
   - GitOps workflow
   - Testes automatizados antes do deploy

## 🛑 Remover os recursos

Para remover todos os recursos criados:

```bash
kubectl delete -f k8s/00-namespaces.yaml
```

Ou remover individualmente por namespace:

```bash
kubectl delete ns actajus observability monitoring
```


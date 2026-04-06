# ActaJus - Sistema de Observabilidade

Este projeto inclui uma infraestrutura completa de observabilidade com:

- **Jaeger**: Tracing distribuído
- **Prometheus**: Coleta de métricas
- **Grafana**: Visualização de métricas e logs
- **Loki**: Armazenamento e consulta de logs
- **PostgreSQL**: Banco de dados
- **NATS**: Sistema de mensagens

## Pré-requisitos

- Docker
- Docker Compose
- Git

## Como subir os serviços

### 1. Clone o repositório (se necessário)

```bash
git clone <seu-repositorio>
cd actajus
```

### 2. Suba todos os serviços

```bash
./start-observability.sh
```

### 3. Acesse os serviços

| Serviço    | URL                      | Credenciais      |
| ---------- | ------------------------ | ---------------- |
| Aplicação  | <http://localhost:8080>  | -                |
| Jaeger UI  | <http://localhost:16686> | -                |
| Prometheus | <http://localhost:9090>  | -                |
| Grafana    | <http://localhost:3000>  | admin / admin123 |
| Loki       | <http://localhost:3100>  | -                |

## Configuração do Grafana

### Adicionando Loki como fonte de dados

1. Acesse o Grafana em <http://localhost:3000>
2. Clique em "Administration" (ícone de engrenagem) → "Data Sources"
3. Clique em "Add data source"
4. Selecione "Loki"
5. Preencha as configurações:

```
HTTP URL: http://loki:3100
```

6. Clique em "Save & Test"

### Adicionando Prometheus como fonte de dados

1. Acesse o Grafana em <http://localhost:3000>
2. Clique em "Administration" (ícone de engrenagem) → "Data Sources"
3. Clique em "Add data source"
4. Selecione "Prometheus"
5. Preencha as configurações:

```
HTTP URL: http://prometheus:9090
```

6. Clique em "Save & Test"

## Testando o sistema

### 1. Faça requisições para a API

```bash
# Teste de login
curl -X POST http://localhost:8080/signin \
  -H "Content-Type: application/json" \
  -d '{"cpf":"12345678900", "password":"senha123"}'

# Teste de recuperação de senha
curl -X POST http://localhost:8080/auth/email \
  -H "Content-Type: application/json" \
  -d '{"cpf":"12345678900"}'

# Teste de solicitação de recuperação de senha
curl -X POST http://localhost:8080/auth/recover \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@example.com"}'
```

### 2. Verifique os dados no Jaeger

- Acesse <http://localhost:16686>
- Selecione o serviço "actajus-api"
- Execute requisições e veja os traces aparecerem

### 3. Verifique as métricas no Prometheus

- Acesse <http://localhost:9090>
- Verifique métricas como:
  - `http_requests_total`
  - `http_request_duration_seconds`
  - `auth_success_total`
  - `auth_failures_total`

### 4. Verifique os logs no Grafana

- Acesse <http://localhost:3000>
- Após configurar o Loki como fonte de dados
- Vá para "Explore" e selecione a fonte Loki
- Use queries como `{service_name="actajus-api"} |= "error"` para encontrar erros

## Arquitetura

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───▶│   ActaJus   │───▶│   Jaeger    │
│   App       │    │    API      │    │   (Tracing) │
└─────────────┘    └─────────────┘    └─────────────┘
                         │
                   ┌─────▼─────┐
                   │           │
              ┌────▼   DB    ▼────┐
              │  PostgreSQL   │   │
              └───────────────┘   │
                         │        │
                   ┌─────▼─────┐  │
                   │   NATS    │  │
                   │(Messaging)│  │
                   └───────────┘  │
                         │        │
              ┌──────────▼────────▼──┐
              │      Grafana         │
              │   (Visualization)    │
              └──────────────────────┘
```

## Variáveis de Ambiente

O sistema pode ser configurado com as seguintes variáveis de ambiente:

| Variável             | Descrição                    | Padrão      |
| -------------------- | ---------------------------- | ----------- |
| TRACING_ENABLED      | Habilita/desabilita tracing  | true        |
| TRACING_SERVICE_NAME | Nome do serviço para tracing | actajus-api |
| TRACING_AGENT_HOST   | Host do agente Jaeger        | jaeger      |
| TRACING_AGENT_PORT   | Porta do agente Jaeger       | 6831        |
| SERVER_PORT          | Porta do servidor HTTP       | 8080        |
| DB_HOST              | Host do banco de dados       | postgres    |
| DB_PORT              | Porta do banco de dados      | 5432        |

## URLs de Serviços

| Serviço | URL | Credenciais |
|---------|-----|-------------|
| Aplicação | http://localhost:8080 | - |
| Jaeger UI | http://localhost:16686 | - |
| Prometheus | http://localhost:9090 | - |
| Grafana | http://localhost:3000 | admin / admin123 |
| Loki | http://localhost:3100 | - |
| pgAdmin4 | http://localhost:5050 | admin@actajus.com.br / admin123 |

## Parando os serviços

```bash
./stop-observability.sh
```

Para remover volumes e dados persistentes:

```bash
docker-compose down -v
```

## Troubleshooting

### Serviços não sobem

- Verifique se o Docker está rodando
- Verifique se as portas 8080, 3000, 3100, 9090, 16686 estão disponíveis

### Logs da aplicação

```bash
docker-compose logs -f app
```

### Problemas de conexão com banco de dados

- Verifique se o PostgreSQL subiu corretamente
- Verifique as variáveis de ambiente no docker-compose.yml

## Estrutura de diretórios

```
actajus/
├── docker-compose.yml          # Configuração do Docker Compose
├── Dockerfile                  # Imagem da aplicação
├── start-observability.sh      # Script para subir os serviços
├── stop-observability.sh       # Script para parar os serviços
├── prometheus/
│   └── prometheus.yml         # Configuração do Prometheus
├── loki/
│   └── loki-config.yml        # Configuração do Loki
└── ...
```

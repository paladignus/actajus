# E2E

Suite E2E do dashboard WEB usando Playwright sobre o servidor HTTP real do projeto.

## Instalar

```bash
pnpm --dir web install
pnpm --dir web e2e:install
```

## Rodar smoke publico

```bash
pnpm --dir web e2e
```

Por padrao, o Playwright sobe o servidor com:

```bash
go run ../cmd/http/server.go
```

e espera o health check em `http://127.0.0.1:8080/health`.

## Reusar servidor ja rodando

```bash
ACTAJUS_E2E_REUSE_SERVER=true ACTAJUS_E2E_BASE_URL=http://127.0.0.1:8080 pnpm --dir web e2e
```

## Fluxo autenticado

O teste autenticado e opcional e depende de credenciais reais:

```bash
ACTAJUS_E2E_EMAIL=admin@example.com \
ACTAJUS_E2E_PASSWORD=secret \
pnpm --dir web e2e
```

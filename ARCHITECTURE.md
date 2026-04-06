# Architecture Rules

Este documento descreve as regras arquiteturais vigentes para a parte modular do projeto em `internal/module/*`.

## Escopo

- Ignorar `internal/application`, `internal/domain` e `internal/infrastructure` da raiz.
- A arquitetura alvo do projeto está nos módulos, em `internal/module/*`.

## Estrutura de Módulo

Cada módulo deve, quando aplicável, seguir a separação:

- `domain`
- `application`
- `infrastructure`
- `presentation`

Nem todo módulo precisa ter a mesma complexidade, mas a direção de dependência deve ser preservada.

## Dependency Rule

As dependências devem apontar para dentro.

Permitido:

- `application` -> `domain`
- `infrastructure` -> `application`, `domain`, `internal/shared/*`
- `presentation` -> `application`, `internal/shared/*`
- `module.go` -> composição interna do próprio módulo

Proibido entre módulos:

- importar `application` de outro módulo
- importar `presentation` de outro módulo
- importar `infrastructure` de outro módulo

Essa regra é verificada por teste em [internal/module/architecture_test.go](/home/marcelo/Workspace/go/projects/actajus/internal/module/architecture_test.go:1).

## Bounded Contexts

Cada módulo é um bounded context.

Consequências práticas:

- um módulo não deve reutilizar adapters de outro módulo
- um módulo não deve reutilizar repositories concretos de outro módulo
- se precisar falar com outro contexto, exponha porta própria no contexto consumidor
- tipos compartilhados entre contexts devem ser explícitos como shared kernel, não “atalhos”

## Shared Kernel

O que pode ser compartilhado deve morar em `internal/shared/*`.

Exemplos atuais:

- logger
- unit of work
- serialization
- messaging contracts
- paginação compartilhada
- value objects compartilhados

Se algo só faz sentido para um módulo, não pertence ao shared.

## Composition Roots

Os entrypoints não devem concentrar wiring detalhado de infraestrutura.

Hoje o projeto usa:

- runtime base compartilhado em [pkg/bootstrap/runtime.go](/home/marcelo/Workspace/go/projects/actajus/pkg/bootstrap/runtime.go:1)
- runtime especializado de notificação em [pkg/bootstrap/notification.go](/home/marcelo/Workspace/go/projects/actajus/pkg/bootstrap/notification.go:1)

Objetivo:

- `cmd/*` monta aplicação
- `pkg/bootstrap/*` concentra infraestrutura compartilhável
- módulos continuam responsáveis apenas por composição interna

## Contracts

Adapters protobuf devem ter testes de contrato.

Cobertura mínima esperada:

- request protobuf -> command/query
- read model -> response protobuf
- tipos opcionais e timestamps

Referências atuais:

- [company adapter tests](/home/marcelo/Workspace/go/projects/actajus/internal/module/company/presentation/grpc/adapter/company_adapter_test.go:1)
- [identity adapter tests](/home/marcelo/Workspace/go/projects/actajus/internal/module/identity/presentation/grpc/adapter/auth_adapter_test.go:1)

## Persistence

Repositórios críticos devem ter testes comportamentais com `pgxmock`.

Cobertura mínima esperada:

- `Create`
- `Update`
- `Delete` quando existir
- leitura principal
- mapeamento de erros relevantes do Postgres

Referências atuais:

- [company persistence tests](/home/marcelo/Workspace/go/projects/actajus/internal/module/company/infrastructure/persistence/database/postgres/address_repository_test.go:1)
- [person persistence tests](/home/marcelo/Workspace/go/projects/actajus/internal/module/person/infrastructure/persistence/database/postgres/person_test.go:1)

## Naming

Usar nomes semânticos de camada.

Evitar:

- `dto` genérico quando o conceito real é `command`, `query`, `message`, `event` ou `pagination`

Preferir:

- `command`
- `query`
- `readmodel`
- `message`
- `event`
- `pagination`

## Practical Standard

Ao adicionar novo código:

1. definir a responsabilidade no módulo correto
2. evitar importar camadas internas de outro módulo
3. preferir porta local a reuso de implementação concreta externa
4. adicionar teste de contrato se tocar adapter
5. adicionar teste de persistência se tocar SQL relevante

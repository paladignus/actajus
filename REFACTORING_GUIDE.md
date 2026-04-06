# 📋 Guia de Refatoração - Organização de DTOs

## **Problema Atual**

```
application/
├── dto/                    ← Mistura comandos, queries e read models
│   ├── login_command.go
│   ├── auth_tokens.go
│   └── ...
└── model/                  ← Modelos internos (confuso com dto)
    └── login_normalized.go
```

## **Solução Proposta**

```
application/
├── command/                ← Comandos (write operations)
│   ├── login_command.go
│   ├── change_password_command.go
│   └── create_company_command.go
├── query/                  ← Queries (read operations)
│   ├── find_company_by_id_query.go
│   └── list_companies_query.go
├── readmodel/              ← Read Models (resultados)
│   ├── auth_tokens_readmodel.go
│   ├── company_readmodel.go
│   └── user_readmodel.go
└── mapper/                 ← Mappers entre camadas
    ├── command_mapper.go
    └── readmodel_mapper.go
```

## **Exemplo de Migração**

### **ANTES**

```go
// application/dto/login_command.go
type LoginCommand struct {
    Email    string
    Password string
    IP       string
}

// application/dto/auth_tokens.go
type AuthTokensReadModel struct {
    IDSession    int64
    IDUser       int64
    AccessToken  string
    RefreshToken string
}

// application/model/login_normalized.go
type LoginNormalized struct {
    Email     string
    Password  string
    IP        string
    UserAgent string
}
```

### **DEPOIS**

```go
// application/command/login_command.go
type LoginCommand struct {
    Email    string
    Password string
    IP       string
    UserAgent string
}

// application/readmodel/auth_tokens_readmodel.go
type AuthTokensReadModel struct {
    IDSession    int64
    IDUser       int64
    AccessToken  string
    RefreshToken string
    ExpiresAt    time.Time
}

// application/mapper/login_mapper.go
type LoginMapper struct {
    validator *validation.Validator
}

func (m *LoginMapper) ToNormalized(cmd command.LoginCommand) (LoginNormalized, error) {
    // Validação e normalização
}
```

## **Benefícios**

1. **Clareza** - Cada pasta tem responsabilidade única
2. **Descoberta** - Fácil encontrar comandos vs read models
3. **CQRS** - Separação clara entre write e read
4. **Manutenção** - Menos confusão entre tipos

## **Script de Migração**

```bash
# Criar novas pastas
mkdir -p application/{command,query,readmodel,mapper}

# Mover arquivos (exemplo para identity)
mv application/dto/*_command.go application/command/
mv application/dto/*_readmodel.go application/readmodel/
mv application/dto/*_query.go application/query/
mv application/model/* application/mapper/

# Remover pastas vazias
rmdir application/dto application/model
```

## **Status por Módulo**

| Módulo | dto | model | Prioridade |
|--------|-----|-------|------------|
| identity | ✅ | ✅ | Alta |
| company | ✅ | ❌ | Média |
| person | ✅ | ❌ | Baixa |
| address | ✅ | ❌ | Baixa |
| phone | ✅ | ❌ | Baixa |
| email | ✅ | ❌ | Baixa |
| social_media | ✅ | ❌ | Baixa |
| notification | ✅ | ❌ | Baixa |

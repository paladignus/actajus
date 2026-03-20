# 📋 Guia de Refatoração - Items 3 e 5

## **Item 3: Use Cases com muitas dependências**

### **Problema**

Use cases como `Login` têm 8+ dependências, violando o princípio da responsabilidade única e dificultando testes.

```go
// ANTES - 8 dependências!
type Login struct {
    user       repository.UserRepository
    session    repository.SessionRepository
    hasher     service.PasswordHasher
    refresh    service.RefreshTokenService
    access     service.AccessTokenService
    clock      service.Clock
    config     config.AuthConfig
    mapper     mapper.AuthMapper
    projection mapper.AuthProjectionMapper
}
```

### **Solução**

Agrupar dependências relacionadas em **facets** ou **services**:

```go
// DEPOIS - 3 dependências principais
type LoginDeps struct {
    Authn      AuthnService      // user + session + hasher
    Tokens     TokenService      // refresh + access + clock + config
    Projection TokenProjection   // mapper + projection
}

type Login struct {
    deps LoginDeps
}
```

### **Benefícios**

1. **Menos dependências** - De 8 para 3
2. **Maior coesão** - Serviços relacionados agrupados
3. **Testabilidade** - Mocks mais simples
4. **Reuso** - AuthnService pode ser usado em outros usecases

---

## **Item 5: Infrastructure no main.go**

### **Problema**

```go
// cmd/main.go
func main() {
    hasher := security.NewArgon2idPasswordHasher()
    db, _ := sharedPostgres.NewConnection(...)
    
    // Lógica de infraestrutura misturada com bootstrap
    http.HandleFunc(...)
    http.ListenAndServe(...)
}
```

### **Solução**

Mover para **pkg/di** (Dependency Injection):

```
cmd/
├── http/
│   └── server.go       ← Ponto de entrada HTTP
├── grpc/
│   └── server.go       ← Ponto de entrada gRPC
└── main.go             ← Mínimo, delega para servers

pkg/
└── di/
    ├── container.go    ← DI Container
    ├── database.go     ← Database setup
    ├── modules.go      ← Module initializers
    └── config.go       ← Config loading
```

### **Estrutura do DI Container**

```go
// pkg/di/container.go
type Container struct {
    Config      *config.Config
    DB          *sql.DB
    Logger      Logger
    Modules     Modules
}

type Modules struct {
    Identity    identity.Module
    Company     company.Module
    // ...
}

func NewContainer() (*Container, error) {
    cfg := config.Load()
    logger := logger.New(cfg)
    db := database.New(cfg.Database)
    
    modules, err := initializeModules(db, logger, cfg)
    if err != nil {
        return nil, err
    }
    
    return &Container{cfg, db, logger, modules}, nil
}
```

### **Bootstrap limpo**

```go
// cmd/http/server.go
func main() {
    container, err := di.NewContainer()
    if err != nil {
        log.Fatal(err)
    }
    defer container.Close()
    
    server := httpserver.New(container)
    server.Start()
}
```

# sIOC - Simple Inversion of Control

sIOC is a Go library that provides a simple dependency injection system for Go applications.

## Versões Disponíveis

### v1 (Recomendada) 🚀
A versão mais recente com API moderna, containers isolados e thread-safety.

- **Documentação**: [v1.md](./doc/v1.md)
- **Características**: Containers isolados, thread-safe, generics, API limpa
- **Uso**: `import "github.com/sergiodii/sioc/v1"`

### v0 (Legacy)
Versão inicial com padrão singleton global.

- **Documentação**: [v0.md](./doc/v0.md)
- **Características**: Singleton global, API simples
- **Uso**: `import "github.com/sergiodii/sioc/v0"`

## Quick Start (v1)

```go
package main

import (
    "fmt"
    "github.com/sergiodii/sioc/v1"
)

type UserService struct {
    Name string
}

func main() {
    // Cria um container
    container := sioc.NewContainer()
    
    // Registra um serviço
    sioc.Inject(&UserService{Name: "admin"}, container)
    
    // Resolve o serviço
    service := sioc.Get[*UserService](container)
    fmt.Println(service.Name) // Output: admin
}
```

## Installation

```bash
# Para v1 (recomendada)
go get github.com/sergiodii/sioc/v1

# Para v0 (legacy)
go get github.com/sergiodii/sioc/v0
```

## Documentação Completa

- **[v1 - Documentação Completa](./doc/v1.md)** - Versão recomendada com containers isolados
- **[v0 - Documentação Completa](./doc/v0.md)** - Versão legacy com singleton global

## Exemplos Básicos

### v1 - Exemplo com Inicialização

```go
package main

import (
    "fmt"
    "github.com/sergiodii/sioc/v1"
)

type Database struct {
    Connected bool
}

func (d *Database) Init() {
    d.Connected = true
}

type UserService struct {
    db *Database
}

func (u *UserService) Init(db *Database) {
    u.db = db
}

func main() {
    container := sioc.NewContainer()
    
    // Registra as dependências
    sioc.Inject(&Database{}, container)
    sioc.Inject(&UserService{}, container)
    
    // Inicializa todas as dependências
    sioc.Init(container)
    
    // Resolve o serviço
    userService := sioc.Get[*UserService](container)
    fmt.Printf("Database connected: %v\n", userService.db.Connected)
}
```

### v0 - Exemplo Básico

```go
package main

import (
    "fmt"
    "github.com/sergiodii/sioc/v0"
)

type Service struct {
    Name string
}

func (s *Service) Init() {
    fmt.Println("Service initialized")
}

func main() {
    sioc.Start()
    sioc.Register(&Service{Name: "test"})
    sioc.Init()
    
    service := sioc.Get[*Service]()
    fmt.Println(service.Name)
}
```

## Características Principais

### v1 (Recomendada)
- ✅ **Containers Isolados**: Múltiplos contextos independentes
- ✅ **Thread Safety**: Operações concorrentes seguras
- ✅ **Generics**: Suporte completo a tipos genéricos
- ✅ **API Limpa**: Interface moderna e intuitiva
- ✅ **Testabilidade**: Fácil criação de containers para testes
- ✅ **Sanitização**: Nomes de serviços automaticamente limpos

### v0 (Legacy)
- ⚠️ **Singleton Global**: Todas as dependências gerenciadas globalmente
- ⚠️ **API Simples**: Interface básica e direta
- ⚠️ **Compatibilidade**: Mantida para projetos existentes

## Migração

Para migrar da v0 para v1, consulte o guia de migração na [documentação da v1](./doc/v1.md#migração-da-v0).


## Contribuição

Contribuições são bem-vindas! Por favor, abra uma issue ou pull request.

## Licença

MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.
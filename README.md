# SIOC - Simple Injection Of Control

SIOC is a Go library that provides a simple dependency injection system.

## Installation

```bash
go get github.com/sergiodii/sioc
```

## Usage 

### Registering an instance and getting an instance


You can register any instance of any type, but it must be a pointer.
register the instance using the Register function.

```go
sioc.Register(instance)
```

After registering the instance, you can get the instance using the Get function.

```go
result := sioc.Get[T]()
```

You should pass the type of the instance to the Get function. You can pass a interface or a struct.

```go
import "github.com/sergiodii/sioc"


type Service struct {
}

func (s *Service) DoSomething() {
	fmt.Println("Doing something")
}

func main() {
	sioc.Register(&Service{})
	executeService()
}

// executeService is a function that executes the service.
// call service in any part of your code
func executeService() {

    // Get the service instance, pass the type of the instance
	service := sioc.Get[Service]()
	service.DoSomething()
}

```

### Initializing as a constructor

As we all know, there is no constructor in Go, but we can simulate the behavior of a constructor that automatically injects the dependencies that are in the parameter of this function.

I present the "Constructor", Init.

```go
sioc.Init()
```

Every structure you create and have a function that has the name Init, it will be executed when you call the Init of the sIOC.

```go

import "github.com/sergiodii/sioc"

type Service struct {
}

func (s *Service) Init() {
	fmt.Println("Service initialized")
}

func main() {
	sioc.Register(&Service{})
	sioc.Init()
    // "Service initialized" will be printed
}
```

With this, you can create a structure and inject the dependencies you need in the constructor of the structure.

```go
type Service struct {
    db *sql.DB
}

func (s *Service) Init() {
    s.db = sql.NewDB()
}
```

Inject the dependency in the constructor of the structure.

you can create a function Init and inject the dependencies inside it, through the parameter of this function.

```go
type Service struct {
    db *sql.DB
}

func (s *Service) Init(db *sql.DB) {
    s.db = db
}
```

Ready! Now you can use the Service structure and the db dependency injected in the constructor.

```go
import "github.com/sergiodii/sioc"

type Service struct {
    db *sql.DB
}

func (s *Service) Init(db *sql.DB) {
    s.db = db
}


func main() {
	sioc.Register(&Service{})
	sioc.Register(&sql.DB{})
	sioc.Init()

    // the sIOC will inject the db dependency in the constructor of the Service structure
}
```

The order of registration does not matter, as the sIOC will call the dependencies in order, that is, if you register a dependency after registering the structure that depends on it, the sIOC will still find the dependency and inject it into the structure correctly.


license MIT
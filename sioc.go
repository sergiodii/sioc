package sioc

import v0_injection "github.com/sergiodii/sioc/v0"

// Get returns an instance of type T from the dependency injection container.
// If no instance of type T is found, it will panic with a fatal error.
func Get[T any]() T {
	return v0_injection.Get[T]()
}

// Init initializes the dependency injection container.
// It must be called after using the Register function.
func Init() {
	v0_injection.Init()
}

// Register registers an instance of type T in the dependency injection container.
func Register[T any](instance T) {
	v0_injection.Inject(instance)
}

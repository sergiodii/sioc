package sioc

import (
	"sync"

	"github.com/sergiodii/sioc/extension/text"
)

// ServiceContainer defines the interface for a dependency injection container.
type ServiceContainer interface {
	Register(serviceKey string, serviceInstance any)
	Resolve(serviceKey string) (any, bool)
	ListAll() []any
	Count() int
}

// serviceRegistry implements the ServiceContainer interface using sync.Map for concurrency safety.
type serviceRegistry struct {
	services sync.Map
}

// NewContainer creates a new, empty service container instance.
func NewContainer() ServiceContainer {
	return &serviceRegistry{}
}

// Register stores a service instance in the container under the given key.
func (sr *serviceRegistry) Register(serviceKey string, serviceInstance any) {
	sr.services.Store(text.Sanitize(serviceKey), serviceInstance)
}

// Resolve retrieves a service instance by key. Returns (nil, false) if not found.
func (sr *serviceRegistry) Resolve(serviceKey string) (any, bool) {
	return sr.services.Load(text.Sanitize(serviceKey))
}

// ListAll returns a slice of all registered service instances.
func (sr *serviceRegistry) ListAll() []any {
	var serviceList []any
	sr.services.Range(func(_, serviceInstance any) bool {
		serviceList = append(serviceList, serviceInstance)
		return true
	})
	return serviceList
}

// Count returns the number of registered service instances.
func (sr *serviceRegistry) Count() int {
	serviceCount := 0
	sr.services.Range(func(_, _ any) bool {
		serviceCount++
		return true
	})
	return serviceCount
}

package sioc

import (
	"reflect"

	"github.com/sergiodii/sioc/extension/text"
)

// serviceWrapper is a generic wrapper for service instances.
type serviceWrapper[T any] struct {
	serviceInstance T
	serviceName     string
}

// NewServiceWrapper creates a new service wrapper for type T.
func NewServiceWrapper[T any]() ServiceWrapper[T] {
	return &serviceWrapper[T]{}
}

// SetService sets the service instance and its name in the wrapper.
func (sw *serviceWrapper[T]) SetService(serviceInstance T) ServiceWrapper[T] {
	sw.serviceName = sw.sanitizeServiceName(reflect.TypeOf(serviceInstance).String())
	sw.serviceInstance = serviceInstance
	return sw
}

// Backward compatibility: AddInstance is an alias for SetService.
func (sw *serviceWrapper[T]) AddInstance(instance T) ServiceWrapper[T] {
	return sw.SetService(instance)
}

// GetService returns the stored service instance.
func (sw *serviceWrapper[T]) GetService() T {
	return sw.serviceInstance
}

// Backward compatibility: GetInstance is an alias for GetService.
func (sw *serviceWrapper[T]) GetInstance() T {
	return sw.GetService()
}

// CreateNewService returns a copy of the stored service instance (if possible).
func (sw *serviceWrapper[T]) CreateNewService() T {
	newService := new(T)
	*newService = sw.serviceInstance
	return *newService
}

// Backward compatibility: GetNewInstance is an alias for CreateNewService.
func (sw *serviceWrapper[T]) GetNewInstance() T {
	return sw.CreateNewService()
}

// MatchesServiceName checks if the given name matches the wrapper's service name.
func (sw *serviceWrapper[T]) MatchesServiceName(serviceName string) bool {
	return sw.serviceName == sw.sanitizeServiceName(serviceName)
}

// Backward compatibility: MatchWithName is an alias for MatchesServiceName.
func (sw *serviceWrapper[T]) MatchWithName(name string) bool {
	return sw.MatchesServiceName(name)
}

// sanitizeServiceName sanitizes the service name for consistent storage.
func (sw *serviceWrapper[T]) sanitizeServiceName(serviceName string) string {
	return text.Sanitize(serviceName)
}

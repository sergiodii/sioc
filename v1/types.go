package sioc

// ServiceProvider is a marker interface for types that can provide services.
type ServiceProvider interface {
	ProvideService() any
}

// ServiceWrapper is a generic interface for dependency wrappers.
type ServiceWrapper[T any] interface {
	GetService() T
	CreateNewService() T
	MatchesServiceName(serviceName string) bool
	SetService(serviceInstance T) ServiceWrapper[T]
}

// Backward compatibility type aliases
type Container = ServiceContainer

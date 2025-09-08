package sioc

import (
	"log"
	"reflect"
	"runtime"
	"strings"
)

// NewInjector creates a new service wrapper (backward compatibility).
func NewInjector[T any]() ServiceWrapper[T] {
	return NewServiceWrapper[T]()
}

// Get retrieves a service instance of type T from the container.
// It checks for both direct type and interface implementations.
func Get[T any](serviceContainer ServiceContainer) T {
	targetType := reflect.TypeOf((*T)(nil)).Elem()

	// Try direct lookup by type name
	if serviceInstance, found := serviceContainer.Resolve(targetType.String()); found {
		if wrapper, ok := serviceInstance.(ServiceWrapper[T]); ok {
			return wrapper.GetService()
		}
		if wrapperPtr, ok := serviceInstance.(ServiceWrapper[*T]); ok {
			return *wrapperPtr.GetService()
		}
		// Also try ServiceWrapper[any] for backward compatibility
		if wrapperAny, ok := serviceInstance.(ServiceWrapper[any]); ok {
			service := wrapperAny.GetService()
			if typedService, ok := service.(T); ok {
				return typedService
			}
		}
	}

	// Try to find by interface implementation or pointer match
	for _, registeredService := range serviceContainer.ListAll() {
		// Try ServiceWrapper[any] first (most common case)
		if wrapperAny, ok := registeredService.(ServiceWrapper[any]); ok {
			serviceInstance := wrapperAny.GetService()
			// Check direct type match
			if typedService, ok := serviceInstance.(T); ok {
				return typedService
			}
			// Check interface implementation
			if targetType.Kind() == reflect.Interface && reflect.TypeOf(serviceInstance).Implements(targetType) {
				if typedService, ok := serviceInstance.(T); ok {
					return typedService
				}
			}
		}

		// Try specific type wrappers
		if wrapper, ok := registeredService.(ServiceWrapper[T]); ok {
			serviceInstance := wrapper.GetService()
			if targetType.Kind() == reflect.Interface && reflect.TypeOf(serviceInstance).Implements(targetType) {
				return serviceInstance
			}
		}
		if wrapperPtr, ok := registeredService.(ServiceWrapper[*T]); ok {
			serviceInstance := wrapperPtr.GetService()
			if targetType.Kind() == reflect.Interface && reflect.TypeOf(serviceInstance).Implements(targetType) {
				return *serviceInstance
			}
		}
	}

	log.Fatalf("Service of type %s not found in container", targetType)
	var emptyService T
	return emptyService
}

// ResolveService is an alias for Get for more descriptive naming.
func ResolveService[T any](serviceContainer ServiceContainer) T {
	return Get[T](serviceContainer)
}

// Inject registers a service instance in the container, wrapping it in a ServiceWrapper.
func Inject(serviceInstance any, serviceContainer ServiceContainer) {
	wrapper := NewServiceWrapper[any]()
	wrapper.SetService(serviceInstance)
	serviceContainer.Register(reflect.TypeOf(serviceInstance).String(), wrapper)
}

// RegisterService is an alias for Inject for more descriptive naming.
func RegisterService(serviceInstance any, serviceContainer ServiceContainer) {
	Inject(serviceInstance, serviceContainer)
}

// GetFunctionName returns the name of a function from its value.
func GetFunctionName(functionValue interface{}) string {
	functionPointer := reflect.ValueOf(functionValue).Pointer()
	functionName := runtime.FuncForPC(functionPointer).Name()
	nameParts := strings.Split(functionName, ".")
	return nameParts[len(nameParts)-1]
}

// ExtractFunctionName is an alias for GetFunctionName for more descriptive naming.
func ExtractFunctionName(functionValue interface{}) string {
	return GetFunctionName(functionValue)
}

// Init calls the Init method on all registered services that have it, resolving dependencies.
func Init(serviceContainer ServiceContainer) {
	dependencyMap := make(map[reflect.Type]ServiceWrapper[any])
	for _, registeredService := range serviceContainer.ListAll() {
		if wrapper, ok := registeredService.(ServiceWrapper[any]); ok {
			dependencyMap[reflect.TypeOf(wrapper.GetService())] = wrapper
		}
	}

	for _, registeredService := range serviceContainer.ListAll() {
		wrapper, ok := registeredService.(ServiceWrapper[any])
		if !ok {
			continue
		}
		serviceInstance := wrapper.GetService()
		initializationMethod := reflect.ValueOf(serviceInstance).MethodByName("Init")
		if !initializationMethod.IsValid() {
			continue
		}

		methodType := initializationMethod.Type()
		methodParams := make([]reflect.Value, methodType.NumIn())
		for paramIndex := 0; paramIndex < methodType.NumIn(); paramIndex++ {
			parameterType := methodType.In(paramIndex)
			if dependency, exists := dependencyMap[parameterType]; exists {
				methodParams[paramIndex] = reflect.ValueOf(dependency.GetService())
			} else {
				log.Fatalf("Dependency not found for Init of %v: %v", reflect.TypeOf(serviceInstance), parameterType)
			}
		}
		initializationMethod.Call(methodParams)
	}
}

// InitializeServices is an alias for Init for more descriptive naming.
func InitializeServices(serviceContainer ServiceContainer) {
	Init(serviceContainer)
}

package sioc

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

// Get retrieves a service instance of type T from the container.
// It checks for both direct type and interface implementations.
func Get[T any](serviceContainer ServiceContainer) T {
	targetType := reflect.TypeOf((*T)(nil)).Elem()

	if serviceInstance, found := serviceContainer.Resolve(targetType.String()); found {
		if wrapper, ok := serviceInstance.(ServiceWrapper[T]); ok {
			fmt.Println("Achei aqui 1")
			return wrapper.GetService()
		}
		if wrapperPtr, ok := serviceInstance.(ServiceWrapper[*T]); ok {
			fmt.Println("Achei aqui 2")
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

			// Verifica se o tipo é um ponteiro
			if targetType.Kind() == reflect.Ptr {
				typeT := targetType.Elem()
				// Verifica se o tipo do serviço é um ponteiro para o tipo desejado
				if reflect.TypeOf(serviceInstance) == reflect.PtrTo(typeT) {
					return serviceInstance.(T)
				}
			} else {
				// Para tipos não ponteiros, verifica o tipo base
				if reflect.TypeOf(serviceInstance) == targetType {
					return serviceInstance.(T)
				}

				// if item.(v1_interfaces.Injector[T]).MatchWithName("*" + typeT.String()) {
				// 	return *item.(v1_interfaces.Injector[*T]).GetInstance()
				// }
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

// Inject registers a service instance in the container, wrapping it in a ServiceWrapper.
func Inject(serviceInstance any, serviceContainer ServiceContainer) {
	wrapper := NewServiceWrapper[any]()
	wrapper.SetService(serviceInstance)
	serviceContainer.Register(reflect.TypeOf(serviceInstance).String(), wrapper)
}

// GetFunctionName returns the name of a function from its value.
func GetFunctionName(functionValue interface{}) string {
	functionPointer := reflect.ValueOf(functionValue).Pointer()
	functionName := runtime.FuncForPC(functionPointer).Name()
	nameParts := strings.Split(functionName, ".")
	return nameParts[len(nameParts)-1]
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
			dependency, exists := dependencyMap[parameterType]
			if !exists {
				continue
			}
			methodParams[paramIndex] = reflect.ValueOf(dependency.GetService())
		}
		initializationMethod.Call(methodParams)
	}
}

// func Init2(container Container) {

// 	// First step: Map all required dependencies
// 	dependencyMap := make(map[reflect.Type]Injector[interface{}])
// 	for _, item := range container.ListAll() {
// 		dependencyMap[reflect.TypeOf(item.(Injector[interface{}]).GetInstance())] = item.(Injector[interface{}])
// 	}

// 	for _, item := range container.ListAll() {
// 		instance := item.(v1_interfaces.Injector[interface{}]).GetInstance()
// 		init := reflect.ValueOf(instance).MethodByName("Init")

// 		if init.IsValid() {
// 			// Check Init method parameters
// 			initType := init.Type()
// 			if initType.NumIn() > 0 {

// 				// Prepare required parameters
// 				params := make([]reflect.Value, initType.NumIn())

// 				for i := 0; i < initType.NumIn(); i++ {
// 					paramType := initType.In(i)
// 					preKeyName := "offers_platform_core_injection.InitializeNewInstanceTo"
// 					isPreKey := paramType.String() == preKeyName
// 					isNewInstance := (i != 0) && initType.In(i-1).String() == preKeyName
// 					if isPreKey && !isNewInstance {
// 						n := InitializeNewInstanceTo(NEW)
// 						params[i] = reflect.ValueOf(n)
// 						continue
// 					}

// 					if dep, exists := dependencyMap[paramType]; exists {
// 						if isNewInstance {
// 							params[i] = reflect.ValueOf(dep.GetNewInstance())
// 						} else {
// 							params[i] = reflect.ValueOf(dep.GetInstance())
// 						}
// 						continue
// 					}
// 					// Verifica se o parâmetro é do tipo Module
// 					// isModule := paramType == reflect.TypeOf(Module{}) || paramType == reflect.TypeOf(&Module{}) || (len(paramType.Name()) != 0 && paramType.Implements(reflect.TypeOf(Module{}).Elem()))
// 					dependenciaEncontrada := false
// 					for ref, instance := range dependencyMap {

// 						if len(paramType.Name()) != 0 && ref.Implements(paramType) {
// 							if isNewInstance {
// 								params[i] = reflect.ValueOf(instance.GetNewInstance())
// 							} else {
// 								params[i] = reflect.ValueOf(instance.GetInstance())
// 							}
// 							dependenciaEncontrada = true
// 							break
// 						}

// 					}

// 					if !dependenciaEncontrada {
// 						log.Fatalf("Dependência não encontrada para Init de %v: %v", reflect.TypeOf(instance), paramType)
// 					}

// 				}

// 				init.Call(params)
// 			} else {
// 				// If no parameters, call normally
// 				init.Call(nil)
// 			}
// 		}

// 		fmt.Println("Offer Platform Handler: ", reflect.ValueOf(instance).Type().String(), ", started")
// 	}
// }

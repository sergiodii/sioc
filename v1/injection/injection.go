package v1_injection

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"

	v1_container "github.com/sergiodii/sIOC/v1/container"
	v1_interfaces "github.com/sergiodii/sIOC/v1/interfaces"
)

func Get[T any](container v1_container.Container) T {
	typeT := reflect.TypeOf((*T)(nil))
	instance, ok := container.Get(typeT.String())
	s := typeT.String()
	fmt.Println(s)
	if ok {
		r := instance.(v1_interfaces.Injector[*T]).GetInstance()
		return *r
	}

	for _, item := range container.GetAll() {

		typeT := reflect.TypeOf((*T)(nil)).Elem()

		// Verifica se o tipo é uma interface
		if typeT.Kind() == reflect.Interface {
			if reflect.TypeOf(item.(v1_interfaces.Injector[T]).GetInstance()).Implements(typeT) {
				return item.(v1_interfaces.Injector[T]).GetInstance()
			}
			continue
		}

		// Verifica se o tipo é um ponteiro
		if typeT.Kind() == reflect.Ptr {
			if item.(v1_interfaces.Injector[T]).MatchWithName(typeT.String()) {
				return item.(v1_interfaces.Injector[T]).GetInstance()
			}
		} else {
			// Para tipos não ponteiros, verifica o tipo base
			if item.(v1_interfaces.Injector[T]).MatchWithName("*" + typeT.String()) {
				return *item.(v1_interfaces.Injector[*T]).GetInstance()
			}
		}
	}

	log.Fatalf("Instância do tipo %s não encontrada", typeT)
	var empty T
	return empty
}

func Inject(cls interface{}, container v1_container.Container) {

	newInj := NewInjector[interface{}]()

	newInj.AddInstance(cls)

	container.Set(reflect.TypeOf(cls).String(), newInj)

}

func GetFunctionName(temp interface{}) string {
	p := reflect.ValueOf(temp).Pointer()
	k := runtime.FuncForPC(p).Name()
	strs := strings.Split((k), ".")
	return strs[len(strs)-1]
}

func Init(container v1_container.Container) {

	// First step: Map all required dependencies
	dependencyMap := make(map[reflect.Type]v1_interfaces.Injector[interface{}])
	for _, item := range container.GetAll() {
		dependencyMap[reflect.TypeOf(item.(v1_interfaces.Injector[interface{}]).GetInstance())] = item.(v1_interfaces.Injector[interface{}])
	}

	for _, item := range container.GetAll() {
		instance := item.(v1_interfaces.Injector[interface{}]).GetInstance()
		init := reflect.ValueOf(instance).MethodByName("Init")

		if init.IsValid() {
			// Check Init method parameters
			initType := init.Type()
			if initType.NumIn() > 0 {

				// Prepare required parameters
				params := make([]reflect.Value, initType.NumIn())

				for i := 0; i < initType.NumIn(); i++ {
					paramType := initType.In(i)
					preKeyName := "offers_platform_core_injection.InitializeNewInstanceTo"
					isPreKey := paramType.String() == preKeyName
					isNewInstance := (i != 0) && initType.In(i-1).String() == preKeyName
					if isPreKey && !isNewInstance {
						n := InitializeNewInstanceTo(NEW)
						params[i] = reflect.ValueOf(n)
						continue
					}

					if dep, exists := dependencyMap[paramType]; exists {
						if isNewInstance {
							params[i] = reflect.ValueOf(dep.GetNewInstance())
						} else {
							params[i] = reflect.ValueOf(dep.GetInstance())
						}
						continue
					}
					// Verifica se o parâmetro é do tipo Module
					// isModule := paramType == reflect.TypeOf(Module{}) || paramType == reflect.TypeOf(&Module{}) || (len(paramType.Name()) != 0 && paramType.Implements(reflect.TypeOf(Module{}).Elem()))
					dependenciaEncontrada := false
					for ref, instance := range dependencyMap {

						if len(paramType.Name()) != 0 && ref.Implements(paramType) {
							if isNewInstance {
								params[i] = reflect.ValueOf(instance.GetNewInstance())
							} else {
								params[i] = reflect.ValueOf(instance.GetInstance())
							}
							dependenciaEncontrada = true
							break
						}

					}

					if !dependenciaEncontrada {
						log.Fatalf("Dependência não encontrada para Init de %v: %v", reflect.TypeOf(instance), paramType)
					}

				}

				init.Call(params)
			} else {
				// If no parameters, call normally
				init.Call(nil)
			}
		}

		fmt.Println("Offer Platform Handler: ", reflect.ValueOf(instance).Type().String(), ", started")
	}
}

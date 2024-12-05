package v0_injection

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	// v0_injection "github.com/sergiodii/sioc/v0"
)

var injOnce sync.Once
var List *[]Injector[interface{}]
var inj *injector[Injector[interface{}]]

type injector[T any] struct {
	List []*T
}

func Start() {
	injOnce.Do(func() {
		inj = &injector[Injector[interface{}]]{}
	})
}

func Get[T any]() T {
	Start()
	var element T
	typeT := reflect.TypeOf((*T)(nil)).Elem()

	for _, item := range inj.List {
		if instance := tryGetInstance[T](item, typeT); instance != nil {
			return *instance
		}
	}

	log.Fatalf("Instance type %s not found", reflect.TypeOf((*T)(nil)))
	return element
}

func tryGetInstance[T any](item *Injector[interface{}], typeT reflect.Type) *T {
	if typeT.Kind() == reflect.Interface {
		if reflect.TypeOf(item.GetInstance()).Implements(typeT) {
			result := item.GetInstance().(T)
			return &result
		}
		return nil
	}

	if typeT.Kind() == reflect.Ptr {
		if item.MatchWithName(typeT.String()) {
			result := item.GetInstance().(T)
			return &result
		}
	} else if item.MatchWithName("*" + typeT.String()) {
		result := *item.GetInstance().(*T)
		return &result
	}

	return nil
}

func Len() int {
	return len(inj.List)
}

// Inject injects a class into the dependency injection container.
// It will panic if the class is not a pointer.
func Register(cls interface{}) {
	Start()

	fmt.Println("InitializeNewInstanceTo", reflect.TypeOf(InitializeNewInstanceTo("")).String())

	if reflect.TypeOf(cls).Kind() != reflect.Ptr {
		log.Fatalf("%s is not a pointer", reflect.TypeOf(cls).String())
	}

	injectInstance(cls)
	injectFromIInjector(cls)

}

func injectInstance(cls interface{}) {
	newInj := NewInjector[interface{}]()
	newInj.AddInstance(cls)
	inj.List = append(inj.List, &newInj)
}

func injectFromIInjector(cls interface{}) {
	iInjectorImplementation := reflect.TypeOf((*IInjector)(nil)).Elem()
	if reflect.TypeOf(cls).Implements(iInjectorImplementation) {
		newInj := NewInjector[interface{}]()
		newInj.AddInstance(cls.(IInjector).InjectorInsertion())
		inj.List = append(inj.List, &newInj)
	}
}

func GetFunctionName(temp interface{}) string {
	p := reflect.ValueOf(temp).Pointer()
	k := runtime.FuncForPC(p).Name()
	strs := strings.Split((k), ".")
	return strs[len(strs)-1]
}

type Status struct {
	initialized bool
	instance    *Injector[interface{}]
}

func Init() {
	Start()
	dependencyMap := buildDependencyMap()
	initializeInjectors(dependencyMap)
}

func buildDependencyMap() map[reflect.Type]*Injector[interface{}] {
	dependencyMap := make(map[reflect.Type]*Injector[interface{}])
	for _, item := range inj.List {
		dependencyMap[reflect.TypeOf(item.GetInstance())] = item
	}
	return dependencyMap
}

func initializeInjectors(dependencyMap map[reflect.Type]*Injector[interface{}]) {
	for _, item := range inj.List {
		initStart(item, dependencyMap)
		fmt.Println("Injection: ", item.incjetionName, ", started")
	}
}

func ClearList() {
	inj.List = []*Injector[interface{}]{}
}

func initStart(inj *Injector[interface{}], m map[reflect.Type]*Injector[interface{}]) {
	if inj.Initialized() {
		return
	}

	init := reflect.ValueOf(inj.GetInstance()).MethodByName("Init")
	if !init.IsValid() {
		return
	}

	initType := init.Type()
	if initType.NumIn() > 0 {
		params := prepareInitParams(inj, initType, m)
		init.Call(params)
	} else {
		init.Call(nil)
	}

	inj.SetInitialization()
}

func prepareInitParams(inj *Injector[interface{}], initType reflect.Type, m map[reflect.Type]*Injector[interface{}]) []reflect.Value {
	params := make([]reflect.Value, initType.NumIn())

	for i := 0; i < initType.NumIn(); i++ {
		paramType := initType.In(i)
		if handlePreKeyParam(paramType, i, initType, params) {
			continue
		}

		if dep, exists := m[paramType]; exists {
			executeDependency(dep, m, isNewInstance(i, initType), params, i)
			continue
		}

		findAndExecuteDependency(paramType, m, i, initType, params)
	}

	return params
}

func handlePreKeyParam(paramType reflect.Type, i int, initType reflect.Type, params []reflect.Value) bool {
	preKeyType := reflect.TypeOf(InitializeNewInstanceTo(""))
	isPreKey := paramType == preKeyType
	isNewInstance := (i != 0) && initType.In(i-1) == preKeyType

	if isPreKey && !isNewInstance {
		n := InitializeNewInstanceTo(NEW)
		params[i] = reflect.ValueOf(n)
		return true
	}
	return false
}

func isNewInstance(i int, initType reflect.Type) bool {
	return (i != 0) && initType.In(i-1).String() == "offers_platform_core_injection.InitializeNewInstanceTo"
}

func findAndExecuteDependency(paramType reflect.Type, m map[reflect.Type]*Injector[interface{}], i int, initType reflect.Type, params []reflect.Value) {
	for ref, dep := range m {
		if matchDependency(ref, paramType, dep, m, i, initType, params) {
			return
		}
	}
	log.Fatalf("Dependency not found for Init of %v: %v", reflect.TypeOf(i), paramType)
}

func matchDependency(ref reflect.Type, paramType reflect.Type, dep *Injector[interface{}], m map[reflect.Type]*Injector[interface{}], i int, initType reflect.Type, params []reflect.Value) bool {
	if NewInjector[interface{}]().transformName(ref.String()) == paramType.String() {
		executeDependency(dep, m, isNewInstance(i, initType), params, i)
		return true
	}

	if (len(paramType.Name()) != 0 && paramType.Kind() == reflect.Interface) && ref.Implements(paramType) {
		executeDependency(dep, m, isNewInstance(i, initType), params, i)
		return true
	}

	return false
}

func executeDependency(dep *Injector[interface{}], m map[reflect.Type]*Injector[interface{}], isNewInstance bool, params []reflect.Value, i int) {
	if !dep.initialized {
		initStart(dep, m)
	}
	if isNewInstance {
		params[i] = reflect.ValueOf(dep.GetNewInstance())
	} else {
		params[i] = reflect.ValueOf(dep.GetInstance())
	}
}

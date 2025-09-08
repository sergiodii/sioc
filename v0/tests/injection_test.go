package sioc_test

import (
	"os"
	"testing"

	"github.com/sergiodii/sioc/v0"
)

type TestStruct struct {
	Name string
}

type TestStructWithInit struct {
	Initialized bool
}

func (t *TestStructWithInit) Init() {
	t.Initialized = true
}

type TestInjector struct{}

func (t *TestInjector) InjectorInsertion() any {
	return &TestStruct{Name: "injected"}
}

type ITestInterface interface {
	GetName() string
}

type TestStructWithInterface struct {
	Name string
}

func (t *TestStructWithInterface) GetName() string {
	return t.Name
}

func TestStartShouldInitializeInjector(t *testing.T) {
	sioc.Start()
	if sioc.Len() != 0 {
		t.Error("Expected empty injector after Start()")
	}
}

func TestInjectShouldAddInstance(t *testing.T) {
	sioc.Start()
	testStruct := &TestStruct{Name: "test"}
	sioc.Register(testStruct)

	if sioc.Len() != 1 {
		t.Error("Expected one instance after Inject()")
	}
	sioc.ClearList()
}

func TestGetShouldReturnInjectedInstance(t *testing.T) {
	sioc.Start()
	testStruct := &TestStruct{Name: "test"}
	sioc.Register(testStruct)
	result := sioc.Get[*TestStruct]()
	if result.Name != "test" {
		t.Error("Expected to get injected instance")
	}
	sioc.ClearList()
}

func TestInjectWithInjectorInterface(t *testing.T) {
	sioc.Start()
	testInjector := &TestInjector{}
	sioc.Register(testInjector)

	if sioc.Len() != 2 {
		t.Error("Expected two instances after injecting IInjector, has: ", sioc.Len())
	}
	sioc.ClearList()
}

func TestGetInjectedFromInjector(t *testing.T) {
	sioc.Start()
	testInjector := &TestInjector{}
	sioc.Register(testInjector)
	result := sioc.Get[TestStruct]()
	if result.Name != "injected" {
		t.Error("Expected to get instance from IInjector")
	}
}

func TestInitShouldCallInitMethod(t *testing.T) {
	sioc.Start()
	testStruct := &TestStructWithInit{}
	sioc.Register(testStruct)
	sioc.Init()
	if !testStruct.Initialized {
		t.Error("Expected Init() to be called")
	}
}

func TestGetWithInterface(t *testing.T) {
	sioc.Start()
	testStruct := &TestStructWithInterface{Name: "interface"}
	sioc.Register(testStruct)
	result := sioc.Get[ITestInterface]()
	if result.GetName() != "interface" {
		t.Error("Expected to get instance implementing interface")
	}
}

func TestGetFunctionName(t *testing.T) {
	testFunc := func() {}
	name := sioc.GetFunctionName(testFunc)
	if name == "" {
		t.Error("Expected non-empty function name")
	}
}

func TestInjectorMatchWithName(t *testing.T) {
	injector := sioc.NewInjector[interface{}]()
	testStruct := &TestStruct{}
	injector.AddInstance(testStruct)
	if !injector.MatchWithName("*sioc_test.TestStruct") {
		t.Error("Expected injector to match with type name")
	}
}

func TestGetInstanceFromInjector(t *testing.T) {
	injector := sioc.NewInjector[interface{}]()
	testStruct := &TestStruct{Name: "test"}
	injector.AddInstance(testStruct)
	result := injector.GetInstance()
	if result.(*TestStruct).Name != "test" {
		t.Error("Expected to get correct instance from injector")
	}
}

type TestStructWithDependency struct {
	dependency  *TestStruct
	initialized bool
}

func (t *TestStructWithDependency) Init(dep *TestStruct) {
	t.dependency = dep
	t.initialized = true
}

type TestStructWithMultipleDeps struct {
	dep1        *TestStruct
	dep2        ITestInterface
	initialized bool
}

func (t *TestStructWithMultipleDeps) Init(dep1 *TestStruct, dep2 ITestInterface) {
	t.dep1 = dep1
	t.dep2 = dep2
	t.initialized = true
}

func TestInitWithNoDependencies(t *testing.T) {
	sioc.Start()
	testStruct := &TestStructWithInit{}
	sioc.Register(testStruct)
	sioc.Init()
	if !testStruct.Initialized {
		t.Error("Expected Init() to be called without dependencies")
	}
}

func TestInitWithOneDependency(t *testing.T) {
	sioc.Start()
	dep := &TestStruct{Name: "dependency"}
	testStruct := &TestStructWithDependency{}

	sioc.Register(dep)
	sioc.Register(testStruct)

	sioc.Init()

	if !testStruct.initialized {
		t.Error("Expected Init() to be called")
	}
	if testStruct.dependency == nil || testStruct.dependency.Name != "dependency" {
		t.Error("Expected dependency to be injected correctly")
	}
}

func TestInitWithMultipleDependencies(t *testing.T) {
	sioc.Start()
	dep1 := &TestStruct{Name: "dep1"}
	dep2 := &TestStructWithInterface{Name: "dep2"}
	testStruct := &TestStructWithMultipleDeps{}

	sioc.Register(dep1)
	sioc.Register(dep2)
	sioc.Register(testStruct)

	sioc.Init()

	if !testStruct.initialized {
		t.Error("Expected Init() to be called")
	}
	if testStruct.dep1 == nil || testStruct.dep1.Name != "dep1" {
		t.Error("Expected dep1 to be injected correctly")
	}
	if testStruct.dep2 == nil || testStruct.dep2.GetName() != "dep2" {
		t.Error("Expected dep2 to be injected correctly")
	}
}

func TestInitOrder(t *testing.T) {
	sioc.Start()
	dep := &TestStructWithInit{}
	testStruct := &TestStructWithInit{}

	sioc.Register(dep)
	sioc.Register(testStruct)

	sioc.Init()

	if !dep.Initialized || !testStruct.Initialized {
		t.Error("Expected both structures to be initialized")
	}
}

type TestStructWithDependencyValue struct {
	Value string
}

type TestStructDependent struct {
	Dep *TestStructWithDependencyValue
}

func (t *TestStructDependent) Init(depNew *TestStructWithDependencyValue) {
	t.Dep = depNew
}

func TestValueChangeReflectsDependency1(t *testing.T) {
	sioc.Start()
	dep := &TestStructWithDependencyValue{Value: "initial-value"}
	dependent := &TestStructDependent{}

	sioc.Register(dep)
	sioc.Register(dependent)
	sioc.Init()

	dep.Value = "new-value"
	if dependent.Dep.Value != "new-value" {
		t.Error("Expected value change to reflect in dependency")
	}
}

func TestValueChangeReflectsDependency2(t *testing.T) {
	sioc.Start()
	dep := &TestStructWithDependencyValue{Value: "test1"}
	dependent := &TestStructDependent{}

	sioc.Register(dep)
	sioc.Register(dependent)
	sioc.Init()

	dep.Value = "test2"
	if dependent.Dep.Value != dep.Value {
		t.Error("Values should be equal after change")
	}
}

func TestValueChangeReflectsDependency3(t *testing.T) {
	os.Setenv("NODE_ENV", "test")
	sioc.Start()
	dep := &TestStructWithDependencyValue{Value: "abc"}
	dependent1 := &TestStructDependent{}
	dependent2 := &TestStructDependent{}

	sioc.Register(dep)
	sioc.Register(dependent1)
	sioc.Register(dependent2)
	sioc.Init()

	dep.Value = "xyz"
	if dependent1.Dep.Value != "xyz" || dependent2.Dep.Value != "xyz" {
		t.Error("Change should reflect in multiple dependencies")
	}
}

func TestValueChangeReflectsDependency4(t *testing.T) {
	sioc.Start()
	dep := &TestStructWithDependencyValue{Value: "initial"}
	dependent := &TestStructDependent{}

	sioc.Register(dep)
	sioc.Register(dependent)
	sioc.Init()

	originalValue := dep.Value
	dep.Value = "modified"
	dep.Value = originalValue

	if dependent.Dep.Value != originalValue {
		t.Error("Value should return to original after multiple changes")
	}
}

func TestValueChangeReflectsDependency5(t *testing.T) {
	sioc.Start()
	dep := &TestStructWithDependencyValue{Value: ""}
	dependent := &TestStructDependent{}

	sioc.Register(dep)
	sioc.Register(dependent)
	sioc.Init()

	values := []string{"test1", "test2", "test3"}
	for _, value := range values {
		dep.Value = value
		if dependent.Dep.Value != value {
			t.Errorf("Expected value %s, got %s", value, dependent.Dep.Value)
		}
	}
}

type TestNewInstance struct {
	Initialized bool
	a           *TestStructWithDependencyValue
}

func (t *TestNewInstance) Init(_ sioc.InitializeNewInstanceTo, a *TestStructWithDependencyValue) {
	t.Initialized = true
	t.a = a
	a.Value = "new value"
}

func TestNewInstanceUsingNewInstanceTo(t *testing.T) {
	os.Setenv("NODE_ENV", "test")
	sioc.Start()
	depA := &TestStructWithDependencyValue{Value: "test"}
	depB := &TestStructDependent{}

	sioc.Register(&TestNewInstance{})
	sioc.Register(depB)
	sioc.Register(depA)
	sioc.Init()

	testModule := sioc.Get[TestNewInstance]()

	if !testModule.Initialized {
		t.Error("Expected module to be initialized")
	}

	if testModule.a.Value != "new value" {
		t.Error("Expected value of a to be changed")
	}
}

type TestNewInstanceTwo struct {
	Initialized bool
	a           *TestStructWithDependencyValue
	b           *TestStructWithDependencyValue
}

func (t *TestNewInstanceTwo) Init(a *TestStructWithDependencyValue, _ sioc.InitializeNewInstanceTo, b *TestStructWithDependencyValue) {
	t.Initialized = true
	t.a = a
	t.b = b
	b.Value = "new value"
}

func TestNewInstanceUsingNewInstanceToB(t *testing.T) {
	os.Setenv("NODE_ENV", "test")
	sioc.Start()
	depA := &TestStructWithDependencyValue{Value: "test"}
	depB := &TestStructWithDependencyValue{Value: "test"}

	sioc.Register(&TestNewInstanceTwo{})
	sioc.Register(depA)
	sioc.Register(depB)
	sioc.Init()

	testModule := sioc.Get[TestNewInstanceTwo]()

	if !testModule.Initialized {
		t.Error("Expected module to be initialized")
	}

	if testModule.b.Value != "new value" {
		t.Error("Expected value of a to be changed")
	}
}

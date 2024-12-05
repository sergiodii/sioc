package offers_platform_core_injection_test

import (
	offers_platform_core_injection "marketing-api/internal/offers_platform/core/injection"
	"os"
	"testing"
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
	offers_platform_core_injection.Start()
	if offers_platform_core_injection.Len() != 0 {
		t.Error("Expected empty injector after Start()")
	}
}

func TestInjectShouldAddInstance(t *testing.T) {
	offers_platform_core_injection.Start()
	testStruct := &TestStruct{Name: "test"}
	offers_platform_core_injection.Inject(testStruct)

	if offers_platform_core_injection.Len() != 1 {
		t.Error("Expected one instance after Inject()")
	}
	offers_platform_core_injection.ClearList()
}

func TestGetShouldReturnInjectedInstance(t *testing.T) {
	offers_platform_core_injection.Start()
	testStruct := &TestStruct{Name: "test"}
	offers_platform_core_injection.Inject(testStruct)
	result := offers_platform_core_injection.Get[*TestStruct]()
	if result.Name != "test" {
		t.Error("Expected to get injected instance")
	}
	offers_platform_core_injection.ClearList()
}

func TestInjectWithInjectorInterface(t *testing.T) {
	offers_platform_core_injection.Start()
	testInjector := &TestInjector{}
	offers_platform_core_injection.Inject(testInjector)

	if offers_platform_core_injection.Len() != 2 {
		t.Error("Expected two instances after injecting IInjector, has: ", offers_platform_core_injection.Len())
	}
	offers_platform_core_injection.ClearList()
}

func TestGetInjectedFromInjector(t *testing.T) {
	offers_platform_core_injection.Start()
	testInjector := &TestInjector{}
	offers_platform_core_injection.Inject(testInjector)
	result := offers_platform_core_injection.Get[TestStruct]()
	if result.Name != "injected" {
		t.Error("Expected to get instance from IInjector")
	}
}

func TestInitShouldCallInitMethod(t *testing.T) {
	offers_platform_core_injection.Start()
	testStruct := &TestStructWithInit{}
	offers_platform_core_injection.Inject(testStruct)
	offers_platform_core_injection.Init()
	if !testStruct.Initialized {
		t.Error("Expected Init() to be called")
	}
}

func TestGetWithInterface(t *testing.T) {
	offers_platform_core_injection.Start()
	testStruct := &TestStructWithInterface{Name: "interface"}
	offers_platform_core_injection.Inject(testStruct)
	result := offers_platform_core_injection.Get[ITestInterface]()
	if result.GetName() != "interface" {
		t.Error("Expected to get instance implementing interface")
	}
}

func TestGetFunctionName(t *testing.T) {
	testFunc := func() {}
	name := offers_platform_core_injection.GetFunctionName(testFunc)
	if name == "" {
		t.Error("Expected non-empty function name")
	}
}

func TestInjectorMatchWithName(t *testing.T) {
	injector := offers_platform_core_injection.NewInjector[interface{}]()
	testStruct := &TestStruct{}
	injector.AddInstance(testStruct)
	if !injector.MatchWithName("*offers_platform_core_injection_test.TestStruct") {
		t.Error("Expected injector to match with type name")
	}
}

func TestGetInstanceFromInjector(t *testing.T) {
	injector := offers_platform_core_injection.NewInjector[interface{}]()
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
	offers_platform_core_injection.Start()
	testStruct := &TestStructWithInit{}
	offers_platform_core_injection.Inject(testStruct)
	offers_platform_core_injection.Init()
	if !testStruct.Initialized {
		t.Error("Expected Init() to be called without dependencies")
	}
}

func TestInitWithOneDependency(t *testing.T) {
	offers_platform_core_injection.Start()
	dep := &TestStruct{Name: "dependency"}
	testStruct := &TestStructWithDependency{}

	offers_platform_core_injection.Inject(dep)
	offers_platform_core_injection.Inject(testStruct)

	offers_platform_core_injection.Init()

	if !testStruct.initialized {
		t.Error("Expected Init() to be called")
	}
	if testStruct.dependency == nil || testStruct.dependency.Name != "dependency" {
		t.Error("Expected dependency to be injected correctly")
	}
}

func TestInitWithMultipleDependencies(t *testing.T) {
	offers_platform_core_injection.Start()
	dep1 := &TestStruct{Name: "dep1"}
	dep2 := &TestStructWithInterface{Name: "dep2"}
	testStruct := &TestStructWithMultipleDeps{}

	offers_platform_core_injection.Inject(dep1)
	offers_platform_core_injection.Inject(dep2)
	offers_platform_core_injection.Inject(testStruct)

	offers_platform_core_injection.Init()

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
	offers_platform_core_injection.Start()
	dep := &TestStructWithInit{}
	testStruct := &TestStructWithInit{}

	offers_platform_core_injection.Inject(dep)
	offers_platform_core_injection.Inject(testStruct)

	offers_platform_core_injection.Init()

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
	offers_platform_core_injection.Start()
	dep := &TestStructWithDependencyValue{Value: "initial-value"}
	dependent := &TestStructDependent{}

	offers_platform_core_injection.Inject(dep)
	offers_platform_core_injection.Inject(dependent)
	offers_platform_core_injection.Init()

	dep.Value = "new-value"
	if dependent.Dep.Value != "new-value" {
		t.Error("Expected value change to reflect in dependency")
	}
}

func TestValueChangeReflectsDependency2(t *testing.T) {
	offers_platform_core_injection.Start()
	dep := &TestStructWithDependencyValue{Value: "test1"}
	dependent := &TestStructDependent{}

	offers_platform_core_injection.Inject(dep)
	offers_platform_core_injection.Inject(dependent)
	offers_platform_core_injection.Init()

	dep.Value = "test2"
	if dependent.Dep.Value != dep.Value {
		t.Error("Values should be equal after change")
	}
}

func TestValueChangeReflectsDependency3(t *testing.T) {
	os.Setenv("NODE_ENV", "test")
	offers_platform_core_injection.Start()
	dep := &TestStructWithDependencyValue{Value: "abc"}
	dependent1 := &TestStructDependent{}
	dependent2 := &TestStructDependent{}

	offers_platform_core_injection.Inject(dep)
	offers_platform_core_injection.Inject(dependent1)
	offers_platform_core_injection.Inject(dependent2)
	offers_platform_core_injection.Init()

	dep.Value = "xyz"
	if dependent1.Dep.Value != "xyz" || dependent2.Dep.Value != "xyz" {
		t.Error("Change should reflect in multiple dependencies")
	}
}

func TestValueChangeReflectsDependency4(t *testing.T) {
	offers_platform_core_injection.Start()
	dep := &TestStructWithDependencyValue{Value: "initial"}
	dependent := &TestStructDependent{}

	offers_platform_core_injection.Inject(dep)
	offers_platform_core_injection.Inject(dependent)
	offers_platform_core_injection.Init()

	originalValue := dep.Value
	dep.Value = "modified"
	dep.Value = originalValue

	if dependent.Dep.Value != originalValue {
		t.Error("Value should return to original after multiple changes")
	}
}

func TestValueChangeReflectsDependency5(t *testing.T) {
	offers_platform_core_injection.Start()
	dep := &TestStructWithDependencyValue{Value: ""}
	dependent := &TestStructDependent{}

	offers_platform_core_injection.Inject(dep)
	offers_platform_core_injection.Inject(dependent)
	offers_platform_core_injection.Init()

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

func (t *TestNewInstance) Init(_ offers_platform_core_injection.InitializeNewInstanceTo, a *TestStructWithDependencyValue) {
	t.Initialized = true
	t.a = a
	a.Value = "new value"
}

func TestNewInstanceUsingNewInstanceTo(t *testing.T) {
	os.Setenv("NODE_ENV", "test")
	offers_platform_core_injection.Start()
	depA := &TestStructWithDependencyValue{Value: "test"}
	depB := &TestStructDependent{}

	offers_platform_core_injection.Inject(&TestNewInstance{})
	offers_platform_core_injection.Inject(depB)
	offers_platform_core_injection.Inject(depA)
	offers_platform_core_injection.Init()

	testModule := offers_platform_core_injection.Get[TestNewInstance]()

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

func (t *TestNewInstanceTwo) Init(a *TestStructWithDependencyValue, _ offers_platform_core_injection.InitializeNewInstanceTo, b *TestStructWithDependencyValue) {
	t.Initialized = true
	t.a = a
	t.b = b
	b.Value = "new value"
}

func TestNewInstanceUsingNewInstanceToB(t *testing.T) {
	os.Setenv("NODE_ENV", "test")
	offers_platform_core_injection.Start()
	depA := &TestStructWithDependencyValue{Value: "test"}
	depB := &TestStructWithDependencyValue{Value: "test"}

	offers_platform_core_injection.Inject(&TestNewInstanceTwo{})
	offers_platform_core_injection.Inject(depA)
	offers_platform_core_injection.Inject(depB)
	offers_platform_core_injection.Init()

	testModule := offers_platform_core_injection.Get[TestNewInstanceTwo]()

	if !testModule.Initialized {
		t.Error("Expected module to be initialized")
	}

	if testModule.b.Value != "new value" {
		t.Error("Expected value of a to be changed")
	}
}

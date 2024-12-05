package v1_injection_tests

import (
	"fmt"
	"testing"

	v1_container "github.com/sergiodii/sioc/v1/container"
	v1_injection "github.com/sergiodii/sioc/v1/injection"
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
	container := v1_container.New()
	if container.Len() != 0 {
		t.Error("Expected empty injector after Start()")
	}
}

func TestInjectShouldAddInstance(t *testing.T) {
	container := v1_container.New()
	testStruct := &TestStruct{Name: "test"}
	v1_injection.Inject(testStruct, container)
	if container.Len() != 1 {
		t.Error("Expected one instance after Inject()")
	}
}

func TestGetShouldReturnInjectedInstance(t *testing.T) {
	container := v1_container.New()
	testStruct := &TestStruct{Name: "test"}
	v1_injection.Inject(testStruct, container)
	result := v1_injection.Get[TestStruct](container)
	if result.Name != "test" {
		t.Error("Expected to get injected instance")
	}
}

func TestInjectWithInjectorInterface(t *testing.T) {
	container := v1_container.New()
	testInjector := &TestInjector{}
	v1_injection.Inject(testInjector, container)
	if container.Len() != 2 {
		t.Error("Expected two instances after injecting IInjector")
	}
}

func TestGetInjectedFromInjector(t *testing.T) {
	container := v1_container.New()
	testInjector := &TestInjector{}
	v1_injection.Inject(testInjector, container)
	result := v1_injection.Get[TestStruct](container)
	if result.Name != "injected" {
		t.Error("Expected to get instance from IInjector")
	}
}

func TestInitShouldCallInitMethod(t *testing.T) {
	container := v1_container.New()
	testStruct := &TestStructWithInit{}
	v1_injection.Inject(testStruct, container)
	v1_injection.Init(container)
	if !testStruct.Initialized {
		t.Error("Expected Init() to be called")
	}
}

func TestGetWithInterface(t *testing.T) {
	container := v1_container.New()
	testStruct := &TestStructWithInterface{Name: "interface"}
	v1_injection.Inject(testStruct, container)
	result := v1_injection.Get[ITestInterface](container)
	if result.GetName() != "interface" {
		t.Error("Expected to get instance implementing interface")
	}
}

func TestMultipleInjects(t *testing.T) {
	container := v1_container.New()
	for i := 0; i < 5; i++ {
		v1_injection.Inject(&TestStruct{Name: fmt.Sprintf("test%d", i)}, container)
	}
	if container.Len() != 5 {
		t.Error("Expected 5 instances after multiple injects")
	}
}

func TestGetFunctionName(t *testing.T) {
	testFunc := func() {}
	name := v1_injection.GetFunctionName(testFunc)
	if name == "" {
		t.Error("Expected non-empty function name")
	}
}

func TestInjectorMatchWithName(t *testing.T) {
	injector := v1_injection.NewInjector[interface{}]()
	testStruct := &TestStruct{}
	injector.AddInstance(testStruct)
	if !injector.MatchWithName("*v1_injection_test.TestStruct") {
		t.Error("Expected injector to match with type name")
	}
}

func TestGetInstanceFromInjector(t *testing.T) {
	injector := v1_injection.NewInjector[interface{}]()
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
	container := v1_container.New()
	testStruct := &TestStructWithInit{}
	v1_injection.Inject(testStruct, container)
	v1_injection.Init(container)
	if !testStruct.Initialized {
		t.Error("Esperava que Init() fosse chamado sem dependências")
	}
}

func TestInitWithOneDependency(t *testing.T) {
	container := v1_container.New()
	dep := &TestStruct{Name: "dependency"}
	testStruct := &TestStructWithDependency{}

	v1_injection.Inject(dep, container)
	v1_injection.Inject(testStruct, container)

	v1_injection.Init(container)

	if !testStruct.initialized {
		t.Error("Esperava que Init() fosse chamado")
	}
	if testStruct.dependency == nil || testStruct.dependency.Name != "dependency" {
		t.Error("Esperava que a dependência fosse injetada corretamente")
	}
}

func TestInitWithMultipleDependencies(t *testing.T) {
	container := v1_container.New()
	dep1 := &TestStruct{Name: "dep1"}
	dep2 := &TestStructWithInterface{Name: "dep2"}
	testStruct := &TestStructWithMultipleDeps{}

	v1_injection.Inject(dep1, container)
	v1_injection.Inject(dep2, container)
	v1_injection.Inject(testStruct, container)

	v1_injection.Init(container)

	if !testStruct.initialized {
		t.Error("Esperava que Init() fosse chamado")
	}
	if testStruct.dep1 == nil || testStruct.dep1.Name != "dep1" {
		t.Error("Esperava que dep1 fosse injetada corretamente")
	}
	if testStruct.dep2 == nil || testStruct.dep2.GetName() != "dep2" {
		t.Error("Esperava que dep2 fosse injetada corretamente")
	}
}

func TestInitOrder(t *testing.T) {
	container := v1_container.New()
	dep := &TestStructWithInit{}
	testStruct := &TestStructWithInit{}

	v1_injection.Inject(dep, container)
	v1_injection.Inject(testStruct, container)

	v1_injection.Init(container)

	if !dep.Initialized || !testStruct.Initialized {
		t.Error("Esperava que ambas as estruturas fossem inicializadas")
	}
}

type TestStructComDependencia struct {
	Valor string
}

type TestStructDependente struct {
	Dep *TestStructComDependencia
}

func (t *TestStructDependente) Init(depNew *TestStructComDependencia) {
	t.Dep = depNew
}

func TestAlteracaoValorRefleteDependencia1(t *testing.T) {
	container := v1_container.New()
	dep := &TestStructComDependencia{Valor: "valor-inicial"}
	dependente := &TestStructDependente{}

	v1_injection.Inject(dep, container)
	v1_injection.Inject(dependente, container)
	v1_injection.Init(container)

	dep.Valor = "novo-valor"
	if dependente.Dep.Valor != "novo-valor" {
		t.Error("Esperava que a alteração do valor refletisse na dependência")
	}
}

func TestAlteracaoValorRefleteDependencia2(t *testing.T) {
	container := v1_container.New()
	dep := &TestStructComDependencia{Valor: "teste1"}
	dependente := &TestStructDependente{}

	v1_injection.Inject(dep, container)
	v1_injection.Inject(dependente, container)
	v1_injection.Init(container)

	dep.Valor = "teste2"
	if dependente.Dep.Valor != dep.Valor {
		t.Error("Valores deveriam ser iguais após alteração")
	}
}

func TestAlteracaoValorRefleteDependencia3(t *testing.T) {
	container := v1_container.New()
	dep := &TestStructComDependencia{Valor: "abc"}
	dependente1 := &TestStructDependente{}
	dependente2 := &TestStructDependente{}

	v1_injection.Inject(dep, container)
	v1_injection.Inject(dependente1, container)
	v1_injection.Inject(dependente2, container)
	v1_injection.Init(container)

	dep.Valor = "xyz"
	if dependente1.Dep.Valor != "xyz" || dependente2.Dep.Valor != "xyz" {
		t.Error("Alteração deveria refletir em múltiplas dependências")
	}
}

func TestAlteracaoValorRefleteDependencia4(t *testing.T) {
	container := v1_container.New()
	dep := &TestStructComDependencia{Valor: "inicial"}
	dependente := &TestStructDependente{}

	v1_injection.Inject(dep, container)
	v1_injection.Inject(dependente, container)
	v1_injection.Init(container)

	valorOriginal := dep.Valor
	dep.Valor = "modificado"
	dep.Valor = valorOriginal

	if dependente.Dep.Valor != valorOriginal {
		t.Error("Valor deveria voltar ao original após múltiplas alterações")
	}
}

func TestAlteracaoValorRefleteDependencia5(t *testing.T) {
	container := v1_container.New()
	dep := &TestStructComDependencia{Valor: ""}
	dependente := &TestStructDependente{}

	v1_injection.Inject(dep, container)
	v1_injection.Inject(dependente, container)
	v1_injection.Init(container)

	valores := []string{"teste1", "teste2", "teste3"}
	for _, valor := range valores {
		dep.Valor = valor
		if dependente.Dep.Valor != valor {
			t.Errorf("Valor esperado %s, obtido %s", valor, dependente.Dep.Valor)
		}
	}
}

type TestNewInstance struct {
	Initialized bool
	a           *TestStructComDependencia
}

func (t *TestNewInstance) Init(_ v1_injection.InitializeNewInstanceTo, a *TestStructComDependencia) {
	t.Initialized = true
	t.a = a
	a.Valor = "novo valor"
}

func TestNewInstanceUsingNewInstanceTo(t *testing.T) {
	container := v1_container.New()
	depA := &TestStructComDependencia{Valor: "teste"}
	depB := &TestStructDependente{}

	v1_injection.Inject(depA, container)
	v1_injection.Inject(depB, container)
	v1_injection.Inject(&TestNewInstance{}, container)
	v1_injection.Init(container)

	testModule := v1_injection.Get[TestNewInstance](container)

	if !testModule.Initialized {
		t.Error("Esperava que o módulo fosse inicializado")
	}

	if testModule.a.Valor != "novo valor" {
		t.Error("Esperava que o valor de a fosse alterado")
	}
}

type TestNewInstanceTwo struct {
	Initialized bool
	a           *TestStructComDependencia
	b           *TestStructComDependencia
}

func (t *TestNewInstanceTwo) Init(a *TestStructComDependencia, _ v1_injection.InitializeNewInstanceTo, b *TestStructComDependencia) {
	t.Initialized = true
	t.a = a
	t.b = b
	b.Valor = "novo valor"
}

func TestNewInstanceUsingNewInstanceToB(t *testing.T) {
	container := v1_container.New()
	depA := &TestStructComDependencia{Valor: "teste"}
	depB := &TestStructComDependencia{Valor: "teste"}

	v1_injection.Inject(depA, container)
	v1_injection.Inject(depB, container)
	v1_injection.Inject(&TestNewInstanceTwo{}, container)
	v1_injection.Init(container)

	testModule := v1_injection.Get[TestNewInstanceTwo](container)

	if !testModule.Initialized {
		t.Error("Esperava que o módulo fosse inicializado")
	}

	if testModule.b.Valor != "novo valor" {
		t.Error("Esperava que o valor de a fosse alterado")
	}
}

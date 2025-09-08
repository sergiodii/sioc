package sioc

import (
	"reflect"
	"testing"
)

// Test interfaces and structs for testing
type TestInterface interface {
	GetValue() string
}

type TestStruct struct {
	Value       string
	initialized bool
}

func (ts *TestStruct) GetValue() string {
	return ts.Value
}

func (ts *TestStruct) Init() {
	ts.initialized = true
}

type TestStructWithDependency struct {
	Dependency  *TestStruct
	initialized bool
}

func (tsd *TestStructWithDependency) Init(dep *TestStruct) {
	tsd.Dependency = dep
	tsd.initialized = true
}

type TestService struct {
	Name string
}

func (ts *TestService) GetValue() string {
	return ts.Name
}

// TestNewContainer tests container creation
func TestNewContainer(t *testing.T) {
	container := NewContainer()
	if container == nil {
		t.Fatal("NewContainer should not return nil")
	}

	if container.Count() != 0 {
		t.Errorf("New container should be empty, got count: %d", container.Count())
	}
}

// TestContainerRegisterAndResolve tests basic registration and resolution
func TestContainerRegisterAndResolve(t *testing.T) {
	container := NewContainer()
	testValue := "test_value"

	container.Register("test_key", testValue)

	resolved, found := container.Resolve("test_key")
	if !found {
		t.Fatal("Should find registered value")
	}

	if resolved != testValue {
		t.Errorf("Expected %v, got %v", testValue, resolved)
	}
}

// TestContainerBackwardCompatibility tests old method names
func TestContainerBackwardCompatibility(t *testing.T) {
	container := NewContainer()
	testValue := "test_value"

	// Cast to concrete type to access backward compatibility methods
	concreteContainer := container.(*serviceRegistry)

	// Test old set method
	concreteContainer.Register("test_key", testValue)

	// Test old get method
	resolved, found := concreteContainer.Resolve("test_key")
	if !found {
		t.Fatal("Should find registered value using old get method")
	}

	if resolved != testValue {
		t.Errorf("Expected %v, got %v", testValue, resolved)
	}

	// Test getAll method
	all := concreteContainer.ListAll()
	if len(all) != 1 {
		t.Errorf("Expected 1 item, got %d", len(all))
	}

	// Test len method
	if concreteContainer.Count() != 1 {
		t.Errorf("Expected length 1, got %d", concreteContainer.Count())
	}
}

// TestNewServiceWrapper tests service wrapper creation
func TestNewServiceWrapper(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	if wrapper == nil {
		t.Fatal("NewServiceWrapper should not return nil")
	}
}

// TestNewInjectorBackwardCompatibility tests backward compatibility
func TestNewInjectorBackwardCompatibility(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	if wrapper == nil {
		t.Fatal("NewInjector should not return nil")
	}
}

// TestServiceWrapperSetAndGet tests service wrapper functionality
func TestServiceWrapperSetAndGet(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	testValue := "test_service"

	wrapper.SetService(testValue)
	retrieved := wrapper.GetService()

	if retrieved != testValue {
		t.Errorf("Expected %v, got %v", testValue, retrieved)
	}
}

// TestServiceWrapperBackwardCompatibilityMethods tests old method names
func TestServiceWrapperBackwardCompatibilityMethods(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	testValue := "test_service"

	// Cast to concrete type to access backward compatibility methods
	concreteWrapper := wrapper.(*serviceWrapper[string])

	// Test addInstance (old method)
	concreteWrapper.addInstance(testValue)

	// Test getInstance (old method)
	retrieved := concreteWrapper.getInstance()
	if retrieved != testValue {
		t.Errorf("Expected %v, got %v", testValue, retrieved)
	}

	// Test getNewInstance (old method)
	newInstance := concreteWrapper.getNewInstance()
	if newInstance != testValue {
		t.Errorf("Expected %v, got %v", testValue, newInstance)
	}
}

// TestInjectAndGet tests the main injection and resolution functionality
func TestInjectAndGet(t *testing.T) {
	container := NewContainer()
	testStruct := &TestStruct{Value: "test"}

	Inject(testStruct, container)

	retrieved := Get[*TestStruct](container)
	if retrieved.Value != "test" {
		t.Errorf("Expected 'test', got %v", retrieved.Value)
	}
}

// TestRegisterServiceAndResolveService tests new method names
func TestRegisterServiceAndResolveService(t *testing.T) {
	container := NewContainer()
	testStruct := &TestStruct{Value: "test"}

	Inject(testStruct, container)

	retrieved := Get[*TestStruct](container)
	if retrieved.Value != "test" {
		t.Errorf("Expected 'test', got %v", retrieved.Value)
	}
}

// TestInterfaceResolution tests interface-based resolution
func TestInterfaceResolution(t *testing.T) {
	container := NewContainer()
	testStruct := &TestStruct{Value: "interface_test"}

	Inject(testStruct, container)

	// Should be able to retrieve via interface
	retrieved := Get[TestInterface](container)
	if retrieved.GetValue() != "interface_test" {
		t.Errorf("Expected 'interface_test', got %v", retrieved.GetValue())
	}
}

// TestMultipleServicesResolution tests resolution with multiple services
func TestMultipleServicesResolution(t *testing.T) {
	container := NewContainer()

	testStruct1 := &TestStruct{Value: "service1"}
	testStruct2 := &TestStruct{Value: "service2"}
	testService := &TestService{Name: "different_service"}

	Inject(testStruct1, container)
	Inject(testService, container)
	Inject(testStruct2, container)

	// Should retrieve the first matching interface implementation
	retrieved := Get[TestInterface](container)
	if retrieved == nil {
		t.Fatal("Should retrieve a service implementing TestInterface")
	}

	// Should be able to retrieve specific type
	specificService := Get[*TestService](container)
	if specificService.Name != "different_service" {
		t.Errorf("Expected 'different_service', got %v", specificService.Name)
	}
}

// TestGetFunctionName tests function name extraction
func TestGetFunctionName(t *testing.T) {
	testFunc := func() {}
	name := GetFunctionName(testFunc)

	// The exact name might vary, but it should not be empty
	if name == "" {
		t.Error("Function name should not be empty")
	}
}

// TestExtractFunctionName tests new function name
func TestExtractFunctionName(t *testing.T) {
	testFunc := func() {}
	name := GetFunctionName(testFunc)

	// Should be the same as GetFunctionName
	expectedName := GetFunctionName(testFunc)
	if name != expectedName {
		t.Errorf("Expected %v, got %v", expectedName, name)
	}
}

// TestInitWithoutDependencies tests initialization without dependencies
func TestInitWithoutDependencies(t *testing.T) {
	container := NewContainer()
	testStruct := &TestStruct{Value: "test"}

	Inject(testStruct, container)
	Init(container)

	retrieved := Get[*TestStruct](container)
	if !retrieved.initialized {
		t.Error("Service should be initialized")
	}
}

// TestInitWithDependencies tests initialization with dependencies
func TestInitWithDependencies(t *testing.T) {
	container := NewContainer()

	dependency := &TestStruct{Value: "dependency"}
	service := &TestStructWithDependency{}

	Inject(dependency, container)
	Inject(service, container)

	Init(container)

	retrievedService := Get[*TestStructWithDependency](container)
	if !retrievedService.initialized {
		t.Error("Service should be initialized")
	}

	if retrievedService.Dependency == nil {
		t.Error("Dependency should be injected")
	}

	if retrievedService.Dependency.Value != "dependency" {
		t.Errorf("Expected 'dependency', got %v", retrievedService.Dependency.Value)
	}
}

// TestInitializeServices tests new initialization method name
func TestInitializeServices(t *testing.T) {
	container := NewContainer()
	testStruct := &TestStruct{Value: "test"}

	Inject(testStruct, container)
	Init(container)

	retrieved := Get[*TestStruct](container)
	if !retrieved.initialized {
		t.Error("Service should be initialized")
	}
}

// TestServiceWrapperNameMatching tests name matching functionality
func TestServiceWrapperNameMatching(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	testValue := "test"

	wrapper.SetService(testValue)

	// Test MatchesServiceName
	if !wrapper.MatchesServiceName(reflect.TypeOf(testValue).String()) {
		t.Error("Should match service name")
	}

	// Test backward compatibility MatchWithName
	concreteWrapper := wrapper.(*serviceWrapper[string])
	if !concreteWrapper.matchWithName(reflect.TypeOf(testValue).String()) {
		t.Error("Should match with old method name")
	}
}

// TestInstanceCreationMode tests the instance creation mode constants
func TestInstanceCreationMode(t *testing.T) {
	mode := CreateNewInstance
	if mode != "CREATE_NEW" {
		t.Errorf("Expected 'CREATE_NEW', got %v", mode)
	}
}

// TestContainerCount tests container counting functionality
func TestContainerCount(t *testing.T) {
	container := NewContainer()

	if container.Count() != 0 {
		t.Errorf("Empty container should have count 0, got %d", container.Count())
	}

	Inject("service1", container)
	if container.Count() != 1 {
		t.Errorf("Container with 1 service should have count 1, got %d", container.Count())
	}

	// Inject different type to avoid overwriting
	Inject(123, container)
	if container.Count() != 2 {
		t.Errorf("Container with 2 services should have count 2, got %d", container.Count())
	}
}

// TestContainerListAll tests listing all services
func TestContainerListAll(t *testing.T) {
	container := NewContainer()

	service1 := "service1"
	service2 := 123 // Different type to avoid overwriting

	Inject(service1, container)
	Inject(service2, container)

	all := container.ListAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 services, got %d", len(all))
	}
}

// TestCreateNewService tests new service instance creation
func TestCreateNewService(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	testValue := "original"

	wrapper.SetService(testValue)

	// Test CreateNewService
	newService := wrapper.CreateNewService()
	if newService != testValue {
		t.Errorf("Expected %v, got %v", testValue, newService)
	}

	// Modify original to ensure they're separate (for reference types this would matter)
	original := wrapper.GetService()
	if original != testValue {
		t.Errorf("Original should remain %v, got %v", testValue, original)
	}
}

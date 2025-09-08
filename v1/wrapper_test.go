package sioc

import (
	"reflect"
	"testing"
)

// TestServiceWrapperInterface tests the ServiceWrapper interface methods
func TestServiceWrapperInterface(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	testService := "test_service"

	// Test SetService
	result := wrapper.SetService(testService)
	if result != wrapper {
		t.Error("SetService should return the same wrapper")
	}

	// Test GetService
	retrieved := wrapper.GetService()
	if retrieved != testService {
		t.Errorf("Expected %v, got %v", testService, retrieved)
	}

	// Test MatchesServiceName
	typeName := reflect.TypeOf(testService).String()
	if !wrapper.MatchesServiceName(typeName) {
		t.Error("Should match service name")
	}

	// Test CreateNewService
	newService := wrapper.CreateNewService()
	if newService != testService {
		t.Errorf("Expected %v, got %v", testService, newService)
	}
}

// TestServiceWrapperWithStruct tests wrapper with struct types
func TestServiceWrapperWithStruct(t *testing.T) {
	type TestStruct struct {
		Name string
		ID   int
	}

	wrapper := NewServiceWrapper[*TestStruct]()
	testStruct := &TestStruct{Name: "test", ID: 42}

	wrapper.SetService(testStruct)

	retrieved := wrapper.GetService()
	if retrieved.Name != "test" || retrieved.ID != 42 {
		t.Errorf("Expected {Name: test, ID: 42}, got %+v", retrieved)
	}

	// Test CreateNewService creates a copy
	newStruct := wrapper.CreateNewService()
	// For pointer types, CreateNewService returns the same pointer value,
	// not a deep copy. This is expected behavior.
	if newStruct.Name != "test" || newStruct.ID != 42 {
		t.Errorf("Expected copy with {Name: test, ID: 42}, got %+v", newStruct)
	}
}

// TestServiceWrapperWithInterface tests wrapper with interface types
func TestServiceWrapperWithInterface(t *testing.T) {
	type Writer interface {
		Write([]byte) (int, error)
	}

	wrapper := NewServiceWrapper[Writer]()

	// This should compile without issues
	if wrapper == nil {
		t.Fatal("Should create wrapper for interface type")
	}
}

// TestServiceWrapperBackwardCompatibilityMethodsDetailed tests all backward compatibility methods
func TestServiceWrapperBackwardCompatibilityMethodsDetailed(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	concreteWrapper := wrapper.(*serviceWrapper[string])
	testService := "test_service"

	// Test addInstance (backward compatibility)
	result := concreteWrapper.addInstance(testService)
	if result != wrapper {
		t.Error("addInstance should return the same wrapper")
	}

	// Test getInstance (backward compatibility)
	retrieved := concreteWrapper.getInstance()
	if retrieved != testService {
		t.Errorf("Expected %v, got %v", testService, retrieved)
	}

	// Test getNewInstance (backward compatibility)
	newInstance := concreteWrapper.getNewInstance()
	if newInstance != testService {
		t.Errorf("Expected %v, got %v", testService, newInstance)
	}

	// Test matchWithName (backward compatibility)
	typeName := reflect.TypeOf(testService).String()
	if !concreteWrapper.matchWithName(typeName) {
		t.Error("Should match with old method name")
	}
}

// TestServiceWrapperNameSanitization tests that service names are properly sanitized
func TestServiceWrapperNameSanitization(t *testing.T) {
	wrapper := NewServiceWrapper[string]()
	concreteWrapper := wrapper.(*serviceWrapper[string])

	testService := "test"
	wrapper.SetService(testService)

	// Test with sanitized name
	typeName := reflect.TypeOf(testService).String()
	sanitized := concreteWrapper.sanitizeServiceName(typeName)

	if !wrapper.MatchesServiceName(sanitized) {
		t.Error("Should match sanitized name")
	}

	// Test that internal name is sanitized
	if concreteWrapper.serviceName != sanitized {
		t.Errorf("Expected internal name to be %v, got %v", sanitized, concreteWrapper.serviceName)
	}
}

// TestServiceWrapperGenericTypes tests wrapper with various generic types
func TestServiceWrapperGenericTypes(t *testing.T) {
	// Test with slice
	sliceWrapper := NewServiceWrapper[[]string]()
	testSlice := []string{"a", "b", "c"}
	sliceWrapper.SetService(testSlice)

	retrieved := sliceWrapper.GetService()
	if len(retrieved) != 3 || retrieved[0] != "a" {
		t.Errorf("Expected slice [a b c], got %v", retrieved)
	}

	// Test with map
	mapWrapper := NewServiceWrapper[map[string]int]()
	testMap := map[string]int{"key": 42}
	mapWrapper.SetService(testMap)

	retrievedMap := mapWrapper.GetService()
	if retrievedMap["key"] != 42 {
		t.Errorf("Expected map[key:42], got %v", retrievedMap)
	}

	// Test with function
	funcWrapper := NewServiceWrapper[func() string]()
	testFunc := func() string { return "hello" }
	funcWrapper.SetService(testFunc)

	retrievedFunc := funcWrapper.GetService()
	if retrievedFunc() != "hello" {
		t.Errorf("Expected function returning 'hello', got %v", retrievedFunc())
	}
}

// TestServiceWrapperZeroValues tests wrapper behavior with zero values
func TestServiceWrapperZeroValues(t *testing.T) {
	wrapper := NewServiceWrapper[string]()

	// Before setting any service, should return zero value
	retrieved := wrapper.GetService()
	if retrieved != "" {
		t.Errorf("Expected empty string (zero value), got %v", retrieved)
	}

	// Test with zero value explicitly set
	wrapper.SetService("")
	retrieved = wrapper.GetService()
	if retrieved != "" {
		t.Errorf("Expected empty string, got %v", retrieved)
	}

	// Test with pointer wrapper
	ptrWrapper := NewServiceWrapper[*string]()
	retrieved2 := ptrWrapper.GetService()
	if retrieved2 != nil {
		t.Errorf("Expected nil pointer, got %v", retrieved2)
	}
}

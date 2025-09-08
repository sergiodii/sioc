package sioc

import (
	"testing"
)

// TestServiceContainerInterface tests the ServiceContainer interface methods
func TestServiceContainerInterface(t *testing.T) {
	container := NewContainer()

	// Test Register and Resolve
	testService := "test_service"
	container.Register("test_key", testService)

	resolved, found := container.Resolve("test_key")
	if !found {
		t.Fatal("Should find registered service")
	}

	if resolved != testService {
		t.Errorf("Expected %v, got %v", testService, resolved)
	}
}

// TestServiceContainerConcurrency tests concurrent access to container
func TestServiceContainerConcurrency(t *testing.T) {
	container := NewContainer()

	// Test concurrent writes and reads
	done := make(chan bool, 2)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			container.Register("key", i)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			container.Resolve("key")
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Should not panic or race
}

// TestServiceContainerEmpty tests empty container behavior
func TestServiceContainerEmpty(t *testing.T) {
	container := NewContainer()

	if container.Count() != 0 {
		t.Errorf("Empty container should have count 0, got %d", container.Count())
	}

	all := container.ListAll()
	if len(all) != 0 {
		t.Errorf("Empty container should return empty list, got %d items", len(all))
	}

	_, found := container.Resolve("nonexistent")
	if found {
		t.Error("Should not find nonexistent key")
	}
}

// TestServiceContainerOverwrite tests overwriting services
func TestServiceContainerOverwrite(t *testing.T) {
	container := NewContainer()

	// Register first service
	container.Register("key", "first")
	if container.Count() != 1 {
		t.Errorf("Expected count 1, got %d", container.Count())
	}

	// Overwrite with second service
	container.Register("key", "second")
	if container.Count() != 1 {
		t.Errorf("Expected count still 1 after overwrite, got %d", container.Count())
	}

	// Should get the second service
	resolved, found := container.Resolve("key")
	if !found {
		t.Fatal("Should find service")
	}

	if resolved != "second" {
		t.Errorf("Expected 'second', got %v", resolved)
	}
}

// TestServiceRegistryBackwardCompatibility tests all backward compatibility methods
func TestServiceRegistryBackwardCompatibility(t *testing.T) {
	container := NewContainer()
	registry := container.(*serviceRegistry)

	// Test set
	registry.Register("key1", "value1")

	// Test get
	value, found := registry.Resolve("key1")
	if !found || value != "value1" {
		t.Errorf("Expected 'value1', got %v (found: %t)", value, found)
	}

	// Test with multiple values
	registry.Register("key2", "value2")

	// Test getAll
	all := registry.ListAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 items, got %d", len(all))
	}

	// Test len
	if registry.Count() != 2 {
		t.Errorf("Expected length 2, got %d", registry.Count())
	}
}

// TestServiceContainerNameSanitization tests that service names are sanitized
func TestServiceContainerNameSanitization(t *testing.T) {
	container := NewContainer()

	// Register with a name that needs sanitization
	container.Register("  Test Key  ", "test_value")

	// Should be able to resolve with the same unsanitized name
	resolved, found := container.Resolve("  Test Key  ")
	if !found {
		t.Fatal("Should find service with sanitized name")
	}

	if resolved != "test_value" {
		t.Errorf("Expected 'test_value', got %v", resolved)
	}
}

package kademlia

import (
	"testing"
)

// TestNewStorage verifies that a new Storage instance is created correctly
func TestNewStorage(t *testing.T) {
	storage := NewStorage()

	if storage == nil {
		t.Fatalf("Expected new Storage instance, got nil")
	}

	if len(storage.data) != 0 {
		t.Errorf("Expected storage to be initialized with empty data, but got %d items", len(storage.data))
	}
}

// TestStoreValue verifies that a value is correctly stored in the Storage
func TestStoreValue(t *testing.T) {
	storage := NewStorage()
	kademliaID := NewRandomKademliaID()
	value := "test-value"

	storage.StoreValue(*kademliaID, value)

	if storedValue, exists := storage.data[*kademliaID]; !exists {
		t.Fatalf("Expected value to be stored, but it was not found")
	} else if storedValue != value {
		t.Errorf("Expected stored value to be %s, but got %s", value, storedValue)
	}
}

// TestGetValue verifies that a value can be retrieved from the Storage
func TestGetValue(t *testing.T) {
	storage := NewStorage()
	kademliaID := NewRandomKademliaID()
	value := "test-value"

	storage.StoreValue(*kademliaID, value)

	retrievedValue, exists := storage.GetValue(*kademliaID)

	if !exists {
		t.Fatalf("Expected value to be found, but it was not")
	}

	if retrievedValue != value {
		t.Errorf("Expected retrieved value to be %s, but got %s", value, retrievedValue)
	}
}

// TestGetValueNotFound verifies that a non-existent value returns false
func TestGetValueNotFound(t *testing.T) {
	storage := NewStorage()
	kademliaID := NewRandomKademliaID()

	_, exists := storage.GetValue(*kademliaID)

	if exists {
		t.Fatalf("Did not expect a value to be found, but one was")
	}
}

// TestDeleteValue verifies that a value can be deleted from the Storage
func TestDeleteValue(t *testing.T) {
	storage := NewStorage()
	kademliaID := NewRandomKademliaID()
	value := "test-value"

	storage.StoreValue(*kademliaID, value)
	storage.DeleteValue(*kademliaID)

	if _, exists := storage.data[*kademliaID]; exists {
		t.Fatalf("Expected value to be deleted, but it still exists")
	}
}

// TestHasValue verifies that the existence of a value is correctly determined
func TestHasValue(t *testing.T) {
	storage := NewStorage()
	kademliaID := NewRandomKademliaID()
	value := "test-value"

	storage.StoreValue(*kademliaID, value)

	if !storage.HasValue(*kademliaID) {
		t.Errorf("Expected HasValue to return true, but it returned false")
	}

	nonExistentID := NewRandomKademliaID()
	if storage.HasValue(*nonExistentID) {
		t.Errorf("Expected HasValue to return false for a non-existent key, but it returned true")
	}
}

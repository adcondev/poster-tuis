package service

import (
	"reflect"
	"testing"
)

func TestNewManager(t *testing.T) {
	variant := Variant{
		ID:           "test-id",
		Family:       "test-family",
		Variant:      "Local",
		RegistryName: "TestService",
		DisplayName:  "Test Service Display",
		ExeName:      "test.exe",
		Binary:       []byte("test binary data"),
	}

	manager := NewManager(variant)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if !reflect.DeepEqual(manager.variant, variant) {
		t.Errorf("NewManager() = %v, want variant %v", manager.variant, variant)
	}
}

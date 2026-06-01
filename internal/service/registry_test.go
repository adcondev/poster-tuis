package service

import (
	"reflect"
	"sort"
	"testing"
)

func TestGetFamilyNames(t *testing.T) {
	// Call the function
	names := GetFamilyNames()

	// Get the registry to compare
	reg := GetServiceRegistry()

	// 1. Verify the length
	if len(names) != len(reg) {
		t.Fatalf("GetFamilyNames() returned %d items, expected %d", len(names), len(reg))
	}

	// 2. Verify all keys from the registry are present
	expectedNames := make([]string, 0, len(reg))
	for name := range reg {
		expectedNames = append(expectedNames, name)
	}

	// 3. Verify it is sorted
	sort.Strings(expectedNames)

	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf("GetFamilyNames() = %v, want %v", names, expectedNames)
	}
}

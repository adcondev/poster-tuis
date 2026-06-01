package service

import (
	"strings"
	"testing"
)

func TestIsValidServiceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Happy paths
		{"Valid alphabetic lowercase", "myservice", true},
		{"Valid alphabetic uppercase", "MYSERVICE", true},
		{"Valid mixed case", "MyService", true},
		{"Valid with numbers", "MyService123", true},
		{"Valid with underscores", "My_Service_123", true},
		{"Valid with hyphens", "my-service-123", true},

		// Edge cases
		{"Empty string", "", false},
		{"Single character", "a", true},
		{"Maximum length (256)", strings.Repeat("a", 256), true},
		{"Too long (257)", strings.Repeat("a", 257), false},

		// Invalid characters
		{"Contains space", "my service", false},
		{"Contains dot", "my.service", false},
		{"Contains slash", "my/service", false},
		{"Contains backslash", `my\service`, false},
		{"Contains colon", "my:service", false},
		{"Contains path traversal", "../myservice", false},
		{"Contains special chars", "my@service", false},
		{"Contains exclamation", "myservice!", false},
		{"Contains unicode", "myservicé", false},
		{"Starts with space", " myservice", false},
		{"Ends with space", "myservice ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidServiceName(tt.input)
			if result != tt.expected {
				t.Errorf("isValidServiceName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

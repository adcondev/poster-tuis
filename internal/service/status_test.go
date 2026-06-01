package service

import (
	"testing"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{
			name:     "StatusStopPending",
			status:   StatusStopPending,
			expected: "[~] DETENIÉNDOSE...",
		},
		{
			name:     "StatusStartPending",
			status:   StatusStartPending,
			expected: "[~] INICIÁNDOSE...",
		},
		{
			name:     "StatusRunning",
			status:   StatusRunning,
			expected: "[+] EN EJECUCIÓN",
		},
		{
			name:     "StatusStopped",
			status:   StatusStopped,
			expected: "[.] DETENIDO",
		},
		{
			name:     "StatusNotInstalled",
			status:   StatusNotInstalled,
			expected: "[-] NO INSTALADO",
		},
		{
			name:     "StatusUnknown",
			status:   StatusUnknown,
			expected: statusUnknownString,
		},
		{
			name:     "Undefined status (fallback to default)",
			status:   Status(999),
			expected: statusUnknownString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Status.String() for %s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

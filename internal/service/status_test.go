package service

import (
	"testing"
)

func TestFamilyStatus_GetActiveStatus(t *testing.T) {
	tests := []struct {
		name         string
		localStatus  Status
		remoteStatus Status
		expected     Status
	}{
		{
			name:         "Neither installed",
			localStatus:  StatusNotInstalled,
			remoteStatus: StatusNotInstalled,
			expected:     StatusNotInstalled,
		},
		{
			name:         "Local is running",
			localStatus:  StatusRunning,
			remoteStatus: StatusNotInstalled,
			expected:     StatusRunning,
		},
		{
			name:         "Local is stopped",
			localStatus:  StatusStopped,
			remoteStatus: StatusNotInstalled,
			expected:     StatusStopped,
		},
		{
			name:         "Remote is running",
			localStatus:  StatusNotInstalled,
			remoteStatus: StatusRunning,
			expected:     StatusRunning,
		},
		{
			name:         "Remote is stopped",
			localStatus:  StatusNotInstalled,
			remoteStatus: StatusStopped,
			expected:     StatusStopped,
		},
		{
			name:         "Both installed - Local takes priority",
			localStatus:  StatusRunning,
			remoteStatus: StatusStopped,
			expected:     StatusRunning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FamilyStatus{
				LocalStatus:  tt.localStatus,
				RemoteStatus: tt.remoteStatus,
			}
			got := fs.GetActiveStatus()
			if got != tt.expected {
				t.Errorf("FamilyStatus.GetActiveStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

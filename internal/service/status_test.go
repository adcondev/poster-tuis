package service

import "testing"

func TestFamilyStatus_GetInstalledVariant(t *testing.T) {
	tests := []struct {
		name         string
		localStatus  Status
		remoteStatus Status
		want         string
	}{
		{
			name:         "Neither installed",
			localStatus:  StatusNotInstalled,
			remoteStatus: StatusNotInstalled,
			want:         "",
		},
		{
			name:         "Only Local installed (running)",
			localStatus:  StatusRunning,
			remoteStatus: StatusNotInstalled,
			want:         Local,
		},
		{
			name:         "Only Local installed (stopped)",
			localStatus:  StatusStopped,
			remoteStatus: StatusNotInstalled,
			want:         Local,
		},
		{
			name:         "Only Remote installed (running)",
			localStatus:  StatusNotInstalled,
			remoteStatus: StatusRunning,
			want:         Remoto,
		},
		{
			name:         "Only Remote installed (stopped)",
			localStatus:  StatusNotInstalled,
			remoteStatus: StatusStopped,
			want:         Remoto,
		},
		{
			name:         "Both installed (Conflict)",
			localStatus:  StatusStopped,
			remoteStatus: StatusRunning,
			want:         "Conflict",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FamilyStatus{
				LocalStatus:  tt.localStatus,
				RemoteStatus: tt.remoteStatus,
			}
			if got := fs.GetInstalledVariant(); got != tt.want {
				t.Errorf("FamilyStatus.GetInstalledVariant() = %v, want %v", got, tt.want)
			}
		})
	}
}

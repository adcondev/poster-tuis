package service

import (
	"strings"
	"time"
)

// ══════════════════════════════════════════════════════════════
// Status Types
// ══════════════════════════════════════════════════════════════

const (
	// Local is the identifier for the "Local" variant of a service family
	Local = "Local"
	// Remoto is the identifier for the "Remoto" variant of a service family
	Remoto = "Remoto"
)

// Status represents the current state of a Windows service
type Status int

const (
	// StatusNotInstalled indicates the service is not present in the system
	StatusNotInstalled Status = iota
	// StatusStopped indicates the service is installed but not running
	StatusStopped
	// StatusRunning indicates the service is currently running
	StatusRunning
	// StatusStopPending indicates the service is in the process of stopping
	StatusStopPending
	// StatusStartPending indicates the service is in the process of starting
	StatusStartPending
	// StatusUnknown indicates an unrecognized or error state when querying the service
	StatusUnknown
)

// String returns a formatted status string for display
func (s Status) String() string {
	switch s {
	case StatusStopPending:
		return "[~] DETENIÉNDOSE..."
	case StatusStartPending:
		return "[~] INICIÁNDOSE..."
	case StatusRunning:
		return "[+] EN EJECUCIÓN"
	case StatusStopped:
		return "[.] DETENIDO"
	case StatusNotInstalled:
		return "[-] NO INSTALADO"
	default:
		return "[?] ESTADO DESCONOCIDO"
	}
}

// ══════════════════════════════════════════════════════════════
// Family Status (Mutual Exclusivity Tracking)
// ══════════════════════════════════════════════════════════════

// FamilyStatus represents the combined state of both variants in a family.
// This is the KEY structure for enforcing mutual exclusivity:
// only one variant (Local OR Remoto) can be installed at a time.
type FamilyStatus struct {
	LocalStatus  Status
	RemoteStatus Status
}

// GetInstalledVariant returns which variant is currently installed.
// Returns: "Local", "Remoto", or "" if neither is installed.
// CRITICAL: Used by the UI to determine which menu options to show.
func (fs FamilyStatus) GetInstalledVariant() string {
	localInstalled := fs.LocalStatus != StatusNotInstalled
	remoteInstalled := fs.RemoteStatus != StatusNotInstalled

	if localInstalled && remoteInstalled {
		// Both installed — shouldn't happen. Return "Local" but this is a conflict.
		return "Conflict" // or return "CONFLICT" for special handling
	}
	if localInstalled {
		return Local
	}
	if remoteInstalled {
		return Remoto
	}
	return ""
}

// GetActiveStatus returns the status of the currently installed variant.
// Returns StatusNotInstalled if no variant is installed.
func (fs FamilyStatus) GetActiveStatus() Status {
	if fs.LocalStatus != StatusNotInstalled {
		return fs.LocalStatus
	}
	if fs.RemoteStatus != StatusNotInstalled {
		return fs.RemoteStatus
	}
	return StatusNotInstalled
}

// ══════════════════════════════════════════════════════════════
// Status Checking
// ══════════════════════════════════════════════════════════════

// CheckStatus queries the Windows service control manager for the
// current state of this manager's service variant.
func (m *Manager) CheckStatus() Status {
	output, err := secureScRun("query", m.variant.RegistryName)
	outStr := string(output)

	if err != nil {
		// If sc indicates the service doesn't exist, treat as not installed.
		if strings.Contains(outStr, "1060") ||
			strings.Contains(strings.ToLower(outStr), "does not exist") ||
			strings.Contains(strings.ToLower(outStr), "not found") {
			return StatusNotInstalled
		}
		// Unexpected error from sc; return unknown so callers can handle/log if needed.
		return StatusUnknown
	}

	switch {
	case strings.Contains(outStr, "RUNNING"):
		return StatusRunning
	case strings.Contains(outStr, "STOPPED"):
		return StatusStopped
	case strings.Contains(outStr, "STOP_PENDING"):
		return StatusStopPending
	case strings.Contains(outStr, "START_PENDING"):
		return StatusStartPending
	default:
		return StatusUnknown
	}
}

// CheckFamilyStatus checks the status of both variants in a family
// and returns their combined status. This is used for mutual exclusivity
// enforcement in the UI layer.
func CheckFamilyStatus(variants []Variant) FamilyStatus {
	fs := FamilyStatus{
		LocalStatus:  StatusNotInstalled,
		RemoteStatus: StatusNotInstalled,
	}

	for _, v := range variants {
		mgr := NewManager(v)
		status := mgr.CheckStatus()

		switch v.Variant {
		case Local:
			fs.LocalStatus = status
		case Remoto:
			fs.RemoteStatus = status
		}
	}

	return fs
}

// WaitForStatus polls the service status until it reaches the expected state
// or times out. Returns true if the expected status was reached.
func (m *Manager) WaitForStatus(expectedStatus Status, timeout time.Duration) bool {
	// Poll less aggressively — 500ms is plenty for service state changes
	const pollInterval = 500 * time.Millisecond

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		currentStatus := m.CheckStatus()
		if currentStatus == expectedStatus {
			return true
		}
		// If status is not a transitional state and not the expected one, bail early
		if currentStatus != StatusStopPending && currentStatus != StatusStartPending &&
			currentStatus != expectedStatus {
			// e.g., waiting for STOPPED but service is NOT_INSTALLED — that's fine too
			if expectedStatus == StatusStopped && currentStatus == StatusNotInstalled {
				return true
			}
		}
		time.Sleep(pollInterval)
	}
	return false
}

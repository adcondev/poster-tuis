package service

import (
	"os/exec"
	"strings"
	"time"
)

// ══════════════════════════════════════════════════════════════
// Status Types
// ══════════════════════════════════════════════════════════════

// Status represents the current state of a Windows service
type Status int

const (
	StatusNotInstalled Status = iota
	StatusStopped
	StatusRunning
	StatusStopPending
	StatusStartPending
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
// only one variant (Local OR Remote) can be installed at a time.
type FamilyStatus struct {
	LocalStatus  Status
	RemoteStatus Status
}

// GetInstalledVariant returns which variant is currently installed.
// Returns: "Local", "Remote", or "" if neither is installed.
// CRITICAL: Used by the UI to determine which menu options to show.
func (fs FamilyStatus) GetInstalledVariant() string {
	localInstalled := fs.LocalStatus != StatusNotInstalled
	remoteInstalled := fs.RemoteStatus != StatusNotInstalled

	if localInstalled && remoteInstalled {
		// Both installed — shouldn't happen. Return "Local" but this is a conflict.
		return "Local" // or return "CONFLICT" for special handling
	}
	if localInstalled {
		return "Local"
	}
	if remoteInstalled {
		return "Remote"
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
	cmd := exec.Command("sc", "query", m.variant.RegistryName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return StatusNotInstalled
	}

	outputStr := string(output)
	switch {
	case strings.Contains(outputStr, "RUNNING"):
		return StatusRunning
	case strings.Contains(outputStr, "STOPPED"):
		return StatusStopped
	case strings.Contains(outputStr, "STOP_PENDING"):
		return StatusStopPending
	case strings.Contains(outputStr, "START_PENDING"):
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
		case "Local":
			fs.LocalStatus = status
		case "Remote":
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

package service

import (
	"os/exec"
	"strings"
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
	StatusUnknown
)

// String returns a formatted status string for display
func (s Status) String() string {
	switch s {
	case StatusRunning:
		return "[+] EN EJECUCIÓN"
	case StatusStopped:
		return "[.] DETENIDO"
	case StatusNotInstalled:
		return "[-] NO INSTALADO"
	default:
		return "[?] DESCONOCIDO"
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
	if fs.LocalStatus != StatusNotInstalled {
		return "Local"
	}
	if fs.RemoteStatus != StatusNotInstalled {
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
	default:
		return StatusUnknown
	}
}

// CheckFamilyStatus checks the status of both variants in a family
// and returns their combined status. This is used for mutual exclusivity
// enforcement in the UI layer.
func CheckFamilyStatus(variants []ServiceVariant) FamilyStatus {
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

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ══════════════════════════════════════════════════════════════
// Log File Operations
// ══════════════════════════════════════════════════════════════

// GetLogPath returns the full path to the service's log file.
// Pattern: %PROGRAMDATA%\{RegistryName}\{RegistryName}.log
func (m *Manager) GetLogPath() string {
	return filepath.Join(
		os.Getenv("PROGRAMDATA"),
		m.variant.RegistryName,
		m.variant.RegistryName+".log",
	)
}

// GetLogDir returns the directory containing log files.
// Pattern: %PROGRAMDATA%\{RegistryName}
func (m *Manager) GetLogDir() string {
	return filepath.Join(os.Getenv("PROGRAMDATA"), m.variant.RegistryName)
}

func secureLaunch(exe, target, allowedBase string) error {
	// Ensure target exists
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("target does not exist: %s", target)
	} else if err != nil {
		return fmt.Errorf("stat failed: %w", err)
	}

	absTarget, err := filepath.Abs(filepath.Clean(target))
	if err != nil {
		return fmt.Errorf("unable to resolve target path: %w", err)
	}

	absBase, err := filepath.Abs(filepath.Clean(allowedBase))
	if err != nil {
		return fmt.Errorf("unable to resolve base path: %w", err)
	}

	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return fmt.Errorf("unable to evaluate path relation: %w", err)
	}
	if strings.HasPrefix(rel, "..") || rel == "." && absTarget != absBase {
		return fmt.Errorf("path %s is outside allowed base %s", absTarget, absBase)
	}

	exePath, err := exec.LookPath(exe)
	if err != nil {
		return fmt.Errorf("executable not found: %w", err)
	}

	// SECURE: Validate and launch with absolute paths to prevent command injection
	// Eliminamos el context porque Notepad/Explorer deben vivir de forma independiente.
	//nolint:gosec,noctx // inputs validated, lookpath used, detached GUI app explicitly needs no context
	cmd := exec.Command(exePath, absTarget)
	return cmd.Start()
}

// OpenLogFile opens the log file in Notepad
func (m *Manager) OpenLogFile() error {
	logPath := m.GetLogPath()
	// validate and launch
	return secureLaunch("notepad.exe", logPath, m.GetLogDir())
}

// OpenLogDir opens the log directory in Windows Explorer
func (m *Manager) OpenLogDir() error {
	logDir := m.GetLogDir()
	// validate and launch
	return secureLaunch("explorer", logDir, logDir)
}

package service

import (
	"os"
	"os/exec"
	"path/filepath"
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

// OpenLogFile opens the log file in Notepad
func (m *Manager) OpenLogFile() error {
	return exec.Command("notepad.exe", m.GetLogPath()).Start()
}

// OpenLogDir opens the log directory in Windows Explorer
func (m *Manager) OpenLogDir() error {
	return exec.Command("explorer", m.GetLogDir()).Start()
}

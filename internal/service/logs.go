package service

import (
	"fmt"
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
	logPath := m.GetLogPath()
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo de logs no existe: %s", logPath)
	}
	return exec.Command("notepad.exe", logPath).Start()
}

// OpenLogDir opens the log directory in Windows Explorer
func (m *Manager) OpenLogDir() error {
	logDir := m.GetLogDir()
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return fmt.Errorf("la carpeta de logs no existe: %s", logDir)
	}
	return exec.Command("explorer", logDir).Start()
}

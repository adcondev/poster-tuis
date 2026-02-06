package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ══════════════════════════════════════════════════════════════
// Service Manager
// ══════════════════════════════════════════════════════════════

// Manager handles Windows service lifecycle operations for a specific variant
type Manager struct {
	variant ServiceVariant
}

// NewManager creates a manager for a specific service variant
func NewManager(variant ServiceVariant) *Manager {
	return &Manager{variant: variant}
}

// ══════════════════════════════════════════════════════════════
// Install / Uninstall
// ══════════════════════════════════════════════════════════════

// Install creates the Windows service: writes the embedded binary to disk
// and registers it with the service control manager.
func (m *Manager) Install() error {
	targetDir := filepath.Join(os.Getenv("ProgramFiles"), m.variant.RegistryName)
	targetPath := filepath.Join(targetDir, m.variant.ExeName)

	// 1. Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("crear directorio: %w", err)
	}

	// 2. Write embedded binary to disk
	if err := os.WriteFile(targetPath, m.variant.Binary, 0755); err != nil {
		return fmt.Errorf("extraer binario: %w", err)
	}

	// 3. Register service with sc.exe
	// SECURE: Quoted binPath to prevent Unquoted Service Path vulnerability
	cmd := exec.Command("sc", "create", m.variant.RegistryName,
		fmt.Sprintf("binPath=\"%s\"", targetPath),
		"start=auto",
		fmt.Sprintf("DisplayName=%s", m.variant.DisplayName))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("sc create: %s", strings.TrimSpace(string(output)))
	}

	// 4. Configure failure recovery (restart on failure)
	_ = exec.Command("sc", "failure", m.variant.RegistryName,
		"reset=86400",
		"actions=restart/5000/restart/5000/restart/5000").Run()

	return nil
}

// Uninstall stops the service, removes it from the registry,
// and deletes the binary files from disk.
func (m *Manager) Uninstall() error {
	// Stop first (ignore errors — might not be running)
	_ = exec.Command("sc", "stop", m.variant.RegistryName).Run()
	time.Sleep(2 * time.Second)

	// Delete service from registry
	cmd := exec.Command("sc", "delete", m.variant.RegistryName)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "1060") {
			return fmt.Errorf("servicio no instalado")
		}
		return fmt.Errorf("sc delete: %s", strings.TrimSpace(outputStr))
	}

	// Remove binary files from disk
	targetDir := filepath.Join(os.Getenv("ProgramFiles"), m.variant.RegistryName)
	return os.RemoveAll(targetDir)
}

// ══════════════════════════════════════════════════════════════
// Start / Stop / Restart
// ══════════════════════════════════════════════════════════════

// Start starts the Windows service
func (m *Manager) Start() error {
	cmd := exec.Command("sc", "start", m.variant.RegistryName)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "1056") {
			return fmt.Errorf("el servicio ya está en ejecución")
		} else if strings.Contains(outputStr, "1060") {
			return fmt.Errorf("el servicio no está instalado")
		}
		return fmt.Errorf("sc start: %s", strings.TrimSpace(outputStr))
	}
	return nil
}

// Stop stops the Windows service
func (m *Manager) Stop() error {
	cmd := exec.Command("sc", "stop", m.variant.RegistryName)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "1062") {
			return fmt.Errorf("el servicio no está en ejecución")
		}
		return fmt.Errorf("sc stop: %s", strings.TrimSpace(outputStr))
	}
	return nil
}

// Restart stops and starts the service with a delay between operations
func (m *Manager) Restart() error {
	// Stop — ignore "not running" errors
	if err := m.Stop(); err != nil {
		if !strings.Contains(err.Error(), "1062") {
			return err
		}
	}
	time.Sleep(2 * time.Second)
	return m.Start()
}

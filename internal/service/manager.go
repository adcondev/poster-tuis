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
	variant Variant
}

// NewManager creates a manager for a specific service variant
func NewManager(variant Variant) *Manager {
	return &Manager{variant: variant}
}

// validateServiceVariantFields checks that ServiceVariant fields are safe
// and contain only expected characters to prevent command injection
func validateServiceVariantFields(variant Variant) error {
	// Check RegistryName
	if !isValidServiceName(variant.RegistryName) {
		return fmt.Errorf("invalid RegistryName: contains unsafe characters")
	}
	// Check DisplayName
	if !isValidDisplayName(variant.DisplayName) {
		return fmt.Errorf("invalid DisplayName: contains unsafe characters")
	}
	// Check ExeName
	if !isValidFileName(variant.ExeName) {
		return fmt.Errorf("invalid ExeName: contains unsafe characters")
	}
	return nil
}

// isValidServiceName validates that a service name contains only alphanumeric, underscores, and hyphens
func isValidServiceName(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

// isValidDisplayName validates that a display name doesn't contain dangerous characters
func isValidDisplayName(name string) bool {
	if name == "" {
		return false
	}
	// Allow alphanumeric, spaces, and common punctuation, but block shell metacharacters
	for _, c := range name {
		if c == '"' || c == '\'' || c == '`' || c == '$' || c == '&' || c == '|' || c == ';' || c == '\n' || c == '\r' {
			return false
		}
	}
	return true
}

// isValidFileName validates that a file name is safe
func isValidFileName(name string) bool {
	if name == "" || hasPathTraversal(name) {
		return false
	}
	for _, c := range name {
		if c == '"' || c == '\'' || c == '`' || c == '$' || c == '&' || c == '|' || c == ';' || c == '\n' || c == '\r' {
			return false
		}
	}
	return true
}

// hasPathTraversal checks if a string contains path traversal sequences
func hasPathTraversal(s string) bool {
	return strings.Contains(s, "..") || strings.Contains(s, "/") || strings.Contains(s, "\\")
}

// ══════════════════════════════════════════════════════════════
// Install / Uninstall
// ══════════════════════════════════════════════════════════════

// Install creates the Windows service: writes the embedded binary to disk
// and registers it with the service control manager.
func (m *Manager) Install() error {
	// Validate ServiceVariant fields before proceeding
	if err := validateServiceVariantFields(m.variant); err != nil {
		return fmt.Errorf("validación de campos: %w", err)
	}

	// Pre-check: fail fast if already registered
	currentStatus := m.CheckStatus()
	if currentStatus != StatusNotInstalled {
		return fmt.Errorf("el servicio ya está registrado (estado: %s) — desinstale primero", currentStatus)
	}

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
		fmt.Sprintf("DisplayName=\"%s\"", m.variant.DisplayName))

	// 4. More descriptive on common failures
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := strings.TrimSpace(string(output))
		_ = os.RemoveAll(targetDir)
		if strings.Contains(outputStr, "1073") {
			return fmt.Errorf("el servicio ya existe en el registro de Windows (use Desinstalar primero)")
		}
		return fmt.Errorf("no se pudo registrar el servicio: %s", outputStr)
	}

	// 5. Configure failure recovery (restart on failure)
	_ = exec.Command("sc", "failure", m.variant.RegistryName,
		"reset=86400",
		"actions=restart/5000/restart/5000/restart/5000").Run()

	return nil
}

// Uninstall stops the service, removes it from the registry,
// and deletes the binary files from disk.
func (m *Manager) Uninstall() error {
	// Step 1: Attempt to stop the service
	stopErr := exec.Command("sc", "stop", m.variant.RegistryName).Run()

	// Step 2: Wait for STOPPED state with proper timeout
	if stopErr == nil {
		stopped := m.WaitForStatus(StatusStopped, 15*time.Second)
		if !stopped {
			// Force-kill the service process as a last resort
			_ = exec.Command("taskkill", "/F", "/FI",
				fmt.Sprintf("SERVICES eq %s", m.variant.RegistryName)).Run()
			// Wait again briefly after force-kill
			m.WaitForStatus(StatusStopped, 5*time.Second)
		}
	}

	// Step 3: Delete service from registry
	cmd := exec.Command("sc", "delete", m.variant.RegistryName)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "1060") {
			return fmt.Errorf("el servicio no está instalado")
		}
		if strings.Contains(outputStr, "1072") {
			// Service marked for deletion — will complete after process exits
			// This is not a hard failure; inform the user
			return fmt.Errorf("servicio marcado para eliminación (se completará al cerrar el proceso)")
		}
		return fmt.Errorf("sc delete: %s", strings.TrimSpace(outputStr))
	}

	// Step 4: Remove binary files from disk
	targetDir := filepath.Join(os.Getenv("ProgramFiles"), m.variant.RegistryName)
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("no se pudieron eliminar los archivos: %w (puede que el proceso aún esté activo)", err)
	}

	return nil
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
			return fmt.Errorf("el servicio '%s' no está en ejecución", m.variant.DisplayName)
		}
		return fmt.Errorf("no se pudo detener '%s': %s", m.variant.DisplayName, strings.TrimSpace(outputStr))
	}
	return nil
}

func (m *Manager) Restart() error {
	currentStatus := m.CheckStatus()

	// Only try to stop if actually running or in a running-like state
	if currentStatus == StatusRunning || currentStatus == StatusStartPending {
		if err := m.Stop(); err != nil {
			return fmt.Errorf("no se pudo detener el servicio para reiniciar: %w", err)
		}
	}

	if !m.WaitForStatus(StatusStopped, 15*time.Second) {
		return fmt.Errorf("el servicio no se detuvo a tiempo para reiniciar")
	}

	return m.Start()
}

// Package service provides a secure implementation of Windows service management for the embedded agent.
package service

import (
	"context"
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
		switch {
		case c >= 'a' && c <= 'z':
			// Is lowercase, is valid
		case c >= 'A' && c <= 'Z':
			// Is uppercase, is valid
		case c >= '0' && c <= '9':
			// Is digit, is valid
		case c == '_', c == '-':
			// Is underscore or hyphen, is valid
		default:
			// Si no entra en ninguna de las categorías anteriores, es inválido
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
	// Bloqueamos barras y path traversal solo para el nombre de archivo (ExeName)
	if name == "" || hasPathTraversal(name) || strings.ContainsAny(name, `/\`) {
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
	// Solo comprobamos secuencias de escape de directorio
	return strings.Contains(s, "..")
}

// ══════════════════════════════════════════════════════════════
// Install / Uninstall
// ══════════════════════════════════════════════════════════════

// secureScRun runs an sc action (start/stop/delete/etc.) with a validated service name.
func secureScRun(action, regName string, extraArgs ...string) ([]byte, error) {
	if !isValidServiceName(regName) {
		return nil, fmt.Errorf("invalid RegistryName")
	}
	scPath, err := exec.LookPath("sc")
	if err != nil {
		return nil, fmt.Errorf("sc executable not found: %w", err)
	}
	args := append([]string{action, regName}, extraArgs...)

	// Use a context with timeout to avoid hanging subprocesses
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//nolint:gosec // inputs validated and sc resolved via LookPath
	cmd := exec.CommandContext(ctx, scPath, args...)
	return cmd.CombinedOutput()
}

// secureTaskKillByService runs taskkill to force kill processes by service filter.
// Validates service name and resolves executable via LookPath.
func secureTaskKillByService(regName string) error {
	if !isValidServiceName(regName) {
		return fmt.Errorf("invalid RegistryName")
	}
	taskkillPath, err := exec.LookPath("taskkill")
	if err != nil {
		return fmt.Errorf("taskkill not found: %w", err)
	}
	filter := fmt.Sprintf("SERVICES eq %s", regName)
	args := []string{"/F", "/FI", filter}

	// Use a short timeout for taskkill
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//nolint:gosec // inputs validated and taskkill resolved via LookPath
	cmd := exec.CommandContext(ctx, taskkillPath, args...)
	return cmd.Run()
}

// resolveAndEnsure resolves and cleans both base and target and ensures
// target is inside base. Returns absolute base and target on success.
func resolveAndEnsure(base, target string) (absBase, absTarget string, err error) {
	absBase, err = filepath.Abs(filepath.Clean(base))
	if err != nil {
		err = fmt.Errorf("resolve base: %w", err)
		return
	}
	absTarget, err = filepath.Abs(filepath.Clean(target))
	if err != nil {
		err = fmt.Errorf("resolve target: %w", err)
		return
	}
	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		err = fmt.Errorf("evaluate relation: %w", err)
		return
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		err = fmt.Errorf("path %s is outside allowed base %s", absTarget, absBase)
		return
	}
	return
}

// secureScFailure validates the service name and runs the failure config
func secureScFailure(regName string) error {
	if !isValidServiceName(regName) {
		return fmt.Errorf("invalid RegistryName")
	}

	scPath, err := exec.LookPath("sc")
	if err != nil {
		return fmt.Errorf("sc executable not found: %w", err)
	}

	// Separamos llaves de valores
	args := []string{
		"failure",
		regName,
		"reset=", "86400",
		"actions=", "restart/5000/restart/5000/restart/5000",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//nolint:gosec // inputs validated
	cmd := exec.CommandContext(ctx, scPath, args...)
	return cmd.Run()
}

// secureScCreate validates inputs and runs `sc create` safely without cmd.exe
func secureScCreate(regName, binPath, displayName string) ([]byte, error) {
	if !isValidServiceName(regName) {
		return nil, fmt.Errorf("invalid RegistryName")
	}
	if !isValidDisplayName(displayName) {
		return nil, fmt.Errorf("invalid DisplayName")
	}
	if hasPathTraversal(binPath) {
		return nil, fmt.Errorf("invalid binPath")
	}

	scPath, err := exec.LookPath("sc")
	if err != nil {
		return nil, fmt.Errorf("sc executable not found: %w", err)
	}

	// ¡EL TRUCO DE ORO! Separamos la llave del valor por comas.
	// Go generará internamente: sc create RegName binPath= "C:\Ruta\..." start= auto DisplayName= "Nombre..."
	args := []string{
		"create",
		regName,
		"binPath=", binPath,
		"start=", "auto",
		"DisplayName=", displayName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//nolint:gosec // inputs are strictly validated above, safe from command injection
	cmd := exec.CommandContext(ctx, scPath, args...)
	return cmd.CombinedOutput()
}

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

	// Prepare safe absolute paths and validate file name
	programFiles := os.Getenv("ProgramFiles")
	if programFiles == "" {
		return fmt.Errorf("environment variable `ProgramFiles` is empty")
	}

	targetDir := filepath.Join(programFiles, m.variant.RegistryName)
	targetPath := filepath.Join(targetDir, m.variant.ExeName)

	// Ensure ExeName doesn't contain path separators (extra safety)
	if strings.ContainsAny(m.variant.ExeName, `\/`) {
		return fmt.Errorf("invalid ExeName: contains path separator")
	}

	// Resolve and ensure `targetDir` is inside `ProgramFiles`
	_, absTargetDir, err := resolveAndEnsure(programFiles, targetDir)
	if err != nil {
		return fmt.Errorf("invalid target directory: %w", err)
	}

	// Resolve and ensure `targetPath` is inside the resolved `targetDir`
	_, absTargetPath, err := resolveAndEnsure(absTargetDir, targetPath)
	if err != nil {
		return fmt.Errorf("invalid target path: %w", err)
	}

	// 1. Create target directory (using validated absolute path)
	//nolint:gosec // We have validated the path, so this is not vulnerable to injection
	if err := os.MkdirAll(absTargetDir, 0750); err != nil {
		return fmt.Errorf("crear directorio: %w", err)
	}

	// 2. Write embedded binary to disk (using validated absolute path)
	//nolint:gosec // We have validated the path, so this is not vulnerable to injection
	if err := os.WriteFile(absTargetPath, m.variant.Binary, 0600); err != nil {
		return fmt.Errorf("extraer binario: %w", err)
	}

	// 3. Register service with sc.exe using the validated absolute binary path
	output, err := secureScCreate(m.variant.RegistryName, absTargetPath, m.variant.DisplayName)
	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		//nolint:gosec // We have validated the path, so this is not vulnerable to injection
		_ = os.RemoveAll(absTargetDir) // clean up using validated absolute dir
		if strings.Contains(outputStr, "1073") {
			return fmt.Errorf("el servicio ya existe en el registro de Windows (use Desinstalar primero)")
		}
		// CLAVE: Mostramos el error de Go (%v) para no volar a ciegas
		return fmt.Errorf("fallo del sistema (%w) - Salida sc: '%s'", err, outputStr)
	}

	// 4. Configure failure recovery (restart on failure)
	_ = secureScFailure(m.variant.RegistryName)

	return nil
}

// Uninstall stops the service, removes it from the registry,
// and deletes the binary files from disk.
func (m *Manager) Uninstall() error {
	// Step 1: Attempt to stop the service
	_, stopErr := secureScRun("stop", m.variant.RegistryName)

	// Step 2: Wait for STOPPED state with proper timeout
	if stopErr == nil {
		stopped := m.WaitForStatus(StatusStopped, 15*time.Second)
		if !stopped {
			// Force-kill the service process as a last resort
			_ = secureTaskKillByService(m.variant.RegistryName)
			// Wait again briefly after force-kill
			m.WaitForStatus(StatusStopped, 5*time.Second)
		}
	}

	// Step 3: Delete service from registry
	output, err := secureScRun("delete", m.variant.RegistryName)
	if err != nil {
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

	// Step 4: Remove binary files from disk (validate and resolve ProgramFiles path)
	programFiles := os.Getenv("ProgramFiles")
	if programFiles == "" {
		return fmt.Errorf("environment variable `ProgramFiles` is empty")
	}
	targetDir := filepath.Join(programFiles, m.variant.RegistryName)
	_, absTargetDir, err := resolveAndEnsure(programFiles, targetDir)
	if err != nil {
		return fmt.Errorf("invalid target directory: %w", err)
	}
	//nolint:gosec // We have validated the path, so this is not vulnerable to injection
	if err := os.RemoveAll(absTargetDir); err != nil {
		return fmt.Errorf("no se pudieron eliminar los archivos: %w (puede que el proceso aún esté activo)", err)
	}

	return nil
}

// ══════════════════════════════════════════════════════════════
// Start / Stop / Restart
// ══════════════════════════════════════════════════════════════

// Start starts the Windows service
func (m *Manager) Start() error {
	output, err := secureScRun("start", m.variant.RegistryName)
	if err != nil {
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
	output, err := secureScRun("stop", m.variant.RegistryName)
	if err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "1062") {
			return fmt.Errorf("el servicio '%s' no está en ejecución", m.variant.DisplayName)
		}
		return fmt.Errorf("no se pudo detener '%s': %s", m.variant.DisplayName, strings.TrimSpace(outputStr))
	}
	return nil
}

// Restart restarts the Windows service. It first checks the current status and only attempts to stop if it's running or in a pending start state.
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

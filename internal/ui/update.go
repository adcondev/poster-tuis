package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/adcondev/poster-tuis/internal/service"
)

// ══════════════════════════════════════════════════════════════
// Main Update Function
// ══════════════════════════════════════════════════════════════

// Update handles all incoming messages and delegates to screen-specific handlers
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg), nil

	case statusUpdateMsg:
		return m.handleStatusUpdate(msg)

	case progressMsg:
		if m.processing {
			m.progressPercent += float64(msg)
			if m.progressPercent >= 1.0 {
				m.progressPercent = 1.0
			}
			return m, simulateProgress()
		}

	case operationDoneMsg:
		m.processing = false
		m.result = msg.message
		m.success = msg.success
		m.previousScreen = m.currentScreen
		m.currentScreen = screenResult
		m.progressPercent = 0
		return m, m.refreshStatusCmd()

	case tea.KeyMsg:
		switch m.currentScreen {
		case screenDashboard:
			return m.handleDashboardKey(msg)
		case screenFamily:
			return m.handleFamilyKey(msg)
		case screenLogs:
			return m.handleLogsKey(msg)
		case screenResult:
			return m.handleResultKey(msg)
		case screenConfirm:
			return m.handleConfirmKey(msg)
		case screenProcessing:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}

	case spinner.TickMsg:
		if m.processing {
			newSpinner, cmd := m.spinner.Update(msg)
			m.spinner = newSpinner
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// ══════════════════════════════════════════════════════════════
// Window Resize Handler
// ══════════════════════════════════════════════════════════════

func (m Model) handleWindowResize(msg tea.WindowSizeMsg) Model {
	m.width = msg.Width
	m.height = msg.Height
	m.ready = true

	headerHeight := 13
	footerHeight := 3

	listHeight := m.height - headerHeight - footerHeight
	if listHeight < 8 {
		listHeight = 8
	}

	m.list.SetSize(m.width-4, listHeight)
	m.progress.Width = m.width - 10
	if m.progress.Width < 20 {
		m.progress.Width = 20
	}
	m.help.Width = m.width

	return m
}

// ══════════════════════════════════════════════════════════════
// Status Update Handler
// ══════════════════════════════════════════════════════════════

func (m Model) handleStatusUpdate(msg statusUpdateMsg) (Model, tea.Cmd) {
	// If the message has actual data, use it; otherwise do a fresh check
	if msg.statuses != nil {
		m.familyStatuses = msg.statuses
	} else {
		// Periodic refresh: trigger actual check
		return m, m.refreshStatusCmd()
	}

	// Rebuild current menu to reflect updated statuses
	switch m.currentScreen {
	case screenDashboard:
		items := buildDashboardItems(m.familyStatuses)
		m.list.SetItems(items)
	case screenFamily:
		if m.selectedFamily != "" {
			fs := m.familyStatuses[m.selectedFamily]
			items := buildFamilyMenuItems(m.selectedFamily, fs)
			m.list.SetItems(items)
		}
	}

	// Schedule next periodic refresh
	return m, periodicRefreshCmd()
}

// ══════════════════════════════════════════════════════════════
// Key Handlers by Screen
// ══════════════════════════════════════════════════════════════

func (m Model) handleDashboardKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Quit) {
		return m, tea.Quit
	}

	if key.Matches(msg, m.keys.Enter) {
		item := m.list.SelectedItem()
		if item == nil {
			// No item selected; ignore Enter.
			return m, nil
		}

		selected, ok := item.(menuItem)
		if !ok {
			// Unexpected item type; ignore Enter.
			return m, nil
		}

		if selected.data == "quit" {
			return m, tea.Quit
		}

		// Navigate to family management screen
		return m.goToFamilyMenu(selected.data, selected.title)
	}

	// Delegate to list for navigation (up/down)
	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m Model) handleFamilyKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ESC/Q on family screen → back to dashboard
	if msg.String() == "esc" || msg.String() == "q" {
		return m.goToDashboard()
	}

	// "r" shortcut for restart (only when service is running)
	if key.Matches(msg, m.keys.Restart) {
		mgr := m.getActiveManager()
		if mgr != nil {
			fs := m.familyStatuses[m.selectedFamily]
			if fs.GetActiveStatus() == service.StatusRunning {
				return m.executeAction("Reiniciar Servicio", mgr.Restart)
			}
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Enter) {
		item := m.list.SelectedItem()
		if item == nil {
			// No item selected; ignore Enter.
			return m, nil
		}

		selected, ok := item.(menuItem)
		if !ok {
			// Unexpected item type; ignore Enter.
			return m, nil
		}

		switch selected.data {
		case "back":
			return m.goToDashboard()

		case "install-local":
			return m.confirmInstall("Local")

		case "install-remote":
			return m.confirmInstall("Remote")

		case "start":
			mgr := m.getActiveManager()
			if mgr == nil {
				m.statusMessage = "No hay servicio activo instalado."
				return m, nil
			}
			return m.executeAction("Iniciar Servicio", mgr.Start)

		case "stop":
			mgr := m.getActiveManager()
			if mgr == nil {
				m.statusMessage = "No hay servicio activo instalado."
				return m, nil
			}
			return m.executeAction("Detener Servicio", mgr.Stop)

		case "restart":
			mgr := m.getActiveManager()
			if mgr == nil {
				m.statusMessage = "No hay servicio activo instalado."
				return m, nil
			}
			return m.executeAction("Reiniciar Servicio", mgr.Restart)

		case "uninstall":
			return m.confirmUninstall()

		case "logs":
			return m.goToLogsMenu()
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m Model) handleLogsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ESC/Q on logs screen → back to family menu
	if msg.String() == "esc" || msg.String() == "q" {
		return m.returnToFamilyMenu()
	}

	if key.Matches(msg, m.keys.Enter) {
		item := m.list.SelectedItem()
		if item == nil {
			// No item selected; ignore Enter.
			return m, nil
		}

		selected, ok := item.(menuItem)
		if !ok {
			// Unexpected item type; ignore Enter.
			return m, nil
		}

		switch selected.data {
		case "open-file":
			mgr := m.getActiveManager()
			if mgr == nil {
				m.statusMessage = "No hay servicio activo instalado; no se pueden abrir los logs."
				return m, nil
			}
			if err := mgr.OpenLogFile(); err != nil {
				m.statusMessage = fmt.Sprintf("Error al abrir el archivo de logs: %v", err)
				return m, nil
			}
			m.statusMessage = "Abriendo archivo de logs..."
			return m, nil

		case "open-dir":
			mgr := m.getActiveManager()
			if mgr == nil {
				m.statusMessage = "No hay servicio activo instalado; no se puede abrir la carpeta de logs."
				return m, nil
			}
			if err := mgr.OpenLogDir(); err != nil {
				m.statusMessage = fmt.Sprintf("Error al abrir la carpeta de logs: %v", err)
				return m, nil
			}
			m.statusMessage = "Abriendo carpeta de logs..."
			return m, nil

		case "back":
			return m.returnToFamilyMenu()
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList
	return m, cmd
}

func (m Model) handleResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" || msg.String() == "esc" {
		m.result = ""
		m.statusMessage = ""

		// Navigate back to the appropriate screen
		switch m.previousScreen {
		case screenFamily, screenProcessing:
			updated, cmd := m.returnToFamilyMenu()
			return updated, tea.Batch(cmd, m.refreshStatusCmd())
		default:
			updated, cmd := m.goToDashboard()
			return updated, tea.Batch(cmd, m.refreshStatusCmd())
		}
	}
	return m, nil
}

func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s", "S":
		m.currentScreen = screenProcessing
		m.processing = true
		return m, tea.Batch(
			m.spinner.Tick,
			m.confirmCallback,
			simulateProgress(),
		)
	case "n", "N", "esc":
		m.confirmAction = ""
		m.confirmCallback = nil
		return m.returnToFamilyMenu()
	}
	return m, nil
}

// ══════════════════════════════════════════════════════════════
// Action Helpers
// ══════════════════════════════════════════════════════════════

// confirmInstall shows a confirmation dialog for installing a variant
func (m Model) confirmInstall(variant string) (Model, tea.Cmd) {
	m.confirmAction = fmt.Sprintf("¿Instalar versión %s de %s?",
		variant, capitalize(m.selectedFamily))

	m.confirmCallback = func() tea.Msg {
		variantID := fmt.Sprintf("%s-%s", m.selectedFamily, strings.ToLower(variant))
		mgr := m.managers[variantID]

		if err := mgr.Install(); err != nil {
			return operationDoneMsg{
				success: false,
				message: fmt.Sprintf("[X] Error al instalar: %v", err),
			}
		}

		// Auto-start after install
		if err := mgr.Start(); err != nil {
			return operationDoneMsg{
				success: true,
				message: fmt.Sprintf("[+] %s instalado (inicio manual requerido: %v)", variant, err),
			}
		}

		return operationDoneMsg{
			success: true,
			message: fmt.Sprintf("[+] %s instalado y iniciado correctamente", variant),
		}
	}

	m.previousScreen = screenFamily
	m.currentScreen = screenConfirm
	return m, nil
}

// confirmUninstall shows a confirmation dialog for uninstalling the active variant
func (m Model) confirmUninstall() (Model, tea.Cmd) {
	fs := m.familyStatuses[m.selectedFamily]
	installed := fs.GetInstalledVariant()

	m.confirmAction = fmt.Sprintf("¿Desinstalar %s de %s?",
		installed, capitalize(m.selectedFamily))

	m.confirmCallback = func() tea.Msg {
		mgr := m.getActiveManager()
		if mgr == nil {
			return operationDoneMsg{
				success: false,
				message: "[X] No se encontró el servicio instalado",
			}
		}

		if err := mgr.Uninstall(); err != nil {
			return operationDoneMsg{
				success: false,
				message: fmt.Sprintf("[X] Error al desinstalar: %v", err),
			}
		}

		return operationDoneMsg{
			success: true,
			message: fmt.Sprintf("[-] %s desinstalado correctamente", installed),
		}
	}

	m.previousScreen = screenFamily
	m.currentScreen = screenConfirm
	return m, nil
}

// executeAction wraps a service operation with a loading/processing screen
func (m Model) executeAction(actionName string, fn func() error) (Model, tea.Cmd) {
	m.processing = true
	m.currentScreen = screenProcessing

	cmd := func() tea.Msg {
		if err := fn(); err != nil {
			return operationDoneMsg{
				success: false,
				message: fmt.Sprintf("[X] %s falló: %v", actionName, err),
			}
		}
		return operationDoneMsg{
			success: true,
			message: fmt.Sprintf("[OK] %s completado", actionName),
		}
	}

	return m, tea.Batch(m.spinner.Tick, cmd, simulateProgress())
}

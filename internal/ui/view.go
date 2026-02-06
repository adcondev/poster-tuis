package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/adcondev/poster-tuis/internal/config"
	"github.com/adcondev/poster-tuis/internal/service"
)

// ══════════════════════════════════════════════════════════════
// Main View Function
// ══════════════════════════════════════════════════════════════

// View renders the current screen (required by bubbletea.Model interface)
func (m Model) View() string {
	if !m.ready {
		return "Inicializando..."
	}

	switch m.currentScreen {
	case screenDashboard:
		return m.viewDashboard()
	case screenFamily:
		return m.viewFamily()
	case screenLogs:
		return m.viewLogs()
	case screenProcessing:
		return m.viewProcessing()
	case screenResult:
		return m.viewResult()
	case screenConfirm:
		return m.viewConfirm()
	default:
		return "Estado desconocido"
	}
}

// ══════════════════════════════════════════════════════════════
// Screen Renderers
// ══════════════════════════════════════════════════════════════

func (m Model) viewDashboard() string {
	var b strings.Builder

	b.WriteString(bannerStyle.Render(config.GetBanner()) + "\n\n")
	b.WriteString(titleStyle.Render("SELECCIONE UNA FAMILIA DE SERVICIOS") + "\n\n")

	b.WriteString(m.list.View())

	// Health summary bar
	healthSummary := m.renderHealthSummary()
	b.WriteString("\n" + statusBarStyle.Render(healthSummary))

	b.WriteString("\n" + m.help.View(m.keys))

	if m.statusMessage != "" {
		b.WriteString("\n" + infoStyle.Render(m.statusMessage))
	}

	return b.String()
}

func (m Model) viewFamily() string {
	var b strings.Builder

	fs := m.familyStatuses[m.selectedFamily]
	installed := fs.GetInstalledVariant()

	b.WriteString(bannerStyle.Render(config.GetBanner()) + "\n\n")

	if installed == "" {
		b.WriteString(titleStyle.Render(
			fmt.Sprintf("%s - SIN INSTALAR", strings.ToUpper(m.selectedFamily))) + "\n")
		b.WriteString(infoStyle.Render(
			"Seleccione una versión para instalar (Local o Remote)") + "\n\n")
	} else {
		status := fs.GetActiveStatus()
		b.WriteString(titleStyle.Render(
			fmt.Sprintf("%s - %s", strings.ToUpper(m.selectedFamily), installed)) + "\n")
		b.WriteString(statusBarStyle.Render(
			fmt.Sprintf("Estado: %s", status.String())) + "\n\n")
	}

	b.WriteString(m.list.View())
	b.WriteString("\n" + m.help.View(m.keys))

	if m.statusMessage != "" {
		b.WriteString("\n" + successStyle.Render(m.statusMessage))
	}

	return b.String()
}

func (m Model) viewLogs() string {
	var b strings.Builder

	b.WriteString(bannerStyle.Render(config.GetBanner()) + "\n\n")
	b.WriteString(statusBarStyle.Render(
		fmt.Sprintf("[#] GESTIÓN DE LOGS - %s", strings.ToUpper(m.selectedFamily))) + "\n\n")

	b.WriteString(m.list.View())

	b.WriteString("\n" + infoStyle.Render("[ESC] Volver al menú de servicio"))

	if m.statusMessage != "" {
		b.WriteString("\n" + successStyle.Render(m.statusMessage))
	}

	return b.String()
}

func (m Model) viewProcessing() string {
	var b strings.Builder

	b.WriteString(bannerStyle.Render(config.GetBanner()) + "\n\n")

	spinnerView := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.spinner.View(),
		" Procesando operación...",
	)
	b.WriteString(spinnerView + "\n\n")

	if m.progressPercent > 0 {
		b.WriteString(m.progress.ViewAs(m.progressPercent))
		b.WriteString(fmt.Sprintf("\n%.0f%% completado", m.progressPercent*100))
	}

	pulseStyle := lipgloss.NewStyle().Foreground(secondaryColor)
	b.WriteString("\n\n" + pulseStyle.Render("[~] Por favor espere..."))

	return b.String()
}

func (m Model) viewResult() string {
	var b strings.Builder

	b.WriteString(bannerStyle.Render(config.GetBanner()) + "\n\n")

	var boxStyle lipgloss.Style
	if m.success {
		boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2).
			Width(m.width - 10)
	} else {
		boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			Width(m.width - 10)
	}

	b.WriteString(boxStyle.Render(m.result))
	b.WriteString("\n\n" + infoStyle.Render("Presione Enter para continuar..."))

	return b.String()
}

func (m Model) viewConfirm() string {
	var b strings.Builder

	b.WriteString(bannerStyle.Render(config.GetBanner()) + "\n\n")

	confirmBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(warningColor).
		Padding(1, 2).
		Width(m.width - 10).
		Align(lipgloss.Center)

	content := fmt.Sprintf("[!] CONFIRMACIÓN\n\n%s\n\n", m.confirmAction)
	content += successStyle.Render("[S]í") + "    " + warningStyle.Render("[N]o")

	b.WriteString(confirmBox.Render(content))

	return b.String()
}

// ══════════════════════════════════════════════════════════════
// Helper Renderers
// ══════════════════════════════════════════════════════════════

// renderHealthSummary generates the status bar showing all families' health
func (m Model) renderHealthSummary() string {
	var parts []string

	for _, family := range service.GetFamilyNames() {
		fs := m.familyStatuses[family]
		installed := fs.GetInstalledVariant()

		if installed == "" {
			parts = append(parts, fmt.Sprintf("%s: [-]", family))
		} else {
			status := fs.GetActiveStatus()

			icon := "?"
			switch status {
			case service.StatusRunning:
				icon = "+"
			case service.StatusStopped:
				icon = "."
				// default: icon stays "?" — safe for StatusUnknown
			}

			parts = append(parts, fmt.Sprintf("%s: [%s] %s", family, icon, installed))
		}
	}

	return strings.Join(parts, " | ")
}

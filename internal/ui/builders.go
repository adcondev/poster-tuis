package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"

	"github.com/adcondev/poster-tuis/internal/service"
)

// ══════════════════════════════════════════════════════════════
// Menu Item Type
// ══════════════════════════════════════════════════════════════

// menuItem implements list.Item interface for bubbletea lists
type menuItem struct {
	title       string
	description string
	icon        string
	data        string // Generic data field for routing/action identification
}

func (i menuItem) Title() string       { return i.icon + " " + i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

// ══════════════════════════════════════════════════════════════
// Dashboard Menu Builder
// ══════════════════════════════════════════════════════════════

// buildDashboardItems creates the main menu showing both service families
func buildDashboardItems(statuses map[string]service.FamilyStatus) []list.Item {
	return []list.Item{
		menuItem{
			title:       "Scale Service",
			description: formatFamilyStatus(statuses["scale"]),
			icon:        "[1]",
			data:        "scale",
		},
		menuItem{
			title:       "Ticket Service",
			description: formatFamilyStatus(statuses["ticket"]),
			icon:        "[2]",
			data:        "ticket",
		},
		menuItem{
			title:       "Salir",
			description: "Cerrar el instalador",
			icon:        "[Q]",
			data:        "quit",
		},
	}
}

// formatFamilyStatus generates a human-readable status summary for the dashboard
func formatFamilyStatus(fs service.FamilyStatus) string {
	installed := fs.GetInstalledVariant()
	if installed == "" {
		return "No instalado"
	}

	status := fs.GetActiveStatus()
	return fmt.Sprintf("%s - %s", installed, status.String())
}

// ══════════════════════════════════════════════════════════════
// Family Menu Builder (CRITICAL — Enforces Mutual Exclusivity)
// ══════════════════════════════════════════════════════════════

// buildFamilyMenuItems creates the menu for managing a specific service family.
// This function ENFORCES mutual exclusivity by controlling which options appear:
// - If neither variant is installed → show install options for Local AND Remote
// - If one variant is installed → show lifecycle actions for that variant ONLY
// - The install option for the OTHER variant NEVER appears when one is installed
func buildFamilyMenuItems(family string, fs service.FamilyStatus) []list.Item {
	installed := fs.GetInstalledVariant()

	// ── SCENARIO A: Clean slate (neither variant installed) ──
	if installed == "" {
		return []list.Item{
			menuItem{
				title:       "Instalar Versión LOCAL",
				description: "Instala el servicio para uso en este equipo",
				icon:        "[1]",
				data:        "install-local",
			},
			menuItem{
				title:       "Instalar Versión REMOTA",
				description: "Instala el servicio para acceso desde red (LAN)",
				icon:        "[2]",
				data:        "install-remote",
			},
			menuItem{
				title:       "Volver",
				description: "Regresar al menú principal",
				icon:        "[<]",
				data:        "back",
			},
		}
	}

	// ── SCENARIO B: One variant is installed ──
	status := fs.GetActiveStatus()
	items := []list.Item{}

	// Conditional actions based on service state
	if status == service.StatusStopped {
		items = append(items, menuItem{
			title:       "Iniciar Servicio",
			description: "Inicia el servicio detenido",
			icon:        "[>]",
			data:        "start",
		})
	}

	if status == service.StatusRunning {
		items = append(items, menuItem{
			title:       "Detener Servicio",
			description: "Detiene el servicio en ejecución",
			icon:        "[.]",
			data:        "stop",
		})
		items = append(items, menuItem{
			title:       "Reiniciar Servicio",
			description: "Reinicia el servicio (detiene e inicia)",
			icon:        "[*]",
			data:        "restart",
		})
	}

	// Always available when a variant is installed
	items = append(items,
		menuItem{
			title:       "Ver Logs",
			description: "Abrir archivo o carpeta de logs",
			icon:        "[#]",
			data:        "logs",
		},
		menuItem{
			title:       fmt.Sprintf("Desinstalar %s", installed),
			description: "Elimina completamente el servicio del sistema",
			icon:        "[-]",
			data:        "uninstall",
		},
		menuItem{
			title:       "Volver",
			description: "Regresar al menú principal",
			icon:        "[<]",
			data:        "back",
		},
	)

	return items
}

// ══════════════════════════════════════════════════════════════
// Logs Menu Builder
// ══════════════════════════════════════════════════════════════

// buildLogsMenuItems creates the log management submenu
func buildLogsMenuItems() []list.Item {
	return []list.Item{
		menuItem{
			title:       "Abrir Archivo de Logs",
			description: "Abre el archivo .log en Notepad",
			icon:        "[#]",
			data:        "open-file",
		},
		menuItem{
			title:       "Abrir Carpeta de Logs",
			description: "Abre la ubicación de los logs en Explorer",
			icon:        "[D]",
			data:        "open-dir",
		},
		menuItem{
			title:       "Volver",
			description: "Regresar al menú de servicio",
			icon:        "[<]",
			data:        "back",
		},
	}
}

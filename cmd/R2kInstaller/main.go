package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/adcondev/poster-tuis/internal/ui"
)

// ══════════════════════════════════════════════════════════════
// Admin Check
// ══════════════════════════════════════════════════════════════

// isAdmin checks if the program is running with administrator privileges
// by attempting to open PHYSICALDRIVE0, which requires elevated access.
func isAdmin() bool {
	f, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
			os.Exit(1)
		}
	}(f)
	return true
}

// ══════════════════════════════════════════════════════════════
// Main Entry Point
// ══════════════════════════════════════════════════════════════

func main() {
	// Enforce admin privileges — required for sc.exe operations
	if !isAdmin() {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")).
			Bold(true)
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa2f7"))

		fmt.Println(errorStyle.Render("[!] Permisos de Administrador Requeridos"))
		fmt.Println(infoStyle.Render("\n[i] Instrucciones:"))
		fmt.Println("1. Cierre esta ventana")
		fmt.Println("2. Clic derecho en el instalador")
		fmt.Println("3. Seleccione 'Ejecutar como administrador'")
		fmt.Println("\nPresione Enter para salir...")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}

	// Start TUI application
	p := tea.NewProgram(
		ui.InitialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/adcondev/poster-tuis/internal/service"
)

// ══════════════════════════════════════════════════════════════
// Model Definition
// ══════════════════════════════════════════════════════════════

// Model is the top-level bubbletea model for the TUI application
type Model struct {
	// Navigation state
	currentScreen  screen
	previousScreen screen
	selectedFamily string // "scale" or "ticket"

	// Service data
	registry       map[string][]service.ServiceVariant
	managers       map[string]*service.Manager // Indexed by variant.ID
	familyStatuses map[string]service.FamilyStatus

	// UI components
	list     list.Model
	spinner  spinner.Model
	progress progress.Model
	help     help.Model
	keys     keyMap

	// Operation state
	processing      bool
	result          string
	success         bool
	confirmAction   string
	confirmCallback tea.Cmd
	progressPercent float64
	statusMessage   string

	// Dimensions
	width  int
	height int
	ready  bool
}

// ══════════════════════════════════════════════════════════════
// Key Map
// ══════════════════════════════════════════════════════════════

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Help    key.Binding
	Quit    key.Binding
	Restart key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Enter, k.Restart},
		{k.Help, k.Quit},
	}
}

var defaultKeys = keyMap{
	Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "arriba")),
	Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "abajo")),
	Enter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "seleccionar")),
	Restart: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reiniciar servicio")),
	Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "ayuda")),
	Quit:    key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "salir")),
}

// ══════════════════════════════════════════════════════════════
// Message Types
// ══════════════════════════════════════════════════════════════

type statusUpdateMsg struct {
	statuses map[string]service.FamilyStatus
}

type operationDoneMsg struct {
	success bool
	message string
}

type progressMsg float64

// ══════════════════════════════════════════════════════════════
// Initialization
// ══════════════════════════════════════════════════════════════

// InitialModel creates the initial model with all services registered
func InitialModel() Model {
	// Get service registry
	registry := service.GetServiceRegistry()

	// Create managers for all variants
	managers := make(map[string]*service.Manager)
	for _, variants := range registry {
		for _, variant := range variants {
			managers[variant.ID] = service.NewManager(variant)
		}
	}

	// Initialize family statuses (initial check)
	familyStatuses := make(map[string]service.FamilyStatus)
	for _, family := range service.GetFamilyNames() {
		familyStatuses[family] = service.CheckFamilyStatus(registry[family])
	}

	// ── Initialize UI components ──

	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = spinnerStyle

	p := progress.New(
		progress.WithScaledGradient("#CC0033", "#33A0FF"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	h := help.New()
	h.Styles.ShortKey = helpKeyStyle
	h.Styles.ShortDesc = helpDescStyle
	h.Styles.FullKey = helpKeyStyle
	h.Styles.FullDesc = helpDescStyle

	// Build dashboard menu
	dashboardItems := buildDashboardItems(familyStatuses)

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedStyle
	delegate.Styles.SelectedDesc = selectedStyle.Faint(true)
	delegate.Styles.NormalTitle = normalStyle
	delegate.Styles.NormalDesc = normalStyle.Faint(true)
	delegate.Styles.DimmedTitle = disabledStyle
	delegate.Styles.DimmedDesc = disabledStyle.Faint(true)

	l := list.New(dashboardItems, delegate, 80, 20)
	l.Title = "Gestor de Servicios"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return Model{
		currentScreen:  screenDashboard,
		registry:       registry,
		managers:       managers,
		familyStatuses: familyStatuses,
		list:           l,
		spinner:        s,
		progress:       p,
		help:           h,
		keys:           defaultKeys,
		ready:          false,
	}
}

// Init starts background tasks (required by bubbletea.Model interface)
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.refreshStatusCmd(),
	)
}

// ══════════════════════════════════════════════════════════════
// Helper Functions
// ══════════════════════════════════════════════════════════════

// getActiveManager returns the manager for the currently installed variant
// of the selected family. Returns nil if no variant is installed.
func (m Model) getActiveManager() *service.Manager {
	fs := m.familyStatuses[m.selectedFamily]
	installed := fs.GetInstalledVariant()

	if installed == "" {
		return nil
	}

	variantID := fmt.Sprintf("%s-%s", m.selectedFamily, strings.ToLower(installed))
	return m.managers[variantID]
}

// refreshStatusCmd checks the status of all service families in the background
func (m Model) refreshStatusCmd() tea.Cmd {
	return func() tea.Msg {
		statuses := make(map[string]service.FamilyStatus)

		for _, family := range service.GetFamilyNames() {
			statuses[family] = service.CheckFamilyStatus(m.registry[family])
		}

		return statusUpdateMsg{statuses: statuses}
	}
}

// periodicRefreshCmd sets up periodic status polling
func periodicRefreshCmd() tea.Cmd {
	return tea.Every(5*time.Second, func(t time.Time) tea.Msg {
		// Return an empty statusUpdateMsg to trigger a refresh
		// The actual check happens in the Update handler
		return statusUpdateMsg{}
	})
}

// ── Navigation helpers ──

// goToDashboard rebuilds the dashboard menu and navigates to it
func (m Model) goToDashboard() (Model, tea.Cmd) {
	items := buildDashboardItems(m.familyStatuses)
	m.list.SetItems(items)
	m.list.Title = "Gestor de Servicios"
	m.currentScreen = screenDashboard
	m.selectedFamily = ""
	m.statusMessage = ""
	return m, nil
}

// goToFamilyMenu navigates to the family management screen
func (m Model) goToFamilyMenu(family, familyTitle string) (Model, tea.Cmd) {
	m.selectedFamily = family
	m.previousScreen = screenDashboard
	m.currentScreen = screenFamily

	fs := m.familyStatuses[family]
	items := buildFamilyMenuItems(family, fs)
	m.list.SetItems(items)
	m.list.Title = fmt.Sprintf("Gestión - %s", familyTitle)

	return m, nil
}

// goToLogsMenu navigates to the logs management submenu
func (m Model) goToLogsMenu() (Model, tea.Cmd) {
	items := buildLogsMenuItems()
	m.list.SetItems(items)
	m.list.Title = fmt.Sprintf("Logs - %s", capitalize(m.selectedFamily))
	m.previousScreen = screenFamily
	m.currentScreen = screenLogs
	return m, nil
}

// returnToFamilyMenu rebuilds the family menu and navigates back to it
func (m Model) returnToFamilyMenu() (Model, tea.Cmd) {
	fs := m.familyStatuses[m.selectedFamily]
	items := buildFamilyMenuItems(m.selectedFamily, fs)
	m.list.SetItems(items)
	m.list.Title = fmt.Sprintf("Gestión - %s", capitalize(m.selectedFamily))
	m.currentScreen = screenFamily
	m.statusMessage = ""
	return m, nil
}

// simulateProgress creates a progress animation command
func simulateProgress() tea.Cmd {
	return tea.Every(100*time.Millisecond, func(t time.Time) tea.Msg {
		return progressMsg(0.1)
	})
}

// capitalize returns a string with the first letter uppercased
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

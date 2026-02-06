package ui

// ══════════════════════════════════════════════════════════════
// Screen State Constants
// ══════════════════════════════════════════════════════════════
// Defines the strict state machine for UI navigation.
//
// Navigation flow:
//   screenDashboard → screenFamily → screenLogs
//                                  → screenProcessing → screenResult
//                                  → screenConfirm → screenProcessing → screenResult

type screen int

const (
	screenDashboard  screen = iota // Family selector (main menu)
	screenFamily                   // Service operations for selected family
	screenLogs                     // Log management submenu
	screenProcessing               // Blocking operation indicator (with spinner/progress)
	screenResult                   // Operation result display (success/error)
	screenConfirm                  // Yes/No confirmation dialog
)

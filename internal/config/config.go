package config

import "fmt"

// Variables injected by Taskfile (ldflags)
var (
	BuildDate           string
	BuildTime           string
	ScaleIDLocal        string
	ScaleIDRemote       string
	ScaleDisplayLocal   string
	ScaleDisplayRemote  string
	TicketIDLocal       string
	TicketIDRemote      string
	TicketDisplayLocal  string
	TicketDisplayRemote string
)

func GetBanner() string {
	// Build the "Build:" line dynamically and then pad/truncate it to fit the banner width.
	buildInfo := fmt.Sprintf("Build: %s %s", BuildDate, BuildTime)

	return fmt.Sprintf(`
╔═════════════════════════════════════════════╗
║        SERVICE FAMILY MANAGER v2.0          ║
║        %s        ║
║                                             ║
║     Gestión de Servicios Red2000            ║
║     - Scale Service (Local/Remoto)          ║
║     - Ticket Service (Local/Remoto)         ║
║                                             ║
╚═════════════════════════════════════════════╝`,
		buildInfo,
	)
}

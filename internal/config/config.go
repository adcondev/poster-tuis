package config

import "fmt"

// Variables injected by Taskfile (ldflags)
var (
	BuildDate      string
	BuildTime      string
	ScaleBaseName  string
	TicketBaseName string
)

const (
	ServiceSuffix = "Servicio"
)

func GenerateServiceNames(baseName, variant string) (registryName, displayName, exeName string) {
	registryName = fmt.Sprintf("%s%s_%s", baseName, ServiceSuffix, variant)
	displayName = fmt.Sprintf("%s %s (%s)", baseName, ServiceSuffix, variant)
	exeName = fmt.Sprintf("%s.exe", registryName)
	return
}

func GetBanner() string {
	return fmt.Sprintf(`
╔═════════════════════════════════════════════╗
║        SERVICE FAMILY MANAGER v2.0          ║
║       Build: %s %s         ║
║                                             ║
║     Gestión de Servicios Red2000            ║
║     - Scale Service (Local/Remote)          ║
║     - Ticket Service (Local/Remote)         ║
║                                             ║
╚═════════════════════════════════════════════╝`,
		BuildDate, BuildTime,
	)
}

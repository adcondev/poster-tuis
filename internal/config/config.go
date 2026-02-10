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
	// Build the "Build:" line dynamically and then pad/truncate it to fit the banner width.
	buildInfo := fmt.Sprintf("Build: %s %s", BuildDate, BuildTime)

	return fmt.Sprintf(`
╔═════════════════════════════════════════════╗
║        SERVICE FAMILY MANAGER v2.0          ║
║         %s         ║
║                                             ║
║     Gestión de Servicios Red2000            ║
║     - Scale Service (Local/Remoto)          ║
║     - Ticket Service (Local/Remoto)         ║
║                                             ║
╚═════════════════════════════════════════════╝`,
		buildInfo,
	)
}

// Package config provides configuration values and utilities for the POS installer application.
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

// Define ANSI color codes for the "dope" look
var colors = map[string]string{
	"reset":  "\033[0m",
	"cyan":   "\033[36m",
	"blue":   "\033[34m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"white":  "\033[97m",
	"bold":   "\033[1m",
}

// GetBanner returns a stylized banner string with build info and ASCII art logo
func GetBanner() string {
	// Build Info (ejemplo)
	buildStr := fmt.Sprintf("Build: %s", BuildDate)

	return fmt.Sprintf(`%s╔══════════════════════════════════════════════════════╗%s
%s║    ██████╗ ██████╗ ██╗  ██╗                          ║%s
%s║    ██╔══██╗╚════██╗██║ ██╔╝   %sInstalador v3.0.0%s      %s║%s
%s║    ██████╔╝ █████╔╝█████╔╝    %s%s%s      %s║%s
%s║    ██╔══██╗██╔═══╝ ██╔═██╗                           ║%s
%s║    ██║  ██║███████╗██║  ██╗   -> Ticket v2.0.0       %s║%s
%s║    ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   -> Scale  v2.0.0       ║%s
%s║                                                      ║%s
%s║            Instalador de Servicios POS               ║%s
%s║                 (C) 2025 Red2000                     ║%s
%s╚══════════════════════════════════════════════════════╝%s`,
		// Border Colors
		colors["blue"], colors["reset"],

		// R2K Logo (Cyan)
		colors["cyan"], colors["reset"],
		colors["cyan"], colors["white"], colors["cyan"], colors["blue"], colors["reset"],
		colors["cyan"], colors["reset"], colors["yellow"], buildStr, colors["blue"], colors["reset"],
		colors["cyan"], colors["reset"],
		colors["cyan"], colors["blue"], colors["reset"],
		colors["cyan"], colors["reset"],

		// Footer
		colors["blue"], colors["reset"],
		colors["blue"], colors["reset"],
		colors["blue"], colors["reset"],
		colors["blue"], colors["reset"],
	)
}

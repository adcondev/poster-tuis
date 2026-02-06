package service

import (
	"github.com/adcondev/poster-tuis/internal/assets"
	"github.com/adcondev/poster-tuis/internal/config"
)

// ══════════════════════════════════════════════════════════════
// Service Variant Definition
// ══════════════════════════════════════════════════════════════

// Variant represents a specific variant (Local/Remote) of a service family
type Variant struct {
	ID           string // Unique identifier: "scale-local", "ticket-remote"
	Family       string // Family name: "scale", "ticket"
	Variant      string // Variant type: "Local", "Remote"
	RegistryName string // Windows service registry name
	DisplayName  string // Human-readable display name
	ExeName      string // Binary filename on disk
	Binary       []byte // Embedded binary data
}

// ══════════════════════════════════════════════════════════════
// Registry Functions
// ══════════════════════════════════════════════════════════════

// GetServiceRegistry returns all service families with their variants.
// Returns map with keys: "scale", "ticket"
// Each key maps to a slice of 2 variants (Local, Remote).
func GetServiceRegistry() map[string][]Variant {
	// Generate naming for Scale family
	scaleLocalReg, scaleLocalDisp, scaleLocalExe := config.GenerateServiceNames(config.ScaleBaseName, "Local")
	scaleRemoteReg, scaleRemoteDisp, scaleRemoteExe := config.GenerateServiceNames(config.ScaleBaseName, "Remote")

	// Generate naming for Ticket family
	ticketLocalReg, ticketLocalDisp, ticketLocalExe := config.GenerateServiceNames(config.TicketBaseName, "Local")
	ticketRemoteReg, ticketRemoteDisp, ticketRemoteExe := config.GenerateServiceNames(config.TicketBaseName, "Remote")

	return map[string][]Variant{
		"scale": {
			{
				ID:           "scale-local",
				Family:       "scale",
				Variant:      "Local",
				RegistryName: scaleLocalReg,
				DisplayName:  scaleLocalDisp,
				ExeName:      scaleLocalExe,
				Binary:       assets.BasculaLocalBinary,
			},
			{
				ID:           "scale-remote",
				Family:       "scale",
				Variant:      "Remote",
				RegistryName: scaleRemoteReg,
				DisplayName:  scaleRemoteDisp,
				ExeName:      scaleRemoteExe,
				Binary:       assets.BasculaRemoteBinary,
			},
		},
		"ticket": {
			{
				ID:           "ticket-local",
				Family:       "ticket",
				Variant:      "Local",
				RegistryName: ticketLocalReg,
				DisplayName:  ticketLocalDisp,
				ExeName:      ticketLocalExe,
				Binary:       assets.TicketLocalBinary,
			},
			{
				ID:           "ticket-remote",
				Family:       "ticket",
				Variant:      "Remote",
				RegistryName: ticketRemoteReg,
				DisplayName:  ticketRemoteDisp,
				ExeName:      ticketRemoteExe,
				Binary:       assets.TicketRemoteBinary,
			},
		},
	}
}

// GetFamilyNames returns the ordered list of service family identifiers
func GetFamilyNames() []string {
	return []string{"scale", "ticket"}
}

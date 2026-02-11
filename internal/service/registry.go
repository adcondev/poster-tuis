package service

import (
	"github.com/adcondev/poster-tuis/internal/assets"
	"github.com/adcondev/poster-tuis/internal/config"
)

// ══════════════════════════════════════════════════════════════
// Service Variant Definition
// ══════════════════════════════════════════════════════════════

// Variant represents a specific variant (Local/Remoto) of a service family
type Variant struct {
	ID           string // Unique identifier: "scale-local", "ticket-remote"
	Family       string // Family name: "scale", "ticket"
	Variant      string // Variant type: "Local", "Remoto"
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
// Each key maps to a slice of 2 variants (Local, Remoto).
func GetServiceRegistry() map[string][]Variant {
	// Helper to generate Display Name and Exe Name
	// We use the ID from config as the filename base to ensure consistency
	makeVariant := func(id, family, variantStr, registryID, displayName string, binary []byte) Variant {
		exeName := registryID + ".exe"

		return Variant{
			ID:           id,
			Family:       family,
			Variant:      variantStr,
			RegistryName: registryID, // <--- DIRECTLY FROM TASKFILE
			DisplayName:  displayName,
			ExeName:      exeName,
			Binary:       binary,
		}
	}

	return map[string][]Variant{
		"scale": {
			makeVariant("scale-local", "scale", "Local", config.ScaleIDLocal, config.ScaleDisplayLocal, assets.BasculaLocalBinary),
			// ID is "scale-remoto" so UI lookup matches "Remoto"
			makeVariant("scale-remoto", "scale", "Remoto", config.ScaleIDRemote, config.ScaleDisplayRemote, assets.BasculaRemoteBinary),
		},
		"ticket": {
			makeVariant("ticket-local", "ticket", "Local", config.TicketIDLocal, config.TicketDisplayLocal, assets.TicketLocalBinary),
			makeVariant("ticket-remoto", "ticket", "Remoto", config.TicketIDRemote, config.TicketDisplayRemote, assets.TicketRemoteBinary),
		},
	}
}

// GetFamilyNames returns the ordered list of service family identifiers
func GetFamilyNames() []string {
	return []string{"scale", "ticket"}
}

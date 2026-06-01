package service

import (
	"testing"
	"github.com/adcondev/poster-tuis/internal/config"
	"github.com/adcondev/poster-tuis/internal/assets"
)

func TestGetServiceRegistry(t *testing.T) {
	// Set up config vars for testing
	config.ScaleIDLocal = "ScaleLocal"
	config.ScaleIDRemote = "ScaleRemote"
	config.ScaleDisplayLocal = "Scale Local Display"
	config.ScaleDisplayRemote = "Scale Remote Display"
	config.TicketIDLocal = "TicketLocal"
	config.TicketIDRemote = "TicketRemote"
	config.TicketDisplayLocal = "Ticket Local Display"
	config.TicketDisplayRemote = "Ticket Remote Display"

	// Mock assets so Binary isn't empty
	assets.BasculaLocalBinary = []byte("mock")
	assets.BasculaRemoteBinary = []byte("mock")
	assets.TicketLocalBinary = []byte("mock")
	assets.TicketRemoteBinary = []byte("mock")

	registry := GetServiceRegistry()

	families := GetFamilyNames()
	if len(registry) != len(families) {
		t.Errorf("expected %d families in registry, got %d", len(families), len(registry))
	}

	for _, family := range families {
		variants, ok := registry[family]
		if !ok {
			t.Errorf("expected family %s in registry", family)
		}
		if len(variants) != 2 {
			t.Errorf("expected 2 variants for family %s, got %d", family, len(variants))
		}

		for _, variant := range variants {
			if variant.Family != family {
				t.Errorf("expected family %s, got %s", family, variant.Family)
			}
			if variant.Variant != Local && variant.Variant != Remoto {
				t.Errorf("expected variant Local or Remoto, got %s", variant.Variant)
			}
			if variant.RegistryName == "" {
				t.Errorf("expected non-empty RegistryName")
			}
			if variant.DisplayName == "" {
				t.Errorf("expected non-empty DisplayName")
			}
			if variant.ExeName == "" {
				t.Errorf("expected non-empty ExeName")
			}
			if len(variant.Binary) == 0 {
				t.Errorf("expected non-empty Binary")
			}

			// Use the existing validation function
			err := validateServiceVariantFields(variant)
			if err != nil {
				t.Errorf("failed validation for variant %s: %v", variant.ID, err)
			}
		}
	}
}

func TestGetFamilyNames(t *testing.T) {
	families := GetFamilyNames()
	expected := []string{"scale", "ticket"}
	if len(families) != len(expected) {
		t.Fatalf("expected %d families, got %d", len(expected), len(families))
	}
	for i, f := range families {
		if f != expected[i] {
			t.Errorf("expected family at index %d to be %s, got %s", i, expected[i], f)
		}
	}
}

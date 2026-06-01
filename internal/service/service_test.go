package service

import (
	"strings"
	"testing"
)

func TestValidateServiceVariantFields(t *testing.T) {
	tests := []struct {
		name        string
		variant     Variant
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid Variant",
			variant: Variant{
				RegistryName: "Service_Name-1",
				ExeName:      "service.exe",
				DisplayName:  "My Service Display Name",
			},
			wantErr: false,
		},
		{
			name: "Invalid RegistryName",
			variant: Variant{
				RegistryName: "Service!Name",
				ExeName:      "service.exe",
				DisplayName:  "My Service",
			},
			wantErr:     true,
			errContains: "invalid RegistryName",
		},
		{
			name: "Empty RegistryName",
			variant: Variant{
				RegistryName: "",
				ExeName:      "service.exe",
				DisplayName:  "My Service",
			},
			wantErr:     true,
			errContains: "invalid RegistryName",
		},
		{
			name: "Invalid ExeName Path Traversal",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "../service.exe",
				DisplayName:  "My Service",
			},
			wantErr:     true,
			errContains: "invalid ExeName",
		},
		{
			name: "Invalid ExeName Forward Slash",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "dir/service.exe",
				DisplayName:  "My Service",
			},
			wantErr:     true,
			errContains: "invalid ExeName",
		},
		{
			name: "Invalid ExeName Backslash",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "dir\\service.exe",
				DisplayName:  "My Service",
			},
			wantErr:     true,
			errContains: "invalid ExeName",
		},
		{
			name: "Invalid ExeName Shell Metachar",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "service$.exe",
				DisplayName:  "My Service",
			},
			wantErr:     true,
			errContains: "invalid ExeName",
		},
		{
			name: "Empty ExeName",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "",
				DisplayName:  "My Service",
			},
			wantErr:     true,
			errContains: "invalid ExeName",
		},
		{
			name: "Invalid DisplayName Shell Metachar",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "service.exe",
				DisplayName:  "My Service $",
			},
			wantErr:     true,
			errContains: "invalid DisplayName",
		},
		{
			name: "Invalid DisplayName Newline",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "service.exe",
				DisplayName:  "My Service \n",
			},
			wantErr:     true,
			errContains: "invalid DisplayName",
		},
		{
			name: "Empty DisplayName",
			variant: Variant{
				RegistryName: "ServiceName",
				ExeName:      "service.exe",
				DisplayName:  "",
			},
			wantErr:     true,
			errContains: "invalid DisplayName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServiceVariantFields(tt.variant)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateServiceVariantFields() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validateServiceVariantFields() error = %v, expected it to contain %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("validateServiceVariantFields() unexpected error = %v", err)
				}
			}
		})
	}
}

package service

import (
	"strings"
	"testing"
)

const (
	errInvalidDisplayName = "invalid DisplayName"
	errInvalidExeName     = "invalid ExeName"
	errInvalidRegistryName = "invalid RegistryName"
	testServiceName       = "ServiceName"
	testMyService         = "My Service"
	testServiceExe        = "service.exe"
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
				ExeName:      testServiceExe,
				DisplayName:  "My Service Display Name",
			},
			wantErr: false,
		},
		{
			name: "Invalid RegistryName",
			variant: Variant{
				RegistryName: "Service!Name",
				ExeName:      testServiceExe,
				DisplayName:  testMyService,
			},
			wantErr:     true,
			errContains: errInvalidRegistryName,
		},
		{
			name: "Empty RegistryName",
			variant: Variant{
				RegistryName: "",
				ExeName:      testServiceExe,
				DisplayName:  testMyService,
			},
			wantErr:     true,
			errContains: errInvalidRegistryName,
		},
		{
			name: "Invalid ExeName Path Traversal",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      "../service.exe",
				DisplayName:  testMyService,
			},
			wantErr:     true,
			errContains: errInvalidExeName,
		},
		{
			name: "Invalid ExeName Forward Slash",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      "dir/service.exe",
				DisplayName:  testMyService,
			},
			wantErr:     true,
			errContains: errInvalidExeName,
		},
		{
			name: "Invalid ExeName Backslash",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      "dir\\service.exe",
				DisplayName:  testMyService,
			},
			wantErr:     true,
			errContains: errInvalidExeName,
		},
		{
			name: "Invalid ExeName Shell Metachar",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      "service$.exe",
				DisplayName:  testMyService,
			},
			wantErr:     true,
			errContains: errInvalidExeName,
		},
		{
			name: "Empty ExeName",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      "",
				DisplayName:  testMyService,
			},
			wantErr:     true,
			errContains: errInvalidExeName,
		},
		{
			name: "Invalid DisplayName Shell Metachar",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      testServiceExe,
				DisplayName:  "My Service $",
			},
			wantErr:     true,
			errContains: errInvalidDisplayName,
		},
		{
			name: "Invalid DisplayName Newline",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      testServiceExe,
				DisplayName:  "My Service \n",
			},
			wantErr:     true,
			errContains: errInvalidDisplayName,
		},
		{
			name: "Empty DisplayName",
			variant: Variant{
				RegistryName: testServiceName,
				ExeName:      testServiceExe,
				DisplayName:  "",
			},
			wantErr:     true,
			errContains: errInvalidDisplayName,
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

package typespec_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/stretchr/testify/assert"
)

func TestTypeSpecConfig_GetPackageRelativePath(t *testing.T) {
	tests := []struct {
		name     string
		config   typespec.TypeSpecConfig
		expected string
	}{
		{
			name: "Package path from module",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{
							"module": "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents",
						},
					},
				},
			},
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Package path from module with placeholder",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{
							"module":      "github.com/Azure/azure-sdk-for-go/{service-dir}/armcompute",
							"service-dir": "sdk/resourcemanager/compute",
						},
					},
				},
			},
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Empty package path",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{},
					},
				},
			},
			expected: "",
		},
		{
			name: "Package path from service and package dir",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{
							"module":      "github.com/Azure/azure-sdk-for-go/{service-dir}/azadmin",
							"service-dir": "sdk/security/keyvault",
							"package-dir": "azadmin/backup",
						},
					},
				},
			},
			expected: "sdk/security/keyvault/azadmin/backup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.GetPackageRelativePath())
		})
	}
}

func TestTypeSpecConfig_GetModuleRelativePath(t *testing.T) {
	tests := []struct {
		name     string
		config   typespec.TypeSpecConfig
		expected string
	}{
		{
			name: "Normal",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{
							"module": "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents",
						},
					},
				},
			},
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Module with placeholder",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{
							"module":      "github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}",
							"service-dir": "sdk/resourcemanager/compute",
							"package-dir": "armcompute",
						},
					},
				},
			},
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Empty module",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{},
					},
				},
			},
			expected: "",
		},
		{
			name: "Module different from package path",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{
							"module":      "github.com/Azure/azure-sdk-for-go/{service-dir}/azadmin",
							"service-dir": "sdk/security/keyvault",
							"package-dir": "azadmin/backup",
						},
					},
				},
			},
			expected: "sdk/security/keyvault/azadmin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.GetModuleRelativePath())
		})
	}
}

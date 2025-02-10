// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/stretchr/testify/assert"
)

func TestTypeSpecConfig_GetPackageModuleRelativePath(t *testing.T) {
	tests := []struct {
		name     string
		config   typespec.TypeSpecConfig
		expected string
	}{
		{
			name: "Valid module path",
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
			name: "Module path with placeholders",
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
			name: "Module path without module key",
			config: typespec.TypeSpecConfig{
				TypeSpecProjectSchema: typespec.TypeSpecProjectSchema{
					Options: map[string]any{
						"@azure-tools/typespec-go": map[string]any{},
					},
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.GetPackageModuleRelativePath())
		})
	}
}

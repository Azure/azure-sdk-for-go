// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/stretchr/testify/assert"
)

func TestTypeSpecConfig_GetPackageRelativePath(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected string
	}{
		{
			name: "Package path from module",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents"`,
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Package path from module with placeholder",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/{service-dir}/armcompute"
    service-dir: "sdk/resourcemanager/compute"`,
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Package path from service and package dir",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/{service-dir}/azadmin"
    service-dir: "sdk/security/keyvault"
    package-dir: "azadmin/backup"`,
			expected: "sdk/security/keyvault/azadmin/backup",
		},
		{
			name: "Package path from emitter output dir",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents"
    emitter-output-dir: "{output-dir}/sdk/messaging/eventgrid/azsystemevents"`,
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Package path from emitter output dir with placeholder",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/{service-dir}/azsystemevents"
    service-dir: "sdk/messaging/eventgrid"
    emitter-output-dir: "{output-dir}/{service-dir}/azsystemevents"`,
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Package path from containing module when no package-dir specified",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"`,
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Package path from service-dir and package-dir with containing module",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin"
    service-dir: "sdk/security/keyvault"
    package-dir: "azadmin/backup"`,
			expected: "sdk/security/keyvault/azadmin/backup",
		},
		{
			name: "Package path from emitter output dir with containing module",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
    emitter-output-dir: "{output-dir}/sdk/resourcemanager/compute/armcompute"`,
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Package path from emitter output dir with placeholder and containing module",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/{service-dir}/armnetwork"
    service-dir: "sdk/resourcemanager/network"
    emitter-output-dir: "{output-dir}/{service-dir}/armnetwork"`,
			expected: "sdk/resourcemanager/network/armnetwork",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file with YAML content
			tmpFile, err := os.CreateTemp("", "tspconfig_*.yaml")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.yaml)
			assert.NoError(t, err)
			tmpFile.Close()

			// Parse config from file
			config, err := typespec.ParseTypeSpecConfig(tmpFile.Name())
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, config.GetPackageRelativePath())
		})
	}
}

func TestTypeSpecConfig_GetModuleRelativePath(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected string
	}{
		{
			name: "Normal",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents"`,
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Module with placeholder",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}"
    service-dir: "sdk/resourcemanager/compute"
    package-dir: "armcompute"`,
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Module different from package path",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/{service-dir}/azadmin"
    service-dir: "sdk/security/keyvault"
    package-dir: "azadmin/backup"`,
			expected: "sdk/security/keyvault/azadmin",
		},
		{
			name: "Module with major version suffix removed",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"`,
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Module with v2 major version suffix removed",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents/v2"`,
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Module with double-digit major version suffix removed",
			yaml: `options:
  "@azure-tools/typespec-go":
    module: "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/v12"`,
			expected: "sdk/storage/azblob",
		},
		{
			name: "Containing module takes precedence over module",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"`,
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Containing module with placeholder",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}"
    service-dir: "sdk/resourcemanager/storage"
    package-dir: "armstorage"`,
			expected: "sdk/resourcemanager/storage/armstorage",
		},
		{
			name: "Containing module with major version suffix removed",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"`,
			expected: "sdk/resourcemanager/compute/armcompute",
		},
		{
			name: "Containing module with v2 major version suffix removed",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/azsystemevents/v2"`,
			expected: "sdk/messaging/eventgrid/azsystemevents",
		},
		{
			name: "Containing module with double-digit major version suffix removed",
			yaml: `options:
  "@azure-tools/typespec-go":
    containing-module: "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/v12"`,
			expected: "sdk/storage/azblob",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file with YAML content
			tmpFile, err := os.CreateTemp("", "tspconfig_*.yaml")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.yaml)
			assert.NoError(t, err)
			tmpFile.Close()

			// Parse config from file
			config, err := typespec.ParseTypeSpecConfig(tmpFile.Name())
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, config.GetModuleRelativePath())
		})
	}
}

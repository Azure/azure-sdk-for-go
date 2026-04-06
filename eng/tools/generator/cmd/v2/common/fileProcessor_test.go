// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestFindModule(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	sdkRoot := utils.NormalizePath(cwd)
	sdkRepo, err := repo.OpenSDKRepository(sdkRoot)
	assert.NoError(t, err)
	module, err := FindModuleDirByGoMod(fmt.Sprintf("%s/%s", filepath.ToSlash(sdkRepo.Root()), "sdk/security/keyvault/azadmin/settings"))
	assert.NoError(t, err)
	moduleRelativePath, err := filepath.Rel(sdkRepo.Root(), module)
	assert.NoError(t, err)
	assert.Equal(t, "sdk/security/keyvault/azadmin", filepath.ToSlash(moduleRelativePath))
}

func TestUpdateReadMeClientFactory(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "tmp")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(`
	clientFactory, err := NewClientFactory(<subscription ID>, cred, nil)
	clientFactory, err := NewClientFactory(<subscription ID>, cred, &options)
	`), 0644)
	assert.NoError(t, err)

	// Test case 1: without subscription ID
	err = os.WriteFile(filepath.Join(tmpDir, "client_factory.go"), []byte(`
	func NewClientFactory(credential azcore.TokenCredential, options *arm.ClientOptions) *ClientFactory {
	`), 0644)
	assert.NoError(t, err)

	err = UpdateReadmeClientFactory(tmpDir)
	assert.NoError(t, err)

	noSubContent, err := os.ReadFile(filepath.Join(tmpDir, "README.md"))
	assert.NoError(t, err)
	assert.Contains(t, string(noSubContent), "NewClientFactory(cred, nil)")
	assert.Contains(t, string(noSubContent), "NewClientFactory(cred, &options)")

	// Test case 2: with subscription ID
	err = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(`
	clientFactory, err := NewClientFactory(cred, nil)
	clientFactory, err := NewClientFactory(cred, &options)
	`), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "client_factory.go"), []byte(`
	func NewClientFactory(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *ClientFactory {
	`), 0644)
	assert.NoError(t, err)

	err = UpdateReadmeClientFactory(tmpDir)
	assert.NoError(t, err)

	withSubContent, err := os.ReadFile(filepath.Join(tmpDir, "README.md"))
	assert.NoError(t, err)
	assert.Contains(t, string(withSubContent), `NewClientFactory(<subscription ID>, cred, nil)`)
	assert.Contains(t, string(withSubContent), `NewClientFactory(<subscription ID>, cred, &options)`)
}

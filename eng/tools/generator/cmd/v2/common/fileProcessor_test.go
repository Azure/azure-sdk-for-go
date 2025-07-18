// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateNewVersion(t *testing.T) {
	fixChange := &changelog.Changelog{Modified: &report.Package{}}
	breakingChange := &changelog.Changelog{RemovedPackage: true, Modified: &report.Package{}}
	additiveChange := &changelog.Changelog{Modified: &report.Package{AdditiveChanges: &report.AdditiveChanges{Added: &delta.Content{Content: exports.Content{Consts: map[string]exports.Const{"test": {}}}}}}}

	// previous 0.x.x
	// fix with stable
	newVersion, prl, err := CalculateNewVersion(fixChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, FirstGALabel, prl)

	// fix with beat
	newVersion, prl, err = CalculateNewVersion(fixChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.5.1")
	assert.Equal(t, BetaLabel, prl)

	// breaking with stable
	newVersion, prl, err = CalculateNewVersion(breakingChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, FirstGABreakingChangeLabel, prl)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.6.0")
	assert.Equal(t, BetaBreakingChangeLabel, prl)

	// additive with stable
	newVersion, prl, err = CalculateNewVersion(additiveChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, FirstGALabel, prl)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.6.0")
	assert.Equal(t, BetaLabel, prl)

	// previous 1.2.0
	// fix with stable
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.1")
	assert.Equal(t, StableLabel, prl)

	// fix with beat
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.1-beta.1")
	assert.Equal(t, BetaLabel, prl)

	// breaking with stable
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "2.0.0")
	assert.Equal(t, StableBreakingChangeLabel, prl)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "2.0.0-beta.1")
	assert.Equal(t, BetaBreakingChangeLabel, prl)

	// additive with stable
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.3.0")
	assert.Equal(t, StableLabel, prl)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.3.0-beta.1")
	assert.Equal(t, BetaLabel, prl)

	// previous 1.2.0-beta.1
	// fix with stable
	_, _, err = CalculateNewVersion(fixChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// fix with beat
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, BetaLabel, prl)

	// breaking with stable
	_, _, err = CalculateNewVersion(breakingChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, BetaBreakingChangeLabel, prl)

	// additive with stable
	_, _, err = CalculateNewVersion(additiveChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, BetaLabel, prl)
}

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

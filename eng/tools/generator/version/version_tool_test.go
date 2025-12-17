// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, utils.FirstGALabel, prl)

	// fix with beta
	newVersion, prl, err = CalculateNewVersion(fixChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.5.1")
	assert.Equal(t, utils.BetaLabel, prl)

	// breaking with stable
	newVersion, prl, err = CalculateNewVersion(breakingChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, utils.FirstGABreakingChangeLabel, prl)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.6.0")
	assert.Equal(t, utils.BetaBreakingChangeLabel, prl)

	// additive with stable
	newVersion, prl, err = CalculateNewVersion(additiveChange, "0.5.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.0.0")
	assert.Equal(t, utils.FirstGALabel, prl)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "0.5.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "0.6.0")
	assert.Equal(t, utils.BetaLabel, prl)

	// previous 1.2.0
	// fix with stable
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.1")
	assert.Equal(t, utils.StableLabel, prl)

	// fix with beta
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.1-beta.1")
	assert.Equal(t, utils.BetaLabel, prl)

	// breaking with stable
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "2.0.0")
	assert.Equal(t, utils.StableBreakingChangeLabel, prl)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "2.0.0-beta.1")
	assert.Equal(t, utils.BetaBreakingChangeLabel, prl)

	// additive with stable
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0", false)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.3.0")
	assert.Equal(t, utils.StableLabel, prl)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.3.0-beta.1")
	assert.Equal(t, utils.BetaLabel, prl)

	// previous 1.2.0-beta.1
	// fix with stable
	_, _, err = CalculateNewVersion(fixChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// fix with beta
	newVersion, prl, err = CalculateNewVersion(fixChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, utils.BetaLabel, prl)

	// breaking with stable
	_, _, err = CalculateNewVersion(breakingChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// breaking with beta
	newVersion, prl, err = CalculateNewVersion(breakingChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, utils.BetaBreakingChangeLabel, prl)

	// additive with stable
	_, _, err = CalculateNewVersion(additiveChange, "1.2.0-beta.1", false)
	assert.NotEmpty(t, err)

	// additive with beta
	newVersion, prl, err = CalculateNewVersion(additiveChange, "1.2.0-beta.1", true)
	assert.NoError(t, err)
	assert.Equal(t, newVersion.String(), "1.2.0-beta.2")
	assert.Equal(t, utils.BetaLabel, prl)
}

func TestContainsPreviewAPIVersion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-preview-api")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("Package with preview API", func(t *testing.T) {
		packageDir := filepath.Join(tempDir, "preview")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		goFile := filepath.Join(packageDir, "client.go")
		goContent := `package preview

func NewClient() *Client {
	return &Client{
		options: map[string]string{
			"api-version": "2023-01-01-preview",
		},
	}
}
`
		err = os.WriteFile(goFile, []byte(goContent), 0644)
		require.NoError(t, err)

		hasPreview, err := containsPreviewAPIVersion(packageDir)
		assert.NoError(t, err)
		assert.True(t, hasPreview)
	})

	t.Run("Package without preview API", func(t *testing.T) {
		packageDir := filepath.Join(tempDir, "stable")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		goFile := filepath.Join(packageDir, "client.go")
		goContent := `package stable

func NewClient() *Client {
	return &Client{
		options: map[string]string{
			"api-version": "2023-01-01",
		},
	}
}
`
		err = os.WriteFile(goFile, []byte(goContent), 0644)
		require.NoError(t, err)

		hasPreview, err := containsPreviewAPIVersion(packageDir)
		assert.NoError(t, err)
		assert.False(t, hasPreview)
	})

	t.Run("Package with multiple API versions including preview", func(t *testing.T) {
		packageDir := filepath.Join(tempDir, "mixed")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		goFile := filepath.Join(packageDir, "client.go")
		goContent := `package mixed

func NewClient() *Client {
	return &Client{
		options: map[string]string{
			"api-version": "2023-06-01-preview",
		},
	}
}

func NewAnotherClient() *AnotherClient {
	return &AnotherClient{
		options: map[string]string{
			"api-version": "2023-01-01",
		},
	}
}
`
		err = os.WriteFile(goFile, []byte(goContent), 0644)
		require.NoError(t, err)

		hasPreview, err := containsPreviewAPIVersion(packageDir)
		assert.NoError(t, err)
		assert.True(t, hasPreview)
	})

	t.Run("Empty package directory", func(t *testing.T) {
		packageDir := filepath.Join(tempDir, "empty")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		hasPreview, err := containsPreviewAPIVersion(packageDir)
		assert.NoError(t, err)
		assert.False(t, hasPreview)
	})

	t.Run("Non-existent directory", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "nonexistent")

		hasPreview, err := containsPreviewAPIVersion(nonExistentDir)
		assert.Error(t, err)
		assert.False(t, hasPreview)
	})
}

func TestUpdateImportPaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-update-imports")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock SDK repository structure
	sdkDir := filepath.Join(tempDir, "sdk")
	err = os.MkdirAll(sdkDir, 0755)
	require.NoError(t, err)

	// Mock SDKRepository
	mockRepo := &mockSDKRepo{root: tempDir}

	t.Run("Update imports from v1 to v2", func(t *testing.T) {
		packageDir := filepath.Join(sdkDir, "resourcemanager", "compute", "armcompute")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		// Create a Go file with v1 imports
		goFile := filepath.Join(packageDir, "test.go")
		goContent := `package armcompute

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/fake"
)
`
		err = os.WriteFile(goFile, []byte(goContent), 0644)
		require.NoError(t, err)

		// Create version from semver
		version, err := semver.NewVersion("2.0.0")
		require.NoError(t, err)

		err = UpdateImportPaths(packageDir, version, mockRepo)
		assert.NoError(t, err)

		// Read the updated file
		updatedContent, err := os.ReadFile(goFile)
		require.NoError(t, err)

		// Verify imports were updated to v2
		assert.Contains(t, string(updatedContent), "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v2")
		assert.Contains(t, string(updatedContent), "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v2/fake")
		assert.NotContains(t, string(updatedContent), "\"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute\"")
	})

	t.Run("Update imports from v2 to v3", func(t *testing.T) {
		packageDir := filepath.Join(sdkDir, "resourcemanager", "storage", "armstorage")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		// Create a Go file with v2 imports
		goFile := filepath.Join(packageDir, "test.go")
		goContent := `package armstorage

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v2/fake"
)
`
		err = os.WriteFile(goFile, []byte(goContent), 0644)
		require.NoError(t, err)

		// Create version from semver
		version, err := semver.NewVersion("3.0.0")
		require.NoError(t, err)

		err = UpdateImportPaths(packageDir, version, mockRepo)
		assert.NoError(t, err)

		// Read the updated file
		updatedContent, err := os.ReadFile(goFile)
		require.NoError(t, err)

		// Verify imports were updated to v3
		assert.Contains(t, string(updatedContent), "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v3")
		assert.Contains(t, string(updatedContent), "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v3/fake")
		assert.NotContains(t, string(updatedContent), "/v2")
	})

	t.Run("No changes for v1 modules", func(t *testing.T) {
		packageDir := filepath.Join(sdkDir, "data", "azcosmos")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		// Create a Go file with imports
		goFile := filepath.Join(packageDir, "test.go")
		originalContent := `package azcosmos

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)
`
		err = os.WriteFile(goFile, []byte(originalContent), 0644)
		require.NoError(t, err)

		// Create v1 version
		version, err := semver.NewVersion("1.5.0")
		require.NoError(t, err)

		err = UpdateImportPaths(packageDir, version, mockRepo)
		assert.NoError(t, err)

		// Read the file - should be unchanged
		updatedContent, err := os.ReadFile(goFile)
		require.NoError(t, err)

		// Verify no changes for v1
		assert.Equal(t, originalContent, string(updatedContent))
	})

	t.Run("Update multiple files in package", func(t *testing.T) {
		packageDir := filepath.Join(sdkDir, "resourcemanager", "network", "armnetwork")
		err = os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		// Create multiple Go files
		clientFile := filepath.Join(packageDir, "client.go")
		clientContent := `package armnetwork

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

type Client struct {}
`
		err = os.WriteFile(clientFile, []byte(clientContent), 0644)
		require.NoError(t, err)

		modelsFile := filepath.Join(packageDir, "models.go")
		modelsContent := `package armnetwork

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/fake"
)

type Model struct {}
`
		err = os.WriteFile(modelsFile, []byte(modelsContent), 0644)
		require.NoError(t, err)

		// Create version from semver
		version, err := semver.NewVersion("4.0.0")
		require.NoError(t, err)

		err = UpdateImportPaths(packageDir, version, mockRepo)
		assert.NoError(t, err)

		// Verify both files were updated
		updatedClient, err := os.ReadFile(clientFile)
		require.NoError(t, err)
		assert.Contains(t, string(updatedClient), "armnetwork/v4")

		updatedModels, err := os.ReadFile(modelsFile)
		require.NoError(t, err)
		assert.Contains(t, string(updatedModels), "armnetwork/v4")
		assert.Contains(t, string(updatedModels), "armnetwork/v4/fake")
	})

	t.Run("Handle subpackage imports correctly", func(t *testing.T) {
		packageDir := filepath.Join(sdkDir, "messaging", "azservicebus")
		subDir := filepath.Join(packageDir, "internal")
		err = os.MkdirAll(subDir, 0755)
		require.NoError(t, err)

		// Create a Go file with subpackage imports
		goFile := filepath.Join(packageDir, "test.go")
		goContent := `package azservicebus

import (
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
)
`
		err = os.WriteFile(goFile, []byte(goContent), 0644)
		require.NoError(t, err)

		// Create version from semver
		version, err := semver.NewVersion("2.0.0")
		require.NoError(t, err)

		err = UpdateImportPaths(packageDir, version, mockRepo)
		assert.NoError(t, err)

		// Read the updated file
		updatedContent, err := os.ReadFile(goFile)
		require.NoError(t, err)

		// Verify both base and subpackage imports were updated
		assert.Contains(t, string(updatedContent), "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/v2")
		assert.Contains(t, string(updatedContent), "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/v2/internal")
	})
}

// mockSDKRepo is a mock implementation of SDKRepository for testing
type mockSDKRepo struct {
	root string
}

func (m *mockSDKRepo) Root() string {
	return m.root
}

func (m *mockSDKRepo) CreateReleaseBranch(releaseBranchName string) error {
	return nil
}

func (m *mockSDKRepo) AddReleaseCommit(rpName, namespaceName, specHash, version string) error {
	return nil
}

func (m *mockSDKRepo) Add(path string) error {
	return nil
}

func (m *mockSDKRepo) Commit(message string) error {
	return nil
}

func (m *mockSDKRepo) Checkout(opt *repo.CheckoutOptions) error {
	return nil
}

func (m *mockSDKRepo) CheckoutTag(tag string) error {
	return nil
}

func (m *mockSDKRepo) CreateBranch(branch *repo.Branch) error {
	return nil
}

func (m *mockSDKRepo) DeleteBranch(name string) error {
	return nil
}

func (m *mockSDKRepo) CherryPick(commit string) error {
	return nil
}

func (m *mockSDKRepo) Stash() error {
	return nil
}

func (m *mockSDKRepo) StashPop() error {
	return nil
}

func (m *mockSDKRepo) Head() (*plumbing.Reference, error) {
	return nil, nil
}

func (m *mockSDKRepo) Tags() (storer.ReferenceIter, error) {
	return nil, nil
}

func (m *mockSDKRepo) Remotes() ([]*git.Remote, error) {
	return nil, nil
}

func (m *mockSDKRepo) DeleteRemote(name string) error {
	return nil
}

func (m *mockSDKRepo) CreateRemote(c *config.RemoteConfig) (*git.Remote, error) {
	return nil, nil
}

func (m *mockSDKRepo) Fetch(o *git.FetchOptions) error {
	return nil
}

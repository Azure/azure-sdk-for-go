// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/delta"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/report"
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

	// fix with beat
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

	// fix with beat
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

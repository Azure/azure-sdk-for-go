// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package metadata

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/stretchr/testify/require"
)

func TestInferPackageTitle(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		expected    string
	}{
		{name: "ArmPrefix", packageName: "armcompute", expected: "Compute"},
		{name: "AzPrefix", packageName: "azblob", expected: "Blob"},
		{name: "NoPrefix", packageName: "storage", expected: "Storage"},
		{name: "Empty", packageName: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferPackageTitle(tt.packageName)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestParseSingleTemplate_CIYml(t *testing.T) {
	// Create a temporary template file
	tmpDir := t.TempDir()
	tplPath := filepath.Join(tmpDir, "ci.yml.tpl")
	tplContent := `# NOTE: Please refer to https://aka.ms/azsdk/engsys/ci-yaml before editing this file.
trigger:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

pr:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

extends:
  template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml
  parameters:
    ServiceDirectory: '{{.serviceDir}}'
`
	err := os.WriteFile(tplPath, []byte(tplContent), 0644)
	require.NoError(t, err)

	outputPath := filepath.Join(tmpDir, "ci.yml")
	data := map[string]string{
		"moduleRelativePath": "sdk/resourcemanager/compute/armcompute",
		"serviceDir":         "resourcemanager/compute/armcompute",
	}

	err = typespec.ParseSingleTemplate(tplPath, outputPath, data)
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, "# NOTE: Please refer to https://aka.ms/azsdk/engsys/ci-yaml before editing this file.")
	require.Contains(t, contentStr, "sdk/resourcemanager/compute/armcompute/")
	require.Contains(t, contentStr, "ServiceDirectory: 'resourcemanager/compute/armcompute'")
	require.Contains(t, contentStr, "template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml")
}

func TestParseSingleTemplate_DataPlane(t *testing.T) {
	tmpDir := t.TempDir()
	tplPath := filepath.Join(tmpDir, "ci.yml.tpl")
	tplContent := `trigger:
  paths:
    include:
    - {{.moduleRelativePath}}/
extends:
  template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml
  parameters:
    ServiceDirectory: '{{.serviceDir}}'
`
	err := os.WriteFile(tplPath, []byte(tplContent), 0644)
	require.NoError(t, err)

	outputPath := filepath.Join(tmpDir, "ci.yml")
	data := map[string]string{
		"moduleRelativePath": "sdk/messaging/azeventhubs",
		"serviceDir":         "messaging/azeventhubs",
	}

	err = typespec.ParseSingleTemplate(tplPath, outputPath, data)
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, "sdk/messaging/azeventhubs/")
	require.Contains(t, contentStr, "ServiceDirectory: 'messaging/azeventhubs'")
}

func TestProcessMetadata_MgmtPlaneNewPackage(t *testing.T) {
	// Create a mock SDK root with .git directory
	sdkRoot := t.TempDir()
	err := os.MkdirAll(filepath.Join(sdkRoot, ".git"), 0755)
	require.NoError(t, err)

	// Create the template directory with ci.yml.tpl and README.md.tpl
	templateDir := filepath.Join(sdkRoot, "eng", "tools", "generator", "template", "typespec")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	ciTpl := `# NOTE: Please refer to https://aka.ms/azsdk/engsys/ci-yaml before editing this file.
trigger:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

pr:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

extends:
  template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml
  parameters:
    ServiceDirectory: '{{.serviceDir}}'
`
	err = os.WriteFile(filepath.Join(templateDir, "ci.yml.tpl"), []byte(ciTpl), 0644)
	require.NoError(t, err)

	readmeTpl := `# Azure {{.packageTitle}} Module for Go

The ` + "`{{.packageName}}`" + ` module provides operations for working with Azure {{.packageTitle}}.

[Source code](https://github.com/Azure/azure-sdk-for-go/tree/main/{{.moduleRelativePath}})
`
	err = os.WriteFile(filepath.Join(templateDir, "README.md.tpl"), []byte(readmeTpl), 0644)
	require.NoError(t, err)

	// Create the package directory
	packageDir := filepath.Join(sdkRoot, "sdk", "resourcemanager", "compute", "armcompute")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"), 0644)
	require.NoError(t, err)

	// Process metadata
	result, err := processMetadata(packageDir, "Compute", false)
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Contains(t, result.CreatedFiles, "ci.yml")
	require.Contains(t, result.CreatedFiles, "README.md")

	// Verify ci.yml was created
	ciContent, err := os.ReadFile(filepath.Join(packageDir, "ci.yml"))
	require.NoError(t, err)
	require.Contains(t, string(ciContent), "sdk/resourcemanager/compute/armcompute/")
	require.Contains(t, string(ciContent), "ServiceDirectory: 'resourcemanager/compute/armcompute'")

	// Verify README.md was created
	readmeContent, err := os.ReadFile(filepath.Join(packageDir, "README.md"))
	require.NoError(t, err)
	require.Contains(t, string(readmeContent), "Azure Compute Module for Go")
	require.Contains(t, string(readmeContent), "armcompute")
}

func TestProcessMetadata_DataPlaneNewPackage(t *testing.T) {
	// Create a mock SDK root with .git directory
	sdkRoot := t.TempDir()
	err := os.MkdirAll(filepath.Join(sdkRoot, ".git"), 0755)
	require.NoError(t, err)

	// Create the template directory with ci.yml.tpl
	templateDir := filepath.Join(sdkRoot, "eng", "tools", "generator", "template", "typespec")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	ciTpl := `# NOTE: Please refer to https://aka.ms/azsdk/engsys/ci-yaml before editing this file.
trigger:
  paths:
    include:
    - {{.moduleRelativePath}}/
pr:
  paths:
    include:
    - {{.moduleRelativePath}}/
extends:
  template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml
  parameters:
    ServiceDirectory: '{{.serviceDir}}'
`
	err = os.WriteFile(filepath.Join(templateDir, "ci.yml.tpl"), []byte(ciTpl), 0644)
	require.NoError(t, err)

	// Create the package directory (data plane - no README.md should be created)
	packageDir := filepath.Join(sdkRoot, "sdk", "storage", "azblob")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"), 0644)
	require.NoError(t, err)

	// Process metadata
	result, err := processMetadata(packageDir, "", false)
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Contains(t, result.CreatedFiles, "ci.yml")
	// README.md should NOT be created for data plane packages
	for _, f := range result.CreatedFiles {
		require.NotEqual(t, "README.md", f)
	}

	// Verify ci.yml was created
	ciContent, err := os.ReadFile(filepath.Join(packageDir, "ci.yml"))
	require.NoError(t, err)
	require.Contains(t, string(ciContent), "sdk/storage/azblob/")
	require.Contains(t, string(ciContent), "ServiceDirectory: 'storage/azblob'")
}

func TestProcessMetadata_ExistingFiles(t *testing.T) {
	// Create a mock SDK root with .git directory
	sdkRoot := t.TempDir()
	err := os.MkdirAll(filepath.Join(sdkRoot, ".git"), 0755)
	require.NoError(t, err)

	// Create the template directory with templates
	templateDir := filepath.Join(sdkRoot, "eng", "tools", "generator", "template", "typespec")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(templateDir, "ci.yml.tpl"), []byte("# {{.serviceDir}}"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(templateDir, "README.md.tpl"), []byte("# {{.packageTitle}}"), 0644)
	require.NoError(t, err)

	// Create the package directory with existing files
	packageDir := filepath.Join(sdkRoot, "sdk", "resourcemanager", "compute", "armcompute")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module test"), 0644)
	require.NoError(t, err)

	existingCIContent := "# existing ci.yml"
	err = os.WriteFile(filepath.Join(packageDir, "ci.yml"), []byte(existingCIContent), 0644)
	require.NoError(t, err)

	existingReadmeContent := "# existing README"
	err = os.WriteFile(filepath.Join(packageDir, "README.md"), []byte(existingReadmeContent), 0644)
	require.NoError(t, err)

	// Process metadata - should skip existing files
	result, err := processMetadata(packageDir, "", false)
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Empty(t, result.CreatedFiles)
	require.Contains(t, result.SkippedFiles, "ci.yml")
	require.Contains(t, result.SkippedFiles, "README.md")

	// Verify existing files were not modified
	ciContent, err := os.ReadFile(filepath.Join(packageDir, "ci.yml"))
	require.NoError(t, err)
	require.Equal(t, existingCIContent, string(ciContent))

	readmeContent, err := os.ReadFile(filepath.Join(packageDir, "README.md"))
	require.NoError(t, err)
	require.Equal(t, existingReadmeContent, string(readmeContent))
}

func TestProcessMetadata_NoSDKRoot(t *testing.T) {
	tmpDir := t.TempDir()
	packageDir := filepath.Join(tmpDir, "pkg")
	err := os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module test"), 0644)
	require.NoError(t, err)

	result, err := processMetadata(packageDir, "", false)
	require.NoError(t, err)
	require.False(t, result.Success)
	require.Contains(t, result.Message, "Failed to find SDK root")
}

func TestCIYmlTemplate_Format(t *testing.T) {
	// Create a template file matching the real ci.yml.tpl
	tmpDir := t.TempDir()
	tplPath := filepath.Join(tmpDir, "ci.yml.tpl")
	tplContent := `# NOTE: Please refer to https://aka.ms/azsdk/engsys/ci-yaml before editing this file.
trigger:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

pr:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

extends:
  template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml
  parameters:
    ServiceDirectory: '{{.serviceDir}}'
`
	err := os.WriteFile(tplPath, []byte(tplContent), 0644)
	require.NoError(t, err)

	outputPath := filepath.Join(tmpDir, "ci.yml")
	data := map[string]string{
		"moduleRelativePath": "sdk/resourcemanager/network/armnetwork",
		"serviceDir":         "resourcemanager/network/armnetwork",
	}

	err = typespec.ParseSingleTemplate(tplPath, outputPath, data)
	require.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")

	// Verify the overall structure
	require.True(t, strings.HasPrefix(lines[0], "# NOTE:"))

	// Count occurrences of the path
	pathCount := strings.Count(string(content), "sdk/resourcemanager/network/armnetwork/")
	require.Equal(t, 2, pathCount, "Path should appear twice (trigger and PR sections)")

	serviceDirCount := strings.Count(string(content), "ServiceDirectory: 'resourcemanager/network/armnetwork'")
	require.Equal(t, 1, serviceDirCount, "ServiceDirectory should appear once")
}

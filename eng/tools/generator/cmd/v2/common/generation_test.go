// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/require"
)

func TestTypeSpecOnBoardGeneratorPreservesReadme(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	sdkRepo, err := repo.OpenSDKRepository(cwd)
	require.NoError(t, err)
	modulePath := t.TempDir()
	generator := TypeSpecOnBoardGenerator{
		TypeSpecCommonGenerator: &TypeSpecCommonGenerator{
			GenerateContext:    &GenerateContext{SDKPath: sdkRepo.Root()},
			ModulePath:         modulePath,
			ModuleRelativePath: "sdk/resourcemanager/networksecurity/armnetworksecurity",
			Version:            semver.MustParse("0.1.0"),
		},
	}
	params := &GenerateParam{RPName: "networksecurity", NamespaceName: "armnetworksecurity"}
	require.NoError(t, generator.PreGenerate(params))
	readmePath := filepath.Join(modulePath, "README.md")
	readme, err := os.ReadFile(readmePath)
	require.NoError(t, err)
	require.Contains(t, string(readme), "# Azure Networksecurity Module for Go")

	customReadme := "# Azure Network Security Module for Go\n\n" +
		"The `armnetworksecurity` module provides operations for working with Azure Network Security.\n\n" +
		"go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/networksecurity/armnetworksecurity\n"
	require.NoError(t, os.WriteFile(readmePath, []byte(customReadme), 0644))
	for range 2 {
		require.NoError(t, generator.PreGenerate(params))
		readme, err = os.ReadFile(readmePath)
		require.NoError(t, err)
		require.Equal(t, customReadme, string(readme))
	}
}

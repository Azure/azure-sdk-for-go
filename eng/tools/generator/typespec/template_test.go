// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTypeSpecTemplatesPreservesReadme(t *testing.T) {
	templateDir := t.TempDir()
	outputDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "README.md.tpl"), []byte("# Azure {{.packageTitle}} Module for Go\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "version.go.tpl"), []byte("{{.packageVersion}}"), 0644))

	data := map[string]any{"packageTitle": "Networksecurity", "packageVersion": "0.1.0", "releaseDate": ""}
	require.NoError(t, ParseTypeSpecTemplates(templateDir, outputDir, data, nil))
	readmePath := filepath.Join(outputDir, "README.md")
	readme, err := os.ReadFile(readmePath)
	require.NoError(t, err)
	require.Equal(t, "# Azure Networksecurity Module for Go\n", string(readme))

	customReadme := "# Azure Network Security Module for Go\n\nCustomer-authored documentation.\n"
	require.NoError(t, os.WriteFile(readmePath, []byte(customReadme), 0644))
	data["packageVersion"] = "0.2.0"
	for range 2 {
		require.NoError(t, ParseTypeSpecTemplates(templateDir, outputDir, data, nil))
		readme, err = os.ReadFile(readmePath)
		require.NoError(t, err)
		require.Equal(t, customReadme, string(readme))
		version, err := os.ReadFile(filepath.Join(outputDir, "version.go"))
		require.NoError(t, err)
		require.Equal(t, "0.2.0", string(version))
	}
}

func TestParseTypeSpecTemplatesReadmeIsDirectory(t *testing.T) {
	templateDir := t.TempDir()
	outputDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "README.md.tpl"), []byte("README"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(outputDir, "README.md"), 0755))

	err := ParseTypeSpecTemplates(templateDir, outputDir, map[string]any{"releaseDate": ""}, nil)
	require.Error(t, err)
}

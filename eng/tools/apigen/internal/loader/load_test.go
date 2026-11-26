// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"apigen/internal/loader"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := writeTestModule(t, map[string]string{
		"go.mod": "module example.com/sample\n\ngo 1.26.0\n",
		"sample.go": `package sample

import "example.com/sample/internal/model"

type Item = model.Item
`,
		"sample_test.go": `package sample

var TestOnly = undefinedInTests
`,
		"excluded.go": `//go:build ignore

package sample

var Excluded = undefinedInExcludedFile
`,
		"internal/model/model.go": `package model

type Item struct {
	Name string
}
`,
		"sub/sub.go":        "package sub\n\nconst Name = \"sub\"\n",
		"nested/go.mod":     "module example.com/nested\n\ngo 1.26.0\n",
		"nested/invalid.go": "package nested\n\nvar Value = undefinedInNestedModule\n",
	})

	module, err := loader.Load(t.Context(), dir)
	require.NoError(t, err)
	assert.Equal(t, dir, module.Dir)
	assert.Equal(t, "example.com/sample", module.ModFile.Module.Mod.Path)
	require.NotNil(t, module.ModFile.Go)
	assert.Equal(t, "1.26.0", module.ModFile.Go.Version)

	var paths []string
	for _, pkg := range module.Packages {
		paths = append(paths, pkg.PkgPath)
		assert.NotEmpty(t, pkg.Syntax, "syntax missing for %s", pkg.PkgPath)
		require.NotNil(t, pkg.Types, "types missing for %s", pkg.PkgPath)
		assert.NotNil(t, pkg.TypesInfo, "type information missing for %s", pkg.PkgPath)
	}
	require.Equal(t, []string{
		"example.com/sample",
		"example.com/sample/internal/model",
		"example.com/sample/sub",
	}, paths)

	root := module.Packages[0]
	assert.NotNil(t, root.Types.Scope().Lookup("Item"))
	assert.Nil(t, root.Types.Scope().Lookup("TestOnly"))
	assert.Nil(t, root.Types.Scope().Lookup("Excluded"))
	dependency := root.Imports["example.com/sample/internal/model"]
	require.NotNil(t, dependency)
	assert.NotEmpty(t, dependency.Syntax)
	assert.NotNil(t, dependency.TypesInfo)
}

func TestLoadErrors(t *testing.T) {
	t.Setenv("GOWORK", "off")
	tests := []struct {
		name    string
		files   map[string]string
		message string
	}{
		{
			name:    "missing module file",
			files:   map[string]string{},
			message: "read module file",
		},
		{
			name:    "missing module directive",
			files:   map[string]string{"go.mod": "go 1.26.0\n"},
			message: "module directive is missing",
		},
		{
			name:    "invalid module file",
			files:   map[string]string{"go.mod": "module example.com/sample\n\ngo invalid\n"},
			message: "parse module file",
		},
		{
			name:    "no packages",
			files:   map[string]string{"go.mod": "module example.com/sample\n\ngo 1.26.0\n"},
			message: "no non-test Go packages",
		},
		{
			name: "type error",
			files: map[string]string{
				"go.mod":    "module example.com/sample\n\ngo 1.26.0\n",
				"sample.go": "package sample\n\nvar Value = undefinedValue\n",
			},
			message: "undefinedValue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := writeTestModule(t, tt.files)
			module, err := loader.Load(t.Context(), dir)
			require.ErrorContains(t, err, tt.message)
			assert.Nil(t, module)
		})
	}
}

func writeTestModule(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	return dir
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"apigen/internal/loader"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestModuleVersion(t *testing.T) {
	for _, tt := range []struct {
		name     string
		files    map[string]string
		expected string
		message  string
	}{
		{
			name: "module constant in a block",
			files: map[string]string{"constants.go": `package sample
const (
	moduleVersion = "v1.2.3"
	ServerVersion = "v8.0.21"
)
`},
			expected: "1.2.3",
		},
		{
			name: "typed prerelease",
			files: map[string]string{"version.go": `package sample
const version string = "v1.2.3-beta.1"
`},
			expected: "1.2.3-beta.1",
		},
		{
			name: "service library version",
			files: map[string]string{"version.go": `package sample
const serviceLibVersion = "v1.6.0-beta.4"
`},
			expected: "1.6.0-beta.4",
		},
		{
			name: "exported module version",
			files: map[string]string{"version.go": `package sample
const ModuleVersion = "v1.9.0-beta.1"
`},
			expected: "1.9.0-beta.1",
		},
		{
			name: "exported version with build metadata",
			files: map[string]string{"constants.go": `package sample
const Version = "v2.0.0+build.123"
`},
			expected: "2.0.0+build.123",
		},
		{
			name: "internal constant and repeated identical version",
			files: map[string]string{
				"version.go": `package sample
const version = "v1.2.3"
`,
				"internal/shared/constants.go": `package shared
const Version = "v1.2.3"
`,
			},
			expected: "1.2.3",
		},
		{
			name: "unrelated declarations do not supply a version",
			files: map[string]string{
				"constants.go": `package sample
const ServerVersion = "v8.0.21"
const text = "moduleVersion = \"v9.9.9\""
// moduleVersion = "v9.9.9"
func f() { const version = "v9.9.9"; _ = version }
`,
				"sample.go": `package sample
const version = "v9.9.9"
`,
			},
			message: "no module version constant",
		},
		{
			name: "missing prefix",
			files: map[string]string{"version.go": `package sample
const version = "1.2.3"
`},
			message: `should be "v1.2.3", not "1.2.3"`,
		},
		{
			name: "invalid semantic version",
			files: map[string]string{"version.go": `package sample
const version = "v1.2.03"
`},
			message: "must have the form vMAJOR.MINOR.PATCH",
		},
		{
			name: "abbreviated version",
			files: map[string]string{"version.go": `package sample
const version = "v1.2"
`},
			message: "must have the form vMAJOR.MINOR.PATCH",
		},
		{
			name: "nonstring version",
			files: map[string]string{"version.go": `package sample
const version = 123
`},
			message: "must be a string",
		},
		{
			name: "variable is not a constant",
			files: map[string]string{"version.go": `package sample
var version = "v1.2.3"
`},
			message: "no module version constant",
		},
		{
			name: "conflicting versions",
			files: map[string]string{
				"version.go": `package sample
const version = "v1.2.3"
`,
				"internal/shared/constants.go": `package shared
const Version = "v9.9.9"
`,
			},
			message: "conflicting module versions",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			loaded := versionTestModule(t, tt.files)
			version, err := moduleVersion(loaded)
			if tt.message != "" {
				require.ErrorContains(t, err, tt.message)
				assert.Empty(t, version)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, version)
			}
		})
	}
}

func versionTestModule(t *testing.T, files map[string]string) *loader.Result {
	t.Helper()
	dir := t.TempDir()
	fileSet := token.NewFileSet()
	syntax := make(map[string][]*ast.File)
	for filename, source := range files {
		path := filepath.Join(dir, filepath.FromSlash(filename))
		file, err := parser.ParseFile(fileSet, path, source, 0)
		require.NoError(t, err)
		syntax[filepath.Dir(path)] = append(syntax[filepath.Dir(path)], file)
	}
	loaded := &loader.Result{Dir: dir}
	for directory, files := range syntax {
		pkg, err := new(types.Config).Check("example.com/"+filepath.Base(directory), fileSet, files, nil)
		require.NoError(t, err)
		loaded.Packages = append(loaded.Packages, &packages.Package{Types: pkg, Fset: fileSet})
	}
	return loaded
}

func TestPackageRelativePath(t *testing.T) {
	for _, gitDirectory := range []bool{false, true} {
		name := "Git worktree file"
		if gitDirectory {
			name = "Git directory"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			module := filepath.Join(root, "sdk", "service", "module")
			require.NoError(t, os.MkdirAll(module, 0o755))
			if gitDirectory {
				require.NoError(t, os.Mkdir(filepath.Join(root, ".git"), 0o755))
			} else {
				require.NoError(t, os.WriteFile(filepath.Join(root, ".git"), []byte("gitdir: elsewhere\n"), 0o644))
			}
			path, err := packageRelativePath(module, "")
			require.NoError(t, err)
			assert.Equal(t, "sdk/service/module", path)
		})
	}
	t.Run("explicit root without Git", func(t *testing.T) {
		root := t.TempDir()
		module := filepath.Join(root, "sdk", "module")
		require.NoError(t, os.MkdirAll(module, 0o755))
		path, err := packageRelativePath(module, root)
		require.NoError(t, err)
		assert.Equal(t, "sdk/module", path)
		path, err = packageRelativePath(root, root)
		require.NoError(t, err)
		assert.Equal(t, ".", path)
	})
	t.Run("outside root", func(t *testing.T) {
		_, err := packageRelativePath(t.TempDir(), t.TempDir())
		require.ErrorContains(t, err, "outside repository root")
	})
	t.Run("no Git root", func(t *testing.T) {
		_, err := packageRelativePath(t.TempDir(), "")
		require.ErrorContains(t, err, "specify --repo-root")
	})
	t.Run("missing explicit root", func(t *testing.T) {
		_, err := packageRelativePath(t.TempDir(), filepath.Join(t.TempDir(), "missing"))
		require.ErrorContains(t, err, "resolve repository root")
	})
}

func TestReviewMetadata(t *testing.T) {
	loaded := versionTestModule(t, map[string]string{
		"version.go": "package sample\nconst version = \"v1.2.3\"\n",
	})
	metadata, err := reviewMetadata(loaded, loaded.Dir, []byte("abc"))
	require.NoError(t, err)
	assert.Equal(t, "apiHash: ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad\n"+
		"packageVersion: 1.2.3\npackageRelativePath: \".\"\n"+
		"apigenVersion: "+apigenVersion+"\ngoVersion: "+strings.TrimPrefix(runtime.Version(), "go")+"\n", string(metadata))

	changed, err := reviewMetadata(loaded, loaded.Dir, []byte("abc\n"))
	require.NoError(t, err)
	assert.NotEqual(t, metadata, changed, "the hash must include the exact report bytes, including the final newline")
}

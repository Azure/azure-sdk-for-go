// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := filepath.Join(t.TempDir(), "module with spaces")
	require.NoError(t, os.Mkdir(dir, 0o755))
	require.NoError(t, os.Mkdir(filepath.Join(filepath.Dir(dir), ".git"), 0o755))
	require.NoError(t, os.CopyFS(dir, fstest.MapFS{
		"go.mod":           {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
		"sample.go":        {Data: []byte("package sample\n\ntype Client struct{ Name string }\n")},
		"version.go":       {Data: []byte("package sample\n\nconst version = \"v1.2.3-beta.1\"\n")},
		"api.md":           {Data: []byte("existing draft\n")},
		"api.metadata.yml": {Data: []byte("existing metadata\n")},
	}))
	var expected string
	t.Run("default output in current directory", func(t *testing.T) {
		t.Chdir(t.TempDir())
		command := newRootCommand()
		command.SetArgs([]string{dir})
		var output, diagnostics bytes.Buffer
		command.SetOut(&output)
		command.SetErr(&diagnostics)
		require.NoError(t, command.ExecuteContext(t.Context()))
		assert.Empty(t, output.String())
		content, err := os.ReadFile("api.md")
		require.NoError(t, err)
		expected = string(content)
		assert.Contains(t, expected, "# sample\n")
		assert.Contains(t, expected, "## Package sample\n")
		assert.Contains(t, expected, "type Client struct")
		assert.Empty(t, diagnostics.String())
		generatedMetadata, err := os.ReadFile("api.metadata.yml")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("apiHash: %x\npackageVersion: 1.2.3-beta.1\npackageRelativePath: \"module with spaces\"\napigenVersion: %s\ngoVersion: %s\n",
			sha256.Sum256(content), apigenVersion, strings.TrimPrefix(runtime.Version(), "go")), string(generatedMetadata))
		draft, err := os.ReadFile(filepath.Join(dir, "api.md"))
		require.NoError(t, err)
		assert.Equal(t, "existing draft\n", string(draft))
		metadata, err := os.ReadFile(filepath.Join(dir, "api.metadata.yml"))
		require.NoError(t, err)
		assert.Equal(t, "existing metadata\n", string(metadata))
	})
	t.Run("replace bundle in relative output directory", func(t *testing.T) {
		t.Chdir(t.TempDir())
		directory := "review bundle"
		require.NoError(t, os.Mkdir(directory, 0o755))
		destination := filepath.Join(directory, "api.md")
		require.NoError(t, os.WriteFile(destination, []byte("old report\n"), 0o640))
		metadataPath := filepath.Join(filepath.Dir(destination), "api.metadata.yml")
		require.NoError(t, os.WriteFile(metadataPath, []byte("old metadata\n"), 0o600))
		original, err := os.Stat(destination)
		require.NoError(t, err)
		originalMetadata, err := os.Stat(metadataPath)
		require.NoError(t, err)
		command := newRootCommand()
		command.SetArgs([]string{dir, "-o", directory})
		var output bytes.Buffer
		command.SetOut(&output)
		require.NoError(t, command.ExecuteContext(t.Context()))
		assert.Empty(t, output.String())
		report, err := os.ReadFile(destination)
		require.NoError(t, err)
		assert.Equal(t, expected, string(report))
		updated, err := os.Stat(destination)
		require.NoError(t, err)
		assert.Equal(t, original.Mode().Perm(), updated.Mode().Perm())
		metadata, err := os.ReadFile(metadataPath)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("apiHash: %x\npackageVersion: 1.2.3-beta.1\npackageRelativePath: \"module with spaces\"\napigenVersion: %s\ngoVersion: %s\n",
			sha256.Sum256([]byte(expected)), apigenVersion, strings.TrimPrefix(runtime.Version(), "go")), string(metadata))
		updatedMetadata, err := os.Stat(metadataPath)
		require.NoError(t, err)
		assert.Equal(t, originalMetadata.Mode().Perm(), updatedMetadata.Mode().Perm())
		files, err := filepath.Glob(filepath.Join(filepath.Dir(destination), ".apigen-*.tmp"))
		require.NoError(t, err)
		assert.Empty(t, files)
	})
	t.Run("explicit repository root", func(t *testing.T) {
		directory := t.TempDir()
		command := newRootCommand()
		command.SetArgs([]string{dir, "--output", directory, "--repo-root", dir})
		require.NoError(t, command.ExecuteContext(t.Context()))
		metadata, err := os.ReadFile(filepath.Join(directory, "api.metadata.yml"))
		require.NoError(t, err)
		assert.Contains(t, string(metadata), "packageRelativePath: \".\"\n")
	})
}

func TestCommandErrors(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    []string
		message string
	}{
		{name: "missing module", args: []string{}, message: "accepts 1 arg(s), received 0"},
		{name: "extra module", args: []string{"first", "second"}, message: "accepts 1 arg(s), received 2"},
		{name: "unknown flag", args: []string{"--unknown"}, message: "unknown flag"},
		{name: "empty output", args: []string{"module", "--output="}, message: "output must be"},
		{name: "stdout output", args: []string{"module", "--output", "-"}, message: "stdout output is not supported"},
		{name: "missing module file", args: []string{t.TempDir()}, message: "read module file"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			command := newRootCommand()
			command.SetArgs(tt.args)
			var output bytes.Buffer
			command.SetOut(&output)
			require.ErrorContains(t, command.ExecuteContext(t.Context()), tt.message)
			assert.Empty(t, output.String())
		})
	}
}

func TestBundleFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not available on Windows")
	}
	dir := t.TempDir()
	require.NoError(t, writeBundleFiles(dir, []byte("report"), []byte("metadata")))
	for _, name := range []string{"api.md", "api.metadata.yml"} {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.Zero(t, info.Mode().Perm()&^0o600, "%s has mode %04o", name, info.Mode().Perm())
	}

	report := filepath.Join(dir, "api.md")
	metadata := filepath.Join(dir, "api.metadata.yml")
	require.NoError(t, os.Chmod(report, 0o640))
	require.NoError(t, os.Chmod(metadata, 0o400))
	require.NoError(t, writeBundleFiles(dir, []byte("updated report"), []byte("updated metadata")))

	reportInfo, err := os.Stat(report)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o640), reportInfo.Mode().Perm())
	metadataInfo, err := os.Stat(metadata)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o400), metadataInfo.Mode().Perm())
}

func TestBundlePreparationFailurePreservesDestination(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "api.md")
	require.NoError(t, os.WriteFile(path, []byte("existing report\n"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "api.metadata.yml"), 0o755))
	err := writeBundleFiles(dir, []byte("new report\n"), []byte("new metadata\n"))
	require.ErrorContains(t, err, "not a regular file")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "existing report\n", string(content))
	files, err := filepath.Glob(filepath.Join(dir, ".apigen-*.tmp"))
	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestReportFileRejectsDirectory(t *testing.T) {
	directory := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(directory, "api.md"), 0o755))
	err := writeBundleFiles(directory, nil, nil)
	require.ErrorContains(t, err, "not a regular file")
}

func TestBundleOutputDirectoryErrors(t *testing.T) {
	t.Run("file path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "api.md")
		require.NoError(t, os.WriteFile(path, []byte("existing report\n"), 0o644))
		require.ErrorContains(t, writeBundleFiles(path, nil, nil), "not a directory")
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "existing report\n", string(content))
	})
	t.Run("missing directory", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing")
		err := writeBundleFiles(path, nil, nil)
		require.ErrorContains(t, err, "inspect output directory")
		assert.ErrorIs(t, err, os.ErrNotExist)
		_, err = os.Stat(path)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestCommandMetadataFailurePreservesBundle(t *testing.T) {
	t.Setenv("GOWORK", "off")
	for _, tt := range []struct {
		name        string
		declaration string
		outsideRoot bool
		message     string
	}{
		{name: "missing version", message: "no module version constant"},
		{name: "invalid version", declaration: `const version = "1.2.3"`, message: `should be "v1.2.3"`},
		{name: "invalid repository root", declaration: `const version = "v1.2.3"`, outsideRoot: true, message: "outside repository root"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.CopyFS(dir, fstest.MapFS{
				"go.mod":           {Data: []byte("module example.com/sample\n\ngo 1.26.0\n")},
				"sample.go":        {Data: []byte("package sample\n\ntype Value string\n")},
				"version.go":       {Data: []byte("package sample\n" + tt.declaration + "\n")},
				"api.md":           {Data: []byte("existing report\n")},
				"api.metadata.yml": {Data: []byte("existing metadata\n")},
			}))
			t.Chdir(dir)
			root := dir
			if tt.outsideRoot {
				root = t.TempDir()
			}
			command := newRootCommand()
			command.SetArgs([]string{dir, "--repo-root", root})
			var output bytes.Buffer
			command.SetOut(&output)
			require.ErrorContains(t, command.ExecuteContext(t.Context()), tt.message)
			assert.Empty(t, output.String())
			for name, expected := range map[string]string{
				"api.md": "existing report\n", "api.metadata.yml": "existing metadata\n",
			} {
				content, err := os.ReadFile(filepath.Join(dir, name))
				require.NoError(t, err)
				assert.Equal(t, expected, string(content))
			}
		})
	}
}

func TestHashCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api.md")
	require.NoError(t, os.WriteFile(path, []byte("abc"), 0o644))
	command := newRootCommand()
	command.SetArgs([]string{"hash", path})
	var output bytes.Buffer
	command.SetOut(&output)
	require.NoError(t, command.ExecuteContext(t.Context()))
	assert.Equal(t, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad\n", output.String())

	command = newRootCommand()
	command.SetArgs([]string{"hash", filepath.Join(t.TempDir(), "missing.md")})
	command.SetOut(&output)
	require.ErrorContains(t, command.ExecuteContext(t.Context()), "read API report")
}

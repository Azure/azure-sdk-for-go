// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type serverTests struct {
	suite.Suite
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTests))
}

func (s *serverTests) SetupSuite() {
	// Ignore manual start in pipeline tests, we always want to exercise install
	require.NoError(s.T(), os.Setenv(proxyManualStartEnv, "false"))
}

func (s *serverTests) TestProxyDownloadFile() {
	file, err := getTestProxyDownloadFile()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), file)
}

func (s *serverTests) TestEnsureTestProxyInstalled() {
	cwd, err := os.Getwd()
	require.NoError(s.T(), err)
	gitRoot, err := getGitRoot(cwd)
	require.NoError(s.T(), err)

	proxyDir := filepath.Join(os.TempDir(), ".proxy")
	proxyVersion, err := getProxyVersion(gitRoot)
	require.NoError(s.T(), err)

	err = os.RemoveAll(proxyDir)
	require.NoError(s.T(), err)
	err = os.MkdirAll(proxyDir, 0755)
	require.NoError(s.T(), err)

	proxyPath := filepath.Join(proxyDir, "Azure.Sdk.Tools.TestProxy")
	if runtime.GOOS == "windows" {
		proxyPath += ".exe"
	}

	// Test download proxy
	err = ensureTestProxyInstalled(proxyVersion, proxyPath, proxyDir, "")
	require.NoError(s.T(), err)

	stat1, err := os.Stat(proxyPath)
	require.NoError(s.T(), err)

	// Test cached proxy
	err = ensureTestProxyInstalled(proxyVersion, proxyPath, proxyDir, "")
	require.NoError(s.T(), err)

	stat2, err := os.Stat(proxyPath)
	require.NoError(s.T(), err)

	require.Equal(s.T(), stat1.ModTime(), stat2.ModTime(), "Expected proxy download to be cached")
}

func (s *serverTests) TestExtractInsecurePath() {
	s.T().Run("tar", func(t *testing.T) {
		td := t.TempDir()
		p := filepath.Join(td, "test.tar.gz")
		f, err := os.Create(p)
		require.NoError(t, err)
		zw := gzip.NewWriter(f)
		tw := tar.NewWriter(zw)
		b := []byte("_")
		err = tw.WriteHeader(&tar.Header{
			Name: filepath.Join("..", "file"),
			Size: int64(len(b)),
		})
		require.NoError(t, err)
		_, err = tw.Write(b)
		require.NoError(t, err)
		require.NoError(t, tw.Close())
		require.NoError(t, zw.Close())
		require.NoError(t, f.Close())

		err = installTestProxy(p, td, td)
		require.ErrorContains(t, err, "illegal file path")
	})
	s.T().Run("zip", func(t *testing.T) {
		td := t.TempDir()
		p := filepath.Join(td, "test.zip")
		f, err := os.Create(p)
		require.NoError(t, err)
		defer func() {
			require.NoError(s.T(), f.Close())
		}()
		zw := zip.NewWriter(f)
		w, err := zw.Create("../file")
		require.NoError(t, err)
		_, err = w.Write([]byte("_"))
		require.NoError(t, err)
		require.NoError(t, zw.Close())

		err = installTestProxy(p, td, td)
		require.ErrorContains(t, err, "illegal file path")
	})
}

// TestExtractSiblingPrefixPath covers entries that do not traverse upwards out
// of the destination's parent, but do land in a sibling directory whose name
// begins with the destination's name. "out-evil" shares the textual prefix
// "out", so a containment check written as a plain string prefix accepts it.
func (s *serverTests) TestExtractSiblingPrefixPath() {
	// filepath.Join cleans the "..", so the resolved path is <td>/out-evil/file
	// rather than anything containing "..".
	const entryName = "../out-evil/file"

	s.T().Run("tar", func(t *testing.T) {
		td := t.TempDir()
		dest := filepath.Join(td, "out")
		require.NoError(t, os.MkdirAll(dest, 0755))

		p := filepath.Join(td, "test.tar.gz")
		f, err := os.Create(p)
		require.NoError(t, err)
		zw := gzip.NewWriter(f)
		tw := tar.NewWriter(zw)
		b := []byte("_")
		err = tw.WriteHeader(&tar.Header{
			Name: entryName,
			Size: int64(len(b)),
		})
		require.NoError(t, err)
		_, err = tw.Write(b)
		require.NoError(t, err)
		require.NoError(t, tw.Close())
		require.NoError(t, zw.Close())
		require.NoError(t, f.Close())

		err = extractTestProxyArchive(p, dest)
		require.ErrorContains(t, err, "illegal file path")
		require.NoFileExists(t, filepath.Join(td, "out-evil", "file"))
	})

	s.T().Run("zip", func(t *testing.T) {
		td := t.TempDir()
		dest := filepath.Join(td, "out")
		require.NoError(t, os.MkdirAll(dest, 0755))

		p := filepath.Join(td, "test.zip")
		f, err := os.Create(p)
		require.NoError(t, err)
		zw := zip.NewWriter(f)
		w, err := zw.Create(entryName)
		require.NoError(t, err)
		_, err = w.Write([]byte("_"))
		require.NoError(t, err)
		require.NoError(t, zw.Close())
		require.NoError(t, f.Close())

		err = extractTestProxyZip(p, dest)
		require.ErrorContains(t, err, "illegal file path")
		require.NoFileExists(t, filepath.Join(td, "out-evil", "file"))
	})
}

func TestResolveExtractPath(t *testing.T) {
	dir := filepath.Join("tmp", "out")

	for _, tt := range []struct {
		name    string
		entry   string
		wantErr bool
	}{
		{name: "plain file", entry: "file"},
		{name: "nested file", entry: filepath.Join("a", "b", "file")},
		{name: "archive root", entry: "."},
		{name: "interior dot dot resolving inside", entry: filepath.Join("a", "..", "file")},
		{name: "parent", entry: "..", wantErr: true},
		{name: "traversal", entry: filepath.Join("..", "file"), wantErr: true},
		{name: "deep traversal", entry: filepath.Join("..", "..", "file"), wantErr: true},
		{name: "sibling sharing a prefix", entry: filepath.Join("..", "out-evil", "file"), wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveExtractPath(dir, tt.entry)
			if tt.wantErr {
				require.ErrorContains(t, err, "illegal file path")
				require.Empty(t, got)
				return
			}
			require.NoError(t, err)
			// Everything permitted must resolve to dir itself or below it.
			rel, err := filepath.Rel(dir, got)
			require.NoError(t, err)
			require.False(t, rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)))
		})
	}
}

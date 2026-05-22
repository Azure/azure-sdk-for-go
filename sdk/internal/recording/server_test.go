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

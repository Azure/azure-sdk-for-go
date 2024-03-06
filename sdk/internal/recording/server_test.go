//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"archive/zip"
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
	os.Setenv(proxyManualStartEnv, "false")
}

func (s *serverTests) TestProxyDownloadFile() {
	file, err := getTestProxyDownloadFile()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), file)
}

func (s *serverTests) TestExtractTestProxyZip() {
	zipFile, err := os.CreateTemp("", "test-extract-*.zip")
	require.NoError(s.T(), err)
	defer zipFile.Close()

	// Create a new zip archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
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

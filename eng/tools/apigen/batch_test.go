// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAllAPIReports(t *testing.T) {
	powershell, err := exec.LookPath("pwsh")
	if err != nil {
		t.Skip("PowerShell 7 is required to test the batch script")
	}
	t.Setenv("GOWORK", "off")
	binary := filepath.Join(t.TempDir(), "apigen")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.CommandContext(t.Context(), "go", "build", "-o", binary, ".")
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))
	script, err := filepath.Abs("Generate-AllApiReports.ps1")
	require.NoError(t, err)

	t.Run("continues after failures and reports every failed module", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "SDK modules with spaces")
		require.NoError(t, os.Mkdir(root, 0o755))
		require.NoError(t, os.Mkdir(filepath.Join(root, ".git"), 0o755))
		require.NoError(t, os.CopyFS(root, fstest.MapFS{
			"a-good/go.mod":     {Data: []byte("module example.com/first\n\ngo 1.26.0\n")},
			"a-good/main.go":    {Data: []byte("package first\n\ntype Value string\n")},
			"a-good/version.go": {Data: []byte("package first\n\nconst moduleVersion = \"v1.2.3\"\n")},
			"b-bad/go.mod":      {Data: []byte("module example.com/bad\n\ngo 1.26.0\n")},
			"b-bad/main.go":     {Data: []byte("package bad\n\nvar Value = undefinedValue\n")},
			"b-bad/api.md":      {Data: []byte("existing report\n")},
			"c-good/go.mod":     {Data: []byte("module example.com/last\n\ngo 1.26.0\n")},
			"c-good/main.go":    {Data: []byte("package last\n\ntype Value int\n")},
			"c-good/version.go": {Data: []byte("package last\n\nconst version = \"v2.0.0-beta.1\"\n")},
			"d-bad/go.mod":      {Data: []byte("module example.com/invalid\n\ngo invalid\n")},
			"d-bad/main.go":     {Data: []byte("package invalid\n")},
		}))
		command := exec.CommandContext(t.Context(), powershell, "-NoLogo", "-NoProfile", "-File", script,
			"-SdkRoot", root, "-ApigenPath", binary, "-Concurrency", "2")
		output, err := command.CombinedOutput()
		var failure *exec.ExitError
		require.True(t, errors.As(err, &failure), "%v\n%s", err, output)
		assert.Equal(t, 1, failure.ExitCode())
		assert.Contains(t, string(output), "Failed to generate API review bundles for 2 of 4 modules:")
		assert.Contains(t, string(output), filepath.Join(root, "b-bad"))
		assert.Contains(t, string(output), filepath.Join(root, "d-bad"))
		assert.Contains(t, strings.ReplaceAll(string(output), "\r\n", "\n"),
			"Failed to generate API review bundles for 2 of 4 modules:\n"+
				filepath.Join(root, "b-bad")+"\n"+filepath.Join(root, "d-bad")+"\n")
		for _, name := range []string{"a-good", "c-good"} {
			content, err := os.ReadFile(filepath.Join(root, name, "api.md"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "## Module")
			metadata, err := os.ReadFile(filepath.Join(root, name, "api.metadata.yml"))
			require.NoError(t, err)
			assert.Contains(t, string(metadata), "apiHash: ")
		}
		content, err := os.ReadFile(filepath.Join(root, "b-bad", "api.md"))
		require.NoError(t, err)
		assert.Equal(t, "existing report\n", string(content))
		_, err = os.Stat(filepath.Join(root, "d-bad", "api.md"))
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("reports success and includes nested modules", func(t *testing.T) {
		root := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(root, ".git"), 0o755))
		require.NoError(t, os.CopyFS(root, fstest.MapFS{
			"go.mod":            {Data: []byte("module example.com/outer\n\ngo 1.26.0\n")},
			"outer.go":          {Data: []byte("package outer\n\ntype Value int\n")},
			"version.go":        {Data: []byte("package outer\n\nconst version = \"v1.2.3\"\n")},
			"nested/go.mod":     {Data: []byte("module example.com/inner\n\ngo 1.26.0\n")},
			"nested/inner.go":   {Data: []byte("package inner\n\ntype Value string\n")},
			"nested/version.go": {Data: []byte("package inner\n\nconst version = \"v2.3.4\"\n")},
		}))
		command := exec.CommandContext(t.Context(), powershell, "-NoLogo", "-NoProfile", "-File", script,
			"-SdkRoot", root, "-ApigenPath", binary)
		output, err := command.CombinedOutput()
		require.NoError(t, err, string(output))
		assert.Contains(t, string(output), "Successfully generated API review bundles for all 2 modules.")
		for _, name := range []string{"", "nested"} {
			content, err := os.ReadFile(filepath.Join(root, name, "api.md"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "## Module")
			metadata, err := os.ReadFile(filepath.Join(root, name, "api.metadata.yml"))
			require.NoError(t, err)
			expectedPath := "."
			expectedVersion := "1.2.3"
			if name != "" {
				expectedPath, expectedVersion = "nested", "2.3.4"
			}
			assert.Contains(t, string(metadata), "packageRelativePath: \""+expectedPath+"\"\n")
			assert.Contains(t, string(metadata), "packageVersion: "+expectedVersion+"\n")
		}
	})
	t.Run("empty root is not reported as success", func(t *testing.T) {
		command := exec.CommandContext(t.Context(), powershell, "-NoLogo", "-NoProfile", "-File", script,
			"-SdkRoot", t.TempDir(), "-ApigenPath", binary)
		output, err := command.CombinedOutput()
		require.Error(t, err)
		assert.Contains(t, string(output), "No Go modules found")
		assert.NotContains(t, string(output), "Successfully generated")
	})
	t.Run("rejects nonpositive concurrency", func(t *testing.T) {
		for _, concurrency := range []string{"0", "-1"} {
			command := exec.CommandContext(t.Context(), powershell, "-NoLogo", "-NoProfile", "-File", script,
				"-SdkRoot", t.TempDir(), "-ApigenPath", binary, "-Concurrency", concurrency)
			output, err := command.CombinedOutput()
			require.Error(t, err, string(output))
			assert.Contains(t, string(output), "Concurrency")
			assert.NotContains(t, string(output), "Generating API reports")
		}
	})
}

func TestGenerateAllAPIReportsConcurrency(t *testing.T) {
	powershell, err := exec.LookPath("pwsh")
	if err != nil {
		t.Skip("PowerShell 7 is required to test the batch script")
	}
	script, err := filepath.Abs("Generate-AllApiReports.ps1")
	require.NoError(t, err)
	for _, concurrency := range []int{1, 2} {
		t.Run(fmt.Sprintf("limit %d", concurrency), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
			defer cancel()
			started := make(chan struct{}, 4)
			release := make(chan struct{})
			var mutex sync.Mutex
			active, maximum := 0, 0
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				mutex.Lock()
				active++
				maximum = max(maximum, active)
				mutex.Unlock()
				started <- struct{}{}
				select {
				case <-release:
				case <-request.Context().Done():
				}
				mutex.Lock()
				active--
				mutex.Unlock()
				writer.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()
			defer close(release)
			t.Setenv("APIGEN_CONCURRENCY_PROBE_URL", server.URL)

			root := t.TempDir()
			require.NoError(t, os.CopyFS(root, fstest.MapFS{
				"a/go.mod": {Data: []byte("module example.com/a\n")},
				"b/go.mod": {Data: []byte("module example.com/b\n")},
				"c/go.mod": {Data: []byte("module example.com/c\n")},
				"d/go.mod": {Data: []byte("module example.com/d\n")},
				"probe.ps1": {Data: []byte(`param([string] $ModuleDirectory, [string] $o)
$ErrorActionPreference = 'Stop'
Invoke-WebRequest -Uri $env:APIGEN_CONCURRENCY_PROBE_URL -TimeoutSec 30 | Out-Null
Set-Content -LiteralPath (Join-Path $o 'api.md') -Value 'report'
Set-Content -LiteralPath (Join-Path $o 'api.metadata.yml') -Value 'metadata'
exit 0
`)},
			}))
			command := exec.CommandContext(ctx, powershell, "-NoLogo", "-NoProfile", "-File", script,
				"-SdkRoot", root, "-ApigenPath", filepath.Join(root, "probe.ps1"),
				"-Concurrency", fmt.Sprint(concurrency))
			finished := make(chan struct{})
			var output []byte
			var commandError error
			go func() {
				output, commandError = command.CombinedOutput()
				close(finished)
			}()
			defer func() {
				cancel()
				<-finished
			}()

			// Hold each wave until every permitted task has actually started.
			for completed := 0; completed < 4; completed += concurrency {
				for range concurrency {
					select {
					case <-started:
					case <-finished:
						require.FailNow(t, "script exited before starting the expected tasks", "%v\n%s", commandError, output)
					case <-ctx.Done():
						require.FailNow(t, "timed out waiting for concurrent tasks", "%v", ctx.Err())
					}
				}
				for range concurrency {
					select {
					case release <- struct{}{}:
					case <-ctx.Done():
						require.FailNow(t, "timed out releasing concurrent tasks", "%v", ctx.Err())
					}
				}
			}
			<-finished
			require.NoError(t, commandError, string(output))
			assert.Contains(t, string(output), "Successfully generated API review bundles for all 4 modules.")
			mutex.Lock()
			observedMaximum := maximum
			mutex.Unlock()
			assert.Equal(t, concurrency, observedMaximum, "tasks must overlap without exceeding the configured limit")
		})
	}
}

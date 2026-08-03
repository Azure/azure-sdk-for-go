// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveRunInvocationFromConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "workloads.json")

	err := os.WriteFile(configPath, []byte(`{
		"defaultWorkload": "wl-upload",
		"workloads": {
			"wl-upload": {
				"test": "UploadBlobTest",
				"parameters": {
					"duration": 5,
					"warmup": 2,
					"debug": true
				}
			}
		}
	}`), 0o600)
	require.NoError(t, err)

	invocation, err := resolveRunInvocation([]string{"--config", configPath})
	require.NoError(t, err)
	require.Equal(t, "UploadBlobTest", invocation.TestName)
	require.Equal(t, "wl-upload", invocation.Workload)
	require.True(t, invocation.UsesConfig)
	require.Contains(t, invocation.ConfigArgs, "--duration=5")
	require.Contains(t, invocation.ConfigArgs, "--warmup=2")
	require.Contains(t, invocation.ConfigArgs, "--debug=true")
}

func TestResolveRunInvocationCLITakesPrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "workloads.json")

	err := os.WriteFile(configPath, []byte(`{
		"defaultWorkload": "wl-upload",
		"workloads": {
			"wl-upload": {
				"test": "UploadBlobTest",
				"parameters": {
					"duration": 5
				}
			}
		}
	}`), 0o600)
	require.NoError(t, err)

	invocation, err := resolveRunInvocation([]string{"--config", configPath, "--duration", "9"})
	require.NoError(t, err)
	require.Equal(t, []string{"--duration", "9"}, invocation.CLIArgs)
	invocationArgs := append([]string{}, invocation.ConfigArgs...)
	invocationArgs = append(invocationArgs, invocation.CLIArgs...)
	require.Contains(t, invocationArgs, "--duration=5")
	require.Contains(t, invocationArgs, "--duration")
	require.Contains(t, invocationArgs, "9")
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewWorkTreeWithWorktreeConfigExtension guards against dependency upgrades that break
// opening repositories initialized with `git sparse-checkout init`, which is how the SDK
// automation pipelines check out azure-sdk-for-go. That command enables the `worktreeConfig`
// extension, and go-git v5.17.0 through at least v5.19.1 reject such repositories.
func TestNewWorkTreeWithWorktreeConfigExtension(t *testing.T) {
	dir := t.TempDir()

	for _, args := range [][]string{
		{"init"},
		{"sparse-checkout", "init"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "git %v failed: %s", args, string(output))
	}

	wt, err := NewWorkTree(dir)
	require.NoError(t, err)
	require.NotNil(t, wt)
}

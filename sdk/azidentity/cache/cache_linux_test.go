// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCacheDir(t *testing.T) {
	t.Run("no XDG_CACHE_HOME", func(t *testing.T) {
		home, err := os.UserHomeDir()
		require.NoError(t, err)
		expected := filepath.Join(home, ".cache")
		t.Setenv("XDG_CACHE_HOME", "")
		actual, err := cacheDir()
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
	t.Run("XDG_CACHE_HOME", func(t *testing.T) {
		expected := string(filepath.Separator) + "foo"
		t.Setenv("XDG_CACHE_HOME", expected)
		actual, err := cacheDir()
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestKeyExistsButNotFile(t *testing.T) {
	expected := []byte(t.Name())
	a, err := storage(t.Name())
	require.NoError(t, err)
	err = a.Write(ctx, append([]byte("not"), expected...))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, a.Delete(ctx)) })

	p, err := cacheFilePath(t.Name())
	require.NoError(t, err)
	require.NoError(t, os.Remove(p))

	b, err := newKeyring(t.Name())
	require.NoError(t, err)
	data, err := b.Read(ctx)
	require.NoError(t, err)
	require.Empty(t, data)

	err = b.Write(ctx, expected)
	require.NoError(t, err)
	actual, err := b.Read(ctx)
	require.NoError(t, err)
	require.EqualValues(t, expected, actual)
}

func TestReadWriteDelete(t *testing.T) {
	for _, test := range []struct {
		expected   []byte
		desc, name string
	}{
		{desc: "write empty slice"},
		{desc: "write then read", expected: []byte("expected")},
		{desc: "file exists but not key", expected: []byte("expected"), name: t.Name()},
	} {
		t.Run(test.desc, func(t *testing.T) {
			dir, err := cacheDir()
			require.NoError(t, err)
			name := test.name
			if name == "" {
				// a UUID name ensures the file and key don't exist
				name = uuid.NewString()
			} else {
				// Write the file to simulate a cache encrypted by a lost key. In this
				// case Read should return nil and Write should overwrite the file.
				p := filepath.Join(dir, ".IdentityService", name)
				err = os.MkdirAll(filepath.Dir(p), 0600)
				require.NoError(t, err)
				err = os.WriteFile(p, []byte("eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiNjI3OTEzODA4In0..gPRNjqd4HcrlFxJdEEaFeA.Pqpr_IYG7e1lt6KPoE0v_A.i9h5iJWw9bT217I5M2Ufrg"), 0600)
				require.NoError(t, err)
			}
			k, err := newKeyring(name)
			require.NoError(t, err)

			actual, err := k.Read(ctx)
			require.NoError(t, err)
			require.Empty(t, actual)

			cp := make([]byte, len(test.expected))
			copy(cp, test.expected)
			err = k.Write(ctx, cp)
			require.NoError(t, err)
			if len(test.expected) > 0 {
				b, err := os.ReadFile(k.file)
				require.NoError(t, err)
				require.NotEqual(t, test.expected, b, "file content isn't encrypted")
			}

			actual, err = k.Read(ctx)
			require.NoError(t, err)
			require.EqualValues(t, test.expected, actual)

			require.NoError(t, k.Delete(ctx))
			require.NoFileExists(t, k.file)

			_, err = k.Read(ctx)
			require.NoError(t, err)
		})
	}
}

func TestTwoInstances(t *testing.T) {
	for _, deleteFile := range []bool{false, true} {
		s := "key and file exist"
		if deleteFile {
			s = "key exists but not file"
		}
		t.Run(s, func(t *testing.T) {
			name := uuid.NewString()
			a, err := newKeyring(name)
			require.NoError(t, err)
			expected := []byte(t.Name())
			err = a.Write(ctx, expected)
			require.NoError(t, err)

			if deleteFile {
				require.NoError(t, os.Remove(a.file))
			}

			b, err := newKeyring(name)
			require.NoError(t, err)
			var actual []byte
			if deleteFile {
				// a should be able to read the file written by b
				err = b.Write(ctx, expected)
				require.NoError(t, err)
				actual, err = a.Read(ctx)
			} else {
				// b should be able to read the file with the key created by a
				actual, err = b.Read(ctx)
			}
			require.NoError(t, err)
			require.EqualValues(t, expected, actual)

			require.NoError(t, a.Delete(ctx))
			// neither the file nor key should exist, however b shouldn't return an error
			require.NoError(t, b.Delete(ctx))
		})
	}
}

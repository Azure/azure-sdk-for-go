// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeDocumentExactSize(t *testing.T) {
	body, err := makeDocument("campaign-owner", 17)
	require.NoError(t, err)
	require.Len(t, body, documentBytes)
	require.NoError(t, validateDocument(body, "campaign-owner", 17))
}

func TestKeySequenceRoundRobin(t *testing.T) {
	require.Equal(t, 0, keyIndex(0))
	require.Equal(t, 1023, keyIndex(1023))
	require.Equal(t, 0, keyIndex(1024))
	require.Equal(t, 1, keyIndex(1025))
}

func TestOwnerGuard(t *testing.T) {
	body, err := makeDocument("campaign-owner", 2)
	require.NoError(t, err)
	owned, err := ownedBy(body, "other-owner")
	require.NoError(t, err)
	require.False(t, owned)
}

func TestRejectSecretArguments(t *testing.T) {
	require.Error(t, rejectSecretArguments([]string{"run", "--key=secret"}))
	t.Setenv("AZURE_COSMOS_KEY", "secret")
	require.Error(t, rejectSecretArguments([]string{"run"}))
}

func TestReadSecretFD(t *testing.T) {
	readFile, writeFile, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = readFile.Close()
		_ = writeFile.Close()
	})
	_, err = writeFile.WriteString("secret-value\n")
	require.NoError(t, err)
	require.NoError(t, writeFile.Close())

	value, err := readSecretFD("3")
	if readFile.Fd() != 3 {
		value, err = readSecretFD(fdString(readFile.Fd()))
	}
	require.NoError(t, err)
	require.Equal(t, []byte("secret-value"), value)
	clearBytes(value)
	require.Equal(t, make([]byte, len(value)), value)
}

func fdString(fd uintptr) string {
	const digits = "0123456789"
	if fd == 0 {
		return "0"
	}
	var reversed [32]byte
	index := len(reversed)
	for fd > 0 {
		index--
		reversed[index] = digits[fd%10]
		fd /= 10
	}
	return string(reversed[index:])
}

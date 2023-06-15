//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilesystemPermissions_String(t *testing.T) {
	testdata := []struct {
		input    FilesystemPermissions
		expected string
	}{
		{input: FilesystemPermissions{Read: true}, expected: "r"},
		{input: FilesystemPermissions{Add: true}, expected: "a"},
		{input: FilesystemPermissions{Create: true}, expected: "c"},
		{input: FilesystemPermissions{Write: true}, expected: "w"},
		{input: FilesystemPermissions{Delete: true}, expected: "d"},
		{input: FilesystemPermissions{List: true}, expected: "l"},
		{input: FilesystemPermissions{Move: true}, expected: "m"},
		{input: FilesystemPermissions{Execute: true}, expected: "e"},
		{input: FilesystemPermissions{ModifyOwnership: true}, expected: "o"},
		{input: FilesystemPermissions{ModifyPermissions: true}, expected: "p"},
		{input: FilesystemPermissions{
			Read:              true,
			Add:               true,
			Create:            true,
			Write:             true,
			Delete:            true,
			List:              true,
			Move:              true,
			Execute:           true,
			ModifyOwnership:   true,
			ModifyPermissions: true,
		}, expected: "racwdlmeop"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestFilesystemPermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected FilesystemPermissions
	}{
		{expected: FilesystemPermissions{Read: true}, input: "r"},
		{expected: FilesystemPermissions{Add: true}, input: "a"},
		{expected: FilesystemPermissions{Create: true}, input: "c"},
		{expected: FilesystemPermissions{Write: true}, input: "w"},
		{expected: FilesystemPermissions{Delete: true}, input: "d"},
		{expected: FilesystemPermissions{List: true}, input: "l"},
		{expected: FilesystemPermissions{Move: true}, input: "m"},
		{expected: FilesystemPermissions{Execute: true}, input: "e"},
		{expected: FilesystemPermissions{ModifyOwnership: true}, input: "o"},
		{expected: FilesystemPermissions{ModifyPermissions: true}, input: "p"},
		{expected: FilesystemPermissions{
			Read:              true,
			Add:               true,
			Create:            true,
			Write:             true,
			Delete:            true,
			List:              true,
			Move:              true,
			Execute:           true,
			ModifyOwnership:   true,
			ModifyPermissions: true,
		}, input: "racwdlmeop"},
		{expected: FilesystemPermissions{
			Read:              true,
			Add:               true,
			Create:            true,
			Write:             true,
			Delete:            true,
			List:              true,
			Move:              true,
			Execute:           true,
			ModifyOwnership:   true,
			ModifyPermissions: true,
		}, input: "cpwmreodal"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseFilesystemPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestFilesystemPermissions_ParseNegative(t *testing.T) {
	_, err := parseFilesystemPermissions("cpwmreodalz") // Here 'z' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "122")
}

func TestFilePermissions_String(t *testing.T) {
	testdata := []struct {
		input    FilePermissions
		expected string
	}{
		{input: FilePermissions{Read: true}, expected: "r"},
		{input: FilePermissions{Add: true}, expected: "a"},
		{input: FilePermissions{Create: true}, expected: "c"},
		{input: FilePermissions{Write: true}, expected: "w"},
		{input: FilePermissions{Delete: true}, expected: "d"},
		{input: FilePermissions{List: true}, expected: "l"},
		{input: FilePermissions{Move: true}, expected: "m"},
		{input: FilePermissions{Execute: true}, expected: "e"},
		{input: FilePermissions{Ownership: true}, expected: "o"},
		{input: FilePermissions{Permissions: true}, expected: "p"},
		{input: FilePermissions{
			Read:        true,
			Add:         true,
			Create:      true,
			Write:       true,
			Delete:      true,
			List:        true,
			Move:        true,
			Execute:     true,
			Ownership:   true,
			Permissions: true,
		}, expected: "racwdlmeop"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestFilePermissions_Parse(t *testing.T) {
	testdata := []struct {
		expected FilePermissions
		input    string
	}{
		{expected: FilePermissions{Read: true}, input: "r"},
		{expected: FilePermissions{Add: true}, input: "a"},
		{expected: FilePermissions{Create: true}, input: "c"},
		{expected: FilePermissions{Write: true}, input: "w"},
		{expected: FilePermissions{Delete: true}, input: "d"},
		{expected: FilePermissions{List: true}, input: "l"},
		{expected: FilePermissions{Move: true}, input: "m"},
		{expected: FilePermissions{Execute: true}, input: "e"},
		{expected: FilePermissions{Ownership: true}, input: "o"},
		{expected: FilePermissions{Permissions: true}, input: "p"},
		{expected: FilePermissions{
			Read:        true,
			Add:         true,
			Create:      true,
			Write:       true,
			Delete:      true,
			List:        true,
			Move:        true,
			Execute:     true,
			Ownership:   true,
			Permissions: true,
		}, input: "racwdlmeop"},
		{expected: FilePermissions{
			Read:        true,
			Add:         true,
			Create:      true,
			Write:       true,
			Delete:      true,
			List:        true,
			Move:        true,
			Execute:     true,
			Ownership:   true,
			Permissions: true,
		}, input: "apwecrdlmo"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parsePathPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestParsePermissionsNegative(t *testing.T) {
	_, err := parsePathPermissions("awecrdlfmo") // Here 'f' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "102")
}

func TestGetCanonicalName(t *testing.T) {
	testdata := []struct {
		inputAccount   string
		inputContainer string
		inputBlob      string
		inputDirectory string
		expected       string
	}{
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", expected: "/blob/fakestorageaccount/fakestoragecontainer"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestorageblob", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestorageblob"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestoragevirtualdir/fakestorageblob", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestoragevirtualdir/fakestorageblob"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestoragevirtualdir\\fakestorageblob", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestoragevirtualdir/fakestorageblob"},
		{inputAccount: "fakestorageaccount", inputContainer: "fakestoragecontainer", inputBlob: "fakestoragedirectory", expected: "/blob/fakestorageaccount/fakestoragecontainer/fakestoragedirectory"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getCanonicalName(c.inputAccount, c.inputContainer, c.inputBlob, c.inputDirectory))
	}
}

func TestGetDirectoryDepth(t *testing.T) {
	testdata := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "myfile", expected: "1"},
		{input: "mydirectory", expected: "1"},
		{input: "mydirectory/myfile", expected: "2"},
		{input: "mydirectory/mysubdirectory", expected: "2"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getDirectoryDepth(c.input))
	}
}

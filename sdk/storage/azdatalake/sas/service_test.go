//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFileSystemPermissions_String(t *testing.T) {
	testdata := []struct {
		input    FileSystemPermissions
		expected string
	}{
		{input: FileSystemPermissions{Read: true}, expected: "r"},
		{input: FileSystemPermissions{Add: true}, expected: "a"},
		{input: FileSystemPermissions{Create: true}, expected: "c"},
		{input: FileSystemPermissions{Write: true}, expected: "w"},
		{input: FileSystemPermissions{Delete: true}, expected: "d"},
		{input: FileSystemPermissions{List: true}, expected: "l"},
		{input: FileSystemPermissions{Move: true}, expected: "m"},
		{input: FileSystemPermissions{Execute: true}, expected: "e"},
		{input: FileSystemPermissions{ModifyOwnership: true}, expected: "o"},
		{input: FileSystemPermissions{ModifyPermissions: true}, expected: "p"},
		{input: FileSystemPermissions{
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

func TestFileSystemPermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected FileSystemPermissions
	}{
		{expected: FileSystemPermissions{Read: true}, input: "r"},
		{expected: FileSystemPermissions{Add: true}, input: "a"},
		{expected: FileSystemPermissions{Create: true}, input: "c"},
		{expected: FileSystemPermissions{Write: true}, input: "w"},
		{expected: FileSystemPermissions{Delete: true}, input: "d"},
		{expected: FileSystemPermissions{List: true}, input: "l"},
		{expected: FileSystemPermissions{Move: true}, input: "m"},
		{expected: FileSystemPermissions{Execute: true}, input: "e"},
		{expected: FileSystemPermissions{ModifyOwnership: true}, input: "o"},
		{expected: FileSystemPermissions{ModifyPermissions: true}, input: "p"},
		{expected: FileSystemPermissions{
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
		{expected: FileSystemPermissions{
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
		permissions, err := parseFileSystemPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestFileSystemPermissions_ParseNegative(t *testing.T) {
	_, err := parseFileSystemPermissions("cpwmreodalz") // Here 'z' is invalid
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

func TestDirectoryPermissions_String(t *testing.T) {
	testdata := []struct {
		input    DirectoryPermissions
		expected string
	}{
		{input: DirectoryPermissions{Read: true}, expected: "r"},
		{input: DirectoryPermissions{Add: true}, expected: "a"},
		{input: DirectoryPermissions{Create: true}, expected: "c"},
		{input: DirectoryPermissions{Write: true}, expected: "w"},
		{input: DirectoryPermissions{Delete: true}, expected: "d"},
		{input: DirectoryPermissions{List: true}, expected: "l"},
		{input: DirectoryPermissions{Move: true}, expected: "m"},
		{input: DirectoryPermissions{Execute: true}, expected: "e"},
		{input: DirectoryPermissions{Ownership: true}, expected: "o"},
		{input: DirectoryPermissions{Permissions: true}, expected: "p"},
		{input: DirectoryPermissions{
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

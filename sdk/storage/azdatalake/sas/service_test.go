// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		{input: FileSystemPermissions{Tag: true}, expected: "t"},
		{input: FileSystemPermissions{
			Read:              true,
			Add:               true,
			Create:            true,
			Write:             true,
			Delete:            true,
			List:              true,
			Tag:               true,
			Move:              true,
			Execute:           true,
			ModifyOwnership:   true,
			ModifyPermissions: true,
		}, expected: "racwdltmeop"},
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
		{expected: FileSystemPermissions{Tag: true}, input: "t"},
		{expected: FileSystemPermissions{
			Read:              true,
			Add:               true,
			Create:            true,
			Write:             true,
			Delete:            true,
			List:              true,
			Tag:               true,
			Move:              true,
			Execute:           true,
			ModifyOwnership:   true,
			ModifyPermissions: true,
		}, input: "racwdltmeop"},
		{expected: FileSystemPermissions{
			Read:              true,
			Add:               true,
			Create:            true,
			Write:             true,
			Delete:            true,
			List:              true,
			Tag:               true,
			Move:              true,
			Execute:           true,
			ModifyOwnership:   true,
			ModifyPermissions: true,
		}, input: "cpwmtreodal"}, // Wrong order parses correctly
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
		{input: FilePermissions{Tag: true}, expected: "t"},
		{input: FilePermissions{
			Read:        true,
			Add:         true,
			Create:      true,
			Write:       true,
			Delete:      true,
			List:        true,
			Tag:         true,
			Move:        true,
			Execute:     true,
			Ownership:   true,
			Permissions: true,
		}, expected: "racwdltmeop"},
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
		{expected: FilePermissions{Tag: true}, input: "t"},
		{expected: FilePermissions{
			Read:        true,
			Add:         true,
			Create:      true,
			Write:       true,
			Delete:      true,
			List:        true,
			Tag:         true,
			Move:        true,
			Execute:     true,
			Ownership:   true,
			Permissions: true,
		}, input: "racwdltmeop"},
		{expected: FilePermissions{
			Read:        true,
			Add:         true,
			Create:      true,
			Write:       true,
			Delete:      true,
			List:        true,
			Tag:         true,
			Move:        true,
			Execute:     true,
			Ownership:   true,
			Permissions: true,
		}, input: "apwecrdtlmo"}, // Wrong order parses correctly
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
		{input: DirectoryPermissions{Tag: true}, expected: "t"},
		{input: DirectoryPermissions{
			Read:        true,
			Add:         true,
			Create:      true,
			Write:       true,
			Delete:      true,
			List:        true,
			Tag:         true,
			Move:        true,
			Execute:     true,
			Ownership:   true,
			Permissions: true,
		}, expected: "racwdltmeop"},
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
		{input: "foo/bar/hello", expected: "3"},
		{input: "a/b/c/d", expected: "4"},
		{input: "/", expected: "0"},
		{input: "foo/", expected: "1"},
		{input: "/foo", expected: "1"},
		{input: "/foo/bar/", expected: "2"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getDirectoryDepth(c.input))
	}
}

func TestFormatSignedRequestHeaders(t *testing.T) {
	testdata := []struct {
		desc                  string
		input                 map[string]string
		expectedNames         string
		expectedCanonicalized string
	}{
		{
			desc:                  "nil map",
			input:                 nil,
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "empty map",
			input:                 map[string]string{},
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "single header",
			input:                 map[string]string{"abra": "cadabra"},
			expectedNames:         "abra",
			expectedCanonicalized: "abra:cadabra\n",
		},
	}
	for _, c := range testdata {
		t.Run(c.desc, func(t *testing.T) {
			names, canonicalized := formatSignedRequestHeaders(c.input)
			require.Equal(t, c.expectedNames, names)
			require.Equal(t, c.expectedCanonicalized, canonicalized)
		})
	}

	// Multi-key test: we can't predict map iteration order, so validate both parts independently
	t.Run("multiple headers", func(t *testing.T) {
		input := map[string]string{"foo": "123", "bar": "456"}
		names, canonicalized := formatSignedRequestHeaders(input)

		// names should contain both keys comma-separated
		require.Contains(t, names, "foo")
		require.Contains(t, names, "bar")
		require.Contains(t, names, ",")

		// canonicalized should contain both key:value pairs each ending with \n
		require.Contains(t, canonicalized, "foo:123\n")
		require.Contains(t, canonicalized, "bar:456\n")
	})
}

func TestFormatSignedRequestQueryParameters(t *testing.T) {
	testdata := []struct {
		desc                  string
		input                 map[string]string
		expectedNames         string
		expectedCanonicalized string
	}{
		{
			desc:                  "nil map",
			input:                 nil,
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "empty map",
			input:                 map[string]string{},
			expectedNames:         "",
			expectedCanonicalized: "",
		},
		{
			desc:                  "single param",
			input:                 map[string]string{"foo": "123"},
			expectedNames:         "foo",
			expectedCanonicalized: "\nfoo:123",
		},
	}
	for _, c := range testdata {
		t.Run(c.desc, func(t *testing.T) {
			names, canonicalized := formatSignedRequestQueryParameters(c.input)
			require.Equal(t, c.expectedNames, names)
			require.Equal(t, c.expectedCanonicalized, canonicalized)
		})
	}

	// Multi-key test: validate both parts independently due to map iteration order
	t.Run("multiple params", func(t *testing.T) {
		input := map[string]string{"hello": "world", "abra": "cadabra"}
		names, canonicalized := formatSignedRequestQueryParameters(input)

		// names should contain both keys comma-separated
		require.Contains(t, names, "hello")
		require.Contains(t, names, "abra")
		require.Contains(t, names, ",")

		// canonicalized should contain both key:value pairs each prefixed with \n
		require.Contains(t, canonicalized, "\nhello:world")
		require.Contains(t, canonicalized, "\nabra:cadabra")
	})
}

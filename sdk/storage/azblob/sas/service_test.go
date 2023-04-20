//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContainerPermissions_String(t *testing.T) {
	testdata := []struct {
		input    ContainerPermissions
		expected string
	}{
		{input: ContainerPermissions{Read: true}, expected: "r"},
		{input: ContainerPermissions{Add: true}, expected: "a"},
		{input: ContainerPermissions{Create: true}, expected: "c"},
		{input: ContainerPermissions{Write: true}, expected: "w"},
		{input: ContainerPermissions{Delete: true}, expected: "d"},
		{input: ContainerPermissions{DeletePreviousVersion: true}, expected: "x"},
		{input: ContainerPermissions{List: true}, expected: "l"},
		{input: ContainerPermissions{Tag: true}, expected: "t"},
		{input: ContainerPermissions{FilterByTags: true}, expected: "f"},
		{input: ContainerPermissions{Move: true}, expected: "m"},
		{input: ContainerPermissions{Execute: true}, expected: "e"},
		{input: ContainerPermissions{ModifyOwnership: true}, expected: "o"},
		{input: ContainerPermissions{ModifyPermissions: true}, expected: "p"},
		{input: ContainerPermissions{SetImmutabilityPolicy: true}, expected: "i"},
		{input: ContainerPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			List:                  true,
			Tag:                   true,
			FilterByTags:          true,
			Move:                  true,
			Execute:               true,
			ModifyOwnership:       true,
			ModifyPermissions:     true,
			SetImmutabilityPolicy: true,
		}, expected: "racwdxltfmeopi"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestContainerPermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected ContainerPermissions
	}{
		{expected: ContainerPermissions{Read: true}, input: "r"},
		{expected: ContainerPermissions{Add: true}, input: "a"},
		{expected: ContainerPermissions{Create: true}, input: "c"},
		{expected: ContainerPermissions{Write: true}, input: "w"},
		{expected: ContainerPermissions{Delete: true}, input: "d"},
		{expected: ContainerPermissions{DeletePreviousVersion: true}, input: "x"},
		{expected: ContainerPermissions{List: true}, input: "l"},
		{expected: ContainerPermissions{Tag: true}, input: "t"},
		{expected: ContainerPermissions{FilterByTags: true}, input: "f"},
		{expected: ContainerPermissions{Move: true}, input: "m"},
		{expected: ContainerPermissions{Execute: true}, input: "e"},
		{expected: ContainerPermissions{ModifyOwnership: true}, input: "o"},
		{expected: ContainerPermissions{ModifyPermissions: true}, input: "p"},
		{expected: ContainerPermissions{SetImmutabilityPolicy: true}, input: "i"},
		{expected: ContainerPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			List:                  true,
			Tag:                   true,
			FilterByTags:          true,
			Move:                  true,
			Execute:               true,
			ModifyOwnership:       true,
			ModifyPermissions:     true,
			SetImmutabilityPolicy: true,
		}, input: "racwdxltfmeopi"},
		{expected: ContainerPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			List:                  true,
			Tag:                   true,
			FilterByTags:          true,
			Move:                  true,
			Execute:               true,
			ModifyOwnership:       true,
			ModifyPermissions:     true,
			SetImmutabilityPolicy: true,
		}, input: "ctpwxfmreodail"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseContainerPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestContainerPermissions_ParseNegative(t *testing.T) {
	_, err := parseContainerPermissions("cpwxtfmreodailz") // Here 'z' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "122")
}

func TestBlobPermissions_String(t *testing.T) {
	testdata := []struct {
		input    BlobPermissions
		expected string
	}{
		{input: BlobPermissions{Read: true}, expected: "r"},
		{input: BlobPermissions{Add: true}, expected: "a"},
		{input: BlobPermissions{Create: true}, expected: "c"},
		{input: BlobPermissions{Write: true}, expected: "w"},
		{input: BlobPermissions{Delete: true}, expected: "d"},
		{input: BlobPermissions{DeletePreviousVersion: true}, expected: "x"},
		{input: BlobPermissions{PermanentDelete: true}, expected: "y"},
		{input: BlobPermissions{List: true}, expected: "l"},
		{input: BlobPermissions{Tag: true}, expected: "t"},
		{input: BlobPermissions{Move: true}, expected: "m"},
		{input: BlobPermissions{Execute: true}, expected: "e"},
		{input: BlobPermissions{Ownership: true}, expected: "o"},
		{input: BlobPermissions{Permissions: true}, expected: "p"},
		{input: BlobPermissions{SetImmutabilityPolicy: true}, expected: "i"},
		{input: BlobPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Tag:                   true,
			Move:                  true,
			Execute:               true,
			Ownership:             true,
			Permissions:           true,
			SetImmutabilityPolicy: true,
		}, expected: "racwdxyltmeopi"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestBlobPermissions_Parse(t *testing.T) {
	testdata := []struct {
		expected BlobPermissions
		input    string
	}{
		{expected: BlobPermissions{Read: true}, input: "r"},
		{expected: BlobPermissions{Add: true}, input: "a"},
		{expected: BlobPermissions{Create: true}, input: "c"},
		{expected: BlobPermissions{Write: true}, input: "w"},
		{expected: BlobPermissions{Delete: true}, input: "d"},
		{expected: BlobPermissions{DeletePreviousVersion: true}, input: "x"},
		{expected: BlobPermissions{PermanentDelete: true}, input: "y"},
		{expected: BlobPermissions{List: true}, input: "l"},
		{expected: BlobPermissions{Tag: true}, input: "t"},
		{expected: BlobPermissions{Move: true}, input: "m"},
		{expected: BlobPermissions{Execute: true}, input: "e"},
		{expected: BlobPermissions{Ownership: true}, input: "o"},
		{expected: BlobPermissions{Permissions: true}, input: "p"},
		{expected: BlobPermissions{SetImmutabilityPolicy: true}, input: "i"},
		{expected: BlobPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Tag:                   true,
			Move:                  true,
			Execute:               true,
			Ownership:             true,
			Permissions:           true,
			SetImmutabilityPolicy: true,
		}, input: "racwdxyltmeopi"},
		{expected: BlobPermissions{
			Read:                  true,
			Add:                   true,
			Create:                true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Tag:                   true,
			Move:                  true,
			Execute:               true,
			Ownership:             true,
			Permissions:           true,
			SetImmutabilityPolicy: true,
		}, input: "apwecxrdlmyiot"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseBlobPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestBlobPermissions_ParseNegative(t *testing.T) {
	_, err := parseBlobPermissions("apwecxrdlfmyiot") // Here 'f' is invalid
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

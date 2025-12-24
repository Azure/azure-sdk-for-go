// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSharePermissions_String(t *testing.T) {
	testdata := []struct {
		input    SharePermissions
		expected string
	}{
		{input: SharePermissions{Read: true}, expected: "r"},
		{input: SharePermissions{Create: true}, expected: "c"},
		{input: SharePermissions{Write: true}, expected: "w"},
		{input: SharePermissions{Delete: true}, expected: "d"},
		{input: SharePermissions{List: true}, expected: "l"},
		{input: SharePermissions{
			Read:   true,
			Create: true,
			Write:  true,
			Delete: true,
			List:   true,
		}, expected: "rcwdl"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestSharePermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected SharePermissions
	}{
		{expected: SharePermissions{Read: true}, input: "r"},
		{expected: SharePermissions{Create: true}, input: "c"},
		{expected: SharePermissions{Write: true}, input: "w"},
		{expected: SharePermissions{Delete: true}, input: "d"},
		{expected: SharePermissions{List: true}, input: "l"},
		{expected: SharePermissions{
			Read:   true,
			Create: true,
			Write:  true,
			Delete: true,
			List:   true,
		}, input: "rcwdl"},
		{expected: SharePermissions{
			Read:   true,
			Create: true,
			Write:  true,
			Delete: true,
			List:   true,
		}, input: "cwrdl"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseSharePermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestSharePermissions_ParseNegative(t *testing.T) {
	_, err := parseSharePermissions("cwtrdl") // Here 't' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "116")
}

func TestFilePermissions_String(t *testing.T) {
	testdata := []struct {
		input    FilePermissions
		expected string
	}{
		{input: FilePermissions{Read: true}, expected: "r"},
		{input: FilePermissions{Create: true}, expected: "c"},
		{input: FilePermissions{Write: true}, expected: "w"},
		{input: FilePermissions{Delete: true}, expected: "d"},
		{input: FilePermissions{
			Read:   true,
			Create: true,
			Write:  true,
			Delete: true,
		}, expected: "rcwd"},
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
		{expected: FilePermissions{Create: true}, input: "c"},
		{expected: FilePermissions{Write: true}, input: "w"},
		{expected: FilePermissions{Delete: true}, input: "d"},
		{expected: FilePermissions{
			Read:   true,
			Create: true,
			Write:  true,
			Delete: true,
		}, input: "rcwd"},
		{expected: FilePermissions{
			Read:   true,
			Create: true,
			Write:  true,
			Delete: true,
		}, input: "wcrd"}, // Wrong order parses correctly
	}
	for _, c := range testdata {
		permissions, err := parseFilePermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestFilePermissions_ParseNegative(t *testing.T) {
	_, err := parseFilePermissions("wcrdf") // Here 'f' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "102")
}

func TestGetCanonicalName(t *testing.T) {
	testdata := []struct {
		inputAccount  string
		inputShare    string
		inputFilePath string
		expected      string
	}{
		{inputAccount: "fakestorageaccount", inputShare: "fakestorageshare", expected: "/file/fakestorageaccount/fakestorageshare"},
		{inputAccount: "fakestorageaccount", inputShare: "fakestorageshare", inputFilePath: "fakestoragefile", expected: "/file/fakestorageaccount/fakestorageshare/fakestoragefile"},
		{inputAccount: "fakestorageaccount", inputShare: "fakestorageshare", inputFilePath: "fakestoragedirectory/fakestoragefile", expected: "/file/fakestorageaccount/fakestorageshare/fakestoragedirectory/fakestoragefile"},
		{inputAccount: "fakestorageaccount", inputShare: "fakestorageshare", inputFilePath: "fakestoragedirectory\\fakestoragefile", expected: "/file/fakestorageaccount/fakestorageshare/fakestoragedirectory/fakestoragefile"},
		{inputAccount: "fakestorageaccount", inputShare: "fakestorageshare", inputFilePath: "fakestoragedirectory", expected: "/file/fakestorageaccount/fakestorageshare/fakestoragedirectory"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, getCanonicalName(c.inputAccount, c.inputShare, c.inputFilePath))
	}
}

func TestFileSignatureValues_SignWithSharedKey(t *testing.T) {
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.Nil(t, err, "error creating valid shared key credentials.")

	expiryDate, err := time.Parse("2006-01-02", "2023-07-20")
	require.Nil(t, err, "error creating valid expiry date.")

	testdata := []struct {
		object        SignatureValues
		expected      QueryParameters
		expectedError error
	}{
		{
			object:        SignatureValues{ShareName: "fakestorageshare", Permissions: "r", ExpiryTime: expiryDate},
			expected:      QueryParameters{version: Version, permissions: "r", expiryTime: expiryDate, resource: "s"},
			expectedError: nil,
		},
		{
			object:        SignatureValues{ShareName: "fakestorageshare", Permissions: "", ExpiryTime: expiryDate},
			expected:      QueryParameters{},
			expectedError: errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
		{
			object:        SignatureValues{ShareName: "fakestorageshare", Permissions: "r", ExpiryTime: time.Time{}},
			expected:      QueryParameters{},
			expectedError: errors.New("service SAS is missing at least one of these: ExpiryTime or Permissions"),
		},
		{
			object:        SignatureValues{ShareName: "fakestorageshare", Permissions: "", ExpiryTime: time.Time{}, Identifier: "fakepolicyname"},
			expected:      QueryParameters{version: Version, resource: "s", identifier: "fakepolicyname"},
			expectedError: nil,
		},
	}
	for _, c := range testdata {
		act, err := c.object.SignWithSharedKey(cred)
		require.Equal(t, c.expectedError, err)
		// ignore signature value
		act.signature = ""
		require.Equal(t, c.expected, act)
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/stretchr/testify/require"
)

func TestAccountPermissions_String(t *testing.T) {
	testdata := []struct {
		input    AccountPermissions
		expected string
	}{
		{input: AccountPermissions{Read: true}, expected: "r"},
		{input: AccountPermissions{Write: true}, expected: "w"},
		{input: AccountPermissions{Delete: true}, expected: "d"},
		{input: AccountPermissions{DeletePreviousVersion: true}, expected: "x"},
		{input: AccountPermissions{PermanentDelete: true}, expected: "y"},
		{input: AccountPermissions{List: true}, expected: "l"},
		{input: AccountPermissions{Add: true}, expected: "a"},
		{input: AccountPermissions{Create: true}, expected: "c"},
		{input: AccountPermissions{Update: true}, expected: "u"},
		{input: AccountPermissions{Process: true}, expected: "p"},
		{input: AccountPermissions{Tag: true}, expected: "t"},
		{input: AccountPermissions{FilterByTags: true}, expected: "f"},
		{input: AccountPermissions{SetImmutabilityPolicy: true}, expected: "i"},
		{input: AccountPermissions{
			Read:                  true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Add:                   true,
			Create:                true,
			Update:                true,
			Process:               true,
			Tag:                   true,
			FilterByTags:          true,
			SetImmutabilityPolicy: true,
		}, expected: "rwdxylacupfti"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestAccountPermissions_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected AccountPermissions
	}{
		{expected: AccountPermissions{Read: true}, input: "r"},
		{expected: AccountPermissions{Write: true}, input: "w"},
		{expected: AccountPermissions{Delete: true}, input: "d"},
		{expected: AccountPermissions{DeletePreviousVersion: true}, input: "x"},
		{expected: AccountPermissions{PermanentDelete: true}, input: "y"},
		{expected: AccountPermissions{List: true}, input: "l"},
		{expected: AccountPermissions{Add: true}, input: "a"},
		{expected: AccountPermissions{Create: true}, input: "c"},
		{expected: AccountPermissions{Update: true}, input: "u"},
		{expected: AccountPermissions{Process: true}, input: "p"},
		{expected: AccountPermissions{Tag: true}, input: "t"},
		{expected: AccountPermissions{FilterByTags: true}, input: "f"},
		{expected: AccountPermissions{SetImmutabilityPolicy: true}, input: "i"},
		{expected: AccountPermissions{
			Read:                  true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Add:                   true,
			Create:                true,
			Update:                true,
			Process:               true,
			Tag:                   true,
			FilterByTags:          true,
			SetImmutabilityPolicy: true,
		}, input: "rwdxylacupfti"},
		{expected: AccountPermissions{
			Read:                  true,
			Write:                 true,
			Delete:                true,
			DeletePreviousVersion: true,
			PermanentDelete:       true,
			List:                  true,
			Add:                   true,
			Create:                true,
			Update:                true,
			Process:               true,
			Tag:                   true,
			FilterByTags:          true,
			SetImmutabilityPolicy: true,
		}, input: "trwlapdixfycu"},
	}
	for _, c := range testdata {
		permissions, err := parseAccountPermissions(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestAccountPermissions_ParseNegative(t *testing.T) {
	_, err := parseAccountPermissions("trwlapdixfycuz") // Here 'z' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "122")
}

func TestAccountResourceTypes_String(t *testing.T) {
	testdata := []struct {
		input    AccountResourceTypes
		expected string
	}{
		{input: AccountResourceTypes{Service: true}, expected: "s"},
		{input: AccountResourceTypes{Container: true}, expected: "c"},
		{input: AccountResourceTypes{Object: true}, expected: "o"},
		{input: AccountResourceTypes{
			Service:   true,
			Container: true,
			Object:    true,
		}, expected: "sco"},
	}
	for _, c := range testdata {
		require.Equal(t, c.expected, c.input.String())
	}
}

func TestAccountResourceTypes_Parse(t *testing.T) {
	testdata := []struct {
		input    string
		expected AccountResourceTypes
	}{
		{expected: AccountResourceTypes{Service: true}, input: "s"},
		{expected: AccountResourceTypes{Container: true}, input: "c"},
		{expected: AccountResourceTypes{Object: true}, input: "o"},
		{expected: AccountResourceTypes{
			Service:   true,
			Container: true,
			Object:    true,
		}, input: "sco"},
		{expected: AccountResourceTypes{
			Service:   true,
			Container: true,
			Object:    true,
		}, input: "osc"},
	}
	for _, c := range testdata {
		permissions, err := parseAccountResourceTypes(c.input)
		require.Nil(t, err)
		require.Equal(t, c.expected, permissions)
	}
}

func TestAccountResourceTypes_ParseNegative(t *testing.T) {
	_, err := parseAccountResourceTypes("scoz") // Here 'z' is invalid
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "122")
}

// TODO: Sign With Shared Key
// Negative Case
// Version not provided
// SignWithSharedKey tests

func TestAccountSignatureValues_SignWithSharedKeyNonUTCTimes(t *testing.T) {
	cred, err := exported.NewSharedKeyCredential("fakeaccountname", "AKIAIOSFODNN7EXAMPLE")
	require.NoError(t, err)

	sign := func(t *testing.T, start, expiry time.Time) QueryParameters {
		t.Helper()
		qp, err := AccountSignatureValues{
			Version:       "2020-10-02",
			Protocol:      ProtocolHTTPS,
			StartTime:     start,
			ExpiryTime:    expiry,
			Permissions:   "r",
			ResourceTypes: "o",
		}.SignWithSharedKey(cred)
		require.NoError(t, err)
		return qp
	}

	start := time.Date(2026, 10, 15, 11, 0, 0, 0, time.UTC)
	expiry := start.Add(time.Hour)
	expected := sign(t, start, expiry)
	for _, loc := range []*time.Location{time.FixedZone("UTC+09:00", 9*60*60), time.FixedZone("UTC-05:00", -5*60*60)} {
		t.Run(loc.String(), func(t *testing.T) {
			actual := sign(t, start.In(loc), expiry.In(loc))
			require.Equal(t, expected.Encode(), actual.Encode())
		})
	}
}

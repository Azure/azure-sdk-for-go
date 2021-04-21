// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"path/filepath"
	"strings"
	"testing"
)

func Test_GetAliasPath(t *testing.T) {
	testCases := []struct {
		original string
		expected string
	}{
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "cdn", "mgmt", "2015-06-01", "cdn"),
			filepath.Join("cdn", "mgmt", "cdn"),
		},
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "cdn", "mgmt", "2015-06-01", "cdn", "cdnapi"),
			filepath.Join("cdn", "mgmt", "cdn", "cdnapi"),
		},
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "cdn", "mgmt", "2015-06-01", "cdn", "v2"),
			filepath.Join("cdn", "mgmt", "cdn"),
		},
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "keyvault", "2016-10-01", "keyvault"),
			filepath.Join("keyvault", "keyvault"),
		},
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "keyvault", "2016-10-01", "keyvault", "v21"),
			filepath.Join("keyvault", "keyvault"),
		},
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "keyvault", "2016-10-01", "keyvault", "v21", "keyvaultapi"),
			filepath.Join("keyvault", "keyvault", "keyvaultapi"),
		},
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "keyvault", "mgmt", "2016-10-01", "keyvault"),
			filepath.Join("keyvault", "mgmt", "keyvault"),
		},
		{
			filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "datalake", "analytics", "2016-11-01-preview", "catalog"),
			filepath.Join("datalake", "analytics", "catalog"),
		},
	}

	const consistentSeperator = "/"

	pathNorm := func(location string) string {
		return strings.Replace(location, `\`, consistentSeperator, -1)
	}

	pathsEqual := func(left, right string) bool {
		left, right = pathNorm(left), pathNorm(right)
		pieceWiseLeft, pieceWiseRight := strings.Split(left, consistentSeperator), strings.Split(right, consistentSeperator)

		if len(pieceWiseLeft) != len(pieceWiseRight) {
			return false
		}

		for i, lval := range pieceWiseLeft {
			rval := pieceWiseRight[i]
			if lval != rval {
				return false
			}
		}
		return true
	}

	for _, tc := range testCases {
		t.Run(tc.original, func(t *testing.T) {
			got, err := getAliasPath(tc.original)
			if err != nil {
				t.Error(err)
			}

			if !pathsEqual(tc.expected, got) {
				t.Logf("\ngot: \t%q\nwant:\t%q", pathNorm(got), pathNorm(tc.expected))
				t.Fail()
			}
		})
	}
}

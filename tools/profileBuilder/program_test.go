package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func Test_getAliasPath(t *testing.T) {
	const profileName = "p1"
	testCases := []struct {
		original string
		expected string
	}{
		{
			filepath.Join("services", "cdn", "mgmt", "2015-06-01", "cdn"),
			filepath.Join(profileName, "cdn", "mgmt", "cdn"),
		},
		{
			filepath.Join("services", "keyvault", "2016-10-01", "keyvault"),
			filepath.Join(profileName, "keyvault", "keyvault"),
		},
		{
			filepath.Join("services", "keyvault", "mgmt", "2016-10-01", "keyvault"),
			filepath.Join(profileName, "keyvault", "mgmt", "keyvault"),
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
			got, err := getAliasPath(tc.original, profileName)
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

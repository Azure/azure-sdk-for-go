// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sasurl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppend(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		sasQuery string
		expected string
	}{
		{"no existing query", "https://acct.blob.core.windows.net/c/b", "sv=x&sig=y", "https://acct.blob.core.windows.net/c/b?sv=x&sig=y"},
		{"existing versionid", "https://acct.blob.core.windows.net/c/b?versionid=v1", "sv=x&sig=y", "https://acct.blob.core.windows.net/c/b?versionid=v1&sv=x&sig=y"},
		{"existing snapshot", "https://acct.blob.core.windows.net/c/b?snapshot=s1", "sv=x&sig=y", "https://acct.blob.core.windows.net/c/b?snapshot=s1&sv=x&sig=y"},
		{"existing sharesnapshot", "https://acct.file.core.windows.net/s/f?sharesnapshot=s1", "sv=x&sig=y", "https://acct.file.core.windows.net/s/f?sharesnapshot=s1&sv=x&sig=y"},
		{"empty sas query", "https://acct.blob.core.windows.net/c/b?versionid=v1", "", "https://acct.blob.core.windows.net/c/b?versionid=v1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, Append(tt.url, tt.sasQuery))
		})
	}
}

func TestAppendToAccountURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		sasQuery string
		expected string
	}{
		{"no trailing slash", "https://acct.blob.core.windows.net", "sv=x", "https://acct.blob.core.windows.net/?sv=x"},
		{"trailing slash", "https://acct.blob.core.windows.net/", "sv=x", "https://acct.blob.core.windows.net/?sv=x"},
		{"query without trailing slash", "https://acct.blob.core.windows.net?customparam=value", "sv=x", "https://acct.blob.core.windows.net/?customparam=value&sv=x"},
		{"query with trailing slash", "https://acct.blob.core.windows.net/?customparam=value", "sv=x", "https://acct.blob.core.windows.net/?customparam=value&sv=x"},
		{"empty sas query", "https://acct.blob.core.windows.net", "", "https://acct.blob.core.windows.net/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, AppendToAccountURL(tt.url, tt.sasQuery))
		})
	}
}

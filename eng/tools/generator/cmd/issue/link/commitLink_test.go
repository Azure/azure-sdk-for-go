// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package link_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/link"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
)

func TestCommitLink_GetCommitHash(t *testing.T) {
	testdata := []struct {
		input    string
		expected string
	}{
		{
			input:    "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
			expected: "c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
		},
		{
			input:    "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c#diff-708c2fb843b022cac4af8c6f996a527440c1e0d328abb81f54670747bf14ab1a",
			expected: "c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
		},
		{
			input:    "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c/somethingElse",
			expected: "c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
		},
	}
	ctx := context.Background()
	client := query.NewClient()
	for _, c := range testdata {
		l := link.NewCommitLink(ctx, client, "", c.input)
		hash, err := l.(link.CommitHashLink).GetCommitHash()
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if hash != c.expected {
			t.Fatalf("expected %s but got %s", c.expected, hash)
		}
	}
}

func TestNewCommitLink(t *testing.T) {
	testdata := []struct {
		input    string
		expected string
	}{
		{
			input:    "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
			expected: "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
		},
		{
			input:    "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c#diff-708c2fb843b022cac4af8c6f996a527440c1e0d328abb81f54670747bf14ab1a",
			expected: "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
		},
		{
			input:    "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c/somethingElse",
			expected: "commit/c40f8fafa2c9b1bc3bb81e1feea4ff715d48e00c",
		},
	}
	ctx := context.Background()
	client := query.NewClient()
	for _, c := range testdata {
		l := link.NewCommitLink(ctx, client, "", c.input)
		if l.GetReleaseLink() != c.expected {
			t.Fatalf("expected %s but got %s", c.expected, l.GetReleaseLink())
		}
	}
}

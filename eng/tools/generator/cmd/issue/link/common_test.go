// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package link_test

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/link"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
)

func TestGetReadmeFromPath(t *testing.T) {
	token := os.Getenv("ACCESS_TOKEN")
	if len(token) == 0 {
		// skip this test when the token is not specified
		t.Skip()
	}
	ctx := context.Background()
	client := query.NewClientWithAccessToken(ctx, token)

	testdata := []struct {
		input    string
		expected link.Readme
	}{
		{
			input:    "specification/compute/resource-manager/Microsoft.Compute/stable/2020-12-01/compute.json",
			expected: "specification/compute/resource-manager/readme.md",
		},
		{
			input:    "specification/compute/resource-manager/Microsoft.Compute/stable/2020-12-01",
			expected: "specification/compute/resource-manager/readme.md",
		},
	}

	for _, c := range testdata {
		readme, err := link.GetReadmeFromPath(ctx, client, c.input)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if readme != c.expected {
			t.Fatalf("expected %s but got %s", c.expected, readme)
		}
	}
}

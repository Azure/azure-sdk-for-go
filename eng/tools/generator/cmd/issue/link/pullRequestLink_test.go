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

func TestPullRequestLink_Resolve(t *testing.T) {
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
		code     link.Code
	}{
		{
			input:    "pull/13024",
			expected: "specification/communication/resource-manager/readme.md",
			code:     link.CodeSuccess,
		},
	}

	for _, c := range testdata {
		l := link.NewPullRequestLink(ctx, client, "", c.input)
		result, err := l.Resolve()
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if result.GetCode() != c.code {
			t.Fatalf("expect code %v but got %v", c.code, result.GetCode())
		}
		if result.GetReadme() != c.expected {
			t.Fatalf("expect %s but got %s", c.expected, result.GetReadme())
		}
	}
}

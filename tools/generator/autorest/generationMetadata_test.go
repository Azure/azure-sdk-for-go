// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest_test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
)

func TestParse(t *testing.T) {
	tests := []struct {
		changelog string
		expected  autorest.GenerationMetadata
	}{
		{
			changelog: "Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/compute/resource-manager/readme.md tag: `package-2020-06-30`\n\nCode generator @microsoft.azure/autorest.go@2.1.168\n",
			expected: autorest.GenerationMetadata{
				CommitHash:     "3c764635e7d442b3e74caf593029fcd440b3ef82",
				Readme:         "specification/compute/resource-manager/readme.md",
				Tag:            "package-2020-06-30",
				CodeGenVersion: "@microsoft.azure/autorest.go@2.1.168",
			},
		},
	}

	for _, c := range tests {
		reader := strings.NewReader(c.changelog)
		m, err := autorest.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if m.String() != c.expected.String() {
			t.Fatalf("expect %+v, but got %+v", c.expected, *m)
		}
	}
}

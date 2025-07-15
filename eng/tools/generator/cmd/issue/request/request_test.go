// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package request_test

// func TestParseReadmeFromLink(t *testing.T) {
// 	testData := []struct {
// 		name     string
// 		input    string
// 		expected link.Readme
// 	}{
// 		{
// 			name:     "parse PR link",
// 			input:    "https://github.com/Azure/azure-rest-api-specs/pull/11377",
// 			expected: "specification/security/resource-manager/readme.md",
// 		},
// 		{
// 			name:     "parse github directory",
// 			input:    "https://github.com/Azure/azure-rest-api-specs/tree/636ade9c4b88ad7e408d85f564ddedbc638ae524/specification/monitor/resource-manager",
// 			expected: "specification/monitor/resource-manager/readme.md",
// 		},
// 		{
// 			name:     "parse github file",
// 			input:    "https://github.com/Azure/azure-rest-api-specs/blob/636ade9c4b88ad7e408d85f564ddedbc638ae524/specification/storage/resource-manager/readme.md",
// 			expected: "specification/storage/resource-manager/readme.md",
// 		},
// 	}

// 	for _, c := range testData {
// 		t.Logf("Testing %s...", c.name)
// 		client := &query.Client{
// 			Client: github.NewClient(nil),
// 		}
// 		result, err := request.ParseReadmeFromLink(context.Background(), client, request.ReleaseRequestIssue{
// 			IssueLink:  "",
// 			TargetLink: c.input,
// 		})
// 		if err != nil {
// 			t.Fatalf("unexpected error: %+v", err)
// 		}
// 		if result.GetReadme() != c.expected {
// 			t.Fatalf("expected %q but got %q", c.expected, result.GetReadme())
// 		}
// 	}
// }

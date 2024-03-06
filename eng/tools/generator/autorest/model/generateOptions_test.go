// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
)

func TestNewGenerateOptionsFrom(t *testing.T) {
	testdata := []struct {
		input    string
		expected *model.GenerateOptions
	}{
		{
			input: `{
  "autorestArguments": [
    "--use=@microsoft.azure/autorest.go@2.1.178",
    "--go",
    "--verbose",
    "--go-sdk-folder=.",
    "--multiapi",
    "--use-onever",
    "--version=V2",
    "--go.license-header=MICROSOFT_MIT_NO_VERSION"
  ],
  "additionalOptions": [
    "--enum-prefix"
  ]
}`,
			expected: &model.GenerateOptions{
				AutorestArguments: []string{
					"--use=@microsoft.azure/autorest.go@2.1.178",
					"--go",
					"--verbose",
					"--go-sdk-folder=.",
					"--multiapi",
					"--use-onever",
					"--version=V2",
					"--go.license-header=MICROSOFT_MIT_NO_VERSION",
				},
				AdditionalOptions: []string{
					"--enum-prefix",
				},
			},
		},
	}

	for _, c := range testdata {
		options, err := model.NewGenerateOptionsFrom(strings.NewReader(c.input))
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !reflect.DeepEqual(*options, *c.expected) {
			t.Fatalf("expected %+v, but got %+v", *c.expected, *options)
		}
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model_test

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
)

func TestNewOption(t *testing.T) {
	testcase := []struct {
		input    string
		expected model.Option
	}{
		{
			input:    "specification/compute/resource-manager/readme.md",
			expected: model.NewArgument("specification/compute/resource-manager/readme.md"),
		},
		{
			input:    "--multiapi",
			expected: model.NewFlagOption("multiapi"),
		},
		{
			input:    "--tag=package-2020-01-01",
			expected: model.NewKeyValueOption("tag", "package-2020-01-01"),
		},
	}

	for _, c := range testcase {
		o, err := model.NewOption(c.input)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !reflect.DeepEqual(o, c.expected) {
			t.Fatalf("expecting %+v, but got %+v", c.expected, o)
		}
	}
}

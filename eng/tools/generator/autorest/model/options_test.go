// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model_test

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
)

func TestParseOptions(t *testing.T) {
	testdata := []struct {
		input    []string
		expected []model.Option
	}{
		{
			input: []string{
				"--use=@microsoft.azure/autorest.go@2.1.178",
				"--go",
				"--verbose",
				"--go-sdk-folder=.",
				"--multiapi",
				"--use-onever",
				"--version=V2",
				"--go.license-header=MICROSOFT_MIT_NO_VERSION",
			},
			expected: []model.Option{
				model.NewKeyValueOption("use", "@microsoft.azure/autorest.go@2.1.178"),
				model.NewFlagOption("go"),
				model.NewFlagOption("verbose"),
				model.NewKeyValueOption("go-sdk-folder", "."),
				model.NewFlagOption("multiapi"),
				model.NewFlagOption("use-onever"),
				model.NewKeyValueOption("version", "V2"),
				model.NewKeyValueOption("go.license-header", "MICROSOFT_MIT_NO_VERSION"),
			},
		},
	}

	for _, c := range testdata {
		options, err := model.ParseOptions(c.input)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !reflect.DeepEqual(options.Arguments(), c.expected) {
			t.Fatalf("expected %+v, but got %+v", c.expected, options.Arguments())
		}
	}
}

func TestLocalOptions_MergeOptions(t *testing.T) {
	testcase := []struct {
		description string
		input       model.Options
		newOptions  []model.Option
		expected    model.Options
	}{
		{
			description: "merge a new argument",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewArgument("nothing"),
			},
			expected: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
				model.NewArgument("nothing"),
			),
		},
		{
			description: "merge multiple arguments",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewArgument("nothing"),
				model.NewArgument("anything"),
			},
			expected: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
				model.NewArgument("nothing"),
				model.NewArgument("anything"),
			),
		},
		{
			description: "merge a new flag (unique)",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewFlagOption("nothing"),
			},
			expected: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
				model.NewFlagOption("nothing"),
			),
		},
		{
			description: "merge a new flag (duplicate)",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewFlagOption("a"),
			},
			expected: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
		},
		{
			description: "merge multiple flags",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewFlagOption("a"),
				model.NewFlagOption("nothing"),
			},
			expected: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
				model.NewFlagOption("nothing"),
			),
		},
		{
			description: "merge a key value option (unique)",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewKeyValueOption("d", "something"),
			},
			expected: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
				model.NewKeyValueOption("d", "something"),
			),
		},
		{
			description: "merge a key value option (duplicate with key value)",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewKeyValueOption("b", "something"),
			},
			expected: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "something"),
				model.NewArgument("something"),
			),
		},
		{
			description: "merge a key value option (duplicate with flag)",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewKeyValueOption("a", "something"),
			},
			expected: model.NewOptions(
				model.NewKeyValueOption("a", "something"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
		},
		{
			description: "merge multiple key value options",
			input: model.NewOptions(
				model.NewFlagOption("a"),
				model.NewKeyValueOption("b", "c"),
				model.NewArgument("something"),
			),
			newOptions: []model.Option{
				model.NewKeyValueOption("a", "something"),
				model.NewKeyValueOption("b", "anything"),
				model.NewKeyValueOption("d", "nothing"),
			},
			expected: model.NewOptions(
				model.NewKeyValueOption("a", "something"),
				model.NewKeyValueOption("b", "anything"),
				model.NewArgument("something"),
				model.NewKeyValueOption("d", "nothing"),
			),
		},
	}

	for _, c := range testcase {
		t.Log(c.description)
		r := c.input.MergeOptions(c.newOptions...)
		if !reflect.DeepEqual(r, c.expected) {
			t.Fatalf("expected %+v, but got %+v", c.expected, r)
		}
	}
}

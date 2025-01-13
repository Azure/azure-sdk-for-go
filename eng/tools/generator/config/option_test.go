// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config_test

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
)

func TestNewOption(t *testing.T) {
	testcase := []struct {
		input    string
		expected config.Option
	}{
		{
			input:    "specification/compute/resource-manager/readme.md",
			expected: config.NewArgument("specification/compute/resource-manager/readme.md"),
		},
		{
			input:    "--multiapi",
			expected: config.NewFlagOption("multiapi"),
		},
		{
			input:    "--tag=package-2020-01-01",
			expected: config.NewKeyValueOption("tag", "package-2020-01-01"),
		},
	}

	for _, c := range testcase {
		o, err := config.NewOption(c.input)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !reflect.DeepEqual(o, c.expected) {
			t.Fatalf("expecting %+v, but got %+v", c.expected, o)
		}
	}
}

func TestParseOptions(t *testing.T) {
	testdata := []struct {
		input    []string
		expected []config.Option
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
			expected: []config.Option{
				config.NewKeyValueOption("use", "@microsoft.azure/autorest.go@2.1.178"),
				config.NewFlagOption("go"),
				config.NewFlagOption("verbose"),
				config.NewKeyValueOption("go-sdk-folder", "."),
				config.NewFlagOption("multiapi"),
				config.NewFlagOption("use-onever"),
				config.NewKeyValueOption("version", "V2"),
				config.NewKeyValueOption("go.license-header", "MICROSOFT_MIT_NO_VERSION"),
			},
		},
	}

	for _, c := range testdata {
		options, err := config.ParseOptions(c.input)
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
		input       config.Options
		newOptions  []config.Option
		expected    config.Options
	}{
		{
			description: "merge a new argument",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewArgument("nothing"),
			},
			expected: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
				config.NewArgument("nothing"),
			),
		},
		{
			description: "merge multiple arguments",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewArgument("nothing"),
				config.NewArgument("anything"),
			},
			expected: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
				config.NewArgument("nothing"),
				config.NewArgument("anything"),
			),
		},
		{
			description: "merge a new flag (unique)",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewFlagOption("nothing"),
			},
			expected: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
				config.NewFlagOption("nothing"),
			),
		},
		{
			description: "merge a new flag (duplicate)",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewFlagOption("a"),
			},
			expected: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
		},
		{
			description: "merge multiple flags",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewFlagOption("a"),
				config.NewFlagOption("nothing"),
			},
			expected: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
				config.NewFlagOption("nothing"),
			),
		},
		{
			description: "merge a key value option (unique)",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewKeyValueOption("d", "something"),
			},
			expected: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
				config.NewKeyValueOption("d", "something"),
			),
		},
		{
			description: "merge a key value option (duplicate with key value)",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewKeyValueOption("b", "something"),
			},
			expected: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "something"),
				config.NewArgument("something"),
			),
		},
		{
			description: "merge a key value option (duplicate with flag)",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewKeyValueOption("a", "something"),
			},
			expected: config.NewOptions(
				config.NewKeyValueOption("a", "something"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
		},
		{
			description: "merge multiple key value options",
			input: config.NewOptions(
				config.NewFlagOption("a"),
				config.NewKeyValueOption("b", "c"),
				config.NewArgument("something"),
			),
			newOptions: []config.Option{
				config.NewKeyValueOption("a", "something"),
				config.NewKeyValueOption("b", "anything"),
				config.NewKeyValueOption("d", "nothing"),
			},
			expected: config.NewOptions(
				config.NewKeyValueOption("a", "something"),
				config.NewKeyValueOption("b", "anything"),
				config.NewArgument("something"),
				config.NewKeyValueOption("d", "nothing"),
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

package model_ext

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
)

func TestOptions_MergeOptions(t *testing.T) {
	testcase := []struct {
		description string
		input       Options
		newOptions  []model.Option
		expected    *Options
	}{
		{
			description: "merge a new argument",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewArgument("nothing"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
					model.NewArgument("nothing"),
				},
			},
		},
		{
			description: "merge multiple arguments",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewArgument("nothing"),
				model.NewArgument("anything"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
					model.NewArgument("nothing"),
					model.NewArgument("anything"),
				},
			},
		},
		{
			description: "merge a new flag (unique)",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewFlagOption("nothing"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
					model.NewFlagOption("nothing"),
				},
			},
		},
		{
			description: "merge a new flag (duplicate)",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewFlagOption("a"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
		},
		{
			description: "merge multiple flags",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewFlagOption("a"),
				model.NewFlagOption("nothing"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
					model.NewFlagOption("nothing"),
				},
			},
		},
		{
			description: "merge a key value option (unique)",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewKeyValueOption("d", "something"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
					model.NewKeyValueOption("d", "something"),
				},
			},
		},
		{
			description: "merge a key value option (duplicate with key value)",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewKeyValueOption("b", "something"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "something"),
					model.NewArgument("something"),
				},
			},
		},
		{
			description: "merge a key value option (duplicate with flag)",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewKeyValueOption("a", "something"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewKeyValueOption("a", "something"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
		},
		{
			description: "merge multiple key value options",
			input: Options{
				autorestArguments: []model.Option{
					model.NewFlagOption("a"),
					model.NewKeyValueOption("b", "c"),
					model.NewArgument("something"),
				},
			},
			newOptions: []model.Option{
				model.NewKeyValueOption("a", "something"),
				model.NewKeyValueOption("b", "anything"),
				model.NewKeyValueOption("d", "nothing"),
			},
			expected: &Options{
				autorestArguments: []model.Option{
					model.NewKeyValueOption("a", "something"),
					model.NewKeyValueOption("b", "anything"),
					model.NewArgument("something"),
					model.NewKeyValueOption("d", "nothing"),
				},
			},
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

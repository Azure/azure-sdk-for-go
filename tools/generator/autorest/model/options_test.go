package model_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
)

func TestNewOptionsFrom(t *testing.T) {
	testdata := []struct {
		input    string
		expected []model.Option
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
  ]
}`,
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
		options, err := model.NewOptionsFrom(strings.NewReader(c.input))
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !reflect.DeepEqual(options.Arguments(), c.expected) {
			t.Fatalf("expected %+v, but got %+v", c.expected, options.Arguments())
		}
	}
}

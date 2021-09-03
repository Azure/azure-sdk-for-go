// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
)

func TestReadBatchTags(t *testing.T) {
	testdata := []struct {
		readmePath string
		expected   []string
		err        string
	}{
		{
			readmePath: "./testdata/proper_readme.go.md",
			expected: []string{
				"package-2020-01",
				"package-2019-01-preview",
			},
		},
		{
			readmePath: "./testdata/multiple_multiapi_readme.go.md",
			err:        "multiple multiapi section found on line 14 and 19, we should only have one",
		},
		{
			readmePath: "./testdata/missing_multiapi_readme.go.md",
			err:        "cannot find multiapi section",
		},
		{
			readmePath: "./testdata/too_short_readme.go.md",
			err:        "multiapi section cannot be parsed",
		},
		{
			readmePath: "./testdata/missing_batch_readme.go.md",
			err:        "multiapi section should begin with `batch:`",
		},
	}

	for _, c := range testdata {
		t.Logf("Testing %s", c.readmePath)
		reader, err := os.Open(c.readmePath)
		if err != nil {
			t.Fatalf("unexpected error when opening readme file: %+v", err)
		}

		tags, err := autorest.ReadBatchTags(reader)
		if err != nil {
			if c.expected != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if err.Error() != c.err {
				t.Fatalf("expected error message '%s', but got error message '%s'", c.err, err.Error())
			}
		} else {
			if !reflect.DeepEqual(c.expected, tags) {
				t.Fatalf("expected %+v, but got %+v", c.expected, tags)
			}
		}
	}
}

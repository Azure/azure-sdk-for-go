package autorest

import (
	"reflect"
	"testing"
)

func TestGetChangedPackages(t *testing.T) {
	testData := []struct {
		description  string
		changedFiles []string
		expected     map[string][]string
	}{
		{
			description: "one file changed in one package",
			changedFiles: []string{
				"../testdata/2020-10-29/foo/models.go",
			},
			expected: map[string][]string{
				"../testdata/2020-10-29/foo": {
					"../testdata/2020-10-29/foo/models.go",
				},
			},
		},
		{
			description: "two files changed in one package",
			changedFiles: []string{
				"../testdata/2020-10-29/foo/models.go",
				"../testdata/2020-10-29/foo/version.go",
			},
			expected: map[string][]string{
				"../testdata/2020-10-29/foo": {
					"../testdata/2020-10-29/foo/models.go",
					"../testdata/2020-10-29/foo/version.go",
				},
			},
		},
		{
			description: "multiple files changed in two packages",
			changedFiles: []string{
				"../testdata/2020-10-29/foo/models.go",
				"../testdata/2020-10-29/foo/version.go",
				"../testdata/2020-10-30/foo/models.go",
			},
			expected: map[string][]string{
				"../testdata/2020-10-29/foo": {
					"../testdata/2020-10-29/foo/models.go",
					"../testdata/2020-10-29/foo/version.go",
				},
				"../testdata/2020-10-30/foo": {
					"../testdata/2020-10-30/foo/models.go",
				},
			},
		},
		{
			description: "one directory untracked and one file changed in one package",
			changedFiles: []string{
				"../testdata/2020-10-29/foo",
				"../testdata/2020-10-30/foo/models.go",
			},
			expected: map[string][]string{
				"../testdata/2020-10-29/foo": {
					"../testdata/2020-10-29/foo/client.go",
					"../testdata/2020-10-29/foo/models.go",
					"../testdata/2020-10-29/foo/version.go",
				},
				"../testdata/2020-10-30/foo": {
					"../testdata/2020-10-30/foo/models.go",
				},
			},
		},
		{
			description: "two untracked directories",
			changedFiles: []string{
				"../testdata/2020-10-29/foo",
				"../testdata/2020-10-30/foo",
			},
			expected: map[string][]string{
				"../testdata/2020-10-29/foo": {
					"../testdata/2020-10-29/foo/client.go",
					"../testdata/2020-10-29/foo/models.go",
					"../testdata/2020-10-29/foo/version.go",
				},
				"../testdata/2020-10-30/foo": {
					"../testdata/2020-10-30/foo/client.go",
					"../testdata/2020-10-30/foo/models.go",
					"../testdata/2020-10-30/foo/version.go",
				},
			},
		},
		{
			description: "one untracked directory that contains one packages",
			changedFiles: []string{
				"../testdata/2020-10-29",
			},
			expected: map[string][]string{
				"../testdata/2020-10-29/foo": {
					"../testdata/2020-10-29/foo/client.go",
					"../testdata/2020-10-29/foo/models.go",
					"../testdata/2020-10-29/foo/version.go",
				},
			},
		},
		{
			description: "one untracked directory that contains multiple packages",
			changedFiles: []string{
				"../testdata/",
			},
			expected: map[string][]string{
				"../testdata/2020-10-29/foo": {
					"../testdata/2020-10-29/foo/client.go",
					"../testdata/2020-10-29/foo/models.go",
					"../testdata/2020-10-29/foo/version.go",
				},
				"../testdata/2020-10-30/foo": {
					"../testdata/2020-10-30/foo/client.go",
					"../testdata/2020-10-30/foo/models.go",
					"../testdata/2020-10-30/foo/version.go",
				},
			},
		},
	}

	for _, c := range testData {
		t.Logf("testing %s", c.description)
		r, err := GetChangedPackages(c.changedFiles)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !mapDeepEqual(r, c.expected) {
			t.Fatalf("expected %+v, but got %+v", c.expected, r)
		}
	}
}

func TestExpandChangedDirectories(t *testing.T) {
	testData := []struct {
		description string
		input       []string
		expected    []string
	}{
		{
			description: "only files",
			input: []string{
				"../testdata/2020-10-29/foo/client.go",
			},
			expected: []string{
				"../testdata/2020-10-29/foo/client.go",
			},
		},
		{
			description: "only directories",
			input: []string{
				"../testdata/2020-10-29/foo",
			},
			expected: []string{
				"../testdata/2020-10-29/foo/client.go",
				"../testdata/2020-10-29/foo/models.go",
				"../testdata/2020-10-29/foo/version.go",
			},
		},
		{
			description: "both directories and files",
			input: []string{
				"../testdata/2020-10-29/foo",
				"../testdata/2020-10-30/foo/models.go",
			},
			expected: []string{
				"../testdata/2020-10-29/foo/client.go",
				"../testdata/2020-10-29/foo/models.go",
				"../testdata/2020-10-29/foo/version.go",
				"../testdata/2020-10-30/foo/models.go",
			},
		},
		{
			description: "multiple hierarchy of directories but only one sub-directory",
			input: []string{
				"../testdata/2020-10-29",
			},
			expected: []string{
				"../testdata/2020-10-29/foo/client.go",
				"../testdata/2020-10-29/foo/models.go",
				"../testdata/2020-10-29/foo/version.go",
			},
		},
		{
			description: "multiple hierarchy of directories",
			input: []string{
				"../testdata",
			},
			expected: []string{
				"../testdata/2020-10-29/foo/client.go",
				"../testdata/2020-10-29/foo/models.go",
				"../testdata/2020-10-29/foo/version.go",
				"../testdata/2020-10-30/foo/client.go",
				"../testdata/2020-10-30/foo/models.go",
				"../testdata/2020-10-30/foo/version.go",
			},
		},
	}

	for _, c := range testData {
		t.Logf("testing %s", c.description)
		r, err := ExpandChangedDirectories(c.input)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if !reflect.DeepEqual(r, c.expected) {
			t.Fatalf("expected %v but got %v", c.expected, r)
		}
	}
}

// subsetOf return true if m2 is the subset of m1 (every key in m1 exists in m2 and the corresponding values are the same)
func subsetOf(m1, m2 map[string][]string) bool {
	for k, v := range m1 {
		if v2, ok := m2[k]; !ok || !reflect.DeepEqual(v, v2) {
			return false
		}
	}
	return true
}

func mapDeepEqual(m1, m2 map[string][]string) bool {
	return subsetOf(m1, m2) && subsetOf(m2, m1)
}

func TestIsValidPackage(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "../../../services/compute/mgmt/2020-06-01/compute",
			expected: true,
		},
		{
			input:    "../../../services/compute/mgmt/2020-06-01",
			expected: false,
		},
		{
			input:    "../../../storage",
			expected: false,
		},
		{
			input:    "../../../profiles",
			expected: false,
		},
	}

	for _, c := range testData {
		r := IsValidPackage(c.input)
		if r != c.expected {
			t.Fatalf("expected %v but got %v", c.expected, r)
		}
	}
}

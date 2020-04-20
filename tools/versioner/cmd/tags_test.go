// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/report"
	"github.com/Azure/azure-sdk-for-go/tools/internal/modinfo"
)

func Test_getTags(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	tags, err := getTags(cwd, "v10")
	if err != nil {
		t.Fatalf("failed to get tags: %v", err)
	}
	if l := len(tags); l != 11 {
		t.Fatalf("expected 11 tags, got %d", l)
	}
	found := false
	for _, tag := range tags {
		if tag == "v10.1.0-beta" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("didn't find tag v10.1.0-beta")
	}
}

func Test_getTagPrefix(t *testing.T) {
	testData := []struct {
		dir      string
		expected string
		errored  bool
	}{
		{
			dir:      filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "redis", "mgmt", "2018-03-01", "redis"),
			expected: "services/redis/mgmt/2018-03-01/redis",
		},
		{
			dir:      filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "redis", "mgmt", "2018-03-01", "redis", "v2"),
			expected: "services/redis/mgmt/2018-03-01/redis/v2",
		},
		{
			dir:     filepath.Join("work", "src", "github.com", "something", "else"),
			errored: true,
		},
	}

	const repoRoot = "github.com/Azure/azure-sdk-for-go"

	for _, c := range testData {
		p, err := getTagPrefix(c.dir, repoRoot)
		if err != nil && !c.errored {
			t.Fatalf("unexpected error for case '%s': %+v", c.dir, err)
		}
		if err == nil && c.errored {
			t.Fatalf("expected error but got nothing")
		}
		if p != c.expected {
			t.Fatalf("expected '%s' but got '%s'", c.expected, p)
		}
	}
}

type mockModInfo struct {
	isPreview bool
	exports   bool
	breaks    bool
	newModule bool
}

func (mock mockModInfo) DestDir() string {
	return ""
}

func (mock mockModInfo) IsARMPackage() bool {
	return false
}

func (mock mockModInfo) IsPreviewPackage() bool {
	return mock.isPreview
}

func (mock mockModInfo) NewExports() bool {
	return mock.exports
}

func (mock mockModInfo) BreakingChanges() bool {
	return mock.breaks
}

func (mock mockModInfo) VersionSuffix() bool {
	return mock.breaks && !mock.isPreview
}

func (mock mockModInfo) NewModule() bool {
	return mock.newModule
}

func (mock mockModInfo) GenerateReport() report.Package {
	// not needed by tests
	return report.Package{}
}

func TestCalculateModuleTag(t *testing.T) {
	testData := []struct {
		name     string
		baseline string
		pkg      modinfo.Provider
		hookFunc TagsHookFunc
		errored  bool
		expected string
	}{
		{
			name:     "major version v1",
			baseline: filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				newModule: true,
			},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{}, nil
			},
			expected: "services/foo/v1.0.0",
		},
		{
			name:     "major version v2",
			baseline: filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				breaks: true,
			},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{
					"services/foo/v1.0.0",
					"services/foo/v1.1.0",
				}, nil
			},
			expected: "services/foo/v2.0.0",
		},
		{
			name:     "minor v1",
			baseline: filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				exports: true,
			},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{
					"services/foo/v1.0.0",
					"services/foo/v1.0.1",
				}, nil
			},
			expected: "services/foo/v1.1.0",
		},
		{
			name:     "patch v2",
			baseline: filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg:      mockModInfo{},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{
					"services/foo/v1.0.0",
					"services/foo/v1.0.1",
					"services/foo/v2.0.0",
				}, nil
			},
			expected: "services/foo/v2.0.1",
		},
		{
			name: "major v3",
			baseline: filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				exports: false,
				breaks:  true,
			},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{
					"services/foo/v1.0.0",
					"services/foo/v2.0.0",
					"services/foo/v2.1.0",
					"services/foo/v2.1.1",
				}, nil
			},
			expected: "services/foo/v3.0.0",
		},
		{
			name: "new preview package",
			baseline:filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				isPreview: true,
				exports:   true,
				breaks:   true,
				newModule: true,
			},
			hookFunc: func(root string, tagPrefix string) (	[]string, error) {
				return []string{}, nil
			},
			expected: "services/foo/v0.0.0",
		},
		{
			name: "preview package breaking change",
			baseline:filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				isPreview: true,
				breaks: true,
			},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{
					"services/foo/v0.0.0",
					"services/foo/v0.1.0",
					"services/foo/v0.1.1",
				}, nil
			},
			expected: "services/foo/v0.2.0",
		},
		{
			name: "preview package incremental change",
			baseline:filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				isPreview: true,
				exports: true,
			},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{
					"services/foo/v0.0.0",
					"services/foo/v0.1.0",
					"services/foo/v0.1.1",
				}, nil
			},
			expected: "services/foo/v0.2.0",
		},
		{
			name: "preview package patch",
			baseline:filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo"),
			pkg: mockModInfo{
				isPreview: true,
			},
			hookFunc: func(root string, tagPrefix string) ([]string, error) {
				return []string{
					"services/foo/v0.0.0",
					"services/foo/v0.1.0",
					"services/foo/v0.1.1",
				}, nil
			},
			expected: "services/foo/v0.1.2",
		},
	}

	const repoRoot = "github.com/Azure/azure-sdk-for-go"

	for _, c := range testData {
		t.Logf("Testing %s", c.name)
		versionSetting, _ := parseVersionSetting()
		tag, err := calculateModuleTag(c.baseline, versionSetting, repoRoot, c.pkg, c.hookFunc)
		if err != nil && !c.errored {
			t.Fatalf("unexpected error: %+v", err)
		}
		if err == nil && c.errored {
			t.Fatalf("expected error but got nothing")
		}
		if tag != c.expected {
			t.Fatalf("expected '%s' but got '%s'", c.expected, tag)
		}
	}
}

func TestGetLatestSemver(t *testing.T) {
	testData := []struct{
		tags []string
		tagPrefix string
		expected string
	}{
		{
			tags: []string{},
			expected: "",
		},
		{
			tags: []string{
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.0.0",
			},
			tagPrefix: "services/foo",
			expected: "1.0.0",
		},
		{
			tags: []string{
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.1.0",
			},
			tagPrefix: "services/foo",
			expected: "1.1.0",
		},
		{
			tags: []string{
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.1.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.0.0",
			},
			tagPrefix: "services/foo",
			expected: "2.0.0",
		},
		{
			tags: []string{
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.1.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.0.1",
			},
			tagPrefix: "services/foo",
			expected: "2.0.1",
		},
		{
			tags: []string{
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.1.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.0.1",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.1.0",
			},
			tagPrefix: "services/foo",
			expected: "2.1.0",
		},
		{
			tags: []string{
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v1.1.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.0.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.0.1",
				"github.com/Azure/azure-sdk-for-go/services/foo/v2.1.0",
				"github.com/Azure/azure-sdk-for-go/services/foo/v3.0.0",
			},
			tagPrefix: "services/foo",
			expected: "3.0.0",
		},
	}

	for _, c := range testData {
		v, err := getLatestSemver(c.tags, c.tagPrefix)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if v == nil {
			if c.expected != "" {
				t.Fatalf("expected %s but got nothing", c.expected)
			}
			continue
		}
		if v.String() != c.expected {
			t.Fatalf("expected %s but got %s", c.expected, v.String())
		}
	}
}

func TestFindVersionSuffixInTag(t *testing.T) {
	testData := []struct{
		tag string
		expected string
	}{
		{
			tag: "services/foo/v1.0.0",
			expected: "",
		},
		{
			tag: "services/foo/v1.2.100",
			expected: "",
		},
		{
			tag: "services/foo/v2.0.0",
			expected: "v2",
		},
		{
			tag: "services/foo/v1000.60.100",
			expected: "v1000",
		},
		{
			tag: "services/foov1999/v3.0.1",
			expected: "v3",
		},
		{
			tag: "services/foo/v1999/v3.0.1",
			expected: "v3",
		},
	}

	for _, c := range testData {
		ver := findVersionSuffixInTag(c.tag)
		if ver != c.expected {
			t.Fatalf("expected '%s' but got '%s'", c.expected, ver)
		}
	}
}

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

package modinfo

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetModuleInfo(t *testing.T) {
	type expected struct {
		breakingChange bool
		newExports     bool
		versionSuffix  bool
		destDir        string
	}
	testData := []struct {
		name string
		baseline string
		stage string
		expected
	}{
		{
			name: "scenario a",
			baseline: "../../testdata/scenarioa/foo",
			stage: "../../testdata/scenarioa/foo/stage",
			expected: expected{
				breakingChange: false,
				newExports:     true,
				versionSuffix:  false,
				destDir:        "../../testdata/scenarioa/foo",
			},
		},
		{
			name: "scenario b",
			baseline: "../../testdata/scenariob/foo",
			stage: "../../testdata/scenariob/foo/stage",
			expected: expected{
				breakingChange: true,
				newExports:     true,
				versionSuffix:  true,
				destDir:        "../../testdata/scenariob/foo/v2",
			},
		},
		{
			name: "scenario c",
			baseline: "../../testdata/scenarioc/foo",
			stage: "../../testdata/scenarioc/foo/stage",
			expected: expected{
				breakingChange: false,
				newExports:     false,
				versionSuffix:  false,
				destDir:        "../../testdata/scenarioc/foo",
			},
		},
		{
			name: "scenario d",
			baseline: "../../testdata/scenariod/foo",
			stage: "../../testdata/scenariod/foo/stage",
			expected: expected{
				breakingChange: true,
				newExports:     false,
				versionSuffix:  true,
				destDir:        "../../testdata/scenariod/foo/v2",
			},
		},
		{
			name: "scenario e",
			baseline: "../../testdata/scenarioe/foo/v2",
			stage: "../../testdata/scenarioe/foo/stage",
			expected: expected{
				breakingChange: false,
				newExports:     true,
				versionSuffix:  true,
				destDir:        "../../testdata/scenarioe/foo/v2",
			},
		},
		{
			name: "scenario f",
			baseline: "../../testdata/scenariof/foo",
			stage: "../../testdata/scenariof/foo/stage",
			expected: expected{
				breakingChange: false,
				newExports:     true,
				versionSuffix:  false,
				destDir:        "../../testdata/scenariof/foo",
			},
		},
	}

	for _, c := range testData {
		t.Logf("Testing %s", c.name)
		mod, err := GetModuleInfo(c.baseline, c.stage)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if mod.BreakingChanges() != c.expected.breakingChange {
			t.Fatalf("breaking changes: expected %t but got %t", c.expected.breakingChange, mod.BreakingChanges())
		}
		if mod.NewExports() != c.expected.newExports {
			t.Fatalf("new exports: expected %t but got %t", c.expected.newExports, mod.NewExports())
		}
		if mod.VersionSuffix() != c.expected.versionSuffix {
			t.Fatalf("version suffix: expected %t but got %t", c.expected.versionSuffix, mod.VersionSuffix())
		}
		destDir, err := filepath.Abs(mod.DestDir())
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		expectedDest, err := filepath.Abs(c.expected.destDir)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if destDir != expectedDest {
			t.Fatalf("destDir: expected %v but got %v", expectedDest, destDir)
		}
	}
}

func Test_sortModuleTagsBySemver(t *testing.T) {
	before := []string{
		"v1.0.0",
		"v1.0.1",
		"v1.1.0",
		"v10.0.0",
		"v11.1.1",
		"v2.0.0",
		"v20.2.3",
		"v3.1.0",
	}
	sortModuleTagsBySemver(before)
	after := []string{
		"v1.0.0",
		"v1.0.1",
		"v1.1.0",
		"v2.0.0",
		"v3.1.0",
		"v10.0.0",
		"v11.1.1",
		"v20.2.3",
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("sort order doesn't match, expected '%v' got '%v'", after, before)
	}
}

func TestIncrementModuleVersion(t *testing.T) {
	v := IncrementModuleVersion("")
	if v != "v2" {
		t.Fatalf("expected v2 got %s", v)
	}
	v = IncrementModuleVersion("v2")
	if v != "v3" {
		t.Fatalf("expected v3 got %s", v)
	}
	v = IncrementModuleVersion("v10")
	if v != "v11" {
		t.Fatalf("expected v11 got %s", v)
	}
}

func TestCreateModuleNameFromPath(t *testing.T) {
	n, err := CreateModuleNameFromPath(filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "foo", "apiver", "foo"))
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	const expected = "github.com/Azure/azure-sdk-for-go/services/foo/apiver/foo"
	if n != expected {
		t.Fatalf("expected '%s' got '%s'", expected, n)
	}
}

func TestCreateModuleNameFromPathFail(t *testing.T) {
	n, err := CreateModuleNameFromPath(filepath.Join("work", "src", "github.com", "other", "project", "foo", "bar"))
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if n != "" {
		t.Fatalf("expected empty module name, got %s", n)
	}
}

func TestIsValidModuleVersion(t *testing.T) {
	if !IsValidModuleVersion("v10.21.23") {
		t.Fatal("unexpected invalid module version")
	}
	if IsValidModuleVersion("1.2.3") {
		t.Fatal("unexpected valid module version, missing v")
	}
	if IsValidModuleVersion("v11.563") {
		t.Fatal("unexpected valid module version, missing patch")
	}
}

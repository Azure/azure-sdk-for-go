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

package pkgs

import (
	"path/filepath"
	"strings"
	"testing"
)

const (
	testRoot = "../../testpkgs"
)

var (
	expected = map[string]string{
		"/scenrioa/foo":                    "foo",
		"/scenriob/foo":                    "foo",
		"/scenriob/foo/v2":                 "foo",
		"/scenrioc/mgmt/2019-10-11/foo":    "foo",
		"/scenriod/mgmt/2019-10-11/foo":    "foo",
		"/scenriod/mgmt/2019-10-11/foo/v2": "foo",
		"/scenrioe/mgmt/2019-10-11/foo":    "foo",
		"/scenrioe/mgmt/2019-10-11/foo/v2": "foo",
		"/scenrioe/mgmt/2019-10-11/foo/v3": "foo",
	}
)

func TestGetPkgs(t *testing.T) {
	root, err := filepath.Abs(testRoot)
	if err != nil {
		t.Fatalf("failed to get absolute path: %+v", err)
	}
	pkgs, err := GetPkgs(root)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}
	if len(pkgs) != len(expected) {
		t.Fatalf("expected %d packages, but got %d", len(expected), len(pkgs))
	}
	for _, pkg := range pkgs {
		if pkgName, ok := expected[pkg.Dest]; !ok {
			t.Fatalf("got pkg path '%s', but not found in expected", pkg.Dest)
		} else if !strings.EqualFold(pkgName, pkg.Package.Name) {
			t.Fatalf("expected package of '%s' in path '%s', but got '%s'", pkgName, pkg.Dest, pkg.Package.Name)
		}
	}
}

func TestPkg_GetApiVersion(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "/netapp/mgmt/2019-07-01/netapp",
			expected: "2019-07-01",
		},
		{
			input:    "/preview/resources/mgmt/2017-08-31-preview/managementgroups",
			expected: "2017-08-31-preview",
		},
		{
			input:    "/batch/2018-12-01.8.0/batch",
			expected: "2018-12-01.8.0",
		},
		{
			input:    "/cognitiveservices/v1.0/imagesearch",
			expected: "v1.0",
		},
		{
			input:    "/servicefabric/6.3/servicefabric",
			expected: "6.3",
		},
		{
			input:    "/keyvault/2015-06-01/keyvault",
			expected: "2015-06-01",
		},
		{
			input:    "/preview/datalake/analytics/2017-09-01-preview/job",
			expected: "2017-09-01-preview",
		},
	}
	for _, c := range cases {
		p := Pkg{Dest: c.input}
		api, err := p.GetAPIVersion()
		if err != nil {
			t.Fatalf("failed to get api version: %+v", err)
		}
		if !strings.EqualFold(api, c.expected) {
			t.Fatalf("expected: %s, but got %s", c.expected, api)
		}
	}
}

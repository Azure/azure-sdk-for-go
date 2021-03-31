// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

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

func Test_getPkgs(t *testing.T) {
	rootDir, err := filepath.Abs(testRoot)
	if err != nil {
		t.Fatalf("failed to get absolute path: %+v", err)
	}
	pkgs, err := getPkgs(rootDir)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}
	if len(pkgs) != len(expected) {
		t.Fatalf("expected %d packages, but got %d", len(expected), len(pkgs))
	}
	for _, pkg := range pkgs {
		if pkgName, ok := expected[pkg.d]; !ok {
			t.Fatalf("got pkg path '%s', but not found in expected", pkg.d)
		} else if !strings.EqualFold(pkgName, pkg.p.Name) {
			t.Fatalf("expected package of '%s' in path '%s', but got '%s'", pkgName, pkg.d, pkg.p.Name)
		}
	}
}

func Test_verifyDirectoryStructure(t *testing.T) {
	rootDir, err := filepath.Abs(testRoot)
	if err != nil {
		t.Fatalf("failed to get absolute path: %+v", err)
	}
	pkgs, err := getPkgs(rootDir)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}
	for _, pkg := range pkgs {
		if err := verifyDirectoryStructure(pkg); err != nil {
			t.Fatalf("failed to verify directory structure: %+v", err)
		}
	}
}

func Test_verifyPkgMatchesDir(t *testing.T) {
	rootDir, err := filepath.Abs(testRoot)
	if err != nil {
		t.Fatalf("failed to get absolute path: %+v", err)
	}
	pkgs, err := getPkgs(rootDir)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}
	for _, pkg := range pkgs {
		if err := verifyPkgMatchesDir(pkg); err != nil {
			t.Fatalf("failed to verify package directory name: %+v", err)
		}
	}
}

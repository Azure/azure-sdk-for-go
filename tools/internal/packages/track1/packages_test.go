// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package track1

import (
	"path/filepath"
	"strings"
	"testing"
)

const (
	testRoot = "./testpkgs"
)

var (
	expected = map[string]string{
		"scenrioa/foo":                 "foo",
		"scenriob/mgmt/2019-10-11/foo": "foo",
	}
)

func TestList(t *testing.T) {
	root, err := filepath.Abs(testRoot)
	if err != nil {
		t.Fatalf("failed to get absolute path: %+v", err)
	}
	pkgs, err := List(root)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}
	if len(pkgs) != len(expected) {
		t.Fatalf("expected %d packages, but got %d", len(expected), len(pkgs))
	}
	for _, pkg := range pkgs {
		if pkgName, ok := expected[pkg.Path()]; !ok {
			t.Fatalf("got pkg path '%s', but not found in expected", pkg.Path())
		} else if !strings.EqualFold(pkgName, pkg.Name()) {
			t.Fatalf("expected package of '%s' in path '%s', but got '%s'", pkgName, pkg.Path(), pkg.Name())
		}
	}
}

func TestVerifier_Verify(t *testing.T) {
	root, err := filepath.Abs(testRoot)
	if err != nil {
		t.Fatalf("failed to get absolute path: %+v", err)
	}
	pkgs, err := List(root)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}

	verifier := GetDefaultVerifier()
	for _, pkg := range pkgs {
		if errors := verifier.Verify(pkg); len(errors) != 0 {
			t.Fatalf("failed to verify packages: %+v", errors)
		}
	}
}

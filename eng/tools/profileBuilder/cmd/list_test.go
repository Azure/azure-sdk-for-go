// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder/model"
)

func Test_updateModuleVersions(t *testing.T) {
	ld := model.ListDefinition{
		Include: []string{
			"../../testdata/scenarioa/foo",
			"../../testdata/scenariod/foo",
			"../../testdata/scenarioe/foo/v2",
		},
	}
	updateModuleVersions(&ld)
	expected := []string{
		"../../testdata/scenarioa/foo",
		"../../testdata/scenariod/foo/v2",
		"../../testdata/scenarioe/foo/v2",
	}
	if !reflect.DeepEqual(ld.Include, expected) {
		t.Fatalf("expected '%v' got '%v'", expected, ld.Include)
	}
}

func Test_getLatestModVer(t *testing.T) {
	d, err := getLatestModVer("../../testdata/scenarioa/foo")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if d != "" {
		t.Fatalf("expected empty string got '%s'", d)
	}
	d, err = getLatestModVer("../../testdata/scenariod/foo")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if d != "v2" {
		t.Fatalf("expected 'v2' string got '%s'", d)
	}
}

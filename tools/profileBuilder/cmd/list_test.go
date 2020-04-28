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
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/profileBuilder/model"
)

func Test_updateModuleVersions(t *testing.T) {
	// current module implementation is without major sub-directories
	ld := model.ListDefinition{
		Include: []string{
			"../../testdata/scenarioa/foo",
			"../../testdata/scenariod/foo",
			"../../testdata/scenarioe/foo",
		},
	}
	updateModuleVersions(&ld)
	expected := []string{
		"../../testdata/scenarioa/foo",
		"../../testdata/scenariod/foo",
		"../../testdata/scenarioe/foo",
	}
	if !reflect.DeepEqual(ld.Include, expected) {
		t.Fatalf("expected '%v' got '%v'", expected, ld.Include)
	}
}

func Test_getLatestModVer(t *testing.T) {
	// current module implementation is without major sub-directories, therefore we cannot get the module version just by the directory structure here
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
	if d != "" {
		t.Fatalf("expected empty string got '%s'", d)
	}
}

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
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/internal/pkgs"
)

const (
	testRoot = "../../testpkgs"
)

func Test_verifyDirectoryStructure(t *testing.T) {
	rootDir, err := filepath.Abs(testRoot)
	if err != nil {
		t.Fatalf("failed to get absolute path: %+v", err)
	}
	ps, err := pkgs.GetPkgs(rootDir)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}
	for _, pkg := range ps {
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
	ps, err := pkgs.GetPkgs(rootDir)
	if err != nil {
		t.Fatalf("failed to get packages: %+v", err)
	}
	for _, pkg := range ps {
		if err := verifyPkgMatchesDir(pkg); err != nil {
			t.Fatalf("failed to verify package directory name: %+v", err)
		}
	}
}

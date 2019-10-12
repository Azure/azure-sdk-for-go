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
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func cleanTestData() {
	cmd := exec.Command("git", "clean", "-xdf", "../../testdata")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(string(output))
	}
	cmd = exec.Command("git", "checkout", "--", "../../testdata")
	output, err = cmd.CombinedOutput()
	if err != nil {
		panic(string(output))
	}
}

func Test_findAllSubDirectories(t *testing.T) {
	cleanTestData()
	defer cleanTestData()
	expected := []string{
		"../../testdata/scenarioa/foo/stage",
		"../../testdata/scenariob/foo/stage",
		"../../testdata/scenarioc/foo/stage",
		"../../testdata/scenariod/foo/stage",
		"../../testdata/scenarioe/foo/stage",
		"../../testdata/scenariof/foo/stage",
		"../../testdata/scenariog/foo/mgmt/2019-10-11/foo/stage",
		"../../testdata/scenarioh/foo/mgmt/2019-10-11/foo/stage",
	}
	root, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatalf("error when get absolute path of root: %+v", err)
	}
	stages, err := findAllSubDirectories(root, "stage")
	if err != nil {
		t.Fatalf("error when listing all stage folders: %+v", err)
	}
	if len(stages) != len(expected) {
		t.Fatalf("expected %d stages folders, but got %d", len(expected), len(stages))
	}
	for i, stage := range stages {
		e, err := filepath.Abs(expected[i])
		if err != nil {
			t.Fatalf("error when parsing expected results '%s'(%d)", expected[i], i)
		}
		if !strings.EqualFold(stage, e) {
			t.Fatalf("expected folder '%s', but got '%s'", e, stage)
		}
	}
}

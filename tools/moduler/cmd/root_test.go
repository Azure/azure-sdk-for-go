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
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/versioner/cmd"
)

var (
	tagsFullList = []string{
		"tools/testdata/scenarioa/foo/v1.0.0",
		"tools/testdata/scenariob/foo/v1.0.0",
		"tools/testdata/scenariob/foo/v1.1.0",
		"tools/testdata/scenarioc/foo/v1.0.0",
		"tools/testdata/scenariod/foo/v1.0.0",
		"tools/testdata/scenariod/foo/v1.0.1",
		"tools/testdata/scenariod/foo/v1.1.0",
		"tools/testdata/scenariod/foo/v1.2.0",
		"tools/testdata/scenariod/foo/v2.0.0",
		"tools/testdata/scenariod/foo/v2.1.0",
		"tools/testdata/scenariod/foo/v2.1.1",
		"tools/testdata/scenarioe/foo/v1.0.0",
		"tools/testdata/scenarioe/foo/v1.1.0",
		"tools/testdata/scenarioe/foo/v2.0.0",
		"tools/testdata/scenarioe/foo/v2.1.0",
	}
)

func prepareTestData(root string, tagsHook cmd.TagsHookFunc) error {
	if _, err := cmd.ExecuteVersioner(root, tagsHook); err != nil {
		return fmt.Errorf("failed to prepare test data: %+v", err)
	}
	return nil
}

func cleanTestData() {
	c := exec.Command("git", "clean", "-xdf", "../../testdata")
	if output, err := c.CombinedOutput(); err != nil {
		panic(string(output))
	}
	c = exec.Command("git", "checkout", "--", "../../testdata")
	if output, err := c.CombinedOutput(); err != nil {
		panic(string(output))
	}
}

func Test_findAllVersionFiles(t *testing.T) {
	cleanTestData()
	if err := prepareTestData("../../testdata", func(repoPath string, tagPrefix string) ([]string, error) {
		return tagsFullList, nil
	}); err != nil {
		t.Fatalf("failed to prepare test data: %+v", err)
	}
	defer cleanTestData()
	expected := []string{
		"../../testdata/scenarioa/foo/version.go",
		"../../testdata/scenariob/foo/version.go",
		"../../testdata/scenariob/foo/v2/version.go",
		"../../testdata/scenarioc/foo/version.go",
		"../../testdata/scenariod/foo/version.go",
		"../../testdata/scenariod/foo/v2/version.go",
		"../../testdata/scenariod/foo/v3/version.go",
		"../../testdata/scenarioe/foo/version.go",
		"../../testdata/scenarioe/foo/v2/version.go",
		"../../testdata/scenariof/foo/version.go",
		"../../testdata/scenariog/foo/mgmt/2019-10-11/foo/version.go",
		"../../testdata/scenariog/foo/mgmt/2019-10-11/foo/v2/version.go",
		"../../testdata/scenarioh/foo/mgmt/2019-10-11/foo/version.go",
		"../../testdata/scenarioh/foo/mgmt/2019-10-11/foo/v2/version.go",
	}
	root, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatalf("failed to get absolute path of root: %+v", err)
	}
	files, err := findAllFiles(root, "version.go")
	if len(expected) != len(files) {
		t.Fatalf("expected %d version files, but got %d", len(expected), len(files))
	}
	fileSet, err := makeSet(files)
	if err != nil {
		t.Fatalf("failed to make file set of version files: %+v", err)
	}
	for _, e := range expected {
		absE, err := filepath.Abs(e)
		if err != nil {
			t.Fatalf("failed to get absolute path of expected file '%s': %+v", e, err)
		}
		if _, ok := fileSet[absE]; !ok {
			t.Fatalf("expected folder '%s' not found", absE)
		}
	}
}

func makeSet(array []string) (map[string]bool, error) {
	result := make(map[string]bool, len(array))
	for _, item := range array {
		if _, ok := result[item]; ok {
			return nil, fmt.Errorf("failed to make set, contains duplicate items '%s'", item)
		}
		result[item] = true
	}
	return result, nil
}

func Test_readNewTagInFile(t *testing.T) {
	cleanTestData()
	if err := prepareTestData("../../testdata/scenarioa", func(repoPath string, tagPrefix string) ([]string, error) {
		return []string{
			"tools/testdata/scenarioa/foo/v1.0.0",
		}, nil
	}); err != nil {
		t.Fatalf("failed to prepare test data: %+v", err)
	}
	defer cleanTestData()
	expected := "tools/testdata/scenarioa/foo/v1.1.0"
	p, err := filepath.Abs("../../testdata/scenarioa/foo/version.go")
	if err != nil {
		t.Fatalf("failed to get absolute path of version file: %+v", err)
	}
	tag, err := readNewTagInFile(p)
	if !strings.EqualFold(expected, tag) {
		t.Fatalf("expected tag '%s', but got '%s'", expected, tag)
	}
}

func Test_readNewTags(t *testing.T) {
	cleanTestData()
	if err := prepareTestData("../../testdata", func(repoPath string, tagPrefix string) ([]string, error) {
		tags := make([]string, 0)
		for _, tag := range tagsFullList {
			if strings.HasPrefix(tag, tagPrefix) {
				tags = append(tags, tag)
			}
		}
		return tags, nil
	}); err != nil {
		t.Fatalf("failed to prepare test data: %+v", err)
	}
	defer cleanTestData()
	expected := []string{
		"tools/testdata/scenarioa/foo/v1.1.0",
		"tools/testdata/scenariob/foo/v2.0.0",
		"tools/testdata/scenarioc/foo/v1.0.1",
		"tools/testdata/scenariod/foo/v3.0.0",
		"tools/testdata/scenarioe/foo/v2.2.0",
		"tools/testdata/scenariof/foo/v1.0.0",
		"tools/testdata/scenariog/foo/mgmt/2019-10-11/foo/v2.0.0",
		"tools/testdata/scenarioh/foo/mgmt/2019-10-11/foo/v2.0.0",
	}
	root, err := filepath.Abs("../../testdata")
	if err != nil {
		t.Fatalf("failed to get absolute path of root: %+v", err)
	}
	newTags, err := readNewTags(root)
	if err != nil {
		t.Fatalf("failed to read new tags: %+v", err)
	}
	if len(expected) != len(newTags) {
		t.Fatalf("expected %d new tags, but got %d", len(expected), len(newTags))
	}
	tagSet, err := makeSet(newTags)
	if err != nil {
		t.Fatalf("failed to make file set of tags: %+v", err)
	}
	for _, tag := range expected {
		if _, ok := tagSet[tag]; !ok {
			t.Fatalf("expected tag '%s' not found", tag)
		}
	}
}

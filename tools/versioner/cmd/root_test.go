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
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/versioner/internal/modinfo"
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

func Test_sortModuleTagsBySemver(t *testing.T) {
	before := []string{
		"foo/v1.0.0",
		"foo/v1.0.1",
		"foo/v1.1.0",
		"foo/v10.0.0",
		"foo/v11.1.1",
		"foo/v2.0.0",
		"foo/v20.2.3",
		"foo/v3.1.0",
	}
	sortModuleTagsBySemver(before)
	after := []string{
		"foo/v1.0.0",
		"foo/v1.0.1",
		"foo/v1.1.0",
		"foo/v2.0.0",
		"foo/v3.1.0",
		"foo/v10.0.0",
		"foo/v11.1.1",
		"foo/v20.2.3",
	}
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("sort order doesn't match, expected '%v' got '%v'", after, before)
	}
}

func Test_getTagPrefix(t *testing.T) {
	p, err := getTagPrefix("/work/src/github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis")
	if err != nil {
		t.Fatal("failed to get tag prefix")
	}
	if p != "services/redis/mgmt/2018-03-01/redis" {
		t.Fatalf("wrong value '%s' for tag prefix", p)
	}
	p, err = getTagPrefix("/work/src/github.com/something/else")
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if p != "" {
		t.Fatalf("unexpected tag '%s'", p)
	}
}

type mockModInfo struct {
	dir     string
	exports bool
	breaks  bool
}

func (mock mockModInfo) DestDir() string {
	return mock.dir
}

func (mock mockModInfo) NewExports() bool {
	return mock.exports
}

func (mock mockModInfo) BreakingChanges() bool {
	return mock.breaks
}

func (mock mockModInfo) VersionSuffix() bool {
	return modinfo.HasVersionSuffix(mock.dir)
}

func Test_calculateModuleTagMajorV1(t *testing.T) {
	pkg := mockModInfo{
		dir: "/work/src/github.com/Azure/azure-sdk-for-go/services/foo",
	}
	tag, err := calculateModuleTag([]string{}, pkg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag != "services/foo/v1.0.0" {
		t.Fatalf("bad tag '%s", tag)
	}
}

func Test_calculateModuleTagMajorV2(t *testing.T) {
	tags := []string{
		"services/foo/v1.0.0",
		"services/foo/v1.1.0",
	}
	pkg := mockModInfo{
		dir:    "/work/src/github.com/Azure/azure-sdk-for-go/services/foo/v2",
		breaks: true,
	}
	tag, err := calculateModuleTag(tags, pkg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag != "services/foo/v2.0.0" {
		t.Fatalf("bad tag '%s", tag)
	}
}

func Test_calculateModuleTagMajorV2Invalid(t *testing.T) {
	tags := []string{
		"services/foo/v1.0.0",
		"services/foo/v1.1.0",
	}
	pkg := mockModInfo{
		dir:    "/work/src/github.com/Azure/azure-sdk-for-go/services/foo", // missing /v2 suffix
		breaks: true,
	}
	tag, err := calculateModuleTag(tags, pkg)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if tag != "" {
		t.Fatal("expected no tag")
	}
}

func Test_calculateModuleTagMinorV1(t *testing.T) {
	tags := []string{
		"services/foo/v1.0.0",
		"services/foo/v1.0.1",
	}
	pkg := mockModInfo{
		dir:     "/work/src/github.com/Azure/azure-sdk-for-go/services/foo",
		exports: true,
	}
	tag, err := calculateModuleTag(tags, pkg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag != "services/foo/v1.1.0" {
		t.Fatalf("bad tag '%s", tag)
	}
}

func Test_calculateModuleTagMinorV2(t *testing.T) {
	tags := []string{
		"services/foo/v1.0.0",
		"services/foo/v1.0.1",
		"services/foo/v2.0.0",
	}
	pkg := mockModInfo{
		dir: "/work/src/github.com/Azure/azure-sdk-for-go/services/foo/v2",
	}
	tag, err := calculateModuleTag(tags, pkg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag != "services/foo/v2.0.1" {
		t.Fatalf("bad tag '%s", tag)
	}
}

func Test_findLatestMajorVersion(t *testing.T) {
	ver, err := findLatestMajorVersion("./testdata/moda/stage")
	if err != nil {
		t.Fatalf("failed to find LMV: %v", err)
	}
	if ver != "testdata/moda" {
		t.Fatalf("bad LMV %s", ver)
	}
	ver, err = findLatestMajorVersion("./testdata/modb/stage")
	if err != nil {
		t.Fatalf("failed to find LMV: %v", err)
	}
	if ver != "testdata/modb/v2" {
		t.Fatalf("bad LMV %s", ver)
	}
}

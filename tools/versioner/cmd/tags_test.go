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
	"path/filepath"
	"testing"
)

func Test_getTags(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		// travis does a shallow clone so tag count is not consistent
		t.SkipNow()
	}
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

func Test_getTagPrefix(t *testing.T) {
	p, err := getTagPrefix(filepath.Join("work", "src", "github.com", "Azure", "azure-sdk-for-go", "services", "redis", "mgmt", "2018-03-01", "redis"))
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

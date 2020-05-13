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

package dirs

import (
	"github.com/Azure/azure-sdk-for-go/tools/apidiff/ioext"
	"os"
	"path/filepath"
	"testing"
)

func TestGetSubdirs(t *testing.T) {
	sd, err := GetSubdirs("../..")
	if err != nil {
		t.Fatalf("failed to get subdirs: %v", err)
	}
	if len(sd) == 0 {
		t.Fatal("unexpected zero length subdirs")
	}
}

func TestGetSubdirsEmpty(t *testing.T) {
	sd, err := GetSubdirs(".")
	if err != nil {
		t.Fatalf("failed to get subdirs: %v", err)
	}
	if len(sd) != 0 {
		t.Fatal("expected zero length subdirs")
	}
}

func TestGetSubdirsNoExist(t *testing.T) {
	sd, err := GetSubdirs("../thisdoesntexist")
	if err == nil {
		t.Fatal("expected nil error")
	}
	if sd != nil {
		t.Fatal("expected nil subdirs")
	}
}

func TestDeepCompare(t *testing.T) {
	profiles, err := filepath.Abs("../../../profiles")
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	temp := profiles + "1temp"
	defer func() {
		err := os.RemoveAll(temp)
		if err != nil {
			t.Fatalf("cannot remove the temp directory: %+v", err)
		}
	}()
	if err := ioext.CopyDir(profiles, temp); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	identical, err := DeepCompare(profiles, temp)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if !identical {
		t.Fatalf("expected %v but got %v", true, identical)
	}
}

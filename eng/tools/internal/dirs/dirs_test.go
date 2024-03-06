// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package dirs

import "testing"

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

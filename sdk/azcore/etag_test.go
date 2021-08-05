// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import "testing"

func TestEtagEquals(t *testing.T) {
	e1 := NewEtag("tag")
	if e1.String() != "tag" {
		t.Fatalf("ETag values are not equal")
	}

	e2 := NewEtag("\"tag\"")
	if e2.String() != "\"tag\"" {
		t.Fatalf("ETag values are not equal")
	}

	e3 := NewEtag("W/\"weakETag\"")
	if e3.String() != "W/\"weakETag\"" {
		t.Fatalf("ETag values are not equal")
	}
	if !e3.IsWeak() {
		t.Fatalf("ETag is expected to be weak")
	}

	strongETag := NewEtag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	if strongETag.String() != "\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"" {
		t.Fatalf("Etag values are not equal")
	}
}

func TestEtagWeak(t *testing.T) {
	
}
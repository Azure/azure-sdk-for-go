//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import "testing"

func TestKindFromToken(t *testing.T) {
	const tk = `{ "type": "pollerID;kind" }`
	k, err := KindFromToken("pollerID", tk)
	if err != nil {
		t.Fatal(err)
	}
	if k != "kind" {
		t.Fatalf("unexpected kind %s", k)
	}
	k, err = KindFromToken("mismatched", tk)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
}

func TestKindFromTokenInvalid(t *testing.T) {
	const tk1 = `{ "missing": "type" }`
	k, err := KindFromToken("mismatched", tk1)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
	const tk2 = `{ "type": false }`
	k, err = KindFromToken("mismatched", tk2)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
	const tk3 = `{ "type": "pollerID;kind;extra" }`
	k, err = KindFromToken("mismatched", tk3)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
}

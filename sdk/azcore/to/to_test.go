//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package to

import (
	"testing"
)

func TestPtr(t *testing.T) {
	b := true
	pb := Ptr(b)
	if pb == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pb != b {
		t.Fatalf("got %v, want %v", *pb, b)
	}
}

func TestSliceOfPtrs(t *testing.T) {
	arr := SliceOfPtrs[int]()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = SliceOfPtrs(1, 2, 3, 4, 5)
	for i, v := range arr {
		if *v != i+1 {
			t.Fatal("values don't match")
		}
	}
}

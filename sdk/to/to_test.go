// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package to

import "testing"

func TestBoolPtr(t *testing.T) {
	b := true
	pb := BoolPtr(b)
	if pb == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pb != b {
		t.Fatalf("got %v, want %v", *pb, b)
	}
}

func TestFloat32Ptr(t *testing.T) {
	f32 := float32(3.1415926)
	pf32 := Float32Ptr(f32)
	if pf32 == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pf32 != f32 {
		t.Fatalf("got %v, want %v", *pf32, f32)
	}
}

func TestFloat64Ptr(t *testing.T) {
	f64 := float64(2.71828182845904)
	pf64 := Float64Ptr(f64)
	if pf64 == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pf64 != f64 {
		t.Fatalf("got %v, want %v", *pf64, f64)
	}
}

func TestInt32Ptr(t *testing.T) {
	i32 := int32(123456789)
	pi32 := Int32Ptr(i32)
	if pi32 == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pi32 != i32 {
		t.Fatalf("got %v, want %v", *pi32, i32)
	}
}

func TestInt64Ptr(t *testing.T) {
	i64 := int64(9876543210)
	pi64 := Int64Ptr(i64)
	if pi64 == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pi64 != i64 {
		t.Fatalf("got %v, want %v", *pi64, i64)
	}
}

func TestStringPtr(t *testing.T) {
	s := "the string"
	ps := StringPtr(s)
	if ps == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *ps != s {
		t.Fatalf("got %v, want %v", *ps, s)
	}
}

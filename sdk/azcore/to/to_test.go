//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package to

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

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

func TestTimePtr(t *testing.T) {
	tt := time.Now()
	pt := TimePtr(tt)
	if pt == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pt != tt {
		t.Fatalf("got %v, want %v", *pt, tt)
	}
}

func TestInt32PtrArray(t *testing.T) {
	arr := Int32PtrArray()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = Int32PtrArray(1, 2, 3, 4, 5)
	for i, v := range arr {
		if *v != int32(i+1) {
			t.Fatal("values don't match")
		}
	}
}

func TestInt64PtrArray(t *testing.T) {
	arr := Int64PtrArray()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = Int64PtrArray(1, 2, 3, 4, 5)
	for i, v := range arr {
		if *v != int64(i+1) {
			t.Fatal("values don't match")
		}
	}
}

func TestFloat32PtrArray(t *testing.T) {
	arr := Float32PtrArray()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = Float32PtrArray(1.1, 2.2, 3.3, 4.4, 5.5)
	for i, v := range arr {
		f, err := strconv.ParseFloat(fmt.Sprintf("%d.%d", i+1, i+1), 32)
		if err != nil {
			t.Fatal(err)
		}
		if *v != float32(f) {
			t.Fatal("values don't match")
		}
	}
}

func TestFloat64PtrArray(t *testing.T) {
	arr := Float64PtrArray()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = Float64PtrArray(1.1, 2.2, 3.3, 4.4, 5.5)
	for i, v := range arr {
		f, err := strconv.ParseFloat(fmt.Sprintf("%d.%d", i+1, i+1), 64)
		if err != nil {
			t.Fatal(err)
		}
		if *v != f {
			t.Fatal("values don't match")
		}
	}
}

func TestBoolPtrArray(t *testing.T) {
	arr := BoolPtrArray()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = BoolPtrArray(true, false, true)
	curr := true
	for _, v := range arr {
		if *v != curr {
			t.Fatal("values don'p match")
		}
		curr = !curr
	}
}

func TestStringPtrArray(t *testing.T) {
	arr := StringPtrArray()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = StringPtrArray("one", "", "three")
	if !reflect.DeepEqual(arr, []*string{StringPtr("one"), StringPtr(""), StringPtr("three")}) {
		t.Fatal("values don't match")
	}
}

func TestTimePtrArray(t *testing.T) {
	arr := TimePtrArray()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	t1 := time.Now()
	t2 := time.Time{}
	t3 := t1.Add(24 * time.Hour)
	arr = TimePtrArray(t1, t2, t3)
	if !reflect.DeepEqual(arr, []*time.Time{&t1, &t2, &t3}) {
		t.Fatal("values don't match")
	}
}

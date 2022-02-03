//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"reflect"
	"testing"
)

func TestNullValue(t *testing.T) {
	v := NullValue[*string]()
	vv := NullValue[*string]()
	if v != vv {
		t.Fatal("null values should match for the same types")
	}
}

func TestIsNullValue(t *testing.T) {
	if IsNullValue("") {
		t.Fatal("string literal can't be a null value")
	}
	s := ""
	if IsNullValue(&s) {
		t.Fatal("&s isn't a null value")
	}
	var i *int
	if IsNullValue(i) {
		t.Fatal("i isn't a null value")
	}
	i = NullValue[*int]()
	if !IsNullValue(i) {
		t.Fatal("expected null value for i")
	}
	i2 := 0
	i = &i2
	if IsNullValue(i) {
		t.Fatal("i should no longer be null value")
	}
}

func TestNullValueMapSlice(t *testing.T) {
	v := NullValue[[]string]()
	vv := NullValue[[]string]()
	if reflect.TypeOf(v) != reflect.TypeOf(vv) {
		t.Fatal("null values should match for the same types")
	}
	m := NullValue[map[string]int]()
	if reflect.TypeOf(v) == reflect.TypeOf(m) {
		t.Fatal("null values for string and int should not match")
	}
}

func TestIsNullValueMapSlice(t *testing.T) {
	if IsNullValue([]string{}) {
		t.Fatal("slice literal can't be a null value")
	}
	if IsNullValue(map[int]string{}) {
		t.Fatal("map literal can't be a null value")
	}
	s := NullValue[[]int]()
	if !IsNullValue(s) {
		t.Fatal("expected null value for s")
	}
	m := NullValue[map[string]interface{}]()
	if !IsNullValue(m) {
		t.Fatal("expected null value for s")
	}

	type nullFields struct {
		Map   map[string]int
		Slice []string
	}

	nf := nullFields{}
	if IsNullValue(nf.Map) {
		t.Fatal("unexpected null map")
	}
	if IsNullValue(nf.Slice) {
		t.Fatal("unexpected null slice")
	}

	nf.Map = map[string]int{}
	nf.Slice = []string{}
	if IsNullValue(nf.Map) {
		t.Fatal("unexpected null map")
	}
	if IsNullValue(nf.Slice) {
		t.Fatal("unexpected null slice")
	}

	nf.Map = NullValue[map[string]int]()
	nf.Slice = NullValue[[]string]()
	if !IsNullValue(nf.Map) {
		t.Fatal("expected null map")
	}
	if !IsNullValue(nf.Slice) {
		t.Fatal("expected null slice")
	}
}

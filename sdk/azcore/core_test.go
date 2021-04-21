// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import "testing"

func TestNullValue(t *testing.T) {
	v := NullValue("")
	if _, ok := v.(*string); !ok {
		t.Fatalf("unexpected type %T", v)
	}
	vv := NullValue((*string)(nil))
	if _, ok := vv.(*string); !ok {
		t.Fatalf("unexpected type %T", vv)
	}
	if v != vv {
		t.Fatal("null values should match for the same types")
	}
	i := NullValue(1)
	if _, ok := i.(*int); !ok {
		t.Fatalf("unexpected type %T", v)
	}
	if v == i {
		t.Fatal("null values for string and int should not match")
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
	i = NullValue(0).(*int)
	if !IsNullValue(i) {
		t.Fatal("expected null value for i")
	}
	i2 := 0
	i = &i2
	if IsNullValue(i) {
		t.Fatal("i should no longer be null value")
	}
}

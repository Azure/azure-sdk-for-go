// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
	"testing"
)

func TestEtagEquals(t *testing.T) {
	e1 := NewETag("tag")
	if e1.String() != "tag" {
		t.Fatalf("ETag values are not equal")
	}

	e2 := NewETag("\"tag\"")
	if e2.String() != "\"tag\"" {
		t.Fatalf("ETag values are not equal")
	}

	e3 := NewETag("W/\"weakETag\"")
	if e3.String() != "W/\"weakETag\"" {
		t.Fatalf("ETag values are not equal")
	}
	if !e3.IsWeak() {
		t.Fatalf("ETag is expected to be weak")
	}

	strongETag := NewETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	if strongETag.String() != "\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"" {
		t.Fatalf("Etag values are not equal")
	}

	if ETagAny().IsWeak() {
		t.Fatalf("ETagAny should not be weak")
	}
}

func TestETagWeak(t *testing.T) {
	et1 := NewETag("tag")
	if !et1.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et2 := NewETag("\"tag\"")
	if !et2.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et3 := NewETag("W/\"weakETag\"")
	if !et3.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et4 := NewETag("W/\"\"")
	if !et4.IsWeak() {
		t.Fatalf("Expected to be weak")
	}

	et5 := ETagAny()
	if !et5.IsWeak() {
		t.Fatalf("Expected to be weak")
	}
}

/*

func TestEtagStars(t *testing.T) {
	anyETag := ETagAny()
	star := NewETag("*")
	// weakStar := NewETag("W\"*\"")
	quotedStar := NewETag("\"*\"")

	strongETag := NewETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")

	if anyETag.Equals(*anyETag, Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if anyETag.Equals(*ETagAny(), Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if anyETag.Equals(*strongETag, Strong) {
		t.Fatalf("Expected etags to be equal")
	}

	// expectEqual(t, star, star)
	// expectEqual(t, star, ETagAny())
	// expectEqual(t, star, anyETag)

	// expectNotEqual(t, star, weakStar)
	// expectNotEqual(t, weakStar, ETagAny())
	// expectNotEqual(t, quotedStar, weakStar)

	expectNotEqual(t, star, quotedStar)
	expectEqual(t, anyETag, star)
}
*/

func expectEqual(t *testing.T, left *ETag, right *ETag) {
	if !left.Equals(*right, Strong) {
		fmt.Println(*left.value, *right.value)
		t.Fatalf("Expected etags to be equal")
	}
}

func expectNotEqual(t *testing.T, left *ETag, right *ETag) {
	if left.Equals(*right, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
}

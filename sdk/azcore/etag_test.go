// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
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

func TestEtagEquality(t *testing.T) {
	weakTag := NewETag("W/\"\"")
	weakTag1 := NewETag("W/\"1\"")
	weakTag2 := NewETag("W/\"Two\"")
	strongTag1 := NewETag("\"1\"")
	strongTag2 := NewETag("\"Two\"")
	strongTagValidChars := NewETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	weakTagValidChars := NewETag("W/\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")

	if weakTag.Equals(*weakTag, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag1.Equals(*weakTag1, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag2.Equals(*weakTag2, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTagValidChars.Equals(*weakTagValidChars, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if !strongTag1.Equals(*strongTag1, Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if !strongTag2.Equals(*strongTag2, Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if !strongTagValidChars.Equals(*strongTagValidChars, Strong) {
		t.Fatalf("Expected etags to be equal")
	}

	if weakTag.Equals(*weakTag1, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTagValidChars.Equals(*strongTagValidChars, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}

	if weakTag1.Equals(*weakTag2, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag1.Equals(*strongTag1, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag2.Equals(*strongTag2, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
}

func TestEtagAny(t *testing.T) {
	anyETag := ETagAny()
	star := NewETag("*")
	weakStar := NewETag("W\"*\"")
	quotedStart := NewETag("\"*\"")

	if !anyETag.Equals(*anyETag, Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if !anyETag.Equals(*ETagAny(), Strong) {
		t.Fatalf("Expected etags to be equal")
	}

	if !star.Equals(*star, Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if !star.Equals(*ETagAny(), Strong) {
		t.Fatalf("Expected etags to be equal")
	}
	if !star.Equals(*anyETag, Strong) {
		t.Fatalf("Expected etags to be equal")
	}

	if star.Equals(*weakStar, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakStar.Equals(*ETagAny(), Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if quotedStart.Equals(*weakStar, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}


	if star.Equals(*quotedStart, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
	if !ETagAny().Equals(*star, Strong) {
		t.Fatalf("Expected etags to not be equal")
	}
}

func TestEtagWeakComparison(t *testing.T) {
	// W/""
	weakTag := NewETag("W/\"\"");
	// W/"1"
	weakTag1 := NewETag("W/\"1\"");
	// W/"Two"
	weakTagTwo := NewETag("W/\"Two\"");
	// W/"two"
	weakTagtwo := NewETag("W/\"two\"");
	// "1"
	strongTag1 := NewETag("\"1\"");
	// "Two"
	strongTagTwo := NewETag("\"Two\"");
	// "two"
	strongTagtwo := NewETag("\"two\"");

	if !weakTag.Equals(*weakTag, Weak) {
		t.Fatalf("Expected etags to be equal")
	}
	if !weakTag1.Equals(*weakTag1, Weak) {
		t.Fatalf("Expected etags to be equal")
	}
	if !weakTagTwo.Equals(*weakTagTwo, Weak) {
		t.Fatalf("Expected etags to be equal")
	}
	if !weakTagtwo.Equals(*weakTagtwo, Weak) {
		t.Fatalf("Expected etags to be equal")
	}
	if !strongTag1.Equals(*strongTag1, Weak) {
		t.Fatalf("Expected etags to be equal")
	}
	if !strongTagTwo.Equals(*strongTagTwo, Weak) {
		t.Fatalf("Expected etags to be equal")
	}
	if !strongTagtwo.Equals(*strongTagtwo, Weak) {
		t.Fatalf("Expected etags to be equal")
	}

	if weakTag.Equals(*weakTag1, Weak) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTag1.Equals(*weakTagTwo, Weak) {
		t.Fatalf("Expected etags to not be equal")
	}

	if !weakTag1.Equals(*strongTag1, Weak) {
		t.Fatalf("Expected etags to be equal")
	}
	if !weakTagTwo.Equals(*strongTagTwo, Weak) {
		t.Fatalf("Expected etags to be equal")
	}

	if strongTagTwo.Equals(*weakTag1, Weak) {
		t.Fatalf("Expected etags to not be equal")
	}
	if strongTagTwo.Equals(*weakTagtwo, Weak) {
		t.Fatalf("Expected etags to not be equal")
	}

	if strongTagTwo.Equals(*strongTagtwo, Weak) {
		t.Fatalf("Expected etags to not be equal")
	}
	if weakTagTwo.Equals(*weakTagtwo, Weak) {
		t.Fatalf("Expected etags to not be equal")
	}
}

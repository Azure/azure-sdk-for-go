//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createETag(s string) ETag {
	return ETag(s)
}

func TestETagEquals(t *testing.T) {
	e1 := createETag("tag")
	require.Equal(t, "tag", string(e1))

	e2 := createETag("\"tag\"")
	require.Equal(t, "\"tag\"", string(e2))

	e3 := createETag("W/\"weakETag\"")
	require.Equal(t, "W/\"weakETag\"", string(e3))
	require.Truef(t, e3.IsWeak(), "ETag is expected to be weak")

	strongETag := createETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	require.Equal(t, "\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"", string(strongETag))

	require.Falsef(t, ETagAny.IsWeak(), "ETagAny should not be weak")
}

func TestETagWeak(t *testing.T) {
	et1 := createETag("tag")
	require.Falsef(t, et1.IsWeak(), "expected etag to be strong")

	et2 := createETag("\"tag\"")
	require.Falsef(t, et2.IsWeak(), "expected etag to be strong")

	et3 := createETag("W/\"weakETag\"")
	require.Truef(t, et3.IsWeak(), "expected etag to be weak")

	et4 := createETag("W/\"\"")
	require.Truef(t, et4.IsWeak(), "expected etag to be weak")

	et5 := ETagAny
	require.Falsef(t, et5.IsWeak(), "expected etag to be strong")
}

func TestETagEquality(t *testing.T) {
	weakTag := createETag("W/\"\"")
	weakTag1 := createETag("W/\"1\"")
	weakTag2 := createETag("W/\"Two\"")
	strongTag1 := createETag("\"1\"")
	strongTag2 := createETag("\"Two\"")
	strongTagValidChars := createETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	weakTagValidChars := createETag("W/\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")

	require.Falsef(t, weakTag.Equals(weakTag), "Expected etags to not be equal")
	require.Falsef(t, weakTag1.Equals(weakTag1), "Expected etags to not be equal")
	require.Falsef(t, weakTag2.Equals(weakTag2), "Expected etags to not be equal")
	require.Falsef(t, weakTagValidChars.Equals(weakTagValidChars), "Expected etags to not be equal")

	require.Truef(t, strongTag1.Equals(strongTag1), "Expected etags to be equal")
	require.Truef(t, strongTag2.Equals(strongTag2), "Expected etags to be equal")
	require.Truef(t, strongTagValidChars.Equals(strongTagValidChars), "Expected etags to be equal")

	require.Falsef(t, weakTag1.Equals(weakTag), "Expected etags to not be equal")
	require.Falsef(t, strongTagValidChars.Equals(weakTagValidChars), "Expected etags to not be equal")
	require.Falsef(t, weakTag2.Equals(weakTag1), "Expected etags to not be equal")
	require.Falsef(t, strongTag1.Equals(weakTag1), "Expected etags to not be equal")
	require.Falsef(t, strongTag2.Equals(weakTag2), "Expected etags to not be equal")
}

func TestEtagAny(t *testing.T) {
	anyETag := ETagAny
	star := createETag("*")
	weakStar := createETag("W\"*\"")
	quotedStar := createETag("\"*\"")

	require.Truef(t, anyETag.Equals(anyETag), "Expected etags to be equal")
	require.Truef(t, ETagAny.Equals(anyETag), "Expected etags to be equal")

	require.Truef(t, star.Equals(star), "Expected etags to be equal")
	require.Truef(t, ETagAny.Equals(star), "Expected etags to be equal")
	require.Truef(t, anyETag.Equals(star), "Expected etags to be equal")

	require.Falsef(t, star.Equals(weakStar), "Expected etags to be equal")
	require.Falsef(t, ETagAny.Equals(weakStar), "Expected etags to be equal")
	require.Falsef(t, weakStar.Equals(quotedStar), "Expected etags to be equal")

	require.Falsef(t, quotedStar.Equals(star), "Expected etags to be equal")

	require.Truef(t, star.Equals(ETagAny), "Expected etags to be equal")
}

func TestETagWeakComparison(t *testing.T) {
	// W/""
	weakTag := createETag("W/\"\"")
	// W/"1"
	weakTag1 := createETag("W/\"1\"")
	// W/"Two"
	weakTagTwo := createETag("W/\"Two\"")
	// W/"two"
	weakTagtwo := createETag("W/\"two\"")
	// "1"
	strongTag1 := createETag("\"1\"")
	// "Two"
	strongTagTwo := createETag("\"Two\"")
	// "two"
	strongTagtwo := createETag("\"two\"")

	require.Truef(t, weakTag.WeakEquals(weakTag), "expected etags to be equal")
	require.Truef(t, weakTag1.WeakEquals(weakTag1), "expected etags to be equal")
	require.Truef(t, weakTagTwo.WeakEquals(weakTagTwo), "expected etags to be equal")
	require.Truef(t, weakTagtwo.WeakEquals(weakTagtwo), "expected etags to be equal")
	require.Truef(t, strongTag1.WeakEquals(strongTag1), "expected etags to be equal")
	require.Truef(t, strongTagTwo.WeakEquals(strongTagTwo), "expected etags to be equal")
	require.Truef(t, strongTagtwo.WeakEquals(strongTagtwo), "expected etags to be equal")

	require.Falsef(t, weakTag1.WeakEquals(weakTag), "Expected etags to not be equal")
	require.Falsef(t, weakTagTwo.WeakEquals(weakTag1), "Expected etags to not be equal")

	require.Truef(t, strongTag1.WeakEquals(weakTag1), "expected etags to be equal")
	require.Truef(t, strongTagTwo.WeakEquals(weakTagTwo), "expected etags to be equal")

	require.Falsef(t, weakTag1.WeakEquals(strongTagTwo), "Expected etags to not be equal")
	require.Falsef(t, weakTagtwo.WeakEquals(strongTagTwo), "Expected etags to not be equal")

	require.Falsef(t, strongTagtwo.WeakEquals(strongTagTwo), "Expected etags to not be equal")
	require.Falsef(t, weakTagtwo.WeakEquals(weakTagTwo), "Expected etags to not be equal")
}

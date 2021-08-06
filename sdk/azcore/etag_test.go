// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createETag(s string) ETag {
	return ETag(&s)
}

func TestETagEquals(t *testing.T) {
	e1 := createETag("tag")
	require.Equal(t, *e1, "tag")

	e2 := createETag("\"tag\"")
	require.Equal(t, *e2, "\"tag\"")

	e3 := createETag("W/\"weakETag\"")
	require.Equal(t, *e3, "W/\"weakETag\"")
	require.Truef(t, IsWeak(e3), "ETag is ecpected to be weak")

	strongETag := createETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	require.Equal(t, *strongETag, "\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")

	require.Falsef(t, IsWeak(ETagAny()), "ETagAny should not be weak")
}

func TestETagWeak(t *testing.T) {
	et1 := createETag("tag")
	require.Falsef(t, IsWeak(et1), "expected etag to be strong")

	et2 := createETag("\"tag\"")
	require.Falsef(t, IsWeak(et2), "expected etag to be strong")

	et3 := createETag("W/\"weakETag\"")
	require.Truef(t, IsWeak(et3), "expected etag to be weak")

	et4 := createETag("W/\"\"")
	require.Truef(t, IsWeak(et4), "expected etag to be weak")

	et5 := ETagAny()
	require.Falsef(t, IsWeak(et5), "expected etag to be strong")
}

func TestETagEquality(t *testing.T) {
	weakTag := createETag("W/\"\"")
	weakTag1 := createETag("W/\"1\"")
	weakTag2 := createETag("W/\"Two\"")
	strongTag1 := createETag("\"1\"")
	strongTag2 := createETag("\"Two\"")
	strongTagValidChars := createETag("\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")
	weakTagValidChars := createETag("W/\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"")

	require.Falsef(t, StrongEquals(weakTag, weakTag), "Expected etags to not be equal")
	require.Falsef(t, StrongEquals(weakTag1, weakTag1), "Expected etags to not be equal")
	require.Falsef(t, StrongEquals(weakTag2, weakTag2), "Expected etags to not be equal")
	require.Falsef(t, StrongEquals(weakTagValidChars, weakTagValidChars), "Expected etags to not be equal")

	require.Truef(t, StrongEquals(strongTag1, strongTag1), "Expected etags to be equal")
	require.Truef(t, StrongEquals(strongTag2, strongTag2), "Expected etags to be equal")
	require.Truef(t, StrongEquals(strongTagValidChars, strongTagValidChars), "Expected etags to be equal")

	require.Falsef(t, StrongEquals(weakTag, weakTag1), "Expected etags to not be equal")
	require.Falsef(t, StrongEquals(weakTagValidChars, strongTagValidChars), "Expected etags to not be equal")
	require.Falsef(t, StrongEquals(weakTag1, weakTag2), "Expected etags to not be equal")
	require.Falsef(t, StrongEquals(weakTag1, strongTag1), "Expected etags to not be equal")
	require.Falsef(t, StrongEquals(weakTag2, strongTag2), "Expected etags to not be equal")
}

func TestEtagAny(t *testing.T) {
	anyETag := ETagAny()
	star := createETag("*")
	weakStar := createETag("W\"*\"")
	quotedStart := createETag("\"*\"")

	require.Truef(t, StrongEquals(anyETag, anyETag), "Expected etags to be equal")
	require.Truef(t, StrongEquals(anyETag, ETagAny()), "Expected etags to be equal")

	require.Truef(t, StrongEquals(star, star), "Expected etags to be equal")
	require.Truef(t, StrongEquals(star, ETagAny()), "Expected etags to be equal")
	require.Truef(t, StrongEquals(star, anyETag), "Expected etags to be equal")

	require.Falsef(t, StrongEquals(weakStar, star), "Expected etags to be equal")
	require.Falsef(t, StrongEquals(weakStar, ETagAny()), "Expected etags to be equal")
	require.Falsef(t, StrongEquals(quotedStart, weakStar), "Expected etags to be equal")

	require.Falsef(t, StrongEquals(star, quotedStart), "Expected etags to be equal")
	require.Falsef(t, StrongEquals(ETagAny(), star), "Expected etags to be equal")
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

	require.Truef(t, WeakEquals(weakTag, weakTag), "expected etags to be equal")
	require.Truef(t, WeakEquals(weakTag1, weakTag1), "expected etags to be equal")
	require.Truef(t, WeakEquals(weakTagTwo, weakTagTwo), "expected etags to be equal")
	require.Truef(t, WeakEquals(weakTagtwo, weakTagtwo), "expected etags to be equal")
	require.Truef(t, WeakEquals(strongTag1, strongTag1), "expected etags to be equal")
	require.Truef(t, WeakEquals(strongTagTwo, strongTagTwo), "expected etags to be equal")
	require.Truef(t, WeakEquals(strongTagtwo, strongTagtwo), "expected etags to be equal")

	require.Falsef(t, WeakEquals(weakTag, weakTag1), "Expected etags to not be equal")
	require.Falsef(t, WeakEquals(weakTag1, weakTagTwo), "Expected etags to not be equal")

	require.Truef(t, WeakEquals(weakTag1, strongTag1), "expected etags to be equal")
	require.Truef(t, WeakEquals(weakTagTwo, strongTagTwo), "expected etags to be equal")

	require.Falsef(t, WeakEquals(strongTagTwo, weakTag1), "Expected etags to not be equal")
	require.Falsef(t, WeakEquals(strongTagTwo, weakTagtwo), "Expected etags to not be equal")

	require.Falsef(t, WeakEquals(strongTagTwo, strongTagtwo), "Expected etags to not be equal")
	require.Falsef(t, WeakEquals(weakTagTwo, weakTagtwo), "Expected etags to not be equal")
}

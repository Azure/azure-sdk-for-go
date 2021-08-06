// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"strings"
)

// ETag is a property used for optimistic concurrency during updates
// ETag is a validator based on https://tools.ietf.org/html/rfc7232#section-2.3.2
type ETag *string

// StrongEquals does a strong comparison of two ETags. StrongEquals returns true when both
// ETags are not weak and the values of the underlying strings are equal. If both ETags are "nil" they are considered equal
func StrongEquals(a, b ETag) bool {
	if !HasValue(a) || !HasValue(b) {
		return a == b
	}
	return !IsWeak(a) && !IsWeak(b) && *a == *b
}

// WeakEquals does a weak compariosn of two ETags. Two ETags are equivalent if their opaque-tags match
// character-by-character, regardless of either or both being tagged as "weak". If both ETags are "nil" they are considered equal
func WeakEquals(a, b ETag) bool {
	if !HasValue(a) || !HasValue(b) {
		return a == b
	}

	getStart := func(e ETag) int {
		if IsWeak(e) {
			return 2
		}
		return 0
	}
	aStart := getStart(a)
	bStart := getStart(b)

	aVal := (*a)[aStart:]
	bVal := (*b)[bStart:]

	return aVal == bVal
}

// HasValue returns whether an ETag is present
func HasValue(e ETag) bool {
	return e != nil
}

// IsWeak specifies whether the ETag is strong or weak.
func IsWeak(e ETag) bool {
	return HasValue(e) && len(*e) >= 4 && strings.HasPrefix(*e, "W/\"") && strings.HasSuffix(*e, "\"")
}

// ETagAny returns a new ETag that represents everything, the value is "*"
func ETagAny() ETag {
	any := "*"
	return ETag(&any)
}

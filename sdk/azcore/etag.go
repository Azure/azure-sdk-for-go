// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"strings"
)

// ETag is a property used for optimistic concurrency during updates
// ETag is a validator based on https://tools.ietf.org/html/rfc7232#section-2.3.2
type ETag struct {
	value *string
}

// NewETag creates a new ETag struct
func NewETag(value string) *ETag {
	return &ETag{value: &value}
}

// ComparisonType specifies what type of comparison to use, Strong or Weak
type ComparisonType string

const (
	Strong ComparisonType = "strong"
	Weak   ComparisonType = "weak"
)

// Equals determines whether two entity-tags are equal. There are two types of comparison: Strong and Weak.
// Strong comparison: two ETags are equivalent if both are not weak and their opaque-tags match character-by-character
// Weak Comparison: two ETags are equivalent if their opaque-tags match character-by-character regardless of either or both being tagged as "weak"
func (e ETag) Equals(right ETag, comparisonKind ComparisonType) bool {
	// ETags are != if one value is null
	if *e.value == "" || *right.value == "" {
		// If both are null, they are considered equal
		return *e.value == *right.value
	}

	if comparisonKind == Strong {
		return !e.IsWeak() && !right.IsWeak() && *e.value == *right.value
	}

	leftStart := e.getStart()
	rightStart := right.getStart()

	leftValue := (*e.value)[leftStart:]
	rightValue := (*right.value)[rightStart:]

	return leftValue == rightValue
}

func (e ETag) getStart() int {
	if e.IsWeak() {
		return 2
	}
	return 0
}

// String returns a string representation of an ETag
func (e ETag) String() string {
	return *e.value
}

// HasValue returns whether an ETag is present
func (e ETag) HasValue() bool {
	return e.value != nil
}

// IsWeak specifies whether the ETag is strong or weak.
func (e ETag) IsWeak() bool {
	return e.value != nil && len(*e.value) >= 4 && strings.HasPrefix(*e.value, "W/\"") && strings.HasSuffix(*e.value, "\"")
}

// ETagAny returns a new ETag that represents everything, the value is "*"
func ETagAny() *ETag {
	any := "*"
	return &ETag{value: &any}
}

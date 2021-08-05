// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
	"strings"
)

// ETag is a property used for optimistic concurrency during updates
type ETag struct {
	value *string
}

func NewETag(value string) *ETag {
	return &ETag{value: &value}
}

type ComparisonType string

const (
	Strong ComparisonType = "strong"
	Weak   ComparisonType = "weak"
)

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

func (e ETag) String() string {
	return *e.value
}

func (e ETag) HasValue() bool {
	return e.value != nil
}

func (e ETag) IsWeak() bool {
	// fmt.Println(e.value != nil)
	// fmt.Println(len(*e.value) >= 4)
	// fmt.Println(strings.HasPrefix(*e.value, "W/\""))
	// fmt.Println(strings.HasSuffix(*e.value, "\""))
	return e.value != nil && len(*e.value) >= 4 && strings.HasPrefix(*e.value, "W/\"") && strings.HasSuffix(*e.value, "\"")
}

func ETagAny() *ETag {
	any := "*"
	return &ETag{value: &any}
}

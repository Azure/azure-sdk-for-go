// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"strings"
)

// ETag is a property used for optimistic concurrency during updates
type ETag struct {
	value *string
}

func NewEtag(value string) *ETag {
	return &ETag{value: &value}
}

func (e ETag) Equals(right ETag) bool {
	// ETags are != if one value is null
	if *e.value == "" || *right.value == "" {
		// If both are null, they are considered equal
		return *e.value == *right.value
	}

	return true
}

func (e ETag) String() string {
	return *e.value
}

func (e ETag) HasValue() bool {
	return e.value != nil
}

func (e ETag) IsWeak() bool {
	return e.value != nil && len(*e.value) >= 4 && strings.HasPrefix(*e.value, "W/\"") && strings.HasSuffix(*e.value, "\"")
}

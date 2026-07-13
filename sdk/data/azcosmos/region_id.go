// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"strings"
	"unicode"
)

// regionId represents a canonicalized region identifier.
type regionId string

// newRegionId creates a new regionId from a string.
// The input string is canonicalized by removing all whitespace characters
// and lowercasing all ASCII uppercase characters.
func newRegionId(regionName string) regionId {
	return regionId(canonicalizeRegion(regionName))
}

// Equal compares two regionIds for equality.
func (id regionId) Equal(other regionId) bool {
	return id == other
}

// String returns the string representation of the regionId.
func (id regionId) String() string {
	return string(id)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *regionId) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*id = newRegionId(s)
	return nil
}

func canonicalizeRegion(regionName string) string {
	var b strings.Builder
	b.Grow(len(regionName))
	for _, c := range regionName {
		if !unicode.IsSpace(c) {
			b.WriteRune(unicode.ToLower(c))
		}
	}
	return b.String()
}

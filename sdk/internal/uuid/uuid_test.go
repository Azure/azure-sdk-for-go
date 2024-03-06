//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package uuid

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	u, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if reflect.ValueOf(u).IsZero() {
		t.Fatal("unexpected zero-value UUID")
	}
	s := u.String()
	match, err := regexp.MatchString(`[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}`, s)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("invalid UUID string %s", s)
	}
}

func TestParse(t *testing.T) {
	testCases := []string{
		"72d0f24f-82be-4016-729d-31fd13bd681e",
		"{72d0f24f-82be-4016-729d-31fd13bd681e}",
	}
	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			u, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if reflect.ValueOf(u).IsZero() {
				t.Fatal("unexpected zero-value UUID")
			}
			if len(input) > 36 {
				// strip off the {} as String() doesn't output them
				input = input[1:37]
			}
			if s := u.String(); s != input {
				t.Fatalf("didn't round trip: %s", s)
			}
		})
	}
}

func TestParseFail(t *testing.T) {
	testCases := []string{
		"72d0f24f-82be-4016-729d-31fd13bd681",
		"{72d0f24f-82be+4016-729d-31fd13bd681e}",
	}
	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err == nil {
				t.Fatalf("unexpected nil error for: %s", input)
			}
		})
	}
}

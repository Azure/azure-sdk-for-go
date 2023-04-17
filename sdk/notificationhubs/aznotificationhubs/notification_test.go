//go:build go1.20
// +build go1.20

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznotificationhubs

import "testing"

func TestCreateTagExpression(t *testing.T) {
	tags := []string{"language_en", "country_US"}
	tagExpression := CreateTagExpression(tags)

	if tagExpression != "language_en||country_US" {
		t.Fatalf(`CreateTagExpression %v`, tagExpression)
	}
}

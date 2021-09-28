//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import "testing"

var (
	testData = map[string]*ResourceIdentifier{
		"/subscriptions/00000000-0000-0000-0000-000000000000": &ResourceIdentifier{
			Parent:            RootResourceIdentifier,
			SubscriptionId:    "00000000-0000-0000-0000-000000000000",
			Provider:          "",
			ResourceGroupName: "",
			ResourceType:      SubscriptionResourceType,
			Name:              "00000000-0000-0000-0000-000000000000",
			IsChild:           false,
		},
	}
)

func TestParseResourceIdentifier(t *testing.T) {
	for input, expected := range testData {
		id, err := ParseResourceIdentifier(input)
		if err != nil && expected != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if err != nil && expected == nil {
			continue
		}
		if *id != *expected {
			t.Fatalf("resource id not identical, get %v, expected %v", *id, *expected)
		}
	}
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeserializeORSPolicies(t *testing.T) {

	headers := map[string]*string{
		"x-ms-or-111_111":   to.Ptr("Completed"),
		"x-ms-or-111_222":   to.Ptr("Failed"),
		"x-ms-or-222_111":   to.Ptr("Completed"),
		"x-ms-or-222_222":   to.Ptr("Failed"),
		"x-ms-or-policy-id": to.Ptr("333"),     // to be ignored
		"x-ms-not-related":  to.Ptr("garbage"), // to be ignored
	}

	result := deserializeORSPolicies(headers)
	require.NotEmpty(t, result)
	rules0, rules1 := *result[0].Rules, *result[1].Rules
	require.Len(t, result, 2)
	require.Len(t, rules0, 2)
	require.Len(t, rules1, 2)

	if rules0[0].RuleID == "111" {
		require.Equal(t, rules0[0].Status, "Completed")
	} else {
		require.Equal(t, rules0[0].Status, "Failed")
	}

	if rules0[1].RuleID == "222" {
		require.Equal(t, rules0[1].Status, "Failed")
	} else {
		require.Equal(t, rules0[1].Status, "Completed")
	}

	if rules1[0].RuleID == "111" {
		require.Equal(t, rules1[0].Status, "Completed")
	} else {
		require.Equal(t, rules1[0].Status, "Failed")
	}

	if rules1[1].RuleID == "222" {
		require.Equal(t, rules1[1].Status, "Failed")
	} else {
		require.Equal(t, rules1[1].Status, "Completed")
	}
}

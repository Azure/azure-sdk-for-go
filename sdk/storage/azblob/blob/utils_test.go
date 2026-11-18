// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseContentRangeTotal(t *testing.T) {
	require.Equal(t, int64(1000), parseContentRangeTotal("bytes 0-999/1000"))
	require.Equal(t, int64(5000000), parseContentRangeTotal("bytes 0-4194303/5000000"))
	require.Equal(t, int64(100), parseContentRangeTotal("bytes 50-99/100"))
	require.Equal(t, int64(0), parseContentRangeTotal("invalid"))
	require.Equal(t, int64(0), parseContentRangeTotal(""))
}

func TestParseContentRangeLength(t *testing.T) {
	require.Equal(t, int64(1000), parseContentRangeLength("bytes 0-999/1000"))
	require.Equal(t, int64(4194304), parseContentRangeLength("bytes 0-4194303/5000000"))
	require.Equal(t, int64(50), parseContentRangeLength("bytes 50-99/100"))
	require.Equal(t, int64(0), parseContentRangeLength("invalid"))
}

func TestDownloadStreamObjectReplicationRules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("x-ms-or-abc_123", "Complete")
		w.Header().Set("x-ms-or-abc_def", "Failed")
		w.Header().Set("x-ms-or-fed_123", "Failed")
		w.Header().Set("x-ms-or-policy-id", "destination-policy")
		w.Header().Set("x-ms-meta-extra_key", "metadata")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithNoCredential(server.URL+"/container/blob", nil)
	require.NoError(t, err)
	response, err := client.DownloadStream(t.Context(), nil)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Len(t, response.ObjectReplicationRules, 2)
	rules := make(map[string][]ObjectReplicationRules)
	for _, policy := range response.ObjectReplicationRules {
		rules[*policy.PolicyID] = *policy.Rules
	}
	require.ElementsMatch(t, []ObjectReplicationRules{
		{RuleID: "123", Status: "Complete"},
		{RuleID: "def", Status: "Failed"},
	}, rules["abc"])
	require.Equal(t, []ObjectReplicationRules{{RuleID: "123", Status: "Failed"}}, rules["fed"])
}

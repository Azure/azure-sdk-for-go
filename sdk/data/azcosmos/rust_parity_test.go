// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBinaryEncodingRustDefaultsAndEnvironment(t *testing.T) {
	for _, value := range []string{"1", "true", "YES", " on ", "false", "", "nope"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", value)
			want := value == "1" || value == "true" || value == "YES" || value == " on "
			require.Equal(t, want, resolveClientBinaryEncoding(nil).enabled())
			require.True(t, resolveClientBinaryEncoding(&BinaryEncodingOptions{}).enabled())
			require.False(t, resolveClientBinaryEncoding(&BinaryEncodingOptions{Enabled: to(false)}).enabled())
		})
	}
}

func TestClientCopiesDedicatedBinaryOptions(t *testing.T) {
	options := &BinaryEncodingOptions{Enabled: to(false), RequestTextResponse: true}
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", mustKeyCredential(t), &ClientOptions{BinaryEncoding: options})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	*options.Enabled = true
	options.RequestTextResponse = false
	require.False(t, client.binaryEncoding.enabled())
	require.True(t, client.binaryEncoding.RequestTextResponse)
}

func TestPatchTrackingRustRetentionAndUUID(t *testing.T) {
	for _, id := range []PatchTrackingID{
		"00112233-4455-6677-8899-AABBCCDDEEFF",
		"00112233445566778899aabbccddeeff",
		"{00112233-4455-6677-8899-aabbccddeeff}",
		"urn:uuid:00112233-4455-6677-8899-aabbccddeeff",
	} {
		normalized, err := id.normalized()
		require.NoError(t, err)
		require.Equal(t, PatchTrackingID("00112233-4455-6677-8899-aabbccddeeff"), normalized)
	}
	require.NoError(t, (PatchItemOptions{TrackingRetention: to(time.Duration(math.MaxInt64))}).validateTracking())
}

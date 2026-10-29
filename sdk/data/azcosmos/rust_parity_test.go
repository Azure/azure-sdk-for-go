// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBinaryEncodingRustDefaultsAndEnvironment(t *testing.T) {
	t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", "")
	require.NoError(t, os.Unsetenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED"))
	require.False(t, runtimeBinaryEncoding(nil).Disabled)
	for _, value := range []string{"1", "true", "YES", " on ", "false", "", "nope"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("AZURE_COSMOS_BINARY_ENCODING_ENABLED", value)
			want := value == "1" || value == "true" || value == "YES" || value == " on "
			require.Equal(t, !want, runtimeBinaryEncoding(nil).Disabled)
			require.False(t, runtimeBinaryEncoding(&BinaryEncodingOptions{}).Disabled)
			require.True(t, runtimeBinaryEncoding(&BinaryEncodingOptions{Disabled: true}).Disabled)
		})
	}
}

func TestClientCopiesOperationBinaryOptions(t *testing.T) {
	options := &BinaryEncodingOptions{Disabled: true}
	client, err := newTestClientWithKey("https://myaccount.documents.azure.com", mustKeyCredential(t), &ClientOptions{Operation: OperationOptions{BinaryEncoding: options}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	options.Disabled = false
	require.True(t, client.options.Operation.BinaryEncoding.Disabled)
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

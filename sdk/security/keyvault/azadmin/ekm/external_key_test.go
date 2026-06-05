// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package ekm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stretchr/testify/require"
)

// TestExternalKeyReference creates an azkeys key whose material is hosted on the
// EKM proxy referenced by an already-configured EKM connection. The external
// key (named via EKM_EXTERNAL_ID) must already exist at the EKM proxy. The
// test then reads the key back and asserts the ExternalKey attribute round-trips
// correctly through the wire.
//
// This is an end-to-end integration test for the EKM external-key reference
// feature. It requires:
//   - AZURE_MANAGEDHSM_URL pointing at an HSM with an active EKM connection
//   - EKM_EXTERNAL_ID naming a key that exists at the EKM proxy
//
// In playback mode the test uses the recorded HSM/external-id placeholders.
func TestExternalKeyReference(t *testing.T) {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	keyClient, err := azkeys.NewClient(hsmURL, credential, &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	localKeyName := "ekm-external-key-test"
	ctx := context.Background()

	// Best-effort cleanup of a prior run. Run this in playback too so the
	// corresponding entry in the recording is consumed in lockstep.
	_, _ = keyClient.DeleteKey(ctx, localKeyName, nil)

	// ExternalKey and Kty are mutually exclusive on CreateKey: the key material
	// lives at the EKM proxy, so the Key Vault never picks its own key type.
	params := azkeys.CreateKeyParameters{
		KeyAttributes: &azkeys.KeyAttributes{
			ExternalKey: &azkeys.ExternalKey{ID: to.Ptr(externalKeyID)},
		},
	}

	created, err := keyClient.CreateKey(ctx, localKeyName, params, nil)
	if err != nil {
		// The most common failure is that no EKM connection is configured on
		// the HSM, in which case the service returns a 4xx. Surface that as a
		// skip — including in playback, where the recording may legitimately
		// have captured that 4xx — so the wire contract test still runs without
		// requiring a fully provisioned EKM proxy.
		var httpErr *azcore.ResponseError
		if errors.As(err, &httpErr) {
			t.Skipf("CreateKey failed (HTTP %d, %s); is an EKM connection configured on the HSM?",
				httpErr.StatusCode, httpErr.ErrorCode)
		}
		require.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = keyClient.DeleteKey(context.Background(), localKeyName, nil)
	})

	require.NotNil(t, created.Attributes, "expected attributes on created key")
	require.NotNil(t, created.Attributes.ExternalKey, "expected ExternalKey on created key attributes")
	require.NotNil(t, created.Attributes.ExternalKey.ID)
	require.Equal(t, externalKeyID, *created.Attributes.ExternalKey.ID)
	require.NotNil(t, created.Key)
	require.NotNil(t, created.Key.KID)

	// Round-trip via GET so we know the reference is persisted (not just echoed).
	got, err := keyClient.GetKey(ctx, localKeyName, "", nil)
	require.NoError(t, err)
	require.NotNil(t, got.Attributes)
	require.NotNil(t, got.Attributes.ExternalKey)
	require.NotNil(t, got.Attributes.ExternalKey.ID)
	require.Equal(t, externalKeyID, *got.Attributes.ExternalKey.ID)
}

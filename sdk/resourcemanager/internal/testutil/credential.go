//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

// FakeCredential is an empty credential for testing.
type FakeCredential struct {
}

// GetToken provide a fake access token.
func (c *FakeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "FakeToken", ExpiresOn: time.Now().Add(time.Hour * 24).UTC()}, nil
}

// GetCredAndClientOptions will create a credential and a client options for test application.
// They can be used in any Azure resource management client.
// The client options will initialize the transport for recording client add recording policy to the pipeline.
// In the record mode, the credential will be a DefaultAzureCredential which combines several common credentials.
// In the playback mode, the credential will be a fake credential which will bypass truly authorization.
func GetCredAndClientOptions(t *testing.T) (azcore.TokenCredential, *arm.ClientOptions) {
	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	if err != nil {
		t.Fatalf("Failed to create recording client: %v", err)
	}

	options := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			PerCallPolicies: []policy.Policy{p},
			Transport:       client,
		},
	}

	var cred azcore.TokenCredential
	if recording.GetRecordMode() != recording.PlaybackMode {
		cred, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			t.Fatalf("Failed to create credential: %v", err)
		}
	} else {
		cred = &FakeCredential{}
	}

	return cred, options
}

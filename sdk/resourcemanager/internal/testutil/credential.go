//go:build go1.16
// +build go1.16

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

type FakeCredential struct {
}

func (c *FakeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{Token: "FakeToken", ExpiresOn: time.Now().Add(time.Hour * 24).UTC()}, nil
}

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

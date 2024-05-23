// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package credential

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

var recordMode = recording.GetRecordMode()

type NewOptions struct{}

// New constructs an [azcore.TokenCredential] for use in tests
func New(*NewOptions) (azcore.TokenCredential, error) {
	if recordMode == recording.PlaybackMode {
		return &fakeCredential{}, nil
	}
	if s := os.Getenv("AZURE_SERVICE_DIRECTORY"); s != "" {
		// New-TestResources.ps1 has configured this environment, possibly with service principal details
		clientID := os.Getenv(s + "_CLIENT_ID")
		secret := os.Getenv(s + "_CLIENT_SECRET")
		tenant := os.Getenv(s + "_TENANT_ID")
		if clientID != "" && secret != "" && tenant != "" {
			return azidentity.NewClientSecretCredential(tenant, clientID, secret, nil)
		}
	}
	return azidentity.NewDefaultAzureCredential(nil)
}

type fakeCredential struct{}

func (fakeCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{ExpiresOn: time.Now().Add(time.Hour).UTC(), Token: recording.SanitizedValue}, nil
}

var _ azcore.TokenCredential = (*fakeCredential)(nil)

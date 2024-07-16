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

// New constructs a credential for use in tests. In playback mode, it returns [Fake], which always
// provides valid tokens. In live and record modes, it returns a credential based on the environment:
//   - [azidentity.AzurePipelinesCredential] when running in an Azure Pipelines AzurePowerShell task
//   - [azidentity.ClientSecretCredential], if New-TestResources.ps1 provided service principal details
//   - [azidentity.DefaultAzureCredential] otherwise
func New(*NewOptions) (azcore.TokenCredential, error) {
	if recordMode == recording.PlaybackMode {
		return &Fake{}, nil
	}
	accessToken := os.Getenv("SYSTEM_ACCESSTOKEN")
	clientID := os.Getenv("AZURESUBSCRIPTION_CLIENT_ID")
	connectionID := os.Getenv("AZURESUBSCRIPTION_SERVICE_CONNECTION_ID")
	tenant := os.Getenv("AZURESUBSCRIPTION_TENANT_ID")
	if accessToken != "" && clientID != "" && connectionID != "" && tenant != "" {
		return azidentity.NewAzurePipelinesCredential(tenant, clientID, connectionID, accessToken, nil)
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

// Fake always returns a valid token. Use this type to fake authentication in tests
// that never send a real request. For live or recorded tests, call [New] instead.
type Fake struct{}

// GetToken returns a fake access token
func (Fake) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{ExpiresOn: time.Now().Add(time.Hour).UTC(), Token: recording.SanitizedValue}, nil
}

var _ azcore.TokenCredential = (*Fake)(nil)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testutil

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
)

// FakeCredential is an empty credential for testing.
//
// Deprecated: use Fake from github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential
type FakeCredential = credential.Fake

// GetCredAndClientOptions will create a credential and a client options for test application.
// They can be used in any Azure resource management client.
// The client options will initialize the transport for recording client add recording policy to the pipeline.
// In the record mode, the credential will be a DefaultAzureCredential which combines several common credentials.
// In the playback mode, the credential will be a fake credential which will bypass truly authorization.
func GetCredAndClientOptions(t *testing.T) (azcore.TokenCredential, *arm.ClientOptions) {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	if err != nil {
		t.Fatalf("Failed to create recording transport: %v", err)
	}

	options := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: transport,
		},
	}

	cred, err := credential.New(nil)
	if err != nil {
		t.Fatalf("Failed to create credential: %v", err)
	}

	return cred, options
}

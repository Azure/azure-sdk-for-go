// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package credential

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	before := recordMode
	for _, pipeline := range []bool{true, false} {
		name := "dev/"
		if pipeline {
			name = "pipeline/"
		}
		for _, mode := range []string{recording.LiveMode, recording.PlaybackMode, recording.RecordingMode} {
			t.Run(name+mode, func(t *testing.T) {
				recordMode = mode
				defer func() { recordMode = before }()
				if pipeline {
					// set environment variables as would New-TestResources.ps1
					t.Setenv("AZURE_SERVICE_DIRECTORY", t.Name())
					for _, v := range []string{"_CLIENT_ID", "_CLIENT_SECRET", "_TENANT_ID"} {
						t.Setenv(t.Name()+v, recording.SanitizedValue)
					}
				}
				cred, err := New(nil)
				require.NoError(t, err)
				require.NotNil(t, cred)
				if mode == recording.PlaybackMode {
					tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"scope"}})
					require.NoError(t, err)
					require.Equal(t, recording.SanitizedValue, tk.Token)
					require.True(t, tk.ExpiresOn.After(time.Now()))
				} else {
					if pipeline {
						require.IsType(t, &azidentity.ClientSecretCredential{}, cred)
					} else {
						require.IsType(t, &azidentity.DefaultAzureCredential{}, cred)
					}
				}
			})
		}
	}
}

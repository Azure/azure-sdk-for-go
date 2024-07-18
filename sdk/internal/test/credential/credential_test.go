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
	for _, test := range []struct {
		name     string
		env      map[string]string
		expected any
	}{
		{
			name: "pipeline/SP/",
			env: map[string]string{
				"AZURE_SERVICE_DIRECTORY":   t.Name(),
				t.Name() + "_CLIENT_ID":     "clientID",
				t.Name() + "_CLIENT_SECRET": "secret",
				t.Name() + "_TENANT_ID":     "tenant",
			},
			expected: &azidentity.ClientSecretCredential{},
		},
		{
			name: "pipeline/WIF/",
			env: map[string]string{
				"AZURESUBSCRIPTION_CLIENT_ID":             "clientID",
				"AZURESUBSCRIPTION_SERVICE_CONNECTION_ID": "connectionID",
				"AZURESUBSCRIPTION_TENANT_ID":             "tenant",
				"SYSTEM_ACCESSTOKEN":                      "token",
				"SYSTEM_OIDCREQUESTURI":                   "https://localhost",
			},
			expected: &azidentity.AzurePipelinesCredential{},
		},
		{
			name:     "dev/",
			expected: &azidentity.DefaultAzureCredential{},
		},
	} {
		for _, mode := range []string{recording.LiveMode, recording.PlaybackMode, recording.RecordingMode} {
			t.Run(test.name+mode, func(t *testing.T) {
				recordMode = mode
				defer func() { recordMode = before }()
				for k, v := range test.env {
					t.Setenv(k, v)
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
					require.IsType(t, test.expected, cred)
				}
			})
		}
	}
}

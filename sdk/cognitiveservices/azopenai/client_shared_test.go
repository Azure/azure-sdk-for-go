// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	endpoint string
	apiKey   string
)

const fakeEndpoint = "https://recordedhost/"
const fakeAPIKey = "redacted"

func init() {
	if recording.GetRecordMode() == recording.PlaybackMode {
		endpoint = fakeEndpoint
		apiKey = fakeAPIKey
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Failed to load .env file: %s\n", err)
			os.Exit(1)
		}

		endpoint = os.Getenv("AOAI_ENDPOINT")
		apiKey = os.Getenv("AOAI_API_KEY")
	}
}

func newRecordingTransporter(t *testing.T) policy.Transporter {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	err = recording.Start(t, "sdk/cognitiveservices/azopenai/testdata", nil)
	require.NoError(t, err)

	if recording.GetRecordMode() != recording.PlaybackMode {
		err = recording.AddHeaderRegexSanitizer("Api-Key", fakeAPIKey, "", nil)
		require.NoError(t, err)

		// "RequestUri": "https://openai-shared.openai.azure.com/openai/deployments/text-davinci-003/completions?api-version=2023-03-15-preview",
		err = recording.AddURISanitizer(fakeEndpoint, endpoint, nil)
		require.NoError(t, err)
	}

	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})

	return transport
}

func newClientOptionsForTest(t *testing.T) *ClientOptions {
	co := &ClientOptions{}
	co.Transport = newRecordingTransporter(t)
	return co
}

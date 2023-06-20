// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	endpoint                 string // env: AOAI_ENDPOINT
	apiKey                   string // env: AOAI_API_KEY
	streamingModelDeployment string // env: AOAI_STREAMING_MODEL_DEPLOYMENT

	openAIKey      string // env: OPENAI_API_KEY
	openAIEndpoint string // env: OPENAI_ENDPOINT
)

const fakeEndpoint = "https://recordedhost/"
const fakeAPIKey = "redacted"

func init() {
	if recording.GetRecordMode() == recording.PlaybackMode {
		endpoint = fakeEndpoint
		apiKey = fakeAPIKey
		openAIKey = fakeAPIKey
		openAIEndpoint = fakeEndpoint
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Failed to load .env file: %s\n", err)
			os.Exit(1)
		}

		endpoint = os.Getenv("AOAI_ENDPOINT")

		if endpoint != "" && !strings.HasSuffix(endpoint, "/") {
			// (this just makes recording replacement easier)
			endpoint += "/"
		}

		apiKey = os.Getenv("AOAI_API_KEY")

		// Ex: text-davinci-003
		streamingModelDeployment = os.Getenv("AOAI_STREAMING_MODEL_DEPLOYMENT")

		openAIKey = os.Getenv("OPENAI_API_KEY")
		openAIEndpoint = os.Getenv("OPENAI_ENDPOINT")

		if openAIEndpoint != "" && !strings.HasSuffix(openAIEndpoint, "/") {
			// (this just makes recording replacement easier)
			openAIEndpoint += "/"
		}
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
		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(endpoint), nil)
		require.NoError(t, err)

		if openAIEndpoint != "" {
			err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(openAIEndpoint), nil)
			require.NoError(t, err)
		}
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

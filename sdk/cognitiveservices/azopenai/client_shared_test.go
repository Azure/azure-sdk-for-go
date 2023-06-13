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

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load .env file: %s\n", err)
		os.Exit(1)
	}
	endpoint = os.Getenv("AOAI_ENDPOINT")
	apiKey = os.Getenv("AOAI_API_KEY")
}

func newRecordingTransporter(t *testing.T) policy.Transporter {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	err = recording.Start(t, "sdk/cognitive/azopenai/testdata", nil)
	require.NoError(t, err)

	err = recording.AddHeaderRegexSanitizer("Api-Key", `redacted`, "", nil)
	require.NoError(t, err)

	err = recording.AddURISanitizer("https://aoai", "https://[^/]+", nil)
	require.NoError(t, err)

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

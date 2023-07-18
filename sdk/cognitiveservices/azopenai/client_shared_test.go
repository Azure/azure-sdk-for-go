// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	endpoint                       string // env: AOAI_ENDPOINT
	apiKey                         string // env: AOAI_API_KEY
	completionsModelDeployment     string // env: AOAI_COMPLETIONS_MODEL_DEPLOYMENT
	chatCompletionsModelDeployment string // env: AOAI_CHAT_COMPLETIONS_MODEL_DEPLOYMENT

	canaryEndpoint                       string // env: AOAI_ENDPOINT_CANARY
	canaryAPIKey                         string // env: AOAI_API_KEY_CANARY
	canaryCompletionsModelDeployment     string // env: AOAI_COMPLETIONS_MODEL_DEPLOYMENT_CANARY
	canaryChatCompletionsModelDeployment string // env: AOAI_CHAT_COMPLETIONS_MODEL_DEPLOYMENT_CANARY

	openAIKey                  string // env: OPENAI_API_KEY
	openAIEndpoint             string // env: OPENAI_ENDPOINT
	openAICompletionsModel     string // env: OPENAI_CHAT_COMPLETIONS_MODEL
	openAIChatCompletionsModel string // env: OPENAI_COMPLETIONS_MODEL
)

func getVars(suffix string) (endpoint, apiKey, completionsModelDeployment, chatCompletionsModelDeployment string) {
	endpoint = os.Getenv("AOAI_ENDPOINT" + suffix)

	if endpoint != "" && !strings.HasSuffix(endpoint, "/") {
		// (this just makes recording replacement easier)
		endpoint += "/"
	}

	apiKey = os.Getenv("AOAI_API_KEY" + suffix)
	completionsModelDeployment = os.Getenv("AOAI_COMPLETIONS_MODEL_DEPLOYMENT" + suffix)

	// ex: gpt-4-0613
	chatCompletionsModelDeployment = os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_DEPLOYMENT" + suffix)

	return
}

const fakeEndpoint = "https://recordedhost/"
const fakeAPIKey = "redacted"

func initEnvVars() {
	if recording.GetRecordMode() == recording.PlaybackMode {
		endpoint = fakeEndpoint
		apiKey = fakeAPIKey
		openAIKey = fakeAPIKey
		openAIEndpoint = fakeEndpoint

		canaryEndpoint = fakeEndpoint
		canaryAPIKey = fakeAPIKey
		canaryCompletionsModelDeployment = ""
		canaryChatCompletionsModelDeployment = "gpt-4"

		completionsModelDeployment = "text-davinci-003"
		openAICompletionsModel = "text-davinci-003"

		chatCompletionsModelDeployment = "gpt-4-0613"
		openAIChatCompletionsModel = "gpt-4"
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Failed to load .env file: %s\n", err)
			os.Exit(1)
		}

		endpoint, apiKey, completionsModelDeployment, chatCompletionsModelDeployment = getVars("")
		canaryEndpoint, canaryAPIKey, canaryCompletionsModelDeployment, canaryChatCompletionsModelDeployment = getVars("_CANARY")

		openAIKey = os.Getenv("OPENAI_API_KEY")
		openAIEndpoint = os.Getenv("OPENAI_ENDPOINT")
		openAICompletionsModel = os.Getenv("OPENAI_COMPLETIONS_MODEL")
		openAIChatCompletionsModel = os.Getenv("OPENAI_CHAT_COMPLETIONS_MODEL") // ex: gpt-4-0613

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

		err = recording.AddHeaderRegexSanitizer("User-Agent", "fake-user-agent", ".*", nil)
		require.NoError(t, err)

		// "RequestUri": "https://openai-shared.openai.azure.com/openai/deployments/text-davinci-003/completions?api-version=2023-03-15-preview",
		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(endpoint), nil)
		require.NoError(t, err)

		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(canaryEndpoint), nil)
		require.NoError(t, err)

		err = recording.AddURISanitizer("/openai/operations/images/00000000-AAAA-BBBB-CCCC-DDDDDDDDDDDD", "/openai/operations/images/[A-Za-z-0-9]+", nil)
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

func newClientOptionsForTest(t *testing.T) *azopenai.ClientOptions {
	co := &azopenai.ClientOptions{}

	if recording.GetRecordMode() == recording.LiveMode {
		keyLogPath := os.Getenv("SSLKEYLOGFILE")

		if keyLogPath == "" {
			return nil
		}

		keyLogWriter, err := os.OpenFile(keyLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = keyLogWriter.Close()
		})

		tp := http.DefaultTransport.(*http.Transport).Clone()
		tp.TLSClientConfig = &tls.Config{
			KeyLogWriter: keyLogWriter,
		}

		co.Transport = &http.Client{Transport: tp}
	} else {
		co.Transport = newRecordingTransporter(t)
	}

	return co
}

// newAzureOpenAIClientForTest can create a client pointing to the "canary" endpoint (basically - leading fixes or features)
// or the current deployed endpoint.
func newAzureOpenAIClientForTest(t *testing.T, modelDeploymentID string, useCanary bool) *azopenai.Client {
	var apiKey = apiKey
	var endpoint = endpoint

	if useCanary {
		apiKey = canaryAPIKey
		endpoint = canaryEndpoint
	}

	cred, err := azopenai.NewKeyCredential(apiKey)
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(endpoint, cred, modelDeploymentID, newClientOptionsForTest(t))
	require.NoError(t, err)

	return client
}

func newOpenAIClientForTest(t *testing.T) *azopenai.Client {
	if openAIKey == "" {
		t.Skipf("OPENAI_API_KEY not defined, skipping OpenAI public endpoint test")
	}

	cred, err := azopenai.NewKeyCredential(openAIKey)
	require.NoError(t, err)

	// we get rate limited quite a bit.
	options := newClientOptionsForTest(t)

	if options == nil {
		options = &azopenai.ClientOptions{}
	}

	options.Retry = policy.RetryOptions{
		MaxRetries:    60,
		RetryDelay:    time.Second,
		MaxRetryDelay: time.Second,
	}

	chatClient, err := azopenai.NewClientForOpenAI(openAIEndpoint, cred, options)
	require.NoError(t, err)

	return chatClient
}

// newBogusAzureOpenAIClient creates a client that uses an invalid key, which will cause Azure OpenAI to return
// a failure.
func newBogusAzureOpenAIClient(t *testing.T, modelDeploymentID string) *azopenai.Client {
	cred, err := azopenai.NewKeyCredential("bogus-api-key")
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(endpoint, cred, modelDeploymentID, newClientOptionsForTest(t))
	require.NoError(t, err)
	return client
}

// newBogusOpenAIClient creates a client that uses an invalid key, which will cause OpenAI to return
// a failure.
func newBogusOpenAIClient(t *testing.T) *azopenai.Client {
	cred, err := azopenai.NewKeyCredential("bogus-api-key")
	require.NoError(t, err)

	client, err := azopenai.NewClientForOpenAI(openAIEndpoint, cred, newClientOptionsForTest(t))
	require.NoError(t, err)
	return client
}

func assertResponseIsError(t *testing.T, err error) {
	t.Helper()

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)

	// we sometimes get rate limited but (for this kind of test) it's actually okay
	require.Truef(t, respErr.StatusCode == http.StatusUnauthorized || respErr.StatusCode == http.StatusTooManyRequests, "An acceptable error comes back (actual: %d)", respErr.StatusCode)
}

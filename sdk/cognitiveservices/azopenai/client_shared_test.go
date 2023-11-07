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
	azureOpenAI       testVars
	azureOpenAICanary testVars
	openAI            testVars
)

type testVars struct {
	Endpoint        string // env: AOAI_ENDPOINT, OPENAI_ENDPOINT
	APIKey          string // env: AOAI_API_KEY, OPENAI_API_KEY
	Completions     string // env: AOAI_COMPLETIONS_MODEL_DEPLOYMENT, OPENAI_COMPLETIONS_MODEL
	ChatCompletions string // env: AOAI_CHAT_COMPLETIONS_MODEL_DEPLOYMENT, OPENAI_CHAT_COMPLETIONS_MODEL
	Embeddings      string // env: AOAI_EMBEDDINGS_MODEL_DEPLOYMENT, OPENAI_EMBEDDINGS_MODEL
	Azure           bool
}

func newTestVars(prefix string, isCanary bool) testVars {
	getRequired := func(name string) string {
		v := os.Getenv(name)

		if v == "" {
			panic(fmt.Sprintf("Env variable %s is missing", name))
		}

		return v
	}

	azure := prefix == "AOAI"

	canarySuffix := ""
	deplSuffix := ""

	if azure {
		deplSuffix += "_DEPLOYMENT"
	}

	if isCanary {
		canarySuffix += "_CANARY"
	}

	tv := testVars{
		Endpoint: getRequired(prefix + "_ENDPOINT" + canarySuffix),
		APIKey:   getRequired(prefix + "_API_KEY" + canarySuffix),

		Completions: getRequired(prefix + "_COMPLETIONS_MODEL" + deplSuffix + canarySuffix),

		// ex: gpt-4-0613
		ChatCompletions: getRequired(prefix + "_CHAT_COMPLETIONS_MODEL" + deplSuffix + canarySuffix),

		// ex: embedding
		Embeddings: getRequired(prefix + "_EMBEDDINGS_MODEL" + deplSuffix + canarySuffix),

		Azure: azure,
	}

	if tv.Endpoint != "" && !strings.HasSuffix(tv.Endpoint, "/") {
		// (this just makes recording replacement easier)
		tv.Endpoint += "/"
	}

	return tv
}

const fakeEndpoint = "https://recordedhost/"
const fakeAPIKey = "redacted"

func initEnvVars() {
	if recording.GetRecordMode() == recording.PlaybackMode {
		azureOpenAI.Azure = true
		azureOpenAI.Endpoint = fakeEndpoint
		azureOpenAI.APIKey = fakeAPIKey
		openAI.APIKey = fakeAPIKey
		openAI.Endpoint = fakeEndpoint

		azureOpenAICanary.Azure = true
		azureOpenAICanary.Endpoint = fakeEndpoint
		azureOpenAICanary.APIKey = fakeAPIKey
		azureOpenAICanary.Completions = ""
		azureOpenAICanary.ChatCompletions = "gpt-4"

		azureOpenAI.Completions = "text-davinci-003"
		openAI.Completions = "text-davinci-003"

		azureOpenAI.ChatCompletions = "gpt-4-0613"
		openAI.ChatCompletions = "gpt-4-0613"

		openAI.Embeddings = "text-similarity-curie-001"
		azureOpenAI.Embeddings = "embedding"
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Failed to load .env file: %s\n", err)
			os.Exit(1)
		}

		azureOpenAI = newTestVars("AOAI", false)
		azureOpenAICanary = newTestVars("AOAI", true)
		openAI = newTestVars("OPENAI", false)
	}
}

func newRecordingTransporter(t *testing.T) policy.Transporter {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	err = recording.Start(t, RecordingDirectory, nil)
	require.NoError(t, err)

	if recording.GetRecordMode() != recording.PlaybackMode {
		err = recording.AddHeaderRegexSanitizer("Api-Key", fakeAPIKey, "", nil)
		require.NoError(t, err)

		err = recording.AddHeaderRegexSanitizer("User-Agent", "fake-user-agent", ".*", nil)
		require.NoError(t, err)

		// "RequestUri": "https://openai-shared.openai.azure.com/openai/deployments/text-davinci-003/completions?api-version=2023-03-15-preview",
		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(azureOpenAI.Endpoint), nil)
		require.NoError(t, err)

		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(azureOpenAICanary.Endpoint), nil)
		require.NoError(t, err)

		err = recording.AddURISanitizer("/openai/operations/images/00000000-AAAA-BBBB-CCCC-DDDDDDDDDDDD", "/openai/operations/images/[A-Za-z-0-9]+", nil)
		require.NoError(t, err)

		if openAI.Endpoint != "" {
			err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(openAI.Endpoint), nil)
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
func newAzureOpenAIClientForTest(t *testing.T, tv testVars) *azopenai.Client {
	cred, err := azopenai.NewKeyCredential(tv.APIKey)
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(tv.Endpoint, cred, newClientOptionsForTest(t))
	require.NoError(t, err)

	return client
}

func newOpenAIClientForTest(t *testing.T) *azopenai.Client {
	if openAI.APIKey == "" {
		t.Skipf("OPENAI_API_KEY not defined, skipping OpenAI public endpoint test")
	}

	cred, err := azopenai.NewKeyCredential(openAI.APIKey)
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

	chatClient, err := azopenai.NewClientForOpenAI(openAI.Endpoint, cred, options)
	require.NoError(t, err)

	return chatClient
}

// newBogusAzureOpenAIClient creates a client that uses an invalid key, which will cause Azure OpenAI to return
// a failure.
func newBogusAzureOpenAIClient(t *testing.T) *azopenai.Client {
	cred, err := azopenai.NewKeyCredential("bogus-api-key")
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAI.Endpoint, cred, newClientOptionsForTest(t))
	require.NoError(t, err)
	return client
}

// newBogusOpenAIClient creates a client that uses an invalid key, which will cause OpenAI to return
// a failure.
func newBogusOpenAIClient(t *testing.T) *azopenai.Client {
	cred, err := azopenai.NewKeyCredential("bogus-api-key")
	require.NoError(t, err)

	client, err := azopenai.NewClientForOpenAI(openAI.Endpoint, cred, newClientOptionsForTest(t))
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

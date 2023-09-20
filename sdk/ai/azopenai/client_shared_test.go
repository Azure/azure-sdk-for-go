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

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	azureOpenAI        testVars
	azureOpenAICanary  testVars
	openAI             testVars
	azureWhisper       endpoint
	azureWhisperModel  string
	openAIWhisper      endpoint
	openAIWhisperModel string
)

type endpoint struct {
	URL    string
	APIKey string
	Azure  bool
}

func newTestClient(t *testing.T, ep endpoint) *azopenai.Client {
	if ep.Azure {
		cred, err := azopenai.NewKeyCredential(ep.APIKey)
		require.NoError(t, err)

		client, err := azopenai.NewClientWithKeyCredential(ep.URL, cred, newClientOptionsForTest(t))
		require.NoError(t, err)

		return client
	} else {
		if ep.APIKey == "" {
			t.Skipf("OPENAI_API_KEY not defined, skipping OpenAI public endpoint test")
		}

		cred, err := azopenai.NewKeyCredential(ep.APIKey)
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

		client, err := azopenai.NewClientForOpenAI(ep.URL, cred, options)
		require.NoError(t, err)

		return client
	}
}

type testVars struct {
	Endpoint        endpoint
	Completions     string
	ChatCompletions string
	Embeddings      string
	Cognitive       azopenai.AzureCognitiveSearchChatExtensionConfiguration
	Azure           bool
}

func newTestVars(prefix string, isCanary bool) testVars {
	azure := prefix == "AOAI"
	suffix := ""

	if isCanary {
		suffix += "_CANARY"
	}

	tv := testVars{
		Endpoint: endpoint{
			URL:    getRequired(prefix + "_ENDPOINT" + suffix),
			APIKey: getRequired(prefix + "_API_KEY" + suffix),
			Azure:  azure,
		},
		Completions: getRequired(prefix + "_COMPLETIONS_MODEL" + suffix),

		// ex: gpt-4-0613
		ChatCompletions: getRequired(prefix + "_CHAT_COMPLETIONS_MODEL" + suffix),

		// ex: embedding
		Embeddings: getRequired(prefix + "_EMBEDDINGS_MODEL" + suffix),

		Azure: azure,

		Cognitive: azopenai.AzureCognitiveSearchChatExtensionConfiguration{
			Endpoint:  to.Ptr(getRequired("COGNITIVE_SEARCH_API_ENDPOINT")),
			IndexName: to.Ptr(getRequired("COGNITIVE_SEARCH_API_INDEX")),
			Key:       to.Ptr(getRequired("COGNITIVE_SEARCH_API_KEY")),
		},
	}

	if tv.Endpoint.URL != "" && !strings.HasSuffix(tv.Endpoint.URL, "/") {
		// (this just makes recording replacement easier)
		tv.Endpoint.URL += "/"
	}

	return tv
}

const fakeEndpoint = "https://fake-recorded-host.microsoft.com/"
const fakeAPIKey = "redacted"
const fakeCognitiveEndpoint = "https://fake-cognitive-endpoint.microsoft.com"
const fakeCognitiveIndexName = "index"

func initEnvVars() {
	if recording.GetRecordMode() == recording.PlaybackMode {
		azureOpenAI.Endpoint = endpoint{
			URL:    fakeEndpoint,
			APIKey: fakeAPIKey,
			Azure:  true,
		}

		azureOpenAICanary.Endpoint = endpoint{
			URL:    fakeEndpoint,
			APIKey: fakeAPIKey,
			Azure:  true,
		}

		azureWhisper = endpoint{
			URL:    fakeEndpoint,
			APIKey: fakeAPIKey,
			Azure:  true,
		}

		azureWhisperModel = "whisper-deployment"

		openAI.Endpoint = endpoint{
			APIKey: fakeAPIKey,
			URL:    fakeEndpoint,
		}

		openAIWhisperModel = "whisper-1"

		azureOpenAICanary.Completions = ""
		azureOpenAICanary.ChatCompletions = "gpt-4"

		azureOpenAI.Completions = "text-davinci-003"
		openAI.Completions = "text-davinci-003"

		azureOpenAI.ChatCompletions = "gpt-4-0613"
		openAI.ChatCompletions = "gpt-4-0613"

		openAI.Embeddings = "text-embedding-ada-002"
		azureOpenAI.Embeddings = "text-embedding-ada-002"

		azureOpenAI.Cognitive = azopenai.AzureCognitiveSearchChatExtensionConfiguration{
			Endpoint:  to.Ptr(fakeCognitiveEndpoint),
			IndexName: to.Ptr(fakeCognitiveIndexName),
			Key:       to.Ptr(fakeAPIKey),
		}
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Failed to load .env file: %s\n", err)
		}

		azureOpenAI = newTestVars("AOAI", false)
		azureOpenAICanary = newTestVars("AOAI", true)
		openAI = newTestVars("OPENAI", false)

		azureWhisper = endpoint{
			URL:    getRequired("AOAI_ENDPOINT_WHISPER"),
			APIKey: getRequired("AOAI_API_KEY_WHISPER"),
			Azure:  true,
		}
		azureWhisperModel = getRequired("AOAI_MODEL_WHISPER")

		openAIWhisper = endpoint{
			URL:    getRequired("OPENAI_ENDPOINT"),
			APIKey: getRequired("OPENAI_API_KEY"),
			Azure:  true,
		}
		openAIWhisperModel = "whisper-1"
	}
}

func newRecordingTransporter(t *testing.T) policy.Transporter {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	err = recording.Start(t, "sdk/ai/azopenai/testdata", nil)
	require.NoError(t, err)

	if recording.GetRecordMode() != recording.PlaybackMode {
		err = recording.AddHeaderRegexSanitizer("Api-Key", fakeAPIKey, "", nil)
		require.NoError(t, err)

		err = recording.AddHeaderRegexSanitizer("User-Agent", "fake-user-agent", ".*", nil)
		require.NoError(t, err)

		// "RequestUri": "https://openai-shared.openai.azure.com/openai/deployments/text-davinci-003/completions?api-version=2023-03-15-preview",
		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(azureOpenAI.Endpoint.URL), nil)
		require.NoError(t, err)

		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(azureOpenAICanary.Endpoint.URL), nil)
		require.NoError(t, err)

		err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(azureWhisper.URL), nil)
		require.NoError(t, err)

		err = recording.AddURISanitizer("/openai/operations/images/00000000-AAAA-BBBB-CCCC-DDDDDDDDDDDD", "/openai/operations/images/[A-Za-z-0-9]+", nil)
		require.NoError(t, err)

		if openAI.Endpoint.URL != "" {
			err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(openAI.Endpoint.URL), nil)
			require.NoError(t, err)
		}

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"endpoint": "%s"`, fakeCognitiveEndpoint),
			fmt.Sprintf(`"endpoint":\s*"%s"`, *azureOpenAI.Cognitive.Endpoint), nil)
		require.NoError(t, err)

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"indexName": "%s"`, fakeCognitiveIndexName),
			fmt.Sprintf(`"indexName":\s*"%s"`, *azureOpenAI.Cognitive.IndexName), nil)
		require.NoError(t, err)

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"key": "%s"`, fakeAPIKey),
			fmt.Sprintf(`"key":\s*"%s"`, *azureOpenAI.Cognitive.Key), nil)
		require.NoError(t, err)
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
	return newTestClient(t, tv.Endpoint)
}

func newOpenAIClientForTest(t *testing.T) *azopenai.Client {
	return newTestClient(t, openAI.Endpoint)
}

// newBogusAzureOpenAIClient creates a client that uses an invalid key, which will cause Azure OpenAI to return
// a failure.
func newBogusAzureOpenAIClient(t *testing.T) *azopenai.Client {
	cred, err := azopenai.NewKeyCredential("bogus-api-key")
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAI.Endpoint.URL, cred, newClientOptionsForTest(t))
	require.NoError(t, err)

	return client
}

// newBogusOpenAIClient creates a client that uses an invalid key, which will cause OpenAI to return
// a failure.
func newBogusOpenAIClient(t *testing.T) *azopenai.Client {
	cred, err := azopenai.NewKeyCredential("bogus-api-key")
	require.NoError(t, err)

	client, err := azopenai.NewClientForOpenAI(openAI.Endpoint.URL, cred, newClientOptionsForTest(t))
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

func getRequired(name string) string {
	v := os.Getenv(name)

	if v == "" {
		panic(fmt.Sprintf("Env variable %s is missing", name))
	}

	return v
}

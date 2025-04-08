// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiextensions"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/openai/openai-go/option"
	"github.com/stretchr/testify/require"
)

const apiVersion = "2025-03-01-preview"

type endpoint struct {
	URL    string
	APIKey string
	Azure  bool
}

type testVars struct {
	Assistants                            endpointWithModel
	ChatCompletions                       endpointWithModel
	ChatCompletionsLegacyFunctions        endpointWithModel
	ChatCompletionsOYD                    endpointWithModel // azure only
	ChatCompletionsRAI                    endpointWithModel // azure only
	ChatCompletionsWithJSONResponseFormat endpointWithModel
	Cognitive                             azopenaiextensions.AzureSearchChatExtensionConfiguration
	Completions                           endpointWithModel
	DallE                                 endpointWithModel
	Embeddings                            endpointWithModel
	Speech                                endpointWithModel
	TextEmbedding3Small                   endpointWithModel
	Vision                                endpointWithModel
	Whisper                               endpointWithModel
}

type endpointWithModel struct {
	Endpoint endpoint
	Model    string
}

// getEnvVariable is recording.GetEnvVariable but it panics if the
// value isn't found, rather than falling back to the playback value.
func getEnvVariable(varName string, playbackValue string) string {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return playbackValue
	}

	val := os.Getenv(varName)

	if val == "" {
		panic(fmt.Sprintf("Missing required environment variable %s", varName))
	}

	return val
}

func getEndpoint(ev string, azure bool) string {
	fakeEP := fakeAzureEndpoint

	if !azure {
		fakeEP = fakeOpenAIEndpoint
	}

	v := getEnvVariable(ev, fakeEP)

	if !strings.HasSuffix(v, "/") {
		// (this just makes recording replacement easier)
		v += "/"
	}

	return v
}

var azureOpenAI = func() testVars {
	servers := struct {
		USEast         endpoint
		USNorthCentral endpoint
		USEast2        endpoint
		SWECentral     endpoint
		OpenAI         endpoint
	}{
		USEast: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_USEAST", true),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USEAST_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		USEast2: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_USEAST2", true),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USEAST2_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		USNorthCentral: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_USNORTHCENTRAL", true),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USNORTHCENTRAL_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		SWECentral: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_SWECENTRAL", true),
			APIKey: getEnvVariable("AOAI_ENDPOINT_SWECENTRAL_API_KEY", fakeAPIKey),
			Azure:  true,
		},
	}

	newTestVarsFn := func() testVars {
		return testVars{
			Assistants: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4o-0806",
			},
			ChatCompletions: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4-0613",
			},
			ChatCompletionsLegacyFunctions: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4-0613",
			},
			ChatCompletionsOYD: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4-0613",
			},
			ChatCompletionsRAI: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4-0613",
			},
			ChatCompletionsWithJSONResponseFormat: endpointWithModel{
				Endpoint: servers.SWECentral,
				Model:    "gpt-4-1106-preview",
			},
			Completions: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-35-turbo-instruct",
			},
			DallE: endpointWithModel{
				Endpoint: servers.SWECentral,
				Model:    "dall-e-3",
			},
			Embeddings: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "text-embedding-ada-002",
			},
			Speech: endpointWithModel{
				Endpoint: servers.SWECentral,
				Model:    "tts",
			},
			TextEmbedding3Small: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "text-embedding-3-small",
			},
			Vision: endpointWithModel{
				Endpoint: servers.SWECentral,
				Model:    "gpt-4-vision-preview",
			},
			Whisper: endpointWithModel{
				Endpoint: servers.USNorthCentral,
				Model:    "whisper",
			},
			Cognitive: azopenaiextensions.AzureSearchChatExtensionConfiguration{
				Parameters: &azopenaiextensions.AzureSearchChatExtensionParameters{
					Endpoint:       to.Ptr(getEnvVariable("COGNITIVE_SEARCH_API_ENDPOINT", fakeCognitiveEndpoint)),
					IndexName:      to.Ptr(getEnvVariable("COGNITIVE_SEARCH_API_INDEX", fakeCognitiveIndexName)),
					Authentication: &azopenaiextensions.OnYourDataSystemAssignedManagedIdentityAuthenticationOptions{},
				},
			},
		}
	}

	azureTestVars := newTestVarsFn()

	if recording.GetRecordMode() == recording.LiveMode {
		// these are for the examples - we don't want to mention regions or anything in them so the
		// env variables have a more friendly naming scheme.
		remaps := map[string]endpointWithModel{
			"CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS": azureTestVars.ChatCompletionsLegacyFunctions,
			"CHAT_COMPLETIONS_RAI":                    azureTestVars.ChatCompletionsRAI,
			"CHAT_COMPLETIONS":                        azureTestVars.ChatCompletions,
			"COMPLETIONS":                             azureTestVars.Completions,
			"DALLE":                                   azureTestVars.DallE,
			"EMBEDDINGS":                              azureTestVars.Embeddings,
			// these resources are oversubscribed and occasionally fail in live testing.
			// "VISION":                                  azureTestVars.Vision,
			// "WHISPER": azureTestVars.Whisper,
		}

		for area, epm := range remaps {
			os.Setenv("AOAI_"+area+"_ENDPOINT", epm.Endpoint.URL)
			os.Setenv("AOAI_"+area+"_API_KEY", epm.Endpoint.APIKey)
			os.Setenv("AOAI_"+area+"_MODEL", epm.Model)
		}
	}

	return azureTestVars
}()

type stainlessTestClientOptions struct {
	UseAPIKey bool
}

func getRecordingOptions(t *testing.T) *recording.RecordingOptions {
	var port int
	val := os.Getenv("PROXY_PORT")

	if len(val) > 0 {
		parsedPort, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			panic(fmt.Sprintf("Invalid proxy port %s", val))
		}
		port = int(parsedPort)
	} else {
		port = os.Getpid()%10000 + 20000
	}
	return &recording.RecordingOptions{
		UseHTTPS:     true,
		ProxyPort:    int(port),
		TestInstance: t,
	}
}

func newStainlessTestClient(t *testing.T, ep endpoint) openai.Client {
	return newStainlessTestClientWithOptions(t, ep, nil)
}

const fakeAzureEndpoint = "https://Sanitized.openai.azure.com/"
const fakeOpenAIEndpoint = "https://Sanitized.openai.com/v1"
const fakeAPIKey = "redacted"
const fakeCognitiveEndpoint = "https://Sanitized.openai.azure.com"
const fakeCognitiveIndexName = "index"

// newRecordingTransporter sets up our recording policy to sanitize endpoints and any parts of the response that might
// involve UUIDs that would make the response/request inconsistent.
func newRecordingTransporter(t *testing.T) policy.Transporter {
	defaultOptions := getRecordingOptions(t)
	t.Logf("Using test proxy on port %d", defaultOptions.ProxyPort)

	transport, err := recording.NewRecordingHTTPClient(t, defaultOptions)
	require.NoError(t, err)

	err = recording.Start(t, RecordingDirectory, defaultOptions)
	require.NoError(t, err)

	if recording.GetRecordMode() != recording.PlaybackMode {
		err = recording.AddHeaderRegexSanitizer("Api-Key", fakeAPIKey, "", defaultOptions)
		require.NoError(t, err)

		err = recording.AddHeaderRegexSanitizer("User-Agent", "fake-user-agent", "", defaultOptions)
		require.NoError(t, err)

		err = recording.AddURISanitizer("/openai/operations/images/00000000-AAAA-BBBB-CCCC-DDDDDDDDDDDD", "/openai/operations/images/[A-Za-z-0-9]+", defaultOptions)
		require.NoError(t, err)

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"endpoint": "%s"`, fakeCognitiveEndpoint),
			fmt.Sprintf(`"endpoint":\s*"%s"`, *azureOpenAI.Cognitive.Parameters.Endpoint), defaultOptions)
		require.NoError(t, err)

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"index_name": "%s"`, fakeCognitiveIndexName),
			`"index_name":\s*".+?"`, defaultOptions)
		require.NoError(t, err)
	}

	t.Cleanup(func() {
		err := recording.Stop(t, defaultOptions)
		require.NoError(t, err)
	})

	return transport
}

type recordingRoundTripper struct {
	transport policy.Transporter
}

func (d *recordingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return d.transport.Do(req)
}

func newStainlessTestClientWithOptions(t *testing.T, ep endpoint, options *stainlessTestClientOptions) openai.Client {
	var client *http.Client
	if recording.GetRecordMode() == recording.LiveMode {
		client = &http.Client{}
	} else {
		transport := newRecordingTransporter(t)
		client = &http.Client{
			Transport: &recordingRoundTripper{transport: transport},
		}
	}

	if options != nil && options.UseAPIKey {
		return openai.NewClient(
			azure.WithEndpoint(ep.URL, apiVersion),
			azure.WithAPIKey(ep.APIKey),
			option.WithHTTPClient(client),
		)
	}

	tokenCredential, err := credential.New(nil)
	require.NoError(t, err)

	return openai.NewClient(
		azure.WithEndpoint(ep.URL, apiVersion),
		azure.WithTokenCredential(tokenCredential),
		option.WithHTTPClient(client),
	)
}

func newStainlessChatCompletionService(t *testing.T, ep endpoint) openai.ChatCompletionService {
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skip("Skipping tests in playback mode")
	}

	tokenCredential, err := credential.New(nil)
	require.NoError(t, err)
	return openai.NewChatCompletionService(azure.WithEndpoint(ep.URL, apiVersion),
		azure.WithTokenCredential(tokenCredential),
	)
}

func skipNowIfThrottled(t *testing.T, err error) {
	if respErr := (*azcore.ResponseError)(nil); errors.As(err, &respErr) && respErr.StatusCode == http.StatusTooManyRequests {
		t.Skipf("OpenAI resource overloaded, skipping this test")
	}
}

// customRequireNoError checks the error but allows throttling errors to account for resources that are
// constrained.
func customRequireNoError(t *testing.T, err error) {
	if err == nil {
		return
	}

	var respErr *openai.Error

	switch {
	case errors.As(err, &respErr) && respErr.StatusCode == http.StatusTooManyRequests:
		t.Skip("Skipping test because of throttling (http.StatusTooManyRequests)")
		return
	// If you're using OYD, then the response error (from Azure OpenAI) will be a 400, but the underlying text will mention
	// that it's 429'd.
	// 	  "code": 400,
	// 	  "message": "Server responded with status 429. Error message: {'error': {'code': '429', 'message': 'Rate limit is exceeded. Try again in 1 seconds.'}}"
	case errors.As(err, &respErr) && respErr.StatusCode == http.StatusBadRequest && strings.Contains(err.Error(), "Rate limit is exceeded"):
		t.Skip("Skipping test because of throttling in OYD resource")
		return
	case errors.Is(err, context.DeadlineExceeded):
		t.Skip("Skipping test because of throttling (DeadlineExceeded)")
		return
	}

	require.NoError(t, err)
}

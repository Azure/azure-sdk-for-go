// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/option"
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
	Cognitive                             azopenai.AzureSearchChatExtensionConfiguration
	Completions                           endpointWithModel
	DallE                                 endpointWithModel
	Embeddings                            endpointWithModel
	Speech                                endpointWithModel
	TextEmbedding3Small                   endpointWithModel
	Vision                                endpointWithModel
	Whisper                               endpointWithModel
	Reasoning                             endpointWithModel
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

func getEndpoint(ev string) string {
	v := getEnvVariable(ev, fakeAzureEndpoint)

	if !strings.HasSuffix(v, "/") {
		// (this just makes recording replacement easier)
		v += "/"
	}

	return v
}

var azureOpenAI = func() testVars {
	if recording.GetRecordMode() != recording.PlaybackMode {
		// check if some of the variables are already in the environment - this'll happen with
		// live testing.
		if os.Getenv("COGNITIVE_SEARCH_API_ENDPOINT") == "" {
			if err := godotenv.Load(); err != nil {
				panic(fmt.Errorf("Failed to load .env file: %w", err))
			} else {
				log.Printf(".env file loaded")
			}
		} else {
			log.Printf(".env file loading skipped - variables already in environment")
		}
	} else {
		log.Printf(".env file loading skipped, since we're in playback mode")
	}

	servers := struct {
		USEast         endpoint
		USNorthCentral endpoint
		USEast2        endpoint
		SWECentral     endpoint
		OpenAI         endpoint
	}{
		USEast: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_USEAST"),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USEAST_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		USEast2: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_USEAST2"),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USEAST2_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		USNorthCentral: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_USNORTHCENTRAL"),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USNORTHCENTRAL_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		SWECentral: endpoint{
			URL:    getEndpoint("AOAI_ENDPOINT_SWECENTRAL"),
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
				Model:    "gpt-4",
			},
			ChatCompletionsLegacyFunctions: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4",
			},
			ChatCompletionsOYD: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4",
			},
			ChatCompletionsRAI: endpointWithModel{
				Endpoint: servers.USEast,
				Model:    "gpt-4",
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
			Reasoning: endpointWithModel{
				Endpoint: servers.SWECentral,
				Model:    "o3-2025-04-16",
			},
			Cognitive: azopenai.AzureSearchChatExtensionConfiguration{
				Parameters: &azopenai.AzureSearchChatExtensionParameters{
					Endpoint:       to.Ptr(getEnvVariable("COGNITIVE_SEARCH_API_ENDPOINT", fakeCognitiveEndpoint)),
					IndexName:      to.Ptr(getEnvVariable("COGNITIVE_SEARCH_API_INDEX", fakeCognitiveIndexName)),
					Authentication: &azopenai.OnYourDataSystemAssignedManagedIdentityAuthenticationOptions{},
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
			_ = os.Setenv("AOAI_"+area+"_ENDPOINT", epm.Endpoint.URL)
			_ = os.Setenv("AOAI_"+area+"_API_KEY", epm.Endpoint.APIKey)
			_ = os.Setenv("AOAI_"+area+"_MODEL", epm.Model)
		}
	}

	return azureTestVars
}()

type stainlessTestClientOptions struct {
	UseAPIKey bool
	// UseV1Endpoint controls which endpoint style we use for the created client.
	//    - If true, we use the /openai/v1 style endpoint. See the [api-doc] for what parts of the OpenAI are implemented.
	//    - If false, we use the older style Azure OpenAI endpoints, which contain a deployment in the URL
	//
	// [api-doc]: https://github.com/MicrosoftDocs/azure-ai-docs/blob/main/articles/ai-foundry/openai/latest.md
	UseV1Endpoint bool
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

func newStainlessTestClientWithAzureURL(t *testing.T, ep endpoint) openai.Client {
	return newStainlessTestClientWithOptions(t, ep, &stainlessTestClientOptions{
		UseV1Endpoint: false,
	})
}

func newStainlessTestClientWithV1URL(t *testing.T, ep endpoint) openai.Client {
	return newStainlessTestClientWithOptions(t, ep, &stainlessTestClientOptions{
		UseV1Endpoint: true,
	})
}

const fakeAzureEndpoint = "https://Sanitized.openai.azure.com/"
const fakeAPIKey = "redacted"
const fakeCognitiveEndpoint = "https://Sanitized.openai.azure.com"
const fakeCognitiveIndexName = "index"

func configureTestProxy(options recording.RecordingOptions) error {
	if err := recording.SetDefaultMatcher(nil, &recording.SetDefaultMatcherOptions{
		RecordingOptions: options,
		ExcludedHeaders: []string{
			"X-Stainless-Arch",
			"X-Stainless-Lang",
			"X-Stainless-Os",
			"X-Stainless-Package-Version",
			"X-Stainless-Retry-Count",
			"X-Stainless-Runtime",
			"X-Stainless-Runtime-Version",
		},
	}); err != nil {
		return err
	}

	if err := recording.AddHeaderRegexSanitizer("Api-Key", fakeAPIKey, "", &options); err != nil {
		return err
	}

	if err := recording.AddHeaderRegexSanitizer("User-Agent", "fake-user-agent", "", &options); err != nil {
		return err
	}

	if err := recording.AddURISanitizer("/openai/operations/images/00000000-AAAA-BBBB-CCCC-DDDDDDDDDDDD", "/openai/operations/images/[A-Za-z-0-9]+", &options); err != nil {
		return err
	}

	if err := recording.AddGeneralRegexSanitizer(
		fmt.Sprintf(`"endpoint": "%s"`, fakeCognitiveEndpoint),
		`"endpoint":\s*"[^"]+"`, &options); err != nil {
		return err
	}

	if err := recording.AddGeneralRegexSanitizer(
		fmt.Sprintf(`"index_name": "%s"`, fakeCognitiveIndexName),
		`"index_name":\s*"[^"]+"`, &options); err != nil {
		return err
	}

	return nil
}

// newRecordingTransporter sets up our recording policy to sanitize endpoints and any parts of the response that might
// involve UUIDs that would make the response/request inconsistent.
func newRecordingTransporter(t *testing.T) policy.Transporter {
	defaultOptions := getRecordingOptions(t)
	t.Logf("Using test proxy on port %d", defaultOptions.ProxyPort)

	transport, err := recording.NewRecordingHTTPClient(t, defaultOptions)
	require.NoError(t, err)

	// if we're creating more than one client in a test (for instance, TestClient_GetAudioSpeech!)
	// then we don't want to start or stop recording again.
	if recording.GetRecordingId(t) == "" {
		err = recording.Start(t, RecordingDirectory, defaultOptions)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := recording.Stop(t, defaultOptions)
			require.NoError(t, err)
		})
	}

	return transport
}

type recordingRoundTripper struct {
	transport policy.Transporter
}

func (d *recordingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return d.transport.Do(req)
}

func newStainlessTestClientWithOptions(t *testing.T, ep endpoint, options *stainlessTestClientOptions) openai.Client {
	if options == nil {
		options = &stainlessTestClientOptions{}
	}

	var client *http.Client
	if recording.GetRecordMode() == recording.LiveMode {
		client = &http.Client{}
	} else {
		transport := newRecordingTransporter(t)
		client = &http.Client{
			Transport: &recordingRoundTripper{transport: transport},
		}
	}

	endpointOption := azure.WithEndpoint(ep.URL, apiVersion)

	if options.UseV1Endpoint {
		endpointOption = option.WithBaseURL(ep.URL + "openai/v1")
	}

	if options.UseAPIKey {
		return openai.NewClient(
			endpointOption,
			azure.WithAPIKey(ep.APIKey),
			option.WithHTTPClient(client),
		)
	}

	tokenCredential, err := credential.New(nil)
	require.NoError(t, err)

	return openai.NewClient(
		endpointOption,
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
	t.Helper()

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

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiextensions"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
	"github.com/stretchr/testify/require"
)

const apiVersion = "2024-07-01-preview"

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
func getEnvVariable(varName string) string {
	val := os.Getenv(varName)

	if val == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic(fmt.Sprintf("Missing required environment variable %s", varName))
		}
	}

	return val
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
			URL:    getEnvVariable("AOAI_ENDPOINT_USEAST"),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USEAST_API_KEY"),
			Azure:  true,
		},
		USEast2: endpoint{
			URL:    getEnvVariable("AOAI_ENDPOINT_USEAST2"),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USEAST2_API_KEY"),
			Azure:  true,
		},
		USNorthCentral: endpoint{
			URL:    getEnvVariable("AOAI_ENDPOINT_USNORTHCENTRAL"),
			APIKey: getEnvVariable("AOAI_ENDPOINT_USNORTHCENTRAL_API_KEY"),
			Azure:  true,
		},
		SWECentral: endpoint{
			URL:    getEnvVariable("AOAI_ENDPOINT_SWECENTRAL"),
			APIKey: getEnvVariable("AOAI_ENDPOINT_SWECENTRAL_API_KEY"),
			Azure:  true,
		},
	}

	newTestVarsFn := func() testVars {
		return testVars{
			Assistants: endpointWithModel{
				Endpoint: servers.SWECentral,
				Model:    "gpt-4-1106-preview",
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
					Endpoint:       to.Ptr(getEnvVariable("COGNITIVE_SEARCH_API_ENDPOINT")),
					IndexName:      to.Ptr(getEnvVariable("COGNITIVE_SEARCH_API_INDEX")),
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

func newStainlessTestClient(t *testing.T, ep endpoint) *openai.Client {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("Skipping tests in playback mode")
		return nil
	}

	tokenCredential, err := credential.New(nil)
	require.NoError(t, err)

	return openai.NewClient(
		azure.WithEndpoint(ep.URL, apiVersion),
		azure.WithTokenCredential(tokenCredential),
	)
}

func newStainlessChatCompletionService(t *testing.T, ep endpoint) *openai.ChatCompletionService {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("Skipping tests in playback mode")
		return nil
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
func customRequireNoError(t *testing.T, err error, throttlingAllowed bool) {
	if err == nil {
		return
	}

	if throttlingAllowed {

		if respErr := (*openai.Error)(nil); errors.As(err, &respErr) && respErr.StatusCode == http.StatusTooManyRequests {
			t.Skip("Skipping test because of throttling (http.StatusTooManyRequests)")
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			t.Skip("Skipping test because of throttling (DeadlineExceeded)")
			return
		}
	}

	require.NoError(t, err)
}

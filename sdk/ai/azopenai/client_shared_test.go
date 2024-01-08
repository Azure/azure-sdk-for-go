// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	azureOpenAI testVars
	openAI      testVars
)

type endpoint struct {
	URL    string
	APIKey string
	Azure  bool
}

type testVars struct {
	Endpoint                       endpoint
	Completions                    string
	ChatCompletions                string
	ChatCompletionsLegacyFunctions string
	Embeddings                     string
	Cognitive                      azopenai.AzureCognitiveSearchChatExtensionConfiguration
	Whisper                        endpointWithModel
	DallE                          endpointWithModel
	Vision                         endpointWithModel

	ChatCompletionsRAI endpointWithModel // at the moment this is Azure only

	// "own your data" - bringing in Azure resources as part of a chat completions
	// request.
	ChatCompletionsOYD endpointWithModel
}

type endpointWithModel struct {
	Endpoint endpoint
	Model    string
}

type testClientOption func(opt *azopenai.ClientOptions)

func withForgivingRetryOption() testClientOption {
	return func(opt *azopenai.ClientOptions) {
		opt.Retry = policy.RetryOptions{
			MaxRetries: 10,
		}
	}
}

// newTestClient creates a client enabled for HTTP recording, if needed.
// See [newRecordingTransporter] for sanitization code.
func newTestClient(t *testing.T, ep endpoint, options ...testClientOption) *azopenai.Client {
	clientOptions := newClientOptionsForTest(t)

	for _, opt := range options {
		opt(clientOptions)
	}

	if ep.Azure {
		cred := azcore.NewKeyCredential(ep.APIKey)

		client, err := azopenai.NewClientWithKeyCredential(ep.URL, cred, clientOptions)
		require.NoError(t, err)

		return client
	} else {
		if ep.APIKey == "" {
			t.Skipf("OPENAI_API_KEY not defined, skipping OpenAI public endpoint test")
		}

		cred := azcore.NewKeyCredential(ep.APIKey)

		client, err := azopenai.NewClientForOpenAI(ep.URL, cred, clientOptions)
		require.NoError(t, err)

		return client
	}
}

// getEndpointWithModel retrieves details for an endpoint and a model.
// - res - the resource type for a particular endpoint. Ex: "DALLE".
//
// For example, if azure is true we'll load these environment values based on res:
//   - AOAI_DALLE_ENDPOINT
//   - AOAI_DALLE_API_KEY
//   - AOAI_DALLE_MODEL
//
// if azure is false we'll load these environment values based on res:
//   - OPENAI_ENDPOINT
//   - OPENAI_API_KEY
//   - OPENAI_DALLE_MODEL
func getEndpointWithModel(res string, isAzure bool) endpointWithModel {
	var ep endpointWithModel
	if isAzure {
		// during development resources are often shifted between different
		// internal Azure OpenAI resources.
		ep = endpointWithModel{
			Endpoint: endpoint{
				URL:    getRequired("AOAI_" + res + "_ENDPOINT"),
				APIKey: getRequired("AOAI_" + res + "_API_KEY"),
				Azure:  true,
			},
			Model: getRequired("AOAI_" + res + "_MODEL"),
		}
	} else {
		ep = endpointWithModel{
			Endpoint: endpoint{
				URL:    getRequired("OPENAI_ENDPOINT"),
				APIKey: getRequired("OPENAI_API_KEY"),
				Azure:  false,
			},
			Model: getRequired("OPENAI_" + res + "_MODEL"),
		}
	}

	if !strings.HasSuffix(ep.Endpoint.URL, "/") {
		// (this just makes recording replacement easier)
		ep.Endpoint.URL += "/"
	}

	return ep
}

func newTestVars(prefix string) testVars {
	azure := prefix == "AOAI"

	tv := testVars{
		Endpoint: endpoint{
			URL:    getRequired(prefix + "_ENDPOINT"),
			APIKey: getRequired(prefix + "_API_KEY"),
			Azure:  azure,
		},
		Completions:                    getRequired(prefix + "_COMPLETIONS_MODEL"),
		ChatCompletions:                getRequired(prefix + "_CHAT_COMPLETIONS_MODEL"),
		ChatCompletionsLegacyFunctions: getRequired(prefix + "_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS"),
		Embeddings:                     getRequired(prefix + "_EMBEDDINGS_MODEL"),

		Cognitive: azopenai.AzureCognitiveSearchChatExtensionConfiguration{
			Parameters: &azopenai.AzureCognitiveSearchChatExtensionParameters{
				Endpoint:  to.Ptr(getRequired("COGNITIVE_SEARCH_API_ENDPOINT")),
				IndexName: to.Ptr(getRequired("COGNITIVE_SEARCH_API_INDEX")),
				Authentication: &azopenai.OnYourDataAPIKeyAuthenticationOptions{
					Key: to.Ptr(getRequired("COGNITIVE_SEARCH_API_KEY")),
				},
			},
		},

		DallE:   getEndpointWithModel("DALLE", azure),
		Whisper: getEndpointWithModel("WHISPER", azure),
		Vision:  getEndpointWithModel("VISION", azure),
	}

	if azure {
		tv.ChatCompletionsRAI = getEndpointWithModel("CHAT_COMPLETIONS_RAI", azure)
		tv.ChatCompletionsOYD = getEndpointWithModel("OYD", azure)
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
		// Setup our variables so our requests are consistent with what we recorded.
		// Endpoints are sanitized using the recording policy
		azureOpenAI.Endpoint = endpoint{
			URL:    fakeEndpoint,
			APIKey: fakeAPIKey,
			Azure:  true,
		}

		azureOpenAI.Whisper = endpointWithModel{
			Endpoint: azureOpenAI.Endpoint,
			Model:    "whisper-deployment",
		}

		azureOpenAI.ChatCompletionsRAI = endpointWithModel{
			Endpoint: azureOpenAI.Endpoint,
			Model:    "gpt-4",
		}

		azureOpenAI.ChatCompletionsOYD = endpointWithModel{
			Endpoint: azureOpenAI.Endpoint,
			Model:    "gpt-4",
		}

		azureOpenAI.DallE = endpointWithModel{
			Endpoint: azureOpenAI.Endpoint,
			Model:    "dall-e-3",
		}

		azureOpenAI.Vision = endpointWithModel{
			Endpoint: azureOpenAI.Endpoint,
			Model:    "gpt-4-vision-preview",
		}

		openAI.Endpoint = endpoint{
			APIKey: fakeAPIKey,
			URL:    fakeEndpoint,
		}

		openAI.Whisper = endpointWithModel{
			Endpoint: endpoint{
				APIKey: fakeAPIKey,
				URL:    fakeEndpoint,
			},
			Model: "whisper-1",
		}

		openAI.DallE = endpointWithModel{
			Endpoint: openAI.Endpoint,
			Model:    "dall-e-3",
		}

		openAI.Vision = azureOpenAI.Vision

		azureOpenAI.Completions = "gpt-35-turbo"
		openAI.Completions = "gpt-35-turbo"

		azureOpenAI.ChatCompletions = "gpt-35-turbo-0613"
		azureOpenAI.ChatCompletionsLegacyFunctions = "gpt-4-0613"
		openAI.ChatCompletions = "gpt-4-0613"
		openAI.ChatCompletionsLegacyFunctions = "gpt-4-0613"

		openAI.Embeddings = "text-embedding-ada-002"
		azureOpenAI.Embeddings = "text-embedding-ada-002"

		azureOpenAI.Cognitive = azopenai.AzureCognitiveSearchChatExtensionConfiguration{
			Parameters: &azopenai.AzureCognitiveSearchChatExtensionParameters{
				Endpoint:  to.Ptr(fakeCognitiveEndpoint),
				IndexName: to.Ptr(fakeCognitiveIndexName),
				Authentication: &azopenai.OnYourDataAPIKeyAuthenticationOptions{
					Key: to.Ptr(fakeAPIKey),
				},
			},
		}
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Failed to load .env file: %s\n", err)
		}

		azureOpenAI = newTestVars("AOAI")
		openAI = newTestVars("OPENAI")
	}
}

type MultipartRecordingPolicy struct {
}

func (mrp MultipartRecordingPolicy) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}

// newRecordingTransporter sets up our recording policy to sanitize endpoints and any parts of the response that might
// involve UUIDs that would make the response/request inconsistent.
func newRecordingTransporter(t *testing.T) policy.Transporter {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	err = recording.Start(t, RecordingDirectory, nil)
	require.NoError(t, err)

	if recording.GetRecordMode() != recording.PlaybackMode {
		err = recording.AddHeaderRegexSanitizer("Api-Key", fakeAPIKey, "", nil)
		require.NoError(t, err)

		err = recording.AddHeaderRegexSanitizer("User-Agent", "fake-user-agent", "", nil)
		require.NoError(t, err)

		endpoints := []string{
			azureOpenAI.Endpoint.URL,
			azureOpenAI.ChatCompletionsRAI.Endpoint.URL,
			azureOpenAI.Whisper.Endpoint.URL,
			azureOpenAI.DallE.Endpoint.URL,
			azureOpenAI.Vision.Endpoint.URL,
			azureOpenAI.ChatCompletionsOYD.Endpoint.URL,
			azureOpenAI.ChatCompletionsRAI.Endpoint.URL,
		}

		for _, ep := range endpoints {
			err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(ep), nil)
			require.NoError(t, err)
		}

		err = recording.AddURISanitizer("/openai/operations/images/00000000-AAAA-BBBB-CCCC-DDDDDDDDDDDD", "/openai/operations/images/[A-Za-z-0-9]+", nil)
		require.NoError(t, err)

		if openAI.Endpoint.URL != "" {
			err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(openAI.Endpoint.URL), nil)
			require.NoError(t, err)
		}

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"endpoint": "%s"`, fakeCognitiveEndpoint),
			fmt.Sprintf(`"endpoint":\s*"%s"`, *azureOpenAI.Cognitive.Parameters.Endpoint), nil)
		require.NoError(t, err)

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"indexName": "%s"`, fakeCognitiveIndexName),
			fmt.Sprintf(`"indexName":\s*"%s"`, *azureOpenAI.Cognitive.Parameters.IndexName), nil)
		require.NoError(t, err)

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"key": "%s"`, fakeAPIKey),
			fmt.Sprintf(`"key":\s*"%s"`, *azureOpenAI.Cognitive.Parameters.Authentication.(*azopenai.OnYourDataAPIKeyAuthenticationOptions).Key), nil)
		require.NoError(t, err)
	}

	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})

	return transport
}

// newClientOptionsForTest creates options that enable a few optional things:
//   - If we're recording then it injects the recording transporter
//   - If `SSLKEYLOGFILE` is set in the environment it'll automatically setup
//     the HTTP policy so it writes the keylog for that client to a file. You can
//     use this with WireShark to decrypt and view a network trace.
func newClientOptionsForTest(t *testing.T) *azopenai.ClientOptions {
	co := &azopenai.ClientOptions{}

	// Useful when debugging responses.
	co.Logging.IncludeBody = true

	if recording.GetRecordMode() == recording.LiveMode {
		keyLogPath := os.Getenv("SSLKEYLOGFILE")

		if keyLogPath == "" {
			return co
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
		co.PerRetryPolicies = append(co.PerRetryPolicies, &mimeTypeRecordingPolicy{})
		co.Transport = newRecordingTransporter(t)
	}

	return co
}

func newAzureOpenAIClientForTest(t *testing.T, tv testVars) *azopenai.Client {
	return newTestClient(t, tv.Endpoint)
}

func newOpenAIClientForTest(t *testing.T) *azopenai.Client {
	return newTestClient(t, openAI.Endpoint)
}

// newBogusAzureOpenAIClient creates a client that uses an invalid key, which will cause Azure OpenAI to return
// a failure.
func newBogusAzureOpenAIClient(t *testing.T) *azopenai.Client {
	cred := azcore.NewKeyCredential("bogus-api-key")

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAI.Endpoint.URL, cred, newClientOptionsForTest(t))
	require.NoError(t, err)

	return client
}

// newBogusOpenAIClient creates a client that uses an invalid key, which will cause OpenAI to return
// a failure.
func newBogusOpenAIClient(t *testing.T) *azopenai.Client {
	cred := azcore.NewKeyCredential("bogus-api-key")

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

func skipNowIfThrottled(t *testing.T, err error) {
	if respErr := (*azcore.ResponseError)(nil); errors.As(err, &respErr) && respErr.StatusCode == http.StatusTooManyRequests {
		t.Skipf("OpenAI resource overloaded, skipping this test")
	}
}

type mimeTypeRecordingPolicy struct{}

// Do changes out the boundary for a multipart message. This makes it simpler to write
// recordings.
func (mrp *mimeTypeRecordingPolicy) Do(req *policy.Request) (*http.Response, error) {
	if recording.GetRecordMode() == recording.LiveMode {
		// this is strictly to make the IDs in the multipart body stable for test recordings.
		return req.Next()
	}

	// we'll fix up the multipart to make it more predictable for test recordings.
	//    Content-Type: multipart/form-data; boundary=787c880ce3dd11f9b6384d625c399c8490fc8989ceb6b7d208ec7426c12e
	mediaType, params, err := mime.ParseMediaType(req.Raw().Header[http.CanonicalHeaderKey("Content-type")][0])

	if err != nil || mediaType != "multipart/form-data" {
		// we'll just assume our policy doesn't apply here.
		return req.Next()
	}

	origBoundary := params["boundary"]

	if origBoundary == "" {
		return nil, errors.New("Invalid use of this policy - no boundary was passed as part of the multipart mime type")
	}

	params["boundary"] = "boundary-for-recordings"

	// now let's update the body itself - we'll just do a simple string replacement. The entire purpose of the boundary string is to provide a
	// separator, which is distinct from the content.
	body := req.Body()
	defer body.Close()

	origBody, err := io.ReadAll(body)

	if err != nil {
		return nil, err
	}

	newBody := bytes.ReplaceAll(origBody, []byte(origBoundary), []byte("boundary-for-recordings"))

	if err := req.SetBody(streaming.NopCloser(bytes.NewReader(newBody)), mime.FormatMediaType(mediaType, params)); err != nil {
		return nil, err
	}

	return req.Next()
}

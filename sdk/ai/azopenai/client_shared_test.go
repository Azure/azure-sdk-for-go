// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
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

type endpoint struct {
	URL    string
	APIKey string
	Azure  bool
}

type testVars struct {
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
}

type endpointWithModel struct {
	Endpoint endpoint
	Model    string
}

func ifAzure[T string | endpoint](azure bool, forAzure T, forOpenAI T) T {
	if azure {
		return forAzure
	}
	return forOpenAI
}

var azureOpenAI, openAI, servers = func() (testVars, testVars, []string) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Failed to load .env file: %s\n", err)
		}
	}

	servers := struct {
		USEast         endpoint
		USNorthCentral endpoint
		USEast2        endpoint
		SWECentral     endpoint
		OpenAI         endpoint
	}{
		OpenAI: endpoint{
			URL:    getEndpoint("OPENAI_ENDPOINT"), // ex: https://api.openai.com/v1/
			APIKey: recording.GetEnvVariable("OPENAI_API_KEY", fakeAPIKey),
			Azure:  false,
		},
		USEast: endpoint{
			URL:    getEndpoint("ENDPOINT_USEAST"),
			APIKey: recording.GetEnvVariable("ENDPOINT_USEAST_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		USEast2: endpoint{
			URL:    getEndpoint("ENDPOINT_USEAST2"),
			APIKey: recording.GetEnvVariable("ENDPOINT_USEAST2_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		USNorthCentral: endpoint{
			URL:    getEndpoint("ENDPOINT_USNORTHCENTRAL"),
			APIKey: recording.GetEnvVariable("ENDPOINT_USNORTHCENTRAL_API_KEY", fakeAPIKey),
			Azure:  true,
		},
		SWECentral: endpoint{
			URL:    getEndpoint("ENDPOINT_SWECENTRAL"),
			APIKey: recording.GetEnvVariable("ENDPOINT_SWECENTRAL_API_KEY", fakeAPIKey),
			Azure:  true,
		},
	}

	// used when we setup the recording policy
	endpoints := []string{
		servers.OpenAI.URL,
		servers.USEast.URL,
		servers.USEast2.URL,
		servers.USNorthCentral.URL,
		servers.SWECentral.URL,
	}

	newTestVarsFn := func(azure bool) testVars {
		return testVars{
			ChatCompletions: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "gpt-4-0613", "gpt-4-0613"),
			},
			ChatCompletionsLegacyFunctions: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "gpt-4-0613", "gpt-4-0613"),
			},
			ChatCompletionsOYD: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "gpt-4-0613", ""), // azure only
			},
			ChatCompletionsRAI: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "gpt-4-0613", ""), // azure only
			},
			ChatCompletionsWithJSONResponseFormat: endpointWithModel{
				Endpoint: ifAzure(azure, servers.SWECentral, servers.OpenAI),
				Model:    ifAzure(azure, "gpt-4-1106-preview", "gpt-3.5-turbo-1106"),
			},
			Completions: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "gpt-35-turbo-instruct", "gpt-3.5-turbo-instruct"),
			},
			DallE: endpointWithModel{
				Endpoint: ifAzure(azure, servers.SWECentral, servers.OpenAI),
				Model:    ifAzure(azure, "dall-e-3", "dall-e-3"),
			},
			Embeddings: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "text-embedding-ada-002", "text-embedding-ada-002"),
			},
			Speech: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "tts-1", "tts-1"),
			},
			TextEmbedding3Small: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast, servers.OpenAI),
				Model:    ifAzure(azure, "text-embedding-3-small", "text-embedding-3-small"),
			},
			Vision: endpointWithModel{
				Endpoint: ifAzure(azure, servers.SWECentral, servers.OpenAI),
				Model:    ifAzure(azure, "gpt-4-vision-preview", "gpt-4-vision-preview"),
			},
			Whisper: endpointWithModel{
				Endpoint: ifAzure(azure, servers.USEast2, servers.OpenAI),
				Model:    ifAzure(azure, "whisper-deployment", "whisper-1"),
			},
			Cognitive: azopenai.AzureSearchChatExtensionConfiguration{
				Parameters: &azopenai.AzureSearchChatExtensionParameters{
					Endpoint:  to.Ptr(recording.GetEnvVariable("COGNITIVE_SEARCH_API_ENDPOINT", fakeCognitiveEndpoint)),
					IndexName: to.Ptr(recording.GetEnvVariable("COGNITIVE_SEARCH_API_INDEX", fakeCognitiveIndexName)),
					Authentication: &azopenai.OnYourDataAPIKeyAuthenticationOptions{
						Key: to.Ptr(recording.GetEnvVariable("COGNITIVE_SEARCH_API_KEY", fakeAPIKey)),
					},
				},
			},
		}
	}

	return newTestVarsFn(true), newTestVarsFn(false), endpoints
}()

type testClientOption func(opt *azopenai.ClientOptions)

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

const fakeEndpoint = "https://fake-recorded-host.microsoft.com/"
const fakeAPIKey = "redacted"
const fakeCognitiveEndpoint = "https://fake-cognitive-endpoint.microsoft.com"
const fakeCognitiveIndexName = "index"

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

		for _, ep := range servers {
			err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(ep), nil)
			require.NoError(t, err)
		}

		err = recording.AddURISanitizer("/openai/operations/images/00000000-AAAA-BBBB-CCCC-DDDDDDDDDDDD", "/openai/operations/images/[A-Za-z-0-9]+", nil)
		require.NoError(t, err)

		// there's only one OpenAI endpoint
		if openAI.ChatCompletions.Endpoint.URL != "" {
			err = recording.AddURISanitizer(fakeEndpoint, regexp.QuoteMeta(openAI.ChatCompletions.Endpoint.URL), nil)
			require.NoError(t, err)
		}

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"endpoint": "%s"`, fakeCognitiveEndpoint),
			fmt.Sprintf(`"endpoint":\s*"%s"`, *azureOpenAI.Cognitive.Parameters.Endpoint), nil)
		require.NoError(t, err)

		err = recording.AddGeneralRegexSanitizer(
			fmt.Sprintf(`"index_name": "%s"`, fakeCognitiveIndexName),
			fmt.Sprintf(`"index_name":\s*"%s"`, *azureOpenAI.Cognitive.Parameters.IndexName), nil)
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

// newBogusAzureOpenAIClient creates a client that uses an invalid key, which will cause Azure OpenAI to return
// a failure.
func newBogusAzureOpenAIClient(t *testing.T) *azopenai.Client {
	cred := azcore.NewKeyCredential("bogus-api-key")

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAI.Completions.Endpoint.URL, cred, newClientOptionsForTest(t))
	require.NoError(t, err)

	return client
}

// newBogusOpenAIClient creates a client that uses an invalid key, which will cause OpenAI to return
// a failure.
func newBogusOpenAIClient(t *testing.T) *azopenai.Client {
	cred := azcore.NewKeyCredential("bogus-api-key")

	client, err := azopenai.NewClientForOpenAI(openAI.Completions.Endpoint.URL, cred, newClientOptionsForTest(t))
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

func getEndpoint(ev string) string {
	v := recording.GetEnvVariable(ev, fakeEndpoint)

	if !strings.HasSuffix(v, "/") {
		// (this just makes recording replacement easier)
		v += "/"
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

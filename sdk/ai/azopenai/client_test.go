//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

func TestClient_OpenAI_InvalidModel(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode || testing.Short() {
		t.Skip()
	}

	chatClient := newTestClient(t, openAI.ChatCompletions.Endpoint)

	_, err := chatClient.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestSystemMessage{
				Content: to.Ptr("hello"),
			},
		},
		DeploymentName: to.Ptr("non-existent-model"),
	}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
	require.Contains(t, respErr.Error(), "The model `non-existent-model` does not exist")
}

func TestClient_EmptyOptionsChecking(t *testing.T) {
	// I'm ignoring these in the methods so if they ever actually have relevant config
	// you should revisit them and make sure they're used properly.
	emptyOptionsType := []any{
		azopenai.GenerateSpeechFromTextOptions{},
		azopenai.GetChatCompletionsOptions{},
		azopenai.GetCompletionsOptions{},
		azopenai.GetEmbeddingsOptions{},
		azopenai.GetImageGenerationsOptions{},
	}

	for _, v := range emptyOptionsType {
		fields := reflect.VisibleFields(reflect.TypeOf(v))
		require.Emptyf(t, fields, "%T is ignored in our function signatures because it's empty", v)
	}
}

func TestClient_InsecureHTTPAllowed(t *testing.T) {
	const fakeID = "fake-id"

	hf := func(resp http.ResponseWriter, req *http.Request) {
		// _just_ enough of a response to prove we made it through the pipeline.
		_, err := resp.Write([]byte(fmt.Sprintf("{ \"id\": \"%s\" }", fakeID)))
		require.NoError(t, err)
	}

	urlCh := make(chan string)

	go func() {
		// start an HTTP service
		server := httptest.NewServer(http.HandlerFunc(hf))
		urlCh <- server.URL
		t.Cleanup(server.Close)
	}()

	url := <-urlCh
	t.Logf(url)

	t.Run("DefaultsToHTTPSOnly", func(t *testing.T) {
		client, err := azopenai.NewClientForOpenAI(url, azcore.NewKeyCredential("fake-key"), nil)
		require.NoError(t, err)

		resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{},
		}, nil)
		require.Empty(t, resp)
		require.EqualError(t, err, "authenticated requests are not permitted for non TLS protected (https) endpoints")

		client, err = azopenai.NewClientWithKeyCredential(url, azcore.NewKeyCredential("fake-key"), nil)
		require.NoError(t, err)

		resp, err = client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{},
		}, nil)
		require.Empty(t, resp)
		require.EqualError(t, err, "authenticated requests are not permitted for non TLS protected (https) endpoints")

		fakeCred := &credential.Fake{}

		client, err = azopenai.NewClient(url, fakeCred, nil)
		require.NoError(t, err)

		resp, err = client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{},
		}, nil)
		require.Empty(t, resp)
		require.EqualError(t, err, "authenticated requests are not permitted for non TLS protected (https) endpoints")
	})

	t.Run("InsecureAllowCredentialWithHTTP", func(t *testing.T) {
		clientOptions := &azopenai.ClientOptions{
			ClientOptions: policy.ClientOptions{
				InsecureAllowCredentialWithHTTP: true,
			},
		}

		client, err := azopenai.NewClientForOpenAI(url, azcore.NewKeyCredential("fake-key"), clientOptions)
		require.NoError(t, err)

		resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{},
		}, nil)
		require.NoError(t, err)
		require.Equal(t, fakeID, *resp.ID)

		client, err = azopenai.NewClientWithKeyCredential(url, azcore.NewKeyCredential("fake-key"), clientOptions)
		require.NoError(t, err)

		resp, err = client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{},
		}, nil)
		require.NoError(t, err)
		require.Equal(t, fakeID, *resp.ID)

		fakeCred := &credential.Fake{}

		client, err = azopenai.NewClient(url, fakeCred, clientOptions)
		require.NoError(t, err)

		resp, err = client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{},
		}, nil)
		require.NoError(t, err)
		require.Equal(t, fakeID, *resp.ID)
	})
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
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

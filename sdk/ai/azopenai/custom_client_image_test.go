//go:build go1.21
// +build go1.21

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/require"
)

func TestImageGeneration_AzureOpenAI(t *testing.T) {
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skipf("Ignoring poller-based test")
	}

	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.DallE.Endpoint)
	// testImageGeneration(t, client, azureOpenAI.DallE.Model, azopenai.ImageGenerationResponseFormatURL, true)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resp, err := client.Images.Generate(ctx, openai.ImageGenerateParams{
		// saw this prompt in a thread about trying to _prevent_ Dall-E3 from rewriting your
		// propmt. When this is revised you'll see the text in the
		Prompt:         "acrylic painting of a sunflower with bees",
		Size:           openai.ImageGenerateParamsSize1024x1792,
		ResponseFormat: openai.ImageGenerateParamsResponseFormatURL,
		Model:          openai.ImageModel(azureOpenAI.DallE.Model),
	})
	customRequireNoError(t, err)

	if recording.GetRecordMode() == recording.LiveMode {
		headResp, err := http.DefaultClient.Head(resp.Data[0].URL)
		require.NoError(t, err)

		err = headResp.Body.Close()
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, headResp.StatusCode)
		require.NotEmpty(t, resp.Data[0].RevisedPrompt)
	}
}

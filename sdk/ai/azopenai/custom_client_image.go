//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// CreateImageOptions contains the optional parameters for the Client.CreateImage method.
type CreateImageOptions struct {
	// placeholder for future optional parameters
}

// CreateImageResponse contains the response from method Client.CreateImage.
type CreateImageResponse struct {
	ImageGenerations
}

// CreateImage creates an image using the Dall-E API.
func (client *Client) CreateImage(ctx context.Context, body ImageGenerationOptions, options *CreateImageOptions) (CreateImageResponse, error) {
	// on Azure the image generation API is a poller. This is a temporary state so we're abstracting it away
	// until it becomes a sync endpoint.
	if client.azure {
		return generateImageWithAzure(client, ctx, body)
	}

	return generateImageWithOpenAI(ctx, client, body)
}

func generateImageWithAzure(client *Client, ctx context.Context, body ImageGenerationOptions) (CreateImageResponse, error) {
	resp, err := client.beginAzureBatchImageGeneration(ctx, body, nil)

	if err != nil {
		return CreateImageResponse{}, err
	}

	v, err := resp.PollUntilDone(ctx, nil)

	if err != nil {
		return CreateImageResponse{}, err
	}

	return CreateImageResponse{
		ImageGenerations: *v.Result,
	}, nil
}

func generateImageWithOpenAI(ctx context.Context, client *Client, body ImageGenerationOptions) (CreateImageResponse, error) {
	urlPath := "/images/generations"
	req, err := runtime.NewRequest(ctx, http.MethodPost, client.formatURL(urlPath, ""))
	if err != nil {
		return CreateImageResponse{}, err
	}
	reqQP := req.Raw().URL.Query()
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}

	if err := runtime.MarshalAsJSON(req, body); err != nil {
		return CreateImageResponse{}, err
	}

	resp, err := client.internal.Pipeline().Do(req)

	if err != nil {
		return CreateImageResponse{}, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CreateImageResponse{}, runtime.NewResponseError(resp)
	}

	var gens *ImageGenerations

	if err := runtime.UnmarshalAsJSON(resp, &gens); err != nil {
		return CreateImageResponse{}, err
	}

	return CreateImageResponse{
		ImageGenerations: *gens,
	}, err
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

// GetAudioTranscriptionOptions contains the optional parameters for the [Client.GetAudioTranscription] method.
type GetAudioTranscriptionOptions struct {
	// placeholder for future optional parameters
}

// GetAudioTranscriptionResponse contains the response from method [Client.GetAudioTranscription].
type GetAudioTranscriptionResponse struct {
	AudioTranscription
}

// GetAudioTranscription gets transcribed text and associated metadata from provided spoken audio data. Audio will
// be transcribed in the written language corresponding to the language it was spoken in. Gets transcribed text
// and associated metadata from provided spoken audio data. Audio will be transcribed in the written language corresponding
// to the language it was spoken in.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
//   - body - contains parameters to specify audio data to transcribe and control the transcription.
//   - options - optional parameters for this method.
func (client *Client) GetAudioTranscription(ctx context.Context, body AudioTranscriptionOptions, options *GetAudioTranscriptionOptions) (GetAudioTranscriptionResponse, error) {
	resp, err := client.getAudioTranscriptionInternal(ctx, streaming.NopCloser(bytes.NewReader(body.File)), &getAudioTranscriptionInternalOptions{
		Filename:       body.Filename,
		Language:       body.Language,
		DeploymentName: body.DeploymentName,
		Prompt:         body.Prompt,
		ResponseFormat: body.ResponseFormat,
		Temperature:    body.Temperature,
	})

	if err != nil {
		return GetAudioTranscriptionResponse{}, err
	}

	return GetAudioTranscriptionResponse(resp), nil
}

// GetAudioTranslationOptions contains the optional parameters for the [Client.GetAudioTranslation] method.
type GetAudioTranslationOptions struct {
	// placeholder for future optional parameters
}

// GetAudioTranslationResponse contains the response from method [Client.GetAudioTranslation].
type GetAudioTranslationResponse struct {
	AudioTranslation
}

// GetAudioTranslation gets English language transcribed text and associated metadata from provided spoken audio
// data. Gets English language transcribed text and associated metadata from provided spoken audio data.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
//   - body - contains parameters to specify audio data to translate and control the translation.
//   - options - optional parameters for this method.
func (client *Client) GetAudioTranslation(ctx context.Context, body AudioTranslationOptions, options *GetAudioTranslationOptions) (GetAudioTranslationResponse, error) {
	resp, err := client.getAudioTranslationInternal(ctx, streaming.NopCloser(bytes.NewReader(body.File)), &getAudioTranslationInternalOptions{
		Filename:       body.Filename,
		DeploymentName: body.DeploymentName,
		Prompt:         body.Prompt,
		ResponseFormat: body.ResponseFormat,
		Temperature:    body.Temperature,
	})

	if err != nil {
		return GetAudioTranslationResponse{}, err
	}

	return GetAudioTranslationResponse(resp), nil
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azopenai

import "io"

// CancelBatchResponse contains the response from method Client.CancelBatch.
type CancelBatchResponse struct {
	// The Batch object.
	Batch
}

// CreateBatchResponse contains the response from method Client.CreateBatch.
type CreateBatchResponse struct {
	// The Batch object.
	Batch
}

// DeleteFileResponse contains the response from method Client.DeleteFile.
type DeleteFileResponse struct {
	// A status response from a file deletion operation.
	FileDeletionStatus
}

// GenerateSpeechFromTextResponse contains the response from method Client.GenerateSpeechFromText.
type GenerateSpeechFromTextResponse struct {
	// Body contains the streaming response.
	Body io.ReadCloser
}


// getAudioTranscriptionInternalResponse contains the response from method Client.getAudioTranscriptionInternal.
type getAudioTranscriptionInternalResponse struct {
	// Result information for an operation that transcribed spoken audio into written text.
	AudioTranscription
}


// getAudioTranslationInternalResponse contains the response from method Client.getAudioTranslationInternal.
type getAudioTranslationInternalResponse struct {
	// Result information for an operation that translated spoken audio into written text.
	AudioTranslation
}

// GetBatchResponse contains the response from method Client.GetBatch.
type GetBatchResponse struct {
	// The Batch object.
	Batch
}

// GetChatCompletionsResponse contains the response from method Client.GetChatCompletions.
type GetChatCompletionsResponse struct {
	// Representation of the response data from a chat completions request.
// Completions support a wide variety of tasks and generate text that continues from or "completes"
// provided prompt data.
	ChatCompletions
}

// GetCompletionsResponse contains the response from method Client.GetCompletions.
type GetCompletionsResponse struct {
	// Representation of the response data from a completions request.
// Completions support a wide variety of tasks and generate text that continues from or "completes"
// provided prompt data.
	Completions
}

// GetEmbeddingsResponse contains the response from method Client.GetEmbeddings.
type GetEmbeddingsResponse struct {
	// Representation of the response data from an embeddings request.
// Embeddings measure the relatedness of text strings and are commonly used for search, clustering,
// recommendations, and other similar scenarios.
	Embeddings
}

// GetFileContentResponse contains the response from method Client.GetFileContent.
type GetFileContentResponse struct {
	Value []byte
}

// GetFileResponse contains the response from method Client.GetFile.
type GetFileResponse struct {
	// Represents an assistant that can call the model and use tools.
	File
}

// GetImageGenerationsResponse contains the response from method Client.GetImageGenerations.
type GetImageGenerationsResponse struct {
	// The result of a successful image generation operation.
	ImageGenerations
}

// ListBatchesResponse contains the response from method Client.ListBatches.
type ListBatchesResponse struct {
	// A list of paginated Batch objects.
	ListBatchesPage
}

// ListFilesResponse contains the response from method Client.ListFiles.
type ListFilesResponse struct {
	// The response data from a file list operation.
	FileListResponse
}

// UploadFileResponse contains the response from method Client.UploadFile.
type UploadFileResponse struct {
	// Represents an assistant that can call the model and use tools.
	File
}


//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// API Reference: https://platform.openai.com/docs/api-reference/assistants-streaming

// CreateThreadAndRunStreamResponse contains the response from [CreateThreadAndRunStream].
type CreateThreadAndRunStreamResponse struct {
	// Stream can be used to stream response events.
	Stream *EventReader[StreamEvent]
}

// CreateThreadAndRunStreamOptions contains the optional parameters for [CreateThreadAndRunStream].
type CreateThreadAndRunStreamOptions struct {
	// for future expansion
}

// CreateThreadAndRunStream is the equivalent of [CreateThreadAndRun], but it returns a stream of responses instead of a
// single response.
func (client *Client) CreateThreadAndRunStream(ctx context.Context, body CreateAndRunThreadBody, _ *CreateThreadAndRunStreamOptions) (CreateThreadAndRunStreamResponse, error) {
	// enable streaming.
	// https: //platform.openai.com/docs/api-reference/runs/createThreadAndRun#runs-createthreadandrun-stream
	body.stream = to.Ptr(true)

	var err error

	req, err := client.createThreadAndRunCreateRequest(ctx, body, nil)
	if err != nil {
		return CreateThreadAndRunStreamResponse{}, err
	}

	runtime.SkipBodyDownload(req) // we'll handle this.

	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CreateThreadAndRunStreamResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CreateThreadAndRunStreamResponse{}, err
	}

	return CreateThreadAndRunStreamResponse{
		Stream: newEventReader(httpResp.Body, unmarshalStreamEvent),
	}, err
}

// API Reference: https://platform.openai.com/docs/api-reference/runs/createRun

// CreateRunStreamResponse contains the response from [CreateRunStream].
type CreateRunStreamResponse struct {
	// Stream can be used to stream response events.
	Stream *EventReader[StreamEvent]
}

// CreateRunStreamOptions contains the optional parameters for [CreateRunStream].
type CreateRunStreamOptions struct {
	CreateRunOptions
	// for future expansion
}

// CreateRunStream is the equivalent of [CreateRun], but it returns a stream of responses instead of a
// single response.
func (client *Client) CreateRunStream(ctx context.Context, threadID string, body CreateRunBody, options *CreateRunStreamOptions) (CreateRunStreamResponse, error) {
	var err error
	body.stream = to.Ptr(true)

	req, err := client.createRunCreateRequest(ctx, threadID, body, &options.CreateRunOptions)
	if err != nil {
		return CreateRunStreamResponse{}, err
	}

	runtime.SkipBodyDownload(req) // we'll handle this.

	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CreateRunStreamResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CreateRunStreamResponse{}, err
	}

	return CreateRunStreamResponse{
		Stream: newEventReader[StreamEvent](httpResp.Body, unmarshalStreamEvent),
	}, err
}

// SubmitToolOutputsToRunStreamResponse contains the response from [SubmitToolOutputsToRunStream].
type SubmitToolOutputsToRunStreamResponse struct {
	// Stream can be used to stream response events.
	Stream *EventReader[StreamEvent]
}

// SubmitToolOutputsToRunStreamOptions contains the optional parameters for [SubmitToolOutputsToRunStream].
type SubmitToolOutputsToRunStreamOptions struct {
}

// SubmitToolOutputsToRunStream is the equivalent of [SubmitToolOutputsToRun], but it returns a stream of responses instead of a
// single response.
func (client *Client) SubmitToolOutputsToRunStream(ctx context.Context, threadID string, runID string, body SubmitToolOutputsToRunBody, options *SubmitToolOutputsToRunOptions) (SubmitToolOutputsToRunStreamResponse, error) {
	var err error

	body.stream = to.Ptr(true)

	req, err := client.submitToolOutputsToRunCreateRequest(ctx, threadID, runID, body, options)
	if err != nil {
		return SubmitToolOutputsToRunStreamResponse{}, err
	}

	runtime.SkipBodyDownload(req) // we'll handle this.

	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SubmitToolOutputsToRunStreamResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return SubmitToolOutputsToRunStreamResponse{}, err
	}

	return SubmitToolOutputsToRunStreamResponse{
		Stream: newEventReader[StreamEvent](httpResp.Body, unmarshalStreamEvent),
	}, nil
}

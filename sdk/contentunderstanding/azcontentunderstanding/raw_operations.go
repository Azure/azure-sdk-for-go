// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

func setJSONBody(req *policy.Request, body json.RawMessage, contentType string) error {
	if !json.Valid(body) || !bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) {
		return errors.New("request body must be a JSON object")
	}
	return req.SetBody(streaming.NopCloser(bytes.NewReader(body)), contentType)
}

// BeginCreateAnalyzerJSON creates an analyzer from a JSON definition without
// discarding unmodeled properties. Capture the HTTP response for raw output.
func (c *Client) BeginCreateAnalyzerJSON(ctx context.Context, analyzerID string, definition json.RawMessage, options *ClientBeginCreateAnalyzerOptions) (*runtime.Poller[ClientCreateAnalyzerResponse], error) {
	if options != nil && options.ResumeToken != "" {
		return c.BeginCreateAnalyzer(ctx, analyzerID, ContentAnalyzer{}, options)
	}
	req, err := c.createAnalyzerCreateRequest(ctx, analyzerID, ContentAnalyzer{}, options)
	if err != nil {
		return nil, err
	}
	if err := setJSONBody(req, definition, "application/json"); err != nil {
		return nil, err
	}
	resp, err := c.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}
	return runtime.NewPoller[ClientCreateAnalyzerResponse](resp, c.internal.Pipeline(), nil)
}

// UpdateAnalyzerJSON applies a JSON merge patch, preserving explicit null, false,
// zero and properties unknown to this library.
func (c *Client) UpdateAnalyzerJSON(ctx context.Context, analyzerID string, patch json.RawMessage, options *ClientUpdateAnalyzerOptions) (ClientUpdateAnalyzerResponse, error) {
	req, err := c.updateAnalyzerCreateRequest(ctx, analyzerID, ContentAnalyzer{}, options)
	if err != nil {
		return ClientUpdateAnalyzerResponse{}, err
	}
	if err := setJSONBody(req, patch, "application/merge-patch+json"); err != nil {
		return ClientUpdateAnalyzerResponse{}, err
	}
	resp, err := c.internal.Pipeline().Do(req)
	if err != nil {
		return ClientUpdateAnalyzerResponse{}, err
	}
	return c.updateAnalyzerHandleResponse(resp, http.StatusOK)
}

// BeginAnalyzeJSON submits a raw JSON analysis request. Query options and
// stringEncoding are handled identically to BeginAnalyze.
func (c *Client) BeginAnalyzeJSON(ctx context.Context, analyzerID string, body json.RawMessage, stringEncoding string, options *ClientBeginAnalyzeOptions) (*runtime.Poller[ClientAnalyzeResponse], error) {
	if options != nil && options.ResumeToken != "" {
		return c.BeginAnalyze(ctx, analyzerID, nil, stringEncoding, options)
	}
	req, err := c.analyzeCreateRequest(ctx, analyzerID, nil, stringEncoding, options)
	if err != nil {
		return nil, err
	}
	if err := setJSONBody(req, body, "application/json"); err != nil {
		return nil, err
	}
	resp, err := c.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return runtime.NewPoller[ClientAnalyzeResponse](resp, c.internal.Pipeline(), &runtime.NewPollerOptions[ClientAnalyzeResponse]{OperationLocationResultPath: "result"})
}

// AnalyzeInlineJSON submits a raw JSON inline request and validates its terminal status.
func (c *Client) AnalyzeInlineJSON(ctx context.Context, analyzerID string, body json.RawMessage, stringEncoding string, options *ClientAnalyzeInlineOptions) (ClientAnalyzeInlineResponse, error) {
	req, err := c.analyzeInlineCreateRequest(ctx, analyzerID, nil, stringEncoding, options)
	if err != nil {
		return ClientAnalyzeInlineResponse{}, err
	}
	if err := setJSONBody(req, body, "application/json"); err != nil {
		return ClientAnalyzeInlineResponse{}, err
	}
	resp, err := c.internal.Pipeline().Do(req)
	if err != nil {
		return ClientAnalyzeInlineResponse{}, err
	}
	return c.analyzeInlineHandleResponse(resp, http.StatusOK)
}

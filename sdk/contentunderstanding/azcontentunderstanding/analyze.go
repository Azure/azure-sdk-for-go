// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func validateAnalyze(body AnalyzeRequest) error {
	if len(body.Inputs) == 0 {
		return errors.New("at least one analysis input is required")
	}
	for _, input := range body.Inputs {
		if (input.URL != "") == (len(input.Data) != 0) {
			return errors.New("each analysis input requires exactly one non-empty URL or Data")
		}
	}
	return nil
}

// BeginAnalyze submits JSON inputs. PollUntilDone accepts a Frequency of one second
// for interactive use; nil polling options use the azcore default.
func (c *Client) BeginAnalyze(ctx context.Context, analyzerID string, body AnalyzeRequest) (*runtime.Poller[AnalyzeResponse], error) {
	return c.beginAnalyze(ctx, analyzerID, body, nil, "")
}

// BeginAnalyzeBinary submits binary input. The content type defaults to application/octet-stream.
// The stream must support seeking so retries can resend it.
func (c *Client) BeginAnalyzeBinary(ctx context.Context, analyzerID string, body io.ReadSeekCloser, contentType string) (*runtime.Poller[AnalyzeResponse], error) {
	if body == nil {
		return nil, errors.New("binary input is required")
	}
	return c.beginAnalyze(ctx, analyzerID, AnalyzeRequest{}, body, contentType)
}

func (c *Client) beginAnalyze(ctx context.Context, analyzerID string, input AnalyzeRequest, binary io.ReadSeekCloser, contentType string) (*runtime.Poller[AnalyzeResponse], error) {
	resp, err := c.analyze(ctx, analyzerID, input, binary, contentType, false)
	if err != nil {
		return nil, err
	}
	return runtime.NewPoller[AnalyzeResponse](resp, c.internal.Pipeline(), &runtime.NewPollerOptions[AnalyzeResponse]{
		FinalStateVia: runtime.FinalStateViaOpLocation,
	})
}

// AnalyzeInline analyzes JSON inputs without polling. Requires a preview API version.
func (c *Client) AnalyzeInline(ctx context.Context, analyzerID string, body AnalyzeRequest) (AnalyzeResponse, error) {
	return c.analyzeInline(ctx, analyzerID, body, nil, "")
}

// AnalyzeBinaryInline analyzes binary input without polling. Requires a preview API version.
func (c *Client) AnalyzeBinaryInline(ctx context.Context, analyzerID string, body io.ReadSeekCloser, contentType string) (AnalyzeResponse, error) {
	if body == nil {
		return AnalyzeResponse{}, errors.New("binary input is required")
	}
	return c.analyzeInline(ctx, analyzerID, AnalyzeRequest{}, body, contentType)
}

func (c *Client) analyzeInline(ctx context.Context, analyzerID string, input AnalyzeRequest, binary io.ReadSeekCloser, contentType string) (AnalyzeResponse, error) {
	if c.apiVersion == APIVersionGA {
		return AnalyzeResponse{}, errors.New("inline analysis requires a preview API version")
	}
	resp, err := c.analyze(ctx, analyzerID, input, binary, contentType, true)
	if err != nil {
		return AnalyzeResponse{}, err
	}
	var state struct {
		Status string `json:"status"`
		Error  struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if _, err := runtime.Payload(resp); err != nil {
		return AnalyzeResponse{}, err
	}
	if err := runtime.UnmarshalAsJSON(resp, &state); err != nil {
		return AnalyzeResponse{}, err
	}
	if state.Status != "Succeeded" {
		code := state.Error.Code
		if code == "" {
			code = "UnexpectedOperationStatus"
		}
		return AnalyzeResponse{}, &azcore.ResponseError{StatusCode: resp.StatusCode, ErrorCode: code, RawResponse: resp}
	}
	return decodeResponse[AnalyzeResponse](resp)
}

func (c *Client) analyze(ctx context.Context, analyzerID string, input AnalyzeRequest, binary io.ReadSeekCloser, contentType string, inline bool) (*http.Response, error) {
	path, err := analyzerPath(analyzerID)
	if err != nil {
		return nil, err
	}
	action := ":analyze"
	if binary != nil {
		action += "Binary"
	}
	if inline {
		action += "Inline"
	}
	query := url.Values{"stringEncoding": {"codePoint"}}
	if binary == nil {
		if err := validateAnalyze(input); err != nil {
			return nil, err
		}
		req, err := c.jsonRequest(ctx, http.MethodPost, path+action, query, input, "application/json")
		if err != nil {
			return nil, err
		}
		if inline {
			return c.send(req, http.StatusOK)
		}
		return c.send(req, http.StatusAccepted)
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	req, err := c.request(ctx, http.MethodPost, path+action, query, nil, "")
	if err != nil {
		return nil, err
	}
	if err := req.SetBody(binary, contentType); err != nil {
		return nil, err
	}
	if inline {
		return c.send(req, http.StatusOK)
	}
	return c.send(req, http.StatusAccepted)
}

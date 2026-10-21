// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// GetAnalyzer retrieves an analyzer definition.
func (c *Client) GetAnalyzer(ctx context.Context, analyzerID string) (Analyzer, error) {
	path, err := analyzerPath(analyzerID)
	if err != nil {
		return Analyzer{}, err
	}
	req, err := c.request(ctx, http.MethodGet, path, nil, nil, "")
	if err != nil {
		return Analyzer{}, err
	}
	resp, err := c.send(req, http.StatusOK)
	if err != nil {
		return Analyzer{}, err
	}
	return decodeResponse[Analyzer](resp)
}

// NewListAnalyzersPager iterates over all analyzer pages.
func (c *Client) NewListAnalyzersPager() *runtime.Pager[ListAnalyzersResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListAnalyzersResponse]{
		More: func(page ListAnalyzersResponse) bool { return page.NextLink != "" },
		Fetcher: func(ctx context.Context, page *ListAnalyzersResponse) (ListAnalyzersResponse, error) {
			req, err := c.request(ctx, http.MethodGet, "/analyzers", nil, nil, "")
			if err != nil {
				return ListAnalyzersResponse{}, err
			}
			if page != nil {
				next, err := url.Parse(page.NextLink)
				if err != nil {
					return ListAnalyzersResponse{}, err
				}
				next = req.Raw().URL.ResolveReference(next)
				if next.Scheme != req.Raw().URL.Scheme || next.Host != req.Raw().URL.Host || next.User != nil {
					return ListAnalyzersResponse{}, errors.New("analyzer nextLink must remain on the resource endpoint")
				}
				req, err = runtime.NewRequest(ctx, http.MethodGet, next.String())
				if err != nil {
					return ListAnalyzersResponse{}, err
				}
			}
			resp, err := c.send(req, http.StatusOK)
			if err != nil {
				return ListAnalyzersResponse{}, err
			}
			return decodeResponse[ListAnalyzersResponse](resp)
		},
	})
}

// BeginCreateAnalyzer creates an analyzer without replacing an existing definition.
// definition is an analyzer JSON object; unknown schema/config properties are retained.
func (c *Client) BeginCreateAnalyzer(ctx context.Context, analyzerID string, definition json.RawMessage) (*runtime.Poller[Analyzer], error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(definition, &object); err != nil || object == nil {
		return nil, errors.New("analyzer definition must be a JSON object")
	}
	path, err := analyzerPath(analyzerID)
	if err != nil {
		return nil, err
	}
	req, err := c.request(ctx, http.MethodPut, path, url.Values{"allowReplace": {"false"}}, definition, "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.send(req, http.StatusCreated, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return runtime.NewPoller[Analyzer](resp, c.internal.Pipeline(), &runtime.NewPollerOptions[Analyzer]{
		FinalStateVia: runtime.FinalStateViaOpLocation, OperationLocationResultPath: "result",
	})
}

// DeleteAnalyzer deletes an analyzer. Missing resources follow service DELETE semantics.
func (c *Client) DeleteAnalyzer(ctx context.Context, analyzerID string) error {
	path, err := analyzerPath(analyzerID)
	if err != nil {
		return err
	}
	req, err := c.request(ctx, http.MethodDelete, path, nil, nil, "")
	if err != nil {
		return err
	}
	resp, err := c.send(req, http.StatusNoContent)
	if err == nil {
		runtime.Drain(resp)
	}
	return err
}

// BeginCopyAnalyzer starts a same-resource or cross-resource copy without replacement.
// For cross-resource copy call this method on the destination client after granting
// authorization on the source client.
func (c *Client) BeginCopyAnalyzer(ctx context.Context, analyzerID string, body CopyAnalyzerRequest) (*runtime.Poller[Analyzer], error) {
	if body.SourceAnalyzerID == "" {
		return nil, errors.New("source analyzer ID is required")
	}
	path, err := analyzerPath(analyzerID)
	if err != nil {
		return nil, err
	}
	req, err := c.jsonRequest(ctx, http.MethodPost, path+":copy", url.Values{"allowReplace": {"false"}}, body, "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.send(req, http.StatusCreated, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return runtime.NewPoller[Analyzer](resp, c.internal.Pipeline(), &runtime.NewPollerOptions[Analyzer]{
		FinalStateVia: runtime.FinalStateViaOpLocation, OperationLocationResultPath: "result",
	})
}

// GrantCopyAuthorization permits a destination resource to copy this source analyzer.
func (c *Client) GrantCopyAuthorization(ctx context.Context, analyzerID string, body GrantCopyAuthorizationRequest) (CopyAuthorization, error) {
	if body.TargetAzureResourceID == "" {
		return CopyAuthorization{}, errors.New("target Azure resource ID is required")
	}
	path, err := analyzerPath(analyzerID)
	if err != nil {
		return CopyAuthorization{}, err
	}
	req, err := c.jsonRequest(ctx, http.MethodPost, path+":grantCopyAuthorization", nil, body, "application/json")
	if err != nil {
		return CopyAuthorization{}, err
	}
	resp, err := c.send(req, http.StatusOK)
	if err != nil {
		return CopyAuthorization{}, err
	}
	return decodeResponse[CopyAuthorization](resp)
}

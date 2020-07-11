// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armeventhub

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"time"
)

// EhNamespacePoller provides polling facilities until the operation completes
type EhNamespacePoller interface {
	Done() bool
	Poll(ctx context.Context) (*http.Response, error)
	FinalResponse(ctx context.Context) (*EhNamespaceResponse, error)
	ResumeToken() (string, error)
}

type ehNamespacePoller struct {
	// the client for making the request
	pipeline azcore.Pipeline
	pt       pollingTracker
}

// Done returns true if there was an error or polling has reached a terminal state
func (p *ehNamespacePoller) Done() bool {
	return p.pt.hasTerminated()
}

// Poll will send poll the service endpoint and return an http.Response or error received from the service
func (p *ehNamespacePoller) Poll(ctx context.Context) (*http.Response, error) {
	if lroPollDone(ctx, p.pipeline, p.pt) {
		return p.pt.latestResponse().Response, p.pt.pollingError()
	}
	return nil, p.pt.pollingError()
}

func (p *ehNamespacePoller) FinalResponse(ctx context.Context) (*EhNamespaceResponse, error) {
	if !p.Done() {
		return nil, errors.New("cannot return a final response from a poller in a non-terminal state")
	}
	if p.pt.pollerMethodVerb() == http.MethodPut || p.pt.pollerMethodVerb() == http.MethodPatch {
		res, err := p.handleResponse(p.pt.latestResponse())
		if err != nil {
			return nil, err
		}
		if res != nil && (*res.EhNamespace != EhNamespace{}) {
			return res, nil
		}
	}
	// checking if there was a FinalStateVia configuration to re-route the final GET
	// request to the value specified in the FinalStateVia property on the poller
	err := p.pt.setFinalState()
	if err != nil {
		return nil, err
	}
	if p.pt.finalGetURL() == "" {
		// we can end up in this situation if the async operation returns a 200
		// with no polling URLs.  in that case return the response which should
		// contain the JSON payload (only do this for successful terminal cases).
		if lr := p.pt.latestResponse(); lr != nil && p.pt.hasSucceeded() {
			result, err := p.handleResponse(lr)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, errors.New("missing URL for retrieving result")
	}
	u, err := url.Parse(p.pt.finalGetURL())
	if err != nil {
		return nil, err
	}
	req := azcore.NewRequest(http.MethodGet, *u)
	if err != nil {
		return nil, err
	}
	resp, err := p.pipeline.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	return p.handleResponse(resp)
}

// ResumeToken generates the string token that can be used with the ResumeEhNamespacePoller method
// on the client to create a new poller from the data held in the current poller type
func (p *ehNamespacePoller) ResumeToken() (string, error) {
	if p.pt.hasTerminated() {
		return "", errors.New("cannot create a ResumeToken from a poller in a terminal state")
	}
	js, err := json.Marshal(p.pt)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func (p *ehNamespacePoller) pollUntilDone(ctx context.Context, frequency time.Duration) (*EhNamespaceResponse, error) {
	// initial check for a retry-after header existing on the initial response
	if retryAfter := azcore.RetryAfter(p.pt.latestResponse().Response); retryAfter > 0 {
		err := delay(ctx, retryAfter)
		if err != nil {
			return nil, err
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := p.Poll(ctx)
		if err != nil {
			return nil, err
		}
		if p.Done() {
			break
		}
		d := frequency
		if retryAfter := azcore.RetryAfter(resp); retryAfter > 0 {
			d = retryAfter
		}
		err = delay(ctx, d)
		if err != nil {
			return nil, err
		}
	}
	return p.FinalResponse(ctx)
}

func (p *ehNamespacePoller) handleResponse(resp *azcore.Response) (*EhNamespaceResponse, error) {
	result := EhNamespaceResponse{RawResponse: resp.Response}
	if resp.HasStatusCode(http.StatusNoContent) {
		return &result, nil
	}
	if !resp.HasStatusCode(pollingCodes[:]...) {
		return nil, p.pt.handleError(resp)
	}
	return &result, resp.UnmarshalAsJSON(&result.EhNamespace)
}

// HTTPPoller provides polling facilities until the operation completes
type HTTPPoller interface {
	Done() bool
	Poll(ctx context.Context) (*http.Response, error)
	FinalResponse() *http.Response
	ResumeToken() (string, error)
}

type httpPoller struct {
	// the client for making the request
	pipeline azcore.Pipeline
	pt       pollingTracker
}

// Done returns true if there was an error or polling has reached a terminal state
func (p *httpPoller) Done() bool {
	return p.pt.hasTerminated()
}

// Poll will send poll the service endpoint and return an http.Response or error received from the service
func (p *httpPoller) Poll(ctx context.Context) (*http.Response, error) {
	if lroPollDone(ctx, p.pipeline, p.pt) {
		return p.pt.latestResponse().Response, p.pt.pollingError()
	}
	return nil, p.pt.pollingError()
}

func (p *httpPoller) FinalResponse() *http.Response {
	return p.pt.latestResponse().Response
}

// ResumeToken generates the string token that can be used with the ResumeHTTPPoller method
// on the client to create a new poller from the data held in the current poller type
func (p *httpPoller) ResumeToken() (string, error) {
	if p.pt.hasTerminated() {
		return "", errors.New("cannot create a ResumeToken from a poller in a terminal state")
	}
	js, err := json.Marshal(p.pt)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func (p *httpPoller) pollUntilDone(ctx context.Context, frequency time.Duration) (*http.Response, error) {
	// initial check for a retry-after header existing on the initial response
	if retryAfter := azcore.RetryAfter(p.pt.latestResponse().Response); retryAfter > 0 {
		err := delay(ctx, retryAfter)
		if err != nil {
			return nil, err
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := p.Poll(ctx)
		if err != nil {
			return nil, err
		}
		if p.Done() {
			break
		}
		d := frequency
		if retryAfter := azcore.RetryAfter(resp); retryAfter > 0 {
			d = retryAfter
		}
		err = delay(ctx, d)
		if err != nil {
			return nil, err
		}
	}
	return p.FinalResponse(), nil
}

func delay(ctx context.Context, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

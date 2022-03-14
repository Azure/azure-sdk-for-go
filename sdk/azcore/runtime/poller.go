//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/op"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// NewPoller creates a Poller based on the provided initial response.
// pollerID - a unique identifier for an LRO, it's usually the client.Method string.
func NewPoller[T any](pollerID string, resp *http.Response, pl pipeline.Pipeline, rt *T) (*Poller[T], error) {
	defer resp.Body.Close()
	// this is a back-stop in case the swagger is incorrect (i.e. missing one or more status codes for success).
	// ideally the codegen should return an error if the initial response failed and not even create a poller.
	if !pollers.StatusCodeValid(resp) {
		return nil, errors.New("the operation failed or was cancelled")
	}
	// determine the polling method
	var lro pollers.Operation
	var err error
	// op poller must be checked first as it can also have a location header
	if op.Applicable(resp) {
		lro, err = op.New(resp, pollerID)
	} else if loc.Applicable(resp) {
		lro, err = loc.New(resp, pollerID)
	} else {
		lro = &pollers.NopPoller{}
	}
	if err != nil {
		return nil, err
	}
	return &Poller[T]{
		pt: pollers.NewPoller(lro, resp, pl),
		rt: rt,
	}, nil
}

// NewPollerFromResumeToken creates a Poller from a resume token string.
// pollerID - a unique identifier for an LRO, it's usually the client.Method string.
func NewPollerFromResumeToken[T any](pollerID string, token string, pl pipeline.Pipeline, rt *T) (*Poller[T], error) {
	kind, err := pollers.KindFromToken(pollerID, token)
	if err != nil {
		return nil, err
	}
	// now rehydrate the poller based on the encoded poller type
	var lro pollers.Operation
	switch kind {
	case loc.Kind:
		log.Writef(log.EventLRO, "Resuming %s poller.", loc.Kind)
		lro = &loc.Poller{}
	case op.Kind:
		log.Writef(log.EventLRO, "Resuming %s poller.", op.Kind)
		lro = &op.Poller{}
	default:
		return nil, fmt.Errorf("unhandled poller type %s", kind)
	}
	if err = json.Unmarshal([]byte(token), lro); err != nil {
		return nil, err
	}
	return &Poller[T]{
		pt: pollers.NewPoller(lro, nil, pl),
		rt: rt,
	}, nil
}

// Poller encapsulates a long-running operation, providing polling facilities until the operation reaches a terminal state.
type Poller[T any] struct {
	pt *pollers.Poller
	rt *T
}

// PollUntilDone will poll the service endpoint until a terminal state is reached, an error is received, or the context expires.
// freq: the time to wait between intervals in absence of a Retry-After header. Allowed minimum is one second.
func (p *Poller[T]) PollUntilDone(ctx context.Context, freq time.Duration) (T, error) {
	var resp T
	if p.rt != nil {
		resp = *p.rt
	}
	_, err := p.pt.PollUntilDone(ctx, freq, &resp)
	return resp, err
}

// Poll fetches the latest state of the LRO.  It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP
// response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is
// updated and the error is returned.
// If the LRO has not reached a terminal state, the poller's state is updated and
// the latest HTTP response is returned.
// If Poll fails, the poller's state is unmodified and the error is returned.
// Calling Poll on an LRO that has reached a terminal state will return the final
// HTTP response or error.
func (p *Poller[T]) Poll(ctx context.Context) (*http.Response, error) {
	return p.pt.Poll(ctx)
}

// Done returns true if the LRO has reached a terminal state.
func (p *Poller[T]) Done() bool {
	return p.pt.Done()
}

// Result returns the result of the LRO and is meant to be used in conjunction with Poll and Done.
// Calling this on an LRO in a non-terminal state will return an error.
func (p *Poller[T]) Result(ctx context.Context) (T, error) {
	var resp T
	if p.rt != nil {
		resp = *p.rt
	}
	_, err := p.pt.FinalResponse(ctx, &resp)
	return resp, err
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *Poller[T]) ResumeToken() (string, error) {
	return p.pt.ResumeToken()
}

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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// FinalStateVia is the enumerated type for the possible final-state-via values.
type FinalStateVia = pollers.FinalStateVia

const (
	// FinalStateViaAzureAsyncOp indicates the final payload comes from the Azure-AsyncOperation URL.
	FinalStateViaAzureAsyncOp = pollers.FinalStateViaAzureAsyncOp

	// FinalStateViaLocation indicates the final payload comes from the Location URL.
	FinalStateViaLocation = pollers.FinalStateViaLocation

	// FinalStateViaOriginalURI indicates the final payload comes from the original URL.
	FinalStateViaOriginalURI = pollers.FinalStateViaOriginalURI

	// FinalStateViaOpLocation indicates the final payload comes from the Operation-Location URL.
	FinalStateViaOpLocation = pollers.FinalStateViaOpLocation
)

// NewPollerOptions contains the optional parameters for NewPoller.
type NewPollerOptions[T any] struct {
	// FinalStateVia contains the final-state-via value for the LRO.
	FinalStateVia FinalStateVia

	// Response contains a preconstructed response type.
	// The final payload will be unmarshaled into it and returned.
	Response *T
}

// NewPoller creates a Poller based on the provided initial response.
func NewPoller[T any](resp *http.Response, pl pipeline.Pipeline, options *NewPollerOptions[T]) (*Poller[T], error) {
	if options == nil {
		options = &NewPollerOptions[T]{}
	}
	defer resp.Body.Close()
	// this is a back-stop in case the swagger is incorrect (i.e. missing one or more status codes for success).
	// ideally the codegen should return an error if the initial response failed and not even create a poller.
	if !pollers.StatusCodeValid(resp) {
		return nil, errors.New("the LRO failed or was cancelled")
	}
	tName, err := pollers.PollerTypeName[T]()
	if err != nil {
		return nil, err
	}
	// determine the polling method
	var lro pollers.Operation
	if async.Applicable(resp) {
		lro, err = async.New(resp, options.FinalStateVia, tName)
	} else if loc.Applicable(resp) {
		lro, err = loc.New(resp, tName)
	} else if body.Applicable(resp) {
		// must test body poller last as it's a subset of the other pollers.
		// TODO: this is ambiguous for PATCH/PUT if it returns a 200 with no polling headers (sync completion)
		lro, err = body.New(resp, tName)
	} else if m := resp.Request.Method; resp.StatusCode == http.StatusAccepted && (m == http.MethodDelete || m == http.MethodPost) {
		// if we get here it means we have a 202 with no polling headers.
		// for DELETE and POST this is a hard error per ARM RPC spec.
		return nil, errors.New("response is missing polling URL")
	} else {
		lro = &pollers.NopPoller{}
	}
	if err != nil {
		return nil, err
	}
	return &Poller[T]{
		pt: pollers.NewPoller(lro, resp, pl),
		rt: options.Response,
	}, nil
}

// NewPollerFromResumeTokenOptions contains the optional parameters for NewPollerFromResumeToken.
type NewPollerFromResumeTokenOptions[T any] struct {
	// Response contains a preconstructed response type.
	// The final payload will be unmarshaled into it and returned.
	Response *T
}

// NewPollerFromResumeToken creates a Poller from a resume token string.
func NewPollerFromResumeToken[T any](token string, pl pipeline.Pipeline, options *NewPollerFromResumeTokenOptions[T]) (*Poller[T], error) {
	if options == nil {
		options = &NewPollerFromResumeTokenOptions[T]{}
	}
	tName, err := pollers.PollerTypeName[T]()
	if err != nil {
		return nil, err
	}
	kind, err := pollers.KindFromToken(tName, token)
	if err != nil {
		return nil, err
	}
	// now rehydrate the poller based on the encoded poller type
	var lro pollers.Operation
	switch kind {
	case async.Kind:
		log.Writef(log.EventLRO, "Resuming %s poller.", async.Kind)
		lro = &async.Poller{}
	case loc.Kind:
		log.Writef(log.EventLRO, "Resuming %s poller.", loc.Kind)
		lro = &loc.Poller{}
	case body.Kind:
		log.Writef(log.EventLRO, "Resuming %s poller.", body.Kind)
		lro = &body.Poller{}
	default:
		return nil, fmt.Errorf("unhandled poller type %s", kind)
	}
	if err = json.Unmarshal([]byte(token), lro); err != nil {
		return nil, err
	}
	return &Poller[T]{
		pt: pollers.NewPoller(lro, nil, pl),
		rt: options.Response,
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

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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/armloc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/op"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
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

	// Handler[T] contains a custom polling implementation.
	Handler PollingHandler[T]
}

// NewPoller creates a Poller based on the provided initial response.
func NewPoller[T any](resp *http.Response, pl exported.Pipeline, options *NewPollerOptions[T]) (*Poller[T], error) {
	if options == nil {
		options = &NewPollerOptions[T]{}
	}
	if options.Handler != nil {
		return &Poller[T]{
			op:     options.Handler,
			resp:   resp,
			result: options.Response,
		}, nil
	}

	defer resp.Body.Close()
	// this is a back-stop in case the swagger is incorrect (i.e. missing one or more status codes for success).
	// ideally the codegen should return an error if the initial response failed and not even create a poller.
	if !pollers.StatusCodeValid(resp) {
		return nil, errors.New("the operation failed or was cancelled")
	}

	// determine the polling method
	var opr PollingHandler[T]
	var err error
	if async.Applicable(resp) {
		// async poller must be checked first as it can also have a location header
		opr, err = async.New[T](pl, resp, options.FinalStateVia)
	} else if op.Applicable(resp) {
		// op poller must be checked before loc as it can also have a location header
		opr, err = op.New[T](pl, resp, options.FinalStateVia)
	} else if armloc.Applicable(resp) {
		opr, err = armloc.New[T](pl, resp)
	} else if loc.Applicable(resp) {
		opr, err = loc.New[T](pl, resp)
	} else if body.Applicable(resp) {
		// must test body poller last as it's a subset of the other pollers.
		// TODO: this is ambiguous for PATCH/PUT if it returns a 200 with no polling headers (sync completion)
		opr, err = body.New[T](pl, resp)
	} else if m := resp.Request.Method; resp.StatusCode == http.StatusAccepted && (m == http.MethodDelete || m == http.MethodPost) {
		// if we get here it means we have a 202 with no polling headers.
		// for DELETE and POST this is a hard error per ARM RPC spec.
		return nil, errors.New("response is missing polling URL")
	} else {
		opr, err = pollers.NewNopPoller[T](resp)
	}

	if err != nil {
		return nil, err
	}
	return &Poller[T]{
		op:     opr,
		resp:   resp,
		result: options.Response,
	}, nil
}

// NewPollerFromResumeTokenOptions contains the optional parameters for NewPollerFromResumeToken.
type NewPollerFromResumeTokenOptions[T any] struct {
	// Response contains a preconstructed response type.
	// The final payload will be unmarshaled into it and returned.
	Response *T

	// Handler[T] contains a custom polling implementation.
	Handler PollingHandler[T]
}

// NewPollerFromResumeToken creates a Poller from a resume token string.
func NewPollerFromResumeToken[T any](token string, pl exported.Pipeline, options *NewPollerFromResumeTokenOptions[T]) (*Poller[T], error) {
	if options == nil {
		options = &NewPollerFromResumeTokenOptions[T]{}
	}

	if err := pollers.IsTokenValid[T](token); err != nil {
		return nil, err
	}
	raw, err := pollers.ExtractToken(token)
	if err != nil {
		return nil, err
	}
	var asJSON map[string]interface{}
	if err := json.Unmarshal(raw, &asJSON); err != nil {
		return nil, err
	}

	opr := options.Handler
	// now rehydrate the poller based on the encoded poller type
	if async.CanResume(asJSON) {
		opr, _ = async.New[T](pl, nil, "")
	} else if armloc.CanResume(asJSON) {
		opr, _ = armloc.New[T](pl, nil)
	} else if body.CanResume(asJSON) {
		opr, _ = body.New[T](pl, nil)
	} else if loc.CanResume(asJSON) {
		opr, _ = loc.New[T](pl, nil)
	} else if op.CanResume(asJSON) {
		opr, _ = op.New[T](pl, nil, "")
	} else if opr != nil {
		log.Writef(log.EventLRO, "Resuming custom poller %T.", opr)
	} else {
		return nil, fmt.Errorf("unhandled poller token %s", string(raw))
	}
	if err := json.Unmarshal(raw, &opr); err != nil {
		return nil, err
	}
	return &Poller[T]{
		op:     opr,
		result: options.Response,
	}, nil
}

// PollingHandler[T] abstracts the differences among poller implementations.
type PollingHandler[T any] interface {
	// Done returns true if the LRO has reached a terminal state.
	Done() bool

	// Poll fetches the latest state of the LRO.
	Poll(context.Context) (*http.Response, error)

	// Result is called once the LRO has reached a terminal state. It returns the result of the operation.
	// The out parameter is an optional, preconstructed response type to receive the final payload.
	Result(ctx context.Context, out *T) (T, error)
}

// Poller encapsulates a long-running operation, providing polling facilities until the operation reaches a terminal state.
type Poller[T any] struct {
	op     PollingHandler[T]
	resp   *http.Response
	err    error
	result *T
}

// PollUntilDoneOptions contains the optional values for the Poller[T].PollUntilDone() method.
type PollUntilDoneOptions struct {
	// Frequency is the time to wait between polling intervals in absence of a Retry-After header. Allowed minimum is one second.
	// Pass zero to accept the default value (30s).
	Frequency time.Duration
}

// PollUntilDone will poll the service endpoint until a terminal state is reached, an error is received, or the context expires.
// options: pass nil to accept the default values.
// NOTE: the default polling frequency is 30 seconds which works well for most operations.  However, some operations might
// benefit from a shorter or longer duration.
func (p *Poller[T]) PollUntilDone(ctx context.Context, options *PollUntilDoneOptions) (T, error) {
	if options == nil {
		options = &PollUntilDoneOptions{}
	}
	cp := *options
	if cp.Frequency == 0 {
		cp.Frequency = 30 * time.Second
	}

	if cp.Frequency < time.Second {
		return *new(T), errors.New("polling frequency minimum is one second")
	}

	start := time.Now()
	logPollUntilDoneExit := func(v interface{}) {
		log.Writef(log.EventLRO, "END PollUntilDone() for %T: %v, total time: %s", p.op, v, time.Since(start))
	}
	log.Writef(log.EventLRO, "BEGIN PollUntilDone() for %T", p.op)
	if p.resp != nil {
		// initial check for a retry-after header existing on the initial response
		if retryAfter := shared.RetryAfter(p.resp); retryAfter > 0 {
			log.Writef(log.EventLRO, "initial Retry-After delay for %s", retryAfter.String())
			if err := shared.Delay(ctx, retryAfter); err != nil {
				logPollUntilDoneExit(err)
				return *new(T), err
			}
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := p.Poll(ctx)
		if err != nil {
			logPollUntilDoneExit(err)
			return *new(T), err
		}
		if p.Done() {
			logPollUntilDoneExit("succeeded")
			return p.Result(ctx)
		}
		d := cp.Frequency
		if retryAfter := shared.RetryAfter(resp); retryAfter > 0 {
			log.Writef(log.EventLRO, "Retry-After delay for %s", retryAfter.String())
			d = retryAfter
		} else {
			log.Writef(log.EventLRO, "delay for %s", d.String())
		}
		if err = shared.Delay(ctx, d); err != nil {
			logPollUntilDoneExit(err)
			return *new(T), err
		}
	}
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
	if p.Done() {
		// the LRO has reached a terminal state, don't poll again
		if p.err != nil {
			return nil, p.err
		}
		return p.resp, nil
	}
	p.resp, p.err = p.op.Poll(ctx)
	return p.resp, p.err
}

// Done returns true if the LRO has reached a terminal state.
func (p *Poller[T]) Done() bool {
	return p.op.Done()
}

// Result returns the result of the LRO and is meant to be used in conjunction with Poll and Done.
// Calling this on an LRO in a non-terminal state will return an error.
func (p *Poller[T]) Result(ctx context.Context) (T, error) {
	if !p.Done() {
		return *new(T), errors.New("cannot return a final response from a poller in a non-terminal state")
	}
	return p.op.Result(ctx, p.result)
}

// ResumeToken returns a value representing the poller that can be used to resume
// the LRO at a later time. ResumeTokens are unique per service operation.
func (p *Poller[T]) ResumeToken() (string, error) {
	if p.Done() {
		return "", errors.New("cannot create a ResumeToken from a poller in a terminal state")
	}
	tk, err := pollers.NewResumeToken[T](p.op)
	if err != nil {
		return "", err
	}
	return tk, err
}

// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ErrorUnmarshaller is the func to invoke when the endpoint returns an error response that requires unmarshalling.
type ErrorUnmarshaller func(*azcore.Response) error

// NewLROPoller creates an LROPoller based on the provided initial response.
// pollerID - a unique identifier for an LRO.  it's usually the client.Method string.
// NOTE: this is only meant for internal use in generated code.
func NewLROPoller(pollerID string, finalState string, resp *azcore.Response, pl azcore.Pipeline, eu ErrorUnmarshaller) (*LROPoller, error) {
	// this is a back-stop in case the swagger is incorrect (i.e. missing one or more status codes for success).
	// ideally the codegen should return an error if the initial response failed and not even create a poller.
	if !lroStatusCodeValid(resp) {
		return nil, errors.New("the LRO failed or was cancelled")
	}
	// determine the polling method
	var lro lroPoller
	var err error
	if async.Applicable(resp) {
		lro, err = async.New(resp, finalState, pollerID)
	} else if loc.Applicable(resp) {
		lro, err = loc.New(resp, pollerID)
	} else if body.Applicable(resp) {
		// must test body poller last as it's a subset of the other pollers.
		// TODO: this is ambiguous for PATCH/PUT if it returns a 200 with no polling headers (sync completion)
		lro, err = body.New(resp, pollerID)
	} else if m := resp.Request.Method; resp.StatusCode == http.StatusAccepted && (m == http.MethodDelete || m == http.MethodPost) {
		// if we get here it means we have a 202 with no polling headers.
		// for DELETE and POST this is a hard error per ARM RPC spec.
		return nil, errors.New("response is missing polling URL")
	} else {
		lro = &nopPoller{}
	}
	if err != nil {
		return nil, err
	}
	return &LROPoller{lro: lro, Pipeline: pl, eu: eu, resp: resp}, nil
}

// NewLROPollerFromResumeToken creates an LROPoller from a resume token string.
// pollerID - a unique identifier for an LRO.  it's usually the client.Method string.
// NOTE: this is only meant for internal use in generated code.
func NewLROPollerFromResumeToken(pollerID string, token string, pl azcore.Pipeline, eu ErrorUnmarshaller) (*LROPoller, error) {
	// unmarshal into JSON object to determine the poller type
	obj := map[string]interface{}{}
	err := json.Unmarshal([]byte(token), &obj)
	if err != nil {
		return nil, err
	}
	t, ok := obj["type"]
	if !ok {
		return nil, errors.New("missing type field")
	}
	tt, ok := t.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type format %T", t)
	}
	ttID, ttKind, err := pollers.DecodeID(tt)
	if err != nil {
		return nil, err
	}
	// ensure poller types match
	if ttID != pollerID {
		return nil, fmt.Errorf("cannot resume from this poller token.  expected %s, received %s", pollerID, ttID)
	}
	// now rehydrate the poller based on the encoded poller type
	var lro lroPoller
	switch ttKind {
	case "async":
		azcore.Log().Write(azcore.LogLongRunningOperation, "Resuming Azure-AsyncOperation poller.")
		lro = &async.Poller{}
	case "loc":
		azcore.Log().Write(azcore.LogLongRunningOperation, "Resuming Location poller.")
		lro = &loc.Poller{}
	case "body":
		azcore.Log().Write(azcore.LogLongRunningOperation, "Resuming Body poller.")
		lro = &body.Poller{}
	default:
		return nil, fmt.Errorf("unhandled poller type %s", ttKind)
	}
	if err = json.Unmarshal([]byte(token), lro); err != nil {
		return nil, err
	}
	return &LROPoller{lro: lro, Pipeline: pl, eu: eu}, nil
}

// LROPoller encapsulates state and logic for polling on long-running operations.
// NOTE: this is only meant for internal use in generated code.
type LROPoller struct {
	Pipeline azcore.Pipeline
	lro      lroPoller
	eu       ErrorUnmarshaller
	resp     *azcore.Response
	err      error
}

// Done returns true if the LRO has reached a terminal state.
func (l *LROPoller) Done() bool {
	if l.err != nil {
		return true
	}
	return l.lro.Done()
}

// Poll sends a polling request to the polling endpoint and returns the response or error.
func (l *LROPoller) Poll(ctx context.Context) (*http.Response, error) {
	if l.Done() {
		// the LRO has reached a terminal state, don't poll again
		if l.resp != nil {
			return l.resp.Response, nil
		}
		return nil, l.err
	}
	req, err := azcore.NewRequest(ctx, http.MethodGet, l.lro.URL())
	if err != nil {
		return nil, err
	}
	resp, err := l.Pipeline.Do(req)
	if err != nil {
		// don't update the poller for failed requests
		return nil, err
	}
	defer resp.Body.Close()
	if !lroStatusCodeValid(resp) {
		// the LRO failed.  unmarshall the error and update state
		l.err = l.eu(resp)
		l.resp = nil
		return nil, l.err
	}
	if err = l.lro.Update(resp); err != nil {
		return nil, err
	}
	l.resp = resp
	azcore.Log().Writef(azcore.LogLongRunningOperation, "Status %s", l.lro.Status())
	if pollers.Failed(l.lro.Status()) {
		l.err = l.eu(resp)
		l.resp = nil
		return nil, l.err
	}
	return l.resp.Response, nil
}

// ResumeToken returns a token string that can be used to resume a poller that has not yet reached a terminal state.
func (l *LROPoller) ResumeToken() (string, error) {
	if l.Done() {
		return "", errors.New("cannot create a ResumeToken from a poller in a terminal state")
	}
	b, err := json.Marshal(l.lro)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FinalResponse will perform a final GET request and return the final HTTP response for the polling
// operation and unmarshall the content of the payload into the respType interface that is provided.
func (l *LROPoller) FinalResponse(ctx context.Context, respType interface{}) (*http.Response, error) {
	if !l.Done() {
		return nil, errors.New("cannot return a final response from a poller in a non-terminal state")
	}
	// if there's nothing to unmarshall into or no response body just return the final response
	if respType == nil {
		return l.resp.Response, nil
	} else if l.resp.StatusCode == http.StatusNoContent || l.resp.ContentLength == 0 {
		azcore.Log().Write(azcore.LogLongRunningOperation, "final response specifies a response type but no payload was received")
		return l.resp.Response, nil
	}
	if u := l.lro.FinalGetURL(); u != "" {
		azcore.Log().Write(azcore.LogLongRunningOperation, "Performing final GET.")
		req, err := azcore.NewRequest(ctx, http.MethodGet, u)
		if err != nil {
			return nil, err
		}
		resp, err := l.Pipeline.Do(req)
		if err != nil {
			return nil, err
		}
		if !lroStatusCodeValid(resp) {
			return nil, l.eu(resp)
		}
		l.resp = resp
	}
	body, err := ioutil.ReadAll(l.resp.Body)
	l.resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, respType); err != nil {
		return nil, err
	}
	return l.resp.Response, nil
}

// PollUntilDone will handle the entire span of the polling operation until a terminal state is reached,
// then return the final HTTP response for the polling operation and unmarshal the content of the payload
// into the respType interface that is provided.
// freq - the time to wait between polling intervals if the endpoint doesn't send a Retry-After header.
//        A good starting value is 30 seconds.  Note that some resources might benefit from a different value.
func (l *LROPoller) PollUntilDone(ctx context.Context, freq time.Duration, respType interface{}) (*http.Response, error) {
	start := time.Now()
	logPollUntilDoneExit := func(v interface{}) {
		azcore.Log().Writef(azcore.LogLongRunningOperation, "END PollUntilDone() for %T: %v, total time: %s", l.lro, v, time.Since(start))
	}
	azcore.Log().Writef(azcore.LogLongRunningOperation, "BEGIN PollUntilDone() for %T", l.lro)
	if l.resp != nil {
		// initial check for a retry-after header existing on the initial response
		if retryAfter := azcore.RetryAfter(l.resp.Response); retryAfter > 0 {
			azcore.Log().Writef(azcore.LogLongRunningOperation, "initial Retry-After delay for %s", retryAfter.String())
			if err := delay(ctx, retryAfter); err != nil {
				logPollUntilDoneExit(err)
				return nil, err
			}
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := l.Poll(ctx)
		if err != nil {
			logPollUntilDoneExit(err)
			return nil, err
		}
		if l.Done() {
			logPollUntilDoneExit(l.lro.Status())
			return l.FinalResponse(ctx, respType)
		}
		d := freq
		if retryAfter := azcore.RetryAfter(resp); retryAfter > 0 {
			azcore.Log().Writef(azcore.LogLongRunningOperation, "Retry-After delay for %s", retryAfter.String())
			d = retryAfter
		} else {
			azcore.Log().Writef(azcore.LogLongRunningOperation, "delay for %s", d.String())
		}
		if err = delay(ctx, d); err != nil {
			logPollUntilDoneExit(err)
			return nil, err
		}
	}
}

var _ azcore.Poller = (*LROPoller)(nil)

// abstracts the differences between concrete poller types
type lroPoller interface {
	Done() bool
	Update(resp *azcore.Response) error
	FinalGetURL() string
	URL() string
	Status() string
}

// ====================================================================================================

// used if the operation synchronously completed
type nopPoller struct{}

func (*nopPoller) URL() string {
	return ""
}

func (*nopPoller) Done() bool {
	return true
}

func (*nopPoller) Succeeded() bool {
	return true
}

func (*nopPoller) Update(*azcore.Response) error {
	return nil
}

func (*nopPoller) FinalGetURL() string {
	return ""
}

func (*nopPoller) Status() string {
	return pollers.StatusSucceeded
}

// returns true if the LRO response contains a valid HTTP status code
func lroStatusCodeValid(resp *azcore.Response) bool {
	return resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusNoContent)
}

func delay(ctx context.Context, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

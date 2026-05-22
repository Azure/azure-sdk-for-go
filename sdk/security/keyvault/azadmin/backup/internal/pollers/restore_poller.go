// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pollers

// this file contains handwritten additions to the generated code
// custom poller handler for restore operations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/poller"
)

func pollHelper(ctx context.Context, endpoint string, pl runtime.Pipeline, update func(resp *http.Response) (string, error)) error {
	req, err := runtime.NewRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return err
	}
	resp, err := pl.Do(req)
	if err != nil {
		return err
	}
	_, err = update(resp)
	if err != nil {
		return err
	}

	return nil
}

func resultHelper[T any](resp *http.Response, failed bool, out *T) error {
	// short-circuit the simple success case with no response body to unmarshal
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	if !poller.StatusCodeValid(resp) || failed {
		// the LRO failed.  unmarshall the error and update state
		return runtime.NewResponseError(resp)
	}

	// success case
	payload, err := runtime.Payload(resp)
	if err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}

	if err = json.Unmarshal(payload, out); err != nil {
		return err
	}
	return nil
}

// Poller is an LRO poller that uses the Azure-AsyncOperation pattern.
type restorePoller[T any] struct {
	pl runtime.Pipeline

	resp *http.Response

	// The URL from Azure-AsyncOperation header.
	AsyncURL string `json:"asyncURL"`

	// The URL from Location header.
	LocURL string `json:"locURL"`

	// The URL from the initial LRO request.
	OrigURL string `json:"origURL"`

	// The HTTP method from the initial LRO request.
	Method string `json:"method"`

	// The value of final-state-via from swagger, can be the empty string.
	FinalState runtime.FinalStateVia `json:"finalState"`

	// The LRO's current state.
	CurState string `json:"state"`
}

// NewRestorePoller creates a new Poller from the provided initial response and final-state type.
// Pass nil for response to create an empty Poller for rehydration.
func NewRestorePoller[T any](pl runtime.Pipeline, resp *http.Response, finalState runtime.FinalStateVia) (*restorePoller[T], error) {
	if resp == nil {
		return &restorePoller[T]{pl: pl}, nil
	}
	asyncURL := resp.Header.Get("Azure-AsyncOperation")
	if asyncURL == "" {
		return nil, errors.New("response is missing Azure-AsyncOperation header")
	}
	if !poller.IsValidURL(asyncURL) {
		return nil, fmt.Errorf("invalid polling URL %s", asyncURL)
	}
	// check for provisioning state.  if the operation is a RELO
	// and terminates synchronously this will prevent extra polling.
	// it's ok if there's no provisioning state.
	state, _ := poller.GetProvisioningState(resp)
	if state == "" {
		state = poller.StatusInProgress
	}
	p := &restorePoller[T]{
		pl:         pl,
		resp:       resp,
		AsyncURL:   asyncURL,
		LocURL:     resp.Header.Get("Location"),
		OrigURL:    resp.Request.URL.String(),
		Method:     resp.Request.Method,
		FinalState: finalState,
		CurState:   state,
	}
	return p, nil
}

// Done returns true if the LRO is in a terminal state.
func (p *restorePoller[T]) Done() bool {
	return poller.IsTerminalState(p.CurState)
}

// Poll retrieves the current state of the LRO.
func (p *restorePoller[T]) Poll(ctx context.Context) (*http.Response, error) {
	err := pollHelper(ctx, p.AsyncURL, p.pl, func(resp *http.Response) (string, error) {
		if !poller.StatusCodeValid(resp) {
			p.resp = resp
			return "", runtime.NewResponseError(resp)
		}
		state, err := poller.GetStatus(resp)
		if err != nil {
			return "", err
		} else if state == "" {
			return "", errors.New("the response did not contain a status")
		}
		p.resp = resp
		p.CurState = state
		return p.CurState, nil
	})
	if err != nil {
		return nil, err
	}
	return p.resp, nil
}

func (p *restorePoller[T]) Result(ctx context.Context, out *T) error {
	if p.resp.StatusCode == http.StatusNoContent {
		return nil
	} else if poller.Failed(p.CurState) {
		return runtime.NewResponseError(p.resp)
	}

	return resultHelper(p.resp, poller.Failed(p.CurState), out)
}

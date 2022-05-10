//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armloc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Kind is the identifier of this type in a resume token.
const kind = "ARM-Location"

// Applicable returns true if the LRO is using Location.
func Applicable(resp *http.Response) bool {
	return resp.StatusCode == http.StatusAccepted && resp.Header.Get(shared.HeaderLocation) != ""
}

// CanResume returns true if the token can rehydrate this poller type.
func CanResume(token map[string]interface{}) bool {
	t, ok := token["type"]
	if !ok {
		return false
	}
	tt, ok := t.(string)
	if !ok {
		return false
	}
	return tt == kind
}

// Poller is an LRO poller that uses the Location pattern.
type Poller[T any] struct {
	pl exported.Pipeline

	resp *http.Response

	// The poller's type, used for resume token processing.
	Type string `json:"type"`

	// The URL for polling.
	PollURL string `json:"pollURL"`

	// The LRO's current state.
	CurState string `json:"state"`
}

// New creates a new Poller from the provided initial response.
// Pass nil for response to create an empty Poller for rehydration.
func New[T any](pl exported.Pipeline, resp *http.Response) (*Poller[T], error) {
	if resp == nil {
		log.Write(log.EventLRO, "Resuming Location poller.")
		return &Poller[T]{pl: pl}, nil
	}
	log.Write(log.EventLRO, "Using Location poller.")
	locURL := resp.Header.Get(shared.HeaderLocation)
	if locURL == "" {
		return nil, errors.New("response is missing Location header")
	}
	if !pollers.IsValidURL(locURL) {
		return nil, fmt.Errorf("invalid polling URL %s", locURL)
	}
	p := &Poller[T]{
		pl:       pl,
		resp:     resp,
		Type:     kind,
		PollURL:  locURL,
		CurState: pollers.StatusInProgress,
	}
	return p, nil
}

// Done implements the Done method for the Operation interface.
func (p *Poller[T]) Done() bool {
	return pollers.IsTerminalState(p.CurState)
}

// Poll implements the Poll method for the Operation interface.
func (p *Poller[T]) Poll(ctx context.Context) (*http.Response, error) {
	err := pollers.PollHelper(ctx, p.PollURL, p.pl, func(resp *http.Response) (string, error) {
		// location polling can return an updated polling URL
		if h := resp.Header.Get(shared.HeaderLocation); h != "" {
			p.PollURL = h
		}
		if exported.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
			// if a 200/201 returns a provisioning state, use that instead
			state, err := pollers.GetProvisioningState(resp)
			if err != nil && !errors.Is(err, pollers.ErrNoBody) {
				return "", err
			}
			if state != "" {
				p.CurState = state
			} else {
				// a 200/201 with no provisioning state indicates success
				p.CurState = pollers.StatusSucceeded
			}
		} else if resp.StatusCode == http.StatusNoContent {
			p.CurState = pollers.StatusSucceeded
		} else if resp.StatusCode > 399 && resp.StatusCode < 500 {
			p.CurState = pollers.StatusFailed
		}
		p.resp = resp
		// a 202 falls through, means the LRO is still in progress and we don't check for provisioning state
		return p.CurState, nil
	})
	if err != nil {
		return nil, err
	}
	return p.resp, nil
}

// Result implements the Result method for the Operation interface.
func (p *Poller[T]) Result(ctx context.Context, out *T) error {
	return pollers.ResultHelper(p.resp, pollers.Failed(p.CurState), out)
}

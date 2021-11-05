//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package body

import (
	"errors"
	"net/http"

	armpollers "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Kind is the identifier of this type in a resume token.
const Kind = "Body"

// Applicable returns true if the LRO is using no headers, just provisioning state.
// This is only applicable to PATCH and PUT methods and assumes no polling headers.
func Applicable(resp *http.Response) bool {
	// we can't check for absense of headers due to some misbehaving services
	// like redis that return a Location header but don't actually use that protocol
	return resp.Request.Method == http.MethodPatch || resp.Request.Method == http.MethodPut
}

// Poller is an LRO poller that uses the Body pattern.
type Poller struct {
	// The poller's type, used for resume token processing.
	Type string `json:"type"`

	// The URL for polling.
	PollURL string `json:"pollURL"`

	// The LRO's current state.
	CurState string `json:"state"`
}

// New creates a new Poller from the provided initial response.
func New(resp *http.Response, pollerID string) (*Poller, error) {
	log.Write(log.EventLRO, "Using Body poller.")
	p := &Poller{
		Type:    pollers.MakeID(pollerID, Kind),
		PollURL: resp.Request.URL.String(),
	}
	// default initial state to InProgress.  depending on the HTTP
	// status code and provisioning state, we might change the value.
	curState := pollers.StatusInProgress
	provState, err := armpollers.GetProvisioningState(resp)
	if err != nil && !errors.Is(err, shared.ErrNoBody) {
		return nil, err
	}
	if resp.StatusCode == http.StatusCreated && provState != "" {
		// absense of provisioning state is ok for a 201, means the operation is in progress
		curState = provState
	} else if resp.StatusCode == http.StatusOK {
		if provState != "" {
			curState = provState
		} else if provState == "" {
			// for a 200, absense of provisioning state indicates success
			curState = pollers.StatusSucceeded
		}
	} else if resp.StatusCode == http.StatusNoContent {
		curState = pollers.StatusSucceeded
	}
	p.CurState = curState
	return p, nil
}

// URL returns the polling URL.
func (p *Poller) URL() string {
	return p.PollURL
}

// Done returns true if the LRO has reached a terminal state.
func (p *Poller) Done() bool {
	return pollers.IsTerminalState(p.Status())
}

// Update updates the Poller from the polling response.
func (p *Poller) Update(resp *http.Response) error {
	if resp.StatusCode == http.StatusNoContent {
		p.CurState = pollers.StatusSucceeded
		return nil
	}
	state, err := armpollers.GetProvisioningState(resp)
	if errors.Is(err, shared.ErrNoBody) {
		// a missing response body in non-204 case is an error
		return err
	} else if state == "" {
		// a response body without provisioning state is considered terminal success
		state = pollers.StatusSucceeded
	} else if err != nil {
		return err
	}
	p.CurState = state
	return nil
}

// FinalGetURL returns the empty string as no final GET is required for this poller type.
func (*Poller) FinalGetURL() string {
	return ""
}

// Status returns the status of the LRO.
func (p *Poller) Status() string {
	return p.CurState
}

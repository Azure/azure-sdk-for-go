// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package body

import (
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	stateSucceeded  = "Succeeded"
	stateInProgress = "InProgress"
)

// Applicable returns true if the LRO is using no headers, just provisioning state.
// This is only applicable to PATCH and PUT methods and assumes no polling headers.
func Applicable(resp *azcore.Response) bool {
	// we can't check for absense of headers due to some misbehaving services
	// like redis that return a Location header but don't actually use that protocol
	return resp.Request.Method == http.MethodPatch || resp.Request.Method == http.MethodPut
}

// Poller is an LRO poller that uses the Body pattern.
type Poller struct {
	Type     string `json:"type"`
	PollURL  string `json:"pollURL"`
	CurState string `json:"state"`
}

// New creates a new Poller from the provided initial response.
func New(resp *azcore.Response, pollerID string) (*Poller, error) {
	azcore.Log().Write(azcore.LogLongRunningOperation, "Using Body poller.")
	p := &Poller{
		Type:    pollers.MakeID(pollerID, "body"),
		PollURL: resp.Request.URL.String(),
	}
	// default initial state to InProgress.  depending on the HTTP
	// status code and provisioning state, we might change the value.
	curState := stateInProgress
	provState, err := pollers.GetProvisioningState(resp)
	if err != nil && !errors.Is(err, pollers.ErrNoBody) && !errors.Is(err, pollers.ErrNoProvisioningState) {
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
			curState = stateSucceeded
		}
	} else if resp.StatusCode == http.StatusNoContent {
		curState = stateSucceeded
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
func (p *Poller) Update(resp *azcore.Response) error {
	if resp.StatusCode == http.StatusNoContent {
		p.CurState = stateSucceeded
		return nil
	}
	state, err := pollers.GetProvisioningState(resp)
	if errors.Is(err, pollers.ErrNoBody) {
		// a missing response body in non-204 case is an error
		return err
	} else if errors.Is(err, pollers.ErrNoProvisioningState) {
		// a response body without provisioning state is considered terminal success
		state = stateSucceeded
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

// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package loc

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Applicable returns true if the LRO is using Location.
func Applicable(resp *azcore.Response) bool {
	return resp.StatusCode == http.StatusAccepted && resp.Header.Get(pollers.HeaderLocation) != ""
}

// Poller is an LRO poller that uses the Location pattern.
type Poller struct {
	Type     string `json:"type"`
	PollURL  string `json:"pollURL"`
	CurState string `json:"state"`
}

// New creates a new Poller from the provided initial response.
func New(resp *azcore.Response, pollerID string) (*Poller, error) {
	azcore.Log().Write(azcore.LogLongRunningOperation, "Using Location poller.")
	p := &Poller{
		Type:     pollers.MakeID(pollerID, "loc"),
		PollURL:  resp.Header.Get(pollers.HeaderLocation),
		CurState: "InProgress",
	}
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
	// location polling can return an updated polling URL
	if h := resp.Header.Get(pollers.HeaderLocation); h != "" {
		p.PollURL = h
	}
	// any 2xx code other than 202 indicates success
	if resp.HasStatusCode(http.StatusOK, http.StatusCreated, http.StatusNoContent) {
		p.CurState = "Succeeded"
	} else if resp.StatusCode > 399 && resp.StatusCode < 500 {
		p.CurState = "Failed"
	}
	return nil
}

// FinalGetURL returns the empty string as no final GET is required for this poller type.
func (p *Poller) FinalGetURL() string {
	return ""
}

// Status returns the status of the LRO.
func (p *Poller) Status() string {
	return p.CurState
}

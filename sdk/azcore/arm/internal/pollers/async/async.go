//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package async

import (
	"errors"
	"fmt"
	"net/http"

	armpollers "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Kind is the identifier of this type in a resume token.
const Kind = "Azure-AsyncOperation"

const (
	finalStateAsync = "azure-async-operation"
	finalStateLoc   = "location" //nolint
	finalStateOrig  = "original-uri"
)

// Applicable returns true if the LRO is using Azure-AsyncOperation.
func Applicable(resp *http.Response) bool {
	return resp.Header.Get(shared.HeaderAzureAsync) != ""
}

// Poller is an LRO poller that uses the Azure-AsyncOperation pattern.
type Poller struct {
	// The poller's type, used for resume token processing.
	Type string `json:"type"`

	// The URL from Azure-AsyncOperation header.
	AsyncURL string `json:"asyncURL"`

	// The URL from Location header.
	LocURL string `json:"locURL"`

	// The URL from the initial LRO request.
	OrigURL string `json:"origURL"`

	// The HTTP method from the initial LRO request.
	Method string `json:"method"`

	// The value of final-state-via from swagger, can be the empty string.
	FinalState string `json:"finalState"`

	// The LRO's current state.
	CurState string `json:"state"`
}

// New creates a new Poller from the provided initial response and final-state type.
func New(resp *http.Response, finalState string, pollerID string) (*Poller, error) {
	log.Write(log.EventLRO, "Using Azure-AsyncOperation poller.")
	asyncURL := resp.Header.Get(shared.HeaderAzureAsync)
	if asyncURL == "" {
		return nil, errors.New("response is missing Azure-AsyncOperation header")
	}
	if !pollers.IsValidURL(asyncURL) {
		return nil, fmt.Errorf("invalid polling URL %s", asyncURL)
	}
	p := &Poller{
		Type:       pollers.MakeID(pollerID, Kind),
		AsyncURL:   asyncURL,
		LocURL:     resp.Header.Get(shared.HeaderLocation),
		OrigURL:    resp.Request.URL.String(),
		Method:     resp.Request.Method,
		FinalState: finalState,
	}
	// check for provisioning state
	state, err := armpollers.GetProvisioningState(resp)
	if errors.Is(err, shared.ErrNoBody) || state == "" {
		// NOTE: the ARM RPC spec explicitly states that for async PUT the initial response MUST
		// contain a provisioning state.  to maintain compat with track 1 and other implementations
		// we are explicitly relaxing this requirement.
		/*if resp.Request.Method == http.MethodPut {
			// initial response for a PUT requires a provisioning state
			return nil, err
		}*/
		// for DELETE/PATCH/POST, provisioning state is optional
		state = pollers.StatusInProgress
	} else if err != nil {
		return nil, err
	}
	p.CurState = state
	return p, nil
}

// Done returns true if the LRO has reached a terminal state.
func (p *Poller) Done() bool {
	return pollers.IsTerminalState(p.Status())
}

// Update updates the Poller from the polling response.
func (p *Poller) Update(resp *http.Response) error {
	state, err := armpollers.GetStatus(resp)
	if err != nil {
		return err
	} else if state == "" {
		return errors.New("the response did not contain a status")
	}
	p.CurState = state
	return nil
}

// FinalGetURL returns the URL to perform a final GET for the payload, or the empty string if not required.
func (p *Poller) FinalGetURL() string {
	if p.Method == http.MethodPatch || p.Method == http.MethodPut {
		// for PATCH and PUT, the final GET is on the original resource URL
		return p.OrigURL
	} else if p.Method == http.MethodPost {
		if p.FinalState == finalStateAsync {
			return ""
		} else if p.FinalState == finalStateOrig {
			return p.OrigURL
		} else if p.LocURL != "" {
			// ideally FinalState would be set to "location" but it isn't always.
			// must check last due to more permissive condition.
			return p.LocURL
		}
	}
	return ""
}

// URL returns the polling URL.
func (p *Poller) URL() string {
	return p.AsyncURL
}

// Status returns the status of the LRO.
func (p *Poller) Status() string {
	return p.CurState
}

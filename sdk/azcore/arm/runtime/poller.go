//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// NewPoller creates a Poller based on the provided initial response.
// pollerID - a unique identifier for an LRO.  it's usually the client.Method string.
func NewPoller(pollerID string, finalState string, resp *http.Response, pl pipeline.Pipeline, eu func(*http.Response) error) (*pollers.Poller, error) {
	defer resp.Body.Close()
	// this is a back-stop in case the swagger is incorrect (i.e. missing one or more status codes for success).
	// ideally the codegen should return an error if the initial response failed and not even create a poller.
	if !pollers.StatusCodeValid(resp) {
		return nil, errors.New("the LRO failed or was cancelled")
	}
	// determine the polling method
	var lro pollers.Operation
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
		lro = &pollers.NopPoller{}
	}
	if err != nil {
		return nil, err
	}
	return pollers.NewPoller(lro, resp, pl, eu), nil
}

// NewPollerFromResumeToken creates a Poller from a resume token string.
// pollerID - a unique identifier for an LRO.  it's usually the client.Method string.
func NewPollerFromResumeToken(pollerID string, token string, pl pipeline.Pipeline, eu func(*http.Response) error) (*pollers.Poller, error) {
	kind, err := pollers.KindFromToken(pollerID, token)
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
	return pollers.NewPoller(lro, nil, pl, eu), nil
}

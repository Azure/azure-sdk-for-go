//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

// NOTE: this file will be deleted in a future release

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewPollerOptions contains the optional parameters for NewPoller.
type NewPollerOptions[T any] struct {
	// FinalStateVia contains the final-state-via value for the LRO.
	FinalStateVia runtime.FinalStateVia

	// Response contains a preconstructed response type.
	// The final payload will be unmarshaled into it and returned.
	Response *T
}

// NewPoller creates a Poller based on the provided initial response.
func NewPoller[T any](resp *http.Response, pl runtime.Pipeline, options *NewPollerOptions[T]) (*runtime.Poller[T], error) {
	if options == nil {
		options = &NewPollerOptions[T]{}
	}
	defer resp.Body.Close()
	// this is a back-stop in case the swagger is incorrect (i.e. missing one or more status codes for success).
	// ideally the codegen should return an error if the initial response failed and not even create a poller.
	if !pollers.StatusCodeValid(resp) {
		return nil, errors.New("the LRO failed or was cancelled")
	}

	// determine the polling method
	var op runtime.PollerMethod[T]
	var err error
	if async.Applicable(resp) {
		op, err = async.New[T](pl, resp, options.FinalStateVia)
	} else if loc.Applicable(resp) {
		op, err = loc.New[T](pl, resp)
	} else if body.Applicable(resp) {
		// must test body poller last as it's a subset of the other pollers.
		// TODO: this is ambiguous for PATCH/PUT if it returns a 200 with no polling headers (sync completion)
		op, err = body.New[T](pl, resp)
	} else if m := resp.Request.Method; resp.StatusCode == http.StatusAccepted && (m == http.MethodDelete || m == http.MethodPost) {
		// if we get here it means we have a 202 with no polling headers.
		// for DELETE and POST this is a hard error per ARM RPC spec.
		return nil, errors.New("response is missing polling URL")
	} else {
		op, err = pollers.NewNopPoller[T](resp)
	}
	if err != nil {
		return nil, err
	}
	return runtime.NewPoller(resp, pl, &runtime.NewPollerOptions[T]{
		FinalStateVia: options.FinalStateVia,
		Response:      options.Response,
		PollerMethod:  op,
	})
}

// NewPollerFromResumeTokenOptions contains the optional parameters for NewPollerFromResumeToken.
type NewPollerFromResumeTokenOptions[T any] struct {
	// Response contains a preconstructed response type.
	// The final payload will be unmarshaled into it and returned.
	Response *T
}

// NewPollerFromResumeToken creates a Poller from a resume token string.
func NewPollerFromResumeToken[T any](token string, pl runtime.Pipeline, options *NewPollerFromResumeTokenOptions[T]) (*runtime.Poller[T], error) {
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

	// now rehydrate the poller based on the encoded poller type
	var op runtime.PollerMethod[T]
	if async.CanResume(asJSON) {
		op, _ = async.New[T](pl, nil, "")
	} else if loc.CanResume(asJSON) {
		op, _ = loc.New[T](pl, nil)
	} else if body.CanResume(asJSON) {
		op, _ = body.New[T](pl, nil)
	} else {
		return nil, fmt.Errorf("unhandled poller token %+v", asJSON)
	}
	if err := json.Unmarshal(raw, &op); err != nil {
		return nil, err
	}
	return runtime.NewPoller(nil, pl, &runtime.NewPollerOptions[T]{
		Response:     options.Response,
		PollerMethod: op,
	})
}

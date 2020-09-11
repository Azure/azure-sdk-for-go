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
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	headerAsyncOperation = "Azure-AsyncOperation"
	headerLocation       = "Location"
)

const (
	operationInProgress string = "InProgress"
	operationCanceled   string = "Canceled"
	operationFailed     string = "Failed"
	operationSucceeded  string = "Succeeded"
)

var pollingCodes = [...]int{http.StatusNoContent, http.StatusAccepted, http.StatusCreated, http.StatusOK}

// NewPoller creates a polling tracker based on the verb of the original request and returns
// the polling tracker implementation for the method verb or an error.
// NOTE: this is only meant for internal use in generated code.
func NewPoller(pollerType string, finalState string, resp *azcore.Response, errorHandler methodErrorHandler) (Poller, error) {
	switch strings.ToUpper(resp.Request.Method) {
	case http.MethodDelete:
		pt := &pollingTrackerDelete{pollingTrackerBase: pollingTrackerBase{PollerType: pollerType, FinalStateVia: finalState, OriginalURI: resp.Request.URL.String(), resp: resp, errorHandler: errorHandler}}
		if err := pt.initializeState(); err != nil {
			return pt, err
		}
		// this initializes the polling header values, we do this during creation in case the
		// initial response send us invalid values; this way the API call will return a non-nil
		// error (not doing this means the error shows up in Future.Done)
		return pt, pt.updatePollingMethod()
	case http.MethodPatch:
		pt := &pollingTrackerPatch{pollingTrackerBase: pollingTrackerBase{PollerType: pollerType, FinalStateVia: finalState, OriginalURI: resp.Request.URL.String(), resp: resp, errorHandler: errorHandler}}
		if err := pt.initializeState(); err != nil {
			return pt, err
		}
		// this initializes the polling header values, we do this during creation in case the
		// initial response send us invalid values; this way the API call will return a non-nil
		// error (not doing this means the error shows up in Future.Done)
		return pt, pt.updatePollingMethod()
	case http.MethodPost:
		pt := &pollingTrackerPost{pollingTrackerBase: pollingTrackerBase{PollerType: pollerType, FinalStateVia: finalState, OriginalURI: resp.Request.URL.String(), resp: resp, errorHandler: errorHandler}}
		if err := pt.initializeState(); err != nil {
			return pt, err
		}
		// this initializes the polling header values, we do this during creation in case the
		// initial response send us invalid values; this way the API call will return a non-nil
		// error (not doing this means the error shows up in Future.Done)
		return pt, pt.updatePollingMethod()
	case http.MethodPut:
		pt := &pollingTrackerPut{pollingTrackerBase: pollingTrackerBase{PollerType: pollerType, FinalStateVia: finalState, OriginalURI: resp.Request.URL.String(), resp: resp, errorHandler: errorHandler}}
		if err := pt.initializeState(); err != nil {
			return pt, err
		}
		// this initializes the polling header values, we do this during creation in case the
		// initial response send us invalid values; this way the API call will return a non-nil
		// error (not doing this means the error shows up in Future.Done)
		return pt, pt.updatePollingMethod()
	default:
		return nil, fmt.Errorf("unsupported HTTP method %s", resp.Request.Method)
	}
}

// NewPollerFromResumeToken creates a polling tracker from a resume token string.
// NOTE: this is only meant for internal use in generated code.
func NewPollerFromResumeToken(pollerType string, token string, errorHandler methodErrorHandler) (Poller, error) {
	// unmarshal into JSON object to determine the tracker type
	obj := map[string]interface{}{}
	err := json.Unmarshal([]byte(token), &obj)
	if err != nil {
		return nil, err
	}
	if obj["pollerType"] != pollerType {
		return nil, fmt.Errorf("Cannot resume from this poller type. Expected: %s, Received: %s", pollerType, obj["pollerType"])
	}
	if obj["method"] == nil {
		return nil, fmt.Errorf("Token is missing 'method' property")
	}
	method := obj["method"].(string)
	switch strings.ToUpper(method) {
	case http.MethodDelete:
		pt := &pollingTrackerDelete{pollingTrackerBase: pollingTrackerBase{errorHandler: errorHandler}}
		// now unmarshal into the tracker
		err = json.Unmarshal([]byte(token), &pt)
		if err != nil {
			return nil, err
		}
		return pt, nil
	case http.MethodPatch:
		pt := &pollingTrackerPatch{pollingTrackerBase: pollingTrackerBase{errorHandler: errorHandler}}
		// now unmarshal into the tracker
		err = json.Unmarshal([]byte(token), &pt)
		if err != nil {
			return nil, err
		}
		return pt, nil
	case http.MethodPost:
		pt := &pollingTrackerPost{pollingTrackerBase: pollingTrackerBase{errorHandler: errorHandler}}
		// now unmarshal into the tracker
		err = json.Unmarshal([]byte(token), &pt)
		if err != nil {
			return nil, err
		}
		return pt, nil
	case http.MethodPut:
		pt := &pollingTrackerPut{pollingTrackerBase: pollingTrackerBase{errorHandler: errorHandler}}
		// now unmarshal into the tracker
		err = json.Unmarshal([]byte(token), &pt)
		if err != nil {
			return nil, err
		}
		return pt, nil
	default:
		return nil, fmt.Errorf("unsupported method '%s'", method)
	}
}

// Poller defines the methods that will be called internally in the generated code for long-running operations.
// NOTE: this is only meant for internal use in generated code.
type Poller interface {
	// Done signals if the polling operation has reached a terminal state.
	Done() bool
	// Poll sends a polling request to the service endpoint and returns the http.Response received from the endpoint or an error.
	Poll(ctx context.Context, p azcore.Pipeline) (*http.Response, error)
	// FinalResponse will perform a final GET and return the final http response for the polling operation and unmarshal the content of the payload into the respType interface that is provided.
	FinalResponse(ctx context.Context, pipeline azcore.Pipeline, respType interface{}) (*http.Response, error)
	// ResumeToken returns a token string that can be used to resume polling on a poller that has not yet reached a terminal state.
	ResumeToken() (string, error)
	// PollUntilDone will handle the entire span of the polling operation until a terminal state is reached, then return the final http response for the polling operation and unmarshal the content of the payload into the respType interface that is provided.
	PollUntilDone(ctx context.Context, frequency time.Duration, pipeline azcore.Pipeline, respType interface{}) (*http.Response, error)
}

// Done returns true if there was an error or polling has reached a terminal state.
func (pt *pollingTrackerBase) Done() bool {
	return pt.hasTerminated()
}

// FinalResponse will perform a final GET and return the final response for the polling operation.
func (pt *pollingTrackerBase) FinalResponse(ctx context.Context, pipeline azcore.Pipeline, respType interface{}) (*http.Response, error) {
	if !pt.Done() {
		return nil, errors.New("cannot return a final response from a poller in a non-terminal state")
	}
	// if respType is nil, this indicates that the request was made from an HTTPPoller
	if respType == nil {
		return pt.latestResponse().Response, nil
	}
	if pt.pollerMethodVerb() == http.MethodPut || pt.pollerMethodVerb() == http.MethodPatch {
		res, err := pt.handleResponse(pt.latestResponse(), respType)
		if err != nil {
			return nil, err
		}
		if res != nil && !reflect.Indirect(reflect.ValueOf(respType)).IsZero() {
			return res, nil
		}
	}
	// checking if there was a FinalStateVia configuration to re-route the final GET
	// request to the value specified in the FinalStateVia property on the poller
	err := pt.setFinalState()
	if err != nil {
		return nil, err
	}
	if pt.finalGetURL() == "" {
		// we can end up in this situation if the async operation returns a 200
		// with no polling URLs.  in that case return the response which should
		// contain the JSON payload (only do this for successful terminal cases).
		if lr := pt.latestResponse(); lr != nil && pt.hasSucceeded() {
			result, err := pt.handleResponse(lr, respType)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, errors.New("missing URL for retrieving result")
	}
	req, err := azcore.NewRequest(ctx, http.MethodGet, pt.finalGetURL())
	if err != nil {
		return nil, err
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		return nil, err
	}
	return pt.handleResponse(resp, respType)
}

// ResumeToken generates the string token that can be used with the Resume<PollerType> method
// on the client to create a new poller from the data held in the current poller type.
func (pt *pollingTrackerBase) ResumeToken() (string, error) {
	if pt.hasTerminated() {
		return "", errors.New("cannot create a ResumeToken from a poller in a terminal state")
	}
	js, err := json.Marshal(pt)
	if err != nil {
		return "", err
	}
	return string(js), nil
}

func (pt *pollingTrackerBase) handleResponse(resp *azcore.Response, respType interface{}) (*http.Response, error) {
	if resp.HasStatusCode(http.StatusNoContent) {
		return resp.Response, nil
	}
	if !resp.HasStatusCode(pollingCodes[:]...) {
		return nil, pt.handleError(resp)
	}
	return resp.Response, resp.UnmarshalAsJSON(respType)
}

func delay(ctx context.Context, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// done queries the service to see if the operation has completed.
func lroPollDone(ctx context.Context, p azcore.Pipeline, pt pollingTracker) bool {
	if pt.hasTerminated() {
		return true
	}
	if err := pt.pollForStatus(ctx, p); err != nil {
		return true
	}
	if err := pt.checkForErrors(); err != nil {
		return true
	}
	if err := pt.updatePollingState(pt.provisioningStateApplicable()); err != nil {
		return true
	}
	if err := pt.initPollingMethod(); err != nil {
		return true
	}
	if err := pt.updatePollingMethod(); err != nil {
		return true
	}
	return pt.hasTerminated()
}

type pollingTracker interface {
	// these methods can differ per tracker

	// checks the response headers and status code to determine the polling mechanism
	updatePollingMethod() error

	// checks the response for tracker-specific error conditions
	checkForErrors() error

	// returns true if provisioning state should be checked
	provisioningStateApplicable() bool

	// methods common to all trackers

	// initializes a tracker's polling URL and method, called for each iteration.
	// these values can be overridden by each polling tracker as required.
	initPollingMethod() error

	// initializes the tracker's internal state, call this when the tracker is created
	initializeState() error

	// makes an HTTP request to check the status of the LRO
	pollForStatus(ctx context.Context, client azcore.Pipeline) error

	// updates internal tracker state, call this after each call to pollForStatus
	updatePollingState(provStateApl bool) error

	// returns the error response from the service, can be nil
	pollingError() error

	// returns the polling method being used
	pollingMethod() pollingMethodType

	// returns the state of the LRO as returned from the service
	pollingStatus() string

	// returns the URL used for polling status
	pollingURL() string

	// returns the URL used for the final GET to retrieve the resource
	finalGetURL() string

	// returns true if the LRO is in a terminal state
	hasTerminated() bool

	// returns true if the LRO is in a failed terminal state
	hasFailed() bool

	// returns true if the LRO is in a successful terminal state
	hasSucceeded() bool

	// returns the cached HTTP response after a call to pollForStatus(), can be nil
	latestResponse() *azcore.Response

	// converts an *azcore.Response to an error
	handleError(resp *azcore.Response) error

	// sets the FinalGetURI to the value pointed to in FinalStateVia
	setFinalState() error

	// returns the verb used with the initial request
	pollerMethodVerb() string
}

type methodErrorHandler func(resp *azcore.Response) error

type pollingTrackerBase struct {
	// resp is the last response, either from the submission of the LRO or from polling
	resp *azcore.Response

	// PollerType is the name of the type of poller that is created
	PollerType string `json:"pollerType"`

	// errorHandler is the method to invoke to unmarshall an error response
	errorHandler methodErrorHandler

	// method is the HTTP verb, this is needed for deserialization
	Method string `json:"method"`

	// rawBody is the raw JSON response body
	rawBody map[string]interface{}

	// denotes if polling is using async-operation or location header
	Pm pollingMethodType `json:"pollingMethod"`

	// the URL to poll for status
	URI string `json:"pollingURI"`

	// the state of the LRO as returned from the service
	State string `json:"lroState"`

	// the URL to GET for the final result
	FinalGetURI string `json:"resultURI"`

	// stores the name of the header that the final get should be performed on,
	// can be empty which will go to default behavior
	FinalStateVia string `json:"finalStateVia"`

	// the original request URL of the initial request for the polling operation
	OriginalURI string `json:"originalURI"`

	// used to hold an error object returned from the service
	Err error `json:"error,omitempty"`
}

func (pt *pollingTrackerBase) initializeState() error {
	// determine the initial polling state based on response body and/or HTTP status
	// code.  this is applicable to the initial LRO response, not polling responses!
	pt.Method = pt.resp.Request.Method
	if err := pt.updateRawBody(); err != nil {
		pt.Err = err
		return err
	}
	switch pt.resp.StatusCode {
	case http.StatusOK:
		if ps := pt.getProvisioningState(); ps != nil {
			pt.State = *ps
			if pt.hasFailed() {
				pt.updateErrorFromResponse()
				return pt.pollingError()
			}
		} else {
			pt.State = operationSucceeded
		}
	case http.StatusCreated:
		if ps := pt.getProvisioningState(); ps != nil {
			pt.State = *ps
		} else {
			pt.State = operationInProgress
		}
	case http.StatusAccepted:
		pt.State = operationInProgress
	case http.StatusNoContent:
		pt.State = operationSucceeded
	default:
		pt.State = operationFailed
		pt.updateErrorFromResponse()
		return pt.pollingError()
	}
	return pt.initPollingMethod()
}

func (pt *pollingTrackerBase) getProvisioningState() *string {
	if pt.rawBody != nil && pt.rawBody["properties"] != nil {
		p := pt.rawBody["properties"].(map[string]interface{})
		if ps := p["provisioningState"]; ps != nil {
			s := ps.(string)
			return &s
		}
	}
	return nil
}

func (pt *pollingTrackerBase) updateRawBody() error {
	pt.rawBody = map[string]interface{}{}
	if pt.resp.ContentLength != 0 {
		defer pt.resp.Body.Close()
		b, err := ioutil.ReadAll(pt.resp.Body)
		if err != nil {
			pt.Err = err
			return pt.Err
		}
		// observed in 204 responses over HTTP/2.0; the content length is -1 but body is empty
		if len(b) == 0 {
			return nil
		}
		if err = json.Unmarshal(b, &pt.rawBody); err != nil {
			pt.Err = err
			return pt.Err
		}
	}
	return nil
}

func (pt *pollingTrackerBase) pollForStatus(ctx context.Context, client azcore.Pipeline) error {
	req, err := azcore.NewRequest(ctx, http.MethodGet, pt.URI)
	if err != nil {
		pt.Err = err
		return err
	}
	resp, err := client.Do(req)
	pt.resp = resp
	if err != nil {
		pt.Err = err
		return pt.Err
	}
	if pt.resp.HasStatusCode(pollingCodes[:]...) {
		// reset the service error on success case
		pt.Err = nil
		err = pt.updateRawBody()
	} else {
		// check response body for error content
		pt.updateErrorFromResponse()
		err = pt.pollingError()
	}
	return err
}

// attempts to unmarshal a ServiceError type from the response body.
// if that fails then make a best attempt at creating something meaningful.
// NOTE: this assumes that the async operation has failed.
func (pt *pollingTrackerBase) updateErrorFromResponse() {
	pt.Err = pt.errorHandler(pt.resp)
}

func (pt *pollingTrackerBase) updatePollingState(provStateApl bool) error {
	if pt.Pm == pollingAsyncOperation && pt.rawBody["status"] != nil {
		pt.State = pt.rawBody["status"].(string)
	} else {
		if pt.resp.StatusCode == http.StatusAccepted {
			pt.State = operationInProgress
		} else if provStateApl {
			if ps := pt.getProvisioningState(); ps != nil {
				pt.State = *ps
			} else {
				pt.State = operationSucceeded
			}
		} else {
			pt.Err = fmt.Errorf("the response from the async operation has an invalid status code: %d", pt.resp.StatusCode)
			return pt.Err
		}
	}
	// if the operation has failed update the error state
	if pt.hasFailed() {
		pt.updateErrorFromResponse()
	}
	return nil
}

func (pt *pollingTrackerBase) pollingError() error {
	return pt.Err
}

func (pt *pollingTrackerBase) pollingMethod() pollingMethodType {
	return pt.Pm
}

func (pt *pollingTrackerBase) pollingStatus() string {
	return pt.State
}

func (pt *pollingTrackerBase) pollingURL() string {
	return pt.URI
}

func (pt *pollingTrackerBase) finalGetURL() string {
	return pt.FinalGetURI
}

func (pt *pollingTrackerBase) hasTerminated() bool {
	return strings.EqualFold(pt.State, operationCanceled) || strings.EqualFold(pt.State, operationFailed) || strings.EqualFold(pt.State, operationSucceeded)
}

func (pt *pollingTrackerBase) hasFailed() bool {
	return strings.EqualFold(pt.State, operationCanceled) || strings.EqualFold(pt.State, operationFailed)
}

func (pt *pollingTrackerBase) hasSucceeded() bool {
	return strings.EqualFold(pt.State, operationSucceeded)
}

func (pt *pollingTrackerBase) latestResponse() *azcore.Response {
	return pt.resp
}

// error checking common to all trackers
func (pt *pollingTrackerBase) baseCheckForErrors() error {
	// for Azure-AsyncOperations the response body cannot be nil or empty
	if pt.Pm == pollingAsyncOperation {
		if pt.resp.Body == nil || pt.resp.ContentLength == 0 {
			pt.Err = errors.New("for Azure-AsyncOperation response body cannot be nil")
			return pt.Err
		}
		if pt.rawBody["status"] == nil {
			pt.Err = errors.New("missing status property in Azure-AsyncOperation response body")
			return pt.Err
		}
	}
	return nil
}

// default initialization of polling URL/method.  each verb tracker will update this as required.
func (pt *pollingTrackerBase) initPollingMethod() error {
	if ao, err := getURLFromAsyncOpHeader(pt.resp); err != nil {
		pt.Err = err
		return err
	} else if ao != "" {
		pt.URI = ao
		pt.Pm = pollingAsyncOperation
		return nil
	}
	if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
		pt.Err = err
		return err
	} else if lh != "" {
		pt.URI = lh
		pt.Pm = pollingLocation
		return nil
	}
	// it's ok if we didn't find a polling header, this will be handled elsewhere
	return nil
}

func (pt *pollingTrackerBase) handleError(resp *azcore.Response) error {
	return pt.errorHandler(resp)
}

func (pt *pollingTrackerBase) setFinalState() error {
	if len(pt.FinalStateVia) == 0 {
		return nil
	}
	if pt.FinalStateVia == "azure-async-operation" {
		ao, err := getURLFromAsyncOpHeader(pt.latestResponse())
		if err != nil {
			return err
		}
		pt.FinalGetURI = ao
	} else if pt.FinalStateVia == "location" {
		lh, err := getURLFromLocationHeader(pt.latestResponse())
		if err != nil {
			return err
		}
		pt.FinalGetURI = lh
	} else if pt.FinalStateVia == "original-uri" {
		pt.FinalGetURI = pt.OriginalURI
	}
	return nil
}

func (pt *pollingTrackerBase) pollerMethodVerb() string {
	return pt.Method
}

// DELETE

type pollingTrackerDelete struct {
	pollingTrackerBase
}

// Poll will send poll the service endpoint and return an http.Response or error received from the service.
func (pt *pollingTrackerDelete) Poll(ctx context.Context, pipeline azcore.Pipeline) (*http.Response, error) {
	if lroPollDone(ctx, pipeline, pt) {
		return pt.latestResponse().Response, pt.pollingError()
	}
	return nil, pt.pollingError()
}

// PollUntilDone will handle the entire span of the polling operation until a terminal state is reached.
// It will return the final http.Response received or an error, it will also unmarshal the response payload
// in the respType interface that is passed into the function.
func (pt *pollingTrackerDelete) PollUntilDone(ctx context.Context, frequency time.Duration, pipeline azcore.Pipeline, respType interface{}) (*http.Response, error) {
	// initial check for a retry-after header existing on the initial response
	if retryAfter := azcore.RetryAfter(pt.latestResponse().Response); retryAfter > 0 {
		err := delay(ctx, retryAfter)
		if err != nil {
			return nil, err
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := pt.Poll(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		if pt.Done() {
			break
		}
		d := frequency
		if retryAfter := azcore.RetryAfter(resp); retryAfter > 0 {
			d = retryAfter
		}
		err = delay(ctx, d)
		if err != nil {
			return nil, err
		}
	}
	return pt.FinalResponse(ctx, pipeline, respType)
}

func (pt *pollingTrackerDelete) updatePollingMethod() error {
	// for 201 the Location header is required
	if pt.resp.StatusCode == http.StatusCreated {
		if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
			pt.Err = err
			return err
		} else if lh == "" {
			pt.Err = errors.New("missing Location header in 201 response")
			return pt.Err
		} else {
			pt.URI = lh
		}
		pt.Pm = pollingLocation
		pt.FinalGetURI = pt.URI
	}
	// for 202 prefer the Azure-AsyncOperation header but fall back to Location if necessary
	if pt.resp.StatusCode == http.StatusAccepted {
		ao, err := getURLFromAsyncOpHeader(pt.resp)
		if err != nil {
			pt.Err = err
			return err
		} else if ao != "" {
			pt.URI = ao
			pt.Pm = pollingAsyncOperation
		}
		// if the Location header is invalid and we already have a polling URL
		// then we don't care if the Location header URL is malformed.
		if lh, err := getURLFromLocationHeader(pt.resp); err != nil && pt.URI == "" {
			pt.Err = err
			return err
		} else if lh != "" {
			if ao == "" {
				pt.URI = lh
				pt.Pm = pollingLocation
			}
			// when both headers are returned we use the value in the Location header for the final GET
			pt.FinalGetURI = lh
		}
		// make sure a polling URL was found
		if pt.URI == "" {
			pt.Err = errors.New("didn't get any suitable polling URLs in 202 response")
			return pt.Err
		}
	}
	return nil
}

func (pt *pollingTrackerDelete) checkForErrors() error {
	return pt.baseCheckForErrors()
}

func (pt *pollingTrackerDelete) provisioningStateApplicable() bool {
	return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusNoContent
}

// PATCH

type pollingTrackerPatch struct {
	pollingTrackerBase
}

// Poll will send poll the service endpoint and return an http.Response or error received from the service.
func (pt *pollingTrackerPatch) Poll(ctx context.Context, pipeline azcore.Pipeline) (*http.Response, error) {
	if lroPollDone(ctx, pipeline, pt) {
		return pt.latestResponse().Response, pt.pollingError()
	}
	return nil, pt.pollingError()
}

// PollUntilDone will handle the entire span of the polling operation until a terminal state is reached.
// It will return the final http.Response received or an error, it will also unmarshal the response payload
// in the respType interface that is passed into the function.
func (pt *pollingTrackerPatch) PollUntilDone(ctx context.Context, frequency time.Duration, pipeline azcore.Pipeline, respType interface{}) (*http.Response, error) {
	// initial check for a retry-after header existing on the initial response
	if retryAfter := azcore.RetryAfter(pt.latestResponse().Response); retryAfter > 0 {
		err := delay(ctx, retryAfter)
		if err != nil {
			return nil, err
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := pt.Poll(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		if pt.Done() {
			break
		}
		d := frequency
		if retryAfter := azcore.RetryAfter(resp); retryAfter > 0 {
			d = retryAfter
		}
		err = delay(ctx, d)
		if err != nil {
			return nil, err
		}
	}
	return pt.FinalResponse(ctx, pipeline, respType)
}

func (pt *pollingTrackerPatch) updatePollingMethod() error {
	// by default we can use the original URL for polling and final GET
	if pt.URI == "" {
		pt.URI = pt.resp.Request.URL.String()
	}
	if pt.FinalGetURI == "" {
		pt.FinalGetURI = pt.resp.Request.URL.String()
	}
	if pt.Pm == pollingUnknown {
		pt.Pm = pollingRequestURI
	}
	// for 201 it's permissible for no headers to be returned
	if pt.resp.StatusCode == http.StatusCreated {
		if ao, err := getURLFromAsyncOpHeader(pt.resp); err != nil {
			pt.Err = err
			return err
		} else if ao != "" {
			pt.URI = ao
			pt.Pm = pollingAsyncOperation
		}
	}
	// for 202 prefer the Azure-AsyncOperation header but fall back to Location if necessary
	// note the absence of the "final GET" mechanism for PATCH
	if pt.resp.StatusCode == http.StatusAccepted {
		ao, err := getURLFromAsyncOpHeader(pt.resp)
		if err != nil {
			pt.Err = err
			return err
		} else if ao != "" {
			pt.URI = ao
			pt.Pm = pollingAsyncOperation
		}
		if ao == "" {
			if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
				pt.Err = err
				return err
			} else if lh == "" {
				pt.Err = errors.New("didn't get any suitable polling URLs in 202 response")
				return pt.Err
			} else {
				pt.URI = lh
				pt.Pm = pollingLocation
			}
		}
	}
	return nil
}

func (pt *pollingTrackerPatch) checkForErrors() error {
	return pt.baseCheckForErrors()
}

func (pt *pollingTrackerPatch) provisioningStateApplicable() bool {
	return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusCreated
}

// POST

type pollingTrackerPost struct {
	pollingTrackerBase
}

// Poll will send poll the service endpoint and return an http.Response or error received from the service.
func (pt *pollingTrackerPost) Poll(ctx context.Context, pipeline azcore.Pipeline) (*http.Response, error) {
	if lroPollDone(ctx, pipeline, pt) {
		return pt.latestResponse().Response, pt.pollingError()
	}
	return nil, pt.pollingError()
}

// PollUntilDone will handle the entire span of the polling operation until a terminal state is reached.
// It will return the final http.Response received or an error, it will also unmarshal the response payload
// in the respType interface that is passed into the function.
func (pt *pollingTrackerPost) PollUntilDone(ctx context.Context, frequency time.Duration, pipeline azcore.Pipeline, respType interface{}) (*http.Response, error) {
	// initial check for a retry-after header existing on the initial response
	if retryAfter := azcore.RetryAfter(pt.latestResponse().Response); retryAfter > 0 {
		err := delay(ctx, retryAfter)
		if err != nil {
			return nil, err
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := pt.Poll(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		if pt.Done() {
			break
		}
		d := frequency
		if retryAfter := azcore.RetryAfter(resp); retryAfter > 0 {
			d = retryAfter
		}
		err = delay(ctx, d)
		if err != nil {
			return nil, err
		}
	}
	return pt.FinalResponse(ctx, pipeline, respType)
}

func (pt *pollingTrackerPost) updatePollingMethod() error {
	// 201 requires Location header
	if pt.resp.StatusCode == http.StatusCreated {
		if lh, err := getURLFromLocationHeader(pt.resp); err != nil {
			pt.Err = err
			return err
		} else if lh == "" {
			pt.Err = errors.New("missing Location header in 201 response")
			return pt.Err
		} else {
			pt.URI = lh
			pt.FinalGetURI = lh
			pt.Pm = pollingLocation
		}
	}
	// for 202 prefer the Azure-AsyncOperation header but fall back to Location if necessary
	if pt.resp.StatusCode == http.StatusAccepted {
		ao, err := getURLFromAsyncOpHeader(pt.resp)
		if err != nil {
			pt.Err = err
			return err
		} else if ao != "" {
			pt.URI = ao
			pt.Pm = pollingAsyncOperation
		}
		// if the Location header is invalid and we already have a polling URL
		// then we don't care if the Location header URL is malformed.
		if lh, err := getURLFromLocationHeader(pt.resp); err != nil && pt.URI == "" {
			pt.Err = err
			return err
		} else if lh != "" {
			if ao == "" {
				pt.URI = lh
				pt.Pm = pollingLocation
			}
			// when both headers are returned we use the value in the Location header for the final GET
			pt.FinalGetURI = lh
		}
		// make sure a polling URL was found
		if pt.URI == "" {
			pt.Err = errors.New("didn't get any suitable polling URLs in 202 response")
			return pt.Err
		}
	}
	return nil
}

func (pt *pollingTrackerPost) checkForErrors() error {
	return pt.baseCheckForErrors()
}

func (pt *pollingTrackerPost) provisioningStateApplicable() bool {
	return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusNoContent
}

// PUT

type pollingTrackerPut struct {
	pollingTrackerBase
}

// Poll will send poll the service endpoint and return an http.Response or error received from the service.
func (pt *pollingTrackerPut) Poll(ctx context.Context, pipeline azcore.Pipeline) (*http.Response, error) {
	if lroPollDone(ctx, pipeline, pt) {
		return pt.latestResponse().Response, pt.pollingError()
	}
	return nil, pt.pollingError()
}

// PollUntilDone will handle the entire span of the polling operation until a terminal state is reached.
// It will return the final http.Response received or an error, it will also unmarshal the response payload
// in the respType interface that is passed into the function.
func (pt *pollingTrackerPut) PollUntilDone(ctx context.Context, frequency time.Duration, pipeline azcore.Pipeline, respType interface{}) (*http.Response, error) {
	// p.resp should only be nil when calling PollUntilDone from a poller that was instantiated from a resume token string
	if pt.resp != nil {
		// initial check for a retry-after header existing on the initial response
		if retryAfter := azcore.RetryAfter(pt.latestResponse().Response); retryAfter > 0 {
			err := delay(ctx, retryAfter)
			if err != nil {
				return nil, err
			}
		}
	}
	// begin polling the endpoint until a terminal state is reached
	for {
		resp, err := pt.Poll(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		if pt.Done() {
			break
		}
		d := frequency
		if retryAfter := azcore.RetryAfter(resp); retryAfter > 0 {
			d = retryAfter
		}
		err = delay(ctx, d)
		if err != nil {
			return nil, err
		}
	}
	return pt.FinalResponse(ctx, pipeline, respType)
}

func (pt *pollingTrackerPut) updatePollingMethod() error {
	// by default we can use the original URL for polling and final GET
	if pt.URI == "" {
		pt.URI = pt.resp.Request.URL.String()
	}
	if pt.FinalGetURI == "" {
		pt.FinalGetURI = pt.resp.Request.URL.String()
	}
	if pt.Pm == pollingUnknown {
		pt.Pm = pollingRequestURI
	}
	// for 201 it's permissible for no headers to be returned
	if pt.resp.StatusCode == http.StatusCreated {
		if ao, err := getURLFromAsyncOpHeader(pt.resp); err != nil {
			pt.Err = err
			return err
		} else if ao != "" {
			pt.URI = ao
			pt.Pm = pollingAsyncOperation
		}
	}
	// for 202 prefer the Azure-AsyncOperation header but fall back to Location if necessary
	if pt.resp.StatusCode == http.StatusAccepted {
		ao, err := getURLFromAsyncOpHeader(pt.resp)
		if err != nil {
			pt.Err = err
			return err
		} else if ao != "" {
			pt.URI = ao
			pt.Pm = pollingAsyncOperation
		}
		// if the Location header is invalid and we already have a polling URL
		// then we don't care if the Location header URL is malformed.
		if lh, err := getURLFromLocationHeader(pt.resp); err != nil && pt.URI == "" {
			pt.Err = err
			return err
		} else if lh != "" {
			if ao == "" {
				pt.URI = lh
				pt.Pm = pollingLocation
			}
		}
		// make sure a polling URL was found
		if pt.URI == "" {
			pt.Err = errors.New("didn't get any suitable polling URLs in 202 response")
			return pt.Err
		}
	}
	return nil
}

func (pt *pollingTrackerPut) checkForErrors() error {
	err := pt.baseCheckForErrors()
	if err != nil {
		pt.Err = err
		return err
	}
	// if there are no LRO headers then the body cannot be empty
	ao, err := getURLFromAsyncOpHeader(pt.resp)
	if err != nil {
		pt.Err = err
		return err
	}
	lh, err := getURLFromLocationHeader(pt.resp)
	if err != nil {
		pt.Err = err
		return err
	}
	if ao == "" && lh == "" && len(pt.rawBody) == 0 {
		pt.Err = errors.New("the response did not contain a body")
		return pt.Err
	}
	return nil
}

func (pt *pollingTrackerPut) provisioningStateApplicable() bool {
	return pt.resp.StatusCode == http.StatusOK || pt.resp.StatusCode == http.StatusCreated
}

// gets the polling URL from the Azure-AsyncOperation header.
// ensures the URL is well-formed and absolute.
func getURLFromAsyncOpHeader(resp *azcore.Response) (string, error) {
	s := resp.Header.Get(http.CanonicalHeaderKey(headerAsyncOperation))
	if s == "" {
		return "", nil
	}
	if !isValidURL(s) {
		return "", fmt.Errorf("invalid polling URL '%s'", s)
	}
	return s, nil
}

// gets the polling URL from the Location header.
// ensures the URL is well-formed and absolute.
func getURLFromLocationHeader(resp *azcore.Response) (string, error) {
	s := resp.Header.Get(http.CanonicalHeaderKey(headerLocation))
	if s == "" {
		return "", nil
	}
	if !isValidURL(s) {
		return "", fmt.Errorf("invalid polling URL '%s'", s)
	}
	return s, nil
}

// verify that the URL is valid and absolute
func isValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs()
}

// pollingMethodType defines a type used for enumerating polling mechanisms.
type pollingMethodType string

const (
	// pollingAsyncOperation indicates the polling method uses the Azure-AsyncOperation header.
	pollingAsyncOperation pollingMethodType = "AsyncOperation"

	// pollingLocation indicates the polling method uses the Location header.
	pollingLocation pollingMethodType = "Location"

	// pollingRequestURI indicates the polling method uses the original request URI.
	pollingRequestURI pollingMethodType = "RequestURI"

	// pollingUnknown indicates an unknown polling method and is the default value.
	pollingUnknown pollingMethodType = ""
)

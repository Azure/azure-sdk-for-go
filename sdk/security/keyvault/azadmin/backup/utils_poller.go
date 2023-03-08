//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package backup

// this file contains handwritten additions to the generated code
// code to support the custom poller handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/exported"
)

// the well-known set of LRO status/provisioning state values.
const (
	statusSucceeded  = "Succeeded"
	statusCanceled   = "Canceled"
	statusFailed     = "Failed"
	statusInProgress = "InProgress"
)

// isTerminalState returns true if the LRO's state is terminal.
func isTerminalState(s string) bool {
	return strings.EqualFold(s, statusSucceeded) || strings.EqualFold(s, statusFailed) || strings.EqualFold(s, statusCanceled)
}

// failed returns true if the LRO's state is terminal failure.
func failed(s string) bool {
	return strings.EqualFold(s, statusFailed) || strings.EqualFold(s, statusCanceled)
}

// returns true if the LRO response contains a valid HTTP status code
func statusCodeValid(resp *http.Response) bool {
	return exported.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusNoContent)
}

// isValidURL verifies that the URL is valid and absolute.
func isValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs()
}

// errNoBody is returned if the response didn't contain a body.
var errNoBody = errors.New("the response did not contain a body")

// getJSON reads the response body into a raw JSON object.
// It returns ErrNoBody if there was no content.
func getJSON(resp *http.Response) (map[string]any, error) {
	body, err := exported.Payload(resp)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errNoBody
	}
	// unmarshall the body to get the value
	var jsonBody map[string]any
	if err = json.Unmarshal(body, &jsonBody); err != nil {
		return nil, err
	}
	return jsonBody, nil
}

// provisioningState returns the provisioning state from the response or the empty string.
func provisioningState(jsonBody map[string]any) string {
	jsonProps, ok := jsonBody["properties"]
	if !ok {
		return ""
	}
	props, ok := jsonProps.(map[string]any)
	if !ok {
		return ""
	}
	rawPs, ok := props["provisioningState"]
	if !ok {
		return ""
	}
	ps, ok := rawPs.(string)
	if !ok {
		return ""
	}
	return ps
}

// status returns the status from the response or the empty string.
func status(jsonBody map[string]any) string {
	rawStatus, ok := jsonBody["status"]
	if !ok {
		return ""
	}
	status, ok := rawStatus.(string)
	if !ok {
		return ""
	}
	return status
}

// getStatus returns the LRO's status from the response body.
// Typically used for Azure-AsyncOperation flows.
// If there is no status in the response body the empty string is returned.
func getStatus(resp *http.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	return status(jsonBody), nil
}

// getProvisioningState returns the LRO's state from the response body.
// If there is no state in the response body the empty string is returned.
func getProvisioningState(resp *http.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	return provisioningState(jsonBody), nil
}

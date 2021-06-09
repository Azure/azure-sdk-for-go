// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	HeaderAzureAsync = "Azure-AsyncOperation"
	HeaderLocation   = "Location"
)

const (
	statusSucceeded = "succeeded"
	statusCanceled  = "canceled"
	statusFailed    = "failed"
)

// reads the response body into a raw JSON object.
// returns an empty object if there was no content.
func getJSON(resp *azcore.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if len(body) == 0 {
		return nil, ErrNoBody
	}
	// put the body back so it's available to others
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	// unmarshall the body to get the value
	var jsonBody map[string]interface{}
	if err = json.Unmarshal(body, &jsonBody); err != nil {
		return nil, err
	}
	return jsonBody, nil
}

// provisioningState returns the provisioning state from the response or the empty string.
func provisioningState(jsonBody map[string]interface{}) string {
	jsonProps, ok := jsonBody["properties"]
	if !ok {
		return ""
	}
	props, ok := jsonProps.(map[string]interface{})
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
func status(jsonBody map[string]interface{}) string {
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

// IsTerminalState returns true if the LRO's state is terminal.
func IsTerminalState(s string) bool {
	return strings.EqualFold(s, statusSucceeded) || strings.EqualFold(s, statusFailed) || strings.EqualFold(s, statusCanceled)
}

// Failed returns true if the LRO's state is terminal failure.
func Failed(s string) bool {
	return strings.EqualFold(s, statusFailed) || strings.EqualFold(s, statusCanceled)
}

// GetStatus returns the LRO's status from the response body.
// Typically used for Azure-AsyncOperation flows.
// If there is no status in the response body ErrNoStatus is returned.
func GetStatus(resp *azcore.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	if s := status(jsonBody); s != "" {
		return s, nil
	}
	return "", ErrNoStatus
}

// GetProvisioningState returns the LRO's state from the response body.
// If there is no state in the response body ErrNoProvisioningState is returned.
func GetProvisioningState(resp *azcore.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	if ps := provisioningState(jsonBody); ps != "" {
		return ps, nil
	}
	return "", ErrNoProvisioningState
}

// IsValidURL verifies that the URL is valid and absolute.
func IsValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs()
}

// MakeID returns the unique poller identifier in the format pollerID;poller.
func MakeID(pollerID string, kind string) string {
	return fmt.Sprintf("%s;%s", pollerID, kind)
}

// ErrNoBody is returned if the response didn't contain a body.
var ErrNoBody = errors.New("the response did not contain a body")

// ErrNoStatus is returned if the response body didn't contain a status.
var ErrNoStatus = errors.New("the response did not contain a status")

// ErrNoProvisioningState is returned if the response body didn't contain a provisioning state.
var ErrNoProvisioningState = errors.New("the response did not contain a provisioning state")

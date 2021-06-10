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
	StatusSucceeded  = "Succeeded"
	StatusCanceled   = "Canceled"
	StatusFailed     = "Failed"
	StatusInProgress = "InProgress"
)

// reads the response body into a raw JSON object.
// returns ErrNoBody if there was no content.
func getJSON(resp *azcore.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
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
	return strings.EqualFold(s, StatusSucceeded) || strings.EqualFold(s, StatusFailed) || strings.EqualFold(s, StatusCanceled)
}

// Failed returns true if the LRO's state is terminal failure.
func Failed(s string) bool {
	return strings.EqualFold(s, StatusFailed) || strings.EqualFold(s, StatusCanceled)
}

// GetStatus returns the LRO's status from the response body.
// Typically used for Azure-AsyncOperation flows.
// If there is no status in the response body the empty string is returned.
func GetStatus(resp *azcore.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	return status(jsonBody), nil
}

// GetProvisioningState returns the LRO's state from the response body.
// If there is no state in the response body the empty string is returned.
func GetProvisioningState(resp *azcore.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	return provisioningState(jsonBody), nil
}

// IsValidURL verifies that the URL is valid and absolute.
func IsValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs()
}

const idSeparator = ";"

// MakeID returns the poller ID from the provided values.
func MakeID(pollerID string, kind string) string {
	return fmt.Sprintf("%s%s%s", pollerID, idSeparator, kind)
}

// DecodeID decodes the poller ID, returning [pollerID, kind] or an error.
func DecodeID(tk string) (string, string, error) {
	raw := strings.Split(tk, idSeparator)
	// strings.Split will include any/all whitespace strings, we want to omit those
	parts := []string{}
	for _, r := range raw {
		if s := strings.TrimSpace(r); s != "" {
			parts = append(parts, s)
		}
	}
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid token %s", tk)
	}
	return parts[0], parts[1], nil
}

// ErrNoBody is returned if the response didn't contain a body.
var ErrNoBody = errors.New("the response did not contain a body")

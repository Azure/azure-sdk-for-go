//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

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

// GetStatus returns the LRO's status from the response body.
// Typically used for Azure-AsyncOperation flows.
// If there is no status in the response body the empty string is returned.
func GetStatus(resp *http.Response) (string, error) {
	jsonBody, err := shared.GetJSON(resp)
	if err != nil {
		return "", err
	}
	return status(jsonBody), nil
}

// GetProvisioningState returns the LRO's state from the response body.
// If there is no state in the response body the empty string is returned.
func GetProvisioningState(resp *http.Response) (string, error) {
	jsonBody, err := shared.GetJSON(resp)
	if err != nil {
		return "", err
	}
	return provisioningState(jsonBody), nil
}

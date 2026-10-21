// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"context"
	"errors"
	"net/http"
)

// GetDefaults reads model deployment mappings. DefaultsNotSet remains a service error.
func (c *Client) GetDefaults(ctx context.Context) (Defaults, error) {
	req, err := c.request(ctx, http.MethodGet, "/defaults", nil, nil, "")
	if err != nil {
		return Defaults{}, err
	}
	resp, err := c.send(req, http.StatusOK)
	if err != nil {
		return Defaults{}, err
	}
	return decodeResponse[Defaults](resp)
}

// UpdateDefaults patches model deployment mappings. A nil value removes a mapping.
// Merge or replacement policy and prebuilt aliases belong to the caller.
func (c *Client) UpdateDefaults(ctx context.Context, mappings map[string]*string) (Defaults, error) {
	if mappings == nil {
		return Defaults{}, errors.New("model deployment mappings are required")
	}
	body := struct {
		ModelDeployments map[string]*string `json:"modelDeployments"`
	}{mappings}
	req, err := c.jsonRequest(ctx, http.MethodPatch, "/defaults", nil, body, "application/merge-patch+json")
	if err != nil {
		return Defaults{}, err
	}
	resp, err := c.send(req, http.StatusOK)
	if err != nil {
		return Defaults{}, err
	}
	return decodeResponse[Defaults](resp)
}

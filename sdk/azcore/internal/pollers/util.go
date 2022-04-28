//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
)

// the well-known set of LRO status/provisioning state values.
const (
	StatusSucceeded  = "Succeeded"
	StatusCanceled   = "Canceled"
	StatusFailed     = "Failed"
	StatusInProgress = "InProgress"
)

// OperationState contains the set of non-terminal and terminal states for an LRO.
type OperationState int

const (
	OperationStateInProgress OperationState = 1
	OperationStateSucceeded  OperationState = 2
	OperationStateFailed     OperationState = 3
)

// String implements the fmt.Stringer interface for the OperationState type.
func (o OperationState) String() string {
	switch o {
	case OperationStateInProgress:
		return "InProgress"
	case OperationStateSucceeded:
		return "Succeeded"
	case OperationStateFailed:
		return "Failed"
	default:
		return fmt.Sprintf("unknown state %d", o)
	}
}

// Operation abstracts the differences among long-running operation implementations.
type Operation interface {
	// State returns the current state of the LRO.
	// Calls to Update() will update the state as required.
	State() OperationState

	// Update provides the implementation with the latest HTTP response so it can react accordingly.
	Update(resp *http.Response) error

	// FinalGetURL returns the URL to GET when the LRO has reached a terminal, success state.
	// Can return the empty string to indicate no final GET is required.  This usually indicates
	// that the final response payload (if applicable) was within the terminal success response.
	FinalGetURL() string

	// URL returns the polling URL.
	URL() string
}

// IsTerminalState returns true if the LRO's state is terminal.
func IsTerminalState(s string) bool {
	return strings.EqualFold(s, StatusSucceeded) || strings.EqualFold(s, StatusFailed) || strings.EqualFold(s, StatusCanceled)
}

// Failed returns true if the LRO's state is terminal failure.
func Failed(s string) bool {
	return strings.EqualFold(s, StatusFailed) || strings.EqualFold(s, StatusCanceled)
}

// returns true if the LRO response contains a valid HTTP status code
func StatusCodeValid(resp *http.Response) bool {
	return exported.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusNoContent)
}

// IsValidURL verifies that the URL is valid and absolute.
func IsValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs()
}

const idSeparator = ";"

// PollerTypeName returns the type name to use when constructing the poller ID.
// An error is returned if the generic type has no name (e.g. struct{}).
func PollerTypeName[T any]() (string, error) {
	tt := shared.TypeOfT[T]()
	var n string
	if tt.Kind() == reflect.Pointer {
		n = "*"
		tt = tt.Elem()
	}
	n += tt.Name()
	if n == "" {
		return "", errors.New("nameless types are not allowed")
	}
	return n, nil
}

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

// GetJSON reads the response body into a raw JSON object.
// It returns ErrNoBody if there was no content.
func GetJSON(resp *http.Response) (map[string]interface{}, error) {
	body, err := exported.Payload(resp)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, ErrNoBody
	}
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

// GetStatus returns the LRO's status from the response body.
// Typically used for Azure-AsyncOperation flows.
// If there is no status in the response body the empty string is returned.
func GetStatus(resp *http.Response) (string, error) {
	jsonBody, err := GetJSON(resp)
	if err != nil {
		return "", err
	}
	return status(jsonBody), nil
}

// GetProvisioningState returns the LRO's state from the response body.
// If there is no state in the response body the empty string is returned.
func GetProvisioningState(resp *http.Response) (string, error) {
	jsonBody, err := GetJSON(resp)
	if err != nil {
		return "", err
	}
	return provisioningState(jsonBody), nil
}

// used if the operation synchronously completed
type NopPoller struct{}

func (*NopPoller) URL() string {
	return ""
}

func (*NopPoller) State() OperationState {
	return OperationStateSucceeded
}

func (*NopPoller) Update(*http.Response) error {
	return nil
}

func (*NopPoller) FinalGetURL() string {
	return ""
}

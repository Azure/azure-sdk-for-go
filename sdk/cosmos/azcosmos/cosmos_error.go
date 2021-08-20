// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// cosmosError is used as base error for any error response from the Cosmos service.
type cosmosError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func newCosmosError(response *azcore.Response) error {
	var cError cosmosError
	bytesRead, err := response.Payload()
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesRead, &cError)
	if err != nil {
		return errors.New(string(bytesRead))
	}

	return &cError
}

func (e *cosmosError) Error() string {
	if e.Code == "" && e.Message == "" {
		return ""
	}
	return fmt.Sprintf("Code: %v, Message %v", e.Code, e.Message)
}

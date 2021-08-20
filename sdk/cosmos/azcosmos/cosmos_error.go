// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// cosmosError is used as base error for any error response from the Cosmos service.
type cosmosError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func newCosmosError(response *azcore.Response) error {
	defer response.Body.Close()

	var cError cosmosError
	err := json.NewDecoder(response.Body).Decode(&cError)
	switch {
	case err == io.EOF:
		return errors.New("request failed")
	case err != nil:
		return err
	}

	return &cError
}

func (e *cosmosError) Error() string {
	return fmt.Sprintf("Code: %v, Message %v", e.Code, e.Message)
}

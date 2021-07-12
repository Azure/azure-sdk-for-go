// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// A CosmosClient is used to interact with the Azure Cosmos DB database service.
type CosmosClient struct{}

func NewCosmosClient(endpoint string, cred azcore.Credential, options *CosmosClientOptions) (*CosmosClient, error) {
	return &CosmosClient{}, errors.New("not implemented")
}

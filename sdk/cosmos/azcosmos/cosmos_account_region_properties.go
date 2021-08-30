// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// The AccountRegion struct represents an Azure Cosmos DB database account in a specific region.
type AccountRegion struct {
	// Gets the name of the database account in the Azure Cosmos DB service.
	Name string `json:"Constants.Properties.Name"`
	// Gets the URL of the database account in the Azure Cosmos DB service.
	Endpoint string `json:"Constants.Properties.DatabaseAccountEndpoint"`
}

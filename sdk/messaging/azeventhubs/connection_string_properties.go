// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"

// ConnectionStringProperties are the properties of a connection string
// as returned by [NewConnectionStringProperties].
type ConnectionStringProperties = exported.ConnectionStringProperties

// NewConnectionStringProperties takes a connection string from the Azure portal and returns the
// parsed representation The method will return an error if the Endpoint, SharedAccessKeyName
// or SharedAccessKey is empty.
func NewConnectionStringProperties(connStr string) (ConnectionStringProperties, error) {
	return exported.NewConnectionStringProperties(connStr)
}

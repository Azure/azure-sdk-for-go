//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azwebpubsub

// AddToGroupsRequest - The request object containing targets groups and a connection filter
type AddToGroupsRequest struct {
	// An OData filter which target connections satisfy
	Filter *string

	// A list of groups which target connections will be added into
	Groups []string
}

// TokenResponse - The response object containing the token for the client
type TokenResponse struct {
	// The token value for the WebSocket client to connect to the service
	Token *string
}

// RemoveFromGroupsRequest - The request object containing targets groups and a connection filter
type RemoveFromGroupsRequest struct {
	// An OData filter which target connections satisfy
	Filter *string

	// A list of groups which target connections will be removed from
	Groups []string
}

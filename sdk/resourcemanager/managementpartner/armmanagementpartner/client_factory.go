//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagementpartner

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// ClientFactory is a client factory used to create any client in this module.
// Don't use this type directly, use NewClientFactory instead.
type ClientFactory struct {
	credential azcore.TokenCredential
	options    *arm.ClientOptions
}

// NewClientFactory creates a new instance of ClientFactory with the specified values.
// The parameter values will be propagated to any client created from this factory.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewClientFactory(credential azcore.TokenCredential, options *arm.ClientOptions) (*ClientFactory, error) {
	_, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	return &ClientFactory{
		credential: credential,
		options:    options.Clone(),
	}, nil
}

// NewOperationClient creates a new instance of OperationClient.
func (c *ClientFactory) NewOperationClient() *OperationClient {
	subClient, _ := NewOperationClient(c.credential, c.options)
	return subClient
}

// NewPartnerClient creates a new instance of PartnerClient.
func (c *ClientFactory) NewPartnerClient() *PartnerClient {
	subClient, _ := NewPartnerClient(c.credential, c.options)
	return subClient
}

// NewPartnersClient creates a new instance of PartnersClient.
func (c *ClientFactory) NewPartnersClient() *PartnersClient {
	subClient, _ := NewPartnersClient(c.credential, c.options)
	return subClient
}

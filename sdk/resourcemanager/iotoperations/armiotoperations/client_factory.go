// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armiotoperations

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// ClientFactory is a client factory used to create any client in this module.
// Don't use this type directly, use NewClientFactory instead.
type ClientFactory struct {
	subscriptionID string
	internal       *arm.Client
}

// NewClientFactory creates a new instance of ClientFactory with the specified values.
// The parameter values will be propagated to any client created from this factory.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewClientFactory(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClientFactory, error) {
	internal, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	return &ClientFactory{
		subscriptionID: subscriptionID,
		internal:       internal,
	}, nil
}

// NewBrokerAuthenticationClient creates a new instance of BrokerAuthenticationClient.
func (c *ClientFactory) NewBrokerAuthenticationClient() *BrokerAuthenticationClient {
	return &BrokerAuthenticationClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBrokerAuthorizationClient creates a new instance of BrokerAuthorizationClient.
func (c *ClientFactory) NewBrokerAuthorizationClient() *BrokerAuthorizationClient {
	return &BrokerAuthorizationClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBrokerClient creates a new instance of BrokerClient.
func (c *ClientFactory) NewBrokerClient() *BrokerClient {
	return &BrokerClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBrokerListenerClient creates a new instance of BrokerListenerClient.
func (c *ClientFactory) NewBrokerListenerClient() *BrokerListenerClient {
	return &BrokerListenerClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewDataflowClient creates a new instance of DataflowClient.
func (c *ClientFactory) NewDataflowClient() *DataflowClient {
	return &DataflowClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewDataflowEndpointClient creates a new instance of DataflowEndpointClient.
func (c *ClientFactory) NewDataflowEndpointClient() *DataflowEndpointClient {
	return &DataflowEndpointClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewDataflowProfileClient creates a new instance of DataflowProfileClient.
func (c *ClientFactory) NewDataflowProfileClient() *DataflowProfileClient {
	return &DataflowProfileClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewInstanceClient creates a new instance of InstanceClient.
func (c *ClientFactory) NewInstanceClient() *InstanceClient {
	return &InstanceClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewOperationsClient creates a new instance of OperationsClient.
func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	return &OperationsClient{
		internal: c.internal,
	}
}

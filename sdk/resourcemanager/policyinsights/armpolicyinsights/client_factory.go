// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpolicyinsights

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
//   - subscriptionID - Microsoft Azure subscription ID.
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

// NewAttestationsClient creates a new instance of AttestationsClient.
func (c *ClientFactory) NewAttestationsClient() *AttestationsClient {
	return &AttestationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewComponentPolicyStatesClient creates a new instance of ComponentPolicyStatesClient.
func (c *ClientFactory) NewComponentPolicyStatesClient() *ComponentPolicyStatesClient {
	return &ComponentPolicyStatesClient{
		internal: c.internal,
	}
}

// NewOperationsClient creates a new instance of OperationsClient.
func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	return &OperationsClient{
		internal: c.internal,
	}
}

// NewPolicyEventsClient creates a new instance of PolicyEventsClient.
func (c *ClientFactory) NewPolicyEventsClient() *PolicyEventsClient {
	return &PolicyEventsClient{
		internal: c.internal,
	}
}

// NewPolicyMetadataClient creates a new instance of PolicyMetadataClient.
func (c *ClientFactory) NewPolicyMetadataClient() *PolicyMetadataClient {
	return &PolicyMetadataClient{
		internal: c.internal,
	}
}

// NewPolicyRestrictionsClient creates a new instance of PolicyRestrictionsClient.
func (c *ClientFactory) NewPolicyRestrictionsClient() *PolicyRestrictionsClient {
	return &PolicyRestrictionsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewPolicyStatesClient creates a new instance of PolicyStatesClient.
func (c *ClientFactory) NewPolicyStatesClient() *PolicyStatesClient {
	return &PolicyStatesClient{
		internal: c.internal,
	}
}

// NewPolicyTrackedResourcesClient creates a new instance of PolicyTrackedResourcesClient.
func (c *ClientFactory) NewPolicyTrackedResourcesClient() *PolicyTrackedResourcesClient {
	return &PolicyTrackedResourcesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewRemediationsClient creates a new instance of RemediationsClient.
func (c *ClientFactory) NewRemediationsClient() *RemediationsClient {
	return &RemediationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

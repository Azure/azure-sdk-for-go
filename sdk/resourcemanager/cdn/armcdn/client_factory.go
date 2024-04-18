//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcdn

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
//   - subscriptionID - Azure Subscription ID.
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

// NewAFDCustomDomainsClient creates a new instance of AFDCustomDomainsClient.
func (c *ClientFactory) NewAFDCustomDomainsClient() *AFDCustomDomainsClient {
	return &AFDCustomDomainsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAFDEndpointsClient creates a new instance of AFDEndpointsClient.
func (c *ClientFactory) NewAFDEndpointsClient() *AFDEndpointsClient {
	return &AFDEndpointsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAFDOriginGroupsClient creates a new instance of AFDOriginGroupsClient.
func (c *ClientFactory) NewAFDOriginGroupsClient() *AFDOriginGroupsClient {
	return &AFDOriginGroupsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAFDOriginsClient creates a new instance of AFDOriginsClient.
func (c *ClientFactory) NewAFDOriginsClient() *AFDOriginsClient {
	return &AFDOriginsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAFDProfilesClient creates a new instance of AFDProfilesClient.
func (c *ClientFactory) NewAFDProfilesClient() *AFDProfilesClient {
	return &AFDProfilesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewCustomDomainsClient creates a new instance of CustomDomainsClient.
func (c *ClientFactory) NewCustomDomainsClient() *CustomDomainsClient {
	return &CustomDomainsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewEdgeNodesClient creates a new instance of EdgeNodesClient.
func (c *ClientFactory) NewEdgeNodesClient() *EdgeNodesClient {
	return &EdgeNodesClient{
		internal: c.internal,
	}
}

// NewEndpointsClient creates a new instance of EndpointsClient.
func (c *ClientFactory) NewEndpointsClient() *EndpointsClient {
	return &EndpointsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewLogAnalyticsClient creates a new instance of LogAnalyticsClient.
func (c *ClientFactory) NewLogAnalyticsClient() *LogAnalyticsClient {
	return &LogAnalyticsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewManagedRuleSetsClient creates a new instance of ManagedRuleSetsClient.
func (c *ClientFactory) NewManagedRuleSetsClient() *ManagedRuleSetsClient {
	return &ManagedRuleSetsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewManagementClient creates a new instance of ManagementClient.
func (c *ClientFactory) NewManagementClient() *ManagementClient {
	return &ManagementClient{
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

// NewOriginGroupsClient creates a new instance of OriginGroupsClient.
func (c *ClientFactory) NewOriginGroupsClient() *OriginGroupsClient {
	return &OriginGroupsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewOriginsClient creates a new instance of OriginsClient.
func (c *ClientFactory) NewOriginsClient() *OriginsClient {
	return &OriginsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewPoliciesClient creates a new instance of PoliciesClient.
func (c *ClientFactory) NewPoliciesClient() *PoliciesClient {
	return &PoliciesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewProfilesClient creates a new instance of ProfilesClient.
func (c *ClientFactory) NewProfilesClient() *ProfilesClient {
	return &ProfilesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewResourceUsageClient creates a new instance of ResourceUsageClient.
func (c *ClientFactory) NewResourceUsageClient() *ResourceUsageClient {
	return &ResourceUsageClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewRoutesClient creates a new instance of RoutesClient.
func (c *ClientFactory) NewRoutesClient() *RoutesClient {
	return &RoutesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewRuleSetsClient creates a new instance of RuleSetsClient.
func (c *ClientFactory) NewRuleSetsClient() *RuleSetsClient {
	return &RuleSetsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewRulesClient creates a new instance of RulesClient.
func (c *ClientFactory) NewRulesClient() *RulesClient {
	return &RulesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewSecretsClient creates a new instance of SecretsClient.
func (c *ClientFactory) NewSecretsClient() *SecretsClient {
	return &SecretsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewSecurityPoliciesClient creates a new instance of SecurityPoliciesClient.
func (c *ClientFactory) NewSecurityPoliciesClient() *SecurityPoliciesClient {
	return &SecurityPoliciesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

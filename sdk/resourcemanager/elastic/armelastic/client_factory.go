//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by Microsoft (R) AutoRest Code Generator.Changes may cause incorrect behavior and will be lost if the code
// is regenerated.
// Code generated by @autorest/go. DO NOT EDIT.

package armelastic

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// ClientFactory is a client factory used to create any client in this module.
// Don't use this type directly, use NewClientFactory instead.
type ClientFactory struct {
	subscriptionID string
	credential azcore.TokenCredential
	options *arm.ClientOptions
}

// NewClientFactory creates a new instance of ClientFactory with the specified values.
// The parameter values will be propagated to any client created from this factory.
//   - subscriptionID - The Azure subscription ID. This is a GUID-formatted string (e.g. 00000000-0000-0000-0000-000000000000)
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewClientFactory(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClientFactory, error) {
	_, err := arm.NewClient(moduleName+".ClientFactory", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	return &ClientFactory{
		subscriptionID: 	subscriptionID,		credential: credential,
		options: options.Clone(),
	}, nil
}

func (c *ClientFactory) NewAllTrafficFiltersClient() *AllTrafficFiltersClient {
	subClient, _ := NewAllTrafficFiltersClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewAssociateTrafficFilterClient() *AssociateTrafficFilterClient {
	subClient, _ := NewAssociateTrafficFilterClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewCreateAndAssociateIPFilterClient() *CreateAndAssociateIPFilterClient {
	subClient, _ := NewCreateAndAssociateIPFilterClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewCreateAndAssociatePLFilterClient() *CreateAndAssociatePLFilterClient {
	subClient, _ := NewCreateAndAssociatePLFilterClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewDeploymentInfoClient() *DeploymentInfoClient {
	subClient, _ := NewDeploymentInfoClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewDetachAndDeleteTrafficFilterClient() *DetachAndDeleteTrafficFilterClient {
	subClient, _ := NewDetachAndDeleteTrafficFilterClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewDetachTrafficFilterClient() *DetachTrafficFilterClient {
	subClient, _ := NewDetachTrafficFilterClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewExternalUserClient() *ExternalUserClient {
	subClient, _ := NewExternalUserClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewListAssociatedTrafficFiltersClient() *ListAssociatedTrafficFiltersClient {
	subClient, _ := NewListAssociatedTrafficFiltersClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewMonitorClient() *MonitorClient {
	subClient, _ := NewMonitorClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewMonitoredResourcesClient() *MonitoredResourcesClient {
	subClient, _ := NewMonitoredResourcesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewMonitorsClient() *MonitorsClient {
	subClient, _ := NewMonitorsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	subClient, _ := NewOperationsClient(c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewOrganizationsClient() *OrganizationsClient {
	subClient, _ := NewOrganizationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewTagRulesClient() *TagRulesClient {
	subClient, _ := NewTagRulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewTrafficFiltersClient() *TrafficFiltersClient {
	subClient, _ := NewTrafficFiltersClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewUpgradableVersionsClient() *UpgradableVersionsClient {
	subClient, _ := NewUpgradableVersionsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVMCollectionClient() *VMCollectionClient {
	subClient, _ := NewVMCollectionClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVMHostClient() *VMHostClient {
	subClient, _ := NewVMHostClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVMIngestionClient() *VMIngestionClient {
	subClient, _ := NewVMIngestionClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVersionsClient() *VersionsClient {
	subClient, _ := NewVersionsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}


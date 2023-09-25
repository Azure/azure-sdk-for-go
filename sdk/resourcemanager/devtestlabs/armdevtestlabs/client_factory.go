//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdevtestlabs

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
//   - subscriptionID - The subscription ID.
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

func (c *ClientFactory) NewArmTemplatesClient() *ArmTemplatesClient {
	subClient, _ := NewArmTemplatesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewArtifactSourcesClient() *ArtifactSourcesClient {
	subClient, _ := NewArtifactSourcesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewArtifactsClient() *ArtifactsClient {
	subClient, _ := NewArtifactsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewCostsClient() *CostsClient {
	subClient, _ := NewCostsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewCustomImagesClient() *CustomImagesClient {
	subClient, _ := NewCustomImagesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewDisksClient() *DisksClient {
	subClient, _ := NewDisksClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewEnvironmentsClient() *EnvironmentsClient {
	subClient, _ := NewEnvironmentsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewFormulasClient() *FormulasClient {
	subClient, _ := NewFormulasClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewGalleryImagesClient() *GalleryImagesClient {
	subClient, _ := NewGalleryImagesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewGlobalSchedulesClient() *GlobalSchedulesClient {
	subClient, _ := NewGlobalSchedulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewLabsClient() *LabsClient {
	subClient, _ := NewLabsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewNotificationChannelsClient() *NotificationChannelsClient {
	subClient, _ := NewNotificationChannelsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	subClient, _ := NewOperationsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPoliciesClient() *PoliciesClient {
	subClient, _ := NewPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPolicySetsClient() *PolicySetsClient {
	subClient, _ := NewPolicySetsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewProviderOperationsClient() *ProviderOperationsClient {
	subClient, _ := NewProviderOperationsClient(c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSchedulesClient() *SchedulesClient {
	subClient, _ := NewSchedulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSecretsClient() *SecretsClient {
	subClient, _ := NewSecretsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewServiceFabricSchedulesClient() *ServiceFabricSchedulesClient {
	subClient, _ := NewServiceFabricSchedulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewServiceFabricsClient() *ServiceFabricsClient {
	subClient, _ := NewServiceFabricsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewServiceRunnersClient() *ServiceRunnersClient {
	subClient, _ := NewServiceRunnersClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewUsersClient() *UsersClient {
	subClient, _ := NewUsersClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVirtualMachineSchedulesClient() *VirtualMachineSchedulesClient {
	subClient, _ := NewVirtualMachineSchedulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVirtualMachinesClient() *VirtualMachinesClient {
	subClient, _ := NewVirtualMachinesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVirtualNetworksClient() *VirtualNetworksClient {
	subClient, _ := NewVirtualNetworksClient(c.subscriptionID, c.credential, c.options)
	return subClient
}


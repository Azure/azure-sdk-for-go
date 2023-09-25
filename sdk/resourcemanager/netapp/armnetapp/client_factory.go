//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetapp

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
//   - subscriptionID - The ID of the target subscription.
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

func (c *ClientFactory) NewAccountBackupsClient() *AccountBackupsClient {
	subClient, _ := NewAccountBackupsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewAccountsClient() *AccountsClient {
	subClient, _ := NewAccountsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewBackupPoliciesClient() *BackupPoliciesClient {
	subClient, _ := NewBackupPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewBackupsClient() *BackupsClient {
	subClient, _ := NewBackupsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	subClient, _ := NewOperationsClient(c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewPoolsClient() *PoolsClient {
	subClient, _ := NewPoolsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewResourceClient() *ResourceClient {
	subClient, _ := NewResourceClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewResourceQuotaLimitsClient() *ResourceQuotaLimitsClient {
	subClient, _ := NewResourceQuotaLimitsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSnapshotPoliciesClient() *SnapshotPoliciesClient {
	subClient, _ := NewSnapshotPoliciesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSnapshotsClient() *SnapshotsClient {
	subClient, _ := NewSnapshotsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewSubvolumesClient() *SubvolumesClient {
	subClient, _ := NewSubvolumesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVolumeGroupsClient() *VolumeGroupsClient {
	subClient, _ := NewVolumeGroupsClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVolumeQuotaRulesClient() *VolumeQuotaRulesClient {
	subClient, _ := NewVolumeQuotaRulesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}

func (c *ClientFactory) NewVolumesClient() *VolumesClient {
	subClient, _ := NewVolumesClient(c.subscriptionID, c.credential, c.options)
	return subClient
}


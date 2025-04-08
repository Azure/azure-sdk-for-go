// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpostgresqlflexibleservers

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

// NewAdministratorsClient creates a new instance of AdministratorsClient.
func (c *ClientFactory) NewAdministratorsClient() *AdministratorsClient {
	return &AdministratorsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBackupsClient creates a new instance of BackupsClient.
func (c *ClientFactory) NewBackupsClient() *BackupsClient {
	return &BackupsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewCheckNameAvailabilityClient creates a new instance of CheckNameAvailabilityClient.
func (c *ClientFactory) NewCheckNameAvailabilityClient() *CheckNameAvailabilityClient {
	return &CheckNameAvailabilityClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewCheckNameAvailabilityWithLocationClient creates a new instance of CheckNameAvailabilityWithLocationClient.
func (c *ClientFactory) NewCheckNameAvailabilityWithLocationClient() *CheckNameAvailabilityWithLocationClient {
	return &CheckNameAvailabilityWithLocationClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewConfigurationsClient creates a new instance of ConfigurationsClient.
func (c *ClientFactory) NewConfigurationsClient() *ConfigurationsClient {
	return &ConfigurationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewDatabasesClient creates a new instance of DatabasesClient.
func (c *ClientFactory) NewDatabasesClient() *DatabasesClient {
	return &DatabasesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewFirewallRulesClient creates a new instance of FirewallRulesClient.
func (c *ClientFactory) NewFirewallRulesClient() *FirewallRulesClient {
	return &FirewallRulesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewFlexibleServerClient creates a new instance of FlexibleServerClient.
func (c *ClientFactory) NewFlexibleServerClient() *FlexibleServerClient {
	return &FlexibleServerClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewGetPrivateDNSZoneSuffixClient creates a new instance of GetPrivateDNSZoneSuffixClient.
func (c *ClientFactory) NewGetPrivateDNSZoneSuffixClient() *GetPrivateDNSZoneSuffixClient {
	return &GetPrivateDNSZoneSuffixClient{
		internal: c.internal,
	}
}

// NewLocationBasedCapabilitiesClient creates a new instance of LocationBasedCapabilitiesClient.
func (c *ClientFactory) NewLocationBasedCapabilitiesClient() *LocationBasedCapabilitiesClient {
	return &LocationBasedCapabilitiesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewLogFilesClient creates a new instance of LogFilesClient.
func (c *ClientFactory) NewLogFilesClient() *LogFilesClient {
	return &LogFilesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewLtrBackupOperationsClient creates a new instance of LtrBackupOperationsClient.
func (c *ClientFactory) NewLtrBackupOperationsClient() *LtrBackupOperationsClient {
	return &LtrBackupOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewMigrationsClient creates a new instance of MigrationsClient.
func (c *ClientFactory) NewMigrationsClient() *MigrationsClient {
	return &MigrationsClient{
		internal: c.internal,
	}
}

// NewOperationsClient creates a new instance of OperationsClient.
func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	return &OperationsClient{
		internal: c.internal,
	}
}

// NewPostgreSQLServerManagementClient creates a new instance of PostgreSQLServerManagementClient.
func (c *ClientFactory) NewPostgreSQLServerManagementClient() *PostgreSQLServerManagementClient {
	return &PostgreSQLServerManagementClient{
		internal: c.internal,
	}
}

// NewPrivateEndpointConnectionClient creates a new instance of PrivateEndpointConnectionClient.
func (c *ClientFactory) NewPrivateEndpointConnectionClient() *PrivateEndpointConnectionClient {
	return &PrivateEndpointConnectionClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewPrivateEndpointConnectionsClient creates a new instance of PrivateEndpointConnectionsClient.
func (c *ClientFactory) NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient {
	return &PrivateEndpointConnectionsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewPrivateLinkResourcesClient creates a new instance of PrivateLinkResourcesClient.
func (c *ClientFactory) NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient {
	return &PrivateLinkResourcesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewReplicasClient creates a new instance of ReplicasClient.
func (c *ClientFactory) NewReplicasClient() *ReplicasClient {
	return &ReplicasClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewServerCapabilitiesClient creates a new instance of ServerCapabilitiesClient.
func (c *ClientFactory) NewServerCapabilitiesClient() *ServerCapabilitiesClient {
	return &ServerCapabilitiesClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewServerThreatProtectionSettingsClient creates a new instance of ServerThreatProtectionSettingsClient.
func (c *ClientFactory) NewServerThreatProtectionSettingsClient() *ServerThreatProtectionSettingsClient {
	return &ServerThreatProtectionSettingsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewServersClient creates a new instance of ServersClient.
func (c *ClientFactory) NewServersClient() *ServersClient {
	return &ServersClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewVirtualEndpointsClient creates a new instance of VirtualEndpointsClient.
func (c *ClientFactory) NewVirtualEndpointsClient() *VirtualEndpointsClient {
	return &VirtualEndpointsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewVirtualNetworkSubnetUsageClient creates a new instance of VirtualNetworkSubnetUsageClient.
func (c *ClientFactory) NewVirtualNetworkSubnetUsageClient() *VirtualNetworkSubnetUsageClient {
	return &VirtualNetworkSubnetUsageClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

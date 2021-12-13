# Unreleased

## Breaking Changes

### Removed Funcs

1. MoveResourcePropertiesMoveStatus.MarshalJSON() ([]byte, error)
1. MoveResourceStatus.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. SummaryItem

#### Removed Struct Fields

1. MoveResourceCollection.Summary
1. MoveResourcePropertiesMoveStatus.TargetID
1. MoveResourceStatus.TargetID
1. PublicIPAddressResourceSettings.FQDN

### Signature Changes

#### Funcs

1. UnresolvedDependenciesClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, DependencyLevel, string, string
	- Returns
		- From: UnresolvedDependencyCollection, error
		- To: UnresolvedDependencyCollectionPage, error
1. UnresolvedDependenciesClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, DependencyLevel, string, string

## Additive Changes

### New Constants

1. DependencyLevel.Descendant
1. DependencyLevel.Direct
1. MoveState.DeleteSourcePending
1. MoveState.ResourceMoveCompleted
1. ResourceType.ResourceTypeMicrosoftComputediskEncryptionSets
1. ResourceType.ResourceTypeMicrosoftKeyVaultvaults

### New Funcs

1. *UnresolvedDependencyCollectionIterator.Next() error
1. *UnresolvedDependencyCollectionIterator.NextWithContext(context.Context) error
1. *UnresolvedDependencyCollectionPage.Next() error
1. *UnresolvedDependencyCollectionPage.NextWithContext(context.Context) error
1. AvailabilitySetResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. AvailabilitySetResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsAvailabilitySetResourceSettings() (*AvailabilitySetResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsBasicResourceSettings() (BasicResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsLoadBalancerResourceSettings() (*LoadBalancerResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsNetworkInterfaceResourceSettings() (*NetworkInterfaceResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsNetworkSecurityGroupResourceSettings() (*NetworkSecurityGroupResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsPublicIPAddressResourceSettings() (*PublicIPAddressResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsResourceGroupResourceSettings() (*ResourceGroupResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsResourceSettings() (*ResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsSQLDatabaseResourceSettings() (*SQLDatabaseResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsSQLElasticPoolResourceSettings() (*SQLElasticPoolResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsSQLServerResourceSettings() (*SQLServerResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsVirtualMachineResourceSettings() (*VirtualMachineResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.AsVirtualNetworkResourceSettings() (*VirtualNetworkResourceSettings, bool)
1. DiskEncryptionSetResourceSettings.MarshalJSON() ([]byte, error)
1. KeyVaultResourceSettings.AsAvailabilitySetResourceSettings() (*AvailabilitySetResourceSettings, bool)
1. KeyVaultResourceSettings.AsBasicResourceSettings() (BasicResourceSettings, bool)
1. KeyVaultResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. KeyVaultResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. KeyVaultResourceSettings.AsLoadBalancerResourceSettings() (*LoadBalancerResourceSettings, bool)
1. KeyVaultResourceSettings.AsNetworkInterfaceResourceSettings() (*NetworkInterfaceResourceSettings, bool)
1. KeyVaultResourceSettings.AsNetworkSecurityGroupResourceSettings() (*NetworkSecurityGroupResourceSettings, bool)
1. KeyVaultResourceSettings.AsPublicIPAddressResourceSettings() (*PublicIPAddressResourceSettings, bool)
1. KeyVaultResourceSettings.AsResourceGroupResourceSettings() (*ResourceGroupResourceSettings, bool)
1. KeyVaultResourceSettings.AsResourceSettings() (*ResourceSettings, bool)
1. KeyVaultResourceSettings.AsSQLDatabaseResourceSettings() (*SQLDatabaseResourceSettings, bool)
1. KeyVaultResourceSettings.AsSQLElasticPoolResourceSettings() (*SQLElasticPoolResourceSettings, bool)
1. KeyVaultResourceSettings.AsSQLServerResourceSettings() (*SQLServerResourceSettings, bool)
1. KeyVaultResourceSettings.AsVirtualMachineResourceSettings() (*VirtualMachineResourceSettings, bool)
1. KeyVaultResourceSettings.AsVirtualNetworkResourceSettings() (*VirtualNetworkResourceSettings, bool)
1. KeyVaultResourceSettings.MarshalJSON() ([]byte, error)
1. LoadBalancerResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. LoadBalancerResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. MoveCollectionProperties.MarshalJSON() ([]byte, error)
1. MoveCollectionsClient.ListRequiredFor(context.Context, string, string, string) (RequiredForResourcesCollection, error)
1. MoveCollectionsClient.ListRequiredForPreparer(context.Context, string, string, string) (*http.Request, error)
1. MoveCollectionsClient.ListRequiredForResponder(*http.Response) (RequiredForResourcesCollection, error)
1. MoveCollectionsClient.ListRequiredForSender(*http.Request) (*http.Response, error)
1. MoveResourceCollection.MarshalJSON() ([]byte, error)
1. NetworkInterfaceResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. NetworkInterfaceResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. NetworkSecurityGroupResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. NetworkSecurityGroupResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. NewUnresolvedDependencyCollectionIterator(UnresolvedDependencyCollectionPage) UnresolvedDependencyCollectionIterator
1. NewUnresolvedDependencyCollectionPage(UnresolvedDependencyCollection, func(context.Context, UnresolvedDependencyCollection) (UnresolvedDependencyCollection, error)) UnresolvedDependencyCollectionPage
1. PossibleDependencyLevelValues() []DependencyLevel
1. PublicIPAddressResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. PublicIPAddressResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. ResourceGroupResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. ResourceGroupResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. ResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. ResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. SQLDatabaseResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. SQLDatabaseResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. SQLElasticPoolResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. SQLElasticPoolResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. SQLServerResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. SQLServerResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. UnresolvedDependenciesClient.GetComplete(context.Context, string, string, DependencyLevel, string, string) (UnresolvedDependencyCollectionIterator, error)
1. UnresolvedDependencyCollection.IsEmpty() bool
1. UnresolvedDependencyCollection.MarshalJSON() ([]byte, error)
1. UnresolvedDependencyCollectionIterator.NotDone() bool
1. UnresolvedDependencyCollectionIterator.Response() UnresolvedDependencyCollection
1. UnresolvedDependencyCollectionIterator.Value() UnresolvedDependency
1. UnresolvedDependencyCollectionPage.NotDone() bool
1. UnresolvedDependencyCollectionPage.Response() UnresolvedDependencyCollection
1. UnresolvedDependencyCollectionPage.Values() []UnresolvedDependency
1. VirtualMachineResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. VirtualMachineResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)
1. VirtualNetworkResourceSettings.AsDiskEncryptionSetResourceSettings() (*DiskEncryptionSetResourceSettings, bool)
1. VirtualNetworkResourceSettings.AsKeyVaultResourceSettings() (*KeyVaultResourceSettings, bool)

### Struct Changes

#### New Structs

1. DiskEncryptionSetResourceSettings
1. KeyVaultResourceSettings
1. MoveCollectionPropertiesErrors
1. NsgReference
1. PublicIPReference
1. RequiredForResourcesCollection
1. Summary
1. SummaryCollection
1. UnresolvedDependenciesFilter
1. UnresolvedDependenciesFilterProperties
1. UnresolvedDependencyCollectionIterator
1. UnresolvedDependencyCollectionPage

#### New Struct Fields

1. MoveCollection.Etag
1. MoveCollectionProperties.Errors
1. MoveResourceCollection.SummaryCollection
1. MoveResourceCollection.TotalCount
1. MoveResourceProperties.IsResolveRequired
1. NicIPConfigurationResourceSettings.LoadBalancerNatRules
1. NicIPConfigurationResourceSettings.PublicIP
1. PublicIPAddressResourceSettings.Fqdn
1. SubnetResourceSettings.NetworkSecurityGroup
1. UnresolvedDependencyCollection.SummaryCollection
1. UnresolvedDependencyCollection.TotalCount

package cdn

import (
	 original "github.com/Azure/azure-sdk-for-go/service/cdn/management/2016-10-02/cdn"
)

type (
	 ManagementClient = original.ManagementClient
	 CustomDomainsClient = original.CustomDomainsClient
	 EdgeNodesClient = original.EdgeNodesClient
	 EndpointsClient = original.EndpointsClient
	 CustomDomainResourceState = original.CustomDomainResourceState
	 CustomHTTPSProvisioningState = original.CustomHTTPSProvisioningState
	 EndpointResourceState = original.EndpointResourceState
	 GeoFilterActions = original.GeoFilterActions
	 OriginResourceState = original.OriginResourceState
	 ProfileResourceState = original.ProfileResourceState
	 QueryStringCachingBehavior = original.QueryStringCachingBehavior
	 ResourceType = original.ResourceType
	 SkuName = original.SkuName
	 CheckNameAvailabilityInput = original.CheckNameAvailabilityInput
	 CheckNameAvailabilityOutput = original.CheckNameAvailabilityOutput
	 CidrIPAddress = original.CidrIPAddress
	 CustomDomain = original.CustomDomain
	 CustomDomainListResult = original.CustomDomainListResult
	 CustomDomainParameters = original.CustomDomainParameters
	 CustomDomainProperties = original.CustomDomainProperties
	 CustomDomainPropertiesParameters = original.CustomDomainPropertiesParameters
	 DeepCreatedOrigin = original.DeepCreatedOrigin
	 DeepCreatedOriginProperties = original.DeepCreatedOriginProperties
	 EdgeNode = original.EdgeNode
	 EdgeNodeProperties = original.EdgeNodeProperties
	 EdgenodeResult = original.EdgenodeResult
	 Endpoint = original.Endpoint
	 EndpointListResult = original.EndpointListResult
	 EndpointProperties = original.EndpointProperties
	 EndpointPropertiesUpdateParameters = original.EndpointPropertiesUpdateParameters
	 EndpointUpdateParameters = original.EndpointUpdateParameters
	 ErrorResponse = original.ErrorResponse
	 GeoFilter = original.GeoFilter
	 IPAddressGroup = original.IPAddressGroup
	 LoadParameters = original.LoadParameters
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 Origin = original.Origin
	 OriginListResult = original.OriginListResult
	 OriginProperties = original.OriginProperties
	 OriginPropertiesParameters = original.OriginPropertiesParameters
	 OriginUpdateParameters = original.OriginUpdateParameters
	 Profile = original.Profile
	 ProfileListResult = original.ProfileListResult
	 ProfileProperties = original.ProfileProperties
	 ProfileUpdateParameters = original.ProfileUpdateParameters
	 PurgeParameters = original.PurgeParameters
	 Resource = original.Resource
	 ResourceUsage = original.ResourceUsage
	 ResourceUsageListResult = original.ResourceUsageListResult
	 Sku = original.Sku
	 SsoURI = original.SsoURI
	 ValidateCustomDomainInput = original.ValidateCustomDomainInput
	 ValidateCustomDomainOutput = original.ValidateCustomDomainOutput
	 OriginsClient = original.OriginsClient
	 ProfilesClient = original.ProfilesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Active = original.Active
	 Creating = original.Creating
	 Deleting = original.Deleting
	 Disabled = original.Disabled
	 Disabling = original.Disabling
	 Enabled = original.Enabled
	 Enabling = original.Enabling
	 Failed = original.Failed
	 EndpointResourceStateCreating = original.EndpointResourceStateCreating
	 EndpointResourceStateDeleting = original.EndpointResourceStateDeleting
	 EndpointResourceStateRunning = original.EndpointResourceStateRunning
	 EndpointResourceStateStarting = original.EndpointResourceStateStarting
	 EndpointResourceStateStopped = original.EndpointResourceStateStopped
	 EndpointResourceStateStopping = original.EndpointResourceStateStopping
	 Allow = original.Allow
	 Block = original.Block
	 OriginResourceStateActive = original.OriginResourceStateActive
	 OriginResourceStateCreating = original.OriginResourceStateCreating
	 OriginResourceStateDeleting = original.OriginResourceStateDeleting
	 ProfileResourceStateActive = original.ProfileResourceStateActive
	 ProfileResourceStateCreating = original.ProfileResourceStateCreating
	 ProfileResourceStateDeleting = original.ProfileResourceStateDeleting
	 ProfileResourceStateDisabled = original.ProfileResourceStateDisabled
	 BypassCaching = original.BypassCaching
	 IgnoreQueryString = original.IgnoreQueryString
	 NotSet = original.NotSet
	 UseQueryString = original.UseQueryString
	 MicrosoftCdnProfilesEndpoints = original.MicrosoftCdnProfilesEndpoints
	 CustomVerizon = original.CustomVerizon
	 PremiumVerizon = original.PremiumVerizon
	 StandardAkamai = original.StandardAkamai
	 StandardChinaCdn = original.StandardChinaCdn
	 StandardVerizon = original.StandardVerizon
)

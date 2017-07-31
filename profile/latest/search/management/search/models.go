package search

import (
	 original "github.com/Azure/azure-sdk-for-go/service/search/management/2015-08-19/search"
)

type (
	 QueryKeysClient = original.QueryKeysClient
	 ServicesClient = original.ServicesClient
	 AdminKeysClient = original.AdminKeysClient
	 ManagementClient = original.ManagementClient
	 AdminKeyKind = original.AdminKeyKind
	 HostingMode = original.HostingMode
	 ProvisioningState = original.ProvisioningState
	 ServiceStatus = original.ServiceStatus
	 SkuName = original.SkuName
	 UnavailableNameReason = original.UnavailableNameReason
	 AdminKeyResult = original.AdminKeyResult
	 CheckNameAvailabilityInput = original.CheckNameAvailabilityInput
	 CheckNameAvailabilityOutput = original.CheckNameAvailabilityOutput
	 CloudError = original.CloudError
	 CloudErrorBody = original.CloudErrorBody
	 ListQueryKeysResult = original.ListQueryKeysResult
	 QueryKey = original.QueryKey
	 Resource = original.Resource
	 Service = original.Service
	 ServiceListResult = original.ServiceListResult
	 ServiceProperties = original.ServiceProperties
	 Sku = original.Sku
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Primary = original.Primary
	 Secondary = original.Secondary
	 Default = original.Default
	 HighDensity = original.HighDensity
	 Failed = original.Failed
	 Provisioning = original.Provisioning
	 Succeeded = original.Succeeded
	 ServiceStatusDegraded = original.ServiceStatusDegraded
	 ServiceStatusDeleting = original.ServiceStatusDeleting
	 ServiceStatusDisabled = original.ServiceStatusDisabled
	 ServiceStatusError = original.ServiceStatusError
	 ServiceStatusProvisioning = original.ServiceStatusProvisioning
	 ServiceStatusRunning = original.ServiceStatusRunning
	 Basic = original.Basic
	 Free = original.Free
	 Standard = original.Standard
	 Standard2 = original.Standard2
	 Standard3 = original.Standard3
	 AlreadyExists = original.AlreadyExists
	 Invalid = original.Invalid
)

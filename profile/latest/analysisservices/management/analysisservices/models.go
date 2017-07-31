package analysisservices

import (
	 original "github.com/Azure/azure-sdk-for-go/service/analysisservices/management/2016-05-16/analysisservices"
)

type (
	 ConnectionMode = original.ConnectionMode
	 ProvisioningState = original.ProvisioningState
	 SkuTier = original.SkuTier
	 State = original.State
	 Resource = original.Resource
	 ResourceSku = original.ResourceSku
	 Server = original.Server
	 ServerAdministrators = original.ServerAdministrators
	 ServerMutableProperties = original.ServerMutableProperties
	 ServerProperties = original.ServerProperties
	 Servers = original.Servers
	 ServerUpdateParameters = original.ServerUpdateParameters
	 SkuDetailsForExistingResource = original.SkuDetailsForExistingResource
	 SkuEnumerationForExistingResourceResult = original.SkuEnumerationForExistingResourceResult
	 SkuEnumerationForNewResourceResult = original.SkuEnumerationForNewResourceResult
	 ServersClient = original.ServersClient
	 ManagementClient = original.ManagementClient
)

const (
	 All = original.All
	 ReadOnly = original.ReadOnly
	 Deleting = original.Deleting
	 Failed = original.Failed
	 Paused = original.Paused
	 Pausing = original.Pausing
	 Preparing = original.Preparing
	 Provisioning = original.Provisioning
	 Resuming = original.Resuming
	 Scaling = original.Scaling
	 Succeeded = original.Succeeded
	 Suspended = original.Suspended
	 Suspending = original.Suspending
	 Updating = original.Updating
	 Basic = original.Basic
	 Development = original.Development
	 Standard = original.Standard
	 StateDeleting = original.StateDeleting
	 StateFailed = original.StateFailed
	 StatePaused = original.StatePaused
	 StatePausing = original.StatePausing
	 StatePreparing = original.StatePreparing
	 StateProvisioning = original.StateProvisioning
	 StateResuming = original.StateResuming
	 StateScaling = original.StateScaling
	 StateSucceeded = original.StateSucceeded
	 StateSuspended = original.StateSuspended
	 StateSuspending = original.StateSuspending
	 StateUpdating = original.StateUpdating
	 DefaultBaseURI = original.DefaultBaseURI
)

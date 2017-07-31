package managedapplications

import (
	 original "github.com/Azure/azure-sdk-for-go/service/resources/management/2016-09-01-preview/managedapplications"
)

type (
	 AppliancesClient = original.AppliancesClient
	 ManagementClient = original.ManagementClient
	 ApplianceArtifactType = original.ApplianceArtifactType
	 ApplianceLockLevel = original.ApplianceLockLevel
	 ProvisioningState = original.ProvisioningState
	 ResourceIdentityType = original.ResourceIdentityType
	 Appliance = original.Appliance
	 ApplianceArtifact = original.ApplianceArtifact
	 ApplianceDefinition = original.ApplianceDefinition
	 ApplianceDefinitionListResult = original.ApplianceDefinitionListResult
	 ApplianceDefinitionProperties = original.ApplianceDefinitionProperties
	 ApplianceListResult = original.ApplianceListResult
	 AppliancePatchable = original.AppliancePatchable
	 ApplianceProperties = original.ApplianceProperties
	 AppliancePropertiesPatchable = original.AppliancePropertiesPatchable
	 ApplianceProviderAuthorization = original.ApplianceProviderAuthorization
	 ErrorResponse = original.ErrorResponse
	 GenericResource = original.GenericResource
	 Identity = original.Identity
	 Plan = original.Plan
	 PlanPatchable = original.PlanPatchable
	 Resource = original.Resource
	 Sku = original.Sku
	 ApplianceDefinitionsClient = original.ApplianceDefinitionsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Custom = original.Custom
	 Template = original.Template
	 CanNotDelete = original.CanNotDelete
	 None = original.None
	 ReadOnly = original.ReadOnly
	 Accepted = original.Accepted
	 Canceled = original.Canceled
	 Created = original.Created
	 Creating = original.Creating
	 Deleted = original.Deleted
	 Deleting = original.Deleting
	 Failed = original.Failed
	 Ready = original.Ready
	 Running = original.Running
	 Succeeded = original.Succeeded
	 Updating = original.Updating
	 SystemAssigned = original.SystemAssigned
)

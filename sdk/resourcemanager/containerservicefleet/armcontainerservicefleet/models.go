// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armcontainerservicefleet

import "time"

// APIServerAccessProfile - Access profile for the Fleet hub API server.
type APIServerAccessProfile struct {
	// Whether to create the Fleet hub as a private cluster or not.
	EnablePrivateCluster *bool

	// Whether to enable apiserver vnet integration for the Fleet hub or not.
	EnableVnetIntegration *bool

	// The subnet to be used when apiserver vnet integration is enabled. It is required when creating a new Fleet with BYO vnet.
	SubnetID *string
}

// AgentProfile - Agent profile for the Fleet hub.
type AgentProfile struct {
	// The ID of the subnet which the Fleet hub node will join on startup. If this is not specified, a vnet and subnet will be
	// generated and used.
	SubnetID *string

	// The virtual machine size of the Fleet hub.
	VMSize *string
}

// AutoUpgradeNodeImageSelection - The node image upgrade to be applied to the target clusters in auto upgrade.
type AutoUpgradeNodeImageSelection struct {
	// REQUIRED; The node image upgrade type.
	Type *AutoUpgradeNodeImageSelectionType
}

// AutoUpgradeProfile - The AutoUpgradeProfile resource.
type AutoUpgradeProfile struct {
	// The resource-specific properties for this resource.
	Properties *AutoUpgradeProfileProperties

	// READ-ONLY; The name of the AutoUpgradeProfile resource.
	Name *string

	// READ-ONLY; If eTag is provided in the response body, it may also be provided as a header per the normal etag convention.
	// Entity tags are used for comparing two or more entities from the same requested resource. HTTP/1.1 uses entity tags in
	// the etag (section 14.19), If-Match (section 14.24), If-None-Match (section 14.26), and If-Range (section 14.27) header
	// fields.
	ETag *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AutoUpgradeProfileListResult - The response of a AutoUpgradeProfile list operation.
type AutoUpgradeProfileListResult struct {
	// REQUIRED; The AutoUpgradeProfile items on this page
	Value []*AutoUpgradeProfile

	// The link to the next page of items
	NextLink *string
}

// AutoUpgradeProfileProperties - The properties of the AutoUpgradeProfile.
type AutoUpgradeProfileProperties struct {
	// REQUIRED; Configures how auto-upgrade will be run.
	Channel *UpgradeChannel

	// The status of the auto upgrade profile.
	AutoUpgradeProfileStatus *AutoUpgradeProfileStatus

	// If set to False: the auto upgrade has effect - target managed clusters will be upgraded on schedule.
	// If set to True: the auto upgrade has no effect - no upgrade will be run on the target managed clusters.
	// This is a boolean and not an enum because enabled/disabled are all available states of the auto upgrade profile.
	// By default, this is set to False.
	Disabled *bool

	// The node image upgrade to be applied to the target clusters in auto upgrade.
	NodeImageSelection *AutoUpgradeNodeImageSelection

	// The resource id of the UpdateStrategy resource to reference. If not specified, the auto upgrade will run on all clusters
	// which are members of the fleet.
	UpdateStrategyID *string

	// READ-ONLY; The provisioning state of the AutoUpgradeProfile resource.
	ProvisioningState *AutoUpgradeProfileProvisioningState
}

// AutoUpgradeProfileStatus is the status of an auto upgrade profile.
type AutoUpgradeProfileStatus struct {
	// READ-ONLY; The error details of the last trigger.
	LastTriggerError *ErrorDetail

	// READ-ONLY; The status of the last AutoUpgrade trigger.
	LastTriggerStatus *AutoUpgradeLastTriggerStatus

	// READ-ONLY; The target Kubernetes version or node image versions of the last trigger.
	LastTriggerUpgradeVersions []*string

	// READ-ONLY; The UTC time of the last attempt to automatically create and start an UpdateRun as triggered by the release
	// of new versions.
	LastTriggeredAt *time.Time
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info *ErrorAdditionalInfoInfo

	// READ-ONLY; The additional info type.
	Type *string
}

type ErrorAdditionalInfoInfo struct {
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo

	// READ-ONLY; The error code.
	Code *string

	// READ-ONLY; The error details.
	Details []*ErrorDetail

	// READ-ONLY; The error message.
	Message *string

	// READ-ONLY; The error target.
	Target *string
}

// Fleet - The Fleet resource.
type Fleet struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// READ-ONLY; The name of the Fleet resource.
	Name *string

	// Managed identity.
	Identity *ManagedServiceIdentity

	// The resource-specific properties for this resource.
	Properties *FleetProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; If eTag is provided in the response body, it may also be provided as a header per the normal etag convention.
	// Entity tags are used for comparing two or more entities from the same requested resource. HTTP/1.1 uses entity tags in
	// the etag (section 14.19), If-Match (section 14.24), If-None-Match (section 14.26), and If-Range (section 14.27) header
	// fields.
	ETag *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// FleetCredentialResult - One credential result item.
type FleetCredentialResult struct {
	// READ-ONLY; The name of the credential.
	Name *string

	// READ-ONLY; Base64-encoded Kubernetes configuration file.
	Value []byte
}

// FleetCredentialResults - The Credential results response.
type FleetCredentialResults struct {
	// READ-ONLY; Array of base64-encoded Kubernetes configuration files.
	Kubeconfigs []*FleetCredentialResult
}

// FleetHubProfile - The FleetHubProfile configures the fleet hub.
type FleetHubProfile struct {
	// The access profile for the Fleet hub API server.
	APIServerAccessProfile *APIServerAccessProfile

	// The agent profile for the Fleet hub.
	AgentProfile *AgentProfile

	// DNS prefix used to create the FQDN for the Fleet hub.
	DNSPrefix *string

	// READ-ONLY; The FQDN of the Fleet hub.
	Fqdn *string

	// READ-ONLY; The Kubernetes version of the Fleet hub.
	KubernetesVersion *string

	// READ-ONLY; The Azure Portal FQDN of the Fleet hub.
	PortalFqdn *string
}

// FleetListResult - The response of a Fleet list operation.
type FleetListResult struct {
	// REQUIRED; The Fleet items on this page
	Value []*Fleet

	// The link to the next page of items
	NextLink *string
}

// FleetMember - A member of the Fleet. It contains a reference to an existing Kubernetes cluster on Azure.
type FleetMember struct {
	// The resource-specific properties for this resource.
	Properties *FleetMemberProperties

	// READ-ONLY; The name of the Fleet member resource.
	Name *string

	// READ-ONLY; If eTag is provided in the response body, it may also be provided as a header per the normal etag convention.
	// Entity tags are used for comparing two or more entities from the same requested resource. HTTP/1.1 uses entity tags in
	// the etag (section 14.19), If-Match (section 14.24), If-None-Match (section 14.26), and If-Range (section 14.27) header
	// fields.
	ETag *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// FleetMemberListResult - The response of a FleetMember list operation.
type FleetMemberListResult struct {
	// REQUIRED; The FleetMember items on this page
	Value []*FleetMember

	// The link to the next page of items
	NextLink *string
}

// FleetMemberProperties - A member of the Fleet. It contains a reference to an existing Kubernetes cluster on Azure.
type FleetMemberProperties struct {
	// REQUIRED; The ARM resource id of the cluster that joins the Fleet. Must be a valid Azure resource id. e.g.: '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{clusterName}'.
	ClusterResourceID *string

	// The group this member belongs to for multi-cluster update management.
	Group *string

	// READ-ONLY; The status of the last operation.
	ProvisioningState *FleetMemberProvisioningState

	// READ-ONLY; Status information of the last operation for fleet member.
	Status *FleetMemberStatus
}

// FleetMemberStatus - Status information for the fleet member
type FleetMemberStatus struct {
	// READ-ONLY; The last operation error of the fleet member
	LastOperationError *ErrorDetail

	// READ-ONLY; The last operation ID for the fleet member
	LastOperationID *string
}

// FleetMemberUpdate - The type used for update operations of the FleetMember.
type FleetMemberUpdate struct {
	// The resource-specific properties for this resource.
	Properties *FleetMemberUpdateProperties
}

// FleetMemberUpdateProperties - The updatable properties of the FleetMember.
type FleetMemberUpdateProperties struct {
	// The group this member belongs to for multi-cluster update management.
	Group *string
}

// FleetPatch - Properties of a Fleet that can be patched.
type FleetPatch struct {
	// Managed identity.
	Identity *ManagedServiceIdentity

	// Resource tags.
	Tags map[string]*string
}

// FleetProperties - Fleet properties.
type FleetProperties struct {
	// The FleetHubProfile configures the Fleet's hub.
	HubProfile *FleetHubProfile

	// READ-ONLY; The status of the last operation.
	ProvisioningState *FleetProvisioningState

	// READ-ONLY; Status information for the fleet.
	Status *FleetStatus
}

// FleetStatus - Status information for the fleet.
type FleetStatus struct {
	// READ-ONLY; The last operation error for the fleet.
	LastOperationError *ErrorDetail

	// READ-ONLY; The last operation ID for the fleet.
	LastOperationID *string
}

// FleetUpdateStrategy - Defines a multi-stage process to perform update operations across members of a Fleet.
type FleetUpdateStrategy struct {
	// The resource-specific properties for this resource.
	Properties *FleetUpdateStrategyProperties

	// READ-ONLY; The name of the UpdateStrategy resource.
	Name *string

	// READ-ONLY; If eTag is provided in the response body, it may also be provided as a header per the normal etag convention.
	// Entity tags are used for comparing two or more entities from the same requested resource. HTTP/1.1 uses entity tags in
	// the etag (section 14.19), If-Match (section 14.24), If-None-Match (section 14.26), and If-Range (section 14.27) header
	// fields.
	ETag *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// FleetUpdateStrategyListResult - The response of a FleetUpdateStrategy list operation.
type FleetUpdateStrategyListResult struct {
	// REQUIRED; The FleetUpdateStrategy items on this page
	Value []*FleetUpdateStrategy

	// The link to the next page of items
	NextLink *string
}

// FleetUpdateStrategyProperties - The properties of the UpdateStrategy.
type FleetUpdateStrategyProperties struct {
	// REQUIRED; Defines the update sequence of the clusters.
	Strategy *UpdateRunStrategy

	// READ-ONLY; The provisioning state of the UpdateStrategy resource.
	ProvisioningState *FleetUpdateStrategyProvisioningState
}

// ManagedClusterUpdate - The update to be applied to the ManagedClusters.
type ManagedClusterUpdate struct {
	// REQUIRED; The upgrade to apply to the ManagedClusters.
	Upgrade *ManagedClusterUpgradeSpec

	// The node image upgrade to be applied to the target nodes in update run.
	NodeImageSelection *NodeImageSelection
}

// ManagedClusterUpgradeSpec - The upgrade to apply to a ManagedCluster.
type ManagedClusterUpgradeSpec struct {
	// REQUIRED; ManagedClusterUpgradeType is the type of upgrade to be applied.
	Type *ManagedClusterUpgradeType

	// The Kubernetes version to upgrade the member clusters to.
	KubernetesVersion *string
}

// ManagedServiceIdentity - Managed service identity (system assigned and/or user assigned identities)
type ManagedServiceIdentity struct {
	// REQUIRED; The type of managed identity assigned to this resource.
	Type *ManagedServiceIdentityType

	// The identities assigned to this resource by the user.
	UserAssignedIdentities map[string]*UserAssignedIdentity

	// READ-ONLY; The service principal ID of the system assigned identity. This property will only be provided for a system assigned
	// identity.
	PrincipalID *string

	// READ-ONLY; The tenant ID of the system assigned identity. This property will only be provided for a system assigned identity.
	TenantID *string
}

// MemberUpdateStatus - The status of a member update operation.
type MemberUpdateStatus struct {
	// READ-ONLY; The Azure resource id of the target Kubernetes cluster.
	ClusterResourceID *string

	// READ-ONLY; The status message after processing the member update operation.
	Message *string

	// READ-ONLY; The name of the FleetMember.
	Name *string

	// READ-ONLY; The operation resource id of the latest attempt to perform the operation.
	OperationID *string

	// READ-ONLY; The status of the MemberUpdate operation.
	Status *UpdateStatus
}

// NodeImageSelection - The node image upgrade to be applied to the target nodes in update run.
type NodeImageSelection struct {
	// REQUIRED; The node image upgrade type.
	Type *NodeImageSelectionType

	// Custom node image versions to upgrade the nodes to. This field is required if node image selection type is Custom. Otherwise,
	// it must be empty. For each node image family (e.g., 'AKSUbuntu-1804gen2containerd'), this field can contain at most one
	// version (e.g., only one of 'AKSUbuntu-1804gen2containerd-2023.01.12' or 'AKSUbuntu-1804gen2containerd-2023.02.12', not
	// both). If the nodes belong to a family without a matching image version in this field, they are not upgraded.
	CustomNodeImageVersions []*NodeImageVersion
}

// NodeImageSelectionStatus - The node image upgrade specs for the update run.
type NodeImageSelectionStatus struct {
	// READ-ONLY; The image versions to upgrade the nodes to.
	SelectedNodeImageVersions []*NodeImageVersion
}

// NodeImageVersion - The node upgrade image version.
type NodeImageVersion struct {
	// READ-ONLY; The image version to upgrade the nodes to (e.g., 'AKSUbuntu-1804gen2containerd-2022.12.13').
	Version *string
}

// Operation - Details of a REST API operation, returned from the Resource Provider Operations API
type Operation struct {
	// Localized display information for this particular operation.
	Display *OperationDisplay

	// READ-ONLY; Extensible enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
	ActionType *ActionType

	// READ-ONLY; Whether the operation applies to data-plane. This is "true" for data-plane operations and "false" for Azure
	// Resource Manager/control-plane operations.
	IsDataAction *bool

	// READ-ONLY; The name of the operation, as per Resource-Based Access Control (RBAC). Examples: "Microsoft.Compute/virtualMachines/write",
	// "Microsoft.Compute/virtualMachines/capture/action"
	Name *string

	// READ-ONLY; The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
	// value is "user,system"
	Origin *Origin
}

// OperationDisplay - Localized display information for and operation.
type OperationDisplay struct {
	// READ-ONLY; The short, localized friendly description of the operation; suitable for tool tips and detailed views.
	Description *string

	// READ-ONLY; The concise, localized friendly name for the operation; suitable for dropdowns. E.g. "Create or Update Virtual
	// Machine", "Restart Virtual Machine".
	Operation *string

	// READ-ONLY; The localized friendly form of the resource provider name, e.g. "Microsoft Monitoring Insights" or "Microsoft
	// Compute".
	Provider *string

	// READ-ONLY; The localized friendly name of the resource type related to this operation. E.g. "Virtual Machines" or "Job
	// Schedule Collections".
	Resource *string
}

// OperationListResult - A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to
// get the next set of results.
type OperationListResult struct {
	// REQUIRED; The Operation items on this page
	Value []*Operation

	// The link to the next page of items
	NextLink *string
}

// SkipProperties - The properties of a skip operation containing multiple skip requests.
type SkipProperties struct {
	// REQUIRED; The targets to skip.
	Targets []*SkipTarget
}

// SkipTarget - The definition of a single skip request.
type SkipTarget struct {
	// REQUIRED; The skip target's name.
	// To skip a member/group/stage, use the member/group/stage's name;
	// Tp skip an after stage wait, use the parent stage's name.
	Name *string

	// REQUIRED; The skip target type.
	Type *TargetType
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// The identity that created the resource.
	CreatedBy *string

	// The type of identity that created the resource.
	CreatedByType *CreatedByType

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time

	// The identity that last modified the resource.
	LastModifiedBy *string

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType
}

// UpdateGroup - A group to be updated.
type UpdateGroup struct {
	// REQUIRED; Name of the group.
	// It must match a group name of an existing fleet member.
	Name *string
}

// UpdateGroupStatus - The status of a UpdateGroup.
type UpdateGroupStatus struct {
	// READ-ONLY; The list of member this UpdateGroup updates.
	Members []*MemberUpdateStatus

	// READ-ONLY; The name of the UpdateGroup.
	Name *string

	// READ-ONLY; The status of the UpdateGroup.
	Status *UpdateStatus
}

// UpdateRun - A multi-stage process to perform update operations across members of a Fleet.
type UpdateRun struct {
	// The resource-specific properties for this resource.
	Properties *UpdateRunProperties

	// READ-ONLY; The name of the UpdateRun resource.
	Name *string

	// READ-ONLY; If eTag is provided in the response body, it may also be provided as a header per the normal etag convention.
	// Entity tags are used for comparing two or more entities from the same requested resource. HTTP/1.1 uses entity tags in
	// the etag (section 14.19), If-Match (section 14.24), If-None-Match (section 14.26), and If-Range (section 14.27) header
	// fields.
	ETag *string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; Azure Resource Manager metadata containing createdBy and modifiedBy information.
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// UpdateRunListResult - The response of a UpdateRun list operation.
type UpdateRunListResult struct {
	// REQUIRED; The UpdateRun items on this page
	Value []*UpdateRun

	// The link to the next page of items
	NextLink *string
}

// UpdateRunProperties - The properties of the UpdateRun.
type UpdateRunProperties struct {
	// REQUIRED; The update to be applied to all clusters in the UpdateRun. The managedClusterUpdate can be modified until the
	// run is started.
	ManagedClusterUpdate *ManagedClusterUpdate

	// The strategy defines the order in which the clusters will be updated.
	// If not set, all members will be updated sequentially. The UpdateRun status will show a single UpdateStage and a single
	// UpdateGroup targeting all members.
	// The strategy of the UpdateRun can be modified until the run is started.
	Strategy *UpdateRunStrategy

	// The resource id of the FleetUpdateStrategy resource to reference.
	// When creating a new run, there are three ways to define a strategy for the run:
	// 1. Define a new strategy in place: Set the "strategy" field.
	// 2. Use an existing strategy: Set the "updateStrategyId" field. (since 2023-08-15-preview)
	// 3. Use the default strategy to update all the members one by one: Leave both "updateStrategyId" and "strategy" unset. (since
	// 2023-08-15-preview)
	// Setting both "updateStrategyId" and "strategy" is invalid.
	// UpdateRuns created by "updateStrategyId" snapshot the referenced UpdateStrategy at the time of creation and store it in
	// the "strategy" field.
	// Subsequent changes to the referenced FleetUpdateStrategy resource do not propagate.
	// UpdateRunStrategy changes can be made directly on the "strategy" field before launching the UpdateRun.
	UpdateStrategyID *string

	// READ-ONLY; AutoUpgradeProfileId is the id of an auto upgrade profile resource.
	AutoUpgradeProfileID *string

	// READ-ONLY; The provisioning state of the UpdateRun resource.
	ProvisioningState *UpdateRunProvisioningState

	// READ-ONLY; The status of the UpdateRun.
	Status *UpdateRunStatus
}

// UpdateRunStatus - The status of a UpdateRun.
type UpdateRunStatus struct {
	// READ-ONLY; The node image upgrade specs for the update run. It is only set in update run when `NodeImageSelection.type`
	// is `Consistent`.
	NodeImageSelection *NodeImageSelectionStatus

	// READ-ONLY; The stages composing an update run. Stages are run sequentially withing an UpdateRun.
	Stages []*UpdateStageStatus

	// READ-ONLY; The status of the UpdateRun.
	Status *UpdateStatus
}

// UpdateRunStrategy - Defines the update sequence of the clusters via stages and groups.
// Stages within a run are executed sequentially one after another.
// Groups within a stage are executed in parallel.
// Member clusters within a group are updated sequentially one after another.
// A valid strategy contains no duplicate groups within or across stages.
type UpdateRunStrategy struct {
	// REQUIRED; The list of stages that compose this update run. Min size: 1.
	Stages []*UpdateStage
}

// UpdateStage - Defines a stage which contains the groups to update and the steps to take (e.g., wait for a time period)
// before starting the next stage.
type UpdateStage struct {
	// REQUIRED; The name of the stage. Must be unique within the UpdateRun.
	Name *string

	// The time in seconds to wait at the end of this stage before starting the next one. Defaults to 0 seconds if unspecified.
	AfterStageWaitInSeconds *int32

	// Defines the groups to be executed in parallel in this stage. Duplicate groups are not allowed. Min size: 1.
	Groups []*UpdateGroup
}

// UpdateStageStatus - The status of a UpdateStage.
type UpdateStageStatus struct {
	// READ-ONLY; The status of the wait period configured on the UpdateStage.
	AfterStageWaitStatus *WaitStatus

	// READ-ONLY; The list of groups to be updated as part of this UpdateStage.
	Groups []*UpdateGroupStatus

	// READ-ONLY; The name of the UpdateStage.
	Name *string

	// READ-ONLY; The status of the UpdateStage.
	Status *UpdateStatus
}

// UpdateStatus - The status for an operation or group of operations.
type UpdateStatus struct {
	// READ-ONLY; The time the operation or group was completed.
	CompletedTime *time.Time

	// READ-ONLY; The error details when a failure is encountered.
	Error *ErrorDetail

	// READ-ONLY; The time the operation or group was started.
	StartTime *time.Time

	// READ-ONLY; The State of the operation or group.
	State *UpdateState
}

// UserAssignedIdentity - User assigned identity properties
type UserAssignedIdentity struct {
	// READ-ONLY; The client ID of the assigned identity.
	ClientID *string

	// READ-ONLY; The principal ID of the assigned identity.
	PrincipalID *string
}

// WaitStatus - The status of the wait duration.
type WaitStatus struct {
	// READ-ONLY; The status of the wait duration.
	Status *UpdateStatus

	// READ-ONLY; The wait duration configured in seconds.
	WaitDurationInSeconds *int32
}

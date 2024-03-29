//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridnetwork

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridnetwork/armhybridnetwork"
	moduleVersion = "v2.0.0"
)

// ActionType - Enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
type ActionType string

const (
	ActionTypeInternal ActionType = "Internal"
)

// PossibleActionTypeValues returns the possible values for the ActionType const type.
func PossibleActionTypeValues() []ActionType {
	return []ActionType{
		ActionTypeInternal,
	}
}

// ApplicationEnablement - The application enablement.
type ApplicationEnablement string

const (
	ApplicationEnablementDisabled ApplicationEnablement = "Disabled"
	ApplicationEnablementEnabled  ApplicationEnablement = "Enabled"
	ApplicationEnablementUnknown  ApplicationEnablement = "Unknown"
)

// PossibleApplicationEnablementValues returns the possible values for the ApplicationEnablement const type.
func PossibleApplicationEnablementValues() []ApplicationEnablement {
	return []ApplicationEnablement{
		ApplicationEnablementDisabled,
		ApplicationEnablementEnabled,
		ApplicationEnablementUnknown,
	}
}

// ArtifactManifestState - The artifact manifest state.
type ArtifactManifestState string

const (
	ArtifactManifestStateSucceeded        ArtifactManifestState = "Succeeded"
	ArtifactManifestStateUnknown          ArtifactManifestState = "Unknown"
	ArtifactManifestStateUploaded         ArtifactManifestState = "Uploaded"
	ArtifactManifestStateUploading        ArtifactManifestState = "Uploading"
	ArtifactManifestStateValidating       ArtifactManifestState = "Validating"
	ArtifactManifestStateValidationFailed ArtifactManifestState = "ValidationFailed"
)

// PossibleArtifactManifestStateValues returns the possible values for the ArtifactManifestState const type.
func PossibleArtifactManifestStateValues() []ArtifactManifestState {
	return []ArtifactManifestState{
		ArtifactManifestStateSucceeded,
		ArtifactManifestStateUnknown,
		ArtifactManifestStateUploaded,
		ArtifactManifestStateUploading,
		ArtifactManifestStateValidating,
		ArtifactManifestStateValidationFailed,
	}
}

// ArtifactReplicationStrategy - The replication strategy.
type ArtifactReplicationStrategy string

const (
	ArtifactReplicationStrategySingleReplication ArtifactReplicationStrategy = "SingleReplication"
	ArtifactReplicationStrategyUnknown           ArtifactReplicationStrategy = "Unknown"
)

// PossibleArtifactReplicationStrategyValues returns the possible values for the ArtifactReplicationStrategy const type.
func PossibleArtifactReplicationStrategyValues() []ArtifactReplicationStrategy {
	return []ArtifactReplicationStrategy{
		ArtifactReplicationStrategySingleReplication,
		ArtifactReplicationStrategyUnknown,
	}
}

// ArtifactState - The artifact state.
type ArtifactState string

const (
	ArtifactStateActive     ArtifactState = "Active"
	ArtifactStateDeprecated ArtifactState = "Deprecated"
	ArtifactStatePreview    ArtifactState = "Preview"
	ArtifactStateUnknown    ArtifactState = "Unknown"
)

// PossibleArtifactStateValues returns the possible values for the ArtifactState const type.
func PossibleArtifactStateValues() []ArtifactState {
	return []ArtifactState{
		ArtifactStateActive,
		ArtifactStateDeprecated,
		ArtifactStatePreview,
		ArtifactStateUnknown,
	}
}

// ArtifactStoreType - The artifact store type.
type ArtifactStoreType string

const (
	ArtifactStoreTypeAzureContainerRegistry ArtifactStoreType = "AzureContainerRegistry"
	ArtifactStoreTypeAzureStorageAccount    ArtifactStoreType = "AzureStorageAccount"
	ArtifactStoreTypeUnknown                ArtifactStoreType = "Unknown"
)

// PossibleArtifactStoreTypeValues returns the possible values for the ArtifactStoreType const type.
func PossibleArtifactStoreTypeValues() []ArtifactStoreType {
	return []ArtifactStoreType{
		ArtifactStoreTypeAzureContainerRegistry,
		ArtifactStoreTypeAzureStorageAccount,
		ArtifactStoreTypeUnknown,
	}
}

// ArtifactType - The artifact type.
type ArtifactType string

const (
	ArtifactTypeArmTemplate  ArtifactType = "ArmTemplate"
	ArtifactTypeImageFile    ArtifactType = "ImageFile"
	ArtifactTypeOCIArtifact  ArtifactType = "OCIArtifact"
	ArtifactTypeUnknown      ArtifactType = "Unknown"
	ArtifactTypeVhdImageFile ArtifactType = "VhdImageFile"
)

// PossibleArtifactTypeValues returns the possible values for the ArtifactType const type.
func PossibleArtifactTypeValues() []ArtifactType {
	return []ArtifactType{
		ArtifactTypeArmTemplate,
		ArtifactTypeImageFile,
		ArtifactTypeOCIArtifact,
		ArtifactTypeUnknown,
		ArtifactTypeVhdImageFile,
	}
}

// AzureArcKubernetesArtifactType - The artifact type.
type AzureArcKubernetesArtifactType string

const (
	AzureArcKubernetesArtifactTypeHelmPackage AzureArcKubernetesArtifactType = "HelmPackage"
	AzureArcKubernetesArtifactTypeUnknown     AzureArcKubernetesArtifactType = "Unknown"
)

// PossibleAzureArcKubernetesArtifactTypeValues returns the possible values for the AzureArcKubernetesArtifactType const type.
func PossibleAzureArcKubernetesArtifactTypeValues() []AzureArcKubernetesArtifactType {
	return []AzureArcKubernetesArtifactType{
		AzureArcKubernetesArtifactTypeHelmPackage,
		AzureArcKubernetesArtifactTypeUnknown,
	}
}

// AzureCoreArtifactType - The artifact type.
type AzureCoreArtifactType string

const (
	AzureCoreArtifactTypeArmTemplate  AzureCoreArtifactType = "ArmTemplate"
	AzureCoreArtifactTypeUnknown      AzureCoreArtifactType = "Unknown"
	AzureCoreArtifactTypeVhdImageFile AzureCoreArtifactType = "VhdImageFile"
)

// PossibleAzureCoreArtifactTypeValues returns the possible values for the AzureCoreArtifactType const type.
func PossibleAzureCoreArtifactTypeValues() []AzureCoreArtifactType {
	return []AzureCoreArtifactType{
		AzureCoreArtifactTypeArmTemplate,
		AzureCoreArtifactTypeUnknown,
		AzureCoreArtifactTypeVhdImageFile,
	}
}

// AzureOperatorNexusArtifactType - The artifact type.
type AzureOperatorNexusArtifactType string

const (
	AzureOperatorNexusArtifactTypeArmTemplate AzureOperatorNexusArtifactType = "ArmTemplate"
	AzureOperatorNexusArtifactTypeImageFile   AzureOperatorNexusArtifactType = "ImageFile"
	AzureOperatorNexusArtifactTypeUnknown     AzureOperatorNexusArtifactType = "Unknown"
)

// PossibleAzureOperatorNexusArtifactTypeValues returns the possible values for the AzureOperatorNexusArtifactType const type.
func PossibleAzureOperatorNexusArtifactTypeValues() []AzureOperatorNexusArtifactType {
	return []AzureOperatorNexusArtifactType{
		AzureOperatorNexusArtifactTypeArmTemplate,
		AzureOperatorNexusArtifactTypeImageFile,
		AzureOperatorNexusArtifactTypeUnknown,
	}
}

// ConfigurationGroupValueConfigurationType - The secret type which indicates if secret or not.
type ConfigurationGroupValueConfigurationType string

const (
	ConfigurationGroupValueConfigurationTypeOpen    ConfigurationGroupValueConfigurationType = "Open"
	ConfigurationGroupValueConfigurationTypeSecret  ConfigurationGroupValueConfigurationType = "Secret"
	ConfigurationGroupValueConfigurationTypeUnknown ConfigurationGroupValueConfigurationType = "Unknown"
)

// PossibleConfigurationGroupValueConfigurationTypeValues returns the possible values for the ConfigurationGroupValueConfigurationType const type.
func PossibleConfigurationGroupValueConfigurationTypeValues() []ConfigurationGroupValueConfigurationType {
	return []ConfigurationGroupValueConfigurationType{
		ConfigurationGroupValueConfigurationTypeOpen,
		ConfigurationGroupValueConfigurationTypeSecret,
		ConfigurationGroupValueConfigurationTypeUnknown,
	}
}

// ContainerizedNetworkFunctionNFVIType - The network function type.
type ContainerizedNetworkFunctionNFVIType string

const (
	ContainerizedNetworkFunctionNFVITypeAzureArcKubernetes ContainerizedNetworkFunctionNFVIType = "AzureArcKubernetes"
	ContainerizedNetworkFunctionNFVITypeUnknown            ContainerizedNetworkFunctionNFVIType = "Unknown"
)

// PossibleContainerizedNetworkFunctionNFVITypeValues returns the possible values for the ContainerizedNetworkFunctionNFVIType const type.
func PossibleContainerizedNetworkFunctionNFVITypeValues() []ContainerizedNetworkFunctionNFVIType {
	return []ContainerizedNetworkFunctionNFVIType{
		ContainerizedNetworkFunctionNFVITypeAzureArcKubernetes,
		ContainerizedNetworkFunctionNFVITypeUnknown,
	}
}

// CreatedByType - The type of identity that created the resource.
type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns the possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{
		CreatedByTypeApplication,
		CreatedByTypeKey,
		CreatedByTypeManagedIdentity,
		CreatedByTypeUser,
	}
}

// CredentialType - The credential type.
type CredentialType string

const (
	CredentialTypeAzureContainerRegistryScopedToken CredentialType = "AzureContainerRegistryScopedToken"
	CredentialTypeAzureStorageAccountToken          CredentialType = "AzureStorageAccountToken"
	CredentialTypeUnknown                           CredentialType = "Unknown"
)

// PossibleCredentialTypeValues returns the possible values for the CredentialType const type.
func PossibleCredentialTypeValues() []CredentialType {
	return []CredentialType{
		CredentialTypeAzureContainerRegistryScopedToken,
		CredentialTypeAzureStorageAccountToken,
		CredentialTypeUnknown,
	}
}

// HTTPMethod - The http method of the request.
type HTTPMethod string

const (
	HTTPMethodDelete  HTTPMethod = "Delete"
	HTTPMethodGet     HTTPMethod = "Get"
	HTTPMethodPatch   HTTPMethod = "Patch"
	HTTPMethodPost    HTTPMethod = "Post"
	HTTPMethodPut     HTTPMethod = "Put"
	HTTPMethodUnknown HTTPMethod = "Unknown"
)

// PossibleHTTPMethodValues returns the possible values for the HTTPMethod const type.
func PossibleHTTPMethodValues() []HTTPMethod {
	return []HTTPMethod{
		HTTPMethodDelete,
		HTTPMethodGet,
		HTTPMethodPatch,
		HTTPMethodPost,
		HTTPMethodPut,
		HTTPMethodUnknown,
	}
}

// IDType - The resource reference arm id type.
type IDType string

const (
	IDTypeOpen    IDType = "Open"
	IDTypeSecret  IDType = "Secret"
	IDTypeUnknown IDType = "Unknown"
)

// PossibleIDTypeValues returns the possible values for the IDType const type.
func PossibleIDTypeValues() []IDType {
	return []IDType{
		IDTypeOpen,
		IDTypeSecret,
		IDTypeUnknown,
	}
}

// ManagedServiceIdentityType - Type of managed service identity (where both SystemAssigned and UserAssigned types are allowed).
type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeNone                       ManagedServiceIdentityType = "None"
	ManagedServiceIdentityTypeSystemAssigned             ManagedServiceIdentityType = "SystemAssigned"
	ManagedServiceIdentityTypeSystemAssignedUserAssigned ManagedServiceIdentityType = "SystemAssigned,UserAssigned"
	ManagedServiceIdentityTypeUserAssigned               ManagedServiceIdentityType = "UserAssigned"
)

// PossibleManagedServiceIdentityTypeValues returns the possible values for the ManagedServiceIdentityType const type.
func PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType {
	return []ManagedServiceIdentityType{
		ManagedServiceIdentityTypeNone,
		ManagedServiceIdentityTypeSystemAssigned,
		ManagedServiceIdentityTypeSystemAssignedUserAssigned,
		ManagedServiceIdentityTypeUserAssigned,
	}
}

// NFVIType - The NFVI type.
type NFVIType string

const (
	NFVITypeAzureArcKubernetes NFVIType = "AzureArcKubernetes"
	NFVITypeAzureCore          NFVIType = "AzureCore"
	NFVITypeAzureOperatorNexus NFVIType = "AzureOperatorNexus"
	NFVITypeUnknown            NFVIType = "Unknown"
)

// PossibleNFVITypeValues returns the possible values for the NFVIType const type.
func PossibleNFVITypeValues() []NFVIType {
	return []NFVIType{
		NFVITypeAzureArcKubernetes,
		NFVITypeAzureCore,
		NFVITypeAzureOperatorNexus,
		NFVITypeUnknown,
	}
}

// NetworkFunctionConfigurationType - The secret type which indicates if secret or not.
type NetworkFunctionConfigurationType string

const (
	NetworkFunctionConfigurationTypeOpen    NetworkFunctionConfigurationType = "Open"
	NetworkFunctionConfigurationTypeSecret  NetworkFunctionConfigurationType = "Secret"
	NetworkFunctionConfigurationTypeUnknown NetworkFunctionConfigurationType = "Unknown"
)

// PossibleNetworkFunctionConfigurationTypeValues returns the possible values for the NetworkFunctionConfigurationType const type.
func PossibleNetworkFunctionConfigurationTypeValues() []NetworkFunctionConfigurationType {
	return []NetworkFunctionConfigurationType{
		NetworkFunctionConfigurationTypeOpen,
		NetworkFunctionConfigurationTypeSecret,
		NetworkFunctionConfigurationTypeUnknown,
	}
}

// NetworkFunctionType - The network function type.
type NetworkFunctionType string

const (
	NetworkFunctionTypeContainerizedNetworkFunction NetworkFunctionType = "ContainerizedNetworkFunction"
	NetworkFunctionTypeUnknown                      NetworkFunctionType = "Unknown"
	NetworkFunctionTypeVirtualNetworkFunction       NetworkFunctionType = "VirtualNetworkFunction"
)

// PossibleNetworkFunctionTypeValues returns the possible values for the NetworkFunctionType const type.
func PossibleNetworkFunctionTypeValues() []NetworkFunctionType {
	return []NetworkFunctionType{
		NetworkFunctionTypeContainerizedNetworkFunction,
		NetworkFunctionTypeUnknown,
		NetworkFunctionTypeVirtualNetworkFunction,
	}
}

// Origin - The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
// value is "user,system"
type Origin string

const (
	OriginSystem     Origin = "system"
	OriginUser       Origin = "user"
	OriginUserSystem Origin = "user,system"
)

// PossibleOriginValues returns the possible values for the Origin const type.
func PossibleOriginValues() []Origin {
	return []Origin{
		OriginSystem,
		OriginUser,
		OriginUserSystem,
	}
}

// PodEventType - The type of pod event.
type PodEventType string

const (
	PodEventTypeNormal  PodEventType = "Normal"
	PodEventTypeWarning PodEventType = "Warning"
)

// PossiblePodEventTypeValues returns the possible values for the PodEventType const type.
func PossiblePodEventTypeValues() []PodEventType {
	return []PodEventType{
		PodEventTypeNormal,
		PodEventTypeWarning,
	}
}

// PodStatus - The status of a Pod.
type PodStatus string

const (
	PodStatusFailed      PodStatus = "Failed"
	PodStatusNotReady    PodStatus = "NotReady"
	PodStatusPending     PodStatus = "Pending"
	PodStatusRunning     PodStatus = "Running"
	PodStatusSucceeded   PodStatus = "Succeeded"
	PodStatusTerminating PodStatus = "Terminating"
	PodStatusUnknown     PodStatus = "Unknown"
)

// PossiblePodStatusValues returns the possible values for the PodStatus const type.
func PossiblePodStatusValues() []PodStatus {
	return []PodStatus{
		PodStatusFailed,
		PodStatusNotReady,
		PodStatusPending,
		PodStatusRunning,
		PodStatusSucceeded,
		PodStatusTerminating,
		PodStatusUnknown,
	}
}

// ProvisioningState - The current provisioning state.
type ProvisioningState string

const (
	ProvisioningStateAccepted   ProvisioningState = "Accepted"
	ProvisioningStateCanceled   ProvisioningState = "Canceled"
	ProvisioningStateConverging ProvisioningState = "Converging"
	ProvisioningStateDeleted    ProvisioningState = "Deleted"
	ProvisioningStateDeleting   ProvisioningState = "Deleting"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
	ProvisioningStateUnknown    ProvisioningState = "Unknown"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{
		ProvisioningStateAccepted,
		ProvisioningStateCanceled,
		ProvisioningStateConverging,
		ProvisioningStateDeleted,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateSucceeded,
		ProvisioningStateUnknown,
	}
}

// PublisherScope - Publisher Scope.
type PublisherScope string

const (
	PublisherScopePrivate PublisherScope = "Private"
	PublisherScopeUnknown PublisherScope = "Unknown"
)

// PossiblePublisherScopeValues returns the possible values for the PublisherScope const type.
func PossiblePublisherScopeValues() []PublisherScope {
	return []PublisherScope{
		PublisherScopePrivate,
		PublisherScopeUnknown,
	}
}

// SKUName - Name of this Sku
type SKUName string

const (
	SKUNameBasic    SKUName = "Basic"
	SKUNameStandard SKUName = "Standard"
)

// PossibleSKUNameValues returns the possible values for the SKUName const type.
func PossibleSKUNameValues() []SKUName {
	return []SKUName{
		SKUNameBasic,
		SKUNameStandard,
	}
}

// SKUTier - The SKU tier based on the SKU name.
type SKUTier string

const (
	SKUTierBasic    SKUTier = "Basic"
	SKUTierStandard SKUTier = "Standard"
)

// PossibleSKUTierValues returns the possible values for the SKUTier const type.
func PossibleSKUTierValues() []SKUTier {
	return []SKUTier{
		SKUTierBasic,
		SKUTierStandard,
	}
}

// Status - The component resource deployment status.
type Status string

const (
	StatusDeployed        Status = "Deployed"
	StatusDownloading     Status = "Downloading"
	StatusFailed          Status = "Failed"
	StatusInstalling      Status = "Installing"
	StatusPendingInstall  Status = "Pending-Install"
	StatusPendingRollback Status = "Pending-Rollback"
	StatusPendingUpgrade  Status = "Pending-Upgrade"
	StatusReinstalling    Status = "Reinstalling"
	StatusRollingback     Status = "Rollingback"
	StatusSuperseded      Status = "Superseded"
	StatusUninstalled     Status = "Uninstalled"
	StatusUninstalling    Status = "Uninstalling"
	StatusUnknown         Status = "Unknown"
	StatusUpgrading       Status = "Upgrading"
)

// PossibleStatusValues returns the possible values for the Status const type.
func PossibleStatusValues() []Status {
	return []Status{
		StatusDeployed,
		StatusDownloading,
		StatusFailed,
		StatusInstalling,
		StatusPendingInstall,
		StatusPendingRollback,
		StatusPendingUpgrade,
		StatusReinstalling,
		StatusRollingback,
		StatusSuperseded,
		StatusUninstalled,
		StatusUninstalling,
		StatusUnknown,
		StatusUpgrading,
	}
}

// TemplateType - The template type.
type TemplateType string

const (
	TemplateTypeArmTemplate TemplateType = "ArmTemplate"
	TemplateTypeUnknown     TemplateType = "Unknown"
)

// PossibleTemplateTypeValues returns the possible values for the TemplateType const type.
func PossibleTemplateTypeValues() []TemplateType {
	return []TemplateType{
		TemplateTypeArmTemplate,
		TemplateTypeUnknown,
	}
}

// Type - The resource element template type.
type Type string

const (
	TypeArmResourceDefinition     Type = "ArmResourceDefinition"
	TypeNetworkFunctionDefinition Type = "NetworkFunctionDefinition"
	TypeUnknown                   Type = "Unknown"
)

// PossibleTypeValues returns the possible values for the Type const type.
func PossibleTypeValues() []Type {
	return []Type{
		TypeArmResourceDefinition,
		TypeNetworkFunctionDefinition,
		TypeUnknown,
	}
}

// VersionState - The configuration group schema state.
type VersionState string

const (
	VersionStateActive           VersionState = "Active"
	VersionStateDeprecated       VersionState = "Deprecated"
	VersionStatePreview          VersionState = "Preview"
	VersionStateUnknown          VersionState = "Unknown"
	VersionStateValidating       VersionState = "Validating"
	VersionStateValidationFailed VersionState = "ValidationFailed"
)

// PossibleVersionStateValues returns the possible values for the VersionState const type.
func PossibleVersionStateValues() []VersionState {
	return []VersionState{
		VersionStateActive,
		VersionStateDeprecated,
		VersionStatePreview,
		VersionStateUnknown,
		VersionStateValidating,
		VersionStateValidationFailed,
	}
}

// VirtualNetworkFunctionNFVIType - The network function type.
type VirtualNetworkFunctionNFVIType string

const (
	VirtualNetworkFunctionNFVITypeAzureCore          VirtualNetworkFunctionNFVIType = "AzureCore"
	VirtualNetworkFunctionNFVITypeAzureOperatorNexus VirtualNetworkFunctionNFVIType = "AzureOperatorNexus"
	VirtualNetworkFunctionNFVITypeUnknown            VirtualNetworkFunctionNFVIType = "Unknown"
)

// PossibleVirtualNetworkFunctionNFVITypeValues returns the possible values for the VirtualNetworkFunctionNFVIType const type.
func PossibleVirtualNetworkFunctionNFVITypeValues() []VirtualNetworkFunctionNFVIType {
	return []VirtualNetworkFunctionNFVIType{
		VirtualNetworkFunctionNFVITypeAzureCore,
		VirtualNetworkFunctionNFVITypeAzureOperatorNexus,
		VirtualNetworkFunctionNFVITypeUnknown,
	}
}

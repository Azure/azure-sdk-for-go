//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcustomerinsights

const (
	moduleName    = "armcustomerinsights"
	moduleVersion = "v1.0.0"
)

// CalculationWindowTypes - The calculation window.
type CalculationWindowTypes string

const (
	CalculationWindowTypesLifetime CalculationWindowTypes = "Lifetime"
	CalculationWindowTypesHour     CalculationWindowTypes = "Hour"
	CalculationWindowTypesDay      CalculationWindowTypes = "Day"
	CalculationWindowTypesWeek     CalculationWindowTypes = "Week"
	CalculationWindowTypesMonth    CalculationWindowTypes = "Month"
)

// PossibleCalculationWindowTypesValues returns the possible values for the CalculationWindowTypes const type.
func PossibleCalculationWindowTypesValues() []CalculationWindowTypes {
	return []CalculationWindowTypes{
		CalculationWindowTypesLifetime,
		CalculationWindowTypesHour,
		CalculationWindowTypesDay,
		CalculationWindowTypesWeek,
		CalculationWindowTypesMonth,
	}
}

// CanonicalPropertyValueType - Type of canonical property value.
type CanonicalPropertyValueType string

const (
	CanonicalPropertyValueTypeCategorical        CanonicalPropertyValueType = "Categorical"
	CanonicalPropertyValueTypeDerivedCategorical CanonicalPropertyValueType = "DerivedCategorical"
	CanonicalPropertyValueTypeDerivedNumeric     CanonicalPropertyValueType = "DerivedNumeric"
	CanonicalPropertyValueTypeNumeric            CanonicalPropertyValueType = "Numeric"
)

// PossibleCanonicalPropertyValueTypeValues returns the possible values for the CanonicalPropertyValueType const type.
func PossibleCanonicalPropertyValueTypeValues() []CanonicalPropertyValueType {
	return []CanonicalPropertyValueType{
		CanonicalPropertyValueTypeCategorical,
		CanonicalPropertyValueTypeDerivedCategorical,
		CanonicalPropertyValueTypeDerivedNumeric,
		CanonicalPropertyValueTypeNumeric,
	}
}

// CardinalityTypes - The Relationship Cardinality.
type CardinalityTypes string

const (
	CardinalityTypesOneToOne   CardinalityTypes = "OneToOne"
	CardinalityTypesOneToMany  CardinalityTypes = "OneToMany"
	CardinalityTypesManyToMany CardinalityTypes = "ManyToMany"
)

// PossibleCardinalityTypesValues returns the possible values for the CardinalityTypes const type.
func PossibleCardinalityTypesValues() []CardinalityTypes {
	return []CardinalityTypes{
		CardinalityTypesOneToOne,
		CardinalityTypesOneToMany,
		CardinalityTypesManyToMany,
	}
}

// CompletionOperationTypes - The type of completion operation.
type CompletionOperationTypes string

const (
	CompletionOperationTypesDoNothing  CompletionOperationTypes = "DoNothing"
	CompletionOperationTypesDeleteFile CompletionOperationTypes = "DeleteFile"
	CompletionOperationTypesMoveFile   CompletionOperationTypes = "MoveFile"
)

// PossibleCompletionOperationTypesValues returns the possible values for the CompletionOperationTypes const type.
func PossibleCompletionOperationTypesValues() []CompletionOperationTypes {
	return []CompletionOperationTypes{
		CompletionOperationTypesDoNothing,
		CompletionOperationTypesDeleteFile,
		CompletionOperationTypesMoveFile,
	}
}

// ConnectorMappingStates - State of connector mapping.
type ConnectorMappingStates string

const (
	ConnectorMappingStatesCreating ConnectorMappingStates = "Creating"
	ConnectorMappingStatesCreated  ConnectorMappingStates = "Created"
	ConnectorMappingStatesFailed   ConnectorMappingStates = "Failed"
	ConnectorMappingStatesReady    ConnectorMappingStates = "Ready"
	ConnectorMappingStatesRunning  ConnectorMappingStates = "Running"
	ConnectorMappingStatesStopped  ConnectorMappingStates = "Stopped"
	ConnectorMappingStatesExpiring ConnectorMappingStates = "Expiring"
)

// PossibleConnectorMappingStatesValues returns the possible values for the ConnectorMappingStates const type.
func PossibleConnectorMappingStatesValues() []ConnectorMappingStates {
	return []ConnectorMappingStates{
		ConnectorMappingStatesCreating,
		ConnectorMappingStatesCreated,
		ConnectorMappingStatesFailed,
		ConnectorMappingStatesReady,
		ConnectorMappingStatesRunning,
		ConnectorMappingStatesStopped,
		ConnectorMappingStatesExpiring,
	}
}

// ConnectorStates - State of connector.
type ConnectorStates string

const (
	ConnectorStatesCreating ConnectorStates = "Creating"
	ConnectorStatesCreated  ConnectorStates = "Created"
	ConnectorStatesReady    ConnectorStates = "Ready"
	ConnectorStatesExpiring ConnectorStates = "Expiring"
	ConnectorStatesDeleting ConnectorStates = "Deleting"
	ConnectorStatesFailed   ConnectorStates = "Failed"
)

// PossibleConnectorStatesValues returns the possible values for the ConnectorStates const type.
func PossibleConnectorStatesValues() []ConnectorStates {
	return []ConnectorStates{
		ConnectorStatesCreating,
		ConnectorStatesCreated,
		ConnectorStatesReady,
		ConnectorStatesExpiring,
		ConnectorStatesDeleting,
		ConnectorStatesFailed,
	}
}

// ConnectorTypes - Type of connector.
type ConnectorTypes string

const (
	ConnectorTypesAzureBlob      ConnectorTypes = "AzureBlob"
	ConnectorTypesCRM            ConnectorTypes = "CRM"
	ConnectorTypesExchangeOnline ConnectorTypes = "ExchangeOnline"
	ConnectorTypesNone           ConnectorTypes = "None"
	ConnectorTypesOutbound       ConnectorTypes = "Outbound"
	ConnectorTypesSalesforce     ConnectorTypes = "Salesforce"
)

// PossibleConnectorTypesValues returns the possible values for the ConnectorTypes const type.
func PossibleConnectorTypesValues() []ConnectorTypes {
	return []ConnectorTypes{
		ConnectorTypesAzureBlob,
		ConnectorTypesCRM,
		ConnectorTypesExchangeOnline,
		ConnectorTypesNone,
		ConnectorTypesOutbound,
		ConnectorTypesSalesforce,
	}
}

// DataSourceType - The data source type.
type DataSourceType string

const (
	DataSourceTypeConnector       DataSourceType = "Connector"
	DataSourceTypeLinkInteraction DataSourceType = "LinkInteraction"
	DataSourceTypeSystemDefault   DataSourceType = "SystemDefault"
)

// PossibleDataSourceTypeValues returns the possible values for the DataSourceType const type.
func PossibleDataSourceTypeValues() []DataSourceType {
	return []DataSourceType{
		DataSourceTypeConnector,
		DataSourceTypeLinkInteraction,
		DataSourceTypeSystemDefault,
	}
}

// EntityType - Type of source entity.
type EntityType string

const (
	EntityTypeNone         EntityType = "None"
	EntityTypeProfile      EntityType = "Profile"
	EntityTypeInteraction  EntityType = "Interaction"
	EntityTypeRelationship EntityType = "Relationship"
)

// PossibleEntityTypeValues returns the possible values for the EntityType const type.
func PossibleEntityTypeValues() []EntityType {
	return []EntityType{
		EntityTypeNone,
		EntityTypeProfile,
		EntityTypeInteraction,
		EntityTypeRelationship,
	}
}

// EntityTypes - Type of entity.
type EntityTypes string

const (
	EntityTypesNone         EntityTypes = "None"
	EntityTypesProfile      EntityTypes = "Profile"
	EntityTypesInteraction  EntityTypes = "Interaction"
	EntityTypesRelationship EntityTypes = "Relationship"
)

// PossibleEntityTypesValues returns the possible values for the EntityTypes const type.
func PossibleEntityTypesValues() []EntityTypes {
	return []EntityTypes{
		EntityTypesNone,
		EntityTypesProfile,
		EntityTypesInteraction,
		EntityTypesRelationship,
	}
}

// ErrorManagementTypes - The type of error management to use for the mapping.
type ErrorManagementTypes string

const (
	ErrorManagementTypesRejectAndContinue ErrorManagementTypes = "RejectAndContinue"
	ErrorManagementTypesStopImport        ErrorManagementTypes = "StopImport"
	ErrorManagementTypesRejectUntilLimit  ErrorManagementTypes = "RejectUntilLimit"
)

// PossibleErrorManagementTypesValues returns the possible values for the ErrorManagementTypes const type.
func PossibleErrorManagementTypesValues() []ErrorManagementTypes {
	return []ErrorManagementTypes{
		ErrorManagementTypesRejectAndContinue,
		ErrorManagementTypesStopImport,
		ErrorManagementTypesRejectUntilLimit,
	}
}

// FrequencyTypes - The frequency to update.
type FrequencyTypes string

const (
	FrequencyTypesMinute FrequencyTypes = "Minute"
	FrequencyTypesHour   FrequencyTypes = "Hour"
	FrequencyTypesDay    FrequencyTypes = "Day"
	FrequencyTypesWeek   FrequencyTypes = "Week"
	FrequencyTypesMonth  FrequencyTypes = "Month"
)

// PossibleFrequencyTypesValues returns the possible values for the FrequencyTypes const type.
func PossibleFrequencyTypesValues() []FrequencyTypes {
	return []FrequencyTypes{
		FrequencyTypesMinute,
		FrequencyTypesHour,
		FrequencyTypesDay,
		FrequencyTypesWeek,
		FrequencyTypesMonth,
	}
}

// InstanceOperationType - Determines whether this link is supposed to create or delete instances if Link is NOT Reference
// Only.
type InstanceOperationType string

const (
	InstanceOperationTypeUpsert InstanceOperationType = "Upsert"
	InstanceOperationTypeDelete InstanceOperationType = "Delete"
)

// PossibleInstanceOperationTypeValues returns the possible values for the InstanceOperationType const type.
func PossibleInstanceOperationTypeValues() []InstanceOperationType {
	return []InstanceOperationType{
		InstanceOperationTypeUpsert,
		InstanceOperationTypeDelete,
	}
}

// KpiFunctions - The computation function for the KPI.
type KpiFunctions string

const (
	KpiFunctionsSum           KpiFunctions = "Sum"
	KpiFunctionsAvg           KpiFunctions = "Avg"
	KpiFunctionsMin           KpiFunctions = "Min"
	KpiFunctionsMax           KpiFunctions = "Max"
	KpiFunctionsLast          KpiFunctions = "Last"
	KpiFunctionsCount         KpiFunctions = "Count"
	KpiFunctionsNone          KpiFunctions = "None"
	KpiFunctionsCountDistinct KpiFunctions = "CountDistinct"
)

// PossibleKpiFunctionsValues returns the possible values for the KpiFunctions const type.
func PossibleKpiFunctionsValues() []KpiFunctions {
	return []KpiFunctions{
		KpiFunctionsSum,
		KpiFunctionsAvg,
		KpiFunctionsMin,
		KpiFunctionsMax,
		KpiFunctionsLast,
		KpiFunctionsCount,
		KpiFunctionsNone,
		KpiFunctionsCountDistinct,
	}
}

// LinkTypes - Link type.
type LinkTypes string

const (
	LinkTypesUpdateAlways LinkTypes = "UpdateAlways"
	LinkTypesCopyIfNull   LinkTypes = "CopyIfNull"
)

// PossibleLinkTypesValues returns the possible values for the LinkTypes const type.
func PossibleLinkTypesValues() []LinkTypes {
	return []LinkTypes{
		LinkTypesUpdateAlways,
		LinkTypesCopyIfNull,
	}
}

// PermissionTypes - Supported permission types.
type PermissionTypes string

const (
	PermissionTypesRead   PermissionTypes = "Read"
	PermissionTypesWrite  PermissionTypes = "Write"
	PermissionTypesManage PermissionTypes = "Manage"
)

// PossiblePermissionTypesValues returns the possible values for the PermissionTypes const type.
func PossiblePermissionTypesValues() []PermissionTypes {
	return []PermissionTypes{
		PermissionTypesRead,
		PermissionTypesWrite,
		PermissionTypesManage,
	}
}

// PredictionModelLifeCycle - Prediction model life cycle. When prediction is in PendingModelConfirmation status, it is allowed
// to update the status to PendingFeaturing or Active through API.
type PredictionModelLifeCycle string

const (
	PredictionModelLifeCycleActive                   PredictionModelLifeCycle = "Active"
	PredictionModelLifeCycleDeleted                  PredictionModelLifeCycle = "Deleted"
	PredictionModelLifeCycleDiscovering              PredictionModelLifeCycle = "Discovering"
	PredictionModelLifeCycleEvaluating               PredictionModelLifeCycle = "Evaluating"
	PredictionModelLifeCycleEvaluatingFailed         PredictionModelLifeCycle = "EvaluatingFailed"
	PredictionModelLifeCycleFailed                   PredictionModelLifeCycle = "Failed"
	PredictionModelLifeCycleFeaturing                PredictionModelLifeCycle = "Featuring"
	PredictionModelLifeCycleFeaturingFailed          PredictionModelLifeCycle = "FeaturingFailed"
	PredictionModelLifeCycleHumanIntervention        PredictionModelLifeCycle = "HumanIntervention"
	PredictionModelLifeCycleNew                      PredictionModelLifeCycle = "New"
	PredictionModelLifeCyclePendingDiscovering       PredictionModelLifeCycle = "PendingDiscovering"
	PredictionModelLifeCyclePendingFeaturing         PredictionModelLifeCycle = "PendingFeaturing"
	PredictionModelLifeCyclePendingModelConfirmation PredictionModelLifeCycle = "PendingModelConfirmation"
	PredictionModelLifeCyclePendingTraining          PredictionModelLifeCycle = "PendingTraining"
	PredictionModelLifeCycleProvisioning             PredictionModelLifeCycle = "Provisioning"
	PredictionModelLifeCycleProvisioningFailed       PredictionModelLifeCycle = "ProvisioningFailed"
	PredictionModelLifeCycleTraining                 PredictionModelLifeCycle = "Training"
	PredictionModelLifeCycleTrainingFailed           PredictionModelLifeCycle = "TrainingFailed"
)

// PossiblePredictionModelLifeCycleValues returns the possible values for the PredictionModelLifeCycle const type.
func PossiblePredictionModelLifeCycleValues() []PredictionModelLifeCycle {
	return []PredictionModelLifeCycle{
		PredictionModelLifeCycleActive,
		PredictionModelLifeCycleDeleted,
		PredictionModelLifeCycleDiscovering,
		PredictionModelLifeCycleEvaluating,
		PredictionModelLifeCycleEvaluatingFailed,
		PredictionModelLifeCycleFailed,
		PredictionModelLifeCycleFeaturing,
		PredictionModelLifeCycleFeaturingFailed,
		PredictionModelLifeCycleHumanIntervention,
		PredictionModelLifeCycleNew,
		PredictionModelLifeCyclePendingDiscovering,
		PredictionModelLifeCyclePendingFeaturing,
		PredictionModelLifeCyclePendingModelConfirmation,
		PredictionModelLifeCyclePendingTraining,
		PredictionModelLifeCycleProvisioning,
		PredictionModelLifeCycleProvisioningFailed,
		PredictionModelLifeCycleTraining,
		PredictionModelLifeCycleTrainingFailed,
	}
}

// ProvisioningStates - Provisioning state.
type ProvisioningStates string

const (
	ProvisioningStatesDeleting          ProvisioningStates = "Deleting"
	ProvisioningStatesExpiring          ProvisioningStates = "Expiring"
	ProvisioningStatesFailed            ProvisioningStates = "Failed"
	ProvisioningStatesHumanIntervention ProvisioningStates = "HumanIntervention"
	ProvisioningStatesProvisioning      ProvisioningStates = "Provisioning"
	ProvisioningStatesSucceeded         ProvisioningStates = "Succeeded"
)

// PossibleProvisioningStatesValues returns the possible values for the ProvisioningStates const type.
func PossibleProvisioningStatesValues() []ProvisioningStates {
	return []ProvisioningStates{
		ProvisioningStatesDeleting,
		ProvisioningStatesExpiring,
		ProvisioningStatesFailed,
		ProvisioningStatesHumanIntervention,
		ProvisioningStatesProvisioning,
		ProvisioningStatesSucceeded,
	}
}

// RoleTypes - Type of roles.
type RoleTypes string

const (
	RoleTypesAdmin        RoleTypes = "Admin"
	RoleTypesReader       RoleTypes = "Reader"
	RoleTypesManageAdmin  RoleTypes = "ManageAdmin"
	RoleTypesManageReader RoleTypes = "ManageReader"
	RoleTypesDataAdmin    RoleTypes = "DataAdmin"
	RoleTypesDataReader   RoleTypes = "DataReader"
)

// PossibleRoleTypesValues returns the possible values for the RoleTypes const type.
func PossibleRoleTypesValues() []RoleTypes {
	return []RoleTypes{
		RoleTypesAdmin,
		RoleTypesReader,
		RoleTypesManageAdmin,
		RoleTypesManageReader,
		RoleTypesDataAdmin,
		RoleTypesDataReader,
	}
}

// Status - The data source status.
type Status string

const (
	StatusActive  Status = "Active"
	StatusDeleted Status = "Deleted"
	StatusNone    Status = "None"
)

// PossibleStatusValues returns the possible values for the Status const type.
func PossibleStatusValues() []Status {
	return []Status{
		StatusActive,
		StatusDeleted,
		StatusNone,
	}
}

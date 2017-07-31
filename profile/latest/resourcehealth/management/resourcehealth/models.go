package resourcehealth

import (
	 original "github.com/Azure/azure-sdk-for-go/service/resourcehealth/management/2017-07-01/resourcehealth"
)

type (
	 AvailabilityStateValues = original.AvailabilityStateValues
	 ReasonChronicityTypes = original.ReasonChronicityTypes
	 AvailabilityStatus = original.AvailabilityStatus
	 AvailabilityStatusProperties = original.AvailabilityStatusProperties
	 AvailabilityStatusPropertiesRecentlyResolvedState = original.AvailabilityStatusPropertiesRecentlyResolvedState
	 AvailabilityStatusListResult = original.AvailabilityStatusListResult
	 ErrorResponse = original.ErrorResponse
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 RecommendedAction = original.RecommendedAction
	 ServiceImpactingEvent = original.ServiceImpactingEvent
	 ServiceImpactingEventIncidentProperties = original.ServiceImpactingEventIncidentProperties
	 ServiceImpactingEventStatus = original.ServiceImpactingEventStatus
	 OperationsClient = original.OperationsClient
	 AvailabilityStatusesClient = original.AvailabilityStatusesClient
	 ManagementClient = original.ManagementClient
)

const (
	 Available = original.Available
	 Unavailable = original.Unavailable
	 Unknown = original.Unknown
	 Persistent = original.Persistent
	 Transient = original.Transient
	 DefaultBaseURI = original.DefaultBaseURI
)

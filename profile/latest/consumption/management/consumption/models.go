package consumption

import (
	 original "github.com/Azure/azure-sdk-for-go/service/consumption/management/2017-04-24-preview/consumption"
)

type (
	 ManagementClient = original.ManagementClient
	 ErrorDetails = original.ErrorDetails
	 ErrorResponse = original.ErrorResponse
	 MeterDetails = original.MeterDetails
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 Resource = original.Resource
	 UsageDetail = original.UsageDetail
	 UsageDetailProperties = original.UsageDetailProperties
	 UsageDetailsListResult = original.UsageDetailsListResult
	 OperationsClient = original.OperationsClient
	 UsageDetailsClient = original.UsageDetailsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)

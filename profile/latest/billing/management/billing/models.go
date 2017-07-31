package billing

import (
	 original "github.com/Azure/azure-sdk-for-go/service/billing/management/2017-04-24-preview/billing"
)

type (
	 DownloadURL = original.DownloadURL
	 ErrorDetails = original.ErrorDetails
	 ErrorResponse = original.ErrorResponse
	 Invoice = original.Invoice
	 InvoiceProperties = original.InvoiceProperties
	 InvoicesListResult = original.InvoicesListResult
	 Operation = original.Operation
	 OperationDisplay = original.OperationDisplay
	 OperationListResult = original.OperationListResult
	 Period = original.Period
	 PeriodProperties = original.PeriodProperties
	 PeriodsListResult = original.PeriodsListResult
	 Resource = original.Resource
	 OperationsClient = original.OperationsClient
	 PeriodsClient = original.PeriodsClient
	 ManagementClient = original.ManagementClient
	 InvoicesClient = original.InvoicesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)

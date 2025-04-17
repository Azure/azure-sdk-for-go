// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armcarbonoptimization

// CarbonServiceClientQueryCarbonEmissionDataAvailableDateRangeResponse contains the response from method CarbonServiceClient.QueryCarbonEmissionDataAvailableDateRange.
type CarbonServiceClientQueryCarbonEmissionDataAvailableDateRangeResponse struct {
	// Response for available date range of carbon emission data
	CarbonEmissionDataAvailableDateRange
}

// CarbonServiceClientQueryCarbonEmissionReportsResponse contains the response from method CarbonServiceClient.NewQueryCarbonEmissionReportsPager.
type CarbonServiceClientQueryCarbonEmissionReportsResponse struct {
	// List of carbon emission results
	CarbonEmissionDataListResult
}

// OperationsClientListResponse contains the response from method OperationsClient.NewListPager.
type OperationsClientListResponse struct {
	// A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to get the next set of results.
	OperationListResult
}

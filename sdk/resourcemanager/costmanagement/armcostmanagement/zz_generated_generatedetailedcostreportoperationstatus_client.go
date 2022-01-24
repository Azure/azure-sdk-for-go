//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcostmanagement

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// GenerateDetailedCostReportOperationStatusClient contains the methods for the GenerateDetailedCostReportOperationStatus group.
// Don't use this type directly, use NewGenerateDetailedCostReportOperationStatusClient() instead.
type GenerateDetailedCostReportOperationStatusClient struct {
	host string
	pl   runtime.Pipeline
}

// NewGenerateDetailedCostReportOperationStatusClient creates a new instance of GenerateDetailedCostReportOperationStatusClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewGenerateDetailedCostReportOperationStatusClient(credential azcore.TokenCredential, options *arm.ClientOptions) *GenerateDetailedCostReportOperationStatusClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &GenerateDetailedCostReportOperationStatusClient{
		host: string(cp.Endpoint),
		pl:   armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// Get - Get the status of the specified operation. This link is provided in the GenerateDetailedCostReport creation request
// response header.
// If the operation fails it returns an *azcore.ResponseError type.
// operationID - The target operation Id.
// scope - The scope associated with usage details operations. This includes '/subscriptions/{subscriptionId}/' for subscription
// scope, '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}' for
// Billing Account scope, '/providers/Microsoft.Billing/departments/{departmentId}' for Department scope, '/providers/Microsoft.Billing/enrollmentAccounts/{enrollmentAccountId}'
// for EnrollmentAccount
// scope. Also, Modern Commerce Account scopes are '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}' for billingAccount
// scope,
// '/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}' for billingProfile
// scope,
// 'providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}/invoiceSections/{invoiceSectionId}'
// for invoiceSection scope, and
// 'providers/Microsoft.Billing/billingAccounts/{billingAccountId}/customers/{customerId}' specific for partners.
// options - GenerateDetailedCostReportOperationStatusClientGetOptions contains the optional parameters for the GenerateDetailedCostReportOperationStatusClient.Get
// method.
func (client *GenerateDetailedCostReportOperationStatusClient) Get(ctx context.Context, operationID string, scope string, options *GenerateDetailedCostReportOperationStatusClientGetOptions) (GenerateDetailedCostReportOperationStatusClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, operationID, scope, options)
	if err != nil {
		return GenerateDetailedCostReportOperationStatusClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return GenerateDetailedCostReportOperationStatusClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return GenerateDetailedCostReportOperationStatusClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *GenerateDetailedCostReportOperationStatusClient) getCreateRequest(ctx context.Context, operationID string, scope string, options *GenerateDetailedCostReportOperationStatusClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.CostManagement/operationStatus/{operationId}"
	if operationID == "" {
		return nil, errors.New("parameter operationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationId}", url.PathEscape(operationID))
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *GenerateDetailedCostReportOperationStatusClient) getHandleResponse(resp *http.Response) (GenerateDetailedCostReportOperationStatusClientGetResponse, error) {
	result := GenerateDetailedCostReportOperationStatusClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.GenerateDetailedCostReportOperationStatuses); err != nil {
		return GenerateDetailedCostReportOperationStatusClientGetResponse{}, err
	}
	return result, nil
}

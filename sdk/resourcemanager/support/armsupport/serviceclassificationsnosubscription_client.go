//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsupport

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// ServiceClassificationsNoSubscriptionClient contains the methods for the ServiceClassificationsNoSubscription group.
// Don't use this type directly, use NewServiceClassificationsNoSubscriptionClient() instead.
type ServiceClassificationsNoSubscriptionClient struct {
	internal *arm.Client
}

// NewServiceClassificationsNoSubscriptionClient creates a new instance of ServiceClassificationsNoSubscriptionClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewServiceClassificationsNoSubscriptionClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ServiceClassificationsNoSubscriptionClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ServiceClassificationsNoSubscriptionClient{
		internal: cl,
	}
	return client, nil
}

// ClassifyServices - Classify the list of right Azure services.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01-preview
//   - serviceClassificationRequest - Input to check.
//   - options - ServiceClassificationsNoSubscriptionClientClassifyServicesOptions contains the optional parameters for the ServiceClassificationsNoSubscriptionClient.ClassifyServices
//     method.
func (client *ServiceClassificationsNoSubscriptionClient) ClassifyServices(ctx context.Context, serviceClassificationRequest ServiceClassificationRequest, options *ServiceClassificationsNoSubscriptionClientClassifyServicesOptions) (ServiceClassificationsNoSubscriptionClientClassifyServicesResponse, error) {
	var err error
	const operationName = "ServiceClassificationsNoSubscriptionClient.ClassifyServices"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.classifyServicesCreateRequest(ctx, serviceClassificationRequest, options)
	if err != nil {
		return ServiceClassificationsNoSubscriptionClientClassifyServicesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ServiceClassificationsNoSubscriptionClientClassifyServicesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ServiceClassificationsNoSubscriptionClientClassifyServicesResponse{}, err
	}
	resp, err := client.classifyServicesHandleResponse(httpResp)
	return resp, err
}

// classifyServicesCreateRequest creates the ClassifyServices request.
func (client *ServiceClassificationsNoSubscriptionClient) classifyServicesCreateRequest(ctx context.Context, serviceClassificationRequest ServiceClassificationRequest, options *ServiceClassificationsNoSubscriptionClientClassifyServicesOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Support/classifyServices"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, serviceClassificationRequest); err != nil {
		return nil, err
	}
	return req, nil
}

// classifyServicesHandleResponse handles the ClassifyServices response.
func (client *ServiceClassificationsNoSubscriptionClient) classifyServicesHandleResponse(resp *http.Response) (ServiceClassificationsNoSubscriptionClientClassifyServicesResponse, error) {
	result := ServiceClassificationsNoSubscriptionClientClassifyServicesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ServiceClassificationOutput); err != nil {
		return ServiceClassificationsNoSubscriptionClientClassifyServicesResponse{}, err
	}
	return result, nil
}

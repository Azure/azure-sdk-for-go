//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsupport

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ProblemClassificationsNoSubscriptionClient contains the methods for the ProblemClassificationsNoSubscription group.
// Don't use this type directly, use NewProblemClassificationsNoSubscriptionClient() instead.
type ProblemClassificationsNoSubscriptionClient struct {
	internal *arm.Client
}

// NewProblemClassificationsNoSubscriptionClient creates a new instance of ProblemClassificationsNoSubscriptionClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewProblemClassificationsNoSubscriptionClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ProblemClassificationsNoSubscriptionClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ProblemClassificationsNoSubscriptionClient{
		internal: cl,
	}
	return client, nil
}

// ClassifyProblems - Classify the right problem classifications (categories) available for a specific Azure service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01-preview
//   - problemServiceName - Name of the Azure service for which the problem classifications need to be retrieved.
//   - problemClassificationsClassificationInput - Input to check.
//   - options - ProblemClassificationsNoSubscriptionClientClassifyProblemsOptions contains the optional parameters for the ProblemClassificationsNoSubscriptionClient.ClassifyProblems
//     method.
func (client *ProblemClassificationsNoSubscriptionClient) ClassifyProblems(ctx context.Context, problemServiceName string, problemClassificationsClassificationInput ProblemClassificationsClassificationInput, options *ProblemClassificationsNoSubscriptionClientClassifyProblemsOptions) (ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse, error) {
	var err error
	const operationName = "ProblemClassificationsNoSubscriptionClient.ClassifyProblems"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.classifyProblemsCreateRequest(ctx, problemServiceName, problemClassificationsClassificationInput, options)
	if err != nil {
		return ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse{}, err
	}
	resp, err := client.classifyProblemsHandleResponse(httpResp)
	return resp, err
}

// classifyProblemsCreateRequest creates the ClassifyProblems request.
func (client *ProblemClassificationsNoSubscriptionClient) classifyProblemsCreateRequest(ctx context.Context, problemServiceName string, problemClassificationsClassificationInput ProblemClassificationsClassificationInput, options *ProblemClassificationsNoSubscriptionClientClassifyProblemsOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Support/services/{problemServiceName}/classifyProblems"
	if problemServiceName == "" {
		return nil, errors.New("parameter problemServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{problemServiceName}", url.PathEscape(problemServiceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, problemClassificationsClassificationInput); err != nil {
		return nil, err
	}
	return req, nil
}

// classifyProblemsHandleResponse handles the ClassifyProblems response.
func (client *ProblemClassificationsNoSubscriptionClient) classifyProblemsHandleResponse(resp *http.Response) (ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse, error) {
	result := ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ProblemClassificationsClassificationOutput); err != nil {
		return ProblemClassificationsNoSubscriptionClientClassifyProblemsResponse{}, err
	}
	return result, nil
}

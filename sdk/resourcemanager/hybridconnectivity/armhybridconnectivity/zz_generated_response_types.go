//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridconnectivity

import "net/http"

// EndpointsCreateOrUpdateResponse contains the response from method Endpoints.CreateOrUpdate.
type EndpointsCreateOrUpdateResponse struct {
	EndpointsCreateOrUpdateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// EndpointsCreateOrUpdateResult contains the result from method Endpoints.CreateOrUpdate.
type EndpointsCreateOrUpdateResult struct {
	EndpointResource
}

// EndpointsDeleteResponse contains the response from method Endpoints.Delete.
type EndpointsDeleteResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// EndpointsGetResponse contains the response from method Endpoints.Get.
type EndpointsGetResponse struct {
	EndpointsGetResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// EndpointsGetResult contains the result from method Endpoints.Get.
type EndpointsGetResult struct {
	EndpointResource
}

// EndpointsListCredentialsResponse contains the response from method Endpoints.ListCredentials.
type EndpointsListCredentialsResponse struct {
	EndpointsListCredentialsResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// EndpointsListCredentialsResult contains the result from method Endpoints.ListCredentials.
type EndpointsListCredentialsResult struct {
	EndpointAccessResource
}

// EndpointsListResponse contains the response from method Endpoints.List.
type EndpointsListResponse struct {
	EndpointsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// EndpointsListResult contains the result from method Endpoints.List.
type EndpointsListResult struct {
	EndpointsList
}

// EndpointsUpdateResponse contains the response from method Endpoints.Update.
type EndpointsUpdateResponse struct {
	EndpointsUpdateResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// EndpointsUpdateResult contains the result from method Endpoints.Update.
type EndpointsUpdateResult struct {
	EndpointResource
}

// OperationsListResponse contains the response from method Operations.List.
type OperationsListResponse struct {
	OperationsListResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// OperationsListResult contains the result from method Operations.List.
type OperationsListResult struct {
	OperationListResult
}

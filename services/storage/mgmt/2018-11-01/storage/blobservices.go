package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// BlobServicesClient is the the Azure Storage Management API.
type BlobServicesClient struct {
	BaseClient
}

// NewBlobServicesClient creates an instance of the BlobServicesClient client.
func NewBlobServicesClient(subscriptionID string) BlobServicesClient {
	return NewBlobServicesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewBlobServicesClientWithBaseURI creates an instance of the BlobServicesClient client using a custom endpoint.  Use
// this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewBlobServicesClientWithBaseURI(baseURI string, subscriptionID string) BlobServicesClient {
	return BlobServicesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// GetServiceProperties gets the properties of a storage account’s Blob service, including properties for Storage
// Analytics and CORS (Cross-Origin Resource Sharing) rules.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// accountName - the name of the storage account within the specified resource group. Storage account names
// must be between 3 and 24 characters in length and use numbers and lower-case letters only.
func (client BlobServicesClient) GetServiceProperties(ctx context.Context, resourceGroupName string, accountName string) (result BlobServiceProperties, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BlobServicesClient.GetServiceProperties")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("storage.BlobServicesClient", "GetServiceProperties", err.Error())
	}

	req, err := client.GetServicePropertiesPreparer(ctx, resourceGroupName, accountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.BlobServicesClient", "GetServiceProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetServicePropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storage.BlobServicesClient", "GetServiceProperties", resp, "Failure sending request")
		return
	}

	result, err = client.GetServicePropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.BlobServicesClient", "GetServiceProperties", resp, "Failure responding to request")
		return
	}

	return
}

// GetServicePropertiesPreparer prepares the GetServiceProperties request.
func (client BlobServicesClient) GetServicePropertiesPreparer(ctx context.Context, resourceGroupName string, accountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"BlobServicesName":  autorest.Encode("path", "default"),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/blobServices/{BlobServicesName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetServicePropertiesSender sends the GetServiceProperties request. The method will close the
// http.Response Body if it receives an error.
func (client BlobServicesClient) GetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetServicePropertiesResponder handles the response to the GetServiceProperties request. The method always
// closes the http.Response Body.
func (client BlobServicesClient) GetServicePropertiesResponder(resp *http.Response) (result BlobServiceProperties, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// SetServiceProperties sets the properties of a storage account’s Blob service, including properties for Storage
// Analytics and CORS (Cross-Origin Resource Sharing) rules.
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// accountName - the name of the storage account within the specified resource group. Storage account names
// must be between 3 and 24 characters in length and use numbers and lower-case letters only.
// parameters - the properties of a storage account’s Blob service, including properties for Storage Analytics
// and CORS (Cross-Origin Resource Sharing) rules.
func (client BlobServicesClient) SetServiceProperties(ctx context.Context, resourceGroupName string, accountName string, parameters BlobServiceProperties) (result BlobServiceProperties, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BlobServicesClient.SetServiceProperties")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: accountName,
			Constraints: []validation.Constraint{{Target: "accountName", Name: validation.MaxLength, Rule: 24, Chain: nil},
				{Target: "accountName", Name: validation.MinLength, Rule: 3, Chain: nil}}},
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.BlobServicePropertiesProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "parameters.BlobServicePropertiesProperties.DeleteRetentionPolicy", Name: validation.Null, Rule: false,
					Chain: []validation.Constraint{{Target: "parameters.BlobServicePropertiesProperties.DeleteRetentionPolicy.Days", Name: validation.Null, Rule: false,
						Chain: []validation.Constraint{{Target: "parameters.BlobServicePropertiesProperties.DeleteRetentionPolicy.Days", Name: validation.InclusiveMaximum, Rule: int64(365), Chain: nil},
							{Target: "parameters.BlobServicePropertiesProperties.DeleteRetentionPolicy.Days", Name: validation.InclusiveMinimum, Rule: int64(1), Chain: nil},
						}},
					}},
				}}}}}); err != nil {
		return result, validation.NewError("storage.BlobServicesClient", "SetServiceProperties", err.Error())
	}

	req, err := client.SetServicePropertiesPreparer(ctx, resourceGroupName, accountName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.BlobServicesClient", "SetServiceProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetServicePropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storage.BlobServicesClient", "SetServiceProperties", resp, "Failure sending request")
		return
	}

	result, err = client.SetServicePropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.BlobServicesClient", "SetServiceProperties", resp, "Failure responding to request")
		return
	}

	return
}

// SetServicePropertiesPreparer prepares the SetServiceProperties request.
func (client BlobServicesClient) SetServicePropertiesPreparer(ctx context.Context, resourceGroupName string, accountName string, parameters BlobServiceProperties) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"BlobServicesName":  autorest.Encode("path", "default"),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/blobServices/{BlobServicesName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetServicePropertiesSender sends the SetServiceProperties request. The method will close the
// http.Response Body if it receives an error.
func (client BlobServicesClient) SetServicePropertiesSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// SetServicePropertiesResponder handles the response to the SetServiceProperties request. The method always
// closes the http.Response Body.
func (client BlobServicesClient) SetServicePropertiesResponder(resp *http.Response) (result BlobServiceProperties, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

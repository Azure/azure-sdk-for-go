package iotspaces

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
    "context"
    "encoding/json"
    "github.com/Azure/go-autorest/autorest"
    "github.com/Azure/go-autorest/autorest/azure"
    "github.com/Azure/go-autorest/autorest/to"
    "github.com/Azure/go-autorest/tracing"
    "net/http"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go//services/preview/iotspaces/mgmt/2017-10-01-preview/iotspaces"

        // NameUnavailabilityReason enumerates the values for name unavailability reason.
    type NameUnavailabilityReason string

    const (
                // AlreadyExists ...
        AlreadyExists NameUnavailabilityReason = "AlreadyExists"
                // Invalid ...
        Invalid NameUnavailabilityReason = "Invalid"
            )
    // PossibleNameUnavailabilityReasonValues returns an array of possible values for the NameUnavailabilityReason const type.
    func PossibleNameUnavailabilityReasonValues() []NameUnavailabilityReason {
        return []NameUnavailabilityReason{AlreadyExists,Invalid}
    }

        // ProvisioningState enumerates the values for provisioning state.
    type ProvisioningState string

    const (
                // Canceled ...
        Canceled ProvisioningState = "Canceled"
                // Deleting ...
        Deleting ProvisioningState = "Deleting"
                // Failed ...
        Failed ProvisioningState = "Failed"
                // Provisioning ...
        Provisioning ProvisioningState = "Provisioning"
                // Succeeded ...
        Succeeded ProvisioningState = "Succeeded"
            )
    // PossibleProvisioningStateValues returns an array of possible values for the ProvisioningState const type.
    func PossibleProvisioningStateValues() []ProvisioningState {
        return []ProvisioningState{Canceled,Deleting,Failed,Provisioning,Succeeded}
    }

        // Sku enumerates the values for sku.
    type Sku string

    const (
                // F1 ...
        F1 Sku = "F1"
                // S1 ...
        S1 Sku = "S1"
                // S2 ...
        S2 Sku = "S2"
                // S3 ...
        S3 Sku = "S3"
            )
    // PossibleSkuValues returns an array of possible values for the Sku const type.
    func PossibleSkuValues() []Sku {
        return []Sku{F1,S1,S2,S3}
    }

            // CreateOrUpdateFuture an abstraction for monitoring and retrieving the results of a long-running
            // operation.
            type CreateOrUpdateFuture struct {
                azure.Future
            }
        // Result returns the result of the asynchronous operation.
        // If the operation has not completed it will return an error.
        func (future *CreateOrUpdateFuture) Result(client Client) (d Description, err error) {
        var done bool
        done, err = future.Done(client)
        if err != nil {
        err = autorest.NewErrorWithError(err, "iotspaces.CreateOrUpdateFuture", "Result", future.Response(), "Polling failure")
        return
        }
        if !done {
        err = azure.NewAsyncOpIncompleteError("iotspaces.CreateOrUpdateFuture")
        return
        }
            sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            if d.Response.Response, err = future.GetResult(sender); err == nil && d.Response.Response.StatusCode != http.StatusNoContent {
            d, err = client.CreateOrUpdateResponder(d.Response.Response)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.CreateOrUpdateFuture", "Result", d.Response.Response, "Failure responding to request")
            }
            }
            return
        }

            // DeleteFuture an abstraction for monitoring and retrieving the results of a long-running operation.
            type DeleteFuture struct {
                azure.Future
            }
        // Result returns the result of the asynchronous operation.
        // If the operation has not completed it will return an error.
        func (future *DeleteFuture) Result(client Client) (d Description, err error) {
        var done bool
        done, err = future.Done(client)
        if err != nil {
        err = autorest.NewErrorWithError(err, "iotspaces.DeleteFuture", "Result", future.Response(), "Polling failure")
        return
        }
        if !done {
        err = azure.NewAsyncOpIncompleteError("iotspaces.DeleteFuture")
        return
        }
            sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            if d.Response.Response, err = future.GetResult(sender); err == nil && d.Response.Response.StatusCode != http.StatusNoContent {
            d, err = client.DeleteResponder(d.Response.Response)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.DeleteFuture", "Result", d.Response.Response, "Failure responding to request")
            }
            }
            return
        }

            // Description the description of the IoTSpaces service.
            type Description struct {
            autorest.Response `json:"-"`
            // Properties - The common properties of a IoTSpaces service.
            Properties *Properties `json:"properties,omitempty"`
            // Sku - A valid instance SKU.
            Sku *SkuInfo `json:"sku,omitempty"`
            // ID - The resource identifier.
            ID *string `json:"id,omitempty"`
            // Name - The resource name.
            Name *string `json:"name,omitempty"`
            // Type - The resource type.
            Type *string `json:"type,omitempty"`
            // Location - The resource location.
            Location *string `json:"location,omitempty"`
            // Tags - The resource tags.
            Tags map[string]*string `json:"tags"`
            }

        // MarshalJSON is the custom marshaler for Description.
        func (d Description)MarshalJSON() ([]byte, error){
        objectMap := make(map[string]interface{})
                if(d.Properties != nil) {
                objectMap["properties"] = d.Properties
                }
                if(d.Sku != nil) {
                objectMap["sku"] = d.Sku
                }
                if(d.ID != nil) {
                objectMap["id"] = d.ID
                }
                if(d.Name != nil) {
                objectMap["name"] = d.Name
                }
                if(d.Type != nil) {
                objectMap["type"] = d.Type
                }
                if(d.Location != nil) {
                objectMap["location"] = d.Location
                }
                if(d.Tags != nil) {
                objectMap["tags"] = d.Tags
                }
                return json.Marshal(objectMap)
        }

            // DescriptionListResult a list of IoTSpaces description objects with a next link.
            type DescriptionListResult struct {
            autorest.Response `json:"-"`
            // NextLink - The link used to get the next page of IoTSpaces description objects.
            NextLink *string `json:"nextLink,omitempty"`
            // Value - A list of IoTSpaces description objects.
            Value *[]Description `json:"value,omitempty"`
            }

            // DescriptionListResultIterator provides access to a complete listing of Description values.
            type DescriptionListResultIterator struct {
                i int
                page DescriptionListResultPage
            }
        // NextWithContext advances to the next value.  If there was an error making
        // the request the iterator does not advance and the error is returned.
        func (iter * DescriptionListResultIterator) NextWithContext(ctx context.Context) (err error) {
        if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DescriptionListResultIterator.NextWithContext")
        defer func() {
        sc := -1
        if iter.Response().Response.Response != nil {
        sc = iter.Response().Response.Response.StatusCode
        }
        tracing.EndSpan(ctx, sc, err)
        }()
        }
        iter.i++
        if iter.i < len(iter. page.Values()) {
        return nil
        }
        err = iter.page.NextWithContext(ctx)
        if err != nil {
        iter. i--
        return err
        }
        iter.i = 0
        return nil
        }
        // Next advances to the next value.  If there was an error making
        // the request the iterator does not advance and the error is returned.
        // Deprecated: Use NextWithContext() instead.
        func (iter * DescriptionListResultIterator) Next() error {
        return iter.NextWithContext(context.Background())
        }
        // NotDone returns true if the enumeration should be started or is not yet complete.
        func (iter DescriptionListResultIterator) NotDone() bool {
        return iter.page.NotDone() && iter.i < len(iter. page.Values())
        }
        // Response returns the raw server response from the last page request.
        func (iter DescriptionListResultIterator) Response() DescriptionListResult {
        return iter.page.Response()
        }
        // Value returns the current value or a zero-initialized value if the
        // iterator has advanced beyond the end of the collection.
        func (iter DescriptionListResultIterator) Value() Description {
        if !iter.page.NotDone() {
        return Description{}
        }
        return iter.page.Values()[iter.i]
        }


                // IsEmpty returns true if the ListResult contains no values.
                func (dlr DescriptionListResult) IsEmpty() bool {
                return dlr.Value == nil || len(*dlr.Value) == 0
                }

                    // descriptionListResultPreparer prepares a request to retrieve the next set of results.
                    // It returns nil if no more results exist.
                    func (dlr DescriptionListResult) descriptionListResultPreparer(ctx context.Context) (*http.Request, error) {
                    if dlr.NextLink == nil || len(to.String(dlr.NextLink)) < 1 {
                    return nil, nil
                    }
                    return autorest.Prepare((&http.Request{}).WithContext(ctx),
                    autorest.AsJSON(),
                    autorest.AsGet(),
                    autorest.WithBaseURL(to.String( dlr.NextLink)));
                    }

            // DescriptionListResultPage contains a page of Description values.
            type DescriptionListResultPage struct {
                fn func(context.Context, DescriptionListResult) (DescriptionListResult, error)
                dlr DescriptionListResult
            }

        // NextWithContext advances to the next page of values.  If there was an error making
        // the request the page does not advance and the error is returned.
        func (page * DescriptionListResultPage) NextWithContext(ctx context.Context) (err error) {
        if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/DescriptionListResultPage.NextWithContext")
        defer func() {
        sc := -1
        if page.Response().Response.Response != nil {
        sc = page.Response().Response.Response.StatusCode
        }
        tracing.EndSpan(ctx, sc, err)
        }()
        }
        next, err := page.fn(ctx, page.dlr)
        if err != nil {
        return err
        }
        page.dlr = next
        return nil
        }

        // Next advances to the next page of values.  If there was an error making
        // the request the page does not advance and the error is returned.
        // Deprecated: Use NextWithContext() instead.
        func (page * DescriptionListResultPage) Next() error {
        return page.NextWithContext(context.Background())
        }
        // NotDone returns true if the page enumeration should be started or is not yet complete.
        func (page DescriptionListResultPage) NotDone() bool {
        return !page.dlr.IsEmpty()
        }
        // Response returns the raw server response from the last page request.
        func (page DescriptionListResultPage) Response() DescriptionListResult {
        return page.dlr
        }
        // Values returns the slice of values for the current page or nil if there are no values.
        func (page DescriptionListResultPage) Values() []Description {
        if page.dlr.IsEmpty() {
        return nil
        }
        return *page.dlr.Value
        }

            // ErrorDetails error details.
            type ErrorDetails struct {
            // Code - The error code.
            Code *string `json:"code,omitempty"`
            // Message - The error message.
            Message *string `json:"message,omitempty"`
            // Target - The target of the particular error.
            Target *string `json:"target,omitempty"`
            }

            // NameAvailabilityInfo the properties indicating whether a given IoTSpaces service name is available.
            type NameAvailabilityInfo struct {
            autorest.Response `json:"-"`
            // NameAvailable - The value which indicates whether the provided name is available.
            NameAvailable *bool `json:"nameAvailable,omitempty"`
            // Reason - The reason for unavailability. Possible values include: 'Invalid', 'AlreadyExists'
            Reason NameUnavailabilityReason `json:"reason,omitempty"`
            // Message - The detailed reason message.
            Message *string `json:"message,omitempty"`
            }

            // Operation ioTSpaces service REST API operation
            type Operation struct {
            // Name - Operation name: {provider}/{resource}/{read | write | action | delete}
            Name *string `json:"name,omitempty"`
            Display *OperationDisplay `json:"display,omitempty"`
            }

            // OperationDisplay the object that represents the operation.
            type OperationDisplay struct {
            // Provider - Service provider: Microsoft IoTSpaces
            Provider *string `json:"provider,omitempty"`
            // Resource - Resource Type: IoTSpaces
            Resource *string `json:"resource,omitempty"`
            // Operation - Name of the operation
            Operation *string `json:"operation,omitempty"`
            // Description - Friendly description for the operation,
            Description *string `json:"description,omitempty"`
            }

            // OperationInputs input values.
            type OperationInputs struct {
            // Name - The name of the IoTSpaces service instance to check.
            Name *string `json:"name,omitempty"`
            }

            // OperationListResult a list of IoTSpaces service operations. It contains a list of operations and a URL
            // link to get the next set of results.
            type OperationListResult struct {
            autorest.Response `json:"-"`
            // NextLink - The link used to get the next page of IoTSpaces description objects.
            NextLink *string `json:"nextLink,omitempty"`
            // Value - A list of IoT spaces operations supported by the Microsoft.IoTSpaces resource provider.
            Value *[]Operation `json:"value,omitempty"`
            }

            // OperationListResultIterator provides access to a complete listing of Operation values.
            type OperationListResultIterator struct {
                i int
                page OperationListResultPage
            }
        // NextWithContext advances to the next value.  If there was an error making
        // the request the iterator does not advance and the error is returned.
        func (iter * OperationListResultIterator) NextWithContext(ctx context.Context) (err error) {
        if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/OperationListResultIterator.NextWithContext")
        defer func() {
        sc := -1
        if iter.Response().Response.Response != nil {
        sc = iter.Response().Response.Response.StatusCode
        }
        tracing.EndSpan(ctx, sc, err)
        }()
        }
        iter.i++
        if iter.i < len(iter. page.Values()) {
        return nil
        }
        err = iter.page.NextWithContext(ctx)
        if err != nil {
        iter. i--
        return err
        }
        iter.i = 0
        return nil
        }
        // Next advances to the next value.  If there was an error making
        // the request the iterator does not advance and the error is returned.
        // Deprecated: Use NextWithContext() instead.
        func (iter * OperationListResultIterator) Next() error {
        return iter.NextWithContext(context.Background())
        }
        // NotDone returns true if the enumeration should be started or is not yet complete.
        func (iter OperationListResultIterator) NotDone() bool {
        return iter.page.NotDone() && iter.i < len(iter. page.Values())
        }
        // Response returns the raw server response from the last page request.
        func (iter OperationListResultIterator) Response() OperationListResult {
        return iter.page.Response()
        }
        // Value returns the current value or a zero-initialized value if the
        // iterator has advanced beyond the end of the collection.
        func (iter OperationListResultIterator) Value() Operation {
        if !iter.page.NotDone() {
        return Operation{}
        }
        return iter.page.Values()[iter.i]
        }


                // IsEmpty returns true if the ListResult contains no values.
                func (olr OperationListResult) IsEmpty() bool {
                return olr.Value == nil || len(*olr.Value) == 0
                }

                    // operationListResultPreparer prepares a request to retrieve the next set of results.
                    // It returns nil if no more results exist.
                    func (olr OperationListResult) operationListResultPreparer(ctx context.Context) (*http.Request, error) {
                    if olr.NextLink == nil || len(to.String(olr.NextLink)) < 1 {
                    return nil, nil
                    }
                    return autorest.Prepare((&http.Request{}).WithContext(ctx),
                    autorest.AsJSON(),
                    autorest.AsGet(),
                    autorest.WithBaseURL(to.String( olr.NextLink)));
                    }

            // OperationListResultPage contains a page of Operation values.
            type OperationListResultPage struct {
                fn func(context.Context, OperationListResult) (OperationListResult, error)
                olr OperationListResult
            }

        // NextWithContext advances to the next page of values.  If there was an error making
        // the request the page does not advance and the error is returned.
        func (page * OperationListResultPage) NextWithContext(ctx context.Context) (err error) {
        if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/OperationListResultPage.NextWithContext")
        defer func() {
        sc := -1
        if page.Response().Response.Response != nil {
        sc = page.Response().Response.Response.StatusCode
        }
        tracing.EndSpan(ctx, sc, err)
        }()
        }
        next, err := page.fn(ctx, page.olr)
        if err != nil {
        return err
        }
        page.olr = next
        return nil
        }

        // Next advances to the next page of values.  If there was an error making
        // the request the page does not advance and the error is returned.
        // Deprecated: Use NextWithContext() instead.
        func (page * OperationListResultPage) Next() error {
        return page.NextWithContext(context.Background())
        }
        // NotDone returns true if the page enumeration should be started or is not yet complete.
        func (page OperationListResultPage) NotDone() bool {
        return !page.olr.IsEmpty()
        }
        // Response returns the raw server response from the last page request.
        func (page OperationListResultPage) Response() OperationListResult {
        return page.olr
        }
        // Values returns the slice of values for the current page or nil if there are no values.
        func (page OperationListResultPage) Values() []Operation {
        if page.olr.IsEmpty() {
        return nil
        }
        return *page.olr.Value
        }

            // PatchDescription the description of the IoTSpaces service.
            type PatchDescription struct {
            // Tags - Instance tags
            Tags map[string]*string `json:"tags"`
            // Properties - The common properties of an IoTSpaces service.
            Properties *Properties `json:"properties,omitempty"`
            }

        // MarshalJSON is the custom marshaler for PatchDescription.
        func (pd PatchDescription)MarshalJSON() ([]byte, error){
        objectMap := make(map[string]interface{})
                if(pd.Tags != nil) {
                objectMap["tags"] = pd.Tags
                }
                if(pd.Properties != nil) {
                objectMap["properties"] = pd.Properties
                }
                return json.Marshal(objectMap)
        }

            // Properties the properties of an IoTSpaces instance.
            type Properties struct {
            // ProvisioningState - The provisioning state. Possible values include: 'Provisioning', 'Deleting', 'Succeeded', 'Failed', 'Canceled'
            ProvisioningState ProvisioningState `json:"provisioningState,omitempty"`
            // ManagementAPIURL - The management Api endpoint.
            ManagementAPIURL *string `json:"managementApiUrl,omitempty"`
            // WebPortalURL - The management UI endpoint.
            WebPortalURL *string `json:"webPortalUrl,omitempty"`
            // StorageContainer - The properties of the designated storage container.
            StorageContainer *StorageContainerProperties `json:"storageContainer,omitempty"`
            }

            // Resource the common properties of an IoTSpaces service.
            type Resource struct {
            // ID - The resource identifier.
            ID *string `json:"id,omitempty"`
            // Name - The resource name.
            Name *string `json:"name,omitempty"`
            // Type - The resource type.
            Type *string `json:"type,omitempty"`
            // Location - The resource location.
            Location *string `json:"location,omitempty"`
            // Tags - The resource tags.
            Tags map[string]*string `json:"tags"`
            }

        // MarshalJSON is the custom marshaler for Resource.
        func (r Resource)MarshalJSON() ([]byte, error){
        objectMap := make(map[string]interface{})
                if(r.ID != nil) {
                objectMap["id"] = r.ID
                }
                if(r.Name != nil) {
                objectMap["name"] = r.Name
                }
                if(r.Type != nil) {
                objectMap["type"] = r.Type
                }
                if(r.Location != nil) {
                objectMap["location"] = r.Location
                }
                if(r.Tags != nil) {
                objectMap["tags"] = r.Tags
                }
                return json.Marshal(objectMap)
        }

            // SkuInfo information about the SKU of the IoTSpaces instance.
            type SkuInfo struct {
            // Name - The name of the SKU. Possible values include: 'F1', 'S1', 'S2', 'S3'
            Name Sku `json:"name,omitempty"`
            }

            // StorageContainerProperties the properties of the Azure Storage Container for file archive.
            type StorageContainerProperties struct {
            // ConnectionString - The connection string of the storage account.
            ConnectionString *string `json:"connectionString,omitempty"`
            // SubscriptionID - The subscription identifier of the storage account.
            SubscriptionID *string `json:"subscriptionId,omitempty"`
            // ResourceGroup - The name of the resource group of the storage account.
            ResourceGroup *string `json:"resourceGroup,omitempty"`
            // ContainerName - The name of storage container in the storage account.
            ContainerName *string `json:"containerName,omitempty"`
            }

            // UpdateFuture an abstraction for monitoring and retrieving the results of a long-running operation.
            type UpdateFuture struct {
                azure.Future
            }
        // Result returns the result of the asynchronous operation.
        // If the operation has not completed it will return an error.
        func (future *UpdateFuture) Result(client Client) (d Description, err error) {
        var done bool
        done, err = future.Done(client)
        if err != nil {
        err = autorest.NewErrorWithError(err, "iotspaces.UpdateFuture", "Result", future.Response(), "Polling failure")
        return
        }
        if !done {
        err = azure.NewAsyncOpIncompleteError("iotspaces.UpdateFuture")
        return
        }
            sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
            if d.Response.Response, err = future.GetResult(sender); err == nil && d.Response.Response.StatusCode != http.StatusNoContent {
            d, err = client.UpdateResponder(d.Response.Response)
            if err != nil {
            err = autorest.NewErrorWithError(err, "iotspaces.UpdateFuture", "Result", d.Response.Response, "Failure responding to request")
            }
            }
            return
        }


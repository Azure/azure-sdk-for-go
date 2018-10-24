package managementgroups

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
    "github.com/Azure/go-autorest/autorest/date"
    "github.com/Azure/go-autorest/autorest/to"
    "github.com/Azure/go-autorest/tracing"
    "github.com/satori/go.uuid"
    "net/http"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go//services/preview/resources/mgmt/2017-08-31-preview/managementgroups"

        // ChildType enumerates the values for child type.
    type ChildType string

    const (
                // Account ...
        Account ChildType = "Account"
                // Department ...
        Department ChildType = "Department"
                // Enrollment ...
        Enrollment ChildType = "Enrollment"
                // Subscription ...
        Subscription ChildType = "Subscription"
            )
    // PossibleChildTypeValues returns an array of possible values for the ChildType const type.
    func PossibleChildTypeValues() []ChildType {
        return []ChildType{Account,Department,Enrollment,Subscription}
    }

        // ChildType1 enumerates the values for child type 1.
    type ChildType1 string

    const (
                // ChildType1Account ...
        ChildType1Account ChildType1 = "Account"
                // ChildType1Department ...
        ChildType1Department ChildType1 = "Department"
                // ChildType1Enrollment ...
        ChildType1Enrollment ChildType1 = "Enrollment"
                // ChildType1Subscription ...
        ChildType1Subscription ChildType1 = "Subscription"
            )
    // PossibleChildType1Values returns an array of possible values for the ChildType1 const type.
    func PossibleChildType1Values() []ChildType1 {
        return []ChildType1{ChildType1Account,ChildType1Department,ChildType1Enrollment,ChildType1Subscription}
    }

        // ManagementGroupType enumerates the values for management group type.
    type ManagementGroupType string

    const (
                // ManagementGroupTypeAccount ...
        ManagementGroupTypeAccount ManagementGroupType = "Account"
                // ManagementGroupTypeDepartment ...
        ManagementGroupTypeDepartment ManagementGroupType = "Department"
                // ManagementGroupTypeEnrollment ...
        ManagementGroupTypeEnrollment ManagementGroupType = "Enrollment"
                // ManagementGroupTypeSubscription ...
        ManagementGroupTypeSubscription ManagementGroupType = "Subscription"
            )
    // PossibleManagementGroupTypeValues returns an array of possible values for the ManagementGroupType const type.
    func PossibleManagementGroupTypeValues() []ManagementGroupType {
        return []ManagementGroupType{ManagementGroupTypeAccount,ManagementGroupTypeDepartment,ManagementGroupTypeEnrollment,ManagementGroupTypeSubscription}
    }

            // ChildInfo the unique identifier (ID) of a management group.
            type ChildInfo struct {
            // ChildType - Possible values include: 'Enrollment', 'Department', 'Account', 'Subscription'
            ChildType ChildType `json:"childType,omitempty"`
            // ChildID - The ID of the child resource (management group or subscription). E.g. /providers/Microsoft.Management/managementGroups/40000000-0000-0000-0000-000000000000
            ChildID *string `json:"childId,omitempty"`
            // DisplayName - The friendly name of the child resource.
            DisplayName *string `json:"displayName,omitempty"`
            // TenantID - (Optional) The AAD Tenant ID associated with the child resource.
            TenantID *uuid.UUID `json:"tenantId,omitempty"`
            }

            // DetailsProperties the details properties of a management group.
            type DetailsProperties struct {
            // Version - The version number of the object.
            Version *float64 `json:"version,omitempty"`
            // UpdatedTime - The date and time when this object was last updated.
            UpdatedTime *date.Time `json:"updatedTime,omitempty"`
            // UpdatedBy - The identity of the principal or process that updated the object.
            UpdatedBy *string `json:"updatedBy,omitempty"`
            Parent *ParentGroupInfo `json:"parent,omitempty"`
            // ManagementGroupType - Possible values include: 'ManagementGroupTypeEnrollment', 'ManagementGroupTypeDepartment', 'ManagementGroupTypeAccount', 'ManagementGroupTypeSubscription'
            ManagementGroupType ManagementGroupType `json:"managementGroupType,omitempty"`
            }

            // ErrorDetails the details of the error.
            type ErrorDetails struct {
            // Code - One of a server-defined set of error codes.
            Code *string `json:"code,omitempty"`
            // Message - A human-readable representation of the error.
            Message *string `json:"message,omitempty"`
            // Target - (Optional) The target of the error.
            Target *string `json:"target,omitempty"`
            }

            // ErrorResponse the error object.
            type ErrorResponse struct {
            Error *ErrorDetails `json:"error,omitempty"`
            }

            // Info the management group.
            type Info struct {
            // ID - The ID of the management group. E.g. /providers/Microsoft.Management/managementGroups/20000000-0000-0000-0000-000000000000
            ID *string `json:"id,omitempty"`
            // Type - The type of the resource. E.g. /providers/Microsoft.Management/managementGroups
            Type *string `json:"type,omitempty"`
            // Name - The name of the management group. E.g. 20000000-0000-0000-0000-000000000000
            Name *uuid.UUID `json:"name,omitempty"`
            *InfoProperties `json:"properties,omitempty"`
            }

        // MarshalJSON is the custom marshaler for Info.
        func (i Info)MarshalJSON() ([]byte, error){
        objectMap := make(map[string]interface{})
                if(i.ID != nil) {
                objectMap["id"] = i.ID
                }
                if(i.Type != nil) {
                objectMap["type"] = i.Type
                }
                if(i.Name != nil) {
                objectMap["name"] = i.Name
                }
                if(i.InfoProperties != nil) {
                objectMap["properties"] = i.InfoProperties
                }
                return json.Marshal(objectMap)
        }
        // UnmarshalJSON is the custom unmarshaler for Info struct.
        func (i *Info) UnmarshalJSON(body []byte) error {
        var m map[string]*json.RawMessage
        err := json.Unmarshal(body, &m)
        if err != nil {
        return err
        }
        for k, v := range  m {
        switch k {
                case "id":
    if v != nil {
        var ID string
        err = json.Unmarshal(*v, &ID)
    if err != nil {
    return err
    }
        i.ID = &ID
    }
                case "type":
    if v != nil {
        var typeVar string
        err = json.Unmarshal(*v, &typeVar)
    if err != nil {
    return err
    }
        i.Type = &typeVar
    }
                case "name":
    if v != nil {
        var name uuid.UUID
        err = json.Unmarshal(*v, &name)
    if err != nil {
    return err
    }
        i.Name = &name
    }
                case "properties":
    if v != nil {
        var infoProperties InfoProperties
        err = json.Unmarshal(*v, &infoProperties)
    if err != nil {
    return err
    }
        i.InfoProperties = &infoProperties
    }
            }
        }

        return nil
        }

            // InfoProperties the generic properties of a management group.
            type InfoProperties struct {
            // TenantID - The AAD Tenant ID associated with the management group. E.g. 10000000-0000-0000-0000-000000000000
            TenantID *uuid.UUID `json:"tenantId,omitempty"`
            // DisplayName - The friendly name of the management group.
            DisplayName *string `json:"displayName,omitempty"`
            }

            // ListResult the result of listing management groups.
            type ListResult struct {
            autorest.Response `json:"-"`
            // Value - The list of management groups.
            Value *[]Info `json:"value,omitempty"`
            // NextLink - The URL to the next page of results.
            NextLink *string `json:"nextLink,omitempty"`
            }

            // ListResultIterator provides access to a complete listing of Info values.
            type ListResultIterator struct {
                i int
                page ListResultPage
            }
        // NextWithContext advances to the next value.  If there was an error making
        // the request the iterator does not advance and the error is returned.
        func (iter * ListResultIterator) NextWithContext(ctx context.Context) (err error) {
        if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ListResultIterator.NextWithContext")
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
        func (iter * ListResultIterator) Next() error {
        return iter.NextWithContext(context.Background())
        }
        // NotDone returns true if the enumeration should be started or is not yet complete.
        func (iter ListResultIterator) NotDone() bool {
        return iter.page.NotDone() && iter.i < len(iter. page.Values())
        }
        // Response returns the raw server response from the last page request.
        func (iter ListResultIterator) Response() ListResult {
        return iter.page.Response()
        }
        // Value returns the current value or a zero-initialized value if the
        // iterator has advanced beyond the end of the collection.
        func (iter ListResultIterator) Value() Info {
        if !iter.page.NotDone() {
        return Info{}
        }
        return iter.page.Values()[iter.i]
        }


                // IsEmpty returns true if the ListResult contains no values.
                func (lr ListResult) IsEmpty() bool {
                return lr.Value == nil || len(*lr.Value) == 0
                }

                    // listResultPreparer prepares a request to retrieve the next set of results.
                    // It returns nil if no more results exist.
                    func (lr ListResult) listResultPreparer(ctx context.Context) (*http.Request, error) {
                    if lr.NextLink == nil || len(to.String(lr.NextLink)) < 1 {
                    return nil, nil
                    }
                    return autorest.Prepare((&http.Request{}).WithContext(ctx),
                    autorest.AsJSON(),
                    autorest.AsGet(),
                    autorest.WithBaseURL(to.String( lr.NextLink)));
                    }

            // ListResultPage contains a page of Info values.
            type ListResultPage struct {
                fn func(context.Context, ListResult) (ListResult, error)
                lr ListResult
            }

        // NextWithContext advances to the next page of values.  If there was an error making
        // the request the page does not advance and the error is returned.
        func (page * ListResultPage) NextWithContext(ctx context.Context) (err error) {
        if tracing.IsEnabled() {
        ctx = tracing.StartSpan(ctx, fqdn + "/ListResultPage.NextWithContext")
        defer func() {
        sc := -1
        if page.Response().Response.Response != nil {
        sc = page.Response().Response.Response.StatusCode
        }
        tracing.EndSpan(ctx, sc, err)
        }()
        }
        next, err := page.fn(ctx, page.lr)
        if err != nil {
        return err
        }
        page.lr = next
        return nil
        }

        // Next advances to the next page of values.  If there was an error making
        // the request the page does not advance and the error is returned.
        // Deprecated: Use NextWithContext() instead.
        func (page * ListResultPage) Next() error {
        return page.NextWithContext(context.Background())
        }
        // NotDone returns true if the page enumeration should be started or is not yet complete.
        func (page ListResultPage) NotDone() bool {
        return !page.lr.IsEmpty()
        }
        // Response returns the raw server response from the last page request.
        func (page ListResultPage) Response() ListResult {
        return page.lr
        }
        // Values returns the slice of values for the current page or nil if there are no values.
        func (page ListResultPage) Values() []Info {
        if page.lr.IsEmpty() {
        return nil
        }
        return *page.lr.Value
        }

            // ManagementGroup the management group details.
            type ManagementGroup struct {
            // ID - The ID of the management group. E.g. /providers/Microsoft.Management/managementGroups/20000000-0000-0000-0000-000000000000
            ID *string `json:"id,omitempty"`
            // Type - The type of the resource. E.g. /providers/Microsoft.Management/managementGroups
            Type *string `json:"type,omitempty"`
            // Name - The name of the management group. E.g. 20000000-0000-0000-0000-000000000000
            Name *uuid.UUID `json:"name,omitempty"`
            *Properties `json:"properties,omitempty"`
            }

        // MarshalJSON is the custom marshaler for ManagementGroup.
        func (mg ManagementGroup)MarshalJSON() ([]byte, error){
        objectMap := make(map[string]interface{})
                if(mg.ID != nil) {
                objectMap["id"] = mg.ID
                }
                if(mg.Type != nil) {
                objectMap["type"] = mg.Type
                }
                if(mg.Name != nil) {
                objectMap["name"] = mg.Name
                }
                if(mg.Properties != nil) {
                objectMap["properties"] = mg.Properties
                }
                return json.Marshal(objectMap)
        }
        // UnmarshalJSON is the custom unmarshaler for ManagementGroup struct.
        func (mg *ManagementGroup) UnmarshalJSON(body []byte) error {
        var m map[string]*json.RawMessage
        err := json.Unmarshal(body, &m)
        if err != nil {
        return err
        }
        for k, v := range  m {
        switch k {
                case "id":
    if v != nil {
        var ID string
        err = json.Unmarshal(*v, &ID)
    if err != nil {
    return err
    }
        mg.ID = &ID
    }
                case "type":
    if v != nil {
        var typeVar string
        err = json.Unmarshal(*v, &typeVar)
    if err != nil {
    return err
    }
        mg.Type = &typeVar
    }
                case "name":
    if v != nil {
        var name uuid.UUID
        err = json.Unmarshal(*v, &name)
    if err != nil {
    return err
    }
        mg.Name = &name
    }
                case "properties":
    if v != nil {
        var properties Properties
        err = json.Unmarshal(*v, &properties)
    if err != nil {
    return err
    }
        mg.Properties = &properties
    }
            }
        }

        return nil
        }

            // Operation a Management REST API operation.
            type Operation struct {
            // Name - Operation name: {provider}/{resource}/{operation}.
            Name *string `json:"name,omitempty"`
            // Display - The object that represents the operation.
            Display *OperationDisplay `json:"display,omitempty"`
            }

            // OperationDisplay the object that represents the operation.
            type OperationDisplay struct {
            // Provider - Service provider: Microsoft.Management.
            Provider *string `json:"provider,omitempty"`
            // Resource - Resource on which the operation is performed: Invoice, etc.
            Resource *string `json:"resource,omitempty"`
            // Operation - Operation type: Read, write, delete, etc.
            Operation *string `json:"operation,omitempty"`
            }

            // OperationListResult result listing  operations. It contains a list of operations and a URL link to get
            // the next set of results.
            type OperationListResult struct {
            autorest.Response `json:"-"`
            // Value - List of management operations supported by the Microsoft.Management resource provider.
            Value *[]Operation `json:"value,omitempty"`
            // NextLink - URL to get the next set of operation list results if there are any.
            NextLink *string `json:"nextLink,omitempty"`
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

            // ParentGroupInfo (Optional) The ID of the parent management group.
            type ParentGroupInfo struct {
            // ParentID - The ID of the parent management group. E.g. /providers/Microsoft.Management/managementGroups/30000000-0000-0000-0000-000000000000
            ParentID *string `json:"parentId,omitempty"`
            // DisplayName - The friendly name of the management group.
            DisplayName *string `json:"displayName,omitempty"`
            }

            // Properties the generic properties of a management group.
            type Properties struct {
            // TenantID - The AAD Tenant ID associated with the management group. E.g. 10000000-0000-0000-0000-000000000000
            TenantID *uuid.UUID `json:"tenantId,omitempty"`
            // DisplayName - The friendly name of the management group.
            DisplayName *string `json:"displayName,omitempty"`
            Details *DetailsProperties `json:"details,omitempty"`
            }

            // PropertiesWithChildren the generic properties of a management group.
            type PropertiesWithChildren struct {
            // TenantID - The AAD Tenant ID associated with the management group. E.g. 10000000-0000-0000-0000-000000000000
            TenantID *uuid.UUID `json:"tenantId,omitempty"`
            // DisplayName - The friendly name of the management group.
            DisplayName *string `json:"displayName,omitempty"`
            Details *DetailsProperties `json:"details,omitempty"`
            // Children - The list of children.
            Children *[]ChildInfo `json:"children,omitempty"`
            }

            // PropertiesWithHierarchy the generic properties of a management group.
            type PropertiesWithHierarchy struct {
            // TenantID - The AAD Tenant ID associated with the management group. E.g. 10000000-0000-0000-0000-000000000000
            TenantID *uuid.UUID `json:"tenantId,omitempty"`
            // DisplayName - The friendly name of the management group.
            DisplayName *string `json:"displayName,omitempty"`
            Details *DetailsProperties `json:"details,omitempty"`
            // Children - The list of children.
            Children *[]RecursiveChildInfo `json:"children,omitempty"`
            }

            // RecursiveChildInfo the unique identifier (ID) of a management group.
            type RecursiveChildInfo struct {
            // ChildType - Possible values include: 'ChildType1Enrollment', 'ChildType1Department', 'ChildType1Account', 'ChildType1Subscription'
            ChildType ChildType1 `json:"childType,omitempty"`
            // ChildID - The ID of the child resource (management group or subscription). E.g. /providers/Microsoft.Management/managementGroups/40000000-0000-0000-0000-000000000000
            ChildID *string `json:"childId,omitempty"`
            // DisplayName - The friendly name of the child resource.
            DisplayName *string `json:"displayName,omitempty"`
            // TenantID - (Optional) The AAD Tenant ID associated with the child resource.
            TenantID *uuid.UUID `json:"tenantId,omitempty"`
            // Children - The list of children.
            Children *[]RecursiveChildInfo `json:"children,omitempty"`
            }

            // WithChildren the management group details.
            type WithChildren struct {
            // ID - The ID of the management group. E.g. /providers/Microsoft.Management/managementGroups/20000000-0000-0000-0000-000000000000
            ID *string `json:"id,omitempty"`
            // Type - The type of the resource. E.g. /providers/Microsoft.Management/managementGroups
            Type *string `json:"type,omitempty"`
            // Name - The name of the management group. E.g. 20000000-0000-0000-0000-000000000000
            Name *uuid.UUID `json:"name,omitempty"`
            *PropertiesWithChildren `json:"properties,omitempty"`
            }

        // MarshalJSON is the custom marshaler for WithChildren.
        func (wc WithChildren)MarshalJSON() ([]byte, error){
        objectMap := make(map[string]interface{})
                if(wc.ID != nil) {
                objectMap["id"] = wc.ID
                }
                if(wc.Type != nil) {
                objectMap["type"] = wc.Type
                }
                if(wc.Name != nil) {
                objectMap["name"] = wc.Name
                }
                if(wc.PropertiesWithChildren != nil) {
                objectMap["properties"] = wc.PropertiesWithChildren
                }
                return json.Marshal(objectMap)
        }
        // UnmarshalJSON is the custom unmarshaler for WithChildren struct.
        func (wc *WithChildren) UnmarshalJSON(body []byte) error {
        var m map[string]*json.RawMessage
        err := json.Unmarshal(body, &m)
        if err != nil {
        return err
        }
        for k, v := range  m {
        switch k {
                case "id":
    if v != nil {
        var ID string
        err = json.Unmarshal(*v, &ID)
    if err != nil {
    return err
    }
        wc.ID = &ID
    }
                case "type":
    if v != nil {
        var typeVar string
        err = json.Unmarshal(*v, &typeVar)
    if err != nil {
    return err
    }
        wc.Type = &typeVar
    }
                case "name":
    if v != nil {
        var name uuid.UUID
        err = json.Unmarshal(*v, &name)
    if err != nil {
    return err
    }
        wc.Name = &name
    }
                case "properties":
    if v != nil {
        var propertiesWithChildren PropertiesWithChildren
        err = json.Unmarshal(*v, &propertiesWithChildren)
    if err != nil {
    return err
    }
        wc.PropertiesWithChildren = &propertiesWithChildren
    }
            }
        }

        return nil
        }

            // WithHierarchy the management group details.
            type WithHierarchy struct {
            autorest.Response `json:"-"`
            // ID - The ID of the management group. E.g. /providers/Microsoft.Management/managementGroups/20000000-0000-0000-0000-000000000000
            ID *string `json:"id,omitempty"`
            // Type - The type of the resource. E.g. /providers/Microsoft.Management/managementGroups
            Type *string `json:"type,omitempty"`
            // Name - The name of the management group. E.g. 20000000-0000-0000-0000-000000000000
            Name *uuid.UUID `json:"name,omitempty"`
            *PropertiesWithHierarchy `json:"properties,omitempty"`
            }

        // MarshalJSON is the custom marshaler for WithHierarchy.
        func (wh WithHierarchy)MarshalJSON() ([]byte, error){
        objectMap := make(map[string]interface{})
                if(wh.ID != nil) {
                objectMap["id"] = wh.ID
                }
                if(wh.Type != nil) {
                objectMap["type"] = wh.Type
                }
                if(wh.Name != nil) {
                objectMap["name"] = wh.Name
                }
                if(wh.PropertiesWithHierarchy != nil) {
                objectMap["properties"] = wh.PropertiesWithHierarchy
                }
                return json.Marshal(objectMap)
        }
        // UnmarshalJSON is the custom unmarshaler for WithHierarchy struct.
        func (wh *WithHierarchy) UnmarshalJSON(body []byte) error {
        var m map[string]*json.RawMessage
        err := json.Unmarshal(body, &m)
        if err != nil {
        return err
        }
        for k, v := range  m {
        switch k {
                case "id":
    if v != nil {
        var ID string
        err = json.Unmarshal(*v, &ID)
    if err != nil {
    return err
    }
        wh.ID = &ID
    }
                case "type":
    if v != nil {
        var typeVar string
        err = json.Unmarshal(*v, &typeVar)
    if err != nil {
    return err
    }
        wh.Type = &typeVar
    }
                case "name":
    if v != nil {
        var name uuid.UUID
        err = json.Unmarshal(*v, &name)
    if err != nil {
    return err
    }
        wh.Name = &name
    }
                case "properties":
    if v != nil {
        var propertiesWithHierarchy PropertiesWithHierarchy
        err = json.Unmarshal(*v, &propertiesWithHierarchy)
    if err != nil {
    return err
    }
        wh.PropertiesWithHierarchy = &propertiesWithHierarchy
    }
            }
        }

        return nil
        }


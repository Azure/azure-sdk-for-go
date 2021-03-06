package managementgroupsapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2017-08-31-preview/managementgroups"
)

// ClientAPI contains the set of methods on the Client type.
type ClientAPI interface {
	Get(ctx context.Context, expand string, recurse *bool) (result managementgroups.WithHierarchy, err error)
	List(ctx context.Context, skiptoken string) (result managementgroups.ListResultPage, err error)
	ListComplete(ctx context.Context, skiptoken string) (result managementgroups.ListResultIterator, err error)
}

var _ ClientAPI = (*managementgroups.Client)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result managementgroups.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result managementgroups.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*managementgroups.OperationsClient)(nil)

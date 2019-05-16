package hanaonazureapi

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
	"github.com/Azure/azure-sdk-for-go/services/preview/hanaonazure/mgmt/2017-11-03-preview/hanaonazure"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result hanaonazure.OperationList, err error)
}

var _ OperationsClientAPI = (*hanaonazure.OperationsClient)(nil)

// HanaInstancesClientAPI contains the set of methods on the HanaInstancesClient type.
type HanaInstancesClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, hanaInstanceName string) (result hanaonazure.HanaInstancesCreateFuture, err error)
	EnableMonitoring(ctx context.Context, resourceGroupName string, hanaInstanceName string, monitoringParameter hanaonazure.MonitoringDetails) (result hanaonazure.HanaInstancesEnableMonitoringFuture, err error)
	Get(ctx context.Context, resourceGroupName string, hanaInstanceName string) (result hanaonazure.HanaInstance, err error)
	List(ctx context.Context) (result hanaonazure.HanaInstancesListResultPage, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result hanaonazure.HanaInstancesListResultPage, err error)
	Restart(ctx context.Context, resourceGroupName string, hanaInstanceName string) (result hanaonazure.HanaInstancesRestartFuture, err error)
	Update(ctx context.Context, resourceGroupName string, hanaInstanceName string, tagsParameter hanaonazure.Tags) (result hanaonazure.HanaInstance, err error)
}

var _ HanaInstancesClientAPI = (*hanaonazure.HanaInstancesClient)(nil)

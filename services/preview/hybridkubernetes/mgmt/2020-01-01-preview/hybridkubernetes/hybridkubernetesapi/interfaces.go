package hybridkubernetesapi

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
	"github.com/Azure/azure-sdk-for-go/services/preview/hybridkubernetes/mgmt/2020-01-01-preview/hybridkubernetes"
)

// ConnectedClusterClientAPI contains the set of methods on the ConnectedClusterClient type.
type ConnectedClusterClientAPI interface {
	Create(ctx context.Context, resourceGroupName string, clusterName string, connectedCluster hybridkubernetes.ConnectedCluster) (result hybridkubernetes.ConnectedClusterCreateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, clusterName string) (result hybridkubernetes.ConnectedClusterDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, clusterName string) (result hybridkubernetes.ConnectedCluster, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result hybridkubernetes.ConnectedClusterListPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result hybridkubernetes.ConnectedClusterListIterator, err error)
	ListBySubscription(ctx context.Context) (result hybridkubernetes.ConnectedClusterListPage, err error)
	ListBySubscriptionComplete(ctx context.Context) (result hybridkubernetes.ConnectedClusterListIterator, err error)
	ListClusterUserCredentials(ctx context.Context, resourceGroupName string, clusterName string, clientAuthenticationDetails *hybridkubernetes.AuthenticationDetails) (result hybridkubernetes.CredentialResults, err error)
	Update(ctx context.Context, resourceGroupName string, clusterName string, connectedClusterPatch hybridkubernetes.ConnectedClusterPatch) (result hybridkubernetes.ConnectedCluster, err error)
}

var _ ConnectedClusterClientAPI = (*hybridkubernetes.ConnectedClusterClient)(nil)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	Get(ctx context.Context) (result hybridkubernetes.OperationListPage, err error)
	GetComplete(ctx context.Context) (result hybridkubernetes.OperationListIterator, err error)
}

var _ OperationsClientAPI = (*hybridkubernetes.OperationsClient)(nil)

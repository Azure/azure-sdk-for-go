package avsapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
	"github.com/Azure/go-autorest/autorest"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result avs.OperationListPage, err error)
	ListComplete(ctx context.Context) (result avs.OperationListIterator, err error)
}

var _ OperationsClientAPI = (*avs.OperationsClient)(nil)

// LocationsClientAPI contains the set of methods on the LocationsClient type.
type LocationsClientAPI interface {
	CheckQuotaAvailability(ctx context.Context, location string) (result avs.Quota, err error)
	CheckTrialAvailability(ctx context.Context, location string) (result avs.Trial, err error)
}

var _ LocationsClientAPI = (*avs.LocationsClient)(nil)

// PrivateCloudsClientAPI contains the set of methods on the PrivateCloudsClient type.
type PrivateCloudsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, privateCloudName string, privateCloud avs.PrivateCloud) (result avs.PrivateCloudsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.PrivateCloudsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.PrivateCloud, err error)
	List(ctx context.Context, resourceGroupName string) (result avs.PrivateCloudListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string) (result avs.PrivateCloudListIterator, err error)
	ListAdminCredentials(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.AdminCredentials, err error)
	ListInSubscription(ctx context.Context) (result avs.PrivateCloudListPage, err error)
	ListInSubscriptionComplete(ctx context.Context) (result avs.PrivateCloudListIterator, err error)
	Update(ctx context.Context, resourceGroupName string, privateCloudName string, privateCloudUpdate avs.PrivateCloudUpdate) (result avs.PrivateCloudsUpdateFuture, err error)
}

var _ PrivateCloudsClientAPI = (*avs.PrivateCloudsClient)(nil)

// ClustersClientAPI contains the set of methods on the ClustersClient type.
type ClustersClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, privateCloudName string, clusterName string, cluster avs.Cluster) (result avs.ClustersCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, privateCloudName string, clusterName string) (result avs.ClustersDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, privateCloudName string, clusterName string) (result avs.Cluster, err error)
	List(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.ClusterListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.ClusterListIterator, err error)
	Update(ctx context.Context, resourceGroupName string, privateCloudName string, clusterName string, clusterUpdate avs.ClusterUpdate) (result avs.ClustersUpdateFuture, err error)
}

var _ ClustersClientAPI = (*avs.ClustersClient)(nil)

// HcxEnterpriseSitesClientAPI contains the set of methods on the HcxEnterpriseSitesClient type.
type HcxEnterpriseSitesClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, privateCloudName string, hcxEnterpriseSiteName string, hcxEnterpriseSite avs.HcxEnterpriseSite) (result avs.HcxEnterpriseSite, err error)
	Delete(ctx context.Context, resourceGroupName string, privateCloudName string, hcxEnterpriseSiteName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, privateCloudName string, hcxEnterpriseSiteName string) (result avs.HcxEnterpriseSite, err error)
	List(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.HcxEnterpriseSiteListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.HcxEnterpriseSiteListIterator, err error)
}

var _ HcxEnterpriseSitesClientAPI = (*avs.HcxEnterpriseSitesClient)(nil)

// AuthorizationsClientAPI contains the set of methods on the AuthorizationsClient type.
type AuthorizationsClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, privateCloudName string, authorizationName string, authorization avs.ExpressRouteAuthorization) (result avs.AuthorizationsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, privateCloudName string, authorizationName string) (result avs.AuthorizationsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, privateCloudName string, authorizationName string) (result avs.ExpressRouteAuthorization, err error)
	List(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.ExpressRouteAuthorizationListPage, err error)
	ListComplete(ctx context.Context, resourceGroupName string, privateCloudName string) (result avs.ExpressRouteAuthorizationListIterator, err error)
}

var _ AuthorizationsClientAPI = (*avs.AuthorizationsClient)(nil)

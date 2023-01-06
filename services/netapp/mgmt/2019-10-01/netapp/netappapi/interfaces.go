// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package netappapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-10-01/netapp"
	"github.com/Azure/go-autorest/autorest"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result netapp.OperationListResult, err error)
}

var _ OperationsClientAPI = (*netapp.OperationsClient)(nil)

// ResourceClientAPI contains the set of methods on the ResourceClient type.
type ResourceClientAPI interface {
	CheckFilePathAvailability(ctx context.Context, body netapp.ResourceNameAvailabilityRequest, location string) (result netapp.ResourceNameAvailability, err error)
	CheckNameAvailability(ctx context.Context, body netapp.ResourceNameAvailabilityRequest, location string) (result netapp.ResourceNameAvailability, err error)
}

var _ ResourceClientAPI = (*netapp.ResourceClient)(nil)

// AccountsClientAPI contains the set of methods on the AccountsClient type.
type AccountsClientAPI interface {
	CreateOrUpdate(ctx context.Context, body netapp.Account, resourceGroupName string, accountName string) (result netapp.AccountsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, accountName string) (result netapp.AccountsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, accountName string) (result netapp.Account, err error)
	List(ctx context.Context, resourceGroupName string) (result netapp.AccountList, err error)
	Update(ctx context.Context, body netapp.AccountPatch, resourceGroupName string, accountName string) (result netapp.Account, err error)
}

var _ AccountsClientAPI = (*netapp.AccountsClient)(nil)

// PoolsClientAPI contains the set of methods on the PoolsClient type.
type PoolsClientAPI interface {
	CreateOrUpdate(ctx context.Context, body netapp.CapacityPool, resourceGroupName string, accountName string, poolName string) (result netapp.PoolsCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, accountName string, poolName string) (result netapp.PoolsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, accountName string, poolName string) (result netapp.CapacityPool, err error)
	List(ctx context.Context, resourceGroupName string, accountName string) (result netapp.CapacityPoolList, err error)
	Update(ctx context.Context, body netapp.CapacityPoolPatch, resourceGroupName string, accountName string, poolName string) (result netapp.CapacityPool, err error)
}

var _ PoolsClientAPI = (*netapp.PoolsClient)(nil)

// VolumesClientAPI contains the set of methods on the VolumesClient type.
type VolumesClientAPI interface {
	AuthorizeReplication(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, body netapp.AuthorizeRequest) (result autorest.Response, err error)
	BreakReplication(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result autorest.Response, err error)
	CreateOrUpdate(ctx context.Context, body netapp.Volume, resourceGroupName string, accountName string, poolName string, volumeName string) (result netapp.VolumesCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result netapp.VolumesDeleteFuture, err error)
	DeleteReplication(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result autorest.Response, err error)
	Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result netapp.Volume, err error)
	List(ctx context.Context, resourceGroupName string, accountName string, poolName string) (result netapp.VolumeList, err error)
	ReplicationStatusMethod(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result netapp.ReplicationStatus, err error)
	ResyncReplication(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result autorest.Response, err error)
	Update(ctx context.Context, body netapp.VolumePatch, resourceGroupName string, accountName string, poolName string, volumeName string) (result netapp.Volume, err error)
}

var _ VolumesClientAPI = (*netapp.VolumesClient)(nil)

// MountTargetsClientAPI contains the set of methods on the MountTargetsClient type.
type MountTargetsClientAPI interface {
	List(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result netapp.MountTargetList, err error)
}

var _ MountTargetsClientAPI = (*netapp.MountTargetsClient)(nil)

// SnapshotsClientAPI contains the set of methods on the SnapshotsClient type.
type SnapshotsClientAPI interface {
	Create(ctx context.Context, body netapp.Snapshot, resourceGroupName string, accountName string, poolName string, volumeName string, snapshotName string) (result netapp.SnapshotsCreateFuture, err error)
	Delete(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, snapshotName string) (result netapp.SnapshotsDeleteFuture, err error)
	Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, snapshotName string) (result netapp.Snapshot, err error)
	List(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string) (result netapp.SnapshotsList, err error)
	Update(ctx context.Context, body netapp.SnapshotPatch, resourceGroupName string, accountName string, poolName string, volumeName string, snapshotName string) (result netapp.Snapshot, err error)
}

var _ SnapshotsClientAPI = (*netapp.SnapshotsClient)(nil)

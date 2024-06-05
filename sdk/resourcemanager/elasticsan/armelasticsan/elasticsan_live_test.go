//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armelasticsan_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/elasticsan/armelasticsan"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type ElasticsanTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	elasticSanName    string
	elasticsanId      string
	snapshotName      string
	volumeGroupName   string
	volumeName        string
	location          string
	resourceGroupName string
	subscriptionId    string

	volumeId string
}

func (testsuite *ElasticsanTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.elasticSanName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "elastics", 14, true)
	testsuite.snapshotName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "snapshot", 14, true)
	testsuite.volumeGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "volumegr", 14, true)
	testsuite.volumeName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "volumena", 14, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *ElasticsanTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestElasticsanTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticsanTestSuite))
}

func (testsuite *ElasticsanTestSuite) Prepare() {
	var err error
	// From step ElasticSans_Create
	fmt.Println("Call operation: ElasticSans_Create")
	elasticSansClient, err := armelasticsan.NewElasticSansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	elasticSansClientCreateResponsePoller, err := elasticSansClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, armelasticsan.ElasticSan{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key9316": to.Ptr("ihndtieqibtob"),
		},
		Properties: &armelasticsan.Properties{
			AvailabilityZones: []*string{
				to.Ptr("1")},
			BaseSizeTiB:             to.Ptr[int64](5),
			ExtendedCapacitySizeTiB: to.Ptr[int64](25),
			PublicNetworkAccess:     to.Ptr(armelasticsan.PublicNetworkAccessEnabled),
			SKU: &armelasticsan.SKU{
				Name: to.Ptr(armelasticsan.SKUNamePremiumLRS),
				Tier: to.Ptr(armelasticsan.SKUTierPremium),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	var elasticSansClientCreateResponse *armelasticsan.ElasticSansClientCreateResponse
	elasticSansClientCreateResponse, err = testutil.PollForTest(testsuite.ctx, elasticSansClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.elasticsanId = *elasticSansClientCreateResponse.ID

	// From step VolumeGroups_Create
	fmt.Println("Call operation: VolumeGroups_Create")
	volumeGroupsClient, err := armelasticsan.NewVolumeGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	volumeGroupsClientCreateResponsePoller, err := volumeGroupsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, armelasticsan.VolumeGroup{}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, volumeGroupsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Volumes_Create
	fmt.Println("Call operation: Volumes_Create")
	volumesClient, err := armelasticsan.NewVolumesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	volumesClientCreateResponsePoller, err := volumesClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, testsuite.volumeName, armelasticsan.Volume{
		Properties: &armelasticsan.VolumeProperties{
			// CreationData: &armelasticsan.SourceCreationData{
			// 	CreateSource: to.Ptr(armelasticsan.VolumeCreateOptionNone),
			// 	SourceID:     to.Ptr("ARM Id of Resource"),
			// },
			// ManagedBy: &armelasticsan.ManagedByInfo{
			// 	ResourceID: to.Ptr("mtkeip"),
			// },
			SizeGiB: to.Ptr[int64](9),
		},
	}, nil)
	testsuite.Require().NoError(err)
	volumesClientCreateResponse, err := testutil.PollForTest(testsuite.ctx, volumesClientCreateResponsePoller)
	testsuite.Require().NoError(err)
	testsuite.volumeId = *volumesClientCreateResponse.ID
}

// Microsoft.ElasticSan/elasticSans/{elasticSanName}
func (testsuite *ElasticsanTestSuite) TestElasticSans() {
	var err error
	// From step ElasticSans_ListBySubscription
	fmt.Println("Call operation: ElasticSans_ListBySubscription")
	elasticSansClient, err := armelasticsan.NewElasticSansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	elasticSansClientNewListBySubscriptionPager := elasticSansClient.NewListBySubscriptionPager(nil)
	for elasticSansClientNewListBySubscriptionPager.More() {
		_, err := elasticSansClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ElasticSans_ListByResourceGroup
	fmt.Println("Call operation: ElasticSans_ListByResourceGroup")
	elasticSansClientNewListByResourceGroupPager := elasticSansClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for elasticSansClientNewListByResourceGroupPager.More() {
		_, err := elasticSansClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step ElasticSans_Get
	fmt.Println("Call operation: ElasticSans_Get")
	_, err = elasticSansClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, nil)
	testsuite.Require().NoError(err)

	// From step ElasticSans_Update
	fmt.Println("Call operation: ElasticSans_Update")
	elasticSansClientUpdateResponsePoller, err := elasticSansClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, armelasticsan.Update{
		Properties: &armelasticsan.UpdateProperties{
			BaseSizeTiB:             to.Ptr[int64](13),
			ExtendedCapacitySizeTiB: to.Ptr[int64](29),
			PublicNetworkAccess:     to.Ptr(armelasticsan.PublicNetworkAccessEnabled),
		},
		Tags: map[string]*string{
			"key1931": to.Ptr("yhjwkgmrrwrcoxblgwgzjqusch"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, elasticSansClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ElasticSan/elasticSans/{elasticSanName}/volumegroups/{volumeGroupName}
func (testsuite *ElasticsanTestSuite) TestVolumeGroups() {
	var err error
	// From step VolumeGroups_ListByElasticSan
	fmt.Println("Call operation: VolumeGroups_ListByElasticSan")
	volumeGroupsClient, err := armelasticsan.NewVolumeGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	volumeGroupsClientNewListByElasticSanPager := volumeGroupsClient.NewListByElasticSanPager(testsuite.resourceGroupName, testsuite.elasticSanName, nil)
	for volumeGroupsClientNewListByElasticSanPager.More() {
		_, err := volumeGroupsClientNewListByElasticSanPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VolumeGroups_Get
	fmt.Println("Call operation: VolumeGroups_Get")
	_, err = volumeGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, nil)
	testsuite.Require().NoError(err)

	// From step VolumeGroups_Update
	fmt.Println("Call operation: VolumeGroups_Update")
	volumeGroupsClientUpdateResponsePoller, err := volumeGroupsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, armelasticsan.VolumeGroupUpdate{}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, volumeGroupsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ElasticSan/elasticSans/{elasticSanName}/volumegroups/{volumeGroupName}/volumes/{volumeName}
func (testsuite *ElasticsanTestSuite) TestVolumes() {
	var err error
	// From step Volumes_ListByVolumeGroup
	fmt.Println("Call operation: Volumes_ListByVolumeGroup")
	volumesClient, err := armelasticsan.NewVolumesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	volumesClientNewListByVolumeGroupPager := volumesClient.NewListByVolumeGroupPager(testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, nil)
	for volumesClientNewListByVolumeGroupPager.More() {
		_, err := volumesClientNewListByVolumeGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Volumes_Get
	fmt.Println("Call operation: Volumes_Get")
	_, err = volumesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, testsuite.volumeName, nil)
	testsuite.Require().NoError(err)

	// From step Volumes_Update
	fmt.Println("Call operation: Volumes_Update")
	volumesClientUpdateResponsePoller, err := volumesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, testsuite.volumeName, armelasticsan.VolumeUpdate{
		Properties: &armelasticsan.VolumeUpdateProperties{
			SizeGiB: to.Ptr[int64](11),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, volumesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ElasticSan/elasticSans/{elasticSanName}/volumegroups/{volumeGroupName}/snapshots/{snapshotName}
func (testsuite *ElasticsanTestSuite) TestVolumeSnapshots() {
	var err error
	// From step VolumeSnapshots_Create
	fmt.Println("Call operation: VolumeSnapshots_Create")
	volumeSnapshotsClient, err := armelasticsan.NewVolumeSnapshotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	volumeSnapshotsClientCreateResponsePoller, err := volumeSnapshotsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, testsuite.snapshotName, armelasticsan.Snapshot{
		Properties: &armelasticsan.SnapshotProperties{
			CreationData: &armelasticsan.SnapshotCreationData{
				SourceID: to.Ptr(testsuite.volumeId),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, volumeSnapshotsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step VolumeSnapshots_ListByVolumeGroup
	fmt.Println("Call operation: VolumeSnapshots_ListByVolumeGroup")
	volumeSnapshotsClientNewListByVolumeGroupPager := volumeSnapshotsClient.NewListByVolumeGroupPager(testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, &armelasticsan.VolumeSnapshotsClientListByVolumeGroupOptions{Filter: to.Ptr("volumeName eq " + testsuite.volumeName)})
	for volumeSnapshotsClientNewListByVolumeGroupPager.More() {
		_, err := volumeSnapshotsClientNewListByVolumeGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step VolumeSnapshots_Get
	fmt.Println("Call operation: VolumeSnapshots_Get")
	_, err = volumeSnapshotsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, testsuite.snapshotName, nil)
	testsuite.Require().NoError(err)

	// From step VolumeSnapshots_Delete
	fmt.Println("Call operation: VolumeSnapshots_Delete")
	volumeSnapshotsClientDeleteResponsePoller, err := volumeSnapshotsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, testsuite.snapshotName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, volumeSnapshotsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.ElasticSan/operations
func (testsuite *ElasticsanTestSuite) TestOperations() {
	var err error
	// From step Operations_List
	fmt.Println("Call operation: Operations_List")
	operationsClient, err := armelasticsan.NewOperationsClient(testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	operationsClientNewListPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPager.More() {
		_, err := operationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

// Microsoft.ElasticSan/skus
func (testsuite *ElasticsanTestSuite) TestSkus() {
	var err error
	// From step Skus_List
	fmt.Println("Call operation: Skus_List")
	sKUsClient, err := armelasticsan.NewSKUsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	sKUsClientNewListPager := sKUsClient.NewListPager(&armelasticsan.SKUsClientListOptions{Filter: to.Ptr("obwwdrkq")})
	for sKUsClientNewListPager.More() {
		_, err := sKUsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}
}

func (testsuite *ElasticsanTestSuite) Cleanup() {
	var err error
	// From step Volumes_Delete
	fmt.Println("Call operation: Volumes_Delete")
	volumesClient, err := armelasticsan.NewVolumesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	volumesClientDeleteResponsePoller, err := volumesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, testsuite.volumeName, &armelasticsan.VolumesClientBeginDeleteOptions{XMSDeleteSnapshots: to.Ptr(armelasticsan.XMSDeleteSnapshotsTrue),
		XMSForceDelete: to.Ptr(armelasticsan.XMSForceDeleteTrue),
	})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, volumesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step VolumeGroups_Delete
	fmt.Println("Call operation: VolumeGroups_Delete")
	volumeGroupsClient, err := armelasticsan.NewVolumeGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	volumeGroupsClientDeleteResponsePoller, err := volumeGroupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, testsuite.volumeGroupName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, volumeGroupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// From step ElasticSans_Delete
	fmt.Println("Call operation: ElasticSans_Delete")
	elasticSansClient, err := armelasticsan.NewElasticSansClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	elasticSansClientDeleteResponsePoller, err := elasticSansClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.elasticSanName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, elasticSansClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

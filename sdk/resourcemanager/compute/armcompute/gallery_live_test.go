//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type GalleryTestSuite struct {
	suite.Suite

	ctx                    context.Context
	cred                   azcore.TokenCredential
	options                *arm.ClientOptions
	galleryApplicationName string
	galleryImageName       string
	galleryName            string
	location               string
	resourceGroupName      string
	subscriptionId         string
}

func (testsuite *GalleryTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.galleryApplicationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "galleryapp", 16, false)
	testsuite.galleryImageName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "galleryima", 16, false)
	testsuite.galleryName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "gallerynam", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *GalleryTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestGalleryTestSuite(t *testing.T) {
	suite.Run(t, new(GalleryTestSuite))
}

func (testsuite *GalleryTestSuite) Prepare() {
	var err error
	// From step Galleries_CreateOrUpdate
	fmt.Println("Call operation: Galleries_CreateOrUpdate")
	galleriesClient, err := armcompute.NewGalleriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	galleriesClientCreateOrUpdateResponsePoller, err := galleriesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, armcompute.Gallery{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.GalleryProperties{
			Description: to.Ptr("This is the gallery description."),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleriesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/galleries/{galleryName}
func (testsuite *GalleryTestSuite) TestGalleries() {
	var err error
	// From step Galleries_List
	fmt.Println("Call operation: Galleries_List")
	galleriesClient, err := armcompute.NewGalleriesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	galleriesClientNewListPager := galleriesClient.NewListPager(nil)
	for galleriesClientNewListPager.More() {
		_, err := galleriesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Galleries_ListByResourceGroup
	fmt.Println("Call operation: Galleries_ListByResourceGroup")
	galleriesClientNewListByResourceGroupPager := galleriesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for galleriesClientNewListByResourceGroupPager.More() {
		_, err := galleriesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Galleries_Get
	fmt.Println("Call operation: Galleries_Get")
	_, err = galleriesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, &armcompute.GalleriesClientGetOptions{Select: nil,
		Expand: nil,
	})
	testsuite.Require().NoError(err)

	// From step Galleries_Update
	fmt.Println("Call operation: Galleries_Update")
	galleriesClientUpdateResponsePoller, err := galleriesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, armcompute.GalleryUpdate{
		Properties: &armcompute.GalleryProperties{
			Description: to.Ptr("This is the gallery description."),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleriesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/galleries/{galleryName}/images/{galleryImageName}
func (testsuite *GalleryTestSuite) TestGalleryImages() {
	var err error
	// From step GalleryImages_CreateOrUpdate
	fmt.Println("Call operation: GalleryImages_CreateOrUpdate")
	galleryImagesClient, err := armcompute.NewGalleryImagesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	galleryImagesClientCreateOrUpdateResponsePoller, err := galleryImagesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryImageName, armcompute.GalleryImage{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.GalleryImageProperties{
			HyperVGeneration: to.Ptr(armcompute.HyperVGenerationV1),
			Identifier: &armcompute.GalleryImageIdentifier{
				Offer:     to.Ptr("myOfferName"),
				Publisher: to.Ptr("myPublisherName"),
				SKU:       to.Ptr("mySkuName"),
			},
			OSState: to.Ptr(armcompute.OperatingSystemStateTypesGeneralized),
			OSType:  to.Ptr(armcompute.OperatingSystemTypesWindows),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleryImagesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step GalleryImages_ListByGallery
	fmt.Println("Call operation: GalleryImages_ListByGallery")
	galleryImagesClientNewListByGalleryPager := galleryImagesClient.NewListByGalleryPager(testsuite.resourceGroupName, testsuite.galleryName, nil)
	for galleryImagesClientNewListByGalleryPager.More() {
		_, err := galleryImagesClientNewListByGalleryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step GalleryImages_Get
	fmt.Println("Call operation: GalleryImages_Get")
	_, err = galleryImagesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryImageName, nil)
	testsuite.Require().NoError(err)

	// From step GalleryImages_Update
	fmt.Println("Call operation: GalleryImages_Update")
	galleryImagesClientUpdateResponsePoller, err := galleryImagesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryImageName, armcompute.GalleryImageUpdate{
		Tags: map[string]*string{
			"0": to.Ptr("[object Object]"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleryImagesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step GalleryImages_Delete
	fmt.Println("Call operation: GalleryImages_Delete")
	galleryImagesClientDeleteResponsePoller, err := galleryImagesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryImageName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleryImagesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/galleries/{galleryName}/applications/{galleryApplicationName}
func (testsuite *GalleryTestSuite) TestGalleryApplications() {
	var err error
	// From step GalleryApplications_CreateOrUpdate
	fmt.Println("Call operation: GalleryApplications_CreateOrUpdate")
	galleryApplicationsClient, err := armcompute.NewGalleryApplicationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	galleryApplicationsClientCreateOrUpdateResponsePoller, err := galleryApplicationsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryApplicationName, armcompute.GalleryApplication{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.GalleryApplicationProperties{
			Description: to.Ptr("This is the gallery application description."),
			CustomActions: []*armcompute.GalleryApplicationCustomAction{
				{
					Name:        to.Ptr("myCustomAction"),
					Description: to.Ptr("This is the custom action description."),
					Parameters: []*armcompute.GalleryApplicationCustomActionParameter{
						{
							Name:         to.Ptr("myCustomActionParameter"),
							Type:         to.Ptr(armcompute.GalleryApplicationCustomActionParameterTypeString),
							Description:  to.Ptr("This is the description of the parameter"),
							DefaultValue: to.Ptr("default value of parameter."),
							Required:     to.Ptr(false),
						}},
					Script: to.Ptr("myCustomActionScript"),
				}},
			Eula:                to.Ptr("This is the gallery application EULA."),
			PrivacyStatementURI: to.Ptr("myPrivacyStatementUri}"),
			ReleaseNoteURI:      to.Ptr("myReleaseNoteUri"),
			SupportedOSType:     to.Ptr(armcompute.OperatingSystemTypesWindows),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleryApplicationsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step GalleryApplications_ListByGallery
	fmt.Println("Call operation: GalleryApplications_ListByGallery")
	galleryApplicationsClientNewListByGalleryPager := galleryApplicationsClient.NewListByGalleryPager(testsuite.resourceGroupName, testsuite.galleryName, nil)
	for galleryApplicationsClientNewListByGalleryPager.More() {
		_, err := galleryApplicationsClientNewListByGalleryPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step GalleryApplications_Get
	fmt.Println("Call operation: GalleryApplications_Get")
	_, err = galleryApplicationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryApplicationName, nil)
	testsuite.Require().NoError(err)

	// From step GalleryApplications_Update
	fmt.Println("Call operation: GalleryApplications_Update")
	galleryApplicationsClientUpdateResponsePoller, err := galleryApplicationsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryApplicationName, armcompute.GalleryApplicationUpdate{
		Properties: &armcompute.GalleryApplicationProperties{
			Description: to.Ptr("This is the gallery application description."),
			CustomActions: []*armcompute.GalleryApplicationCustomAction{
				{
					Name:        to.Ptr("myCustomAction"),
					Description: to.Ptr("This is the custom action description."),
					Parameters: []*armcompute.GalleryApplicationCustomActionParameter{
						{
							Name:         to.Ptr("myCustomActionParameter"),
							Type:         to.Ptr(armcompute.GalleryApplicationCustomActionParameterTypeString),
							Description:  to.Ptr("This is the description of the parameter"),
							DefaultValue: to.Ptr("default value of parameter."),
							Required:     to.Ptr(false),
						}},
					Script: to.Ptr("myCustomActionScript"),
				}},
			Eula:                to.Ptr("This is the gallery application EULA."),
			PrivacyStatementURI: to.Ptr("myPrivacyStatementUri}"),
			ReleaseNoteURI:      to.Ptr("myReleaseNoteUri"),
			SupportedOSType:     to.Ptr(armcompute.OperatingSystemTypesWindows),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleryApplicationsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step GalleryApplications_Delete
	fmt.Println("Call operation: GalleryApplications_Delete")
	galleryApplicationsClientDeleteResponsePoller, err := galleryApplicationsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.galleryName, testsuite.galleryApplicationName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, galleryApplicationsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

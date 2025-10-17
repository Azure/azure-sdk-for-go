//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v3"
	"github.com/stretchr/testify/suite"
)

type BlobTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *BlobTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName = "blobaccountnam"
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *BlobTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestBlobTestSuite(t *testing.T) {
	suite.Run(t, new(BlobTestSuite))
}

func (testsuite *BlobTestSuite) Prepare() {
	var err error
	// From step StorageAccount_Create
	fmt.Println("Call operation: StorageAccounts_Create")
	accountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountsClientCreateResponsePoller, err := accountsClient.BeginCreate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.AccountCreateParameters{
		Kind:     to.Ptr(armstorage.KindStorageV2),
		Location: to.Ptr(testsuite.location),
		Properties: &armstorage.AccountPropertiesCreateParameters{
			AllowBlobPublicAccess:        to.Ptr(false),
			AllowSharedKeyAccess:         to.Ptr(true),
			DefaultToOAuthAuthentication: to.Ptr(false),
			Encryption: &armstorage.Encryption{
				KeySource:                       to.Ptr(armstorage.KeySourceMicrosoftStorage),
				RequireInfrastructureEncryption: to.Ptr(false),
				Services: &armstorage.EncryptionServices{
					Blob: &armstorage.EncryptionService{
						Enabled: to.Ptr(true),
						KeyType: to.Ptr(armstorage.KeyTypeAccount),
					},
					File: &armstorage.EncryptionService{
						Enabled: to.Ptr(true),
						KeyType: to.Ptr(armstorage.KeyTypeAccount),
					},
				},
			},
			IsHnsEnabled:  to.Ptr(true),
			IsSftpEnabled: to.Ptr(true),
			KeyPolicy: &armstorage.KeyPolicy{
				KeyExpirationPeriodInDays: to.Ptr[int32](20),
			},
			MinimumTLSVersion: to.Ptr(armstorage.MinimumTLSVersionTLS12),
			RoutingPreference: &armstorage.RoutingPreference{
				PublishInternetEndpoints:  to.Ptr(true),
				PublishMicrosoftEndpoints: to.Ptr(true),
				RoutingChoice:             to.Ptr(armstorage.RoutingChoiceMicrosoftRouting),
			},
			SasPolicy: &armstorage.SasPolicy{
				ExpirationAction:    to.Ptr(armstorage.ExpirationActionLog),
				SasExpirationPeriod: to.Ptr("1.15:59:59"),
			},
		},
		SKU: &armstorage.SKU{
			Name: to.Ptr(armstorage.SKUNameStandardGRS),
		},
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, accountsClientCreateResponsePoller)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/blobServices/{BlobServicesName}
func (testsuite *BlobTestSuite) TestBlobServices() {
	var err error
	// From step BlobServices_SetServiceProperties
	fmt.Println("Call operation: BlobServices_SetServiceProperties")
	blobServicesClient, err := armstorage.NewBlobServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = blobServicesClient.SetServiceProperties(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.BlobServiceProperties{}, nil)
	testsuite.Require().NoError(err)

	// From step BlobServices_List
	fmt.Println("Call operation: BlobServices_List")
	blobServicesClientNewListPager := blobServicesClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for blobServicesClientNewListPager.More() {
		_, err := blobServicesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BlobServices_GetServiceProperties
	fmt.Println("Call operation: BlobServices_GetServiceProperties")
	_, err = blobServicesClient.GetServiceProperties(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/blobServices/default/containers/{containerName}
func (testsuite *BlobTestSuite) TestBlobContainers() {
	containerName := "containerna"
	var etag string
	var err error
	// From step BlobContainers_Create
	fmt.Println("Call operation: BlobContainers_Create")
	blobContainersClient, err := armstorage.NewBlobContainersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = blobContainersClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, armstorage.BlobContainer{}, nil)
	testsuite.Require().NoError(err)

	// From step BlobContainers_List
	fmt.Println("Call operation: BlobContainers_List")
	blobContainersClientNewListPager := blobContainersClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armstorage.BlobContainersClientListOptions{Maxpagesize: nil,
		Filter:  nil,
		Include: nil,
	})
	for blobContainersClientNewListPager.More() {
		_, err := blobContainersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step BlobContainers_Get
	fmt.Println("Call operation: BlobContainers_Get")
	_, err = blobContainersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, nil)
	testsuite.Require().NoError(err)

	// From step BlobContainers_Update
	fmt.Println("Call operation: BlobContainers_Update")
	_, err = blobContainersClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, armstorage.BlobContainer{
		ContainerProperties: &armstorage.ContainerProperties{
			Metadata: map[string]*string{
				"metadata": to.Ptr("true"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step BlobContainers_Lease
	fmt.Println("Call operation: BlobContainers_Lease")
	_, err = blobContainersClient.Lease(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, &armstorage.BlobContainersClientLeaseOptions{Parameters: &armstorage.LeaseContainerRequest{
		Action:        to.Ptr(armstorage.LeaseContainerRequestActionAcquire),
		LeaseDuration: to.Ptr[int32](-1),
	},
	})
	testsuite.Require().NoError(err)

	// From step BlobContainers_Lease_Break
	fmt.Println("Call operation: BlobContainers_Lease")
	_, err = blobContainersClient.Lease(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, &armstorage.BlobContainersClientLeaseOptions{Parameters: &armstorage.LeaseContainerRequest{
		Action:        to.Ptr(armstorage.LeaseContainerRequestActionBreak),
		LeaseDuration: to.Ptr[int32](-1),
	},
	})
	testsuite.Require().NoError(err)

	// From step BlobContainers_SetLegalHold
	fmt.Println("Call operation: BlobContainers_SetLegalHold")
	_, err = blobContainersClient.SetLegalHold(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, armstorage.LegalHold{
		Tags: []*string{
			to.Ptr("tag1"),
			to.Ptr("tag2"),
			to.Ptr("tag3")},
	}, nil)
	testsuite.Require().NoError(err)

	// From step BlobContainers_ClearLegalHold
	fmt.Println("Call operation: BlobContainers_ClearLegalHold")
	_, err = blobContainersClient.ClearLegalHold(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, armstorage.LegalHold{
		Tags: []*string{
			to.Ptr("tag1"),
			to.Ptr("tag2"),
			to.Ptr("tag3")},
	}, nil)
	testsuite.Require().NoError(err)

	// From step BlobContainers_CreateOrUpdateImmutabilityPolicy
	fmt.Println("Call operation: BlobContainers_CreateOrUpdateImmutabilityPolicy")
	blobContainersClientCreateOrUpdateImmutabilityPolicyResponse, err := blobContainersClient.CreateOrUpdateImmutabilityPolicy(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, &armstorage.BlobContainersClientCreateOrUpdateImmutabilityPolicyOptions{
		Parameters: &armstorage.ImmutabilityPolicy{
			Properties: &armstorage.ImmutabilityPolicyProperty{
				AllowProtectedAppendWrites:            to.Ptr(true),
				ImmutabilityPeriodSinceCreationInDays: to.Ptr[int32](3),
			},
		},
	})
	testsuite.Require().NoError(err)
	etag = *blobContainersClientCreateOrUpdateImmutabilityPolicyResponse.Etag

	// From step BlobContainers_GetImmutabilityPolicy
	fmt.Println("Call operation: BlobContainers_GetImmutabilityPolicy")
	_, err = blobContainersClient.GetImmutabilityPolicy(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, &armstorage.BlobContainersClientGetImmutabilityPolicyOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step BlobContainers_DeleteImmutabilityPolicy
	fmt.Println("Call operation: BlobContainers_DeleteImmutabilityPolicy")
	_, err = blobContainersClient.DeleteImmutabilityPolicy(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, etag, nil)
	testsuite.Require().NoError(err)

	// From step BlobContainers_Delete
	fmt.Println("Call operation: BlobContainers_Delete")
	_, err = blobContainersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, containerName, nil)
	testsuite.Require().NoError(err)
}

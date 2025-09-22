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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v2"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *FileTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName = "fileaccountnam1"
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *FileTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestFileTestSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}

func (testsuite *FileTestSuite) Prepare() {
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

// Microsoft.Storage/storageAccounts/{accountName}/fileServices/{FileServicesName}
func (testsuite *FileTestSuite) TestFileServices() {
	var err error
	// From step FileServices_SetServiceProperties
	fmt.Println("Call operation: FileServices_SetServiceProperties")
	fileServicesClient, err := armstorage.NewFileServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = fileServicesClient.SetServiceProperties(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.FileServiceProperties{}, nil)
	testsuite.Require().NoError(err)

	// From step FileServices_List
	fmt.Println("Call operation: FileServices_List")
	_, err = fileServicesClient.List(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step FileServices_GetServiceProperties
	fmt.Println("Call operation: FileServices_GetServiceProperties")
	_, err = fileServicesClient.GetServiceProperties(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/fileServices/default/shares/{shareName}
func (testsuite *FileTestSuite) TestFileShares() {
	shareName := "sharename"
	var err error
	// From step FileShares_Create
	fmt.Println("Call operation: FileShares_Create")
	fileSharesClient, err := armstorage.NewFileSharesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = fileSharesClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, shareName, armstorage.FileShare{}, &armstorage.FileSharesClientCreateOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step FileShares_List
	fmt.Println("Call operation: FileShares_List")
	fileSharesClientNewListPager := fileSharesClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armstorage.FileSharesClientListOptions{Maxpagesize: nil,
		Filter: nil,
		Expand: nil,
	})
	for fileSharesClientNewListPager.More() {
		_, err := fileSharesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FileShares_Get
	fmt.Println("Call operation: FileShares_Get")
	_, err = fileSharesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, shareName, &armstorage.FileSharesClientGetOptions{Expand: nil,
		XMSSnapshot: nil,
	})
	testsuite.Require().NoError(err)

	// From step FileShares_Update
	fmt.Println("Call operation: FileShares_Update")
	_, err = fileSharesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, shareName, armstorage.FileShare{
		FileShareProperties: &armstorage.FileShareProperties{
			Metadata: map[string]*string{
				"type": to.Ptr("image"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step FileShares_Lease
	fmt.Println("Call operation: FileShares_Lease")
	_, err = fileSharesClient.Lease(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, shareName, &armstorage.FileSharesClientLeaseOptions{XMSSnapshot: nil,
		Parameters: &armstorage.LeaseShareRequest{
			Action:        to.Ptr(armstorage.LeaseShareActionAcquire),
			LeaseDuration: to.Ptr[int32](-1),
		},
	})
	testsuite.Require().NoError(err)

	// From step FileShares_Lease_Break
	fmt.Println("Call operation: FileShares_Lease")
	_, err = fileSharesClient.Lease(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, shareName, &armstorage.FileSharesClientLeaseOptions{XMSSnapshot: nil,
		Parameters: &armstorage.LeaseShareRequest{
			Action:        to.Ptr(armstorage.LeaseShareActionBreak),
			LeaseDuration: to.Ptr[int32](-1),
		},
	})
	testsuite.Require().NoError(err)

	// From step FileShares_Delete
	fmt.Println("Call operation: FileShares_Delete")
	_, err = fileSharesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, shareName, &armstorage.FileSharesClientDeleteOptions{XMSSnapshot: nil,
		Include: nil,
	})
	testsuite.Require().NoError(err)
}

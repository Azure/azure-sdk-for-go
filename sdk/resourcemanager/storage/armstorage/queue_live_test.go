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

type QueueTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	accountName       string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *QueueTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName = "queueaccountnam"
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *QueueTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}

func (testsuite *QueueTestSuite) Prepare() {
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

// Microsoft.Storage/storageAccounts/{accountName}/queueServices/{queueServiceName}
func (testsuite *QueueTestSuite) TestQueueServices() {
	var err error
	// From step QueueServices_SetServiceProperties
	fmt.Println("Call operation: QueueServices_SetServiceProperties")
	queueServicesClient, err := armstorage.NewQueueServicesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = queueServicesClient.SetServiceProperties(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armstorage.QueueServiceProperties{}, nil)
	testsuite.Require().NoError(err)

	// From step QueueServices_List
	fmt.Println("Call operation: QueueServices_List")
	_, err = queueServicesClient.List(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step QueueServices_GetServiceProperties
	fmt.Println("Call operation: QueueServices_GetServiceProperties")
	_, err = queueServicesClient.GetServiceProperties(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Storage/storageAccounts/{accountName}/queueServices/default/queues/{queueName}
func (testsuite *QueueTestSuite) TestQueue() {
	queueName := "queuename"
	var err error
	// From step Queue_Create
	fmt.Println("Call operation: Queue_Create")
	queueClient, err := armstorage.NewQueueClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = queueClient.Create(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, queueName, armstorage.Queue{}, nil)
	testsuite.Require().NoError(err)

	// From step Queue_List
	fmt.Println("Call operation: Queue_List")
	queueClientNewListPager := queueClient.NewListPager(testsuite.resourceGroupName, testsuite.accountName, &armstorage.QueueClientListOptions{Maxpagesize: nil,
		Filter: nil,
	})
	for queueClientNewListPager.More() {
		_, err := queueClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Queue_Get
	fmt.Println("Call operation: Queue_Get")
	_, err = queueClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, queueName, nil)
	testsuite.Require().NoError(err)

	// From step Queue_Update
	fmt.Println("Call operation: Queue_Update")
	_, err = queueClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, queueName, armstorage.Queue{}, nil)
	testsuite.Require().NoError(err)

	// From step Queue_Delete
	fmt.Println("Call operation: Queue_Delete")
	_, err = queueClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, queueName, nil)
	testsuite.Require().NoError(err)
}

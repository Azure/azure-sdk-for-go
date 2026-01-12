// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armmaps_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/maps/armmaps/v2"
	"github.com/stretchr/testify/suite"
)

type MapsTestSuite struct {
	suite.Suite

	ctx			context.Context
	cred			azcore.TokenCredential
	options			*arm.ClientOptions
	accountName		string
	armEndpoint		string
	creatorName		string
	location		string
	resourceGroupName	string
	subscriptionId		string
}

func (testsuite *MapsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.accountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "accountn", 14, false)
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.creatorName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "creatorn", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus2")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *MapsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestMapsTestSuite(t *testing.T) {
	suite.Run(t, new(MapsTestSuite))
}

func (testsuite *MapsTestSuite) Prepare() {
	var err error
	// From step Accounts_CreateOrUpdate
	fmt.Println("Call operation: Accounts_CreateOrUpdate")
	accountsClient, err := armmaps.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = accountsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmaps.Account{
		Location:	to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"test": to.Ptr("true"),
		},
		Kind:	to.Ptr(armmaps.KindGen2),
		Properties: &armmaps.AccountProperties{
			Cors: &armmaps.CorsRules{
				CorsRules: []*armmaps.CorsRule{
					{
						AllowedOrigins: []*string{
							to.Ptr("http://www.contoso.com"),
							to.Ptr("http://www.fabrikam.com")},
					}},
			},
			DisableLocalAuth:	to.Ptr(false),
		},
		SKU: &armmaps.SKU{
			Name: to.Ptr(armmaps.NameG2),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Maps/accounts/{accountName}
func (testsuite *MapsTestSuite) TestAccounts() {
	var err error
	// From step Accounts_ListBySubscription
	fmt.Println("Call operation: Accounts_ListBySubscription")
	accountsClient, err := armmaps.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	accountsClientNewListBySubscriptionPager := accountsClient.NewListBySubscriptionPager(nil)
	for accountsClientNewListBySubscriptionPager.More() {
		_, err := accountsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Accounts_ListByResourceGroup
	fmt.Println("Call operation: Accounts_ListByResourceGroup")
	accountsClientNewListByResourceGroupPager := accountsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for accountsClientNewListByResourceGroupPager.More() {
		_, err := accountsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Accounts_Get
	fmt.Println("Call operation: Accounts_Get")
	_, err = accountsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step Accounts_ListSas
	fmt.Println("Call operation: Accounts_ListSas")
	_, err = accountsClient.ListSas(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmaps.AccountSasParameters{
		Expiry:			to.Ptr("2017-05-24T11:42:03.1567373Z"),
		MaxRatePerSecond:	to.Ptr[int32](500),
		PrincipalID:		to.Ptr("e917f87b-324d-4728-98ed-e31d311a7d65"),
		Regions: []*string{
			to.Ptr("eastus")},
		SigningKey:	to.Ptr(armmaps.SigningKeyPrimaryKey),
		Start:		to.Ptr("2017-05-24T10:42:03.1567373Z"),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Accounts_RegenerateKeys
	fmt.Println("Call operation: Accounts_RegenerateKeys")
	_, err = accountsClient.RegenerateKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmaps.KeySpecification{
		KeyType: to.Ptr(armmaps.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step Accounts_ListKeys
	fmt.Println("Call operation: Accounts_ListKeys")
	_, err = accountsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)

	// From step Accounts_Update
	fmt.Println("Call operation: Accounts_Update")
	_, err = accountsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, armmaps.AccountUpdateParameters{
		Tags: map[string]*string{
			"specialTag": to.Ptr("true"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// The current provisioningState must transition to a terminal state before the resource can be updated.
	recordMode := recording.GetEnvVariable("AZURE_RECORD_MODE", "playback")
	if recordMode == "record" {
		time.Sleep(60 * time.Second)
	}
}

// Microsoft.Maps/accounts/{accountName}/creators/{creatorName}
func (testsuite *MapsTestSuite) TestCreators() {
	var err error
	// From step Creators_CreateOrUpdate
	fmt.Println("Call operation: Creators_CreateOrUpdate")
	creatorsClient, err := armmaps.NewCreatorsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = creatorsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.creatorName, armmaps.Creator{
		Location:	to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"test": to.Ptr("true"),
		},
		Properties: &armmaps.CreatorProperties{
			StorageUnits: to.Ptr[int32](5),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Creators_ListByAccount
	fmt.Println("Call operation: Creators_ListByAccount")
	creatorsClientNewListByAccountPager := creatorsClient.NewListByAccountPager(testsuite.resourceGroupName, testsuite.accountName, nil)
	for creatorsClientNewListByAccountPager.More() {
		_, err := creatorsClientNewListByAccountPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Creators_Get
	fmt.Println("Call operation: Creators_Get")
	_, err = creatorsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.creatorName, nil)
	testsuite.Require().NoError(err)

	// From step Creators_Update
	fmt.Println("Call operation: Creators_Update")
	_, err = creatorsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.creatorName, armmaps.CreatorUpdateParameters{
		Properties: &armmaps.CreatorProperties{
			StorageUnits: to.Ptr[int32](10),
		},
		Tags: map[string]*string{
			"specialTag": to.Ptr("true"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step Creators_Delete
	fmt.Println("Call operation: Creators_Delete")
	_, err = creatorsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, testsuite.creatorName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *MapsTestSuite) Cleanup() {
	var err error
	// From step Accounts_Delete
	fmt.Println("Call operation: Accounts_Delete")
	accountsClient, err := armmaps.NewAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = accountsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.accountName, nil)
	testsuite.Require().NoError(err)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armneonpostgres_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/neonpostgres/armneonpostgres"
	"github.com/stretchr/testify/suite"
)

const (
	ResourceLocation = "eastus2"
)

type NeonpostgresTestSuite struct {
	suite.Suite
	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	armEndpoint       string
	location          string
	resourceGroupName string
	subscriptionId    string
	organizationName  string
}

func (testsuite *NeonpostgresTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.organizationName = "testogname"
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	fmt.Println("testsuite.resourceGroupName:", testsuite.resourceGroupName)
	testsuite.Prepare()
}

func (testsuite *NeonpostgresTestSuite) TearDownSuite() {
	testsuite.CleanUp()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	time.Sleep(time.Second * 3)
	testutil.StopRecording(testsuite.T())
}

func TestNeonpostgresTestSuite(t *testing.T) {
	suite.Run(t, new(NeonpostgresTestSuite))
}

func (testsuite *NeonpostgresTestSuite) TestOptionList() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armneonpostgres.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			fmt.Println("page value", v)
		}
		page = armneonpostgres.OperationsClientListResponse{
			OperationListResult: armneonpostgres.OperationListResult{
				Value: []*armneonpostgres.Operation{
					{
						Name:         to.Ptr("ekvuojrpbuyogrikfuzyghio"),
						IsDataAction: to.Ptr(true),
						Display: &armneonpostgres.OperationDisplay{
							Provider:    to.Ptr("lvppuskqcdcoejwuvsqrloczvnouy"),
							Resource:    to.Ptr("rvvd"),
							Operation:   to.Ptr("odjjqnodcgszczpsdrhrpwmqssrybr"),
							Description: to.Ptr("bwmstukaiaoisiu"),
						},
						Origin:     to.Ptr(armneonpostgres.OriginUser),
						ActionType: to.Ptr(armneonpostgres.ActionTypeInternal),
					},
				},
				NextLink: to.Ptr("https://contoso.com/nextlink"),
			},
		}
	}
}

func (testsuite *NeonpostgresTestSuite) Prepare() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armneonpostgres.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	_, err = clientFactory.NewOrganizationsClient().BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.organizationName, armneonpostgres.OrganizationResource{
		Properties: &armneonpostgres.OrganizationProperties{
			MarketplaceDetails: &armneonpostgres.MarketplaceDetails{
				SubscriptionID:     to.Ptr(testsuite.subscriptionId),
				SubscriptionStatus: to.Ptr(armneonpostgres.MarketplaceSubscriptionStatusPendingFulfillmentStart),
				OfferDetails: &armneonpostgres.OfferDetails{
					PublisherID: to.Ptr("hporaxnopmolttlnkbarw"),
					OfferID:     to.Ptr("bunyeeupoedueofwrzej"),
					PlanID:      to.Ptr("nlbfiwtslenfwek"),
					PlanName:    to.Ptr("ljbmgpkfqklaufacbpml"),
					TermUnit:    to.Ptr("qbcq"),
					TermID:      to.Ptr("aedlchikwqckuploswthvshe"),
				},
			},
			UserDetails: &armneonpostgres.UserDetails{
				FirstName:    to.Ptr("buwwe"),
				LastName:     to.Ptr("escynjpynkoox"),
				EmailAddress: to.Ptr("3i_%@w8-y.H-p.tvj.dG"),
				Upn:          to.Ptr("fwedjamgwwrotcjaucuzdwycfjdqn"),
				PhoneNumber:  to.Ptr("dlrqoowumy"),
			},
			CompanyDetails: &armneonpostgres.CompanyDetails{
				CompanyName:       to.Ptr("uxn"),
				Country:           to.Ptr("lpajqzptqchuko"),
				OfficeAddress:     to.Ptr("chpkrlpmfslmawgunjxdllzcrctykq"),
				BusinessPhone:     to.Ptr("hbeb"),
				Domain:            to.Ptr("krjldeakhwiepvs"),
				NumberOfEmployees: to.Ptr[int64](23),
			},
			PartnerOrganizationProperties: &armneonpostgres.PartnerOrganizationProperties{
				OrganizationID:   to.Ptr("nrhvoqzulowcunhmvwfgjcaibvwcl"),
				OrganizationName: to.Ptr("2__.-"),
				SingleSignOnProperties: &armneonpostgres.SingleSignOnProperties{
					SingleSignOnState: to.Ptr(armneonpostgres.SingleSignOnStatesInitial),
					EnterpriseAppID:   to.Ptr("fpibacregjfncfdsojs"),
					SingleSignOnURL:   to.Ptr("tmojh"),
					AADDomains: []*string{
						to.Ptr("kndszgrwzbvvlssvkej"),
					},
				},
			},
		},
		Tags: map[string]*string{
			"key2099": to.Ptr("omjjymaqtrqzksxszhzgyl"),
		},
		Location: to.Ptr(testsuite.location),
		Name:     to.Ptr(testsuite.organizationName),
	}, nil)
	testsuite.Require().NoError(err)
	// _, err = testutil.PollForTest(testsuite.ctx, poller)
	// testsuite.Require().NoError(err)

}

func (testsuite *NeonpostgresTestSuite) CleanUp() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armneonpostgres.NewClientFactory(testsuite.subscriptionId, cred, testsuite.options)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewOrganizationsClient().BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.organizationName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, poller)
	testsuite.Require().NoError(err)
}

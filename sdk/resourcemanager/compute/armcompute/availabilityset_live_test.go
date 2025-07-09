//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type AvailabilitySetTestSuite struct {
	suite.Suite

	ctx                 context.Context
	cred                azcore.TokenCredential
	options             *arm.ClientOptions
	location            string
	resourceGroupName   string
	subscriptionId      string
	availabilitySetName string
}

func (testsuite *AvailabilitySetTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.availabilitySetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "availabili", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *AvailabilitySetTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestAvailabilitySetTestSuite(t *testing.T) {
	suite.Run(t, new(AvailabilitySetTestSuite))
}

// Microsoft.Compute/availabilitySets
func (testsuite *AvailabilitySetTestSuite) TestAvailabilitySets() {
	var err error
	// From step AvailabilitySets_CreateOrUpdate
	fmt.Println("Call operation: AvailabilitySets_CreateOrUpdate")
	availabilitySetsClient, err := armcompute.NewAvailabilitySetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = availabilitySetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, armcompute.AvailabilitySet{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.AvailabilitySetProperties{
			PlatformFaultDomainCount:  to.Ptr[int32](2),
			PlatformUpdateDomainCount: to.Ptr[int32](20),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step AvailabilitySets_ListBySubscription
	fmt.Println("Call operation: AvailabilitySets_ListBySubscription")
	availabilitySetsClientNewListBySubscriptionPager := availabilitySetsClient.NewListBySubscriptionPager(&armcompute.AvailabilitySetsClientListBySubscriptionOptions{Expand: nil})
	for availabilitySetsClientNewListBySubscriptionPager.More() {
		_, err := availabilitySetsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AvailabilitySets_List
	fmt.Println("Call operation: AvailabilitySets_List")
	availabilitySetsClientNewListPager := availabilitySetsClient.NewListPager(testsuite.resourceGroupName, nil)
	for availabilitySetsClientNewListPager.More() {
		_, err := availabilitySetsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AvailabilitySets_Get
	fmt.Println("Call operation: AvailabilitySets_Get")
	_, err = availabilitySetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, nil)
	testsuite.Require().NoError(err)

	// From step AvailabilitySets_ListAvailableSizes
	fmt.Println("Call operation: AvailabilitySets_ListAvailableSizes")
	availabilitySetsClientNewListAvailableSizesPager := availabilitySetsClient.NewListAvailableSizesPager(testsuite.resourceGroupName, testsuite.availabilitySetName, nil)
	for availabilitySetsClientNewListAvailableSizesPager.More() {
		_, err := availabilitySetsClientNewListAvailableSizesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step AvailabilitySets_Update
	fmt.Println("Call operation: AvailabilitySets_Update")
	_, err = availabilitySetsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, armcompute.AvailabilitySetUpdate{}, nil)
	testsuite.Require().NoError(err)

	// From step AvailabilitySets_Delete
	fmt.Println("Call operation: AvailabilitySets_Delete")
	_, err = availabilitySetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.availabilitySetName, nil)
	testsuite.Require().NoError(err)
}

func GetCred(options azcore.ClientOptions) (azcore.TokenCredential, error) {
	accessToken := os.Getenv("SYSTEM_ACCESSTOKEN")
	clientID := os.Getenv("AZURESUBSCRIPTION_CLIENT_ID")
	connectionID := os.Getenv("AZURESUBSCRIPTION_SERVICE_CONNECTION_ID")
	tenant := os.Getenv("AZURESUBSCRIPTION_TENANT_ID")
	if accessToken != "" && clientID != "" && connectionID != "" && tenant != "" {
		return azidentity.NewAzurePipelinesCredential(tenant, clientID, connectionID, accessToken, &azidentity.AzurePipelinesCredentialOptions{
			ClientOptions: options,
		})
	}
	if s := os.Getenv("AZURE_SERVICE_DIRECTORY"); s != "" {
		// New-TestResources.ps1 has configured this environment, possibly with service principal details
		clientID := os.Getenv(s + "_CLIENT_ID")
		secret := os.Getenv(s + "_CLIENT_SECRET")
		tenant := os.Getenv(s + "_TENANT_ID")
		if clientID != "" && secret != "" && tenant != "" {
			return azidentity.NewClientSecretCredential(tenant, clientID, secret, &azidentity.ClientSecretCredentialOptions{
				ClientOptions: options,
			})
		}
	}
	return azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: options,
	})
}

// Added to test IPv6 connectivity in Live mode for ARM clients
func (testsuite *AvailabilitySetTestSuite) TestAvailabilitySets_Live_IPv6() {
	if recording.GetRecordMode() != recording.LiveMode {
		testsuite.T().Skip()
	}

	ipv6Client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
				dialer := net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}
				return dialer.DialContext(ctx, "tcp6", addr)
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	cred, err := GetCred(azcore.ClientOptions{Transport: ipv6Client})
	testsuite.Require().NoError(err)

	availabilitySetsClient, err := armcompute.NewAvailabilitySetsClient(testsuite.subscriptionId, cred, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: ipv6Client,
		},
	})
	testsuite.Require().NoError(err)

	availabilitySetName := testsuite.availabilitySetName + "IPv6"
	_, err = availabilitySetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, availabilitySetName, armcompute.AvailabilitySet{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.AvailabilitySetProperties{
			PlatformFaultDomainCount:  to.Ptr[int32](2),
			PlatformUpdateDomainCount: to.Ptr[int32](20),
		},
	}, nil)
	testsuite.Require().NoError(err)

	_, err = availabilitySetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, availabilitySetName, nil)
	testsuite.Require().NoError(err)

	availabilitySetsClientNewListAvailableSizesPager := availabilitySetsClient.NewListAvailableSizesPager(testsuite.resourceGroupName, availabilitySetName, nil)
	for availabilitySetsClientNewListAvailableSizesPager.More() {
		_, err := availabilitySetsClientNewListAvailableSizesPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	_, err = availabilitySetsClient.Update(testsuite.ctx, testsuite.resourceGroupName, availabilitySetName, armcompute.AvailabilitySetUpdate{}, nil)
	testsuite.Require().NoError(err)

	_, err = availabilitySetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, availabilitySetName, nil)
	testsuite.Require().NoError(err)
}

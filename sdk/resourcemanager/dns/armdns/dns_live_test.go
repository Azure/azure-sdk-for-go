// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type DnsTestSuite struct {
	suite.Suite

	ctx                   context.Context
	cred                  azcore.TokenCredential
	options               *arm.ClientOptions
	armEndpoint           string
	relativeRecordSetName string
	zoneName              string
	location              string
	resourceGroupName     string
	subscriptionId        string
}

func (testsuite *DnsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.armEndpoint = "https://management.azure.com"
	testsuite.relativeRecordSetName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "recordsetna", 17, false)
	testsuite.zoneName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "zonename", 14, false)
	testsuite.zoneName += ".com"
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DnsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDnsTestSuite(t *testing.T) {
	suite.Run(t, new(DnsTestSuite))
}

func (testsuite *DnsTestSuite) Prepare() {
	var err error
	// From step Zones_CreateOrUpdate
	fmt.Println("Call operation: Zones_CreateOrUpdate")
	zonesClient, err := armdns.NewZonesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = zonesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, armdns.Zone{
		Location: to.Ptr("Global"),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, &armdns.ZonesClientCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)
}

// Microsoft.Network/dnsZones/{zoneName}
func (testsuite *DnsTestSuite) TestZones() {
	var err error
	// From step Zones_List
	fmt.Println("Call operation: Zones_List")
	zonesClient, err := armdns.NewZonesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	zonesClientNewListPager := zonesClient.NewListPager(&armdns.ZonesClientListOptions{Top: nil})
	for zonesClientNewListPager.More() {
		_, err := zonesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Zones_ListByResourceGroup
	fmt.Println("Call operation: Zones_ListByResourceGroup")
	zonesClientNewListByResourceGroupPager := zonesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armdns.ZonesClientListByResourceGroupOptions{Top: nil})
	for zonesClientNewListByResourceGroupPager.More() {
		_, err := zonesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Zones_Get
	fmt.Println("Call operation: Zones_Get")
	_, err = zonesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, nil)
	testsuite.Require().NoError(err)

	// From step Zones_Update
	fmt.Println("Call operation: Zones_Update")
	_, err = zonesClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, armdns.ZoneUpdate{
		Tags: map[string]*string{
			"key2": to.Ptr("value2"),
		},
	}, &armdns.ZonesClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.Network/dnsZones/{zoneName}/{recordType}/{relativeRecordSetName}
func (testsuite *DnsTestSuite) TestRecordSets() {
	var err error
	// From step RecordSets_CreateOrUpdate
	fmt.Println("Call operation: RecordSets_CreateOrUpdate")
	recordSetsClient, err := armdns.NewRecordSetsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = recordSetsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, testsuite.relativeRecordSetName, armdns.RecordTypeA, armdns.RecordSet{
		Properties: &armdns.RecordSetProperties{
			ARecords: []*armdns.ARecord{
				{
					IPv4Address: to.Ptr("127.0.0.1"),
				}},
			TTL: to.Ptr[int64](3600),
			Metadata: map[string]*string{
				"key1": to.Ptr("value1"),
			},
		},
	}, &armdns.RecordSetsClientCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
	testsuite.Require().NoError(err)

	// From step RecordSets_ListAllByDnsZone
	fmt.Println("Call operation: RecordSets_ListAllByDnsZone")
	recordSetsClientNewListAllByDNSZonePager := recordSetsClient.NewListAllByDNSZonePager(testsuite.resourceGroupName, testsuite.zoneName, &armdns.RecordSetsClientListAllByDNSZoneOptions{Top: nil,
		RecordSetNameSuffix: nil,
	})
	for recordSetsClientNewListAllByDNSZonePager.More() {
		_, err := recordSetsClientNewListAllByDNSZonePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RecordSets_ListByType
	fmt.Println("Call operation: RecordSets_ListByType")
	recordSetsClientNewListByTypePager := recordSetsClient.NewListByTypePager(testsuite.resourceGroupName, testsuite.zoneName, armdns.RecordTypeA, &armdns.RecordSetsClientListByTypeOptions{Top: nil,
		Recordsetnamesuffix: nil,
	})
	for recordSetsClientNewListByTypePager.More() {
		_, err := recordSetsClientNewListByTypePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RecordSets_ListByDnsZone
	fmt.Println("Call operation: RecordSets_ListByDnsZone")
	recordSetsClientNewListByDNSZonePager := recordSetsClient.NewListByDNSZonePager(testsuite.resourceGroupName, testsuite.zoneName, &armdns.RecordSetsClientListByDNSZoneOptions{Top: nil,
		Recordsetnamesuffix: nil,
	})
	for recordSetsClientNewListByDNSZonePager.More() {
		_, err := recordSetsClientNewListByDNSZonePager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step RecordSets_Get
	fmt.Println("Call operation: RecordSets_Get")
	_, err = recordSetsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, testsuite.relativeRecordSetName, armdns.RecordTypeA, nil)
	testsuite.Require().NoError(err)

	// From step RecordSets_Update
	fmt.Println("Call operation: RecordSets_Update")
	_, err = recordSetsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, testsuite.relativeRecordSetName, armdns.RecordTypeA, armdns.RecordSet{
		Properties: &armdns.RecordSetProperties{
			Metadata: map[string]*string{
				"key2": to.Ptr("value2"),
			},
		},
	}, &armdns.RecordSetsClientUpdateOptions{IfMatch: nil})
	testsuite.Require().NoError(err)

	// From step RecordSets_Delete
	fmt.Println("Call operation: RecordSets_Delete")
	_, err = recordSetsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, testsuite.relativeRecordSetName, armdns.RecordTypeA, &armdns.RecordSetsClientDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
}

// Microsoft.Network/getDnsResourceReference
func (testsuite *DnsTestSuite) TestDnsResourceReference() {
	var err error
	// From step DnsResourceReference_GetByTargetResources
	fmt.Println("Call operation: DnsResourceReference_GetByTargetResources")
	resourceReferenceClient, err := armdns.NewResourceReferenceClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = resourceReferenceClient.GetByTargetResources(testsuite.ctx, armdns.ResourceReferenceRequest{
		Properties: &armdns.ResourceReferenceRequestProperties{
			TargetResources: []*armdns.SubResource{
				{
					ID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/rg1/providers/Microsoft.Network/trafficManagerProfiles/testpp2"),
				}},
		},
	}, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *DnsTestSuite) Cleanup() {
	var err error
	// From step Zones_Delete
	fmt.Println("Call operation: Zones_Delete")
	zonesClient, err := armdns.NewZonesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	zonesClientDeleteResponsePoller, err := zonesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.zoneName, &armdns.ZonesClientBeginDeleteOptions{IfMatch: nil})
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, zonesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}

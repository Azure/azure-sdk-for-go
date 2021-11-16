package armnetwork_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVirtualNetworksClient_BeginCreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"createVN","westus")
	rgName := *rg.Name
	defer clean()

	// create virtual network
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID,cred,opt)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	vnResp, err := vnPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vnResp.Name, vnName)
}

func TestVirtualNetworksClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"getVN","westus")
	rgName := *rg.Name
	defer clean()

	// create virtual network
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID,cred,nil)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	vnResp, err := vnPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vnResp.Name, vnName)

	// get virtual network
	vnResp2,err := vnClient.Get(context.Background(), rgName, vnName, nil)
	require.NoError(t, err)
	require.Equal(t,*vnResp2.Name,vnName)
}

func TestVirtualNetworksClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"listVN","westus")
	rgName := *rg.Name
	defer clean()

	// virtual network create
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID,cred,nil)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	vnResp, err := vnPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vnResp.Name, vnName)

	//virtual network list
	listPager := vnClient.List(rgName,nil)
	require.Equal(t, true,listPager.NextPage(context.Background()))
}

func TestVirtualNetworksClient_UpdateTags(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"updateVN","westus")
	rgName := *rg.Name
	defer clean()

	// virtual network create
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID,cred,nil)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	vnResp, err := vnPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vnResp.Name, vnName)

	//virtual network update tags
	tagResp,err := vnClient.UpdateTags(
		context.Background(),
		rgName,
		vnName,
		armnetwork.TagsObject{
			Tags: map[string]*string{
				"tag1": to.StringPtr("value1"),
				"tag2": to.StringPtr("value2"),
			},
		},
		nil,
		)
	require.NoError(t, err)
	require.Equal(t, *tagResp.Name,vnName)
}

func TestVirtualNetworksClient_BeginDelete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"deleteVN","westus")
	rgName := *rg.Name
	defer clean()

	// virtual network create
	vnClient := armnetwork.NewVirtualNetworksClient(subscriptionID,cred,nil)
	vnName, err := createRandomName(t, "network")
	require.NoError(t, err)
	vnPoller, err := vnClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		armnetwork.VirtualNetwork{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.StringPtr("10.1.0.0/16"),
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	vnResp, err := vnPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vnResp.Name, vnName)

	//virtual network delete
	delPoller,err := vnClient.BeginDelete(context.Background(),rgName,vnName,nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(context.Background(),10*time.Second)
	require.Equal(t, delResp.RawResponse.StatusCode,200)
}
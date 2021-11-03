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


func TestSubnetsClient_BeginCreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"createSUB","westus")
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

	// create subnet
	subClient := armnetwork.NewSubnetsClient(subscriptionID,cred,nil)
	subName, err := createRandomName(t, "subnet")
	require.NoError(t, err)
	subPoller, err := subClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.1.10.0/24"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	subResp, err := subPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *subResp.Name, subName)
}

func TestSubnetsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"getSUB","westus")
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

	// create subnet
	subClient := armnetwork.NewSubnetsClient(subscriptionID,cred,nil)
	subName, err := createRandomName(t, "subnet")
	require.NoError(t, err)
	subPoller, err := subClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.1.10.0/24"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	subResp, err := subPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *subResp.Name, subName)

	// get subnet
	getResp,err := subClient.Get(context.Background(),rgName,vnName,subName,nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Name,subName)
}

func TestSubnetsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"listSUB","westus")
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

	// create subnet
	subClient := armnetwork.NewSubnetsClient(subscriptionID,cred,nil)
	subName, err := createRandomName(t, "subnet")
	require.NoError(t, err)
	subPoller, err := subClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.1.10.0/24"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	subResp, err := subPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *subResp.Name, subName)

	// list subnet
	listPager := subClient.List(rgName,vnName,nil)
	require.Equal(t, listPager.NextPage(context.Background()),true)
}

func TestSubnetsClient_BeginDelete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"deleteSUB","westus")
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

	// create subnet
	subClient := armnetwork.NewSubnetsClient(subscriptionID,cred,nil)
	subName, err := createRandomName(t, "subnet")
	require.NoError(t, err)
	subPoller, err := subClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vnName,
		subName,
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.1.10.0/24"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	subResp, err := subPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *subResp.Name, subName)

	// delete subnet
	delPoller,err := subClient.BeginDelete(context.Background(),rgName,vnName,subName,nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(context.Background(),10*time.Second)
	require.NoError(t, err)
	require.Equal(t, delResp.RawResponse.StatusCode,200)
}

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


func TestIPGroupsClient_BeginCreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"createIPG","westus")
	rgName := *rg.Name
	defer clean()

	// create ip group
	ipgClient := armnetwork.NewIPGroupsClient(subscriptionID,cred,nil)
	ipgName, err := createRandomName(t, "ipg")
	require.NoError(t, err)
	ipgPoller, err := ipgClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		ipgName,
		armnetwork.IPGroup{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.IPGroupPropertiesFormat{
				IPAddresses: []*string{
					to.StringPtr("13.64.39.16/32"),
					to.StringPtr("40.74.146.80/31"),
					to.StringPtr("40.74.147.32/28"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := ipgPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *resp.Name, ipgName)
}

func TestIPGroupsClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"getIPG","westus")
	rgName := *rg.Name
	defer clean()

	// create ip group
	ipgClient := armnetwork.NewIPGroupsClient(subscriptionID,cred,nil)
	ipgName, err := createRandomName(t, "ipg")
	require.NoError(t, err)
	ipgPoller, err := ipgClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		ipgName,
		armnetwork.IPGroup{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.IPGroupPropertiesFormat{
				IPAddresses: []*string{
					to.StringPtr("13.64.39.16/32"),
					to.StringPtr("40.74.146.80/31"),
					to.StringPtr("40.74.147.32/28"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := ipgPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *resp.Name, ipgName)

	// get ip group
	getResp,err := ipgClient.Get(context.Background(),rgName,ipgName,nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Name,ipgName)
}

func TestIPGroupsClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"listIPG","westus")
	rgName := *rg.Name
	defer clean()

	// create ip group
	ipgClient := armnetwork.NewIPGroupsClient(subscriptionID,cred,nil)
	ipgName, err := createRandomName(t, "ipg")
	require.NoError(t, err)
	ipgPoller, err := ipgClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		ipgName,
		armnetwork.IPGroup{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.IPGroupPropertiesFormat{
				IPAddresses: []*string{
					to.StringPtr("13.64.39.16/32"),
					to.StringPtr("40.74.146.80/31"),
					to.StringPtr("40.74.147.32/28"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := ipgPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *resp.Name, ipgName)

	// list ip group
	listPager := ipgClient.List(nil)
	require.Equal(t, listPager.NextPage(context.Background()),true)
}

func TestIPGroupsClient_BeginDelete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, _ := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rg,clean := createResourceGroup(t,cred,subscriptionID,"deleteIPG","westus")
	rgName := *rg.Name
	defer clean()

	// create ip group
	ipgClient := armnetwork.NewIPGroupsClient(subscriptionID,cred,nil)
	ipgName, err := createRandomName(t, "ipg")
	require.NoError(t, err)
	ipgPoller, err := ipgClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		ipgName,
		armnetwork.IPGroup{
			Resource: armnetwork.Resource{
				Location: to.StringPtr("westus"),
			},
			Properties: &armnetwork.IPGroupPropertiesFormat{
				IPAddresses: []*string{
					to.StringPtr("13.64.39.16/32"),
					to.StringPtr("40.74.146.80/31"),
					to.StringPtr("40.74.147.32/28"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp, err := ipgPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *resp.Name, ipgName)

	// delete ip group
	delPoller,err := ipgClient.BeginDelete(context.Background(),rgName,ipgName,nil)
	require.NoError(t, err)
	delResp,err := delPoller.PollUntilDone(context.Background(),10*time.Second)
	require.Equal(t, delResp.RawResponse.StatusCode,200)
}
package armcompute_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDedicatedHostsClient_CreateOrUpdate(t *testing.T) {
	//stop := startTest(t)
	//defer stop()
	//
	//cred, opt := authenticateTest(t)
	//conn := arm.NewDefaultConnection(cred, opt)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	cred,err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	//conn := arm.NewDefaultConnection(cred,nil)
	//subscriptionID,ok := os.LookupEnv("AZURE_SUBSCRIPTION_ID")
	//require.Equal(t, true,ok)

	// create resource group
	rgName, err := createRandomName(t, "testRP")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(subscriptionID,cred,nil)
	_, err = rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	defer cleanup(t, rgClient, rgName)
	require.NoError(t, err)

	// create dedicated host group
	dhgClient := armcompute.NewDedicatedHostGroupsClient(subscriptionID,cred,nil)
	dhgName, err := createRandomName(t, "dhg")
	require.NoError(t, err)
	dhgResp, err := dhgClient.CreateOrUpdate(
		context.Background(),
		rgName,
		dhgName,
		armcompute.DedicatedHostGroup{
			Resource: armcompute.Resource{
				Location: to.StringPtr("eastus"),
			},
			Properties: &armcompute.DedicatedHostGroupProperties{
				PlatformFaultDomainCount: to.Int32Ptr(3),
			},
			Zones: []*string{to.StringPtr("1")},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, *dhgResp.Name, dhgName)

	// create dedicated host
	dhClient := armcompute.NewDedicatedHostsClient(subscriptionID,cred,nil)
	dhName, err := createRandomName(t, "dh")
	require.NoError(t, err)
	dhPoller, err := dhClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		dhgName,
		dhName,
		armcompute.DedicatedHost{
			Resource: armcompute.Resource{
				Location: to.StringPtr("eastus"),
			},
			Properties: &armcompute.DedicatedHostProperties{
				PlatformFaultDomain: to.Int32Ptr(1),
			},
			SKU: &armcompute.SKU{
				Name: to.StringPtr("DSv3-Type1"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	dhResp, err := dhPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *dhResp.Name, dhName)
}

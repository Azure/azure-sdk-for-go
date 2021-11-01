package armcompute_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func cleanup(t *testing.T, client *armresources.ResourceGroupsClient, resourceGroupName string) {
	_, err := client.BeginDelete(context.Background(), resourceGroupName, nil)
	require.NoError(t, err)
}

func TestAvailabilitySetsClient_CreateOrUpdate(t *testing.T) {
	//stop := startTest(t)
	//defer stop()

	//cred, opt := authenticateTest(t)
	//conn := arm.NewDefaultConnection(cred, opt)
	//subscriptionID := recording.GetEnvVariable(t, "AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")

	cred,err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	conn := arm.NewDefaultConnection(cred,nil)
	subscriptionID,ok := os.LookupEnv("AZURE_SUBSCRIPTION_ID")
	require.Equal(t, true,ok)

	// create resource group
	rgName, err := createRandomName(t, "testRP")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(conn, subscriptionID)
	_, err = rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus"),
	}, nil)
	defer cleanup(t, rgClient, rgName)
	require.NoError(t, err)

	// create availability sets
	client := armcompute.NewAvailabilitySetsClient(conn, subscriptionID)
	name, err := createRandomName(t, "set")
	require.NoError(t, err)
	resp, err := client.CreateOrUpdate(
		context.Background(),
		rgName,
		name,
		armcompute.AvailabilitySet{
			Resource: armcompute.Resource{
				Location: to.StringPtr("westus"),
			},
			SKU: &armcompute.SKU{
				Name: to.StringPtr(string(armcompute.AvailabilitySetSKUTypesAligned)),
			},
			Properties: &armcompute.AvailabilitySetProperties{
				PlatformFaultDomainCount:  to.Int32Ptr(1),
				PlatformUpdateDomainCount: to.Int32Ptr(1),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, *resp.Name, name)
}

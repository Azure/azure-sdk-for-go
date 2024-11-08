package armsql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/stretchr/testify/assert"
)

const (
	SubscriptionID    = "faa080af-c1d8-40ad-9cce-e1a450ca5b57"
	ResourceGroupName = "judytest02"
	ResourceLocation  = "eastus2"
	SqlServerName     = "sql-test001"
	DatabaseName      = "test01"
)

func TestArmMysql(t *testing.T) {

}

func GetAzureCredentail() (cred *azidentity.DefaultAzureCredential) {
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Println("get default credentail fail, error:", err.Error())
		return
	}
	fmt.Println("get default credential")
	return cred
}

func CreateAzureResourceGroup() (group *armresources.ResourceGroup, err error) {
	cred := GetAzureCredentail()
	clientFactory, err := armresources.NewClientFactory(SubscriptionID, cred, nil)
	if err != nil {
		fmt.Println("create resource client factory fail, error:", err.Error())
		return
	}
	client := clientFactory.NewResourceGroupsClient()
	createOrUpdateGroupResponse, err := client.CreateOrUpdate(context.TODO(), ResourceGroupName, armresources.ResourceGroup{
		Location: to.Ptr(ResourceLocation),
	}, nil)
	if err != nil {
		fmt.Println("create resource group fail, error:", err.Error())
		return
	}
	group = &createOrUpdateGroupResponse.ResourceGroup
	return
}

func createServer(ctx context.Context, serversClient *armsql.ServersClient) (*armsql.Server, error) {

	pollerResp, err := serversClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		serverName,
		armsql.Server{
			Location: to.Ptr(location),
			Properties: &armsql.ServerProperties{
				AdministratorLogin:         to.Ptr("dummylogin"),
				AdministratorLoginPassword: to.Ptr("QWE123!@#"),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

func getServer(ctx context.Context, serversClient *armsql.ServersClient) (*armsql.Server, error) {

	resp, err := serversClient.Get(ctx, ResourceGroupName, SqlServerName, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

func cleanup(ctx context.Context, resourceGroupClient *armresources.ResourceGroupsClient) error {

	pollerResp, err := resourceGroupClient.BeginDelete(ctx, ResourceGroupName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func TestCreateOrUpdateGroupOfSubscriptionId(t *testing.T) {
	// subsriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	subsriptionId := "faa080af-c1d8-40ad-9cce-e1a450ca5b57"
	assert.NotEmpty(t, subsriptionId)
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	fmt.Println("get default credential")

	// new client factory
	clientFactory, err := armresources.NewClientFactory(subsriptionId, cred, nil)
	assert.Nil(t, err)
	assert.NotNil(t, clientFactory)
	client := clientFactory.NewResourceGroupsClient()
	assert.NotNil(t, client)
	ctx := context.Background()

	createOrUpdateGroupResponse, err := client.CreateOrUpdate(ctx, ResourceGroupName, armresources.ResourceGroup{
		Location: to.Ptr(ResourceLocation),
	}, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createOrUpdateGroupResponse)
	fmt.Println("create resource group client")

	// check whether create new group successfully
	checkGroupExistResponse, err := client.CheckExistence(ctx, ResourceGroupName, nil)
	assert.Nil(t, err)
	assert.True(t, checkGroupExistResponse.Success)
	fmt.Println("create new resource group ", ResourceGroupName, " of ", subsriptionId, "successfully")

	// NewTestDeleteGroupOfSubscriptionId(t, ctx, client)

}

func NewTestDeleteGroupOfSubscriptionId(t *testing.T, ctx context.Context, client *armresources.ResourceGroupsClient) {
	// delete resource group created newly
	resourceGroupsClientDeleteResponse, err := client.BeginDelete(ctx, ResourceGroupName, nil)
	assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	response, err := resourceGroupsClientDeleteResponse.Poll(ctx)
	assert.Nil(t, err)
	assert.True(t, response.StatusCode >= 200 && response.StatusCode < 300)
	fmt.Println("delete resource group successfully")
}

func TestCreateOrUpdateArmMysqlResourceOnGroup(t *testing.T) {
	// make sure that the group has been created
	TestCreateOrUpdateGroupOfSubscriptionId(t)

	cred := GetAzureCredentail()
	assert.NotNil(t, cred)
	// create new msql resource under group
	sqlClientFactory, err := armsql.NewClientFactory(SubscriptionID, cred, nil)
	assert.Nil(t, err)
	assert.NotNil(t, sqlClientFactory)

	serverclient := sqlClientFactory.NewServersClient()
	assert.NotNil(t, serverclient)
	databasesClient := sqlClientFactory.NewDatabasesClient()
	assert.NotNil(t, databasesClient)

	// resourceGroup, err := CreateAzureResourceGroup()
	// assert.Nil(t, err)
	server, err := createServer(context.Background(), serverclient)
	assert.Nil(t, err)
	assert.NotNil(t, server)
	fmt.Println("create serverId", *server.ID)

	server, err = getServer(context.Background(), serverclient)
	assert.Nil(t, err)
	assert.NotNil(t, server)
	fmt.Println("get server:", *server.ID)

	// cleanup(context.Background(), resourceGroup)
	// PUT https://management.azure.com/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/flexibleServers/{serverName}?api-version=2024-06-01-preview

}

func TestGetListOfArmMysqlResourceOnGroup(t *testing.T) {

}

func TestDeleteArmMysqlResourceOnGroup(t *testing.T) {

}

func TestArmMysqlInstance(t *testing.T) {

}
